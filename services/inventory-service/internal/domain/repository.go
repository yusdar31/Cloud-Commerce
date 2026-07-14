package domain

import "context"

// InventoryRepository defines the interface for inventory persistence.
type InventoryRepository interface {
	// Inventory operations
	Create(ctx context.Context, inv *Inventory) error
	FindByID(ctx context.Context, id string) (*Inventory, error)
	FindByProductID(ctx context.Context, tenantID, productID string) (*Inventory, error)
	FindBySKU(ctx context.Context, tenantID, sku string) (*Inventory, error)
	List(ctx context.Context, tenantID string, offset, limit int) ([]*Inventory, error)
	Update(ctx context.Context, inv *Inventory) error
	UpdateQuantity(ctx context.Context, id string, quantity int) error
	Delete(ctx context.Context, id string) error
	
	// Stock operations with optimistic locking
	Reserve(ctx context.Context, tenantID, productID string, quantity int) error
	Release(ctx context.Context, tenantID, productID string, quantity int) error
	
	// Low stock query
	FindLowStock(ctx context.Context, tenantID string) ([]*Inventory, error)
}

// ReservationRepository defines the interface for reservation persistence.
type ReservationRepository interface {
	Create(ctx context.Context, res *Reservation) error
	FindByID(ctx context.Context, id string) (*Reservation, error)
	FindByOrderID(ctx context.Context, orderID string) ([]*Reservation, error)
	FindActiveByProduct(ctx context.Context, tenantID, productID string) ([]*Reservation, error)
	UpdateStatus(ctx context.Context, id string, status ReservationStatus) error
	ReleaseExpired(ctx context.Context) (int, error)
	Delete(ctx context.Context, id string) error
}

// StockMovementRepository defines the interface for stock movement audit trail.
type StockMovementRepository interface {
	Create(ctx context.Context, movement *StockMovement) error
	FindByProductID(ctx context.Context, tenantID, productID string, offset, limit int) ([]*StockMovement, error)
	FindByReference(ctx context.Context, tenantID, referenceID, referenceType string) ([]*StockMovement, error)
}
