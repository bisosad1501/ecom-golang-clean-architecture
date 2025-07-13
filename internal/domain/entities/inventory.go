package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// WarehouseStatus represents warehouse status
type WarehouseStatus string

const (
	WarehouseStatusActive     WarehouseStatus = "active"
	WarehouseStatusInactive   WarehouseStatus = "inactive"
	WarehouseStatusMaintenance WarehouseStatus = "maintenance"
	WarehouseStatusClosed     WarehouseStatus = "closed"
)

// WarehouseType represents warehouse type
type WarehouseType string

const (
	WarehouseTypeDistribution WarehouseType = "distribution"
	WarehouseTypeFulfillment  WarehouseType = "fulfillment"
	WarehouseTypeStorage      WarehouseType = "storage"
	WarehouseTypeRetail       WarehouseType = "retail"
	WarehouseTypeColdStorage  WarehouseType = "cold_storage"
)

// WarehouseRole represents user role in warehouse
type WarehouseRole string

const (
	WarehouseRoleManager     WarehouseRole = "manager"
	WarehouseRoleSupervisor  WarehouseRole = "supervisor"
	WarehouseRoleWorker      WarehouseRole = "worker"
	WarehouseRoleDriver      WarehouseRole = "driver"
	WarehouseRoleSecurity    WarehouseRole = "security"
	WarehouseRoleMaintenance WarehouseRole = "maintenance"
)

// InventoryMovementType represents the type of inventory movement
type InventoryMovementType string

const (
	InventoryMovementTypeIn       InventoryMovementType = "in"        // Stock increase
	InventoryMovementTypeOut      InventoryMovementType = "out"       // Stock decrease
	InventoryMovementTypeAdjust   InventoryMovementType = "adjust"    // Stock adjustment
	InventoryMovementTypeReserve  InventoryMovementType = "reserve"   // Stock reservation
	InventoryMovementTypeRelease  InventoryMovementType = "release"   // Release reservation
	InventoryMovementTypeReturn   InventoryMovementType = "return"    // Return to stock
	InventoryMovementTypeDamaged  InventoryMovementType = "damaged"   // Damaged goods
	InventoryMovementTypeExpired  InventoryMovementType = "expired"   // Expired goods
)

// InventoryMovementReason represents the reason for inventory movement
type InventoryMovementReason string

const (
	InventoryReasonPurchase     InventoryMovementReason = "purchase"      // New stock purchase
	InventoryReasonSale         InventoryMovementReason = "sale"          // Product sold
	InventoryReasonReturn       InventoryMovementReason = "return"        // Customer return
	InventoryReasonDamage       InventoryMovementReason = "damage"        // Damaged goods
	InventoryReasonExpiry       InventoryMovementReason = "expiry"        // Expired goods
	InventoryReasonAdjustment   InventoryMovementReason = "adjustment"    // Manual adjustment
	InventoryReasonReservation  InventoryMovementReason = "reservation"   // Order reservation
	InventoryReasonCancellation InventoryMovementReason = "cancellation"  // Order cancellation
	InventoryReasonTransfer     InventoryMovementReason = "transfer"      // Warehouse transfer
)

