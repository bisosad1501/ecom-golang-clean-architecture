package handlers

import (
	"net/http"
	"strconv"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/usecases"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OrderHandler handles order-related HTTP requests
type OrderHandler struct {
	orderUseCase usecases.OrderUseCase
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderUseCase usecases.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
	}
}

// CreateOrder handles creating a new order
// @Summary Create a new order
// @Description Create a new order from user's cart
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.CreateOrderRequest true "Create order request"
// @Success 201 {object} usecases.OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
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

	order, err := h.orderUseCase.CreateOrder(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Order created successfully",
		Data:    order,
	})
}

// GetOrder handles getting an order by ID
// @Summary Get order by ID
// @Description Get a single order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} usecases.OrderResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid order ID",
		})
		return
	}

	order, err := h.orderUseCase.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: order,
	})
}

// GetOrderPublic godoc
// @Summary Get order details (public access for success page)
// @Description Get order details without authentication for success page
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} usecases.OrderResponse
// @Failure 404 {object} ErrorResponse
// @Router /orders/{id}/public [get]
func (h *OrderHandler) GetOrderPublic(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid order ID",
		})
		return
	}

	order, err := h.orderUseCase.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Order not found",
		})
		return
	}

	// Return full order info for public access (since we have order ID, it's safe)
	// Users who have the order ID should be able to see all non-sensitive details
	orderResponse := map[string]interface{}{
		"id":               order.ID,
		"order_number":     order.OrderNumber,
		"status":           order.Status,
		"payment_status":   order.PaymentStatus,
		"subtotal":         order.Subtotal,
		"tax_amount":       order.TaxAmount,
		"shipping_amount":  order.ShippingAmount,
		"discount_amount":  order.DiscountAmount,
		"total":            order.Total,
		"currency":         order.Currency,
		"notes":            order.Notes,
		"item_count":       order.ItemCount,
		"can_be_cancelled": order.CanBeCancelled,
		"can_be_refunded":  order.CanBeRefunded,
		"created_at":       order.CreatedAt,
		"updated_at":       order.UpdatedAt,
		"items":            order.Items,
		"shipping_address": order.ShippingAddress,
		"billing_address":  order.BillingAddress,
		"user":             order.User,
		"payment":          order.Payment,
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Order retrieved successfully",
		Data:    orderResponse,
	})
}

// GetUserOrders handles getting user's orders
// @Summary Get user's orders
// @Description Get current user's order history
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} usecases.OrderResponse
// @Failure 401 {object} ErrorResponse
// @Router /orders [get]
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	orders, err := h.orderUseCase.GetUserOrders(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: orders,
	})
}

// CancelOrder handles cancelling an order
// @Summary Cancel order
// @Description Cancel an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} usecases.OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid order ID",
		})
		return
	}

	order, err := h.orderUseCase.CancelOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Order cancelled successfully",
		Data:    order,
	})
}

// GetOrders handles getting list of orders (admin only)
// @Summary Get orders list
// @Description Get list of all orders with filters (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Order status"
// @Param payment_status query string false "Payment status"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} usecases.OrderResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/orders [get]
func (h *OrderHandler) GetOrders(c *gin.Context) {
	req := usecases.GetOrdersRequest{
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
		Limit:     10,
		Offset:    0,
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := entities.OrderStatus(statusStr)
		req.Status = &status
	}

	if paymentStatusStr := c.Query("payment_status"); paymentStatusStr != "" {
		paymentStatus := entities.PaymentStatus(paymentStatusStr)
		req.PaymentStatus = &paymentStatus
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = offset
		}
	}

	orders, err := h.orderUseCase.GetOrders(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: orders,
	})
}

// UpdateOrderStatus handles updating order status (admin only)
// @Summary Update order status
// @Description Update the status of an order (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Param request body map[string]string true "Status update request"
// @Success 200 {object} usecases.OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/orders/{id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid order ID",
		})
		return
	}

	var req struct {
		Status string `json:"status" validate:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	status := entities.OrderStatus(req.Status)
	order, err := h.orderUseCase.UpdateOrderStatus(c.Request.Context(), orderID, status)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Order status updated successfully",
		Data:    order,
	})
}

// GetOrderBySessionID handles getting an order by session ID
// @Summary Get order by session ID
// @Description Get a single order by its checkout session ID
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id query string true "Session ID"
// @Success 200 {object} usecases.OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /orders/by-session [get]
func (h *OrderHandler) GetOrderBySessionID(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Session ID is required",
		})
		return
	}

	// Get user ID from token for authorization
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	order, err := h.orderUseCase.GetOrderBySessionID(c.Request.Context(), sessionID, userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: order,
	})
}
