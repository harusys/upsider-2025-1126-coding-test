package invoice_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/harusys/super-shiharai-kun/internal/controller/invoice"
	"github.com/harusys/super-shiharai-kun/internal/controller/middleware"
	"github.com/harusys/super-shiharai-kun/internal/domain"
	"github.com/harusys/super-shiharai-kun/internal/domain/entity"
	usecase "github.com/harusys/super-shiharai-kun/internal/usecase/invoice"
	"github.com/harusys/super-shiharai-kun/internal/usecase/invoice/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupRouter(handler *invoice.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()

	// Mock auth middleware to inject company_id
	r.Use(func(c *gin.Context) {
		c.Set(middleware.CompanyIDKey, int64(1))
		c.Next()
	})

	r.POST("/invoices", handler.Create)
	r.GET("/invoices", handler.List)
	r.GET("/invoices/:id", handler.GetByID)

	return r
}

func TestHandler_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       map[string]any
		prepare    func(m *mock.MockUsecase)
		wantStatus int
		wantError  string
	}{
		{
			name: "success",
			body: map[string]any{
				"vendor_id":              1,
				"vendor_bank_account_id": 1,
				"issue_date":             "2024-01-15",
				"payment_amount":         10000,
				"due_date":               "2024-02-15",
			},
			prepare: func(m *mock.MockUsecase) {
				issueDate, _ := time.Parse("2006-01-02", "2024-01-15")
				dueDate, _ := time.Parse("2006-01-02", "2024-02-15")

				m.EXPECT().
					Create(gomock.Any(), &usecase.CreateInput{
						CompanyID:           1,
						VendorID:            1,
						VendorBankAccountID: 1,
						IssueDate:           issueDate,
						PaymentAmount:       10000,
						DueDate:             dueDate,
					}).
					Return(&entity.Invoice{
						ID:                  1,
						CompanyID:           1,
						VendorID:            1,
						VendorBankAccountID: 1,
						IssueDate:           issueDate,
						PaymentAmount:       10000,
						Fee:                 400,
						FeeRate:             decimal.RequireFromString("0.04"),
						Tax:                 40,
						TaxRate:             decimal.RequireFromString("0.10"),
						TotalAmount:         10440,
						DueDate:             dueDate,
						Status:              entity.InvoiceStatusPending,
					}, nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "vendor not found",
			body: map[string]any{
				"vendor_id":              999,
				"vendor_bank_account_id": 1,
				"issue_date":             "2024-01-15",
				"payment_amount":         10000,
				"due_date":               "2024-02-15",
			},
			prepare: func(m *mock.MockUsecase) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, domain.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantError:  "vendor or bank account not found",
		},
		{
			name: "invalid request - missing vendor_id",
			body: map[string]any{
				"vendor_bank_account_id": 1,
				"issue_date":             "2024-01-15",
				"payment_amount":         10000,
				"due_date":               "2024-02-15",
			},
			prepare:    func(_ *mock.MockUsecase) {},
			wantStatus: http.StatusBadRequest,
			wantError:  "validation error",
		},
		{
			name: "invalid request - invalid date format",
			body: map[string]any{
				"vendor_id":              1,
				"vendor_bank_account_id": 1,
				"issue_date":             "2024/01/15",
				"payment_amount":         10000,
				"due_date":               "2024-02-15",
			},
			prepare:    func(_ *mock.MockUsecase) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock.NewMockUsecase(ctrl)
			tt.prepare(mockUsecase)

			handler := invoice.NewHandler(mockUsecase, validator.New())
			r := setupRouter(handler)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/invoices", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantError != "" {
				var resp map[string]any

				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Equal(t, tt.wantError, resp["error"])
			}
		})
	}
}

func TestHandler_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		query      string
		prepare    func(m *mock.MockUsecase)
		wantStatus int
		wantCount  int
	}{
		{
			name:  "success - list all",
			query: "",
			prepare: func(m *mock.MockUsecase) {
				m.EXPECT().
					List(gomock.Any(), &usecase.ListInput{
						CompanyID: 1,
					}).
					Return([]*entity.Invoice{
						{
							ID:        1,
							CompanyID: 1,
							FeeRate:   decimal.NewFromInt(0),
							TaxRate:   decimal.NewFromInt(0),
						},
						{
							ID:        2,
							CompanyID: 1,
							FeeRate:   decimal.NewFromInt(0),
							TaxRate:   decimal.NewFromInt(0),
						},
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name:  "success - empty list",
			query: "",
			prepare: func(m *mock.MockUsecase) {
				m.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return([]*entity.Invoice{}, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name:  "invalid date format",
			query: "?start_date=2024/01/01",
			prepare: func(_ *mock.MockUsecase) {
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock.NewMockUsecase(ctrl)
			tt.prepare(mockUsecase)

			handler := invoice.NewHandler(mockUsecase, validator.New())
			r := setupRouter(handler)

			req := httptest.NewRequest(http.MethodGet, "/invoices"+tt.query, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var resp []map[string]any

				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Len(t, resp, tt.wantCount)
			}
		})
	}
}

func TestHandler_GetByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		invoiceID  string
		prepare    func(m *mock.MockUsecase)
		wantStatus int
		wantError  string
	}{
		{
			name:      "success",
			invoiceID: "1",
			prepare: func(m *mock.MockUsecase) {
				m.EXPECT().
					GetByID(gomock.Any(), int64(1), int64(1)).
					Return(&entity.Invoice{
						ID:        1,
						CompanyID: 1,
						FeeRate:   decimal.NewFromInt(0),
						TaxRate:   decimal.NewFromInt(0),
					}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:      "not found",
			invoiceID: "999",
			prepare: func(m *mock.MockUsecase) {
				m.EXPECT().
					GetByID(gomock.Any(), int64(1), int64(999)).
					Return(nil, domain.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantError:  "invoice not found",
		},
		{
			name:       "invalid id",
			invoiceID:  "invalid",
			prepare:    func(_ *mock.MockUsecase) {},
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid invoice id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock.NewMockUsecase(ctrl)
			tt.prepare(mockUsecase)

			handler := invoice.NewHandler(mockUsecase, validator.New())
			r := setupRouter(handler)

			req := httptest.NewRequest(http.MethodGet, "/invoices/"+tt.invoiceID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantError != "" {
				var resp map[string]any

				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Equal(t, tt.wantError, resp["error"])
			}
		})
	}
}
