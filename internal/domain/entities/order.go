package entities

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Get payment timeout (minutes) from env/config
func getOrderTimeoutMinutes() int {
	val := os.Getenv("ORDER_PAYMENT_TIMEOUT_MINUTES")
	if val != "" {
		if minutes, err := strconv.Atoi(val); err == nil {
			return minutes
		}
	}
	return 30 // default 30 phút
}

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

// PaymentStatus is now defined in payment.go to avoid duplication

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
	PaymentMethod     PaymentMethod     `json:"payment_method" gorm:"default:'credit_card'"` // Store payment method
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
	ReservedUntil     *time.Time `json:"reserved_until" gorm:"index"`  // Index for cleanup jobs
	PaymentTimeout    *time.Time `json:"payment_timeout" gorm:"index"` // Index for cleanup jobs

	// Audit fields
	Version        int        `json:"version" gorm:"default:1"` // For optimistic locking
	LastModifiedBy *uuid.UUID `json:"last_modified_by" gorm:"type:uuid"`

	// Relationships
	Payments    []Payment    `json:"payments" gorm:"foreignKey:OrderID"`
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
	Weight      float64   `json:"weight" gorm:"default:0"` // Individual item weight for shipping calculation
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for OrderItem entity
func (OrderItem) TableName() string {
	return "order_items"
}

// Validate validates order item data
func (oi *OrderItem) Validate() error {
	if oi.ProductID == uuid.Nil {
		return fmt.Errorf("product ID is required")
	}
	if oi.ProductName == "" {
		return fmt.Errorf("product name is required")
	}
	if oi.ProductSKU == "" {
		return fmt.Errorf("product SKU is required")
	}
	if oi.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}
	if oi.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	if oi.Total < 0 {
		return fmt.Errorf("total cannot be negative")
	}

	// Verify that total matches price * quantity with floating point tolerance
	expectedTotal := oi.Price * float64(oi.Quantity)
	const epsilon = 0.01
	if math.Abs(oi.Total - expectedTotal) > epsilon {
		return fmt.Errorf("total %.2f does not match price %.2f * quantity %d = %.2f",
			oi.Total, oi.Price, oi.Quantity, expectedTotal)
	}

	return nil
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

// Validate validates order address data
func (a *OrderAddress) Validate() error {
	if a.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if a.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if a.Address1 == "" {
		return fmt.Errorf("address line 1 is required")
	}
	if a.City == "" {
		return fmt.Errorf("city is required")
	}
	if a.State == "" {
		return fmt.Errorf("state is required")
	}
	if a.ZipCode == "" {
		return fmt.Errorf("zip code is required")
	}
	if a.Country == "" {
		return fmt.Errorf("country is required")
	}

	// Validate zip code format (basic validation)
	if len(a.ZipCode) < 3 || len(a.ZipCode) > 20 {
		return fmt.Errorf("zip code must be between 3 and 20 characters")
	}

	return nil
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
	// Can refund if payment is completed and order is not already refunded/cancelled
	return o.PaymentStatus == PaymentStatusPaid &&
		o.Status != OrderStatusCancelled &&
		o.Status != OrderStatusRefunded &&
		o.Status != OrderStatusReturned
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

// ValidateTimeouts validates and sets default timeouts if not set
func (o *Order) ValidateTimeouts() {
	if o.ReservedUntil == nil && o.InventoryReserved {
		o.SetReservationTimeout(30) // 30 minutes for stock reservation
	}
	if o.PaymentTimeout == nil && o.Status == OrderStatusPending {
		o.SetPaymentTimeout(24) // 24 hours for payment
	}
}

// IncrementVersion increments the version for optimistic locking
func (o *Order) IncrementVersion() {
	o.Version++
	o.UpdatedAt = time.Now()
}

// ReleaseReservation releases the inventory reservation
func (o *Order) ReleaseReservation() {
	o.InventoryReserved = false
	o.ReservedUntil = nil
	// Optionally, update stock here if needed
}

// Validate validates order data
func (o *Order) Validate() error {
	// Validate required fields
	if o.OrderNumber == "" {
		return fmt.Errorf("order number is required")
	}
	if o.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}
	if len(o.Items) == 0 {
		return fmt.Errorf("order must have at least one item")
	}

	// Validate financial fields
	if o.Subtotal < 0 {
		return fmt.Errorf("subtotal cannot be negative")
	}
	if o.TaxAmount < 0 {
		return fmt.Errorf("tax amount cannot be negative")
	}
	if o.ShippingAmount < 0 {
		return fmt.Errorf("shipping amount cannot be negative")
	}
	if o.DiscountAmount < 0 {
		return fmt.Errorf("discount amount cannot be negative")
	}
	if o.TipAmount < 0 {
		return fmt.Errorf("tip amount cannot be negative")
	}
	if o.Total < 0 {
		return fmt.Errorf("total cannot be negative")
	}

	// Validate total calculation with floating point tolerance
	expectedTotal := o.Subtotal + o.TaxAmount + o.ShippingAmount + o.TipAmount - o.DiscountAmount
	const epsilon = 0.01
	if math.Abs(o.Total - expectedTotal) > epsilon {
		return fmt.Errorf("total %.2f does not match calculated total %.2f", o.Total, expectedTotal)
	}

	// Validate currency
	if o.Currency == "" {
		return fmt.Errorf("currency is required")
	}

	// Validate order items
	for i, item := range o.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("item %d validation failed: %w", i, err)
		}
	}

	// Validate addresses if present
	if o.ShippingAddress != nil {
		if err := o.ShippingAddress.Validate(); err != nil {
			return fmt.Errorf("shipping address validation failed: %w", err)
		}
	}
	if o.BillingAddress != nil {
		if err := o.BillingAddress.Validate(); err != nil {
			return fmt.Errorf("billing address validation failed: %w", err)
		}
	}

	return nil
}

