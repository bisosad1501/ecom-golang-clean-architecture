package database

import (
	"context"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *gorm.DB) repositories.NotificationRepository {
	return &notificationRepository{db: db}
}

// Create creates a new notification
func (r *notificationRepository) Create(ctx context.Context, notification *entities.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

// GetByID gets a notification by ID
func (r *notificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Notification, error) {
	var notification entities.Notification
	err := r.db.WithContext(ctx).First(&notification, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// GetByUser gets notifications for a user
func (r *notificationRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Notification, error) {
	var notifications []*entities.Notification
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

// GetUnreadByUser gets unread notifications for a user
func (r *notificationRepository) GetUnreadByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Notification, error) {
	var notifications []*entities.Notification
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND read_at IS NULL", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

// CountUnreadByUser counts unread notifications for a user
func (r *notificationRepository) CountUnreadByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

// MarkAsRead marks a notification as read
func (r *notificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"read_at":    time.Now(),
			"updated_at": time.Now(),
		}).Error
}

// MarkAllAsReadByUser marks all notifications as read for a user
func (r *notificationRepository) MarkAllAsReadByUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Updates(map[string]interface{}{
			"read_at":    time.Now(),
			"updated_at": time.Now(),
		}).Error
}

// Update updates a notification
func (r *notificationRepository) Update(ctx context.Context, notification *entities.Notification) error {
	notification.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(notification).Error
}

// Delete deletes a notification
func (r *notificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Notification{}, "id = ?", id).Error
}

// DeleteByUser deletes all notifications for a user
func (r *notificationRepository) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Notification{}, "user_id = ?", userID).Error
}

// List lists notifications with filters
func (r *notificationRepository) List(ctx context.Context, filters repositories.NotificationFilters) ([]*entities.Notification, error) {
	var notifications []*entities.Notification
	query := r.db.WithContext(ctx)

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}

	if filters.IsRead != nil {
		if *filters.IsRead {
			query = query.Where("read_at IS NOT NULL")
		} else {
			query = query.Where("read_at IS NULL")
		}
	}

	if filters.CreatedAfter != nil {
		query = query.Where("created_at >= ?", *filters.CreatedAfter)
	}

	if filters.CreatedBefore != nil {
		query = query.Where("created_at <= ?", *filters.CreatedBefore)
	}

	// Apply sorting
	switch filters.SortBy {
	case "created_at":
		if filters.SortOrder == "desc" {
			query = query.Order("created_at DESC")
		} else {
			query = query.Order("created_at ASC")
		}
	case "type":
		if filters.SortOrder == "desc" {
			query = query.Order("type DESC")
		} else {
			query = query.Order("type ASC")
		}
	default:
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&notifications).Error
	return notifications, err
}

// Count counts notifications with filters
func (r *notificationRepository) Count(ctx context.Context, filters repositories.NotificationFilters) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Notification{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}

	if filters.IsRead != nil {
		if *filters.IsRead {
			query = query.Where("read_at IS NOT NULL")
		} else {
			query = query.Where("read_at IS NULL")
		}
	}

	if filters.CreatedAfter != nil {
		query = query.Where("created_at >= ?", *filters.CreatedAfter)
	}

	if filters.CreatedBefore != nil {
		query = query.Where("created_at <= ?", *filters.CreatedBefore)
	}

	err := query.Count(&count).Error
	return count, err
}

