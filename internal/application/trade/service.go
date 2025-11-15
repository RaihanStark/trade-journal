package trade

import (
	"context"
	"time"

	"github.com/raihanstark/trade-journal/internal/domain/account"
	"github.com/raihanstark/trade-journal/internal/domain/trade"
)

type Service struct {
	repo        trade.Repository
	accountRepo account.Repository
}

func NewService(repo trade.Repository, accountRepo account.Repository) *Service {
	return &Service{
		repo:        repo,
		accountRepo: accountRepo,
	}
}

func (s *Service) CreateTrade(ctx context.Context, userID int64, req CreateTradeRequest) (*TradeDTO, error) {
	// Parse date and time
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, err
	}

	tradeTime, err := time.Parse("15:04", req.Time)
	if err != nil {
		return nil, err
	}

	t := &trade.Trade{
		UserID:     userID,
		AccountID:  req.AccountID,
		Date:       date,
		Time:       tradeTime,
		Pair:       req.Pair,
		Type:       trade.TradeType(req.Type),
		Entry:      req.Entry,
		Exit:       req.Exit,
		Lots:       req.Lots,
		StopLoss:   req.StopLoss,
		TakeProfit: req.TakeProfit,
		Notes:      req.Notes,
		Mistakes:   req.Mistakes,
		Amount:     req.Amount,
	}

	// Calculate metrics (pips, P/L, R:R, status)
	CalculateTradeMetrics(t)

	// Convert strategy IDs to Strategy objects
	var strategies []trade.Strategy
	for _, id := range req.StrategyIDs {
		strategies = append(strategies, trade.Strategy{ID: id})
	}
	t.Strategies = strategies

	created, err := s.repo.Create(ctx, t)
	if err != nil {
		return nil, err
	}

	// Update account balance for DEPOSIT and WITHDRAW trades
	if (t.Type == trade.TradeTypeDeposit || t.Type == trade.TradeTypeWithdraw) && t.AccountID != nil && t.Amount != nil {
		amount := *t.Amount
		if t.Type == trade.TradeTypeWithdraw {
			amount = -amount // Withdraw reduces balance
		}
		_, err = s.accountRepo.UpdateBalance(ctx, *t.AccountID, userID, amount)
		if err != nil {
			// Log error but don't fail the trade creation
			// You might want to handle this differently in production
		}
	}

	return s.toDTO(created), nil
}

func (s *Service) GetUserTrades(ctx context.Context, userID int64) ([]*TradeDTO, error) {
	trades, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*TradeDTO, len(trades))
	for i, t := range trades {
		dtos[i] = s.toDTO(t)
	}

	return dtos, nil
}

func (s *Service) GetTrade(ctx context.Context, id int64, userID int64) (*TradeDTO, error) {
	t, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return s.toDTO(t), nil
}

func (s *Service) UpdateTrade(ctx context.Context, id int64, userID int64, req UpdateTradeRequest) (*TradeDTO, error) {
	// Parse date and time
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, err
	}

	tradeTime, err := time.Parse("15:04", req.Time)
	if err != nil {
		return nil, err
	}

	t := &trade.Trade{
		ID:         id,
		UserID:     userID,
		AccountID:  req.AccountID,
		Date:       date,
		Time:       tradeTime,
		Pair:       req.Pair,
		Type:       trade.TradeType(req.Type),
		Entry:      req.Entry,
		Exit:       req.Exit,
		Lots:       req.Lots,
		StopLoss:   req.StopLoss,
		TakeProfit: req.TakeProfit,
		Notes:      req.Notes,
		Mistakes:   req.Mistakes,
		Amount:     req.Amount,
	}

	// Calculate metrics (pips, P/L, R:R, status)
	CalculateTradeMetrics(t)

	// Convert strategy IDs to Strategy objects
	var strategies []trade.Strategy
	for _, id := range req.StrategyIDs {
		strategies = append(strategies, trade.Strategy{ID: id})
	}
	t.Strategies = strategies

	updated, err := s.repo.Update(ctx, t)
	if err != nil {
		return nil, err
	}

	return s.toDTO(updated), nil
}

func (s *Service) DeleteTrade(ctx context.Context, id int64, userID int64) error {
	// Get the trade first to check if it's a deposit/withdrawal
	t, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	// Revert account balance for DEPOSIT and WITHDRAW trades before deleting
	if (t.Type == trade.TradeTypeDeposit || t.Type == trade.TradeTypeWithdraw) && t.AccountID != nil && t.Amount != nil {
		amount := *t.Amount
		// Reverse the transaction
		if t.Type == trade.TradeTypeDeposit {
			amount = -amount // Revert deposit by subtracting
		}
		// For withdraw, amount stays positive to add it back
		_, err = s.accountRepo.UpdateBalance(ctx, *t.AccountID, userID, amount)
		if err != nil {
			// Log error but continue with deletion
			// You might want to handle this differently in production
		}
	}

	return s.repo.Delete(ctx, id, userID)
}

func (s *Service) toDTO(t *trade.Trade) *TradeDTO {
	strategies := make([]Strategy, len(t.Strategies))
	for i, s := range t.Strategies {
		strategies[i] = Strategy{
			ID:   s.ID,
			Name: s.Name,
		}
	}

	return &TradeDTO{
		ID:         t.ID,
		AccountID:  t.AccountID,
		Date:       t.Date.Format("2006-01-02"),
		Time:       t.Time.Format("15:04"),
		Pair:       t.Pair,
		Type:       string(t.Type),
		Entry:      t.Entry,
		Exit:       t.Exit,
		Lots:       t.Lots,
		Pips:       t.Pips,
		PL:         t.PL,
		RR:         t.RR,
		Status:     string(t.Status),
		StopLoss:   t.StopLoss,
		TakeProfit: t.TakeProfit,
		Notes:      t.Notes,
		Mistakes:   t.Mistakes,
		Amount:     t.Amount,
		Strategies: strategies,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}
