//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

package repository

import (
	"context"
	"time"

	"github.com/harusys/super-shiharai-kun/internal/domain/entity"
)

// InvoiceRepository defines the interface for invoice data access.
type InvoiceRepository interface {
	GetByID(ctx context.Context, id int64) (*entity.Invoice, error)
	GetByIDAndCompanyID(ctx context.Context, id, companyID int64) (*entity.Invoice, error)
	GetByCompanyID(ctx context.Context, companyID int64) ([]*entity.Invoice, error)
	GetByCompanyIDAndDateRange(ctx context.Context, companyID int64, startDate, endDate time.Time) ([]*entity.Invoice, error)
	Create(ctx context.Context, invoice *entity.Invoice) (*entity.Invoice, error)
	UpdateStatus(ctx context.Context, id int64, status entity.InvoiceStatus) (*entity.Invoice, error)
}
