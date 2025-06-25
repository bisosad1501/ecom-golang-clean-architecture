package usecases

import (
	"context"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

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
	
	// Reports
	GetPaymentReport(ctx context.Context, req PaymentReportRequest) (*PaymentReportResponse, error)
}

type paymentUseCase struct {
	paymentRepo       repositories.PaymentRepository
	orderRepo         repositories.OrderRepository
	userRepo          repositories.UserRepository
	stripeService     PaymentGatewayService
	paypalService     PaymentGatewayService
	notificationUseCase NotificationUseCase
}

// NewPaymentUseCase creates a new payment use case
func NewPaymentUseCase(
	paymentRepo repositories.PaymentRepository,
	orderRepo repositories.OrderRepository,
	userRepo repositories.UserRepository,
	stripeService PaymentGatewayService,
	paypalService PaymentGatewayService,
	notificationUseCase NotificationUseCase,
) PaymentUseCase {
	return &paymentUseCase{
		paymentRepo:         paymentRepo,
		orderRepo:           orderRepo,
		userRepo:            userRepo,
		stripeService:       stripeService,
		paypalService:       paypalService,
		notificationUseCase: notificationUseCase,
	}
}

// Payment gateway service interface
type PaymentGatewayService interface {
	ProcessPayment(ctx context.Context, req PaymentGatewayRequest) (*PaymentGatewayResponse, error)
	ProcessRefund(ctx context.Context, transactionID string, amount float64) (*RefundGatewayResponse, error)
	VerifyWebhook(payload []byte, signature string) (bool, error)
	ParseWebhook(payload []byte) (*WebhookEvent, error)
}

// Request/Response types
type ProcessPaymentRequest struct {
	OrderID         uuid.UUID                 `json:"order_id" validate:"required"`
	Amount          float64                   `json:"amount" validate:"required,gt=0"`
	Currency        string                    `json:"currency" validate:"required"`
	Method          entities.PaymentMethod    `json:"method" validate:"required"`
	PaymentToken    string                    `json:"payment_token,omitempty"`
	PaymentMethodID *uuid.UUID                `json:"payment_method_id,omitempty"`
	BillingAddress  *BillingAddressRequest    `json:"billing_address,omitempty"`
	Metadata        map[string]interface{}    `json:"metadata,omitempty"`
}

