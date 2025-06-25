package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// PayPalService implements payment processing with PayPal
type PayPalService struct {
	clientID     string
	clientSecret string
	baseURL      string
	httpClient   *http.Client
}

// NewPayPalService creates a new PayPal service
func NewPayPalService(clientID, clientSecret string, sandbox bool) *PayPalService {
	baseURL := "https://api.paypal.com"
	if sandbox {
		baseURL = "https://api.sandbox.paypal.com"
	}

	return &PayPalService{
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      baseURL,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

// PayPalAccessTokenResponse represents PayPal access token response
type PayPalAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// PayPalPaymentRequest represents PayPal payment request
type PayPalPaymentRequest struct {
	Intent string `json:"intent"`
	Payer  struct {
		PaymentMethod string `json:"payment_method"`
	} `json:"payer"`
	Transactions []struct {
		Amount struct {
			Total    string `json:"total"`
			Currency string `json:"currency"`
		} `json:"amount"`
		Description string `json:"description"`
	} `json:"transactions"`
	RedirectURLs struct {
		ReturnURL string `json:"return_url"`
		CancelURL string `json:"cancel_url"`
	} `json:"redirect_urls"`
}

// PayPalPaymentResponse represents PayPal payment response
type PayPalPaymentResponse struct {
	ID    string `json:"id"`
	State string `json:"state"`
	Links []struct {
		Href   string `json:"href"`
		Rel    string `json:"rel"`
		Method string `json:"method"`
	} `json:"links"`
}

// getAccessToken gets an access token from PayPal
func (p *PayPalService) getAccessToken(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/v1/oauth2/token", p.baseURL)
	
	data := strings.NewReader("grant_type=client_credentials")
	req, err := http.NewRequestWithContext(ctx, "POST", url, data)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(p.clientID, p.clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("PayPal auth failed with status: %d", resp.StatusCode)
	}

	var tokenResp PayPalAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

// ProcessPayment processes a payment through PayPal
func (p *PayPalService) ProcessPayment(ctx context.Context, req PaymentGatewayRequest) (*PaymentGatewayResponse, error) {
	// Get access token
	token, err := p.getAccessToken(ctx)
	if err != nil {
		return &PaymentGatewayResponse{
			Success: false,
			Message: fmt.Sprintf("PayPal auth failed: %v", err),
		}, err
	}

	// Create payment request
	paymentReq := PayPalPaymentRequest{
		Intent: "sale",
	}
	paymentReq.Payer.PaymentMethod = "paypal"
	
	transaction := struct {
		Amount struct {
			Total    string `json:"total"`
			Currency string `json:"currency"`
		} `json:"amount"`
		Description string `json:"description"`
	}{}
	
	transaction.Amount.Total = fmt.Sprintf("%.2f", req.Amount)
	transaction.Amount.Currency = req.Currency
	transaction.Description = req.Description
	
	paymentReq.Transactions = []struct {
		Amount struct {
			Total    string `json:"total"`
			Currency string `json:"currency"`
		} `json:"amount"`
		Description string `json:"description"`
	}{transaction}

	// Set redirect URLs (these would be configurable in a real app)
	paymentReq.RedirectURLs.ReturnURL = "https://example.com/return"
	paymentReq.RedirectURLs.CancelURL = "https://example.com/cancel"

	// Make payment request
	jsonData, err := json.Marshal(paymentReq)
	if err != nil {
		return &PaymentGatewayResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to marshal payment request: %v", err),
		}, err
	}

	url := fmt.Sprintf("%s/v1/payments/payment", p.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return &PaymentGatewayResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create request: %v", err),
		}, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return &PaymentGatewayResponse{
			Success: false,
			Message: fmt.Sprintf("PayPal request failed: %v", err),
		}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &PaymentGatewayResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to read response: %v", err),
		}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return &PaymentGatewayResponse{
			Success: false,
			Message: fmt.Sprintf("PayPal payment failed with status: %d, body: %s", resp.StatusCode, string(body)),
		}, fmt.Errorf("PayPal payment failed")
	}

	var paymentResp PayPalPaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return &PaymentGatewayResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to parse response: %v", err),
		}, err
	}

	// For PayPal, the payment is created but not completed until user approves
	success := paymentResp.State == "created" || paymentResp.State == "approved"

	return &PaymentGatewayResponse{
		Success:       success,
		TransactionID: paymentResp.ID,
		ExternalID:    paymentResp.ID,
		Message:       paymentResp.State,
		Status:        paymentResp.State,
	}, nil
}

// ProcessRefund processes a refund through PayPal
func (p *PayPalService) ProcessRefund(ctx context.Context, req RefundGatewayRequest) (*RefundGatewayResponse, error) {
	// Get access token
	token, err := p.getAccessToken(ctx)
	if err != nil {
		return &RefundGatewayResponse{
			Success: false,
			Message: fmt.Sprintf("PayPal auth failed: %v", err),
		}, err
	}

	// Create refund request
	refundReq := map[string]interface{}{
		"amount": map[string]string{
			"total":    fmt.Sprintf("%.2f", req.Amount),
			"currency": "USD", // This should be configurable
		},
	}

	if req.Reason != "" {
		refundReq["reason"] = req.Reason
	}

	jsonData, err := json.Marshal(refundReq)
	if err != nil {
		return &RefundGatewayResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to marshal refund request: %v", err),
		}, err
	}

	// Note: This is a simplified refund URL - in reality you'd need the sale ID
	url := fmt.Sprintf("%s/v1/payments/sale/%s/refund", p.baseURL, req.TransactionID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return &RefundGatewayResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create request: %v", err),
		}, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return &RefundGatewayResponse{
			Success: false,
			Message: fmt.Sprintf("PayPal refund request failed: %v", err),
		}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RefundGatewayResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to read response: %v", err),
		}, err
	}

	success := resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK

	var refundResp map[string]interface{}
	if err := json.Unmarshal(body, &refundResp); err == nil {
		if refundID, ok := refundResp["id"].(string); ok {
			return &RefundGatewayResponse{
				Success:  success,
				RefundID: refundID,
				Message:  "Refund processed",
				Status:   "completed",
			}, nil
		}
	}

	return &RefundGatewayResponse{
		Success: success,
		Message: "Refund processed",
		Status:  "completed",
	}, nil
}

// CreateCheckoutSession creates a PayPal checkout session (placeholder implementation)
func (p *PayPalService) CreateCheckoutSession(ctx context.Context, req CheckoutSessionRequest) (*CheckoutSessionResponse, error) {
	// PayPal doesn't have the same "checkout session" concept as Stripe
	// This is a placeholder implementation that would redirect to PayPal's payment flow

	// In a real implementation, you would:
	// 1. Create a PayPal payment
	// 2. Return the approval URL for the user to complete payment

	return &CheckoutSessionResponse{
		Success:    false,
		Message:    "PayPal checkout sessions not implemented yet",
	}, fmt.Errorf("PayPal checkout sessions not implemented")
}
