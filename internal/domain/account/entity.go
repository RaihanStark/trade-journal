package account

import "time"

// AccountType represents the type of trading account
type AccountType string

const (
	AccountTypeDemo AccountType = "demo"
	AccountTypeLive AccountType = "live"
)

// Account represents a trading account in the system
type Account struct {
	ID             int64
	UserID         int64
	Name           string
	Broker         string
	AccountNumber  string
	AccountType    AccountType
	Currency       string
	CurrentBalance float64
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewAccount creates a new account instance
func NewAccount(userID int64, name, broker, accountNumber string, accountType AccountType, currency string) *Account {
	now := time.Now()
	return &Account{
		UserID:        userID,
		Name:          name,
		Broker:        broker,
		AccountNumber: accountNumber,
		AccountType:   accountType,
		Currency:      currency,
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
