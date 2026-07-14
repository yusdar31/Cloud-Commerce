package transport

import (
	"time"

	"github.com/cloudcommerce/store-service/internal/domain"
	"github.com/go-playground/validator/v10"
)

// CreateStoreRequest represents the HTTP body for creating a store.
type CreateStoreRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Slug        string `json:"slug" validate:"omitempty,min=3,max=50"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	LogoURL     string `json:"logo_url" validate:"omitempty,url"`
	BannerURL   string `json:"banner_url" validate:"omitempty,url"`
}

// UpdateStoreRequest represents the HTTP body for updating a store.
type UpdateStoreRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Slug        string `json:"slug" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	LogoURL     string `json:"logo_url" validate:"omitempty,url"`
	BannerURL   string `json:"banner_url" validate:"omitempty,url"`
	Status      string `json:"status" validate:"omitempty,oneof=active suspended"`
}

// StoreResponse represents the serialization structure for a store.
type StoreResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	OwnerID     string    `json:"owner_id"`
	Description string    `json:"description,omitempty"`
	LogoURL     string    `json:"logo_url,omitempty"`
	BannerURL   string    `json:"banner_url,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     int       `json:"version"`
}

// Validate validates structs using validator tags.
func Validate(req interface{}) error {
	v := validator.New()
	return v.Struct(req)
}

// ToStoreResponse converts a domain Tenant to a StoreResponse DTO.
func ToStoreResponse(t *domain.Tenant) StoreResponse {
	return StoreResponse{
		ID:          t.ID,
		Name:        t.Name,
		Slug:        t.Slug,
		OwnerID:     t.OwnerID,
		Description: t.Description,
		LogoURL:     t.LogoURL,
		BannerURL:   t.BannerURL,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		Version:     t.Version,
	}
}
