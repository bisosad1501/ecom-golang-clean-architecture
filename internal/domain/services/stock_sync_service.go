package services

import (
	"context"
	"fmt"

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

	// Update product stock to match inventory
	oldStock := product.Stock
	if err := s.productRepo.UpdateStock(ctx, productID, inventory.QuantityOnHand); err != nil {
		return fmt.Errorf("failed to update product stock for %s: %w", productID, err)
	}

	fmt.Printf("✅ Synced product %s stock: %d -> %d (from inventory)\n", 
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
