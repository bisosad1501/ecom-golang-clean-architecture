package services

import (
	"context"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
)

// StockReservationService handles stock reservation business logic
type StockReservationService interface {
	// Reserve stock for an order
	ReserveStockForOrder(ctx context.Context, orderID, userID uuid.UUID, items []entities.CartItem) error

	// Confirm reservations (convert to actual stock reduction)
	ConfirmReservations(ctx context.Context, orderID uuid.UUID) error

	// Release reservations
	ReleaseReservations(ctx context.Context, orderID uuid.UUID) error

	// Check if stock can be reserved
	CanReserveStock(ctx context.Context, productID uuid.UUID, quantity int) (bool, error)

	// Get available stock (actual stock - reserved stock)
	GetAvailableStock(ctx context.Context, productID uuid.UUID) (int, error)

	// Cleanup expired reservations
	CleanupExpiredReservations(ctx context.Context) error

	// Extend reservation timeout
	ExtendReservation(ctx context.Context, orderID uuid.UUID, minutes int) error
}

type stockReservationService struct {
	reservationRepo repositories.StockReservationRepository
	productRepo     repositories.ProductRepository
	inventoryRepo   repositories.InventoryRepository
}

// NewStockReservationService creates a new stock reservation service
func NewStockReservationService(
	reservationRepo repositories.StockReservationRepository,
	productRepo repositories.ProductRepository,
	inventoryRepo repositories.InventoryRepository,
) StockReservationService {
	return &stockReservationService{
		reservationRepo: reservationRepo,
		productRepo:     productRepo,
		inventoryRepo:   inventoryRepo,
	}
}

// ReserveStockForOrder reserves stock for an order
func (s *stockReservationService) ReserveStockForOrder(ctx context.Context, orderID, userID uuid.UUID, items []entities.CartItem) error {
	var reservations []*entities.StockReservation

	// Check availability and create reservations
	for _, item := range items {
		// Check if we can reserve the stock
		canReserve, err := s.CanReserveStock(ctx, item.ProductID, item.Quantity)
		if err != nil {
			return fmt.Errorf("failed to check stock availability for product %s: %w", item.ProductID, err)
		}

		if !canReserve {
			return fmt.Errorf("insufficient stock for product %s: requested %d", item.ProductID, item.Quantity)
		}

		// Create reservation
		reservation := &entities.StockReservation{
			ID:        uuid.New(),
			ProductID: item.ProductID,
			OrderID:   &orderID,
			UserID:    userID,
			Quantity:  item.Quantity,
			Type:      entities.ReservationTypeOrder,
			Status:    entities.ReservationStatusActive,
			Notes:     fmt.Sprintf("Reserved for order %s", orderID.String()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Set expiration (30 minutes default)
		reservation.SetExpiration(30)

		reservations = append(reservations, reservation)
	}

	// Create all reservations in batch
	if err := s.reservationRepo.CreateBatch(ctx, reservations); err != nil {
		return fmt.Errorf("failed to create stock reservations: %w", err)
	}

	return nil
}

// ConfirmReservations confirms reservations and reduces actual stock
func (s *stockReservationService) ConfirmReservations(ctx context.Context, orderID uuid.UUID) error {
	// Get all reservations for the order
	reservations, err := s.reservationRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get reservations for order %s: %w", orderID, err)
	}

	// Confirm each reservation and reduce actual stock
	for _, reservation := range reservations {
		if !reservation.CanBeConfirmed() {
			continue // Skip expired or already confirmed reservations
		}

		// Get product and reduce stock
		product, err := s.productRepo.GetByID(ctx, reservation.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product %s: %w", reservation.ProductID, err)
		}

		// Reduce actual stock
		if err := product.ReduceStock(reservation.Quantity); err != nil {
			return fmt.Errorf("failed to reduce stock for product %s: %w", product.Name, err)
		}

		// Update product stock in database
		if err := s.productRepo.UpdateStock(ctx, reservation.ProductID, product.Stock); err != nil {
			return fmt.Errorf("failed to update stock for product %s: %w", product.Name, err)
		}

		// Record inventory movement for stock reduction
		if s.inventoryRepo != nil {
			// Try to get inventory record for the product
			inventory, err := s.inventoryRepo.GetByProductID(ctx, reservation.ProductID)
			if err == nil {
				// Update inventory stock levels
				if err := s.inventoryRepo.UpdateStock(ctx, inventory.ID, -reservation.Quantity, "order_confirmed"); err != nil {
					fmt.Printf("Warning: Failed to update inventory for product %s: %v\n", reservation.ProductID, err)
				}
			}
		}

		// Confirm the reservation
		reservation.Confirm()
		if err := s.reservationRepo.Update(ctx, reservation); err != nil {
			return fmt.Errorf("failed to confirm reservation %s: %w", reservation.ID, err)
		}
	}

	return nil
}

// ReleaseReservations releases all reservations for an order
func (s *stockReservationService) ReleaseReservations(ctx context.Context, orderID uuid.UUID) error {
	return s.reservationRepo.ReleaseByOrderID(ctx, orderID)
}

// CanReserveStock checks if stock can be reserved
func (s *stockReservationService) CanReserveStock(ctx context.Context, productID uuid.UUID, quantity int) (bool, error) {
	// Get product
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return false, err
	}

	// Check if product is available
	if !product.IsAvailable() {
		return false, nil
	}

	// Get available stock
	availableStock, err := s.GetAvailableStock(ctx, productID)
	if err != nil {
		return false, err
	}

	return availableStock >= quantity, nil
}

// GetAvailableStock gets available stock (actual stock - reserved stock)
func (s *stockReservationService) GetAvailableStock(ctx context.Context, productID uuid.UUID) (int, error) {
	// Get product stock
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return 0, err
	}

	// Get total reserved stock
	reservedStock, err := s.reservationRepo.GetTotalReservedStock(ctx, productID)
	if err != nil {
		return 0, err
	}

	// Use max to ensure available stock is not negative
	availableStock := max(0, product.Stock-reservedStock)

	return availableStock, nil
}

// CleanupExpiredReservations cleans up expired reservations
func (s *stockReservationService) CleanupExpiredReservations(ctx context.Context) error {
	return s.reservationRepo.ReleaseExpiredReservations(ctx)
}

// ExtendReservation extends reservation timeout for an order
func (s *stockReservationService) ExtendReservation(ctx context.Context, orderID uuid.UUID, minutes int) error {
	reservations, err := s.reservationRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get reservations for order %s: %w", orderID, err)
	}

	for _, reservation := range reservations {
		if reservation.IsActive() {
			reservation.ExtendExpiration(minutes)
			if err := s.reservationRepo.Update(ctx, reservation); err != nil {
				return fmt.Errorf("failed to extend reservation %s: %w", reservation.ID, err)
			}
		}
	}

	return nil
}
