package entities

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending        OrderStatus = "pending"
	OrderStatusConfirmed      OrderStatus = "confirmed"
	OrderStatusProcessing     OrderStatus = "processing"
	OrderStatusReadyToShip    OrderStatus = "ready_to_ship"
	OrderStatusShipped        OrderStatus = "shipped"
	OrderStatusOutForDelivery OrderStatus = "out_for_delivery"
	OrderStatusDelivered      OrderStatus = "delivered"
	OrderStatusCancelled      OrderStatus = "cancelled"
	OrderStatusRefunded       OrderStatus = "refunded"
	OrderStatusReturned       OrderStatus = "returned"
	OrderStatusExchanged      OrderStatus = "exchanged"
)

// FulfillmentStatus represents the fulfillment status of an order
type FulfillmentStatus string

const (
	FulfillmentStatusPending    FulfillmentStatus = "pending"
	FulfillmentStatusProcessing FulfillmentStatus = "processing"
	FulfillmentStatusPacked     FulfillmentStatus = "packed"
	FulfillmentStatusShipped    FulfillmentStatus = "shipped"
	FulfillmentStatusDelivered  FulfillmentStatus = "delivered"
	FulfillmentStatusReturned   FulfillmentStatus = "returned"
	FulfillmentStatusCancelled  FulfillmentStatus = "cancelled"
)

// OrderPriority represents the priority level of an order
type OrderPriority string

const (
	OrderPriorityLow      OrderPriority = "low"
	OrderPriorityNormal   OrderPriority = "normal"
	OrderPriorityHigh     OrderPriority = "high"
	OrderPriorityUrgent   OrderPriority = "urgent"
	OrderPriorityCritical OrderPriority = "critical"
)

// OrderSource represents where the order came from
type OrderSource string

const (
	OrderSourceWeb    OrderSource = "web"
	OrderSourceMobile OrderSource = "mobile"
	OrderSourceAdmin  OrderSource = "admin"
	OrderSourceAPI    OrderSource = "api"
	OrderSourcePhone  OrderSource = "phone"
	OrderSourceEmail  OrderSource = "email"
	OrderSourceSocial OrderSource = "social"
)

// CustomerType represents the type of customer
type CustomerType string

const (
	CustomerTypeGuest      CustomerType = "guest"
	CustomerTypeRegistered CustomerType = "registered"
	CustomerTypeVIP        CustomerType = "vip"
	CustomerTypeWholesale  CustomerType = "wholesale"
	CustomerTypeCorporate  CustomerType = "corporate"
)

