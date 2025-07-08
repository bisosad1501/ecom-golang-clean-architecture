package payment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/refund"
	"github.com/stripe/stripe-go/v76/webhook"
)

// StripeService implements payment processing with Stripe
type StripeService struct {
	apiKey        string
	webhookSecret string
}

// NewStripeService creates a new Stripe service
func NewStripeService(apiKey string) *StripeService {
	stripe.Key = apiKey
	return &StripeService{
		apiKey: apiKey,
	}
}

// NewStripeServiceWithWebhook creates a new Stripe service with webhook support
func NewStripeServiceWithWebhook(apiKey, webhookSecret string) *StripeService {
	stripe.Key = apiKey
	return &StripeService{
		apiKey:        apiKey,
		webhookSecret: webhookSecret,
	}
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

// HandleWebhook processes Stripe webhooks
func (s *StripeService) HandleWebhook(ctx context.Context, payload []byte, signature string) (*WebhookEvent, error) {
	// Verify webhook signature if webhook secret is configured
	var event stripe.Event
	var err error

	if s.webhookSecret != "" {
		event, err = webhook.ConstructEvent(payload, signature, s.webhookSecret)
		if err != nil {
			return nil, fmt.Errorf("webhook signature verification failed: %v", err)
		}
	} else {
		// For development/testing - parse without verification
		err = json.Unmarshal(payload, &event)
		if err != nil {
			return nil, fmt.Errorf("failed to parse webhook payload: %v", err)
		}
	}

	// Create webhook event response
	webhookEvent := &WebhookEvent{
		ID:   event.ID,
		Type: string(event.Type),
	}

	// Process different event types
	switch event.Type {
	case "checkout.session.completed":
		// Handle successful checkout session
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			return nil, fmt.Errorf("failed to parse checkout session: %v", err)
		}

		webhookEvent.Data = map[string]interface{}{
			"session_id":     session.ID,
			"payment_status": session.PaymentStatus,
			"amount_total":   session.AmountTotal,
			"currency":       session.Currency,
			"metadata":       session.Metadata,
		}

	case "payment_intent.succeeded":
		// Handle successful payment intent
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			return nil, fmt.Errorf("failed to parse payment intent: %v", err)
		}

		webhookEvent.Data = map[string]interface{}{
			"payment_intent_id": paymentIntent.ID,
			"amount":            paymentIntent.Amount,
			"currency":          paymentIntent.Currency,
			"status":            paymentIntent.Status,
			"metadata":          paymentIntent.Metadata,
		}

	case "payment_intent.payment_failed":
		// Handle failed payment intent
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			return nil, fmt.Errorf("failed to parse payment intent: %v", err)
		}

		webhookEvent.Data = map[string]interface{}{
			"payment_intent_id":  paymentIntent.ID,
			"amount":             paymentIntent.Amount,
			"currency":           paymentIntent.Currency,
			"status":             paymentIntent.Status,
			"last_payment_error": paymentIntent.LastPaymentError,
			"metadata":           paymentIntent.Metadata,
		}

	default:
		// For other event types, just store the raw data
		webhookEvent.Data = map[string]interface{}{
			"raw": event.Data.Raw,
		}
	}

	return webhookEvent, nil
}
