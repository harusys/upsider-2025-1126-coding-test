package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/harusys/super-shiharai-kun/internal/domain"
	"github.com/harusys/super-shiharai-kun/internal/domain/entity"
	"github.com/harusys/super-shiharai-kun/internal/domain/repository"
	"github.com/harusys/super-shiharai-kun/internal/infrastructure/database/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type invoiceRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

// NewInvoiceRepository creates a new InvoiceRepository.
func NewInvoiceRepository(pool *pgxpool.Pool) repository.InvoiceRepository {
	return &invoiceRepository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r *invoiceRepository) GetByID(ctx context.Context, id int64) (*entity.Invoice, error) {
	invoice, err := r.queries.GetInvoiceByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return toInvoiceEntity(&invoice), nil
}

func (r *invoiceRepository) GetByIDAndCompanyID(
	ctx context.Context,
	id, companyID int64,
) (*entity.Invoice, error) {
	invoice, err := r.queries.GetInvoiceByIDAndCompanyID(ctx, sqlc.GetInvoiceByIDAndCompanyIDParams{
		ID:        id,
		CompanyID: companyID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return toInvoiceEntity(&invoice), nil
}

func (r *invoiceRepository) GetByCompanyID(
	ctx context.Context,
	companyID int64,
) ([]*entity.Invoice, error) {
	invoices, err := r.queries.GetInvoicesByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Invoice, len(invoices))
	for i, inv := range invoices {
		result[i] = toInvoiceEntity(&inv)
	}

	return result, nil
}

func (r *invoiceRepository) GetByCompanyIDAndDateRange(
	ctx context.Context,
	companyID int64,
	startDate, endDate time.Time,
) ([]*entity.Invoice, error) {
	invoices, err := r.queries.GetInvoicesByCompanyIDAndDateRange(
		ctx,
		sqlc.GetInvoicesByCompanyIDAndDateRangeParams{
			CompanyID: companyID,
			DueDate:   toPgDate(startDate),
			DueDate_2: toPgDate(endDate),
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Invoice, len(invoices))
	for i, inv := range invoices {
		result[i] = toInvoiceEntity(&inv)
	}

	return result, nil
}

func (r *invoiceRepository) Create(
	ctx context.Context,
	invoice *entity.Invoice,
) (*entity.Invoice, error) {
	created, err := r.queries.CreateInvoice(ctx, sqlc.CreateInvoiceParams{
		CompanyID:           invoice.CompanyID,
		VendorID:            invoice.VendorID,
		VendorBankAccountID: invoice.VendorBankAccountID,
		IssueDate:           toPgDate(invoice.IssueDate),
		PaymentAmount:       invoice.PaymentAmount,
		Fee:                 invoice.Fee,
		FeeRate:             invoice.FeeRate,
		Tax:                 invoice.Tax,
		TaxRate:             invoice.TaxRate,
		TotalAmount:         invoice.TotalAmount,
		DueDate:             toPgDate(invoice.DueDate),
		Status:              string(invoice.Status),
	})
	if err != nil {
		return nil, err
	}

	return toInvoiceEntity(&created), nil
}

func (r *invoiceRepository) UpdateStatus(
	ctx context.Context,
	id int64,
	status entity.InvoiceStatus,
) (*entity.Invoice, error) {
	updated, err := r.queries.UpdateInvoiceStatus(ctx, sqlc.UpdateInvoiceStatusParams{
		ID:     id,
		Status: string(status),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return toInvoiceEntity(&updated), nil
}

func toInvoiceEntity(i *sqlc.Invoice) *entity.Invoice {
	return &entity.Invoice{
		ID:                  i.ID,
		CompanyID:           i.CompanyID,
		VendorID:            i.VendorID,
		VendorBankAccountID: i.VendorBankAccountID,
		IssueDate:           i.IssueDate.Time,
		PaymentAmount:       i.PaymentAmount,
		Fee:                 i.Fee,
		FeeRate:             i.FeeRate,
		Tax:                 i.Tax,
		TaxRate:             i.TaxRate,
		TotalAmount:         i.TotalAmount,
		DueDate:             i.DueDate.Time,
		Status:              entity.InvoiceStatus(i.Status),
		CreatedAt:           i.CreatedAt.Time,
		UpdatedAt:           i.UpdatedAt.Time,
	}
}

func toPgDate(t time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  t,
		Valid: true,
	}
}
