package trade

import (
	"context"
	"testing"
	"time"

	accountApp "github.com/raihanstark/trade-journal/internal/application/account"
	strategyApp "github.com/raihanstark/trade-journal/internal/application/strategy"
	"github.com/raihanstark/trade-journal/internal/domain/user"
	"github.com/raihanstark/trade-journal/internal/infrastructure/persistence"
	"github.com/raihanstark/trade-journal/internal/testutil"
)

// Integration tests for trade service
// These test the CRITICAL balance update logic with a real database
// This is where we verify the bug fix actually works end-to-end!

func TestTradeService_CreateTrade_Deposit_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	tradeRepo := persistence.NewTradeRepository(pg.Queries)
	accountRepo := persistence.NewAccountRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)

	tradeService := NewService(tradeRepo, accountRepo)
	accountService := accountApp.NewService(accountRepo)

	ctx := context.Background()

	t.Run("deposit trade updates account balance", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("deposit@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create account with $0 balance
		accountReq := accountApp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account, _ := accountService.CreateAccount(ctx, createdUser.ID, accountReq)

		// Verify initial balance is $0
		if account.CurrentBalance != 0.0 {
			t.Fatalf("expected initial balance 0, got %.2f", account.CurrentBalance)
		}

		// Create deposit trade for $1000
		amount := 1000.0
		tradeReq := CreateTradeRequest{
			AccountID: &account.ID,
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04"),
			Pair:      "USD",
			Type:      "DEPOSIT",
			Entry:     0,
			Lots:      0,
			Amount:    &amount,
		}

		_, err := tradeService.CreateTrade(ctx, createdUser.ID, tradeReq)
		if err != nil {
			t.Fatalf("failed to create deposit trade: %v", err)
		}

		// Verify balance was updated to $1000
		var balance float64
		err = pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", account.ID).Scan(&balance)
		if err != nil {
			t.Fatalf("failed to query balance: %v", err)
		}

		if balance != 1000.0 {
			t.Errorf("expected balance 1000 after deposit, got %.2f", balance)
		}
	})
}

func TestTradeService_CreateTrade_Withdraw_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	tradeRepo := persistence.NewTradeRepository(pg.Queries)
	accountRepo := persistence.NewAccountRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)

	tradeService := NewService(tradeRepo, accountRepo)
	accountService := accountApp.NewService(accountRepo)

	ctx := context.Background()

	t.Run("withdraw trade reduces account balance", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("withdraw@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create account
		accountReq := accountApp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account, _ := accountService.CreateAccount(ctx, createdUser.ID, accountReq)

		// Add initial balance of $1000
		accountRepo.UpdateBalance(ctx, account.ID, createdUser.ID, 1000.0)

		// Create withdraw trade for $300
		amount := 300.0
		tradeReq := CreateTradeRequest{
			AccountID: &account.ID,
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04"),
			Pair:      "USD",
			Type:      "WITHDRAW",
			Entry:     0,
			Lots:      0,
			Amount:    &amount,
		}

		_, err := tradeService.CreateTrade(ctx, createdUser.ID, tradeReq)
		if err != nil {
			t.Fatalf("failed to create withdraw trade: %v", err)
		}

		// Verify balance is now $700 (1000 - 300)
		var balance float64
		err = pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", account.ID).Scan(&balance)
		if err != nil {
			t.Fatalf("failed to query balance: %v", err)
		}

		if balance != 700.0 {
			t.Errorf("expected balance 700 after withdraw, got %.2f", balance)
		}
	})
}

