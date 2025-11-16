package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/raihanstark/trade-journal/internal/testutil"
)

func TestE2E_Analytics_GetUserAnalytics(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	// Register and get token
	authToken := registerUser(t, e, "analytics@example.com", "password123")

	// Create account
	accountID := createAccount(t, e, authToken, "Test Account")

	t.Run("calculates analytics from real trades", func(t *testing.T) {
		// Create test trades with different P/L values
		trades := []struct {
			tradeType string
			exit      float64
		}{
			{"BUY", 1.1100},  // +100 pips * 1 lot * 10 = +1000 (profit)
			{"SELL", 1.1050}, // Price up 50 pips = loss for SELL = -500
			{"BUY", 1.1200},  // +200 pips * 1 lot * 10 = +2000 (profit)
			{"SELL", 1.1075}, // Price up 75 pips = loss for SELL = -750
		}

		for _, trade := range trades {
			createTrade(t, e, authToken, accountID, trade.tradeType, 1.1000, &trade.exit, "2025-01-01")
		}

		// Call analytics endpoint
		req := httptest.NewRequest(http.MethodGet, "/api/analytics", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var analytics map[string]any
		json.Unmarshal(rec.Body.Bytes(), &analytics)

		// Verify analytics calculations
		// Total P/L: 1000 - 500 + 2000 - 750 = 1750
		if analytics["total_pl"] != 1750.0 {
			t.Errorf("total_pl = %v, want 1750", analytics["total_pl"])
		}

		if analytics["total_trades"] != 4.0 {
			t.Errorf("total_trades = %v, want 4", analytics["total_trades"])
		}

		if analytics["winning_trades"] != 2.0 {
			t.Errorf("winning_trades = %v, want 2", analytics["winning_trades"])
		}

		if analytics["losing_trades"] != 2.0 {
			t.Errorf("losing_trades = %v, want 2", analytics["losing_trades"])
		}

		// Win rate: 2/4 * 100 = 50%
		if analytics["win_rate"] != 50.0 {
			t.Errorf("win_rate = %v, want 50", analytics["win_rate"])
		}

		// Avg win: (1000 + 2000) / 2 = 1500
		if analytics["avg_win"] != 1500.0 {
			t.Errorf("avg_win = %v, want 1500", analytics["avg_win"])
		}

		// Avg loss: -(500 + 750) / 2 = -625
		if analytics["avg_loss"] != -625.0 {
			t.Errorf("avg_loss = %v, want -625", analytics["avg_loss"])
		}

		// Profit factor: 3000 / 1250 = 2.4
		if analytics["profit_factor"] != 2.4 {
			t.Errorf("profit_factor = %v, want 2.4", analytics["profit_factor"])
		}

		if analytics["largest_win"] != 2000.0 {
			t.Errorf("largest_win = %v, want 2000", analytics["largest_win"])
		}

		if analytics["largest_loss"] != -750.0 {
			t.Errorf("largest_loss = %v, want -750", analytics["largest_loss"])
		}
	})
}

func TestE2E_Analytics_FiltersOpenTrades(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	// Register and get token
	authToken := registerUser(t, e, "openfilter@example.com", "password123")

	// Create account
	accountID := createAccount(t, e, authToken, "Test Account")

	t.Run("filters out open trades from analytics", func(t *testing.T) {
		// Create closed trade
		exit := 1.1100
		createTrade(t, e, authToken, accountID, "BUY", 1.1000, &exit, "2025-01-01")

		// Create open trade (no exit)
		createTrade(t, e, authToken, accountID, "SELL", 1.1100, nil, "2025-01-01")

		// Call analytics endpoint
		req := httptest.NewRequest(http.MethodGet, "/api/analytics", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var analytics map[string]any
		json.Unmarshal(rec.Body.Bytes(), &analytics)

		// Should only count the closed trade
		if analytics["total_trades"] != 1.0 {
			t.Errorf("total_trades = %v, want 1 (open trades filtered)", analytics["total_trades"])
		}

		// P/L: +100 pips * 1 lot * 10 = 1000
		if analytics["total_pl"] != 1000.0 {
			t.Errorf("total_pl = %v, want 1000", analytics["total_pl"])
		}
	})
}

func TestE2E_Analytics_FiltersDepositsAndWithdrawals(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	// Register and get token
	authToken := registerUser(t, e, "depositfilter@example.com", "password123")

	// Create account
	accountID := createAccount(t, e, authToken, "Test Account")

	t.Run("filters out deposits and withdrawals from analytics", func(t *testing.T) {
		// Create BUY trade
		exit := 1.1100
		createTrade(t, e, authToken, accountID, "BUY", 1.1000, &exit, "2025-01-01")

		// Create DEPOSIT (should be filtered)
		createDeposit(t, e, authToken, accountID, 1000.0)

		// Create WITHDRAW (should be filtered)
		createWithdrawal(t, e, authToken, accountID, 500.0)

		// Call analytics endpoint
		req := httptest.NewRequest(http.MethodGet, "/api/analytics", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var analytics map[string]any
		json.Unmarshal(rec.Body.Bytes(), &analytics)

		// Should only count the BUY trade
		if analytics["total_trades"] != 1.0 {
			t.Errorf("total_trades = %v, want 1 (deposits/withdrawals filtered)", analytics["total_trades"])
		}

		// P/L: +100 pips * 1 lot * 10 = 1000
		if analytics["total_pl"] != 1000.0 {
			t.Errorf("total_pl = %v, want 1000", analytics["total_pl"])
		}
	})
}

func TestE2E_Analytics_NoTrades(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	// Register and get token
	authToken := registerUser(t, e, "notrades@example.com", "password123")

	t.Run("returns zero analytics when no trades exist", func(t *testing.T) {
		// Call analytics endpoint without creating any trades
		req := httptest.NewRequest(http.MethodGet, "/api/analytics", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var analytics map[string]any
		json.Unmarshal(rec.Body.Bytes(), &analytics)

		// Verify zero values
		if analytics["total_trades"] != 0.0 {
			t.Errorf("total_trades = %v, want 0", analytics["total_trades"])
		}

		if analytics["total_pl"] != 0.0 {
			t.Errorf("total_pl = %v, want 0", analytics["total_pl"])
		}

		if analytics["win_rate"] != 0.0 {
			t.Errorf("win_rate = %v, want 0", analytics["win_rate"])
		}

		if analytics["profit_factor"] != 0.0 {
			t.Errorf("profit_factor = %v, want 0", analytics["profit_factor"])
		}
	})
}

// Helper functions
func createTrade(t *testing.T, e *echo.Echo, token string, accountID int, tradeType string, entry float64, exit *float64, date string) {
	t.Helper()

	payload := map[string]any{
		"account_id": accountID,
		"date":       date,
		"time":       "10:00",
		"pair":       "EUR/USD",
		"type":       tradeType,
		"entry":      entry,
		"lots":       1.0,
	}

	if exit != nil {
		payload["exit"] = *exit
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/trades", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("failed to create trade: %d - %s", rec.Code, rec.Body.String())
	}
}

func createDeposit(t *testing.T, e *echo.Echo, token string, accountID int, amount float64) {
	t.Helper()

	payload := map[string]any{
		"account_id": accountID,
		"date":       "2024-01-02",
		"time":       "11:00",
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
		t.Fatalf("failed to create deposit: %d - %s", rec.Code, rec.Body.String())
	}
}

func createWithdrawal(t *testing.T, e *echo.Echo, token string, accountID int, amount float64) {
	t.Helper()

	payload := map[string]any{
		"account_id": accountID,
		"date":       "2024-01-03",
		"time":       "12:00",
		"type":       "WITHDRAW",
		"amount":     amount,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/trades", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("failed to create withdrawal: %d - %s", rec.Code, rec.Body.String())
	}
}
