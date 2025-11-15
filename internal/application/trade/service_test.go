package trade

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/raihanstark/trade-journal/internal/domain/account"
	tradedom "github.com/raihanstark/trade-journal/internal/domain/trade"
)

// TradeRepositorySpy records calls to the trade repository
type TradeRepositorySpy struct {
	CreateCalls  []*tradedom.Trade
	GetByIDCalls []GetByIDCall
	UpdateCalls  []*tradedom.Trade
	DeleteCalls  []DeleteCall

	CreateResult  *tradedom.Trade
	CreateError   error
	GetByIDResult *tradedom.Trade
	GetByIDError  error
	UpdateResult  *tradedom.Trade
	UpdateError   error
	DeleteError   error
}

type GetByIDCall struct {
	ID     int64
	UserID int64
}

type DeleteCall struct {
	ID     int64
	UserID int64
}

func (s *TradeRepositorySpy) Create(ctx context.Context, trade *tradedom.Trade) (*tradedom.Trade, error) {
	s.CreateCalls = append(s.CreateCalls, trade)
	return s.CreateResult, s.CreateError
}

func (s *TradeRepositorySpy) GetByID(ctx context.Context, id int64, userID int64) (*tradedom.Trade, error) {
	s.GetByIDCalls = append(s.GetByIDCalls, GetByIDCall{ID: id, UserID: userID})
	return s.GetByIDResult, s.GetByIDError
}

func (s *TradeRepositorySpy) GetByUserID(ctx context.Context, userID int64) ([]*tradedom.Trade, error) {
	return nil, errors.New("not implemented")
}

func (s *TradeRepositorySpy) Update(ctx context.Context, trade *tradedom.Trade) (*tradedom.Trade, error) {
	s.UpdateCalls = append(s.UpdateCalls, trade)
	return s.UpdateResult, s.UpdateError
}

func (s *TradeRepositorySpy) Delete(ctx context.Context, id int64, userID int64) error {
	s.DeleteCalls = append(s.DeleteCalls, DeleteCall{ID: id, UserID: userID})
	return s.DeleteError
}

// AccountRepositorySpy records calls to the account repository
type AccountRepositorySpy struct {
	UpdateBalanceCalls []UpdateBalanceCall

	UpdateBalanceResult *account.Account
	UpdateBalanceError  error
}

type UpdateBalanceCall struct {
	ID     int64
	UserID int64
	Amount float64
}

func (s *AccountRepositorySpy) UpdateBalance(ctx context.Context, id int64, userID int64, amount float64) (*account.Account, error) {
	s.UpdateBalanceCalls = append(s.UpdateBalanceCalls, UpdateBalanceCall{
		ID:     id,
		UserID: userID,
		Amount: amount,
	})
	return s.UpdateBalanceResult, s.UpdateBalanceError
}

func (s *AccountRepositorySpy) Create(ctx context.Context, acc *account.Account) (*account.Account, error) {
	return nil, errors.New("not implemented")
}

func (s *AccountRepositorySpy) GetByID(ctx context.Context, id int64, userID int64) (*account.Account, error) {
	return nil, errors.New("not implemented")
}

func (s *AccountRepositorySpy) GetByUserID(ctx context.Context, userID int64) ([]*account.Account, error) {
	return nil, errors.New("not implemented")
}

func (s *AccountRepositorySpy) Update(ctx context.Context, acc *account.Account) (*account.Account, error) {
	return nil, errors.New("not implemented")
}

func (s *AccountRepositorySpy) Delete(ctx context.Context, id int64, userID int64) error {
	return errors.New("not implemented")
}

func TestService_CreateTrade_Validation(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)

	t.Run("account_id is required for creating trade", func(t *testing.T) {
		tradeSpy := &TradeRepositorySpy{}
		accountSpy := &AccountRepositorySpy{}
		service := NewService(tradeSpy, accountSpy)

		amount := 1000.0
		_, err := service.CreateTrade(ctx, userID, CreateTradeRequest{
			AccountID: nil, // Missing account_id
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04"),
			Type:      "DEPOSIT",
			Amount:    &amount,
		})

		// Assert error
		if err == nil {
			t.Fatal("expected error when account_id is nil, got nil")
		}
		if err.Error() != "account_id is required" {
			t.Errorf("expected error 'account_id is required', got '%v'", err)
		}

		// Verify no repository calls were made
		if len(tradeSpy.CreateCalls) != 0 {
			t.Errorf("expected 0 calls to Create when validation fails, got %d", len(tradeSpy.CreateCalls))
		}
	})
}

