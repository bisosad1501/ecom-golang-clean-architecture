package entities

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ShippingMethodType represents the type of shipping method
type ShippingMethodType string

const (
	ShippingMethodStandard  ShippingMethodType = "standard"
	ShippingMethodExpress   ShippingMethodType = "express"
	ShippingMethodOvernight ShippingMethodType = "overnight"
	ShippingMethodSameDay   ShippingMethodType = "same_day"
	ShippingMethodPickup    ShippingMethodType = "pickup"
	ShippingMethodFree      ShippingMethodType = "free"
)

// ShippingMethod represents available shipping methods
type ShippingMethod struct {
	ID                uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name              string             `json:"name" gorm:"not null" validate:"required"`
	Description       string             `json:"description"`
	Type              ShippingMethodType `json:"type" gorm:"not null"`
	Carrier           string             `json:"carrier" gorm:"not null"` // UPS, FedEx, USPS, DHL, etc.
	
	// Pricing
	BaseCost          float64 `json:"base_cost" gorm:"default:0"`
	CostPerKg         float64 `json:"cost_per_kg" gorm:"default:0"`
	CostPerKm         float64 `json:"cost_per_km" gorm:"default:0"`
	FreeShippingMin   float64 `json:"free_shipping_min" gorm:"default:0"` // Minimum order for free shipping
	
	// Delivery time
	MinDeliveryDays   int `json:"min_delivery_days" gorm:"default:1"`
	MaxDeliveryDays   int `json:"max_delivery_days" gorm:"default:7"`
	
	// Restrictions
	MaxWeight         float64 `json:"max_weight" gorm:"default:0"`        // 0 = no limit
	MaxDimensions     string  `json:"max_dimensions"`                     // LxWxH format
	RestrictedItems   string  `json:"restricted_items"`                   // JSON array of restricted item types
	
	// Coverage
	DomesticOnly      bool   `json:"domestic_only" gorm:"default:true"`
	SupportedCountries string `json:"supported_countries"`               // JSON array of country codes
	SupportedZones    string `json:"supported_zones"`                   // JSON array of shipping zones
	
	// Status
	IsActive          bool      `json:"is_active" gorm:"default:true"`
	IsDefault         bool      `json:"is_default" gorm:"default:false"`
	SortOrder         int       `json:"sort_order" gorm:"default:0"`
	
	// Metadata
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for ShippingMethod entity
func (ShippingMethod) TableName() string {
	return "shipping_methods"
}

// Validate validates shipping method data
func (sm *ShippingMethod) Validate() error {
	if sm.Name == "" {
		return fmt.Errorf("shipping method name is required")
	}
	if sm.Carrier == "" {
		return fmt.Errorf("carrier is required")
	}
	if sm.BaseCost < 0 {
		return fmt.Errorf("base cost cannot be negative")
	}
	if sm.CostPerKg < 0 {
		return fmt.Errorf("cost per kg cannot be negative")
	}
	if sm.CostPerKm < 0 {
		return fmt.Errorf("cost per km cannot be negative")
	}
	if sm.FreeShippingMin < 0 {
		return fmt.Errorf("free shipping minimum cannot be negative")
	}
	if sm.MinDeliveryDays < 0 {
		return fmt.Errorf("minimum delivery days cannot be negative")
	}
	if sm.MaxDeliveryDays < sm.MinDeliveryDays {
		return fmt.Errorf("maximum delivery days cannot be less than minimum delivery days")
	}
	if sm.MaxWeight < 0 {
		return fmt.Errorf("maximum weight cannot be negative")
	}

	// Validate dimensions format if present
	if sm.MaxDimensions != "" {
		if err := validateDimensionsFormat(sm.MaxDimensions); err != nil {
			return fmt.Errorf("invalid max dimensions format: %w", err)
		}
	}

	return nil
}

// CalculateCost calculates shipping cost based on weight, distance, and order value
func (sm *ShippingMethod) CalculateCost(weight float64, distance float64, orderValue float64) float64 {
	// Validate inputs
	if weight < 0 {
		weight = 0
	}
	if distance < 0 {
		distance = 100 // Default distance for unknown locations
	}
	if orderValue < 0 {
		orderValue = 0
	}

	// Check if method is available for this weight
	if !sm.IsAvailableForWeight(weight) {
		return -1 // Return -1 to indicate method not available
	}

	// Check if eligible for free shipping (exact match or greater)
	if sm.FreeShippingMin > 0 && orderValue >= sm.FreeShippingMin {
		return 0
	}

	cost := sm.BaseCost

	// Add weight-based cost
	if sm.CostPerKg > 0 && weight > 0 {
		cost += weight * sm.CostPerKg
	}

	// Add distance-based cost
	if sm.CostPerKm > 0 && distance > 0 {
		cost += distance * sm.CostPerKm
	}

	// Ensure cost is not negative
	if cost < 0 {
		cost = 0
	}

	// Round to 2 decimal places
	return float64(int(cost*100+0.5)) / 100
}

// IsAvailableForWeight checks if method supports the given weight
func (sm *ShippingMethod) IsAvailableForWeight(weight float64) bool {
	return sm.MaxWeight == 0 || weight <= sm.MaxWeight
}

// EstimateDeliveryTime estimates delivery time based on distance and method type
func (sm *ShippingMethod) EstimateDeliveryTime(distance float64) (minDays, maxDays int) {
	// Base delivery time from method configuration
	minDays = sm.MinDeliveryDays
	maxDays = sm.MaxDeliveryDays

	// Adjust based on distance (for domestic shipping)
	if distance > 0 {
		// Add extra days for long distances
		if distance > 1000 { // > 1000km
			extraDays := int(distance/1000) // 1 extra day per 1000km
			if extraDays > 3 {
				extraDays = 3 // Cap at 3 extra days
			}
			minDays += extraDays
			maxDays += extraDays
		}
	}

	// Adjust based on method type
	switch sm.Type {
	case ShippingMethodSameDay:
		minDays = 0
		maxDays = 0
	case ShippingMethodOvernight:
		minDays = 1
		maxDays = 1
	case ShippingMethodExpress:
		if minDays > 3 {
			minDays = 3
		}
		if maxDays > 5 {
			maxDays = 5
		}
	case ShippingMethodStandard:
		// Use configured values
	case ShippingMethodPickup:
		minDays = 0
		maxDays = 1
	}

	// Ensure minimum values
	if minDays < 0 {
		minDays = 0
	}
	if maxDays < minDays {
		maxDays = minDays
	}

	return minDays, maxDays
}

// GetEstimatedDeliveryDate calculates estimated delivery date
func (sm *ShippingMethod) GetEstimatedDeliveryDate(distance float64) (minDate, maxDate time.Time) {
	minDays, maxDays := sm.EstimateDeliveryTime(distance)

	now := time.Now()

	// Skip weekends for business days calculation
	minDate = addBusinessDays(now, minDays)
	maxDate = addBusinessDays(now, maxDays)

	return minDate, maxDate
}

// addBusinessDays adds business days to a date (skipping weekends)
func addBusinessDays(date time.Time, days int) time.Time {
	for days > 0 {
		date = date.AddDate(0, 0, 1)
		// Skip weekends
		if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
			days--
		}
	}
	return date
}

