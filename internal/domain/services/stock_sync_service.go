package services

import (
	"context"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
)

// StockSyncService handles synchronization between Product and Inventory stock
type StockSyncService interface {
	// SyncProductWithInventory syncs product stock with inventory (Inventory is source of truth)
	SyncProductWithInventory(ctx context.Context, productID uuid.UUID) error

	// SyncAllProducts syncs all products with their inventory
	SyncAllProducts(ctx context.Context) error

	// ValidateStockConsistency validates that product stock matches inventory
	ValidateStockConsistency(ctx context.Context, productID uuid.UUID) (bool, error)

	// AutoSyncIfNeeded automatically syncs product with inventory if inconsistency is detected
	AutoSyncIfNeeded(ctx context.Context, productID uuid.UUID) error

	// ValidateAndFixAllInconsistencies validates and fixes all stock inconsistencies
	ValidateAndFixAllInconsistencies(ctx context.Context) error
}

type stockSyncService struct {
	productRepo   repositories.ProductRepository
	inventoryRepo repositories.InventoryRepository
}

// NewStockSyncService creates a new stock sync service
func NewStockSyncService(
	productRepo repositories.ProductRepository,
	inventoryRepo repositories.InventoryRepository,
) StockSyncService {
	return &stockSyncService{
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
	}
}

// SyncProductWithInventory syncs product stock with inventory (Inventory is source of truth)
func (s *stockSyncService) SyncProductWithInventory(ctx context.Context, productID uuid.UUID) error {
	// Get inventory (source of truth)
	inventory, err := s.inventoryRepo.GetByProductID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get inventory for product %s: %w", productID, err)
	}

	// Get current product
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get product %s: %w", productID, err)
	}

	// Check if sync is needed
	if product.Stock == inventory.QuantityOnHand {
		return nil // Already in sync
	}

	// Validate inventory data consistency
	if inventory.QuantityOnHand < 0 {
		return fmt.Errorf("invalid inventory data for product %s: negative quantity on hand (%d)",
			productID, inventory.QuantityOnHand)
	}

	if inventory.QuantityReserved < 0 {
		return fmt.Errorf("invalid inventory data for product %s: negative quantity reserved (%d)",
			productID, inventory.QuantityReserved)
	}

	if inventory.QuantityReserved > inventory.QuantityOnHand {
		fmt.Printf("⚠️ Warning: Reserved quantity (%d) exceeds on-hand quantity (%d) for product %s\n",
			inventory.QuantityReserved, inventory.QuantityOnHand, productID)
		// Fix the inconsistency by recalculating available stock
		inventory.QuantityAvailable = 0
		if err := s.inventoryRepo.Update(ctx, inventory); err != nil {
			fmt.Printf("❌ Failed to fix inventory inconsistency for product %s: %v\n", productID, err)
		}
	}

	// Update product stock to match inventory (Product.Stock is cached value)
	oldStock := product.Stock
	if err := s.productRepo.UpdateStock(ctx, productID, inventory.QuantityOnHand); err != nil {
		return fmt.Errorf("failed to sync product stock from inventory: %w", err)
	}

	fmt.Printf("✅ Synced product %s stock: %d -> %d (from inventory source of truth)\n",
		productID, oldStock, inventory.QuantityOnHand)

	return nil
}

// SyncAllProducts syncs all products with their inventory
func (s *stockSyncService) SyncAllProducts(ctx context.Context) error {
	// Get all products
	products, err := s.productRepo.List(ctx, 0, 1000) // Get first 1000 products
	if err != nil {
		return fmt.Errorf("failed to get products: %w", err)
	}

	var syncErrors []error
	syncedCount := 0

	for _, product := range products {
		if err := s.SyncProductWithInventory(ctx, product.ID); err != nil {
			syncErrors = append(syncErrors, fmt.Errorf("product %s: %w", product.ID, err))
		} else {
			syncedCount++
		}
	}

	fmt.Printf("✅ Synced %d products successfully\n", syncedCount)

	if len(syncErrors) > 0 {
		fmt.Printf("❌ Failed to sync %d products:\n", len(syncErrors))
		for _, err := range syncErrors {
			fmt.Printf("  - %v\n", err)
		}
		return fmt.Errorf("failed to sync %d products", len(syncErrors))
	}

	return nil
}

// ValidateStockConsistency validates that product stock matches inventory
func (s *stockSyncService) ValidateStockConsistency(ctx context.Context, productID uuid.UUID) (bool, error) {
	// Get product
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return false, fmt.Errorf("failed to get product %s: %w", productID, err)
	}

	// Get inventory
	inventory, err := s.inventoryRepo.GetByProductID(ctx, productID)
	if err != nil {
		return false, fmt.Errorf("failed to get inventory for product %s: %w", productID, err)
	}

	// Check consistency
	isConsistent := product.Stock == inventory.QuantityOnHand

	if !isConsistent {
		fmt.Printf("❌ Stock inconsistency for product %s: Product.Stock=%d, Inventory.QuantityOnHand=%d\n",
			productID, product.Stock, inventory.QuantityOnHand)
	}

	return isConsistent, nil
}

// AutoSyncIfNeeded automatically syncs product with inventory if inconsistency is detected
func (s *stockSyncService) AutoSyncIfNeeded(ctx context.Context, productID uuid.UUID) error {
	isConsistent, err := s.ValidateStockConsistency(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to validate stock consistency: %w", err)
	}

	if !isConsistent {
		fmt.Printf("🔄 Stock inconsistency detected for product %s, auto-syncing...\n", productID)
		return s.SyncProductWithInventory(ctx, productID)
	}

	return nil
}

// ValidateAndFixAllInconsistencies validates and fixes all stock inconsistencies
func (s *stockSyncService) ValidateAndFixAllInconsistencies(ctx context.Context) error {
	fmt.Printf("🔍 Starting validation and fix of all stock inconsistencies...\n")
	startTime := time.Now()

	// Get all products
	products, err := s.productRepo.List(ctx, 0, 1000) // Get first 1000 products
	if err != nil {
		return fmt.Errorf("failed to get products: %w", err)
	}

	var inconsistentProducts []uuid.UUID
	var fixErrors []error
	fixedCount := 0

	for _, product := range products {
		isConsistent, err := s.ValidateStockConsistency(ctx, product.ID)
		if err != nil {
			fixErrors = append(fixErrors, fmt.Errorf("product %s validation error: %w", product.ID, err))
			continue
		}

		if !isConsistent {
			inconsistentProducts = append(inconsistentProducts, product.ID)
			if err := s.SyncProductWithInventory(ctx, product.ID); err != nil {
				fixErrors = append(fixErrors, fmt.Errorf("product %s sync error: %w", product.ID, err))
			} else {
				fixedCount++
			}
		}
	}

	duration := time.Since(startTime)

	if len(inconsistentProducts) == 0 {
		fmt.Printf("✅ All products are consistent (checked %d products in %v)\n", len(products), duration)
		return nil
	}

	fmt.Printf("🔧 Fixed %d/%d inconsistent products in %v\n", fixedCount, len(inconsistentProducts), duration)

	if len(fixErrors) > 0 {
		fmt.Printf("❌ %d errors occurred during fix process\n", len(fixErrors))
		return fmt.Errorf("validation and fix completed with %d errors", len(fixErrors))
	}

	return nil
}
