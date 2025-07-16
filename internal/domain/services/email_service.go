package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

// EmailService defines the interface for email operations
type EmailService interface {
	// Send email operations
	SendEmail(ctx context.Context, email *entities.Email) error
	SendTemplateEmail(ctx context.Context, templateName string, to, toName string, data map[string]interface{}) error
	
	// Bulk operations
	SendBulkEmails(ctx context.Context, emails []*entities.Email) error
	
	// Template operations
	RenderTemplate(ctx context.Context, templateName string, data map[string]interface{}) (subject, bodyText, bodyHTML string, err error)
	
	// Delivery tracking
	TrackDelivery(ctx context.Context, externalID string, status entities.EmailStatus) error
	TrackOpen(ctx context.Context, emailID uuid.UUID) error
	TrackClick(ctx context.Context, emailID uuid.UUID, url string) error
	
	// Retry operations
	RetryFailedEmails(ctx context.Context) error
	
	// Validation
	ValidateEmailAddress(email string) error
}

// EmailProvider defines the interface for email providers (SMTP, SendGrid, etc.)
type EmailProvider interface {
	SendEmail(ctx context.Context, email *entities.Email) (externalID string, err error)
	SendBulkEmails(ctx context.Context, emails []*entities.Email) (results map[uuid.UUID]string, err error)
	ValidateConfiguration() error
}

type emailService struct {
	emailRepo         repositories.EmailRepository
	templateRepo      repositories.EmailTemplateRepository
	subscriptionRepo  repositories.EmailSubscriptionRepository
	provider          EmailProvider
	defaultFromEmail  string
	defaultFromName   string
}

// NewEmailService creates a new email service
func NewEmailService(
	emailRepo repositories.EmailRepository,
	templateRepo repositories.EmailTemplateRepository,
	subscriptionRepo repositories.EmailSubscriptionRepository,
	provider EmailProvider,
	defaultFromEmail, defaultFromName string,
) EmailService {
	return &emailService{
		emailRepo:        emailRepo,
		templateRepo:     templateRepo,
		subscriptionRepo: subscriptionRepo,
		provider:         provider,
		defaultFromEmail: defaultFromEmail,
		defaultFromName:  defaultFromName,
	}
}

// SendEmail sends an email
func (s *emailService) SendEmail(ctx context.Context, email *entities.Email) error {
	// Validate email address
	if err := s.ValidateEmailAddress(email.ToEmail); err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}

	// Set defaults if not provided
	if email.FromEmail == "" {
		email.FromEmail = s.defaultFromEmail
	}
	if email.FromName == "" {
		email.FromName = s.defaultFromName
	}
	if email.Priority == "" {
		email.Priority = entities.EmailPriorityNormal
	}

	// Check subscription preferences if user is provided
	if email.UserID != nil {
		subscription, err := s.subscriptionRepo.GetByUserID(ctx, *email.UserID)
		if err == nil && !subscription.IsSubscribedTo(email.Type) {
			return fmt.Errorf("user is not subscribed to %s emails", email.Type)
		}
	}

	// Save email to database
	if err := s.emailRepo.Create(ctx, email); err != nil {
		return fmt.Errorf("failed to save email: %w", err)
	}

	// Send email via provider
	externalID, err := s.provider.SendEmail(ctx, email)
	if err != nil {
		email.MarkAsFailed(err.Error())
		_ = s.emailRepo.Update(ctx, email)
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Mark as sent
	email.MarkAsSent(externalID)
	if err := s.emailRepo.Update(ctx, email); err != nil {
		return fmt.Errorf("failed to update email status: %w", err)
	}

	return nil
}

