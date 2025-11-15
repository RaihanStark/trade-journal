package account

import "context"

// Repository defines the interface for account data access
type Repository interface {
	Create(ctx context.Context, account *Account) (*Account, error)
	GetByID(ctx context.Context, id int64, userID int64) (*Account, error)
	GetByUserID(ctx context.Context, userID int64) ([]*Account, error)
	Update(ctx context.Context, account *Account) (*Account, error)
	Delete(ctx context.Context, id int64, userID int64) error
}
