package entities

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeEmail NotificationType = "email"
	NotificationTypeSMS   NotificationType = "sms"
	NotificationTypePush  NotificationType = "push"
	NotificationTypeInApp NotificationType = "in_app"
)

// NotificationCategory represents the category of notification
type NotificationCategory string

const (
	NotificationCategoryOrder     NotificationCategory = "order"
	NotificationCategoryPayment   NotificationCategory = "payment"
	NotificationCategoryShipping  NotificationCategory = "shipping"
	NotificationCategoryPromotion NotificationCategory = "promotion"
	NotificationCategoryAccount   NotificationCategory = "account"
	NotificationCategorySystem    NotificationCategory = "system"
	NotificationCategoryMarketing NotificationCategory = "marketing"
	NotificationCategoryReview    NotificationCategory = "review"
	NotificationCategoryInventory NotificationCategory = "inventory"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "pending"
	NotificationStatusSent      NotificationStatus = "sent"
	NotificationStatusDelivered NotificationStatus = "delivered"
	NotificationStatusFailed    NotificationStatus = "failed"
	NotificationStatusRead      NotificationStatus = "read"
)

// NotificationPriority represents the priority of a notification
type NotificationPriority string

const (
	NotificationPriorityLow      NotificationPriority = "low"
	NotificationPriorityNormal   NotificationPriority = "normal"
	NotificationPriorityHigh     NotificationPriority = "high"
	NotificationPriorityCritical NotificationPriority = "critical"
)

