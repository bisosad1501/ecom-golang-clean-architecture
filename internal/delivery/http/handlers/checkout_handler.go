package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ecom-golang-clean-architecture/internal/usecases"
)

// CheckoutHandler handles checkout-related HTTP requests
type CheckoutHandler struct {
	checkoutUseCase usecases.CheckoutUseCase
}

// NewCheckoutHandler creates a new checkout handler
func NewCheckoutHandler(checkoutUseCase usecases.CheckoutUseCase) *CheckoutHandler {
	return &CheckoutHandler{
		checkoutUseCase: checkoutUseCase,
	}
}

// CreateCheckoutSession handles creating a checkout session for online payments
// @Summary Create checkout session
// @Description Create a checkout session for online payments (credit card, PayPal, etc.)
// @Tags checkout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.CreateNewCheckoutSessionRequest true "Create checkout session request"
// @Success 201 {object} usecases.NewCheckoutSessionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /checkout/session [post]
func (h *CheckoutHandler) CreateCheckoutSession(c *gin.Context) {
	// Debug logging
	authHeader := c.GetHeader("Authorization")
	fmt.Printf("🔍 CreateCheckoutSession - Auth header: %s\n", authHeader)

	userIDInterface, exists := c.Get("user_id")
	fmt.Printf("🔍 CreateCheckoutSession - User ID exists: %v\n", exists)
	if exists {
		fmt.Printf("🔍 CreateCheckoutSession - User ID: %v\n", userIDInterface)
	}

	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
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

	var req usecases.CreateNewCheckoutSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Validate request fields
	if err := validateCreateCheckoutSessionRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request data",
			Details: err.Error(),
		})
		return
	}

	session, err := h.checkoutUseCase.CreateCheckoutSession(c.Request.Context(), userID, req)
	if err != nil {
		statusCode := getErrorStatusCode(err)
		c.JSON(statusCode, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Checkout session created successfully",
		Data:    session,
	})
}

// GetCheckoutSession handles getting a checkout session
// @Summary Get checkout session
// @Description Get checkout session by session ID
// @Tags checkout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path string true "Session ID"
// @Success 200 {object} usecases.NewCheckoutSessionResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /checkout/session/{session_id} [get]
func (h *CheckoutHandler) GetCheckoutSession(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Session ID is required",
		})
		return
	}

	session, err := h.checkoutUseCase.GetCheckoutSession(c.Request.Context(), sessionID)
	if err != nil {
		statusCode := getErrorStatusCode(err)
		c.JSON(statusCode, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Checkout session retrieved successfully",
		Data:    session,
	})
}

// CompleteCheckoutSession handles completing a checkout session after payment
// @Summary Complete checkout session
// @Description Complete checkout session after successful payment
// @Tags checkout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path string true "Session ID"
// @Success 200 {object} usecases.OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /checkout/session/{session_id}/complete [post]
func (h *CheckoutHandler) CompleteCheckoutSession(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Session ID is required",
		})
		return
	}

	order, err := h.checkoutUseCase.CompleteCheckoutSession(c.Request.Context(), sessionID)
	if err != nil {
		statusCode := getErrorStatusCode(err)
		c.JSON(statusCode, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Order created successfully",
		Data:    order,
	})
}

// CancelCheckoutSession handles cancelling a checkout session
// @Summary Cancel checkout session
// @Description Cancel an active checkout session
// @Tags checkout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path string true "Session ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /checkout/session/{session_id}/cancel [post]
func (h *CheckoutHandler) CancelCheckoutSession(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Session ID is required",
		})
		return
	}

	err := h.checkoutUseCase.CancelCheckoutSession(c.Request.Context(), sessionID)
	if err != nil {
		statusCode := getErrorStatusCode(err)
		c.JSON(statusCode, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Checkout session cancelled successfully",
	})
}

