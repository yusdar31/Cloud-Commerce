package postgres

import (
	"context"
	"fmt"

	"github.com/cloudcommerce/inventory-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InventoryRepository is a PostgreSQL implementation of domain.InventoryRepository.
type InventoryRepository struct {
	pool *pgxpool.Pool
}

// NewInventoryRepository creates a new inventory repository.
func NewInventoryRepository(pool *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{pool: pool}
}

// Create creates a new inventory record.
func (r *InventoryRepository) Create(ctx context.Context, inv *domain.Inventory) error {
	query := `
		INSERT INTO inventory (
			id, tenant_id, product_id, sku, quantity, reserved,
			low_stock_threshold, created_at, updated_at, created_by, version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.pool.Exec(ctx, query,
		inv.ID, inv.TenantID, inv.ProductID, inv.SKU, inv.Quantity, inv.Reserved,
		inv.LowStockThreshold, inv.CreatedAt, inv.UpdatedAt, inv.CreatedBy, inv.Version,
	)
	if err != nil {
		return fmt.Errorf("failed to create inventory: %w", err)
	}
	return nil
}

// FindByID finds inventory by ID.
func (r *InventoryRepository) FindByID(ctx context.Context, id string) (*domain.Inventory, error) {
	query := `
		SELECT id, tenant_id, product_id, sku, quantity, reserved, available,
			   low_stock_threshold, created_at, updated_at, deleted_at,
			   created_by, updated_by, version
		FROM inventory
		WHERE id = $1 AND deleted_at IS NULL
	`
	var inv domain.Inventory
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&inv.ID, &inv.TenantID, &inv.ProductID, &inv.SKU, &inv.Quantity, &inv.Reserved, &inv.Available,
		&inv.LowStockThreshold, &inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt,
		&inv.CreatedBy, &inv.UpdatedBy, &inv.Version,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrInventoryNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory: %w", err)
	}
	return &inv, nil
}

// FindByProductID finds inventory by product ID.
func (r *InventoryRepository) FindByProductID(ctx context.Context, tenantID, productID string) (*domain.Inventory, error) {
	query := `
		SELECT id, tenant_id, product_id, sku, quantity, reserved, available,
			   low_stock_threshold, created_at, updated_at, deleted_at,
			   created_by, updated_by, version
		FROM inventory
		WHERE tenant_id = $1 AND product_id = $2 AND deleted_at IS NULL
	`
	var inv domain.Inventory
	err := r.pool.QueryRow(ctx, query, tenantID, productID).Scan(
		&inv.ID, &inv.TenantID, &inv.ProductID, &inv.SKU, &inv.Quantity, &inv.Reserved, &inv.Available,
		&inv.LowStockThreshold, &inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt,
		&inv.CreatedBy, &inv.UpdatedBy, &inv.Version,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrInventoryNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory: %w", err)
	}
	return &inv, nil
}

// FindBySKU finds inventory by SKU.
func (r *InventoryRepository) FindBySKU(ctx context.Context, tenantID, sku string) (*domain.Inventory, error) {
	query := `
		SELECT id, tenant_id, product_id, sku, quantity, reserved, available,
			   low_stock_threshold, created_at, updated_at, deleted_at,
			   created_by, updated_by, version
		FROM inventory
		WHERE tenant_id = $1 AND sku = $2 AND deleted_at IS NULL
	`
	var inv domain.Inventory
	err := r.pool.QueryRow(ctx, query, tenantID, sku).Scan(
		&inv.ID, &inv.TenantID, &inv.ProductID, &inv.SKU, &inv.Quantity, &inv.Reserved, &inv.Available,
		&inv.LowStockThreshold, &inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt,
		&inv.CreatedBy, &inv.UpdatedBy, &inv.Version,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrInventoryNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory: %w", err)
	}
	return &inv, nil
}

// List lists inventory with pagination.
func (r *InventoryRepository) List(ctx context.Context, tenantID string, offset, limit int) ([]*domain.Inventory, error) {
	query := `
		SELECT id, tenant_id, product_id, sku, quantity, reserved, available,
			   low_stock_threshold, created_at, updated_at, deleted_at,
			   created_by, updated_by, version
		FROM inventory
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory: %w", err)
	}
	defer rows.Close()

	var inventories []*domain.Inventory
	for rows.Next() {
		var inv domain.Inventory
		if err := rows.Scan(
			&inv.ID, &inv.TenantID, &inv.ProductID, &inv.SKU, &inv.Quantity, &inv.Reserved, &inv.Available,
			&inv.LowStockThreshold, &inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt,
			&inv.CreatedBy, &inv.UpdatedBy, &inv.Version,
		); err != nil {
			return nil, fmt.Errorf("failed to scan inventory: %w", err)
		}
		inventories = append(inventories, &inv)
	}

	return inventories, nil
}

