package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/domain/services"
	"ecom-golang-clean-architecture/internal/infrastructure/database"
	"ecom-golang-clean-architecture/internal/infrastructure/payment"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentGatewayService defines the interface for payment gateway services
type PaymentGatewayService interface {
	ProcessPayment(ctx context.Context, req payment.PaymentGatewayRequest) (*payment.PaymentGatewayResponse, error)
	ProcessRefund(ctx context.Context, req payment.RefundGatewayRequest) (*payment.RefundGatewayResponse, error)
	CreateCheckoutSession(ctx context.Context, req payment.CheckoutSessionRequest) (*payment.CheckoutSessionResponse, error)
}

// Type aliases for convenience
type PaymentGatewayRequest = payment.PaymentGatewayRequest
type PaymentGatewayResponse = payment.PaymentGatewayResponse
type RefundGatewayRequest = payment.RefundGatewayRequest
type RefundGatewayResponse = payment.RefundGatewayResponse
type CheckoutSessionRequest = payment.CheckoutSessionRequest
type CheckoutSessionResponse = payment.CheckoutSessionResponse

// PaymentUseCase defines payment use cases
type PaymentUseCase interface {
	// Payment processing
	ProcessPayment(ctx context.Context, req ProcessPaymentRequest) (*PaymentResponse, error)
	GetPayment(ctx context.Context, id uuid.UUID) (*PaymentResponse, error)
	GetOrderPayments(ctx context.Context, orderID uuid.UUID) ([]*PaymentResponse, error)
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status entities.PaymentStatus, transactionID string) (*PaymentResponse, error)

	// Refunds
	ProcessRefund(ctx context.Context, req ProcessRefundRequest) (*RefundResponse, error)
	GetRefunds(ctx context.Context, paymentID uuid.UUID) ([]*RefundResponse, error)
	ApproveRefund(ctx context.Context, refundID uuid.UUID, approvedBy uuid.UUID) (*RefundResponse, error)
	RejectRefund(ctx context.Context, refundID uuid.UUID, reason string) error
	GetPendingRefunds(ctx context.Context, limit, offset int) ([]*RefundResponse, error)

	// Payment methods
	SavePaymentMethod(ctx context.Context, req SavePaymentMethodRequest) (*PaymentMethodResponse, error)
	GetUserPaymentMethods(ctx context.Context, userID uuid.UUID) ([]*PaymentMethodResponse, error)
	DeletePaymentMethod(ctx context.Context, id uuid.UUID) error
	SetDefaultPaymentMethod(ctx context.Context, userID, methodID uuid.UUID) error

	// Webhooks
	HandleWebhook(ctx context.Context, provider string, payload []byte, signature string) error

	// Payment confirmation (fallback method)
	ConfirmPaymentSuccess(ctx context.Context, orderID, userID uuid.UUID, sessionID string) error

	// Reports
	GetPaymentReport(ctx context.Context, req PaymentReportRequest) (*PaymentReportResponse, error)

	// Stripe Checkout
	CreateCheckoutSession(ctx context.Context, req CreateCheckoutSessionRequest) (*CreateCheckoutSessionResponse, error)
}

type paymentUseCase struct {
	paymentRepo             repositories.PaymentRepository
	paymentMethodRepo       repositories.PaymentMethodRepository
	orderRepo               repositories.OrderRepository
	userRepo                repositories.UserRepository
	stripeService           PaymentGatewayService
	paypalService           PaymentGatewayService
	notificationUseCase     NotificationUseCase
	stockReservationService services.StockReservationService
	orderEventService       services.OrderEventService
	txManager               *database.TransactionManager
}

// NewPaymentUseCase creates a new payment use case
func NewPaymentUseCase(
	paymentRepo repositories.PaymentRepository,
	paymentMethodRepo repositories.PaymentMethodRepository,
	orderRepo repositories.OrderRepository,
	userRepo repositories.UserRepository,
	stripeService PaymentGatewayService,
	paypalService PaymentGatewayService,
	notificationUseCase NotificationUseCase,
	stockReservationService services.StockReservationService,
	orderEventService services.OrderEventService,
	txManager *database.TransactionManager,
) PaymentUseCase {
	return &paymentUseCase{
		paymentRepo:             paymentRepo,
		paymentMethodRepo:       paymentMethodRepo,
		orderRepo:               orderRepo,
		userRepo:                userRepo,
		stripeService:           stripeService,
		paypalService:           paypalService,
		notificationUseCase:     notificationUseCase,
		stockReservationService: stockReservationService,
		orderEventService:       orderEventService,
		txManager:               txManager,
	}
}

