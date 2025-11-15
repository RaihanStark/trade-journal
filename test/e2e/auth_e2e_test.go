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

func TestE2E_Auth_RegisterAndLogin(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	var token string

	// Register
	t.Run("register new user", func(t *testing.T) {
		payload := map[string]string{
			"email":    "user@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
		}

		var response map[string]any
		json.Unmarshal(rec.Body.Bytes(), &response)

		if response["token"] == nil {
			t.Error("expected token in response")
		}
		token = response["token"].(string)
	})

	// Login
	t.Run("login with registered user", func(t *testing.T) {
		payload := map[string]string{
			"email":    "user@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var response map[string]any
		json.Unmarshal(rec.Body.Bytes(), &response)

		if response["token"] == nil {
			t.Error("expected token in response")
		}
		if response["token"].(string) == token {
			t.Log("tokens match (this is fine - just different generation)")
		}
	})
}

func TestE2E_Auth_InvalidCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	pg := testutil.SetupTestDatabase(t)
	testutil.TruncateTables(t, pg.DB)

	e := setupTestServer(t, pg.DB, pg.Queries)

	// Register user first
	payload := map[string]string{
		"email":    "test@example.com",
		"password": "correctpassword",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Try login with wrong password
	t.Run("login with wrong password fails", func(t *testing.T) {
		payload := map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}
