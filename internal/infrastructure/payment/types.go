package payment

import (
	"time"
)

// PaymentGatewayRequest represents a payment request to a gateway
type PaymentGatewayRequest struct {
	Amount          float64           `json:"amount"`
	Currency        string            `json:"currency"`
	PaymentMethodID string            `json:"payment_method_id,omitempty"`
	PaymentToken    string            `json:"payment_token,omitempty"`
	Description     string            `json:"description,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// PaymentGatewayResponse represents a payment response from a gateway
type PaymentGatewayResponse struct {
	Success       bool   `json:"success"`
	TransactionID string `json:"transaction_id"`
	ExternalID    string `json:"external_id"`
	Message       string `json:"message"`
	Status        string `json:"status"`
}

// RefundGatewayRequest represents a refund request to a gateway
type RefundGatewayRequest struct {
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Reason        string  `json:"reason,omitempty"`
}

// RefundGatewayResponse represents a refund response from a gateway
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
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// CheckoutSessionResponse represents a checkout session response
type CheckoutSessionResponse struct {
	Success    bool   `json:"success"`
	SessionID  string `json:"session_id"`
	SessionURL string `json:"session_url"`
	Message    string `json:"message"`
}

// WebhookEvent represents a webhook event from payment providers
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
}

// WebhookEventType represents different types of webhook events
type WebhookEventType string

const (
	WebhookEventPaymentSucceeded    WebhookEventType = "payment.succeeded"
	WebhookEventPaymentFailed       WebhookEventType = "payment.failed"
	WebhookEventCheckoutCompleted   WebhookEventType = "checkout.session.completed"
	WebhookEventRefundProcessed     WebhookEventType = "refund.processed"
	WebhookEventSubscriptionCreated WebhookEventType = "subscription.created"
	WebhookEventSubscriptionUpdated WebhookEventType = "subscription.updated"
	WebhookEventSubscriptionDeleted WebhookEventType = "subscription.deleted"
)

// PaymentProvider represents different payment providers
type PaymentProvider string

const (
	PaymentProviderStripe PaymentProvider = "stripe"
	PaymentProviderPayPal PaymentProvider = "paypal"
)
