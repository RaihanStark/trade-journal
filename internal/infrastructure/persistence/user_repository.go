package persistence

import (
	"context"

	"github.com/raihanstark/trade-journal/internal/db"
	"github.com/raihanstark/trade-journal/internal/domain/user"
)

// UserRepository implements user.Repository using sqlc
type UserRepository struct {
	queries *db.Queries
}

// NewUserRepository creates a new user repository
func NewUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{
		queries: queries,
	}
}

// Create creates a new user in the database
func (r *UserRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	result, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	})
	if err != nil {
		return nil, err
	}

	return &user.User{
		ID:           int64(result.ID),
		Email:        result.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    result.CreatedAt.Time,
		UpdatedAt:    result.UpdatedAt.Time,
	}, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	result, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &user.User{
		ID:           int64(result.ID),
		Email:        result.Email,
		PasswordHash: result.PasswordHash,
		CreatedAt:    result.CreatedAt.Time,
		UpdatedAt:    result.UpdatedAt.Time,
	}, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*user.User, error) {
	result, err := r.queries.GetUserByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return &user.User{
		ID:        int64(result.ID),
		Email:     result.Email,
		CreatedAt: result.CreatedAt.Time,
		UpdatedAt: result.UpdatedAt.Time,
	}, nil
}
