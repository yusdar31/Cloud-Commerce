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
	ccpostgres "github.com/cloudcommerce/shared-go/postgres"
	"github.com/cloudcommerce/user-service/internal/application"
	"github.com/cloudcommerce/user-service/internal/domain"
	"github.com/cloudcommerce/user-service/internal/infrastructure/postgres"
	"github.com/cloudcommerce/user-service/internal/infrastructure/security"
	"github.com/cloudcommerce/user-service/internal/transport"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load("identity-service")
	logger := logging.New("identity-service", cfg.Environment)

	// Setup database connection
	ctx := context.Background()
	databaseURL := cfg.DatabaseURL
	if databaseURL == "" {
		databaseURL = fmt.Sprintf("postgres://postgres:postgres@localhost:5432/identity_db?sslmode=disable")
	}

	pool, err := ccpostgres.NewPool(ctx, databaseURL, cfg.DBMaxConns(), cfg.DBMinConns())
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer ccpostgres.Close(pool)

	logger.Info("database connected", "url", maskURL(databaseURL))

	// Setup JWT manager
	jwtManager, err := ccjwt.NewManager(cfg.JWTSecret, cfg.JWTExpiry)
	if err != nil {
		logger.Error("failed to create JWT manager", "error", err)
		os.Exit(1)
	}

	// Wire up dependencies
	userRepo := postgres.NewUserRepository(pool)
	tokenRepo := postgres.NewRefreshTokenRepository(pool)
	hasher := security.NewBcryptHasher()

	authService := application.NewAuthService(userRepo, tokenRepo, hasher, jwtManager, logger)
	authHandler := transport.NewAuthHandler(authService)

	// Setup router
	router := transport.SetupRouter(authHandler, jwtManager, "identity-service")
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
		logger.Info("identity service starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down identity service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("identity service stopped")

	// Suppress unused import warnings for domain package
	_ = domain.ErrUserNotFound
	_ = (*pgxpool.Pool)(nil)
}

// maskURL hides credentials in a database URL for logging.
func maskURL(url string) string {
	// Simple masking: replace password portion with ***
	// In production, use proper URL parsing
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
