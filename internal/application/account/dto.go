package account

// CreateAccountRequest represents the data required to create a new account
type CreateAccountRequest struct {
	Name          string `json:"name" validate:"required"`
	Broker        string `json:"broker" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	AccountType   string `json:"account_type" validate:"required,oneof=demo live"`
	Currency      string `json:"currency" validate:"required"`
	IsActive      bool   `json:"is_active"`
}

// UpdateAccountRequest represents the data required to update an account
type UpdateAccountRequest struct {
	Name          string `json:"name" validate:"required"`
	Broker        string `json:"broker" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	AccountType   string `json:"account_type" validate:"required,oneof=demo live"`
	Currency      string `json:"currency" validate:"required"`
	IsActive      bool   `json:"is_active"`
}

// AccountDTO represents account data transfer object
type AccountDTO struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Broker         string  `json:"broker"`
	AccountNumber  string  `json:"account_number"`
	AccountType    string  `json:"account_type"`
	Currency       string  `json:"currency"`
	CurrentBalance float64 `json:"current_balance"`
	IsActive       bool    `json:"is_active"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}
