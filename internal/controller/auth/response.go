package auth

import "time"

// TokenResponse is the response body containing authentication tokens.
type TokenResponse struct {
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

// ErrorResponse is the standard error response body.
type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

// NewErrorResponse creates a new ErrorResponse.
func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Error: message,
	}
}

// NewValidationErrorResponse creates a new ErrorResponse for validation errors.
func NewValidationErrorResponse(details map[string]string) *ErrorResponse {
	return &ErrorResponse{
		Error:   "validation error",
		Details: details,
	}
}
