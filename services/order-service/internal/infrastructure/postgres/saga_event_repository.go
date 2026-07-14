package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudcommerce/order-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SagaEventRepository is a PostgreSQL implementation of domain.SagaEventRepository.
type SagaEventRepository struct {
	pool *pgxpool.Pool
}

// NewSagaEventRepository creates a new saga event repository.
func NewSagaEventRepository(pool *pgxpool.Pool) *SagaEventRepository {
	return &SagaEventRepository{pool: pool}
}

// Create creates a new saga event record.
func (r *SagaEventRepository) Create(ctx context.Context, event *domain.SagaEvent) error {
	query := `
		INSERT INTO saga_events (id, order_id, event_type, event_data, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		event.ID, event.OrderID, event.EventType, event.EventData,
		event.Status, event.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create saga event: %w", err)
	}
	return nil
}

// FindByOrderID finds all saga events for an order.
func (r *SagaEventRepository) FindByOrderID(ctx context.Context, orderID string) ([]*domain.SagaEvent, error) {
	query := `
		SELECT id, order_id, event_type, event_data, status, error_message, created_at, processed_at
		FROM saga_events
		WHERE order_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find saga events: %w", err)
	}
	defer rows.Close()

	var events []*domain.SagaEvent
	for rows.Next() {
		var e domain.SagaEvent
		if err := rows.Scan(
			&e.ID, &e.OrderID, &e.EventType, &e.EventData,
			&e.Status, &e.ErrorMessage, &e.CreatedAt, &e.ProcessedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan saga event: %w", err)
		}
		events = append(events, &e)
	}
	return events, nil
}

// MarkProcessed marks a saga event as processed.
func (r *SagaEventRepository) MarkProcessed(ctx context.Context, id string) error {
	query := `
		UPDATE saga_events
		SET status = 'processed', processed_at = $1
		WHERE id = $2
	`
	_, err := r.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to mark saga event as processed: %w", err)
	}
	return nil
}
