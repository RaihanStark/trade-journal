package auth

import (
	"context"
	"database/sql"
	"testing"

	"github.com/raihanstark/trade-journal/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// UserRepositorySpy is a spy implementation of user.Repository
type UserRepositorySpy struct {
	// Recorded calls
	CreateCalls     []*user.User
	GetByEmailCalls []string
	GetByIDCalls    []int64

	// Configured responses
	CreateResult     *user.User
	CreateError      error
	GetByEmailResult *user.User
	GetByEmailError  error
	GetByIDResult    *user.User
	GetByIDError     error
}

func (s *UserRepositorySpy) Create(ctx context.Context, u *user.User) (*user.User, error) {
	s.CreateCalls = append(s.CreateCalls, u)
	return s.CreateResult, s.CreateError
}

func (s *UserRepositorySpy) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	s.GetByEmailCalls = append(s.GetByEmailCalls, email)
	return s.GetByEmailResult, s.GetByEmailError
}

func (s *UserRepositorySpy) GetByID(ctx context.Context, id int64) (*user.User, error) {
	s.GetByIDCalls = append(s.GetByIDCalls, id)
	return s.GetByIDResult, s.GetByIDError
}

// TokenGeneratorSpy is a spy implementation of TokenGenerator
type TokenGeneratorSpy struct {
	// Recorded calls
	GenerateCalls []GenerateCall

	// Configured responses
	GenerateResult string
	GenerateError  error
}

type GenerateCall struct {
	UserID int64
	Email  string
}

func (s *TokenGeneratorSpy) Generate(userID int64, email string) (string, error) {
	s.GenerateCalls = append(s.GenerateCalls, GenerateCall{
		UserID: userID,
		Email:  email,
	})
	return s.GenerateResult, s.GenerateError
}

func TestService_Register(t *testing.T) {
	ctx := context.Background()

	t.Run("successful registration", func(t *testing.T) {
		userSpy := &UserRepositorySpy{
			CreateResult: &user.User{
				ID:    1,
				Email: "test@example.com",
			},
		}
		tokenSpy := &TokenGeneratorSpy{
			GenerateResult: "test-token-123",
		}
		service := NewService(userSpy, tokenSpy)

		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		result, err := service.Register(ctx, req)

		// Assert no error
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Assert user repository was called
		if len(userSpy.CreateCalls) != 1 {
			t.Fatalf("expected 1 call to Create, got %d", len(userSpy.CreateCalls))
		}

		// Assert email is correct
		createdUser := userSpy.CreateCalls[0]
		if createdUser.Email != "test@example.com" {
			t.Errorf("expected email 'test@example.com', got %s", createdUser.Email)
		}

		// Assert password is hashed (not plain text)
		if createdUser.PasswordHash == "password123" {
			t.Error("password should be hashed, not plain text")
		}

		// Assert password hash can be verified
		err = bcrypt.CompareHashAndPassword([]byte(createdUser.PasswordHash), []byte("password123"))
		if err != nil {
			t.Error("password hash verification failed")
		}

		// Assert token generator was called
		if len(tokenSpy.GenerateCalls) != 1 {
			t.Fatalf("expected 1 call to Generate, got %d", len(tokenSpy.GenerateCalls))
		}

		tokenCall := tokenSpy.GenerateCalls[0]
		if tokenCall.UserID != 1 {
			t.Errorf("expected userID 1, got %d", tokenCall.UserID)
		}
		if tokenCall.Email != "test@example.com" {
			t.Errorf("expected email 'test@example.com', got %s", tokenCall.Email)
		}

		// Assert response
		if result.Token != "test-token-123" {
			t.Errorf("expected token 'test-token-123', got %s", result.Token)
		}
		if result.User.ID != 1 {
			t.Errorf("expected user ID 1, got %d", result.User.ID)
		}
		if result.User.Email != "test@example.com" {
			t.Errorf("expected email 'test@example.com', got %s", result.User.Email)
		}
	})

	t.Run("registration fails with duplicate email", func(t *testing.T) {
		userSpy := &UserRepositorySpy{
			CreateError: sql.ErrConnDone, // Simulate DB error
		}
		tokenSpy := &TokenGeneratorSpy{}
		service := NewService(userSpy, tokenSpy)

		req := RegisterRequest{
			Email:    "existing@example.com",
			Password: "password123",
		}

		result, err := service.Register(ctx, req)

		// Assert error
		if err != ErrEmailAlreadyExists {
			t.Errorf("expected ErrEmailAlreadyExists, got %v", err)
		}

		// Assert result is nil
		if result != nil {
			t.Error("expected nil result on error")
		}

		// Assert token generator was not called
		if len(tokenSpy.GenerateCalls) != 0 {
			t.Error("token generator should not be called on registration failure")
		}
	})

	t.Run("registration fails when token generation fails", func(t *testing.T) {
		userSpy := &UserRepositorySpy{
			CreateResult: &user.User{
				ID:    1,
				Email: "test@example.com",
			},
		}
		tokenSpy := &TokenGeneratorSpy{
			GenerateError: sql.ErrConnDone, // Simulate token generation error
		}
		service := NewService(userSpy, tokenSpy)

		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		result, err := service.Register(ctx, req)

		// Assert error
		if err != ErrGeneratingToken {
			t.Errorf("expected ErrGeneratingToken, got %v", err)
		}

		// Assert result is nil
		if result != nil {
			t.Error("expected nil result on error")
		}
	})
}