// ShippingZone represents shipping zones for rate calculation
type ShippingZone struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null" validate:"required"`
	Description string    `json:"description"`
	
	// Geographic coverage
	Countries   string `json:"countries"`   // JSON array of country codes
	States      string `json:"states"`      // JSON array of state codes
	ZipCodes    string `json:"zip_codes"`   // JSON array of zip code patterns
	
	// Zone settings
	IsDefault   bool      `json:"is_default" gorm:"default:false"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	
	// Metadata
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	Rates []ShippingRate `json:"rates,omitempty" gorm:"foreignKey:ZoneID"`
}

// TableName returns the table name for ShippingZone entity
func (ShippingZone) TableName() string {
	return "shipping_zones"
}

// ShippingRate represents shipping rates for different zones and methods
type ShippingRate struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ZoneID           uuid.UUID      `json:"zone_id" gorm:"type:uuid;not null;index"`
	Zone             ShippingZone   `json:"zone,omitempty" gorm:"foreignKey:ZoneID"`
	ShippingMethodID uuid.UUID      `json:"shipping_method_id" gorm:"type:uuid;not null;index"`
	ShippingMethod   ShippingMethod `json:"shipping_method,omitempty" gorm:"foreignKey:ShippingMethodID"`
	
	// Rate structure
	MinWeight        float64 `json:"min_weight" gorm:"default:0"`
	MaxWeight        float64 `json:"max_weight" gorm:"default:0"`        // 0 = no limit
	MinOrderValue    float64 `json:"min_order_value" gorm:"default:0"`
	MaxOrderValue    float64 `json:"max_order_value" gorm:"default:0"`   // 0 = no limit
	
	// Pricing
	BaseCost         float64 `json:"base_cost" gorm:"default:0"`
	CostPerKg        float64 `json:"cost_per_kg" gorm:"default:0"`
	FreeShippingMin  float64 `json:"free_shipping_min" gorm:"default:0"`
	
	// Status
	IsActive         bool      `json:"is_active" gorm:"default:true"`
	
	// Metadata
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for ShippingRate entity
func (ShippingRate) TableName() string {
	return "shipping_rates"
}

// CalculateCost calculates shipping cost for this rate
func (sr *ShippingRate) CalculateCost(weight float64, orderValue float64) float64 {
	// Validate inputs
	if weight < 0 {
		weight = 0
	}
	if orderValue < 0 {
		orderValue = 0
	}

	// Check if rate is applicable
	if !sr.IsApplicable(weight, orderValue) {
		return -1 // Return -1 to indicate rate not applicable
	}

	// Check if eligible for free shipping
	if sr.FreeShippingMin > 0 && orderValue >= sr.FreeShippingMin {
		return 0
	}

	cost := sr.BaseCost

	// Add weight-based cost
	if sr.CostPerKg > 0 && weight > 0 {
		cost += weight * sr.CostPerKg
	}

	// Ensure cost is not negative
	if cost < 0 {
		cost = 0
	}

	// Round to 2 decimal places
	return float64(int(cost*100+0.5)) / 100
}

// IsApplicable checks if rate applies to given weight and order value
func (sr *ShippingRate) IsApplicable(weight float64, orderValue float64) bool {
	// Check weight range
	if sr.MinWeight > 0 && weight < sr.MinWeight {
		return false
	}
	if sr.MaxWeight > 0 && weight > sr.MaxWeight {
		return false
	}
	
	// Check order value range
	if sr.MinOrderValue > 0 && orderValue < sr.MinOrderValue {
		return false
	}
	if sr.MaxOrderValue > 0 && orderValue > sr.MaxOrderValue {
		return false
	}
	
	return true
}

// ShipmentStatus represents the status of a shipment
type ShipmentStatus string

const (
	ShipmentStatusPending    ShipmentStatus = "pending"
	ShipmentStatusProcessing ShipmentStatus = "processing"
	ShipmentStatusShipped    ShipmentStatus = "shipped"
	ShipmentStatusInTransit  ShipmentStatus = "in_transit"
	ShipmentStatusOutForDelivery ShipmentStatus = "out_for_delivery"
	ShipmentStatusDelivered  ShipmentStatus = "delivered"
	ShipmentStatusFailed     ShipmentStatus = "failed"
	ShipmentStatusReturned   ShipmentStatus = "returned"
	ShipmentStatusCancelled  ShipmentStatus = "cancelled"
)

// Shipment represents a shipment for an order
type Shipment struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID          uuid.UUID      `json:"order_id" gorm:"type:uuid;not null;index"`
	Order            Order          `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	ShippingMethodID uuid.UUID      `json:"shipping_method_id" gorm:"type:uuid;not null"`
	ShippingMethod   ShippingMethod `json:"shipping_method,omitempty" gorm:"foreignKey:ShippingMethodID"`
	
	// Tracking information
	TrackingNumber   string         `json:"tracking_number" gorm:"uniqueIndex"`
	Carrier          string         `json:"carrier" gorm:"not null"`
	Status           ShipmentStatus `json:"status" gorm:"default:'pending'"`
	
	// Shipping details
	Weight           float64   `json:"weight" gorm:"default:0"`
	Dimensions       string    `json:"dimensions"`                    // LxWxH format
	PackageCount     int       `json:"package_count" gorm:"default:1"`
	InsuranceValue   float64   `json:"insurance_value" gorm:"default:0"`
	
	// Addresses (denormalized for tracking)
	FromAddress      string    `json:"from_address" gorm:"type:text"`
	ToAddress        string    `json:"to_address" gorm:"type:text"`
	
	// Costs
	ShippingCost     float64   `json:"shipping_cost" gorm:"default:0"`
	InsuranceCost    float64   `json:"insurance_cost" gorm:"default:0"`
	TotalCost        float64   `json:"total_cost" gorm:"default:0"`
	
	// Dates
	ShippedAt        *time.Time `json:"shipped_at"`
	EstimatedDelivery *time.Time `json:"estimated_delivery"`
	ActualDelivery   *time.Time `json:"actual_delivery"`
	
	// Additional information
	Notes            string    `json:"notes"`
	SpecialInstructions string `json:"special_instructions"`
	SignatureRequired bool     `json:"signature_required" gorm:"default:false"`
	
	// Metadata
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	TrackingEvents []ShipmentTracking `json:"tracking_events,omitempty" gorm:"foreignKey:ShipmentID"`
}

