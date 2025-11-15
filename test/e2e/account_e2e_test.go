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

func TestE2E_Account_CreateAndList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	// Register and get token
	authToken := registerUser(t, e, "accountuser@example.com", "password123")

	var accountID int

	// Create account
	t.Run("create account", func(t *testing.T) {
		payload := map[string]any{
			"name":           "Demo Account",
			"broker":         "IC Markets",
			"account_number": "12345678",
			"account_type":   "demo",
			"currency":       "USD",
			"is_active":      true,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
		}

		var response map[string]any
		json.Unmarshal(rec.Body.Bytes(), &response)
		accountID = int(response["id"].(float64))

		if response["name"] != "Demo Account" {
			t.Errorf("expected name 'Demo Account', got %v", response["name"])
		}
		if response["broker"] != "IC Markets" {
			t.Errorf("expected broker 'IC Markets', got %v", response["broker"])
		}
		if response["current_balance"] != 0.0 {
			t.Errorf("expected initial balance 0, got %v", response["current_balance"])
		}
	})

	// List accounts
	t.Run("list accounts", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/accounts", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var accounts []map[string]any
		json.Unmarshal(rec.Body.Bytes(), &accounts)

		if len(accounts) != 1 {
			t.Errorf("expected 1 account, got %d", len(accounts))
		}
		if len(accounts) > 0 && accounts[0]["name"] != "Demo Account" {
			t.Errorf("expected account name 'Demo Account', got %v", accounts[0]["name"])
		}
	})

	// Verify in database
	t.Run("account saved to database", func(t *testing.T) {
		var name, broker string
		var balance float64
		err := pg.DB.QueryRow("SELECT name, broker, current_balance FROM accounts WHERE id = $1", accountID).
			Scan(&name, &broker, &balance)

		if err != nil {
			t.Fatalf("failed to query account: %v", err)
		}
		if name != "Demo Account" {
			t.Errorf("expected name 'Demo Account', got %s", name)
		}
		if broker != "IC Markets" {
			t.Errorf("expected broker 'IC Markets', got %s", broker)
		}
		if balance != 0.0 {
			t.Errorf("expected balance 0, got %.2f", balance)
		}
	})
}

func TestE2E_Account_UpdateAndDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	// Register and get token
	authToken := registerUser(t, e, "updateuser@example.com", "password123")

	// Create account
	accountID := createAccount(t, e, authToken, "Old Name")

	// Update account
	t.Run("update account", func(t *testing.T) {
		payload := map[string]any{
			"name":           "New Name",
			"broker":         "New Broker",
			"account_number": "87654321",
			"account_type":   "live",
			"currency":       "EUR",
			"is_active":      false,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/accounts/%d", accountID), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var response map[string]any
		json.Unmarshal(rec.Body.Bytes(), &response)

		if response["name"] != "New Name" {
			t.Errorf("expected name 'New Name', got %v", response["name"])
		}
		if response["broker"] != "New Broker" {
			t.Errorf("expected broker 'New Broker', got %v", response["broker"])
		}
		if response["account_type"] != "live" {
			t.Errorf("expected type 'live', got %v", response["account_type"])
		}

		// Verify in database
		var name, broker string
		err := pg.DB.QueryRow("SELECT name, broker FROM accounts WHERE id = $1", accountID).
			Scan(&name, &broker)

		if err != nil {
			t.Fatalf("failed to query account: %v", err)
		}
		if name != "New Name" {
			t.Errorf("expected saved name 'New Name', got %s", name)
		}
		if broker != "New Broker" {
			t.Errorf("expected saved broker 'New Broker', got %s", broker)
		}
	})

	// Delete account
	t.Run("delete account", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/accounts/%d", accountID), nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		// Verify deleted from database
		var count int
		err := pg.DB.QueryRow("SELECT COUNT(*) FROM accounts WHERE id = $1", accountID).Scan(&count)
		if err != nil {
			t.Fatalf("failed to query account count: %v", err)
		}
		if count != 0 {
			t.Errorf("expected account to be deleted, but found %d accounts", count)
		}
	})
}
