package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

// InvoiceStatus represents the status of an invoice.
type InvoiceStatus string

const (
	InvoiceStatusPending    InvoiceStatus = "pending"
	InvoiceStatusProcessing InvoiceStatus = "processing"
	InvoiceStatusPaid       InvoiceStatus = "paid"
	InvoiceStatusError      InvoiceStatus = "error"
)

// Invoice represents an invoice entity.
type Invoice struct {
	ID                  int64
	CompanyID           int64
	VendorID            int64
	VendorBankAccountID int64
	IssueDate           time.Time
	PaymentAmount       int64           // 支払金額
	Fee                 int64           // 手数料
	FeeRate             decimal.Decimal // 手数料率 (default: 0.04)
	Tax                 int64           // 消費税
	TaxRate             decimal.Decimal // 消費税率 (default: 0.10)
	TotalAmount         int64           // 請求金額 (payment_amount + fee + tax)
	DueDate             time.Time       // 支払期日
	Status              InvoiceStatus
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
