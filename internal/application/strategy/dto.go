package strategy

import "time"

// CreateStrategyRequest represents a request to create a new strategy
type CreateStrategyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateStrategyRequest represents a request to update an existing strategy
type UpdateStrategyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// StrategyDTO represents a strategy data transfer object
type StrategyDTO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
