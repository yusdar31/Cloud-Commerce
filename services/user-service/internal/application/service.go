package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/cloudcommerce/shared-go/jwt"
	"github.com/cloudcommerce/user-service/internal/domain"
	"github.com/google/uuid"
)

// AuthService implements authentication use cases.
type AuthService struct {
	userRepo   domain.UserRepository
	tokenRepo  domain.RefreshTokenRepository
	hasher     domain.PasswordHasher
	jwtManager *jwt.Manager
	logger     *slog.Logger
}

// NewAuthService creates a new AuthService.
func NewAuthService(
	userRepo domain.UserRepository,
	tokenRepo domain.RefreshTokenRepository,
	hasher domain.PasswordHasher,
	jwtManager *jwt.Manager,
	logger *slog.Logger,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		hasher:     hasher,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// RegisterInput holds the data needed to register a new user.
type RegisterInput struct {
	Email    string
	Password string
	FullName string
	Phone    string
	Role     domain.Role
}

// Register creates a new user account.
func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	// Check if email already exists
	existing, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, fmt.Errorf("check existing user: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	// Default role is buyer
	role := input.Role
	if role == "" {
		role = domain.RoleBuyer
	}

	// Hash password
	hash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.NewString(),
		Email:        input.Email,
		PasswordHash: hash,
		FullName:     input.FullName,
		Phone:        input.Phone,
		Role:         role,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	s.logger.Info("user registered",
		"user_id", user.ID,
		"email", user.Email,
		"role", user.Role,
	)

	return user, nil
}

// LoginInput holds the data needed to authenticate a user.
type LoginInput struct {
	Email    string
	Password string
}

// Login authenticates a user and returns tokens.
func (s *AuthService) Login(ctx context.Context, input LoginInput) (*domain.AuthTokens, *domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, nil, domain.ErrInvalidCredentials
		}
		return nil, nil, fmt.Errorf("find user: %w", err)
	}

	if !user.IsActive {
		return nil, nil, domain.ErrInvalidCredentials
	}

	if err := s.hasher.Compare(user.PasswordHash, input.Password); err != nil {
		return nil, nil, domain.ErrInvalidCredentials
	}

	tokens, err := s.generateTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	s.logger.Info("user logged in",
		"user_id", user.ID,
		"email", user.Email,
	)

	return tokens, user, nil
}

// Refresh generates new tokens from a valid refresh token.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*domain.AuthTokens, *domain.User, error) {
	stored, err := s.tokenRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		return nil, nil, domain.ErrInvalidToken
	}

	if stored.Revoked {
		return nil, nil, domain.ErrInvalidToken
	}

	if time.Now().After(stored.ExpiresAt) {
		return nil, nil, domain.ErrTokenExpired
	}

	user, err := s.userRepo.FindByID(ctx, stored.UserID)
	if err != nil {
		return nil, nil, domain.ErrInvalidToken
	}

	if !user.IsActive {
		return nil, nil, domain.ErrInvalidCredentials
	}

	// Revoke old token
	if err := s.tokenRepo.Revoke(ctx, refreshToken); err != nil {
		s.logger.Warn("failed to revoke old refresh token", "error", err)
	}

	tokens, err := s.generateTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return tokens, user, nil
}

// Logout revokes all refresh tokens for a user.
func (s *AuthService) Logout(ctx context.Context, userID string) error {
	return s.tokenRepo.RevokeAllForUser(ctx, userID)
}

// GetProfile returns the current user's profile.
func (s *AuthService) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// generateTokens creates a new access token and refresh token.
func (s *AuthService) generateTokens(ctx context.Context, user *domain.User) (*domain.AuthTokens, error) {
	accessToken, err := s.jwtManager.Generate(user.ID, user.TenantID, user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshTokenStr := jwt.GenerateRefreshToken()
	refreshToken := &domain.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
		CreatedAt: time.Now(),
	}

	if err := s.tokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return &domain.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}