func TestTradeService_CreateTrade_ClosedTrade_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	tradeRepo := persistence.NewTradeRepository(pg.Queries)
	accountRepo := persistence.NewAccountRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)

	tradeService := NewService(tradeRepo, accountRepo)
	accountService := accountApp.NewService(accountRepo)

	ctx := context.Background()

	t.Run("closed BUY trade updates balance with P/L", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("closedtrade@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create account with $1000
		accountReq := accountApp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account, _ := accountService.CreateAccount(ctx, createdUser.ID, accountReq)
		accountRepo.UpdateBalance(ctx, account.ID, createdUser.ID, 1000.0)

		// Create closed BUY trade (50 pips profit * 1 lot = $500 profit)
		exit := 1.1050
		stopLoss := 1.0980
		takeProfit := 1.1060
		tradeReq := CreateTradeRequest{
			AccountID:  &account.ID,
			Date:       time.Now().Format("2006-01-02"),
			Time:       time.Now().Format("15:04"),
			Pair:       "EUR/USD",
			Type:       "BUY",
			Entry:      1.1000,
			Exit:       &exit,
			Lots:       1.0,
			StopLoss:   &stopLoss,
			TakeProfit: &takeProfit,
		}

		trade, err := tradeService.CreateTrade(ctx, createdUser.ID, tradeReq)
		if err != nil {
			t.Fatalf("failed to create trade: %v", err)
		}

		// Verify P/L was calculated (50 pips * 1 lot * $10 = $500)
		if trade.PL == nil {
			t.Fatal("expected P/L to be calculated")
		}
		if *trade.PL != 500.0 {
			t.Errorf("expected P/L 500, got %.2f", *trade.PL)
		}

		// Verify balance is now $1500 (1000 + 500)
		var balance float64
		err = pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", account.ID).Scan(&balance)
		if err != nil {
			t.Fatalf("failed to query balance: %v", err)
		}

		if balance != 1500.0 {
			t.Errorf("expected balance 1500 after profitable trade, got %.2f", balance)
		}
	})

	t.Run("open trade does not update balance", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("opentrade@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create account with $1000
		accountReq := accountApp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account, _ := accountService.CreateAccount(ctx, createdUser.ID, accountReq)
		accountRepo.UpdateBalance(ctx, account.ID, createdUser.ID, 1000.0)

		// Create open BUY trade (no exit price)
		stopLoss := 1.0980
		takeProfit := 1.1060
		tradeReq := CreateTradeRequest{
			AccountID:  &account.ID,
			Date:       time.Now().Format("2006-01-02"),
			Time:       time.Now().Format("15:04"),
			Pair:       "EUR/USD",
			Type:       "BUY",
			Entry:      1.1000,
			Lots:       1.0,
			StopLoss:   &stopLoss,
			TakeProfit: &takeProfit,
		}

		trade, err := tradeService.CreateTrade(ctx, createdUser.ID, tradeReq)
		if err != nil {
			t.Fatalf("failed to create trade: %v", err)
		}

		// Verify P/L is nil (open trade)
		if trade.PL != nil {
			t.Error("expected P/L to be nil for open trade")
		}

		// Verify balance is still $1000 (no change)
		var balance float64
		err = pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", account.ID).Scan(&balance)
		if err != nil {
			t.Fatalf("failed to query balance: %v", err)
		}

		if balance != 1000.0 {
			t.Errorf("expected balance 1000 (unchanged) for open trade, got %.2f", balance)
		}
	})

	t.Run("Create trade with negative P/L should update balance with negative amount", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("negativepl@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create account with $1000
		accountReq := accountApp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account, _ := accountService.CreateAccount(ctx, createdUser.ID, accountReq)
		accountRepo.UpdateBalance(ctx, account.ID, createdUser.ID, 1000.0)

		// Create closed BUY trade (50 pips loss * 1 lot = $500 loss)
		exit := 1.0950
		stopLoss := 1.1020
		takeProfit := 1.0980
		tradeReq := CreateTradeRequest{
			AccountID:  &account.ID,
			Date:       time.Now().Format("2006-01-02"),
			Time:       time.Now().Format("15:04"),
			Pair:       "EUR/USD",
			Type:       "BUY",
			Entry:      1.1000,
			Exit:       &exit,
			Lots:       1.0,
			StopLoss:   &stopLoss,
			TakeProfit: &takeProfit,
		}
		trade, _ := tradeService.CreateTrade(ctx, createdUser.ID, tradeReq)

		// Verify P/L was calculated (50 pips * 1 lot * $10 = $500)
		if trade.PL == nil {
			t.Fatal("expected P/L to be calculated")
		}
		if *trade.PL != -500.0 {
			t.Errorf("expected P/L -500, got %.2f", *trade.PL)
		}

		// Verify balance is now $500 (1000 - 500)
		var balance float64
		err := pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", account.ID).Scan(&balance)
		if err != nil {
			t.Fatalf("failed to query balance: %v", err)
		}

		if balance != 500.0 {
			t.Errorf("expected balance 500 after negative P/L trade, got %.2f", balance)
		}
	})
}