// Notification represents a notification to be sent to users
type Notification struct {
	ID          uuid.UUID            `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      *uuid.UUID           `json:"user_id" gorm:"type:uuid;index"`        // null for system-wide notifications
	User        *User                `json:"user,omitempty" gorm:"foreignKey:UserID"`
	
	// Notification details
	Type        NotificationType     `json:"type" gorm:"not null"`
	Category    NotificationCategory `json:"category" gorm:"not null"`
	Priority    NotificationPriority `json:"priority" gorm:"default:'normal'"`
	Status      NotificationStatus   `json:"status" gorm:"default:'pending'"`
	
	// Content
	Title       string               `json:"title" gorm:"not null"`
	Message     string               `json:"message" gorm:"type:text;not null"`
	Data        string               `json:"data" gorm:"type:text"`               // JSON data for additional context
	
	// Delivery details
	Recipient   string               `json:"recipient"`                           // Email address, phone number, etc.
	Subject     string               `json:"subject"`                             // For email notifications
	Template    string               `json:"template"`                            // Template name
	
	// Reference information
	ReferenceType string             `json:"reference_type"`                      // order, payment, product, etc.
	ReferenceID   *uuid.UUID         `json:"reference_id" gorm:"type:uuid;index"`
	
	// Scheduling
	ScheduledAt   *time.Time         `json:"scheduled_at"`                        // When to send
	SentAt        *time.Time         `json:"sent_at"`
	DeliveredAt   *time.Time         `json:"delivered_at"`
	ReadAt        *time.Time         `json:"read_at"`
	
	// Retry information
	RetryCount    int                `json:"retry_count" gorm:"default:0"`
	MaxRetries    int                `json:"max_retries" gorm:"default:3"`
	NextRetryAt   *time.Time         `json:"next_retry_at"`
	
	// Error information
	ErrorMessage  string             `json:"error_message"`
	ErrorCode     string             `json:"error_code"`
	
	// Metadata
	CreatedAt     time.Time          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time          `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Notification entity
func (Notification) TableName() string {
	return "notifications"
}

// IsPending checks if notification is pending
func (n *Notification) IsPending() bool {
	return n.Status == NotificationStatusPending
}

// IsSent checks if notification is sent
func (n *Notification) IsSent() bool {
	return n.Status == NotificationStatusSent || 
		   n.Status == NotificationStatusDelivered ||
		   n.Status == NotificationStatusRead
}

// IsRead checks if notification is read
func (n *Notification) IsRead() bool {
	return n.Status == NotificationStatusRead
}

// CanRetry checks if notification can be retried
func (n *Notification) CanRetry() bool {
	return n.Status == NotificationStatusFailed && 
		   n.RetryCount < n.MaxRetries
}

// MarkAsSent marks notification as sent
func (n *Notification) MarkAsSent() {
	now := time.Now()
	n.Status = NotificationStatusSent
	n.SentAt = &now
	n.UpdatedAt = now
}

// MarkAsDelivered marks notification as delivered
func (n *Notification) MarkAsDelivered() {
	now := time.Now()
	n.Status = NotificationStatusDelivered
	n.DeliveredAt = &now
	n.UpdatedAt = now
}

// MarkAsRead marks notification as read
func (n *Notification) MarkAsRead() {
	now := time.Now()
	n.Status = NotificationStatusRead
	n.ReadAt = &now
	n.UpdatedAt = now
}

// MarkAsFailed marks notification as failed
func (n *Notification) MarkAsFailed(errorMessage, errorCode string) {
	now := time.Now()
	n.Status = NotificationStatusFailed
	n.ErrorMessage = errorMessage
	n.ErrorCode = errorCode
	n.RetryCount++
	
	// Schedule next retry if possible
	if n.CanRetry() {
		// Exponential backoff: 1min, 5min, 15min
		retryDelays := []time.Duration{
			1 * time.Minute,
			5 * time.Minute,
			15 * time.Minute,
		}
		
		if n.RetryCount <= len(retryDelays) {
			nextRetry := now.Add(retryDelays[n.RetryCount-1])
			n.NextRetryAt = &nextRetry
		}
	}
	
	n.UpdatedAt = now
}

// NotificationTemplate represents email/SMS templates
type NotificationTemplate struct {
	ID          uuid.UUID            `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string               `json:"name" gorm:"uniqueIndex;not null" validate:"required"`
	Type        NotificationType     `json:"type" gorm:"not null"`
	Category    NotificationCategory `json:"category" gorm:"not null"`
	
	// Template content
	Subject     string               `json:"subject"`                             // For email
	Body        string               `json:"body" gorm:"type:text;not null"`     // Template body with placeholders
	Variables   string               `json:"variables" gorm:"type:text"`         // JSON array of available variables
	
	// Settings
	IsActive    bool                 `json:"is_active" gorm:"default:true"`
	IsDefault   bool                 `json:"is_default" gorm:"default:false"`
	Language    string               `json:"language" gorm:"default:'en'"`
	
	// Metadata
	Description string               `json:"description"`
	CreatedBy   uuid.UUID            `json:"created_by" gorm:"type:uuid"`
	CreatedAt   time.Time            `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time            `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for NotificationTemplate entity
func (NotificationTemplate) TableName() string {
	return "notification_templates"
}

// NotificationPreference represents user notification preferences
type NotificationPreference struct {
	ID                    uuid.UUID            `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID                uuid.UUID            `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	User                  User                 `json:"user,omitempty" gorm:"foreignKey:UserID"`
	
	// Email preferences
	EmailEnabled          bool                 `json:"email_enabled" gorm:"default:true"`
	EmailOrderUpdates     bool                 `json:"email_order_updates" gorm:"default:true"`
	EmailPaymentUpdates   bool                 `json:"email_payment_updates" gorm:"default:true"`
	EmailShippingUpdates  bool                 `json:"email_shipping_updates" gorm:"default:true"`
	EmailPromotions       bool                 `json:"email_promotions" gorm:"default:true"`
	EmailNewsletter       bool                 `json:"email_newsletter" gorm:"default:false"`
	EmailReviewReminders  bool                 `json:"email_review_reminders" gorm:"default:true"`
	
	// SMS preferences
	SMSEnabled            bool                 `json:"sms_enabled" gorm:"default:false"`
	SMSOrderUpdates       bool                 `json:"sms_order_updates" gorm:"default:false"`
	SMSPaymentUpdates     bool                 `json:"sms_payment_updates" gorm:"default:false"`
	SMSShippingUpdates    bool                 `json:"sms_shipping_updates" gorm:"default:false"`
	SMSSecurityAlerts     bool                 `json:"sms_security_alerts" gorm:"default:true"`
	
	// Push notification preferences
	PushEnabled           bool                 `json:"push_enabled" gorm:"default:true"`
	PushOrderUpdates      bool                 `json:"push_order_updates" gorm:"default:true"`
	PushPaymentUpdates    bool                 `json:"push_payment_updates" gorm:"default:true"`
	PushShippingUpdates   bool                 `json:"push_shipping_updates" gorm:"default:true"`
	PushPromotions        bool                 `json:"push_promotions" gorm:"default:false"`
	PushReviewReminders   bool                 `json:"push_review_reminders" gorm:"default:true"`
	
	// In-app notification preferences
	InAppEnabled          bool                 `json:"in_app_enabled" gorm:"default:true"`
	InAppOrderUpdates     bool                 `json:"in_app_order_updates" gorm:"default:true"`
	InAppPaymentUpdates   bool                 `json:"in_app_payment_updates" gorm:"default:true"`
	InAppShippingUpdates  bool                 `json:"in_app_shipping_updates" gorm:"default:true"`
	InAppPromotions       bool                 `json:"in_app_promotions" gorm:"default:true"`
	InAppSystemUpdates    bool                 `json:"in_app_system_updates" gorm:"default:true"`
	
	// Frequency settings
	DigestFrequency       string               `json:"digest_frequency" gorm:"default:'daily'"` // immediate, daily, weekly
	QuietHoursStart       string               `json:"quiet_hours_start"`                       // HH:MM format
	QuietHoursEnd         string               `json:"quiet_hours_end"`                         // HH:MM format
	Timezone              string               `json:"timezone" gorm:"default:'UTC'"`
	
	// Metadata
	CreatedAt             time.Time            `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time            `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for NotificationPreference entity
func (NotificationPreference) TableName() string {
	return "notification_preferences"
}

// IsNotificationEnabled checks if a specific notification type is enabled
func (np *NotificationPreference) IsNotificationEnabled(notificationType NotificationType, category NotificationCategory) bool {
	switch notificationType {
	case NotificationTypeEmail:
		if !np.EmailEnabled {
			return false
		}
		switch category {
		case NotificationCategoryOrder:
			return np.EmailOrderUpdates
		case NotificationCategoryPayment:
			return np.EmailPaymentUpdates
		case NotificationCategoryShipping:
			return np.EmailShippingUpdates
		case NotificationCategoryPromotion:
			return np.EmailPromotions
		case NotificationCategoryReview:
			return np.EmailReviewReminders
		default:
			return true
		}
		
	case NotificationTypeSMS:
		if !np.SMSEnabled {
			return false
		}
		switch category {
		case NotificationCategoryOrder:
			return np.SMSOrderUpdates
		case NotificationCategoryPayment:
			return np.SMSPaymentUpdates
		case NotificationCategoryShipping:
			return np.SMSShippingUpdates
		case NotificationCategoryAccount:
			return np.SMSSecurityAlerts
		default:
			return false
		}
		
	case NotificationTypePush:
		if !np.PushEnabled {
			return false
		}
		switch category {
		case NotificationCategoryOrder:
			return np.PushOrderUpdates
		case NotificationCategoryPayment:
			return np.PushPaymentUpdates
		case NotificationCategoryShipping:
			return np.PushShippingUpdates
		case NotificationCategoryPromotion:
			return np.PushPromotions
		case NotificationCategoryReview:
			return np.PushReviewReminders
		default:
			return true
		}
		
	case NotificationTypeInApp:
		if !np.InAppEnabled {
			return false
		}
		switch category {
		case NotificationCategoryOrder:
			return np.InAppOrderUpdates
		case NotificationCategoryPayment:
			return np.InAppPaymentUpdates
		case NotificationCategoryShipping:
			return np.InAppShippingUpdates
		case NotificationCategoryPromotion:
			return np.InAppPromotions
		case NotificationCategorySystem:
			return np.InAppSystemUpdates
		default:
			return true
		}
	}
	
	return false
}

// NotificationQueue represents queued notifications for batch processing
type NotificationQueue struct {
	ID            uuid.UUID            `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	NotificationID uuid.UUID           `json:"notification_id" gorm:"type:uuid;not null;index"`
	Notification  Notification         `json:"notification,omitempty" gorm:"foreignKey:NotificationID"`
	
	// Queue details
	Priority      NotificationPriority `json:"priority" gorm:"default:'normal'"`
	ScheduledAt   time.Time            `json:"scheduled_at" gorm:"not null"`
	ProcessedAt   *time.Time           `json:"processed_at"`
	
	// Processing information
	WorkerID      string               `json:"worker_id"`
	ProcessingStartedAt *time.Time     `json:"processing_started_at"`
	
	// Metadata
	CreatedAt     time.Time            `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time            `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for NotificationQueue entity
func (NotificationQueue) TableName() string {
	return "notification_queue"
}

// IsProcessed checks if notification is processed
func (nq *NotificationQueue) IsProcessed() bool {
	return nq.ProcessedAt != nil
}

// MarkAsProcessed marks notification as processed
func (nq *NotificationQueue) MarkAsProcessed(workerID string) {
	now := time.Now()
	nq.ProcessedAt = &now
	nq.WorkerID = workerID
	nq.UpdatedAt = now
}
