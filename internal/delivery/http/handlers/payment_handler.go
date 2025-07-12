package handlers

import (
	"net/http"

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
	var req usecases.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	payment, err := h.paymentUseCase.ProcessPayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
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
	orderIDStr := c.Param("order_id")
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

// SavePaymentMethod saves a user's payment method
func (h *PaymentHandler) SavePaymentMethod(c *gin.Context) {
	var req usecases.SavePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

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
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Details: err.Error(),
		})
		return
	}

	methods, err := h.paymentUseCase.GetUserPaymentMethods(c.Request.Context(), userID)
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
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Details: err.Error(),
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

	if err := h.paymentUseCase.SetDefaultPaymentMethod(c.Request.Context(), userID, methodID); err != nil {
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

	// Validate signature format
	if len(signature) < 10 || len(signature) > 1000 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid signature format",
		})
		return
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
	if userIDStr, exists := c.Get("user_id"); exists {
		if parsedUserID, parseErr := uuid.Parse(userIDStr.(string)); parseErr == nil {
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