func TestTradeService_UpdateTrade_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	tradeRepo := persistence.NewTradeRepository(pg.Queries)
	accountRepo := persistence.NewAccountRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)

	tradeService := NewService(tradeRepo, accountRepo)
	accountService := accountApp.NewService(accountRepo)

	ctx := context.Background()

	t.Run("updating trade applies P/L difference to balance", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("updatetrade@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create account with $1000
		accountReq := accountApp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account, _ := accountService.CreateAccount(ctx, createdUser.ID, accountReq)
		accountRepo.UpdateBalance(ctx, account.ID, createdUser.ID, 1000.0)

		// Create closed trade with 50 pips profit ($500)
		exit1 := 1.1050
		stopLoss := 1.0980
		takeProfit := 1.1100
		tradeReq := CreateTradeRequest{
			AccountID:  &account.ID,
			Date:       time.Now().Format("2006-01-02"),
			Time:       time.Now().Format("15:04"),
			Pair:       "EUR/USD",
			Type:       "BUY",
			Entry:      1.1000,
			Exit:       &exit1,
			Lots:       1.0,
			StopLoss:   &stopLoss,
			TakeProfit: &takeProfit,
		}
		trade, _ := tradeService.CreateTrade(ctx, createdUser.ID, tradeReq)

		// Balance should be $1500 now
		var balance1 float64
		pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", account.ID).Scan(&balance1)
		if balance1 != 1500.0 {
			t.Fatalf("expected balance 1500 after first trade, got %.2f", balance1)
		}

		// Update trade to have 100 pips profit ($1000 instead of $500)
		exit2 := 1.1100
		updateReq := UpdateTradeRequest{
			AccountID:  &account.ID,
			Date:       time.Now().Format("2006-01-02"),
			Time:       time.Now().Format("15:04"),
			Pair:       "EUR/USD",
			Type:       "BUY",
			Entry:      1.1000,
			Exit:       &exit2, // Changed exit
			Lots:       1.0,
			StopLoss:   &stopLoss,
			TakeProfit: &takeProfit,
		}

		_, err := tradeService.UpdateTrade(ctx, trade.ID, createdUser.ID, updateReq)
		if err != nil {
			t.Fatalf("failed to update trade: %v", err)
		}

		// Balance should be $2000 (1000 + 1000, not 1500 + 1000!)
		// This verifies we only apply the DIFFERENCE ($500)
		var balance2 float64
		pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", account.ID).Scan(&balance2)
		if balance2 != 2000.0 {
			t.Errorf("expected balance 2000 after update (1000 + 1000), got %.2f", balance2)
		}
	})

	t.Run("changing trade account reverts old and applies to new", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("changeaccount@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create two accounts
		account1Req := accountApp.CreateAccountRequest{
			Name:          "Account 1",
			Broker:        "Broker A",
			AccountNumber: "111",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account1, _ := accountService.CreateAccount(ctx, createdUser.ID, account1Req)
		accountRepo.UpdateBalance(ctx, account1.ID, createdUser.ID, 1000.0)

		account2Req := accountApp.CreateAccountRequest{
			Name:          "Account 2",
			Broker:        "Broker B",
			AccountNumber: "222",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account2, _ := accountService.CreateAccount(ctx, createdUser.ID, account2Req)
		accountRepo.UpdateBalance(ctx, account2.ID, createdUser.ID, 2000.0)

		// Create trade on account1 with $500 profit
		exit := 1.1050
		stopLoss := 1.0980
		takeProfit := 1.1060
		tradeReq := CreateTradeRequest{
			AccountID:  &account1.ID,
			Date:       time.Now().Format("2006-01-02"),
			Time:       time.Now().Format("15:04"),
			Pair:       "EUR/USD",
			Type:       "BUY",
			Entry:      1.1000,
			Exit:       &exit,
			Lots:       1.0,
			StopLoss:   &stopLoss,
			TakeProfit: &takeProfit,
		}
		trade, _ := tradeService.CreateTrade(ctx, createdUser.ID, tradeReq)

		// Account1 should be $1500 (1000 + 500)
		var balance1Before float64
		pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", account1.ID).Scan(&balance1Before)
		if balance1Before != 1500.0 {
			t.Fatalf("expected account1 balance 1500, got %.2f", balance1Before)
		}

		// Move trade to account2
		updateReq := UpdateTradeRequest{
			AccountID:  &account2.ID, // Changed account!
			Date:       time.Now().Format("2006-01-02"),
			Time:       time.Now().Format("15:04"),
			Pair:       "EUR/USD",
			Type:       "BUY",
			Entry:      1.1000,
			Exit:       &exit,
			Lots:       1.0,
			StopLoss:   &stopLoss,
			TakeProfit: &takeProfit,
		}

		_, err := tradeService.UpdateTrade(ctx, trade.ID, createdUser.ID, updateReq)
		if err != nil {
			t.Fatalf("failed to update trade: %v", err)
		}

		// Account1 should be $1000 (reverted the $500)
		var balance1After float64
		pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", account1.ID).Scan(&balance1After)
		if balance1After != 1000.0 {
			t.Errorf("expected account1 balance 1000 after moving trade, got %.2f", balance1After)
		}

		// Account2 should be $2500 (2000 + 500)
		var balance2After float64
		pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", account2.ID).Scan(&balance2After)
		if balance2After != 2500.0 {
			t.Errorf("expected account2 balance 2500 after receiving trade, got %.2f", balance2After)
		}
	})
}

func TestTradeService_DeleteTrade_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	tradeRepo := persistence.NewTradeRepository(pg.Queries)
	accountRepo := persistence.NewAccountRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)

	tradeService := NewService(tradeRepo, accountRepo)
	accountService := accountApp.NewService(accountRepo)

	ctx := context.Background()

	t.Run("deleting trade reverts balance", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("deletetrade@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create account with $1000
		accountReq := accountApp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account, _ := accountService.CreateAccount(ctx, createdUser.ID, accountReq)
		accountRepo.UpdateBalance(ctx, account.ID, createdUser.ID, 1000.0)

		// Create closed trade with $500 profit
		exit := 1.1050
		stopLoss := 1.0980
		takeProfit := 1.1060
		tradeReq := CreateTradeRequest{
			AccountID:  &account.ID,
			Date:       time.Now().Format("2006-01-02"),
			Time:       time.Now().Format("15:04"),
			Pair:       "EUR/USD",
			Type:       "BUY",
			Entry:      1.1000,
			Exit:       &exit,
			Lots:       1.0,
			StopLoss:   &stopLoss,
			TakeProfit: &takeProfit,
		}
		trade, _ := tradeService.CreateTrade(ctx, createdUser.ID, tradeReq)

		// Balance should be $1500
		var balanceBefore float64
		pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", account.ID).Scan(&balanceBefore)
		if balanceBefore != 1500.0 {
			t.Fatalf("expected balance 1500 before delete, got %.2f", balanceBefore)
		}

		// Delete the trade
		err := tradeService.DeleteTrade(ctx, trade.ID, createdUser.ID)
		if err != nil {
			t.Fatalf("failed to delete trade: %v", err)
		}

		// Balance should be reverted to $1000
		var balanceAfter float64
		pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", account.ID).Scan(&balanceAfter)
		if balanceAfter != 1000.0 {
			t.Errorf("expected balance 1000 after delete, got %.2f", balanceAfter)
		}

		// Verify trade was deleted
		var count int
		pg.DB.QueryRow("SELECT COUNT(*) FROM trades WHERE id = $1", trade.ID).Scan(&count)
		if count != 0 {
			t.Error("expected trade to be deleted from database")
		}
	})
}

