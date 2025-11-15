package auth

import (
	"context"
	"testing"

	"github.com/raihanstark/trade-journal/internal/infrastructure/persistence"
	"github.com/raihanstark/trade-journal/internal/infrastructure/security"
	"github.com/raihanstark/trade-journal/internal/testutil"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	pg := testutil.SetupTestDatabase(t)

	userRepo := persistence.NewUserRepository(pg.Queries)
	tokenGen := security.NewJWTTokenGenerator("test-secret-key")

	service := NewService(userRepo, tokenGen)

	ctx := context.Background()

	t.Run("register creates user in database", func(t *testing.T) {
		// Clean the database before this test
		testutil.TruncateTables(t, pg.DB)

		// TEST: Register a new user
		req := RegisterRequest{
			Email:    "integration@example.com",
			Password: "mypassword123",
		}

		result, err := service.Register(ctx, req)

		// VERIFY: No error
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// VERIFY: Response has correct data
		if result.User.Email != "integration@example.com" {
			t.Errorf("expected email 'integration@example.com', got %s", result.User.Email)
		}
		if result.Token == "" {
			t.Error("expected non-empty token")
		}

		// VERIFY: User was actually saved to database
		var savedEmail string
		var savedPasswordHash string
		err = pg.DB.QueryRow("SELECT email, password_hash FROM users WHERE id = $1", result.User.ID).
			Scan(&savedEmail, &savedPasswordHash)

		if err != nil {
			t.Fatalf("failed to query saved user: %v", err)
		}

		// Check email matches
		if savedEmail != "integration@example.com" {
			t.Errorf("expected saved email 'integration@example.com', got %s", savedEmail)
		}

		// Check password is hashed (not plain text!)
		if savedPasswordHash == "mypassword123" {
			t.Error("password should be hashed in database, not plain text")
		}

		// Check we can verify the password hash
		err = bcrypt.CompareHashAndPassword([]byte(savedPasswordHash), []byte("mypassword123"))
		if err != nil {
			t.Error("saved password hash should match original password")
		}
	})

	t.Run("register fails with duplicate email (database constraint)", func(t *testing.T) {
		// Clean the database before this test
		testutil.TruncateTables(t, pg.DB)

		// First registration - should succeed
		req := RegisterRequest{
			Email:    "duplicate@example.com",
			Password: "password123",
		}

		_, err := service.Register(ctx, req)
		if err != nil {
			t.Fatalf("first registration should succeed, got error: %v", err)
		}

		// Second registration with same email - should fail
		// This tests that the database UNIQUE constraint works
		_, err = service.Register(ctx, req)

		// VERIFY: Should get error
		if err != ErrEmailAlreadyExists {
			t.Errorf("expected ErrEmailAlreadyExists, got %v", err)
		}

		// VERIFY: Only one user in database
		var count int
		err = pg.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", "duplicate@example.com").Scan(&count)
		if err != nil {
			t.Fatalf("failed to count users: %v", err)
		}

		if count != 1 {
			t.Errorf("expected 1 user in database, got %d", count)
		}
	})
}

func TestAuthService_Login_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)

	userRepo := persistence.NewUserRepository(pg.Queries)
	tokenGen := security.NewJWTTokenGenerator("test-secret-key")

	service := NewService(userRepo, tokenGen)

	ctx := context.Background()

	t.Run("login with valid credentials returns token", func(t *testing.T) {
		// Clean the database
		testutil.TruncateTables(t, pg.DB)

		// SETUP: First, register a user
		registerReq := RegisterRequest{
			Email:    "login@example.com",
			Password: "correctpassword",
		}

		registered, err := service.Register(ctx, registerReq)
		if err != nil {
			t.Fatalf("failed to register user: %v", err)
		}

		// TEST: Now try to login with correct credentials
		loginReq := LoginRequest{
			Email:    "login@example.com",
			Password: "correctpassword",
		}

		result, err := service.Login(ctx, loginReq)

		// VERIFY: No error
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// VERIFY: Response has correct data
		if result.User.ID != registered.User.ID {
			t.Errorf("expected user ID %d, got %d", registered.User.ID, result.User.ID)
		}
		if result.User.Email != "login@example.com" {
			t.Errorf("expected email 'login@example.com', got %s", result.User.Email)
		}
		if result.Token == "" {
			t.Error("expected non-empty token")
		}

		// VERIFY: Token is valid
		claims, err := tokenGen.Validate(result.Token)
		if err != nil {
			t.Errorf("token should be valid: %v", err)
		}
		if claims.UserID != registered.User.ID {
			t.Errorf("token userID mismatch: expected %d, got %d", registered.User.ID, claims.UserID)
		}
	})

	t.Run("login fails with non-existent email", func(t *testing.T) {
		// Clean the database
		testutil.TruncateTables(t, pg.DB)

		// TEST: Try to login with email that doesn't exist in database
		loginReq := LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "anypassword",
		}

		result, err := service.Login(ctx, loginReq)

		// VERIFY: Should get error
		if err != ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}

		// VERIFY: No result
		if result != nil {
			t.Error("expected nil result on error")
		}
	})

	t.Run("login fails with incorrect password", func(t *testing.T) {
		// Clean the database
		testutil.TruncateTables(t, pg.DB)

		// SETUP: Register a user
		registerReq := RegisterRequest{
			Email:    "password@example.com",
			Password: "correctpassword",
		}

		_, err := service.Register(ctx, registerReq)
		if err != nil {
			t.Fatalf("failed to register user: %v", err)
		}

		// TEST: Try to login with wrong password
		loginReq := LoginRequest{
			Email:    "password@example.com",
			Password: "wrongpassword",
		}

		result, err := service.Login(ctx, loginReq)

		// VERIFY: Should get error
		if err != ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}

		// VERIFY: No result
		if result != nil {
			t.Error("expected nil result on error")
		}
	})

	t.Run("login retrieves correct user from database", func(t *testing.T) {
		// Clean the database
		testutil.TruncateTables(t, pg.DB)

		// SETUP: Register multiple users
		user1, _ := service.Register(ctx, RegisterRequest{
			Email:    "user1@example.com",
			Password: "password1",
		})

		user2, _ := service.Register(ctx, RegisterRequest{
			Email:    "user2@example.com",
			Password: "password2",
		})

		// TEST: Login as user2
		result, err := service.Login(ctx, LoginRequest{
			Email:    "user2@example.com",
			Password: "password2",
		})

		// VERIFY: Got user2, not user1
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.User.ID != user2.User.ID {
			t.Errorf("expected user2 ID %d, got %d", user2.User.ID, result.User.ID)
		}
		if result.User.ID == user1.User.ID {
			t.Error("should not get user1's data when logging in as user2")
		}
	})
}
