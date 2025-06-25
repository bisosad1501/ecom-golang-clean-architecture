package database

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type warehouseRepository struct {
	db *gorm.DB
}

// NewWarehouseRepository creates a new warehouse repository
func NewWarehouseRepository(db *gorm.DB) repositories.WarehouseRepository {
	return &warehouseRepository{db: db}
}

// Create creates a new warehouse
func (r *warehouseRepository) Create(ctx context.Context, warehouse *entities.Warehouse) error {
	return r.db.WithContext(ctx).Create(warehouse).Error
}

// GetByID gets a warehouse by ID
func (r *warehouseRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Warehouse, error) {
	var warehouse entities.Warehouse
	err := r.db.WithContext(ctx).
		Preload("Address").
		First(&warehouse, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

// GetByCode gets a warehouse by code
func (r *warehouseRepository) GetByCode(ctx context.Context, code string) (*entities.Warehouse, error) {
	var warehouse entities.Warehouse
	err := r.db.WithContext(ctx).
		Preload("Address").
		First(&warehouse, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

// Update updates a warehouse
func (r *warehouseRepository) Update(ctx context.Context, warehouse *entities.Warehouse) error {
	warehouse.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(warehouse).Error
}

// Delete deletes a warehouse
func (r *warehouseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Warehouse{}, "id = ?", id).Error
}

// List lists warehouses with filters
func (r *warehouseRepository) List(ctx context.Context, filters repositories.WarehouseFilters) ([]*entities.Warehouse, error) {
	var warehouses []*entities.Warehouse
	query := r.db.WithContext(ctx).Preload("Address")

	if filters.Name != "" {
		query = query.Where("name LIKE ?", "%"+filters.Name+"%")
	}

	if filters.Code != "" {
		query = query.Where("code LIKE ?", "%"+filters.Code+"%")
	}

	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}

	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}

	if filters.Country != "" {
		query = query.Joins("JOIN addresses ON warehouses.address_id = addresses.id").
			Where("addresses.country = ?", filters.Country)
	}

	if filters.State != "" {
		query = query.Joins("JOIN addresses ON warehouses.address_id = addresses.id").
			Where("addresses.state = ?", filters.State)
	}

	if filters.City != "" {
		query = query.Joins("JOIN addresses ON warehouses.address_id = addresses.id").
			Where("addresses.city = ?", filters.City)
	}

	// Apply sorting
	switch filters.SortBy {
	case "name":
		if filters.SortOrder == "desc" {
			query = query.Order("name DESC")
		} else {
			query = query.Order("name ASC")
		}
	case "code":
		if filters.SortOrder == "desc" {
			query = query.Order("code DESC")
		} else {
			query = query.Order("code ASC")
		}
	case "created_at":
		if filters.SortOrder == "desc" {
			query = query.Order("created_at DESC")
		} else {
			query = query.Order("created_at ASC")
		}
	default:
		query = query.Order("name ASC")
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&warehouses).Error
	return warehouses, err
}

// Count counts warehouses with filters
func (r *warehouseRepository) Count(ctx context.Context, filters repositories.WarehouseFilters) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Warehouse{})

	if filters.Name != "" {
		query = query.Where("name LIKE ?", "%"+filters.Name+"%")
	}

	if filters.Code != "" {
		query = query.Where("code LIKE ?", "%"+filters.Code+"%")
	}

	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}

	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}

	if filters.Country != "" {
		query = query.Joins("JOIN addresses ON warehouses.address_id = addresses.id").
			Where("addresses.country = ?", filters.Country)
	}

	if filters.State != "" {
		query = query.Joins("JOIN addresses ON warehouses.address_id = addresses.id").
			Where("addresses.state = ?", filters.State)
	}

	if filters.City != "" {
		query = query.Joins("JOIN addresses ON warehouses.address_id = addresses.id").
			Where("addresses.city = ?", filters.City)
	}

	err := query.Count(&count).Error
	return count, err
}

// GetActive gets all active warehouses
func (r *warehouseRepository) GetActive(ctx context.Context) ([]*entities.Warehouse, error) {
	var warehouses []*entities.Warehouse
	err := r.db.WithContext(ctx).
		Preload("Address").
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&warehouses).Error
	return warehouses, err
}

