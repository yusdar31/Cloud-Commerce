package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudcommerce/payment-service/internal/application"
	"github.com/cloudcommerce/payment-service/internal/domain"
	"github.com/cloudcommerce/payment-service/internal/infrastructure/postgres"
	"github.com/cloudcommerce/payment-service/internal/transport"
	"github.com/cloudcommerce/shared-go/config"
	"github.com/cloudcommerce/shared-go/idempotency"
	"github.com/cloudcommerce/shared-go/jwt"
	cclogging "github.com/cloudcommerce/shared-go/logging"
	"github.com/cloudcommerce/shared-go/nats"
	ccpostgres "github.com/cloudcommerce/shared-go/postgres"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configuration
	cfg := config.Load("payment-service")
	logger := cclogging.New("payment-service", cfg.Environment)

	// Setup database connection
	ctx := context.Background()
	databaseURL := cfg.DatabaseURL
	if databaseURL == "" {
		databaseURL = fmt.Sprintf("postgres://postgres:postgres@localhost:5432/payment_db?sslmode=disable")
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

	// Setup Redis for idempotency
	var idempotencyStore idempotency.Store
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: "",
		DB:       0,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Warn("redis not available, idempotency middleware disabled", "error", err)
	} else {
		idempotencyStore = idempotency.NewRedisStore(redisClient)
		logger.Info("redis connected for idempotency")
	}

	// Wire up repositories
	paymentRepo := postgres.NewPaymentRepository(pool)
	webhookRepo := postgres.NewWebhookEventRepository(pool)

	// TODO: Wire up payment gateway (Midtrans/Xendit)
	// For now, use nil - will be implemented when integrating real gateway
	var gateway domain.PaymentGateway = nil

	// Wire up service
	paymentService := application.NewPaymentService(
		paymentRepo,
		nil, // transactionRepo - optional for MVP
		webhookRepo,
		gateway,
		publisher,
		logger,
	)

	// Setup HTTP handler
	handler := transport.NewPaymentHandler(paymentService)

	// Setup router
	router := transport.SetupRouter(handler, jwtManager, "payment-service", idempotencyStore, logger)

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
		logger.Info("payment service starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down payment service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("payment service stopped")
}
