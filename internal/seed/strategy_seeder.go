package seed

import (
	"context"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	strategyapp "github.com/raihanstark/trade-journal/internal/application/strategy"
)

// StrategySeeder handles seeding strategy data
type StrategySeeder struct {
	strategyService *strategyapp.Service
}

// NewStrategySeeder creates a new StrategySeeder instance
func NewStrategySeeder(strategyService *strategyapp.Service) *StrategySeeder {
	return &StrategySeeder{
		strategyService: strategyService,
	}
}

// DefaultStrategies returns common trading strategies
var DefaultStrategies = []struct {
	Name        string
	Description string
}{
	{"Trend Following", "Follow the major trend using moving averages and trend indicators"},
	{"Breakout Trading", "Trade breakouts from consolidation zones with volume confirmation"},
	{"Mean Reversion", "Trade when price deviates significantly from the mean"},
	{"Price Action", "Trade based on candlestick patterns and support/resistance levels"},
	{"Scalping", "Quick in and out trades targeting small profits"},
	{"Swing Trading", "Hold positions for several days to capture larger moves"},
	{"News Trading", "Trade based on economic news and events"},
	{"Grid Trading", "Place buy and sell orders at regular intervals"},
}

// SeedForUser creates strategies for a specific user
func (s *StrategySeeder) SeedForUser(ctx context.Context, userID int64, count int) ([]int64, error) {
	var strategyIDs []int64

	for i := 0; i < count && i < len(DefaultStrategies); i++ {
		strategyDTO, err := s.strategyService.CreateStrategy(ctx, userID, strategyapp.CreateStrategyRequest{
			Name:        DefaultStrategies[i].Name,
			Description: DefaultStrategies[i].Description,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create strategy: %w", err)
		}

		strategyIDs = append(strategyIDs, strategyDTO.ID)
	}

	return strategyIDs, nil
}

// SeedRandomForUser creates random strategies for a specific user
func (s *StrategySeeder) SeedRandomForUser(ctx context.Context, userID int64, count int) ([]int64, error) {
	var strategyIDs []int64

	for i := 0; i < count; i++ {
		strategyDTO, err := s.strategyService.CreateStrategy(ctx, userID, strategyapp.CreateStrategyRequest{
			Name:        gofakeit.JobTitle(),
			Description: gofakeit.Sentence(15),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create strategy: %w", err)
		}

		strategyIDs = append(strategyIDs, strategyDTO.ID)
	}

	return strategyIDs, nil
}
