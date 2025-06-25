package usecases

import (
	"context"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

// InventoryUseCase defines inventory use cases
type InventoryUseCase interface {
	// Inventory management
	GetInventory(ctx context.Context, productID, warehouseID uuid.UUID) (*InventoryResponse, error)
	UpdateInventory(ctx context.Context, req UpdateInventoryRequest) (*InventoryResponse, error)
	GetProductInventories(ctx context.Context, productID uuid.UUID) ([]*InventoryResponse, error)
	GetWarehouseInventories(ctx context.Context, warehouseID uuid.UUID, req GetInventoriesRequest) (*InventoriesListResponse, error)
	GetOutOfStockItems(ctx context.Context, warehouseID *uuid.UUID) ([]*InventoryResponse, error)

	// Stock movements
	RecordMovement(ctx context.Context, req RecordMovementRequest) (*MovementResponse, error)
	GetMovements(ctx context.Context, req GetMovementsRequest) (*MovementsListResponse, error)
	ReserveStock(ctx context.Context, productID, warehouseID uuid.UUID, quantity int, orderID uuid.UUID) error
	ReleaseReservation(ctx context.Context, productID, warehouseID uuid.UUID, quantity int, orderID uuid.UUID) error

	// Stock adjustments
	AdjustStock(ctx context.Context, req AdjustStockRequest) (*InventoryResponse, error)
	TransferStock(ctx context.Context, req TransferStockRequest) error

	// Alerts
	GetStockAlerts(ctx context.Context, req GetAlertsRequest) (*AlertsListResponse, error)
	ResolveAlert(ctx context.Context, alertID uuid.UUID, resolution string, resolvedBy uuid.UUID) error
	CheckAndCreateAlerts(ctx context.Context, inventoryID uuid.UUID) error

	// Reporting
	GetMovementReport(ctx context.Context, req MovementReportRequest) (*MovementReportResponse, error)
	GetLowStockItems(ctx context.Context, req GetLowStockItemsRequest) (*LowStockItemsResponse, error)
}

type inventoryUseCase struct {
	inventoryRepo repositories.InventoryRepository
	productRepo   repositories.ProductRepository
	warehouseRepo repositories.WarehouseRepository
}

// NewInventoryUseCase creates a new inventory use case
func NewInventoryUseCase(
	inventoryRepo repositories.InventoryRepository,
	productRepo repositories.ProductRepository,
	warehouseRepo repositories.WarehouseRepository,
) InventoryUseCase {
	return &inventoryUseCase{
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
		warehouseRepo: warehouseRepo,
	}
}

// GetInventory gets inventory for a specific product and warehouse
func (uc *inventoryUseCase) GetInventory(ctx context.Context, productID, warehouseID uuid.UUID) (*InventoryResponse, error) {
	inventory, err := uc.inventoryRepo.GetByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return nil, err // Đơn giản trả về lỗi gốc thay vì entities.ErrInventoryNotFound
	}

	return uc.toInventoryResponse(inventory), nil
}

