package transport

import (
	"github.com/cloudcommerce/payment-service/internal/application"
	"github.com/cloudcommerce/payment-service/internal/domain"
	"github.com/cloudcommerce/shared-go/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// PaymentHandler handles HTTP requests for payments.
type PaymentHandler struct {
	service   *application.PaymentService
	validator *validator.Validate
}

// NewPaymentHandler creates a new payment handler.
func NewPaymentHandler(service *application.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service:   service,
		validator: validator.New(),
	}
}

// Create handles POST /api/v1/payments.
func (h *PaymentHandler) Create(c *gin.Context) {
	var req CreatePaymentRequest
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

	payment, err := h.service.CreatePayment(c.Request.Context(), application.CreatePaymentInput{
		TenantID:      tenantID,
		OrderID:       req.OrderID,
		UserID:        userID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
		Provider:      req.Provider,
	})
	if err != nil {
		response.InternalError(c, "Failed to create payment")
		return
	}

	response.Created(c, toPaymentResponse(payment))
}

// GetByID handles GET /api/v1/payments/:id.
func (h *PaymentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetString("tenant_id")

	payment, err := h.service.GetPaymentRepository().FindByID(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrPaymentNotFound {
			response.NotFound(c, "Payment not found")
			return
		}
		response.InternalError(c, "Failed to get payment")
		return
	}

	if payment.TenantID != tenantID {
		response.Forbidden(c, "Access denied")
		return
	}

	response.OK(c, toPaymentResponse(payment))
}

// GetByOrderID handles GET /api/v1/payments/order/:order_id.
func (h *PaymentHandler) GetByOrderID(c *gin.Context) {
	orderID := c.Param("order_id")
	tenantID := c.GetString("tenant_id")

	payment, err := h.service.GetPaymentRepository().FindByOrderID(c.Request.Context(), orderID)
	if err != nil {
		if err == domain.ErrPaymentNotFound {
			response.NotFound(c, "Payment not found for this order")
			return
		}
		response.InternalError(c, "Failed to get payment")
		return
	}

	if payment.TenantID != tenantID {
		response.Forbidden(c, "Access denied")
		return
	}

	response.OK(c, toPaymentResponse(payment))
}

// List handles GET /api/v1/payments.
func (h *PaymentHandler) List(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	offset := 0
	limit := 20

	payments, err := h.service.GetPaymentRepository().ListByTenant(c.Request.Context(), tenantID, offset, limit)
	if err != nil {
		response.InternalError(c, "Failed to list payments")
		return
	}

	items := make([]PaymentResponse, len(payments))
	for i, p := range payments {
		items[i] = toPaymentResponse(p)
	}

	response.OK(c, items)
}

// HandleWebhook handles POST /webhooks/payments/:provider.
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	_ = c.Param("provider") // TODO: route to provider-specific handler

	// Read raw body
	payload, err := c.GetRawData()
	if err != nil {
		response.InternalError(c, "Failed to read webhook payload")
		return
	}

	// Get signature from header
	signature := c.GetHeader("X-Signature")
	if signature == "" {
		signature = c.GetHeader("X-Midtrans-Signature")
	}

	// Process webhook
	if err := h.service.HandleWebhook(c.Request.Context(), payload, signature); err != nil {
		if err == domain.ErrWebhookAlreadyProcessed {
			// Idempotent - return success
			c.Status(200)
			return
		}
		response.InternalError(c, "Failed to process webhook")
		return
	}

	c.Status(200)
}

// toPaymentResponse converts domain entity to response DTO.
func toPaymentResponse(p *domain.Payment) PaymentResponse {
	return PaymentResponse{
		ID:                    p.ID,
		TenantID:              p.TenantID,
		OrderID:               p.OrderID,
		UserID:                p.UserID,
		Amount:                p.Amount,
		Currency:              p.Currency,
		Status:                string(p.Status),
		PaymentMethod:         p.PaymentMethod,
		PaymentProvider:       p.PaymentProvider,
		ProviderTransactionID: p.ProviderTransactionID,
		ProviderPaymentID:     p.ProviderPaymentID,
		PaymentURL:            p.PaymentURL,
		ExpiresAt:             p.ExpiresAt,
		PaidAt:                p.PaidAt,
		CreatedAt:             p.CreatedAt,
		UpdatedAt:             p.UpdatedAt,
	}
}