func TestService_Login(t *testing.T) {
	ctx := context.Background()

	// Create a hashed password for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	t.Run("successful login with valid credentials", func(t *testing.T) {
		userSpy := &UserRepositorySpy{
			GetByEmailResult: &user.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: string(hashedPassword),
			},
		}
		tokenSpy := &TokenGeneratorSpy{
			GenerateResult: "login-token-456",
		}
		service := NewService(userSpy, tokenSpy)

		req := LoginRequest{
			Email:    "test@example.com",
			Password: "correctpassword",
		}

		result, err := service.Login(ctx, req)

		// Assert no error
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Assert GetByEmail was called
		if len(userSpy.GetByEmailCalls) != 1 {
			t.Fatalf("expected 1 call to GetByEmail, got %d", len(userSpy.GetByEmailCalls))
		}
		if userSpy.GetByEmailCalls[0] != "test@example.com" {
			t.Errorf("expected email 'test@example.com', got %s", userSpy.GetByEmailCalls[0])
		}

		// Assert token generator was called
		if len(tokenSpy.GenerateCalls) != 1 {
			t.Fatalf("expected 1 call to Generate, got %d", len(tokenSpy.GenerateCalls))
		}

		tokenCall := tokenSpy.GenerateCalls[0]
		if tokenCall.UserID != 1 {
			t.Errorf("expected userID 1, got %d", tokenCall.UserID)
		}

		// Assert response
		if result.Token != "login-token-456" {
			t.Errorf("expected token 'login-token-456', got %s", result.Token)
		}
		if result.User.ID != 1 {
			t.Errorf("expected user ID 1, got %d", result.User.ID)
		}
	})

	t.Run("login fails with non-existent email", func(t *testing.T) {
		userSpy := &UserRepositorySpy{
			GetByEmailError: sql.ErrNoRows,
		}
		tokenSpy := &TokenGeneratorSpy{}
		service := NewService(userSpy, tokenSpy)

		req := LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "anypassword",
		}

		result, err := service.Login(ctx, req)

		// Assert error
		if err != ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}

		// Assert result is nil
		if result != nil {
			t.Error("expected nil result on error")
		}

		// Assert token generator was not called
		if len(tokenSpy.GenerateCalls) != 0 {
			t.Error("token generator should not be called on login failure")
		}
	})

	t.Run("login fails with incorrect password", func(t *testing.T) {
		userSpy := &UserRepositorySpy{
			GetByEmailResult: &user.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: string(hashedPassword),
			},
		}
		tokenSpy := &TokenGeneratorSpy{}
		service := NewService(userSpy, tokenSpy)

		req := LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		result, err := service.Login(ctx, req)

		// Assert error
		if err != ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}

		// Assert result is nil
		if result != nil {
			t.Error("expected nil result on error")
		}

		// Assert token generator was not called
		if len(tokenSpy.GenerateCalls) != 0 {
			t.Error("token generator should not be called with wrong password")
		}
	})

	t.Run("login fails when token generation fails", func(t *testing.T) {
		userSpy := &UserRepositorySpy{
			GetByEmailResult: &user.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: string(hashedPassword),
			},
		}
		tokenSpy := &TokenGeneratorSpy{
			GenerateError: sql.ErrConnDone,
		}
		service := NewService(userSpy, tokenSpy)

		req := LoginRequest{
			Email:    "test@example.com",
			Password: "correctpassword",
		}

		result, err := service.Login(ctx, req)

		// Assert error
		if err != ErrGeneratingToken {
			t.Errorf("expected ErrGeneratingToken, got %v", err)
		}

		// Assert result is nil
		if result != nil {
			t.Error("expected nil result on error")
		}
	})
}
