package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/harusys/super-shiharai-kun/internal/infrastructure"
)

// ErrInvalidToken is returned when the token is invalid.
var ErrInvalidToken = errors.New("invalid token")

// ErrExpiredToken is returned when the token is expired.
var ErrExpiredToken = errors.New("expired token")

// Claims represents JWT claims.
type Claims struct {
	UserID    int64 `json:"user_id"`
	CompanyID int64 `json:"company_id"`
	jwt.RegisteredClaims
}

// JWTService provides JWT token generation and validation.
type JWTService struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewJWTService creates a new JWTService.
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey:     []byte(secretKey),
		accessExpiry:  infrastructure.AccessTokenExpiry,
		refreshExpiry: infrastructure.RefreshTokenExpiry,
	}
}

// GenerateAccessToken generates an access token.
func (s *JWTService) GenerateAccessToken(userID, companyID int64) (string, error) {
	return s.generateToken(userID, companyID, s.accessExpiry)
}

// GenerateRefreshToken generates a refresh token.
func (s *JWTService) GenerateRefreshToken(userID, companyID int64) (string, error) {
	return s.generateToken(userID, companyID, s.refreshExpiry)
}

func (s *JWTService) generateToken(userID, companyID int64, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:    userID,
		CompanyID: companyID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secretKey)
}

// ValidateToken validates a token and returns claims.
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}

			return s.secretKey, nil
		},
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}

		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
