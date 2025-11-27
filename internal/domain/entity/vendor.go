package entity

import "time"

// Vendor represents a vendor (payment recipient) entity belonging to a company.
type Vendor struct {
	ID                 int64
	CompanyID          int64
	Name               string
	RepresentativeName string
	PhoneNumber        string
	ZipCode            string
	Address            string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// VendorBankAccount represents a bank account belonging to a vendor.
type VendorBankAccount struct {
	ID                int64
	VendorID          int64
	BankName          string
	BranchName        string
	AccountNumber     string
	AccountHolderName string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
