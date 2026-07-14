package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/cloudcommerce/api-gateway/internal/gateway"
	"github.com/cloudcommerce/api-gateway/internal/middleware"
)

func main() {
	// Load configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-me"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Create gateway router
	configs := gateway.DefaultConfigs()
	router := gateway.NewRouter(configs)

	// Setup Gin
	r := gin.New()

	// Global middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())
	r.Use(gin.Recovery())

	// Rate limiting: 100 requests per minute per IP
	r.Use(middleware.RateLimit(100, time.Minute))

	// Health endpoints
	r.GET("/health", middleware.Health())
	r.GET("/healthz", middleware.Health())

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "api-gateway",
			"version": "0.2.0",
			"routes":  configs,
		})
	})

	// API routes - register each service route with wildcard matching
	for _, cfg := range configs {
		handler := router.HandleProxy(cfg)
		if cfg.Protected {
			// Protected routes require JWT
			r.Any(cfg.BasePath, middleware.JWTAuth(jwtSecret), handler)
			r.Any(cfg.BasePath+"/*anyPath", middleware.JWTAuth(jwtSecret), handler)
		} else {
			// Public routes
			r.Any(cfg.BasePath, handler)
			r.Any(cfg.BasePath+"/*anyPath", handler)
		}
	}

	// Start server
	addr := fmt.Sprintf(":%s", port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("API Gateway running on %s", addr)
		log.Printf("Environment: %s", appEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down API Gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("API Gateway stopped")
}
