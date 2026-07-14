package transport

import "time"

// CreateOrderRequest is the request for creating an order.
type CreateOrderRequest struct {
	Items           []CreateOrderItemRequest `json:"items" validate:"required,min=1"`
	Currency        string                   `json:"currency" validate:"required,len=3"`
	Notes           string                   `json:"notes,omitempty"`
	TaxAmount       int64                    `json:"tax_amount"`
	ShippingAmount  int64                    `json:"shipping_amount"`
	DiscountAmount  int64                    `json:"discount_amount"`
	ShippingAddress *AddressRequest          `json:"shipping_address,omitempty"`
	BillingAddress  *AddressRequest          `json:"billing_address,omitempty"`
}

// CreateOrderItemRequest is a line item in the order request.
type CreateOrderItemRequest struct {
	ProductID      string `json:"product_id" validate:"required"`
	VariantID      string `json:"variant_id,omitempty"`
	ProductName    string `json:"product_name" validate:"required"`
	SKU            string `json:"sku,omitempty"`
	Quantity       int    `json:"quantity" validate:"required,min=1"`
	UnitPrice      int64  `json:"unit_price" validate:"required,min=0"`
	TotalPrice     int64  `json:"total_price" validate:"required,min=0"`
	TaxAmount      int64  `json:"tax_amount"`
	DiscountAmount int64  `json:"discount_amount"`
}

// AddressRequest represents an address.
type AddressRequest struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// OrderResponse is the response for an order.
type OrderResponse struct {
	ID              string              `json:"id"`
	TenantID        string              `json:"tenant_id"`
	UserID          string              `json:"user_id"`
	OrderNumber     string              `json:"order_number"`
	Status          string              `json:"status"`
	Subtotal        int64               `json:"subtotal"`
	TaxAmount       int64               `json:"tax_amount"`
	ShippingAmount  int64               `json:"shipping_amount"`
	DiscountAmount  int64               `json:"discount_amount"`
	TotalAmount     int64               `json:"total_amount"`
	Currency        string              `json:"currency"`
	Notes           string              `json:"notes,omitempty"`
	ShippingAddress interface{}         `json:"shipping_address,omitempty"`
	BillingAddress  interface{}         `json:"billing_address,omitempty"`
	Items           []OrderItemResponse `json:"items"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

// OrderItemResponse is the response for an order item.
type OrderItemResponse struct {
	ID             string `json:"id"`
	ProductID      string `json:"product_id"`
	VariantID      string `json:"variant_id,omitempty"`
	ProductName    string `json:"product_name"`
	SKU            string `json:"sku,omitempty"`
	Quantity       int    `json:"quantity"`
	UnitPrice      int64  `json:"unit_price"`
	TotalPrice     int64  `json:"total_price"`
	TaxAmount      int64  `json:"tax_amount"`
	DiscountAmount int64  `json:"discount_amount"`
}

// OrderListResponse is the response for listing orders.
type OrderListResponse struct {
	Data    []OrderResponse `json:"data"`
	Total   int64           `json:"total"`
	Page    int             `json:"page"`
	PerPage int             `json:"per_page"`
}
