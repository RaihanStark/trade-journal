package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/raihanstark/trade-journal/internal/application/trade"
	"github.com/raihanstark/trade-journal/internal/infrastructure/storage"
)

type TradeHandler struct {
	service *trade.Service
	storage *storage.MinIOStorage
}

func NewTradeHandler(service *trade.Service, storage *storage.MinIOStorage) *TradeHandler {
	return &TradeHandler{
		service: service,
		storage: storage,
	}
}

func (h *TradeHandler) CreateTrade(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	var req trade.CreateTradeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	result, err := h.service.CreateTrade(c.Request().Context(), userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, result)
}

func (h *TradeHandler) GetTrades(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	// Get date filter parameters
	var startDate, endDate *string
	if sd := c.QueryParam("start_date"); sd != "" {
		startDate = &sd
	}
	if ed := c.QueryParam("end_date"); ed != "" {
		endDate = &ed
	}

	// If account_id is provided, get trades by account ID
	if accountID := c.QueryParam("account_id"); accountID != "" {
		accountID, err := strconv.ParseInt(accountID, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid account ID",
			})
		}
		trades, err := h.service.GetTradesByAccountIDWithDateFilter(c.Request().Context(), accountID, userID, startDate, endDate)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, trades)
	}

	// Otherwise, get all trades for the user
	trades, err := h.service.GetUserTradesWithDateFilter(c.Request().Context(), userID, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, trades)
}

func (h *TradeHandler) GetTrade(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid trade ID",
		})
	}

	result, err := h.service.GetTrade(c.Request().Context(), id, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Trade not found",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TradeHandler) UpdateTrade(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid trade ID",
		})
	}

	var req trade.UpdateTradeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	result, err := h.service.UpdateTrade(c.Request().Context(), id, userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TradeHandler) DeleteTrade(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid trade ID",
		})
	}

	if err := h.service.DeleteTrade(c.Request().Context(), id, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *TradeHandler) UploadChart(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	tradeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid trade ID",
		})
	}

	chartType := c.Param("type") // "before" or "after"
	if chartType != "before" && chartType != "after" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid chart type. Must be 'before' or 'after'",
		})
	}

	// Get uploaded file
	file, err := c.FormFile("chart")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No file uploaded",
		})
	}

	// Validate file type (must be image)
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "File must be an image",
		})
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "File size must be less than 5MB",
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to open file",
		})
	}
	defer src.Close()

	// Upload to MinIO
	filename := fmt.Sprintf("trade-%d-chart-%s", tradeID, chartType)
	url, err := h.storage.UploadFile(c.Request().Context(), filename, src, file.Size, contentType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to upload file: %v", err),
		})
	}

	// Update trade in database
	if chartType == "before" {
		_, err = h.service.UpdateChartBefore(c.Request().Context(), tradeID, userID, url)
	} else {
		_, err = h.service.UpdateChartAfter(c.Request().Context(), tradeID, userID, url)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"url":     url,
		"message": fmt.Sprintf("Chart %s uploaded successfully", chartType),
	})
}
