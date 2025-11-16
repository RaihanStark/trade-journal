package analytics

import (
	"context"
	"testing"

	accountapp "github.com/raihanstark/trade-journal/internal/application/account"
	tradeapp "github.com/raihanstark/trade-journal/internal/application/trade"
	"github.com/raihanstark/trade-journal/internal/domain/user"
	"github.com/raihanstark/trade-journal/internal/infrastructure/persistence"
	"github.com/raihanstark/trade-journal/internal/testutil"
)

func TestAnalyticsService_GetUserAnalytics_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	analyticsRepo := persistence.NewAnalyticsRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)
	accountRepo := persistence.NewAccountRepository(pg.Queries)
	tradeRepo := persistence.NewTradeRepository(pg.Queries)

	analyticsService := NewService(analyticsRepo)
	accountService := accountapp.NewService(accountRepo)
	tradeService := tradeapp.NewService(tradeRepo, accountRepo)

	ctx := context.Background()

	t.Run("calculates analytics from real database trades", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create test user
		testUser := user.NewUser("analytics@example.com", "hashedpass")
		createdUser, err := userRepo.Create(ctx, testUser)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		// Create test account using service
		account, err := accountService.CreateAccount(ctx, createdUser.ID, accountapp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		})
		if err != nil {
			t.Fatalf("failed to create test account: %v", err)
		}

		// Create test trades using service
		trades := []struct {
			tradeType string
			exit      float64
		}{
			{"BUY", 1.1100},  // +100 pips * 1 lot * 10 = +1000 (profit)
			{"SELL", 1.1050}, // Price up 50 pips = loss for SELL = -500
			{"BUY", 1.1200},  // +200 pips * 1 lot * 10 = +2000 (profit)
			{"SELL", 1.1075}, // Price up 75 pips = loss for SELL = -750
		}

		for _, trade := range trades {
			_, err := tradeService.CreateTrade(ctx, createdUser.ID, tradeapp.CreateTradeRequest{
				AccountID: &account.ID,
				Date:      "2024-01-01",
				Time:      "10:00",
				Pair:      "EUR/USD",
				Type:      trade.tradeType,
				Entry:     1.1000,
				Exit:      &trade.exit,
				Lots:      1.0,
			})
			if err != nil {
				t.Fatalf("failed to create trade: %v", err)
			}
		}

		// Call analytics service
		dto, err := analyticsService.GetUserAnalytics(ctx, createdUser.ID)

		// Verify no error
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify analytics calculations
		// Total P/L: 1000 - 500 + 2000 - 750 = 1750
		if dto.TotalPL != 1750 {
			t.Errorf("TotalPL = %v, want 1750", dto.TotalPL)
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

		// Avg win: (1000 + 2000) / 2 = 1500
		if dto.AvgWin != 1500 {
			t.Errorf("AvgWin = %v, want 1500", dto.AvgWin)
		}

		// Avg loss: -(500 + 750) / 2 = -625
		if dto.AvgLoss != -625 {
			t.Errorf("AvgLoss = %v, want -625", dto.AvgLoss)
		}

		// Profit factor: 3000 / 1250 = 2.4
		if dto.ProfitFactor != 2.4 {
			t.Errorf("ProfitFactor = %v, want 2.4", dto.ProfitFactor)
		}

		if dto.LargestWin != 2000 {
			t.Errorf("LargestWin = %v, want 2000", dto.LargestWin)
		}

		if dto.LargestLoss != -750 {
			t.Errorf("LargestLoss = %v, want -750", dto.LargestLoss)
		}
	})

	t.Run("filters out open trades from analytics", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create test user
		testUser := user.NewUser("openfilter@example.com", "hashedpass")
		createdUser, err := userRepo.Create(ctx, testUser)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		// Create test account using service
		account, err := accountService.CreateAccount(ctx, createdUser.ID, accountapp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		})
		if err != nil {
			t.Fatalf("failed to create test account: %v", err)
		}

		// Create closed trade using service
		exit := 1.1100
		_, err = tradeService.CreateTrade(ctx, createdUser.ID, tradeapp.CreateTradeRequest{
			AccountID: &account.ID,
			Date:      "2024-01-01",
			Time:      "10:00",
			Pair:      "EUR/USD",
			Type:      "BUY",
			Entry:     1.1000,
			Exit:      &exit,
			Lots:      1.0,
		})
		if err != nil {
			t.Fatalf("failed to create closed trade: %v", err)
		}

		// Create open trade (no exit) using service
		_, err = tradeService.CreateTrade(ctx, createdUser.ID, tradeapp.CreateTradeRequest{
			AccountID: &account.ID,
			Date:      "2024-01-02",
			Time:      "11:00",
			Pair:      "EUR/USD",
			Type:      "SELL",
			Entry:     1.1100,
			Lots:      1.0,
		})
		if err != nil {
			t.Fatalf("failed to create open trade: %v", err)
		}

		// Call analytics service
		dto, err := analyticsService.GetUserAnalytics(ctx, createdUser.ID)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Should only count the closed trade
		if dto.TotalTrades != 1 {
			t.Errorf("TotalTrades = %v, want 1 (open trades filtered)", dto.TotalTrades)
		}

		// P/L: +100 pips * 1 lot * 10 = 1000
		if dto.TotalPL != 1000 {
			t.Errorf("TotalPL = %v, want 1000", dto.TotalPL)
		}
	})

	t.Run("filters out deposits and withdrawals from analytics", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create test user
		testUser := user.NewUser("depositfilter@example.com", "hashedpass")
		createdUser, err := userRepo.Create(ctx, testUser)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		// Create test account using service
		account, err := accountService.CreateAccount(ctx, createdUser.ID, accountapp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		})
		if err != nil {
			t.Fatalf("failed to create test account: %v", err)
		}

		// Create BUY trade using service
		exit := 1.1100
		_, err = tradeService.CreateTrade(ctx, createdUser.ID, tradeapp.CreateTradeRequest{
			AccountID: &account.ID,
			Date:      "2024-01-01",
			Time:      "10:00",
			Pair:      "EUR/USD",
			Type:      "BUY",
			Entry:     1.1000,
			Exit:      &exit,
			Lots:      1.0,
		})
		if err != nil {
			t.Fatalf("failed to create trade: %v", err)
		}

		// Create DEPOSIT using service (should be filtered)
		depositAmount := 1000.0
		_, err = tradeService.CreateTrade(ctx, createdUser.ID, tradeapp.CreateTradeRequest{
			AccountID: &account.ID,
			Date:      "2024-01-02",
			Time:      "11:00",
			Type:      "DEPOSIT",
			Amount:    &depositAmount,
		})
		if err != nil {
			t.Fatalf("failed to create deposit: %v", err)
		}

		// Create WITHDRAW using service (should be filtered)
		withdrawAmount := 500.0
		_, err = tradeService.CreateTrade(ctx, createdUser.ID, tradeapp.CreateTradeRequest{
			AccountID: &account.ID,
			Date:      "2024-01-03",
			Time:      "12:00",
			Type:      "WITHDRAW",
			Amount:    &withdrawAmount,
		})
		if err != nil {
			t.Fatalf("failed to create withdraw: %v", err)
		}

		// Call analytics service
		dto, err := analyticsService.GetUserAnalytics(ctx, createdUser.ID)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Should only count the BUY trade
		if dto.TotalTrades != 1 {
			t.Errorf("TotalTrades = %v, want 1 (deposits/withdrawals filtered)", dto.TotalTrades)
		}

		// P/L: +100 pips * 1 lot * 10 = 1000
		if dto.TotalPL != 1000 {
			t.Errorf("TotalPL = %v, want 1000", dto.TotalPL)
		}
	})

	t.Run("returns zero analytics when no trades exist", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create test user with no trades
		testUser := user.NewUser("notrades@example.com", "hashedpass")
		createdUser, err := userRepo.Create(ctx, testUser)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		// Call analytics service
		dto, err := analyticsService.GetUserAnalytics(ctx, createdUser.ID)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify zero values
		if dto.TotalTrades != 0 {
			t.Errorf("TotalTrades = %v, want 0", dto.TotalTrades)
		}

		if dto.TotalPL != 0 {
			t.Errorf("TotalPL = %v, want 0", dto.TotalPL)
		}

		if dto.WinRate != 0 {
			t.Errorf("WinRate = %v, want 0", dto.WinRate)
		}

		if dto.ProfitFactor != 0 {
			t.Errorf("ProfitFactor = %v, want 0", dto.ProfitFactor)
		}
	})
}
