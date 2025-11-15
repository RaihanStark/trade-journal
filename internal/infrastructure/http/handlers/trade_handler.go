package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/raihanstark/trade-journal/internal/application/trade"
)

type TradeHandler struct {
	service *trade.Service
}

func NewTradeHandler(service *trade.Service) *TradeHandler {
	return &TradeHandler{service: service}
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

	trades, err := h.service.GetUserTrades(c.Request().Context(), userID)
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
