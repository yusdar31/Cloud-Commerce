package transport

import "time"

// CreateInventoryRequest is the request for creating inventory.
type CreateInventoryRequest struct {
	ProductID         string `json:"product_id" validate:"required,uuid"`
	SKU               string `json:"sku" validate:"required"`
	Quantity          int    `json:"quantity" validate:"required,min=0"`
	LowStockThreshold int    `json:"low_stock_threshold" validate:"min=0"`
}

// UpdateInventoryRequest is the request for updating inventory.
type UpdateInventoryRequest struct {
	Quantity          int `json:"quantity" validate:"required,min=0"`
	LowStockThreshold int `json:"low_stock_threshold" validate:"min=0"`
}

// ReserveStockRequest is the request for reserving stock.
type ReserveStockRequest struct {
	OrderID   string `json:"order_id" validate:"required,uuid"`
	ProductID string `json:"product_id" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

// ReleaseStockRequest is the request for releasing stock.
type ReleaseStockRequest struct {
	OrderID string `json:"order_id" validate:"required,uuid"`
	Reason  string `json:"reason"`
}

// InventoryResponse is the response for inventory.
type InventoryResponse struct {
	ID                string     `json:"id"`
	TenantID          string     `json:"tenant_id"`
	ProductID         string     `json:"product_id"`
	SKU               string     `json:"sku"`
	Quantity          int        `json:"quantity"`
	Reserved          int        `json:"reserved"`
	Available         int        `json:"available"`
	LowStockThreshold int        `json:"low_stock_threshold"`
	IsLowStock        bool       `json:"is_low_stock"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// ReservationResponse is the response for reservation.
type ReservationResponse struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	ProductID string    `json:"product_id"`
	SKU       string    `json:"sku"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// StockMovementResponse is the response for stock movement.
type StockMovementResponse struct {
	ID             string    `json:"id"`
	ProductID      string    `json:"product_id"`
	SKU            string    `json:"sku"`
	MovementType   string    `json:"movement_type"`
	QuantityBefore int       `json:"quantity_before"`
	QuantityAfter  int       `json:"quantity_after"`
	QuantityChange int       `json:"quantity_change"`
	ReferenceID    string    `json:"reference_id,omitempty"`
	ReferenceType  string    `json:"reference_type,omitempty"`
	Notes          string    `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by,omitempty"`
}

// ReserveStockResponse is the response for stock reservation.
type ReserveStockResponse struct {
	ReservationID string    `json:"reservation_id"`
	ExpiresAt     time.Time `json:"expires_at"`
	Message       string    `json:"message"`
}
