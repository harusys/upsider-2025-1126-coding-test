package service_test

import (
	"testing"

	"github.com/harusys/super-shiharai-kun/internal/domain"
	"github.com/harusys/super-shiharai-kun/internal/domain/service"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceCalculator_Calculate(t *testing.T) {
	t.Parallel()

	defaultFeeRate := decimal.RequireFromString(domain.DefaultFeeRateStr)
	defaultTaxRate := decimal.RequireFromString(domain.DefaultTaxRateStr)

	tests := []struct {
		name string
		want *service.CalculationResult
	}{
		{
			name: "standard calculation",
			want: &service.CalculationResult{
				PaymentAmount: 10000,
				Fee:           400, // 10000 * 0.04 = 400
				FeeRate:       defaultFeeRate,
				Tax:           40, // 400 * 0.10 = 40
				TaxRate:       defaultTaxRate,
				TotalAmount:   10440, // 10000 + 400 + 40 = 10440
			},
		},
		{
			name: "large amount",
			want: &service.CalculationResult{
				PaymentAmount: 1000000,
				Fee:           40000, // 1000000 * 0.04 = 40000
				FeeRate:       defaultFeeRate,
				Tax:           4000, // 40000 * 0.10 = 4000
				TaxRate:       defaultTaxRate,
				TotalAmount:   1044000,
			},
		},
		{
			name: "small amount with truncation",
			want: &service.CalculationResult{
				PaymentAmount: 123,
				Fee:           4, // 123 * 0.04 = 4.92 -> truncate to 4
				FeeRate:       defaultFeeRate,
				Tax:           0, // 4 * 0.10 = 0.4 -> truncate to 0
				TaxRate:       defaultTaxRate,
				TotalAmount:   127, // 123 + 4 + 0 = 127
			},
		},
		{
			name: "zero amount",
			want: &service.CalculationResult{
				PaymentAmount: 0,
				Fee:           0,
				FeeRate:       defaultFeeRate,
				Tax:           0,
				TaxRate:       defaultTaxRate,
				TotalAmount:   0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			calc := service.NewInvoiceCalculator()
			result := calc.Calculate(tt.want.PaymentAmount)

			assert.Equal(t, tt.want, result)
		})
	}
}

func TestInvoiceCalculator_CalculateWithRates(t *testing.T) {
	t.Parallel()

	calc := service.NewInvoiceCalculator()

	// Custom rates: 5% fee, 8% tax
	feeRate := decimal.RequireFromString("0.05")
	taxRate := decimal.RequireFromString("0.08")

	result := calc.CalculateWithRates(10000, feeRate, taxRate)

	want := &service.CalculationResult{
		PaymentAmount: 10000,
		Fee:           500, // 10000 * 0.05 = 500
		FeeRate:       feeRate,
		Tax:           40, // 500 * 0.08 = 40
		TaxRate:       taxRate,
		TotalAmount:   10540,
	}
	assert.Equal(t, want, result)
}

func TestNewInvoiceCalculatorWithRates(t *testing.T) {
	t.Parallel()

	feeRate := decimal.RequireFromString("0.03")
	taxRate := decimal.RequireFromString("0.08")

	calc := service.NewInvoiceCalculatorWithRates(feeRate, taxRate)
	result := calc.Calculate(10000)

	want := &service.CalculationResult{
		PaymentAmount: 10000,
		Fee:           300, // 10000 * 0.03 = 300
		FeeRate:       feeRate,
		Tax:           24, // 300 * 0.08 = 24
		TaxRate:       taxRate,
		TotalAmount:   10324,
	}
	assert.Equal(t, want, result)
}
