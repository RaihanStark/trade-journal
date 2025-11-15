package strategy

import (
	"context"
	"errors"

	"github.com/raihanstark/trade-journal/internal/domain/strategy"
)

var (
	ErrStrategyNotFound = errors.New("strategy not found")
)

// Service handles strategy business logic
type Service struct {
	repo strategy.Repository
}

// NewService creates a new strategy service
func NewService(repo strategy.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateStrategy creates a new strategy
func (s *Service) CreateStrategy(ctx context.Context, userID int64, req CreateStrategyRequest) (*StrategyDTO, error) {
	strategyEntity := &strategy.Strategy{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}

	created, err := s.repo.Create(ctx, strategyEntity)
	if err != nil {
		return nil, err
	}

	return &StrategyDTO{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}, nil
}

// GetStrategy retrieves a strategy by ID
func (s *Service) GetStrategy(ctx context.Context, id int64, userID int64) (*StrategyDTO, error) {
	strategyEntity, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, ErrStrategyNotFound
	}

	return &StrategyDTO{
		ID:          strategyEntity.ID,
		Name:        strategyEntity.Name,
		Description: strategyEntity.Description,
		CreatedAt:   strategyEntity.CreatedAt,
		UpdatedAt:   strategyEntity.UpdatedAt,
	}, nil
}

// GetUserStrategies retrieves all strategies for a user
func (s *Service) GetUserStrategies(ctx context.Context, userID int64) ([]*StrategyDTO, error) {
	strategies, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*StrategyDTO, len(strategies))
	for i, strategyEntity := range strategies {
		dtos[i] = &StrategyDTO{
			ID:          strategyEntity.ID,
			Name:        strategyEntity.Name,
			Description: strategyEntity.Description,
			CreatedAt:   strategyEntity.CreatedAt,
			UpdatedAt:   strategyEntity.UpdatedAt,
		}
	}

	return dtos, nil
}

// UpdateStrategy updates an existing strategy
func (s *Service) UpdateStrategy(ctx context.Context, id int64, userID int64, req UpdateStrategyRequest) (*StrategyDTO, error) {
	strategyEntity := &strategy.Strategy{
		ID:          id,
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}

	updated, err := s.repo.Update(ctx, strategyEntity)
	if err != nil {
		return nil, ErrStrategyNotFound
	}

	return &StrategyDTO{
		ID:          updated.ID,
		Name:        updated.Name,
		Description: updated.Description,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}, nil
}

// DeleteStrategy deletes a strategy
func (s *Service) DeleteStrategy(ctx context.Context, id int64, userID int64) error {
	return s.repo.Delete(ctx, id, userID)
}
