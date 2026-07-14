package postgres

import (
	"context"
	"fmt"

	"github.com/cloudcommerce/payment-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PaymentRepository is a PostgreSQL implementation of domain.PaymentRepository.
type PaymentRepository struct {
	pool *pgxpool.Pool
}

// NewPaymentRepository creates a new payment repository.
func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{pool: pool}
}

// Create creates a new payment record.
func (r *PaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (
			id, tenant_id, order_id, user_id, amount, currency, status,
			payment_method, payment_provider, provider_transaction_id,
			provider_payment_id, payment_url, expires_at, paid_at,
			created_at, updated_at, created_by, version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`
	_, err := r.pool.Exec(ctx, query,
		payment.ID, payment.TenantID, payment.OrderID, payment.UserID,
		payment.Amount, payment.Currency, payment.Status,
		payment.PaymentMethod, payment.PaymentProvider, payment.ProviderTransactionID,
		payment.ProviderPaymentID, payment.PaymentURL, payment.ExpiresAt, payment.PaidAt,
		payment.CreatedAt, payment.UpdatedAt, payment.CreatedBy, payment.Version,
	)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	return nil
}

// FindByID finds a payment by ID.
func (r *PaymentRepository) FindByID(ctx context.Context, id string) (*domain.Payment, error) {
	query := `
		SELECT id, tenant_id, order_id, user_id, amount, currency, status,
			   payment_method, payment_provider, provider_transaction_id,
			   provider_payment_id, payment_url, expires_at, paid_at,
			   created_at, updated_at, deleted_at, created_by, updated_by, version
		FROM payments
		WHERE id = $1 AND deleted_at IS NULL
	`
	var payment domain.Payment
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&payment.ID, &payment.TenantID, &payment.OrderID, &payment.UserID,
		&payment.Amount, &payment.Currency, &payment.Status,
		&payment.PaymentMethod, &payment.PaymentProvider, &payment.ProviderTransactionID,
		&payment.ProviderPaymentID, &payment.PaymentURL, &payment.ExpiresAt, &payment.PaidAt,
		&payment.CreatedAt, &payment.UpdatedAt, &payment.DeletedAt,
		&payment.CreatedBy, &payment.UpdatedBy, &payment.Version,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrPaymentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find payment: %w", err)
	}
	return &payment, nil
}

// FindByOrderID finds a payment by order ID.
func (r *PaymentRepository) FindByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	query := `
		SELECT id, tenant_id, order_id, user_id, amount, currency, status,
			   payment_method, payment_provider, provider_transaction_id,
			   provider_payment_id, payment_url, expires_at, paid_at,
			   created_at, updated_at, deleted_at, created_by, updated_by, version
		FROM payments
		WHERE order_id = $1 AND deleted_at IS NULL
	`
	var payment domain.Payment
	err := r.pool.QueryRow(ctx, query, orderID).Scan(
		&payment.ID, &payment.TenantID, &payment.OrderID, &payment.UserID,
		&payment.Amount, &payment.Currency, &payment.Status,
		&payment.PaymentMethod, &payment.PaymentProvider, &payment.ProviderTransactionID,
		&payment.ProviderPaymentID, &payment.PaymentURL, &payment.ExpiresAt, &payment.PaidAt,
		&payment.CreatedAt, &payment.UpdatedAt, &payment.DeletedAt,
		&payment.CreatedBy, &payment.UpdatedBy, &payment.Version,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrPaymentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find payment: %w", err)
	}
	return &payment, nil
}

// Update updates a payment record.
func (r *PaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	query := `
		UPDATE payments
		SET status = $1, payment_method = $2, payment_provider = $3,
		    provider_transaction_id = $4, provider_payment_id = $5,
		    payment_url = $6, expires_at = $7, paid_at = $8,
		    updated_at = $9, updated_by = $10, version = version + 1
		WHERE id = $11 AND version = $12 AND deleted_at IS NULL
	`
	result, err := r.pool.Exec(ctx, query,
		payment.Status, payment.PaymentMethod, payment.PaymentProvider,
		payment.ProviderTransactionID, payment.ProviderPaymentID,
		payment.PaymentURL, payment.ExpiresAt, payment.PaidAt,
		payment.UpdatedAt, payment.UpdatedBy, payment.ID, payment.Version,
	)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrPaymentNotFound
	}
	return nil
}

// UpdateStatus updates payment status.
func (r *PaymentRepository) UpdateStatus(ctx context.Context, id string, status domain.PaymentStatus, paidAt *time.Time) error {
	query := `
		UPDATE payments
		SET status = $1, paid_at = $2, updated_at = NOW(), version = version + 1
		WHERE id = $3 AND deleted_at IS NULL
	`
	result, err := r.pool.Exec(ctx, query, status, paidAt, id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrPaymentNotFound
	}
	return nil
}

// ListByTenant lists payments by tenant with pagination.
func (r *PaymentRepository) ListByTenant(ctx context.Context, tenantID string, offset, limit int) ([]*domain.Payment, error) {
	query := `
		SELECT id, tenant_id, order_id, user_id, amount, currency, status,
			   payment_method, payment_provider, provider_transaction_id,
			   provider_payment_id, payment_url, expires_at, paid_at,
			   created_at, updated_at, deleted_at, created_by, updated_by, version
		FROM payments
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(
			&p.ID, &p.TenantID, &p.OrderID, &p.UserID,
			&p.Amount, &p.Currency, &p.Status,
			&p.PaymentMethod, &p.PaymentProvider, &p.ProviderTransactionID,
			&p.ProviderPaymentID, &p.PaymentURL, &p.ExpiresAt, &p.PaidAt,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
			&p.CreatedBy, &p.UpdatedBy, &p.Version,
		); err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, &p)
	}

	return payments, nil
}

// ListByUser lists payments by user with pagination.
func (r *PaymentRepository) ListByUser(ctx context.Context, userID string, offset, limit int) ([]*domain.Payment, error) {
	query := `
		SELECT id, tenant_id, order_id, user_id, amount, currency, status,
			   payment_method, payment_provider, provider_transaction_id,
			   provider_payment_id, payment_url, expires_at, paid_at,
			   created_at, updated_at, deleted_at, created_by, updated_by, version
		FROM payments
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(
			&p.ID, &p.TenantID, &p.OrderID, &p.UserID,
			&p.Amount, &p.Currency, &p.Status,
			&p.PaymentMethod, &p.PaymentProvider, &p.ProviderTransactionID,
			&p.ProviderPaymentID, &p.PaymentURL, &p.ExpiresAt, &p.PaidAt,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
			&p.CreatedBy, &p.UpdatedBy, &p.Version,
		); err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, &p)
	}

	return payments, nil
}

// Delete soft deletes a payment record.
func (r *PaymentRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE payments SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrPaymentNotFound
	}
	return nil
}
