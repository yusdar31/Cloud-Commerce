package transport

import (
	"github.com/cloudcommerce/shared-go/jwt"
	"github.com/cloudcommerce/shared-go/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures all HTTP routes for the Store Service.
func SetupRouter(
	handler *StoreHandler,
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

	// Public routes (no authentication required)
	public := v1.Group("/stores")
	{
		public.GET("/:id", handler.GetByID)
		public.GET("/slug/:slug", handler.GetBySlug)
	}

	// Protected routes (JWT required)
	protected := v1.Group("")
	protected.Use(middleware.JWTAuth(jwtManager))
	{
		protected.POST("/stores", handler.Create)
		protected.GET("/stores/me", handler.GetMe)
		protected.PUT("/stores/me", handler.Update)
	}

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "CloudCommerce Store Service",
			"version": "1.0.0",
		})
	})

	return r
}