// Inventory represents product inventory information
type Inventory struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID         uuid.UUID `json:"product_id" gorm:"type:uuid;not null;uniqueIndex"`
	Product           Product   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	WarehouseID       uuid.UUID `json:"warehouse_id" gorm:"type:uuid;not null;index"`
	Warehouse         Warehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	
	// Stock levels
	QuantityOnHand    int `json:"quantity_on_hand" gorm:"default:0"`     // Physical stock
	QuantityReserved  int `json:"quantity_reserved" gorm:"default:0"`    // Reserved for orders
	QuantityAvailable int `json:"quantity_available" gorm:"default:0"`   // Available for sale
	
	// Thresholds
	ReorderLevel      int `json:"reorder_level" gorm:"default:10"`       // When to reorder
	MaxStockLevel     int `json:"max_stock_level" gorm:"default:1000"`   // Maximum stock
	MinStockLevel     int `json:"min_stock_level" gorm:"default:5"`      // Minimum stock
	
	// Cost information
	AverageCost       float64 `json:"average_cost" gorm:"default:0"`      // Average cost per unit
	LastCost          float64 `json:"last_cost" gorm:"default:0"`         // Last purchase cost
	
	// Tracking
	LastMovementAt    *time.Time `json:"last_movement_at"`
	LastCountAt       *time.Time `json:"last_count_at"`                   // Last physical count
	
	// Status
	IsActive          bool      `json:"is_active" gorm:"default:true"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	Movements []InventoryMovement `json:"movements,omitempty" gorm:"foreignKey:InventoryID"`
	Alerts    []StockAlert        `json:"alerts,omitempty" gorm:"foreignKey:InventoryID"`
}

// TableName returns the table name for Inventory entity
func (Inventory) TableName() string {
	return "inventories"
}

// IsLowStock checks if inventory is below reorder level
func (i *Inventory) IsLowStock() bool {
	return i.QuantityAvailable <= i.ReorderLevel
}

// IsOutOfStock checks if inventory is out of stock
func (i *Inventory) IsOutOfStock() bool {
	return i.QuantityAvailable <= 0
}

// IsOverStock checks if inventory exceeds maximum level
func (i *Inventory) IsOverStock() bool {
	return i.QuantityOnHand > i.MaxStockLevel
}

// CanReserve checks if quantity can be reserved
func (i *Inventory) CanReserve(quantity int) bool {
	return i.QuantityAvailable >= quantity
}

// UpdateAvailableQuantity recalculates available quantity
func (i *Inventory) UpdateAvailableQuantity() {
	i.QuantityAvailable = i.QuantityOnHand - i.QuantityReserved
	if i.QuantityAvailable < 0 {
		i.QuantityAvailable = 0
	}
}

// Validate validates inventory data
func (i *Inventory) Validate() error {
	// Validate required fields
	if i.ProductID == uuid.Nil {
		return fmt.Errorf("product ID is required")
	}
	if i.WarehouseID == uuid.Nil {
		return fmt.Errorf("warehouse ID is required")
	}

	// Validate quantities are non-negative
	if i.QuantityOnHand < 0 {
		return fmt.Errorf("quantity on hand cannot be negative")
	}
	if i.QuantityReserved < 0 {
		return fmt.Errorf("quantity reserved cannot be negative")
	}
	if i.QuantityAvailable < 0 {
		return fmt.Errorf("quantity available cannot be negative")
	}

	// Validate thresholds
	if i.ReorderLevel < 0 {
		return fmt.Errorf("reorder level cannot be negative")
	}
	if i.MaxStockLevel < 0 {
		return fmt.Errorf("max stock level cannot be negative")
	}
	if i.MinStockLevel < 0 {
		return fmt.Errorf("min stock level cannot be negative")
	}

	// Validate threshold relationships
	if i.MinStockLevel > i.ReorderLevel {
		return fmt.Errorf("min stock level (%d) cannot be greater than reorder level (%d)",
			i.MinStockLevel, i.ReorderLevel)
	}
	if i.ReorderLevel > i.MaxStockLevel {
		return fmt.Errorf("reorder level (%d) cannot be greater than max stock level (%d)",
			i.ReorderLevel, i.MaxStockLevel)
	}

	// Validate quantity consistency
	expectedAvailable := i.QuantityOnHand - i.QuantityReserved
	if expectedAvailable < 0 {
		expectedAvailable = 0
	}
	if i.QuantityAvailable != expectedAvailable {
		return fmt.Errorf("quantity available (%d) does not match calculated value (%d) from on hand (%d) - reserved (%d)",
			i.QuantityAvailable, expectedAvailable, i.QuantityOnHand, i.QuantityReserved)
	}

	// Validate reserved quantity doesn't exceed on hand
	if i.QuantityReserved > i.QuantityOnHand {
		return fmt.Errorf("quantity reserved (%d) cannot exceed quantity on hand (%d)",
			i.QuantityReserved, i.QuantityOnHand)
	}

	// Validate cost fields
	if i.AverageCost < 0 {
		return fmt.Errorf("average cost cannot be negative")
	}
	if i.LastCost < 0 {
		return fmt.Errorf("last cost cannot be negative")
	}

	return nil
}

// InventoryMovement represents inventory movement transactions
type InventoryMovement struct {
	ID          uuid.UUID               `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	InventoryID uuid.UUID               `json:"inventory_id" gorm:"type:uuid;not null;index"`
	Inventory   Inventory               `json:"inventory,omitempty" gorm:"foreignKey:InventoryID"`
	
	// Movement details
	Type        InventoryMovementType   `json:"type" gorm:"not null"`
	Reason      InventoryMovementReason `json:"reason" gorm:"not null"`
	Quantity    int                     `json:"quantity" gorm:"not null"`
	UnitCost    float64                 `json:"unit_cost" gorm:"default:0"`
	TotalCost   float64                 `json:"total_cost" gorm:"default:0"`
	
	// Before/after quantities
	QuantityBefore int `json:"quantity_before" gorm:"not null"`
	QuantityAfter  int `json:"quantity_after" gorm:"not null"`
	
	// Reference information
	ReferenceType string     `json:"reference_type"` // order, purchase_order, adjustment, etc.
	ReferenceID   *uuid.UUID `json:"reference_id" gorm:"type:uuid;index"`
	
	// Additional information
	Notes       string    `json:"notes"`
	BatchNumber string    `json:"batch_number"`
	ExpiryDate  *time.Time `json:"expiry_date"`
	
	// Tracking
	CreatedBy uuid.UUID `json:"created_by" gorm:"type:uuid"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for InventoryMovement entity
func (InventoryMovement) TableName() string {
	return "inventory_movements"
}

// IsInbound checks if movement increases stock
func (im *InventoryMovement) IsInbound() bool {
	return im.Type == InventoryMovementTypeIn || 
		   im.Type == InventoryMovementTypeReturn ||
		   im.Type == InventoryMovementTypeRelease
}

// IsOutbound checks if movement decreases stock
func (im *InventoryMovement) IsOutbound() bool {
	return im.Type == InventoryMovementTypeOut ||
		   im.Type == InventoryMovementTypeReserve ||
		   im.Type == InventoryMovementTypeDamaged ||
		   im.Type == InventoryMovementTypeExpired
}

// Validate validates inventory movement data
func (im *InventoryMovement) Validate() error {
	// Validate required fields
	if im.InventoryID == uuid.Nil {
		return fmt.Errorf("inventory ID is required")
	}
	if im.Type == "" {
		return fmt.Errorf("movement type is required")
	}
	if im.Reason == "" {
		return fmt.Errorf("movement reason is required")
	}
	if im.CreatedBy == uuid.Nil {
		return fmt.Errorf("created by is required")
	}

	// Validate quantity
	if im.Quantity == 0 {
		return fmt.Errorf("quantity cannot be zero")
	}

	// Validate before/after quantities are non-negative
	if im.QuantityBefore < 0 {
		return fmt.Errorf("quantity before cannot be negative")
	}
	if im.QuantityAfter < 0 {
		return fmt.Errorf("quantity after cannot be negative")
	}

	// Validate quantity change consistency
	expectedQuantityAfter := im.QuantityBefore
	if im.IsInbound() {
		expectedQuantityAfter += im.Quantity
	} else if im.IsOutbound() {
		expectedQuantityAfter -= im.Quantity
	}

	if im.QuantityAfter != expectedQuantityAfter {
		return fmt.Errorf("quantity after (%d) does not match expected value (%d) based on movement",
			im.QuantityAfter, expectedQuantityAfter)
	}

	// Validate cost fields
	if im.UnitCost < 0 {
		return fmt.Errorf("unit cost cannot be negative")
	}
	if im.TotalCost < 0 {
		return fmt.Errorf("total cost cannot be negative")
	}

	// Validate total cost consistency if both unit cost and quantity are provided
	if im.UnitCost > 0 && im.Quantity > 0 {
		expectedTotalCost := im.UnitCost * float64(im.Quantity)
		if im.TotalCost != expectedTotalCost {
			return fmt.Errorf("total cost (%.2f) does not match unit cost (%.2f) * quantity (%d) = %.2f",
				im.TotalCost, im.UnitCost, im.Quantity, expectedTotalCost)
		}
	}

	return nil
}

// Warehouse represents a storage location
type Warehouse struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code        string    `json:"code" gorm:"uniqueIndex;not null" validate:"required"`
	Name        string    `json:"name" gorm:"not null" validate:"required"`
	Description string    `json:"description"`
	
	// Location information
	Address     string  `json:"address"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	ZipCode     string  `json:"zip_code"`
	Country     string  `json:"country" gorm:"default:'USA'"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	
	// Warehouse details
	Type        string  `json:"type" gorm:"default:'standard'"` // standard, cold_storage, hazmat, etc.
	Capacity    int     `json:"capacity" gorm:"default:0"`      // Total capacity
	IsActive    bool    `json:"is_active" gorm:"default:true"`
	IsDefault   bool    `json:"is_default" gorm:"default:false"`
	
	// Contact information
	ManagerName  string `json:"manager_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	
	// Metadata
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	Inventories []Inventory `json:"inventories,omitempty" gorm:"foreignKey:WarehouseID"`
}

