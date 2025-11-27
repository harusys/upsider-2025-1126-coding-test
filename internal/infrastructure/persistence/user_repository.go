package persistence

import (
	"context"
	"errors"

	"github.com/harusys/super-shiharai-kun/internal/domain"
	"github.com/harusys/super-shiharai-kun/internal/domain/entity"
	"github.com/harusys/super-shiharai-kun/internal/domain/repository"
	"github.com/harusys/super-shiharai-kun/internal/infrastructure/database/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &userRepository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return toUserEntity(&user), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return toUserEntity(&user), nil
}

func (r *userRepository) GetByCompanyID(
	ctx context.Context,
	companyID int64,
) ([]*entity.User, error) {
	users, err := r.queries.GetUsersByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.User, len(users))
	for i, u := range users {
		result[i] = toUserEntity(&u)
	}

	return result, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	created, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		CompanyID:    user.CompanyID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		return nil, err
	}

	return toUserEntity(&created), nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	updated, err := r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return toUserEntity(&updated), nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	return r.queries.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: passwordHash,
	})
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.queries.ExistsUserByEmail(ctx, email)
}

func toUserEntity(u *sqlc.User) *entity.User {
	return &entity.User{
		ID:           u.ID,
		CompanyID:    u.CompanyID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}
}