// CanTransitionTo checks if order can transition to the given status
func (o *Order) CanTransitionTo(newStatus OrderStatus) bool {
	switch o.Status {
	case OrderStatusPending:
		return newStatus == OrderStatusConfirmed || newStatus == OrderStatusCancelled
	case OrderStatusConfirmed:
		return newStatus == OrderStatusProcessing || newStatus == OrderStatusCancelled
	case OrderStatusProcessing:
		return newStatus == OrderStatusReadyToShip || newStatus == OrderStatusCancelled
	case OrderStatusReadyToShip:
		return newStatus == OrderStatusShipped || newStatus == OrderStatusCancelled
	case OrderStatusShipped:
		return newStatus == OrderStatusOutForDelivery || newStatus == OrderStatusDelivered || newStatus == OrderStatusReturned
	case OrderStatusOutForDelivery:
		return newStatus == OrderStatusDelivered || newStatus == OrderStatusReturned
	case OrderStatusDelivered:
		return newStatus == OrderStatusReturned || newStatus == OrderStatusExchanged || newStatus == OrderStatusRefunded
	case OrderStatusCancelled, OrderStatusRefunded, OrderStatusReturned, OrderStatusExchanged:
		return false // Terminal states
	default:
		return false
	}
}

// TransitionTo transitions order to new status with validation
func (o *Order) TransitionTo(newStatus OrderStatus) error {
	if !o.CanTransitionTo(newStatus) {
		return fmt.Errorf("cannot transition from %s to %s", o.Status, newStatus)
	}

	o.Status = newStatus
	o.UpdatedAt = time.Now()

	// Update related fields based on status
	switch newStatus {
	case OrderStatusProcessing:
		if o.ProcessedAt == nil {
			now := time.Now()
			o.ProcessedAt = &now
		}
	case OrderStatusShipped:
		if o.ShippedAt == nil {
			now := time.Now()
			o.ShippedAt = &now
		}
	case OrderStatusDelivered:
		if o.ActualDelivery == nil {
			now := time.Now()
			o.ActualDelivery = &now
		}
	}

	// Sync fulfillment status with order status
	o.syncFulfillmentStatus()

	return nil
}

// syncFulfillmentStatus syncs fulfillment status with order status
func (o *Order) syncFulfillmentStatus() {
	switch o.Status {
	case OrderStatusPending, OrderStatusConfirmed:
		o.FulfillmentStatus = FulfillmentStatusPending
	case OrderStatusProcessing:
		o.FulfillmentStatus = FulfillmentStatusProcessing
	case OrderStatusReadyToShip:
		o.FulfillmentStatus = FulfillmentStatusPacked
	case OrderStatusShipped, OrderStatusOutForDelivery:
		o.FulfillmentStatus = FulfillmentStatusShipped
	case OrderStatusDelivered:
		o.FulfillmentStatus = FulfillmentStatusDelivered
	case OrderStatusReturned:
		o.FulfillmentStatus = FulfillmentStatusReturned
	case OrderStatusCancelled:
		o.FulfillmentStatus = FulfillmentStatusCancelled
	}
}

