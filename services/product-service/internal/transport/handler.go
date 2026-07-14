package transport

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cloudcommerce/product-service/internal/application"
	"github.com/cloudcommerce/product-service/internal/domain"
	"github.com/cloudcommerce/shared-go/response"
	"github.com/gin-gonic/gin"
)

// ProductHandler handles HTTP requests for product endpoints.
type ProductHandler struct {
	productSvc  *application.ProductService
	categorySvc *application.CategoryService
}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler(productSvc *application.ProductService, categorySvc *application.CategoryService) *ProductHandler {
	return &ProductHandler{productSvc: productSvc, categorySvc: categorySvc}
}

// Create handles POST /api/v1/products.
func (h *ProductHandler) Create(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == "" {
		response.Forbidden(c, "Tenant context is required")
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	if err := validateStruct(req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	product, err := h.productSvc.Create(c.Request.Context(), application.CreateProductInput{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Price:       req.Price,
		Currency:    req.Currency,
		CategoryID:  req.CategoryID,
		ImageURL:    req.ImageURL,
		Weight:      req.Weight,
	})
	if err != nil {
		if errors.Is(err, domain.ErrSKUExists) {
			response.Conflict(c, "SKU already exists")
			return
		}
		response.InternalError(c, "Failed to create product")
		return
	}

	response.Created(c, toProductResponse(product))
}

// List handles GET /api/v1/products.
func (h *ProductHandler) List(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == "" {
		response.Forbidden(c, "Tenant context is required")
		return
	}

	limit, offset := getPagination(c)
	status := c.Query("status")

	products, total, err := h.productSvc.List(c.Request.Context(), tenantID, limit, offset, status)
	if err != nil {
		response.InternalError(c, "Failed to list products")
		return
	}

	var items []ProductResponse
	for _, p := range products {
		items = append(items, toProductResponse(p))
	}

	meta := &response.Meta{
		Total:   total,
		PerPage: limit,
		Page:    offset/limit + 1,
	}

	response.OKWithMeta(c, items, meta)
}

// GetByID handles GET /api/v1/products/:id.
func (h *ProductHandler) GetByID(c *gin.Context) {
	tenantID := getTenantID(c)
	id := c.Param("id")

	product, err := h.productSvc.GetByID(c.Request.Context(), id, tenantID)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.NotFound(c, "Product not found")
			return
		}
		response.InternalError(c, "Failed to get product")
		return
	}

	response.OK(c, toProductResponse(product))
}

// Update handles PUT /api/v1/products/:id.
func (h *ProductHandler) Update(c *gin.Context) {
	tenantID := getTenantID(c)
	id := c.Param("id")

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	input := application.UpdateProductInput{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryID:  req.CategoryID,
		ImageURL:    req.ImageURL,
		Weight:      req.Weight,
	}

	if req.Status != nil {
		status := domain.ProductStatus(*req.Status)
		input.Status = &status
	}

	product, err := h.productSvc.Update(c.Request.Context(), id, tenantID, input)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.NotFound(c, "Product not found")
			return
		}
		response.InternalError(c, "Failed to update product")
		return
	}

	response.OK(c, toProductResponse(product))
}

// Delete handles DELETE /api/v1/products/:id.
func (h *ProductHandler) Delete(c *gin.Context) {
	tenantID := getTenantID(c)
	id := c.Param("id")

	if err := h.productSvc.Delete(c.Request.Context(), id, tenantID); err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.NotFound(c, "Product not found")
			return
		}
		response.InternalError(c, "Failed to delete product")
		return
	}

	c.Status(http.StatusNoContent)
}

// Publish handles POST /api/v1/products/:id/publish.
func (h *ProductHandler) Publish(c *gin.Context) {
	tenantID := getTenantID(c)
	id := c.Param("id")

	product, err := h.productSvc.Publish(c.Request.Context(), id, tenantID)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.NotFound(c, "Product not found")
			return
		}
		response.InternalError(c, "Failed to publish product")
		return
	}

	response.OK(c, toProductResponse(product))
}

// Archive handles POST /api/v1/products/:id/archive.
func (h *ProductHandler) Archive(c *gin.Context) {
	tenantID := getTenantID(c)
	id := c.Param("id")

	product, err := h.productSvc.Archive(c.Request.Context(), id, tenantID)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.NotFound(c, "Product not found")
			return
		}
		response.InternalError(c, "Failed to archive product")
		return
	}

	response.OK(c, toProductResponse(product))
}

// === Category Handlers ===

// CreateCategory handles POST /api/v1/categories.
func (h *ProductHandler) CreateCategory(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == "" {
		response.Forbidden(c, "Tenant context is required")
		return
	}

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	if err := validateStruct(req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	category, err := h.categorySvc.Create(c.Request.Context(), application.CreateCategoryInput{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
	})
	if err != nil {
		response.InternalError(c, "Failed to create category")
		return
	}

	response.Created(c, toCategoryResponse(category))
}

// ListCategories handles GET /api/v1/categories.
func (h *ProductHandler) ListCategories(c *gin.Context) {
	tenantID := getTenantID(c)
	if tenantID == "" {
		response.Forbidden(c, "Tenant context is required")
		return
	}

	categories, err := h.categorySvc.List(c.Request.Context(), tenantID)
	if err != nil {
		response.InternalError(c, "Failed to list categories")
		return
	}

	var items []CategoryResponse
	for _, cat := range categories {
		items = append(items, toCategoryResponse(cat))
	}

	response.OK(c, items)
}

// GetCategory handles GET /api/v1/categories/:id.
func (h *ProductHandler) GetCategory(c *gin.Context) {
	tenantID := getTenantID(c)
	id := c.Param("id")

	category, err := h.categorySvc.GetByID(c.Request.Context(), id, tenantID)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			response.NotFound(c, "Category not found")
			return
		}
		response.InternalError(c, "Failed to get category")
		return
	}

	response.OK(c, toCategoryResponse(category))
}

// DeleteCategory handles DELETE /api/v1/categories/:id.
func (h *ProductHandler) DeleteCategory(c *gin.Context) {
	tenantID := getTenantID(c)
	id := c.Param("id")

	if err := h.categorySvc.Delete(c.Request.Context(), id, tenantID); err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			response.NotFound(c, "Category not found")
			return
		}
		response.InternalError(c, "Failed to delete category")
		return
	}

	c.Status(http.StatusNoContent)
}

// Suppress unused import warning
var _ = strconv.Atoi
