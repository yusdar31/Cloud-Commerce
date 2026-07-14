-- +goose Up
-- +goose StatementBegin

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create inventory table
CREATE TABLE IF NOT EXISTS inventory (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id   UUID NOT NULL,
    product_id  UUID NOT NULL,
    sku         VARCHAR(100) NOT NULL,
    quantity    INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reserved    INTEGER NOT NULL DEFAULT 0 CHECK (reserved >= 0),
    available   INTEGER GENERATED ALWAYS AS (quantity - reserved) STORED,
    low_stock_threshold INTEGER NOT NULL DEFAULT 10,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    created_by  UUID,
    updated_by  UUID,
    version     INTEGER NOT NULL DEFAULT 1,
    UNIQUE (tenant_id, product_id),
    UNIQUE (tenant_id, sku)
);

CREATE INDEX IF NOT EXISTS idx_inventory_tenant_id ON inventory(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_inventory_product_id ON inventory(product_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_inventory_sku ON inventory(sku) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_inventory_low_stock ON inventory(tenant_id, available) WHERE available <= low_stock_threshold AND deleted_at IS NULL;

-- Create reservations table
CREATE TABLE IF NOT EXISTS reservations (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id      UUID NOT NULL,
    order_id       UUID NOT NULL,
    product_id     UUID NOT NULL,
    sku            VARCHAR(100) NOT NULL,
    quantity       INTEGER NOT NULL CHECK (quantity > 0),
    status         VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'released', 'fulfilled')),
    expires_at     TIMESTAMPTZ NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (order_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_reservations_tenant_id ON reservations(tenant_id);
CREATE INDEX IF NOT EXISTS idx_reservations_order_id ON reservations(order_id);
CREATE INDEX IF NOT EXISTS idx_reservations_product_id ON reservations(product_id);
CREATE INDEX IF NOT EXISTS idx_reservations_status ON reservations(status);
CREATE INDEX IF NOT EXISTS idx_reservations_expires_at ON reservations(expires_at) WHERE status = 'active';

-- Create stock_movements table (audit trail)
CREATE TABLE IF NOT EXISTS stock_movements (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id   UUID NOT NULL,
    product_id  UUID NOT NULL,
    sku         VARCHAR(100) NOT NULL,
    movement_type VARCHAR(50) NOT NULL, -- 'initial', 'adjustment', 'reservation', 'release', 'sale', 'return'
    quantity_before INTEGER NOT NULL,
    quantity_after  INTEGER NOT NULL,
    quantity_change INTEGER NOT NULL,
    reference_id    UUID, -- order_id, reservation_id, etc
    reference_type  VARCHAR(50), -- 'order', 'reservation', 'manual'
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  UUID
);

CREATE INDEX IF NOT EXISTS idx_stock_movements_tenant_id ON stock_movements(tenant_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_product_id ON stock_movements(product_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_reference ON stock_movements(reference_id, reference_type);
CREATE INDEX IF NOT EXISTS idx_stock_movements_created_at ON stock_movements(created_at DESC);

-- Function to auto-expire reservations
CREATE OR REPLACE FUNCTION expire_old_reservations()
RETURNS void AS $$
BEGIN
    UPDATE reservations
    SET status = 'released',
        updated_at = NOW()
    WHERE status = 'active'
      AND expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd

-- +goose Down
DROP FUNCTION IF EXISTS expire_old_reservations;
DROP TABLE IF EXISTS stock_movements;
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS inventory;
