package application

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/cloudcommerce/shared-go/events"
	"github.com/cloudcommerce/shared-go/nats"
	"github.com/cloudcommerce/store-service/internal/domain"
	"github.com/google/uuid"
)

// CreateStoreInput represents the data required to create a store.
type CreateStoreInput struct {
	Name        string
	Slug        string
	OwnerID     string
	Description string
	LogoURL     string
	BannerURL   string
}

// UpdateStoreInput represents the data required to update store settings.
type UpdateStoreInput struct {
	Name        string
	Slug        string
	Description string
	LogoURL     string
	BannerURL   string
	Status      string
}

// StoreService implements store-service business use cases.
type StoreService struct {
	tenantRepo domain.TenantRepository
	publisher  *nats.Publisher
	logger     *slog.Logger
}

// NewStoreService creates a new StoreService.
func NewStoreService(
	tenantRepo domain.TenantRepository,
	publisher *nats.Publisher,
	logger *slog.Logger,
) *StoreService {
	return &StoreService{
		tenantRepo: tenantRepo,
		publisher:  publisher,
		logger:     logger,
	}
}

var slugRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

// CreateStore creates a new store/tenant.
func (s *StoreService) CreateStore(ctx context.Context, input CreateStoreInput) (*domain.Tenant, error) {
	// Normalize slug
	slug := strings.ToLower(strings.TrimSpace(input.Slug))
	if slug == "" {
		// Generate slug from name
		slug = s.generateSlug(input.Name)
	}

	// Validate slug format
	if !slugRegex.MatchString(slug) {
		return nil, fmt.Errorf("invalid slug format: must contain only letters, numbers, and hyphens")
	}

	// Check if owner already has a store
	existingByOwner, err := s.tenantRepo.FindByOwnerID(ctx, input.OwnerID)
	if err != nil && err != domain.ErrTenantNotFound {
		return nil, fmt.Errorf("check existing store for owner: %w", err)
	}
	if existingByOwner != nil {
		return nil, domain.ErrOwnerAlreadyHasStore
	}

	// Check if slug is already taken
	existingBySlug, err := s.tenantRepo.FindBySlug(ctx, slug)
	if err != nil && err != domain.ErrTenantNotFound {
		return nil, fmt.Errorf("check existing store for slug: %w", err)
	}
	if existingBySlug != nil {
		return nil, domain.ErrSlugAlreadyExists
	}

	tenant := &domain.Tenant{
		ID:          uuid.NewString(),
		Name:        input.Name,
		Slug:        slug,
		OwnerID:     input.OwnerID,
		Description: input.Description,
		LogoURL:     input.LogoURL,
		BannerURL:   input.BannerURL,
		Status:      domain.StatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Version:     1,
	}

	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, err
	}

	s.logger.Info("store created",
		"tenant_id", tenant.ID,
		"name", tenant.Name,
		"slug", tenant.Slug,
		"owner_id", tenant.OwnerID,
	)

	// Publish TenantCreated event to NATS
	if s.publisher != nil {
		event := events.TenantCreatedEvent{
			BaseEvent: events.BaseEvent{
				EventID:   uuid.NewString(),
				Timestamp: time.Now(),
				TenantID:  tenant.ID,
				UserID:    tenant.OwnerID,
			},
			OwnerID:    tenant.OwnerID,
			TenantName: tenant.Name,
			TenantSlug: tenant.Slug,
		}

		if err := s.publisher.Publish(ctx, events.SubjectTenantCreated, event); err != nil {
			s.logger.Error("failed to publish TenantCreated event", "error", err)
			// Don't fail the REST request if event publishing fails, but log it
		} else {
			s.logger.Info("published TenantCreated event", "tenant_id", tenant.ID)
		}
	}

	return tenant, nil
}

// GetStoreByID retrieves a store by its ID.
func (s *StoreService) GetStoreByID(ctx context.Context, id string) (*domain.Tenant, error) {
	return s.tenantRepo.FindByID(ctx, id)
}

// GetStoreBySlug retrieves a store by its slug.
func (s *StoreService) GetStoreBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	return s.tenantRepo.FindBySlug(ctx, strings.ToLower(slug))
}

// GetStoreByOwnerID retrieves a store by its owner's ID.
func (s *StoreService) GetStoreByOwnerID(ctx context.Context, ownerID string) (*domain.Tenant, error) {
	return s.tenantRepo.FindByOwnerID(ctx, ownerID)
}

// UpdateStore updates store configuration/branding.
func (s *StoreService) UpdateStore(ctx context.Context, id string, ownerID string, input UpdateStoreInput) (*domain.Tenant, error) {
	tenant, err := s.tenantRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if tenant.OwnerID != ownerID {
		return nil, fmt.Errorf("unauthorized to update this store")
	}

	// Update slug if changed
	newSlug := strings.ToLower(strings.TrimSpace(input.Slug))
	if newSlug != "" && newSlug != tenant.Slug {
		if !slugRegex.MatchString(newSlug) {
			return nil, fmt.Errorf("invalid slug format: must contain only letters, numbers, and hyphens")
		}

		// Check if new slug is taken
		existing, err := s.tenantRepo.FindBySlug(ctx, newSlug)
		if err != nil && err != domain.ErrTenantNotFound {
			return nil, fmt.Errorf("check new slug: %w", err)
		}
		if existing != nil {
			return nil, domain.ErrSlugAlreadyExists
		}
		tenant.Slug = newSlug
	}

	// Update status if changed and valid
	if input.Status != "" && input.Status != tenant.Status {
		if input.Status != domain.StatusActive && input.Status != domain.StatusSuspended {
			return nil, domain.ErrInvalidTenantStatus
		}
		tenant.Status = input.Status
	}

	tenant.Name = input.Name
	tenant.Description = input.Description
	tenant.LogoURL = input.LogoURL
	tenant.BannerURL = input.BannerURL
	tenant.UpdatedAt = time.Now()

	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return nil, err
	}

	s.logger.Info("store updated",
		"tenant_id", tenant.ID,
		"name", tenant.Name,
		"slug", tenant.Slug,
	)

	return tenant, nil
}

func (s *StoreService) generateSlug(name string) string {
	// Simple slug generator: lowercase and replace non-alphanumeric characters with hyphens
	slug := strings.ToLower(name)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 50 {
		slug = slug[:50]
		slug = strings.Trim(slug, "-")
	}
	if slug == "" {
		slug = "store-" + uuid.NewString()[:8]
	}
	return slug
}
