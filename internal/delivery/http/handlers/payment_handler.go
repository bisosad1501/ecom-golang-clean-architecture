package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/usecases"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	paymentUseCase usecases.PaymentUseCase
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentUseCase usecases.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
	}
}

// ProcessPayment processes a payment
func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	// Check authentication
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Authentication required",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID format",
		})
		return
	}

	var req usecases.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Validate request
	if err := validateProcessPaymentRequest(&req, userID); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid payment request",
			Details: err.Error(),
		})
		return
	}

	payment, err := h.paymentUseCase.ProcessPayment(c.Request.Context(), req)
	if err != nil {
		statusCode := getPaymentErrorStatusCode(err)
		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to process payment",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Payment processed successfully",
		Data:    payment,
	})
}

// UpdatePaymentStatus updates payment status
func (h *PaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid payment ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Status        string `json:"status" binding:"required"`
		TransactionID string `json:"transaction_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Validate status
	validStatuses := []string{"pending", "processing", "paid", "failed", "cancelled", "refunded", "awaiting_payment"}
	isValidStatus := false
	for _, status := range validStatuses {
		if req.Status == status {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid payment status",
		})
		return
	}

	payment, err := h.paymentUseCase.UpdatePaymentStatus(c.Request.Context(), id, entities.PaymentStatus(req.Status), req.TransactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update payment status",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Payment status updated successfully",
		Data:    payment,
	})
}

// GetPayment retrieves a payment by ID
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid payment ID",
			Details: err.Error(),
		})
		return
	}

	payment, err := h.paymentUseCase.GetPayment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Payment not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Payment retrieved successfully",
		Data:    payment,
	})
}

// GetOrderPayments retrieves all payments for an order
func (h *PaymentHandler) GetOrderPayments(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid order ID",
			Details: err.Error(),
		})
		return
	}

	payments, err := h.paymentUseCase.GetOrderPayments(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get order payments",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Order payments retrieved successfully",
		Data:    payments,
	})
}

// ProcessRefund processes a refund
func (h *PaymentHandler) ProcessRefund(c *gin.Context) {
	var req usecases.ProcessRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	refund, err := h.paymentUseCase.ProcessRefund(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to process refund",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Refund processed successfully",
		Data:    refund,
	})
}

// GetRefunds retrieves refunds for a payment
func (h *PaymentHandler) GetRefunds(c *gin.Context) {
	paymentIDStr := c.Param("payment_id")
	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid payment ID",
			Details: err.Error(),
		})
		return
	}

	refunds, err := h.paymentUseCase.GetRefunds(c.Request.Context(), paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get refunds",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Refunds retrieved successfully",
		Data:    refunds,
	})
}

// ApproveRefund approves a pending refund
func (h *PaymentHandler) ApproveRefund(c *gin.Context) {
	refundIDStr := c.Param("refund_id")
	refundID, err := uuid.Parse(refundIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid refund ID",
			Details: err.Error(),
		})
		return
	}

	// Get user from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	approvedBy, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Invalid user ID format",
		})
		return
	}

	refund, err := h.paymentUseCase.ApproveRefund(c.Request.Context(), refundID, approvedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to approve refund",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Refund approved successfully",
		Data:    refund,
	})
}

// RejectRefund rejects a pending refund
func (h *PaymentHandler) RejectRefund(c *gin.Context) {
	refundIDStr := c.Param("refund_id")
	refundID, err := uuid.Parse(refundIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid refund ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Reason string `json:"reason" validate:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	err = h.paymentUseCase.RejectRefund(c.Request.Context(), refundID, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to reject refund",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Refund rejected successfully",
	})
}

// GetPendingRefunds retrieves refunds awaiting approval
func (h *PaymentHandler) GetPendingRefunds(c *gin.Context) {
	limit := 20
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	refunds, err := h.paymentUseCase.GetPendingRefunds(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get pending refunds",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Pending refunds retrieved successfully",
		Data:    refunds,
	})
}

