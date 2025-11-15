package auth

// RegisterRequest represents the data required to register a new user
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest represents the data required to login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string   `json:"token"`
	User  UserDTO  `json:"user"`
}

// UserDTO represents user data transfer object
type UserDTO struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
