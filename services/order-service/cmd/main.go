package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudcommerce/order-service/internal/application"
	"github.com/cloudcommerce/order-service/internal/infrastructure/postgres"
	"github.com/cloudcommerce/order-service/internal/transport"
	"github.com/cloudcommerce/shared-go/config"
	"github.com/cloudcommerce/shared-go/events"
	"github.com/cloudcommerce/shared-go/jwt"
	cclogging "github.com/cloudcommerce/shared-go/logging"
	"github.com/cloudcommerce/shared-go/nats"
	ccpostgres "github.com/cloudcommerce/shared-go/postgres"
)

func main() {
	// Load configuration
	cfg := config.Load("order-service")
	logger := cclogging.New("order-service", cfg.Environment)

	// Setup database connection
	ctx := context.Background()
	databaseURL := cfg.DatabaseURL
	if databaseURL == "" {
		databaseURL = fmt.Sprintf("postgres://postgres:postgres@localhost:5432/order_db?sslmode=disable")
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

	// Setup NATS subscriber for consuming events
	subscriber, err := nats.NewSubscriber(cfg.NATSURL, logger)
	if err != nil {
		logger.Error("failed to create NATS subscriber", "error", err)
		os.Exit(1)
	}
	defer subscriber.Close()

	logger.Info("NATS subscriber created")

	// Setup JWT manager
	jwtManager, err := jwt.NewManager(cfg.JWTSecret, cfg.JWTExpiry)
	if err != nil {
		logger.Error("failed to create JWT manager", "error", err)
		os.Exit(1)
	}

	// Wire up repositories
	orderRepo := postgres.NewOrderRepository(pool)
	itemRepo := postgres.NewOrderItemRepository(pool)
	statusHistoryRepo := postgres.NewOrderStatusHistoryRepository(pool)
	sagaEventRepo := postgres.NewSagaEventRepository(pool)

	// Wire up service
	orderService := application.NewOrderService(
		orderRepo,
		itemRepo,
		statusHistoryRepo,
		sagaEventRepo,
		publisher,
		logger,
	)

	// Setup HTTP handler
	handler := transport.NewOrderHandler(orderService)

	// Setup router
	router := transport.SetupRouter(handler, jwtManager, "order-service")

	// Start event subscribers
	go startEventSubscribers(ctx, subscriber, orderService, logger)

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
		logger.Info("order service starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down order service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("order service stopped")
}

// startEventSubscribers starts NATS event subscribers.
func startEventSubscribers(
	ctx context.Context,
	subscriber *nats.Subscriber,
	orderService *application.OrderService,
	logger *slog.Logger,
) {
	// Subscribe to payment success events
	err := subscriber.Subscribe(ctx, "cloudcommerce.payment.success", "order-payment-success", func(ctx context.Context, data []byte) error {
		var event events.PaymentSuccessEvent
		if err := json.Unmarshal(data, &event); err != nil {
			logger.Info("failed to unmarshal payment success event", "error", err)
			return nil
		}
		return orderService.HandlePaymentSuccess(ctx, event)
	})
	if err != nil {
		logger.Info("failed to subscribe to payment.success", "error", err)
	}

	// Subscribe to payment failed events
	err = subscriber.Subscribe(ctx, "cloudcommerce.payment.failed", "order-payment-failed", func(ctx context.Context, data []byte) error {
		var event events.PaymentFailedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			logger.Info("failed to unmarshal payment failed event", "error", err)
			return nil
		}
		return orderService.HandlePaymentFailed(ctx, event)
	})
	if err != nil {
		logger.Info("failed to subscribe to payment.failed", "error", err)
	}

	// Subscribe to inventory reserved events
	err = subscriber.Subscribe(ctx, "cloudcommerce.inventory.reserved", "order-inventory-reserved", func(ctx context.Context, data []byte) error {
		var event events.InventoryReservedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			logger.Info("failed to unmarshal inventory reserved event", "error", err)
			return nil
		}
		return orderService.HandleInventoryReserved(ctx, event)
	})
	if err != nil {
		logger.Info("failed to subscribe to inventory.reserved", "error", err)
	}

	// Subscribe to inventory reservation failed events
	err = subscriber.Subscribe(ctx, "cloudcommerce.inventory.reservation.failed", "order-inventory-reservation-failed", func(ctx context.Context, data []byte) error {
		var event events.InventoryReservationFailedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			logger.Info("failed to unmarshal inventory reservation failed event", "error", err)
			return nil
		}
		return orderService.HandleInventoryReservationFailed(ctx, event)
	})
	if err != nil {
		logger.Info("failed to subscribe to inventory.reservation.failed", "error", err)
	}

	logger.Info("event subscribers started")
}