// SavePaymentMethod saves a user's payment method
func (h *PaymentHandler) SavePaymentMethod(c *gin.Context) {
	// Get user ID from JWT token context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "User not authenticated",
			Details: "User ID not found in token",
		})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Invalid user ID format",
			Details: "User ID is not a valid UUID",
		})
		return
	}

	var req usecases.SavePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Set user ID from token
	req.UserID = userUUID

	method, err := h.paymentUseCase.SavePaymentMethod(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to save payment method",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Payment method saved successfully",
		Data:    method,
	})
}

// GetUserPaymentMethods retrieves user's payment methods
func (h *PaymentHandler) GetUserPaymentMethods(c *gin.Context) {
	// Get user ID from JWT token context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "User not authenticated",
			Details: "User ID not found in token",
		})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Invalid user ID format",
			Details: "User ID is not a valid UUID",
		})
		return
	}

	methods, err := h.paymentUseCase.GetUserPaymentMethods(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get payment methods",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Payment methods retrieved successfully",
		Data:    methods,
	})
}

// DeletePaymentMethod deletes a payment method
func (h *PaymentHandler) DeletePaymentMethod(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid payment method ID",
			Details: err.Error(),
		})
		return
	}

	if err := h.paymentUseCase.DeletePaymentMethod(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete payment method",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Payment method deleted successfully",
	})
}

// SetDefaultPaymentMethod sets a payment method as default
func (h *PaymentHandler) SetDefaultPaymentMethod(c *gin.Context) {
	// Get user ID from JWT token context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "User not authenticated",
			Details: "User ID not found in token",
		})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Invalid user ID format",
			Details: "User ID is not a valid UUID",
		})
		return
	}

	methodIDStr := c.Param("method_id")
	methodID, err := uuid.Parse(methodIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid method ID",
			Details: err.Error(),
		})
		return
	}

	if err := h.paymentUseCase.SetDefaultPaymentMethod(c.Request.Context(), userUUID, methodID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to set default payment method",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Default payment method set successfully",
	})
}

// HandleWebhook handles payment webhooks
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	provider := c.Param("provider")

	// Validate provider
	validProviders := []string{"stripe", "paypal"}
	isValidProvider := false
	for _, validProvider := range validProviders {
		if provider == validProvider {
			isValidProvider = true
			break
		}
	}
	if !isValidProvider {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid payment provider",
		})
		return
	}

	// Get signature header based on provider
	var signature string
	switch provider {
	case "stripe":
		signature = c.GetHeader("Stripe-Signature")
	case "paypal":
		signature = c.GetHeader("PAYPAL-TRANSMISSION-SIG")
	}

	// Validate signature is present
	if signature == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Missing webhook signature",
		})
		return
	}

	// Validate signature format based on provider
	switch provider {
	case "stripe":
		// Stripe signatures start with "t=" and contain "v1="
		if !strings.Contains(signature, "t=") || !strings.Contains(signature, "v1=") {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid Stripe signature format",
			})
			return
		}
	case "paypal":
		// PayPal signatures are base64 encoded
		if len(signature) < 20 || len(signature) > 500 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid PayPal signature format",
			})
			return
		}
	}

	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Failed to read webhook payload",
		})
		return
	}

	// Validate payload size
	if len(payload) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Empty webhook payload",
		})
		return
	}

	if len(payload) > 1024*1024 { // 1MB limit
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Webhook payload too large",
		})
		return
	}

	// Validate payload is valid JSON
	var jsonPayload map[string]interface{}
	if err := json.Unmarshal(payload, &jsonPayload); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid JSON payload",
		})
		return
	}

	// Process webhook
	if err := h.paymentUseCase.HandleWebhook(c.Request.Context(), provider, payload, signature); err != nil {
		// For security, don't expose internal error details
		// Log the actual error internally but return generic message
		// TODO: Add proper logging here
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Webhook processing failed",
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Webhook processed successfully",
	})
}

// GetPaymentReports returns payment reports
func (h *PaymentHandler) GetPaymentReports(c *gin.Context) {
	var req usecases.PaymentReportRequest

	// Parse query parameters
	if reportType := c.Query("report_type"); reportType != "" {
		req.ReportType = reportType
	}
	if groupBy := c.Query("group_by"); groupBy != "" {
		req.GroupBy = groupBy
	}
	if format := c.Query("format"); format != "" {
		req.Format = format
	}

	report, err := h.paymentUseCase.GetPaymentReport(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get payment reports",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Payment reports retrieved successfully",
		Data:    report,
	})
}