// SendTemplateEmail sends an email using a template
func (s *emailService) SendTemplateEmail(ctx context.Context, templateName string, to, toName string, data map[string]interface{}) error {
	// Render template
	subject, bodyText, bodyHTML, err := s.RenderTemplate(ctx, templateName, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Get template to determine type
	template, err := s.templateRepo.GetByName(ctx, templateName)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	// Create email
	email := &entities.Email{
		ID:           uuid.New(),
		Type:         template.Type,
		Priority:     entities.EmailPriorityNormal,
		Status:       entities.EmailStatusPending,
		ToEmail:      to,
		ToName:       toName,
		FromEmail:    s.defaultFromEmail,
		FromName:     s.defaultFromName,
		Subject:      subject,
		BodyText:     bodyText,
		BodyHTML:     bodyHTML,
		TemplateID:   templateName,
		TemplateData: data,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Extract user ID if provided in data
	if userIDStr, ok := data["user_id"].(string); ok {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			email.UserID = &userID
		}
	}

	// Extract order ID if provided in data
	if orderIDStr, ok := data["order_id"].(string); ok {
		if orderID, err := uuid.Parse(orderIDStr); err == nil {
			email.OrderID = &orderID
		}
	}

	return s.SendEmail(ctx, email)
}

// SendBulkEmails sends multiple emails
func (s *emailService) SendBulkEmails(ctx context.Context, emails []*entities.Email) error {
	if len(emails) == 0 {
		return nil
	}

	// Validate and prepare emails
	validEmails := make([]*entities.Email, 0, len(emails))
	for _, email := range emails {
		if err := s.ValidateEmailAddress(email.ToEmail); err != nil {
			continue // Skip invalid emails
		}

		// Set defaults
		if email.FromEmail == "" {
			email.FromEmail = s.defaultFromEmail
		}
		if email.FromName == "" {
			email.FromName = s.defaultFromName
		}
		if email.Priority == "" {
			email.Priority = entities.EmailPriorityNormal
		}

		validEmails = append(validEmails, email)
	}

	if len(validEmails) == 0 {
		return fmt.Errorf("no valid emails to send")
	}

	// Save emails to database
	for _, email := range validEmails {
		if err := s.emailRepo.Create(ctx, email); err != nil {
			return fmt.Errorf("failed to save email %s: %w", email.ID, err)
		}
	}

	// Send emails via provider
	results, err := s.provider.SendBulkEmails(ctx, validEmails)
	if err != nil {
		return fmt.Errorf("failed to send bulk emails: %w", err)
	}

	// Update email statuses
	for _, email := range validEmails {
		if externalID, ok := results[email.ID]; ok {
			email.MarkAsSent(externalID)
		} else {
			email.MarkAsFailed("failed to send via provider")
		}
		_ = s.emailRepo.Update(ctx, email)
	}

	return nil
}

// RenderTemplate renders an email template with data
func (s *emailService) RenderTemplate(ctx context.Context, templateName string, data map[string]interface{}) (subject, bodyText, bodyHTML string, err error) {
	template, err := s.templateRepo.GetByName(ctx, templateName)
	if err != nil {
		return "", "", "", fmt.Errorf("template not found: %w", err)
	}

	if !template.IsActive {
		return "", "", "", fmt.Errorf("template %s is not active", templateName)
	}

	// Simple template rendering (in production, use a proper template engine)
	subject = s.renderString(template.Subject, data)
	bodyText = s.renderString(template.BodyText, data)
	bodyHTML = s.renderString(template.BodyHTML, data)

	return subject, bodyText, bodyHTML, nil
}

// renderString performs simple variable substitution
func (s *emailService) renderString(template string, data map[string]interface{}) string {
	result := template
	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}

// TrackDelivery tracks email delivery status
func (s *emailService) TrackDelivery(ctx context.Context, externalID string, status entities.EmailStatus) error {
	email, err := s.emailRepo.GetByExternalID(ctx, externalID)
	if err != nil {
		return fmt.Errorf("email not found: %w", err)
	}

	switch status {
	case entities.EmailStatusDelivered:
		email.MarkAsDelivered()
	case entities.EmailStatusBounced:
		email.Status = entities.EmailStatusBounced
		now := time.Now()
		email.BouncedAt = &now
		email.UpdatedAt = now
	case entities.EmailStatusFailed:
		email.MarkAsFailed("delivery failed")
	}

	return s.emailRepo.Update(ctx, email)
}

// TrackOpen tracks email open
func (s *emailService) TrackOpen(ctx context.Context, emailID uuid.UUID) error {
	email, err := s.emailRepo.GetByID(ctx, emailID)
	if err != nil {
		return fmt.Errorf("email not found: %w", err)
	}

	email.MarkAsOpened()
	return s.emailRepo.Update(ctx, email)
}

// TrackClick tracks email click
func (s *emailService) TrackClick(ctx context.Context, emailID uuid.UUID, url string) error {
	email, err := s.emailRepo.GetByID(ctx, emailID)
	if err != nil {
		return fmt.Errorf("email not found: %w", err)
	}

	email.MarkAsClicked()
	
	// Store click metadata
	if email.Metadata == nil {
		email.Metadata = make(map[string]interface{})
	}
	email.Metadata["clicked_url"] = url
	email.Metadata["clicked_at"] = time.Now()

	return s.emailRepo.Update(ctx, email)
}

// RetryFailedEmails retries failed emails that can be retried
func (s *emailService) RetryFailedEmails(ctx context.Context) error {
	emails, err := s.emailRepo.GetRetryableEmails(ctx)
	if err != nil {
		return fmt.Errorf("failed to get retryable emails: %w", err)
	}

	for _, email := range emails {
		if !email.CanRetry() {
			continue
		}

		externalID, err := s.provider.SendEmail(ctx, email)
		if err != nil {
			email.MarkAsFailed(err.Error())
		} else {
			email.MarkAsSent(externalID)
		}
		_ = s.emailRepo.Update(ctx, email)
	}

	return nil
}

// ValidateEmailAddress validates an email address
func (s *emailService) ValidateEmailAddress(email string) error {
	if email == "" {
		return fmt.Errorf("email address is required")
	}

	if len(email) > 254 {
		return fmt.Errorf("email address is too long")
	}

	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email format")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email format")
	}

	localPart := parts[0]
	domainPart := parts[1]

	if len(localPart) == 0 || len(localPart) > 64 {
		return fmt.Errorf("invalid email local part")
	}

	if len(domainPart) == 0 || len(domainPart) > 253 {
		return fmt.Errorf("invalid email domain part")
	}

	if !strings.Contains(domainPart, ".") {
		return fmt.Errorf("invalid email domain format")
	}

	return nil
}