// Request/Response types
type ProcessPaymentRequest struct {
	OrderID         uuid.UUID              `json:"order_id" validate:"required"`
	Amount          float64                `json:"amount" validate:"required,gt=0"`
	Currency        string                 `json:"currency" validate:"required"`
	Method          entities.PaymentMethod `json:"method" validate:"required"`
	PaymentToken    string                 `json:"payment_token,omitempty"`
	PaymentMethodID *uuid.UUID             `json:"payment_method_id,omitempty"`
	BillingAddress  *BillingAddressRequest `json:"billing_address,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

type ProcessRefundRequest struct {
	PaymentID         uuid.UUID                    `json:"payment_id" validate:"required"`
	OrderID           uuid.UUID                    `json:"order_id" validate:"required"`
	Amount            float64                      `json:"amount" validate:"required,gt=0"`
	Reason            entities.RefundReason        `json:"reason" validate:"required"`
	Description       string                       `json:"description,omitempty"`
	Type              entities.RefundType          `json:"type" validate:"required"`
	ForceApproval     bool                         `json:"force_approval,omitempty"`
	ProcessedBy       *uuid.UUID                   `json:"processed_by,omitempty"`
	Metadata          map[string]interface{}       `json:"metadata,omitempty"`
}

type SavePaymentMethodRequest struct {
	UserID            uuid.UUID              `json:"user_id" validate:"required"`
	Type              entities.PaymentMethod `json:"type" validate:"required"`
	Token             string                 `json:"token" validate:"required"`

	// Card information (for card payments)
	Last4             string                 `json:"last4"`
	Brand             string                 `json:"brand"`
	ExpiryMonth       int                    `json:"expiry_month"`
	ExpiryYear        int                    `json:"expiry_year"`

	// Gateway information
	Gateway           string                 `json:"gateway"`
	GatewayCustomerID string                 `json:"gateway_customer_id"`

	// Billing information
	BillingName       string                 `json:"billing_name"`
	BillingEmail      string                 `json:"billing_email"`
	BillingAddress    string                 `json:"billing_address"`

	// Preferences
	IsDefault         bool                   `json:"is_default"`

	// Security
	Fingerprint       string                 `json:"fingerprint"`

	// Metadata
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	MetadataJSON      string                 `json:"-"` // Internal use
	Notes             string                 `json:"notes"`
}

type BillingAddressRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Address1  string `json:"address1" validate:"required"`
	Address2  string `json:"address2,omitempty"`
	City      string `json:"city" validate:"required"`
	State     string `json:"state" validate:"required"`
	ZipCode   string `json:"zip_code" validate:"required"`
	Country   string `json:"country" validate:"required"`
}

type PaymentReportRequest struct {
	ReportType string                  `json:"report_type"`
	DateFrom   *time.Time              `json:"date_from,omitempty"`
	DateTo     *time.Time              `json:"date_to,omitempty"`
	Status     *entities.PaymentStatus `json:"status,omitempty"`
	Method     *entities.PaymentMethod `json:"method,omitempty"`
	GroupBy    string                  `json:"group_by,omitempty" validate:"omitempty,oneof=day week month method status"`
	Format     string                  `json:"format,omitempty" validate:"omitempty,oneof=json csv excel"`
}

// CreateCheckoutSessionRequest represents a request to create a checkout session
type CreateCheckoutSessionRequest struct {
	OrderID     uuid.UUID              `json:"order_id" validate:"required"`
	Amount      float64                `json:"amount" validate:"required,gt=0"`
	Currency    string                 `json:"currency" validate:"required,len=3"`
	Description string                 `json:"description,omitempty"`
	SuccessURL  string                 `json:"success_url" validate:"required,url"`
	CancelURL   string                 `json:"cancel_url" validate:"required,url"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// CreateCheckoutSessionResponse represents a response from creating a checkout session
type CreateCheckoutSessionResponse struct {
	Success    bool   `json:"success"`
	SessionID  string `json:"session_id,omitempty"`
	SessionURL string `json:"session_url,omitempty"`
	Message    string `json:"message"`
}

// Response types
type PaymentResponse struct {
	ID              uuid.UUID              `json:"id"`
	OrderID         uuid.UUID              `json:"order_id"`
	Amount          float64                `json:"amount"`
	Currency        string                 `json:"currency"`
	Method          entities.PaymentMethod `json:"method"`
	Status          entities.PaymentStatus `json:"status"`
	TransactionID   string                 `json:"transaction_id"`
	ExternalID      string                 `json:"external_id"`
	FailureReason   string                 `json:"failure_reason,omitempty"`
	ProcessedAt     *time.Time             `json:"processed_at"`
	RefundedAt      *time.Time             `json:"refunded_at"`
	RefundAmount    float64                `json:"refund_amount"`
	CanBeRefunded   bool                   `json:"can_be_refunded"`
	RemainingRefund float64                `json:"remaining_refund"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type RefundResponse struct {
	ID               uuid.UUID                  `json:"id"`
	PaymentID        uuid.UUID                  `json:"payment_id"`
	OrderID          uuid.UUID                  `json:"order_id"`
	Amount           float64                    `json:"amount"`
	RefundFee        float64                    `json:"refund_fee"`
	NetAmount        float64                    `json:"net_amount"`
	Reason           entities.RefundReason      `json:"reason"`
	Description      string                     `json:"description"`
	Status           entities.RefundStatus      `json:"status"`
	Type             entities.RefundType        `json:"type"`
	TransactionID    string                     `json:"transaction_id"`
	RequiresApproval bool                       `json:"requires_approval"`
	ApprovedBy       *uuid.UUID                 `json:"approved_by"`
	ApprovedAt       *time.Time                 `json:"approved_at"`
	ProcessedAt      *time.Time                 `json:"processed_at"`
	ProcessedBy      *uuid.UUID                 `json:"processed_by"`
	FailureReason    string                     `json:"failure_reason"`
	Metadata         map[string]interface{}     `json:"metadata"`
	CreatedAt        time.Time                  `json:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at"`
}

type PaymentMethodResponse struct {
	ID            uuid.UUID              `json:"id"`
	UserID        uuid.UUID              `json:"user_id"`
	Type          entities.PaymentMethod `json:"type"`
	Last4         string                 `json:"last4,omitempty"`
	Brand         string                 `json:"brand,omitempty"`
	ExpiryMonth   int                    `json:"expiry_month,omitempty"`
	ExpiryYear    int                    `json:"expiry_year,omitempty"`
	Gateway       string                 `json:"gateway,omitempty"`
	BillingName   string                 `json:"billing_name,omitempty"`
	BillingEmail  string                 `json:"billing_email,omitempty"`
	IsDefault     bool                   `json:"is_default"`
	IsActive      bool                   `json:"is_active"`
	IsExpired     bool                   `json:"is_expired"`
	DisplayName   string                 `json:"display_name"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	LastUsedAt    *time.Time             `json:"last_used_at,omitempty"`
}

type BillingAddressResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Address1  string `json:"address1"`
	Address2  string `json:"address2"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zip_code"`
	Country   string `json:"country"`
}

type WebhookEvent struct {
	Type          string                 `json:"type"`
	PaymentID     string                 `json:"payment_id"`
	TransactionID string                 `json:"transaction_id"`
	Status        entities.PaymentStatus `json:"status"`
	Amount        float64                `json:"amount"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

type PaymentReportResponse struct {
	ReportType  string                   `json:"report_type"`
	ReportID    uuid.UUID                `json:"report_id"`
	GeneratedAt time.Time                `json:"generated_at"`
	Format      string                   `json:"format"`
	DownloadURL string                   `json:"download_url,omitempty"`
	Summary     *PaymentReportSummary    `json:"summary"`
	Data        []map[string]interface{} `json:"data,omitempty"`
}

type PaymentReportSummary struct {
	TotalAmount        float64 `json:"total_amount"`
	TotalPayments      int64   `json:"total_payments"`
	SuccessfulPayments int64   `json:"successful_payments"`
	FailedPayments     int64   `json:"failed_payments"`
	RefundAmount       float64 `json:"refund_amount"`
	SuccessRate        float64 `json:"success_rate"`
}

// ProcessPayment processes a payment for an order
func (uc *paymentUseCase) ProcessPayment(ctx context.Context, req ProcessPaymentRequest) (*PaymentResponse, error) {
	// Get order details
	order, err := uc.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, entities.ErrOrderNotFound
	}

	// Validate payment amount (allow partial payments)
	if req.Amount <= 0 {
		return nil, fmt.Errorf("payment amount must be greater than 0")
	}

	if req.Amount > 999999.99 {
		return nil, fmt.Errorf("payment amount cannot exceed $999,999.99")
	}

	// Get existing payments for this order to calculate remaining amount
	existingPayments, err := uc.paymentRepo.GetAllByOrderID(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing payments: %w", err)
	}

	// Calculate total paid amount from successful payments
	totalPaid := 0.0
	if existingPayments != nil {
		for _, p := range existingPayments {
			if p.Status == entities.PaymentStatusPaid || p.Status == entities.PaymentStatusCompleted {
				totalPaid += p.Amount
			}
		}
	}

	// Calculate remaining amount with floating point tolerance
	remainingAmount := order.Total - totalPaid
	const epsilon = 0.01
	if req.Amount > remainingAmount + epsilon {
		return nil, fmt.Errorf("payment amount %.2f exceeds remaining balance %.2f", req.Amount, remainingAmount)
	}

	// Check for duplicate payments (same amount within short time window)
	if existingPayments != nil {
		for _, p := range existingPayments {
			if p.Amount == req.Amount && p.Status == entities.PaymentStatusPending {
				timeDiff := time.Since(p.CreatedAt)
				if timeDiff < 5*time.Minute {
					return nil, fmt.Errorf("duplicate payment detected: same amount %.2f within 5 minutes", req.Amount)
				}
			}
		}
	}

	// Validate currency
	if req.Currency == "" {
		req.Currency = "USD" // Default currency
	}

	if len(req.Currency) != 3 {
		return nil, fmt.Errorf("currency must be a 3-letter ISO code")
	}

	// Validate order status
	if order.Status == entities.OrderStatusCancelled {
		return nil, fmt.Errorf("cannot process payment for cancelled order")
	}

	if order.Status == entities.OrderStatusDelivered {
		return nil, fmt.Errorf("cannot process payment for delivered order")
	}

	if order.Status == entities.OrderStatusRefunded {
		return nil, fmt.Errorf("cannot process payment for refunded order")
	}

	// Determine initial payment status based on method
	var initialStatus entities.PaymentStatus
	if req.Method == entities.PaymentMethodCash {
		initialStatus = entities.PaymentStatusAwaitingPayment
	} else {
		initialStatus = entities.PaymentStatusPending
	}

	// Create payment record
	payment := &entities.Payment{
		ID:        uuid.New(),
		OrderID:   req.OrderID,
		UserID:    order.UserID,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Method:    req.Method,
		Status:    initialStatus,
		Gateway:   uc.getGatewayForMethod(req.Method),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Validate payment entity
	if err := payment.Validate(); err != nil {
		return nil, fmt.Errorf("payment validation failed: %w", err)
	}

	// Save payment record
	if err := uc.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	// Process payment through gateway
	// Convert metadata to string map
	metadata := make(map[string]string)
	for k, v := range req.Metadata {
		if str, ok := v.(string); ok {
			metadata[k] = str
		} else {
			metadata[k] = fmt.Sprintf("%v", v)
		}
	}

	gatewayReq := PaymentGatewayRequest{
		Amount:       req.Amount,
		Currency:     req.Currency,
		PaymentToken: req.PaymentToken,
		Description:  fmt.Sprintf("Payment for order %s", order.OrderNumber),
		Metadata:     metadata,
	}

	if req.PaymentMethodID != nil {
		gatewayReq.PaymentMethodID = req.PaymentMethodID.String()
	}

	var gatewayResp *PaymentGatewayResponse
	switch req.Method {
	case entities.PaymentMethodStripe:
		if uc.stripeService == nil {
			return nil, fmt.Errorf("stripe service not configured")
		}
		gatewayResp, err = uc.stripeService.ProcessPayment(ctx, gatewayReq)
	case entities.PaymentMethodPayPal:
		if uc.paypalService == nil {
			return nil, fmt.Errorf("paypal service not configured")
		}
		gatewayResp, err = uc.paypalService.ProcessPayment(ctx, gatewayReq)
	case entities.PaymentMethodCash:
		// COD payments don't need gateway processing
		// Just mark as awaiting payment (will be marked as paid when delivered)
		gatewayResp = &PaymentGatewayResponse{
			Success:       true,
			TransactionID: fmt.Sprintf("COD-%s", payment.ID.String()[:8]),
			ExternalID:    "",
			Message:       "COD payment created successfully",
		}
	case entities.PaymentMethodCreditCard, entities.PaymentMethodDebitCard:
		// Default to Stripe for credit/debit cards
		if uc.stripeService == nil {
			return nil, fmt.Errorf("stripe service not configured")
		}
		gatewayResp, err = uc.stripeService.ProcessPayment(ctx, gatewayReq)
	default:
		return nil, fmt.Errorf("unsupported payment method: %s", req.Method)
	}

	if err != nil {
		// Categorize error types for better handling
		errorMessage := err.Error()

		// Check if it's a temporary error that can be retried
		isRetryable := isRetryableError(err)

		if isRetryable {
			// For retryable errors, keep status as pending for potential retry
			payment.Status = entities.PaymentStatusPending
			payment.FailureReason = fmt.Sprintf("Retryable error: %s", errorMessage)
		} else {
			// For permanent errors, mark as failed
			payment.MarkAsFailed(errorMessage)
		}

		// Store gateway response for debugging
		if gatewayResp != nil {
			payment.GatewayResponse = gatewayResp.Message
		}

		if updateErr := uc.paymentRepo.Update(ctx, payment); updateErr != nil {
			return nil, fmt.Errorf("failed to update payment after gateway error: %w", updateErr)
		}

		return nil, fmt.Errorf("payment processing failed: %w", err)
	}

	// Update payment with gateway response
	if gatewayResp.Success {
		// For COD payments, keep status as awaiting_payment
		if req.Method == entities.PaymentMethodCash {
			payment.TransactionID = gatewayResp.TransactionID
			payment.ExternalID = gatewayResp.ExternalID
			// Status remains awaiting_payment for COD
		} else {
			// For online payments, mark as processed
			payment.MarkAsProcessed(gatewayResp.TransactionID)
			payment.ExternalID = gatewayResp.ExternalID
		}

		// Reload order with payments to sync payment status
		order, err = uc.orderRepo.GetByID(ctx, order.ID)
		if err != nil {
			return nil, err
		}

		// Sync payment status based on all payments
		order.SyncPaymentStatus(payment.Status)
		if err := uc.orderRepo.Update(ctx, order); err != nil {
			return nil, err
		}

		// Send payment confirmation notification
		if uc.notificationUseCase != nil {
			if err := uc.notificationUseCase.NotifyPaymentReceived(ctx, payment.ID); err != nil {
				return nil, err
			}
		}
	} else {
		payment.MarkAsFailed(gatewayResp.Message)
	}

	// Save updated payment
	if err := uc.paymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	return uc.toPaymentResponse(payment), nil
}

// GetPayment gets a payment by ID
func (uc *paymentUseCase) GetPayment(ctx context.Context, id uuid.UUID) (*PaymentResponse, error) {
	payment, err := uc.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrPaymentNotFound
	}

	return uc.toPaymentResponse(payment), nil
}

// GetOrderPayments gets all payments for an order
func (uc *paymentUseCase) GetOrderPayments(ctx context.Context, orderID uuid.UUID) ([]*PaymentResponse, error) {
	payments, err := uc.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if payments == nil {
		return []*PaymentResponse{}, nil
	}

	// Convert single payment to slice for backward compatibility
	// Note: This suggests the repository method should return []Payment instead of *Payment
	return []*PaymentResponse{uc.toPaymentResponse(payments)}, nil
}

// UpdatePaymentStatus updates payment status and syncs with order
func (uc *paymentUseCase) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status entities.PaymentStatus, transactionID string) (*PaymentResponse, error) {
	// Execute in transaction to ensure consistency
	result, err := uc.txManager.WithTransactionResult(ctx, func(tx *gorm.DB) (interface{}, error) {
		return uc.updatePaymentStatusInTransaction(ctx, id, status, transactionID)
	})
	if err != nil {
		return nil, err
	}
	return result.(*PaymentResponse), nil
}

// updatePaymentStatusInTransaction updates payment status within a transaction
func (uc *paymentUseCase) updatePaymentStatusInTransaction(ctx context.Context, id uuid.UUID, status entities.PaymentStatus, transactionID string) (*PaymentResponse, error) {
	payment, err := uc.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	payment.Status = status
	if transactionID != "" {
		payment.TransactionID = transactionID
	}

	if status == entities.PaymentStatusPaid {
		now := time.Now()
		payment.ProcessedAt = &now
	}

	payment.UpdatedAt = time.Now()

	if err := uc.paymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	// Sync order payment status
	order, err := uc.orderRepo.GetByID(ctx, payment.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order for payment sync: %v", err)
	}

	// Reload order with payments to ensure we have latest payment data
	order, err = uc.orderRepo.GetByID(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload order with payments: %v", err)
	}

	// Sync payment status based on all payments
	oldPaymentStatus := order.PaymentStatus
	order.SyncPaymentStatus(status)

	// Update order status if payment is completed
	if status == entities.PaymentStatusPaid && order.Status == entities.OrderStatusPending && order.IsFullyPaid() {
		order.Status = entities.OrderStatusConfirmed

		// Confirm stock reservations (convert to actual stock reduction)
		if err := uc.stockReservationService.ConfirmReservations(ctx, order.ID); err != nil {
			fmt.Printf("❌ Failed to confirm stock reservations: %v\n", err)
			return nil, fmt.Errorf("failed to confirm stock reservations: %v", err)
		}
		fmt.Printf("✅ Stock reservations confirmed and converted to actual stock reduction\n")

		// Release reservation flags since stock is now actually reduced
		order.ReleaseReservation()
	}

	order.UpdatedAt = time.Now()

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order payment status: %v", err)
	}

	// Log the sync for debugging
	fmt.Printf("✅ Payment status updated: Payment=%s->%s, Order PaymentStatus=%s->%s\n",
		payment.ID, status, oldPaymentStatus, order.PaymentStatus)

	return uc.toPaymentResponse(payment), nil
}

// ProcessRefund processes a refund for a payment
func (uc *paymentUseCase) ProcessRefund(ctx context.Context, req ProcessRefundRequest) (*RefundResponse, error) {
	// Get payment details
	payment, err := uc.paymentRepo.GetByID(ctx, req.PaymentID)
	if err != nil {
		return nil, entities.ErrPaymentNotFound
	}

	// Comprehensive refund validation
	if err := uc.validateRefundRequest(ctx, payment, req); err != nil {
		return nil, err
	}

	// Create refund entity
	refund := &entities.Refund{
		ID:          uuid.New(),
		PaymentID:   req.PaymentID,
		OrderID:     req.OrderID,
		Amount:      req.Amount,
		Reason:      req.Reason,
		Description: req.Description,
		Type:        req.Type,
		Status:      entities.RefundStatusPending,
		Metadata:    req.Metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Calculate refund fee
	refundFee := refund.CalculateRefundFee()
	refund.SetRefundFee(refundFee)

	// Check if approval is required
	requiresApproval := req.ForceApproval || !refund.IsEligibleForAutoApproval()
	refund.RequiresApproval = requiresApproval

	if requiresApproval {
		refund.Status = entities.RefundStatusAwaitingApproval

		// Save refund for approval
		if err := uc.paymentRepo.CreateRefund(ctx, refund); err != nil {
			return nil, err
		}

		return uc.mapRefundToResponse(refund), nil
	}

	// Auto-approve and process
	if req.ProcessedBy != nil {
		refund.MarkAsApproved(*req.ProcessedBy)
	}

	return uc.processApprovedRefund(ctx, payment, refund)
}

// validateRefundRequest validates a refund request
func (uc *paymentUseCase) validateRefundRequest(ctx context.Context, payment *entities.Payment, req ProcessRefundRequest) error {
	// Basic payment validation
	if !payment.CanBeRefunded() {
		return fmt.Errorf("payment cannot be refunded")
	}

	// Time limit validation
	if !payment.CanBeRefundedWithTimeLimit() {
		return entities.ErrRefundTimeExpired
	}

	// Amount validation
	if err := payment.ValidateRefundAmount(req.Amount); err != nil {
		return err
	}

	// Reason validation
	validReasons := []entities.RefundReason{
		entities.RefundReasonDefective,
		entities.RefundReasonNotAsDescribed,
		entities.RefundReasonWrongItem,
		entities.RefundReasonDamaged,
		entities.RefundReasonCustomerRequest,
		entities.RefundReasonDuplicate,
		entities.RefundReasonFraud,
		entities.RefundReasonChargeback,
		entities.RefundReasonOther,
	}

	isValidReason := false
	for _, validReason := range validReasons {
		if req.Reason == validReason {
			isValidReason = true
			break
		}
	}

	if !isValidReason {
		return entities.ErrInvalidRefundReason
	}

	// Check for pending refunds
	existingRefunds, err := uc.paymentRepo.GetRefundsByPaymentID(ctx, req.PaymentID)
	if err != nil {
		return err
	}

	for _, refund := range existingRefunds {
		if refund.Status == entities.RefundStatusPending ||
		   refund.Status == entities.RefundStatusAwaitingApproval ||
		   refund.Status == entities.RefundStatusProcessing {
			return fmt.Errorf("payment has pending refunds")
		}
	}

	return nil
}

// processApprovedRefund processes an approved refund
func (uc *paymentUseCase) processApprovedRefund(ctx context.Context, payment *entities.Payment, refund *entities.Refund) (*RefundResponse, error) {
	refund.MarkAsProcessing()

	// Process refund through gateway
	refundReq := RefundGatewayRequest{
		TransactionID: payment.TransactionID,
		Amount:        refund.Amount,
		Reason:        string(refund.Reason),
	}

	switch payment.Method {
	case entities.PaymentMethodStripe:
		gatewayResp, err := uc.stripeService.ProcessRefund(ctx, refundReq)
		if err != nil {
			refund.MarkAsFailed(fmt.Sprintf("Stripe gateway error: %v", err))
			uc.paymentRepo.UpdateRefund(ctx, refund)
			return nil, fmt.Errorf("gateway refund failed: %v", err)
		}

		if gatewayResp.Success {
			refund.MarkAsCompleted(gatewayResp.RefundID)
		} else {
			refund.MarkAsFailed(gatewayResp.Message)
			uc.paymentRepo.UpdateRefund(ctx, refund)
			return nil, fmt.Errorf("refund failed: %s", gatewayResp.Message)
		}

	case entities.PaymentMethodPayPal:
		gatewayResp, err := uc.paypalService.ProcessRefund(ctx, refundReq)
		if err != nil {
			refund.MarkAsFailed(fmt.Sprintf("PayPal gateway error: %v", err))
			uc.paymentRepo.UpdateRefund(ctx, refund)
			return nil, fmt.Errorf("gateway refund failed: %v", err)
		}

		if gatewayResp.Success {
			refund.MarkAsCompleted(gatewayResp.RefundID)
		} else {
			refund.MarkAsFailed(gatewayResp.Message)
			uc.paymentRepo.UpdateRefund(ctx, refund)
			return nil, fmt.Errorf("refund failed: %s", gatewayResp.Message)
		}

	default:
		refund.MarkAsFailed(fmt.Sprintf("Unsupported gateway: %s", payment.Method))
		uc.paymentRepo.UpdateRefund(ctx, refund)
		return nil, fmt.Errorf("unsupported payment method for refund: %s", payment.Method)
	}

	// Update payment with refund
	if err := payment.AddRefund(refund.Amount); err != nil {
		refund.MarkAsFailed(fmt.Sprintf("Failed to update payment: %v", err))
		uc.paymentRepo.UpdateRefund(ctx, refund)
		return nil, err
	}

	// Save updated payment
	if err := uc.paymentRepo.Update(ctx, payment); err != nil {
		refund.MarkAsFailed(fmt.Sprintf("Failed to save payment: %v", err))
		uc.paymentRepo.UpdateRefund(ctx, refund)
		return nil, err
	}

	// Save completed refund
	if err := uc.paymentRepo.UpdateRefund(ctx, refund); err != nil {
		return nil, err
	}

	return uc.mapRefundToResponse(refund), nil
}

// mapRefundToResponse maps a refund entity to response
func (uc *paymentUseCase) mapRefundToResponse(refund *entities.Refund) *RefundResponse {
	return &RefundResponse{
		ID:               refund.ID,
		PaymentID:        refund.PaymentID,
		OrderID:          refund.OrderID,
		Amount:           refund.Amount,
		RefundFee:        refund.RefundFee,
		NetAmount:        refund.NetAmount,
		Reason:           refund.Reason,
		Description:      refund.Description,
		Status:           refund.Status,
		Type:             refund.Type,
		TransactionID:    refund.TransactionID,
		RequiresApproval: refund.RequiresApproval,
		ApprovedBy:       refund.ApprovedBy,
		ApprovedAt:       refund.ApprovedAt,
		ProcessedAt:      refund.ProcessedAt,
		ProcessedBy:      refund.ProcessedBy,
		FailureReason:    refund.FailureReason,
		Metadata:         refund.Metadata,
		CreatedAt:        refund.CreatedAt,
		UpdatedAt:        refund.UpdatedAt,
	}
}

// ApproveRefund approves a pending refund
func (uc *paymentUseCase) ApproveRefund(ctx context.Context, refundID uuid.UUID, approvedBy uuid.UUID) (*RefundResponse, error) {
	// Get refund
	refund, err := uc.paymentRepo.GetRefund(ctx, refundID)
	if err != nil {
		return nil, entities.ErrRefundNotFound
	}

	// Validate refund status
	if refund.Status != entities.RefundStatusAwaitingApproval {
		return nil, fmt.Errorf("refund is not awaiting approval")
	}

	// Get payment for processing
	payment, err := uc.paymentRepo.GetByID(ctx, refund.PaymentID)
	if err != nil {
		return nil, entities.ErrPaymentNotFound
	}

	// Mark as approved
	refund.MarkAsApproved(approvedBy)

	// Process the approved refund
	return uc.processApprovedRefund(ctx, payment, refund)
}

// RejectRefund rejects a pending refund
func (uc *paymentUseCase) RejectRefund(ctx context.Context, refundID uuid.UUID, reason string) error {
	// Get refund
	refund, err := uc.paymentRepo.GetRefund(ctx, refundID)
	if err != nil {
		return entities.ErrRefundNotFound
	}

	// Validate refund status
	if refund.Status != entities.RefundStatusAwaitingApproval {
		return fmt.Errorf("refund is not awaiting approval")
	}

	// Mark as rejected
	refund.MarkAsRejected(reason)

	// Save updated refund
	return uc.paymentRepo.UpdateRefund(ctx, refund)
}

// GetPendingRefunds retrieves refunds awaiting approval
func (uc *paymentUseCase) GetPendingRefunds(ctx context.Context, limit, offset int) ([]*RefundResponse, error) {
	refunds, err := uc.paymentRepo.GetPendingRefunds(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*RefundResponse, len(refunds))
	for i, refund := range refunds {
		responses[i] = uc.mapRefundToResponse(refund)
	}

	return responses, nil
}

// HandleWebhook handles payment gateway webhooks
func (uc *paymentUseCase) HandleWebhook(ctx context.Context, provider string, payload []byte, signature string) error {
	switch provider {
	case "stripe":
		return uc.handleStripeWebhook(ctx, payload, signature)
	case "paypal":
		return uc.handlePayPalWebhook(ctx, payload, signature)
	default:
		return fmt.Errorf("unsupported payment provider: %s", provider)
	}
}

// handleStripeWebhook processes Stripe webhook events
func (uc *paymentUseCase) handleStripeWebhook(ctx context.Context, payload []byte, signature string) error {
	// Cast to StripeService to access HandleWebhook method
	stripeService, ok := uc.stripeService.(*payment.StripeService)
	if !ok {
		return fmt.Errorf("stripe service not properly configured")
	}

	// Parse webhook event
	webhookEvent, err := stripeService.HandleWebhook(ctx, payload, signature)
	if err != nil {
		return fmt.Errorf("failed to parse stripe webhook: %v", err)
	}

	// Process different event types
	switch webhookEvent.Type {
	case "checkout.session.completed":
		return uc.handleCheckoutSessionCompleted(ctx, webhookEvent)
	case "payment_intent.succeeded":
		return uc.handlePaymentIntentSucceeded(ctx, webhookEvent)
	case "payment_intent.payment_failed":
		return uc.handlePaymentIntentFailed(ctx, webhookEvent)
	default:
		// Log unknown event types but don't fail
		fmt.Printf("Received unknown Stripe webhook event: %s\n", webhookEvent.Type)
		return nil
	}
}

// handleCheckoutSessionCompleted processes successful checkout sessions
func (uc *paymentUseCase) handleCheckoutSessionCompleted(ctx context.Context, event *payment.WebhookEvent) error {
	fmt.Printf("🔔 Processing checkout.session.completed webhook\n")

	sessionID, ok := event.Data["session_id"].(string)
	if !ok {
		fmt.Printf("❌ Missing session_id in webhook data: %+v\n", event.Data)
		return fmt.Errorf("missing session_id in webhook data")
	}

	fmt.Printf("🔍 Looking for payment with session ID: %s\n", sessionID)

	// Execute payment confirmation in a transaction
	return uc.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		return uc.confirmPaymentInTransaction(ctx, sessionID)
	})
}

// confirmPaymentInTransaction handles payment confirmation within a transaction
func (uc *paymentUseCase) confirmPaymentInTransaction(ctx context.Context, sessionID string) error {
	// Find payment by session ID (stored in external_id)
	payment, err := uc.paymentRepo.GetByExternalID(ctx, sessionID)
	if err != nil {
		fmt.Printf("❌ Payment not found for session %s: %v\n", sessionID, err)
		return fmt.Errorf("payment not found for session %s: %v", sessionID, err)
	}

	fmt.Printf("✅ Found payment: ID=%s, OrderID=%s, Status=%s\n", payment.ID, payment.OrderID, payment.Status)

	// Check if payment is already processed
	if payment.Status == entities.PaymentStatusPaid {
		fmt.Printf("ℹ️ Payment already processed, skipping\n")
		return nil
	}

	// Update payment status to paid
	payment.MarkAsProcessed(sessionID)
	if err := uc.paymentRepo.Update(ctx, payment); err != nil {
		fmt.Printf("❌ Failed to update payment status: %v\n", err)
		return fmt.Errorf("failed to update payment status: %v", err)
	}

	fmt.Printf("✅ Payment status updated to: %s\n", payment.Status)

	// Get order details
	order, err := uc.orderRepo.GetByID(ctx, payment.OrderID)
	if err != nil {
		fmt.Printf("❌ Order not found: %v\n", err)
		return fmt.Errorf("order not found: %v", err)
	}

	fmt.Printf("✅ Found order: ID=%s, Number=%s, Status=%s, PaymentStatus=%s\n",
		order.ID, order.OrderNumber, order.Status, order.PaymentStatus)

	// Confirm stock reservations (convert to actual stock reduction)
	if err := uc.stockReservationService.ConfirmReservations(ctx, order.ID); err != nil {
		fmt.Printf("❌ Failed to confirm stock reservations: %v\n", err)
		return fmt.Errorf("failed to confirm stock reservations: %v", err)
	}
	fmt.Printf("✅ Stock reservations confirmed and converted to actual stock reduction\n")

	// Reload order with payments to sync payment status
	order, err = uc.orderRepo.GetByID(ctx, order.ID)
	if err != nil {
		fmt.Printf("❌ Failed to reload order: %v\n", err)
		return fmt.Errorf("failed to reload order: %v", err)
	}

	// Update order payment status
	oldStatus := order.Status
	oldPaymentStatus := order.PaymentStatus

	// Sync payment status based on all payments
	order.SyncPaymentStatus(entities.PaymentStatusPaid)
	// Update order status to confirmed if it was pending and fully paid
	if order.Status == entities.OrderStatusPending && order.IsFullyPaid() {
		order.Status = entities.OrderStatusConfirmed
	}
	// Release reservation flags since stock is now actually reduced
	order.ReleaseReservation()
	order.UpdatedAt = time.Now()

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		fmt.Printf("❌ Failed to update order status: %v\n", err)
		return fmt.Errorf("failed to update order status: %v", err)
	}

	fmt.Printf("✅ Order updated: Status %s→%s, PaymentStatus %s→%s\n",
		oldStatus, order.Status, oldPaymentStatus, order.PaymentStatus)

	// Create payment received event within transaction
	if uc.orderEventService != nil {
		if err := uc.orderEventService.CreatePaymentReceivedEvent(ctx, order.ID, payment.Amount, string(payment.Method), &order.UserID); err != nil {
			fmt.Printf("❌ Failed to create payment received event: %v\n", err)
			// Don't fail the transaction for event creation
		} else {
			fmt.Printf("✅ Payment received event created\n")
		}

		// Create status changed event if status changed
		if oldStatus != order.Status {
			if err := uc.orderEventService.CreateStatusChangedEvent(ctx, order.ID, oldStatus, order.Status, &order.UserID); err != nil {
				fmt.Printf("❌ Failed to create status changed event: %v\n", err)
				// Don't fail the transaction for event creation
			} else {
				fmt.Printf("✅ Status changed event created\n")
			}
		}
	}

	// Send payment confirmation notification (async after transaction)
	go func() {
		if uc.notificationUseCase != nil {
			if err := uc.notificationUseCase.NotifyPaymentReceived(context.Background(), payment.ID); err != nil {
				fmt.Printf("❌ Failed to send payment notification: %v\n", err)
			} else {
				fmt.Printf("✅ Payment notification sent\n")
			}
		}
	}()

	fmt.Printf("🎉 Payment confirmation transaction completed successfully\n")
	return nil
}

// handlePaymentIntentSucceeded processes successful payment intents
func (uc *paymentUseCase) handlePaymentIntentSucceeded(ctx context.Context, event *payment.WebhookEvent) error {
	paymentIntentID, ok := event.Data["payment_intent_id"].(string)
	if !ok {
		return fmt.Errorf("missing payment_intent_id in webhook data")
	}

	// Find payment by payment intent ID
	payment, err := uc.paymentRepo.GetByTransactionID(ctx, paymentIntentID)
	if err != nil {
		return fmt.Errorf("payment not found for payment intent %s: %v", paymentIntentID, err)
	}

	// Update payment status
	payment.MarkAsProcessed(paymentIntentID)
	if err := uc.paymentRepo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment status: %v", err)
	}

	// Update order status
	order, err := uc.orderRepo.GetByID(ctx, payment.OrderID)
	if err != nil {
		return fmt.Errorf("order not found: %v", err)
	}

	// Reload order with payments to sync payment status
	order, err = uc.orderRepo.GetByID(ctx, order.ID)
	if err != nil {
		return fmt.Errorf("failed to reload order: %v", err)
	}

	order.SyncPaymentStatus(entities.PaymentStatusPaid)
	if order.Status == entities.OrderStatusPending && order.IsFullyPaid() {
		order.Status = entities.OrderStatusConfirmed

		// Confirm stock reservations (convert to actual stock reduction)
		if err := uc.stockReservationService.ConfirmReservations(ctx, order.ID); err != nil {
			fmt.Printf("❌ Failed to confirm stock reservations: %v\n", err)
			return fmt.Errorf("failed to confirm stock reservations: %v", err)
		}
		fmt.Printf("✅ Stock reservations confirmed and converted to actual stock reduction\n")

		// Release reservation flags since stock is now actually reduced
		order.ReleaseReservation()
	}
	order.UpdatedAt = time.Now()

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}

	return nil
}

// handlePaymentIntentFailed processes failed payment intents
func (uc *paymentUseCase) handlePaymentIntentFailed(ctx context.Context, event *payment.WebhookEvent) error {
	paymentIntentID, ok := event.Data["payment_intent_id"].(string)
	if !ok {
		return fmt.Errorf("missing payment_intent_id in webhook data")
	}

	// Find payment by payment intent ID
	payment, err := uc.paymentRepo.GetByTransactionID(ctx, paymentIntentID)
	if err != nil {
		return fmt.Errorf("payment not found for payment intent %s: %v", paymentIntentID, err)
	}

	// Get failure reason
	failureReason := "Payment failed"
	if lastError, ok := event.Data["last_payment_error"]; ok {
		if errorMap, ok := lastError.(map[string]interface{}); ok {
			if message, ok := errorMap["message"].(string); ok {
				failureReason = message
			}
		}
	}

	// Update payment status to failed
	payment.MarkAsFailed(failureReason)
	if err := uc.paymentRepo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment status: %v", err)
	}

	// Update order payment status
	order, err := uc.orderRepo.GetByID(ctx, payment.OrderID)
	if err != nil {
		return fmt.Errorf("order not found: %v", err)
	}

	// Release stock reservations since payment failed
	if err := uc.stockReservationService.ReleaseReservations(ctx, order.ID); err != nil {
		fmt.Printf("❌ Failed to release stock reservations for failed payment: %v\n", err)
		// Continue with order update even if reservation release fails
	}

	order.PaymentStatus = entities.PaymentStatusFailed
	order.ReleaseReservation()
	order.UpdatedAt = time.Now()

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}

	return nil
}

// handlePayPalWebhook processes PayPal webhook events
func (uc *paymentUseCase) handlePayPalWebhook(ctx context.Context, payload []byte, signature string) error {
	// TODO: Implement PayPal webhook handling
	return fmt.Errorf("paypal webhook handling not implemented yet")
}

// CreateCheckoutSession creates a Stripe checkout session for hosted payment page
func (uc *paymentUseCase) CreateCheckoutSession(ctx context.Context, req CreateCheckoutSessionRequest) (*CreateCheckoutSessionResponse, error) {
	// Validate order exists
	order, err := uc.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return &CreateCheckoutSessionResponse{
			Success: false,
			Message: "Order not found",
		}, err
	}

	// Check if order is in correct status
	if order.Status != entities.OrderStatusPending {
		return &CreateCheckoutSessionResponse{
			Success: false,
			Message: "Order is not in pending status",
		}, fmt.Errorf("order status is %s, expected pending", order.Status)
	}

	// Convert metadata to string map
	metadata := make(map[string]string)
	for k, v := range req.Metadata {
		if str, ok := v.(string); ok {
			metadata[k] = str
		} else {
			metadata[k] = fmt.Sprintf("%v", v)
		}
	}

	// Add order information to metadata
	metadata["order_id"] = req.OrderID.String()
	metadata["order_number"] = order.OrderNumber

	// Use order currency if request currency is empty
	currency := req.Currency
	if currency == "" {
		currency = order.Currency
	}
	if currency == "" {
		currency = "USD" // Default fallback
	}

	// Ensure description is not empty
	description := req.Description
	if description == "" {
		description = fmt.Sprintf("Payment for Order %s", order.OrderNumber)
	}

	// Create checkout session request
	checkoutReq := CheckoutSessionRequest{
		Amount:      order.Total, // Use order total instead of request amount
		Currency:    currency,
		Description: description,
		OrderID:     req.OrderID.String(),
		SuccessURL:  req.SuccessURL,
		CancelURL:   req.CancelURL,
		Metadata:    metadata,
	}

	// Create checkout session using Stripe service
	if uc.stripeService == nil {
		return &CreateCheckoutSessionResponse{
			Success: false,
			Message: "Stripe service not configured",
		}, fmt.Errorf("stripe service not available")
	}

	checkoutResp, err := uc.stripeService.CreateCheckoutSession(ctx, checkoutReq)
	if err != nil {
		return &CreateCheckoutSessionResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create checkout session: %v", err),
		}, err
	}

	if !checkoutResp.Success {
		return &CreateCheckoutSessionResponse{
			Success: false,
			Message: checkoutResp.Message,
		}, fmt.Errorf("checkout session creation failed")
	}

	// Check if payment record already exists for this order
	existingPayment, err := uc.paymentRepo.GetByOrderID(ctx, req.OrderID)
	if err == nil && existingPayment != nil {
		// Update existing payment record with session ID and correct values
		existingPayment.TransactionID = checkoutResp.SessionID
		existingPayment.ExternalID = checkoutResp.SessionID
		existingPayment.Amount = order.Total
		existingPayment.Currency = currency
		existingPayment.Gateway = "stripe"
		existingPayment.UpdatedAt = time.Now()

		if err := uc.paymentRepo.Update(ctx, existingPayment); err != nil {
			return &CreateCheckoutSessionResponse{
				Success: false,
				Message: "Failed to update payment record",
			}, err
		}
	} else {
		// Create new payment record
		paymentEntity := &entities.Payment{
			ID:            uuid.New(),
			OrderID:       req.OrderID,
			UserID:        order.UserID,
			Amount:        order.Total, // Use order total
			Currency:      currency,    // Use resolved currency
			Method:        entities.PaymentMethodStripe,
			Status:        entities.PaymentStatusPending,
			TransactionID: checkoutResp.SessionID,
			ExternalID:    checkoutResp.SessionID,
			Gateway:       "stripe",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Save payment record
		if err := uc.paymentRepo.Create(ctx, paymentEntity); err != nil {
			return &CreateCheckoutSessionResponse{
				Success: false,
				Message: "Failed to create payment record",
			}, err
		}
	}

	return &CreateCheckoutSessionResponse{
		Success:    true,
		SessionID:  checkoutResp.SessionID,
		SessionURL: checkoutResp.SessionURL,
		Message:    "Checkout session created successfully",
	}, nil
}

// GetPaymentReport gets payment report (placeholder implementation)
func (uc *paymentUseCase) GetPaymentReport(ctx context.Context, req PaymentReportRequest) (*PaymentReportResponse, error) {
	// This is a placeholder implementation
	// In a real implementation, you would generate the report from the database
	return &PaymentReportResponse{
		ReportType:  req.ReportType,
		ReportID:    uuid.New(),
		GeneratedAt: time.Now(),
		Format:      "json",
		Summary: &PaymentReportSummary{
			TotalPayments:      0,
			TotalAmount:        0,
			SuccessfulPayments: 0,
			FailedPayments:     0,
			RefundAmount:       0,
			SuccessRate:        0,
		},
		Data: []map[string]interface{}{},
	}, nil
}

// Helper methods
func (uc *paymentUseCase) toPaymentResponse(payment *entities.Payment) *PaymentResponse {
	return &PaymentResponse{
		ID:              payment.ID,
		OrderID:         payment.OrderID,
		Amount:          payment.Amount,
		Currency:        payment.Currency,
		Method:          payment.Method,
		Status:          payment.Status,
		TransactionID:   payment.TransactionID,
		ExternalID:      payment.ExternalID,
		FailureReason:   payment.FailureReason,
		ProcessedAt:     payment.ProcessedAt,
		RefundedAt:      payment.RefundedAt,
		RefundAmount:    payment.RefundAmount,
		CanBeRefunded:   payment.CanBeRefunded(),
		RemainingRefund: payment.GetRemainingRefundAmount(),
		CreatedAt:       payment.CreatedAt,
		UpdatedAt:       payment.UpdatedAt,
	}
}

// Helper method to convert refund entity to response
func (uc *paymentUseCase) toRefundResponse(refund *entities.Refund) *RefundResponse {
	return uc.mapRefundToResponse(refund)
}

// GetRefunds gets all refunds for a payment
func (uc *paymentUseCase) GetRefunds(ctx context.Context, paymentID uuid.UUID) ([]*RefundResponse, error) {
	refunds, err := uc.paymentRepo.GetRefundsByPaymentID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	responses := make([]*RefundResponse, len(refunds))
	for i, refund := range refunds {
		responses[i] = uc.toRefundResponse(refund)
	}

	return responses, nil
}

// SavePaymentMethod saves a payment method
func (uc *paymentUseCase) SavePaymentMethod(ctx context.Context, req SavePaymentMethodRequest) (*PaymentMethodResponse, error) {
	// Validate user exists
	_, err := uc.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, entities.ErrUserNotFound
	}

	// Check for duplicate fingerprint if provided
	if req.Fingerprint != "" {
		existing, err := uc.paymentMethodRepo.GetByFingerprint(ctx, req.UserID, req.Fingerprint)
		if err == nil && existing != nil {
			return nil, entities.ErrPaymentMethodExists
		}
	}

	// Create payment method entity
	paymentMethod := &entities.PaymentMethodEntity{
		UserID:            req.UserID,
		Type:              req.Type,
		Last4:             req.Last4,
		Brand:             req.Brand,
		ExpiryMonth:       req.ExpiryMonth,
		ExpiryYear:        req.ExpiryYear,
		Gateway:           req.Gateway,
		GatewayToken:      req.Token,
		GatewayCustomerID: req.GatewayCustomerID,
		BillingName:       req.BillingName,
		BillingEmail:      req.BillingEmail,
		BillingAddress:    req.BillingAddress,
		IsDefault:         req.IsDefault,
		IsActive:          true,
		Fingerprint:       req.Fingerprint,
		Metadata:          req.MetadataJSON,
		Notes:             req.Notes,
	}

	// If this is set as default, unset other defaults first
	if req.IsDefault {
		if err := uc.paymentMethodRepo.UnsetDefault(ctx, req.UserID); err != nil {
			return nil, err
		}
	}

	// Save payment method
	if err := uc.paymentMethodRepo.Create(ctx, paymentMethod); err != nil {
		return nil, err
	}

	return uc.toPaymentMethodResponse(paymentMethod), nil
}

// GetUserPaymentMethods gets all payment methods for a user
func (uc *paymentUseCase) GetUserPaymentMethods(ctx context.Context, userID uuid.UUID) ([]*PaymentMethodResponse, error) {
	// Validate user exists
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, entities.ErrUserNotFound
	}

	// Get active payment methods
	paymentMethods, err := uc.paymentMethodRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	responses := make([]*PaymentMethodResponse, len(paymentMethods))
	for i, pm := range paymentMethods {
		responses[i] = uc.toPaymentMethodResponse(pm)
	}

	return responses, nil
}

// DeletePaymentMethod deletes a payment method
func (uc *paymentUseCase) DeletePaymentMethod(ctx context.Context, id uuid.UUID) error {
	// Get payment method to check if it exists and is default
	paymentMethod, err := uc.paymentMethodRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if it's the default payment method
	if paymentMethod.IsDefault {
		// Check if user has other payment methods
		count, err := uc.paymentMethodRepo.Count(ctx, paymentMethod.UserID)
		if err != nil {
			return err
		}

		// If this is the only payment method, allow deletion
		// If there are others, require setting a new default first
		if count > 1 {
			return entities.ErrCannotDeleteDefaultPaymentMethod
		}
	}

	// Deactivate instead of hard delete for audit purposes
	return uc.paymentMethodRepo.Deactivate(ctx, id)
}

// SetDefaultPaymentMethod sets a payment method as default
func (uc *paymentUseCase) SetDefaultPaymentMethod(ctx context.Context, userID, methodID uuid.UUID) error {
	// Validate that the payment method belongs to the user
	paymentMethod, err := uc.paymentMethodRepo.GetByID(ctx, methodID)
	if err != nil {
		return err
	}

	if paymentMethod.UserID != userID {
		return entities.ErrForbidden
	}

	if !paymentMethod.IsActive {
		return entities.ErrPaymentMethodInactive
	}

	// Check if card is expired
	if paymentMethod.IsExpired() {
		return entities.ErrPaymentMethodExpired
	}

	// Set as default (this will unset others automatically)
	return uc.paymentMethodRepo.SetAsDefault(ctx, userID, methodID)
}

// ConfirmPaymentSuccess confirms payment success for an order (fallback method)
func (uc *paymentUseCase) ConfirmPaymentSuccess(ctx context.Context, orderID, userID uuid.UUID, sessionID string) error {
	fmt.Printf("🔄 Fallback payment confirmation: OrderID=%s, UserID=%s, SessionID=%s\n", orderID, userID, sessionID)

	// Get the order
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		fmt.Printf("❌ Order not found: %v\n", err)
		return fmt.Errorf("order not found: %v", err)
	}

	// Verify the order belongs to the user (only if userID is provided)
	if userID != uuid.Nil && order.UserID != userID {
		fmt.Printf("❌ Order does not belong to user: OrderUserID=%s, RequestUserID=%s\n", order.UserID, userID)
		return fmt.Errorf("order does not belong to user")
	}

	fmt.Printf("✅ Order verified: %s (Status: %s, PaymentStatus: %s)\n", order.OrderNumber, order.Status, order.PaymentStatus)

	// Find payment by session ID (stored in external_id)
	payment, err := uc.paymentRepo.GetByExternalID(ctx, sessionID)
	if err != nil {
		fmt.Printf("❌ Payment not found for session %s: %v\n", sessionID, err)
		return fmt.Errorf("payment not found for session %s: %v", sessionID, err)
	}

	// Verify the payment belongs to the order
	if payment.OrderID != orderID {
		fmt.Printf("❌ Payment does not belong to order: PaymentOrderID=%s, RequestOrderID=%s\n", payment.OrderID, orderID)
		return fmt.Errorf("payment does not belong to order")
	}

	fmt.Printf("✅ Payment verified: ID=%s (Status: %s)\n", payment.ID, payment.Status)

	// If payment is already processed, skip
	if payment.Status == entities.PaymentStatusPaid {
		fmt.Printf("ℹ️ Payment already processed, skipping\n")
		return nil
	}

	// Update payment status to paid
	payment.MarkAsProcessed(sessionID)
	if err := uc.paymentRepo.Update(ctx, payment); err != nil {
		fmt.Printf("❌ Failed to update payment status: %v\n", err)
		return fmt.Errorf("failed to update payment status: %v", err)
	}

	fmt.Printf("✅ Payment status updated to: %s\n", payment.Status)

	// Reload order with payments to sync payment status
	order, err = uc.orderRepo.GetByID(ctx, order.ID)
	if err != nil {
		fmt.Printf("❌ Failed to reload order: %v\n", err)
		return fmt.Errorf("failed to reload order: %v", err)
	}

	// Update order status
	oldStatus := order.Status
	oldPaymentStatus := order.PaymentStatus

	order.SyncPaymentStatus(entities.PaymentStatusPaid)
	if order.Status == entities.OrderStatusPending && order.IsFullyPaid() {
		order.Status = entities.OrderStatusConfirmed

		// Confirm stock reservations (convert to actual stock reduction)
		if err := uc.stockReservationService.ConfirmReservations(ctx, order.ID); err != nil {
			fmt.Printf("❌ Failed to confirm stock reservations: %v\n", err)
			return fmt.Errorf("failed to confirm stock reservations: %v", err)
		}
		fmt.Printf("✅ Stock reservations confirmed and converted to actual stock reduction\n")

		// Release reservation flags since stock is now actually reduced
		order.ReleaseReservation()
	}
	order.UpdatedAt = time.Now()

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		fmt.Printf("❌ Failed to update order status: %v\n", err)
		return fmt.Errorf("failed to update order status: %v", err)
	}

	fmt.Printf("✅ Order updated via fallback: Status %s→%s, PaymentStatus %s→%s\n",
		oldStatus, order.Status, oldPaymentStatus, order.PaymentStatus)

	// Send payment confirmation notification if available
	if uc.notificationUseCase != nil {
		if err := uc.notificationUseCase.NotifyPaymentReceived(ctx, payment.ID); err != nil {
			fmt.Printf("❌ Failed to send payment notification: %v\n", err)
		} else {
			fmt.Printf("✅ Payment notification sent\n")
		}
	}

	fmt.Printf("🎉 Fallback payment confirmation completed\n")
	return nil
}

// toPaymentMethodResponse converts PaymentMethodEntity to PaymentMethodResponse
func (uc *paymentUseCase) toPaymentMethodResponse(pm *entities.PaymentMethodEntity) *PaymentMethodResponse {
	return &PaymentMethodResponse{
		ID:            pm.ID,
		UserID:        pm.UserID,
		Type:          pm.Type,
		Last4:         pm.Last4,
		Brand:         pm.Brand,
		ExpiryMonth:   pm.ExpiryMonth,
		ExpiryYear:    pm.ExpiryYear,
		Gateway:       pm.Gateway,
		BillingName:   pm.BillingName,
		BillingEmail:  pm.BillingEmail,
		IsDefault:     pm.IsDefault,
		IsActive:      pm.IsActive,
		IsExpired:     pm.IsExpired(),
		DisplayName:   pm.GetDisplayName(),
		CreatedAt:     pm.CreatedAt,
		UpdatedAt:     pm.UpdatedAt,
		LastUsedAt:    pm.LastUsedAt,
	}
}

// getGatewayForMethod returns the appropriate gateway for a payment method
func (uc *paymentUseCase) getGatewayForMethod(method entities.PaymentMethod) string {
	switch method {
	case entities.PaymentMethodStripe, entities.PaymentMethodCreditCard, entities.PaymentMethodDebitCard:
		return "stripe"
	case entities.PaymentMethodPayPal:
		return "paypal"
	case entities.PaymentMethodApplePay:
		return "apple_pay"
	case entities.PaymentMethodGooglePay:
		return "google_pay"
	case entities.PaymentMethodBankTransfer:
		return "bank_transfer"
	case entities.PaymentMethodCash:
		return "cod"
	default:
		return "stripe" // Default gateway
	}
}

// isRetryableError checks if an error is retryable
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errorMsg := err.Error()

	// Network-related errors that can be retried
	retryableErrors := []string{
		"timeout",
		"connection refused",
		"network unreachable",
		"temporary failure",
		"service unavailable",
		"rate limit",
		"too many requests",
		"internal server error",
		"gateway timeout",
		"bad gateway",
	}

	for _, retryableErr := range retryableErrors {
		if strings.Contains(errorMsg, retryableErr) {
			return true
		}
	}

	return false
}


