package seed

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"
	accountapp "github.com/raihanstark/trade-journal/internal/application/account"
)

// AccountSeeder handles seeding account data
type AccountSeeder struct {
	accountService *accountapp.Service
}

// NewAccountSeeder creates a new AccountSeeder instance
func NewAccountSeeder(accountService *accountapp.Service) *AccountSeeder {
	return &AccountSeeder{
		accountService: accountService,
	}
}

// SeedForUser creates accounts for a specific user
func (s *AccountSeeder) SeedForUser(ctx context.Context, userID int64, count int) ([]int64, error) {
	var accountIDs []int64
	accountTypes := []string{"demo", "live"}
	brokers := []string{"MetaTrader 4", "MetaTrader 5", "cTrader", "Interactive Brokers", "TD Ameritrade"}

	for i := 0; i < count; i++ {
		accountType := accountTypes[i%len(accountTypes)]
		broker := brokers[rand.Intn(len(brokers))]

		accountDTO, err := s.accountService.CreateAccount(ctx, userID, accountapp.CreateAccountRequest{
			Name:          fmt.Sprintf("%s Account %d", accountType, i+1),
			Broker:        broker,
			AccountNumber: gofakeit.UUID(),
			AccountType:   accountType,
			Currency:      "USD",
			IsActive:      true,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create account: %w", err)
		}

		accountIDs = append(accountIDs, accountDTO.ID)
	}

	return accountIDs, nil
}

// SeedRandomForUser creates random accounts for a specific user
func (s *AccountSeeder) SeedRandomForUser(ctx context.Context, userID int64, count int) ([]int64, error) {
	var accountIDs []int64
	accountTypes := []string{"demo", "live"}
	brokers := []string{"MetaTrader 4", "MetaTrader 5", "cTrader", "Interactive Brokers", "TD Ameritrade", "NinjaTrader"}

	for i := 0; i < count; i++ {
		accountType := accountTypes[rand.Intn(len(accountTypes))]
		broker := brokers[rand.Intn(len(brokers))]

		accountDTO, err := s.accountService.CreateAccount(ctx, userID, accountapp.CreateAccountRequest{
			Name:          gofakeit.BS(),
			Broker:        broker,
			AccountNumber: gofakeit.UUID(),
			AccountType:   accountType,
			Currency:      "USD",
			IsActive:      gofakeit.Bool(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create account: %w", err)
		}

		accountIDs = append(accountIDs, accountDTO.ID)
	}

	return accountIDs, nil
}
