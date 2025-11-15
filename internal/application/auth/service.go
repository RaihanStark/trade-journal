package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/raihanstark/trade-journal/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrHashingPassword    = errors.New("failed to hash password")
	ErrGeneratingToken    = errors.New("failed to generate token")
)

// TokenGenerator defines the interface for generating JWT tokens
type TokenGenerator interface {
	Generate(userID int64, email string) (string, error)
}

// Service handles authentication use cases
type Service struct {
	userRepo       user.Repository
	tokenGenerator TokenGenerator
}

// NewService creates a new authentication service
func NewService(userRepo user.Repository, tokenGenerator TokenGenerator) *Service {
	return &Service{
		userRepo:       userRepo,
		tokenGenerator: tokenGenerator,
	}
}

// Register registers a new user
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrHashingPassword
	}

	// Create user entity
	newUser := user.NewUser(req.Email, string(hashedPassword))

	// Save to repository
	createdUser, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Generate token
	token, err := s.tokenGenerator.Generate(createdUser.ID, createdUser.Email)
	if err != nil {
		return nil, ErrGeneratingToken
	}

	return &AuthResponse{
		Token: token,
		User: UserDTO{
			ID:    createdUser.ID,
			Email: createdUser.Email,
		},
	}, nil
}

// Login authenticates a user
func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// Get user by email
	foundUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, err := s.tokenGenerator.Generate(foundUser.ID, foundUser.Email)
	if err != nil {
		return nil, ErrGeneratingToken
	}

	return &AuthResponse{
		Token: token,
		User: UserDTO{
			ID:    foundUser.ID,
			Email: foundUser.Email,
		},
	}, nil
}
