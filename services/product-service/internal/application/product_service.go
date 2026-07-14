package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/cloudcommerce/product-service/internal/domain"
	"github.com/google/uuid"
)

// ProductService implements product use cases.
type ProductService struct {
	repo   domain.ProductRepository
	logger *slog.Logger
}

// NewProductService creates a new ProductService.
func NewProductService(repo domain.ProductRepository, logger *slog.Logger) *ProductService {
	return &ProductService{repo: repo, logger: logger}
}

// CreateProductInput holds the data needed to create a product.
type CreateProductInput struct {
	TenantID    string
	Name        string
	Description string
	SKU         string
	Price       int64
	Currency    string
	CategoryID  *string
	ImageURL    string
	Weight      int
}

// UpdateProductInput holds the data needed to update a product.
type UpdateProductInput struct {
	Name        *string
	Description *string
	Price       *int64
	CategoryID  *string
	ImageURL    *string
	Status      *domain.ProductStatus
	Weight      *int
}

// Create creates a new product.
func (s *ProductService) Create(ctx context.Context, input CreateProductInput) (*domain.Product, error) {
	if input.TenantID == "" {
		return nil, errors.New("tenant_id is required")
	}

	// Check SKU uniqueness
	existing, err := s.repo.FindBySKU(ctx, input.SKU, input.TenantID)
	if err != nil && !errors.Is(err, domain.ErrProductNotFound) {
		return nil, fmt.Errorf("check sku: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrSKUExists
	}

	currency := input.Currency
	if currency == "" {
		currency = "IDR"
	}

	slug := slugify(input.Name)

	product := &domain.Product{
		ID:          uuid.NewString(),
		TenantID:    input.TenantID,
		Name:        input.Name,
		Slug:        slug,
		Description: input.Description,
		SKU:         input.SKU,
		Price:       input.Price,
		Currency:    currency,
		Status:      domain.StatusDraft,
		CategoryID:  input.CategoryID,
		ImageURL:    input.ImageURL,
		Weight:      input.Weight,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	s.logger.Info("product created",
		"product_id", product.ID,
		"tenant_id", product.TenantID,
		"sku", product.SKU,
	)

	return product, nil
}

// GetByID retrieves a product by ID within a tenant.
func (s *ProductService) GetByID(ctx context.Context, id, tenantID string) (*domain.Product, error) {
	product, err := s.repo.FindByID(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// List retrieves a paginated list of products for a tenant.
func (s *ProductService) List(ctx context.Context, tenantID string, limit, offset int, status string) ([]*domain.Product, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	products, total, err := s.repo.List(ctx, tenantID, limit, offset, status)
	if err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}
	return products, total, nil
}

// Update updates an existing product.
func (s *ProductService) Update(ctx context.Context, id, tenantID string, input UpdateProductInput) (*domain.Product, error) {
	product, err := s.repo.FindByID(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		product.Name = *input.Name
		product.Slug = slugify(*input.Name)
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Price != nil {
		product.Price = *input.Price
	}
	if input.CategoryID != nil {
		product.CategoryID = input.CategoryID
	}
	if input.ImageURL != nil {
		product.ImageURL = *input.ImageURL
	}
	if input.Status != nil {
		product.Status = *input.Status
	}
	if input.Weight != nil {
		product.Weight = *input.Weight
	}
	product.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	s.logger.Info("product updated", "product_id", product.ID, "tenant_id", product.TenantID)
	return product, nil
}

// Delete soft-deletes a product.
func (s *ProductService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.repo.SoftDelete(ctx, id, tenantID); err != nil {
		return err
	}
	s.logger.Info("product deleted", "product_id", id, "tenant_id", tenantID)
	return nil
}

// Publish changes a product status to published.
func (s *ProductService) Publish(ctx context.Context, id, tenantID string) (*domain.Product, error) {
	return s.updateStatus(ctx, id, tenantID, domain.StatusPublished)
}

// Archive changes a product status to archived.
func (s *ProductService) Archive(ctx context.Context, id, tenantID string) (*domain.Product, error) {
	return s.updateStatus(ctx, id, tenantID, domain.StatusArchived)
}

func (s *ProductService) updateStatus(ctx context.Context, id, tenantID string, status domain.ProductStatus) (*domain.Product, error) {
	product, err := s.repo.FindByID(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	product.Status = status
	product.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("update product status: %w", err)
	}
	return product, nil
}

// slugify converts a name to a URL-friendly slug.
var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = nonAlnum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "product"
	}
	return s
}