// GetByRegion gets warehouses by region
func (r *warehouseRepository) GetByRegion(ctx context.Context, country, state string) ([]*entities.Warehouse, error) {
	var warehouses []*entities.Warehouse
	query := r.db.WithContext(ctx).
		Preload("Address").
		Joins("JOIN addresses ON warehouses.address_id = addresses.id").
		Where("warehouses.is_active = ?", true)

	if country != "" {
		query = query.Where("addresses.country = ?", country)
	}

	if state != "" {
		query = query.Where("addresses.state = ?", state)
	}

	err := query.Order("warehouses.name ASC").Find(&warehouses).Error
	return warehouses, err
}

// GetNearestWarehouse gets the nearest warehouse to an address
func (r *warehouseRepository) GetNearestWarehouse(ctx context.Context, latitude, longitude float64) (*entities.Warehouse, error) {
	var warehouse entities.Warehouse
	// This is a simplified version - in production, you'd use proper geospatial queries
	err := r.db.WithContext(ctx).
		Preload("Address").
		Where("is_active = ?", true).
		First(&warehouse).Error
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

// GetInventoryByWarehouse gets inventory for a warehouse
func (r *warehouseRepository) GetInventoryByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]*entities.Inventory, error) {
	var inventories []*entities.Inventory
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("warehouse_id = ?", warehouseID).
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&inventories).Error
	return inventories, err
}

// GetLowStockByWarehouse gets low stock items for a warehouse
func (r *warehouseRepository) GetLowStockByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]*entities.Inventory, error) {
	var inventories []*entities.Inventory
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("warehouse_id = ? AND quantity_available <= reorder_level AND quantity_available > 0", warehouseID).
		Order("quantity_available ASC").
		Find(&inventories).Error
	return inventories, err
}

// GetOutOfStockByWarehouse gets out of stock items for a warehouse
func (r *warehouseRepository) GetOutOfStockByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]*entities.Inventory, error) {
	var inventories []*entities.Inventory
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("warehouse_id = ? AND quantity_available = 0", warehouseID).
		Order("updated_at DESC").
		Find(&inventories).Error
	return inventories, err
}

// UpdateStatus updates warehouse status
func (r *warehouseRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.WarehouseStatus) error {
	return r.db.WithContext(ctx).
		Model(&entities.Warehouse{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// GetWarehouseStats gets warehouse statistics
func (r *warehouseRepository) GetWarehouseStats(ctx context.Context, warehouseID uuid.UUID) (*entities.WarehouseStats, error) {
	var stats entities.WarehouseStats
	
	// Get total products
	err := r.db.WithContext(ctx).
		Model(&entities.Inventory{}).
		Where("warehouse_id = ?", warehouseID).
		Count(&stats.TotalProducts).Error
	if err != nil {
		return nil, err
	}

	// Get low stock count
	err = r.db.WithContext(ctx).
		Model(&entities.Inventory{}).
		Where("warehouse_id = ? AND quantity_available <= reorder_level AND quantity_available > 0", warehouseID).
		Count(&stats.LowStockCount).Error
	if err != nil {
		return nil, err
	}

	// Get out of stock count
	err = r.db.WithContext(ctx).
		Model(&entities.Inventory{}).
		Where("warehouse_id = ? AND quantity_available = 0", warehouseID).
		Count(&stats.OutOfStockCount).Error
	if err != nil {
		return nil, err
	}

	// Get total inventory value
	err = r.db.WithContext(ctx).
		Table("inventories").
		Select("COALESCE(SUM(inventories.quantity_on_hand * products.price), 0)").
		Joins("JOIN products ON inventories.product_id = products.id").
		Where("inventories.warehouse_id = ?", warehouseID).
		Scan(&stats.TotalValue).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// Exists checks if a warehouse exists
func (r *warehouseRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Warehouse{}).
		Where("id = ?", id).
		Count(&count).Error
	return count > 0, err
}

// ExistsByCode checks if a warehouse with the given code exists
func (r *warehouseRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Warehouse{}).
		Where("code = ?", code).
		Count(&count).Error
	return count > 0, err
}

// AssignStaff assigns staff to a warehouse
func (r *warehouseRepository) AssignStaff(ctx context.Context, warehouseID, userID uuid.UUID, role entities.WarehouseRole) error {
	// This would typically involve a warehouse_staff table
	// For now, we'll create a simple implementation
	// In a real system, you'd have a proper staff assignment table

	// Create a warehouse staff assignment record
	assignment := map[string]interface{}{
		"id":           uuid.New(),
		"warehouse_id": warehouseID,
		"user_id":      userID,
		"role":         role,
		"assigned_at":  time.Now(),
		"created_at":   time.Now(),
		"updated_at":   time.Now(),
	}

	return r.db.WithContext(ctx).Table("warehouse_staff").Create(assignment).Error
}

// CreateZone creates a warehouse zone
func (r *warehouseRepository) CreateZone(ctx context.Context, zone *entities.WarehouseZone) error {
	zone.ID = uuid.New()
	zone.CreatedAt = time.Now()
	zone.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(zone).Error
}

// DeleteZone deletes a warehouse zone
func (r *warehouseRepository) DeleteZone(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.WarehouseZone{}, "id = ?", id).Error
}

// GetActiveWarehouses gets all active warehouses
func (r *warehouseRepository) GetActiveWarehouses(ctx context.Context) ([]*entities.Warehouse, error) {
	var warehouses []*entities.Warehouse
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&warehouses).Error
	return warehouses, err
}

