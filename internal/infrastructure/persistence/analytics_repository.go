package persistence

import (
	"context"

	"github.com/raihanstark/trade-journal/internal/db"
)

type AnalyticsRepository struct {
	queries *db.Queries
}

func NewAnalyticsRepository(queries *db.Queries) *AnalyticsRepository {
	return &AnalyticsRepository{queries: queries}
}

// GetUserTrades returns all trades for a user (raw data only)
func (r *AnalyticsRepository) GetUserTrades(ctx context.Context, userID int64) ([]db.Trade, error) {
	trades, err := r.queries.GetTradesByUserID(ctx, int32(userID))
	if err != nil {
		return nil, err
	}
	return trades, nil
}