// TableName returns the table name for Shipment entity
func (Shipment) TableName() string {
	return "shipments"
}

// Validate validates shipment data
func (s *Shipment) Validate() error {
	if s.OrderID == uuid.Nil {
		return fmt.Errorf("order ID is required")
	}
	if s.ShippingMethodID == uuid.Nil {
		return fmt.Errorf("shipping method ID is required")
	}
	if s.Carrier == "" {
		return fmt.Errorf("carrier is required")
	}
	if s.Weight < 0 {
		return fmt.Errorf("weight cannot be negative")
	}
	if s.PackageCount <= 0 {
		return fmt.Errorf("package count must be greater than 0")
	}
	if s.InsuranceValue < 0 {
		return fmt.Errorf("insurance value cannot be negative")
	}
	if s.ShippingCost < 0 {
		return fmt.Errorf("shipping cost cannot be negative")
	}
	if s.InsuranceCost < 0 {
		return fmt.Errorf("insurance cost cannot be negative")
	}
	if s.TotalCost < 0 {
		return fmt.Errorf("total cost cannot be negative")
	}

	// Validate dimensions format if present
	if s.Dimensions != "" {
		if err := validateDimensionsFormat(s.Dimensions); err != nil {
			return fmt.Errorf("invalid dimensions format: %w", err)
		}
	}

	return nil
}

