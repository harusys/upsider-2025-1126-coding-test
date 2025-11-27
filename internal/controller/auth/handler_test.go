package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/harusys/super-shiharai-kun/internal/controller/auth"
	usecase "github.com/harusys/super-shiharai-kun/internal/usecase/auth"
	"github.com/harusys/super-shiharai-kun/internal/usecase/auth/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler_Register(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       map[string]any
		prepare    func(m *mock.MockUsecase)
		wantStatus int
		wantBody   map[string]any
	}{
		{
			name: "success",
			body: map[string]any{
				"company_id": 1,
				"name":       "Test User",
				"email":      "test@example.com",
				"password":   "password123",
			},
			prepare: func(m *mock.MockUsecase) {
				m.EXPECT().
					Register(gomock.Any(), &usecase.RegisterInput{
						CompanyID: 1,
						Name:      "Test User",
						Email:     "test@example.com",
						Password:  "password123",
					}).
					Return(&usecase.TokenPair{
						AccessToken:           "access-token",
						AccessTokenExpiresAt:  time.Date(2024, 1, 1, 10, 15, 0, 0, time.UTC),
						RefreshToken:          "refresh-token",
						RefreshTokenExpiresAt: time.Date(2024, 1, 8, 10, 0, 0, 0, time.UTC),
					}, nil)
			},
			wantStatus: http.StatusCreated,
			wantBody: map[string]any{
				"access_token":  "access-token",
				"token_type":    "Bearer",
				"refresh_token": "refresh-token",
			},
		},
		{
			name: "email already exists",
			body: map[string]any{
				"company_id": 1,
				"name":       "Test User",
				"email":      "existing@example.com",
				"password":   "password123",
			},
			prepare: func(m *mock.MockUsecase) {
				m.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Return(nil, usecase.ErrEmailAlreadyExists)
			},
			wantStatus: http.StatusConflict,
			wantBody: map[string]any{
				"error": "email already exists",
			},
		},
		{
			name: "invalid request - missing email",
			body: map[string]any{
				"company_id": 1,
				"name":       "Test User",
				"password":   "password123",
			},
			prepare:    func(_ *mock.MockUsecase) {},
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]any{
				"error": "validation error",
			},
		},
		{
			name: "invalid request - short password",
			body: map[string]any{
				"company_id": 1,
				"name":       "Test User",
				"email":      "test@example.com",
				"password":   "short",
			},
			prepare:    func(_ *mock.MockUsecase) {},
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]any{
				"error": "validation error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock.NewMockUsecase(ctrl)
			tt.prepare(mockUsecase)

			handler := auth.NewHandler(mockUsecase, validator.New())

			gin.SetMode(gin.TestMode)

			r := gin.New()
			r.POST("/register", handler.Register)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var resp map[string]any

			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			for k, v := range tt.wantBody {
				assert.Equal(t, v, resp[k])
			}
		})
	}
}

func TestHandler_Login(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       map[string]any
		prepare    func(m *mock.MockUsecase)
		wantStatus int
		wantBody   map[string]any
	}{
		{
			name: "success",
			body: map[string]any{
				"email":    "test@example.com",
				"password": "password123",
			},
			prepare: func(m *mock.MockUsecase) {
				m.EXPECT().
					Login(gomock.Any(), &usecase.Input{
						Email:    "test@example.com",
						Password: "password123",
					}).
					Return(&usecase.TokenPair{
						AccessToken:           "access-token",
						AccessTokenExpiresAt:  time.Date(2024, 1, 1, 10, 15, 0, 0, time.UTC),
						RefreshToken:          "refresh-token",
						RefreshTokenExpiresAt: time.Date(2024, 1, 8, 10, 0, 0, 0, time.UTC),
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]any{
				"access_token":  "access-token",
				"token_type":    "Bearer",
				"refresh_token": "refresh-token",
			},
		},
		{
			name: "invalid credentials",
			body: map[string]any{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			prepare: func(m *mock.MockUsecase) {
				m.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Return(nil, usecase.ErrInvalidCredentials)
			},
			wantStatus: http.StatusUnauthorized,
			wantBody: map[string]any{
				"error": "invalid credentials",
			},
		},
		{
			name: "invalid request - missing password",
			body: map[string]any{
				"email": "test@example.com",
			},
			prepare:    func(_ *mock.MockUsecase) {},
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]any{
				"error": "validation error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock.NewMockUsecase(ctrl)
			tt.prepare(mockUsecase)

			handler := auth.NewHandler(mockUsecase, validator.New())

			gin.SetMode(gin.TestMode)

			r := gin.New()
			r.POST("/login", handler.Login)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var resp map[string]any

			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			for k, v := range tt.wantBody {
				assert.Equal(t, v, resp[k])
			}
		})
	}
}

func TestHandler_RefreshToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       map[string]any
		prepare    func(m *mock.MockUsecase)
		wantStatus int
		wantBody   map[string]any
	}{
		{
			name: "success",
			body: map[string]any{
				"refresh_token": "valid-refresh-token",
			},
			prepare: func(m *mock.MockUsecase) {
				m.EXPECT().
					RefreshToken(gomock.Any(), "valid-refresh-token").
					Return(&usecase.TokenPair{
						AccessToken:           "new-access-token",
						AccessTokenExpiresAt:  time.Date(2024, 1, 1, 10, 15, 0, 0, time.UTC),
						RefreshToken:          "new-refresh-token",
						RefreshTokenExpiresAt: time.Date(2024, 1, 8, 10, 0, 0, 0, time.UTC),
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]any{
				"access_token":  "new-access-token",
				"token_type":    "Bearer",
				"refresh_token": "new-refresh-token",
			},
		},
		{
			name: "invalid refresh token",
			body: map[string]any{
				"refresh_token": "invalid-token",
			},
			prepare: func(m *mock.MockUsecase) {
				m.EXPECT().
					RefreshToken(gomock.Any(), "invalid-token").
					Return(nil, usecase.ErrInvalidCredentials)
			},
			wantStatus: http.StatusUnauthorized,
			wantBody: map[string]any{
				"error": "invalid or expired refresh token",
			},
		},
		{
			name:       "missing refresh token",
			body:       map[string]any{},
			prepare:    func(_ *mock.MockUsecase) {},
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]any{
				"error": "validation error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock.NewMockUsecase(ctrl)
			tt.prepare(mockUsecase)

			handler := auth.NewHandler(mockUsecase, validator.New())

			gin.SetMode(gin.TestMode)

			r := gin.New()
			r.POST("/refresh", handler.RefreshToken)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var resp map[string]any

			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			for k, v := range tt.wantBody {
				assert.Equal(t, v, resp[k])
			}
		})
	}
}
