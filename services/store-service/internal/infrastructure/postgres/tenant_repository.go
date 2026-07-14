package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cloudcommerce/store-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TenantRepository implements domain.TenantRepository with PostgreSQL.
type TenantRepository struct {
	pool *pgxpool.Pool
}

// NewTenantRepository creates a new TenantRepository.
func NewTenantRepository(pool *pgxpool.Pool) *TenantRepository {
	return &TenantRepository{pool: pool}
}

// Create inserts a new tenant into the database.
func (r *TenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		INSERT INTO tenants (id, name, slug, owner_id, description, logo_url, banner_url, status, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	// Convert empty strings to nil for nullable columns
	var description interface{}
	if tenant.Description != "" {
		description = tenant.Description
	}
	var logoURL interface{}
	if tenant.LogoURL != "" {
		logoURL = tenant.LogoURL
	}
	var bannerURL interface{}
	if tenant.BannerURL != "" {
		bannerURL = tenant.BannerURL
	}

	_, err := r.pool.Exec(ctx, query,
		tenant.ID, tenant.Name, tenant.Slug, tenant.OwnerID, description,
		logoURL, bannerURL, tenant.Status, tenant.CreatedAt, tenant.UpdatedAt, tenant.Version,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			if pgErr.ConstraintName == "tenants_slug_key" {
				return domain.ErrSlugAlreadyExists
			}
			if pgErr.ConstraintName == "tenants_owner_id_key" {
				return domain.ErrOwnerAlreadyHasStore
			}
			return domain.ErrSlugAlreadyExists // default fallback
		}
		return fmt.Errorf("insert tenant: %w", err)
	}
	return nil
}

// FindByID retrieves a tenant by their ID.
func (r *TenantRepository) FindByID(ctx context.Context, id string) (*domain.Tenant, error) {
	query := `
		SELECT id, name, slug, owner_id, description, logo_url, banner_url, status, created_at, updated_at, deleted_at, version
		FROM tenants
		WHERE id = $1 AND deleted_at IS NULL
	`
	tenant, err := r.scanTenant(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTenantNotFound
		}
		return nil, fmt.Errorf("find tenant by id: %w", err)
	}
	return tenant, nil
}

// FindBySlug retrieves a tenant by their slug.
func (r *TenantRepository) FindBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	query := `
		SELECT id, name, slug, owner_id, description, logo_url, banner_url, status, created_at, updated_at, deleted_at, version
		FROM tenants
		WHERE slug = $1 AND deleted_at IS NULL
	`
	tenant, err := r.scanTenant(ctx, query, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTenantNotFound
		}
		return nil, fmt.Errorf("find tenant by slug: %w", err)
	}
	return tenant, nil
}

// FindByOwnerID retrieves a tenant by their owner's ID.
func (r *TenantRepository) FindByOwnerID(ctx context.Context, ownerID string) (*domain.Tenant, error) {
	query := `
		SELECT id, name, slug, owner_id, description, logo_url, banner_url, status, created_at, updated_at, deleted_at, version
		FROM tenants
		WHERE owner_id = $1 AND deleted_at IS NULL
	`
	tenant, err := r.scanTenant(ctx, query, ownerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTenantNotFound
		}
		return nil, fmt.Errorf("find tenant by owner_id: %w", err)
	}
	return tenant, nil
}

// Update updates an existing tenant.
func (r *TenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		UPDATE tenants
		SET name = $2, slug = $3, description = $4, logo_url = $5, banner_url = $6, status = $7, updated_at = $8, version = version + 1
		WHERE id = $1 AND deleted_at IS NULL
	`
	var description interface{}
	if tenant.Description != "" {
		description = tenant.Description
	}
	var logoURL interface{}
	if tenant.LogoURL != "" {
		logoURL = tenant.LogoURL
	}
	var bannerURL interface{}
	if tenant.BannerURL != "" {
		bannerURL = tenant.BannerURL
	}

	result, err := r.pool.Exec(ctx, query,
		tenant.ID, tenant.Name, tenant.Slug, description, logoURL, bannerURL, tenant.Status, time.Now(),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			if pgErr.ConstraintName == "tenants_slug_key" {
				return domain.ErrSlugAlreadyExists
			}
		}
		return fmt.Errorf("update tenant: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrTenantNotFound
	}
	return nil
}

func (r *TenantRepository) scanTenant(ctx context.Context, query string, args ...interface{}) (*domain.Tenant, error) {
	row := r.pool.QueryRow(ctx, query, args...)

	var tenant domain.Tenant
	var deletedAt sql.NullTime
	var description sql.NullString
	var logoURL sql.NullString
	var bannerURL sql.NullString

	err := row.Scan(
		&tenant.ID, &tenant.Name, &tenant.Slug, &tenant.OwnerID, &description,
		&logoURL, &bannerURL, &tenant.Status, &tenant.CreatedAt, &tenant.UpdatedAt, &deletedAt, &tenant.Version,
	)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		tenant.Description = description.String
	}
	if logoURL.Valid {
		tenant.LogoURL = logoURL.String
	}
	if bannerURL.Valid {
		tenant.BannerURL = bannerURL.String
	}
	if deletedAt.Valid {
		tenant.DeletedAt = &deletedAt.Time
	}

	return &tenant, nil
}
