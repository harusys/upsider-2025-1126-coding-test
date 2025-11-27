package auth

import (
	"context"
	"errors"

	"github.com/harusys/super-shiharai-kun/internal/domain"
	"github.com/harusys/super-shiharai-kun/internal/domain/entity"
	"github.com/harusys/super-shiharai-kun/internal/domain/repository"
	"github.com/harusys/super-shiharai-kun/internal/infrastructure"
	"github.com/harusys/super-shiharai-kun/internal/infrastructure/security"
	"github.com/harusys/super-shiharai-kun/pkg/ctxutil"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type usecaseImpl struct {
	userRepo   repository.UserRepository
	jwtService *security.JWTService
}

// NewUsecase creates a new auth Usecase.
func NewUsecase(
	userRepo repository.UserRepository,
	jwtService *security.JWTService,
) Usecase {
	return &usecaseImpl{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (u *usecaseImpl) Register(ctx context.Context, input *RegisterInput) (*TokenPair, error) {
	// Check if email already exists
	exists, err := u.userRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := ctxutil.HashPassword(ctx, input.Password, infrastructure.BcryptCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entity.User{
		CompanyID:    input.CompanyID,
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hashedPassword,
	}

	created, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	return u.generateTokenPair(ctx, created.ID, created.CompanyID)
}

func (u *usecaseImpl) Login(ctx context.Context, input *Input) (*TokenPair, error) {
	// Get user by email
	user, err := u.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	// Verify password
	if !security.CheckPassword(input.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	return u.generateTokenPair(ctx, user.ID, user.CompanyID)
}

func (u *usecaseImpl) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Validate refresh token
	claims, err := u.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Verify user still exists
	_, err = u.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	// Generate new tokens
	return u.generateTokenPair(ctx, claims.UserID, claims.CompanyID)
}

func (u *usecaseImpl) generateTokenPair(
	ctx context.Context,
	userID, companyID int64,
) (*TokenPair, error) {
	accessToken, accessExpiresAt, err := u.jwtService.GenerateAccessToken(ctx, userID, companyID)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiresAt, err := u.jwtService.GenerateRefreshToken(ctx, userID, companyID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}, nil
}
