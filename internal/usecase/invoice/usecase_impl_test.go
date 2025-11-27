package invoice_test

import (
	"context"
	"testing"

	"github.com/harusys/super-shiharai-kun/internal/domain"
	"github.com/harusys/super-shiharai-kun/internal/domain/entity"
	"github.com/harusys/super-shiharai-kun/internal/domain/repository/mock"
	"github.com/harusys/super-shiharai-kun/internal/domain/service"
	"github.com/harusys/super-shiharai-kun/internal/usecase/invoice"
	"github.com/harusys/super-shiharai-kun/pkg/ctxutil/ctxutiltest"
	"github.com/harusys/super-shiharai-kun/pkg/timeutil"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func feeRate() decimal.Decimal {
	return decimal.RequireFromString(domain.DefaultFeeRateStr)
}

func taxRate() decimal.Decimal {
	return decimal.RequireFromString(domain.DefaultTaxRateStr)
}

func TestUsecaseImpl_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   *invoice.CreateInput
		prepare func(ctx context.Context, c *controllers)
		want    *entity.Invoice
		wantErr error
	}{
		{
			name: "success",
			input: &invoice.CreateInput{
				CompanyID:           1,
				VendorID:            1,
				VendorBankAccountID: 1,
				IssueDate:           timeutil.AsiaTokyo(t, "2024-01-15 00:00:00"),
				PaymentAmount:       10000,
				DueDate:             timeutil.AsiaTokyo(t, "2024-02-15 00:00:00"),
			},
			prepare: func(ctx context.Context, c *controllers) {
				c.vendorRepo.EXPECT().
					GetByIDAndCompanyID(ctx, int64(1), int64(1)).
					Return(&entity.Vendor{ID: 1, CompanyID: 1}, nil)
				c.bankAccountRepo.EXPECT().
					GetByIDAndVendorID(ctx, int64(1), int64(1)).
					Return(&entity.VendorBankAccount{ID: 1, VendorID: 1}, nil)
				c.invoiceRepo.EXPECT().
					Create(ctx, &entity.Invoice{
						CompanyID:           1,
						VendorID:            1,
						VendorBankAccountID: 1,
						IssueDate:           timeutil.AsiaTokyo(t, "2024-01-15 00:00:00"),
						PaymentAmount:       10000,
						Fee:                 400,
						FeeRate:             feeRate(),
						Tax:                 40,
						TaxRate:             taxRate(),
						TotalAmount:         10440,
						DueDate:             timeutil.AsiaTokyo(t, "2024-02-15 00:00:00"),
						Status:              entity.InvoiceStatusPending,
					}).
					Return(&entity.Invoice{
						ID:                  1,
						CompanyID:           1,
						VendorID:            1,
						VendorBankAccountID: 1,
						IssueDate:           timeutil.AsiaTokyo(t, "2024-01-15 00:00:00"),
						PaymentAmount:       10000,
						Fee:                 400,
						FeeRate:             feeRate(),
						Tax:                 40,
						TaxRate:             taxRate(),
						TotalAmount:         10440,
						DueDate:             timeutil.AsiaTokyo(t, "2024-02-15 00:00:00"),
						Status:              entity.InvoiceStatusPending,
						CreatedAt:           timeutil.AsiaTokyo(t, "2024-01-15 10:00:00"),
						UpdatedAt:           timeutil.AsiaTokyo(t, "2024-01-15 10:00:00"),
					}, nil)
			},
			want: &entity.Invoice{
				ID:                  1,
				CompanyID:           1,
				VendorID:            1,
				VendorBankAccountID: 1,
				IssueDate:           timeutil.AsiaTokyo(t, "2024-01-15 00:00:00"),
				PaymentAmount:       10000,
				Fee:                 400,
				FeeRate:             feeRate(),
				Tax:                 40,
				TaxRate:             taxRate(),
				TotalAmount:         10440,
				DueDate:             timeutil.AsiaTokyo(t, "2024-02-15 00:00:00"),
				Status:              entity.InvoiceStatusPending,
				CreatedAt:           timeutil.AsiaTokyo(t, "2024-01-15 10:00:00"),
				UpdatedAt:           timeutil.AsiaTokyo(t, "2024-01-15 10:00:00"),
			},
			wantErr: nil,
		},
		{
			name: "vendor not found",
			input: &invoice.CreateInput{
				CompanyID:           1,
				VendorID:            999,
				VendorBankAccountID: 1,
				IssueDate:           timeutil.AsiaTokyo(t, "2024-01-15 00:00:00"),
				PaymentAmount:       10000,
				DueDate:             timeutil.AsiaTokyo(t, "2024-02-15 00:00:00"),
			},
			prepare: func(ctx context.Context, c *controllers) {
				c.vendorRepo.EXPECT().
					GetByIDAndCompanyID(ctx, int64(999), int64(1)).
					Return(nil, domain.ErrNotFound)
			},
			want:    nil,
			wantErr: domain.ErrNotFound,
		},
		{
			name: "bank account not found",
			input: &invoice.CreateInput{
				CompanyID:           1,
				VendorID:            1,
				VendorBankAccountID: 999,
				IssueDate:           timeutil.AsiaTokyo(t, "2024-01-15 00:00:00"),
				PaymentAmount:       10000,
				DueDate:             timeutil.AsiaTokyo(t, "2024-02-15 00:00:00"),
			},
			prepare: func(ctx context.Context, c *controllers) {
				c.vendorRepo.EXPECT().
					GetByIDAndCompanyID(ctx, int64(1), int64(1)).
					Return(&entity.Vendor{ID: 1, CompanyID: 1}, nil)
				c.bankAccountRepo.EXPECT().
					GetByIDAndVendorID(ctx, int64(999), int64(1)).
					Return(nil, domain.ErrNotFound)
			},
			want:    nil,
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, uc, c := newUsecase(t)
			defer c.ctrl.Finish()

			if tt.prepare != nil {
				tt.prepare(ctx, c)
			}

			got, err := uc.Create(ctx, tt.input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUsecaseImpl_List(t *testing.T) {
	t.Parallel()

	sampleInvoices := []*entity.Invoice{
		{ID: 1, CompanyID: 1, PaymentAmount: 10000},
		{ID: 2, CompanyID: 1, PaymentAmount: 20000},
	}

	startDate := timeutil.AsiaTokyo(t, "2023-12-01 00:00:00")
	endDate := timeutil.AsiaTokyo(t, "2024-01-01 00:00:00")

	tests := []struct {
		name    string
		input   *invoice.ListInput
		prepare func(ctx context.Context, c *controllers)
		want    []*entity.Invoice
		wantErr error
	}{
		{
			name: "list all invoices",
			input: &invoice.ListInput{
				CompanyID: 1,
			},
			prepare: func(ctx context.Context, c *controllers) {
				c.invoiceRepo.EXPECT().
					GetByCompanyID(ctx, int64(1)).
					Return(sampleInvoices, nil)
			},
			want:    sampleInvoices,
			wantErr: nil,
		},
		{
			name: "list with date range",
			input: &invoice.ListInput{
				CompanyID: 1,
				StartDate: &startDate,
				EndDate:   &endDate,
			},
			prepare: func(ctx context.Context, c *controllers) {
				c.invoiceRepo.EXPECT().
					GetByCompanyIDAndDateRange(ctx, int64(1), startDate, endDate).
					Return(sampleInvoices[:1], nil)
			},
			want:    sampleInvoices[:1],
			wantErr: nil,
		},
		{
			name: "empty result",
			input: &invoice.ListInput{
				CompanyID: 999,
			},
			prepare: func(ctx context.Context, c *controllers) {
				c.invoiceRepo.EXPECT().
					GetByCompanyID(ctx, int64(999)).
					Return([]*entity.Invoice{}, nil)
			},
			want:    []*entity.Invoice{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, uc, c := newUsecase(t)
			defer c.ctrl.Finish()

			if tt.prepare != nil {
				tt.prepare(ctx, c)
			}

			got, err := uc.List(ctx, tt.input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUsecaseImpl_GetByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		companyID int64
		invoiceID int64
		prepare   func(ctx context.Context, c *controllers)
		want      *entity.Invoice
		wantErr   error
	}{
		{
			name:      "success",
			companyID: 1,
			invoiceID: 1,
			prepare: func(ctx context.Context, c *controllers) {
				c.invoiceRepo.EXPECT().
					GetByIDAndCompanyID(ctx, int64(1), int64(1)).
					Return(&entity.Invoice{ID: 1, CompanyID: 1}, nil)
			},
			want:    &entity.Invoice{ID: 1, CompanyID: 1},
			wantErr: nil,
		},
		{
			name:      "not found",
			companyID: 1,
			invoiceID: 999,
			prepare: func(ctx context.Context, c *controllers) {
				c.invoiceRepo.EXPECT().
					GetByIDAndCompanyID(ctx, int64(999), int64(1)).
					Return(nil, domain.ErrNotFound)
			},
			want:    nil,
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, uc, c := newUsecase(t)
			defer c.ctrl.Finish()

			if tt.prepare != nil {
				tt.prepare(ctx, c)
			}

			got, err := uc.GetByID(ctx, tt.companyID, tt.invoiceID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.want, got)
		})
	}
}

type controllers struct {
	ctrl            *gomock.Controller
	ctxProvider     *ctxutiltest.TestContextProvider
	invoiceRepo     *mock.MockInvoiceRepository
	vendorRepo      *mock.MockVendorRepository
	bankAccountRepo *mock.MockVendorBankAccountRepository
}

func newUsecase(t *testing.T) (context.Context, invoice.Usecase, *controllers) {
	t.Helper()

	ctxProvider := ctxutiltest.TestContextProvider{}
	ctx := ctxutiltest.TestContext(&ctxProvider)

	ctrl := gomock.NewController(t)
	invoiceRepo := mock.NewMockInvoiceRepository(ctrl)
	vendorRepo := mock.NewMockVendorRepository(ctrl)
	bankAccountRepo := mock.NewMockVendorBankAccountRepository(ctrl)
	calculator := service.NewInvoiceCalculator()

	uc := invoice.NewUsecase(
		invoiceRepo,
		vendorRepo,
		bankAccountRepo,
		calculator,
	)

	return ctx, uc, &controllers{
		ctrl:            ctrl,
		ctxProvider:     &ctxProvider,
		invoiceRepo:     invoiceRepo,
		vendorRepo:      vendorRepo,
		bankAccountRepo: bankAccountRepo,
	}
}
