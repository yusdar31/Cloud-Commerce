package idempotency

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Config holds idempotency middleware configuration.
type Config struct {
	Store    Store
	Logger   *slog.Logger
	KeyFunc  func(c *gin.Context) string
	TTL      time.Duration
	SkipFunc func(c *gin.Context) bool
}

// DefaultConfig returns a default idempotency configuration.
func DefaultConfig(store Store, logger *slog.Logger) *Config {
	return &Config{
		Store:  store,
		Logger: logger,
		TTL:    24 * time.Hour,
		KeyFunc: func(c *gin.Context) string {
			body, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			return GenerateKey(c.Request.Method, c.Request.URL.Path, body)
		},
	}
}

// Middleware returns a Gin middleware that enforces idempotency.
func Middleware(cfg *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if this request should skip idempotency
		if cfg.SkipFunc != nil && cfg.SkipFunc(c) {
			c.Next()
			return
		}

		// Generate idempotency key
		key := cfg.KeyFunc(c)
		if key == "" {
			c.Next()
			return
		}

		ctx := c.Request.Context()

		// Check if key already exists
		exists, err := cfg.Store.Exists(ctx, key)
		if err != nil {
			cfg.Logger.Error("failed to check idempotency key",
				"key", key,
				"error", err,
			)
			// Don't block on store errors
			c.Next()
			return
		}

		if exists {
			// Return 409 Conflict for duplicate requests
			c.JSON(http.StatusConflict, gin.H{
				"type":   "https://api.cloudcommerce.com/errors/Conflict",
				"title":  "Idempotency Conflict",
				"status": 409,
				"detail": "This request has already been processed",
			})
			c.Abort()
			return
		}

		// Store the key before processing
		if err := cfg.Store.Set(ctx, key, "pending", cfg.TTL); err != nil {
			cfg.Logger.Error("failed to store idempotency key",
				"key", key,
				"error", err,
			)
			// Don't block on store errors
			c.Next()
			return
		}

		// Process the request
		c.Next()

		// Mark as completed if successful (2xx)
		status := c.Writer.Status()
		if status >= 200 && status < 300 {
			cfg.Store.Set(ctx, key, "completed", cfg.TTL)
		} else {
			// Remove key on failure to allow retry
			cfg.Store.Set(ctx, key, "failed", 5*time.Minute)
		}
	}
}

// WebhookMiddleware returns a specialized middleware for webhook endpoints.
// It uses the provider's event ID for idempotency.
func WebhookMiddleware(store Store, logger *slog.Logger, extractEventID func(c *gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventID := extractEventID(c)
		if eventID == "" {
			c.Next()
			return
		}

		key := "webhook:" + eventID
		ctx := c.Request.Context()

		// Check if already processed
		exists, err := store.Exists(ctx, key)
		if err != nil {
			logger.Error("failed to check webhook idempotency",
				"event_id", eventID,
				"error", err,
			)
			c.Next()
			return
		}

		if exists {
			// Webhook already processed - return success (idempotent)
			c.JSON(http.StatusOK, gin.H{
				"message": "webhook already processed",
			})
			c.Abort()
			return
		}

		// Mark as processing
		store.Set(ctx, key, "processing", 24*time.Hour)

		// Process the webhook
		c.Next()

		// Mark as completed
		status := c.Writer.Status()
		if status >= 200 && status < 300 {
			store.Set(ctx, key, "completed", 24*time.Hour)
		} else {
			// Remove key on failure to allow retry
			store.Set(ctx, key, "failed", 5*time.Minute)
		}
	}
}
