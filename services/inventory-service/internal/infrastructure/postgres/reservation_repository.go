package postgres

import (
	"context"
	"fmt"

	"github.com/cloudcommerce/inventory-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ReservationRepository is a PostgreSQL implementation of domain.ReservationRepository.
type ReservationRepository struct {
	pool *pgxpool.Pool
}

// NewReservationRepository creates a new reservation repository.
func NewReservationRepository(pool *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{pool: pool}
}

// Create creates a new reservation.
func (r *ReservationRepository) Create(ctx context.Context, res *domain.Reservation) error {
	query := `
		INSERT INTO reservations (
			id, tenant_id, order_id, product_id, sku, quantity,
			status, expires_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(ctx, query,
		res.ID, res.TenantID, res.OrderID, res.ProductID, res.SKU, res.Quantity,
		res.Status, res.ExpiresAt, res.CreatedAt, res.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}
	return nil
}

// FindByID finds a reservation by ID.
func (r *ReservationRepository) FindByID(ctx context.Context, id string) (*domain.Reservation, error) {
	query := `
		SELECT id, tenant_id, order_id, product_id, sku, quantity,
		       status, expires_at, created_at, updated_at
		FROM reservations
		WHERE id = $1
	`
	var res domain.Reservation
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&res.ID, &res.TenantID, &res.OrderID, &res.ProductID, &res.SKU, &res.Quantity,
		&res.Status, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrReservationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find reservation: %w", err)
	}
	return &res, nil
}

// FindByOrderID finds all reservations for an order.
func (r *ReservationRepository) FindByOrderID(ctx context.Context, orderID string) ([]*domain.Reservation, error) {
	query := `
		SELECT id, tenant_id, order_id, product_id, sku, quantity,
		       status, expires_at, created_at, updated_at
		FROM reservations
		WHERE order_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find reservations: %w", err)
	}
	defer rows.Close()

	var reservations []*domain.Reservation
	for rows.Next() {
		var res domain.Reservation
		if err := rows.Scan(
			&res.ID, &res.TenantID, &res.OrderID, &res.ProductID, &res.SKU, &res.Quantity,
			&res.Status, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan reservation: %w", err)
		}
		reservations = append(reservations, &res)
	}

	return reservations, nil
}

// FindActiveByProduct finds active reservations for a product.
func (r *ReservationRepository) FindActiveByProduct(ctx context.Context, tenantID, productID string) ([]*domain.Reservation, error) {
	query := `
		SELECT id, tenant_id, order_id, product_id, sku, quantity,
		       status, expires_at, created_at, updated_at
		FROM reservations
		WHERE tenant_id = $1 AND product_id = $2 AND status = $3
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, tenantID, productID, domain.ReservationActive)
	if err != nil {
		return nil, fmt.Errorf("failed to find active reservations: %w", err)
	}
	defer rows.Close()

	var reservations []*domain.Reservation
	for rows.Next() {
		var res domain.Reservation
		if err := rows.Scan(
			&res.ID, &res.TenantID, &res.OrderID, &res.ProductID, &res.SKU, &res.Quantity,
			&res.Status, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan reservation: %w", err)
		}
		reservations = append(reservations, &res)
	}

	return reservations, nil
}

// UpdateStatus updates reservation status.
func (r *ReservationRepository) UpdateStatus(ctx context.Context, id string, status domain.ReservationStatus) error {
	query := `
		UPDATE reservations
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update reservation status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrReservationNotFound
	}
	return nil
}

// ReleaseExpired releases all expired reservations.
func (r *ReservationRepository) ReleaseExpired(ctx context.Context) (int, error) {
	query := `
		UPDATE reservations
		SET status = $1, updated_at = NOW()
		WHERE status = $2 AND expires_at < NOW()
	`
	result, err := r.pool.Exec(ctx, query, domain.ReservationReleased, domain.ReservationActive)
	if err != nil {
		return 0, fmt.Errorf("failed to release expired reservations: %w", err)
	}
	return int(result.RowsAffected()), nil
}

// Delete deletes a reservation.
func (r *ReservationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM reservations WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrReservationNotFound
	}
	return nil
}
