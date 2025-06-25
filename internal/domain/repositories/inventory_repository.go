package repositories

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)



// InventoryRepository defines inventory repository interface
type InventoryRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, inventory *entities.Inventory) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Inventory, error)
	GetByProductID(ctx context.Context, productID uuid.UUID) (*entities.Inventory, error)
	GetByProductAndWarehouse(ctx context.Context, productID, warehouseID uuid.UUID) (*entities.Inventory, error)
	Update(ctx context.Context, inventory *entities.Inventory) error
	Delete(ctx context.Context, id uuid.UUID) error

	// List operations
	List(ctx context.Context, filters InventoryFilters) ([]*entities.Inventory, error)
	Count(ctx context.Context, filters InventoryFilters) (int64, error)

	// Stock operations
	UpdateStock(ctx context.Context, inventoryID uuid.UUID, quantityChange int, reason string) error
	ReserveStock(ctx context.Context, inventoryID uuid.UUID, quantity int) error
	ReleaseReservation(ctx context.Context, inventoryID uuid.UUID, quantity int) error
	GetAvailableStock(ctx context.Context, productID uuid.UUID) (int, error)

	// Movement operations
	CreateMovement(ctx context.Context, movement *entities.InventoryMovement) error
	GetMovements(ctx context.Context, inventoryID uuid.UUID, limit, offset int) ([]*entities.InventoryMovement, error)
	GetMovementsByDateRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*entities.InventoryMovement, error)

	// Alert operations
	CreateAlert(ctx context.Context, alert *entities.StockAlert) error
	GetActiveAlerts(ctx context.Context, limit, offset int) ([]*entities.StockAlert, error)
	ResolveAlert(ctx context.Context, alertID uuid.UUID) error
	GetAlertsByInventory(ctx context.Context, inventoryID uuid.UUID) ([]*entities.StockAlert, error)

	// Stock level operations
	GetLowStockItems(ctx context.Context, limit, offset int) ([]*entities.Inventory, error)
	GetOutOfStockItems(ctx context.Context, limit, offset int) ([]*entities.Inventory, error)
	CountLowStockItems(ctx context.Context) (int64, error)
	CountOutOfStockItems(ctx context.Context) (int64, error)

	// Warehouse operations
	GetInventoryByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]*entities.Inventory, error)
	TransferStock(ctx context.Context, fromInventoryID, toInventoryID uuid.UUID, quantity int) error

	// Reporting
	GetStockReport(ctx context.Context, filters StockReportFilters) (*StockReport, error)
	GetMovementReport(ctx context.Context, filters MovementReportFilters) (*MovementReport, error)
}



// StockReportFilters represents filters for stock reports
type StockReportFilters struct {
	WarehouseID *uuid.UUID
	CategoryID  *uuid.UUID
	DateFrom    *time.Time
	DateTo      *time.Time
}

// MovementReportFilters represents filters for movement reports
type MovementReportFilters struct {
	InventoryID   *uuid.UUID
	WarehouseID   *uuid.UUID
	MovementType  *entities.InventoryMovementType
	DateFrom      *time.Time
	DateTo        *time.Time
}

// StockReport represents stock report data
type StockReport struct {
	TotalItems      int64   `json:"total_items"`
	TotalValue      float64 `json:"total_value"`
	LowStockItems   int64   `json:"low_stock_items"`
	OutOfStockItems int64   `json:"out_of_stock_items"`
	OverStockItems  int64   `json:"over_stock_items"`
	Items           []StockReportItem `json:"items"`
}

// StockReportItem represents individual item in stock report
type StockReportItem struct {
	ProductID         uuid.UUID `json:"product_id"`
	ProductName       string    `json:"product_name"`
	SKU               string    `json:"sku"`
	QuantityOnHand    int       `json:"quantity_on_hand"`
	QuantityReserved  int       `json:"quantity_reserved"`
	QuantityAvailable int       `json:"quantity_available"`
	ReorderLevel      int       `json:"reorder_level"`
	MaxStockLevel     int       `json:"max_stock_level"`
	UnitCost          float64   `json:"unit_cost"`
	TotalValue        float64   `json:"total_value"`
	Status            string    `json:"status"`
}

// MovementReport represents movement report data
type MovementReport struct {
	TotalMovements int64                `json:"total_movements"`
	InboundTotal   int                  `json:"inbound_total"`
	OutboundTotal  int                  `json:"outbound_total"`
	NetChange      int                  `json:"net_change"`
	Movements      []MovementReportItem `json:"movements"`
}

// MovementReportItem represents individual movement in report
type MovementReportItem struct {
	Date         time.Time                        `json:"date"`
	ProductName  string                           `json:"product_name"`
	SKU          string                           `json:"sku"`
	Type         entities.InventoryMovementType   `json:"type"`
	Reason       entities.InventoryMovementReason `json:"reason"`
	Quantity     int                              `json:"quantity"`
	UnitCost     float64                          `json:"unit_cost"`
	TotalCost    float64                          `json:"total_cost"`
	Reference    string                           `json:"reference"`
}