func TestService_CreateTrade_Deposit(t *testing.T) {
	ctx := context.Background()
	accountID := int64(1)
	userID := int64(1)
	amount := 1000.0

	t.Run("deposit updates account balance", func(t *testing.T) {
		tradeSpy := &TradeRepositorySpy{
			CreateResult: &tradedom.Trade{
				ID:        1,
				UserID:    userID,
				AccountID: &accountID,
				Type:      tradedom.TradeTypeDeposit,
				Amount:    &amount,
			},
		}
		accountSpy := &AccountRepositorySpy{}
		service := NewService(tradeSpy, accountSpy)

		_, err := service.CreateTrade(ctx, userID, CreateTradeRequest{
			AccountID: &accountID,
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04"),
			Type:      "DEPOSIT",
			Amount:    &amount,
		})

		// Assert no error
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Assert trade was created
		if len(tradeSpy.CreateCalls) != 1 {
			t.Fatalf("expected 1 call to Create, got %d", len(tradeSpy.CreateCalls))
		}

		// Assert account balance was updated with positive amount
		if len(accountSpy.UpdateBalanceCalls) != 1 {
			t.Fatalf("expected 1 call to UpdateBalance, got %d", len(accountSpy.UpdateBalanceCalls))
		}

		balanceCall := accountSpy.UpdateBalanceCalls[0]
		if balanceCall.ID != accountID {
			t.Errorf("expected account ID %d, got %d", accountID, balanceCall.ID)
		}
		if balanceCall.UserID != userID {
			t.Errorf("expected user ID %d, got %d", userID, balanceCall.UserID)
		}
		if balanceCall.Amount != 1000.0 {
			t.Errorf("expected amount 1000.0, got %.2f", balanceCall.Amount)
		}
	})
}

func TestService_CreateTrade_Withdraw(t *testing.T) {
	ctx := context.Background()
	accountID := int64(1)
	userID := int64(1)
	amount := 500.0

	t.Run("withdraw updates account balance with negative amount", func(t *testing.T) {
		tradeSpy := &TradeRepositorySpy{
			CreateResult: &tradedom.Trade{
				ID:        1,
				UserID:    userID,
				AccountID: &accountID,
				Type:      tradedom.TradeTypeWithdraw,
				Amount:    &amount,
			},
		}
		accountSpy := &AccountRepositorySpy{}
		service := NewService(tradeSpy, accountSpy)

		_, err := service.CreateTrade(ctx, userID, CreateTradeRequest{
			AccountID: &accountID,
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04"),
			Type:      "WITHDRAW",
			Amount:    &amount,
		})

		// Assert no error
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Assert account balance was updated with negative amount
		if len(accountSpy.UpdateBalanceCalls) != 1 {
			t.Fatalf("expected 1 call to UpdateBalance, got %d", len(accountSpy.UpdateBalanceCalls))
		}

		balanceCall := accountSpy.UpdateBalanceCalls[0]
		if balanceCall.Amount != -500.0 {
			t.Errorf("expected amount -500.0 (withdraw), got %.2f", balanceCall.Amount)
		}
	})
}