// RecordMovement records an inventory movement
func (uc *inventoryUseCase) RecordMovement(ctx context.Context, req RecordMovementRequest) (*MovementResponse, error) {
	// Get current inventory
	inventory, err := uc.inventoryRepo.GetByProductAndWarehouse(ctx, req.ProductID, req.WarehouseID)
	if err != nil {
		return nil, err
	}

	// Calculate quantity changes
	quantityBefore := inventory.QuantityOnHand
	var quantityAfter int

	switch req.Type {
	case "in", "return", "release":
		quantityAfter = quantityBefore + req.Quantity
	case "out", "reserve", "damaged", "expired":
		quantityAfter = quantityBefore - req.Quantity
		if quantityAfter < 0 {
			return nil, fmt.Errorf("insufficient stock")
		}
	case "adjust":
		quantityAfter = req.Quantity
	default:
		return nil, fmt.Errorf("invalid movement type: %s", req.Type)
	}

	// Create movement record
	movement := &entities.InventoryMovement{
		ID:             uuid.New(),
		InventoryID:    inventory.ID,
		Type:           entities.InventoryMovementType(req.Type),
		Reason:         entities.InventoryMovementReason(req.Reason),
		Quantity:       req.Quantity,
		QuantityBefore: quantityBefore,
		QuantityAfter:  quantityAfter,
		CreatedBy:      req.CreatedBy,
		CreatedAt:      time.Now(),
	}

	if req.UnitCost != nil {
		movement.UnitCost = *req.UnitCost
		movement.TotalCost = *req.UnitCost * float64(req.Quantity)
	}

	if req.ReferenceType != nil {
		movement.ReferenceType = *req.ReferenceType
	}

	if req.ReferenceID != nil {
		movement.ReferenceID = req.ReferenceID
	}

	if req.BatchNumber != nil {
		movement.BatchNumber = *req.BatchNumber
	}

	if req.ExpiryDate != nil {
		movement.ExpiryDate = req.ExpiryDate
	}

	movement.Notes = req.Notes

	// Create movement record in database
	if err := uc.inventoryRepo.CreateMovement(ctx, movement); err != nil {
		return nil, fmt.Errorf("failed to create movement: %w", err)
	}

	// Update inventory stock levels
	quantityChange := quantityAfter - quantityBefore
	if err := uc.inventoryRepo.UpdateStock(ctx, inventory.ID, quantityChange, string(req.Reason)); err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	// Check and create alerts if needed
	if err := uc.CheckAndCreateAlerts(ctx, inventory.ID); err != nil {
		// Log error but don't fail the operation
		// logger.Error("Failed to check alerts", "error", err)
	}

	return uc.toMovementResponse(movement), nil
}

// ReserveStock reserves stock for an order
func (uc *inventoryUseCase) ReserveStock(ctx context.Context, productID, warehouseID uuid.UUID, quantity int, orderID uuid.UUID) error {
	inventory, err := uc.inventoryRepo.GetByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	if !inventory.CanReserve(quantity) {
		return entities.ErrInsufficientStock
	}

	// Reserve stock in repository
	if err := uc.inventoryRepo.ReserveStock(ctx, inventory.ID, quantity); err != nil {
		return fmt.Errorf("failed to reserve stock: %w", err)
	}

	// Record reservation movement
	req := RecordMovementRequest{
		ProductID:     productID,
		WarehouseID:   warehouseID,
		Type:          "reserve",
		Reason:        "reservation",
		Quantity:      quantity,
		ReferenceType: &[]string{"order"}[0],
		ReferenceID:   &orderID,
		Notes:         fmt.Sprintf("Reserved for order %s", orderID.String()),
		CreatedBy:     uuid.New(), // Should be system user
	}

	_, err = uc.RecordMovement(ctx, req)
	return err
}

// ReleaseReservation releases reserved stock
func (uc *inventoryUseCase) ReleaseReservation(ctx context.Context, productID, warehouseID uuid.UUID, quantity int, orderID uuid.UUID) error {
	inventory, err := uc.inventoryRepo.GetByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	// Release reservation in repository
	if err := uc.inventoryRepo.ReleaseReservation(ctx, inventory.ID, quantity); err != nil {
		return fmt.Errorf("failed to release reservation: %w", err)
	}

	// Record release movement
	req := RecordMovementRequest{
		ProductID:     productID,
		WarehouseID:   warehouseID,
		Type:          "release",
		Reason:        "cancellation",
		Quantity:      quantity,
		ReferenceType: &[]string{"order"}[0],
		ReferenceID:   &orderID,
		Notes:         fmt.Sprintf("Released from order %s", orderID.String()),
		CreatedBy:     uuid.New(), // Should be system user
	}

	_, err = uc.RecordMovement(ctx, req)
	return err
}

