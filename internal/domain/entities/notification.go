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

// NotificationChannel represents the delivery channel
type NotificationChannel string

const (
	NotificationChannelEmail  NotificationChannel = "email"
	NotificationChannelSMS    NotificationChannel = "sms"
	NotificationChannelPush   NotificationChannel = "push"
	NotificationChannelInApp  NotificationChannel = "in_app"
	NotificationChannelWebhook NotificationChannel = "webhook"
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

// NotificationTemplate represents a notification template
type NotificationTemplate struct {
	ID          uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Type        NotificationType `json:"type" gorm:"not null;index"`
	Channel     NotificationChannel `json:"channel" gorm:"not null;index"`
	Name        string           `json:"name" gorm:"not null"`
	Subject     string           `json:"subject,omitempty"`
	Body        string           `json:"body" gorm:"not null"`
	Variables   []string         `json:"variables,omitempty" gorm:"type:jsonb"`
	IsActive    bool             `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for NotificationTemplate entity
func (NotificationTemplate) TableName() string {
	return "notification_templates"
}

// NotificationPreferences represents user notification preferences
type NotificationPreferences struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID            uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	EmailEnabled      bool      `json:"email_enabled" gorm:"default:true"`
	SMSEnabled        bool      `json:"sms_enabled" gorm:"default:false"`
	PushEnabled       bool      `json:"push_enabled" gorm:"default:true"`
	InAppEnabled      bool      `json:"in_app_enabled" gorm:"default:true"`
	OrderUpdates      bool      `json:"order_updates" gorm:"default:true"`
	PromotionalEmails bool      `json:"promotional_emails" gorm:"default:true"`
	SecurityAlerts    bool      `json:"security_alerts" gorm:"default:true"`
	NewsletterEnabled bool      `json:"newsletter_enabled" gorm:"default:false"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for NotificationPreferences entity
func (NotificationPreferences) TableName() string {
	return "notification_preferences"
}

// DeliveryStatus represents notification delivery status
type DeliveryStatus string

const (
	DeliveryStatusPending   DeliveryStatus = "pending"
	DeliveryStatusSent      DeliveryStatus = "sent"
	DeliveryStatusDelivered DeliveryStatus = "delivered"
	DeliveryStatusFailed    DeliveryStatus = "failed"
	DeliveryStatusBounced   DeliveryStatus = "bounced"
	DeliveryStatusOpened    DeliveryStatus = "opened"
	DeliveryStatusClicked   DeliveryStatus = "clicked"
)

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



// IsNotificationEnabled checks if a specific notification type is enabled
func (np *NotificationPreferences) IsNotificationEnabled(notificationType NotificationType, category NotificationCategory) bool {
	switch notificationType {
	case NotificationTypeEmail:
		if !np.EmailEnabled {
			return false
		}
		switch category {
		case NotificationCategoryOrder:
			return np.OrderUpdates
		case NotificationCategoryPayment:
			return np.OrderUpdates
		case NotificationCategoryShipping:
			return np.OrderUpdates
		case NotificationCategoryPromotion:
			return np.PromotionalEmails
		case NotificationCategoryReview:
			return np.OrderUpdates
		default:
			return true
		}
		
	case NotificationTypeSMS:
		if !np.SMSEnabled {
			return false
		}
		switch category {
		case NotificationCategoryOrder:
			return np.OrderUpdates
		case NotificationCategoryPayment:
			return np.OrderUpdates
		case NotificationCategoryShipping:
			return np.OrderUpdates
		case NotificationCategoryAccount:
			return np.SecurityAlerts
		default:
			return false
		}
		
	case NotificationTypePush:
		if !np.PushEnabled {
			return false
		}
		switch category {
		case NotificationCategoryOrder:
			return np.OrderUpdates
		case NotificationCategoryPayment:
			return np.OrderUpdates
		case NotificationCategoryShipping:
			return np.OrderUpdates
		case NotificationCategoryPromotion:
			return np.PromotionalEmails
		case NotificationCategoryReview:
			return np.OrderUpdates
		default:
			return true
		}
		
	case NotificationTypeInApp:
		if !np.InAppEnabled {
			return false
		}
		switch category {
		case NotificationCategoryOrder:
			return np.OrderUpdates
		case NotificationCategoryPayment:
			return np.OrderUpdates
		case NotificationCategoryShipping:
			return np.OrderUpdates
		case NotificationCategoryPromotion:
			return np.PromotionalEmails
		case NotificationCategorySystem:
			return np.OrderUpdates
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
