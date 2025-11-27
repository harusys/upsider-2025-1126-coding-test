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

type companyRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

// NewCompanyRepository creates a new CompanyRepository.
func NewCompanyRepository(pool *pgxpool.Pool) repository.CompanyRepository {
	return &companyRepository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r *companyRepository) GetByID(ctx context.Context, id int64) (*entity.Company, error) {
	company, err := r.queries.GetCompanyByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return toCompanyEntity(&company), nil
}

func (r *companyRepository) Create(
	ctx context.Context,
	company *entity.Company,
) (*entity.Company, error) {
	created, err := r.queries.CreateCompany(ctx, sqlc.CreateCompanyParams{
		Name:               company.Name,
		RepresentativeName: company.RepresentativeName,
		PhoneNumber:        company.PhoneNumber,
		ZipCode:            company.ZipCode,
		Address:            company.Address,
	})
	if err != nil {
		return nil, err
	}

	return toCompanyEntity(&created), nil
}

func (r *companyRepository) Update(
	ctx context.Context,
	company *entity.Company,
) (*entity.Company, error) {
	updated, err := r.queries.UpdateCompany(ctx, sqlc.UpdateCompanyParams{
		ID:                 company.ID,
		Name:               company.Name,
		RepresentativeName: company.RepresentativeName,
		PhoneNumber:        company.PhoneNumber,
		ZipCode:            company.ZipCode,
		Address:            company.Address,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return toCompanyEntity(&updated), nil
}

func toCompanyEntity(c *sqlc.Company) *entity.Company {
	return &entity.Company{
		ID:                 c.ID,
		Name:               c.Name,
		RepresentativeName: c.RepresentativeName,
		PhoneNumber:        c.PhoneNumber,
		ZipCode:            c.ZipCode,
		Address:            c.Address,
		CreatedAt:          c.CreatedAt.Time,
		UpdatedAt:          c.UpdatedAt.Time,
	}
}
