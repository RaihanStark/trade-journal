package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/raihanstark/trade-journal/internal/application/account"
)

// AccountHandler handles account HTTP requests
type AccountHandler struct {
	accountService *account.Service
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(accountService *account.Service) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// CreateAccount handles account creation requests
func (h *AccountHandler) CreateAccount(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	var req account.CreateAccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	acc, err := h.accountService.CreateAccount(c.Request().Context(), userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create account"})
	}

	return c.JSON(http.StatusCreated, acc)
}

// GetAccounts handles fetching all accounts for a user
func (h *AccountHandler) GetAccounts(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	accounts, err := h.accountService.GetUserAccounts(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch accounts"})
	}

	return c.JSON(http.StatusOK, accounts)
}

// GetAccount handles fetching a single account
func (h *AccountHandler) GetAccount(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid account ID"})
	}

	acc, err := h.accountService.GetAccount(c.Request().Context(), id, userID)
	if err != nil {
		if err == account.ErrAccountNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Account not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch account"})
	}

	return c.JSON(http.StatusOK, acc)
}

// UpdateAccount handles account update requests
func (h *AccountHandler) UpdateAccount(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid account ID"})
	}

	var req account.UpdateAccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	acc, err := h.accountService.UpdateAccount(c.Request().Context(), id, userID, req)
	if err != nil {
		if err == account.ErrAccountNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Account not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update account"})
	}

	return c.JSON(http.StatusOK, acc)
}

// DeleteAccount handles account deletion requests
func (h *AccountHandler) DeleteAccount(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid account ID"})
	}

	if err := h.accountService.DeleteAccount(c.Request().Context(), id, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete account"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Account deleted successfully"})
}
