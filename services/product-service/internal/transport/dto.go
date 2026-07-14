package transport

import (
	"strconv"
	"time"

	"github.com/cloudcommerce/product-service/internal/domain"
	"github.com/cloudcommerce/shared-go/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// === Product DTOs ===

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description" validate:"omitempty,max=2000"`
	SKU         string  `json:"sku" validate:"required,min=1,max=100"`
	Price       int64   `json:"price" validate:"required,min=1"`
	Currency    string  `json:"currency" validate:"omitempty,len=3"`
	CategoryID  *string `json:"categoryId,omitempty"`
	ImageURL    string  `json:"imageUrl" validate:"omitempty,url"`
	Weight      int     `json:"weight" validate:"omitempty,min=0"`
}

type UpdateProductRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=2000"`
	Price       *int64  `json:"price,omitempty" validate:"omitempty,min=1"`
	CategoryID  *string `json:"categoryId,omitempty"`
	ImageURL    *string `json:"imageUrl,omitempty" validate:"omitempty,url"`
	Status      *string `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
	Weight      *int    `json:"weight,omitempty" validate:"omitempty,min=0"`
}

type ProductResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	SKU         string    `json:"sku"`
	Price       int64     `json:"price"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	CategoryID  *string   `json:"categoryId,omitempty"`
	ImageURL    string    `json:"imageUrl,omitempty"`
	Weight      int       `json:"weight"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// === Category DTOs ===

type CreateCategoryRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description string  `json:"description" validate:"omitempty,max=500"`
	ParentID    *string `json:"parentId,omitempty"`
}

type CategoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	ParentID    *string   `json:"parentId,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

var validate = validator.New()

func validateStruct(req interface{}) error {
	return validate.Struct(req)
}

func toProductResponse(p *domain.Product) ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		SKU:         p.SKU,
		Price:       p.Price,
		Currency:    p.Currency,
		Status:      string(p.Status),
		CategoryID:  p.CategoryID,
		ImageURL:    p.ImageURL,
		Weight:      p.Weight,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func toCategoryResponse(c *domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		ParentID:    c.ParentID,
		CreatedAt:   c.CreatedAt,
	}
}

func getTenantID(c *gin.Context) string {
	return c.GetString("tenant_id")
}

func getPagination(c *gin.Context) (int, int) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

// Ensure response package is used
var _ = response.OK
