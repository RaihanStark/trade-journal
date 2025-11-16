package analytics

import (
	"context"

	"github.com/raihanstark/trade-journal/internal/db"
)

type Repository interface {
	// GetUserTrades returns raw trade data for a specific user
	GetUserTrades(ctx context.Context, userID int64) ([]db.Trade, error)
}
