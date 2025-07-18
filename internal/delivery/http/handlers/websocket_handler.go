package handlers

import (
	"net/http"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/infrastructure/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub *websocket.Hub
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// HandleNotificationWebSocket handles WebSocket connections for real-time notifications
func (h *WebSocketHandler) HandleNotificationWebSocket(c *gin.Context) {
	h.hub.HandleWebSocket(c)
}

// GetWebSocketStats returns WebSocket connection statistics
func (h *WebSocketHandler) GetWebSocketStats(c *gin.Context) {
	stats := h.hub.GetStats()
	
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "WebSocket statistics retrieved successfully",
		Data:    stats,
	})
}

// GetConnectedUsers returns list of connected users
func (h *WebSocketHandler) GetConnectedUsers(c *gin.Context) {
	users := h.hub.GetConnectedUsers()
	
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Connected users retrieved successfully",
		Data: map[string]interface{}{
			"connected_users": users,
			"count":          len(users),
		},
	})
}

// SendTestNotification sends a test notification to a specific user
func (h *WebSocketHandler) SendTestNotification(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Title    string                 `json:"title" binding:"required"`
		Message  string                 `json:"message" binding:"required"`
		Category string                 `json:"category"`
		Priority string                 `json:"priority"`
		Data     map[string]interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Create test notification (simplified)
	notification := &entities.Notification{
		ID:       uuid.New(),
		UserID:   &userID,
		Type:     entities.NotificationTypeInApp,
		Category: entities.NotificationCategory(req.Category),
		Priority: entities.NotificationPriority(req.Priority),
		Title:    req.Title,
		Message:  req.Message,
		Status:   entities.NotificationStatusSent,
	}

	// Send via WebSocket
	h.hub.SendToUser(userID, notification)

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Test notification sent successfully",
		Data:    notification,
	})
}

// BroadcastTestNotification broadcasts a test notification to all connected users
func (h *WebSocketHandler) BroadcastTestNotification(c *gin.Context) {
	var req struct {
		Title    string                 `json:"title" binding:"required"`
		Message  string                 `json:"message" binding:"required"`
		Category string                 `json:"category"`
		Priority string                 `json:"priority"`
		Data     map[string]interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Create test notification (simplified)
	notification := &entities.Notification{
		ID:       uuid.New(),
		Type:     entities.NotificationTypeInApp,
		Category: entities.NotificationCategory(req.Category),
		Priority: entities.NotificationPriority(req.Priority),
		Title:    req.Title,
		Message:  req.Message,
		Status:   entities.NotificationStatusSent,
	}

	// Broadcast via WebSocket
	h.hub.SendToAll(notification)

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Test notification broadcasted successfully",
		Data:    notification,
	})
}
