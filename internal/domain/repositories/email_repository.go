package repositories

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// EmailRepository defines the interface for email data operations
type EmailRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, email *entities.Email) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Email, error)
	GetByExternalID(ctx context.Context, externalID string) (*entities.Email, error)
	Update(ctx context.Context, email *entities.Email) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// Query operations
	List(ctx context.Context, offset, limit int) ([]*entities.Email, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entities.Email, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.Email, error)
	GetByType(ctx context.Context, emailType entities.EmailType, offset, limit int) ([]*entities.Email, error)
	GetByStatus(ctx context.Context, status entities.EmailStatus, offset, limit int) ([]*entities.Email, error)
	
	// Retry operations
	GetRetryableEmails(ctx context.Context) ([]*entities.Email, error)
	GetFailedEmails(ctx context.Context, since time.Time) ([]*entities.Email, error)
	
	// Analytics operations
	GetEmailStats(ctx context.Context, since time.Time) (*EmailStats, error)
	GetDeliveryRate(ctx context.Context, emailType entities.EmailType, since time.Time) (float64, error)
	GetOpenRate(ctx context.Context, emailType entities.EmailType, since time.Time) (float64, error)
	GetClickRate(ctx context.Context, emailType entities.EmailType, since time.Time) (float64, error)
	
	// Bulk operations
	CreateBatch(ctx context.Context, emails []*entities.Email) error
	UpdateBatch(ctx context.Context, emails []*entities.Email) error
	
	// Search operations
	Search(ctx context.Context, query EmailSearchQuery) ([]*entities.Email, int, error)
}

// EmailTemplateRepository defines the interface for email template data operations
type EmailTemplateRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, template *entities.EmailTemplate) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.EmailTemplate, error)
	GetByName(ctx context.Context, name string) (*entities.EmailTemplate, error)
	Update(ctx context.Context, template *entities.EmailTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// Query operations
	List(ctx context.Context, offset, limit int) ([]*entities.EmailTemplate, error)
	GetByType(ctx context.Context, emailType entities.EmailType) ([]*entities.EmailTemplate, error)
	GetActive(ctx context.Context) ([]*entities.EmailTemplate, error)
	
	// Version operations
	GetLatestVersion(ctx context.Context, name string) (*entities.EmailTemplate, error)
	GetByVersion(ctx context.Context, name string, version int) (*entities.EmailTemplate, error)
}

// EmailSubscriptionRepository defines the interface for email subscription data operations
type EmailSubscriptionRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, subscription *entities.EmailSubscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.EmailSubscription, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.EmailSubscription, error)
	Update(ctx context.Context, subscription *entities.EmailSubscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// Query operations
	List(ctx context.Context, offset, limit int) ([]*entities.EmailSubscription, error)
	GetSubscribedUsers(ctx context.Context, emailType entities.EmailType) ([]uuid.UUID, error)
	GetUnsubscribedUsers(ctx context.Context, emailType entities.EmailType) ([]uuid.UUID, error)
	
	// Bulk operations
	UpdateSubscriptions(ctx context.Context, userID uuid.UUID, subscriptions map[entities.EmailType]bool) error
}

// EmailSearchQuery represents search parameters for emails
type EmailSearchQuery struct {
	// Basic filters
	UserID    *uuid.UUID             `json:"user_id"`
	OrderID   *uuid.UUID             `json:"order_id"`
	Type      *entities.EmailType    `json:"type"`
	Status    *entities.EmailStatus  `json:"status"`
	Priority  *entities.EmailPriority `json:"priority"`
	
	// Email content filters
	ToEmail   string `json:"to_email"`
	Subject   string `json:"subject"`
	
	// Date filters
	CreatedAfter  *time.Time `json:"created_after"`
	CreatedBefore *time.Time `json:"created_before"`
	SentAfter     *time.Time `json:"sent_after"`
	SentBefore    *time.Time `json:"sent_before"`
	
	// Delivery tracking filters
	IsDelivered bool `json:"is_delivered"`
	IsOpened    bool `json:"is_opened"`
	IsClicked   bool `json:"is_clicked"`
	IsBounced   bool `json:"is_bounced"`
	
	// Retry filters
	HasFailed   bool `json:"has_failed"`
	CanRetry    bool `json:"can_retry"`
	
	// Pagination
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	
	// Sorting
	SortBy    string `json:"sort_by"`    // created_at, sent_at, updated_at
	SortOrder string `json:"sort_order"` // asc, desc
}

