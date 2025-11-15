package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/raihanstark/trade-journal/internal/application/strategy"
)

// StrategyHandler handles strategy HTTP requests
type StrategyHandler struct {
	strategyService *strategy.Service
}

// NewStrategyHandler creates a new strategy handler
func NewStrategyHandler(strategyService *strategy.Service) *StrategyHandler {
	return &StrategyHandler{
		strategyService: strategyService,
	}
}

// CreateStrategy handles strategy creation requests
func (h *StrategyHandler) CreateStrategy(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	var req strategy.CreateStrategyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	strat, err := h.strategyService.CreateStrategy(c.Request().Context(), userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create strategy"})
	}

	return c.JSON(http.StatusCreated, strat)
}

// GetStrategies handles fetching all strategies for a user
func (h *StrategyHandler) GetStrategies(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	strategies, err := h.strategyService.GetUserStrategies(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch strategies"})
	}

	return c.JSON(http.StatusOK, strategies)
}

// GetStrategy handles fetching a single strategy
func (h *StrategyHandler) GetStrategy(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid strategy ID"})
	}

	strat, err := h.strategyService.GetStrategy(c.Request().Context(), id, userID)
	if err != nil {
		if err == strategy.ErrStrategyNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Strategy not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch strategy"})
	}

	return c.JSON(http.StatusOK, strat)
}

// UpdateStrategy handles strategy update requests
func (h *StrategyHandler) UpdateStrategy(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid strategy ID"})
	}

	var req strategy.UpdateStrategyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	strat, err := h.strategyService.UpdateStrategy(c.Request().Context(), id, userID, req)
	if err != nil {
		if err == strategy.ErrStrategyNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Strategy not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update strategy"})
	}

	return c.JSON(http.StatusOK, strat)
}

// DeleteStrategy handles strategy deletion requests
func (h *StrategyHandler) DeleteStrategy(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid strategy ID"})
	}

	if err := h.strategyService.DeleteStrategy(c.Request().Context(), id, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete strategy"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Strategy deleted successfully"})
}
