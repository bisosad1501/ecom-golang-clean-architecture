package payment

import (
	"context"
	"fmt"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/refund"
)

// StripeService implements payment processing with Stripe
type StripeService struct {
	apiKey string
}

// NewStripeService creates a new Stripe service
func NewStripeService(apiKey string) *StripeService {
	stripe.Key = apiKey
	return &StripeService{
		apiKey: apiKey,
	}
}

// PaymentGatewayRequest represents a payment request
type PaymentGatewayRequest struct {
	Amount          float64           `json:"amount"`
	Currency        string            `json:"currency"`
	PaymentToken    string            `json:"payment_token"`
	PaymentMethodID string            `json:"payment_method_id"`
	Description     string            `json:"description"`
	Metadata        map[string]string `json:"metadata"`
}

// PaymentGatewayResponse represents a payment response
type PaymentGatewayResponse struct {
	Success       bool   `json:"success"`
	TransactionID string `json:"transaction_id"`
	ExternalID    string `json:"external_id"`
	Message       string `json:"message"`
	Status        string `json:"status"`
}

// RefundGatewayRequest represents a refund request
type RefundGatewayRequest struct {
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Reason        string  `json:"reason"`
}

// RefundGatewayResponse represents a refund response
type RefundGatewayResponse struct {
	Success  bool   `json:"success"`
	RefundID string `json:"refund_id"`
	Message  string `json:"message"`
	Status   string `json:"status"`
}

// CheckoutSessionRequest represents a checkout session request
type CheckoutSessionRequest struct {
	Amount      float64           `json:"amount"`
	Currency    string            `json:"currency"`
	Description string            `json:"description"`
	OrderID     string            `json:"order_id"`
	CustomerID  string            `json:"customer_id,omitempty"`
	SuccessURL  string            `json:"success_url"`
	CancelURL   string            `json:"cancel_url"`
	Metadata    map[string]string `json:"metadata"`
}

// CheckoutSessionResponse represents a checkout session response
type CheckoutSessionResponse struct {
	Success     bool   `json:"success"`
	SessionID   string `json:"session_id"`
	SessionURL  string `json:"session_url"`
	Message     string `json:"message"`
}

// ProcessPayment processes a payment through Stripe
func (s *StripeService) ProcessPayment(ctx context.Context, req PaymentGatewayRequest) (*PaymentGatewayResponse, error) {
	// Convert amount to cents (Stripe uses smallest currency unit)
	amountCents := int64(req.Amount * 100)

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(req.Currency),
		Confirm:  stripe.Bool(true),
	}

	// Set payment method
	if req.PaymentMethodID != "" {
		params.PaymentMethod = stripe.String(req.PaymentMethodID)
	} else if req.PaymentToken != "" {
		params.PaymentMethod = stripe.String(req.PaymentToken)
	}

	// Add description
	if req.Description != "" {
		params.Description = stripe.String(req.Description)
	}

	// Add metadata
	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}

	// Create payment intent
	pi, err := paymentintent.New(params)
	if err != nil {
		return &PaymentGatewayResponse{
			Success: false,
			Message: fmt.Sprintf("Stripe payment failed: %v", err),
		}, err
	}

	// Check payment status
	success := pi.Status == stripe.PaymentIntentStatusSucceeded
	
	return &PaymentGatewayResponse{
		Success:       success,
		TransactionID: pi.ID,
		ExternalID:    pi.ID,
		Message:       string(pi.Status),
		Status:        string(pi.Status),
	}, nil
}

// ProcessRefund processes a refund through Stripe
func (s *StripeService) ProcessRefund(ctx context.Context, req RefundGatewayRequest) (*RefundGatewayResponse, error) {
	// Convert amount to cents
	amountCents := int64(req.Amount * 100)

	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(req.TransactionID),
		Amount:        stripe.Int64(amountCents),
	}

	if req.Reason != "" {
		params.Reason = stripe.String(req.Reason)
	}

	// Create refund
	r, err := refund.New(params)
	if err != nil {
		return &RefundGatewayResponse{
			Success: false,
			Message: fmt.Sprintf("Stripe refund failed: %v", err),
		}, err
	}

	success := r.Status == stripe.RefundStatusSucceeded

	return &RefundGatewayResponse{
		Success:  success,
		RefundID: r.ID,
		Message:  string(r.Status),
		Status:   string(r.Status),
	}, nil
}

// ValidatePaymentMethod validates a payment method
func (s *StripeService) ValidatePaymentMethod(ctx context.Context, paymentMethodID string) error {
	// In a real implementation, you would validate the payment method
	// For now, we'll just check if it's not empty
	if paymentMethodID == "" {
		return fmt.Errorf("payment method ID is required")
	}
	return nil
}

// GetPaymentStatus gets the status of a payment
func (s *StripeService) GetPaymentStatus(ctx context.Context, transactionID string) (string, error) {
	pi, err := paymentintent.Get(transactionID, nil)
	if err != nil {
		return "", err
	}
	return string(pi.Status), nil
}

// CreateCheckoutSession creates a Stripe Checkout Session for hosted payment page
func (s *StripeService) CreateCheckoutSession(ctx context.Context, req CheckoutSessionRequest) (*CheckoutSessionResponse, error) {
	// Convert amount to cents (Stripe uses smallest currency unit)
	amountCents := int64(req.Amount * 100)

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(req.Currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String("Order Payment"),
						Description: stripe.String(req.Description),
					},
					UnitAmount: stripe.Int64(amountCents),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(req.SuccessURL),
		CancelURL:  stripe.String(req.CancelURL),
	}

	// Add customer if provided
	if req.CustomerID != "" {
		params.Customer = stripe.String(req.CustomerID)
	}

	// Add metadata
	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}

	// Create checkout session
	sess, err := session.New(params)
	if err != nil {
		return &CheckoutSessionResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create checkout session: %v", err),
		}, err
	}

	return &CheckoutSessionResponse{
		Success:    true,
		SessionID:  sess.ID,
		SessionURL: sess.URL,
		Message:    "Checkout session created successfully",
	}, nil
}
