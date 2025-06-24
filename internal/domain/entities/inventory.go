package entities

import (
	"time"

	"github.com/google/uuid"
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