// SyncPaymentStatus syncs order payment status with actual payment status
func (o *Order) SyncPaymentStatus(paymentStatus PaymentStatus) {
	o.PaymentStatus = paymentStatus
	o.UpdatedAt = time.Now()

	// Auto-transition order status based on payment status
	if paymentStatus == PaymentStatusCompleted && o.Status == OrderStatusPending {
		o.TransitionTo(OrderStatusConfirmed)
	} else if paymentStatus == PaymentStatusFailed && o.Status == OrderStatusPending {
		o.TransitionTo(OrderStatusCancelled)
	}
}

// Call ValidateTimeouts when creating new order (example constructor)
func NewOrder(userID uuid.UUID, items []OrderItem) *Order {
	order := &Order{
		ID:            uuid.New(),
		UserID:        userID,
		Items:         items,
		Status:        OrderStatusPending,
		PaymentStatus: PaymentStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	order.ValidateTimeouts()
	return order
}

// Notify customer when order is about to expire (pseudo-code, implement in service layer)
func (o *Order) NotifyPaymentExpiring() {
	// Example: send email/SMS/push notification
	// This should be called by a scheduler before o.PaymentTimeout
	// e.g. if time.Until(*o.PaymentTimeout) < 5*time.Minute {
	// sendNotification(o.UserID, "Đơn hàng của bạn sắp hết hạn thanh toán!")
	// }
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

// GetSuccessfulPayments returns all successful payments for this order
func (o *Order) GetSuccessfulPayments() []Payment {
	var successfulPayments []Payment
	for _, payment := range o.Payments {
		if payment.IsSuccessful() {
			successfulPayments = append(successfulPayments, payment)
		}
	}
	return successfulPayments
}

// GetTotalPaidAmount returns the total amount paid for this order
func (o *Order) GetTotalPaidAmount() float64 {
	total := 0.0
	for _, payment := range o.Payments {
		if payment.IsSuccessful() {
			total += payment.Amount
		}
	}
	return total
}

// GetLatestPayment returns the most recent payment for this order
func (o *Order) GetLatestPayment() *Payment {
	if len(o.Payments) == 0 {
		return nil
	}

	latest := &o.Payments[0]
	for i := 1; i < len(o.Payments); i++ {
		if o.Payments[i].CreatedAt.After(latest.CreatedAt) {
			latest = &o.Payments[i]
		}
	}
	return latest
}

// IsFullyPaid checks if the order is fully paid
func (o *Order) IsFullyPaid() bool {
	return o.GetTotalPaidAmount() >= o.Total
}

// IsPartiallyPaid checks if the order is partially paid
func (o *Order) IsPartiallyPaid() bool {
	paidAmount := o.GetTotalPaidAmount()
	return paidAmount > 0 && paidAmount < o.Total
}

// GetRemainingAmount returns the remaining amount to be paid
func (o *Order) GetRemainingAmount() float64 {
	remaining := o.Total - o.GetTotalPaidAmount()
	if remaining < 0 {
		return 0
	}
	return remaining
}

// AutoSyncPaymentStatus automatically synchronizes the order payment status with actual payments
func (o *Order) AutoSyncPaymentStatus() {
	// If no payments exist, determine status based on payment method
	if len(o.Payments) == 0 {
		if o.PaymentMethod == PaymentMethodCash {
			o.PaymentStatus = PaymentStatusAwaitingPayment
		} else {
			o.PaymentStatus = PaymentStatusPending
		}
		return
	}

	// Check payment completion status
	if o.IsFullyPaid() {
		o.PaymentStatus = PaymentStatusPaid
		return
	}

	if o.IsPartiallyPaid() {
		o.PaymentStatus = PaymentStatusPartiallyPaid
		return
	}

	// Check for failed payments
	hasFailedPayments := false
	hasPendingPayments := false
	hasProcessingPayments := false
	isCODOrder := o.PaymentMethod == PaymentMethodCash

	for _, payment := range o.Payments {
		switch payment.Status {
		case PaymentStatusFailed:
			hasFailedPayments = true
		case PaymentStatusPending:
			hasPendingPayments = true
		case PaymentStatusProcessing:
			hasProcessingPayments = true
		}
	}

	// Determine status based on payment states and method
	if hasFailedPayments && !hasPendingPayments && !hasProcessingPayments {
		o.PaymentStatus = PaymentStatusFailed
	} else if hasProcessingPayments {
		o.PaymentStatus = PaymentStatusProcessing
	} else if isCODOrder {
		o.PaymentStatus = PaymentStatusAwaitingPayment
	} else {
		o.PaymentStatus = PaymentStatusPending
	}
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
