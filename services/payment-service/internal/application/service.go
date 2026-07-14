package application

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cloudcommerce/payment-service/internal/domain"
	"github.com/cloudcommerce/shared-go/events"
	"github.com/cloudcommerce/shared-go/nats"
	"github.com/google/uuid"
)

// PaymentService handles payment business logic.
type PaymentService struct {
	paymentRepo     domain.PaymentRepository
	transactionRepo domain.PaymentTransactionRepository
	webhookRepo     domain.WebhookEventRepository
	gateway        domain.PaymentGateway
	publisher      *nats.Publisher
	logger         *slog.Logger
}

// NewPaymentService creates a new payment service.
func NewPaymentService(
	paymentRepo domain.PaymentRepository,
	transactionRepo domain.PaymentTransactionRepository,
	webhookRepo domain.WebhookEventRepository,
	gateway domain.PaymentGateway,
	publisher *nats.Publisher,
	logger *slog.Logger,
) *PaymentService {
	return &PaymentService{
		paymentRepo:     paymentRepo,
		transactionRepo: transactionRepo,
		webhookRepo:     webhookRepo,
		gateway:        gateway,
		publisher:      publisher,
		logger:         logger,
	}
}

// CreatePayment creates a new payment transaction.
func (s *PaymentService) CreatePayment(ctx context.Context, input CreatePaymentInput) (*domain.Payment, error) {
	// Check if payment already exists for this order
	existing, _ := s.paymentRepo.FindByOrderID(ctx, input.OrderID)
	if existing != nil {
		return nil, fmt.Errorf("payment already exists for order %s", input.OrderID)
	}

	// Create payment with gateway
	gatewayResp, err := s.gateway.CreatePayment(ctx, &domain.CreatePaymentRequest{
		OrderID:       input.OrderID,
		Amount:        input.Amount,
		Currency:      input.Currency,
		CustomerEmail: input.CustomerEmail,
		CustomerName:  input.CustomerName,
		PaymentMethod: input.PaymentMethod,
		CallbackURL:   input.CallbackURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create payment with gateway: %w", err)
	}

	// Create payment record
	payment := &domain.Payment{
		ID:                    uuid.NewString(),
		TenantID:              input.TenantID,
		OrderID:               input.OrderID,
		UserID:                input.UserID,
		Amount:                input.Amount,
		Currency:              input.Currency,
		Status:                domain.PaymentPending,
		PaymentMethod:         gatewayResp.PaymentMethod,
		PaymentProvider:       input.Provider,
		ProviderTransactionID: gatewayResp.TransactionID,
		ProviderPaymentID:     gatewayResp.PaymentID,
		PaymentURL:            gatewayResp.PaymentURL,
		ExpiresAt:             gatewayResp.ExpiresAt,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		CreatedBy:             input.UserID,
		Version:               1,
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	// Publish event
	event := events.PaymentPendingEvent{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.NewString(),
			Timestamp: time.Now(),
			TenantID:  input.TenantID,
			UserID:    input.UserID,
		},
		PaymentID:     payment.ID,
		OrderID:       input.OrderID,
		Amount:        input.Amount,
		Currency:      input.Currency,
		PaymentMethod: payment.PaymentMethod,
		PaymentURL:    payment.PaymentURL,
		ExpiresAt:     *gatewayResp.ExpiresAt,
	}
	s.publisher.Publish(ctx, events.SubjectPaymentPending, event)

	s.logger.Info("payment created",
		"payment_id", payment.ID,
		"order_id", input.OrderID,
		"amount", input.Amount,
	)

	return payment, nil
}

// HandleWebhook processes a webhook from payment provider.
func (s *PaymentService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	// Verify webhook
	webhookData, err := s.gateway.VerifyWebhook(payload, signature)
	if err != nil {
		return fmt.Errorf("invalid webhook signature: %w", err)
	}

	// Check idempotency
	exists, err := s.webhookRepo.Exists(ctx, webhookData.EventID)
	if err != nil {
		return fmt.Errorf("failed to check webhook idempotency: %w", err)
	}
	if exists {
		return domain.ErrWebhookAlreadyProcessed
	}

	// Store webhook event for idempotency
	webhookEvent := &domain.WebhookEvent{
		ID:        uuid.NewString(),
		TenantID:  "", // Will be filled from payment lookup
		EventID:   webhookData.EventID,
		EventType: webhookData.EventType,
		Provider:  "midtrans",
		Payload:   webhookData.RawPayload,
		Signature: signature,
		Processed: false,
		CreatedAt: time.Now(),
	}

	// Find payment
	payment, err := s.paymentRepo.FindByID(ctx, webhookData.PaymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	webhookEvent.TenantID = payment.TenantID

	// Store webhook event
	if err := s.webhookRepo.Create(ctx, webhookEvent); err != nil {
		return fmt.Errorf("failed to store webhook event: %w", err)
	}

	// Update payment status
	var newStatus domain.PaymentStatus
	switch webhookData.Status {
	case domain.PaymentPaid:
		newStatus = domain.PaymentPaid
		payment.PaidAt = webhookData.PaidAt
	case domain.PaymentFailed:
		newStatus = domain.PaymentFailed
	default:
		newStatus = webhookData.Status
	}

	if err := s.paymentRepo.UpdateStatus(ctx, payment.ID, newStatus, payment.PaidAt); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	// Mark webhook as processed
	s.webhookRepo.MarkProcessed(ctx, webhookData.EventID)

	// Publish event based on status
	if newStatus == domain.PaymentPaid {
		event := events.PaymentSuccessEvent{
			BaseEvent: events.BaseEvent{
				EventID:   uuid.NewString(),
				Timestamp: time.Now(),
				TenantID:  payment.TenantID,
				UserID:    payment.UserID,
			},
			PaymentID:     payment.ID,
			OrderID:       payment.OrderID,
			Amount:        payment.Amount,
			Currency:      payment.Currency,
			PaymentMethod: payment.PaymentMethod,
			TransactionID: payment.ProviderTransactionID,
			PaidAt:        *payment.PaidAt,
		}
		s.publisher.Publish(ctx, events.SubjectPaymentSuccess, event)
	} else if newStatus == domain.PaymentFailed {
		event := events.PaymentFailedEvent{
			BaseEvent: events.BaseEvent{
				EventID:   uuid.NewString(),
				Timestamp: time.Now(),
				TenantID:  payment.TenantID,
				UserID:    payment.UserID,
			},
			PaymentID: payment.ID,
			OrderID:   payment.OrderID,
			Amount:    payment.Amount,
			Currency:  payment.Currency,
		}
		s.publisher.Publish(ctx, events.SubjectPaymentFailed, event)
	}

	s.logger.Info("webhook processed",
		"event_id", webhookData.EventID,
		"payment_id", payment.ID,
		"status", newStatus,
	)

	return nil
}

// CreatePaymentInput contains data for creating a payment.
type CreatePaymentInput struct {
	TenantID      string
	OrderID       string
	UserID        string
	Amount        int64
	Currency      string
	PaymentMethod string
	Provider      string
	CustomerEmail string
	CustomerName  string
	CallbackURL   string
}
