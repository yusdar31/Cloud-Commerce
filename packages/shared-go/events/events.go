package events

import "time"

// Event subjects for NATS JetStream
const (
	// Order events
	SubjectOrderCreated    = "order.created"
	SubjectOrderUpdated    = "order.updated"
	SubjectOrderCancelled  = "order.cancelled"
	SubjectOrderCompleted  = "order.completed"

	// Payment events
	SubjectPaymentPending   = "payment.pending"
	SubjectPaymentSuccess   = "payment.success"
	SubjectPaymentFailed    = "payment.failed"
	SubjectPaymentRefunded  = "payment.refunded"

	// Inventory events
	SubjectInventoryReserved = "inventory.reserved"
	SubjectInventoryReleased = "inventory.released"
	SubjectInventoryUpdated  = "inventory.updated"
	SubjectStockLow          = "inventory.stock_low"

	// Notification events
	SubjectNotificationEmail = "notification.email"
	SubjectNotificationSMS   = "notification.sms"
)

// BaseEvent contains common fields for all events.
type BaseEvent struct {
	EventID   string    `json:"event_id"`
	Timestamp time.Time `json:"timestamp"`
	TenantID  string    `json:"tenant_id"`
	UserID    string    `json:"user_id,omitempty"`
}

// OrderCreatedEvent is published when a new order is created.
type OrderCreatedEvent struct {
	BaseEvent
	OrderID      string           `json:"order_id"`
	CustomerID   string           `json:"customer_id"`
	Items        []OrderItemEvent `json:"items"`
	TotalAmount  int64            `json:"total_amount"`
	Currency     string           `json:"currency"`
	Status       string           `json:"status"`
}

// OrderItemEvent represents an item in an order.
type OrderItemEvent struct {
	ProductID string `json:"product_id"`
	SKU       string `json:"sku"`
	Quantity  int    `json:"quantity"`
	Price     int64  `json:"price"`
}

// OrderUpdatedEvent is published when an order status changes.
type OrderUpdatedEvent struct {
	BaseEvent
	OrderID     string `json:"order_id"`
	OldStatus   string `json:"old_status"`
	NewStatus   string `json:"new_status"`
	UpdatedBy   string `json:"updated_by"`
	Reason      string `json:"reason,omitempty"`
}

// OrderCancelledEvent is published when an order is cancelled.
type OrderCancelledEvent struct {
	BaseEvent
	OrderID     string `json:"order_id"`
	Reason      string `json:"reason"`
	CancelledBy string `json:"cancelled_by"`
}

// OrderCompletedEvent is published when an order is completed.
type OrderCompletedEvent struct {
	BaseEvent
	OrderID       string    `json:"order_id"`
	CompletedAt   time.Time `json:"completed_at"`
	TotalAmount   int64     `json:"total_amount"`
}

// PaymentPendingEvent is published when payment is initiated.
type PaymentPendingEvent struct {
	BaseEvent
	PaymentID      string `json:"payment_id"`
	OrderID        string `json:"order_id"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	PaymentMethod  string `json:"payment_method"`
	PaymentURL     string `json:"payment_url,omitempty"`
	ExpiresAt      time.Time `json:"expires_at"`
}

// PaymentSuccessEvent is published when payment is confirmed.
type PaymentSuccessEvent struct {
	BaseEvent
	PaymentID        string    `json:"payment_id"`
	OrderID          string    `json:"order_id"`
	Amount           int64     `json:"amount"`
	Currency         string    `json:"currency"`
	PaymentMethod    string    `json:"payment_method"`
	TransactionID    string    `json:"transaction_id"`
	PaidAt           time.Time `json:"paid_at"`
}

// PaymentFailedEvent is published when payment fails.
type PaymentFailedEvent struct {
	BaseEvent
	PaymentID     string `json:"payment_id"`
	OrderID       string `json:"order_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	FailureReason string `json:"failure_reason"`
	ErrorCode     string `json:"error_code,omitempty"`
}

// PaymentRefundedEvent is published when payment is refunded.
type PaymentRefundedEvent struct {
	BaseEvent
	PaymentID     string    `json:"payment_id"`
	OrderID       string    `json:"order_id"`
	RefundAmount  int64     `json:"refund_amount"`
	Currency      string    `json:"currency"`
	RefundReason  string    `json:"refund_reason"`
	RefundedAt    time.Time `json:"refunded_at"`
}

// InventoryReservedEvent is published when stock is reserved for an order.
type InventoryReservedEvent struct {
	BaseEvent
	ReservationID string                  `json:"reservation_id"`
	OrderID       string                  `json:"order_id"`
	Items         []InventoryItemEvent    `json:"items"`
	ExpiresAt     time.Time               `json:"expires_at"`
}

// InventoryItemEvent represents an inventory item.
type InventoryItemEvent struct {
	ProductID      string `json:"product_id"`
	SKU            string `json:"sku"`
	QuantityBefore int    `json:"quantity_before"`
	QuantityAfter  int    `json:"quantity_after"`
	Reserved       int    `json:"reserved"`
}

// InventoryReleasedEvent is published when reserved stock is released.
type InventoryReleasedEvent struct {
	BaseEvent
	ReservationID string                  `json:"reservation_id"`
	OrderID       string                  `json:"order_id"`
	Items         []InventoryItemEvent    `json:"items"`
	Reason        string                  `json:"reason"`
}

// InventoryUpdatedEvent is published when stock levels change.
type InventoryUpdatedEvent struct {
	BaseEvent
	ProductID      string `json:"product_id"`
	SKU            string `json:"sku"`
	OldQuantity    int    `json:"old_quantity"`
	NewQuantity    int    `json:"new_quantity"`
	UpdateType     string `json:"update_type"` // "manual", "reservation", "release", "sale"
}

// StockLowEvent is published when stock falls below threshold.
type StockLowEvent struct {
	BaseEvent
	ProductID   string `json:"product_id"`
	SKU         string `json:"sku"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	Threshold   int    `json:"threshold"`
}

// NotificationEmailEvent is published to send email notifications.
type NotificationEmailEvent struct {
	BaseEvent
	To       []string          `json:"to"`
	Subject  string            `json:"subject"`
	Template string            `json:"template"`
	Data     map[string]string `json:"data"`
}

// NotificationSMSEvent is published to send SMS notifications.
type NotificationSMSEvent struct {
	BaseEvent
	To      string `json:"to"`
	Message string `json:"message"`
}