func TestTradeService_WithStrategies_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	tradeRepo := persistence.NewTradeRepository(pg.Queries)
	accountRepo := persistence.NewAccountRepository(pg.Queries)
	strategyRepo := persistence.NewStrategyRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)

	tradeService := NewService(tradeRepo, accountRepo)
	accountService := accountApp.NewService(accountRepo)
	strategyService := strategyApp.NewService(strategyRepo)

	ctx := context.Background()

	t.Run("trade can be associated with strategies", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("strategies@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create account
		accountReq := accountApp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account, _ := accountService.CreateAccount(ctx, createdUser.ID, accountReq)

		// Create strategies
		strategy1Req := strategyApp.CreateStrategyRequest{
			Name:        "Breakout Strategy",
			Description: "Trade breakouts",
		}
		strategy1, _ := strategyService.CreateStrategy(ctx, createdUser.ID, strategy1Req)

		strategy2Req := strategyApp.CreateStrategyRequest{
			Name:        "Support/Resistance",
			Description: "Trade S/R levels",
		}
		strategy2, _ := strategyService.CreateStrategy(ctx, createdUser.ID, strategy2Req)

		// Create trade with both strategies
		exit := 1.1050
		tradeReq := CreateTradeRequest{
			AccountID:   &account.ID,
			Date:        time.Now().Format("2006-01-02"),
			Time:        time.Now().Format("15:04"),
			Pair:        "EUR/USD",
			Type:        "BUY",
			Entry:       1.1000,
			Exit:        &exit,
			Lots:        1.0,
			StrategyIDs: []int64{strategy1.ID, strategy2.ID},
		}

		trade, err := tradeService.CreateTrade(ctx, createdUser.ID, tradeReq)
		if err != nil {
			t.Fatalf("failed to create trade with strategies: %v", err)
		}

		// Verify trade has both strategies
		if len(trade.Strategies) != 2 {
			t.Fatalf("expected 2 strategies, got %d", len(trade.Strategies))
		}

		// Verify junction table was populated
		var count int
		pg.DB.QueryRow("SELECT COUNT(*) FROM trade_strategies WHERE trade_id = $1", trade.ID).Scan(&count)
		if count != 2 {
			t.Errorf("expected 2 entries in trade_strategies junction table, got %d", count)
		}
	})
}

