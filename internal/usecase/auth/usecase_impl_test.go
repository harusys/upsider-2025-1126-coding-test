package auth_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/harusys/super-shiharai-kun/internal/domain"
	"github.com/harusys/super-shiharai-kun/internal/domain/entity"
	"github.com/harusys/super-shiharai-kun/internal/domain/repository/mock"
	"github.com/harusys/super-shiharai-kun/internal/infrastructure/security"
	"github.com/harusys/super-shiharai-kun/internal/usecase/auth"
	"github.com/harusys/super-shiharai-kun/pkg/ctxutil/ctxutiltest"
	"github.com/harusys/super-shiharai-kun/pkg/timeutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUsecaseImpl_Register(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   *auth.RegisterInput
		prepare func(ctx context.Context, c *controllers)
		want    *auth.TokenPair
		wantErr error
	}{
		{
			name: "success",
			input: &auth.RegisterInput{
				CompanyID: 1,
				Name:      "Test User",
				Email:     "test@example.com",
				Password:  "password123",
			},
			prepare: func(ctx context.Context, c *controllers) {
				c.ctxProvider.SetAsiaTokyo(t, "2024-01-01 01:23:45")

				fixedHash := "hashed_password123"
				c.ctxProvider.PasswordHash = &fixedHash

				c.userRepo.EXPECT().
					ExistsByEmail(ctx, "test@example.com").
					Return(false, nil)
				c.userRepo.EXPECT().
					Create(ctx, &entity.User{
						CompanyID:    1,
						Name:         "Test User",
						Email:        "test@example.com",
						PasswordHash: fixedHash,
					}).
					Return(&entity.User{
						ID:           1,
						CompanyID:    1,
						Name:         "Test User",
						Email:        "test@example.com",
						PasswordHash: fixedHash,
						CreatedAt:    timeutil.AsiaTokyo(t, "2024-01-01 01:23:45"),
						UpdatedAt:    timeutil.AsiaTokyo(t, "2024-01-01 01:23:45"),
					}, nil)
			},
			want: &auth.TokenPair{
				AccessTokenExpiresAt:  timeutil.AsiaTokyo(t, "2024-01-01 01:38:45"),
				RefreshTokenExpiresAt: timeutil.AsiaTokyo(t, "2024-01-08 01:23:45"),
			},
			wantErr: nil,
		},
		{
			name: "email already exists",
			input: &auth.RegisterInput{
				CompanyID: 1,
				Name:      "Test User",
				Email:     "existing@example.com",
				Password:  "password123",
			},
			prepare: func(ctx context.Context, c *controllers) {
				c.userRepo.EXPECT().
					ExistsByEmail(ctx, "existing@example.com").
					Return(true, nil)
			},
			want:    nil,
			wantErr: auth.ErrEmailAlreadyExists,
		},
		{
			name: "repository error on exists check",
			input: &auth.RegisterInput{
				CompanyID: 1,
				Name:      "Test User",
				Email:     "test@example.com",
				Password:  "password123",
			},
			prepare: func(ctx context.Context, c *controllers) {
				c.userRepo.EXPECT().
					ExistsByEmail(ctx, "test@example.com").
					Return(false, errors.New("db error"))
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, uc, c := newUsecase(t)
			defer c.ctrl.Finish()

			if tt.prepare != nil {
				tt.prepare(ctx, c)
			}

			got, err := uc.Register(ctx, tt.input)

			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
				assert.Nil(t, got)

				return
			}

			require.NoError(t, err)
			// トークンはJWT形式であることを検証（3つのドットで区切られた形式）
			assert.Len(t, strings.Split(got.AccessToken, "."), 3)
			assert.Len(t, strings.Split(got.RefreshToken, "."), 3)
			// 有効期限はモック時間から進んでいることを検証
			assert.Equal(t, tt.want.AccessTokenExpiresAt, got.AccessTokenExpiresAt)
			assert.Equal(t, tt.want.RefreshTokenExpiresAt, got.RefreshTokenExpiresAt)
		})
	}
}