// CreateCODOrder handles creating COD orders directly
// @Summary Create COD order
// @Description Create order directly for Cash on Delivery payments
// @Tags checkout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.CreateOrderRequest true "Create COD order request"
// @Success 201 {object} usecases.OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /checkout/cod [post]
func (h *CheckoutHandler) CreateCODOrder(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
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

	var req usecases.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Validate that this is a COD request
	if req.PaymentMethod != "cash" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "This endpoint is only for COD orders",
		})
		return
	}

	// Validate request fields for COD
	if err := validateCreateCODOrderRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request data",
			Details: err.Error(),
		})
		return
	}

	order, err := h.checkoutUseCase.CreateCODOrder(c.Request.Context(), userID, req)
	if err != nil {
		statusCode := getErrorStatusCode(err)
		c.JSON(statusCode, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "COD order created successfully",
		Data:    order,
	})
}

// validateCreateCheckoutSessionRequest validates create checkout session request
func validateCreateCheckoutSessionRequest(req *usecases.CreateNewCheckoutSessionRequest) error {
	// Validate payment method (exclude COD)
	validPaymentMethods := map[string]bool{
		"credit_card":   true,
		"debit_card":    true,
		"paypal":        true,
		"stripe":        true,
		"apple_pay":     true,
		"google_pay":    true,
		"bank_transfer": true,
	}

	if !validPaymentMethods[string(req.PaymentMethod)] {
		return fmt.Errorf("invalid payment method for checkout session: %s", req.PaymentMethod)
	}

	// Validate financial amounts
	if req.TaxRate < 0 || req.TaxRate > 1 {
		return fmt.Errorf("tax rate must be between 0 and 1, got: %.4f", req.TaxRate)
	}

	if req.ShippingCost < 0 {
		return fmt.Errorf("shipping cost cannot be negative, got: %.2f", req.ShippingCost)
	}

	if req.DiscountAmount < 0 {
		return fmt.Errorf("discount amount cannot be negative, got: %.2f", req.DiscountAmount)
	}

	// Validate shipping address (required)
	if req.ShippingAddress.FirstName == "" {
		return fmt.Errorf("shipping address first name is required")
	}
	if req.ShippingAddress.LastName == "" {
		return fmt.Errorf("shipping address last name is required")
	}
	if req.ShippingAddress.Address1 == "" {
		return fmt.Errorf("shipping address line 1 is required")
	}
	if req.ShippingAddress.City == "" {
		return fmt.Errorf("shipping address city is required")
	}
	if req.ShippingAddress.Country == "" {
		return fmt.Errorf("shipping address country is required")
	}

	// Validate country code (should be exactly 2 characters for ISO codes)
	if len(req.ShippingAddress.Country) != 2 {
		return fmt.Errorf("shipping address country must be a 2-letter ISO country code")
	}

	return nil
}

// validateCreateCODOrderRequest validates create COD order request
func validateCreateCODOrderRequest(req *usecases.CreateOrderRequest) error {
	// Only allow cash payment method for COD endpoint
	if req.PaymentMethod != "cash" {
		return fmt.Errorf("only cash payment method is allowed for COD orders: %s", req.PaymentMethod)
	}

	// Validate financial amounts
	if req.TaxRate < 0 || req.TaxRate > 1 {
		return fmt.Errorf("tax rate must be between 0 and 1, got: %.4f", req.TaxRate)
	}

	if req.ShippingCost < 0 {
		return fmt.Errorf("shipping cost cannot be negative, got: %.2f", req.ShippingCost)
	}

	if req.DiscountAmount < 0 {
		return fmt.Errorf("discount amount cannot be negative, got: %.2f", req.DiscountAmount)
	}

	// Validate shipping address (required)
	if req.ShippingAddress.FirstName == "" {
		return fmt.Errorf("shipping address first name is required")
	}
	if req.ShippingAddress.LastName == "" {
		return fmt.Errorf("shipping address last name is required")
	}
	if req.ShippingAddress.Address1 == "" {
		return fmt.Errorf("shipping address line 1 is required")
	}
	if req.ShippingAddress.City == "" {
		return fmt.Errorf("shipping address city is required")
	}
	if req.ShippingAddress.Country == "" {
		return fmt.Errorf("shipping address country is required")
	}

	// Validate country code (should be exactly 2 characters for ISO codes)
	if len(req.ShippingAddress.Country) != 2 {
		return fmt.Errorf("shipping address country must be a 2-letter ISO country code")
	}

	return nil
}
