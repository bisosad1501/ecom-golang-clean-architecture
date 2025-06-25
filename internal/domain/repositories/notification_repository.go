package repositories

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// NotificationRepository defines notification repository interface
type NotificationRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, notification *entities.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Notification, error)
	Update(ctx context.Context, notification *entities.Notification) error
	Delete(ctx context.Context, id uuid.UUID) error

	// User notifications
	GetUserNotifications(ctx context.Context, userID uuid.UUID, filters NotificationFilters) ([]*entities.Notification, error)
	CountUserNotifications(ctx context.Context, userID uuid.UUID, filters NotificationFilters) (int64, error)
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error)
	MarkAsRead(ctx context.Context, notificationID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	MarkAsDelivered(ctx context.Context, notificationID uuid.UUID) error

	// Bulk operations
	CreateBulk(ctx context.Context, notifications []*entities.Notification) error
	MarkMultipleAsRead(ctx context.Context, notificationIDs []uuid.UUID) error
	DeleteOldNotifications(ctx context.Context, olderThan time.Time) error

	// Template operations
	CreateTemplate(ctx context.Context, template *entities.NotificationTemplate) error
	GetTemplate(ctx context.Context, templateID uuid.UUID) (*entities.NotificationTemplate, error)
	GetTemplateByType(ctx context.Context, notificationType entities.NotificationType) (*entities.NotificationTemplate, error)
	UpdateTemplate(ctx context.Context, template *entities.NotificationTemplate) error
	DeleteTemplate(ctx context.Context, templateID uuid.UUID) error
	ListTemplates(ctx context.Context) ([]*entities.NotificationTemplate, error)

	// Preference operations
	GetUserPreferences(ctx context.Context, userID uuid.UUID) (*entities.NotificationPreferences, error)
	UpdateUserPreferences(ctx context.Context, preferences *entities.NotificationPreferences) error
	CreateDefaultPreferences(ctx context.Context, userID uuid.UUID) error

	// Delivery tracking
	GetPendingNotifications(ctx context.Context, channel entities.NotificationChannel, limit int) ([]*entities.Notification, error)
	GetFailedNotifications(ctx context.Context, retryCount int, limit int) ([]*entities.Notification, error)
	UpdateDeliveryStatus(ctx context.Context, notificationID uuid.UUID, status entities.DeliveryStatus, error string) error

	// Analytics
	GetNotificationStats(ctx context.Context, dateFrom, dateTo time.Time) (*NotificationStats, error)
	GetDeliveryStats(ctx context.Context, dateFrom, dateTo time.Time) (*DeliveryStats, error)
	GetEngagementStats(ctx context.Context, dateFrom, dateTo time.Time) (*EngagementStats, error)

	// Admin operations
	GetAllNotifications(ctx context.Context, filters AdminNotificationFilters) ([]*entities.Notification, error)
	CountAllNotifications(ctx context.Context, filters AdminNotificationFilters) (int64, error)
	GetNotificationsByType(ctx context.Context, notificationType entities.NotificationType, limit, offset int) ([]*entities.Notification, error)
}

// AdminNotificationFilters represents filters for admin notification queries
type AdminNotificationFilters struct {
	UserID           *uuid.UUID
	Type             *entities.NotificationType
	Channel          *entities.NotificationChannel
	DeliveryStatus   *entities.DeliveryStatus
	Priority         *entities.NotificationPriority
	DateFrom         *time.Time
	DateTo           *time.Time
	Search           string
	Limit            int
	Offset           int
	SortBy           string
	SortOrder        string
}







// EngagementMetric represents engagement metrics for a specific category
type EngagementMetric struct {
	Sent       int64   `json:"sent"`
	Opened     int64   `json:"opened"`
	Clicked    int64   `json:"clicked"`
	OpenRate   float64 `json:"open_rate"`
	ClickRate  float64 `json:"click_rate"`
}

// DailyNotificationStats represents daily notification statistics
type DailyNotificationStats struct {
	Date      time.Time `json:"date"`
	Sent      int64     `json:"sent"`
	Delivered int64     `json:"delivered"`
	Failed    int64     `json:"failed"`
	Read      int64     `json:"read"`
}