// TableName returns the table name for Warehouse entity
func (Warehouse) TableName() string {
	return "warehouses"
}

// Validate validates warehouse data
func (w *Warehouse) Validate() error {
	// Validate required fields
	if w.Code == "" {
		return fmt.Errorf("warehouse code is required")
	}
	if w.Name == "" {
		return fmt.Errorf("warehouse name is required")
	}

	// Validate capacity
	if w.Capacity < 0 {
		return fmt.Errorf("capacity cannot be negative")
	}

	// Validate coordinates if provided
	if w.Latitude != 0 || w.Longitude != 0 {
		if w.Latitude < -90 || w.Latitude > 90 {
			return fmt.Errorf("latitude must be between -90 and 90, got %.6f", w.Latitude)
		}
		if w.Longitude < -180 || w.Longitude > 180 {
			return fmt.Errorf("longitude must be between -180 and 180, got %.6f", w.Longitude)
		}
	}

	// Validate email format if provided
	if w.Email != "" {
		// Basic email validation
		if len(w.Email) < 5 || !contains(w.Email, "@") || !contains(w.Email, ".") {
			return fmt.Errorf("invalid email format: %s", w.Email)
		}
	}

	// Validate warehouse type
	validTypes := map[string]bool{
		"standard":     true,
		"cold_storage": true,
		"hazmat":       true,
		"distribution": true,
		"fulfillment":  true,
		"retail":       true,
	}
	if w.Type != "" && !validTypes[w.Type] {
		return fmt.Errorf("invalid warehouse type: %s", w.Type)
	}

	return nil
}