func TestService_CreateTrade_ClosedTrade(t *testing.T) {
	ctx := context.Background()
	accountID := int64(1)
	userID := int64(1)
	exit := 1.1050

	t.Run("closed BUY trade updates account balance with P/L", func(t *testing.T) {
		tradeSpy := &TradeRepositorySpy{
			CreateResult: &tradedom.Trade{
				ID:        1,
				UserID:    userID,
				AccountID: &accountID,
				Type:      tradedom.TradeTypeBuy,
				Entry:     1.1000,
				Exit:      &exit,
				Status:    tradedom.TradeStatusClosed,
			},
		}
		accountSpy := &AccountRepositorySpy{}
		service := NewService(tradeSpy, accountSpy)

		_, err := service.CreateTrade(ctx, userID, CreateTradeRequest{
			AccountID: &accountID,
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04"),
			Pair:      "EUR/USD",
			Type:      "BUY",
			Entry:     1.1000,
			Exit:      &exit,
			Lots:      1.0,
		})

		// Assert no error
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Assert account balance was updated with P/L
		if len(accountSpy.UpdateBalanceCalls) != 1 {
			t.Fatalf("expected 1 call to UpdateBalance, got %d", len(accountSpy.UpdateBalanceCalls))
		}

		balanceCall := accountSpy.UpdateBalanceCalls[0]
		// P/L = 50 pips * 1 lot * $10 = $500
		if balanceCall.Amount != 500.0 {
			t.Errorf("expected amount 500.0 (P/L), got %.2f", balanceCall.Amount)
		}
	})

	t.Run("open trade does not update account balance", func(t *testing.T) {
		tradeSpy := &TradeRepositorySpy{
			CreateResult: &tradedom.Trade{
				ID:        1,
				UserID:    userID,
				AccountID: &accountID,
				Type:      tradedom.TradeTypeBuy,
				Entry:     1.1000,
				Status:    tradedom.TradeStatusOpen,
			},
		}
		accountSpy := &AccountRepositorySpy{}
		service := NewService(tradeSpy, accountSpy)

		_, err := service.CreateTrade(ctx, userID, CreateTradeRequest{
			AccountID: &accountID,
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04"),
			Pair:      "EUR/USD",
			Type:      "BUY",
			Entry:     1.1000,
			Lots:      1.0,
		})

		// Assert no error
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Assert account balance was NOT updated
		if len(accountSpy.UpdateBalanceCalls) != 0 {
			t.Errorf("expected 0 calls to UpdateBalance for open trade, got %d", len(accountSpy.UpdateBalanceCalls))
		}
	})

	t.Run("Create trade with negative P/L should update balance with negative amount", func(t *testing.T) {
		exit := float64(0.5000)
		pips := -5000.0
		lotSize := 1.0
		pl := pips * lotSize * 10.0
		tradeSpy := &TradeRepositorySpy{
			CreateResult: &tradedom.Trade{
				ID:        1,
				UserID:    userID,
				AccountID: &accountID,
				Type:      tradedom.TradeTypeBuy,
				Entry:     1.0000,
				Exit:      &exit,
				Lots:      lotSize,
				PL:        &pl,
			},
		}
		accountSpy := &AccountRepositorySpy{}
		service := NewService(tradeSpy, accountSpy)

		_, err := service.CreateTrade(ctx, userID, CreateTradeRequest{
			AccountID: &accountID,
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04"),
			Pair:      "EUR/USD",
			Type:      "BUY",
			Entry:     1.0000,
			Exit:      &exit,
			Lots:      lotSize,
		})

		// Assert no error
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Assert account balance was updated with negative P/L
		if len(accountSpy.UpdateBalanceCalls) != 1 {
			t.Fatalf("expected 1 call to UpdateBalance, got %d", len(accountSpy.UpdateBalanceCalls))
		}

		balanceCall := accountSpy.UpdateBalanceCalls[0]
		if balanceCall.Amount != pl {
			t.Errorf("expected amount %.2f (negative P/L), got %.2f", pl, balanceCall.Amount)
		}
	})
}

