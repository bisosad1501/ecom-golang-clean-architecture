package entities

import (
	"time"

	"github.com/google/uuid"
)

// EmailStatus represents the status of an email
type EmailStatus string

const (
	EmailStatusPending   EmailStatus = "pending"
	EmailStatusSent      EmailStatus = "sent"
	EmailStatusFailed    EmailStatus = "failed"
	EmailStatusDelivered EmailStatus = "delivered"
	EmailStatusBounced   EmailStatus = "bounced"
	EmailStatusOpened    EmailStatus = "opened"
	EmailStatusClicked   EmailStatus = "clicked"
)

// EmailType represents the type of email
type EmailType string

const (
	EmailTypeWelcome           EmailType = "welcome"
	EmailTypeOrderConfirmation EmailType = "order_confirmation"
	EmailTypeOrderShipped      EmailType = "order_shipped"
	EmailTypeOrderDelivered    EmailType = "order_delivered"
	EmailTypeOrderCancelled    EmailType = "order_cancelled"
	EmailTypePasswordReset     EmailType = "password_reset"
	EmailTypeAccountActivation EmailType = "account_activation"
	EmailTypeAbandonedCart     EmailType = "abandoned_cart"
	EmailTypeReviewRequest     EmailType = "review_request"
	EmailTypePromotion         EmailType = "promotion"
	EmailTypeNewsletter        EmailType = "newsletter"
	EmailTypeSupport           EmailType = "support"
	EmailTypeRefund            EmailType = "refund"
	EmailTypeLowStock          EmailType = "low_stock"
)

// EmailPriority represents the priority of an email
type EmailPriority string

const (
	EmailPriorityLow    EmailPriority = "low"
	EmailPriorityNormal EmailPriority = "normal"
	EmailPriorityHigh   EmailPriority = "high"
	EmailPriorityUrgent EmailPriority = "urgent"
)

