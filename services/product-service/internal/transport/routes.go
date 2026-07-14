package transport

import (
	"github.com/cloudcommerce/shared-go/jwt"
	"github.com/cloudcommerce/shared-go/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes for the Catalog Service.
func SetupRouter(
	handler *ProductHandler,
	jwtManager *jwt.Manager,
	serviceName string,
) *gin.Engine {
	r := gin.New()

	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	// Health endpoints
	r.GET("/health", middleware.Health(serviceName))
	r.GET("/healthz", middleware.Health(serviceName))

	// API v1 routes
	v1 := r.Group("/api/v1")

	// All catalog routes require authentication and tenant isolation
	protected := v1.Group("")
	protected.Use(middleware.JWTAuth(jwtManager))
	protected.Use(middleware.TenantIsolation())

	// Product routes
	products := protected.Group("/products")
	{
		products.GET("", handler.List)
		products.POST("", handler.Create)
		products.GET("/:id", handler.GetByID)
		products.PUT("/:id", handler.Update)
		products.DELETE("/:id", handler.Delete)
		products.POST("/:id/publish", handler.Publish)
		products.POST("/:id/archive", handler.Archive)
	}

	// Category routes
	categories := protected.Group("/categories")
	{
		categories.GET("", handler.ListCategories)
		categories.POST("", handler.CreateCategory)
		categories.GET("/:id", handler.GetCategory)
		categories.DELETE("/:id", handler.DeleteCategory)
	}

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "CloudCommerce Catalog Service",
			"version": "1.0.0",
		})
	})

	return r
}
