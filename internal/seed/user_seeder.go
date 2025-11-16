package seed

import (
	"context"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/raihanstark/trade-journal/internal/domain/user"
	"github.com/raihanstark/trade-journal/internal/infrastructure/persistence"
	"golang.org/x/crypto/bcrypt"
)

// UserSeeder handles seeding user data
type UserSeeder struct {
	userRepo *persistence.UserRepository
}

// NewUserSeeder creates a new UserSeeder instance
func NewUserSeeder(userRepo *persistence.UserRepository) *UserSeeder {
	return &UserSeeder{
		userRepo: userRepo,
	}
}

// Seed creates the specified number of test users
func (s *UserSeeder) Seed(ctx context.Context, count int) ([]int64, error) {
	var userIDs []int64

	for i := 1; i <= count; i++ {
		email := fmt.Sprintf("test%d@example.com", i)
		password := "password123"

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}

		u := user.NewUser(email, string(hashedPassword))
		created, err := s.userRepo.Create(ctx, u)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		userIDs = append(userIDs, created.ID)
	}

	return userIDs, nil
}

// SeedRandom creates random users with fake data
func (s *UserSeeder) SeedRandom(ctx context.Context, count int) ([]int64, error) {
	var userIDs []int64

	for i := 0; i < count; i++ {
		email := gofakeit.Email()
		password := "password123"

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}

		u := user.NewUser(email, string(hashedPassword))
		created, err := s.userRepo.Create(ctx, u)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		userIDs = append(userIDs, created.ID)
	}

	return userIDs, nil
}
