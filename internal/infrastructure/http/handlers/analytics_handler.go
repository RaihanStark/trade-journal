package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/raihanstark/trade-journal/internal/application/analytics"
)

type AnalyticsHandler struct {
	service *analytics.Service
}

func NewAnalyticsHandler(service *analytics.Service) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) GetUserAnalytics(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	result, err := h.service.GetUserAnalytics(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}