// GetAll gets all warehouses
func (r *warehouseRepository) GetAll(ctx context.Context) ([]*entities.Warehouse, error) {
	var warehouses []*entities.Warehouse
	err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&warehouses).Error
	return warehouses, err
}

// GetByLocation gets warehouses by location
func (r *warehouseRepository) GetByLocation(ctx context.Context, country, state, city string) ([]*entities.Warehouse, error) {
	var warehouses []*entities.Warehouse
	query := r.db.WithContext(ctx)

	if country != "" {
		query = query.Where("country = ?", country)
	}
	if state != "" {
		query = query.Where("state = ?", state)
	}
	if city != "" {
		query = query.Where("city = ?", city)
	}

	err := query.Order("name ASC").Find(&warehouses).Error
	return warehouses, err
}

// GetInactiveWarehouses gets all inactive warehouses
func (r *warehouseRepository) GetInactiveWarehouses(ctx context.Context) ([]*entities.Warehouse, error) {
	var warehouses []*entities.Warehouse
	err := r.db.WithContext(ctx).
		Where("is_active = ?", false).
		Order("name ASC").
		Find(&warehouses).Error
	return warehouses, err
}

// GetNearestWarehouses gets warehouses nearest to a location
func (r *warehouseRepository) GetNearestWarehouses(ctx context.Context, latitude, longitude float64, radiusKm float64, limit int) ([]*entities.Warehouse, error) {
	var warehouses []*entities.Warehouse
	// This is a simplified implementation - in production you'd use PostGIS or similar
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name ASC").
		Limit(limit).
		Find(&warehouses).Error
	return warehouses, err
}

// GetPerformanceReport gets warehouse performance report
func (r *warehouseRepository) GetPerformanceReport(ctx context.Context, warehouseID uuid.UUID, filters repositories.ReportFilters) (*repositories.PerformanceReport, error) {
	var report repositories.PerformanceReport

	query := r.db.WithContext(ctx).Model(&entities.Order{}).Where("warehouse_id = ?", warehouseID)

	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	// Get total orders processed
	err := query.Count(&report.TotalOrders).Error
	if err != nil {
		return nil, err
	}

	// Get total revenue
	err = query.Select("COALESCE(SUM(total_amount), 0)").
		Where("status = ?", entities.OrderStatusDelivered).
		Scan(&report.TotalRevenue).Error
	if err != nil {
		return nil, err
	}

	return &report, nil
}

// GetStaffWarehouses gets warehouses assigned to a staff member
func (r *warehouseRepository) GetStaffWarehouses(ctx context.Context, userID uuid.UUID) ([]*entities.Warehouse, error) {
	var warehouses []*entities.Warehouse
	// This is a simplified implementation - in production you'd have a staff-warehouse relationship table
	err := r.db.WithContext(ctx).
		Where("manager_id = ? OR created_by = ?", userID, userID).
		Order("name ASC").
		Find(&warehouses).Error
	return warehouses, err
}

