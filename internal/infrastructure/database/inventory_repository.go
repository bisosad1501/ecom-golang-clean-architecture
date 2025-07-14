package database

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type inventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository creates a new inventory repository
func NewInventoryRepository(db *gorm.DB) repositories.InventoryRepository {
	return &inventoryRepository{db: db}
}

// Create creates a new inventory record
func (r *inventoryRepository) Create(ctx context.Context, inventory *entities.Inventory) error {
	return r.db.WithContext(ctx).Create(inventory).Error
}

// GetByID gets an inventory by ID
func (r *inventoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Inventory, error) {
	var inventory entities.Inventory
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Warehouse").
		First(&inventory, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

// GetByProductAndWarehouse gets inventory by product and warehouse
func (r *inventoryRepository) GetByProductAndWarehouse(ctx context.Context, productID, warehouseID uuid.UUID) (*entities.Inventory, error) {
	var inventory entities.Inventory
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Warehouse").
		First(&inventory, "product_id = ? AND warehouse_id = ?", productID, warehouseID).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

// Update updates an inventory record
func (r *inventoryRepository) Update(ctx context.Context, inventory *entities.Inventory) error {
	inventory.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(inventory).Error
}

// UpdateStock updates stock levels
func (r *inventoryRepository) UpdateStock(ctx context.Context, inventoryID uuid.UUID, quantityChange int, reason string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update inventory stock
		err := tx.Model(&entities.Inventory{}).
			Where("id = ?", inventoryID).
			Updates(map[string]interface{}{
				"quantity_on_hand": gorm.Expr("quantity_on_hand + ?", quantityChange),
				"updated_at":       time.Now(),
			}).Error
		if err != nil {
			return err
		}

		// Update available quantity
		return tx.Model(&entities.Inventory{}).
			Where("id = ?", inventoryID).
			Update("quantity_available", gorm.Expr("quantity_on_hand - quantity_reserved")).Error
	})
}

// SyncWithProductStock synchronizes inventory quantity with product stock
func (r *inventoryRepository) SyncWithProductStock(ctx context.Context, inventoryID uuid.UUID, productStock int, reason string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get current inventory state
		var inventory entities.Inventory
		if err := tx.First(&inventory, "id = ?", inventoryID).Error; err != nil {
			return err
		}

		// Calculate the difference
		quantityDifference := productStock - inventory.QuantityOnHand

		// Update inventory to match product stock
		err := tx.Model(&entities.Inventory{}).
			Where("id = ?", inventoryID).
			Updates(map[string]interface{}{
				"quantity_on_hand":   productStock,
				"quantity_available": gorm.Expr("? - quantity_reserved", productStock),
				"last_movement_at":   time.Now(),
				"updated_at":         time.Now(),
			}).Error
		if err != nil {
			return err
		}

		// Create movement record for tracking if there's a difference
		if quantityDifference != 0 {
			movement := &entities.InventoryMovement{
				ID:             uuid.New(),
				InventoryID:    inventoryID,
				Type:           entities.InventoryMovementTypeAdjust,
				Reason:         entities.InventoryMovementReason(reason),
				Quantity:       quantityDifference,
				UnitCost:       inventory.AverageCost,
				TotalCost:      float64(quantityDifference) * inventory.AverageCost,
				QuantityBefore: inventory.QuantityOnHand,
				QuantityAfter:  productStock,
				ReferenceType:  "product_stock_sync",
				Notes:          "Synchronized with product stock after order confirmation",
				CreatedAt:      time.Now(),
			}

			if err := tx.Create(movement).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// ReserveStock reserves stock
func (r *inventoryRepository) ReserveStock(ctx context.Context, inventoryID uuid.UUID, quantity int) error {
	return r.db.WithContext(ctx).Model(&entities.Inventory{}).
		Where("id = ? AND quantity_available >= ?", inventoryID, quantity).
		Updates(map[string]interface{}{
			"quantity_reserved":  gorm.Expr("quantity_reserved + ?", quantity),
			"quantity_available": gorm.Expr("quantity_available - ?", quantity),
			"updated_at":         time.Now(),
		}).Error
}

// ReleaseReservation releases reserved stock
func (r *inventoryRepository) ReleaseReservation(ctx context.Context, inventoryID uuid.UUID, quantity int) error {
	return r.db.WithContext(ctx).Model(&entities.Inventory{}).
		Where("id = ? AND quantity_reserved >= ?", inventoryID, quantity).
		Updates(map[string]interface{}{
			"quantity_reserved":  gorm.Expr("quantity_reserved - ?", quantity),
			"quantity_available": gorm.Expr("quantity_available + ?", quantity),
			"updated_at":         time.Now(),
		}).Error
}

// GetLowStockItems gets items with low stock
func (r *inventoryRepository) GetLowStockItems(ctx context.Context, limit, offset int) ([]*entities.Inventory, error) {
	var inventories []*entities.Inventory
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Warehouse").
		Where("quantity_available <= reorder_level AND quantity_available > 0").
		Order("quantity_available ASC").
		Limit(limit).
		Offset(offset).
		Find(&inventories).Error
	return inventories, err
}

// CountLowStockItems counts items with low stock
func (r *inventoryRepository) CountLowStockItems(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Inventory{}).
		Where("quantity_available <= reorder_level AND quantity_available > 0").
		Count(&count).Error
	return count, err
}

// GetOutOfStockItems gets items that are out of stock
func (r *inventoryRepository) GetOutOfStockItems(ctx context.Context, limit, offset int) ([]*entities.Inventory, error) {
	var inventories []*entities.Inventory
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Warehouse").
		Where("quantity_available = 0").
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&inventories).Error
	return inventories, err
}

// GetInventoryByWarehouse gets inventory by warehouse
func (r *inventoryRepository) GetInventoryByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]*entities.Inventory, error) {
	var inventories []*entities.Inventory
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Warehouse").
		Where("warehouse_id = ?", warehouseID).
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&inventories).Error
	return inventories, err
}

