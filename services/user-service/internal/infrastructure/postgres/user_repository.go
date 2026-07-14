package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cloudcommerce/user-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository implements domain.UserRepository with PostgreSQL.
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Create inserts a new user into the database.
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, phone, role, tenant_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	// Convert empty strings to nil for nullable columns
	var tenantID interface{}
	if user.TenantID != "" {
		tenantID = user.TenantID
	}
	var phone interface{}
	if user.Phone != "" {
		phone = user.Phone
	}

	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FullName, phone,
		string(user.Role), tenantID, user.IsActive, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrEmailAlreadyExists
		}
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

// FindByID retrieves a user by their ID.
func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, phone, role, tenant_id, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	user, err := r.scanUser(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return user, nil
}

// FindByEmail retrieves a user by their email address.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, phone, role, tenant_id, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	user, err := r.scanUser(ctx, query, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return user, nil
}

// Update updates an existing user.
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $2, password_hash = $3, full_name = $4, phone = $5, role = $6,
		    tenant_id = $7, is_active = $8, updated_at = $9
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.pool.Exec(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FullName, user.Phone,
		string(user.Role), user.TenantID, user.IsActive, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

// SoftDelete marks a user as deleted without removing the record.
func (r *UserRepository) SoftDelete(ctx context.Context, id string) error {
	query := `UPDATE users SET deleted_at = $2, updated_at = $2 WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("soft delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) scanUser(ctx context.Context, query string, args ...interface{}) (*domain.User, error) {
	row := r.pool.QueryRow(ctx, query, args...)

	var user domain.User
	var deletedAt sql.NullTime
	var tenantID sql.NullString
	var phone sql.NullString

	err := row.Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &phone,
		&user.Role, &tenantID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}

	if phone.Valid {
		user.Phone = phone.String
	}
	if tenantID.Valid {
		user.TenantID = tenantID.String
	}
	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	return &user, nil
}