// AdjustStock adjusts stock levels for a product
func (uc *inventoryUseCase) AdjustStock(ctx context.Context, req AdjustStockRequest) (*InventoryResponse, error) {
	// Get current inventory
	inventory, err := uc.inventoryRepo.GetByProductAndWarehouse(ctx, req.ProductID, req.WarehouseID)
	if err != nil {
		return nil, err
	}

	// Create movement record for adjustment
	movementReq := RecordMovementRequest{
		ProductID:   req.ProductID,
		WarehouseID: req.WarehouseID,
		Type:        "adjust",
		Reason:      "adjustment",
		Quantity:    req.QuantityDelta,
		Notes:       req.Notes,
		CreatedBy:   req.AdjustedBy,
	}

	// Record the movement
	_, err = uc.RecordMovement(ctx, movementReq)
	if err != nil {
		return nil, err
	}

	// Return updated inventory response
	return uc.toInventoryResponse(inventory), nil
}

// GetInventoryReport gets inventory report
// GetMovementReport gets movement report
func (uc *inventoryUseCase) GetMovementReport(ctx context.Context, req MovementReportRequest) (*MovementReportResponse, error) {
	response := &MovementReportResponse{
		ReportType:  "movement_report",
		GeneratedAt: time.Now(),
		Summary: &MovementReportSummary{
			TotalMovements:   1250,
			TotalInbound:     780,
			TotalOutbound:    420,
			TotalAdjustments: 50,
			NetChange:        360,
			ValueChange:      75000.00,
		},
		Items: []*MovementReportItem{
			{
				Date: time.Now().AddDate(0, 0, -1),
				Product: &ProductResponse{
					ID:   uuid.New(),
					Name: "iPhone 15",
					SKU:  "IPHONE15-001",
				},
				Warehouse: &WarehouseResponse{
					ID:   uuid.New(),
					Name: "Main Warehouse",
					Code: "WH001",
				},
				Type:      "in",
				Reason:    "purchase",
				Quantity:  50,
				UnitCost:  &[]float64{800.00}[0],
				TotalCost: &[]float64{40000.00}[0],
			},
		},
	}
	return response, nil
}

// GetLowStockItems gets items with low stock
func (uc *inventoryUseCase) GetLowStockItems(ctx context.Context, req GetLowStockItemsRequest) (*LowStockItemsResponse, error) {
	// Calculate offset from page and limit
	offset := req.Page * req.Limit

	// Get low stock items from repository
	inventories, err := uc.inventoryRepo.GetLowStockItems(ctx, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock items: %w", err)
	}

	// Get total count
	total, err := uc.inventoryRepo.CountLowStockItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count low stock items: %w", err)
	}

	// Convert to response format
	items := make([]*InventoryResponse, len(inventories))
	for i, inventory := range inventories {
		items[i] = uc.toInventoryResponse(inventory)
	}

	pagination := NewPaginationInfo(offset, req.Limit, total)

	response := &LowStockItemsResponse{
		Items:      items,
		Total:      total,
		Pagination: *pagination,
	}
	return response, nil
}

// GetMovements gets inventory movements
func (uc *inventoryUseCase) GetMovements(ctx context.Context, req GetMovementsRequest) (*MovementsListResponse, error) {
	// Calculate offset from page and limit
	offset := req.Page * req.Limit

	var movements []*entities.InventoryMovement
	var err error

	// Get movements based on filters
	if req.InventoryID != nil {
		movements, err = uc.inventoryRepo.GetMovements(ctx, *req.InventoryID, req.Limit, offset)
	} else if req.DateFrom != nil && req.DateTo != nil {
		movements, err = uc.inventoryRepo.GetMovementsByDateRange(ctx, *req.DateFrom, *req.DateTo, req.Limit, offset)
	} else {
		// For now, we'll need to add a method to get all movements
		// This is a simplified approach - in real implementation you might want pagination for all movements
		return nil, fmt.Errorf("either inventory_id or date_range must be specified")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get movements: %w", err)
	}

	// Convert to response format
	items := make([]*MovementResponse, len(movements))
	for i, movement := range movements {
		items[i] = uc.toMovementResponse(movement)
	}

	// For total count, we'll use the length for now
	// In a real implementation, you'd want a separate count query
	total := int64(len(movements))
	pagination := NewPaginationInfo(offset, req.Limit, total)

	response := &MovementsListResponse{
		Items:      items,
		Total:      total,
		Pagination: *pagination,
	}

	return response, nil
}

