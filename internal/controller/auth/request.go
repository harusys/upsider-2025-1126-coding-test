package auth

// RegisterRequest is the request body for user registration.
type RegisterRequest struct {
	CompanyID int64  `json:"company_id" validate:"required,gt=0"`
	Name      string `json:"name"       validate:"required,min=1,max=255"`
	Email     string `json:"email"      validate:"required,email,max=255"`
	Password  string `json:"password"   validate:"required,min=8,max=72"`
}

// LoginRequest is the request body for user login.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest is the request body for token refresh.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
