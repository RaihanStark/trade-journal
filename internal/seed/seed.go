package seed

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	accountapp "github.com/raihanstark/trade-journal/internal/application/account"
	strategyapp "github.com/raihanstark/trade-journal/internal/application/strategy"
	tradeapp "github.com/raihanstark/trade-journal/internal/application/trade"
	"github.com/raihanstark/trade-journal/internal/db"
	"github.com/raihanstark/trade-journal/internal/infrastructure/persistence"
)

// Config holds the configuration for seeding
type Config struct {
	NumUsers          int
	AccountsPerUser   int
	StrategiesPerUser int
	TradesPerAccount  int
}

// DefaultConfig returns the default seeding configuration
func DefaultConfig() Config {
	return Config{
		NumUsers:          3,
		AccountsPerUser:   2,
		StrategiesPerUser: 4,
		TradesPerAccount:  30,
	}
}

// Seeder orchestrates all seeders
type Seeder struct {
	userSeeder     *UserSeeder
	accountSeeder  *AccountSeeder
	strategySeeder *StrategySeeder
	tradeSeeder    *TradeSeeder
	dbConn         *sql.DB
}

// NewSeeder creates a new Seeder instance
func NewSeeder(dbConn *sql.DB) *Seeder {
	queries := db.New(dbConn)
	userRepository := persistence.NewUserRepository(queries)
	accountRepository := persistence.NewAccountRepository(queries)
	strategyRepository := persistence.NewStrategyRepository(queries)
	tradeRepository := persistence.NewTradeRepository(queries)

	// Initialize services
	accountService := accountapp.NewService(accountRepository)
	strategyService := strategyapp.NewService(strategyRepository)
	tradeService := tradeapp.NewService(tradeRepository, accountRepository)

	return &Seeder{
		userSeeder:     NewUserSeeder(userRepository),
		accountSeeder:  NewAccountSeeder(accountService),
		strategySeeder: NewStrategySeeder(strategyService),
		tradeSeeder:    NewTradeSeeder(tradeService),
		dbConn:         dbConn,
	}
}

// Run executes the full seeding process with the given configuration
func (s *Seeder) Run(ctx context.Context, config Config) error {
	log.Println("Starting database seeding...")

	// Seed users
	log.Println("Seeding users...")
	userIDs, err := s.userSeeder.Seed(ctx, config.NumUsers)
	if err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}
	log.Printf("Created %d users", len(userIDs))

	// Seed accounts for each user
	log.Println("Seeding accounts...")
	accountsByUser := make(map[int64][]int64)
	for _, userID := range userIDs {
		accountIDs, err := s.accountSeeder.SeedForUser(ctx, userID, config.AccountsPerUser)
		if err != nil {
			return fmt.Errorf("failed to seed accounts for user %d: %w", userID, err)
		}
		accountsByUser[userID] = accountIDs
	}
	log.Printf("Created %d accounts per user", config.AccountsPerUser)

	// Seed strategies for each user
	log.Println("Seeding strategies...")
	strategiesByUser := make(map[int64][]int64)
	for _, userID := range userIDs {
		strategyIDs, err := s.strategySeeder.SeedForUser(ctx, userID, config.StrategiesPerUser)
		if err != nil {
			return fmt.Errorf("failed to seed strategies for user %d: %w", userID, err)
		}
		strategiesByUser[userID] = strategyIDs
	}
	log.Printf("Created %d strategies per user", config.StrategiesPerUser)

	// Seed trades for each account
	log.Println("Seeding trades...")
	totalTrades := 0
	for userID, accountIDs := range accountsByUser {
		strategyIDs := strategiesByUser[userID]
		for _, accountID := range accountIDs {
			tradeIDs, err := s.tradeSeeder.SeedForAccount(ctx, userID, accountID, strategyIDs, config.TradesPerAccount)
			if err != nil {
				return fmt.Errorf("failed to seed trades for account %d: %w", accountID, err)
			}
			totalTrades += len(tradeIDs)
		}
	}
	log.Printf("Created %d trades total", totalTrades)

	log.Println("Database seeding completed successfully!")
	return nil
}

// ClearData truncates all tables
func (s *Seeder) ClearData(ctx context.Context) error {
	log.Println("Clearing existing data...")

	queries := []string{
		"TRUNCATE TABLE trade_strategies CASCADE",
		"TRUNCATE TABLE trades CASCADE",
		"TRUNCATE TABLE strategies CASCADE",
		"TRUNCATE TABLE accounts CASCADE",
		"TRUNCATE TABLE users CASCADE",
	}

	for _, query := range queries {
		if _, err := s.dbConn.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("failed to execute %s: %w", query, err)
		}
	}

	return nil
}

// Individual seeder access methods for flexibility

// Users returns the user seeder
func (s *Seeder) Users() *UserSeeder {
	return s.userSeeder
}

// Accounts returns the account seeder
func (s *Seeder) Accounts() *AccountSeeder {
	return s.accountSeeder
}

// Strategies returns the strategy seeder
func (s *Seeder) Strategies() *StrategySeeder {
	return s.strategySeeder
}

// Trades returns the trade seeder
func (s *Seeder) Trades() *TradeSeeder {
	return s.tradeSeeder
}
