package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
)

// SimpleStockService handles stock management with Inventory as single source of truth
// Product.Stock is now just a cached value synced from Inventory.QuantityOnHand
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
// Uses Inventory as source of truth instead of Product.Stock
func (s *simpleStockService) CheckStockAvailability(ctx context.Context, items []entities.CartItem) error {
	for _, item := range items {
		// Get current product for availability check
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product %s: %w", item.ProductID, err)
		}

		// Check if product is available
		if !product.IsAvailable() {
			return fmt.Errorf("product %s is not available", product.Name)
		}

		// Get inventory (source of truth for stock)
		inventory, err := s.inventoryRepo.GetByProductID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get inventory for product %s: %w", item.ProductID, err)
		}

		// Check stock availability from inventory
		if inventory.QuantityAvailable < item.Quantity {
			return fmt.Errorf("insufficient stock for product %s: available=%d, requested=%d",
				product.Name, inventory.QuantityAvailable, item.Quantity)
		}
	}

	return nil
}

// ReduceStock reduces stock for cart items when payment is successful
// Uses Inventory as source of truth, then syncs Product.Stock
func (s *simpleStockService) ReduceStock(ctx context.Context, items []entities.CartItem) error {
	for _, item := range items {
		// Get current product for name
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product %s: %w", item.ProductID, err)
		}

		// Get inventory (source of truth)
		inventory, err := s.inventoryRepo.GetByProductID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get inventory for product %s: %w", item.ProductID, err)
		}

		// Check stock availability one more time (race condition protection)
		if inventory.QuantityAvailable < item.Quantity {
			return fmt.Errorf("insufficient stock for product %s: available=%d, requested=%d",
				product.Name, inventory.QuantityAvailable, item.Quantity)
		}

		// Reduce inventory stock (source of truth)
		oldQuantity := inventory.QuantityOnHand
		inventory.QuantityOnHand -= item.Quantity
		inventory.QuantityAvailable = inventory.QuantityOnHand - inventory.QuantityReserved

		if err := s.inventoryRepo.Update(ctx, inventory); err != nil {
			return fmt.Errorf("failed to update inventory for product %s: %w", item.ProductID, err)
		}

		// Sync product stock from inventory (Product.Stock is now cached value)
		if err := s.productRepo.UpdateStock(ctx, item.ProductID, inventory.QuantityOnHand); err != nil {
			// Log warning but don't fail - inventory is already updated (source of truth)
			fmt.Printf("Warning: Failed to sync product stock for %s: %v\n", item.ProductID, err)
		}

		fmt.Printf("✅ Reduced stock for product %s: %d -> %d (Inventory: %d available)\n",
			product.Name, oldQuantity, inventory.QuantityOnHand, inventory.QuantityAvailable)
	}

	return nil
}

// ReduceStockForOrder reduces stock for order items when payment is confirmed
// Uses Inventory as source of truth, then syncs Product.Stock
func (s *simpleStockService) ReduceStockForOrder(ctx context.Context, items []entities.OrderItem) error {
	for _, item := range items {
		// Get current product for name
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product %s: %w", item.ProductID, err)
		}

		// Get inventory (source of truth)
		inventory, err := s.inventoryRepo.GetByProductID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get inventory for product %s: %w", item.ProductID, err)
		}

		// Check stock availability one more time (race condition protection)
		if inventory.QuantityAvailable < item.Quantity {
			return fmt.Errorf("insufficient stock for product %s: available=%d, requested=%d",
				product.Name, inventory.QuantityAvailable, item.Quantity)
		}

		// Reduce inventory stock (source of truth)
		oldQuantity := inventory.QuantityOnHand
		inventory.QuantityOnHand -= item.Quantity
		inventory.QuantityAvailable = inventory.QuantityOnHand - inventory.QuantityReserved

		if err := s.inventoryRepo.Update(ctx, inventory); err != nil {
			return fmt.Errorf("failed to update inventory for product %s: %w", item.ProductID, err)
		}

		// Sync product stock from inventory (Product.Stock is now cached value)
		if err := s.productRepo.UpdateStock(ctx, item.ProductID, inventory.QuantityOnHand); err != nil {
			// Log warning but don't fail - inventory is already updated (source of truth)
			fmt.Printf("Warning: Failed to sync product stock for %s: %v\n", item.ProductID, err)
		}

		fmt.Printf("✅ Reduced stock for product %s: %d -> %d (Inventory: %d available)\n",
			product.Name, oldQuantity, inventory.QuantityOnHand, inventory.QuantityAvailable)
	}

	return nil
}

// RestoreStock restores stock for order items when order is cancelled/refunded
// Uses Inventory as source of truth, then syncs Product.Stock
func (s *simpleStockService) RestoreStock(ctx context.Context, items []entities.OrderItem) error {
	for _, item := range items {
		// Get current product for name
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product %s: %w", item.ProductID, err)
		}

		// Get inventory (source of truth)
		inventory, err := s.inventoryRepo.GetByProductID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get inventory for product %s: %w", item.ProductID, err)
		}

		// Restore inventory stock (source of truth)
		oldQuantity := inventory.QuantityOnHand
		inventory.QuantityOnHand += item.Quantity
		inventory.QuantityAvailable = inventory.QuantityOnHand - inventory.QuantityReserved

		if err := s.inventoryRepo.Update(ctx, inventory); err != nil {
			return fmt.Errorf("failed to update inventory for product %s: %w", item.ProductID, err)
		}

		// Sync product stock from inventory (Product.Stock is now cached value)
		if err := s.productRepo.UpdateStock(ctx, item.ProductID, inventory.QuantityOnHand); err != nil {
			// Log warning but don't fail - inventory is already updated (source of truth)
			fmt.Printf("Warning: Failed to sync product stock for %s: %v\n", item.ProductID, err)
		}

		fmt.Printf("✅ Restored stock for product %s: %d -> %d (Inventory: %d available)\n",
			product.Name, oldQuantity, inventory.QuantityOnHand, inventory.QuantityAvailable)
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