func TestTradeService_GetTradesByAccountID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	tradeRepo := persistence.NewTradeRepository(pg.Queries)
	accountRepo := persistence.NewAccountRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)

	tradeService := NewService(tradeRepo, accountRepo)
	accountService := accountApp.NewService(accountRepo)

	exit := 1.1050
	stopLoss := 1.0980
	takeProfit := 1.1060
	amount := 1000.0

	ctx := context.Background()

	t.Run("get trades by account ID", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("gettradesbyaccountid@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create account
		accountReq := accountApp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		// Create account 2
		account2Req := accountApp.CreateAccountRequest{
			Name:          "Test Account 2",
			Broker:        "Test Broker 2",
			AccountNumber: "456",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account, _ := accountService.CreateAccount(ctx, createdUser.ID, accountReq)
		account2, _ := accountService.CreateAccount(ctx, createdUser.ID, account2Req)

		// Create trade
		tradeReq := CreateTradeRequest{
			AccountID:  &account.ID,
			Date:       time.Now().Format("2006-01-02"),
			Time:       time.Now().Format("15:04"),
			Pair:       "EUR/USD",
			Type:       "BUY",
			Entry:      1.1000,
			Exit:       &exit,
			Lots:       1.0,
			StopLoss:   &stopLoss,
			TakeProfit: &takeProfit,
			Amount:     &amount,
		}
		_, err := tradeService.CreateTrade(ctx, createdUser.ID, tradeReq)
		if err != nil {
			t.Fatalf("failed to create trade: %v", err)
		}

		// Create trade on account2
		trade2Req := CreateTradeRequest{
			AccountID: &account2.ID,
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04"),
			Pair:      "EUR/USD",
			Type:      "BUY",
			Entry:     1.1000,
		}
		_, err = tradeService.CreateTrade(ctx, createdUser.ID, trade2Req)
		if err != nil {
			t.Fatalf("failed to create trade: %v", err)
		}

		// Get trades by account ID
		trades, err := tradeService.GetTradesByAccountID(ctx, account.ID, createdUser.ID)
		if err != nil {
			t.Fatalf("failed to get trades by account ID: %v", err)
		}

		// verify that only get one trade for account1
		if len(trades) != 1 {
			t.Fatalf("expected 1 trade, got %d", len(trades))
		}

		if trades[0].AccountID == nil {
			t.Fatal("expected account ID to be set, got nil")
		}

		if *trades[0].AccountID != account.ID {
			t.Fatalf("expected account ID %d, got %d", account.ID, *trades[0].AccountID)
		}
	})

}

