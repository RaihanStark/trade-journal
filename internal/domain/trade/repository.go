package trade

import "context"

type Repository interface {
	Create(ctx context.Context, trade *Trade) (*Trade, error)
	GetByID(ctx context.Context, id int64, userID int64) (*Trade, error)
	GetByUserID(ctx context.Context, userID int64) ([]*Trade, error)
	Update(ctx context.Context, trade *Trade) (*Trade, error)
	Delete(ctx context.Context, id int64, userID int64) error
	GetByAccountID(ctx context.Context, accountID int64, userID int64) ([]*Trade, error)
}
