package strategy

import "time"

// Strategy represents a trading strategy entity
type Strategy struct {
	ID          int64
	UserID      int64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
