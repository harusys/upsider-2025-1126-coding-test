//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

package repository

import (
	"context"

	"github.com/harusys/super-shiharai-kun/internal/domain/entity"
)

// VendorRepository defines the interface for vendor data access.
type VendorRepository interface {
	GetByID(ctx context.Context, id int64) (*entity.Vendor, error)
	GetByIDAndCompanyID(ctx context.Context, id, companyID int64) (*entity.Vendor, error)
	GetByCompanyID(ctx context.Context, companyID int64) ([]*entity.Vendor, error)
	Create(ctx context.Context, vendor *entity.Vendor) (*entity.Vendor, error)
	Update(ctx context.Context, vendor *entity.Vendor) (*entity.Vendor, error)
}

// VendorBankAccountRepository defines the interface for vendor bank account data access.
type VendorBankAccountRepository interface {
	GetByID(ctx context.Context, id int64) (*entity.VendorBankAccount, error)
	GetByIDAndVendorID(ctx context.Context, id, vendorID int64) (*entity.VendorBankAccount, error)
	GetByVendorID(ctx context.Context, vendorID int64) ([]*entity.VendorBankAccount, error)
	Create(
		ctx context.Context,
		account *entity.VendorBankAccount,
	) (*entity.VendorBankAccount, error)
	Update(
		ctx context.Context,
		account *entity.VendorBankAccount,
	) (*entity.VendorBankAccount, error)
}
