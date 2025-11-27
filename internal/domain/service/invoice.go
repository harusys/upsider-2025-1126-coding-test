package service

import (
	"github.com/shopspring/decimal"
)

// DefaultFeeRate is the default fee rate (4%).
var DefaultFeeRate = decimal.NewFromFloat(0.04)

// DefaultTaxRate is the default tax rate (10%).
var DefaultTaxRate = decimal.NewFromFloat(0.10)

// InvoiceCalculator provides invoice amount calculation logic.
type InvoiceCalculator struct{}

// NewInvoiceCalculator creates a new InvoiceCalculator.
func NewInvoiceCalculator() *InvoiceCalculator {
	return &InvoiceCalculator{}
}

// CalculationResult holds the result of invoice calculation.
type CalculationResult struct {
	PaymentAmount int64           // 支払金額
	Fee           int64           // 手数料
	FeeRate       decimal.Decimal // 手数料率
	Tax           int64           // 消費税
	TaxRate       decimal.Decimal // 消費税率
	TotalAmount   int64           // 請求金額
}

// Calculate calculates the invoice amounts based on payment amount.
// Formula: total = payment + fee + tax
//
//	fee = payment * feeRate
//	tax = fee * taxRate
//
// Example: payment=10000, feeRate=0.04, taxRate=0.10
//
//	fee = 10000 * 0.04 = 400
//	tax = 400 * 0.10 = 40
//	total = 10000 + 400 + 40 = 10440
func (c *InvoiceCalculator) Calculate(paymentAmount int64) *CalculationResult {
	return c.CalculateWithRates(paymentAmount, DefaultFeeRate, DefaultTaxRate)
}

// CalculateWithRates calculates the invoice amounts with custom rates.
func (c *InvoiceCalculator) CalculateWithRates(paymentAmount int64, feeRate, taxRate decimal.Decimal) *CalculationResult {
	payment := decimal.NewFromInt(paymentAmount)

	// fee = payment * feeRate (truncate to integer)
	fee := payment.Mul(feeRate).Truncate(0)

	// tax = fee * taxRate (truncate to integer)
	tax := fee.Mul(taxRate).Truncate(0)

	// total = payment + fee + tax
	total := payment.Add(fee).Add(tax)

	return &CalculationResult{
		PaymentAmount: paymentAmount,
		Fee:           fee.IntPart(),
		FeeRate:       feeRate,
		Tax:           tax.IntPart(),
		TaxRate:       taxRate,
		TotalAmount:   total.IntPart(),
	}
}