func TestTradeService_DateFiltering_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	tradeRepo := persistence.NewTradeRepository(pg.Queries)
	accountRepo := persistence.NewAccountRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)

	tradeService := NewService(tradeRepo, accountRepo)
	accountService := accountApp.NewService(accountRepo)

	ctx := context.Background()

	t.Run("should filter trades by date range for user", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("datefilter@example.com", "hashedpass")
		createdUser, err := userRepo.Create(ctx, testUser)
		if err != nil {
			t.Fatal(err)
		}

		// Create account
		accountReq := accountApp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account, err := accountService.CreateAccount(ctx, createdUser.ID, accountReq)
		if err != nil {
			t.Fatal(err)
		}

		// Create trades on different dates
		date1 := "2025-01-14"
		date2 := "2025-01-15"
		date3 := "2025-01-16"
		date4 := "2025-01-17"
		tradeTime := "10:00"
		exit1 := 1.1050
		exit2 := 1.1950
		exit3 := 110.50
		exit4 := 0.6950

		// Trade on 2025-01-14 (outside range)
		_, err = tradeService.CreateTrade(ctx, createdUser.ID, CreateTradeRequest{
			AccountID: &account.ID,
			Date:      date1,
			Time:      tradeTime,
			Pair:      "EURUSD",
			Type:      "BUY",
			Entry:     1.1000,
			Exit:      &exit1,
			Lots:      0.1,
		})
		if err != nil {
			t.Fatal(err)
		}

		// Trade on 2025-01-15 (inside range)
		_, err = tradeService.CreateTrade(ctx, createdUser.ID, CreateTradeRequest{
			AccountID: &account.ID,
			Date:      date2,
			Time:      tradeTime,
			Pair:      "GBPUSD",
			Type:      "SELL",
			Entry:     1.2000,
			Exit:      &exit2,
			Lots:      0.2,
		})
		if err != nil {
			t.Fatal(err)
		}

		// Trade on 2025-01-16 (inside range)
		_, err = tradeService.CreateTrade(ctx, createdUser.ID, CreateTradeRequest{
			AccountID: &account.ID,
			Date:      date3,
			Time:      tradeTime,
			Pair:      "USDJPY",
			Type:      "BUY",
			Entry:     110.00,
			Exit:      &exit3,
			Lots:      0.3,
		})
		if err != nil {
			t.Fatal(err)
		}

		// Trade on 2025-01-17 (outside range)
		_, err = tradeService.CreateTrade(ctx, createdUser.ID, CreateTradeRequest{
			AccountID: &account.ID,
			Date:      date4,
			Time:      tradeTime,
			Pair:      "AUDUSD",
			Type:      "SELL",
			Entry:     0.7000,
			Exit:      &exit4,
			Lots:      0.1,
		})
		if err != nil {
			t.Fatal(err)
		}

		// Test with date filter
		startDate := "2025-01-15"
		endDate := "2025-01-16"
		trades, err := tradeService.GetUserTradesWithDateFilter(ctx, createdUser.ID, &startDate, &endDate)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Should only return 2 trades (dates 2025-01-15 and 2025-01-16)
		if len(trades) != 2 {
			t.Fatalf("expected 2 trades, got %d", len(trades))
		}

		// Verify the trades are from the correct dates
		for _, trade := range trades {
			if trade.Date != "2025-01-15" && trade.Date != "2025-01-16" {
				t.Fatalf("expected trade date to be 2025-01-15 or 2025-01-16, got %s", trade.Date)
			}
		}

		// Test without date filter (should return all 4 trades)
		allTrades, err := tradeService.GetUserTradesWithDateFilter(ctx, createdUser.ID, nil, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(allTrades) != 4 {
			t.Fatalf("expected 4 trades, got %d", len(allTrades))
		}
	})

	t.Run("should filter trades by date range for account", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("accountdatefilter@example.com", "hashedpass")
		createdUser, err := userRepo.Create(ctx, testUser)
		if err != nil {
			t.Fatal(err)
		}

		// Create account 1
		accountReq1 := accountApp.CreateAccountRequest{
			Name:          "Test Account 1",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account1, err := accountService.CreateAccount(ctx, createdUser.ID, accountReq1)
		if err != nil {
			t.Fatal(err)
		}

		// Create account 2
		accountReq2 := accountApp.CreateAccountRequest{
			Name:          "Test Account 2",
			Broker:        "Test Broker 2",
			AccountNumber: "456",
			AccountType:   "live",
			Currency:      "USD",
			IsActive:      true,
		}
		account2, err := accountService.CreateAccount(ctx, createdUser.ID, accountReq2)
		if err != nil {
			t.Fatal(err)
		}

		// Create trades on different dates
		date1 := "2025-01-14"
		date2 := "2025-01-15"
		date3 := "2025-01-16"
		tradeTime := "10:00"
		exit1 := 1.1050
		exit2 := 1.1950
		exit3 := 110.50
		exit4 := 0.6950

		// Trade on account 1, date 2025-01-14 (outside range)
		_, err = tradeService.CreateTrade(ctx, createdUser.ID, CreateTradeRequest{
			AccountID: &account1.ID,
			Date:      date1,
			Time:      tradeTime,
			Pair:      "EURUSD",
			Type:      "BUY",
			Entry:     1.1000,
			Exit:      &exit1,
			Lots:      0.1,
		})
		if err != nil {
			t.Fatal(err)
		}

		// Trade on account 1, date 2025-01-15 (inside range)
		_, err = tradeService.CreateTrade(ctx, createdUser.ID, CreateTradeRequest{
			AccountID: &account1.ID,
			Date:      date2,
			Time:      tradeTime,
			Pair:      "GBPUSD",
			Type:      "SELL",
			Entry:     1.2000,
			Exit:      &exit2,
			Lots:      0.2,
		})
		if err != nil {
			t.Fatal(err)
		}

		// Trade on account 2, date 2025-01-15 (inside range but different account)
		_, err = tradeService.CreateTrade(ctx, createdUser.ID, CreateTradeRequest{
			AccountID: &account2.ID,
			Date:      date2,
			Time:      tradeTime,
			Pair:      "USDJPY",
			Type:      "BUY",
			Entry:     110.00,
			Exit:      &exit3,
			Lots:      0.3,
		})
		if err != nil {
			t.Fatal(err)
		}

		// Trade on account 1, date 2025-01-16 (inside range)
		_, err = tradeService.CreateTrade(ctx, createdUser.ID, CreateTradeRequest{
			AccountID: &account1.ID,
			Date:      date3,
			Time:      tradeTime,
			Pair:      "AUDUSD",
			Type:      "SELL",
			Entry:     0.7000,
			Exit:      &exit4,
			Lots:      0.1,
		})
		if err != nil {
			t.Fatal(err)
		}

		// Test with date filter for account 1
		startDate := "2025-01-15"
		endDate := "2025-01-16"
		trades, err := tradeService.GetTradesByAccountIDWithDateFilter(ctx, account1.ID, createdUser.ID, &startDate, &endDate)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Should only return 2 trades from account 1 (dates 2025-01-15 and 2025-01-16)
		if len(trades) != 2 {
			t.Fatalf("expected 2 trades, got %d", len(trades))
		}

		// Verify all trades are from account 1
		for _, trade := range trades {
			if trade.AccountID == nil || *trade.AccountID != account1.ID {
				t.Fatalf("expected trade to be from account %d", account1.ID)
			}
			if trade.Date != "2025-01-15" && trade.Date != "2025-01-16" {
				t.Fatalf("expected trade date to be 2025-01-15 or 2025-01-16, got %s", trade.Date)
			}
		}

		// Test without date filter (should return all 3 trades from account 1)
		allTrades, err := tradeService.GetTradesByAccountIDWithDateFilter(ctx, account1.ID, createdUser.ID, nil, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(allTrades) != 3 {
			t.Fatalf("expected 3 trades, got %d", len(allTrades))
		}
	})
}

