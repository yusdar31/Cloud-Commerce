package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	ccjwt "github.com/cloudcommerce/shared-go/jwt"
)

// CORS returns a gin middleware that handles Cross-Origin Resource Sharing.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-ID, X-Correlation-ID, Idempotency-Key")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// RequestID injects a unique request ID into the context.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Correlation-ID")
		if id == "" {
			id = uuid.NewString()
		}
		c.Set("request_id", id)
		c.Writer.Header().Set("X-Correlation-ID", id)
		c.Next()
	}
}

// Logger returns a structured logging middleware.
func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		attrs := []any{
			"method", method,
			"path", path,
			"status", status,
			"latency_ms", latency.Milliseconds(),
			"request_id", c.GetString("request_id"),
			"client_ip", c.ClientIP(),
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, "errors", c.Errors.String())
		}

		if status >= 500 {
			logger.Error("request completed", attrs...)
		} else if status >= 400 {
			logger.Warn("request completed", attrs...)
		} else {
			logger.Info("request completed", attrs...)
		}
	}
}

// JWTAuth returns a middleware that validates JWT tokens from the Authorization header.
func JWTAuth(manager *ccjwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"type":   "https://api.cloudcommerce.com/errors/Unauthorized",
				"title":  "Unauthorized",
				"status": 401,
				"detail": "Authorization header is required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"type":   "https://api.cloudcommerce.com/errors/Unauthorized",
				"title":  "Unauthorized",
				"status": 401,
				"detail": "Authorization header must be Bearer <token>",
			})
			return
		}

		claims, err := manager.Validate(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"type":   "https://api.cloudcommerce.com/errors/Unauthorized",
				"title":  "Unauthorized",
				"status": 401,
				"detail": "Invalid or expired token",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RequireRole returns a middleware that checks if the user has one of the allowed roles.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")
		for _, r := range roles {
			if userRole == r {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"type":   "https://api.cloudcommerce.com/errors/Forbidden",
			"title":  "Forbidden",
			"status": 403,
			"detail": "You do not have permission to perform this action",
		})
	}
}

// TenantIsolation ensures tenant_id is present in the context.
func TenantIsolation() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		if tenantID == "" {
			// Try header fallback for internal calls
			tenantID = c.GetHeader("X-Tenant-ID")
		}
		if tenantID == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"type":   "https://api.cloudcommerce.com/errors/Forbidden",
				"title":  "Forbidden",
				"status": 403,
				"detail": "Tenant context is required",
			})
			return
		}
		c.Set("tenant_id", tenantID)
		c.Next()
	}
}

// Health returns a simple health check handler.
func Health(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": serviceName,
		})
	}
}