// GetByType gets notifications by type
func (r *notificationRepository) GetByType(ctx context.Context, notificationType entities.NotificationType, limit, offset int) ([]*entities.Notification, error) {
	var notifications []*entities.Notification
	err := r.db.WithContext(ctx).
		Where("type = ?", notificationType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

// CreateBulk creates multiple notifications
func (r *notificationRepository) CreateBulk(ctx context.Context, notifications []*entities.Notification) error {
	return r.db.WithContext(ctx).CreateInBatches(notifications, 100).Error
}

// DeleteOld deletes old notifications
func (r *notificationRepository) DeleteOld(ctx context.Context, olderThan time.Time) error {
	return r.db.WithContext(ctx).
		Delete(&entities.Notification{}, "created_at < ?", olderThan).Error
}

// GetUserPreferences gets user notification preferences
func (r *notificationRepository) GetUserPreferences(ctx context.Context, userID uuid.UUID) (*entities.NotificationPreferences, error) {
	var prefs entities.NotificationPreferences
	err := r.db.WithContext(ctx).First(&prefs, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &prefs, nil
}

// UpdateUserPreferences updates user notification preferences
func (r *notificationRepository) UpdateUserPreferences(ctx context.Context, prefs *entities.NotificationPreferences) error {
	prefs.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(prefs).Error
}

// CreateUserPreferences creates user notification preferences
func (r *notificationRepository) CreateUserPreferences(ctx context.Context, prefs *entities.NotificationPreferences) error {
	return r.db.WithContext(ctx).Create(prefs).Error
}

// CountAllNotifications counts all notifications with filters
func (r *notificationRepository) CountAllNotifications(ctx context.Context, filters repositories.AdminNotificationFilters) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Notification{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}

	if filters.Priority != nil {
		query = query.Where("priority = ?", *filters.Priority)
	}

	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	if filters.Search != "" {
		query = query.Where("title LIKE ? OR message LIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

// CountUserNotifications counts notifications for a specific user with filters
func (r *notificationRepository) CountUserNotifications(ctx context.Context, userID uuid.UUID, filters repositories.NotificationFilters) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Notification{}).Where("user_id = ?", userID)

	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}

	if filters.IsRead != nil {
		if *filters.IsRead {
			query = query.Where("read_at IS NOT NULL")
		} else {
			query = query.Where("read_at IS NULL")
		}
	}

	if filters.Priority != nil {
		query = query.Where("priority = ?", *filters.Priority)
	}

	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	err := query.Count(&count).Error
	return count, err
}

// CreateDefaultPreferences creates default notification preferences for a user
func (r *notificationRepository) CreateDefaultPreferences(ctx context.Context, userID uuid.UUID) error {
	prefs := &entities.NotificationPreferences{
		ID:                uuid.New(),
		UserID:            userID,
		EmailEnabled:      true,
		PushEnabled:       true,
		SMSEnabled:        false,
		InAppEnabled:      true,
		OrderUpdates:      true,
		PromotionalEmails: true,
		SecurityAlerts:    true,
		NewsletterEnabled: false,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	return r.CreateUserPreferences(ctx, prefs)
}

// CreateTemplate creates a notification template
func (r *notificationRepository) CreateTemplate(ctx context.Context, template *entities.NotificationTemplate) error {
	template.ID = uuid.New()
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(template).Error
}

// DeleteOldNotifications deletes old notifications
func (r *notificationRepository) DeleteOldNotifications(ctx context.Context, olderThan time.Time) error {
	return r.db.WithContext(ctx).
		Delete(&entities.Notification{}, "created_at < ?", olderThan).Error
}

// DeleteTemplate deletes a notification template
func (r *notificationRepository) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.NotificationTemplate{}, "id = ?", id).Error
}

// GetAllNotifications gets all notifications with filters (alias for CountAllNotifications logic)
func (r *notificationRepository) GetAllNotifications(ctx context.Context, filters repositories.AdminNotificationFilters) ([]*entities.Notification, error) {
	var notifications []*entities.Notification
	query := r.db.WithContext(ctx).Model(&entities.Notification{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}

	if filters.Priority != nil {
		query = query.Where("priority = ?", *filters.Priority)
	}

	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	if filters.Search != "" {
		query = query.Where("title LIKE ? OR message LIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	// Apply sorting
	sortBy := "created_at"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}
	sortOrder := "DESC"
	if filters.SortOrder != "" {
		sortOrder = filters.SortOrder
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&notifications).Error
	return notifications, err
}

// GetDeliveryStats gets notification delivery statistics
func (r *notificationRepository) GetDeliveryStats(ctx context.Context, from, to time.Time) (*repositories.DeliveryStats, error) {
	var stats repositories.DeliveryStats

	// Get total notifications sent
	err := r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("created_at BETWEEN ? AND ?", from, to).
		Count(&stats.TotalSent).Error
	if err != nil {
		return nil, err
	}

	// Get delivered notifications
	err = r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("created_at BETWEEN ? AND ? AND delivery_status = ?", from, to, "delivered").
		Count(&stats.Delivered).Error
	if err != nil {
		return nil, err
	}

	// Get failed notifications
	err = r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("created_at BETWEEN ? AND ? AND delivery_status = ?", from, to, "failed").
		Count(&stats.Failed).Error
	if err != nil {
		return nil, err
	}

	// Calculate delivery rate
	if stats.TotalSent > 0 {
		stats.DeliveryRate = float64(stats.Delivered) / float64(stats.TotalSent) * 100
	}

	return &stats, nil
}

// GetEngagementStats gets notification engagement statistics
func (r *notificationRepository) GetEngagementStats(ctx context.Context, from, to time.Time) (*repositories.EngagementStats, error) {
	var stats repositories.EngagementStats

	// Get total notifications sent
	err := r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("created_at BETWEEN ? AND ?", from, to).
		Count(&stats.TotalNotifications).Error
	if err != nil {
		return nil, err
	}

	// Get opened notifications
	err = r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("created_at BETWEEN ? AND ? AND read_at IS NOT NULL", from, to).
		Count(&stats.OpenedNotifications).Error
	if err != nil {
		return nil, err
	}

	// Calculate open rate
	if stats.TotalNotifications > 0 {
		stats.OpenRate = float64(stats.OpenedNotifications) / float64(stats.TotalNotifications) * 100
	}

	return &stats, nil
}

// GetFailedNotifications gets failed notifications for retry
func (r *notificationRepository) GetFailedNotifications(ctx context.Context, retryCount int, limit int) ([]*entities.Notification, error) {
	var notifications []*entities.Notification
	err := r.db.WithContext(ctx).
		Where("delivery_status = ? AND retry_count <= ?", "failed", retryCount).
		Order("created_at ASC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

// GetNotificationStats gets notification statistics
func (r *notificationRepository) GetNotificationStats(ctx context.Context, from, to time.Time) (*repositories.NotificationStats, error) {
	var stats repositories.NotificationStats

	// Get total notifications
	err := r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("created_at BETWEEN ? AND ?", from, to).
		Count(&stats.TotalSent).Error
	if err != nil {
		return nil, err
	}

	// Get read notifications
	err = r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("created_at BETWEEN ? AND ? AND read_at IS NOT NULL", from, to).
		Count(&stats.TotalRead).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetNotificationsByType gets notifications by type
func (r *notificationRepository) GetNotificationsByType(ctx context.Context, notificationType entities.NotificationType, limit, offset int) ([]*entities.Notification, error) {
	var notifications []*entities.Notification
	err := r.db.WithContext(ctx).
		Where("type = ?", notificationType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

// GetPendingNotifications gets pending notifications for delivery by channel
func (r *notificationRepository) GetPendingNotifications(ctx context.Context, channel entities.NotificationChannel, limit int) ([]*entities.Notification, error) {
	var notifications []*entities.Notification
	err := r.db.WithContext(ctx).
		Where("(delivery_status = ? OR delivery_status IS NULL) AND channel = ?", "pending", channel).
		Order("created_at ASC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

// GetTemplate gets a notification template by ID
func (r *notificationRepository) GetTemplate(ctx context.Context, id uuid.UUID) (*entities.NotificationTemplate, error) {
	var template entities.NotificationTemplate
	err := r.db.WithContext(ctx).First(&template, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetTemplateByType gets a notification template by type
func (r *notificationRepository) GetTemplateByType(ctx context.Context, templateType entities.NotificationType) (*entities.NotificationTemplate, error) {
	var template entities.NotificationTemplate
	err := r.db.WithContext(ctx).
		Where("type = ? AND is_active = ?", templateType, true).
		First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetUnreadCount gets unread notification count for a user
func (r *notificationRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

// GetUserNotifications gets notifications for a user with filters
func (r *notificationRepository) GetUserNotifications(ctx context.Context, userID uuid.UUID, filters repositories.NotificationFilters) ([]*entities.Notification, error) {
	var notifications []*entities.Notification
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if filters.IsRead != nil {
		if *filters.IsRead {
			query = query.Where("read_at IS NOT NULL")
		} else {
			query = query.Where("read_at IS NULL")
		}
	}

	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}

	if filters.Priority != nil {
		query = query.Where("priority = ?", *filters.Priority)
	}

	err := query.Order("created_at DESC").
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&notifications).Error

	// Set IsRead field based on ReadAt
	for _, n := range notifications {
		n.IsRead = n.ReadAt != nil && !n.ReadAt.IsZero()
	}

	return notifications, err
}

// ListTemplates lists all notification templates
func (r *notificationRepository) ListTemplates(ctx context.Context) ([]*entities.NotificationTemplate, error) {
	var templates []*entities.NotificationTemplate
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&templates).Error
	return templates, err
}

// MarkAllAsRead marks all notifications as read for a user
func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Update("read_at", time.Now()).Error
}

// MarkAsDelivered marks a notification as delivered
func (r *notificationRepository) MarkAsDelivered(ctx context.Context, notificationID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("id = ?", notificationID).
		Updates(map[string]interface{}{
			"delivery_status": "delivered",
			"delivered_at":    time.Now(),
		}).Error
}

// MarkMultipleAsRead marks multiple notifications as read
func (r *notificationRepository) MarkMultipleAsRead(ctx context.Context, notificationIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("id IN (?)", notificationIDs).
		Update("read_at", time.Now()).Error
}

// UpdateDeliveryStatus updates notification delivery status with reason
func (r *notificationRepository) UpdateDeliveryStatus(ctx context.Context, notificationID uuid.UUID, status entities.DeliveryStatus, reason string) error {
	updates := map[string]interface{}{
		"delivery_status": status,
	}
	if reason != "" {
		updates["failure_reason"] = reason
	}
	if status == entities.DeliveryStatusDelivered {
		updates["delivered_at"] = time.Now()
	}

	return r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("id = ?", notificationID).
		Updates(updates).Error
}

// UpdateTemplate updates a notification template
func (r *notificationRepository) UpdateTemplate(ctx context.Context, template *entities.NotificationTemplate) error {
	template.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(template).Error
}

// GetPendingNotificationsForQueue gets pending notifications for queue processing
func (r *notificationRepository) GetPendingNotificationsForQueue(ctx context.Context, limit int) ([]*entities.Notification, error) {
	var notifications []*entities.Notification
	err := r.db.WithContext(ctx).
		Where("status = ?", entities.NotificationStatusPending).
		Where("(next_retry_at IS NULL OR next_retry_at <= ?)", time.Now()).
		Order("created_at ASC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

// GetRetryableNotifications gets notifications that are ready for retry
func (r *notificationRepository) GetRetryableNotifications(ctx context.Context, limit int) ([]*entities.Notification, error) {
	var notifications []*entities.Notification
	err := r.db.WithContext(ctx).
		Where("status = ? AND next_retry_at IS NOT NULL AND next_retry_at <= ?",
			entities.NotificationStatusPending, time.Now()).
		Order("next_retry_at ASC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

// GetPendingCount gets count of pending notifications
func (r *notificationRepository) GetPendingCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Notification{}).
		Where("status = ?", entities.NotificationStatusPending).
		Count(&count).Error
	return count, err
}

// GetProcessingCount gets count of processing notifications
func (r *notificationRepository) GetProcessingCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Notification{}).
		Where("status = ?", entities.NotificationStatusProcessing).
		Count(&count).Error
	return count, err
}

// GetFailedCount gets count of failed notifications
func (r *notificationRepository) GetFailedCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Notification{}).
		Where("status = ?", entities.NotificationStatusFailed).
		Count(&count).Error
	return count, err
}

// GetAdminNotifications gets notifications for admin users (both user-specific and system-wide)
func (r *notificationRepository) GetAdminNotifications(ctx context.Context, userID uuid.UUID, filters repositories.NotificationFilters) ([]*entities.Notification, error) {
	var notifications []*entities.Notification

	query := r.db.WithContext(ctx).Model(&entities.Notification{})

	// Admin sees both their own notifications and system-wide notifications (user_id IS NULL)
	query = query.Where("user_id = ? OR user_id IS NULL", userID)

	// Apply filters
	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}
	if filters.Priority != nil {
		query = query.Where("priority = ?", *filters.Priority)
	}
	if filters.IsRead != nil {
		if *filters.IsRead {
			query = query.Where("read_at IS NOT NULL")
		} else {
			query = query.Where("read_at IS NULL")
		}
	}
	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	// Apply sorting
	sortBy := "created_at"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}
	sortOrder := "DESC"
	if filters.SortOrder != "" {
		sortOrder = filters.SortOrder
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&notifications).Error
	return notifications, err
}

// CountAdminNotifications counts notifications for admin users (both user-specific and system-wide)
func (r *notificationRepository) CountAdminNotifications(ctx context.Context, userID uuid.UUID, filters repositories.NotificationFilters) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&entities.Notification{})

	// Admin sees both their own notifications and system-wide notifications (user_id IS NULL)
	query = query.Where("user_id = ? OR user_id IS NULL", userID)

	// Apply filters
	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}
	if filters.Priority != nil {
		query = query.Where("priority = ?", *filters.Priority)
	}
	if filters.IsRead != nil {
		if *filters.IsRead {
			query = query.Where("read_at IS NOT NULL")
		} else {
			query = query.Where("read_at IS NULL")
		}
	}
	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	err := query.Count(&count).Error
	return count, err
}
