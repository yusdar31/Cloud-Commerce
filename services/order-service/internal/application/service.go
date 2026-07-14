package application

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cloudcommerce/order-service/internal/domain"
	"github.com/cloudcommerce/shared-go/events"
	"github.com/cloudcommerce/shared-go/nats"
	"github.com/google/uuid"
)

// OrderService handles order business logic.
type OrderService struct {
	orderRepo         domain.OrderRepository
	itemRepo          domain.OrderItemRepository
	statusHistoryRepo domain.OrderStatusHistoryRepository
	sagaEventRepo     domain.SagaEventRepository
	publisher         *nats.Publisher
	logger            *slog.Logger
}

// NewOrderService creates a new order service.
func NewOrderService(
	orderRepo domain.OrderRepository,
	itemRepo domain.OrderItemRepository,
	statusHistoryRepo domain.OrderStatusHistoryRepository,
	sagaEventRepo domain.SagaEventRepository,
	publisher *nats.Publisher,
	logger *slog.Logger,
) *OrderService {
	return &OrderService{
		orderRepo:         orderRepo,
		itemRepo:          itemRepo,
		statusHistoryRepo: statusHistoryRepo,
		sagaEventRepo:     sagaEventRepo,
		publisher:         publisher,
		logger:            logger,
	}
}

// GetOrderRepository returns the order repository for direct access.
func (s *OrderService) GetOrderRepository() domain.OrderRepository {
	return s.orderRepo
}

// GetItemRepository returns the item repository for direct access.
func (s *OrderService) GetItemRepository() domain.OrderItemRepository {
	return s.itemRepo
}

