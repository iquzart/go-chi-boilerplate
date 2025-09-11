package dto

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required" example:"user@example.com" format:"email" description:"User's registered email address"`
	Password string `json:"password" binding:"required" example:"P@ssw0rd" description:"User's account password"`
}

// LoginResponse represents the login response payload
type LoginResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6..." description:"JWT access token used for authenticating API requests"`
	RefreshToken string `json:"refresh_token" example:"dGhpc19pc19hX3JlZnJlc2hfdG9rZW4=" description:"Token used to obtain a new access token when the current one expires"`
}