// CreateCheckoutSession creates a Stripe checkout session for hosted payment page
func (h *PaymentHandler) CreateCheckoutSession(c *gin.Context) {
	var req usecases.CreateCheckoutSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Create checkout session
	response, err := h.paymentUseCase.CreateCheckoutSession(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create checkout session",
			Details: err.Error(),
		})
		return
	}

	if !response.Success {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: response.Message,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Checkout session created successfully",
		Data:    response,
	})
}

// ConfirmPaymentSuccess confirms payment success for an order (fallback method)
func (h *PaymentHandler) ConfirmPaymentSuccess(c *gin.Context) {
	var req struct {
		SessionID string `json:"session_id" binding:"required"`
		OrderID   string `json:"order_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Parse order ID
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid order ID format",
		})
		return
	}

	// Get user ID from token (optional for public access)
	var userID uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		if parsedUserID, ok := userIDInterface.(uuid.UUID); ok {
			userID = parsedUserID
		}
	}

	// Confirm payment success
	err = h.paymentUseCase.ConfirmPaymentSuccess(c.Request.Context(), orderID, userID, req.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to confirm payment success",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Payment confirmation processed successfully",
	})
}

// validateProcessPaymentRequest validates process payment request
func validateProcessPaymentRequest(req *usecases.ProcessPaymentRequest, userID uuid.UUID) error {
	// Validate order ID
	if req.OrderID == uuid.Nil {
		return fmt.Errorf("order ID is required")
	}

	// Validate amount
	if req.Amount <= 0 {
		return fmt.Errorf("payment amount must be greater than 0")
	}

	if req.Amount > 999999.99 {
		return fmt.Errorf("payment amount cannot exceed $999,999.99")
	}

	// Validate currency
	if req.Currency == "" {
		req.Currency = "USD" // Set default
	} else if len(req.Currency) != 3 {
		return fmt.Errorf("currency must be a 3-letter ISO code")
	}

	// Validate payment method
	validMethods := map[entities.PaymentMethod]bool{
		entities.PaymentMethodCreditCard:   true,
		entities.PaymentMethodDebitCard:    true,
		entities.PaymentMethodPayPal:       true,
		entities.PaymentMethodStripe:       true,
		entities.PaymentMethodApplePay:     true,
		entities.PaymentMethodGooglePay:    true,
		entities.PaymentMethodBankTransfer: true,
		entities.PaymentMethodCash:         true,
	}

	if !validMethods[req.Method] {
		return fmt.Errorf("invalid payment method: %s", req.Method)
	}

	// Validate payment token for non-COD payments
	if req.Method != entities.PaymentMethodCash && req.PaymentToken == "" {
		return fmt.Errorf("payment token is required for %s payments", req.Method)
	}

	// Validate COD payments
	if req.Method == entities.PaymentMethodCash {
		if req.PaymentToken != "" {
			return fmt.Errorf("payment token should not be provided for COD payments")
		}
	}

	return nil
}

// getPaymentErrorStatusCode returns appropriate HTTP status code for payment errors
func getPaymentErrorStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	errorMsg := strings.ToLower(err.Error())

	// Authentication/Authorization errors
	if strings.Contains(errorMsg, "unauthorized") || strings.Contains(errorMsg, "access denied") {
		return http.StatusUnauthorized
	}

	// Validation errors
	if strings.Contains(errorMsg, "invalid") || strings.Contains(errorMsg, "validation") {
		return http.StatusBadRequest
	}

	// Not found errors
	if strings.Contains(errorMsg, "not found") {
		return http.StatusNotFound
	}

	// Payment specific errors
	if strings.Contains(errorMsg, "insufficient funds") || strings.Contains(errorMsg, "declined") {
		return http.StatusPaymentRequired
	}

	if strings.Contains(errorMsg, "duplicate") || strings.Contains(errorMsg, "already processed") {
		return http.StatusConflict
	}

	// Rate limiting
	if strings.Contains(errorMsg, "rate limit") || strings.Contains(errorMsg, "too many requests") {
		return http.StatusTooManyRequests
	}

	// Default to internal server error
	return http.StatusInternalServerError
}
