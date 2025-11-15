package persistence

import (
	"context"

	"github.com/raihanstark/trade-journal/internal/db"
	"github.com/raihanstark/trade-journal/internal/domain/strategy"
)

// StrategyRepository implements strategy.Repository using sqlc
type StrategyRepository struct {
	queries *db.Queries
}

// NewStrategyRepository creates a new strategy repository
func NewStrategyRepository(queries *db.Queries) *StrategyRepository {
	return &StrategyRepository{
		queries: queries,
	}
}

// Create creates a new strategy in the database
func (r *StrategyRepository) Create(ctx context.Context, s *strategy.Strategy) (*strategy.Strategy, error) {
	result, err := r.queries.CreateStrategy(ctx, db.CreateStrategyParams{
		UserID:      int32(s.UserID),
		Name:        s.Name,
		Description: db.StringToNullString(s.Description),
	})
	if err != nil {
		return nil, err
	}

	return &strategy.Strategy{
		ID:          int64(result.ID),
		UserID:      int64(result.UserID),
		Name:        result.Name,
		Description: db.NullStringToString(result.Description),
		CreatedAt:   result.CreatedAt.Time,
		UpdatedAt:   result.UpdatedAt.Time,
	}, nil
}

// GetByID retrieves a strategy by ID
func (r *StrategyRepository) GetByID(ctx context.Context, id int64, userID int64) (*strategy.Strategy, error) {
	result, err := r.queries.GetStrategyByID(ctx, db.GetStrategyByIDParams{
		ID:     int32(id),
		UserID: int32(userID),
	})
	if err != nil {
		return nil, err
	}

	return &strategy.Strategy{
		ID:          int64(result.ID),
		UserID:      int64(result.UserID),
		Name:        result.Name,
		Description: db.NullStringToString(result.Description),
		CreatedAt:   result.CreatedAt.Time,
		UpdatedAt:   result.UpdatedAt.Time,
	}, nil
}

// GetByUserID retrieves all strategies for a user
func (r *StrategyRepository) GetByUserID(ctx context.Context, userID int64) ([]*strategy.Strategy, error) {
	results, err := r.queries.GetStrategiesByUserID(ctx, int32(userID))
	if err != nil {
		return nil, err
	}

	strategies := make([]*strategy.Strategy, len(results))
	for i, result := range results {
		strategies[i] = &strategy.Strategy{
			ID:          int64(result.ID),
			UserID:      int64(result.UserID),
			Name:        result.Name,
			Description: db.NullStringToString(result.Description),
			CreatedAt:   result.CreatedAt.Time,
			UpdatedAt:   result.UpdatedAt.Time,
		}
	}

	return strategies, nil
}

// Update updates an existing strategy
func (r *StrategyRepository) Update(ctx context.Context, s *strategy.Strategy) (*strategy.Strategy, error) {
	result, err := r.queries.UpdateStrategy(ctx, db.UpdateStrategyParams{
		ID:          int32(s.ID),
		Name:        s.Name,
		Description: db.StringToNullString(s.Description),
		UserID:      int32(s.UserID),
	})
	if err != nil {
		return nil, err
	}

	return &strategy.Strategy{
		ID:          int64(result.ID),
		UserID:      int64(result.UserID),
		Name:        result.Name,
		Description: db.NullStringToString(result.Description),
		CreatedAt:   result.CreatedAt.Time,
		UpdatedAt:   result.UpdatedAt.Time,
	}, nil
}

// Delete deletes a strategy
func (r *StrategyRepository) Delete(ctx context.Context, id int64, userID int64) error {
	return r.queries.DeleteStrategy(ctx, db.DeleteStrategyParams{
		ID:     int32(id),
		UserID: int32(userID),
	})
}
