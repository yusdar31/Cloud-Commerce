package postgres

import (
	"context"
	"fmt"

	"github.com/cloudcommerce/order-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrderItemRepository is a PostgreSQL implementation of domain.OrderItemRepository.
type OrderItemRepository struct {
	pool *pgxpool.Pool
}

// NewOrderItemRepository creates a new order item repository.
func NewOrderItemRepository(pool *pgxpool.Pool) *OrderItemRepository {
	return &OrderItemRepository{pool: pool}
}

// Create creates a new order item record.
func (r *OrderItemRepository) Create(ctx context.Context, item *domain.OrderItem) error {
	query := `
		INSERT INTO order_items (
			id, order_id, tenant_id, product_id, variant_id, product_name, sku,
			quantity, unit_price, total_price, tax_amount, discount_amount,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.pool.Exec(ctx, query,
		item.ID, item.OrderID, item.TenantID, item.ProductID, item.VariantID,
		item.ProductName, item.SKU, item.Quantity, item.UnitPrice, item.TotalPrice,
		item.TaxAmount, item.DiscountAmount, item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create order item: %w", err)
	}
	return nil
}

// FindByOrderID finds all items for an order.
func (r *OrderItemRepository) FindByOrderID(ctx context.Context, orderID string) ([]*domain.OrderItem, error) {
	query := `
		SELECT id, order_id, tenant_id, product_id, variant_id, product_name, sku,
			   quantity, unit_price, total_price, tax_amount, discount_amount,
			   created_at, updated_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order items: %w", err)
	}
	defer rows.Close()

	var items []*domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(
			&item.ID, &item.OrderID, &item.TenantID, &item.ProductID, &item.VariantID,
			&item.ProductName, &item.SKU, &item.Quantity, &item.UnitPrice, &item.TotalPrice,
			&item.TaxAmount, &item.DiscountAmount, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, &item)
	}
	return items, nil
}

// DeleteByOrderID deletes all items for an order.
func (r *OrderItemRepository) DeleteByOrderID(ctx context.Context, orderID string) error {
	query := `DELETE FROM order_items WHERE order_id = $1`
	_, err := r.pool.Exec(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order items: %w", err)
	}
	return nil
}
