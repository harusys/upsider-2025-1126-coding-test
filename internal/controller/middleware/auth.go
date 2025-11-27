package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/harusys/super-shiharai-kun/internal/infrastructure/security"
)

const (
	// UserIDKey is the context key for user ID.
	UserIDKey = "user_id"
	// CompanyIDKey is the context key for company ID.
	CompanyIDKey = "company_id"
)

// ErrorResponse is the standard error response body.
type ErrorResponse struct {
	Error string `json:"error"`
}

// AuthMiddleware creates a JWT authentication middleware.
func AuthMiddleware(jwtService *security.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				ErrorResponse{Error: "missing authorization header"},
			)

			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				ErrorResponse{Error: "invalid authorization header format"},
			)

			return
		}

		token := parts[1]

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				ErrorResponse{Error: "invalid or expired token"},
			)

			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(CompanyIDKey, claims.CompanyID)
		c.Next()
	}
}

// GetUserID retrieves the user ID from the gin context.
func GetUserID(c *gin.Context) int64 {
	userID, _ := c.Get(UserIDKey)
	if id, ok := userID.(int64); ok {
		return id
	}

	return 0
}

// GetCompanyID retrieves the company ID from the gin context.
func GetCompanyID(c *gin.Context) int64 {
	companyID, _ := c.Get(CompanyIDKey)
	if id, ok := companyID.(int64); ok {
		return id
	}

	return 0
}
