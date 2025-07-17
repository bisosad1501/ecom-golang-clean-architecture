package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
)

// SimpleStockService handles simplified stock management
// No complex reservations - just direct stock reduction on payment success
type SimpleStockService interface {
	// Check if stock is available for cart items
	CheckStockAvailability(ctx context.Context, items []entities.CartItem) error

	// Reduce stock when payment is successful
	ReduceStock(ctx context.Context, items []entities.CartItem) error

	// Reduce stock for order items when payment is confirmed
	ReduceStockForOrder(ctx context.Context, items []entities.OrderItem) error

	// Restore stock when order is cancelled/refunded
	RestoreStock(ctx context.Context, items []entities.OrderItem) error

	// Get available stock for a product
	GetAvailableStock(ctx context.Context, productID uuid.UUID) (int, error)
}

type simpleStockService struct {
	productRepo   repositories.ProductRepository
	inventoryRepo repositories.InventoryRepository
}

// NewSimpleStockService creates a new simple stock service
func NewSimpleStockService(
	productRepo repositories.ProductRepository,
	inventoryRepo repositories.InventoryRepository,
) SimpleStockService {
	return &simpleStockService{
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
	}
}

// CheckStockAvailability checks if stock is available for all cart items
func (s *simpleStockService) CheckStockAvailability(ctx context.Context, items []entities.CartItem) error {
	for _, item := range items {
		// Get current product stock
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product %s: %w", item.ProductID, err)
		}

		// Check if product is available
		if !product.IsAvailable() {
			return fmt.Errorf("product %s is not available", product.Name)
		}

		// Check stock availability
		if product.Stock < item.Quantity {
			return fmt.Errorf("insufficient stock for product %s: available=%d, requested=%d",
				product.Name, product.Stock, item.Quantity)
		}
	}

	return nil
}

// ReduceStock reduces stock for cart items when payment is successful
func (s *simpleStockService) ReduceStock(ctx context.Context, items []entities.CartItem) error {
	for _, item := range items {
		// Get current product
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product %s: %w", item.ProductID, err)
		}

		// Check stock availability one more time (race condition protection)
		if product.Stock < item.Quantity {
			return fmt.Errorf("insufficient stock for product %s: available=%d, requested=%d",
				product.Name, product.Stock, item.Quantity)
		}

		// Reduce product stock
		newStock := product.Stock - item.Quantity
		if err := s.productRepo.UpdateStock(ctx, item.ProductID, newStock); err != nil {
			return fmt.Errorf("failed to reduce stock for product %s: %w", item.ProductID, err)
		}

		// Update inventory if exists
		inventory, err := s.inventoryRepo.GetByProductID(ctx, item.ProductID)
		if err == nil {
			// Inventory exists, update it
			inventory.QuantityOnHand -= item.Quantity
			inventory.QuantityAvailable = inventory.QuantityOnHand - inventory.QuantityReserved
			if err := s.inventoryRepo.Update(ctx, inventory); err != nil {
				// Log warning but don't fail - product stock is already updated
				fmt.Printf("Warning: Failed to update inventory for product %s: %v\n", item.ProductID, err)
			}
		}

		fmt.Printf("✅ Reduced stock for product %s: %d -> %d\n", 
			product.Name, product.Stock, newStock)
	}

	return nil
}

// ReduceStockForOrder reduces stock for order items when payment is confirmed
func (s *simpleStockService) ReduceStockForOrder(ctx context.Context, items []entities.OrderItem) error {
	for _, item := range items {
		// Get current product
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product %s: %w", item.ProductID, err)
		}

		// Check stock availability one more time (race condition protection)
		if product.Stock < item.Quantity {
			return fmt.Errorf("insufficient stock for product %s: available=%d, requested=%d",
				product.Name, product.Stock, item.Quantity)
		}

		// Reduce product stock
		newStock := product.Stock - item.Quantity
		if err := s.productRepo.UpdateStock(ctx, item.ProductID, newStock); err != nil {
			return fmt.Errorf("failed to reduce stock for product %s: %w", item.ProductID, err)
		}

		// Update inventory if exists
		inventory, err := s.inventoryRepo.GetByProductID(ctx, item.ProductID)
		if err == nil {
			// Inventory exists, update it
			inventory.QuantityOnHand -= item.Quantity
			inventory.QuantityAvailable = inventory.QuantityOnHand - inventory.QuantityReserved
			if err := s.inventoryRepo.Update(ctx, inventory); err != nil {
				// Log warning but don't fail - product stock is already updated
				fmt.Printf("Warning: Failed to update inventory for product %s: %v\n", item.ProductID, err)
			}
		}

		fmt.Printf("✅ Reduced stock for product %s: %d -> %d\n",
			product.Name, product.Stock, newStock)
	}

	return nil
}

// RestoreStock restores stock for order items when order is cancelled/refunded
func (s *simpleStockService) RestoreStock(ctx context.Context, items []entities.OrderItem) error {
	for _, item := range items {
		// Get current product
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product %s: %w", item.ProductID, err)
		}

		// Restore product stock
		newStock := product.Stock + item.Quantity
		if err := s.productRepo.UpdateStock(ctx, item.ProductID, newStock); err != nil {
			return fmt.Errorf("failed to restore stock for product %s: %w", item.ProductID, err)
		}

		// Update inventory if exists
		inventory, err := s.inventoryRepo.GetByProductID(ctx, item.ProductID)
		if err == nil {
			// Inventory exists, update it
			inventory.QuantityOnHand += item.Quantity
			inventory.QuantityAvailable = inventory.QuantityOnHand - inventory.QuantityReserved
			if err := s.inventoryRepo.Update(ctx, inventory); err != nil {
				// Log warning but don't fail - product stock is already updated
				fmt.Printf("Warning: Failed to update inventory for product %s: %v\n", item.ProductID, err)
			}
		}

		fmt.Printf("✅ Restored stock for product %s: %d -> %d\n", 
			product.Name, product.Stock, newStock)
	}

	return nil
}

// GetAvailableStock gets available stock for a product
func (s *simpleStockService) GetAvailableStock(ctx context.Context, productID uuid.UUID) (int, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return 0, fmt.Errorf("failed to get product %s: %w", productID, err)
	}

	return product.Stock, nil
}
