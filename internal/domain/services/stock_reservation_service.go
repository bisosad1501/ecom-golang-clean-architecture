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

	// Reserve stock for a cart
	ReserveStockForCart(ctx context.Context, reservation *entities.StockReservation) error

	// Confirm reservations (convert to actual stock reduction)
	ConfirmReservations(ctx context.Context, orderID uuid.UUID) error

	// Release reservations
	ReleaseReservations(ctx context.Context, orderID uuid.UUID) error

	// Transfer cart reservations to order atomically
	TransferCartReservationsToOrder(ctx context.Context, userID, orderID uuid.UUID, items []entities.CartItem) error

	// Check if stock can be reserved
	CanReserveStock(ctx context.Context, productID uuid.UUID, quantity int) (bool, error)

	// Atomically check and reserve stock to prevent race conditions
	AtomicReserveStock(ctx context.Context, reservation *entities.StockReservation) error

	// Get available stock (actual stock - reserved stock)
	GetAvailableStock(ctx context.Context, productID uuid.UUID) (int, error)

	// Cleanup expired reservations
	CleanupExpiredReservations(ctx context.Context) error

	// Extend reservation timeout
	ExtendReservation(ctx context.Context, orderID uuid.UUID, minutes int) error

	// Get reservation statistics
	GetReservationStats(ctx context.Context, productID uuid.UUID) (*ReservationStats, error)
}

