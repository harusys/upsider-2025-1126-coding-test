package auth

// TokenResponse is the response body containing authentication tokens.
// Follows OAuth 2.0 (RFC 6749) standard format.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
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
