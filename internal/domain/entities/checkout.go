package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CheckoutSessionStatus represents the status of a checkout session
type CheckoutSessionStatus string

const (
	CheckoutSessionStatusActive    CheckoutSessionStatus = "active"
	CheckoutSessionStatusCompleted CheckoutSessionStatus = "completed"
	CheckoutSessionStatusExpired   CheckoutSessionStatus = "expired"
	CheckoutSessionStatusCancelled CheckoutSessionStatus = "cancelled"
)

// CheckoutSession represents a checkout session before order creation
type CheckoutSession struct {
	ID        uuid.UUID             `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID             `json:"user_id" gorm:"type:uuid;not null;index"`
	User      User                  `json:"user" gorm:"foreignKey:UserID"`
	SessionID string                `json:"session_id" gorm:"uniqueIndex;not null"` // For tracking
	Status    CheckoutSessionStatus `json:"status" gorm:"default:'active'"`

	// Cart snapshot at checkout time
	CartID    uuid.UUID  `json:"cart_id" gorm:"type:uuid;not null"`
	CartItems []CartItem `json:"cart_items" gorm:"serializer:json"` // Snapshot of cart items

	// Address Information
	ShippingAddress *OrderAddress `json:"shipping_address" gorm:"embedded;embeddedPrefix:shipping_"`
	BillingAddress  *OrderAddress `json:"billing_address" gorm:"embedded;embeddedPrefix:billing_"`

	// Payment Information
	PaymentMethod   PaymentMethod `json:"payment_method" gorm:"not null"`
	PaymentIntentID string        `json:"payment_intent_id"` // For Stripe/PayPal

	// Financial Information
	Subtotal       float64 `json:"subtotal" gorm:"not null"`
	TaxAmount      float64 `json:"tax_amount" gorm:"default:0"`
	ShippingAmount float64 `json:"shipping_amount" gorm:"default:0"`
	DiscountAmount float64 `json:"discount_amount" gorm:"default:0"`
	Total          float64 `json:"total" gorm:"not null"`
	Currency       string  `json:"currency" gorm:"default:'USD'"`

	// Tax and shipping details
	TaxRate      float64 `json:"tax_rate" gorm:"default:0"`
	ShippingCost float64 `json:"shipping_cost" gorm:"default:0"`

	// Customer notes
	Notes string `json:"notes"`

	// Timeout and expiration
	ExpiresAt *time.Time `json:"expires_at" gorm:"index"` // For cleanup jobs

	// Result
	OrderID *uuid.UUID `json:"order_id" gorm:"type:uuid"` // Set when order is created

	// Timestamps
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for CheckoutSession entity
func (CheckoutSession) TableName() string {
	return "checkout_sessions"
}

// IsExpired checks if the checkout session has expired
func (cs *CheckoutSession) IsExpired() bool {
	if cs.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*cs.ExpiresAt)
}

// SetExpiration sets the expiration time for the checkout session
func (cs *CheckoutSession) SetExpiration(minutes int) {
	expiresAt := time.Now().Add(time.Duration(minutes) * time.Minute)
	cs.ExpiresAt = &expiresAt
}

// MarkAsCompleted marks the checkout session as completed
func (cs *CheckoutSession) MarkAsCompleted(orderID uuid.UUID) {
	cs.Status = CheckoutSessionStatusCompleted
	cs.OrderID = &orderID
	cs.UpdatedAt = time.Now()
}

// MarkAsExpired marks the checkout session as expired
func (cs *CheckoutSession) MarkAsExpired() {
	cs.Status = CheckoutSessionStatusExpired
	cs.UpdatedAt = time.Now()
}

// MarkAsCancelled marks the checkout session as cancelled
func (cs *CheckoutSession) MarkAsCancelled() {
	cs.Status = CheckoutSessionStatusCancelled
	cs.UpdatedAt = time.Now()
}

// CanBeCompleted checks if the checkout session can be completed
func (cs *CheckoutSession) CanBeCompleted() bool {
	return cs.Status == CheckoutSessionStatusActive && !cs.IsExpired()
}

// Validate validates the checkout session data
func (cs *CheckoutSession) Validate() error {
	// Validate required fields
	if cs.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}
	if cs.CartID == uuid.Nil {
		return fmt.Errorf("cart ID is required")
	}
	if cs.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if len(cs.CartItems) == 0 {
		return fmt.Errorf("cart items are required")
	}

	// Validate financial fields
	if cs.Subtotal < 0 {
		return fmt.Errorf("subtotal cannot be negative")
	}
	if cs.TaxAmount < 0 {
		return fmt.Errorf("tax amount cannot be negative")
	}
	if cs.ShippingAmount < 0 {
		return fmt.Errorf("shipping amount cannot be negative")
	}
	if cs.DiscountAmount < 0 {
		return fmt.Errorf("discount amount cannot be negative")
	}
	if cs.Total < 0 {
		return fmt.Errorf("total cannot be negative")
	}

	// Validate currency
	if cs.Currency == "" {
		return fmt.Errorf("currency is required")
	}

	// Validate shipping address
	if cs.ShippingAddress == nil {
		return fmt.Errorf("shipping address is required")
	}

	return nil
}

// GenerateSessionID generates a unique session ID
func (cs *CheckoutSession) GenerateSessionID() {
	cs.SessionID = fmt.Sprintf("checkout_%s_%d", cs.ID.String()[:8], time.Now().Unix())
}