// Update updates an inventory record.
func (r *InventoryRepository) Update(ctx context.Context, inv *domain.Inventory) error {
	query := `
		UPDATE inventory
		SET quantity = $1, reserved = $2, low_stock_threshold = $3,
		    updated_at = $4, updated_by = $5, version = version + 1
		WHERE id = $6 AND version = $7 AND deleted_at IS NULL
	`
	result, err := r.pool.Exec(ctx, query,
		inv.Quantity, inv.Reserved, inv.LowStockThreshold,
		inv.UpdatedAt, inv.UpdatedBy, inv.ID, inv.Version,
	)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrInventoryNotFound
	}
	return nil
}

// UpdateQuantity updates quantity only.
func (r *InventoryRepository) UpdateQuantity(ctx context.Context, id string, quantity int) error {
	query := `
		UPDATE inventory
		SET quantity = $1, updated_at = NOW(), version = version + 1
		WHERE id = $2 AND deleted_at IS NULL
	`
	result, err := r.pool.Exec(ctx, query, quantity, id)
	if err != nil {
		return fmt.Errorf("failed to update quantity: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrInventoryNotFound
	}
	return nil
}

// Reserve reserves stock (increase reserved count).
func (r *InventoryRepository) Reserve(ctx context.Context, tenantID, productID string, quantity int) error {
	query := `
		UPDATE inventory
		SET reserved = reserved + $1, updated_at = NOW(), version = version + 1
		WHERE tenant_id = $2 AND product_id = $3 AND (quantity - reserved) >= $1 AND deleted_at IS NULL
	`
	result, err := r.pool.Exec(ctx, query, quantity, tenantID, productID)
	if err != nil {
		return fmt.Errorf("failed to reserve stock: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrInsufficientStock
	}
	return nil
}

// Release releases reserved stock (decrease reserved count).
func (r *InventoryRepository) Release(ctx context.Context, tenantID, productID string, quantity int) error {
	query := `
		UPDATE inventory
		SET reserved = reserved - $1, updated_at = NOW(), version = version + 1
		WHERE tenant_id = $2 AND product_id = $3 AND reserved >= $1 AND deleted_at IS NULL
	`
	result, err := r.pool.Exec(ctx, query, quantity, tenantID, productID)
	if err != nil {
		return fmt.Errorf("failed to release stock: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrInventoryNotFound
	}
	return nil
}

// FindLowStock finds inventory items below threshold.
func (r *InventoryRepository) FindLowStock(ctx context.Context, tenantID string) ([]*domain.Inventory, error) {
	query := `
		SELECT id, tenant_id, product_id, sku, quantity, reserved, available,
			   low_stock_threshold, created_at, updated_at, deleted_at,
			   created_by, updated_by, version
		FROM inventory
		WHERE tenant_id = $1 AND available <= low_stock_threshold AND deleted_at IS NULL
		ORDER BY available ASC
	`
	rows, err := r.pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to find low stock: %w", err)
	}
	defer rows.Close()

	var inventories []*domain.Inventory
	for rows.Next() {
		var inv domain.Inventory
		if err := rows.Scan(
			&inv.ID, &inv.TenantID, &inv.ProductID, &inv.SKU, &inv.Quantity, &inv.Reserved, &inv.Available,
			&inv.LowStockThreshold, &inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt,
			&inv.CreatedBy, &inv.UpdatedBy, &inv.Version,
		); err != nil {
			return nil, fmt.Errorf("failed to scan inventory: %w", err)
		}
		inventories = append(inventories, &inv)
	}

	return inventories, nil
}

// Delete soft deletes an inventory record.
func (r *InventoryRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE inventory
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete inventory: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrInventoryNotFound
	}
	return nil
}
