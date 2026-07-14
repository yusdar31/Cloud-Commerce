package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/cloudcommerce/product-service/internal/domain"
	"github.com/google/uuid"
)

// CategoryService implements category use cases.
type CategoryService struct {
	repo   domain.CategoryRepository
	logger *slog.Logger
}

// NewCategoryService creates a new CategoryService.
func NewCategoryService(repo domain.CategoryRepository, logger *slog.Logger) *CategoryService {
	return &CategoryService{repo: repo, logger: logger}
}

// CreateCategoryInput holds the data needed to create a category.
type CreateCategoryInput struct {
	TenantID    string
	Name        string
	Description string
	ParentID    *string
}

// Create creates a new category.
func (s *CategoryService) Create(ctx context.Context, input CreateCategoryInput) (*domain.Category, error) {
	if input.TenantID == "" {
		return nil, errors.New("tenant_id is required")
	}

	category := &domain.Category{
		ID:          uuid.NewString(),
		TenantID:    input.TenantID,
		Name:        input.Name,
		Slug:        slugify(input.Name),
		Description: input.Description,
		ParentID:    input.ParentID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}

	s.logger.Info("category created", "category_id", category.ID, "tenant_id", category.TenantID)
	return category, nil
}

// GetByID retrieves a category by ID within a tenant.
func (s *CategoryService) GetByID(ctx context.Context, id, tenantID string) (*domain.Category, error) {
	return s.repo.FindByID(ctx, id, tenantID)
}

// List retrieves all categories for a tenant.
func (s *CategoryService) List(ctx context.Context, tenantID string) ([]*domain.Category, error) {
	return s.repo.List(ctx, tenantID)
}

// Delete soft-deletes a category.
func (s *CategoryService) Delete(ctx context.Context, id, tenantID string) error {
	return s.repo.SoftDelete(ctx, id, tenantID)
}
