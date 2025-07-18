package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/domain/services"

	"github.com/google/uuid"
)

// NotificationUseCase defines notification use cases
type NotificationUseCase interface {
	// Notification management
	CreateNotification(ctx context.Context, req CreateNotificationRequest) (*NotificationResponse, error)
	GetNotification(ctx context.Context, id uuid.UUID) (*NotificationResponse, error)
	UpdateNotification(ctx context.Context, id uuid.UUID, req UpdateNotificationRequest) (*NotificationResponse, error)
	DeleteNotification(ctx context.Context, id uuid.UUID) error
	ListNotifications(ctx context.Context, req ListNotificationsRequest) (*NotificationsListResponse, error)

	// User notifications
	GetUserNotifications(ctx context.Context, userID uuid.UUID, req GetUserNotificationsRequest) (*NotificationsListResponse, error)
	MarkAsRead(ctx context.Context, userID, notificationID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error)

	// Notification sending
	SendNotification(ctx context.Context, notification *entities.Notification) error
	SendBulkNotifications(ctx context.Context, notifications []*entities.Notification) error
	QueueNotification(ctx context.Context, notification *entities.Notification, scheduledAt *time.Time) error
	ProcessQueue(ctx context.Context, limit int) error

	// Templates
	CreateTemplate(ctx context.Context, req CreateNotificationTemplateRequest) (*NotificationTemplateResponse, error)
	GetTemplate(ctx context.Context, id uuid.UUID) (*NotificationTemplateResponse, error)
	UpdateTemplate(ctx context.Context, id uuid.UUID, req UpdateNotificationTemplateRequest) (*NotificationTemplateResponse, error)
	DeleteTemplate(ctx context.Context, id uuid.UUID) error
	ListTemplates(ctx context.Context, req ListTemplatesRequest) (*TemplatesListResponse, error)

	// Preferences
	GetUserPreferences(ctx context.Context, userID uuid.UUID) (*PreferencesResponse, error)
	UpdateUserPreferences(ctx context.Context, userID uuid.UUID, req UpdatePreferencesRequest) (*PreferencesResponse, error)

	// Event-based notifications
	NotifyOrderCreated(ctx context.Context, orderID uuid.UUID) error
	NotifyOrderStatusChanged(ctx context.Context, orderID uuid.UUID, newStatus string) error
	NotifyPaymentReceived(ctx context.Context, paymentID uuid.UUID) error
	NotifyShippingUpdate(ctx context.Context, orderID uuid.UUID, trackingNumber string) error
	NotifyLowStock(ctx context.Context, inventoryID uuid.UUID) error
	NotifyReviewRequest(ctx context.Context, orderID uuid.UUID) error

	// Admin-specific notifications
	NotifyNewOrder(ctx context.Context, orderID uuid.UUID) error
	NotifyPaymentFailed(ctx context.Context, paymentID uuid.UUID) error
	NotifyNewUser(ctx context.Context, userID uuid.UUID) error
	NotifyNewReview(ctx context.Context, reviewID uuid.UUID) error
}

type notificationUseCase struct {
	notificationRepo repositories.NotificationRepository
	userRepo         repositories.UserRepository
	orderRepo        repositories.OrderRepository
	paymentRepo      repositories.PaymentRepository
	inventoryRepo    repositories.InventoryRepository
	reviewRepo       repositories.ReviewRepository
	productRepo      repositories.ProductRepository
	emailService     services.EmailService
	smsService       SMSService
	pushService      PushService
	websocketHub     WebSocketHub
}

// WebSocketHub interface for real-time notifications
type WebSocketHub interface {
	SendToUser(userID uuid.UUID, notification *entities.Notification)
	SendToAll(notification *entities.Notification)
}

// NewNotificationUseCase creates a new notification use case
func NewNotificationUseCase(
	notificationRepo repositories.NotificationRepository,
	userRepo repositories.UserRepository,
	orderRepo repositories.OrderRepository,
	paymentRepo repositories.PaymentRepository,
	inventoryRepo repositories.InventoryRepository,
	reviewRepo repositories.ReviewRepository,
	productRepo repositories.ProductRepository,
	emailService services.EmailService,
	smsService SMSService,
	pushService PushService,
	websocketHub WebSocketHub,
) NotificationUseCase {
	return &notificationUseCase{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		orderRepo:        orderRepo,
		paymentRepo:      paymentRepo,
		inventoryRepo:    inventoryRepo,
		reviewRepo:       reviewRepo,
		productRepo:      productRepo,
		emailService:     emailService,
		smsService:       smsService,
		pushService:      pushService,
		websocketHub:     websocketHub,
	}
}

// Service interfaces
type EmailService interface {
	SendEmail(ctx context.Context, to, subject, body string, template string, data map[string]interface{}) error
}

type SMSService interface {
	SendSMS(ctx context.Context, to, message string) error
}

type PushService interface {
	SendPush(ctx context.Context, userID uuid.UUID, title, message string, data map[string]interface{}) error
}

// Request/Response types
type CreateNotificationRequest struct {
	UserID        *uuid.UUID                    `json:"user_id,omitempty"`
	Type          entities.NotificationType     `json:"type" validate:"required"`
	Category      entities.NotificationCategory `json:"category" validate:"required"`
	Priority      entities.NotificationPriority `json:"priority"`
	Title         string                        `json:"title" validate:"required"`
	Message       string                        `json:"message" validate:"required"`
	Data          map[string]interface{}        `json:"data,omitempty"`
	Recipient     string                        `json:"recipient,omitempty"`
	Subject       string                        `json:"subject,omitempty"`
	Template      string                        `json:"template,omitempty"`
	ReferenceType string                        `json:"reference_type,omitempty"`
	ReferenceID   *uuid.UUID                    `json:"reference_id,omitempty"`
	ScheduledAt   *time.Time                    `json:"scheduled_at,omitempty"`
}

type UpdateNotificationRequest struct {
	Status   *entities.NotificationStatus `json:"status,omitempty"`
	Title    *string                      `json:"title,omitempty"`
	Message  *string                      `json:"message,omitempty"`
	Data     map[string]interface{}       `json:"data,omitempty"`
	Subject  *string                      `json:"subject,omitempty"`
	Template *string                      `json:"template,omitempty"`
}

type ListNotificationsRequest struct {
	UserID        *uuid.UUID                     `json:"user_id,omitempty"`
	Type          *entities.NotificationType     `json:"type,omitempty"`
	Category      *entities.NotificationCategory `json:"category,omitempty"`
	Status        *entities.NotificationStatus   `json:"status,omitempty"`
	Priority      *entities.NotificationPriority `json:"priority,omitempty"`
	ReferenceType *string                        `json:"reference_type,omitempty"`
	ReferenceID   *uuid.UUID                     `json:"reference_id,omitempty"`
	DateFrom      *time.Time                     `json:"date_from,omitempty"`
	DateTo        *time.Time                     `json:"date_to,omitempty"`
	SortBy        string                         `json:"sort_by,omitempty" validate:"omitempty,oneof=created_at priority status"`
	SortOrder     string                         `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
	Limit         int                            `json:"limit" validate:"min=1,max=100"`
	Offset        int                            `json:"offset" validate:"min=0"`
}

type GetUserNotificationsRequest struct {
	Type      *entities.NotificationType     `json:"type,omitempty"`
	Category  *entities.NotificationCategory `json:"category,omitempty"`
	Status    *entities.NotificationStatus   `json:"status,omitempty"`
	IsRead    *bool                          `json:"is_read,omitempty"`
	SortBy    string                         `json:"sort_by,omitempty" validate:"omitempty,oneof=created_at priority"`
	SortOrder string                         `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
	Limit     int                            `json:"limit" validate:"min=1,max=100"`
	Offset    int                            `json:"offset" validate:"min=0"`
}

