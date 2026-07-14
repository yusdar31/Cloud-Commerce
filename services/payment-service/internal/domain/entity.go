package domain

import (
	"errors"
	"time"
)

// Common domain errors.
var (
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrPaymentAlreadyPaid   = errors.New("payment already paid")
	ErrPaymentAlreadyFailed = errors.New("payment already failed")
	ErrPaymentExpired       = errors.New("payment expired")
	ErrInvalidAmount        = errors.New("invalid amount")
	ErrWebhookAlreadyProcessed = errors.New("webhook already processed")
)

// PaymentStatus represents the state of a payment.
type PaymentStatus string

const (
	PaymentPending    PaymentStatus = "pending"
	PaymentProcessing PaymentStatus = "processing"
	PaymentPaid       PaymentStatus = "paid"
	PaymentFailed     PaymentStatus = "failed"
	PaymentRefunded   PaymentStatus = "refunded"
	PaymentCancelled  PaymentStatus = "cancelled"
)

// Payment represents a payment transaction.
type Payment struct {
	ID                    string
	TenantID              string
	OrderID               string
	UserID                string
	Amount                int64
	Currency              string
	Status                PaymentStatus
	PaymentMethod         string
	PaymentProvider       string
	ProviderTransactionID string
	ProviderPaymentID     string
	PaymentURL            string
	ExpiresAt             *time.Time
	PaidAt                *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
	CreatedBy             string
	UpdatedBy             string
	Version               int
}

// IsPending checks if payment is pending.
func (p *Payment) IsPending() bool {
	return p.Status == PaymentPending
}

// IsPaid checks if payment is paid.
func (p *Payment) IsPaid() bool {
	return p.Status == PaymentPaid
}

// IsExpired checks if payment has expired.
func (p *Payment) IsExpired() bool {
	return p.ExpiresAt != nil && time.Now().After(*p.ExpiresAt)
}

// CanProcess checks if payment can be processed.
func (p *Payment) CanProcess() bool {
	return p.Status == PaymentPending && !p.IsExpired()
}

// PaymentTransaction represents a payment transaction audit trail.
type PaymentTransaction struct {
	ID               string
	TenantID         string
	PaymentID        string
	TransactionType  TransactionType
	Status           string
	Amount           int64
	Currency         string
	ProviderResponse map[string]interface{}
	ErrorMessage     string
	CreatedAt        time.Time
}

// TransactionType represents the type of transaction.
type TransactionType string

const (
	TransactionPayment   TransactionType = "payment"
	TransactionRefund    TransactionType = "refund"
	TransactionChargeback TransactionType = "chargeback"
)

// WebhookEvent represents a webhook event for idempotency.
type WebhookEvent struct {
	ID          string
	TenantID    string
	EventID     string
	EventType   string
	Provider    string
	Payload     map[string]interface{}
	Signature   string
	Processed   bool
	ProcessedAt *time.Time
	CreatedAt   time.Time
}
