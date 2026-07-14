package transport

import (
	"github.com/cloudcommerce/shared-go/jwt"
	"github.com/cloudcommerce/shared-go/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SetupRouter configures all routes for the Identity Service.
func SetupRouter(
	authHandler *AuthHandler,
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

	// Public auth routes (no JWT required)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}

	// Protected routes (JWT required)
	protected := v1.Group("")
	protected.Use(middleware.JWTAuth(jwtManager))
	{
		protected.POST("/auth/logout", authHandler.Logout)
		protected.GET("/users/me", authHandler.GetProfile)
	}

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "CloudCommerce Identity Service",
			"version": "1.0.0",
		})
	})

	return r
}

// GenerateRequestID is a helper for creating unique request IDs.
func GenerateRequestID() string {
	return uuid.NewString()
}
