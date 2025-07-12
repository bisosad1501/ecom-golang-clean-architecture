package usecases

import (
	"context"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/domain/services"
	"ecom-golang-clean-architecture/internal/infrastructure/payment"

	"github.com/google/uuid"
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
	orderRepo               repositories.OrderRepository
	userRepo                repositories.UserRepository
	stripeService           PaymentGatewayService
	paypalService           PaymentGatewayService
	notificationUseCase     NotificationUseCase
	stockReservationService services.StockReservationService
}

// NewPaymentUseCase creates a new payment use case
func NewPaymentUseCase(
	paymentRepo repositories.PaymentRepository,
	orderRepo repositories.OrderRepository,
	userRepo repositories.UserRepository,
	stripeService PaymentGatewayService,
	paypalService PaymentGatewayService,
	notificationUseCase NotificationUseCase,
	stockReservationService services.StockReservationService,
) PaymentUseCase {
	return &paymentUseCase{
		paymentRepo:             paymentRepo,
		orderRepo:               orderRepo,
		userRepo:                userRepo,
		stripeService:           stripeService,
		paypalService:           paypalService,
		notificationUseCase:     notificationUseCase,
		stockReservationService: stockReservationService,
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
	PaymentID uuid.UUID              `json:"payment_id" validate:"required"`
	Amount    float64                `json:"amount" validate:"required,gt=0"`
	Reason    string                 `json:"reason" validate:"required"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type SavePaymentMethodRequest struct {
	UserID         uuid.UUID              `json:"user_id" validate:"required"`
	Type           entities.PaymentMethod `json:"type" validate:"required"`
	Token          string                 `json:"token" validate:"required"`
	IsDefault      bool                   `json:"is_default"`
	BillingAddress *BillingAddressRequest `json:"billing_address,omitempty"`
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
	ID            uuid.UUID  `json:"id"`
	PaymentID     uuid.UUID  `json:"payment_id"`
	Amount        float64    `json:"amount"`
	Reason        string     `json:"reason"`
	Status        string     `json:"status"`
	TransactionID string     `json:"transaction_id"`
	ProcessedAt   *time.Time `json:"processed_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

type PaymentMethodResponse struct {
	ID             uuid.UUID               `json:"id"`
	UserID         uuid.UUID               `json:"user_id"`
	Type           entities.PaymentMethod  `json:"type"`
	Last4          string                  `json:"last4,omitempty"`
	Brand          string                  `json:"brand,omitempty"`
	ExpiryMonth    int                     `json:"expiry_month,omitempty"`
	ExpiryYear     int                     `json:"expiry_year,omitempty"`
	IsDefault      bool                    `json:"is_default"`
	BillingAddress *BillingAddressResponse `json:"billing_address,omitempty"`
	CreatedAt      time.Time               `json:"created_at"`
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

	// Validate payment amount matches order total
	if req.Amount != order.Total {
		return nil, fmt.Errorf("payment amount %.2f does not match order total %.2f", req.Amount, order.Total)
	}

	// Create payment record
	payment := &entities.Payment{
		ID:        uuid.New(),
		OrderID:   req.OrderID,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Method:    req.Method,
		Status:    entities.PaymentStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
	default:
		return nil, fmt.Errorf("unsupported payment method: %s", req.Method)
	}

	if err != nil {
		payment.MarkAsFailed(err.Error())
		if err := uc.paymentRepo.Update(ctx, payment); err != nil {
			return nil, err
		}
		return nil, err
	}

	// Update payment with gateway response
	if gatewayResp.Success {
		payment.MarkAsProcessed(gatewayResp.TransactionID)
		payment.ExternalID = gatewayResp.ExternalID

		// Update order payment status
		order.PaymentStatus = entities.PaymentStatusPaid
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
	payment, err := uc.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if payment == nil {
		return []*PaymentResponse{}, nil
	}

	return []*PaymentResponse{uc.toPaymentResponse(payment)}, nil
}

// UpdatePaymentStatus updates payment status
func (uc *paymentUseCase) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status entities.PaymentStatus, transactionID string) (*PaymentResponse, error) {
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

	return uc.toPaymentResponse(payment), nil
}

// ProcessRefund processes a refund for a payment
func (uc *paymentUseCase) ProcessRefund(ctx context.Context, req ProcessRefundRequest) (*RefundResponse, error) {
	// Get payment details
	payment, err := uc.paymentRepo.GetByID(ctx, req.PaymentID)
	if err != nil {
		return nil, entities.ErrPaymentNotFound
	}

	// Validate refund
	if !payment.CanBeRefunded() {
		return nil, fmt.Errorf("payment cannot be refunded")
	}

	if req.Amount > payment.GetRemainingRefundAmount() {
		return nil, fmt.Errorf("refund amount exceeds remaining refundable amount")
	}

	// Process refund through gateway
	var gatewayResp *RefundGatewayResponse
	refundReq := RefundGatewayRequest{
		TransactionID: payment.TransactionID,
		Amount:        req.Amount,
		Reason:        req.Reason,
	}

	switch payment.Method {
	case entities.PaymentMethodStripe:
		gatewayResp, err = uc.stripeService.ProcessRefund(ctx, refundReq)
	case entities.PaymentMethodPayPal:
		gatewayResp, err = uc.paypalService.ProcessRefund(ctx, refundReq)
	default:
		return nil, fmt.Errorf("unsupported payment method for refund: %s", payment.Method)
	}

	if err != nil {
		return nil, err
	}

	if gatewayResp.Success {
		// Update payment with refund
		if err := payment.AddRefund(req.Amount); err != nil {
			return nil, err
		}

		// Save updated payment
		if err := uc.paymentRepo.Update(ctx, payment); err != nil {
			return nil, err
		}

		// Create refund record
		refund := &entities.Refund{
			ID:            uuid.New(),
			PaymentID:     req.PaymentID,
			Amount:        req.Amount,
			Reason:        req.Reason,
			Status:        entities.RefundStatusCompleted,
			TransactionID: gatewayResp.RefundID,
			ProcessedAt:   &time.Time{},
			CreatedAt:     time.Now(),
		}
		*refund.ProcessedAt = time.Now()

		if err := uc.paymentRepo.CreateRefund(ctx, refund); err != nil {
			return nil, err
		}

		return &RefundResponse{
			ID:            refund.ID,
			PaymentID:     refund.PaymentID,
			Amount:        refund.Amount,
			Reason:        refund.Reason,
			Status:        string(refund.Status),
			TransactionID: refund.TransactionID,
			ProcessedAt:   refund.ProcessedAt,
			CreatedAt:     refund.CreatedAt,
		}, nil
	}

	return nil, fmt.Errorf("refund failed: %s", gatewayResp.Message)
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

	// Find payment by session ID (stored in external_id)
	payment, err := uc.paymentRepo.GetByExternalID(ctx, sessionID)
	if err != nil {
		fmt.Printf("❌ Payment not found for session %s: %v\n", sessionID, err)
		return fmt.Errorf("payment not found for session %s: %v", sessionID, err)
	}

	fmt.Printf("✅ Found payment: ID=%s, OrderID=%s, Status=%s\n", payment.ID, payment.OrderID, payment.Status)

	// Update payment status to paid
	payment.MarkAsProcessed(sessionID)
	if err := uc.paymentRepo.Update(ctx, payment); err != nil {
		fmt.Printf("❌ Failed to update payment status: %v\n", err)
		return fmt.Errorf("failed to update payment status: %v", err)
	}

	fmt.Printf("✅ Payment status updated to: %s\n", payment.Status)

	// Update order status and payment status
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

	// Update order payment status
	oldStatus := order.Status
	oldPaymentStatus := order.PaymentStatus

	order.PaymentStatus = entities.PaymentStatusPaid
	// Update order status to confirmed if it was pending
	if order.Status == entities.OrderStatusPending {
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

	// Send payment confirmation notification
	if uc.notificationUseCase != nil {
		uc.notificationUseCase.NotifyPaymentReceived(ctx, payment.ID)
		fmt.Printf("✅ Payment notification sent\n")
	}

	fmt.Printf("🎉 Webhook processing completed successfully\n")
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

	order.PaymentStatus = entities.PaymentStatusPaid
	if order.Status == entities.OrderStatusPending {
		order.Status = entities.OrderStatusConfirmed
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
	return &RefundResponse{
		ID:            refund.ID,
		PaymentID:     refund.PaymentID,
		Amount:        refund.Amount,
		Reason:        refund.Reason,
		Status:        string(refund.Status),
		TransactionID: refund.TransactionID,
		ProcessedAt:   refund.ProcessedAt,
		CreatedAt:     refund.CreatedAt,
	}
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

// SavePaymentMethod saves a payment method (placeholder implementation)
func (uc *paymentUseCase) SavePaymentMethod(ctx context.Context, req SavePaymentMethodRequest) (*PaymentMethodResponse, error) {
	// This is a placeholder implementation since we don't have a payment method entity
	// In a real implementation, you would save the payment method to the database
	response := &PaymentMethodResponse{
		ID:        uuid.New(),
		UserID:    req.UserID,
		Type:      req.Type,
		Last4:     "****",    // Placeholder since not in request
		Brand:     "unknown", // Placeholder since not in request
		IsDefault: req.IsDefault,
		CreatedAt: time.Now(),
	}

	return response, nil
}

// GetUserPaymentMethods gets all payment methods for a user (placeholder implementation)
func (uc *paymentUseCase) GetUserPaymentMethods(ctx context.Context, userID uuid.UUID) ([]*PaymentMethodResponse, error) {
	// This is a placeholder implementation since we don't have a payment method entity
	// In a real implementation, you would fetch from the database
	return []*PaymentMethodResponse{}, nil
}

// DeletePaymentMethod deletes a payment method (placeholder implementation)
func (uc *paymentUseCase) DeletePaymentMethod(ctx context.Context, id uuid.UUID) error {
	// This is a placeholder implementation since we don't have a payment method entity
	// In a real implementation, you would delete from the database
	return nil
}

// SetDefaultPaymentMethod sets default payment method (placeholder implementation)
func (uc *paymentUseCase) SetDefaultPaymentMethod(ctx context.Context, userID, methodID uuid.UUID) error {
	// This is a placeholder implementation since we don't have a payment method entity
	// In a real implementation, you would update the database
	return nil
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

	// Update order status
	oldStatus := order.Status
	oldPaymentStatus := order.PaymentStatus

	order.PaymentStatus = entities.PaymentStatusPaid
	if order.Status == entities.OrderStatusPending {
		order.Status = entities.OrderStatusConfirmed
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
		uc.notificationUseCase.NotifyPaymentReceived(ctx, payment.ID)
		fmt.Printf("✅ Payment notification sent\n")
	}

	fmt.Printf("🎉 Fallback payment confirmation completed\n")
	return nil
}