// IsDelivered checks if shipment is delivered
func (s *Shipment) IsDelivered() bool {
	return s.Status == ShipmentStatusDelivered
}

// IsInTransit checks if shipment is in transit
func (s *Shipment) IsInTransit() bool {
	return s.Status == ShipmentStatusInTransit || 
		   s.Status == ShipmentStatusOutForDelivery
}

// GetDeliveryDays calculates delivery days
func (s *Shipment) GetDeliveryDays() int {
	if s.ShippedAt == nil || s.ActualDelivery == nil {
		return 0
	}
	return int(s.ActualDelivery.Sub(*s.ShippedAt).Hours() / 24)
}

// CanTransitionTo checks if shipment can transition to the given status
func (s *Shipment) CanTransitionTo(newStatus ShipmentStatus) bool {
	switch s.Status {
	case ShipmentStatusPending:
		return newStatus == ShipmentStatusProcessing || newStatus == ShipmentStatusCancelled
	case ShipmentStatusProcessing:
		return newStatus == ShipmentStatusShipped || newStatus == ShipmentStatusCancelled
	case ShipmentStatusShipped:
		return newStatus == ShipmentStatusInTransit || newStatus == ShipmentStatusDelivered ||
			   newStatus == ShipmentStatusFailed || newStatus == ShipmentStatusReturned
	case ShipmentStatusInTransit:
		return newStatus == ShipmentStatusOutForDelivery || newStatus == ShipmentStatusDelivered ||
			   newStatus == ShipmentStatusFailed || newStatus == ShipmentStatusReturned
	case ShipmentStatusOutForDelivery:
		return newStatus == ShipmentStatusDelivered || newStatus == ShipmentStatusFailed ||
			   newStatus == ShipmentStatusReturned
	case ShipmentStatusDelivered, ShipmentStatusFailed, ShipmentStatusReturned, ShipmentStatusCancelled:
		return false // Terminal states
	default:
		return false
	}
}

