package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudcommerce/inventory-service/internal/application"
	"github.com/cloudcommerce/inventory-service/internal/infrastructure/postgres"
	"github.com/cloudcommerce/inventory-service/internal/transport"
	"github.com/cloudcommerce/shared-go/config"
	"github.com/cloudcommerce/shared-go/jwt"
	cclogging "github.com/cloudcommerce/shared-go/logging"
	"github.com/cloudcommerce/shared-go/nats"
	ccpostgres "github.com/cloudcommerce/shared-go/postgres"
)

func main() {
	// Load configuration
	cfg := config.Load("inventory-service")
	logger := cclogging.New("inventory-service", cfg.Environment)

	// Setup database connection
	ctx := context.Background()
	databaseURL := cfg.DatabaseURL
	if databaseURL == "" {
		databaseURL = fmt.Sprintf("postgres://postgres:postgres@localhost:5432/inventory_db?sslmode=disable")
	}

	pool, err := ccpostgres.NewPool(ctx, databaseURL, cfg.DBMaxConns(), cfg.DBMinConns())
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer ccpostgres.Close(pool)

	logger.Info("database connected")

	// Setup NATS publisher
	publisher, err := nats.NewPublisher(cfg.NATSURL)
	if err != nil {
		logger.Error("failed to connect to NATS", "error", err)
		os.Exit(1)
	}
	defer publisher.Close()

	logger.Info("NATS connected")

	// Setup JWT manager
	jwtManager, err := jwt.NewManager(cfg.JWTSecret, cfg.JWTExpiry)
	if err != nil {
		logger.Error("failed to create JWT manager", "error", err)
		os.Exit(1)
	}

	// Wire up repositories
	inventoryRepo := postgres.NewInventoryRepository(pool)
	reservationRepo := postgres.NewReservationRepository(pool)
	stockMovementRepo := postgres.NewStockMovementRepository(pool)

	// Wire up service
	inventoryService := application.NewInventoryService(
		inventoryRepo,
		reservationRepo,
		stockMovementRepo,
		publisher,
		logger,
	)

	// Setup HTTP handler
	handler := transport.NewInventoryHandler(inventoryService)

	// Setup router
	router := transport.SetupRouter(handler, jwtManager, "inventory-service")

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
		logger.Info("inventory service starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down inventory service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("inventory service stopped")
}