// ReservationStats represents stock reservation statistics
type ReservationStats struct {
	ProductID           uuid.UUID `json:"product_id"`
	TotalStock          int       `json:"total_stock"`
	ReservedStock       int       `json:"reserved_stock"`
	AvailableStock      int       `json:"available_stock"`
	ActiveReservations  int       `json:"active_reservations"`
	ExpiredReservations int       `json:"expired_reservations"`
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
			UserID:    &userID, // Convert to pointer
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

		// Reduce actual stock in Product entity
		if err := product.ReduceStock(reservation.Quantity); err != nil {
			return fmt.Errorf("failed to reduce stock for product %s: %w", product.Name, err)
		}

		// Update product stock in database
		if err := s.productRepo.UpdateStock(ctx, reservation.ProductID, product.Stock); err != nil {
			return fmt.Errorf("failed to update stock for product %s: %w", product.Name, err)
		}

		// Record inventory movement for tracking purposes only (don't double reduce)
		if s.inventoryRepo != nil {
			// Try to get inventory record for the product
			inventory, err := s.inventoryRepo.GetByProductID(ctx, reservation.ProductID)
			if err == nil {
				// Sync inventory quantity with product stock (don't subtract again)
				// This ensures inventory.quantity_on_hand matches product.stock
				if err := s.inventoryRepo.SyncWithProductStock(ctx, inventory.ID, product.Stock, "stock_sync_after_order_confirmation"); err != nil {
					fmt.Printf("Warning: Failed to sync inventory with product stock for product %s: %v\n", reservation.ProductID, err)
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

// AtomicReserveStock atomically checks and reserves stock to prevent race conditions
func (s *stockReservationService) AtomicReserveStock(ctx context.Context, reservation *entities.StockReservation) error {
	// Validate reservation
	if reservation.Quantity <= 0 {
		return fmt.Errorf("reservation quantity must be positive")
	}

	// Use database transaction with row-level locking to prevent race conditions
	return s.reservationRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		// Get product with row-level lock (SELECT FOR UPDATE)
		product, err := s.productRepo.GetByIDForUpdate(txCtx, reservation.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product for stock check: %w", err)
		}

		// Calculate available stock (current stock - active reservations)
		activeReservations, err := s.reservationRepo.GetActiveReservationsByProduct(txCtx, reservation.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get active reservations: %w", err)
		}

		totalReserved := 0
		for _, res := range activeReservations {
			totalReserved += res.Quantity
		}

		availableStock := product.Stock - totalReserved
		if availableStock < reservation.Quantity {
			return fmt.Errorf("insufficient stock for product %s: available %d, requested %d",
				reservation.ProductID, availableStock, reservation.Quantity)
		}

		// Create reservation atomically within the same transaction
		if err := s.reservationRepo.Create(txCtx, reservation); err != nil {
			return fmt.Errorf("failed to create stock reservation: %w", err)
		}

		return nil
	})
}

// TransferCartReservationsToOrder atomically transfers cart reservations to order reservations
func (s *stockReservationService) TransferCartReservationsToOrder(ctx context.Context, userID, orderID uuid.UUID, items []entities.CartItem) error {
	return s.reservationRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		// Get existing cart reservations for this user
		allReservations, err := s.reservationRepo.GetByUserID(txCtx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user reservations: %w", err)
		}

		// Filter for cart reservations
		var cartReservations []*entities.StockReservation
		for _, reservation := range allReservations {
			if reservation.Type == entities.ReservationTypeCart {
				cartReservations = append(cartReservations, reservation)
			}
		}

		// Create map of cart reservations by product ID for quick lookup
		cartReservationMap := make(map[uuid.UUID]*entities.StockReservation)
		for _, reservation := range cartReservations {
			if reservation.Status == entities.ReservationStatusActive {
				cartReservationMap[reservation.ProductID] = reservation
			}
		}

		// Create new order reservations
		var orderReservations []*entities.StockReservation
		for _, item := range items {
			// Check if we have a cart reservation for this product
			cartReservation, hasCartReservation := cartReservationMap[item.ProductID]

			// Verify stock availability for the full quantity needed
			canReserve, err := s.CanReserveStock(txCtx, item.ProductID, item.Quantity)
			if err != nil {
				return fmt.Errorf("failed to check stock availability for product %s: %w", item.ProductID, err)
			}

			if !canReserve {
				return fmt.Errorf("insufficient stock for product %s: requested %d", item.ProductID, item.Quantity)
			}

			// Create order reservation
			orderReservation := &entities.StockReservation{
				ID:        uuid.New(),
				ProductID: item.ProductID,
				OrderID:   &orderID,
				UserID:    &userID,
				Quantity:  item.Quantity,
				Type:      entities.ReservationTypeOrder,
				Status:    entities.ReservationStatusActive,
				Notes:     fmt.Sprintf("Transferred from cart to order %s", orderID.String()),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			orderReservation.SetExpiration(30) // 30 minutes for order reservation

			orderReservations = append(orderReservations, orderReservation)

			// Mark cart reservation for release if it exists
			if hasCartReservation {
				cartReservation.Release()
			}
		}

		// Create order reservations in batch
		if err := s.reservationRepo.CreateBatch(txCtx, orderReservations); err != nil {
			return fmt.Errorf("failed to create order reservations: %w", err)
		}

		// Release cart reservations
		for _, cartReservation := range cartReservationMap {
			if err := s.reservationRepo.Update(txCtx, cartReservation); err != nil {
				// Log warning but don't fail transaction
				continue
			}
		}

		return nil
	})
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

// GetReservationStats gets reservation statistics for a product
func (s *stockReservationService) GetReservationStats(ctx context.Context, productID uuid.UUID) (*ReservationStats, error) {
	// Get product
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Get all reservations for this product
	reservations, err := s.reservationRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservations: %w", err)
	}

	stats := &ReservationStats{
		ProductID:           productID,
		TotalStock:          product.Stock,
		ReservedStock:       0,
		ActiveReservations:  0,
		ExpiredReservations: 0,
	}

	// Calculate statistics
	for _, reservation := range reservations {
		if reservation.IsActive() {
			stats.ReservedStock += reservation.Quantity
			stats.ActiveReservations++
		} else if reservation.IsExpired() {
			stats.ExpiredReservations++
		}
	}

	stats.AvailableStock = stats.TotalStock - stats.ReservedStock
	if stats.AvailableStock < 0 {
		stats.AvailableStock = 0
	}

	return stats, nil
}

// ReserveStockForCart reserves stock for a cart
func (s *stockReservationService) ReserveStockForCart(ctx context.Context, reservation *entities.StockReservation) error {
	// Check if stock can be reserved
	canReserve, err := s.CanReserveStock(ctx, reservation.ProductID, reservation.Quantity)
	if err != nil {
		return fmt.Errorf("failed to check stock availability for product %s: %w", reservation.ProductID, err)
	}

	if !canReserve {
		return fmt.Errorf("insufficient stock for product %s: requested %d", reservation.ProductID, reservation.Quantity)
	}

	// Create reservation
	if err := s.reservationRepo.Create(ctx, reservation); err != nil {
		return fmt.Errorf("failed to create stock reservation for cart: %w", err)
	}

	return nil
}
