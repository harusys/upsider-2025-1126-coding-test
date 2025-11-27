package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/harusys/super-shiharai-kun/internal/controller"
	"github.com/harusys/super-shiharai-kun/internal/domain/service"
	"github.com/harusys/super-shiharai-kun/internal/infrastructure/database"
	"github.com/harusys/super-shiharai-kun/internal/infrastructure/persistence"
	"github.com/harusys/super-shiharai-kun/internal/infrastructure/security"
	"github.com/harusys/super-shiharai-kun/internal/usecase/auth"
	"github.com/harusys/super-shiharai-kun/internal/usecase/invoice"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	pool       *pgxpool.Pool
	router     *gin.Engine
	jwtService *security.JWTService
}

func (s *APITestSuite) SetupSuite() {
	// Skip if no database URL is set
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		s.T().Skip("TEST_DATABASE_URL not set, skipping integration tests")
	}

	ctx := context.Background()

	pool, err := database.NewPool(ctx, dbURL)
	s.Require().NoError(err)

	s.pool = pool

	// Initialize repositories
	userRepo := persistence.NewUserRepository(pool)
	vendorRepo := persistence.NewVendorRepository(pool)
	bankAccountRepo := persistence.NewVendorBankAccountRepository(pool)
	invoiceRepo := persistence.NewInvoiceRepository(pool)

	// Initialize services
	s.jwtService = security.NewJWTService("test-secret-key")
	calculator := service.NewInvoiceCalculator()

	// Initialize usecases
	authUsecase := auth.NewUsecase(userRepo, s.jwtService)
	invoiceUsecase := invoice.NewUsecase(
		invoiceRepo,
		vendorRepo,
		bankAccountRepo,
		calculator,
	)

	// Setup router
	gin.SetMode(gin.TestMode)

	s.router = gin.New()
	controller.SetupRoutes(s.router, &controller.RouterConfig{
		AuthUsecase:    authUsecase,
		InvoiceUsecase: invoiceUsecase,
		JWTService:     s.jwtService,
	})
}

func (s *APITestSuite) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()
	}
}

func (s *APITestSuite) SetupTest() {
	// Clean up test data before each test
	if s.pool != nil {
		ctx := context.Background()
		_, _ = s.pool.Exec(ctx, "DELETE FROM invoices")
		_, _ = s.pool.Exec(ctx, "DELETE FROM vendor_bank_accounts")
		_, _ = s.pool.Exec(ctx, "DELETE FROM vendors")
		_, _ = s.pool.Exec(ctx, "DELETE FROM users")
	}
}

func (s *APITestSuite) TestHealthCheck() {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	w := httptest.NewRecorder()

	// Add health check route for this test
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

func (s *APITestSuite) TestAuthFlow() {
	// 1. Register a new user
	registerBody := map[string]any{
		"company_id": 1,
		"name":       "Test User",
		"email":      "test@example.com",
		"password":   "password123",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)

	var registerResp map[string]any

	err := json.Unmarshal(w.Body.Bytes(), &registerResp)
	s.Require().NoError(err)
	s.NotEmpty(registerResp["access_token"])
	s.NotEmpty(registerResp["refresh_token"])

	// 2. Login with same credentials
	loginBody := map[string]any{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var loginResp map[string]any

	err = json.Unmarshal(w.Body.Bytes(), &loginResp)
	s.Require().NoError(err)
	s.NotEmpty(loginResp["access_token"])

	// 3. Refresh token
	refreshBody := map[string]any{
		"refresh_token": loginResp["refresh_token"],
	}
	body, _ = json.Marshal(refreshBody)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

func (s *APITestSuite) TestProtectedRouteWithoutToken() {
	req := httptest.NewRequest(http.MethodGet, "/api/invoices", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *APITestSuite) TestProtectedRouteWithInvalidToken() {
	req := httptest.NewRequest(http.MethodGet, "/api/invoices", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *APITestSuite) TestInvoiceFlow() {
	// 1. Register and get token
	registerBody := map[string]any{
		"company_id": 1,
		"name":       "Test User",
		"email":      "invoice-test@example.com",
		"password":   "password123",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)
	s.Require().Equal(http.StatusCreated, w.Code)

	var authResp map[string]any

	_ = json.Unmarshal(w.Body.Bytes(), &authResp)

	accessToken := authResp["access_token"].(string)

	// 2. Create vendor and bank account (direct DB insert for test setup)
	ctx := context.Background()

	var vendorID, bankAccountID int64

	err := s.pool.QueryRow(ctx, `
		INSERT INTO vendors (company_id, name, representative_name, phone_number, zip_code, address)
		VALUES (1, 'Test Vendor', 'Rep Name', '03-1234-5678', '100-0001', 'Tokyo')
		RETURNING id
	`).Scan(&vendorID)
	s.Require().NoError(err)

	err = s.pool.QueryRow(ctx, `
		INSERT INTO vendor_bank_accounts (vendor_id, bank_name, branch_name, account_number, account_holder_name)
		VALUES ($1, 'Test Bank', 'Test Branch', '1234567', 'Test Holder')
		RETURNING id
	`, vendorID).Scan(&bankAccountID)
	s.Require().NoError(err)

	// 3. Create invoice
	invoiceBody := map[string]any{
		"vendor_id":              vendorID,
		"vendor_bank_account_id": bankAccountID,
		"issue_date":             "2024-01-15",
		"payment_amount":         10000,
		"due_date":               "2024-02-15",
	}
	body, _ = json.Marshal(invoiceBody)
	req = httptest.NewRequest(http.MethodPost, "/api/invoices", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)

	var invoiceResp map[string]any

	err = json.Unmarshal(w.Body.Bytes(), &invoiceResp)
	s.Require().NoError(err)
	s.InDelta(10000, invoiceResp["payment_amount"].(float64), 0.01)
	s.InDelta(400, invoiceResp["fee"].(float64), 0.01)
	s.InDelta(40, invoiceResp["tax"].(float64), 0.01)
	s.InDelta(10440, invoiceResp["total_amount"].(float64), 0.01)
	s.Equal("pending", invoiceResp["status"])

	invoiceID := int64(invoiceResp["id"].(float64))

	// 4. List invoices
	req = httptest.NewRequest(http.MethodGet, "/api/invoices", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var listResp []map[string]any

	err = json.Unmarshal(w.Body.Bytes(), &listResp)
	s.Require().NoError(err)
	s.Len(listResp, 1)

	// 5. Get invoice by ID
	req = httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/api/invoices/%d", invoiceID),
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var getResp map[string]any

	err = json.Unmarshal(w.Body.Bytes(), &getResp)
	s.Require().NoError(err)
	s.InDelta(invoiceID, getResp["id"].(float64), 0.01)
}

func (s *APITestSuite) TestInvoiceListWithDateFilter() {
	// 1. Register and get token
	registerBody := map[string]any{
		"company_id": 1,
		"name":       "Test User",
		"email":      "filter-test@example.com",
		"password":   "password123",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)
	s.Require().Equal(http.StatusCreated, w.Code)

	var authResp map[string]any

	_ = json.Unmarshal(w.Body.Bytes(), &authResp)
	accessToken := authResp["access_token"].(string)

	// 2. List with date filter (should return empty)
	req = httptest.NewRequest(
		http.MethodGet,
		"/api/invoices?start_date=2024-01-01&end_date=2024-01-31",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var listResp []map[string]any

	err := json.Unmarshal(w.Body.Bytes(), &listResp)
	s.Require().NoError(err)
	s.Empty(listResp)
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