// Helper function for basic string contains check
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// WarehouseZone represents a zone within a warehouse
type WarehouseZone struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WarehouseID uuid.UUID `json:"warehouse_id" gorm:"type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"not null"`
	Code        string    `json:"code" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"` // storage, picking, packing, shipping, receiving
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Warehouse *Warehouse  `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Inventory []Inventory `json:"inventory,omitempty" gorm:"foreignKey:ZoneID"`
}

// TableName returns the table name for WarehouseZone entity
func (WarehouseZone) TableName() string {
	return "warehouse_zones"
}

// StockAlertType represents the type of stock alert
type StockAlertType string

const (
	StockAlertTypeLowStock  StockAlertType = "low_stock"
	StockAlertTypeOutStock  StockAlertType = "out_of_stock"
	StockAlertTypeOverStock StockAlertType = "over_stock"
	StockAlertTypeExpiring  StockAlertType = "expiring"
)

// StockAlertStatus represents the status of a stock alert
type StockAlertStatus string

const (
	StockAlertStatusActive   StockAlertStatus = "active"
	StockAlertStatusResolved StockAlertStatus = "resolved"
	StockAlertStatusIgnored  StockAlertStatus = "ignored"
)

// StockAlert represents inventory alerts
type StockAlert struct {
	ID          uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	InventoryID uuid.UUID        `json:"inventory_id" gorm:"type:uuid;not null;index"`
	Inventory   Inventory        `json:"inventory,omitempty" gorm:"foreignKey:InventoryID"`
	
	// Alert details
	Type        StockAlertType   `json:"type" gorm:"not null"`
	Status      StockAlertStatus `json:"status" gorm:"default:'active'"`
	Message     string           `json:"message" gorm:"not null"`
	Severity    string           `json:"severity" gorm:"default:'medium'"` // low, medium, high, critical
	
	// Threshold information
	CurrentQuantity int `json:"current_quantity"`
	ThresholdValue  int `json:"threshold_value"`
	
	// Resolution
	ResolvedAt *time.Time `json:"resolved_at"`
	ResolvedBy *uuid.UUID `json:"resolved_by" gorm:"type:uuid"`
	Resolution string     `json:"resolution"`
	
	// Metadata
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for StockAlert entity
func (StockAlert) TableName() string {
	return "stock_alerts"
}

// IsActive checks if alert is active
func (sa *StockAlert) IsActive() bool {
	return sa.Status == StockAlertStatusActive
}

// Resolve marks the alert as resolved
func (sa *StockAlert) Resolve(resolvedBy uuid.UUID, resolution string) {
	now := time.Now()
	sa.Status = StockAlertStatusResolved
	sa.ResolvedAt = &now
	sa.ResolvedBy = &resolvedBy
	sa.Resolution = resolution
	sa.UpdatedAt = now
}

// Validate validates stock alert data
func (sa *StockAlert) Validate() error {
	// Validate required fields
	if sa.InventoryID == uuid.Nil {
		return fmt.Errorf("inventory ID is required")
	}
	if sa.Type == "" {
		return fmt.Errorf("alert type is required")
	}
	if sa.Message == "" {
		return fmt.Errorf("alert message is required")
	}

	// Validate severity
	validSeverities := map[string]bool{
		"low":      true,
		"medium":   true,
		"high":     true,
		"critical": true,
	}
	if !validSeverities[sa.Severity] {
		return fmt.Errorf("invalid severity: %s. Must be one of: low, medium, high, critical", sa.Severity)
	}

	// Validate quantities are non-negative
	if sa.CurrentQuantity < 0 {
		return fmt.Errorf("current quantity cannot be negative")
	}
	if sa.ThresholdValue < 0 {
		return fmt.Errorf("threshold value cannot be negative")
	}

	// Validate alert type specific rules
	switch sa.Type {
	case StockAlertTypeLowStock:
		if sa.CurrentQuantity > sa.ThresholdValue {
			return fmt.Errorf("for low stock alert, current quantity (%d) should be <= threshold (%d)",
				sa.CurrentQuantity, sa.ThresholdValue)
		}
	case StockAlertTypeOutStock:
		if sa.CurrentQuantity > 0 {
			return fmt.Errorf("for out of stock alert, current quantity should be 0, got %d", sa.CurrentQuantity)
		}
	case StockAlertTypeOverStock:
		if sa.CurrentQuantity < sa.ThresholdValue {
			return fmt.Errorf("for over stock alert, current quantity (%d) should be >= threshold (%d)",
				sa.CurrentQuantity, sa.ThresholdValue)
		}
	}

	return nil
}

// Supplier represents a product supplier
type Supplier struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code        string    `json:"code" gorm:"uniqueIndex;not null" validate:"required"`
	Name        string    `json:"name" gorm:"not null" validate:"required"`
	Description string    `json:"description"`
	
	// Contact information
	ContactPerson string `json:"contact_person"`
	Email         string `json:"email" validate:"email"`
	Phone         string `json:"phone"`
	Website       string `json:"website"`
	
	// Address information
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country" gorm:"default:'USA'"`
	
	// Business information
	TaxID           string  `json:"tax_id"`
	PaymentTerms    string  `json:"payment_terms" gorm:"default:'Net 30'"`
	CreditLimit     float64 `json:"credit_limit" gorm:"default:0"`
	LeadTimeDays    int     `json:"lead_time_days" gorm:"default:7"`
	MinOrderAmount  float64 `json:"min_order_amount" gorm:"default:0"`
	
	// Status
	IsActive    bool `json:"is_active" gorm:"default:true"`
	IsPreferred bool `json:"is_preferred" gorm:"default:false"`
	
	// Ratings
	QualityRating  float64 `json:"quality_rating" gorm:"default:0"`   // 0-5 scale
	DeliveryRating float64 `json:"delivery_rating" gorm:"default:0"`  // 0-5 scale
	ServiceRating  float64 `json:"service_rating" gorm:"default:0"`   // 0-5 scale
	
	// Metadata
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	Products []Product `json:"products,omitempty" gorm:"many2many:supplier_products;"`
}

