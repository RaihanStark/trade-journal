package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/raihanstark/trade-journal/internal/application/auth"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *auth.Service
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration requests
func (h *AuthHandler) Register(c echo.Context) error {
	var req auth.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	response, err := h.authService.Register(c.Request().Context(), req)
	if err != nil {
		switch err {
		case auth.ErrEmailAlreadyExists:
			return c.JSON(http.StatusConflict, map[string]string{"error": "Email already exists"})
		case auth.ErrHashingPassword:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process password"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
	}

	return c.JSON(http.StatusCreated, response)
}

// Login handles user login requests
func (h *AuthHandler) Login(c echo.Context) error {
	var req auth.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	response, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials:
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
	}

	return c.JSON(http.StatusOK, response)
}
