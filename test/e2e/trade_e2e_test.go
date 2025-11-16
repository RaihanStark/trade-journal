package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/raihanstark/trade-journal/internal/testutil"
)

func TestE2E_Trade_Validation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	// Register and get token
	authToken := registerUser(t, e, "validator@example.com", "password123")

	t.Run("account_id is required for creating trade", func(t *testing.T) {
		payload := map[string]any{
			"date":   "2025-01-15",
			"time":   "10:00",
			"type":   "DEPOSIT",
			"amount": 1000.0,
			// Missing account_id
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/trades", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", rec.Code)
		}

		var response map[string]any
		json.Unmarshal(rec.Body.Bytes(), &response)

		errorMsg := response["error"].(string)
		if errorMsg != "account_id is required" {
			t.Errorf("expected error 'account_id is required', got '%s'", errorMsg)
		}
	})
}

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

	// Create trade with negative P/L
	t.Run("create trade with negative P/L should update balance with negative amount", func(t *testing.T) {
		exit := 0.5000
		payload := map[string]any{
			"account_id": accountID,
			"date":       "2025-01-15",
			"time":       "14:30",
			"pair":       "EUR/USD",
			"type":       "BUY",
			"entry":      1.0000,
			"exit":       &exit,
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
		pl := response["pl"].(float64)
		if pl != -50000.0 {
			t.Errorf("expected P/L -50000, got %.2f", pl)
		}

		// Verify balance updated to -48500 (1000 + 500 - 50000)
		var balance float64
		err := pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", accountID).Scan(&balance)
		if err != nil {
			t.Fatalf("failed to query balance: %v", err)
		}
		if balance != -48500.0 {
			t.Errorf("expected balance -48500, got %.2f", balance)
		}
	})

	// Delete first trade (the +$500 one) and verify balance reverts
	t.Run("deleting trade reverts balance", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/trades/%d", tradeID), nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Fatalf("expected status 204, got %d: %s", rec.Code, rec.Body.String())
		}

		// Verify balance reverted to -49000 (1000 - 50000, after removing the +500 trade)
		var balance float64
		err := pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", accountID).Scan(&balance)
		if err != nil {
			t.Fatalf("failed to query balance: %v", err)
		}
		if balance != -49000.0 {
			t.Errorf("expected balance -49000 after delete, got %.2f", balance)
		}
	})
}

func TestE2E_Trade_FilterTradesByAccount(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	// Register and get token
	authToken := registerUser(t, e, "trader@example.com", "password123")

	// Create account
	accountID := createAccount(t, e, authToken, "Test Account")
	accountID2 := createAccount(t, e, authToken, "Test Account 2")

	// Create trade
	createTrade(t, e, authToken, accountID, "BUY", 1.1000, nil, "2025-01-15")
	createTrade(t, e, authToken, accountID2, "BUY", 1.2, nil, "2025-01-15")

	t.Run("get trades by account 1", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/trades?account_id="+strconv.Itoa(accountID), nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var trades []map[string]any
		json.Unmarshal(rec.Body.Bytes(), &trades)
		if len(trades) != 1 {
			t.Fatalf("expected 1 trade, got %d", len(trades))
		}

		if int(trades[0]["account_id"].(float64)) != accountID {
			t.Fatalf("expected account_id %d, got %v", accountID, trades[0]["account_id"])
		}
	})

	t.Run("get trades by account 2", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/trades?account_id="+strconv.Itoa(accountID2), nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var trades []map[string]any
		json.Unmarshal(rec.Body.Bytes(), &trades)
		if len(trades) != 1 {
			t.Fatalf("expected 1 trade, got %d", len(trades))
		}

		if int(trades[0]["account_id"].(float64)) != accountID2 {
			t.Fatalf("expected account_id %d, got %v", accountID2, trades[0]["account_id"])
		}
	})

	// Test that there's no trades for account 3
	t.Run("get trades by account 3", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/trades?account_id="+strconv.Itoa(3), nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var trades []map[string]any
		json.Unmarshal(rec.Body.Bytes(), &trades)
		if len(trades) != 0 {
			t.Fatalf("expected 0 trades, got %d", len(trades))
		}
	})

}