// TransitionTo transitions shipment to new status with validation
func (s *Shipment) TransitionTo(newStatus ShipmentStatus) error {
	if !s.CanTransitionTo(newStatus) {
		return fmt.Errorf("cannot transition shipment from %s to %s", s.Status, newStatus)
	}

	s.Status = newStatus
	s.UpdatedAt = time.Now()

	// Update related fields based on status
	switch newStatus {
	case ShipmentStatusShipped:
		if s.ShippedAt == nil {
			now := time.Now()
			s.ShippedAt = &now
		}
	case ShipmentStatusDelivered:
		if s.ActualDelivery == nil {
			now := time.Now()
			s.ActualDelivery = &now
		}
	}

	return nil
}

// IsOverdue checks if shipment is overdue based on estimated delivery
func (s *Shipment) IsOverdue() bool {
	if s.EstimatedDelivery == nil {
		return false
	}

	// Only consider overdue if not yet delivered
	if s.Status == ShipmentStatusDelivered {
		return false
	}

	// Check if current time is past estimated delivery + grace period
	gracePeriod := 24 * time.Hour // 1 day grace period
	return time.Now().After(s.EstimatedDelivery.Add(gracePeriod))
}

// UpdateEstimatedDelivery updates estimated delivery based on shipping method
func (s *Shipment) UpdateEstimatedDelivery(distance float64) {
	if s.ShippingMethod.ID != uuid.Nil {
		_, maxDate := s.ShippingMethod.GetEstimatedDeliveryDate(distance)
		s.EstimatedDelivery = &maxDate
	}
}

// ShipmentTracking represents tracking events for a shipment
type ShipmentTracking struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ShipmentID  uuid.UUID      `json:"shipment_id" gorm:"type:uuid;not null;index"`
	Shipment    Shipment       `json:"shipment,omitempty" gorm:"foreignKey:ShipmentID"`
	
	// Event details
	Status      ShipmentStatus `json:"status" gorm:"not null"`
	Location    string         `json:"location"`
	Description string         `json:"description" gorm:"not null"`
	EventTime   time.Time      `json:"event_time" gorm:"not null"`
	
	// Additional information
	Latitude    float64        `json:"latitude"`
	Longitude   float64        `json:"longitude"`
	Notes       string         `json:"notes"`
	
	// Metadata
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for ShipmentTracking entity
func (ShipmentTracking) TableName() string {
	return "shipment_tracking"
}

// ReturnStatus represents the status of a return
type ReturnStatus string

const (
	ReturnStatusRequested ReturnStatus = "requested"
	ReturnStatusApproved  ReturnStatus = "approved"
	ReturnStatusRejected  ReturnStatus = "rejected"
	ReturnStatusShipped   ReturnStatus = "shipped"
	ReturnStatusReceived  ReturnStatus = "received"
	ReturnStatusProcessed ReturnStatus = "processed"
	ReturnStatusCompleted ReturnStatus = "completed"
	ReturnStatusCancelled ReturnStatus = "cancelled"
)

// ReturnReason represents the reason for return
type ReturnReason string

const (
	ReturnReasonDefective     ReturnReason = "defective"
	ReturnReasonWrongItem     ReturnReason = "wrong_item"
	ReturnReasonNotAsDescribed ReturnReason = "not_as_described"
	ReturnReasonDamaged       ReturnReason = "damaged"
	ReturnReasonChangedMind   ReturnReason = "changed_mind"
	ReturnReasonSizeIssue     ReturnReason = "size_issue"
	ReturnReasonOther         ReturnReason = "other"
)

