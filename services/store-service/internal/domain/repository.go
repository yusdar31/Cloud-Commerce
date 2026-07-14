package domain

import (
	"context"
)

// TenantRepository defines the interface for persisting Tenant aggregates.
type TenantRepository interface {
	Create(ctx context.Context, tenant *Tenant) error
	FindByID(ctx context.Context, id string) (*Tenant, error)
	FindBySlug(ctx context.Context, slug string) (*Tenant, error)
	FindByOwnerID(ctx context.Context, ownerID string) (*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) error
}
