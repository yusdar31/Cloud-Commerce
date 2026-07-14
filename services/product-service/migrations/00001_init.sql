-- +goose Up
-- +goose StatementBegin

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id   UUID NOT NULL,
    name        VARCHAR(100) NOT NULL,
    slug        VARCHAR(150) NOT NULL,
    description TEXT,
    parent_id   UUID REFERENCES categories(id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    created_by  UUID,
    updated_by  UUID,
    version     INTEGER NOT NULL DEFAULT 1,
    UNIQUE (tenant_id, slug)
);

CREATE INDEX IF NOT EXISTS idx_categories_tenant_id ON categories(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id) WHERE deleted_at IS NULL;

-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id   UUID NOT NULL,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(300) NOT NULL,
    description TEXT,
    sku         VARCHAR(100) NOT NULL,
    price       BIGINT NOT NULL CHECK (price > 0),
    currency    VARCHAR(3) NOT NULL DEFAULT 'IDR',
    status      VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    image_url   TEXT,
    weight      INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    created_by  UUID,
    updated_by  UUID,
    version     INTEGER NOT NULL DEFAULT 1,
    UNIQUE (tenant_id, sku)
);

CREATE INDEX IF NOT EXISTS idx_products_tenant_id ON products(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_products_status ON products(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_products_slug ON products(slug) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at);
CREATE INDEX IF NOT EXISTS idx_products_tenant_status ON products(tenant_id, status) WHERE deleted_at IS NULL;

-- Create product_images table
CREATE TABLE IF NOT EXISTS product_images (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url        TEXT NOT NULL,
    alt_text   VARCHAR(255),
    position   INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_product_images_product_id ON product_images(product_id);

-- Create product_tags table
CREATE TABLE IF NOT EXISTS product_tags (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    tag        VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (product_id, tag)
);

CREATE INDEX IF NOT EXISTS idx_product_tags_product_id ON product_tags(product_id);
CREATE INDEX IF NOT EXISTS idx_product_tags_tag ON product_tags(tag);

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS product_tags;
DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
