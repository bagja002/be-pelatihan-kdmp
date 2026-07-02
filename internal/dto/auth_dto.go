package dto

// RegisterRequest is the payload for account registration.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=128"`
	Phone    string `json:"phone" validate:"omitempty,max=32"`
}

// LoginRequest is the payload for authentication.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshRequest exchanges a refresh token for a new token pair.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TokenResponse is returned on successful login/refresh.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // access token lifetime in seconds
}

// UserResponse is the safe, serializable view of a user (no password).
type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Phone string `json:"phone,omitempty"`
}
