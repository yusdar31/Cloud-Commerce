package transport

import (
	"encoding/json"
	"net/http"

	"github.com/cloudcommerce/order-service/internal/application"
	"github.com/cloudcommerce/order-service/internal/domain"
	"github.com/cloudcommerce/shared-go/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// OrderHandler handles HTTP requests for orders.
type OrderHandler struct {
	service   *application.OrderService
	validator *validator.Validate
}

// NewOrderHandler creates a new order handler.
func NewOrderHandler(service *application.OrderService) *OrderHandler {
	return &OrderHandler{
		service:   service,
		validator: validator.New(),
	}
}

// Create handles POST /api/v1/orders.
func (h *OrderHandler) Create(c *gin.Context) {
	var req CreateOrderRequest
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

	// Convert items
	items := make([]application.CreateOrderItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = application.CreateOrderItemInput{
			ProductID:      item.ProductID,
			VariantID:      item.VariantID,
			ProductName:    item.ProductName,
			SKU:            item.SKU,
			Quantity:       item.Quantity,
			UnitPrice:      item.UnitPrice,
			TotalPrice:     item.TotalPrice,
			TaxAmount:      item.TaxAmount,
			DiscountAmount: item.DiscountAmount,
		}
	}

	// Marshal addresses
	var shippingAddr, billingAddr []byte
	if req.ShippingAddress != nil {
		shippingAddr, _ = json.Marshal(req.ShippingAddress)
	}
	if req.BillingAddress != nil {
		billingAddr, _ = json.Marshal(req.BillingAddress)
	}

	input := application.CreateOrderInput{
		TenantID:        tenantID,
		UserID:          userID,
		Currency:        req.Currency,
		Notes:           req.Notes,
		TaxAmount:       req.TaxAmount,
		ShippingAmount:  req.ShippingAmount,
		DiscountAmount:  req.DiscountAmount,
		ShippingAddress: shippingAddr,
		BillingAddress:  billingAddr,
		Items:           items,
	}

	order, err := h.service.CreateOrder(c.Request.Context(), input)
	if err != nil {
		response.InternalError(c, "Failed to create order")
		return
	}

	response.Created(c, toOrderResponse(order))
}

// GetByID handles GET /api/v1/orders/:id.
func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetString("tenant_id")

	order, err := h.service.GetOrderRepository().FindByID(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrOrderNotFound {
			response.NotFound(c, "Order not found")
			return
		}
		response.InternalError(c, "Failed to get order")
		return
	}

	if order.TenantID != tenantID {
		response.Forbidden(c, "Access denied")
		return
	}

	// Load items
	items, _ := h.service.GetItemRepository().FindByOrderID(c.Request.Context(), order.ID)
	order.Items = items

	response.OK(c, toOrderResponse(order))
}

// List handles GET /api/v1/orders.
func (h *OrderHandler) List(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	offset := 0
	limit := 20

	orders, err := h.service.GetOrderRepository().ListByTenant(c.Request.Context(), tenantID, offset, limit)
	if err != nil {
		response.InternalError(c, "Failed to list orders")
		return
	}

	items := make([]OrderResponse, len(orders))
	for i, o := range orders {
		orderItems, _ := h.service.GetItemRepository().FindByOrderID(c.Request.Context(), o.ID)
		o.Items = orderItems
		items[i] = toOrderResponse(o)
	}

	response.OK(c, items)
}

// HandlePaymentCallback handles internal payment callback from payment-service.
func (h *OrderHandler) HandlePaymentCallback(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// TODO: Verify internal service auth
	_ = id
	_ = req

	c.JSON(http.StatusOK, gin.H{"message": "callback received"})
}

// toOrderResponse converts domain entity to response DTO.
func toOrderResponse(o *domain.Order) OrderResponse {
	itemResponses := make([]OrderItemResponse, len(o.Items))
	for i, item := range o.Items {
		itemResponses[i] = OrderItemResponse{
			ID:             item.ID,
			ProductID:      item.ProductID,
			VariantID:      item.VariantID,
			ProductName:    item.ProductName,
			SKU:            item.SKU,
			Quantity:       item.Quantity,
			UnitPrice:      item.UnitPrice,
			TotalPrice:     item.TotalPrice,
			TaxAmount:      item.TaxAmount,
			DiscountAmount: item.DiscountAmount,
		}
	}

	var shippingAddr, billingAddr interface{}
	if o.ShippingAddress != nil {
		json.Unmarshal(o.ShippingAddress, &shippingAddr)
	}
	if o.BillingAddress != nil {
		json.Unmarshal(o.BillingAddress, &billingAddr)
	}

	return OrderResponse{
		ID:              o.ID,
		TenantID:        o.TenantID,
		UserID:          o.UserID,
		OrderNumber:     o.OrderNumber,
		Status:          string(o.Status),
		Subtotal:        o.Subtotal,
		TaxAmount:       o.TaxAmount,
		ShippingAmount:  o.ShippingAmount,
		DiscountAmount:  o.DiscountAmount,
		TotalAmount:     o.TotalAmount,
		Currency:        o.Currency,
		Notes:           o.Notes,
		ShippingAddress: shippingAddr,
		BillingAddress:  billingAddr,
		Items:           itemResponses,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}
}
