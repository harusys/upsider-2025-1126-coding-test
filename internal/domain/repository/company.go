//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

package repository

import (
	"context"

	"github.com/harusys/super-shiharai-kun/internal/domain/entity"
)

// CompanyRepository defines the interface for company data access.
type CompanyRepository interface {
	GetByID(ctx context.Context, id int64) (*entity.Company, error)
	Create(ctx context.Context, company *entity.Company) (*entity.Company, error)
	Update(ctx context.Context, company *entity.Company) (*entity.Company, error)
}