type ProcessRefundRequest struct {
	PaymentID uuid.UUID `json:"payment_id" validate:"required"`
	Amount    float64   `json:"amount" validate:"required,gt=0"`
	Reason    string    `json:"reason" validate:"required"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type SavePaymentMethodRequest struct {
	UserID      uuid.UUID                 `json:"user_id" validate:"required"`
	Type        entities.PaymentMethod    `json:"type" validate:"required"`
	Token       string                    `json:"token" validate:"required"`
	IsDefault   bool                      `json:"is_default"`
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
	ReportType string                    `json:"report_type"`
	DateFrom   *time.Time                `json:"date_from,omitempty"`
	DateTo     *time.Time                `json:"date_to,omitempty"`
	Status     *entities.PaymentStatus   `json:"status,omitempty"`
	Method     *entities.PaymentMethod   `json:"method,omitempty"`
	GroupBy    string                    `json:"group_by,omitempty" validate:"omitempty,oneof=day week month method status"`
	Format     string                    `json:"format,omitempty" validate:"omitempty,oneof=json csv excel"`
}

type PaymentGatewayRequest struct {
	Amount          float64                   `json:"amount"`
	Currency        string                    `json:"currency"`
	PaymentToken    string                    `json:"payment_token,omitempty"`
	PaymentMethodID string                    `json:"payment_method_id,omitempty"`
	CustomerID      string                    `json:"customer_id,omitempty"`
	Description     string                    `json:"description,omitempty"`
	Metadata        map[string]interface{}    `json:"metadata,omitempty"`
}

// Response types
type PaymentResponse struct {
	ID              uuid.UUID                 `json:"id"`
	OrderID         uuid.UUID                 `json:"order_id"`
	Amount          float64                   `json:"amount"`
	Currency        string                    `json:"currency"`
	Method          entities.PaymentMethod    `json:"method"`
	Status          entities.PaymentStatus    `json:"status"`
	TransactionID   string                    `json:"transaction_id"`
	ExternalID      string                    `json:"external_id"`
	FailureReason   string                    `json:"failure_reason,omitempty"`
	ProcessedAt     *time.Time                `json:"processed_at"`
	RefundedAt      *time.Time                `json:"refunded_at"`
	RefundAmount    float64                   `json:"refund_amount"`
	CanBeRefunded   bool                      `json:"can_be_refunded"`
	RemainingRefund float64                   `json:"remaining_refund"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

type RefundResponse struct {
	ID            uuid.UUID `json:"id"`
	PaymentID     uuid.UUID `json:"payment_id"`
	Amount        float64   `json:"amount"`
	Reason        string    `json:"reason"`
	Status        string    `json:"status"`
	TransactionID string    `json:"transaction_id"`
	ProcessedAt   *time.Time `json:"processed_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type PaymentMethodResponse struct {
	ID             uuid.UUID                 `json:"id"`
	UserID         uuid.UUID                 `json:"user_id"`
	Type           entities.PaymentMethod    `json:"type"`
	Last4          string                    `json:"last4,omitempty"`
	Brand          string                    `json:"brand,omitempty"`
	ExpiryMonth    int                       `json:"expiry_month,omitempty"`
	ExpiryYear     int                       `json:"expiry_year,omitempty"`
	IsDefault      bool                      `json:"is_default"`
	BillingAddress *BillingAddressResponse   `json:"billing_address,omitempty"`
	CreatedAt      time.Time                 `json:"created_at"`
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

type PaymentGatewayResponse struct {
	Success       bool                      `json:"success"`
	TransactionID string                    `json:"transaction_id"`
	ExternalID    string                    `json:"external_id"`
	Status        entities.PaymentStatus    `json:"status"`
	Message       string                    `json:"message"`
	Metadata      map[string]interface{}    `json:"metadata,omitempty"`
}

type RefundGatewayResponse struct {
	Success       bool   `json:"success"`
	RefundID      string `json:"refund_id"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type WebhookEvent struct {
	Type          string                    `json:"type"`
	PaymentID     string                    `json:"payment_id"`
	TransactionID string                    `json:"transaction_id"`
	Status        entities.PaymentStatus    `json:"status"`
	Amount        float64                   `json:"amount"`
	Metadata      map[string]interface{}    `json:"metadata,omitempty"`
}

type PaymentReportResponse struct {
	ReportType  string                  `json:"report_type"`
	ReportID    uuid.UUID               `json:"report_id"`
	GeneratedAt time.Time               `json:"generated_at"`
	Format      string                  `json:"format"`
	DownloadURL string                  `json:"download_url,omitempty"`
	Summary     *PaymentReportSummary   `json:"summary"`
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
	gatewayReq := PaymentGatewayRequest{
		Amount:          req.Amount,
		Currency:        req.Currency,
		PaymentToken:    req.PaymentToken,
		Description:     fmt.Sprintf("Payment for order %s", order.OrderNumber),
		Metadata:        req.Metadata,
	}

	if req.PaymentMethodID != nil {
		gatewayReq.PaymentMethodID = req.PaymentMethodID.String()
	}

	var gatewayResp *PaymentGatewayResponse
	switch req.Method {
	case entities.PaymentMethodStripe:
		gatewayResp, err = uc.stripeService.ProcessPayment(ctx, gatewayReq)
	case entities.PaymentMethodPayPal:
		gatewayResp, err = uc.paypalService.ProcessPayment(ctx, gatewayReq)
	default:
		return nil, fmt.Errorf("unsupported payment method: %s", req.Method)
	}

	if err != nil {
		payment.MarkAsFailed(err.Error())
		uc.paymentRepo.Update(ctx, payment)
		return nil, err
	}

	// Update payment with gateway response
	if gatewayResp.Success {
		payment.MarkAsProcessed(gatewayResp.TransactionID)
		payment.ExternalID = gatewayResp.ExternalID
		
		// Update order payment status
		order.PaymentStatus = entities.PaymentStatusPaid
		uc.orderRepo.Update(ctx, order)
		
		// Send payment confirmation notification
		if uc.notificationUseCase != nil {
			uc.notificationUseCase.NotifyPaymentReceived(ctx, payment.ID)
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
	switch payment.Method {
	case entities.PaymentMethodStripe:
		gatewayResp, err = uc.stripeService.ProcessRefund(ctx, payment.TransactionID, req.Amount)
	case entities.PaymentMethodPayPal:
		gatewayResp, err = uc.paypalService.ProcessRefund(ctx, payment.TransactionID, req.Amount)
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
	// Choose the appropriate service based on provider
	var service PaymentGatewayService
	switch provider {
	case "stripe":
		service = uc.stripeService
	case "paypal":
		service = uc.paypalService
	default:
		return fmt.Errorf("unsupported payment provider: %s", provider)
	}
	
	// Verify webhook signature
	verified, err := service.VerifyWebhook(payload, signature)
	if err != nil {
		return fmt.Errorf("failed to verify webhook: %w", err)
	}
	
	if !verified {
		return fmt.Errorf("invalid webhook signature")
	}
	
	// Parse webhook event
	event, err := service.ParseWebhook(payload)
	if err != nil {
		return fmt.Errorf("failed to parse webhook: %w", err)
	}
	
	// Process the event based on type
	switch event.Type {
	case "payment.succeeded":
		// Update payment status to paid
		if event.PaymentID != "" {
			_, err = uc.UpdatePaymentStatus(ctx, uuid.MustParse(event.PaymentID), entities.PaymentStatusPaid, event.TransactionID)
		}
	case "payment.failed":
		// Update payment status to failed
		if event.PaymentID != "" {
			_, err = uc.UpdatePaymentStatus(ctx, uuid.MustParse(event.PaymentID), entities.PaymentStatusFailed, "")
		}
	}
	
	return err
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
		ID:       uuid.New(),
		UserID:   req.UserID,
		Type:     req.Type,
		Last4:    "****", // Placeholder since not in request
		Brand:    "unknown", // Placeholder since not in request
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
