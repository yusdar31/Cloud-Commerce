package domain

import (
	"errors"
	"time"
)

// Common domain errors.
var (
	ErrInventoryNotFound      = errors.New("inventory not found")
	ErrInsufficientStock      = errors.New("insufficient stock")
	ErrReservationNotFound    = errors.New("reservation not found")
	ErrReservationExpired     = errors.New("reservation expired")
	ErrInvalidQuantity        = errors.New("invalid quantity")
	ErrProductAlreadyExists   = errors.New("product already exists in inventory")
)

// Inventory represents stock levels for a product.
type Inventory struct {
	ID                 string
	TenantID           string
	ProductID          string
	SKU                string
	Quantity           int
	Reserved           int
	Available          int // Computed: Quantity - Reserved
	LowStockThreshold  int
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
	CreatedBy          string
	UpdatedBy          string
	Version            int
}

// IsLowStock checks if inventory is below threshold.
func (i *Inventory) IsLowStock() bool {
	return i.Available <= i.LowStockThreshold
}

// CanReserve checks if quantity can be reserved.
func (i *Inventory) CanReserve(quantity int) bool {
	return i.Available >= quantity
}

// Reservation represents a temporary stock reservation.
type Reservation struct {
	ID         string
	TenantID   string
	OrderID    string
	ProductID  string
	SKU        string
	Quantity   int
	Status     ReservationStatus
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ReservationStatus represents the state of a reservation.
type ReservationStatus string

const (
	ReservationActive    ReservationStatus = "active"
	ReservationReleased  ReservationStatus = "released"
	ReservationFulfilled ReservationStatus = "fulfilled"
)

// IsExpired checks if reservation has expired.
func (r *Reservation) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// StockMovement represents an audit trail entry for stock changes.
type StockMovement struct {
	ID             string
	TenantID       string
	ProductID      string
	SKU            string
	MovementType   MovementType
	QuantityBefore int
	QuantityAfter  int
	QuantityChange int
	ReferenceID    string
	ReferenceType  string
	Notes          string
	CreatedAt      time.Time
	CreatedBy      string
}

// MovementType represents the type of stock movement.
type MovementType string

const (
	MovementInitial     MovementType = "initial"
	MovementAdjustment  MovementType = "adjustment"
	MovementReservation MovementType = "reservation"
	MovementRelease     MovementType = "release"
	MovementSale        MovementType = "sale"
	MovementReturn      MovementType = "return"
)