// GetWarehouseCapacity gets warehouse capacity information
func (r *warehouseRepository) GetWarehouseCapacity(ctx context.Context, warehouseID uuid.UUID) (*repositories.WarehouseCapacity, error) {
	var capacity repositories.WarehouseCapacity

	// Get warehouse info
	var warehouse entities.Warehouse
	err := r.db.WithContext(ctx).First(&warehouse, "id = ?", warehouseID).Error
	if err != nil {
		return nil, err
	}

	capacity.TotalCapacity = float64(warehouse.Capacity)
	capacity.CapacityUnit = "cubic_meters"

	// Calculate used capacity (simplified)
	var usedCapacity float64
	err = r.db.WithContext(ctx).
		Model(&entities.Inventory{}).
		Select("COALESCE(SUM(quantity_available * 0.1), 0)"). // Assume 0.1 cubic meter per item
		Where("warehouse_id = ?", warehouseID).
		Scan(&usedCapacity).Error
	if err != nil {
		return nil, err
	}

	capacity.UsedCapacity = usedCapacity
	capacity.AvailableCapacity = capacity.TotalCapacity - capacity.UsedCapacity

	return &capacity, nil
}

// GetWarehouseInventory gets inventory items in a warehouse with filters
func (r *warehouseRepository) GetWarehouseInventory(ctx context.Context, warehouseID uuid.UUID, filters repositories.WarehouseInventoryFilters) ([]*entities.Inventory, error) {
	var inventory []*entities.Inventory
	query := r.db.WithContext(ctx).Where("warehouse_id = ?", warehouseID)

	if filters.ProductID != nil {
		query = query.Where("product_id = ?", *filters.ProductID)
	}

	if filters.LowStock {
		query = query.Where("quantity_available <= reorder_level")
	}

	err := query.Preload("Product").
		Order("created_at DESC").
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&inventory).Error
	return inventory, err
}

// GetWarehouseInventoryCount gets inventory count in a warehouse
func (r *warehouseRepository) GetWarehouseInventoryCount(ctx context.Context, warehouseID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Inventory{}).
		Where("warehouse_id = ?", warehouseID).
		Count(&count).Error
	return count, err
}

// GetWarehouseInventoryValue gets total inventory value in a warehouse
func (r *warehouseRepository) GetWarehouseInventoryValue(ctx context.Context, warehouseID uuid.UUID) (float64, error) {
	var value float64
	err := r.db.WithContext(ctx).
		Model(&entities.Inventory{}).
		Select("COALESCE(SUM(quantity_available * unit_cost), 0)").
		Where("warehouse_id = ?", warehouseID).
		Scan(&value).Error
	return value, err
}

// GetWarehouseMetrics gets warehouse metrics with filters
func (r *warehouseRepository) GetWarehouseMetrics(ctx context.Context, warehouseID uuid.UUID, filters repositories.MetricsFilters) (*repositories.WarehouseMetrics, error) {
	var metrics repositories.WarehouseMetrics

	// Get total products
	err := r.db.WithContext(ctx).
		Model(&entities.Inventory{}).
		Where("warehouse_id = ?", warehouseID).
		Count(&metrics.TotalProducts).Error
	if err != nil {
		return nil, err
	}

	// Get low stock count
	err = r.db.WithContext(ctx).
		Model(&entities.Inventory{}).
		Where("warehouse_id = ? AND quantity_available <= reorder_level", warehouseID).
		Count(&metrics.LowStockCount).Error
	if err != nil {
		return nil, err
	}

	// Get out of stock count
	err = r.db.WithContext(ctx).
		Model(&entities.Inventory{}).
		Where("warehouse_id = ? AND quantity_available = 0", warehouseID).
		Count(&metrics.OutOfStockCount).Error
	if err != nil {
		return nil, err
	}

	return &metrics, nil
}

