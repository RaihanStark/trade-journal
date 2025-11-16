package analytics

import (
	"context"

	"github.com/raihanstark/trade-journal/internal/domain/analytics"
)

type Service struct {
	repo       analytics.Repository
	calculator *Calculator
}

func NewService(repo analytics.Repository) *Service {
	return &Service{
		repo:       repo,
		calculator: NewCalculator(),
	}
}

func (s *Service) GetUserAnalytics(ctx context.Context, userID int64) (*AnalyticsDTO, error) {
	// Get raw trade data from repository
	trades, err := s.repo.GetUserTrades(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Calculate analytics using calculator
	analyticsData := s.calculator.CalculateAnalytics(trades)

	// Convert to DTO
	return s.toDTO(analyticsData), nil
}

func (s *Service) toDTO(a *analytics.Analytics) *AnalyticsDTO {
	return &AnalyticsDTO{
		TotalPL:           a.TotalPL,
		WinRate:           a.WinRate,
		TotalTrades:       a.TotalTrades,
		WinningTrades:     a.WinningTrades,
		LosingTrades:      a.LosingTrades,
		AvgWin:            a.AvgWin,
		AvgLoss:           a.AvgLoss,
		ProfitFactor:      a.ProfitFactor,
		SharpeRatio:       a.SharpeRatio,
		MaxDrawdown:       a.MaxDrawdown,
		LargestWin:        a.LargestWin,
		LargestLoss:       a.LargestLoss,
		AvgRR:             a.AvgRR,
		ConsecutiveWins:   a.ConsecutiveWins,
		ConsecutiveLosses: a.ConsecutiveLosses,
		BestStreak:        a.BestStreak,
		WorstStreak:       a.WorstStreak,
	}
}
