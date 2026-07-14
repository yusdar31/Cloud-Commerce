package application

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cloudcommerce/inventory-service/internal/domain"
	"github.com/cloudcommerce/shared-go/events"
	"github.com/cloudcommerce/shared-go/nats"
	"github.com/google/uuid"
)

// InventoryService handles inventory business logic.
type InventoryService struct {
	inventoryRepo       domain.InventoryRepository
	reservationRepo     domain.ReservationRepository
	stockMovementRepo   domain.StockMovementRepository
	publisher           *nats.Publisher
	logger              *slog.Logger
}

// NewInventoryService creates a new inventory service.
func NewInventoryService(
	inventoryRepo domain.InventoryRepository,
	reservationRepo domain.ReservationRepository,
	stockMovementRepo domain.StockMovementRepository,
	publisher *nats.Publisher,
	logger *slog.Logger,
) *InventoryService {
	return &InventoryService{
		inventoryRepo:     inventoryRepo,
		reservationRepo:   reservationRepo,
		stockMovementRepo: stockMovementRepo,
		publisher:         publisher,
		logger:            logger,
	}
}

// CreateInventory creates initial inventory for a product.
func (s *InventoryService) CreateInventory(ctx context.Context, input CreateInventoryInput) (*domain.Inventory, error) {
	if input.Quantity < 0 {
		return nil, domain.ErrInvalidQuantity
	}

	inv := &domain.Inventory{
		ID:                uuid.NewString(),
		TenantID:          input.TenantID,
		ProductID:         input.ProductID,
		SKU:               input.SKU,
		Quantity:          input.Quantity,
		Reserved:          0,
		LowStockThreshold: input.LowStockThreshold,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		CreatedBy:         input.CreatedBy,
		Version:           1,
	}

	if err := s.inventoryRepo.Create(ctx, inv); err != nil {
		return nil, fmt.Errorf("failed to create inventory: %w", err)
	}

	// Record stock movement
	movement := &domain.StockMovement{
		ID:             uuid.NewString(),
		TenantID:       input.TenantID,
		ProductID:      input.ProductID,
		SKU:            input.SKU,
		MovementType:   domain.MovementInitial,
		QuantityBefore: 0,
		QuantityAfter:  input.Quantity,
		QuantityChange: input.Quantity,
		CreatedAt:      time.Now(),
		CreatedBy:      input.CreatedBy,
	}
	s.stockMovementRepo.Create(ctx, movement)

	s.logger.Info("inventory created",
		"inventory_id", inv.ID,
		"product_id", inv.ProductID,
		"quantity", inv.Quantity,
	)

	return inv, nil
}

// ReserveStock reserves stock for an order.
func (s *InventoryService) ReserveStock(ctx context.Context, input ReserveStockInput) (string, error) {
	// Get inventory
	inv, err := s.inventoryRepo.FindByProductID(ctx, input.TenantID, input.ProductID)
	if err != nil {
		return "", err
	}

	// Check availability
	if !inv.CanReserve(input.Quantity) {
		return "", domain.ErrInsufficientStock
	}

	// Create reservation
	reservation := &domain.Reservation{
		ID:        uuid.NewString(),
		TenantID:  input.TenantID,
		OrderID:   input.OrderID,
		ProductID: input.ProductID,
		SKU:       inv.SKU,
		Quantity:  input.Quantity,
		Status:    domain.ReservationActive,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.reservationRepo.Create(ctx, reservation); err != nil {
		return "", fmt.Errorf("failed to create reservation: %w", err)
	}

	// Update reserved quantity
	if err := s.inventoryRepo.Reserve(ctx, input.TenantID, input.ProductID, input.Quantity); err != nil {
		return "", fmt.Errorf("failed to reserve stock: %w", err)
	}

	// Publish event
	event := events.InventoryReservedEvent{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.NewString(),
			Timestamp: time.Now(),
			TenantID:  input.TenantID,
		},
		ReservationID: reservation.ID,
		OrderID:       input.OrderID,
		Items: []events.InventoryItemEvent{
			{
				ProductID:      input.ProductID,
				SKU:            inv.SKU,
				QuantityBefore: inv.Quantity,
				QuantityAfter:  inv.Quantity,
				Reserved:       input.Quantity,
			},
		},
		ExpiresAt: reservation.ExpiresAt,
	}
	s.publisher.Publish(ctx, events.SubjectInventoryReserved, event)

	s.logger.Info("stock reserved",
		"reservation_id", reservation.ID,
		"order_id", input.OrderID,
		"product_id", input.ProductID,
		"quantity", input.Quantity,
	)

	return reservation.ID, nil
}

// ReleaseStock releases a reservation.
func (s *InventoryService) ReleaseStock(ctx context.Context, tenantID, orderID, reason string) error {
	reservations, err := s.reservationRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	for _, res := range reservations {
		if res.Status != domain.ReservationActive {
			continue
		}

		// Release stock
		if err := s.inventoryRepo.Release(ctx, res.TenantID, res.ProductID, res.Quantity); err != nil {
			s.logger.Error("failed to release stock", "error", err, "reservation_id", res.ID)
			continue
		}

		// Update reservation status
		if err := s.reservationRepo.UpdateStatus(ctx, res.ID, domain.ReservationReleased); err != nil {
			s.logger.Error("failed to update reservation", "error", err, "reservation_id", res.ID)
		}

		s.logger.Info("stock released",
			"reservation_id", res.ID,
			"order_id", orderID,
			"product_id", res.ProductID,
			"quantity", res.Quantity,
		)
	}

	// Publish event
	event := events.InventoryReleasedEvent{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.NewString(),
			Timestamp: time.Now(),
			TenantID:  tenantID,
		},
		OrderID: orderID,
		Reason:  reason,
	}
	s.publisher.Publish(ctx, events.SubjectInventoryReleased, event)

	return nil
}

// CreateInventoryInput contains data for creating inventory.
type CreateInventoryInput struct {
	TenantID          string
	ProductID         string
	SKU               string
	Quantity          int
	LowStockThreshold int
	CreatedBy         string
}

// ReserveStockInput contains data for reserving stock.
type ReserveStockInput struct {
	TenantID  string
	OrderID   string
	ProductID string
	Quantity  int
}

// GetInventoryRepository returns the inventory repository for direct queries.
func (s *InventoryService) GetInventoryRepository() domain.InventoryRepository {
	return s.inventoryRepo
}
