package auth

import (
	"errors"
	"net/http"
	"time"

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
//
//	@Summary		ユーザー登録
//	@Description	新規ユーザーを登録し、JWTトークンを発行します
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterRequest	true	"登録リクエスト"
//	@Success		201		{object}	TokenResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		409		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/register [post]
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
		AccessToken:  tokenPair.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Until(tokenPair.AccessTokenExpiresAt).Seconds()),
		RefreshToken: tokenPair.RefreshToken,
	})
}

// Login handles user login.
//
//	@Summary		ログイン
//	@Description	メールアドレスとパスワードで認証し、JWTトークンを発行します
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest	true	"ログインリクエスト"
//	@Success		200		{object}	TokenResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/login [post]
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
		AccessToken:  tokenPair.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Until(tokenPair.AccessTokenExpiresAt).Seconds()),
		RefreshToken: tokenPair.RefreshToken,
	})
}

// RefreshToken handles token refresh.
//
//	@Summary		トークン更新
//	@Description	リフレッシュトークンを使用して新しいアクセストークンを取得します
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RefreshTokenRequest	true	"トークン更新リクエスト"
//	@Success		200		{object}	TokenResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Router			/auth/refresh [post]
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
		AccessToken:  tokenPair.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Until(tokenPair.AccessTokenExpiresAt).Seconds()),
		RefreshToken: tokenPair.RefreshToken,
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
