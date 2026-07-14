package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cloudcommerce/product-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProductRepository implements domain.ProductRepository with PostgreSQL.
type ProductRepository struct {
	pool *pgxpool.Pool
}

// NewProductRepository creates a new ProductRepository.
func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

// Create inserts a new product.
func (r *ProductRepository) Create(ctx context.Context, p *domain.Product) error {
	query := `
		INSERT INTO products (id, tenant_id, name, slug, description, sku, price, currency, status, category_id, image_url, weight, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.pool.Exec(ctx, query,
		p.ID, p.TenantID, p.Name, p.Slug, p.Description, p.SKU, p.Price, p.Currency,
		string(p.Status), p.CategoryID, p.ImageURL, p.Weight, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert product: %w", err)
	}
	return nil
}

// FindByID retrieves a product by ID within a tenant.
func (r *ProductRepository) FindByID(ctx context.Context, id, tenantID string) (*domain.Product, error) {
	query := `
		SELECT id, tenant_id, name, slug, description, sku, price, currency, status, category_id, image_url, weight, created_at, updated_at
		FROM products
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`
	row := r.pool.QueryRow(ctx, query, id, tenantID)
	p, err := scanProduct(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("find product: %w", err)
	}
	return p, nil
}

// FindBySKU retrieves a product by SKU within a tenant.
func (r *ProductRepository) FindBySKU(ctx context.Context, sku, tenantID string) (*domain.Product, error) {
	query := `
		SELECT id, tenant_id, name, slug, description, sku, price, currency, status, category_id, image_url, weight, created_at, updated_at
		FROM products
		WHERE sku = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`
	row := r.pool.QueryRow(ctx, query, sku, tenantID)
	p, err := scanProduct(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("find product by sku: %w", err)
	}
	return p, nil
}

// List retrieves a paginated list of products for a tenant.
func (r *ProductRepository) List(ctx context.Context, tenantID string, limit, offset int, status string) ([]*domain.Product, int64, error) {
	countQuery := `SELECT COUNT(*) FROM products WHERE tenant_id = $1 AND deleted_at IS NULL`
	listQuery := `
		SELECT id, tenant_id, name, slug, description, sku, price, currency, status, category_id, image_url, weight, created_at, updated_at
		FROM products
		WHERE tenant_id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{tenantID}

	if status != "" {
		countQuery += ` AND status = $2`
		listQuery += ` AND status = $2`
		args = append(args, status)
	}

	listQuery += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, p)
	}

	return products, total, nil
}

// Update updates an existing product.
func (r *ProductRepository) Update(ctx context.Context, p *domain.Product) error {
	query := `
		UPDATE products
		SET name = $3, slug = $4, description = $5, sku = $6, price = $7, currency = $8,
		    status = $9, category_id = $10, image_url = $11, weight = $12, updated_at = $13
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`
	result, err := r.pool.Exec(ctx, query,
		p.ID, p.TenantID, p.Name, p.Slug, p.Description, p.SKU, p.Price, p.Currency,
		string(p.Status), p.CategoryID, p.ImageURL, p.Weight, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

// SoftDelete marks a product as deleted.
func (r *ProductRepository) SoftDelete(ctx context.Context, id, tenantID string) error {
	query := `UPDATE products SET deleted_at = $3, updated_at = $3 WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL`
	result, err := r.pool.Exec(ctx, query, id, tenantID, time.Now())
	if err != nil {
		return fmt.Errorf("soft delete product: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

// scanner is an interface for both pgx.Row and pgx.Rows.
type scanner interface {
	Scan(dest ...interface{}) error
}

func scanProduct(s scanner) (*domain.Product, error) {
	var p domain.Product
	var categoryID sql.NullString
	var imageURL sql.NullString
	var description sql.NullString

	err := s.Scan(
		&p.ID, &p.TenantID, &p.Name, &p.Slug, &description, &p.SKU, &p.Price, &p.Currency,
		&p.Status, &categoryID, &imageURL, &p.Weight, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		p.Description = description.String
	}
	if categoryID.Valid {
		p.CategoryID = &categoryID.String
	}
	if imageURL.Valid {
		p.ImageURL = imageURL.String
	}

	return &p, nil
}
