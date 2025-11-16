package analytics

import (
	"context"
	"errors"
	"testing"

	"github.com/raihanstark/trade-journal/internal/db"
)

// AnalyticsRepositorySpy records calls to the analytics repository
type AnalyticsRepositorySpy struct {
	GetUserTradesCalls []int64
	GetUserTradesResult []db.Trade
	GetUserTradesError error
}

func (s *AnalyticsRepositorySpy) GetUserTrades(ctx context.Context, userID int64) ([]db.Trade, error) {
	s.GetUserTradesCalls = append(s.GetUserTradesCalls, userID)
	return s.GetUserTradesResult, s.GetUserTradesError
}

func TestService_GetUserAnalytics_Success(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)

	t.Run("calculates analytics correctly with mixed trades", func(t *testing.T) {
		repoSpy := &AnalyticsRepositorySpy{
			GetUserTradesResult: []db.Trade{
				{Type: db.TradeTypeBUY, Pl: nullString("100")},
				{Type: db.TradeTypeSELL, Pl: nullString("-50")},
				{Type: db.TradeTypeBUY, Pl: nullString("200")},
				{Type: db.TradeTypeSELL, Pl: nullString("-75")},
			},
		}
		service := NewService(repoSpy)

		dto, err := service.GetUserAnalytics(ctx, userID)

		// Assert no error
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Assert repository was called with correct user ID
		if len(repoSpy.GetUserTradesCalls) != 1 {
			t.Fatalf("expected 1 call to GetUserTrades, got %d", len(repoSpy.GetUserTradesCalls))
		}
		if repoSpy.GetUserTradesCalls[0] != userID {
			t.Errorf("expected userID %d, got %d", userID, repoSpy.GetUserTradesCalls[0])
		}

		// Assert analytics calculations
		// Total P/L: 100 - 50 + 200 - 75 = 175
		if dto.TotalPL != 175 {
			t.Errorf("TotalPL = %v, want 175", dto.TotalPL)
		}

		if dto.TotalTrades != 4 {
			t.Errorf("TotalTrades = %v, want 4", dto.TotalTrades)
		}

		if dto.WinningTrades != 2 {
			t.Errorf("WinningTrades = %v, want 2", dto.WinningTrades)
		}

		if dto.LosingTrades != 2 {
			t.Errorf("LosingTrades = %v, want 2", dto.LosingTrades)
		}

		// Win rate: 2/4 * 100 = 50%
		if dto.WinRate != 50 {
			t.Errorf("WinRate = %v, want 50", dto.WinRate)
		}

		// Avg win: (100 + 200) / 2 = 150
		if dto.AvgWin != 150 {
			t.Errorf("AvgWin = %v, want 150", dto.AvgWin)
		}

		// Avg loss: -(50 + 75) / 2 = -62.5
		if dto.AvgLoss != -62.5 {
			t.Errorf("AvgLoss = %v, want -62.5", dto.AvgLoss)
		}

		// Profit factor: 300 / 125 = 2.4
		if dto.ProfitFactor != 2.4 {
			t.Errorf("ProfitFactor = %v, want 2.4", dto.ProfitFactor)
		}

		if dto.LargestWin != 200 {
			t.Errorf("LargestWin = %v, want 200", dto.LargestWin)
		}

		if dto.LargestLoss != -75 {
			t.Errorf("LargestLoss = %v, want -75", dto.LargestLoss)
		}
	})

	t.Run("returns zero analytics when no trades", func(t *testing.T) {
		repoSpy := &AnalyticsRepositorySpy{
			GetUserTradesResult: []db.Trade{},
		}
		service := NewService(repoSpy)

		dto, err := service.GetUserAnalytics(ctx, userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if dto.TotalPL != 0 {
			t.Errorf("TotalPL = %v, want 0", dto.TotalPL)
		}
		if dto.TotalTrades != 0 {
			t.Errorf("TotalTrades = %v, want 0", dto.TotalTrades)
		}
		if dto.WinRate != 0 {
			t.Errorf("WinRate = %v, want 0", dto.WinRate)
		}
	})

	t.Run("filters out open trades", func(t *testing.T) {
		repoSpy := &AnalyticsRepositorySpy{
			GetUserTradesResult: []db.Trade{
				{Type: db.TradeTypeBUY, Pl: nullString("100")},
				{Type: db.TradeTypeSELL, Pl: nullString("")}, // Invalid P/L (open)
				{Type: db.TradeTypeBUY, Pl: nullString("50")},
			},
		}
		service := NewService(repoSpy)

		dto, err := service.GetUserAnalytics(ctx, userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should only count 2 closed trades
		if dto.TotalTrades != 2 {
			t.Errorf("TotalTrades = %v, want 2 (open trades should be filtered)", dto.TotalTrades)
		}

		if dto.TotalPL != 150 {
			t.Errorf("TotalPL = %v, want 150", dto.TotalPL)
		}
	})

	t.Run("filters out deposits and withdrawals", func(t *testing.T) {
		repoSpy := &AnalyticsRepositorySpy{
			GetUserTradesResult: []db.Trade{
				{Type: db.TradeTypeBUY, Pl: nullString("100")},
				{Type: db.TradeTypeDEPOSIT, Pl: nullString("1000")},
				{Type: db.TradeTypeWITHDRAW, Pl: nullString("-500")},
				{Type: db.TradeTypeSELL, Pl: nullString("50")},
			},
		}
		service := NewService(repoSpy)

		dto, err := service.GetUserAnalytics(ctx, userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should only count 2 BUY/SELL trades
		if dto.TotalTrades != 2 {
			t.Errorf("TotalTrades = %v, want 2 (deposits/withdrawals filtered)", dto.TotalTrades)
		}

		// Total P/L should only include BUY/SELL
		if dto.TotalPL != 150 {
			t.Errorf("TotalPL = %v, want 150", dto.TotalPL)
		}
	})
}

func TestService_GetUserAnalytics_Error(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)

	t.Run("propagates repository error", func(t *testing.T) {
		repoSpy := &AnalyticsRepositorySpy{
			GetUserTradesError: errors.New("database connection failed"),
		}
		service := NewService(repoSpy)

		_, err := service.GetUserAnalytics(ctx, userID)

		if err == nil {
			t.Fatal("expected error but got nil")
		}

		if err.Error() != "database connection failed" {
			t.Errorf("expected 'database connection failed', got '%v'", err)
		}
	})
}

func TestService_GetUserAnalytics_OnlyWins(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)

	t.Run("calculates correctly with only winning trades", func(t *testing.T) {
		repoSpy := &AnalyticsRepositorySpy{
			GetUserTradesResult: []db.Trade{
				{Type: db.TradeTypeBUY, Pl: nullString("100")},
				{Type: db.TradeTypeSELL, Pl: nullString("50")},
				{Type: db.TradeTypeBUY, Pl: nullString("75")},
			},
		}
		service := NewService(repoSpy)

		dto, err := service.GetUserAnalytics(ctx, userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if dto.WinningTrades != 3 {
			t.Errorf("WinningTrades = %v, want 3", dto.WinningTrades)
		}
		if dto.LosingTrades != 0 {
			t.Errorf("LosingTrades = %v, want 0", dto.LosingTrades)
		}
		if dto.WinRate != 100 {
			t.Errorf("WinRate = %v, want 100", dto.WinRate)
		}
		if dto.ProfitFactor != 0 {
			t.Errorf("ProfitFactor = %v, want 0 (no losses)", dto.ProfitFactor)
		}
		if dto.TotalPL != 225 {
			t.Errorf("TotalPL = %v, want 225", dto.TotalPL)
		}
	})
}

func TestService_GetUserAnalytics_OnlyLosses(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)

	t.Run("calculates correctly with only losing trades", func(t *testing.T) {
		repoSpy := &AnalyticsRepositorySpy{
			GetUserTradesResult: []db.Trade{
				{Type: db.TradeTypeBUY, Pl: nullString("-100")},
				{Type: db.TradeTypeSELL, Pl: nullString("-50")},
				{Type: db.TradeTypeBUY, Pl: nullString("-75")},
			},
		}
		service := NewService(repoSpy)

		dto, err := service.GetUserAnalytics(ctx, userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if dto.WinningTrades != 0 {
			t.Errorf("WinningTrades = %v, want 0", dto.WinningTrades)
		}
		if dto.LosingTrades != 3 {
			t.Errorf("LosingTrades = %v, want 3", dto.LosingTrades)
		}
		if dto.WinRate != 0 {
			t.Errorf("WinRate = %v, want 0", dto.WinRate)
		}
		if dto.TotalPL != -225 {
			t.Errorf("TotalPL = %v, want -225", dto.TotalPL)
		}
	})
}