// CreateOrder creates a new order with items.
func (s *OrderService) CreateOrder(ctx context.Context, input CreateOrderInput) (*domain.Order, error) {
	// Generate order number
	orderNumber := fmt.Sprintf("ORD-%s-%d",
		time.Now().Format("20060102"),
		time.Now().UnixNano()%1000000,
	)

	// Calculate totals
	var subtotal int64
	for _, item := range input.Items {
		subtotal += item.TotalPrice
	}

	// Create order
	order := &domain.Order{
		ID:              uuid.NewString(),
		TenantID:        input.TenantID,
		UserID:          input.UserID,
		OrderNumber:     orderNumber,
		Status:          domain.OrderPending,
		Subtotal:        subtotal,
		TaxAmount:       input.TaxAmount,
		ShippingAmount:  input.ShippingAmount,
		DiscountAmount:  input.DiscountAmount,
		TotalAmount:     subtotal + input.TaxAmount + input.ShippingAmount - input.DiscountAmount,
		Currency:        input.Currency,
		Notes:           input.Notes,
		ShippingAddress: input.ShippingAddress,
		BillingAddress:  input.BillingAddress,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		CreatedBy:       input.UserID,
		Version:         1,
	}

	// Save order
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items
	for _, itemInput := range input.Items {
		item := &domain.OrderItem{
			ID:            uuid.NewString(),
			OrderID:       order.ID,
			TenantID:      input.TenantID,
			ProductID:     itemInput.ProductID,
			VariantID:     itemInput.VariantID,
			ProductName:   itemInput.ProductName,
			SKU:           itemInput.SKU,
			Quantity:      itemInput.Quantity,
			UnitPrice:     itemInput.UnitPrice,
			TotalPrice:    itemInput.TotalPrice,
			TaxAmount:     itemInput.TaxAmount,
			DiscountAmount: itemInput.DiscountAmount,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := s.itemRepo.Create(ctx, item); err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	// Add status history
	s.addStatusHistory(ctx, order.ID, domain.OrderPending, "Order created", input.UserID)

	// Publish OrderCreated event
	event := events.OrderCreatedEvent{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.NewString(),
			Timestamp: time.Now(),
			TenantID:  input.TenantID,
			UserID:    input.UserID,
		},
		OrderID:      order.ID,
		OrderNumber:  orderNumber,
		TotalAmount:  order.TotalAmount,
		Currency:     order.Currency,
		Items:        toEventItems(input.Items),
	}
	s.publisher.Publish(ctx, events.SubjectOrderCreated, event)

	s.logger.Info("order created",
		"order_id", order.ID,
		"order_number", orderNumber,
		"total", order.TotalAmount,
	)

	return order, nil
}

// HandlePaymentSuccess processes a successful payment event.
func (s *OrderService) HandlePaymentSuccess(ctx context.Context, event events.PaymentSuccessEvent) error {
	order, err := s.orderRepo.FindByID(ctx, event.OrderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	if order.Status != domain.OrderConfirmed {
		s.logger.Warn("order not in confirmed status",
			"order_id", order.ID,
			"status", order.Status,
		)
		return nil
	}

	// Update order status to paid
	if err := s.orderRepo.UpdateStatus(ctx, order.ID, domain.OrderPaid); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	s.addStatusHistory(ctx, order.ID, domain.OrderPaid, "Payment received", "system")

	// Publish OrderPaidEvent for inventory reservation
	paidEvent := events.OrderPaidEvent{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.NewString(),
			Timestamp: time.Now(),
			TenantID:  order.TenantID,
			UserID:    order.UserID,
		},
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		TotalAmount: order.TotalAmount,
	}
	s.publisher.Publish(ctx, events.SubjectOrderPaid, paidEvent)

	s.logger.Info("order paid",
		"order_id", order.ID,
		"payment_id", event.PaymentID,
	)

	return nil
}

// HandlePaymentFailed processes a failed payment event.
func (s *OrderService) HandlePaymentFailed(ctx context.Context, event events.PaymentFailedEvent) error {
	order, err := s.orderRepo.FindByID(ctx, event.OrderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Update order status to failed
	if err := s.orderRepo.UpdateStatus(ctx, order.ID, domain.OrderFailed); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	s.addStatusHistory(ctx, order.ID, domain.OrderFailed, "Payment failed", "system")

	s.logger.Info("order payment failed",
		"order_id", order.ID,
	)

	return nil
}

// HandleInventoryReserved processes an inventory reserved event.
func (s *OrderService) HandleInventoryReserved(ctx context.Context, event events.InventoryReservedEvent) error {
	order, err := s.orderRepo.FindByID(ctx, event.OrderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	if order.Status != domain.OrderPaid {
		s.logger.Warn("order not in paid status for inventory reservation",
			"order_id", order.ID,
			"status", order.Status,
		)
		return nil
	}

	// Update order status to confirmed/processing
	if err := s.orderRepo.UpdateStatus(ctx, order.ID, domain.OrderProcessing); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	s.addStatusHistory(ctx, order.ID, domain.OrderProcessing, "Inventory reserved", "system")

	s.logger.Info("order inventory reserved",
		"order_id", order.ID,
	)

	return nil
}

// HandleInventoryReservationFailed processes a failed inventory reservation.
func (s *OrderService) HandleInventoryReservationFailed(ctx context.Context, event events.InventoryReservationFailedEvent) error {
	order, err := s.orderRepo.FindByID(ctx, event.OrderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Update order status to failed
	if err := s.orderRepo.UpdateStatus(ctx, order.ID, domain.OrderFailed); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	s.addStatusHistory(ctx, order.ID, domain.OrderFailed, "Inventory reservation failed: "+event.Reason, "system")

	// TODO: Trigger payment refund

	s.logger.Info("order inventory reservation failed",
		"order_id", order.ID,
		"reason", event.Reason,
	)

	return nil
}

// addStatusHistory adds a status history record.
func (s *OrderService) addStatusHistory(ctx context.Context, orderID string, status domain.OrderStatus, notes, changedBy string) {
	history := &domain.OrderStatusHistory{
		ID:        uuid.NewString(),
		OrderID:   orderID,
		Status:    status,
		Notes:     notes,
		ChangedBy: changedBy,
		CreatedAt: time.Now(),
	}
	s.statusHistoryRepo.Create(ctx, history)
}

// toEventItems converts input items to event items.
func toEventItems(items []CreateOrderItemInput) []events.OrderItemEvent {
	eventItems := make([]events.OrderItemEvent, len(items))
	for i, item := range items {
		eventItems[i] = events.OrderItemEvent{
			ProductID: item.ProductID,
			SKU:       item.SKU,
			Quantity:  item.Quantity,
			Price:     item.UnitPrice,
		}
	}
	return eventItems
}

// CreateOrderInput contains data for creating an order.
type CreateOrderInput struct {
	TenantID        string
	UserID          string
	Currency        string
	Notes           string
	TaxAmount       int64
	ShippingAmount  int64
	DiscountAmount  int64
	ShippingAddress []byte
	BillingAddress  []byte
	Items           []CreateOrderItemInput
}

// CreateOrderItemInput contains data for creating an order item.
type CreateOrderItemInput struct {
	ProductID      string
	VariantID      string
	ProductName    string
	SKU            string
	Quantity       int
	UnitPrice      int64
	TotalPrice     int64
	TaxAmount      int64
	DiscountAmount int64
}
