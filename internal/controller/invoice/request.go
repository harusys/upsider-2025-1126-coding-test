package invoice

// CreateRequest is the request body for creating an invoice.
type CreateRequest struct {
	VendorID            int64  `json:"vendor_id"              validate:"required,gt=0"`
	VendorBankAccountID int64  `json:"vendor_bank_account_id" validate:"required,gt=0"`
	IssueDate           string `json:"issue_date"             validate:"required,datetime=2006-01-02"`
	PaymentAmount       int64  `json:"payment_amount"         validate:"required,gt=0"`
	DueDate             string `json:"due_date"               validate:"required,datetime=2006-01-02"`
}

// ListRequest is the query parameters for listing invoices.
type ListRequest struct {
	StartDate string `query:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate   string `query:"end_date"   validate:"omitempty,datetime=2006-01-02"`
}
