package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/harusys/super-shiharai-kun/docs/swagger"
	authctrl "github.com/harusys/super-shiharai-kun/internal/controller/auth"
	invoicectrl "github.com/harusys/super-shiharai-kun/internal/controller/invoice"
	"github.com/harusys/super-shiharai-kun/internal/controller/middleware"
	"github.com/harusys/super-shiharai-kun/internal/infrastructure/security"
	"github.com/harusys/super-shiharai-kun/internal/usecase/auth"
	"github.com/harusys/super-shiharai-kun/internal/usecase/invoice"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RouterConfig holds dependencies for setting up routes.
type RouterConfig struct {
	AuthUsecase    auth.Usecase
	InvoiceUsecase invoice.Usecase
	JWTService     *security.JWTService
}

// SetupRoutes configures all API routes.
func SetupRoutes(r *gin.Engine, config *RouterConfig) {
	validate := validator.New()

	authHandler := authctrl.NewHandler(config.AuthUsecase, validate)
	invoiceHandler := invoicectrl.NewHandler(config.InvoiceUsecase, validate)

	api := r.Group("/api")

	// Auth routes (public)
	authGroup := api.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/refresh", authHandler.RefreshToken)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(config.JWTService))

	// Invoice routes
	invoiceGroup := protected.Group("/invoices")
	invoiceGroup.POST("", invoiceHandler.Create)
	invoiceGroup.GET("", invoiceHandler.List)
	invoiceGroup.GET("/:id", invoiceHandler.GetByID)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
