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

type vendorRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

// NewVendorRepository creates a new VendorRepository.
func NewVendorRepository(pool *pgxpool.Pool) repository.VendorRepository {
	return &vendorRepository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r *vendorRepository) GetByID(ctx context.Context, id int64) (*entity.Vendor, error) {
	vendor, err := r.queries.GetVendorByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return toVendorEntity(&vendor), nil
}

func (r *vendorRepository) GetByIDAndCompanyID(
	ctx context.Context,
	id, companyID int64,
) (*entity.Vendor, error) {
	vendor, err := r.queries.GetVendorByIDAndCompanyID(ctx, sqlc.GetVendorByIDAndCompanyIDParams{
		ID:        id,
		CompanyID: companyID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return toVendorEntity(&vendor), nil
}

func (r *vendorRepository) GetByCompanyID(
	ctx context.Context,
	companyID int64,
) ([]*entity.Vendor, error) {
	vendors, err := r.queries.GetVendorsByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Vendor, len(vendors))
	for i, v := range vendors {
		result[i] = toVendorEntity(&v)
	}

	return result, nil
}

func (r *vendorRepository) Create(
	ctx context.Context,
	vendor *entity.Vendor,
) (*entity.Vendor, error) {
	created, err := r.queries.CreateVendor(ctx, sqlc.CreateVendorParams{
		CompanyID:          vendor.CompanyID,
		Name:               vendor.Name,
		RepresentativeName: vendor.RepresentativeName,
		PhoneNumber:        vendor.PhoneNumber,
		ZipCode:            vendor.ZipCode,
		Address:            vendor.Address,
	})
	if err != nil {
		return nil, err
	}

	return toVendorEntity(&created), nil
}

func (r *vendorRepository) Update(
	ctx context.Context,
	vendor *entity.Vendor,
) (*entity.Vendor, error) {
	updated, err := r.queries.UpdateVendor(ctx, sqlc.UpdateVendorParams{
		ID:                 vendor.ID,
		Name:               vendor.Name,
		RepresentativeName: vendor.RepresentativeName,
		PhoneNumber:        vendor.PhoneNumber,
		ZipCode:            vendor.ZipCode,
		Address:            vendor.Address,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return toVendorEntity(&updated), nil
}

func toVendorEntity(v *sqlc.Vendor) *entity.Vendor {
	return &entity.Vendor{
		ID:                 v.ID,
		CompanyID:          v.CompanyID,
		Name:               v.Name,
		RepresentativeName: v.RepresentativeName,
		PhoneNumber:        v.PhoneNumber,
		ZipCode:            v.ZipCode,
		Address:            v.Address,
		CreatedAt:          v.CreatedAt.Time,
		UpdatedAt:          v.UpdatedAt.Time,
	}
}

type vendorBankAccountRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

// NewVendorBankAccountRepository creates a new VendorBankAccountRepository.
func NewVendorBankAccountRepository(pool *pgxpool.Pool) repository.VendorBankAccountRepository {
	return &vendorBankAccountRepository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r *vendorBankAccountRepository) GetByID(
	ctx context.Context,
	id int64,
) (*entity.VendorBankAccount, error) {
	account, err := r.queries.GetVendorBankAccountByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return toVendorBankAccountEntity(&account), nil
}

func (r *vendorBankAccountRepository) GetByIDAndVendorID(
	ctx context.Context,
	id, vendorID int64,
) (*entity.VendorBankAccount, error) {
	account, err := r.queries.GetVendorBankAccountByIDAndVendorID(
		ctx,
		sqlc.GetVendorBankAccountByIDAndVendorIDParams{
			ID:       id,
			VendorID: vendorID,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return toVendorBankAccountEntity(&account), nil
}

func (r *vendorBankAccountRepository) GetByVendorID(
	ctx context.Context,
	vendorID int64,
) ([]*entity.VendorBankAccount, error) {
	accounts, err := r.queries.GetVendorBankAccountsByVendorID(ctx, vendorID)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.VendorBankAccount, len(accounts))
	for i, a := range accounts {
		result[i] = toVendorBankAccountEntity(&a)
	}

	return result, nil
}

func (r *vendorBankAccountRepository) Create(
	ctx context.Context,
	account *entity.VendorBankAccount,
) (*entity.VendorBankAccount, error) {
	created, err := r.queries.CreateVendorBankAccount(ctx, sqlc.CreateVendorBankAccountParams{
		VendorID:          account.VendorID,
		BankName:          account.BankName,
		BranchName:        account.BranchName,
		AccountNumber:     account.AccountNumber,
		AccountHolderName: account.AccountHolderName,
	})
	if err != nil {
		return nil, err
	}

	return toVendorBankAccountEntity(&created), nil
}

func (r *vendorBankAccountRepository) Update(
	ctx context.Context,
	account *entity.VendorBankAccount,
) (*entity.VendorBankAccount, error) {
	updated, err := r.queries.UpdateVendorBankAccount(ctx, sqlc.UpdateVendorBankAccountParams{
		ID:                account.ID,
		BankName:          account.BankName,
		BranchName:        account.BranchName,
		AccountNumber:     account.AccountNumber,
		AccountHolderName: account.AccountHolderName,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return toVendorBankAccountEntity(&updated), nil
}

func toVendorBankAccountEntity(a *sqlc.VendorBankAccount) *entity.VendorBankAccount {
	return &entity.VendorBankAccount{
		ID:                a.ID,
		VendorID:          a.VendorID,
		BankName:          a.BankName,
		BranchName:        a.BranchName,
		AccountNumber:     a.AccountNumber,
		AccountHolderName: a.AccountHolderName,
		CreatedAt:         a.CreatedAt.Time,
		UpdatedAt:         a.UpdatedAt.Time,
	}
}