// List lists inventories with filters
func (r *inventoryRepository) List(ctx context.Context, filters repositories.InventoryFilters) ([]*entities.Inventory, error) {
	var inventories []*entities.Inventory
	query := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Warehouse")

	if filters.ProductID != nil {
		query = query.Where("product_id = ?", *filters.ProductID)
	}

	if filters.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filters.WarehouseID)
	}

	if filters.LowStock {
		query = query.Where("quantity_available <= reorder_level")
	}

	if filters.OutOfStock {
		query = query.Where("quantity_available = 0")
	}

	err := query.Find(&inventories).Error
	return inventories, err
}

// CreateMovement creates an inventory movement record
func (r *inventoryRepository) CreateMovement(ctx context.Context, movement *entities.InventoryMovement) error {
	return r.db.WithContext(ctx).Create(movement).Error
}

// GetMovements gets inventory movements
func (r *inventoryRepository) GetMovements(ctx context.Context, inventoryID uuid.UUID, limit, offset int) ([]*entities.InventoryMovement, error) {
	var movements []*entities.InventoryMovement
	err := r.db.WithContext(ctx).
		Where("inventory_id = ?", inventoryID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&movements).Error
	return movements, err
}

// GetMovementsByDateRange gets movements by date range
func (r *inventoryRepository) GetMovementsByDateRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*entities.InventoryMovement, error) {
	var movements []*entities.InventoryMovement
	err := r.db.WithContext(ctx).
		Where("created_at BETWEEN ? AND ?", from, to).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&movements).Error
	return movements, err
}

// CreateAlert creates a stock alert
func (r *inventoryRepository) CreateAlert(ctx context.Context, alert *entities.StockAlert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

// GetActiveAlerts gets active alerts
func (r *inventoryRepository) GetActiveAlerts(ctx context.Context, limit, offset int) ([]*entities.StockAlert, error) {
	var alerts []*entities.StockAlert
	err := r.db.WithContext(ctx).
		Where("status = ?", entities.StockAlertStatusActive).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&alerts).Error
	return alerts, err
}

