//go:generate mockgen -source=$GOFILE -destination=mock/mock_usecase.go -package=mock

package invoice

import (
	"context"
	"time"

	"github.com/harusys/super-shiharai-kun/internal/domain/entity"
)

// CreateInput is the input for creating an invoice.
type CreateInput struct {
	CompanyID           int64
	VendorID            int64
	VendorBankAccountID int64
	IssueDate           time.Time
	PaymentAmount       int64
	DueDate             time.Time
}

// ListInput is the input for listing invoices.
type ListInput struct {
	CompanyID int64
	StartDate *time.Time
	EndDate   *time.Time
}

// Usecase defines invoice operations.
type Usecase interface {
	// Create creates a new invoice with calculated amounts.
	Create(ctx context.Context, input *CreateInput) (*entity.Invoice, error)
	// List returns invoices for a company, optionally filtered by date range.
	List(ctx context.Context, input *ListInput) ([]*entity.Invoice, error)
	// GetByID returns an invoice by ID (with company authorization check).
	GetByID(ctx context.Context, companyID, invoiceID int64) (*entity.Invoice, error)
}
