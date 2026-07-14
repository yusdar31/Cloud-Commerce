package transport

import (
	"errors"
	"net/http"

	"github.com/cloudcommerce/shared-go/response"
	"github.com/cloudcommerce/user-service/internal/application"
	"github.com/cloudcommerce/user-service/internal/domain"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication endpoints.
type AuthHandler struct {
	authService *application.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *application.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles POST /api/v1/auth/register.
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := Validate(req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	role := domain.Role(req.Role)
	if role == "" {
		role = domain.RoleBuyer
	}

	user, err := h.authService.Register(c.Request.Context(), application.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Phone:    req.Phone,
		Role:     role,
	})
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			response.Conflict(c, "Email already registered")
			return
		}
		response.InternalError(c, "Failed to register user")
		return
	}

	response.Created(c, ToUserResponse(user))
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := Validate(req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	tokens, user, err := h.authService.Login(c.Request.Context(), application.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.Unauthorized(c, "Invalid email or password")
			return
		}
		response.InternalError(c, "Failed to login")
		return
	}

	response.OK(c, AuthResponse{
		Token: ToTokenResponse(tokens),
		User:  ToUserResponse(user),
	})
}

// Refresh handles POST /api/v1/auth/refresh.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := Validate(req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	tokens, user, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) || errors.Is(err, domain.ErrTokenExpired) {
			response.Unauthorized(c, "Invalid or expired refresh token")
			return
		}
		response.InternalError(c, "Failed to refresh token")
		return
	}

	response.OK(c, AuthResponse{
		Token: ToTokenResponse(tokens),
		User:  ToUserResponse(user),
	})
}

// Logout handles POST /api/v1/auth/logout.
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "Authentication required")
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userID); err != nil {
		response.InternalError(c, "Failed to logout")
		return
	}

	c.Status(http.StatusNoContent)
}

// GetProfile handles GET /api/v1/users/me.
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "Authentication required")
		return
	}

	user, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalError(c, "Failed to get profile")
		return
	}

	response.OK(c, ToUserResponse(user))
}
