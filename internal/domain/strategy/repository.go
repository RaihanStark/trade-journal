package strategy

import "context"

// Repository defines the interface for strategy data operations
type Repository interface {
	Create(ctx context.Context, strategy *Strategy) (*Strategy, error)
	GetByID(ctx context.Context, id int64, userID int64) (*Strategy, error)
	GetByUserID(ctx context.Context, userID int64) ([]*Strategy, error)
	Update(ctx context.Context, strategy *Strategy) (*Strategy, error)
	Delete(ctx context.Context, id int64, userID int64) error
}