func TestUsecaseImpl_Login(t *testing.T) {
	t.Parallel()

	hashedPassword, _ := security.HashPassword("password123")

	tests := []struct {
		name    string
		input   *auth.Input
		prepare func(ctx context.Context, c *controllers)
		want    *auth.TokenPair
		wantErr error
	}{
		{
			name: "success",
			input: &auth.Input{
				Email:    "test@example.com",
				Password: "password123",
			},
			prepare: func(ctx context.Context, c *controllers) {
				c.ctxProvider.SetAsiaTokyo(t, "2024-01-01 10:00:00")

				c.userRepo.EXPECT().
					GetByEmail(ctx, "test@example.com").
					Return(&entity.User{
						ID:           1,
						CompanyID:    1,
						Email:        "test@example.com",
						PasswordHash: hashedPassword,
					}, nil)
			},
			want: &auth.TokenPair{
				AccessTokenExpiresAt:  timeutil.AsiaTokyo(t, "2024-01-01 10:15:00"),
				RefreshTokenExpiresAt: timeutil.AsiaTokyo(t, "2024-01-08 10:00:00"),
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			input: &auth.Input{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			prepare: func(ctx context.Context, c *controllers) {
				c.userRepo.EXPECT().
					GetByEmail(ctx, "notfound@example.com").
					Return(nil, domain.ErrNotFound)
			},
			want:    nil,
			wantErr: auth.ErrInvalidCredentials,
		},
		{
			name: "wrong password",
			input: &auth.Input{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			prepare: func(ctx context.Context, c *controllers) {
				c.userRepo.EXPECT().
					GetByEmail(ctx, "test@example.com").
					Return(&entity.User{
						ID:           1,
						CompanyID:    1,
						Email:        "test@example.com",
						PasswordHash: hashedPassword,
					}, nil)
			},
			want:    nil,
			wantErr: auth.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, uc, c := newUsecase(t)
			defer c.ctrl.Finish()

			if tt.prepare != nil {
				tt.prepare(ctx, c)
			}

			got, err := uc.Login(ctx, tt.input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			// トークンはJWT形式であることを検証（3つのドットで区切られた形式）
			assert.Len(t, strings.Split(got.AccessToken, "."), 3)
			assert.Len(t, strings.Split(got.RefreshToken, "."), 3)
			// 有効期限はモック時間から進んでいることを検証
			assert.Equal(t, tt.want.AccessTokenExpiresAt, got.AccessTokenExpiresAt)
			assert.Equal(t, tt.want.RefreshTokenExpiresAt, got.RefreshTokenExpiresAt)
		})
	}
}

func TestUsecaseImpl_RefreshToken(t *testing.T) {
	t.Parallel()

	jwtService := security.NewJWTService("test-secret-key")
	validToken, _, _ := jwtService.GenerateRefreshToken(context.Background(), 1, 1)
	invalidUserToken, _, _ := jwtService.GenerateRefreshToken(context.Background(), 999, 1)

	tests := []struct {
		name         string
		refreshToken string
		prepare      func(ctx context.Context, c *controllers)
		want         *auth.TokenPair
		wantErr      bool
	}{
		{
			name:         "success",
			refreshToken: validToken,
			prepare: func(ctx context.Context, c *controllers) {
				c.ctxProvider.SetAsiaTokyo(t, "2024-01-01 10:00:00")

				c.userRepo.EXPECT().
					GetByID(ctx, int64(1)).
					Return(&entity.User{ID: 1, CompanyID: 1}, nil)
			},
			want: &auth.TokenPair{
				AccessTokenExpiresAt:  timeutil.AsiaTokyo(t, "2024-01-01 10:15:00"),
				RefreshTokenExpiresAt: timeutil.AsiaTokyo(t, "2024-01-08 10:00:00"),
			},
			wantErr: false,
		},
		{
			name:         "invalid token",
			refreshToken: "invalid-token",
			prepare:      func(_ context.Context, _ *controllers) {},
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "user not found",
			refreshToken: invalidUserToken,
			prepare: func(ctx context.Context, c *controllers) {
				c.userRepo.EXPECT().
					GetByID(ctx, int64(999)).
					Return(nil, domain.ErrNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, uc, c := newUsecase(t)
			defer c.ctrl.Finish()

			if tt.prepare != nil {
				tt.prepare(ctx, c)
			}

			got, err := uc.RefreshToken(ctx, tt.refreshToken)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			// トークンはJWT形式であることを検証（3つのドットで区切られた形式）
			assert.Len(t, strings.Split(got.AccessToken, "."), 3)
			assert.Len(t, strings.Split(got.RefreshToken, "."), 3)
			// 有効期限はモック時間から進んでいることを検証
			assert.Equal(t, tt.want.AccessTokenExpiresAt, got.AccessTokenExpiresAt)
			assert.Equal(t, tt.want.RefreshTokenExpiresAt, got.RefreshTokenExpiresAt)
		})
	}
}

type controllers struct {
	ctrl        *gomock.Controller
	ctxProvider *ctxutiltest.TestContextProvider
	userRepo    *mock.MockUserRepository
}

func newUsecase(t *testing.T) (context.Context, auth.Usecase, *controllers) {
	t.Helper()

	ctxProvider := ctxutiltest.TestContextProvider{}
	ctx := ctxutiltest.TestContext(&ctxProvider)

	ctrl := gomock.NewController(t)
	userRepo := mock.NewMockUserRepository(ctrl)
	jwtService := security.NewJWTService("test-secret-key")
	uc := auth.NewUsecase(userRepo, jwtService)

	return ctx, uc, &controllers{
		ctrl,
		&ctxProvider,
		userRepo,
	}
}
