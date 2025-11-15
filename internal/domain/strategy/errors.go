package strategy

import "errors"

var (
	// ErrNotFound is returned when a strategy is not found or access is denied
	ErrNotFound = errors.New("strategy not found")
)