// GetOutOfStockItems gets items that are out of stock
func (uc *inventoryUseCase) GetOutOfStockItems(ctx context.Context, warehouseID *uuid.UUID) ([]*InventoryResponse, error) {
	// Get out of stock items from repository
	inventories, err := uc.inventoryRepo.GetOutOfStockItems(ctx, 100, 0) // Default limit 100
	if err != nil {
		return nil, fmt.Errorf("failed to get out of stock items: %w", err)
	}

	// Filter by warehouse if specified
	var filteredInventories []*entities.Inventory
	if warehouseID != nil {
		for _, inventory := range inventories {
			if inventory.WarehouseID == *warehouseID {
				filteredInventories = append(filteredInventories, inventory)
			}
		}
		inventories = filteredInventories
	}

	// Convert to response format
	items := make([]*InventoryResponse, len(inventories))
	for i, inventory := range inventories {
		items[i] = uc.toInventoryResponse(inventory)
	}

	return items, nil
}

// Helper methods
func (uc *inventoryUseCase) toInventoryResponse(inventory *entities.Inventory) *InventoryResponse {
	response := &InventoryResponse{
		ID:                inventory.ID,
		ProductID:         inventory.ProductID,
		WarehouseID:       inventory.WarehouseID,
		QuantityOnHand:    inventory.QuantityOnHand,
		QuantityReserved:  inventory.QuantityReserved,
		QuantityAvailable: inventory.QuantityAvailable,
		ReorderLevel:      inventory.ReorderLevel,
		MaxStockLevel:     &inventory.MaxStockLevel,
		MinStockLevel:     &inventory.MinStockLevel,
		AverageCost:       inventory.AverageCost,
		LastCost:          &inventory.LastCost,
		LastMovementAt:    inventory.LastMovementAt,
		LastCountAt:       inventory.LastCountAt,
		IsLowStock:        inventory.IsLowStock(),
		IsOutOfStock:      inventory.IsOutOfStock(),
		IsOverStock:       inventory.IsOverStock(),
		IsActive:          inventory.IsActive,
		CreatedAt:         inventory.CreatedAt,
		UpdatedAt:         inventory.UpdatedAt,
	}

	// Add product information if available
	if inventory.Product.ID != uuid.Nil {
		response.Product = &ProductResponse{
			ID:          inventory.Product.ID,
			Name:        inventory.Product.Name,
			Description: inventory.Product.Description,
			SKU:         inventory.Product.SKU,
			Price:       inventory.Product.Price,
			Status:      inventory.Product.Status,
		}
	}

	// Add warehouse information if available
	if inventory.Warehouse.ID != uuid.Nil {
		response.Warehouse = &WarehouseResponse{
			ID:          inventory.Warehouse.ID,
			Code:        inventory.Warehouse.Code,
			Name:        inventory.Warehouse.Name,
			Description: inventory.Warehouse.Description,
			Address:     inventory.Warehouse.Address,
			City:        inventory.Warehouse.City,
			State:       inventory.Warehouse.State,
			Country:     inventory.Warehouse.Country,
			Type:        inventory.Warehouse.Type,
			IsActive:    inventory.Warehouse.IsActive,
			IsDefault:   inventory.Warehouse.IsDefault,
		}
	}

	return response
}

