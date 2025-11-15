package persistence

import (
	"context"

	"github.com/raihanstark/trade-journal/internal/db"
	"github.com/raihanstark/trade-journal/internal/domain/account"
)

// AccountRepository implements account.Repository using sqlc
type AccountRepository struct {
	queries *db.Queries
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(queries *db.Queries) *AccountRepository {
	return &AccountRepository{
		queries: queries,
	}
}

// Create creates a new account in the database
func (r *AccountRepository) Create(ctx context.Context, acc *account.Account) (*account.Account, error) {
	result, err := r.queries.CreateAccount(ctx, db.CreateAccountParams{
		UserID:        int32(acc.UserID),
		Name:          acc.Name,
		Broker:        acc.Broker,
		AccountNumber: acc.AccountNumber,
		AccountType:   string(acc.AccountType),
		Currency:      acc.Currency,
		IsActive:      acc.IsActive,
	})
	if err != nil {
		return nil, err
	}

	return &account.Account{
		ID:            int64(result.ID),
		UserID:        int64(result.UserID),
		Name:          result.Name,
		Broker:        result.Broker,
		AccountNumber: result.AccountNumber,
		AccountType:   account.AccountType(result.AccountType),
		Currency:      result.Currency,
		IsActive:      result.IsActive,
		CreatedAt:     result.CreatedAt.Time,
		UpdatedAt:     result.UpdatedAt.Time,
	}, nil
}

// GetByID retrieves an account by ID
func (r *AccountRepository) GetByID(ctx context.Context, id int64, userID int64) (*account.Account, error) {
	result, err := r.queries.GetAccountByID(ctx, db.GetAccountByIDParams{
		ID:     int32(id),
		UserID: int32(userID),
	})
	if err != nil {
		return nil, err
	}

	return &account.Account{
		ID:            int64(result.ID),
		UserID:        int64(result.UserID),
		Name:          result.Name,
		Broker:        result.Broker,
		AccountNumber: result.AccountNumber,
		AccountType:   account.AccountType(result.AccountType),
		Currency:      result.Currency,
		IsActive:      result.IsActive,
		CreatedAt:     result.CreatedAt.Time,
		UpdatedAt:     result.UpdatedAt.Time,
	}, nil
}

// GetByUserID retrieves all accounts for a user
func (r *AccountRepository) GetByUserID(ctx context.Context, userID int64) ([]*account.Account, error) {
	results, err := r.queries.GetAccountsByUserID(ctx, int32(userID))
	if err != nil {
		return nil, err
	}

	accounts := make([]*account.Account, len(results))
	for i, result := range results {
		accounts[i] = &account.Account{
			ID:            int64(result.ID),
			UserID:        int64(result.UserID),
			Name:          result.Name,
			Broker:        result.Broker,
			AccountNumber: result.AccountNumber,
			AccountType:   account.AccountType(result.AccountType),
			Currency:      result.Currency,
			IsActive:      result.IsActive,
			CreatedAt:     result.CreatedAt.Time,
			UpdatedAt:     result.UpdatedAt.Time,
		}
	}

	return accounts, nil
}

// Update updates an existing account
func (r *AccountRepository) Update(ctx context.Context, acc *account.Account) (*account.Account, error) {
	result, err := r.queries.UpdateAccount(ctx, db.UpdateAccountParams{
		ID:            int32(acc.ID),
		UserID:        int32(acc.UserID),
		Name:          acc.Name,
		Broker:        acc.Broker,
		AccountNumber: acc.AccountNumber,
		AccountType:   string(acc.AccountType),
		Currency:      acc.Currency,
		IsActive:      acc.IsActive,
	})
	if err != nil {
		return nil, err
	}

	return &account.Account{
		ID:            int64(result.ID),
		UserID:        int64(result.UserID),
		Name:          result.Name,
		Broker:        result.Broker,
		AccountNumber: result.AccountNumber,
		AccountType:   account.AccountType(result.AccountType),
		Currency:      result.Currency,
		IsActive:      result.IsActive,
		CreatedAt:     result.CreatedAt.Time,
		UpdatedAt:     result.UpdatedAt.Time,
	}, nil
}

// Delete deletes an account
func (r *AccountRepository) Delete(ctx context.Context, id int64, userID int64) error {
	return r.queries.DeleteAccount(ctx, db.DeleteAccountParams{
		ID:     int32(id),
		UserID: int32(userID),
	})
}
