package invoice

import (
	"context"

	"github.com/harusys/super-shiharai-kun/internal/domain"
	"github.com/harusys/super-shiharai-kun/internal/domain/entity"
	"github.com/harusys/super-shiharai-kun/internal/domain/repository"
	"github.com/harusys/super-shiharai-kun/internal/domain/service"
)

type usecaseImpl struct {
	invoiceRepo       repository.InvoiceRepository
	vendorRepo        repository.VendorRepository
	bankAccountRepo   repository.VendorBankAccountRepository
	invoiceCalculator *service.InvoiceCalculator
}

// NewUsecase creates a new invoice Usecase.
func NewUsecase(
	invoiceRepo repository.InvoiceRepository,
	vendorRepo repository.VendorRepository,
	bankAccountRepo repository.VendorBankAccountRepository,
	invoiceCalculator *service.InvoiceCalculator,
) Usecase {
	return &usecaseImpl{
		invoiceRepo:       invoiceRepo,
		vendorRepo:        vendorRepo,
		bankAccountRepo:   bankAccountRepo,
		invoiceCalculator: invoiceCalculator,
	}
}

func (u *usecaseImpl) Create(
	ctx context.Context,
	input *CreateInput,
) (*entity.Invoice, error) {
	// Verify vendor belongs to company
	_, err := u.vendorRepo.GetByIDAndCompanyID(ctx, input.VendorID, input.CompanyID)
	if err != nil {
		return nil, err
	}

	// Verify bank account belongs to vendor
	_, err = u.bankAccountRepo.GetByIDAndVendorID(ctx, input.VendorBankAccountID, input.VendorID)
	if err != nil {
		return nil, err
	}

	// Calculate invoice amounts
	result := u.invoiceCalculator.Calculate(input.PaymentAmount)

	// Create invoice
	inv := &entity.Invoice{
		CompanyID:           input.CompanyID,
		VendorID:            input.VendorID,
		VendorBankAccountID: input.VendorBankAccountID,
		IssueDate:           input.IssueDate,
		PaymentAmount:       result.PaymentAmount,
		Fee:                 result.Fee,
		FeeRate:             result.FeeRate,
		Tax:                 result.Tax,
		TaxRate:             result.TaxRate,
		TotalAmount:         result.TotalAmount,
		DueDate:             input.DueDate,
		Status:              entity.InvoiceStatusPending,
	}

	return u.invoiceRepo.Create(ctx, inv)
}

func (u *usecaseImpl) List(
	ctx context.Context,
	input *ListInput,
) ([]*entity.Invoice, error) {
	// If date range is specified, use date range query
	if input.StartDate != nil && input.EndDate != nil {
		return u.invoiceRepo.GetByCompanyIDAndDateRange(
			ctx,
			input.CompanyID,
			*input.StartDate,
			*input.EndDate,
		)
	}

	// Otherwise, return all invoices for company
	return u.invoiceRepo.GetByCompanyID(ctx, input.CompanyID)
}

func (u *usecaseImpl) GetByID(
	ctx context.Context,
	companyID, invoiceID int64,
) (*entity.Invoice, error) {
	inv, err := u.invoiceRepo.GetByIDAndCompanyID(ctx, invoiceID, companyID)
	if err != nil {
		return nil, err
	}

	// Ensure company authorization
	if inv.CompanyID != companyID {
		return nil, domain.ErrNotFound
	}

	return inv, nil
}
