package postgres

import (
	"context"
	"fmt"

	"github.com/cloudcommerce/payment-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WebhookEventRepository is a PostgreSQL implementation of domain.WebhookEventRepository.
type WebhookEventRepository struct {
	pool *pgxpool.Pool
}

// NewWebhookEventRepository creates a new webhook event repository.
func NewWebhookEventRepository(pool *pgxpool.Pool) *WebhookEventRepository {
	return &WebhookEventRepository{pool: pool}
}

// Create creates a new webhook event record.
func (r *WebhookEventRepository) Create(ctx context.Context, event *domain.WebhookEvent) error {
	query := `
		INSERT INTO webhook_events (
			id, tenant_id, event_id, event_type, provider,
			payload, signature, processed, processed_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(ctx, query,
		event.ID, event.TenantID, event.EventID, event.EventType, event.Provider,
		event.Payload, event.Signature, event.Processed, event.ProcessedAt, event.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create webhook event: %w", err)
	}
	return nil
}

// FindByEventID finds a webhook event by event ID.
func (r *WebhookEventRepository) FindByEventID(ctx context.Context, eventID string) (*domain.WebhookEvent, error) {
	query := `
		SELECT id, tenant_id, event_id, event_type, provider,
			   payload, signature, processed, processed_at, created_at
		FROM webhook_events
		WHERE event_id = $1
	`
	var event domain.WebhookEvent
	err := r.pool.QueryRow(ctx, query, eventID).Scan(
		&event.ID, &event.TenantID, &event.EventID, &event.EventType, &event.Provider,
		&event.Payload, &event.Signature, &event.Processed, &event.ProcessedAt, &event.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrWebhookAlreadyProcessed
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find webhook event: %w", err)
	}
	return &event, nil
}

// MarkProcessed marks a webhook event as processed.
func (r *WebhookEventRepository) MarkProcessed(ctx context.Context, eventID string) error {
	query := `
		UPDATE webhook_events
		SET processed = TRUE, processed_at = NOW()
		WHERE event_id = $1
	`
	_, err := r.pool.Exec(ctx, query, eventID)
	if err != nil {
		return fmt.Errorf("failed to mark webhook as processed: %w", err)
	}
	return nil
}

// Exists checks if a webhook event exists by event ID.
func (r *WebhookEventRepository) Exists(ctx context.Context, eventID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM webhook_events WHERE event_id = $1)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, eventID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check webhook existence: %w", err)
	}
	return exists, nil
}
