//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

package repository

import (
	"context"

	"github.com/harusys/super-shiharai-kun/internal/domain/entity"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByCompanyID(ctx context.Context, companyID int64) ([]*entity.User, error)
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdatePassword(ctx context.Context, id int64, passwordHash string) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
