package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/harusys/super-shiharai-kun/internal/usecase/auth"
)

// Handler handles authentication endpoints.
type Handler struct {
	usecase   auth.Usecase
	validator *validator.Validate
}

// NewHandler creates a new Handler.
func NewHandler(usecase auth.Usecase, validator *validator.Validate) *Handler {
	return &Handler{
		usecase:   usecase,
		validator: validator,
	}
}

// Register handles user registration.
// POST /api/v1/auth/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid request body"))

		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewValidationErrorResponse(formatValidationErrors(err)))

		return
	}

	input := &auth.RegisterInput{
		CompanyID: req.CompanyID,
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
	}

	tokenPair, err := h.usecase.Register(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, auth.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, NewErrorResponse("email already exists"))

			return
		}

		c.JSON(http.StatusInternalServerError, NewErrorResponse("internal server error"))

		return
	}

	c.JSON(http.StatusCreated, &TokenResponse{
		AccessToken:           tokenPair.AccessToken,
		AccessTokenExpiresAt:  tokenPair.AccessTokenExpiresAt,
		RefreshToken:          tokenPair.RefreshToken,
		RefreshTokenExpiresAt: tokenPair.RefreshTokenExpiresAt,
	})
}

// Login handles user login.
// POST /api/v1/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid request body"))

		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewValidationErrorResponse(formatValidationErrors(err)))

		return
	}

	input := &auth.Input{
		Email:    req.Email,
		Password: req.Password,
	}

	tokenPair, err := h.usecase.Login(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, NewErrorResponse("invalid credentials"))

			return
		}

		c.JSON(http.StatusInternalServerError, NewErrorResponse("internal server error"))

		return
	}

	c.JSON(http.StatusOK, &TokenResponse{
		AccessToken:           tokenPair.AccessToken,
		AccessTokenExpiresAt:  tokenPair.AccessTokenExpiresAt,
		RefreshToken:          tokenPair.RefreshToken,
		RefreshTokenExpiresAt: tokenPair.RefreshTokenExpiresAt,
	})
}

// RefreshToken handles token refresh.
// POST /api/v1/auth/refresh
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid request body"))

		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewValidationErrorResponse(formatValidationErrors(err)))

		return
	}

	tokenPair, err := h.usecase.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, NewErrorResponse("invalid or expired refresh token"))

		return
	}

	c.JSON(http.StatusOK, &TokenResponse{
		AccessToken:           tokenPair.AccessToken,
		AccessTokenExpiresAt:  tokenPair.AccessTokenExpiresAt,
		RefreshToken:          tokenPair.RefreshToken,
		RefreshTokenExpiresAt: tokenPair.RefreshTokenExpiresAt,
	})
}

func formatValidationErrors(err error) map[string]string {
	details := make(map[string]string)

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, e := range validationErrs {
			details[e.Field()] = e.Tag()
		}
	}

	return details
}
