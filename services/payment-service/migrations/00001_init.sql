-- +goose Up
-- +goose StatementBegin

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id        UUID NOT NULL,
    order_id         UUID NOT NULL UNIQUE,
    user_id          UUID NOT NULL,
    amount           BIGINT NOT NULL CHECK (amount > 0),
    currency         VARCHAR(3) NOT NULL DEFAULT 'IDR',
    status           VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'paid', 'failed', 'refunded', 'cancelled')),
    payment_method   VARCHAR(50),
    payment_provider  VARCHAR(50),
    provider_transaction_id VARCHAR(255),
    provider_payment_id     VARCHAR(255),
    payment_url      TEXT,
    expires_at       TIMESTAMPTZ,
    paid_at          TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTAMPS NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ,
    created_by       UUID,
    updated_by       UUID,
    version          INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_payments_tenant_id ON payments(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_payments_provider_transaction ON payments(provider_transaction_id);
CREATE INDEX IF NOT EXISTS idx_payments_created_at ON payments(created_at DESC);

-- Create payment_transactions table (audit trail)
CREATE TABLE IF NOT EXISTS payment_transactions (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id        UUID NOT NULL,
    payment_id       UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    transaction_type VARCHAR(50) NOT NULL, -- 'payment', 'refund', 'chargeback'
    status           VARCHAR(20) NOT NULL,
    amount           BIGINT NOT NULL,
    currency         VARCHAR(3) NOT NULL DEFAULT 'IDR',
    provider_response JSONB,
    error_message    TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payment_transactions_payment_id ON payment_transactions(payment_id);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_created_at ON payment_transactions(created_at DESC);

-- Create webhook_events table (idempotency)
CREATE TABLE IF NOT EXISTS webhook_events (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id        UUID NOT NULL,
    event_id         VARCHAR(255) NOT NULL UNIQUE,
    event_type       VARCHAR(100) NOT NULL,
    provider         VARCHAR(50) NOT NULL,
    payload          JSONB NOT NULL,
    signature        TEXT,
    processed        BOOLEAN NOT NULL DEFAULT FALSE,
    processed_at     TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_webhook_events_event_id ON webhook_events(event_id);
CREATE INDEX IF NOT EXISTS idx_webhook_events_processed ON webhook_events(processed) WHERE processed = FALSE;
CREATE INDEX IF NOT EXISTS idx_webhook_events_created_at ON webhook_events(created_at DESC);

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS webhook_events;
DROP TABLE IF EXISTS payment_transactions;
DROP TABLE IF EXISTS payments;
