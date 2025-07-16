package usecases

import (
	"context"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/domain/services"

	"github.com/google/uuid"
)

// EmailUseCase defines the interface for email business logic
type EmailUseCase interface {
	// Email operations
	SendWelcomeEmail(ctx context.Context, userID uuid.UUID) error
	SendOrderConfirmationEmail(ctx context.Context, orderID uuid.UUID) error
	SendOrderShippedEmail(ctx context.Context, orderID uuid.UUID) error
	SendOrderDeliveredEmail(ctx context.Context, orderID uuid.UUID) error
	SendOrderCancelledEmail(ctx context.Context, orderID uuid.UUID) error
	SendPasswordResetEmail(ctx context.Context, userID uuid.UUID, resetToken string) error
	SendAbandonedCartEmail(ctx context.Context, userID uuid.UUID) error
	SendReviewRequestEmail(ctx context.Context, userID, orderID uuid.UUID) error
	SendLowStockAlert(ctx context.Context, productID uuid.UUID) error

	// Template operations
	CreateTemplate(ctx context.Context, req CreateTemplateRequest) (*TemplateResponse, error)
	UpdateTemplate(ctx context.Context, id uuid.UUID, req UpdateTemplateRequest) (*TemplateResponse, error)
	GetTemplate(ctx context.Context, id uuid.UUID) (*TemplateResponse, error)
	ListTemplates(ctx context.Context, offset, limit int) ([]*TemplateResponse, error)
	DeleteTemplate(ctx context.Context, id uuid.UUID) error

	// Subscription operations
	UpdateSubscriptions(ctx context.Context, userID uuid.UUID, req UpdateSubscriptionsRequest) error
	GetSubscriptions(ctx context.Context, userID uuid.UUID) (*SubscriptionsResponse, error)

	// Analytics operations
	GetEmailStats(ctx context.Context, since time.Time) (*EmailStatsResponse, error)
	GetEmailHistory(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*EmailResponse, error)

	// Admin operations
	RetryFailedEmails(ctx context.Context) error
	GetFailedEmails(ctx context.Context, since time.Time) ([]*EmailResponse, error)
}

type emailUseCase struct {
	emailService     services.EmailService
	emailRepo        repositories.EmailRepository
	templateRepo     repositories.EmailTemplateRepository
	subscriptionRepo repositories.EmailSubscriptionRepository
	userRepo         repositories.UserRepository
	orderRepo        repositories.OrderRepository
	productRepo      repositories.ProductRepository
}

// NewEmailUseCase creates a new email use case
func NewEmailUseCase(
	emailService services.EmailService,
	emailRepo repositories.EmailRepository,
	templateRepo repositories.EmailTemplateRepository,
	subscriptionRepo repositories.EmailSubscriptionRepository,
	userRepo repositories.UserRepository,
	orderRepo repositories.OrderRepository,
	productRepo repositories.ProductRepository,
) EmailUseCase {
	return &emailUseCase{
		emailService:     emailService,
		emailRepo:        emailRepo,
		templateRepo:     templateRepo,
		subscriptionRepo: subscriptionRepo,
		userRepo:         userRepo,
		orderRepo:        orderRepo,
		productRepo:      productRepo,
	}
}

// SendWelcomeEmail sends a welcome email to a new user
func (uc *emailUseCase) SendWelcomeEmail(ctx context.Context, userID uuid.UUID) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	data := map[string]interface{}{
		"user_id":    user.ID.String(),
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	}

	return uc.emailService.SendTemplateEmail(ctx, "welcome", user.Email, user.GetFullName(), data)
}

// SendOrderConfirmationEmail sends order confirmation email
func (uc *emailUseCase) SendOrderConfirmationEmail(ctx context.Context, orderID uuid.UUID) error {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	user, err := uc.userRepo.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	data := map[string]interface{}{
		"user_id":      user.ID.String(),
		"order_id":     order.ID.String(),
		"order_number": order.OrderNumber,
		"first_name":   user.FirstName,
		"total":        order.Total,
		"items_count":  len(order.Items),
	}

	return uc.emailService.SendTemplateEmail(ctx, "order_confirmation", user.Email, user.GetFullName(), data)
}