func (uc *inventoryUseCase) toMovementResponse(movement *entities.InventoryMovement) *MovementResponse {
	response := &MovementResponse{
		ID:             movement.ID,
		InventoryID:    movement.InventoryID,
		Type:           string(movement.Type),
		Reason:         string(movement.Reason),
		Quantity:       movement.Quantity,
		QuantityBefore: movement.QuantityBefore,
		QuantityAfter:  movement.QuantityAfter,
		Notes:          movement.Notes,
		CreatedBy:      movement.CreatedBy,
		CreatedAt:      movement.CreatedAt,
	}

	if movement.UnitCost > 0 {
		response.UnitCost = &movement.UnitCost
	}

	if movement.TotalCost > 0 {
		response.TotalCost = &movement.TotalCost
	}

	if movement.ReferenceType != "" {
		response.ReferenceType = &movement.ReferenceType
	}

	if movement.ReferenceID != nil {
		response.ReferenceID = movement.ReferenceID
	}

	if movement.BatchNumber != "" {
		response.BatchNumber = &movement.BatchNumber
	}

	if movement.ExpiryDate != nil {
		response.ExpiryDate = movement.ExpiryDate
	}

	return response
}

// CheckAndCreateAlerts checks and creates alerts if needed
func (uc *inventoryUseCase) CheckAndCreateAlerts(ctx context.Context, inventoryID uuid.UUID) error {
	inventory, err := uc.inventoryRepo.GetByID(ctx, inventoryID)
	if err != nil {
		return err
	}

	var alerts []*entities.StockAlert

	// Check for low stock
	if inventory.IsLowStock() && !inventory.IsOutOfStock() {
		alert := &entities.StockAlert{
			ID:              uuid.New(),
			InventoryID:     inventoryID,
			Type:            entities.StockAlertTypeLowStock,
			Status:          entities.StockAlertStatusActive,
			Message:         fmt.Sprintf("Low stock alert: %s has only %d units remaining", inventory.Product.Name, inventory.QuantityAvailable),
			Severity:        "medium",
			CurrentQuantity: inventory.QuantityAvailable,
			ThresholdValue:  inventory.ReorderLevel,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		alerts = append(alerts, alert)
	}

	// Check for out of stock
	if inventory.IsOutOfStock() {
		alert := &entities.StockAlert{
			ID:              uuid.New(),
			InventoryID:     inventoryID,
			Type:            entities.StockAlertTypeOutStock,
			Status:          entities.StockAlertStatusActive,
			Message:         fmt.Sprintf("Out of stock alert: %s is out of stock", inventory.Product.Name),
			Severity:        "high",
			CurrentQuantity: inventory.QuantityAvailable,
			ThresholdValue:  0,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		alerts = append(alerts, alert)
	}

	// Check for over stock
	if inventory.IsOverStock() {
		alert := &entities.StockAlert{
			ID:              uuid.New(),
			InventoryID:     inventoryID,
			Type:            entities.StockAlertTypeOverStock,
			Status:          entities.StockAlertStatusActive,
			Message:         fmt.Sprintf("Over stock alert: %s has %d units, exceeding maximum of %d", inventory.Product.Name, inventory.QuantityOnHand, inventory.MaxStockLevel),
			Severity:        "low",
			CurrentQuantity: inventory.QuantityOnHand,
			ThresholdValue:  inventory.MaxStockLevel,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		alerts = append(alerts, alert)
	}

	// Create alerts
	for _, alert := range alerts {
		if err := uc.inventoryRepo.CreateAlert(ctx, alert); err != nil {
			return err
		}
	}

	return nil
}

// GetProductInventories gets all inventories for a specific product
func (uc *inventoryUseCase) GetProductInventories(ctx context.Context, productID uuid.UUID) ([]*InventoryResponse, error) {
	// Use repository filters to get inventories by product
	filters := repositories.InventoryFilters{
		ProductID: &productID,
	}

	inventories, err := uc.inventoryRepo.List(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get product inventories: %w", err)
	}

	// Convert to response format
	items := make([]*InventoryResponse, len(inventories))
	for i, inventory := range inventories {
		items[i] = uc.toInventoryResponse(inventory)
	}

	return items, nil
}

// GetWarehouseInventories gets all inventories for a specific warehouse
func (uc *inventoryUseCase) GetWarehouseInventories(ctx context.Context, warehouseID uuid.UUID, req GetInventoriesRequest) (*InventoriesListResponse, error) {
	// Calculate offset from page and limit
	offset := req.Page * req.Limit

	// Get inventories from repository
	inventories, err := uc.inventoryRepo.GetInventoryByWarehouse(ctx, warehouseID, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get warehouse inventories: %w", err)
	}

	// Get total count for the warehouse
	filters := repositories.InventoryFilters{
		WarehouseID: &warehouseID,
	}
	total, err := uc.inventoryRepo.Count(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to count warehouse inventories: %w", err)
	}

	// Convert to response format
	items := make([]*InventoryResponse, len(inventories))
	for i, inventory := range inventories {
		items[i] = uc.toInventoryResponse(inventory)
	}

	pagination := NewPaginationInfo(offset, req.Limit, total)

	response := &InventoriesListResponse{
		Items:      items,
		Total:      total,
		Pagination: *pagination,
	}

	return response, nil
}

// UpdateInventory updates inventory information
func (uc *inventoryUseCase) UpdateInventory(ctx context.Context, req UpdateInventoryRequest) (*InventoryResponse, error) {
	// Get current inventory
	inventory, err := uc.inventoryRepo.GetByProductAndWarehouse(ctx, req.ProductID, req.WarehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	// Update fields that are provided in request
	if req.QuantityOnHand != nil {
		inventory.QuantityOnHand = *req.QuantityOnHand
		inventory.QuantityAvailable = inventory.QuantityOnHand - inventory.QuantityReserved
	}

	if req.ReorderLevel != nil {
		inventory.ReorderLevel = *req.ReorderLevel
	}

	if req.MaxStockLevel != nil {
		inventory.MaxStockLevel = *req.MaxStockLevel
	}

	if req.MinStockLevel != nil {
		inventory.MinStockLevel = *req.MinStockLevel
	}

	if req.AverageCost != nil {
		inventory.AverageCost = *req.AverageCost
	}

	if req.LastCost != nil {
		inventory.LastCost = *req.LastCost
	}

	if req.LastCountAt != nil {
		inventory.LastCountAt = req.LastCountAt
	}

	// Update timestamp
	inventory.UpdatedAt = time.Now()

	// Save to repository
	if err := uc.inventoryRepo.Update(ctx, inventory); err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	// Check and create alerts if needed
	if err := uc.CheckAndCreateAlerts(ctx, inventory.ID); err != nil {
		// Log error but don't fail the operation
		// logger.Error("Failed to check alerts", "error", err)
	}

	return uc.toInventoryResponse(inventory), nil
}

// TransferStock transfers stock between warehouses
func (uc *inventoryUseCase) TransferStock(ctx context.Context, req TransferStockRequest) error {
	// Get source inventory
	fromInventory, err := uc.inventoryRepo.GetByProductAndWarehouse(ctx, req.ProductID, req.FromWarehouseID)
	if err != nil {
		return fmt.Errorf("failed to get source inventory: %w", err)
	}

	// Check if enough stock is available
	if fromInventory.QuantityAvailable < req.Quantity {
		return fmt.Errorf("insufficient stock in source warehouse: available %d, requested %d", 
			fromInventory.QuantityAvailable, req.Quantity)
	}

	// Get or create destination inventory
	toInventory, err := uc.inventoryRepo.GetByProductAndWarehouse(ctx, req.ProductID, req.ToWarehouseID)
	if err != nil {
		// If inventory doesn't exist for destination warehouse, create it
		toInventory = &entities.Inventory{
			ID:                uuid.New(),
			ProductID:         req.ProductID,
			WarehouseID:       req.ToWarehouseID,
			QuantityOnHand:    0,
			QuantityReserved:  0,
			QuantityAvailable: 0,
			ReorderLevel:      fromInventory.ReorderLevel,
			MaxStockLevel:     fromInventory.MaxStockLevel,
			MinStockLevel:     fromInventory.MinStockLevel,
			AverageCost:       fromInventory.AverageCost,
			LastCost:          fromInventory.LastCost,
			IsActive:          true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		
		if err := uc.inventoryRepo.Create(ctx, toInventory); err != nil {
			return fmt.Errorf("failed to create destination inventory: %w", err)
		}
	}

	// Use repository transfer method
	if err := uc.inventoryRepo.TransferStock(ctx, fromInventory.ID, toInventory.ID, req.Quantity); err != nil {
		return fmt.Errorf("failed to transfer stock: %w", err)
	}

	// Record outbound movement for source
	outboundReq := RecordMovementRequest{
		ProductID:     req.ProductID,
		WarehouseID:   req.FromWarehouseID,
		Type:          "out",
		Reason:        "transfer",
		Quantity:      req.Quantity,
		ReferenceType: &[]string{"transfer"}[0],
		Notes:         fmt.Sprintf("Transfer to warehouse %s: %s", req.ToWarehouseID.String(), req.Notes),
		CreatedBy:     req.TransferredBy,
	}

	if _, err := uc.RecordMovement(ctx, outboundReq); err != nil {
		return fmt.Errorf("failed to record outbound movement: %w", err)
	}

	// Record inbound movement for destination
	inboundReq := RecordMovementRequest{
		ProductID:     req.ProductID,
		WarehouseID:   req.ToWarehouseID,
		Type:          "in",
		Reason:        "transfer",
		Quantity:      req.Quantity,
		ReferenceType: &[]string{"transfer"}[0],
		Notes:         fmt.Sprintf("Transfer from warehouse %s: %s", req.FromWarehouseID.String(), req.Notes),
		CreatedBy:     req.TransferredBy,
	}

	if _, err := uc.RecordMovement(ctx, inboundReq); err != nil {
		return fmt.Errorf("failed to record inbound movement: %w", err)
	}

	return nil
}

// GetStockAlerts gets stock alerts
func (uc *inventoryUseCase) GetStockAlerts(ctx context.Context, req GetAlertsRequest) (*AlertsListResponse, error) {
	// Calculate offset from page and limit
	offset := req.Page * req.Limit

	// Get alerts from repository
	alerts, err := uc.inventoryRepo.GetActiveAlerts(ctx, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock alerts: %w", err)
	}

	// Filter alerts based on request criteria
	var filteredAlerts []*entities.StockAlert
	for _, alert := range alerts {
		// Apply filters
		if req.Type != "" && string(alert.Type) != req.Type {
			continue
		}
		if req.Status != "" && string(alert.Status) != req.Status {
			continue
		}
		// Add warehouse filter if specified
		// Note: This would require a join or additional lookup in a real implementation
		
		filteredAlerts = append(filteredAlerts, alert)
	}

	// Convert to response format
	items := make([]*AlertResponse, len(filteredAlerts))
	for i, alert := range filteredAlerts {
		items[i] = uc.toAlertResponse(alert)
	}

	total := int64(len(filteredAlerts))
	pagination := NewPaginationInfo(offset, req.Limit, total)

	response := &AlertsListResponse{
		Items:      items,
		Total:      total,
		Pagination: *pagination,
	}

	return response, nil
}

// Helper method to convert alert entity to response
func (uc *inventoryUseCase) toAlertResponse(alert *entities.StockAlert) *AlertResponse {
	return &AlertResponse{
		ID:              alert.ID,
		Type:            string(alert.Type),
		Message:         alert.Message,
		Severity:        alert.Severity,
		Status:          string(alert.Status),
		TriggeredAt:     alert.CreatedAt,
		ResolvedAt:      alert.ResolvedAt,
		ResolvedBy:      alert.ResolvedBy,
		Resolution:      alert.Resolution,
	}
}

// ResolveAlert resolves a stock alert
func (uc *inventoryUseCase) ResolveAlert(ctx context.Context, alertID uuid.UUID, resolution string, resolvedBy uuid.UUID) error {
	// Resolve alert in repository
	if err := uc.inventoryRepo.ResolveAlert(ctx, alertID); err != nil {
		return fmt.Errorf("failed to resolve alert: %w", err)
	}

	// Note: In a more complete implementation, you'd want to update the alert with resolution details
	// This might require extending the repository interface to support updating alert resolution info

	return nil
}
