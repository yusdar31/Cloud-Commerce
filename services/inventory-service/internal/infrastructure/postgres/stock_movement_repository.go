package postgres

import (
	"context"
	"fmt"

	"github.com/cloudcommerce/inventory-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StockMovementRepository is a PostgreSQL implementation of domain.StockMovementRepository.
type StockMovementRepository struct {
	pool *pgxpool.Pool
}

// NewStockMovementRepository creates a new stock movement repository.
func NewStockMovementRepository(pool *pgxpool.Pool) *StockMovementRepository {
	return &StockMovementRepository{pool: pool}
}

// Create creates a new stock movement record.
func (r *StockMovementRepository) Create(ctx context.Context, movement *domain.StockMovement) error {
	query := `
		INSERT INTO stock_movements (
			id, tenant_id, product_id, sku, movement_type,
			quantity_before, quantity_after, quantity_change,
			reference_id, reference_type, notes, created_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.pool.Exec(ctx, query,
		movement.ID, movement.TenantID, movement.ProductID, movement.SKU, movement.MovementType,
		movement.QuantityBefore, movement.QuantityAfter, movement.QuantityChange,
		movement.ReferenceID, movement.ReferenceType, movement.Notes, movement.CreatedAt, movement.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to create stock movement: %w", err)
	}
	return nil
}

// FindByProductID finds stock movements for a product with pagination.
func (r *StockMovementRepository) FindByProductID(ctx context.Context, tenantID, productID string, offset, limit int) ([]*domain.StockMovement, error) {
	query := `
		SELECT id, tenant_id, product_id, sku, movement_type,
		       quantity_before, quantity_after, quantity_change,
		       reference_id, reference_type, notes, created_at, created_by
		FROM stock_movements
		WHERE tenant_id = $1 AND product_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := r.pool.Query(ctx, query, tenantID, productID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find stock movements: %w", err)
	}
	defer rows.Close()

	var movements []*domain.StockMovement
	for rows.Next() {
		var m domain.StockMovement
		if err := rows.Scan(
			&m.ID, &m.TenantID, &m.ProductID, &m.SKU, &m.MovementType,
			&m.QuantityBefore, &m.QuantityAfter, &m.QuantityChange,
			&m.ReferenceID, &m.ReferenceType, &m.Notes, &m.CreatedAt, &m.CreatedBy,
		); err != nil {
			return nil, fmt.Errorf("failed to scan stock movement: %w", err)
		}
		movements = append(movements, &m)
	}

	return movements, nil
}

// FindByReference finds stock movements by reference (order_id, reservation_id, etc).
func (r *StockMovementRepository) FindByReference(ctx context.Context, tenantID, referenceID, referenceType string) ([]*domain.StockMovement, error) {
	query := `
		SELECT id, tenant_id, product_id, sku, movement_type,
		       quantity_before, quantity_after, quantity_change,
		       reference_id, reference_type, notes, created_at, created_by
		FROM stock_movements
		WHERE tenant_id = $1 AND reference_id = $2 AND reference_type = $3
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, tenantID, referenceID, referenceType)
	if err != nil {
		return nil, fmt.Errorf("failed to find stock movements by reference: %w", err)
	}
	defer rows.Close()

	var movements []*domain.StockMovement
	for rows.Next() {
		var m domain.StockMovement
		if err := rows.Scan(
			&m.ID, &m.TenantID, &m.ProductID, &m.SKU, &m.MovementType,
			&m.QuantityBefore, &m.QuantityAfter, &m.QuantityChange,
			&m.ReferenceID, &m.ReferenceType, &m.Notes, &m.CreatedAt, &m.CreatedBy,
		); err != nil {
			return nil, fmt.Errorf("failed to scan stock movement: %w", err)
		}
		movements = append(movements, &m)
	}

	return movements, nil
}
