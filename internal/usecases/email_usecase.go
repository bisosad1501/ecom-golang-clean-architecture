package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// EmailUseCase defines email-related operations
type EmailUseCase interface {
	// Email Templates
	GetEmailTemplates(ctx context.Context, req GetEmailTemplatesRequest) (*GetEmailTemplatesResponse, error)
	GetEmailTemplate(ctx context.Context, id uuid.UUID) (*EmailTemplateResponse, error)
	CreateEmailTemplate(ctx context.Context, req CreateEmailTemplateRequest) (*EmailTemplateResponse, error)
	UpdateEmailTemplate(ctx context.Context, id uuid.UUID, req UpdateEmailTemplateRequest) (*EmailTemplateResponse, error)
	DeleteEmailTemplate(ctx context.Context, id uuid.UUID) error

	// Email Sending
	SendEmail(ctx context.Context, req SendEmailRequest) error
	SendBulkEmail(ctx context.Context, req SendBulkEmailRequest) (*EmailBulkResponse, error)
	SendTemplateEmail(ctx context.Context, req SendTemplateEmailRequest) error
	SendAbandonedCartEmail(ctx context.Context, cartID uuid.UUID, userEmail string) error
}

// Request/Response types for Email Templates
type GetEmailTemplatesRequest struct {
	Search    string `json:"search" form:"search"`
	Type      string `json:"type" form:"type"`
	IsActive  *bool  `json:"is_active" form:"is_active"`
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order"`
	Page      int    `json:"page" form:"page"`
	Limit     int    `json:"limit" form:"limit"`
}

type GetEmailTemplatesResponse struct {
	Templates  []EmailTemplateResponse `json:"templates"`
	Total      int64                   `json:"total"`
	Pagination *PaginationInfo         `json:"pagination"`
}

type EmailTemplateResponse struct {
	ID            uuid.UUID               `json:"id"`
	Name          string                  `json:"name"`
	Subject       string                  `json:"subject"`
	Type          string                  `json:"type"`
	Content       string                  `json:"content"`
	PlainText     string                  `json:"plain_text"`
	IsActive      bool                    `json:"is_active"`
	Variables     []EmailVariableResponse `json:"variables"`
	UsageCount    int64                   `json:"usage_count"`
	CreatedBy     uuid.UUID               `json:"created_by"`
	CreatedByName string                  `json:"created_by_name"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

type EmailVariableResponse struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	Required     bool   `json:"required"`
	DefaultValue string `json:"default_value"`
}

type CreateEmailTemplateRequest struct {
	Name      string                  `json:"name" validate:"required"`
	Subject   string                  `json:"subject" validate:"required"`
	Type      string                  `json:"type" validate:"required"`
	Content   string                  `json:"content" validate:"required"`
	PlainText string                  `json:"plain_text"`
	IsActive  bool                    `json:"is_active"`
	Variables []EmailVariableResponse `json:"variables"`
	Settings  map[string]interface{}  `json:"settings"`
	CreatedBy uuid.UUID               `json:"created_by" validate:"required"`
}

type UpdateEmailTemplateRequest struct {
	Name      string                  `json:"name"`
	Subject   string                  `json:"subject"`
	Type      string                  `json:"type"`
	Content   string                  `json:"content"`
	PlainText string                  `json:"plain_text"`
	IsActive  bool                    `json:"is_active"`
	Variables []EmailVariableResponse `json:"variables"`
	Settings  map[string]interface{}  `json:"settings"`
}

// Request/Response types for Email Sending
type SendEmailRequest struct {
	To          []string               `json:"to" validate:"required"`
	CC          []string               `json:"cc"`
	BCC         []string               `json:"bcc"`
	Subject     string                 `json:"subject" validate:"required"`
	Content     string                 `json:"content" validate:"required"`
	PlainText   string                 `json:"plain_text"`
	Attachments []EmailAttachment      `json:"attachments"`
	Variables   map[string]interface{} `json:"variables"`
	Priority    string                 `json:"priority"` // high, normal, low
}

type SendBulkEmailRequest struct {
	Recipients  []EmailRecipient       `json:"recipients" validate:"required"`
	Subject     string                 `json:"subject" validate:"required"`
	Content     string                 `json:"content" validate:"required"`
	PlainText   string                 `json:"plain_text"`
	Attachments []EmailAttachment      `json:"attachments"`
	Variables   map[string]interface{} `json:"variables"`
	TemplateID  *uuid.UUID             `json:"template_id"`
	ScheduledAt *time.Time             `json:"scheduled_at"`
}

type SendTemplateEmailRequest struct {
	TemplateID  uuid.UUID              `json:"template_id" validate:"required"`
	To          []string               `json:"to" validate:"required"`
	CC          []string               `json:"cc"`
	BCC         []string               `json:"bcc"`
	Variables   map[string]interface{} `json:"variables"`
	ScheduledAt *time.Time             `json:"scheduled_at"`
}

type EmailRecipient struct {
	Email     string                 `json:"email" validate:"required,email"`
	Name      string                 `json:"name"`
	Variables map[string]interface{} `json:"variables"`
}

type EmailAttachment struct {
	Filename    string `json:"filename" validate:"required"`
	ContentType string `json:"content_type"`
	Content     []byte `json:"content" validate:"required"`
	Size        int64  `json:"size"`
}

type EmailBulkResponse struct {
	JobID       uuid.UUID  `json:"job_id"`
	TotalEmails int        `json:"total_emails"`
	Status      string     `json:"status"`
	ScheduledAt *time.Time `json:"scheduled_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
