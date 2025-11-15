package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/raihanstark/trade-journal/internal/testutil"
)

func TestE2E_Trade_DepositAndBalanceUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	// Register and get token
	authToken := registerUser(t, e, "trader@example.com", "password123")

	// Create account
	accountID := createAccount(t, e, authToken, "Demo Account")

	// Make deposit
	t.Run("deposit updates account balance", func(t *testing.T) {
		payload := map[string]any{
			"account_id": accountID,
			"date":       "2025-01-15",
			"time":       "10:00",
			"type":       "DEPOSIT",
			"amount":     1000.0,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/trades", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
		}

		// Verify balance updated
		var balance float64
		err := pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", accountID).Scan(&balance)
		if err != nil {
			t.Fatalf("failed to query balance: %v", err)
		}
		if balance != 1000.0 {
			t.Errorf("expected balance 1000, got %.2f", balance)
		}
	})
}

func TestE2E_Trade_ClosedTradeWithProfitAndLoss(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	// Register and get token
	authToken := registerUser(t, e, "trader@example.com", "password123")

	// Create account and deposit
	accountID := createAccount(t, e, authToken, "Live Account")
	makeDeposit(t, e, authToken, accountID, 1000.0)

	var tradeID int

	// Create closed trade with profit
	t.Run("closed trade updates balance with P/L", func(t *testing.T) {
		exit := 1.1050
		payload := map[string]any{
			"account_id": accountID,
			"date":       "2025-01-15",
			"time":       "14:30",
			"pair":       "EUR/USD",
			"type":       "BUY",
			"entry":      1.1000,
			"exit":       exit,
			"lots":       1.0,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/trades", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
		}

		var response map[string]any
		json.Unmarshal(rec.Body.Bytes(), &response)
		tradeID = int(response["id"].(float64))

		// Verify P/L calculated (50 pips * 1 lot * $10 = $500)
		pl := response["pl"].(float64)
		if pl != 500.0 {
			t.Errorf("expected P/L 500, got %.2f", pl)
		}

		// Verify balance updated to $1500
		var balance float64
		err := pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", accountID).Scan(&balance)
		if err != nil {
			t.Fatalf("failed to query balance: %v", err)
		}
		if balance != 1500.0 {
			t.Errorf("expected balance 1500, got %.2f", balance)
		}
	})

	// Delete trade and verify balance reverts
	t.Run("deleting trade reverts balance", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/trades/%d", tradeID), nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Fatalf("expected status 204, got %d: %s", rec.Code, rec.Body.String())
		}

		// Verify balance reverted to $1000
		var balance float64
		err := pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", accountID).Scan(&balance)
		if err != nil {
			t.Fatalf("failed to query balance: %v", err)
		}
		if balance != 1000.0 {
			t.Errorf("expected balance 1000 after delete, got %.2f", balance)
		}
	})
}

// Helper functions
func registerUser(t *testing.T, e *echo.Echo, email, password string) string {
	t.Helper()

	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("failed to register: %d - %s", rec.Code, rec.Body.String())
	}

	var response map[string]any
	json.Unmarshal(rec.Body.Bytes(), &response)
	return response["token"].(string)
}

func createAccount(t *testing.T, e *echo.Echo, token, name string) int {
	t.Helper()

	payload := map[string]any{
		"name":           name,
		"broker":         "IC Markets",
		"account_number": "12345678",
		"account_type":   "demo",
		"currency":       "USD",
		"is_active":      true,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("failed to create account: %d - %s", rec.Code, rec.Body.String())
	}

	var response map[string]any
	json.Unmarshal(rec.Body.Bytes(), &response)
	return int(response["id"].(float64))
}

func makeDeposit(t *testing.T, e *echo.Echo, token string, accountID int, amount float64) {
	t.Helper()

	payload := map[string]any{
		"account_id": accountID,
		"date":       "2025-01-15",
		"time":       "09:00",
		"type":       "DEPOSIT",
		"amount":     amount,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/trades", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("failed to make deposit: %d - %s", rec.Code, rec.Body.String())
	}
}
