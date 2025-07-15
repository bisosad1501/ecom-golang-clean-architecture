package handlers

import (
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
		Data: notification,
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
		Data: notification,
	})
}

// GetUserNotifications gets notifications for current user
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
		Error: "User not authenticated",
		Details: "",
	})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "15"))

	// Validate and normalize pagination for notifications
	page, limit, err := usecases.ValidateAndNormalizePaginationForEntity(page, limit, "notifications")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	req := usecases.GetUserNotificationsRequest{
		Limit:  limit,
		Offset: (page - 1) * limit,
	}

	response, err := h.notificationUseCase.GetUserNotifications(c.Request.Context(), userID.(uuid.UUID), req)
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
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
		Error: "User not authenticated",
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

	err = h.notificationUseCase.MarkAsRead(c.Request.Context(), userID.(uuid.UUID), notificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to mark notification as read",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Notification marked as read",
		Data: nil,
	})
}

// MarkAllAsRead marks all notifications as read for user
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
		Error: "User not authenticated",
		Details: "",
	})
		return
	}

	err := h.notificationUseCase.MarkAllAsRead(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to mark all notifications as read",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "All notifications marked as read",
		Data: nil,
	})
}

// GetUnreadCount gets unread notification count
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
		Error: "User not authenticated",
		Details: "",
	})
		return
	}

	count, err := h.notificationUseCase.GetUnreadCount(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get unread count",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Unread count retrieved successfully",
		Data: gin.H{"count": count},
	})
}

// CreateTemplate creates a notification template
func (h *NotificationHandler) CreateTemplate(c *gin.Context) {
	var req usecases.CreateTemplateRequest
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
		Data: template,
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
		Data: templates,
	})
}