func TestE2E_Trade_FilterTradesByDate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	// Create an account
	authToken := registerUser(t, e, "trader@example.com", "password123")
	accountID := createAccount(t, e, authToken, "Test Account")

	// create trade between 2025-01-14 and 2025-01-16
	createTrade(t, e, authToken, accountID, "BUY", 1.1000, nil, "2025-01-14")
	createTrade(t, e, authToken, accountID, "BUY", 1.1000, nil, "2025-01-15")
	createTrade(t, e, authToken, accountID, "BUY", 1.1000, nil, "2025-01-16")

	t.Run("get trades by date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/trades?start_date=2025-01-15&end_date=2025-01-16", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var trades []map[string]any
		json.Unmarshal(rec.Body.Bytes(), &trades)

		if len(trades) != 2 {
			t.Fatalf("expected 2 trades, got %d", len(trades))
		}

		for _, trade := range trades {
			if trade["date"] != "2025-01-15" && trade["date"] != "2025-01-16" {
				t.Fatalf("expected date 2025-01-15 or 2025-01-16, got %s", trade["date"])
			}
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

func TestE2E_Trade_ChartImageUpload(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	// Register and get token
	authToken := registerUser(t, e, "chartupload@example.com", "password123")

	// Create account
	accountID := createAccount(t, e, authToken, "Chart Test Account")

	// Create a trade
	tradePayload := map[string]any{
		"account_id":  accountID,
		"date":        "2025-01-16",
		"time":        "14:30",
		"pair":        "EURUSD",
		"type":        "BUY",
		"entry":       1.1000,
		"exit":        1.1050,
		"lots":        0.1,
		"stop_loss":   1.0950,
		"take_profit": 1.1100,
		"notes":       "Test trade for chart upload",
	}
	tradeBody, _ := json.Marshal(tradePayload)

	req := httptest.NewRequest(http.MethodPost, "/api/trades", bytes.NewReader(tradeBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("failed to create trade: %d - %s", rec.Code, rec.Body.String())
	}

	var tradeResponse map[string]any
	json.Unmarshal(rec.Body.Bytes(), &tradeResponse)
	tradeID := int64(tradeResponse["id"].(float64))

	t.Run("upload chart before image", func(t *testing.T) {
		// Create a test image (simple 1x1 PNG)
		imageData := createTestImage(t)

		// Create multipart form data
		body := &bytes.Buffer{}
		writer := createMultipartForm(t, body, "chart", "test-before.png", imageData)

		url := fmt.Sprintf("/api/trades/%d/chart/before", tradeID)
		req := httptest.NewRequest(http.MethodPost, url, body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var response map[string]any
		json.Unmarshal(rec.Body.Bytes(), &response)

		// Verify response contains URL
		if response["url"] == nil {
			t.Fatal("expected URL in response")
		}

		chartURL := response["url"].(string)
		if chartURL == "" {
			t.Fatal("expected non-empty chart URL")
		}

		// Verify database updated
		var chartBefore *string
		err := pg.DB.QueryRow("SELECT chart_before FROM trades WHERE id = $1", tradeID).Scan(&chartBefore)
		if err != nil {
			t.Fatalf("failed to query chart_before: %v", err)
		}

		if chartBefore == nil {
			t.Fatal("expected chart_before to be set in database")
		}

		if *chartBefore != chartURL {
			t.Fatalf("expected chart_before to be %s, got %s", chartURL, *chartBefore)
		}
	})

	t.Run("upload chart after image", func(t *testing.T) {
		imageData := createTestImage(t)

		body := &bytes.Buffer{}
		writer := createMultipartForm(t, body, "chart", "test-after.png", imageData)

		url := fmt.Sprintf("/api/trades/%d/chart/after", tradeID)
		req := httptest.NewRequest(http.MethodPost, url, body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var response map[string]any
		json.Unmarshal(rec.Body.Bytes(), &response)

		if response["url"] == nil {
			t.Fatal("expected URL in response")
		}

		// Verify database updated
		var chartAfter *string
		err := pg.DB.QueryRow("SELECT chart_after FROM trades WHERE id = $1", tradeID).Scan(&chartAfter)
		if err != nil {
			t.Fatalf("failed to query chart_after: %v", err)
		}

		if chartAfter == nil {
			t.Fatal("expected chart_after to be set in database")
		}
	})

	t.Run("upload invalid file type returns error", func(t *testing.T) {
		// Create invalid file (text file)
		textData := []byte("This is not an image")

		body := &bytes.Buffer{}
		writer := createMultipartForm(t, body, "chart", "test.txt", textData)

		url := fmt.Sprintf("/api/trades/%d/chart/before", tradeID)
		req := httptest.NewRequest(http.MethodPost, url, body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("upload without authentication returns error", func(t *testing.T) {
		imageData := createTestImage(t)

		body := &bytes.Buffer{}
		writer := createMultipartForm(t, body, "chart", "test.png", imageData)

		url := fmt.Sprintf("/api/trades/%d/chart/before", tradeID)
		req := httptest.NewRequest(http.MethodPost, url, body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		// No authorization header
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", rec.Code)
		}
	})
}