// PaymentStatus represents the payment status of an order
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// Order represents an order in the system
type Order struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderNumber string      `json:"order_number" gorm:"uniqueIndex;not null"`
	UserID      uuid.UUID   `json:"user_id" gorm:"type:uuid;not null;index"`
	User        User        `json:"user" gorm:"foreignKey:UserID"`
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`

	// Order Status & Management
	Status            OrderStatus       `json:"status" gorm:"default:'pending'"`
	FulfillmentStatus FulfillmentStatus `json:"fulfillment_status" gorm:"default:'pending'"`
	PaymentStatus     PaymentStatus     `json:"payment_status" gorm:"default:'pending'"`
	Priority          OrderPriority     `json:"priority" gorm:"default:'normal'"`
	Source            OrderSource       `json:"source" gorm:"default:'web'"`
	CustomerType      CustomerType      `json:"customer_type" gorm:"default:'guest'"`

	// Financial Information
	Subtotal       float64 `json:"subtotal" gorm:"not null"`
	TaxAmount      float64 `json:"tax_amount" gorm:"default:0"`
	ShippingAmount float64 `json:"shipping_amount" gorm:"default:0"`
	DiscountAmount float64 `json:"discount_amount" gorm:"default:0"`
	TipAmount      float64 `json:"tip_amount" gorm:"default:0"`
	Total          float64 `json:"total" gorm:"not null"`
	Currency       string  `json:"currency" gorm:"default:'USD'"`

	// Address Information
	ShippingAddress *OrderAddress `json:"shipping_address" gorm:"embedded;embeddedPrefix:shipping_"`
	BillingAddress  *OrderAddress `json:"billing_address" gorm:"embedded;embeddedPrefix:billing_"`

	// Shipping & Delivery
	ShippingMethod       string     `json:"shipping_method"`
	TrackingNumber       string     `json:"tracking_number"`
	TrackingURL          string     `json:"tracking_url"`
	Carrier              string     `json:"carrier"`
	EstimatedDelivery    *time.Time `json:"estimated_delivery"`
	ActualDelivery       *time.Time `json:"actual_delivery"`
	DeliveryInstructions string     `json:"delivery_instructions" gorm:"type:text"`
	DeliveryAttempts     int        `json:"delivery_attempts" gorm:"default:0"`

	// Customer Information
	CustomerNotes string `json:"customer_notes" gorm:"type:text"`
	AdminNotes    string `json:"admin_notes" gorm:"type:text"`
	InternalNotes string `json:"internal_notes" gorm:"type:text"`

	// Gift Options
	IsGift      bool   `json:"is_gift" gorm:"default:false"`
	GiftMessage string `json:"gift_message" gorm:"type:text"`
	GiftWrap    bool   `json:"gift_wrap" gorm:"default:false"`

	// Business Information
	SalesChannel   string `json:"sales_channel"`
	ReferralSource string `json:"referral_source"`
	CouponCodes    string `json:"coupon_codes" gorm:"type:text"` // JSON array as string
	Tags           string `json:"tags" gorm:"type:text"`         // JSON array as string

	// Fulfillment Information
	WarehouseID *uuid.UUID `json:"warehouse_id" gorm:"type:uuid"`
	PackedAt    *time.Time `json:"packed_at"`
	ShippedAt   *time.Time `json:"shipped_at"`
	ProcessedAt *time.Time `json:"processed_at"`

	// Stock reservation fields
	InventoryReserved bool       `json:"inventory_reserved" gorm:"default:false"`
	ReservedUntil     *time.Time `json:"reserved_until"`
	PaymentTimeout    *time.Time `json:"payment_timeout"`

	// Relationships
	Payment     *Payment     `json:"payment" gorm:"foreignKey:OrderID"`
	OrderEvents []OrderEvent `json:"order_events" gorm:"foreignKey:OrderID"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Order entity
