//go:generate mockgen -source=$GOFILE -destination=mock/mock_usecase.go -package=mock

package auth

import (
	"context"
	"time"
)

// Input is the input for authentication operations.
type Input struct {
	Email    string
	Password string
}

// RegisterInput is the input for user registration.
type RegisterInput struct {
	CompanyID int64
	Name      string
	Email     string
	Password  string
}

// TokenPair holds access and refresh tokens with their expiration times.
type TokenPair struct {
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

// Usecase defines authentication operations.
type Usecase interface {
	// Register creates a new user account.
	Register(ctx context.Context, input *RegisterInput) (*TokenPair, error)
	// Login authenticates a user and returns tokens.
	Login(ctx context.Context, input *Input) (*TokenPair, error)
	// RefreshToken generates new tokens using a refresh token.
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
}