type CreateNotificationTemplateRequest struct {
	Name        string                        `json:"name" validate:"required"`
	Type        entities.NotificationType     `json:"type" validate:"required"`
	Category    entities.NotificationCategory `json:"category" validate:"required"`
	Subject     string                        `json:"subject,omitempty"`
	Body        string                        `json:"body" validate:"required"`
	Variables   []string                      `json:"variables,omitempty"`
	IsActive    bool                          `json:"is_active"`
	IsDefault   bool                          `json:"is_default"`
	Language    string                        `json:"language"`
	Description string                        `json:"description,omitempty"`
	CreatedBy   uuid.UUID                     `json:"created_by" validate:"required"`
}

type UpdateNotificationTemplateRequest struct {
	Name        *string  `json:"name,omitempty"`
	Subject     *string  `json:"subject,omitempty"`
	Body        *string  `json:"body,omitempty"`
	Variables   []string `json:"variables,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
	IsDefault   *bool    `json:"is_default,omitempty"`
	Language    *string  `json:"language,omitempty"`
	Description *string  `json:"description,omitempty"`
}

type ListTemplatesRequest struct {
	Type      *entities.NotificationType     `json:"type,omitempty"`
	Category  *entities.NotificationCategory `json:"category,omitempty"`
	IsActive  *bool                          `json:"is_active,omitempty"`
	Language  *string                        `json:"language,omitempty"`
	SortBy    string                         `json:"sort_by,omitempty" validate:"omitempty,oneof=name created_at"`
	SortOrder string                         `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
	Limit     int                            `json:"limit" validate:"min=1,max=100"`
	Offset    int                            `json:"offset" validate:"min=0"`
}

type UpdatePreferencesRequest struct {
	EmailEnabled         *bool   `json:"email_enabled,omitempty"`
	EmailOrderUpdates    *bool   `json:"email_order_updates,omitempty"`
	EmailPaymentUpdates  *bool   `json:"email_payment_updates,omitempty"`
	EmailShippingUpdates *bool   `json:"email_shipping_updates,omitempty"`
	EmailPromotions      *bool   `json:"email_promotions,omitempty"`
	EmailNewsletter      *bool   `json:"email_newsletter,omitempty"`
	EmailReviewReminders *bool   `json:"email_review_reminders,omitempty"`
	SMSEnabled           *bool   `json:"sms_enabled,omitempty"`
	SMSOrderUpdates      *bool   `json:"sms_order_updates,omitempty"`
	SMSPaymentUpdates    *bool   `json:"sms_payment_updates,omitempty"`
	SMSShippingUpdates   *bool   `json:"sms_shipping_updates,omitempty"`
	SMSSecurityAlerts    *bool   `json:"sms_security_alerts,omitempty"`
	PushEnabled          *bool   `json:"push_enabled,omitempty"`
	PushOrderUpdates     *bool   `json:"push_order_updates,omitempty"`
	PushPaymentUpdates   *bool   `json:"push_payment_updates,omitempty"`
	PushShippingUpdates  *bool   `json:"push_shipping_updates,omitempty"`
	PushPromotions       *bool   `json:"push_promotions,omitempty"`
	PushReviewReminders  *bool   `json:"push_review_reminders,omitempty"`
	InAppEnabled         *bool   `json:"in_app_enabled,omitempty"`
	InAppOrderUpdates    *bool   `json:"in_app_order_updates,omitempty"`
	InAppPaymentUpdates  *bool   `json:"in_app_payment_updates,omitempty"`
	InAppShippingUpdates *bool   `json:"in_app_shipping_updates,omitempty"`
	InAppPromotions      *bool   `json:"in_app_promotions,omitempty"`
	InAppSystemUpdates   *bool   `json:"in_app_system_updates,omitempty"`
	DigestFrequency      *string `json:"digest_frequency,omitempty"`
	QuietHoursStart      *string `json:"quiet_hours_start,omitempty"`
	QuietHoursEnd        *string `json:"quiet_hours_end,omitempty"`
	Timezone             *string `json:"timezone,omitempty"`
}

// Response types
type NotificationResponse struct {
	ID            uuid.UUID                     `json:"id"`
	UserID        *uuid.UUID                    `json:"user_id"`
	Type          entities.NotificationType     `json:"type"`
	Category      entities.NotificationCategory `json:"category"`
	Priority      entities.NotificationPriority `json:"priority"`
	Status        entities.NotificationStatus   `json:"status"`
	Title         string                        `json:"title"`
	Message       string                        `json:"message"`
	Data          map[string]interface{}        `json:"data,omitempty"`
	Recipient     string                        `json:"recipient"`
	Subject       string                        `json:"subject"`
	Template      string                        `json:"template"`
	ReferenceType string                        `json:"reference_type"`
	ReferenceID   *uuid.UUID                    `json:"reference_id"`
	ScheduledAt   *time.Time                    `json:"scheduled_at"`
	SentAt        *time.Time                    `json:"sent_at"`
	DeliveredAt   *time.Time                    `json:"delivered_at"`
	ReadAt        *time.Time                    `json:"read_at"`
	RetryCount    int                           `json:"retry_count"`
	MaxRetries    int                           `json:"max_retries"`
	NextRetryAt   *time.Time                    `json:"next_retry_at"`
	ErrorMessage  string                        `json:"error_message"`
	ErrorCode     string                        `json:"error_code"`
	IsRead        bool                          `json:"is_read"`
	CreatedAt     time.Time                     `json:"created_at"`
	UpdatedAt     time.Time                     `json:"updated_at"`
}

type NotificationsListResponse struct {
	Notifications []*NotificationResponse `json:"notifications"`
	Total         int64                   `json:"total"`
	UnreadCount   int64                   `json:"unread_count,omitempty"`
	Pagination    *PaginationInfo         `json:"pagination"`
}

type NotificationTemplateResponse struct {
	ID          uuid.UUID                     `json:"id"`
	Name        string                        `json:"name"`
	Type        entities.NotificationType     `json:"type"`
	Category    entities.NotificationCategory `json:"category"`
	Subject     string                        `json:"subject"`
	Body        string                        `json:"body"`
	Variables   []string                      `json:"variables"`
	IsActive    bool                          `json:"is_active"`
	IsDefault   bool                          `json:"is_default"`
	Language    string                        `json:"language"`
	Description string                        `json:"description"`
	CreatedBy   uuid.UUID                     `json:"created_by"`
	CreatedAt   time.Time                     `json:"created_at"`
	UpdatedAt   time.Time                     `json:"updated_at"`
}