// SendOrderShippedEmail sends order shipped email
func (uc *emailUseCase) SendOrderShippedEmail(ctx context.Context, orderID uuid.UUID) error {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	user, err := uc.userRepo.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	data := map[string]interface{}{
		"user_id":         user.ID.String(),
		"order_id":        order.ID.String(),
		"order_number":    order.OrderNumber,
		"first_name":      user.FirstName,
		"tracking_number": order.TrackingNumber,
	}

	return uc.emailService.SendTemplateEmail(ctx, "order_shipped", user.Email, user.GetFullName(), data)
}

// SendOrderDeliveredEmail sends order delivered email
func (uc *emailUseCase) SendOrderDeliveredEmail(ctx context.Context, orderID uuid.UUID) error {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	user, err := uc.userRepo.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	data := map[string]interface{}{
		"user_id":      user.ID.String(),
		"order_id":     order.ID.String(),
		"order_number": order.OrderNumber,
		"first_name":   user.FirstName,
	}

	return uc.emailService.SendTemplateEmail(ctx, "order_delivered", user.Email, user.GetFullName(), data)
}

// SendOrderCancelledEmail sends order cancelled email
func (uc *emailUseCase) SendOrderCancelledEmail(ctx context.Context, orderID uuid.UUID) error {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	user, err := uc.userRepo.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	data := map[string]interface{}{
		"user_id":      user.ID.String(),
		"order_id":     order.ID.String(),
		"order_number": order.OrderNumber,
		"first_name":   user.FirstName,
		"total":        order.Total,
	}

	return uc.emailService.SendTemplateEmail(ctx, "order_cancelled", user.Email, user.GetFullName(), data)
}

// SendPasswordResetEmail sends password reset email
func (uc *emailUseCase) SendPasswordResetEmail(ctx context.Context, userID uuid.UUID, resetToken string) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	data := map[string]interface{}{
		"user_id":     user.ID.String(),
		"first_name":  user.FirstName,
		"reset_token": resetToken,
		"reset_url":   fmt.Sprintf("https://yoursite.com/reset-password?token=%s", resetToken),
	}

	return uc.emailService.SendTemplateEmail(ctx, "password_reset", user.Email, user.GetFullName(), data)
}

// SendAbandonedCartEmail sends abandoned cart email
func (uc *emailUseCase) SendAbandonedCartEmail(ctx context.Context, userID uuid.UUID) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	data := map[string]interface{}{
		"user_id":    user.ID.String(),
		"first_name": user.FirstName,
		"cart_url":   "https://yoursite.com/cart",
	}

	return uc.emailService.SendTemplateEmail(ctx, "abandoned_cart", user.Email, user.GetFullName(), data)
}

// SendReviewRequestEmail sends review request email
func (uc *emailUseCase) SendReviewRequestEmail(ctx context.Context, userID, orderID uuid.UUID) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	data := map[string]interface{}{
		"user_id":      user.ID.String(),
		"order_id":     order.ID.String(),
		"order_number": order.OrderNumber,
		"first_name":   user.FirstName,
		"review_url":   fmt.Sprintf("https://yoursite.com/orders/%s/review", order.ID),
	}

	return uc.emailService.SendTemplateEmail(ctx, "review_request", user.Email, user.GetFullName(), data)
}

// SendLowStockAlert sends low stock alert email to admins
func (uc *emailUseCase) SendLowStockAlert(ctx context.Context, productID uuid.UUID) error {
	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Get admin users (you might want to have a specific admin email list)
	data := map[string]interface{}{
		"product_id":    product.ID.String(),
		"product_name":  product.Name,
		"current_stock": product.Stock,
		"threshold":     product.LowStockThreshold,
	}

	// Send to admin email (you should configure this)
	adminEmail := "admin@yoursite.com"
	return uc.emailService.SendTemplateEmail(ctx, "low_stock_alert", adminEmail, "Admin", data)
}

