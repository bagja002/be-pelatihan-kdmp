package dto

// LoginRequest — login admin/super_admin dengan username.
type LoginRequest struct {
	Username string `json:"username" validate:"required,max=128"`
	Password string `json:"password" validate:"required"`
}

// RefreshRequest menukar refresh token dengan pasangan token baru.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TokenResponse dikembalikan saat login/refresh sukses.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// UserResponse — tampilan aman user (tanpa password).
type UserResponse struct {
	ID       uint   `json:"id"`
	Nama     string `json:"nama"`
	Username string `json:"username"`
	Type     string `json:"type"`
	IDSatdik *uint  `json:"idSatdik"`
}
