package domain

import (
	"errors"
	"time"
)

// Common domain errors.
var (
	ErrProductNotFound  = errors.New("product not found")
	ErrCategoryNotFound = errors.New("category not found")
	ErrSKUExists        = errors.New("sku already exists")
	ErrSlugExists       = errors.New("slug already exists")
)

// ProductStatus represents the publish state of a product.
type ProductStatus string

const (
	StatusDraft     ProductStatus = "draft"
	StatusPublished ProductStatus = "published"
	StatusArchived  ProductStatus = "archived"
)

// Product is the aggregate root of the Catalog context.
type Product struct {
	ID          string
	TenantID    string
	Name        string
	Slug        string
	Description string
	SKU         string
	Price       int64 // in smallest currency unit (e.g., cents)
	Currency    string
	Status      ProductStatus
	CategoryID  *string
	ImageURL    string
	Weight      int // in grams
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// Category groups related products.
type Category struct {
	ID          string
	TenantID    string
	Name        string
	Slug        string
	Description string
	ParentID    *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// ProductImage holds additional product images.
type ProductImage struct {
	ID        string
	ProductID string
	URL       string
	AltText   string
	Position  int
	CreatedAt time.Time
}

// ProductTag holds tags for products.
type ProductTag struct {
	ID        string
	ProductID string
	Tag       string
}