// EmailStats represents email statistics
type EmailStats struct {
	// Volume stats
	TotalSent      int64 `json:"total_sent"`
	TotalDelivered int64 `json:"total_delivered"`
	TotalOpened    int64 `json:"total_opened"`
	TotalClicked   int64 `json:"total_clicked"`
	TotalBounced   int64 `json:"total_bounced"`
	TotalFailed    int64 `json:"total_failed"`
	
	// Rate stats
	DeliveryRate float64 `json:"delivery_rate"` // delivered / sent
	OpenRate     float64 `json:"open_rate"`     // opened / delivered
	ClickRate    float64 `json:"click_rate"`    // clicked / delivered
	BounceRate   float64 `json:"bounce_rate"`   // bounced / sent
	FailureRate  float64 `json:"failure_rate"`  // failed / sent
	
	// Type breakdown
	TypeStats map[entities.EmailType]TypeStats `json:"type_stats"`
	
	// Time period
	Since time.Time `json:"since"`
	Until time.Time `json:"until"`
}

// TypeStats represents statistics for a specific email type
type TypeStats struct {
	Sent      int64   `json:"sent"`
	Delivered int64   `json:"delivered"`
	Opened    int64   `json:"opened"`
	Clicked   int64   `json:"clicked"`
	Bounced   int64   `json:"bounced"`
	Failed    int64   `json:"failed"`
	
	DeliveryRate float64 `json:"delivery_rate"`
	OpenRate     float64 `json:"open_rate"`
	ClickRate    float64 `json:"click_rate"`
	BounceRate   float64 `json:"bounce_rate"`
	FailureRate  float64 `json:"failure_rate"`
}

// EmailFilter represents filters for email queries
type EmailFilter struct {
	UserIDs   []uuid.UUID             `json:"user_ids"`
	OrderIDs  []uuid.UUID             `json:"order_ids"`
	Types     []entities.EmailType    `json:"types"`
	Statuses  []entities.EmailStatus  `json:"statuses"`
	Priorities []entities.EmailPriority `json:"priorities"`
	
	// Date ranges
	CreatedAfter  *time.Time `json:"created_after"`
	CreatedBefore *time.Time `json:"created_before"`
	SentAfter     *time.Time `json:"sent_after"`
	SentBefore    *time.Time `json:"sent_before"`
	
	// Content filters
	SubjectContains string `json:"subject_contains"`
	ToEmailContains string `json:"to_email_contains"`
	
	// Delivery status
	DeliveryStatus []string `json:"delivery_status"` // delivered, opened, clicked, bounced
	
	// Retry status
	CanRetry     *bool `json:"can_retry"`
	RetryCount   *int  `json:"retry_count"`
	MaxRetries   *int  `json:"max_retries"`
	
	// External provider
	ExternalProvider string `json:"external_provider"`
	
	// Template
	TemplateID string `json:"template_id"`
}

// EmailSortOptions represents sorting options for email queries
type EmailSortOptions struct {
	Field string `json:"field"` // created_at, sent_at, updated_at, priority
	Order string `json:"order"` // asc, desc
}

// PaginationOptions represents pagination options
type PaginationOptions struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// EmailQueryOptions combines all query options
type EmailQueryOptions struct {
	Filter     *EmailFilter       `json:"filter"`
	Sort       *EmailSortOptions  `json:"sort"`
	Pagination *PaginationOptions `json:"pagination"`
}