// Email represents an email in the system
type Email struct {
	ID          uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Type        EmailType     `json:"type" gorm:"not null;index"`
	Priority    EmailPriority `json:"priority" gorm:"default:'normal'"`
	Status      EmailStatus   `json:"status" gorm:"default:'pending';index"`
	
	// Recipient information
	ToEmail     string `json:"to_email" gorm:"not null;index"`
	ToName      string `json:"to_name"`
	FromEmail   string `json:"from_email" gorm:"not null"`
	FromName    string `json:"from_name"`
	ReplyToEmail string `json:"reply_to_email"`
	
	// Email content
	Subject     string `json:"subject" gorm:"not null"`
	BodyText    string `json:"body_text" gorm:"type:text"`
	BodyHTML    string `json:"body_html" gorm:"type:text"`
	
	// Template information
	TemplateID   string                 `json:"template_id"`
	TemplateData map[string]interface{} `json:"template_data" gorm:"type:jsonb"`
	
	// Reference information
	UserID    *uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	OrderID   *uuid.UUID `json:"order_id" gorm:"type:uuid;index"`
	ProductID *uuid.UUID `json:"product_id" gorm:"type:uuid;index"`
	
	// Delivery tracking
	SentAt       *time.Time `json:"sent_at"`
	DeliveredAt  *time.Time `json:"delivered_at"`
	OpenedAt     *time.Time `json:"opened_at"`
	ClickedAt    *time.Time `json:"clicked_at"`
	BouncedAt    *time.Time `json:"bounced_at"`
	
	// Retry information
	RetryCount   int        `json:"retry_count" gorm:"default:0"`
	MaxRetries   int        `json:"max_retries" gorm:"default:3"`
	NextRetryAt  *time.Time `json:"next_retry_at"`
	
	// Error information
	ErrorMessage string `json:"error_message" gorm:"type:text"`
	
	// External provider information
	ExternalID       string `json:"external_id"`
	ExternalProvider string `json:"external_provider"`
	
	// Metadata
	Metadata  map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	CreatedAt time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Order   *Order   `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

// TableName returns the table name for Email entity
func (Email) TableName() string {
	return "emails"
}

// MarkAsSent marks the email as sent
func (e *Email) MarkAsSent(externalID string) {
	e.Status = EmailStatusSent
	e.ExternalID = externalID
	now := time.Now()
	e.SentAt = &now
	e.UpdatedAt = now
}

// MarkAsDelivered marks the email as delivered
func (e *Email) MarkAsDelivered() {
	e.Status = EmailStatusDelivered
	now := time.Now()
	e.DeliveredAt = &now
	e.UpdatedAt = now
}

// MarkAsFailed marks the email as failed
func (e *Email) MarkAsFailed(errorMessage string) {
	e.Status = EmailStatusFailed
	e.ErrorMessage = errorMessage
	e.RetryCount++
	e.UpdatedAt = time.Now()
	
	// Set next retry time if retries are available
	if e.RetryCount < e.MaxRetries {
		nextRetry := time.Now().Add(time.Duration(e.RetryCount*e.RetryCount) * time.Hour) // Exponential backoff
		e.NextRetryAt = &nextRetry
	}
}

// MarkAsOpened marks the email as opened
func (e *Email) MarkAsOpened() {
	if e.Status == EmailStatusDelivered || e.Status == EmailStatusSent {
		e.Status = EmailStatusOpened
		now := time.Now()
		e.OpenedAt = &now
		e.UpdatedAt = now
	}
}

// MarkAsClicked marks the email as clicked
func (e *Email) MarkAsClicked() {
	e.Status = EmailStatusClicked
	now := time.Now()
	e.ClickedAt = &now
	e.UpdatedAt = now
}

// CanRetry checks if the email can be retried
func (e *Email) CanRetry() bool {
	return e.Status == EmailStatusFailed && e.RetryCount < e.MaxRetries
}

// IsDelivered checks if the email was delivered
func (e *Email) IsDelivered() bool {
	return e.Status == EmailStatusDelivered || e.Status == EmailStatusOpened || e.Status == EmailStatusClicked
}

// GetDeliveryTime returns the time when email was delivered
func (e *Email) GetDeliveryTime() *time.Time {
	if e.DeliveredAt != nil {
		return e.DeliveredAt
	}
	return e.SentAt
}

// EmailTemplate represents an email template
type EmailTemplate struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null"`
	Type        EmailType `json:"type" gorm:"not null;index"`
	Subject     string    `json:"subject" gorm:"not null"`
	BodyText    string    `json:"body_text" gorm:"type:text"`
	BodyHTML    string    `json:"body_html" gorm:"type:text"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Version     int       `json:"version" gorm:"default:1"`
	Description string    `json:"description" gorm:"type:text"`
	
	// Template variables documentation
	Variables map[string]interface{} `json:"variables" gorm:"type:jsonb"`
	
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for EmailTemplate entity
func (EmailTemplate) TableName() string {
	return "email_templates"
}

// EmailSubscription represents user email subscription preferences
type EmailSubscription struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	
	// Subscription preferences
	Newsletter      bool `json:"newsletter" gorm:"default:true"`
	Promotions      bool `json:"promotions" gorm:"default:true"`
	OrderUpdates    bool `json:"order_updates" gorm:"default:true"`
	ReviewRequests  bool `json:"review_requests" gorm:"default:true"`
	AbandonedCart   bool `json:"abandoned_cart" gorm:"default:true"`
	Support         bool `json:"support" gorm:"default:true"`
	
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for EmailSubscription entity
func (EmailSubscription) TableName() string {
	return "email_subscriptions"
}

// IsSubscribedTo checks if user is subscribed to a specific email type
func (es *EmailSubscription) IsSubscribedTo(emailType EmailType) bool {
	switch emailType {
	case EmailTypeNewsletter:
		return es.Newsletter
	case EmailTypePromotion:
		return es.Promotions
	case EmailTypeOrderConfirmation, EmailTypeOrderShipped, EmailTypeOrderDelivered, EmailTypeOrderCancelled:
		return es.OrderUpdates
	case EmailTypeReviewRequest:
		return es.ReviewRequests
	case EmailTypeAbandonedCart:
		return es.AbandonedCart
	case EmailTypeSupport:
		return es.Support
	default:
		return true // Default to subscribed for system emails
	}
}