func (Order) TableName() string {
	return "orders"
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID     uuid.UUID `json:"order_id" gorm:"type:uuid;not null;index"`
	ProductID   uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	Product     Product   `json:"product" gorm:"foreignKey:ProductID"`
	ProductName string    `json:"product_name" gorm:"not null"`
	ProductSKU  string    `json:"product_sku" gorm:"not null"`
	Quantity    int       `json:"quantity" gorm:"not null" validate:"required,gt=0"`
	Price       float64   `json:"price" gorm:"not null"`
	Total       float64   `json:"total" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for OrderItem entity
func (OrderItem) TableName() string {
	return "order_items"
}

// OrderAddress represents an address for orders
type OrderAddress struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Company   string `json:"company"`
	Address1  string `json:"address1" validate:"required"`
	Address2  string `json:"address2"`
	City      string `json:"city" validate:"required"`
	State     string `json:"state" validate:"required"`
	ZipCode   string `json:"zip_code" validate:"required"`
	Country   string `json:"country" validate:"required"`
	Phone     string `json:"phone"`
}

// GetFullName returns the full name from the address
func (a *OrderAddress) GetFullName() string {
	return a.FirstName + " " + a.LastName
}

// OrderEventType represents the type of order event
type OrderEventType string

const (
	OrderEventTypeCreated           OrderEventType = "created"
	OrderEventTypeStatusChanged     OrderEventType = "status_changed"
	OrderEventTypePaymentReceived   OrderEventType = "payment_received"
	OrderEventTypePaymentFailed     OrderEventType = "payment_failed"
	OrderEventTypeShipped           OrderEventType = "shipped"
	OrderEventTypeDelivered         OrderEventType = "delivered"
	OrderEventTypeCancelled         OrderEventType = "cancelled"
	OrderEventTypeRefunded          OrderEventType = "refunded"
	OrderEventTypeReturned          OrderEventType = "returned"
	OrderEventTypeNoteAdded         OrderEventType = "note_added"
	OrderEventTypeTrackingUpdated   OrderEventType = "tracking_updated"
	OrderEventTypeInventoryReserved OrderEventType = "inventory_reserved"
	OrderEventTypeInventoryReleased OrderEventType = "inventory_released"
	OrderEventTypeCustom            OrderEventType = "custom"
)

// OrderEvent represents an event in the order lifecycle
type OrderEvent struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID     uuid.UUID      `json:"order_id" gorm:"type:uuid;not null;index"`
	EventType   OrderEventType `json:"event_type" gorm:"not null"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Data        string         `json:"data" gorm:"type:text"` // JSON data as string
	UserID      *uuid.UUID     `json:"user_id" gorm:"type:uuid"`
	User        *User          `json:"user" gorm:"foreignKey:UserID"`
	IsPublic    bool           `json:"is_public" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for OrderEvent entity
func (OrderEvent) TableName() string {
	return "order_events"
}

// GetFullAddress returns the formatted full address
func (a *OrderAddress) GetFullAddress() string {
	address := a.Address1
	if a.Address2 != "" {
		address += ", " + a.Address2
	}
	address += ", " + a.City + ", " + a.State + " " + a.ZipCode + ", " + a.Country
	return address
}

// CanBeCancelled checks if the order can be cancelled
func (o *Order) CanBeCancelled() bool {
	// Can cancel if order is pending/confirmed and not yet shipped
	if o.Status == OrderStatusShipped || o.Status == OrderStatusOutForDelivery ||
		o.Status == OrderStatusDelivered || o.Status == OrderStatusCancelled ||
		o.Status == OrderStatusRefunded || o.Status == OrderStatusReturned ||
		o.Status == OrderStatusExchanged {
		return false
	}
	return true
}

// CanBeRefunded checks if the order can be refunded
func (o *Order) CanBeRefunded() bool {
	return o.PaymentStatus == PaymentStatusPaid &&
		(o.Status == OrderStatusDelivered || o.Status == OrderStatusShipped)
}

// IsCompleted checks if the order is completed
func (o *Order) IsCompleted() bool {
	return o.Status == OrderStatusDelivered
}

// IsPaid checks if the order is paid
func (o *Order) IsPaid() bool {
	return o.PaymentStatus == PaymentStatusPaid
}

// IsReservationExpired checks if inventory reservation has expired
func (o *Order) IsReservationExpired() bool {
	if o.ReservedUntil == nil {
		return false
	}
	return time.Now().After(*o.ReservedUntil)
}

// IsPaymentExpired checks if payment timeout has expired
func (o *Order) IsPaymentExpired() bool {
	if o.PaymentTimeout == nil {
		return false
	}
	return time.Now().After(*o.PaymentTimeout)
}

// HasInventoryReserved checks if inventory is currently reserved
func (o *Order) HasInventoryReserved() bool {
	return o.InventoryReserved && !o.IsReservationExpired()
}

// SetReservationTimeout sets the reservation timeout (default 30 minutes)
func (o *Order) SetReservationTimeout(minutes int) {
	if minutes <= 0 {
		minutes = 30 // Default 30 minutes
	}
	timeout := time.Now().Add(time.Duration(minutes) * time.Minute)
	o.ReservedUntil = &timeout
	o.InventoryReserved = true
}

// SetPaymentTimeout sets the payment timeout (default 24 hours)
func (o *Order) SetPaymentTimeout(hours int) {
	if hours <= 0 {
		hours = 24 // Default 24 hours
	}
	timeout := time.Now().Add(time.Duration(hours) * time.Hour)
	o.PaymentTimeout = &timeout
}

// ReleaseReservation releases the inventory reservation
func (o *Order) ReleaseReservation() {
	o.InventoryReserved = false
	o.ReservedUntil = nil
}

// GetItemCount returns the total number of items in the order
func (o *Order) GetItemCount() int {
	count := 0
	for _, item := range o.Items {
		count += item.Quantity
	}
	return count
}

// CalculateTotal calculates the total amount of the order
func (o *Order) CalculateTotal() {
	o.Total = o.Subtotal + o.TaxAmount + o.ShippingAmount + o.TipAmount - o.DiscountAmount
}

// CanBeShipped checks if the order can be shipped
func (o *Order) CanBeShipped() bool {
	return o.Status == OrderStatusConfirmed || o.Status == OrderStatusProcessing || o.Status == OrderStatusReadyToShip
}

// CanBeDelivered checks if the order can be marked as delivered
func (o *Order) CanBeDelivered() bool {
	return o.Status == OrderStatusShipped || o.Status == OrderStatusOutForDelivery
}

// CanBeReturned checks if the order can be returned
func (o *Order) CanBeReturned() bool {
	return o.Status == OrderStatusDelivered && o.PaymentStatus == PaymentStatusPaid
}

// IsShipped checks if the order has been shipped
func (o *Order) IsShipped() bool {
	return o.Status == OrderStatusShipped || o.Status == OrderStatusOutForDelivery || o.Status == OrderStatusDelivered
}

// IsDelivered checks if the order has been delivered
func (o *Order) IsDelivered() bool {
	return o.Status == OrderStatusDelivered
}

// HasTracking checks if the order has tracking information
func (o *Order) HasTracking() bool {
	return o.TrackingNumber != ""
}

// IsGiftOrder checks if the order is a gift
func (o *Order) IsGiftOrder() bool {
	return o.IsGift
}

// GetStatusDisplayName returns a human-readable status name
func (o *Order) GetStatusDisplayName() string {
	switch o.Status {
	case OrderStatusPending:
		return "Pending"
	case OrderStatusConfirmed:
		return "Confirmed"
	case OrderStatusProcessing:
		return "Processing"
	case OrderStatusReadyToShip:
		return "Ready to Ship"
	case OrderStatusShipped:
		return "Shipped"
	case OrderStatusOutForDelivery:
		return "Out for Delivery"
	case OrderStatusDelivered:
		return "Delivered"
	case OrderStatusCancelled:
		return "Cancelled"
	case OrderStatusRefunded:
		return "Refunded"
	case OrderStatusReturned:
		return "Returned"
	case OrderStatusExchanged:
		return "Exchanged"
	default:
		return string(o.Status)
	}
}

// GetPriorityDisplayName returns a human-readable priority name
func (o *Order) GetPriorityDisplayName() string {
	switch o.Priority {
	case OrderPriorityLow:
		return "Low"
	case OrderPriorityNormal:
		return "Normal"
	case OrderPriorityHigh:
		return "High"
	case OrderPriorityUrgent:
		return "Urgent"
	case OrderPriorityCritical:
		return "Critical"
	default:
		return string(o.Priority)
	}
}

// SetShipped marks the order as shipped
func (o *Order) SetShipped(trackingNumber, carrier string) {
	o.Status = OrderStatusShipped
	o.FulfillmentStatus = FulfillmentStatusShipped
	o.TrackingNumber = trackingNumber
	o.Carrier = carrier
	now := time.Now()
	o.ShippedAt = &now
	o.UpdatedAt = now
}

// SetDelivered marks the order as delivered
func (o *Order) SetDelivered() {
	o.Status = OrderStatusDelivered
	o.FulfillmentStatus = FulfillmentStatusDelivered
	now := time.Now()
	o.ActualDelivery = &now
	o.UpdatedAt = now
}

// SetProcessing marks the order as processing
func (o *Order) SetProcessing() {
	o.Status = OrderStatusProcessing
	o.FulfillmentStatus = FulfillmentStatusProcessing
	now := time.Now()
	o.ProcessedAt = &now
	o.UpdatedAt = now
}