// ResolveAlert resolves an alert
func (r *inventoryRepository) ResolveAlert(ctx context.Context, alertID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.StockAlert{}).
		Where("id = ?", alertID).
		Updates(map[string]interface{}{
			"status":      entities.StockAlertStatusResolved,
			"resolved_at": time.Now(),
			"updated_at":  time.Now(),
		}).Error
}

// Delete deletes an inventory record
func (r *inventoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Inventory{}, "id = ?", id).Error
}

// Count counts inventory records with filters
func (r *inventoryRepository) Count(ctx context.Context, filters repositories.InventoryFilters) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Inventory{})

	if filters.ProductID != nil {
		query = query.Where("product_id = ?", *filters.ProductID)
	}

	if filters.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filters.WarehouseID)
	}

	if filters.LowStock {
		query = query.Where("quantity_available <= reorder_level")
	}

	if filters.OutOfStock {
		query = query.Where("quantity_available = 0")
	}

	err := query.Count(&count).Error
	return count, err
}

// CountOutOfStockItems counts out of stock items
func (r *inventoryRepository) CountOutOfStockItems(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Inventory{}).
		Where("quantity_available = 0").
		Count(&count).Error
	return count, err
}

// GetAlertsByInventory gets alerts for a specific inventory
func (r *inventoryRepository) GetAlertsByInventory(ctx context.Context, inventoryID uuid.UUID) ([]*entities.StockAlert, error) {
	var alerts []*entities.StockAlert
	err := r.db.WithContext(ctx).
		Where("inventory_id = ?", inventoryID).
		Order("created_at DESC").
		Find(&alerts).Error
	return alerts, err
}

// GetAvailableStock gets available stock for a product across all warehouses
func (r *inventoryRepository) GetAvailableStock(ctx context.Context, productID uuid.UUID) (int, error) {
	var totalStock int64
	err := r.db.WithContext(ctx).
		Model(&entities.Inventory{}).
		Select("SUM(quantity_available)").
		Where("product_id = ?", productID).
		Scan(&totalStock).Error
	return int(totalStock), err
}

