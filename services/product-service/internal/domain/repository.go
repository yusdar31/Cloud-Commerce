package domain

import "context"

// ProductRepository defines the interface for persisting Product aggregates.
type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id, tenantID string) (*Product, error)
	FindBySKU(ctx context.Context, sku, tenantID string) (*Product, error)
	List(ctx context.Context, tenantID string, limit, offset int, status string) ([]*Product, int64, error)
	Update(ctx context.Context, product *Product) error
	SoftDelete(ctx context.Context, id, tenantID string) error
}

// CategoryRepository defines the interface for persisting Category entities.
type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	FindByID(ctx context.Context, id, tenantID string) (*Category, error)
	List(ctx context.Context, tenantID string) ([]*Category, error)
	Update(ctx context.Context, category *Category) error
	SoftDelete(ctx context.Context, id, tenantID string) error
}
