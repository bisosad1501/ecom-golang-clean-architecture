package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PaymentMethod represents the payment method
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodDebitCard  PaymentMethod = "debit_card"
	PaymentMethodPayPal     PaymentMethod = "paypal"
	PaymentMethodStripe     PaymentMethod = "stripe"
	PaymentMethodApplePay   PaymentMethod = "apple_pay"
	PaymentMethodGooglePay  PaymentMethod = "google_pay"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodCash       PaymentMethod = "cash"
)

// PaymentStatus represents the payment status
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"  // Added for frontend compatibility
	PaymentStatusPaid       PaymentStatus = "paid"
	PaymentStatusCompleted  PaymentStatus = "completed"   // Alias for paid
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusRefunded   PaymentStatus = "refunded"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
)

// Payment represents a payment transaction
type Payment struct {
	ID                uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID           uuid.UUID     `json:"order_id" gorm:"type:uuid;not null;uniqueIndex"` // One payment per order
	UserID            uuid.UUID     `json:"user_id" gorm:"type:uuid;index"`
	Amount            float64       `json:"amount" gorm:"not null" validate:"required,gt=0"`
	Currency          string        `json:"currency" gorm:"default:'USD'"`
	Method            PaymentMethod `json:"method" gorm:"not null" validate:"required"`
	Status            PaymentStatus `json:"status" gorm:"default:'pending'"`

	// Transaction identifiers
	TransactionID     string        `json:"transaction_id" gorm:"index"`
	ExternalID        string        `json:"external_id" gorm:"index"`
	PaymentIntentID   string        `json:"payment_intent_id" gorm:"index"` // For Stripe

	// Gateway information
	Gateway           string        `json:"gateway" gorm:"default:'stripe'"` // stripe, paypal, etc.
	GatewayResponse   string        `json:"gateway_response" gorm:"type:text"`

	// Fees and charges
	ProcessingFee     float64       `json:"processing_fee" gorm:"default:0"`
	GatewayFee        float64       `json:"gateway_fee" gorm:"default:0"`
	NetAmount         float64       `json:"net_amount" gorm:"default:0"` // Amount - fees

	// Failure information
	FailureReason     string        `json:"failure_reason"`
	FailureCode       string        `json:"failure_code"`

	// Timestamps
	ProcessedAt       *time.Time    `json:"processed_at"`
	RefundedAt        *time.Time    `json:"refunded_at"`

	// Refund information
	RefundAmount      float64       `json:"refund_amount" gorm:"default:0"`
	RefundReason      string        `json:"refund_reason"`

	// Metadata
	Metadata          string        `json:"metadata" gorm:"type:text"` // JSON metadata
	Notes             string        `json:"notes" gorm:"type:text"`

	CreatedAt         time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time     `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Order             *Order        `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	User              *User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Refunds           []Refund      `json:"refunds,omitempty" gorm:"foreignKey:PaymentID"`
}

// TableName returns the table name for Payment entity
func (Payment) TableName() string {
	return "payments"
}

// IsSuccessful checks if the payment is successful
func (p *Payment) IsSuccessful() bool {
	return p.Status == PaymentStatusPaid
}

// IsFailed checks if the payment failed
func (p *Payment) IsFailed() bool {
	return p.Status == PaymentStatusFailed
}

// IsPending checks if the payment is pending
func (p *Payment) IsPending() bool {
	return p.Status == PaymentStatusPending
}

// IsRefunded checks if the payment is refunded
func (p *Payment) IsRefunded() bool {
	return p.Status == PaymentStatusRefunded
}

// MarkAsProcessed marks the payment as processed
func (p *Payment) MarkAsProcessed(transactionID string) {
	p.Status = PaymentStatusPaid
	p.TransactionID = transactionID
	now := time.Now()
	p.ProcessedAt = &now
	p.UpdatedAt = now

	// Calculate net amount if not set
	if p.NetAmount == 0 {
		p.NetAmount = p.Amount - p.ProcessingFee - p.GatewayFee
	}
}

// SetFees sets the processing and gateway fees
func (p *Payment) SetFees(processingFee, gatewayFee float64) {
	p.ProcessingFee = processingFee
	p.GatewayFee = gatewayFee
	p.NetAmount = p.Amount - processingFee - gatewayFee
	p.UpdatedAt = time.Now()
}

// MarkAsFailed marks the payment as failed
func (p *Payment) MarkAsFailed(reason string) {
	p.Status = PaymentStatusFailed
	p.FailureReason = reason
	p.UpdatedAt = time.Now()
}

// AddRefund adds a refund to the payment
func (p *Payment) AddRefund(amount float64) error {
	if amount <= 0 {
		return ErrInvalidRefundAmount
	}
	
	if p.RefundAmount + amount > p.Amount {
		return ErrRefundAmountExceedsPayment
	}
	
	p.RefundAmount += amount
	
	if p.RefundAmount >= p.Amount {
		p.Status = PaymentStatusRefunded
		now := time.Now()
		p.RefundedAt = &now
	}
	
	p.UpdatedAt = time.Now()
	return nil
}