// Return represents a product return request
type Return struct {
	ID              uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID         uuid.UUID    `json:"order_id" gorm:"type:uuid;not null;index"`
	Order           Order        `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	UserID          uuid.UUID    `json:"user_id" gorm:"type:uuid;not null;index"`
	User            User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	
	// Return details
	ReturnNumber    string       `json:"return_number" gorm:"uniqueIndex;not null"`
	Reason          ReturnReason `json:"reason" gorm:"not null"`
	Status          ReturnStatus `json:"status" gorm:"default:'requested'"`
	Description     string       `json:"description" gorm:"type:text"`
	
	// Items being returned
	Items           []ReturnItem `json:"items,omitempty" gorm:"foreignKey:ReturnID"`
	
	// Financial information
	RefundAmount    float64      `json:"refund_amount" gorm:"default:0"`
	RestockingFee   float64      `json:"restocking_fee" gorm:"default:0"`
	ShippingRefund  float64      `json:"shipping_refund" gorm:"default:0"`
	
	// Tracking
	ReturnShipmentID *uuid.UUID  `json:"return_shipment_id" gorm:"type:uuid"`
	TrackingNumber   string      `json:"tracking_number"`
	
	// Dates
	RequestedAt     time.Time    `json:"requested_at" gorm:"autoCreateTime"`
	ApprovedAt      *time.Time   `json:"approved_at"`
	ReceivedAt      *time.Time   `json:"received_at"`
	ProcessedAt     *time.Time   `json:"processed_at"`
	CompletedAt     *time.Time   `json:"completed_at"`
	
	// Processing information
	ProcessedBy     *uuid.UUID   `json:"processed_by" gorm:"type:uuid"`
	ProcessingNotes string       `json:"processing_notes" gorm:"type:text"`
	
	// Metadata
	CreatedAt       time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Return entity
func (Return) TableName() string {
	return "returns"
}

// IsApproved checks if return is approved
func (r *Return) IsApproved() bool {
	return r.Status == ReturnStatusApproved
}

// IsCompleted checks if return is completed
func (r *Return) IsCompleted() bool {
	return r.Status == ReturnStatusCompleted
}

// validateDimensionsFormat validates dimensions format (LxWxH)
func validateDimensionsFormat(dimensions string) error {
	if dimensions == "" {
		return nil
	}

	// Expected format: "LxWxH" where L, W, H are numbers (can be decimal)
	pattern := `^\d+(\.\d+)?x\d+(\.\d+)?x\d+(\.\d+)?$`
	matched, err := regexp.MatchString(pattern, dimensions)
	if err != nil {
		return fmt.Errorf("regex error: %w", err)
	}

	if !matched {
		return fmt.Errorf("dimensions must be in format 'LxWxH' (e.g., '10x5x3' or '10.5x5.2x3.1')")
	}

	// Parse and validate individual dimensions
	parts := strings.Split(dimensions, "x")
	if len(parts) != 3 {
		return fmt.Errorf("dimensions must have exactly 3 parts (LxWxH)")
	}

	for _, part := range parts {
		value, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return fmt.Errorf("invalid dimension value '%s': %w", part, err)
		}
		if value <= 0 {
			return fmt.Errorf("dimension values must be greater than 0")
		}
		if value > 1000 {
			return fmt.Errorf("dimension values cannot exceed 1000")
		}
	}

	return nil
}

// CanBeProcessed checks if return can be processed
func (r *Return) CanBeProcessed() bool {
	return r.Status == ReturnStatusReceived
}

// ReturnItem represents an item in a return
type ReturnItem struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ReturnID    uuid.UUID `json:"return_id" gorm:"type:uuid;not null;index"`
	Return      Return    `json:"return,omitempty" gorm:"foreignKey:ReturnID"`
	ProductID   uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	Product     Product   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	
	// Item details
	Quantity    int     `json:"quantity" gorm:"not null"`
	UnitPrice   float64 `json:"unit_price" gorm:"not null"`
	TotalPrice  float64 `json:"total_price" gorm:"not null"`
	
	// Return specific
	Reason      ReturnReason `json:"reason" gorm:"not null"`
	Condition   string       `json:"condition"`              // new, used, damaged, etc.
	Notes       string       `json:"notes"`
	
	// Processing
	IsApproved  bool         `json:"is_approved" gorm:"default:false"`
	RefundAmount float64     `json:"refund_amount" gorm:"default:0"`
	
	// Metadata
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for ReturnItem entity
func (ReturnItem) TableName() string {
	return "return_items"
}
