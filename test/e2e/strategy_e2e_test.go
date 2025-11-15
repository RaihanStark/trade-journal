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

func TestE2E_Strategy_CreateAndList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	// Register and get token
	authToken := registerUser(t, e, "strategyuser@example.com", "password123")

	var strategyID int

	// Create strategy
	t.Run("create strategy", func(t *testing.T) {
		payload := map[string]string{
			"name":        "Breakout Strategy",
			"description": "Trade price breakouts at key levels",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/strategies", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
		}

		var response map[string]any
		json.Unmarshal(rec.Body.Bytes(), &response)
		strategyID = int(response["id"].(float64))

		if response["name"] != "Breakout Strategy" {
			t.Errorf("expected name 'Breakout Strategy', got %v", response["name"])
		}
		if response["description"] != "Trade price breakouts at key levels" {
			t.Errorf("expected description 'Trade price breakouts at key levels', got %v", response["description"])
		}
	})

	// List strategies
	t.Run("list strategies", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/strategies", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var strategies []map[string]any
		json.Unmarshal(rec.Body.Bytes(), &strategies)

		if len(strategies) != 1 {
			t.Errorf("expected 1 strategy, got %d", len(strategies))
		}
		if len(strategies) > 0 && strategies[0]["name"] != "Breakout Strategy" {
			t.Errorf("expected strategy name 'Breakout Strategy', got %v", strategies[0]["name"])
		}
	})

	// Verify in database
	t.Run("strategy saved to database", func(t *testing.T) {
		var name, description string
		err := pg.DB.QueryRow("SELECT name, description FROM strategies WHERE id = $1", strategyID).
			Scan(&name, &description)

		if err != nil {
			t.Fatalf("failed to query strategy: %v", err)
		}
		if name != "Breakout Strategy" {
			t.Errorf("expected name 'Breakout Strategy', got %s", name)
		}
		if description != "Trade price breakouts at key levels" {
			t.Errorf("expected description 'Trade price breakouts at key levels', got %s", description)
		}
	})
}

func TestE2E_Strategy_UpdateAndDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	// Register and get token
	authToken := registerUser(t, e, "updatestrategy@example.com", "password123")

	var strategyID int

	// Create strategy
	t.Run("create strategy", func(t *testing.T) {
		payload := map[string]string{
			"name":        "Old Strategy",
			"description": "Old description",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/strategies", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
		}

		var response map[string]any
		json.Unmarshal(rec.Body.Bytes(), &response)
		strategyID = int(response["id"].(float64))
	})

	// Update strategy
	t.Run("update strategy", func(t *testing.T) {
		payload := map[string]string{
			"name":        "New Strategy",
			"description": "New description",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/strategies/%d", strategyID), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var response map[string]any
		json.Unmarshal(rec.Body.Bytes(), &response)

		if response["name"] != "New Strategy" {
			t.Errorf("expected name 'New Strategy', got %v", response["name"])
		}
		if response["description"] != "New description" {
			t.Errorf("expected description 'New description', got %v", response["description"])
		}

		// Verify in database
		var name, description string
		err := pg.DB.QueryRow("SELECT name, description FROM strategies WHERE id = $1", strategyID).
			Scan(&name, &description)

		if err != nil {
			t.Fatalf("failed to query strategy: %v", err)
		}
		if name != "New Strategy" {
			t.Errorf("expected saved name 'New Strategy', got %s", name)
		}
		if description != "New description" {
			t.Errorf("expected saved description 'New description', got %s", description)
		}
	})

	// Delete strategy
	t.Run("delete strategy", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/strategies/%d", strategyID), nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		// Verify deleted from database
		var count int
		err := pg.DB.QueryRow("SELECT COUNT(*) FROM strategies WHERE id = $1", strategyID).Scan(&count)
		if err != nil {
			t.Fatalf("failed to query strategy count: %v", err)
		}
		if count != 0 {
			t.Errorf("expected strategy to be deleted, but found %d strategies", count)
		}
	})
}