// Refund represents a payment refund
type Refund struct {
	ID            uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PaymentID     uuid.UUID     `json:"payment_id" gorm:"type:uuid;not null;index"`
	Amount        float64       `json:"amount" gorm:"not null" validate:"required,gt=0"`
	Reason        string        `json:"reason" gorm:"not null"`
	Status        RefundStatus  `json:"status" gorm:"default:'pending'"`
	TransactionID string        `json:"transaction_id" gorm:"index"`
	ExternalID    string        `json:"external_id" gorm:"index"`
	ProcessedAt   *time.Time    `json:"processed_at"`
	CreatedAt     time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Payment *Payment `json:"payment,omitempty" gorm:"foreignKey:PaymentID"`
}

// RefundStatus represents the refund status
type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusCompleted RefundStatus = "completed"
	RefundStatusFailed    RefundStatus = "failed"
	RefundStatusCancelled RefundStatus = "cancelled"
)

// TableName returns the table name for Refund entity
func (Refund) TableName() string {
	return "refunds"
}

// IsCompleted checks if the refund is completed
func (r *Refund) IsCompleted() bool {
	return r.Status == RefundStatusCompleted
}

// MarkAsCompleted marks the refund as completed
func (r *Refund) MarkAsCompleted(transactionID string) {
	r.Status = RefundStatusCompleted
	r.TransactionID = transactionID
	now := time.Now()
	r.ProcessedAt = &now
	r.UpdatedAt = now
}

// MarkAsFailed marks the refund as failed
func (r *Refund) MarkAsFailed() {
	r.Status = RefundStatusFailed
	r.UpdatedAt = time.Now()
}

// GetRemainingRefundAmount returns the remaining amount that can be refunded
func (p *Payment) GetRemainingRefundAmount() float64 {
	return p.Amount - p.RefundAmount
}

// CanBeRefunded checks if the payment can be refunded
func (p *Payment) CanBeRefunded() bool {
	return p.Status == PaymentStatusPaid && p.RefundAmount < p.Amount
}

// PaymentMethodEntity represents a saved payment method for a user
type PaymentMethodEntity struct {
	ID            uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID     `json:"user_id" gorm:"type:uuid;not null;index"`
	Type          PaymentMethod `json:"type" gorm:"not null" validate:"required"`

	// Card information (for card payments)
	Last4         string        `json:"last4"`                    // Last 4 digits
	Brand         string        `json:"brand"`                    // visa, mastercard, amex, etc.
	ExpiryMonth   int           `json:"expiry_month"`
	ExpiryYear    int           `json:"expiry_year"`

	// Gateway information
	Gateway       string        `json:"gateway" gorm:"default:'stripe'"` // stripe, paypal, etc.
	GatewayToken  string        `json:"gateway_token" gorm:"not null"`   // Token from payment gateway
	GatewayCustomerID string    `json:"gateway_customer_id"`             // Customer ID in gateway

	// Billing information
	BillingName   string        `json:"billing_name"`
	BillingEmail  string        `json:"billing_email"`
	BillingAddress string       `json:"billing_address" gorm:"type:text"`

	// Status and preferences
	IsDefault     bool          `json:"is_default" gorm:"default:false"`
	IsActive      bool          `json:"is_active" gorm:"default:true"`

	// Security
	Fingerprint   string        `json:"fingerprint" gorm:"index"`        // Unique fingerprint to prevent duplicates

	// Metadata
	Metadata      string        `json:"metadata" gorm:"type:text"`       // JSON metadata
	Notes         string        `json:"notes"`

	// Timestamps
	CreatedAt     time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	LastUsedAt    *time.Time    `json:"last_used_at"`

	// Relationships
	User          *User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for PaymentMethodEntity
func (PaymentMethodEntity) TableName() string {
	return "payment_methods"
}

// IsCard checks if this is a card payment method
func (pm *PaymentMethodEntity) IsCard() bool {
	return pm.Type == PaymentMethodCreditCard || pm.Type == PaymentMethodDebitCard
}

// IsExpired checks if the card is expired (for card payments)
func (pm *PaymentMethodEntity) IsExpired() bool {
	if !pm.IsCard() {
		return false
	}

	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	if pm.ExpiryYear < currentYear {
		return true
	}

	if pm.ExpiryYear == currentYear && pm.ExpiryMonth < currentMonth {
		return true
	}

	return false
}

// GetDisplayName returns a display-friendly name for the payment method
func (pm *PaymentMethodEntity) GetDisplayName() string {
	if pm.IsCard() {
		return fmt.Sprintf("%s ending in %s", pm.Brand, pm.Last4)
	}
	return string(pm.Type)
}

// MarkAsUsed updates the last used timestamp
func (pm *PaymentMethodEntity) MarkAsUsed() {
	now := time.Now()
	pm.LastUsedAt = &now
	pm.UpdatedAt = now
}

// SetAsDefault marks this payment method as default
func (pm *PaymentMethodEntity) SetAsDefault() {
	pm.IsDefault = true
	pm.UpdatedAt = time.Now()
}

// Deactivate deactivates the payment method
func (pm *PaymentMethodEntity) Deactivate() {
	pm.IsActive = false
	pm.IsDefault = false
	pm.UpdatedAt = time.Now()
}
