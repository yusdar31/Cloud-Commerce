package domain

import "context"

// OrderRepository defines the interface for order persistence.
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, id string) (*Order, error)
	FindByOrderNumber(ctx context.Context, orderNumber string) (*Order, error)
	Update(ctx context.Context, order *Order) error
	UpdateStatus(ctx context.Context, id string, status OrderStatus) error
	ListByTenant(ctx context.Context, tenantID string, offset, limit int) ([]*Order, error)
	ListByUser(ctx context.Context, userID string, offset, limit int) ([]*Order, error)
	Delete(ctx context.Context, id string) error
}

// OrderItemRepository defines the interface for order item persistence.
type OrderItemRepository interface {
	Create(ctx context.Context, item *OrderItem) error
	FindByOrderID(ctx context.Context, orderID string) ([]*OrderItem, error)
	DeleteByOrderID(ctx context.Context, orderID string) error
}

// OrderStatusHistoryRepository defines the interface for status history.
type OrderStatusHistoryRepository interface {
	Create(ctx context.Context, history *OrderStatusHistory) error
	FindByOrderID(ctx context.Context, orderID string) ([]*OrderStatusHistory, error)
}

// SagaEventRepository defines the interface for saga event log.
type SagaEventRepository interface {
	Create(ctx context.Context, event *SagaEvent) error
	FindByOrderID(ctx context.Context, orderID string) ([]*SagaEvent, error)
	MarkProcessed(ctx context.Context, id string) error
}