// GetWarehouseStaff gets staff assigned to a warehouse
func (r *warehouseRepository) GetWarehouseStaff(ctx context.Context, warehouseID uuid.UUID) ([]*repositories.WarehouseStaff, error) {
	var staff []*repositories.WarehouseStaff
	// This is a simplified implementation - in production you'd have a proper staff-warehouse relationship table
	err := r.db.WithContext(ctx).
		Table("users").
		Select("users.id, users.name, users.email, 'manager' as role, users.created_at as assigned_at, users.is_active").
		Where("users.id IN (SELECT manager_id FROM warehouses WHERE id = ?)", warehouseID).
		Scan(&staff).Error
	return staff, err
}

// GetWarehousesByManager gets warehouses managed by a specific manager
func (r *warehouseRepository) GetWarehousesByManager(ctx context.Context, managerID uuid.UUID) ([]*entities.Warehouse, error) {
	var warehouses []*entities.Warehouse
	err := r.db.WithContext(ctx).
		Where("manager_id = ?", managerID).
		Order("name ASC").
		Find(&warehouses).Error
	return warehouses, err
}

// GetWarehousesByRegion gets warehouses in a specific region
func (r *warehouseRepository) GetWarehousesByRegion(ctx context.Context, region string) ([]*entities.Warehouse, error) {
	var warehouses []*entities.Warehouse
	err := r.db.WithContext(ctx).
		Where("region = ? OR state = ? OR country = ?", region, region, region).
		Order("name ASC").
		Find(&warehouses).Error
	return warehouses, err
}

// GetWarehousesWithAvailableCapacity gets warehouses with available capacity
func (r *warehouseRepository) GetWarehousesWithAvailableCapacity(ctx context.Context, minCapacity float64) ([]*entities.Warehouse, error) {
	var warehouses []*entities.Warehouse
	err := r.db.WithContext(ctx).
		Where("capacity >= ? AND is_active = ?", minCapacity, true).
		Order("capacity DESC").
		Find(&warehouses).Error
	return warehouses, err
}

// GetZones gets all zones in a warehouse
func (r *warehouseRepository) GetZones(ctx context.Context, warehouseID uuid.UUID) ([]*entities.WarehouseZone, error) {
	var zones []*entities.WarehouseZone
	err := r.db.WithContext(ctx).
		Where("warehouse_id = ?", warehouseID).
		Order("name ASC").
		Find(&zones).Error
	return zones, err
}

// RemoveStaff removes staff from a warehouse
func (r *warehouseRepository) RemoveStaff(ctx context.Context, warehouseID, userID uuid.UUID) error {
	// This is a simplified implementation - in production you'd have a proper staff-warehouse relationship table
	return r.db.WithContext(ctx).
		Model(&entities.Warehouse{}).
		Where("id = ? AND manager_id = ?", warehouseID, userID).
		Update("manager_id", nil).Error
}

// Search searches warehouses by query and filters
func (r *warehouseRepository) Search(ctx context.Context, query string, filters repositories.WarehouseFilters) ([]*entities.Warehouse, error) {
	var warehouses []*entities.Warehouse
	dbQuery := r.db.WithContext(ctx)

	if query != "" {
		dbQuery = dbQuery.Where("name LIKE ? OR address LIKE ? OR city LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	if filters.IsActive != nil {
		dbQuery = dbQuery.Where("is_active = ?", *filters.IsActive)
	}

	if filters.Country != "" {
		dbQuery = dbQuery.Where("country = ?", filters.Country)
	}

	if filters.State != "" {
		dbQuery = dbQuery.Where("state = ?", filters.State)
	}

	if filters.City != "" {
		dbQuery = dbQuery.Where("city = ?", filters.City)
	}

	err := dbQuery.Order("name ASC").
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&warehouses).Error
	return warehouses, err
}

// UpdateCapacity updates warehouse capacity
func (r *warehouseRepository) UpdateCapacity(ctx context.Context, warehouseID uuid.UUID, capacity *repositories.WarehouseCapacity) error {
	return r.db.WithContext(ctx).
		Model(&entities.Warehouse{}).
		Where("id = ?", warehouseID).
		Update("capacity", capacity.TotalCapacity).Error
}

// UpdateZone updates a warehouse zone
func (r *warehouseRepository) UpdateZone(ctx context.Context, zone *entities.WarehouseZone) error {
	zone.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(zone).Error
}


