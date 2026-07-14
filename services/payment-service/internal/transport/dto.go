package transport

import "time"

// CreatePaymentRequest is the request for creating a payment.
type CreatePaymentRequest struct {
	OrderID       string `json:"order_id" validate:"required,uuid"`
	Amount        int64  `json:"amount" validate:"required,min=1"`
	Currency      string `json:"currency" validate:"required,len=3"`
	PaymentMethod string `json:"payment_method" validate:"required"`
	Provider      string `json:"provider" validate:"required"`
}

// PaymentResponse is the response for payment.
type PaymentResponse struct {
	ID                    string     `json:"id"`
	TenantID              string     `json:"tenant_id"`
	OrderID               string     `json:"order_id"`
	UserID                string     `json:"user_id"`
	Amount                int64      `json:"amount"`
	Currency              string     `json:"currency"`
	Status                string     `json:"status"`
	PaymentMethod         string     `json:"payment_method"`
	PaymentProvider       string     `json:"payment_provider"`
	ProviderTransactionID string     `json:"provider_transaction_id,omitempty"`
	ProviderPaymentID     string     `json:"provider_payment_id,omitempty"`
	PaymentURL            string     `json:"payment_url,omitempty"`
	ExpiresAt             *time.Time `json:"expires_at,omitempty"`
	PaidAt                *time.Time `json:"paid_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// WebhookRequest represents an incoming webhook payload.
type WebhookRequest struct {
	RawPayload []byte
	Signature  string
}

// PaymentStatusResponse is the response for payment status check.
type PaymentStatusResponse struct {
	PaymentID string `json:"payment_id"`
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
}

// PaymentListResponse is the response for listing payments.
type PaymentListResponse struct {
	Data       []PaymentResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
}
