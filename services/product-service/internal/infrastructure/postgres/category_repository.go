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

// CategoryRepository implements domain.CategoryRepository with PostgreSQL.
type CategoryRepository struct {
	pool *pgxpool.Pool
}

// NewCategoryRepository creates a new CategoryRepository.
func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

// Create inserts a new category.
func (r *CategoryRepository) Create(ctx context.Context, c *domain.Category) error {
	query := `
		INSERT INTO categories (id, tenant_id, name, slug, description, parent_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query,
		c.ID, c.TenantID, c.Name, c.Slug, c.Description, c.ParentID, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert category: %w", err)
	}
	return nil
}

// FindByID retrieves a category by ID within a tenant.
func (r *CategoryRepository) FindByID(ctx context.Context, id, tenantID string) (*domain.Category, error) {
	query := `
		SELECT id, tenant_id, name, slug, description, parent_id, created_at, updated_at
		FROM categories
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`
	row := r.pool.QueryRow(ctx, query, id, tenantID)

	var c domain.Category
	var description sql.NullString
	var parentID sql.NullString

	err := row.Scan(&c.ID, &c.TenantID, &c.Name, &c.Slug, &description, &parentID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("find category: %w", err)
	}

	if description.Valid {
		c.Description = description.String
	}
	if parentID.Valid {
		c.ParentID = &parentID.String
	}

	return &c, nil
}

// List retrieves all categories for a tenant.
func (r *CategoryRepository) List(ctx context.Context, tenantID string) ([]*domain.Category, error) {
	query := `
		SELECT id, tenant_id, name, slug, description, parent_id, created_at, updated_at
		FROM categories
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY name ASC
	`
	rows, err := r.pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("query categories: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var c domain.Category
		var description sql.NullString
		var parentID sql.NullString

		if err := rows.Scan(&c.ID, &c.TenantID, &c.Name, &c.Slug, &description, &parentID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}

		if description.Valid {
			c.Description = description.String
		}
		if parentID.Valid {
			c.ParentID = &parentID.String
		}

		categories = append(categories, &c)
	}

	return categories, nil
}

// Update updates an existing category.
func (r *CategoryRepository) Update(ctx context.Context, c *domain.Category) error {
	query := `
		UPDATE categories
		SET name = $3, slug = $4, description = $5, parent_id = $6, updated_at = $7
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`
	result, err := r.pool.Exec(ctx, query,
		c.ID, c.TenantID, c.Name, c.Slug, c.Description, c.ParentID, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrCategoryNotFound
	}
	return nil
}

// SoftDelete marks a category as deleted.
func (r *CategoryRepository) SoftDelete(ctx context.Context, id, tenantID string) error {
	query := `UPDATE categories SET deleted_at = $3, updated_at = $3 WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL`
	result, err := r.pool.Exec(ctx, query, id, tenantID, time.Now())
	if err != nil {
		return fmt.Errorf("soft delete category: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrCategoryNotFound
	}
	return nil
}
