package trade

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, trade *Trade) (*Trade, error)
	GetByID(ctx context.Context, id int64, userID int64) (*Trade, error)
	GetByUserID(ctx context.Context, userID int64) ([]*Trade, error)
	GetByUserIDAndDateRange(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*Trade, error)
	Update(ctx context.Context, trade *Trade) (*Trade, error)
	Delete(ctx context.Context, id int64, userID int64) error
	GetByAccountID(ctx context.Context, accountID int64, userID int64) ([]*Trade, error)
	GetByAccountIDAndDateRange(ctx context.Context, accountID int64, userID int64, startDate, endDate time.Time) ([]*Trade, error)
}
