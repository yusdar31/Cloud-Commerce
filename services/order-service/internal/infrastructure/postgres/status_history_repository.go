package postgres

import (
	"context"
	"fmt"

	"github.com/cloudcommerce/order-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrderStatusHistoryRepository is a PostgreSQL implementation of domain.OrderStatusHistoryRepository.
type OrderStatusHistoryRepository struct {
	pool *pgxpool.Pool
}

// NewOrderStatusHistoryRepository creates a new status history repository.
func NewOrderStatusHistoryRepository(pool *pgxpool.Pool) *OrderStatusHistoryRepository {
	return &OrderStatusHistoryRepository{pool: pool}
}

// Create creates a new status history record.
func (r *OrderStatusHistoryRepository) Create(ctx context.Context, history *domain.OrderStatusHistory) error {
	query := `
		INSERT INTO order_status_history (id, order_id, status, notes, changed_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		history.ID, history.OrderID, history.Status, history.Notes,
		history.ChangedBy, history.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create status history: %w", err)
	}
	return nil
}

// FindByOrderID finds all status history records for an order.
func (r *OrderStatusHistoryRepository) FindByOrderID(ctx context.Context, orderID string) ([]*domain.OrderStatusHistory, error) {
	query := `
		SELECT id, order_id, status, notes, changed_by, created_at
		FROM order_status_history
		WHERE order_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find status history: %w", err)
	}
	defer rows.Close()

	var histories []*domain.OrderStatusHistory
	for rows.Next() {
		var h domain.OrderStatusHistory
		if err := rows.Scan(
			&h.ID, &h.OrderID, &h.Status, &h.Notes, &h.ChangedBy, &h.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan status history: %w", err)
		}
		histories = append(histories, &h)
	}
	return histories, nil
}
