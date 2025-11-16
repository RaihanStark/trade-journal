package seed

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	tradeapp "github.com/raihanstark/trade-journal/internal/application/trade"
)

// TradeSeeder handles seeding trade data
type TradeSeeder struct {
	tradeService *tradeapp.Service
}

// NewTradeSeeder creates a new TradeSeeder instance
func NewTradeSeeder(tradeService *tradeapp.Service) *TradeSeeder {
	return &TradeSeeder{
		tradeService: tradeService,
	}
}

// SeedForAccount creates trades for a specific account
// This will automatically create an initial deposit and update the account balance
func (s *TradeSeeder) SeedForAccount(ctx context.Context, userID, accountID int64, strategyIDs []int64, count int) ([]int64, error) {
	var tradeIDs []int64

	currencyPairs := []string{"EUR/USD", "GBP/USD", "USD/JPY", "AUD/USD", "USD/CAD", "NZD/USD", "EUR/GBP", "EUR/JPY"}
	tradeTypes := []string{"BUY", "SELL"}

	// Start from 90 days ago
	startDate := time.Now().AddDate(0, 0, -90)

	// Create initial deposit for the account
	initialDepositAmount := gofakeit.Float64Range(1000, 10000)
	depositDTO, err := s.tradeService.CreateTrade(ctx, userID, tradeapp.CreateTradeRequest{
		AccountID: &accountID,
		Date:      startDate.Format("2006-01-02"),
		Time:      "09:00",
		Type:      "DEPOSIT",
		Amount:    &initialDepositAmount,
		Entry:     0,
		Lots:      0,
		Notes:     "Initial account deposit",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create initial deposit: %w", err)
	}
	tradeIDs = append(tradeIDs, depositDTO.ID)

	for i := 0; i < count; i++ {
		// Random date within the last 90 days
		daysAgo := rand.Intn(90)
		tradeDate := startDate.AddDate(0, 0, daysAgo)

		pair := currencyPairs[rand.Intn(len(currencyPairs))]
		tradeType := tradeTypes[rand.Intn(len(tradeTypes))]

		entry := gofakeit.Float64Range(1.0, 1.5)
		lots := gofakeit.Float64Range(0.01, 2.0)

		var exitPtr *float64
		// 60% chance of closed trade
		if rand.Float64() < 0.6 {
			// 55% win rate
			isWin := rand.Float64() < 0.55

			pipValue := gofakeit.Float64Range(5, 50)
			if !isWin {
				pipValue = -pipValue
			}

			exitPrice := entry
			if tradeType == "BUY" {
				exitPrice = entry + (pipValue * 0.0001)
			} else {
				exitPrice = entry - (pipValue * 0.0001)
			}
			exitPtr = &exitPrice
		}

		stopLoss := entry - gofakeit.Float64Range(0.001, 0.01)
		takeProfit := entry + gofakeit.Float64Range(0.001, 0.02)

		if tradeType == "SELL" {
			stopLoss = entry + gofakeit.Float64Range(0.001, 0.01)
			takeProfit = entry - gofakeit.Float64Range(0.001, 0.02)
		}

		notes := ""
		if rand.Float64() < 0.3 {
			notes = gofakeit.Sentence(10)
		}

		mistakes := ""
		if exitPtr != nil && *exitPtr < entry && tradeType == "BUY" {
			if rand.Float64() < 0.4 {
				mistakes = gofakeit.Sentence(8)
			}
		}

		// Select 1-2 random strategies
		selectedStrategyIDs := selectRandomStrategies(strategyIDs)

		tradeDTO, err := s.tradeService.CreateTrade(ctx, userID, tradeapp.CreateTradeRequest{
			AccountID:   &accountID,
			Date:        tradeDate.Format("2006-01-02"),
			Time:        fmt.Sprintf("%02d:%02d", rand.Intn(24), rand.Intn(60)),
			Pair:        pair,
			Type:        tradeType,
			Entry:       entry,
			Exit:        exitPtr,
			Lots:        lots,
			StopLoss:    &stopLoss,
			TakeProfit:  &takeProfit,
			Notes:       notes,
			Mistakes:    mistakes,
			StrategyIDs: selectedStrategyIDs,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create trade: %w", err)
		}

		tradeIDs = append(tradeIDs, tradeDTO.ID)
	}

	return tradeIDs, nil
}

// SeedDeposit creates a deposit trade for an account
func (s *TradeSeeder) SeedDeposit(ctx context.Context, userID, accountID int64, amount float64) (int64, error) {
	depositDTO, err := s.tradeService.CreateTrade(ctx, userID, tradeapp.CreateTradeRequest{
		AccountID: &accountID,
		Date:      time.Now().Format("2006-01-02"),
		Time:      time.Now().Format("15:04"),
		Type:      "DEPOSIT",
		Amount:    &amount,
		Entry:     0,
		Lots:      0,
		Notes:     "Account deposit",
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create deposit: %w", err)
	}
	return depositDTO.ID, nil
}

// SeedWithdrawal creates a withdrawal trade for an account
func (s *TradeSeeder) SeedWithdrawal(ctx context.Context, userID, accountID int64, amount float64) (int64, error) {
	withdrawalDTO, err := s.tradeService.CreateTrade(ctx, userID, tradeapp.CreateTradeRequest{
		AccountID: &accountID,
		Date:      time.Now().Format("2006-01-02"),
		Time:      time.Now().Format("15:04"),
		Type:      "WITHDRAW",
		Amount:    &amount,
		Entry:     0,
		Lots:      0,
		Notes:     "Account withdrawal",
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create withdrawal: %w", err)
	}
	return withdrawalDTO.ID, nil
}

// selectRandomStrategies selects 1-2 random strategies from the list
func selectRandomStrategies(strategyIDs []int64) []int64 {
	if len(strategyIDs) == 0 {
		return []int64{}
	}

	numStrategies := rand.Intn(2) + 1
	if numStrategies > len(strategyIDs) {
		numStrategies = len(strategyIDs)
	}

	selectedStrategyIDs := make([]int64, 0, numStrategies)
	usedIndices := make(map[int]bool)

	for len(selectedStrategyIDs) < numStrategies {
		idx := rand.Intn(len(strategyIDs))
		if !usedIndices[idx] {
			selectedStrategyIDs = append(selectedStrategyIDs, strategyIDs[idx])
			usedIndices[idx] = true
		}
	}

	return selectedStrategyIDs
}