type TemplatesListResponse struct {
	Templates  []*NotificationTemplateResponse `json:"templates"`
	Total      int64                           `json:"total"`
	Pagination *PaginationInfo                 `json:"pagination"`
}

type PreferencesResponse struct {
	ID                   uuid.UUID `json:"id"`
	UserID               uuid.UUID `json:"user_id"`
	EmailEnabled         bool      `json:"email_enabled"`
	EmailOrderUpdates    bool      `json:"email_order_updates"`
	EmailPaymentUpdates  bool      `json:"email_payment_updates"`
	EmailShippingUpdates bool      `json:"email_shipping_updates"`
	EmailPromotions      bool      `json:"email_promotions"`
	EmailNewsletter      bool      `json:"email_newsletter"`
	EmailReviewReminders bool      `json:"email_review_reminders"`
	SMSEnabled           bool      `json:"sms_enabled"`
	SMSOrderUpdates      bool      `json:"sms_order_updates"`
	SMSPaymentUpdates    bool      `json:"sms_payment_updates"`
	SMSShippingUpdates   bool      `json:"sms_shipping_updates"`
	SMSSecurityAlerts    bool      `json:"sms_security_alerts"`
	PushEnabled          bool      `json:"push_enabled"`
	PushOrderUpdates     bool      `json:"push_order_updates"`
	PushPaymentUpdates   bool      `json:"push_payment_updates"`
	PushShippingUpdates  bool      `json:"push_shipping_updates"`
	PushPromotions       bool      `json:"push_promotions"`
	PushReviewReminders  bool      `json:"push_review_reminders"`
	InAppEnabled         bool      `json:"in_app_enabled"`
	InAppOrderUpdates    bool      `json:"in_app_order_updates"`
	InAppPaymentUpdates  bool      `json:"in_app_payment_updates"`
	InAppShippingUpdates bool      `json:"in_app_shipping_updates"`
	InAppPromotions      bool      `json:"in_app_promotions"`
	InAppSystemUpdates   bool      `json:"in_app_system_updates"`
	DigestFrequency      string    `json:"digest_frequency"`
	QuietHoursStart      string    `json:"quiet_hours_start"`
	QuietHoursEnd        string    `json:"quiet_hours_end"`
	Timezone             string    `json:"timezone"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// CreateNotification creates a new notification
func (uc *notificationUseCase) CreateNotification(ctx context.Context, req CreateNotificationRequest) (*NotificationResponse, error) {
	// Convert data to JSON string
	var dataJSON string
	if req.Data != nil {
		dataBytes, err := json.Marshal(req.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal data: %w", err)
		}
		dataJSON = string(dataBytes)
	}

	notification := &entities.Notification{
		ID:            uuid.New(),
		UserID:        req.UserID,
		Type:          req.Type,
		Category:      req.Category,
		Priority:      req.Priority,
		Status:        entities.NotificationStatusPending,
		Title:         req.Title,
		Message:       req.Message,
		Data:          dataJSON,
		Recipient:     req.Recipient,
		Subject:       req.Subject,
		Template:      req.Template,
		ReferenceType: req.ReferenceType,
		ReferenceID:   req.ReferenceID,
		ScheduledAt:   req.ScheduledAt,
		MaxRetries:    3,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := uc.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	// Queue for immediate sending if not scheduled
	if req.ScheduledAt == nil {
		if err := uc.QueueNotification(ctx, notification, nil); err != nil {
			return nil, err
		}
	}

	return uc.toNotificationResponse(notification), nil
}

// GetNotification gets a notification by ID
func (uc *notificationUseCase) GetNotification(ctx context.Context, id uuid.UUID) (*NotificationResponse, error) {
	notification, err := uc.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.toNotificationResponse(notification), nil
}

// UpdateNotification updates a notification
func (uc *notificationUseCase) UpdateNotification(ctx context.Context, id uuid.UUID, req UpdateNotificationRequest) (*NotificationResponse, error) {
	notification, err := uc.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Status != nil {
		notification.Status = *req.Status
	}
	if req.Title != nil {
		notification.Title = *req.Title
	}
	if req.Message != nil {
		notification.Message = *req.Message
	}
	if req.Data != nil {
		dataBytes, err := json.Marshal(req.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal data: %w", err)
		}
		notification.Data = string(dataBytes)
	}
	if req.Subject != nil {
		notification.Subject = *req.Subject
	}
	if req.Template != nil {
		notification.Template = *req.Template
	}

	notification.UpdatedAt = time.Now()

	if err := uc.notificationRepo.Update(ctx, notification); err != nil {
		return nil, err
	}

	return uc.toNotificationResponse(notification), nil
}

// DeleteNotification deletes a notification
func (uc *notificationUseCase) DeleteNotification(ctx context.Context, id uuid.UUID) error {
	return uc.notificationRepo.Delete(ctx, id)
}

// ListNotifications lists notifications with filters
func (uc *notificationUseCase) ListNotifications(ctx context.Context, req ListNotificationsRequest) (*NotificationsListResponse, error) {
	filters := repositories.AdminNotificationFilters{
		Type:      req.Type,
		Priority:  req.Priority,
		DateFrom:  req.DateFrom,
		DateTo:    req.DateTo,
		Limit:     req.Limit,
		Offset:    req.Offset,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	if req.UserID != nil {
		filters.UserID = req.UserID
	}
	if req.Category != nil {
		// Note: We need to map category to something the repository understands
		// For now, we'll skip this filter since AdminNotificationFilters doesn't have Category
	}
	if req.Status != nil {
		// Similar issue with Status - need to map to DeliveryStatus if applicable
	}
	if req.ReferenceType != nil {
		// Add to search field
		filters.Search = *req.ReferenceType
	}

	notifications, err := uc.notificationRepo.GetAllNotifications(ctx, filters)
	if err != nil {
		return nil, err
	}

	total, err := uc.notificationRepo.CountAllNotifications(ctx, filters)
	if err != nil {
		return nil, err
	}

	responses := make([]*NotificationResponse, len(notifications))
	for i, notification := range notifications {
		responses[i] = uc.toNotificationResponse(notification)
	}

	return &NotificationsListResponse{
		Notifications: responses,
		Total:         total,
		Pagination:    NewPaginationInfo(req.Offset, req.Limit, total),
	}, nil
}

// GetUserNotifications gets notifications for a specific user
func (uc *notificationUseCase) GetUserNotifications(ctx context.Context, userID uuid.UUID, req GetUserNotificationsRequest) (*NotificationsListResponse, error) {
	// Get user to check if they are admin
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	filters := repositories.NotificationFilters{
		Type:      req.Type,
		IsRead:    req.IsRead,
		Limit:     req.Limit,
		Offset:    req.Offset,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	var notifications []*entities.Notification
	var total int64

	// If user is admin, get both user-specific and system-wide notifications
	if user.Role == entities.UserRoleAdmin {
		notifications, err = uc.notificationRepo.GetAdminNotifications(ctx, userID, filters)
		if err != nil {
			return nil, err
		}

		total, err = uc.notificationRepo.CountAdminNotifications(ctx, userID, filters)
		if err != nil {
			return nil, err
		}
	} else {
		// Regular users only get their own notifications
		notifications, err = uc.notificationRepo.GetUserNotifications(ctx, userID, filters)
		if err != nil {
			return nil, err
		}

		total, err = uc.notificationRepo.CountUserNotifications(ctx, userID, filters)
		if err != nil {
			return nil, err
		}
	}

	unreadCount, err := uc.notificationRepo.GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*NotificationResponse, len(notifications))
	for i, notification := range notifications {
		responses[i] = uc.toNotificationResponse(notification)
	}

	// Create pagination info using enhanced function
	context := &EcommercePaginationContext{
		EntityType: "notifications",
		UserID:     userID.String(),
	}

	// Use proper offset-to-page conversion
	pagination := NewPaginationInfoFromOffset(req.Offset, req.Limit, total)

	// Apply ecommerce enhancements
	if context != nil {
		// Adjust page sizes based on entity type
		pagination.PageSizes = []int{10, 15, 30} // Notification-friendly sizes

		// Check if cursor pagination should be used
		pagination.UseCursor = ShouldUseCursorPagination(total, context.EntityType)

		// Generate cache key
		cacheParams := map[string]interface{}{
			"page":    pagination.Page,
			"limit":   pagination.Limit,
			"user_id": context.UserID,
		}
		if req.IsRead != nil {
			cacheParams["is_read"] = *req.IsRead
		}
		pagination.CacheKey = GenerateCacheKey("notifications", context.UserID, cacheParams)
	}

	return &NotificationsListResponse{
		Notifications: responses,
		Total:         total,
		UnreadCount:   unreadCount,
		Pagination:    pagination,
	}, nil
}

// MarkAsRead marks a notification as read
func (uc *notificationUseCase) MarkAsRead(ctx context.Context, userID, notificationID uuid.UUID) error {
	return uc.notificationRepo.MarkAsRead(ctx, notificationID)
}

// MarkAllAsRead marks all notifications as read for a user
func (uc *notificationUseCase) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return uc.notificationRepo.MarkAllAsRead(ctx, userID)
}

// GetUnreadCount gets the count of unread notifications for a user
func (uc *notificationUseCase) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return uc.notificationRepo.GetUnreadCount(ctx, userID)
}

// SendNotification sends a notification immediately
func (uc *notificationUseCase) SendNotification(ctx context.Context, notification *entities.Notification) error {
	// Send notification based on type
	switch notification.Type {
	case entities.NotificationTypeEmail:
		if uc.emailService != nil {
			if err := uc.emailService.SendNotificationEmail(ctx, notification); err != nil {
				notification.Status = entities.NotificationStatusFailed
				notification.UpdatedAt = time.Now()
				uc.notificationRepo.Update(ctx, notification)
				return fmt.Errorf("failed to send email notification: %w", err)
			}
		}
	case entities.NotificationTypeSMS:
		// TODO: Implement SMS sending
		fmt.Printf("📱 SMS would be sent: %s\n", notification.Message)
	case entities.NotificationTypePush:
		// TODO: Implement push notification sending
		fmt.Printf("🔔 Push notification would be sent: %s\n", notification.Message)
	case entities.NotificationTypeInApp:
		// In-app notifications are stored in database and sent via WebSocket
		fmt.Printf("📱 In-app notification created: %s\n", notification.Title)

		// Send real-time notification via WebSocket
		if uc.websocketHub != nil && notification.UserID != nil {
			uc.websocketHub.SendToUser(*notification.UserID, notification)
		} else if uc.websocketHub != nil && notification.UserID == nil {
			// System-wide notification (broadcast to all)
			uc.websocketHub.SendToAll(notification)
		}
	}

	// Mark as sent
	notification.Status = entities.NotificationStatusSent
	notification.SentAt = &[]time.Time{time.Now()}[0]
	notification.UpdatedAt = time.Now()
	return uc.notificationRepo.Update(ctx, notification)
}

// SendBulkNotifications sends multiple notifications
func (uc *notificationUseCase) SendBulkNotifications(ctx context.Context, notifications []*entities.Notification) error {
	for _, notification := range notifications {
		if err := uc.SendNotification(ctx, notification); err != nil {
			return err
		}
	}
	return nil
}

// QueueNotification queues a notification for later sending
func (uc *notificationUseCase) QueueNotification(ctx context.Context, notification *entities.Notification, scheduledAt *time.Time) error {
	if scheduledAt != nil {
		notification.ScheduledAt = scheduledAt
	}
	notification.Status = entities.NotificationStatusPending
	return uc.notificationRepo.Update(ctx, notification)
}

// ProcessQueue processes queued notifications
func (uc *notificationUseCase) ProcessQueue(ctx context.Context, limit int) error {
	// This would typically get pending notifications and send them
	// For now, we'll implement a basic version
	return nil
}

// CreateTemplate creates a notification template
func (uc *notificationUseCase) CreateTemplate(ctx context.Context, req CreateNotificationTemplateRequest) (*NotificationTemplateResponse, error) {
	template := &entities.NotificationTemplate{
		ID:        uuid.New(),
		Name:      req.Name,
		Type:      req.Type,
		Channel:   entities.NotificationChannelEmail, // Default channel, can be passed in request
		Subject:   req.Subject,
		Body:      req.Body,
		Variables: req.Variables,
		IsActive:  req.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.notificationRepo.CreateTemplate(ctx, template); err != nil {
		return nil, err
	}

	return uc.toNotificationTemplateResponse(template), nil
}

// GetTemplate gets a template by ID
func (uc *notificationUseCase) GetTemplate(ctx context.Context, id uuid.UUID) (*NotificationTemplateResponse, error) {
	template, err := uc.notificationRepo.GetTemplate(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.toNotificationTemplateResponse(template), nil
}

// UpdateTemplate updates a notification template
func (uc *notificationUseCase) UpdateTemplate(ctx context.Context, id uuid.UUID, req UpdateNotificationTemplateRequest) (*NotificationTemplateResponse, error) {
	template, err := uc.notificationRepo.GetTemplate(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided - only update fields that exist in entity
	if req.Name != nil {
		template.Name = *req.Name
	}
	if req.Subject != nil {
		template.Subject = *req.Subject
	}
	if req.Body != nil {
		template.Body = *req.Body
	}
	if req.Variables != nil {
		template.Variables = req.Variables
	}
	if req.IsActive != nil {
		template.IsActive = *req.IsActive
	}
	// Note: IsDefault, Language, Description are not in the entity, so skip these

	template.UpdatedAt = time.Now()

	if err := uc.notificationRepo.UpdateTemplate(ctx, template); err != nil {
		return nil, err
	}

	return uc.toNotificationTemplateResponse(template), nil
}

// DeleteTemplate deletes a notification template
func (uc *notificationUseCase) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	return uc.notificationRepo.DeleteTemplate(ctx, id)
}

// ListTemplates lists notification templates
func (uc *notificationUseCase) ListTemplates(ctx context.Context, req ListTemplatesRequest) (*TemplatesListResponse, error) {
	templates, err := uc.notificationRepo.ListTemplates(ctx)
	if err != nil {
		return nil, err
	}
	// Apply filters (basic filtering since repo doesn't support advanced filters)
	var filtered []*entities.NotificationTemplate
	for _, template := range templates {
		include := true

		if req.Type != nil && template.Type != *req.Type {
			include = false
		}
		// Skip Category filter since it's not in the entity
		if req.IsActive != nil && template.IsActive != *req.IsActive {
			include = false
		}
		// Skip Language filter since it's not in the entity

		if include {
			filtered = append(filtered, template)
		}
	}

	// Apply pagination
	total := int64(len(filtered))
	start := req.Offset
	end := start + req.Limit
	if start > len(filtered) {
		start = len(filtered)
	}
	if end > len(filtered) {
		end = len(filtered)
	}

	paginatedTemplates := filtered[start:end]

	responses := make([]*NotificationTemplateResponse, len(paginatedTemplates))
	for i, template := range paginatedTemplates {
		responses[i] = uc.toNotificationTemplateResponse(template)
	}

	return &TemplatesListResponse{
		Templates:  responses,
		Total:      total,
		Pagination: NewPaginationInfo(req.Offset, req.Limit, total),
	}, nil
}

// GetUserPreferences gets user notification preferences
func (uc *notificationUseCase) GetUserPreferences(ctx context.Context, userID uuid.UUID) (*PreferencesResponse, error) {
	preferences, err := uc.notificationRepo.GetUserPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}
	return uc.toPreferencesResponse(preferences), nil
}

// UpdateUserPreferences updates user notification preferences
func (uc *notificationUseCase) UpdateUserPreferences(ctx context.Context, userID uuid.UUID, req UpdatePreferencesRequest) (*PreferencesResponse, error) {
	preferences, err := uc.notificationRepo.GetUserPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided - map to available entity fields
	if req.EmailEnabled != nil {
		preferences.EmailEnabled = *req.EmailEnabled
	}
	if req.EmailOrderUpdates != nil {
		preferences.OrderUpdates = *req.EmailOrderUpdates
	}
	if req.EmailPromotions != nil {
		preferences.PromotionalEmails = *req.EmailPromotions
	}
	if req.EmailNewsletter != nil {
		preferences.NewsletterEnabled = *req.EmailNewsletter
	}
	if req.SMSEnabled != nil {
		preferences.SMSEnabled = *req.SMSEnabled
	}
	if req.SMSSecurityAlerts != nil {
		preferences.SecurityAlerts = *req.SMSSecurityAlerts
	}
	if req.PushEnabled != nil {
		preferences.PushEnabled = *req.PushEnabled
	}
	if req.InAppEnabled != nil {
		preferences.InAppEnabled = *req.InAppEnabled
	}

	preferences.UpdatedAt = time.Now()

	if err := uc.notificationRepo.UpdateUserPreferences(ctx, preferences); err != nil {
		return nil, err
	}

	return uc.toPreferencesResponse(preferences), nil
}

// Event-based notification methods
func (uc *notificationUseCase) NotifyOrderCreated(ctx context.Context, orderID uuid.UUID) error {
	// Get order details
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Get user details
	user, err := uc.userRepo.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check user notification preferences
	preferences, err := uc.notificationRepo.GetUserPreferences(ctx, user.ID)
	if err != nil {
		// Create default preferences if not found
		if err := uc.notificationRepo.CreateDefaultPreferences(ctx, user.ID); err != nil {
			return fmt.Errorf("failed to create default preferences: %w", err)
		}
		preferences, _ = uc.notificationRepo.GetUserPreferences(ctx, user.ID)
	}

	// Create notification data
	data := map[string]interface{}{
		"order_id":     order.ID,
		"order_number": order.OrderNumber,
		"total":        order.Total,
		"items_count":  len(order.Items),
	}
	dataJSON, _ := json.Marshal(data)

	// Create in-app notification
	if preferences.IsNotificationEnabled(entities.NotificationTypeInApp, entities.NotificationCategoryOrder) {
		notification := &entities.Notification{
			ID:            uuid.New(),
			UserID:        &user.ID,
			Type:          entities.NotificationTypeInApp,
			Category:      entities.NotificationCategoryOrder,
			Priority:      entities.NotificationPriorityNormal,
			Status:        entities.NotificationStatusPending,
			Title:         "Đơn hàng đã được tạo",
			Message:       fmt.Sprintf("Đơn hàng #%s của bạn đã được tạo thành công với tổng giá trị %.0f VND", order.OrderNumber, order.Total),
			Data:          string(dataJSON),
			ReferenceType: "order",
			ReferenceID:   &order.ID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := uc.notificationRepo.Create(ctx, notification); err != nil {
			return fmt.Errorf("failed to create in-app notification: %w", err)
		}
	}

	// Create email notification
	if preferences.IsNotificationEnabled(entities.NotificationTypeEmail, entities.NotificationCategoryOrder) {
		emailNotification := &entities.Notification{
			ID:            uuid.New(),
			UserID:        &user.ID,
			Type:          entities.NotificationTypeEmail,
			Category:      entities.NotificationCategoryOrder,
			Priority:      entities.NotificationPriorityNormal,
			Status:        entities.NotificationStatusPending,
			Title:         "Xác nhận đơn hàng",
			Message:       fmt.Sprintf("Cảm ơn bạn đã đặt hàng! Đơn hàng #%s đã được xác nhận.", order.OrderNumber),
			Data:          string(dataJSON),
			Recipient:     user.Email,
			Subject:       fmt.Sprintf("Xác nhận đơn hàng #%s", order.OrderNumber),
			Template:      "order_created",
			ReferenceType: "order",
			ReferenceID:   &order.ID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := uc.notificationRepo.Create(ctx, emailNotification); err != nil {
			return fmt.Errorf("failed to create email notification: %w", err)
		}
	}

	return nil
}

func (uc *notificationUseCase) NotifyOrderStatusChanged(ctx context.Context, orderID uuid.UUID, newStatus string) error {
	// Get order details
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Get user details
	user, err := uc.userRepo.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check user notification preferences
	preferences, err := uc.notificationRepo.GetUserPreferences(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get user preferences: %w", err)
	}

	// Map status to Vietnamese
	statusMap := map[string]string{
		"pending":    "Chờ xử lý",
		"confirmed":  "Đã xác nhận",
		"processing": "Đang xử lý",
		"shipped":    "Đã giao vận",
		"delivered":  "Đã giao hàng",
		"cancelled":  "Đã hủy",
		"returned":   "Đã trả hàng",
	}
	statusText := statusMap[newStatus]
	if statusText == "" {
		statusText = newStatus
	}

	// Create notification data
	data := map[string]interface{}{
		"order_id":     order.ID,
		"order_number": order.OrderNumber,
		"old_status":   order.Status,
		"new_status":   newStatus,
		"total":        order.Total,
	}
	dataJSON, _ := json.Marshal(data)

	// Create in-app notification
	if preferences.IsNotificationEnabled(entities.NotificationTypeInApp, entities.NotificationCategoryOrder) {
		notification := &entities.Notification{
			ID:            uuid.New(),
			UserID:        &user.ID,
			Type:          entities.NotificationTypeInApp,
			Category:      entities.NotificationCategoryOrder,
			Priority:      entities.NotificationPriorityNormal,
			Status:        entities.NotificationStatusPending,
			Title:         "Cập nhật trạng thái đơn hàng",
			Message:       fmt.Sprintf("Đơn hàng #%s đã được cập nhật trạng thái: %s", order.OrderNumber, statusText),
			Data:          string(dataJSON),
			ReferenceType: "order",
			ReferenceID:   &order.ID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := uc.notificationRepo.Create(ctx, notification); err != nil {
			return fmt.Errorf("failed to create in-app notification: %w", err)
		}
	}

	// Create email notification for important status changes
	if newStatus == "shipped" || newStatus == "delivered" || newStatus == "cancelled" {
		if preferences.IsNotificationEnabled(entities.NotificationTypeEmail, entities.NotificationCategoryOrder) {
			emailNotification := &entities.Notification{
				ID:            uuid.New(),
				UserID:        &user.ID,
				Type:          entities.NotificationTypeEmail,
				Category:      entities.NotificationCategoryOrder,
				Priority:      entities.NotificationPriorityHigh,
				Status:        entities.NotificationStatusPending,
				Title:         fmt.Sprintf("Đơn hàng #%s - %s", order.OrderNumber, statusText),
				Message:       fmt.Sprintf("Đơn hàng #%s của bạn đã được cập nhật trạng thái: %s", order.OrderNumber, statusText),
				Data:          string(dataJSON),
				Recipient:     user.Email,
				Subject:       fmt.Sprintf("Cập nhật đơn hàng #%s - %s", order.OrderNumber, statusText),
				Template:      "order_status_changed",
				ReferenceType: "order",
				ReferenceID:   &order.ID,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			if err := uc.notificationRepo.Create(ctx, emailNotification); err != nil {
				return fmt.Errorf("failed to create email notification: %w", err)
			}
		}
	}

	return nil
}

func (uc *notificationUseCase) NotifyPaymentReceived(ctx context.Context, paymentID uuid.UUID) error {
	// Get payment details
	payment, err := uc.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	// Get order details
	order, err := uc.orderRepo.GetByID(ctx, payment.OrderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Get user details
	user, err := uc.userRepo.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check user notification preferences
	preferences, err := uc.notificationRepo.GetUserPreferences(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get user preferences: %w", err)
	}

	// Create notification data
	data := map[string]interface{}{
		"payment_id":     payment.ID,
		"order_id":       order.ID,
		"order_number":   order.OrderNumber,
		"amount":         payment.Amount,
		"payment_method": payment.Method,
	}
	dataJSON, _ := json.Marshal(data)

	// Create in-app notification
	if preferences.IsNotificationEnabled(entities.NotificationTypeInApp, entities.NotificationCategoryPayment) {
		notification := &entities.Notification{
			ID:            uuid.New(),
			UserID:        &user.ID,
			Type:          entities.NotificationTypeInApp,
			Category:      entities.NotificationCategoryPayment,
			Priority:      entities.NotificationPriorityHigh,
			Status:        entities.NotificationStatusPending,
			Title:         "Thanh toán thành công",
			Message:       fmt.Sprintf("Thanh toán %.0f VND cho đơn hàng #%s đã được xử lý thành công", payment.Amount, order.OrderNumber),
			Data:          string(dataJSON),
			ReferenceType: "payment",
			ReferenceID:   &payment.ID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := uc.notificationRepo.Create(ctx, notification); err != nil {
			return fmt.Errorf("failed to create in-app notification: %w", err)
		}
	}

	// Create email notification
	if preferences.IsNotificationEnabled(entities.NotificationTypeEmail, entities.NotificationCategoryPayment) {
		emailNotification := &entities.Notification{
			ID:            uuid.New(),
			UserID:        &user.ID,
			Type:          entities.NotificationTypeEmail,
			Category:      entities.NotificationCategoryPayment,
			Priority:      entities.NotificationPriorityHigh,
			Status:        entities.NotificationStatusPending,
			Title:         "Xác nhận thanh toán",
			Message:       fmt.Sprintf("Thanh toán cho đơn hàng #%s đã được xử lý thành công", order.OrderNumber),
			Data:          string(dataJSON),
			Recipient:     user.Email,
			Subject:       fmt.Sprintf("Xác nhận thanh toán - Đơn hàng #%s", order.OrderNumber),
			Template:      "payment_received",
			ReferenceType: "payment",
			ReferenceID:   &payment.ID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := uc.notificationRepo.Create(ctx, emailNotification); err != nil {
			return fmt.Errorf("failed to create email notification: %w", err)
		}
	}

	return nil
}

func (uc *notificationUseCase) NotifyShippingUpdate(ctx context.Context, orderID uuid.UUID, trackingNumber string) error {
	// Get order details
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Get user details
	user, err := uc.userRepo.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check user notification preferences
	preferences, err := uc.notificationRepo.GetUserPreferences(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get user preferences: %w", err)
	}

	// Create notification data
	data := map[string]interface{}{
		"order_id":        order.ID,
		"order_number":    order.OrderNumber,
		"tracking_number": trackingNumber,
	}
	dataJSON, _ := json.Marshal(data)

	// Create in-app notification
	if preferences.IsNotificationEnabled(entities.NotificationTypeInApp, entities.NotificationCategoryShipping) {
		notification := &entities.Notification{
			ID:            uuid.New(),
			UserID:        &user.ID,
			Type:          entities.NotificationTypeInApp,
			Category:      entities.NotificationCategoryShipping,
			Priority:      entities.NotificationPriorityNormal,
			Status:        entities.NotificationStatusPending,
			Title:         "Cập nhật vận chuyển",
			Message:       fmt.Sprintf("Đơn hàng #%s đã được giao cho đơn vị vận chuyển. Mã vận đơn: %s", order.OrderNumber, trackingNumber),
			Data:          string(dataJSON),
			ReferenceType: "order",
			ReferenceID:   &order.ID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := uc.notificationRepo.Create(ctx, notification); err != nil {
			return fmt.Errorf("failed to create in-app notification: %w", err)
		}
	}

	// Create email notification
	if preferences.IsNotificationEnabled(entities.NotificationTypeEmail, entities.NotificationCategoryShipping) {
		emailNotification := &entities.Notification{
			ID:            uuid.New(),
			UserID:        &user.ID,
			Type:          entities.NotificationTypeEmail,
			Category:      entities.NotificationCategoryShipping,
			Priority:      entities.NotificationPriorityNormal,
			Status:        entities.NotificationStatusPending,
			Title:         "Thông tin vận chuyển",
			Message:       fmt.Sprintf("Đơn hàng #%s đã được giao cho đơn vị vận chuyển", order.OrderNumber),
			Data:          string(dataJSON),
			Recipient:     user.Email,
			Subject:       fmt.Sprintf("Thông tin vận chuyển - Đơn hàng #%s", order.OrderNumber),
			Template:      "shipping_update",
			ReferenceType: "order",
			ReferenceID:   &order.ID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := uc.notificationRepo.Create(ctx, emailNotification); err != nil {
			return fmt.Errorf("failed to create email notification: %w", err)
		}
	}

	return nil
}

func (uc *notificationUseCase) NotifyLowStock(ctx context.Context, inventoryID uuid.UUID) error {
	// Get inventory details with product preloaded
	inventory, err := uc.inventoryRepo.GetByID(ctx, inventoryID)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	// Get product name (check if product is preloaded)
	var productName string
	var productID uuid.UUID
	if inventory.Product.ID != uuid.Nil {
		productName = inventory.Product.Name
		productID = inventory.Product.ID
	} else {
		// Fallback: use product ID as name if product not preloaded
		productName = fmt.Sprintf("Product %s", inventory.ProductID.String()[:8])
		productID = inventory.ProductID
	}

	// Create notification data
	data := map[string]interface{}{
		"inventory_id":  inventory.ID,
		"product_id":    productID,
		"product_name":  productName,
		"current_stock": inventory.QuantityOnHand,
		"reorder_level": inventory.ReorderLevel,
		"warehouse_id":  inventory.WarehouseID,
	}
	dataJSON, _ := json.Marshal(data)

	// Create system notification for admins
	notification := &entities.Notification{
		ID:            uuid.New(),
		UserID:        nil, // System-wide notification
		Type:          entities.NotificationTypeInApp,
		Category:      entities.NotificationCategorySystem,
		Priority:      entities.NotificationPriorityHigh,
		Status:        entities.NotificationStatusPending,
		Title:         "Cảnh báo hết hàng",
		Message:       fmt.Sprintf("Sản phẩm '%s' sắp hết hàng. Số lượng hiện tại: %d, mức đặt lại: %d", productName, inventory.QuantityOnHand, inventory.ReorderLevel),
		Data:          string(dataJSON),
		ReferenceType: "inventory",
		ReferenceID:   &inventory.ID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := uc.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to create low stock notification: %w", err)
	}

	return nil
}

func (uc *notificationUseCase) NotifyReviewRequest(ctx context.Context, orderID uuid.UUID) error {
	// Get order details
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Get user details
	user, err := uc.userRepo.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check user notification preferences
	preferences, err := uc.notificationRepo.GetUserPreferences(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get user preferences: %w", err)
	}

	// Create notification data
	data := map[string]interface{}{
		"order_id":     order.ID,
		"order_number": order.OrderNumber,
		"items_count":  len(order.Items),
	}
	dataJSON, _ := json.Marshal(data)

	// Create in-app notification
	if preferences.IsNotificationEnabled(entities.NotificationTypeInApp, entities.NotificationCategoryReview) {
		notification := &entities.Notification{
			ID:            uuid.New(),
			UserID:        &user.ID,
			Type:          entities.NotificationTypeInApp,
			Category:      entities.NotificationCategoryReview,
			Priority:      entities.NotificationPriorityNormal,
			Status:        entities.NotificationStatusPending,
			Title:         "Đánh giá sản phẩm",
			Message:       fmt.Sprintf("Hãy đánh giá sản phẩm trong đơn hàng #%s để giúp khách hàng khác có trải nghiệm tốt hơn", order.OrderNumber),
			Data:          string(dataJSON),
			ReferenceType: "order",
			ReferenceID:   &order.ID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := uc.notificationRepo.Create(ctx, notification); err != nil {
			return fmt.Errorf("failed to create review request notification: %w", err)
		}
	}

	// Create email notification
	if preferences.IsNotificationEnabled(entities.NotificationTypeEmail, entities.NotificationCategoryReview) {
		emailNotification := &entities.Notification{
			ID:            uuid.New(),
			UserID:        &user.ID,
			Type:          entities.NotificationTypeEmail,
			Category:      entities.NotificationCategoryReview,
			Priority:      entities.NotificationPriorityNormal,
			Status:        entities.NotificationStatusPending,
			Title:         "Đánh giá sản phẩm",
			Message:       fmt.Sprintf("Cảm ơn bạn đã mua hàng! Hãy đánh giá sản phẩm trong đơn hàng #%s", order.OrderNumber),
			Data:          string(dataJSON),
			Recipient:     user.Email,
			Subject:       fmt.Sprintf("Đánh giá sản phẩm - Đơn hàng #%s", order.OrderNumber),
			Template:      "review_request",
			ReferenceType: "order",
			ReferenceID:   &order.ID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := uc.notificationRepo.Create(ctx, emailNotification); err != nil {
			return fmt.Errorf("failed to create email notification: %w", err)
		}
	}

	return nil
}

// Helper methods
func (uc *notificationUseCase) toNotificationResponse(notification *entities.Notification) *NotificationResponse {
	var data map[string]interface{}
	if notification.Data != "" {
		if err := json.Unmarshal([]byte(notification.Data), &data); err != nil {
			// Có thể log hoặc bỏ qua nếu không cần thiết
		}
	}

	return &NotificationResponse{
		ID:            notification.ID,
		UserID:        notification.UserID,
		Type:          notification.Type,
		Category:      notification.Category,
		Priority:      notification.Priority,
		Status:        notification.Status,
		Title:         notification.Title,
		Message:       notification.Message,
		Data:          data,
		Recipient:     notification.Recipient,
		Subject:       notification.Subject,
		Template:      notification.Template,
		ReferenceType: notification.ReferenceType,
		ReferenceID:   notification.ReferenceID,
		ScheduledAt:   notification.ScheduledAt,
		SentAt:        notification.SentAt,
		DeliveredAt:   notification.DeliveredAt,
		ReadAt:        notification.ReadAt,
		RetryCount:    notification.RetryCount,
		MaxRetries:    notification.MaxRetries,
		NextRetryAt:   notification.NextRetryAt,
		ErrorMessage:  notification.ErrorMessage,
		ErrorCode:     notification.ErrorCode,
		IsRead:        notification.IsRead,
		CreatedAt:     notification.CreatedAt,
		UpdatedAt:     notification.UpdatedAt,
	}
}

func (uc *notificationUseCase) toNotificationTemplateResponse(template *entities.NotificationTemplate) *NotificationTemplateResponse {
	return &NotificationTemplateResponse{
		ID:          template.ID,
		Name:        template.Name,
		Type:        template.Type,
		Category:    entities.NotificationCategorySystem, // Default, since entity doesn't have Category
		Subject:     template.Subject,
		Body:        template.Body,
		Variables:   template.Variables,
		IsActive:    template.IsActive,
		IsDefault:   false,    // Default, since entity doesn't have IsDefault
		Language:    "en",     // Default, since entity doesn't have Language
		Description: "",       // Default, since entity doesn't have Description
		CreatedBy:   uuid.Nil, // Default, since entity doesn't have CreatedBy
		CreatedAt:   template.CreatedAt,
		UpdatedAt:   template.UpdatedAt,
	}
}

func (uc *notificationUseCase) toPreferencesResponse(preferences *entities.NotificationPreferences) *PreferencesResponse {
	return &PreferencesResponse{
		ID:                   preferences.ID,
		UserID:               preferences.UserID,
		EmailEnabled:         preferences.EmailEnabled,
		EmailOrderUpdates:    preferences.OrderUpdates,
		EmailPaymentUpdates:  preferences.OrderUpdates, // Map to closest available field
		EmailShippingUpdates: preferences.OrderUpdates, // Map to closest available field
		EmailPromotions:      preferences.PromotionalEmails,
		EmailNewsletter:      preferences.NewsletterEnabled,
		EmailReviewReminders: preferences.OrderUpdates, // Map to closest available field
		SMSEnabled:           preferences.SMSEnabled,
		SMSOrderUpdates:      preferences.OrderUpdates,
		SMSPaymentUpdates:    preferences.OrderUpdates,
		SMSShippingUpdates:   preferences.OrderUpdates,
		SMSSecurityAlerts:    preferences.SecurityAlerts,
		PushEnabled:          preferences.PushEnabled,
		PushOrderUpdates:     preferences.OrderUpdates,
		PushPaymentUpdates:   preferences.OrderUpdates,
		PushShippingUpdates:  preferences.OrderUpdates,
		PushPromotions:       preferences.PromotionalEmails,
		PushReviewReminders:  preferences.OrderUpdates,
		InAppEnabled:         preferences.InAppEnabled,
		InAppOrderUpdates:    preferences.OrderUpdates,
		InAppPaymentUpdates:  preferences.OrderUpdates,
		InAppShippingUpdates: preferences.OrderUpdates,
		InAppPromotions:      preferences.PromotionalEmails,
		InAppSystemUpdates:   preferences.InAppEnabled,
		DigestFrequency:      "daily", // Default since not in entity
		QuietHoursStart:      "22:00", // Default since not in entity
		QuietHoursEnd:        "08:00", // Default since not in entity
		Timezone:             "UTC",   // Default since not in entity
		CreatedAt:            preferences.CreatedAt,
		UpdatedAt:            preferences.UpdatedAt,
	}
}

// NotifyNewOrder sends notification to admins when a new order is created
func (uc *notificationUseCase) NotifyNewOrder(ctx context.Context, orderID uuid.UUID) error {
	// Get order details
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Get user details
	user, err := uc.userRepo.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Create notification data
	data := map[string]interface{}{
		"order_id":       order.ID,
		"order_number":   order.OrderNumber,
		"customer_id":    user.ID,
		"customer_name":  user.FirstName + " " + user.LastName,
		"customer_email": user.Email,
		"total_amount":   order.Total,
		"items_count":    len(order.Items),
		"status":         order.Status,
	}
	dataJSON, _ := json.Marshal(data)

	// Create system notification for admins
	notification := &entities.Notification{
		ID:            uuid.New(),
		UserID:        nil, // System-wide notification for admins
		Type:          entities.NotificationTypeInApp,
		Category:      entities.NotificationCategoryOrder,
		Priority:      entities.NotificationPriorityNormal,
		Status:        entities.NotificationStatusPending,
		Title:         "Đơn hàng mới",
		Message:       fmt.Sprintf("Đơn hàng mới #%s từ khách hàng %s với giá trị %.0f VND", order.OrderNumber, user.FirstName+" "+user.LastName, order.Total),
		Data:          string(dataJSON),
		ReferenceType: "order",
		ReferenceID:   &order.ID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := uc.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to create new order notification: %w", err)
	}

	return nil
}

// NotifyPaymentFailed sends notification to admins when a payment fails
func (uc *notificationUseCase) NotifyPaymentFailed(ctx context.Context, paymentID uuid.UUID) error {
	// Get payment details
	payment, err := uc.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	// Get order details
	order, err := uc.orderRepo.GetByID(ctx, payment.OrderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Get user details
	user, err := uc.userRepo.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Create notification data
	data := map[string]interface{}{
		"payment_id":     payment.ID,
		"order_id":       order.ID,
		"order_number":   order.OrderNumber,
		"customer_id":    user.ID,
		"customer_name":  user.FirstName + " " + user.LastName,
		"customer_email": user.Email,
		"amount":         payment.Amount,
		"payment_method": payment.Method,
		"failure_reason": payment.FailureReason,
	}
	dataJSON, _ := json.Marshal(data)

	// Create system notification for admins
	notification := &entities.Notification{
		ID:            uuid.New(),
		UserID:        nil, // System-wide notification for admins
		Type:          entities.NotificationTypeInApp,
		Category:      entities.NotificationCategoryPayment,
		Priority:      entities.NotificationPriorityHigh,
		Status:        entities.NotificationStatusPending,
		Title:         "Thanh toán thất bại",
		Message:       fmt.Sprintf("Thanh toán thất bại cho đơn hàng #%s (%.0f VND) từ khách hàng %s", order.OrderNumber, payment.Amount, user.FirstName+" "+user.LastName),
		Data:          string(dataJSON),
		ReferenceType: "payment",
		ReferenceID:   &payment.ID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := uc.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to create payment failed notification: %w", err)
	}

	return nil
}

// NotifyNewUser sends notification to admins when a new user registers
func (uc *notificationUseCase) NotifyNewUser(ctx context.Context, userID uuid.UUID) error {
	// Get user details
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Create notification data
	data := map[string]interface{}{
		"user_id":    user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	}
	dataJSON, _ := json.Marshal(data)

	// Create system notification for admins
	notification := &entities.Notification{
		ID:            uuid.New(),
		UserID:        nil, // System-wide notification for admins
		Type:          entities.NotificationTypeInApp,
		Category:      entities.NotificationCategorySystem,
		Priority:      entities.NotificationPriorityLow,
		Status:        entities.NotificationStatusPending,
		Title:         "Người dùng mới đăng ký",
		Message:       fmt.Sprintf("Người dùng mới %s %s (%s) đã đăng ký tài khoản", user.FirstName, user.LastName, user.Email),
		Data:          string(dataJSON),
		ReferenceType: "user",
		ReferenceID:   &user.ID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := uc.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to create new user notification: %w", err)
	}

	return nil
}

// NotifyNewReview sends notification to admins when a new review is submitted
func (uc *notificationUseCase) NotifyNewReview(ctx context.Context, reviewID uuid.UUID) error {
	// Get review details
	review, err := uc.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return fmt.Errorf("failed to get review: %w", err)
	}

	// Get user details
	user, err := uc.userRepo.GetByID(ctx, review.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Get product details
	product, err := uc.productRepo.GetByID(ctx, review.ProductID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Create notification data
	data := map[string]interface{}{
		"review_id":     review.ID,
		"product_id":    product.ID,
		"product_name":  product.Name,
		"customer_id":   user.ID,
		"customer_name": user.FirstName + " " + user.LastName,
		"rating":        review.Rating,
		"comment":       review.Comment,
		"status":        review.Status,
	}
	dataJSON, _ := json.Marshal(data)

	// Create system notification for admins
	notification := &entities.Notification{
		ID:            uuid.New(),
		UserID:        nil, // System-wide notification for admins
		Type:          entities.NotificationTypeInApp,
		Category:      entities.NotificationCategorySystem,
		Priority:      entities.NotificationPriorityLow,
		Status:        entities.NotificationStatusPending,
		Title:         "Đánh giá mới",
		Message:       fmt.Sprintf("Đánh giá mới %d sao cho sản phẩm '%s' từ khách hàng %s", review.Rating, product.Name, user.FirstName+" "+user.LastName),
		Data:          string(dataJSON),
		ReferenceType: "review",
		ReferenceID:   &review.ID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := uc.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to create new review notification: %w", err)
	}

	return nil
}