// TableName returns the table name for Supplier entity
func (Supplier) TableName() string {
	return "suppliers"
}

// GetOverallRating calculates overall supplier rating
func (s *Supplier) GetOverallRating() float64 {
	return (s.QualityRating + s.DeliveryRating + s.ServiceRating) / 3
}

// IsReliable checks if supplier is reliable based on ratings
func (s *Supplier) IsReliable() bool {
	return s.GetOverallRating() >= 4.0
}

// Validate validates supplier data
func (s *Supplier) Validate() error {
	// Validate required fields
	if s.Code == "" {
		return fmt.Errorf("supplier code is required")
	}
	if s.Name == "" {
		return fmt.Errorf("supplier name is required")
	}

	// Validate email format if provided
	if s.Email != "" {
		if len(s.Email) < 5 || !contains(s.Email, "@") || !contains(s.Email, ".") {
			return fmt.Errorf("invalid email format: %s", s.Email)
		}
	}

	// Validate financial fields
	if s.CreditLimit < 0 {
		return fmt.Errorf("credit limit cannot be negative")
	}
	if s.MinOrderAmount < 0 {
		return fmt.Errorf("minimum order amount cannot be negative")
	}
	if s.LeadTimeDays < 0 {
		return fmt.Errorf("lead time days cannot be negative")
	}

	// Validate ratings (0-5 scale)
	if s.QualityRating < 0 || s.QualityRating > 5 {
		return fmt.Errorf("quality rating must be between 0 and 5, got %.2f", s.QualityRating)
	}
	if s.DeliveryRating < 0 || s.DeliveryRating > 5 {
		return fmt.Errorf("delivery rating must be between 0 and 5, got %.2f", s.DeliveryRating)
	}
	if s.ServiceRating < 0 || s.ServiceRating > 5 {
		return fmt.Errorf("service rating must be between 0 and 5, got %.2f", s.ServiceRating)
	}

	// Validate payment terms if provided
	if s.PaymentTerms != "" {
		validTerms := map[string]bool{
			"Net 30":     true,
			"Net 60":     true,
			"Net 90":     true,
			"COD":        true,
			"Prepaid":    true,
			"2/10 Net 30": true,
		}
		if !validTerms[s.PaymentTerms] {
			return fmt.Errorf("invalid payment terms: %s", s.PaymentTerms)
		}
	}

	return nil
}

// WarehouseStats represents warehouse statistics
type WarehouseStats struct {
	TotalProducts    int64   `json:"total_products"`
	LowStockCount    int64   `json:"low_stock_count"`
	OutOfStockCount  int64   `json:"out_of_stock_count"`
	TotalValue       float64 `json:"total_value"`
}