// Request/Response types
type CreateTemplateRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Type        entities.EmailType     `json:"type" validate:"required"`
	Subject     string                 `json:"subject" validate:"required"`
	BodyText    string                 `json:"body_text"`
	BodyHTML    string                 `json:"body_html"`
	Description string                 `json:"description"`
	Variables   map[string]interface{} `json:"variables"`
}

type UpdateTemplateRequest struct {
	Subject     *string                `json:"subject"`
	BodyText    *string                `json:"body_text"`
	BodyHTML    *string                `json:"body_html"`
	Description *string                `json:"description"`
	Variables   map[string]interface{} `json:"variables"`
	IsActive    *bool                  `json:"is_active"`
}

type TemplateResponse struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Type        entities.EmailType     `json:"type"`
	Subject     string                 `json:"subject"`
	BodyText    string                 `json:"body_text"`
	BodyHTML    string                 `json:"body_html"`
	IsActive    bool                   `json:"is_active"`
	Version     int                    `json:"version"`
	Description string                 `json:"description"`
	Variables   map[string]interface{} `json:"variables"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type UpdateSubscriptionsRequest struct {
	Newsletter     *bool `json:"newsletter"`
	Promotions     *bool `json:"promotions"`
	OrderUpdates   *bool `json:"order_updates"`
	ReviewRequests *bool `json:"review_requests"`
	AbandonedCart  *bool `json:"abandoned_cart"`
	Support        *bool `json:"support"`
}

type SubscriptionsResponse struct {
	UserID         uuid.UUID `json:"user_id"`
	Newsletter     bool      `json:"newsletter"`
	Promotions     bool      `json:"promotions"`
	OrderUpdates   bool      `json:"order_updates"`
	ReviewRequests bool      `json:"review_requests"`
	AbandonedCart  bool      `json:"abandoned_cart"`
	Support        bool      `json:"support"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type EmailStatsResponse struct {
	TotalSent      int64                                         `json:"total_sent"`
	TotalDelivered int64                                         `json:"total_delivered"`
	TotalOpened    int64                                         `json:"total_opened"`
	TotalClicked   int64                                         `json:"total_clicked"`
	TotalBounced   int64                                         `json:"total_bounced"`
	TotalFailed    int64                                         `json:"total_failed"`
	DeliveryRate   float64                                       `json:"delivery_rate"`
	OpenRate       float64                                       `json:"open_rate"`
	ClickRate      float64                                       `json:"click_rate"`
	BounceRate     float64                                       `json:"bounce_rate"`
	FailureRate    float64                                       `json:"failure_rate"`
	TypeStats      map[entities.EmailType]repositories.TypeStats `json:"type_stats"`
	Since          time.Time                                     `json:"since"`
	Until          time.Time                                     `json:"until"`
}

type EmailResponse struct {
	ID           uuid.UUID              `json:"id"`
	Type         entities.EmailType     `json:"type"`
	Priority     entities.EmailPriority `json:"priority"`
	Status       entities.EmailStatus   `json:"status"`
	ToEmail      string                 `json:"to_email"`
	ToName       string                 `json:"to_name"`
	Subject      string                 `json:"subject"`
	SentAt       *time.Time             `json:"sent_at"`
	DeliveredAt  *time.Time             `json:"delivered_at"`
	OpenedAt     *time.Time             `json:"opened_at"`
	ClickedAt    *time.Time             `json:"clicked_at"`
	RetryCount   int                    `json:"retry_count"`
	ErrorMessage string                 `json:"error_message"`
	CreatedAt    time.Time              `json:"created_at"`
}

