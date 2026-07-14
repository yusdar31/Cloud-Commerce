package domain

import (
	"encoding/json"
	"errors"
	"time"
)

// Common domain errors.
var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderStatus = errors.New("invalid order status transition")
	ErrOrderAlreadyPaid   = errors.New("order already paid")
	ErrInsufficientStock  = errors.New("insufficient stock")
)

// OrderStatus represents the state of an order.
type OrderStatus string

const (
	OrderPending    OrderStatus = "pending"
	OrderConfirmed  OrderStatus = "confirmed"
	OrderPaid       OrderStatus = "paid"
	OrderProcessing OrderStatus = "processing"
	OrderShipped    OrderStatus = "shipped"
	OrderDelivered  OrderStatus = "delivered"
	OrderCancelled  OrderStatus = "cancelled"
	OrderRefunded   OrderStatus = "refunded"
	OrderFailed     OrderStatus = "failed"
)

// ValidTransitions defines valid status transitions.
var ValidTransitions = map[OrderStatus][]OrderStatus{
	OrderPending:    {OrderConfirmed, OrderCancelled, OrderFailed},
	OrderConfirmed:  {OrderPaid, OrderCancelled, OrderFailed},
	OrderPaid:       {OrderProcessing, OrderRefunded},
	OrderProcessing: {OrderShipped, OrderCancelled},
	OrderShipped:    {OrderDelivered},
	OrderDelivered:  {},
	OrderCancelled:  {},
	OrderRefunded:   {},
	OrderFailed:     {OrderPending},
}

// CanTransitionTo checks if a status transition is valid.
func (s OrderStatus) CanTransitionTo(newStatus OrderStatus) bool {
	allowed, ok := ValidTransitions[s]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == newStatus {
			return true
		}
	}
	return false
}

// Order represents an order aggregate.
type Order struct {
	ID               string
	TenantID         string
	UserID           string
	OrderNumber      string
	Status           OrderStatus
	Subtotal         int64
	TaxAmount        int64
	ShippingAmount   int64
	DiscountAmount   int64
	TotalAmount      int64
	Currency         string
	Notes            string
	ShippingAddress  json.RawMessage
	BillingAddress   json.RawMessage
	Items            []*OrderItem
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
	CreatedBy        string
	UpdatedBy        string
	Version          int
}

// OrderItem represents a line item in an order.
type OrderItem struct {
	ID            string
	OrderID       string
	TenantID      string
	ProductID     string
	VariantID     string
	ProductName   string
	SKU           string
	Quantity      int
	UnitPrice     int64
	TotalPrice    int64
	TaxAmount     int64
	DiscountAmount int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// OrderStatusHistory represents a status change record.
type OrderStatusHistory struct {
	ID        string
	OrderID   string
	Status    OrderStatus
	Notes     string
	ChangedBy string
	CreatedAt time.Time
}

// SagaEvent represents a saga orchestration event.
type SagaEvent struct {
	ID           string
	OrderID      string
	EventType    string
	EventData    json.RawMessage
	Status       string
	ErrorMessage string
	CreatedAt    time.Time
	ProcessedAt  *time.Time
}
