package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudcommerce/shared-go/config"
	ccjwt "github.com/cloudcommerce/shared-go/jwt"
	"github.com/cloudcommerce/shared-go/logging"
	"github.com/cloudcommerce/shared-go/middleware"
	"github.com/cloudcommerce/shared-go/nats"
	ccpostgres "github.com/cloudcommerce/shared-go/postgres"
	"github.com/cloudcommerce/store-service/internal/application"
	infraPostgres "github.com/cloudcommerce/store-service/internal/infrastructure/postgres"
	"github.com/cloudcommerce/store-service/internal/transport"
)

func main() {
	cfg := config.Load("store-service")
	logger := logging.New("store-service", cfg.Environment)

	// Setup database connection
	ctx := context.Background()
	databaseURL := cfg.DatabaseURL
	if databaseURL == "" {
		databaseURL = fmt.Sprintf("postgres://postgres:postgres@localhost:5432/tenant_db?sslmode=disable")
	}

	pool, err := ccpostgres.NewPool(ctx, databaseURL, cfg.DBMaxConns(), cfg.DBMinConns())
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer ccpostgres.Close(pool)

	logger.Info("database connected", "url", maskURL(databaseURL))

	// Setup NATS publisher
	var publisher *nats.Publisher
	if cfg.NATSURL != "" {
		publisher, err = nats.NewPublisher(cfg.NATSURL)
		if err != nil {
			logger.Warn("failed to connect to NATS, proceeding without publisher", "error", err)
		} else {
			defer publisher.Close()
			logger.Info("NATS connected", "url", cfg.NATSURL)
		}
	} else {
		logger.Warn("NATSURL not configured, proceeding without publisher")
	}

	// Setup JWT manager
	jwtManager, err := ccjwt.NewManager(cfg.JWTSecret, cfg.JWTExpiry)
	if err != nil {
		logger.Error("failed to create JWT manager", "error", err)
		os.Exit(1)
	}

	// Wire up dependencies
	tenantRepo := infraPostgres.NewTenantRepository(pool)
	storeService := application.NewStoreService(tenantRepo, publisher, logger)
	storeHandler := transport.NewStoreHandler(storeService)

	// Setup router
	router := transport.SetupRouter(storeHandler, jwtManager, "store-service")
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
		logger.Info("store service starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down store service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("store service stopped")
}

// maskURL hides credentials in a database URL for logging.
func maskURL(url string) string {
	return fmt.Sprintf("postgres://***@%s", extractHost(url))
}

func extractHost(url string) string {
	for i := 0; i < len(url); i++ {
		if url[i] == '@' {
			return url[i+1:]
		}
	}
	return "localhost"
}
