package invoice

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/harusys/super-shiharai-kun/internal/controller/middleware"
	"github.com/harusys/super-shiharai-kun/internal/domain"
	"github.com/harusys/super-shiharai-kun/internal/usecase/invoice"
)

// Handler handles invoice endpoints.
type Handler struct {
	usecase   invoice.Usecase
	validator *validator.Validate
}

// NewHandler creates a new Handler.
func NewHandler(usecase invoice.Usecase, validator *validator.Validate) *Handler {
	return &Handler{
		usecase:   usecase,
		validator: validator,
	}
}

// Create handles invoice creation.
//
//	@Summary		請求書作成
//	@Description	新しい請求書データを作成します。手数料・消費税は自動計算されます。
//	@Tags			invoices
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateRequest	true	"請求書作成リクエスト"
//	@Success		201		{object}	Response
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/invoices [post]
func (h *Handler) Create(c *gin.Context) {
	companyID := middleware.GetCompanyID(c)

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid request body"))

		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewValidationErrorResponse(formatValidationErrors(err)))

		return
	}

	issueDate, err := time.Parse("2006-01-02", req.IssueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid issue_date format"))

		return
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid due_date format"))

		return
	}

	input := &invoice.CreateInput{
		CompanyID:           companyID,
		VendorID:            req.VendorID,
		VendorBankAccountID: req.VendorBankAccountID,
		IssueDate:           issueDate,
		PaymentAmount:       req.PaymentAmount,
		DueDate:             dueDate,
	}

	inv, err := h.usecase.Create(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, NewErrorResponse("vendor or bank account not found"))

			return
		}

		c.JSON(http.StatusInternalServerError, NewErrorResponse("internal server error"))

		return
	}

	c.JSON(http.StatusCreated, ToResponse(inv))
}

// List handles listing invoices.
//
//	@Summary		請求書一覧取得
//	@Description	指定期間内に支払いが発生する請求書データの一覧を取得します
//	@Tags			invoices
//	@Accept			json
//	@Produce		json
//	@Param			start_date	query		string	false	"開始日 (YYYY-MM-DD)"
//	@Param			end_date	query		string	false	"終了日 (YYYY-MM-DD)"
//	@Success		200			{array}		Response
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/invoices [get]
func (h *Handler) List(c *gin.Context) {
	companyID := middleware.GetCompanyID(c)

	input := &invoice.ListInput{
		CompanyID: companyID,
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewErrorResponse("invalid start_date format"))

			return
		}

		input.StartDate = &startDate
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewErrorResponse("invalid end_date format"))

			return
		}

		input.EndDate = &endDate
	}

	invoices, err := h.usecase.List(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("internal server error"))

		return
	}

	c.JSON(http.StatusOK, ToResponses(invoices))
}

// GetByID handles getting an invoice by ID.
//
//	@Summary		請求書詳細取得
//	@Description	指定IDの請求書データを取得します
//	@Tags			invoices
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"請求書ID"
//	@Success		200	{object}	Response
//	@Failure		400	{object}	ErrorResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/invoices/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	companyID := middleware.GetCompanyID(c)

	invoiceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid invoice id"))

		return
	}

	inv, err := h.usecase.GetByID(c.Request.Context(), companyID, invoiceID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, NewErrorResponse("invoice not found"))

			return
		}

		c.JSON(http.StatusInternalServerError, NewErrorResponse("internal server error"))

		return
	}

	c.JSON(http.StatusOK, ToResponse(inv))
}

func formatValidationErrors(err error) map[string]string {
	details := make(map[string]string)

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, e := range validationErrs {
			details[e.Field()] = e.Tag()
		}
	}

	return details
}
