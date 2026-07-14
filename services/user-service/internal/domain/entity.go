package domain

import (
	"errors"
	"time"
)

// Common domain errors.
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
)

// Role represents a user role in the system.
type Role string

const (
	RoleSeller Role = "seller"
	RoleBuyer  Role = "buyer"
	RoleAdmin  Role = "admin"
)

// User is the aggregate root of the Identity context.
type User struct {
	ID           string
	Email        string
	PasswordHash string
	FullName     string
	Phone        string
	Role         Role
	TenantID     string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// RefreshToken holds a refresh token record.
type RefreshToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}

// AuthTokens represents the JWT access token and refresh token pair.
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