func TestTradeService_UpdateChart_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	tradeRepo := persistence.NewTradeRepository(pg.Queries)
	accountRepo := persistence.NewAccountRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)

	tradeService := NewService(tradeRepo, accountRepo)
	accountService := accountApp.NewService(accountRepo)

	ctx := context.Background()

	t.Run("updates chart before successfully", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("chartbefore@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create account
		accountReq := accountApp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account, _ := accountService.CreateAccount(ctx, createdUser.ID, accountReq)

		// Create trade
		exit := 1.1050
		tradeReq := CreateTradeRequest{
			AccountID: &account.ID,
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04"),
			Pair:      "EUR/USD",
			Type:      "BUY",
			Entry:     1.1000,
			Exit:      &exit,
			Lots:      1.0,
		}
		trade, err := tradeService.CreateTrade(ctx, createdUser.ID, tradeReq)
		if err != nil {
			t.Fatalf("failed to create trade: %v", err)
		}

		// Verify initially chart_before is nil
		if trade.ChartBefore != nil {
			t.Error("expected chart_before to be nil initially")
		}

		// Update chart before
		chartURL := "http://localhost:9000/trade-journal/1234567890_chart-before.png"
		updatedTrade, err := tradeService.UpdateChartBefore(ctx, trade.ID, createdUser.ID, chartURL)
		if err != nil {
			t.Fatalf("failed to update chart before: %v", err)
		}

		// Verify chart_before was updated
		if updatedTrade.ChartBefore == nil {
			t.Fatal("expected chart_before to be set")
		}
		if *updatedTrade.ChartBefore != chartURL {
			t.Errorf("expected chart_before URL %s, got %s", chartURL, *updatedTrade.ChartBefore)
		}

		// Verify in database
		var dbChartBefore *string
		err = pg.DB.QueryRow("SELECT chart_before FROM trades WHERE id = $1", trade.ID).Scan(&dbChartBefore)
		if err != nil {
			t.Fatalf("failed to query chart_before from database: %v", err)
		}
		if dbChartBefore == nil {
			t.Fatal("expected chart_before in database to be set")
		}
		if *dbChartBefore != chartURL {
			t.Errorf("expected chart_before in database %s, got %s", chartURL, *dbChartBefore)
		}
	})

	t.Run("updates chart after successfully", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("chartafter@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create account
		accountReq := accountApp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account, _ := accountService.CreateAccount(ctx, createdUser.ID, accountReq)

		// Create trade
		exit := 1.1050
		tradeReq := CreateTradeRequest{
			AccountID: &account.ID,
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04"),
			Pair:      "EUR/USD",
			Type:      "BUY",
			Entry:     1.1000,
			Exit:      &exit,
			Lots:      1.0,
		}
		trade, err := tradeService.CreateTrade(ctx, createdUser.ID, tradeReq)
		if err != nil {
			t.Fatalf("failed to create trade: %v", err)
		}

		// Verify initially chart_after is nil
		if trade.ChartAfter != nil {
			t.Error("expected chart_after to be nil initially")
		}

		// Update chart after
		chartURL := "http://localhost:9000/trade-journal/1234567890_chart-after.png"
		updatedTrade, err := tradeService.UpdateChartAfter(ctx, trade.ID, createdUser.ID, chartURL)
		if err != nil {
			t.Fatalf("failed to update chart after: %v", err)
		}

		// Verify chart_after was updated
		if updatedTrade.ChartAfter == nil {
			t.Fatal("expected chart_after to be set")
		}
		if *updatedTrade.ChartAfter != chartURL {
			t.Errorf("expected chart_after URL %s, got %s", chartURL, *updatedTrade.ChartAfter)
		}

		// Verify in database
		var dbChartAfter *string
		err = pg.DB.QueryRow("SELECT chart_after FROM trades WHERE id = $1", trade.ID).Scan(&dbChartAfter)
		if err != nil {
			t.Fatalf("failed to query chart_after from database: %v", err)
		}
		if dbChartAfter == nil {
			t.Fatal("expected chart_after in database to be set")
		}
		if *dbChartAfter != chartURL {
			t.Errorf("expected chart_after in database %s, got %s", chartURL, *dbChartAfter)
		}
	})

	t.Run("can update both charts", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("bothcharts@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create account
		accountReq := accountApp.CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account, _ := accountService.CreateAccount(ctx, createdUser.ID, accountReq)

		// Create trade
		exit := 1.1050
		tradeReq := CreateTradeRequest{
			AccountID: &account.ID,
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04"),
			Pair:      "EUR/USD",
			Type:      "BUY",
			Entry:     1.1000,
			Exit:      &exit,
			Lots:      1.0,
		}
		trade, err := tradeService.CreateTrade(ctx, createdUser.ID, tradeReq)
		if err != nil {
			t.Fatalf("failed to create trade: %v", err)
		}

		// Update chart before
		chartBeforeURL := "http://localhost:9000/trade-journal/before.png"
		_, err = tradeService.UpdateChartBefore(ctx, trade.ID, createdUser.ID, chartBeforeURL)
		if err != nil {
			t.Fatalf("failed to update chart before: %v", err)
		}

		// Update chart after
		chartAfterURL := "http://localhost:9000/trade-journal/after.png"
		_, err = tradeService.UpdateChartAfter(ctx, trade.ID, createdUser.ID, chartAfterURL)
		if err != nil {
			t.Fatalf("failed to update chart after: %v", err)
		}

		// Verify both charts in database
		var dbChartBefore, dbChartAfter *string
		err = pg.DB.QueryRow("SELECT chart_before, chart_after FROM trades WHERE id = $1", trade.ID).Scan(&dbChartBefore, &dbChartAfter)
		if err != nil {
			t.Fatalf("failed to query charts from database: %v", err)
		}

		if dbChartBefore == nil || *dbChartBefore != chartBeforeURL {
			t.Errorf("expected chart_before %s, got %v", chartBeforeURL, dbChartBefore)
		}
		if dbChartAfter == nil || *dbChartAfter != chartAfterURL {
			t.Errorf("expected chart_after %s, got %v", chartAfterURL, dbChartAfter)
		}
	})
}
