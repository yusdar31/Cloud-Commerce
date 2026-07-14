package transport

import (
	"github.com/cloudcommerce/shared-go/jwt"
	"github.com/cloudcommerce/shared-go/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes for the Inventory Service.
func SetupRouter(
	handler *InventoryHandler,
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

	// All inventory routes require authentication and tenant isolation
	protected := v1.Group("")
	protected.Use(middleware.JWTAuth(jwtManager))
	protected.Use(middleware.TenantIsolation())

	// Inventory routes
	inventory := protected.Group("/inventory")
	{
		inventory.POST("", handler.Create)
		inventory.GET("", handler.List)
		inventory.GET("/:id", handler.GetByID)
		inventory.GET("/low-stock", handler.GetLowStock)
		inventory.POST("/reserve", handler.Reserve)
		inventory.POST("/release", handler.Release)
	}

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "CloudCommerce Inventory Service",
			"version": "1.0.0",
		})
	})

	return r
}
