package domain

import (
	"errors"
	"time"
)

// Common domain errors.
var (
	ErrTenantNotFound         = errors.New("tenant not found")
	ErrSlugAlreadyExists      = errors.New("tenant slug already exists")
	ErrOwnerAlreadyHasStore   = errors.New("owner already has a store")
	ErrInvalidTenantStatus    = errors.New("invalid tenant status")
)

// Tenant status values.
const (
	StatusActive    = "active"
	StatusSuspended = "suspended"
)

// Tenant represents a tenant (store) in the system.
type Tenant struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	OwnerID     string     `json:"owner_id"`
	Description string     `json:"description,omitempty"`
	LogoURL     string     `json:"logo_url,omitempty"`
	BannerURL   string     `json:"banner_url,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	CreatedBy   string     `json:"created_by,omitempty"`
	UpdatedBy   string     `json:"updated_by,omitempty"`
	Version     int        `json:"version"`
}
