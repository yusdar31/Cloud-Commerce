package transport

import (
	"time"

	"github.com/cloudcommerce/user-service/internal/domain"
	"github.com/go-playground/validator/v10"
)

// RegisterRequest is the DTO for user registration.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
	Phone    string `json:"phone" validate:"omitempty,min=10,max=20"`
	Role     string `json:"role" validate:"omitempty,oneof=seller buyer admin"`
}

// LoginRequest is the DTO for user login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshRequest is the DTO for token refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// UserResponse is the DTO for user data in API responses.
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone,omitempty"`
	Role      string    `json:"role"`
	TenantID  string    `json:"tenant_id,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// TokenResponse is the DTO for authentication tokens.
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// AuthResponse combines user and token data.
type AuthResponse struct {
	Token TokenResponse `json:"tokens"`
	User  UserResponse  `json:"user"`
}

// Validate validates a struct using validator tags.
func Validate(req interface{}) error {
	v := validator.New()
	return v.Struct(req)
}

// ToUserResponse converts a domain User to a UserResponse DTO.
func ToUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		Phone:     u.Phone,
		Role:      string(u.Role),
		TenantID:  u.TenantID,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}

// ToTokenResponse converts domain AuthTokens to a TokenResponse DTO.
func ToTokenResponse(t *domain.AuthTokens) TokenResponse {
	return TokenResponse{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		ExpiresAt:    t.ExpiresAt,
		TokenType:    "Bearer",
	}
}