func TestService_UpdateTrade_PLDifference(t *testing.T) {
	ctx := context.Background()
	accountID := int64(1)
	userID := int64(1)
	tradeID := int64(1)

	t.Run("updating trade P/L applies difference to balance", func(t *testing.T) {
		// Existing trade: 50 pips * 1 lot * $10 = $500 P/L
		oldExit := 1.1050

		// Updated trade: 100 pips * 1 lot * $10 = $1000 P/L
		newExit := 1.1100

		tradeSpy := &TradeRepositorySpy{
			GetByIDResult: &tradedom.Trade{
				ID:        tradeID,
				UserID:    userID,
				AccountID: &accountID,
				Type:      tradedom.TradeTypeBuy,
				Entry:     1.1000,
				Exit:      &oldExit,
				Lots:      1.0,
			},
			UpdateResult: &tradedom.Trade{
				ID:        tradeID,
				UserID:    userID,
				AccountID: &accountID,
				Type:      tradedom.TradeTypeBuy,
				Entry:     1.1000,
				Exit:      &newExit,
				Lots:      1.0,
			},
		}
		// Set P/L on existing trade (calculated from CalculateTradeMetrics)
		oldPL := 500.0 // 50 pips * 1 lot * $10
		tradeSpy.GetByIDResult.PL = &oldPL

		accountSpy := &AccountRepositorySpy{}
		service := NewService(tradeSpy, accountSpy)

		_, err := service.UpdateTrade(ctx, tradeID, userID, UpdateTradeRequest{
			AccountID: &accountID,
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04"),
			Pair:      "EUR/USD",
			Type:      "BUY",
			Entry:     1.1000,
			Exit:      &newExit,
			Lots:      1.0,
		})

		// Assert no error
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Assert balance was updated with difference ($1000 - $500 = $500)
		if len(accountSpy.UpdateBalanceCalls) != 1 {
			t.Fatalf("expected 1 call to UpdateBalance, got %d", len(accountSpy.UpdateBalanceCalls))
		}

		balanceCall := accountSpy.UpdateBalanceCalls[0]
		expectedDifference := 500.0 // new P/L ($1000) - old P/L ($500)
		if balanceCall.Amount != expectedDifference {
			t.Errorf("expected balance difference %.2f, got %.2f", expectedDifference, balanceCall.Amount)
		}
	})

	t.Run("updating trade with same P/L does not update balance", func(t *testing.T) {
		exit := 1.1050

		tradeSpy := &TradeRepositorySpy{
			GetByIDResult: &tradedom.Trade{
				ID:        tradeID,
				UserID:    userID,
				AccountID: &accountID,
				Type:      tradedom.TradeTypeBuy,
				Entry:     1.1000,
				Exit:      &exit,
				Lots:      1.0,
			},
			UpdateResult: &tradedom.Trade{
				ID:        tradeID,
				UserID:    userID,
				AccountID: &accountID,
				Type:      tradedom.TradeTypeBuy,
				Entry:     1.1000,
				Exit:      &exit,
				Lots:      1.0,
			},
		}
		// Set same P/L on both
		pl := 500.0
		tradeSpy.GetByIDResult.PL = &pl
		accountSpy := &AccountRepositorySpy{}
		service := NewService(tradeSpy, accountSpy)

		_, err := service.UpdateTrade(ctx, tradeID, userID, UpdateTradeRequest{
			AccountID: &accountID,
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04"),
			Pair:      "EUR/USD",
			Type:      "BUY",
			Entry:     1.1000,
			Exit:      &exit,
			Lots:      1.0,
		})

		// Assert no error
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Assert balance was NOT updated (difference is 0)
		if len(accountSpy.UpdateBalanceCalls) != 0 {
			t.Errorf("expected 0 calls to UpdateBalance when P/L unchanged, got %d", len(accountSpy.UpdateBalanceCalls))
		}
	})
}

func TestService_UpdateTrade_AccountChange(t *testing.T) {
	ctx := context.Background()
	oldAccountID := int64(1)
	newAccountID := int64(2)
	userID := int64(1)
	tradeID := int64(1)
	pl := 500.0 // 50 pips * 1 lot * $10
	exit := 1.1050

	t.Run("changing account reverts old balance and applies to new account", func(t *testing.T) {
		tradeSpy := &TradeRepositorySpy{
			GetByIDResult: &tradedom.Trade{
				ID:        tradeID,
				UserID:    userID,
				AccountID: &oldAccountID,
				Type:      tradedom.TradeTypeBuy,
				Entry:     1.1000,
				Exit:      &exit,
				PL:        &pl,
			},
			UpdateResult: &tradedom.Trade{
				ID:        tradeID,
				UserID:    userID,
				AccountID: &newAccountID,
				Type:      tradedom.TradeTypeBuy,
				Entry:     1.1000,
				Exit:      &exit,
				PL:        &pl,
			},
		}
		accountSpy := &AccountRepositorySpy{}
		service := NewService(tradeSpy, accountSpy)

		_, err := service.UpdateTrade(ctx, tradeID, userID, UpdateTradeRequest{
			AccountID: &newAccountID,
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04"),
			Pair:      "EUR/USD",
			Type:      "BUY",
			Entry:     1.1000,
			Exit:      &exit,
			Lots:      1.0,
		})

		// Assert no error
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Assert two balance updates: revert from old, apply to new
		if len(accountSpy.UpdateBalanceCalls) != 2 {
			t.Fatalf("expected 2 calls to UpdateBalance, got %d", len(accountSpy.UpdateBalanceCalls))
		}

		// First call should revert P/L from old account
		revertCall := accountSpy.UpdateBalanceCalls[0]
		if revertCall.ID != oldAccountID {
			t.Errorf("expected first call to old account %d, got %d", oldAccountID, revertCall.ID)
		}
		if revertCall.Amount != -500.0 {
			t.Errorf("expected first call to revert -500.0, got %.2f", revertCall.Amount)
		}

		// Second call should apply P/L to new account
		applyCall := accountSpy.UpdateBalanceCalls[1]
		if applyCall.ID != newAccountID {
			t.Errorf("expected second call to new account %d, got %d", newAccountID, applyCall.ID)
		}
		if applyCall.Amount != 500.0 {
			t.Errorf("expected second call to apply 500.0, got %.2f", applyCall.Amount)
		}
	})
}

