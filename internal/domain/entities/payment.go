package entities

import (
	"fmt"
	"math"
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
	PaymentStatusPending         PaymentStatus = "pending"
	PaymentStatusProcessing      PaymentStatus = "processing"       // Added for frontend compatibility
	PaymentStatusAwaitingPayment PaymentStatus = "awaiting_payment" // For COD orders
	PaymentStatusPaid            PaymentStatus = "paid"
	PaymentStatusCompleted       PaymentStatus = "completed"        // Alias for paid
	PaymentStatusPartiallyPaid   PaymentStatus = "partially_paid"   // For orders with partial payments
	PaymentStatusFailed          PaymentStatus = "failed"
	PaymentStatusRefunded        PaymentStatus = "refunded"
	PaymentStatusCancelled       PaymentStatus = "cancelled"
)

// Payment represents a payment transaction
type Payment struct {
	ID                uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID           uuid.UUID     `json:"order_id" gorm:"type:uuid;not null;index"` // Multiple payments per order allowed
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
	FailedAt          *time.Time    `json:"failed_at"`

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

// Validate validates the payment entity
func (p *Payment) Validate() error {
	// Validate required fields
	if p.OrderID == uuid.Nil {
		return fmt.Errorf("order_id is required")
	}

	if p.UserID == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}

	if p.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if p.Amount > 999999.99 {
		return fmt.Errorf("amount cannot exceed $999,999.99")
	}

	if p.Currency == "" {
		return fmt.Errorf("currency is required")
	}

	if len(p.Currency) != 3 {
		return fmt.Errorf("currency must be a 3-letter ISO code")
	}

	// Validate payment method
	validMethods := []PaymentMethod{
		PaymentMethodCreditCard,
		PaymentMethodDebitCard,
		PaymentMethodPayPal,
		PaymentMethodStripe,
		PaymentMethodApplePay,
		PaymentMethodGooglePay,
		PaymentMethodBankTransfer,
		PaymentMethodCash,
	}

	isValidMethod := false
	for _, method := range validMethods {
		if p.Method == method {
			isValidMethod = true
			break
		}
	}

	if !isValidMethod {
		return fmt.Errorf("invalid payment method: %s", p.Method)
	}

	// Validate payment status
	validStatuses := []PaymentStatus{
		PaymentStatusPending,
		PaymentStatusProcessing,
		PaymentStatusAwaitingPayment,
		PaymentStatusPaid,
		PaymentStatusCompleted,
		PaymentStatusPartiallyPaid,
		PaymentStatusFailed,
		PaymentStatusRefunded,
		PaymentStatusCancelled,
	}

	isValidStatus := false
	for _, status := range validStatuses {
		if p.Status == status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		return fmt.Errorf("invalid payment status: %s", p.Status)
	}

	// Validate fees
	if p.ProcessingFee < 0 {
		return fmt.Errorf("processing_fee cannot be negative")
	}

	if p.GatewayFee < 0 {
		return fmt.Errorf("gateway_fee cannot be negative")
	}

	if p.RefundAmount < 0 {
		return fmt.Errorf("refund_amount cannot be negative")
	}

	if p.RefundAmount > p.Amount {
		return fmt.Errorf("refund_amount cannot exceed payment amount")
	}

	// Auto-calculate net amount if not set
	expectedNetAmount := p.Amount - p.ProcessingFee - p.GatewayFee
	if p.NetAmount == 0 {
		p.NetAmount = expectedNetAmount
	} else {
		// Validate net amount calculation with floating point tolerance
		const epsilon = 0.01
		if math.Abs(p.NetAmount - expectedNetAmount) > epsilon {
			return fmt.Errorf("net_amount %.2f does not match calculated net_amount %.2f", p.NetAmount, expectedNetAmount)
		}
	}

	// Validate COD specific rules
	if p.Method == PaymentMethodCash {
		if p.Status != PaymentStatusAwaitingPayment && p.Status != PaymentStatusPaid && p.Status != PaymentStatusCancelled {
			return fmt.Errorf("COD payments can only have status: awaiting_payment, paid, or cancelled")
		}

		if p.Gateway != "cod" && p.Gateway != "" {
			return fmt.Errorf("COD payments should use 'cod' gateway")
		}
	}

	// Validate status transitions
	if p.Status == PaymentStatusPaid && p.ProcessedAt == nil {
		return fmt.Errorf("paid payments must have processed_at timestamp")
	}

	if p.Status == PaymentStatusRefunded && p.RefundedAt == nil {
		return fmt.Errorf("refunded payments must have refunded_at timestamp")
	}

	return nil
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

// CanTransitionTo checks if payment can transition to the given status
func (p *Payment) CanTransitionTo(newStatus PaymentStatus) bool {
	switch p.Status {
	case PaymentStatusPending:
		return newStatus == PaymentStatusProcessing || newStatus == PaymentStatusPaid ||
			   newStatus == PaymentStatusFailed || newStatus == PaymentStatusCancelled
	case PaymentStatusProcessing:
		return newStatus == PaymentStatusPaid || newStatus == PaymentStatusFailed ||
			   newStatus == PaymentStatusCancelled
	case PaymentStatusAwaitingPayment:
		return newStatus == PaymentStatusPaid || newStatus == PaymentStatusCancelled
	case PaymentStatusPaid, PaymentStatusCompleted:
		return newStatus == PaymentStatusRefunded
	case PaymentStatusPartiallyPaid:
		return newStatus == PaymentStatusPaid || newStatus == PaymentStatusRefunded
	case PaymentStatusFailed, PaymentStatusRefunded, PaymentStatusCancelled:
		return false // Terminal states
	default:
		return false
	}
}

// TransitionTo transitions payment to new status with validation
func (p *Payment) TransitionTo(newStatus PaymentStatus) error {
	if !p.CanTransitionTo(newStatus) {
		return fmt.Errorf("cannot transition payment from %s to %s", p.Status, newStatus)
	}

	p.Status = newStatus
	p.UpdatedAt = time.Now()

	// Update related fields based on status
	switch newStatus {
	case PaymentStatusProcessing:
		// No additional fields to update
	case PaymentStatusPaid, PaymentStatusCompleted:
		if p.ProcessedAt == nil {
			now := time.Now()
			p.ProcessedAt = &now
		}
		// Calculate net amount if not set
		if p.NetAmount == 0 {
			p.NetAmount = p.Amount - p.ProcessingFee - p.GatewayFee
		}
	case PaymentStatusRefunded:
		if p.RefundedAt == nil {
			now := time.Now()
			p.RefundedAt = &now
		}
	case PaymentStatusFailed, PaymentStatusCancelled:
		if p.FailedAt == nil {
			now := time.Now()
			p.FailedAt = &now
		}
	}

	return nil
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

// CalculateNetAmount calculates and updates the net amount
func (p *Payment) CalculateNetAmount() {
	p.NetAmount = p.Amount - p.ProcessingFee - p.GatewayFee
	if p.NetAmount < 0 {
		p.NetAmount = 0
	}
}

// SetFees sets the processing and gateway fees
func (p *Payment) SetFees(processingFee, gatewayFee float64) {
	p.ProcessingFee = processingFee
	p.GatewayFee = gatewayFee
	p.CalculateNetAmount()
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
	OrderID       uuid.UUID     `json:"order_id" gorm:"type:uuid;not null;index"`
	Amount        float64       `json:"amount" gorm:"not null" validate:"required,gt=0"`
	RefundFee     float64       `json:"refund_fee" gorm:"default:0"`
	NetAmount     float64       `json:"net_amount" gorm:"not null"`
	Reason        RefundReason  `json:"reason" gorm:"not null"`
	Description   string        `json:"description" gorm:"type:text"`
	Status        RefundStatus  `json:"status" gorm:"default:'pending'"`
	Type          RefundType    `json:"type" gorm:"default:'full'"`
	TransactionID string        `json:"transaction_id" gorm:"index"`
	ExternalID    string        `json:"external_id" gorm:"index"`

	// Business rules
	RequiresApproval bool      `json:"requires_approval" gorm:"default:false"`
	ApprovedBy       *uuid.UUID `json:"approved_by" gorm:"type:uuid"`
	ApprovedAt       *time.Time `json:"approved_at"`

	// Processing info
	ProcessedAt   *time.Time    `json:"processed_at"`
	ProcessedBy   *uuid.UUID    `json:"processed_by" gorm:"type:uuid"`
	FailureReason string        `json:"failure_reason" gorm:"type:text"`

	// Metadata
	Metadata      map[string]interface{} `json:"metadata" gorm:"type:jsonb"`

	// Timestamps
	CreatedAt     time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Payment *Payment `json:"payment,omitempty" gorm:"foreignKey:PaymentID"`
	Order   *Order   `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// RefundStatus represents the refund status
type RefundStatus string

const (
	RefundStatusPending         RefundStatus = "pending"
	RefundStatusAwaitingApproval RefundStatus = "awaiting_approval"
	RefundStatusApproved        RefundStatus = "approved"
	RefundStatusProcessing      RefundStatus = "processing"
	RefundStatusCompleted       RefundStatus = "completed"
	RefundStatusFailed          RefundStatus = "failed"
	RefundStatusCancelled       RefundStatus = "cancelled"
	RefundStatusRejected        RefundStatus = "rejected"
)

// RefundType represents the type of refund
type RefundType string

const (
	RefundTypeFull    RefundType = "full"
	RefundTypePartial RefundType = "partial"
)

// RefundReason represents the reason for refund
type RefundReason string

const (
	RefundReasonDefective        RefundReason = "defective"
	RefundReasonNotAsDescribed   RefundReason = "not_as_described"
	RefundReasonWrongItem        RefundReason = "wrong_item"
	RefundReasonDamaged          RefundReason = "damaged"
	RefundReasonCustomerRequest  RefundReason = "customer_request"
	RefundReasonDuplicate        RefundReason = "duplicate"
	RefundReasonFraud            RefundReason = "fraud"
	RefundReasonChargeback       RefundReason = "chargeback"
	RefundReasonOther            RefundReason = "other"
)

// Refund business constants
const (
	DefaultRefundTimeLimit = 30 * 24 * time.Hour // 30 days
	MinRefundAmount        = 0.01
	MaxRefundFeePercent    = 0.05 // 5%
	RefundFeeThreshold     = 100.0 // No fee for refunds above this amount
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
func (r *Refund) MarkAsFailed(reason string) {
	r.Status = RefundStatusFailed
	r.FailureReason = reason
	r.UpdatedAt = time.Now()
}

// MarkAsApproved marks the refund as approved
func (r *Refund) MarkAsApproved(approvedBy uuid.UUID) {
	r.Status = RefundStatusApproved
	r.ApprovedBy = &approvedBy
	now := time.Now()
	r.ApprovedAt = &now
	r.UpdatedAt = now
}

// MarkAsRejected marks the refund as rejected
func (r *Refund) MarkAsRejected(reason string) {
	r.Status = RefundStatusRejected
	r.FailureReason = reason
	r.UpdatedAt = time.Now()
}

// MarkAsProcessing marks the refund as processing
func (r *Refund) MarkAsProcessing() {
	r.Status = RefundStatusProcessing
	r.UpdatedAt = time.Now()
}

// CalculateRefundFee calculates the refund fee based on business rules
func (r *Refund) CalculateRefundFee() float64 {
	// No fee for high-value refunds
	if r.Amount >= RefundFeeThreshold {
		return 0
	}

	// Calculate percentage-based fee
	fee := r.Amount * MaxRefundFeePercent

	// Ensure minimum fee
	if fee < 1.0 && r.Amount > 10.0 {
		fee = 1.0
	}

	return fee
}

// SetRefundFee sets the refund fee and calculates net amount
func (r *Refund) SetRefundFee(fee float64) {
	r.RefundFee = fee
	r.NetAmount = r.Amount - fee
}

// IsEligibleForAutoApproval checks if refund can be auto-approved
func (r *Refund) IsEligibleForAutoApproval() bool {
	// Auto-approve small amounts and certain reasons
	autoApprovalReasons := []RefundReason{
		RefundReasonDefective,
		RefundReasonDamaged,
		RefundReasonWrongItem,
	}

	for _, reason := range autoApprovalReasons {
		if r.Reason == reason && r.Amount <= 50.0 {
			return true
		}
	}

	return false
}

// CanBeProcessed checks if refund can be processed
func (r *Refund) CanBeProcessed() bool {
	return r.Status == RefundStatusApproved ||
		   (r.Status == RefundStatusPending && r.IsEligibleForAutoApproval())
}

// GetRemainingRefundAmount returns the remaining amount that can be refunded
func (p *Payment) GetRemainingRefundAmount() float64 {
	return p.Amount - p.RefundAmount
}

// CanBeRefunded checks if the payment can be refunded
func (p *Payment) CanBeRefunded() bool {
	return p.Status == PaymentStatusPaid && p.RefundAmount < p.Amount
}

// CanBeRefundedWithTimeLimit checks if payment can be refunded within time limit
func (p *Payment) CanBeRefundedWithTimeLimit() bool {
	if !p.CanBeRefunded() {
		return false
	}

	// Check time limit
	timeLimit := time.Now().Add(-DefaultRefundTimeLimit)
	return p.CreatedAt.After(timeLimit)
}

// ValidateRefundAmount validates if the refund amount is valid
func (p *Payment) ValidateRefundAmount(amount float64) error {
	if amount <= 0 {
		return ErrInvalidRefundAmount
	}

	if amount < MinRefundAmount {
		return fmt.Errorf("refund amount must be at least %.2f", MinRefundAmount)
	}

	if p.RefundAmount + amount > p.Amount {
		return ErrRefundAmountExceedsPayment
	}

	return nil
}

// GetMaxRefundableAmount returns the maximum amount that can be refunded
func (p *Payment) GetMaxRefundableAmount() float64 {
	return p.Amount - p.RefundAmount
}

// HasPendingRefunds checks if payment has pending refunds
func (p *Payment) HasPendingRefunds() bool {
	// This would need to be implemented with repository query
	// For now, return false - will be implemented in use case
	return false
}

// CalculateRefundImpact calculates the impact of a refund on payment status
func (p *Payment) CalculateRefundImpact(refundAmount float64) PaymentStatus {
	totalRefunded := p.RefundAmount + refundAmount

	if totalRefunded >= p.Amount {
		return PaymentStatusRefunded
	}

	return p.Status // No change
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

// Validate validates payment method data
func (pm *PaymentMethodEntity) Validate() error {
	// Validate required fields
	if pm.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}

	if pm.Type == "" {
		return fmt.Errorf("payment method type is required")
	}

	if pm.GatewayToken == "" {
		return fmt.Errorf("gateway token is required")
	}

	// Validate card-specific fields
	if pm.IsCard() {
		if pm.Last4 == "" {
			return fmt.Errorf("last 4 digits are required for card payments")
		}

		if len(pm.Last4) != 4 {
			return fmt.Errorf("last 4 digits must be exactly 4 characters")
		}

		if pm.Brand == "" {
			return fmt.Errorf("card brand is required for card payments")
		}

		// Validate expiry date
		if pm.ExpiryMonth < 1 || pm.ExpiryMonth > 12 {
			return fmt.Errorf("expiry month must be between 1 and 12")
		}

		currentYear := time.Now().Year()
		if pm.ExpiryYear < currentYear || pm.ExpiryYear > currentYear+20 {
			return fmt.Errorf("expiry year must be between %d and %d", currentYear, currentYear+20)
		}

		// Check if card is expired
		if pm.IsExpired() {
			return fmt.Errorf("card is expired")
		}
	}

	// Validate gateway
	validGateways := map[string]bool{
		"stripe": true,
		"paypal": true,
		"square": true,
		"cod":    true,
	}

	if !validGateways[pm.Gateway] {
		return fmt.Errorf("invalid gateway: %s", pm.Gateway)
	}

	// Validate fingerprint
	if pm.Fingerprint == "" {
		return fmt.Errorf("fingerprint is required for security")
	}

	return nil
}

// MaskSensitiveData masks sensitive data for API responses
func (pm *PaymentMethodEntity) MaskSensitiveData() {
	// Clear sensitive fields
	pm.GatewayToken = "***"
	pm.GatewayCustomerID = "***"
	pm.BillingAddress = "***"
	pm.Fingerprint = "***"
}
