package transport

import (
	"github.com/cloudcommerce/shared-go/jwt"
	"github.com/cloudcommerce/shared-go/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes for the Order Service.
func SetupRouter(
	handler *OrderHandler,
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

	// Order routes
	orders := protected.Group("/orders")
	{
		orders.POST("", handler.Create)
		orders.GET("", handler.List)
		orders.GET("/:id", handler.GetByID)
	}

	// Internal routes (service-to-service)
	internal := v1.Group("/internal")
	{
		// Payment webhook callback
		internal.POST("/orders/:id/payment-callback", handler.HandlePaymentCallback)
	}

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "CloudCommerce Order Service",
			"version": "1.0.0",
		})
	})

	return r
}
