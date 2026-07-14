package transport

import (
	"github.com/cloudcommerce/inventory-service/internal/application"
	"github.com/cloudcommerce/inventory-service/internal/domain"
	"github.com/cloudcommerce/shared-go/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// InventoryHandler handles HTTP requests for inventory.
type InventoryHandler struct {
	service   *application.InventoryService
	validator *validator.Validate
}

// NewInventoryHandler creates a new inventory handler.
func NewInventoryHandler(service *application.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		service:   service,
		validator: validator.New(),
	}
}

// Create handles POST /api/v1/inventory.
func (h *InventoryHandler) Create(c *gin.Context) {
	var req CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	inv, err := h.service.CreateInventory(c.Request.Context(), application.CreateInventoryInput{
		TenantID:          tenantID,
		ProductID:         req.ProductID,
		SKU:               req.SKU,
		Quantity:          req.Quantity,
		LowStockThreshold: req.LowStockThreshold,
		CreatedBy:         userID,
	})
	if err != nil {
		if err == domain.ErrProductAlreadyExists {
			response.Conflict(c, "Inventory already exists for this product")
			return
		}
		response.InternalError(c, "Failed to create inventory")
		return
	}

	response.Created(c, toInventoryResponse(inv))
}

// GetByID handles GET /api/v1/inventory/:id.
func (h *InventoryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetString("tenant_id")

	inv, err := h.service.GetInventoryRepository().FindByID(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrInventoryNotFound {
			response.NotFound(c, "Inventory not found")
			return
		}
		response.InternalError(c, "Failed to get inventory")
		return
	}

	if inv.TenantID != tenantID {
		response.Forbidden(c, "Access denied")
		return
	}

	response.OK(c, toInventoryResponse(inv))
}

// List handles GET /api/v1/inventory.
func (h *InventoryHandler) List(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	offset := 0
	limit := 20

	inventories, err := h.service.GetInventoryRepository().List(c.Request.Context(), tenantID, offset, limit)
	if err != nil {
		response.InternalError(c, "Failed to list inventory")
		return
	}

	items := make([]InventoryResponse, len(inventories))
	for i, inv := range inventories {
		items[i] = toInventoryResponse(inv)
	}

	response.OK(c, items)
}

// Reserve handles POST /api/v1/inventory/reserve.
func (h *InventoryHandler) Reserve(c *gin.Context) {
	var req ReserveStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	tenantID := c.GetString("tenant_id")

	reservationID, err := h.service.ReserveStock(c.Request.Context(), application.ReserveStockInput{
		TenantID:  tenantID,
		OrderID:   req.OrderID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	})
	if err != nil {
		if err == domain.ErrInsufficientStock {
			response.Conflict(c, "Insufficient stock available")
			return
		}
		if err == domain.ErrInventoryNotFound {
			response.NotFound(c, "Product not found in inventory")
			return
		}
		response.InternalError(c, "Failed to reserve stock")
		return
	}

	response.OK(c, ReserveStockResponse{
		ReservationID: reservationID,
		Message:       "Stock reserved successfully",
	})
}

// Release handles POST /api/v1/inventory/release.
func (h *InventoryHandler) Release(c *gin.Context) {
	var req ReleaseStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	tenantID := c.GetString("tenant_id")

	if err := h.service.ReleaseStock(c.Request.Context(), tenantID, req.OrderID, req.Reason); err != nil {
		response.InternalError(c, "Failed to release stock")
		return
	}

	response.OK(c, gin.H{
		"message": "Stock released successfully",
	})
}

// GetLowStock handles GET /api/v1/inventory/low-stock.
func (h *InventoryHandler) GetLowStock(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	inventories, err := h.service.GetInventoryRepository().FindLowStock(c.Request.Context(), tenantID)
	if err != nil {
		response.InternalError(c, "Failed to get low stock items")
		return
	}

	items := make([]InventoryResponse, len(inventories))
	for i, inv := range inventories {
		items[i] = toInventoryResponse(inv)
	}

	response.OK(c, items)
}

// toInventoryResponse converts domain entity to response DTO.
func toInventoryResponse(inv *domain.Inventory) InventoryResponse {
	return InventoryResponse{
		ID:                inv.ID,
		TenantID:          inv.TenantID,
		ProductID:         inv.ProductID,
		SKU:               inv.SKU,
		Quantity:          inv.Quantity,
		Reserved:          inv.Reserved,
		Available:         inv.Available,
		LowStockThreshold: inv.LowStockThreshold,
		IsLowStock:        inv.IsLowStock(),
		CreatedAt:         inv.CreatedAt,
		UpdatedAt:         inv.UpdatedAt,
	}
}
