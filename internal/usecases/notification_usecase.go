package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

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
}

type notificationUseCase struct {
	notificationRepo repositories.NotificationRepository
	userRepo         repositories.UserRepository
	orderRepo        repositories.OrderRepository
	inventoryRepo    repositories.InventoryRepository
	emailService     EmailService
	smsService       SMSService
	pushService      PushService
}

// NewNotificationUseCase creates a new notification use case
func NewNotificationUseCase(
	notificationRepo repositories.NotificationRepository,
	userRepo repositories.UserRepository,
	orderRepo repositories.OrderRepository,
	inventoryRepo repositories.InventoryRepository,
	emailService EmailService,
	smsService SMSService,
	pushService PushService,
) NotificationUseCase {
	return &notificationUseCase{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		orderRepo:        orderRepo,
		inventoryRepo:    inventoryRepo,
		emailService:     emailService,
		smsService:       smsService,
		pushService:      pushService,
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
	filters := repositories.NotificationFilters{
		Type:      req.Type,
		IsRead:    req.IsRead,
		Limit:     req.Limit,
		Offset:    req.Offset,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	notifications, err := uc.notificationRepo.GetUserNotifications(ctx, userID, filters)
	if err != nil {
		return nil, err
	}

	total, err := uc.notificationRepo.CountUserNotifications(ctx, userID, filters)
	if err != nil {
		return nil, err
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
	// Here you would implement the actual sending logic based on notification type
	// For now, we'll just mark it as sent
	notification.Status = entities.NotificationStatusSent
	notification.SentAt = &[]time.Time{time.Now()}[0]
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
	// Implement order created notification logic
	return nil
}

func (uc *notificationUseCase) NotifyOrderStatusChanged(ctx context.Context, orderID uuid.UUID, newStatus string) error {
	// Implement order status changed notification logic
	return nil
}

func (uc *notificationUseCase) NotifyPaymentReceived(ctx context.Context, paymentID uuid.UUID) error {
	// Implement payment received notification logic
	return nil
}

func (uc *notificationUseCase) NotifyShippingUpdate(ctx context.Context, orderID uuid.UUID, trackingNumber string) error {
	// Implement shipping update notification logic
	return nil
}

func (uc *notificationUseCase) NotifyLowStock(ctx context.Context, inventoryID uuid.UUID) error {
	// Implement low stock notification logic
	return nil
}

func (uc *notificationUseCase) NotifyReviewRequest(ctx context.Context, orderID uuid.UUID) error {
	// Implement review request notification logic
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
		IsRead:        notification.IsRead(),
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