// CreateTemplate creates an email template
func (uc *emailUseCase) CreateTemplate(ctx context.Context, req CreateTemplateRequest) (*TemplateResponse, error) {
	template := &entities.EmailTemplate{
		ID:          uuid.New(),
		Name:        req.Name,
		Type:        req.Type,
		Subject:     req.Subject,
		BodyText:    req.BodyText,
		BodyHTML:    req.BodyHTML,
		IsActive:    true,
		Version:     1,
		Description: req.Description,
		Variables:   req.Variables,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.templateRepo.Create(ctx, template); err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	return uc.toTemplateResponse(template), nil
}

// UpdateTemplate updates an email template
func (uc *emailUseCase) UpdateTemplate(ctx context.Context, id uuid.UUID, req UpdateTemplateRequest) (*TemplateResponse, error) {
	template, err := uc.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// Update fields if provided
	if req.Subject != nil {
		template.Subject = *req.Subject
	}
	if req.BodyText != nil {
		template.BodyText = *req.BodyText
	}
	if req.BodyHTML != nil {
		template.BodyHTML = *req.BodyHTML
	}
	if req.Description != nil {
		template.Description = *req.Description
	}
	if req.Variables != nil {
		template.Variables = req.Variables
	}
	if req.IsActive != nil {
		template.IsActive = *req.IsActive
	}

	template.UpdatedAt = time.Now()

	if err := uc.templateRepo.Update(ctx, template); err != nil {
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	return uc.toTemplateResponse(template), nil
}

// GetTemplate gets an email template by ID
func (uc *emailUseCase) GetTemplate(ctx context.Context, id uuid.UUID) (*TemplateResponse, error) {
	template, err := uc.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return uc.toTemplateResponse(template), nil
}

// ListTemplates lists email templates
func (uc *emailUseCase) ListTemplates(ctx context.Context, offset, limit int) ([]*TemplateResponse, error) {
	templates, err := uc.templateRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	responses := make([]*TemplateResponse, len(templates))
	for i, template := range templates {
		responses[i] = uc.toTemplateResponse(template)
	}

	return responses, nil
}

// DeleteTemplate deletes an email template
func (uc *emailUseCase) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	return uc.templateRepo.Delete(ctx, id)
}

// UpdateSubscriptions updates user email subscriptions
func (uc *emailUseCase) UpdateSubscriptions(ctx context.Context, userID uuid.UUID, req UpdateSubscriptionsRequest) error {
	subscription, err := uc.subscriptionRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Create new subscription if not exists
		subscription = &entities.EmailSubscription{
			ID:             uuid.New(),
			UserID:         userID,
			Newsletter:     true,
			Promotions:     true,
			OrderUpdates:   true,
			ReviewRequests: true,
			AbandonedCart:  true,
			Support:        true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
	}

	// Update fields if provided
	if req.Newsletter != nil {
		subscription.Newsletter = *req.Newsletter
	}
	if req.Promotions != nil {
		subscription.Promotions = *req.Promotions
	}
	if req.OrderUpdates != nil {
		subscription.OrderUpdates = *req.OrderUpdates
	}
	if req.ReviewRequests != nil {
		subscription.ReviewRequests = *req.ReviewRequests
	}
	if req.AbandonedCart != nil {
		subscription.AbandonedCart = *req.AbandonedCart
	}
	if req.Support != nil {
		subscription.Support = *req.Support
	}

	subscription.UpdatedAt = time.Now()

	if subscription.CreatedAt.IsZero() {
		return uc.subscriptionRepo.Create(ctx, subscription)
	}
	return uc.subscriptionRepo.Update(ctx, subscription)
}

// GetSubscriptions gets user email subscriptions
func (uc *emailUseCase) GetSubscriptions(ctx context.Context, userID uuid.UUID) (*SubscriptionsResponse, error) {
	subscription, err := uc.subscriptionRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Return default subscriptions if not found
		return &SubscriptionsResponse{
			UserID:         userID,
			Newsletter:     true,
			Promotions:     true,
			OrderUpdates:   true,
			ReviewRequests: true,
			AbandonedCart:  true,
			Support:        true,
			UpdatedAt:      time.Now(),
		}, nil
	}

	return &SubscriptionsResponse{
		UserID:         subscription.UserID,
		Newsletter:     subscription.Newsletter,
		Promotions:     subscription.Promotions,
		OrderUpdates:   subscription.OrderUpdates,
		ReviewRequests: subscription.ReviewRequests,
		AbandonedCart:  subscription.AbandonedCart,
		Support:        subscription.Support,
		UpdatedAt:      subscription.UpdatedAt,
	}, nil
}

