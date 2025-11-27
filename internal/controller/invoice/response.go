package invoice

import (
	"time"

	"github.com/harusys/super-shiharai-kun/internal/domain/entity"
)

// Response is the response body for an invoice.
type Response struct {
	ID                  int64     `json:"id"`
	CompanyID           int64     `json:"company_id"`
	VendorID            int64     `json:"vendor_id"`
	VendorBankAccountID int64     `json:"vendor_bank_account_id"`
	IssueDate           string    `json:"issue_date"`
	PaymentAmount       int64     `json:"payment_amount"`
	Fee                 int64     `json:"fee"`
	FeeRate             string    `json:"fee_rate"`
	Tax                 int64     `json:"tax"`
	TaxRate             string    `json:"tax_rate"`
	TotalAmount         int64     `json:"total_amount"`
	DueDate             string    `json:"due_date"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ToResponse converts an entity.Invoice to Response.
func ToResponse(inv *entity.Invoice) *Response {
	return &Response{
		ID:                  inv.ID,
		CompanyID:           inv.CompanyID,
		VendorID:            inv.VendorID,
		VendorBankAccountID: inv.VendorBankAccountID,
		IssueDate:           inv.IssueDate.Format("2006-01-02"),
		PaymentAmount:       inv.PaymentAmount,
		Fee:                 inv.Fee,
		FeeRate:             inv.FeeRate.String(),
		Tax:                 inv.Tax,
		TaxRate:             inv.TaxRate.String(),
		TotalAmount:         inv.TotalAmount,
		DueDate:             inv.DueDate.Format("2006-01-02"),
		Status:              string(inv.Status),
		CreatedAt:           inv.CreatedAt,
		UpdatedAt:           inv.UpdatedAt,
	}
}

// ToResponses converts a slice of entity.Invoice to a slice of Response.
func ToResponses(invoices []*entity.Invoice) []*Response {
	responses := make([]*Response, len(invoices))
	for i, inv := range invoices {
		responses[i] = ToResponse(inv)
	}

	return responses
}

// ErrorResponse is the standard error response body.
type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

// NewErrorResponse creates a new ErrorResponse.
func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Error: message,
	}
}

// NewValidationErrorResponse creates a new ErrorResponse for validation errors.
func NewValidationErrorResponse(details map[string]string) *ErrorResponse {
	return &ErrorResponse{
		Error:   "validation error",
		Details: details,
	}
}
