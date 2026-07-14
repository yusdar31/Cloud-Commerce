package domain

import (
	"context"
	"time"
)

// PaymentRepository defines the interface for payment persistence.
type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	FindByID(ctx context.Context, id string) (*Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (*Payment, error)
	Update(ctx context.Context, payment *Payment) error
	UpdateStatus(ctx context.Context, id string, status PaymentStatus, paidAt *time.Time) error
	ListByTenant(ctx context.Context, tenantID string, offset, limit int) ([]*Payment, error)
	ListByUser(ctx context.Context, userID string, offset, limit int) ([]*Payment, error)
	Delete(ctx context.Context, id string) error
}

// PaymentTransactionRepository defines the interface for transaction audit trail.
type PaymentTransactionRepository interface {
	Create(ctx context.Context, tx *PaymentTransaction) error
	FindByPaymentID(ctx context.Context, paymentID string, offset, limit int) ([]*PaymentTransaction, error)
}

// WebhookEventRepository defines the interface for webhook idempotency.
type WebhookEventRepository interface {
	Create(ctx context.Context, event *WebhookEvent) error
	FindByEventID(ctx context.Context, eventID string) (*WebhookEvent, error)
	MarkProcessed(ctx context.Context, eventID string) error
	Exists(ctx context.Context, eventID string) (bool, error)
}

// PaymentGateway defines the interface for payment providers (Midtrans, Xendit, etc).
type PaymentGateway interface {
	// CreatePayment creates a payment transaction with the provider
	CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*CreatePaymentResponse, error)
	
	// GetPaymentStatus retrieves payment status from provider
	GetPaymentStatus(ctx context.Context, providerPaymentID string) (*PaymentStatusResponse, error)
	
	// VerifyWebhook verifies webhook signature
	VerifyWebhook(payload []byte, signature string) (*WebhookData, error)
	
	// Refund processes a refund
	Refund(ctx context.Context, req *RefundRequest) (*RefundResponse, error)
}

// CreatePaymentRequest is the request for creating a payment.
type CreatePaymentRequest struct {
	OrderID       string
	Amount        int64
	Currency      string
	CustomerEmail string
	CustomerName  string
	PaymentMethod string
	CallbackURL   string
}

// CreatePaymentResponse is the response from payment creation.
type CreatePaymentResponse struct {
	TransactionID     string
	PaymentID         string
	PaymentURL        string
	ExpiresAt         *time.Time
	PaymentMethod     string
}

// PaymentStatusResponse is the response for payment status.
type PaymentStatusResponse struct {
	TransactionID string
	PaymentID     string
	Status        PaymentStatus
	PaidAt        *time.Time
	PaymentMethod string
}

// WebhookData represents parsed webhook data.
type WebhookData struct {
	EventID          string
	EventType        string
	TransactionID    string
	PaymentID        string
	Status           PaymentStatus
	Amount           int64
	PaymentMethod    string
	PaidAt           *time.Time
	RawPayload       map[string]interface{}
}

// RefundRequest is the request for refund.
type RefundRequest struct {
	PaymentID      string
	Amount         int64
	Reason         string
	RefundKey      string
}

// RefundResponse is the response for refund.
type RefundResponse struct {
	RefundID     string
	Status       string
	RefundedAt   time.Time
}