// GetByProductID gets first inventory record by product ID
func (r *inventoryRepository) GetByProductID(ctx context.Context, productID uuid.UUID) (*entities.Inventory, error) {
	var inventory entities.Inventory
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		First(&inventory).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

// GetMovementReport gets inventory movement report
func (r *inventoryRepository) GetMovementReport(ctx context.Context, filters repositories.MovementReportFilters) (*repositories.MovementReport, error) {
	var report repositories.MovementReport
	var movements []repositories.MovementReportItem

	query := r.db.WithContext(ctx).
		Table("inventory_movements").
		Select("inventory_movements.created_at as date, products.name as product_name, products.sku, inventory_movements.movement_type as type, inventory_movements.reason, inventory_movements.quantity, inventory_movements.reference_id").
		Joins("JOIN inventories ON inventory_movements.inventory_id = inventories.id").
		Joins("JOIN products ON inventories.product_id = products.id")

	if filters.InventoryID != nil {
		query = query.Where("inventory_movements.inventory_id = ?", *filters.InventoryID)
	}

	if filters.WarehouseID != nil {
		query = query.Where("inventories.warehouse_id = ?", *filters.WarehouseID)
	}

	if filters.MovementType != nil {
		query = query.Where("inventory_movements.movement_type = ?", *filters.MovementType)
	}

	if filters.DateFrom != nil {
		query = query.Where("inventory_movements.created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("inventory_movements.created_at <= ?", *filters.DateTo)
	}

	err := query.Order("inventory_movements.created_at DESC").
		Scan(&movements).Error
	if err != nil {
		return nil, err
	}

	// Calculate totals
	var inboundTotal, outboundTotal int
	for _, movement := range movements {
		if movement.Type == entities.InventoryMovementTypeIn {
			inboundTotal += movement.Quantity
		} else if movement.Type == entities.InventoryMovementTypeOut {
			outboundTotal += movement.Quantity
		}
	}

	report.TotalMovements = int64(len(movements))
	report.InboundTotal = inboundTotal
	report.OutboundTotal = outboundTotal
	report.NetChange = inboundTotal - outboundTotal
	report.Movements = movements

	return &report, nil
}

// GetStockReport gets stock report
func (r *inventoryRepository) GetStockReport(ctx context.Context, filters repositories.StockReportFilters) (*repositories.StockReport, error) {
	var report repositories.StockReport
	var items []repositories.StockReportItem

	query := r.db.WithContext(ctx).
		Table("inventories").
		Select("products.name as product_name, products.sku, inventories.quantity_available, inventories.quantity_reserved, inventories.reorder_level, inventories.unit_cost, (inventories.quantity_available * inventories.unit_cost) as total_value, CASE WHEN inventories.quantity_available <= inventories.reorder_level THEN 'low_stock' WHEN inventories.quantity_available = 0 THEN 'out_of_stock' ELSE 'in_stock' END as status").
		Joins("JOIN products ON inventories.product_id = products.id")

	if filters.WarehouseID != nil {
		query = query.Where("inventories.warehouse_id = ?", *filters.WarehouseID)
	}

	if filters.CategoryID != nil {
		query = query.Where("products.category_id = ?", *filters.CategoryID)
	}

	if filters.DateFrom != nil {
		query = query.Where("inventories.updated_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("inventories.updated_at <= ?", *filters.DateTo)
	}

	err := query.Scan(&items).Error
	if err != nil {
		return nil, err
	}

	// Calculate totals
	var totalValue float64
	var lowStockCount, outOfStockCount int64
	for _, item := range items {
		totalValue += item.TotalValue
		if item.Status == "low_stock" {
			lowStockCount++
		} else if item.Status == "out_of_stock" {
			outOfStockCount++
		}
	}

	report.TotalItems = int64(len(items))
	report.TotalValue = totalValue
	report.LowStockItems = lowStockCount
	report.OutOfStockItems = outOfStockCount
	report.Items = items

	return &report, nil
}

// TransferStock transfers stock between warehouses
func (r *inventoryRepository) TransferStock(ctx context.Context, fromInventoryID, toInventoryID uuid.UUID, quantity int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Decrease from source inventory
		err := tx.Model(&entities.Inventory{}).
			Where("id = ? AND quantity_available >= ?", fromInventoryID, quantity).
			Update("quantity_available", gorm.Expr("quantity_available - ?", quantity)).Error
		if err != nil {
			return err
		}

		// Increase to destination inventory
		err = tx.Model(&entities.Inventory{}).
			Where("id = ?", toInventoryID).
			Update("quantity_available", gorm.Expr("quantity_available + ?", quantity)).Error
		if err != nil {
			return err
		}

		// Create movement records
		outMovement := &entities.InventoryMovement{
			ID:          uuid.New(),
			InventoryID: fromInventoryID,
			Type:        entities.InventoryMovementTypeOut,
			Quantity:    quantity,
			Reason:      entities.InventoryReasonTransfer,
			ReferenceID: &toInventoryID,
			CreatedAt:   time.Now(),
		}
		err = tx.Create(outMovement).Error
		if err != nil {
			return err
		}

		inMovement := &entities.InventoryMovement{
			ID:          uuid.New(),
			InventoryID: toInventoryID,
			Type:        entities.InventoryMovementTypeIn,
			Quantity:    quantity,
			Reason:      entities.InventoryReasonTransfer,
			ReferenceID: &fromInventoryID,
			CreatedAt:   time.Now(),
		}
		return tx.Create(inMovement).Error
	})
}
