// Package main provides the entry point for the API server.
//
//	@title						スーパー支払い君.com API
//	@version					1.0
//	@description				請求書管理・支払い処理を行うREST APIサービス
//	@host						localhost:8080
//	@BasePath					/api
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Bearer token authentication
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/harusys/super-shiharai-kun/internal/config"
	"github.com/harusys/super-shiharai-kun/internal/controller"
	"github.com/harusys/super-shiharai-kun/internal/controller/middleware"
	"github.com/harusys/super-shiharai-kun/internal/domain/service"
	"github.com/harusys/super-shiharai-kun/internal/infrastructure/database"
	"github.com/harusys/super-shiharai-kun/internal/infrastructure/persistence"
	"github.com/harusys/super-shiharai-kun/internal/infrastructure/security"
	"github.com/harusys/super-shiharai-kun/internal/usecase/auth"
	"github.com/harusys/super-shiharai-kun/internal/usecase/invoice"
	"github.com/harusys/super-shiharai-kun/pkg/ctxutil"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := run(); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize database connection
	pool, err := database.NewPool(ctx, cfg.DatabaseURL())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	slog.Info("connected to database")

	// Initialize repositories
	userRepo := persistence.NewUserRepository(pool)
	vendorRepo := persistence.NewVendorRepository(pool)
	bankAccountRepo := persistence.NewVendorBankAccountRepository(pool)
	invoiceRepo := persistence.NewInvoiceRepository(pool)

	// Initialize services
	jwtService := security.NewJWTService(cfg.JWTSecret)
	calculator := service.NewInvoiceCalculator()

	// Initialize usecases
	authUsecase := auth.NewUsecase(userRepo, jwtService)
	invoiceUsecase := invoice.NewUsecase(invoiceRepo, vendorRepo, bankAccountRepo, calculator)

	// Setup gin router
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(ginLogger())
	r.Use(middleware.ErrorHandlerMiddleware())

	// Setup routes
	controller.SetupRoutes(r, &controller.RouterConfig{
		AuthUsecase:    authUsecase,
		InvoiceUsecase: invoiceUsecase,
		JWTService:     jwtService,
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	slog.Info("starting server", "addr", addr)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)

	go func() {
		if err := r.Run(addr); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-quit:
		slog.Info("shutting down server...")
		cancel()

		return nil
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	}
}

func ginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		start := ctxutil.Now(ctx)
		path := c.Request.URL.Path

		c.Next()

		slog.Info("request",
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"latency", ctxutil.Now(ctx).Sub(start),
			"client_ip", c.ClientIP(),
		)
	}
}
