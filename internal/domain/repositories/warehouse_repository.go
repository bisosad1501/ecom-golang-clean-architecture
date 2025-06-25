package repositories

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// WarehouseRepository defines warehouse repository interface
type WarehouseRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, warehouse *entities.Warehouse) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Warehouse, error)
	Update(ctx context.Context, warehouse *entities.Warehouse) error
	Delete(ctx context.Context, id uuid.UUID) error

	// List operations
	List(ctx context.Context, filters WarehouseFilters) ([]*entities.Warehouse, error)
	Count(ctx context.Context, filters WarehouseFilters) (int64, error)
	GetAll(ctx context.Context) ([]*entities.Warehouse, error)

	// Location-based operations
	GetByLocation(ctx context.Context, country, state, city string) ([]*entities.Warehouse, error)
	GetNearestWarehouses(ctx context.Context, latitude, longitude float64, radius float64, limit int) ([]*entities.Warehouse, error)
	GetWarehousesByRegion(ctx context.Context, region string) ([]*entities.Warehouse, error)

	// Status operations
	GetActiveWarehouses(ctx context.Context) ([]*entities.Warehouse, error)
	GetInactiveWarehouses(ctx context.Context) ([]*entities.Warehouse, error)
	UpdateStatus(ctx context.Context, warehouseID uuid.UUID, status entities.WarehouseStatus) error

	// Capacity operations
	GetWarehouseCapacity(ctx context.Context, warehouseID uuid.UUID) (*WarehouseCapacity, error)
	UpdateCapacity(ctx context.Context, warehouseID uuid.UUID, capacity *WarehouseCapacity) error
	GetWarehousesWithAvailableCapacity(ctx context.Context, requiredCapacity float64) ([]*entities.Warehouse, error)

	// Inventory operations
	GetWarehouseInventory(ctx context.Context, warehouseID uuid.UUID, filters WarehouseInventoryFilters) ([]*entities.Inventory, error)
	GetWarehouseInventoryCount(ctx context.Context, warehouseID uuid.UUID) (int64, error)
	GetWarehouseInventoryValue(ctx context.Context, warehouseID uuid.UUID) (float64, error)

	// Zone operations
	CreateZone(ctx context.Context, zone *entities.WarehouseZone) error
	GetZones(ctx context.Context, warehouseID uuid.UUID) ([]*entities.WarehouseZone, error)
	UpdateZone(ctx context.Context, zone *entities.WarehouseZone) error
	DeleteZone(ctx context.Context, zoneID uuid.UUID) error

	// Staff operations
	AssignStaff(ctx context.Context, warehouseID, userID uuid.UUID, role entities.WarehouseRole) error
	RemoveStaff(ctx context.Context, warehouseID, userID uuid.UUID) error
	GetWarehouseStaff(ctx context.Context, warehouseID uuid.UUID) ([]*WarehouseStaff, error)
	GetStaffWarehouses(ctx context.Context, userID uuid.UUID) ([]*entities.Warehouse, error)

	// Performance metrics
	GetWarehouseMetrics(ctx context.Context, warehouseID uuid.UUID, filters MetricsFilters) (*WarehouseMetrics, error)
	GetPerformanceReport(ctx context.Context, warehouseID uuid.UUID, filters ReportFilters) (*PerformanceReport, error)

	// Search operations
	Search(ctx context.Context, query string, filters WarehouseFilters) ([]*entities.Warehouse, error)
	GetWarehousesByManager(ctx context.Context, managerID uuid.UUID) ([]*entities.Warehouse, error)
}

















// PerformanceSummary represents performance summary
type PerformanceSummary struct {
	OverallScore        float64 `json:"overall_score"`
	PerformanceGrade    string  `json:"performance_grade"`
	ImprovementAreas    []string `json:"improvement_areas"`
	StrengthAreas       []string `json:"strength_areas"`
	ComparisonToPrevious float64 `json:"comparison_to_previous"`
}

// OrderMetrics represents order-related metrics
type OrderMetrics struct {
	TotalOrders           int64   `json:"total_orders"`
	ProcessedOrders       int64   `json:"processed_orders"`
	PendingOrders         int64   `json:"pending_orders"`
	CancelledOrders       int64   `json:"cancelled_orders"`
	FulfillmentRate       float64 `json:"fulfillment_rate"`
	AverageProcessingTime float64 `json:"average_processing_time"`
	OnTimeDeliveryRate    float64 `json:"on_time_delivery_rate"`
	OrderAccuracyRate     float64 `json:"order_accuracy_rate"`
}

// InventoryMetrics represents inventory-related metrics
type InventoryMetrics struct {
	TotalSKUs           int64   `json:"total_skus"`
	InStockSKUs         int64   `json:"in_stock_skus"`
	OutOfStockSKUs      int64   `json:"out_of_stock_skus"`
	LowStockSKUs        int64   `json:"low_stock_skus"`
	InventoryTurnover   float64 `json:"inventory_turnover"`
	StockAccuracy       float64 `json:"stock_accuracy"`
	InventoryValue      float64 `json:"inventory_value"`
	DeadStockValue      float64 `json:"dead_stock_value"`
}

// StaffMetrics represents staff-related metrics
type StaffMetrics struct {
	TotalStaff        int64   `json:"total_staff"`
	ActiveStaff       int64   `json:"active_staff"`
	ProductivityScore float64 `json:"productivity_score"`
	AttendanceRate    float64 `json:"attendance_rate"`
	TrainingHours     float64 `json:"training_hours"`
	SafetyIncidents   int64   `json:"safety_incidents"`
}

// FinancialMetrics represents financial metrics
type FinancialMetrics struct {
	TotalRevenue    float64 `json:"total_revenue"`
	OperatingCosts  float64 `json:"operating_costs"`
	LaborCosts      float64 `json:"labor_costs"`
	UtilityCosts    float64 `json:"utility_costs"`
	MaintenanceCosts float64 `json:"maintenance_costs"`
	Profitability   float64 `json:"profitability"`
	CostPerOrder    float64 `json:"cost_per_order"`
	ROI             float64 `json:"roi"`
}

// ZonePerformance represents zone performance metrics
type ZonePerformance struct {
	ZoneID          uuid.UUID `json:"zone_id"`
	ZoneName        string    `json:"zone_name"`
	UtilizationRate float64   `json:"utilization_rate"`
	PickingAccuracy float64   `json:"picking_accuracy"`
	ProcessingSpeed float64   `json:"processing_speed"`
	ErrorRate       float64   `json:"error_rate"`
	Efficiency      float64   `json:"efficiency"`
}
