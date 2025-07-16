package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"ecom-golang-clean-architecture/internal/usecases"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	notificationUseCase usecases.NotificationUseCase
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationUseCase usecases.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{
		notificationUseCase: notificationUseCase,
	}
}

// getUserID extracts and validates user ID from context
func (h *NotificationHandler) getUserID(c *gin.Context) (uuid.UUID, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("user not authenticated")
	}

	// Try to convert to UUID directly first
	if userID, ok := userIDInterface.(uuid.UUID); ok {
		return userID, nil
	}

	// If not UUID, try to parse as string
	if userIDStr, ok := userIDInterface.(string); ok {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
		}
		return userID, nil
	}

	return uuid.Nil, fmt.Errorf("user ID has invalid type")
}

// CreateNotification creates a new notification
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var req usecases.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	notification, err := h.notificationUseCase.CreateNotification(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create notification",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Notification created successfully",
		Data:    notification,
	})
}

// GetNotification gets a notification by ID
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid notification ID",
			Details: err.Error(),
		})
		return
	}

	notification, err := h.notificationUseCase.GetNotification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Notification not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Notification retrieved successfully",
		Data:    notification,
	})
}

// GetUserNotifications gets notifications for current user
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   err.Error(),
			Details: "",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "15"))

	// Validate and normalize pagination for notifications
	page, limit, err2 := usecases.ValidateAndNormalizePaginationForEntity(page, limit, "notifications")
	if err2 != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err2.Error(),
		})
		return
	}

	req := usecases.GetUserNotificationsRequest{
		Limit:  limit,
		Offset: (page - 1) * limit,
	}

	response, err := h.notificationUseCase.GetUserNotifications(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get notifications",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       response.Notifications,
		Pagination: response.Pagination,
	})
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   err.Error(),
			Details: "",
		})
		return
	}

	idStr := c.Param("id")
	notificationID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid notification ID",
			Details: err.Error(),
		})
		return
	}

	err = h.notificationUseCase.MarkAsRead(c.Request.Context(), userID, notificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to mark notification as read",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Notification marked as read",
		Data:    nil,
	})
}

// MarkAllAsRead marks all notifications as read for user
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   err.Error(),
			Details: "",
		})
		return
	}

	err = h.notificationUseCase.MarkAllAsRead(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to mark all notifications as read",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "All notifications marked as read",
		Data:    nil,
	})
}

// GetUnreadCount gets unread notification count
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   err.Error(),
			Details: "",
		})
		return
	}

	count, err := h.notificationUseCase.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get unread count",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Unread count retrieved successfully",
		Data:    gin.H{"count": count},
	})
}

// CreateTemplate creates a notification template
func (h *NotificationHandler) CreateTemplate(c *gin.Context) {
	var req usecases.CreateNotificationTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	template, err := h.notificationUseCase.CreateTemplate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create template",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Template created successfully",
		Data:    template,
	})
}

// GetTemplates gets notification templates
func (h *NotificationHandler) GetTemplates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := usecases.ListTemplatesRequest{
		Limit:  limit,
		Offset: (page - 1) * limit,
	}

	templates, err := h.notificationUseCase.ListTemplates(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get templates",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Templates retrieved successfully",
		Data:    templates,
	})
}

// GetUserPreferences gets user notification preferences
func (h *NotificationHandler) GetUserPreferences(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   err.Error(),
			Details: "",
		})
		return
	}

	preferences, err := h.notificationUseCase.GetUserPreferences(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get user preferences",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "User preferences retrieved successfully",
		Data:    preferences,
	})
}

// UpdateUserPreferences updates user notification preferences
func (h *NotificationHandler) UpdateUserPreferences(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   err.Error(),
			Details: "",
		})
		return
	}

	var req struct {
		InAppEnabled         bool `json:"in_app_enabled"`
		EmailEnabled         bool `json:"email_enabled"`
		SMSEnabled           bool `json:"sms_enabled"`
		PushEnabled          bool `json:"push_enabled"`
		EmailOrderUpdates    bool `json:"email_order_updates"`
		EmailPaymentUpdates  bool `json:"email_payment_updates"`
		EmailShippingUpdates bool `json:"email_shipping_updates"`
		EmailPromotions      bool `json:"email_promotions"`
		EmailNewsletter      bool `json:"email_newsletter"`
		SMSOrderUpdates      bool `json:"sms_order_updates"`
		SMSPaymentUpdates    bool `json:"sms_payment_updates"`
		SMSShippingUpdates   bool `json:"sms_shipping_updates"`
		SMSSecurityAlerts    bool `json:"sms_security_alerts"`
		PushOrderUpdates     bool `json:"push_order_updates"`
		PushPaymentUpdates   bool `json:"push_payment_updates"`
		PushShippingUpdates  bool `json:"push_shipping_updates"`
		PushPromotions       bool `json:"push_promotions"`
		InAppOrderUpdates    bool `json:"in_app_order_updates"`
		InAppPaymentUpdates  bool `json:"in_app_payment_updates"`
		InAppShippingUpdates bool `json:"in_app_shipping_updates"`
		InAppPromotions      bool `json:"in_app_promotions"`
		InAppSystemUpdates   bool `json:"in_app_system_updates"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Create update request
	updateReq := usecases.UpdatePreferencesRequest{
		InAppEnabled:         &req.InAppEnabled,
		EmailEnabled:         &req.EmailEnabled,
		SMSEnabled:           &req.SMSEnabled,
		PushEnabled:          &req.PushEnabled,
		EmailOrderUpdates:    &req.EmailOrderUpdates,
		EmailPaymentUpdates:  &req.EmailPaymentUpdates,
		EmailShippingUpdates: &req.EmailShippingUpdates,
		EmailPromotions:      &req.EmailPromotions,
		EmailNewsletter:      &req.EmailNewsletter,
		SMSOrderUpdates:      &req.SMSOrderUpdates,
		SMSPaymentUpdates:    &req.SMSPaymentUpdates,
		SMSShippingUpdates:   &req.SMSShippingUpdates,
		SMSSecurityAlerts:    &req.SMSSecurityAlerts,
		PushOrderUpdates:     &req.PushOrderUpdates,
		PushPaymentUpdates:   &req.PushPaymentUpdates,
		PushShippingUpdates:  &req.PushShippingUpdates,
		PushPromotions:       &req.PushPromotions,
		InAppOrderUpdates:    &req.InAppOrderUpdates,
		InAppPaymentUpdates:  &req.InAppPaymentUpdates,
		InAppShippingUpdates: &req.InAppShippingUpdates,
		InAppPromotions:      &req.InAppPromotions,
		InAppSystemUpdates:   &req.InAppSystemUpdates,
	}

	preferences, err := h.notificationUseCase.UpdateUserPreferences(c.Request.Context(), userID, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update user preferences",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "User preferences updated successfully",
		Data:    preferences,
	})
}
