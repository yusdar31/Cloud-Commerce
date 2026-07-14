package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudcommerce/product-service/internal/application"
	"github.com/cloudcommerce/product-service/internal/infrastructure/postgres"
	"github.com/cloudcommerce/product-service/internal/transport"
	"github.com/cloudcommerce/shared-go/config"
	ccjwt "github.com/cloudcommerce/shared-go/jwt"
	"github.com/cloudcommerce/shared-go/logging"
	"github.com/cloudcommerce/shared-go/middleware"
	ccpostgres "github.com/cloudcommerce/shared-go/postgres"
)

func main() {
	cfg := config.Load("catalog-service")
	logger := logging.New("catalog-service", cfg.Environment)

	// Setup database connection
	ctx := context.Background()
	databaseURL := cfg.DatabaseURL
	if databaseURL == "" {
		databaseURL = fmt.Sprintf("postgres://postgres:postgres@localhost:5432/catalog_db?sslmode=disable")
	}

	pool, err := ccpostgres.NewPool(ctx, databaseURL, cfg.DBMaxConns(), cfg.DBMinConns())
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer ccpostgres.Close(pool)

	logger.Info("database connected")

	// Setup JWT manager (for validating tokens from the Identity Service)
	jwtManager, err := ccjwt.NewManager(cfg.JWTSecret, cfg.JWTExpiry)
	if err != nil {
		logger.Error("failed to create JWT manager", "error", err)
		os.Exit(1)
	}

	// Wire up dependencies
	productRepo := postgres.NewProductRepository(pool)
	categoryRepo := postgres.NewCategoryRepository(pool)

	productSvc := application.NewProductService(productRepo, logger)
	categorySvc := application.NewCategoryService(categoryRepo, logger)

	handler := transport.NewProductHandler(productSvc, categorySvc)

	// Setup router
	router := transport.SetupRouter(handler, jwtManager, "catalog-service")
	router.Use(middleware.Logger(logger))

	// Start server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		logger.Info("catalog service starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down catalog service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("catalog service stopped")
}