func TestService_DeleteTrade_RevertsBalance(t *testing.T) {
	ctx := context.Background()
	accountID := int64(1)
	userID := int64(1)
	tradeID := int64(1)

	t.Run("deleting deposit reverts balance", func(t *testing.T) {
		amount := 1000.0
		tradeSpy := &TradeRepositorySpy{
			GetByIDResult: &tradedom.Trade{
				ID:        tradeID,
				UserID:    userID,
				AccountID: &accountID,
				Type:      tradedom.TradeTypeDeposit,
				Amount:    &amount,
			},
		}
		accountSpy := &AccountRepositorySpy{}
		service := NewService(tradeSpy, accountSpy)

		err := service.DeleteTrade(ctx, tradeID, userID)

		// Assert no error
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Assert balance was reverted (negative of deposit)
		if len(accountSpy.UpdateBalanceCalls) != 1 {
			t.Fatalf("expected 1 call to UpdateBalance, got %d", len(accountSpy.UpdateBalanceCalls))
		}

		balanceCall := accountSpy.UpdateBalanceCalls[0]
		if balanceCall.Amount != -1000.0 {
			t.Errorf("expected balance revert -1000.0, got %.2f", balanceCall.Amount)
		}

		// Assert trade was deleted
		if len(tradeSpy.DeleteCalls) != 1 {
			t.Errorf("expected 1 call to Delete, got %d", len(tradeSpy.DeleteCalls))
		}
	})

	t.Run("deleting withdraw reverts balance", func(t *testing.T) {
		amount := 500.0
		tradeSpy := &TradeRepositorySpy{
			GetByIDResult: &tradedom.Trade{
				ID:        tradeID,
				UserID:    userID,
				AccountID: &accountID,
				Type:      tradedom.TradeTypeWithdraw,
				Amount:    &amount,
			},
		}
		accountSpy := &AccountRepositorySpy{}
		service := NewService(tradeSpy, accountSpy)

		err := service.DeleteTrade(ctx, tradeID, userID)

		// Assert no error
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Assert balance was reverted (positive, adding back the withdrawn amount)
		if len(accountSpy.UpdateBalanceCalls) != 1 {
			t.Fatalf("expected 1 call to UpdateBalance, got %d", len(accountSpy.UpdateBalanceCalls))
		}

		balanceCall := accountSpy.UpdateBalanceCalls[0]
		if balanceCall.Amount != 500.0 {
			t.Errorf("expected balance revert 500.0 (add back withdrawn), got %.2f", balanceCall.Amount)
		}
	})

	t.Run("deleting closed trade reverts P/L", func(t *testing.T) {
		pl := 50.0
		tradeSpy := &TradeRepositorySpy{
			GetByIDResult: &tradedom.Trade{
				ID:        tradeID,
				UserID:    userID,
				AccountID: &accountID,
				Type:      tradedom.TradeTypeBuy,
				PL:        &pl,
			},
		}
		accountSpy := &AccountRepositorySpy{}
		service := NewService(tradeSpy, accountSpy)

		err := service.DeleteTrade(ctx, tradeID, userID)

		// Assert no error
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Assert balance was reverted (negative of P/L)
		if len(accountSpy.UpdateBalanceCalls) != 1 {
			t.Fatalf("expected 1 call to UpdateBalance, got %d", len(accountSpy.UpdateBalanceCalls))
		}

		balanceCall := accountSpy.UpdateBalanceCalls[0]
		if balanceCall.Amount != -50.0 {
			t.Errorf("expected balance revert -50.0, got %.2f", balanceCall.Amount)
		}
	})
}
