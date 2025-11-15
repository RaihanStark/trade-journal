package account

import (
	"context"
	"errors"

	"github.com/raihanstark/trade-journal/internal/domain/account"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrUnauthorized    = errors.New("unauthorized to access this account")
)

// Service handles account use cases
type Service struct {
	accountRepo account.Repository
}

// NewService creates a new account service
func NewService(accountRepo account.Repository) *Service {
	return &Service{
		accountRepo: accountRepo,
	}
}

// CreateAccount creates a new trading account
func (s *Service) CreateAccount(ctx context.Context, userID int64, req CreateAccountRequest) (*AccountDTO, error) {
	// Create domain entity
	acc := account.NewAccount(
		userID,
		req.Name,
		req.Broker,
		req.AccountNumber,
		account.AccountType(req.AccountType),
		req.Currency,
	)
	acc.IsActive = req.IsActive

	// Save to repository
	createdAccount, err := s.accountRepo.Create(ctx, acc)
	if err != nil {
		return nil, err
	}

	return toDTO(createdAccount), nil
}

// GetAccount retrieves an account by ID
func (s *Service) GetAccount(ctx context.Context, id int64, userID int64) (*AccountDTO, error) {
	acc, err := s.accountRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	return toDTO(acc), nil
}

// GetUserAccounts retrieves all accounts for a user
func (s *Service) GetUserAccounts(ctx context.Context, userID int64) ([]*AccountDTO, error) {
	accounts, err := s.accountRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*AccountDTO, len(accounts))
	for i, acc := range accounts {
		dtos[i] = toDTO(acc)
	}

	return dtos, nil
}

// UpdateAccount updates an existing account
func (s *Service) UpdateAccount(ctx context.Context, id int64, userID int64, req UpdateAccountRequest) (*AccountDTO, error) {
	// Get existing account
	existingAccount, err := s.accountRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	// Update fields
	existingAccount.Name = req.Name
	existingAccount.Broker = req.Broker
	existingAccount.AccountNumber = req.AccountNumber
	existingAccount.AccountType = account.AccountType(req.AccountType)
	existingAccount.Currency = req.Currency
	existingAccount.IsActive = req.IsActive

	// Save to repository
	updatedAccount, err := s.accountRepo.Update(ctx, existingAccount)
	if err != nil {
		return nil, err
	}

	return toDTO(updatedAccount), nil
}

// DeleteAccount deletes an account
func (s *Service) DeleteAccount(ctx context.Context, id int64, userID int64) error {
	return s.accountRepo.Delete(ctx, id, userID)
}

// toDTO converts domain entity to DTO
func toDTO(acc *account.Account) *AccountDTO {
	return &AccountDTO{
		ID:             acc.ID,
		Name:           acc.Name,
		Broker:         acc.Broker,
		AccountNumber:  acc.AccountNumber,
		AccountType:    string(acc.AccountType),
		Currency:       acc.Currency,
		CurrentBalance: acc.CurrentBalance,
		IsActive:       acc.IsActive,
		CreatedAt:      acc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      acc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
