package transport

import (
	"github.com/cloudcommerce/shared-go/jwt"
	"github.com/cloudcommerce/shared-go/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes for the Payment Service.
func SetupRouter(
	handler *PaymentHandler,
	jwtManager *jwt.Manager,
	serviceName string,
) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	// Health endpoints
	r.GET("/health", middleware.Health(serviceName))
	r.GET("/healthz", middleware.Health(serviceName))

	// API v1 routes
	v1 := r.Group("/api/v1")

	// Protected routes (JWT + Tenant isolation)
	protected := v1.Group("")
	protected.Use(middleware.JWTAuth(jwtManager))
	protected.Use(middleware.TenantIsolation())

	// Payment routes
	payments := protected.Group("/payments")
	{
		payments.POST("", handler.Create)
		payments.GET("", handler.List)
		payments.GET("/:id", handler.GetByID)
		payments.GET("/order/:order_id", handler.GetByOrderID)
	}

	// Public webhook endpoint (no auth - verified by signature)
	webhooks := r.Group("/webhooks")
	{
		webhooks.POST("/payments/:provider", handler.HandleWebhook)
	}

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "CloudCommerce Payment Service",
			"version": "1.0.0",
		})
	})

	return r
}