// GetEmailStats gets email statistics
func (uc *emailUseCase) GetEmailStats(ctx context.Context, since time.Time) (*EmailStatsResponse, error) {
	stats, err := uc.emailRepo.GetEmailStats(ctx, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get email stats: %w", err)
	}

	return &EmailStatsResponse{
		TotalSent:      stats.TotalSent,
		TotalDelivered: stats.TotalDelivered,
		TotalOpened:    stats.TotalOpened,
		TotalClicked:   stats.TotalClicked,
		TotalBounced:   stats.TotalBounced,
		TotalFailed:    stats.TotalFailed,
		DeliveryRate:   stats.DeliveryRate,
		OpenRate:       stats.OpenRate,
		ClickRate:      stats.ClickRate,
		BounceRate:     stats.BounceRate,
		FailureRate:    stats.FailureRate,
		TypeStats:      stats.TypeStats,
		Since:          since,
		Until:          time.Now(),
	}, nil
}

// GetEmailHistory gets email history for a user
func (uc *emailUseCase) GetEmailHistory(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*EmailResponse, error) {
	emails, err := uc.emailRepo.GetByUserID(ctx, userID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get email history: %w", err)
	}

	responses := make([]*EmailResponse, len(emails))
	for i, email := range emails {
		responses[i] = &EmailResponse{
			ID:           email.ID,
			Type:         email.Type,
			Priority:     email.Priority,
			Status:       email.Status,
			ToEmail:      email.ToEmail,
			ToName:       email.ToName,
			Subject:      email.Subject,
			SentAt:       email.SentAt,
			DeliveredAt:  email.DeliveredAt,
			OpenedAt:     email.OpenedAt,
			ClickedAt:    email.ClickedAt,
			RetryCount:   email.RetryCount,
			ErrorMessage: email.ErrorMessage,
			CreatedAt:    email.CreatedAt,
		}
	}

	return responses, nil
}

// RetryFailedEmails retries failed emails
func (uc *emailUseCase) RetryFailedEmails(ctx context.Context) error {
	return uc.emailService.RetryFailedEmails(ctx)
}

// GetFailedEmails gets failed emails
func (uc *emailUseCase) GetFailedEmails(ctx context.Context, since time.Time) ([]*EmailResponse, error) {
	emails, err := uc.emailRepo.GetFailedEmails(ctx, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get failed emails: %w", err)
	}

	responses := make([]*EmailResponse, len(emails))
	for i, email := range emails {
		responses[i] = &EmailResponse{
			ID:           email.ID,
			Type:         email.Type,
			Priority:     email.Priority,
			Status:       email.Status,
			ToEmail:      email.ToEmail,
			ToName:       email.ToName,
			Subject:      email.Subject,
			SentAt:       email.SentAt,
			DeliveredAt:  email.DeliveredAt,
			OpenedAt:     email.OpenedAt,
			ClickedAt:    email.ClickedAt,
			RetryCount:   email.RetryCount,
			ErrorMessage: email.ErrorMessage,
			CreatedAt:    email.CreatedAt,
		}
	}

	return responses, nil
}

// Helper function to convert template entity to response
func (uc *emailUseCase) toTemplateResponse(template *entities.EmailTemplate) *TemplateResponse {
	return &TemplateResponse{
		ID:          template.ID,
		Name:        template.Name,
		Type:        template.Type,
		Subject:     template.Subject,
		BodyText:    template.BodyText,
		BodyHTML:    template.BodyHTML,
		IsActive:    template.IsActive,
		Version:     template.Version,
		Description: template.Description,
		Variables:   template.Variables,
		CreatedAt:   template.CreatedAt,
		UpdatedAt:   template.UpdatedAt,
	}
}
