package postgres

import (
	"context"
	"fmt"

	"github.com/cloudcommerce/order-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrderRepository is a PostgreSQL implementation of domain.OrderRepository.
type OrderRepository struct {
	pool *pgxpool.Pool
}

// NewOrderRepository creates a new order repository.
func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

// Create creates a new order record.
func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (
			id, tenant_id, user_id, order_number, status, subtotal, tax_amount,
			shipping_amount, discount_amount, total_amount, currency, notes,
			shipping_address, billing_address, created_at, updated_at, created_by, version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`
	_, err := r.pool.Exec(ctx, query,
		order.ID, order.TenantID, order.UserID, order.OrderNumber, order.Status,
		order.Subtotal, order.TaxAmount, order.ShippingAmount, order.DiscountAmount,
		order.TotalAmount, order.Currency, order.Notes,
		order.ShippingAddress, order.BillingAddress,
		order.CreatedAt, order.UpdatedAt, order.CreatedBy, order.Version,
	)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

// FindByID finds an order by ID.
func (r *OrderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	query := `
		SELECT id, tenant_id, user_id, order_number, status, subtotal, tax_amount,
			   shipping_amount, discount_amount, total_amount, currency, notes,
			   shipping_address, billing_address, created_at, updated_at, deleted_at,
			   created_by, updated_by, version
		FROM orders
		WHERE id = $1 AND deleted_at IS NULL
	`
	var order domain.Order
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&order.ID, &order.TenantID, &order.UserID, &order.OrderNumber, &order.Status,
		&order.Subtotal, &order.TaxAmount, &order.ShippingAmount, &order.DiscountAmount,
		&order.TotalAmount, &order.Currency, &order.Notes,
		&order.ShippingAddress, &order.BillingAddress,
		&order.CreatedAt, &order.UpdatedAt, &order.DeletedAt,
		&order.CreatedBy, &order.UpdatedBy, &order.Version,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrOrderNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	return &order, nil
}

// FindByOrderNumber finds an order by order number.
func (r *OrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	query := `
		SELECT id, tenant_id, user_id, order_number, status, subtotal, tax_amount,
			   shipping_amount, discount_amount, total_amount, currency, notes,
			   shipping_address, billing_address, created_at, updated_at, deleted_at,
			   created_by, updated_by, version
		FROM orders
		WHERE order_number = $1 AND deleted_at IS NULL
	`
	var order domain.Order
	err := r.pool.QueryRow(ctx, query, orderNumber).Scan(
		&order.ID, &order.TenantID, &order.UserID, &order.OrderNumber, &order.Status,
		&order.Subtotal, &order.TaxAmount, &order.ShippingAmount, &order.DiscountAmount,
		&order.TotalAmount, &order.Currency, &order.Notes,
		&order.ShippingAddress, &order.BillingAddress,
		&order.CreatedAt, &order.UpdatedAt, &order.DeletedAt,
		&order.CreatedBy, &order.UpdatedBy, &order.Version,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrOrderNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	return &order, nil
}

// Update updates an order record.
func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	query := `
		UPDATE orders
		SET status = $1, notes = $2, shipping_address = $3, billing_address = $4,
		    updated_at = $5, updated_by = $6, version = version + 1
		WHERE id = $7 AND version = $8 AND deleted_at IS NULL
	`
	result, err := r.pool.Exec(ctx, query,
		order.Status, order.Notes, order.ShippingAddress, order.BillingAddress,
		order.UpdatedAt, order.UpdatedBy, order.ID, order.Version,
	)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrOrderNotFound
	}
	return nil
}

// UpdateStatus updates order status.
func (r *OrderRepository) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = NOW(), version = version + 1
		WHERE id = $2 AND deleted_at IS NULL
	`
	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrOrderNotFound
	}
	return nil
}

// ListByTenant lists orders by tenant with pagination.
func (r *OrderRepository) ListByTenant(ctx context.Context, tenantID string, offset, limit int) ([]*domain.Order, error) {
	query := `
		SELECT id, tenant_id, user_id, order_number, status, subtotal, tax_amount,
			   shipping_amount, discount_amount, total_amount, currency, notes,
			   shipping_address, billing_address, created_at, updated_at, deleted_at,
			   created_by, updated_by, version
		FROM orders
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID, &o.TenantID, &o.UserID, &o.OrderNumber, &o.Status,
			&o.Subtotal, &o.TaxAmount, &o.ShippingAmount, &o.DiscountAmount,
			&o.TotalAmount, &o.Currency, &o.Notes,
			&o.ShippingAddress, &o.BillingAddress,
			&o.CreatedAt, &o.UpdatedAt, &o.DeletedAt,
			&o.CreatedBy, &o.UpdatedBy, &o.Version,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &o)
	}
	return orders, nil
}

// ListByUser lists orders by user with pagination.
func (r *OrderRepository) ListByUser(ctx context.Context, userID string, offset, limit int) ([]*domain.Order, error) {
	query := `
		SELECT id, tenant_id, user_id, order_number, status, subtotal, tax_amount,
			   shipping_amount, discount_amount, total_amount, currency, notes,
			   shipping_address, billing_address, created_at, updated_at, deleted_at,
			   created_by, updated_by, version
		FROM orders
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID, &o.TenantID, &o.UserID, &o.OrderNumber, &o.Status,
			&o.Subtotal, &o.TaxAmount, &o.ShippingAmount, &o.DiscountAmount,
			&o.TotalAmount, &o.Currency, &o.Notes,
			&o.ShippingAddress, &o.BillingAddress,
			&o.CreatedAt, &o.UpdatedAt, &o.DeletedAt,
			&o.CreatedBy, &o.UpdatedBy, &o.Version,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &o)
	}
	return orders, nil
}

// Delete soft deletes an order record.
func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE orders SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrOrderNotFound
	}
	return nil
}
