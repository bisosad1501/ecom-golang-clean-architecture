package services

import (
	"context"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
)

// EmailTemplateService handles email template operations
type EmailTemplateService struct {
	templateRepo repositories.EmailTemplateRepository
}

// NewEmailTemplateService creates a new email template service
func NewEmailTemplateService(templateRepo repositories.EmailTemplateRepository) *EmailTemplateService {
	return &EmailTemplateService{
		templateRepo: templateRepo,
	}
}

// InitializeDefaultTemplates creates default email templates if they don't exist
func (s *EmailTemplateService) InitializeDefaultTemplates(ctx context.Context) error {
	templates := s.getDefaultTemplates()

	for _, template := range templates {
		// Check if template already exists
		existing, err := s.templateRepo.GetByName(ctx, template.Name)
		if err == nil && existing != nil {
			continue // Template already exists
		}

		// Create new template
		if err := s.templateRepo.Create(ctx, template); err != nil {
			return fmt.Errorf("failed to create template %s: %w", template.Name, err)
		}
	}

	return nil
}

// getDefaultTemplates returns default email templates
func (s *EmailTemplateService) getDefaultTemplates() []*entities.EmailTemplate {
	return []*entities.EmailTemplate{
		{
			ID:          uuid.New(),
			Name:        "email_verification",
			Type:        entities.EmailTypeWelcome,
			Subject:     "Verify Your Email Address",
			BodyText:    s.getEmailVerificationTextTemplate(),
			BodyHTML:    s.getEmailVerificationHTMLTemplate(),
			IsActive:    true,
			Version:     1,
			Description: "Email verification template for new user registration",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "password_reset",
			Type:        entities.EmailTypePasswordReset,
			Subject:     "Reset Your Password",
			BodyText:    s.getPasswordResetTextTemplate(),
			BodyHTML:    s.getPasswordResetHTMLTemplate(),
			IsActive:    true,
			Version:     1,
			Description: "Password reset template",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "welcome",
			Type:        entities.EmailTypeWelcome,
			Subject:     "Welcome to Our Store!",
			BodyText:    s.getWelcomeTextTemplate(),
			BodyHTML:    s.getWelcomeHTMLTemplate(),
			IsActive:    true,
			Version:     1,
			Description: "Welcome email for new users",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "order_confirmation",
			Type:        entities.EmailTypeOrderConfirmation,
			Subject:     "Order Confirmation - Order #{{.order_number}}",
			BodyText:    s.getOrderConfirmationTextTemplate(),
			BodyHTML:    s.getOrderConfirmationHTMLTemplate(),
			IsActive:    true,
			Version:     1,
			Description: "Order confirmation email",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "order_shipped",
			Type:        entities.EmailTypeOrderShipped,
			Subject:     "Your Order Has Been Shipped - Order #{{.order_number}}",
			BodyText:    s.getOrderShippedTextTemplate(),
			BodyHTML:    s.getOrderShippedHTMLTemplate(),
			IsActive:    true,
			Version:     1,
			Description: "Order shipped notification email",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

// Template content methods
func (s *EmailTemplateService) getEmailVerificationTextTemplate() string {
	return `Hi {{.first_name}},

Thank you for signing up! Please verify your email address by clicking the link below:

{{.verification_link}}

This verification link will expire at {{.expires_at}}.

If you didn't create an account, please ignore this email.

Best regards,
The E-commerce Team`
}

func (s *EmailTemplateService) getEmailVerificationHTMLTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verify Your Email</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #007bff; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .button { display: inline-block; padding: 12px 24px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; }
        .footer { padding: 20px; text-align: center; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Verify Your Email Address</h1>
        </div>
        <div class="content">
            <p>Hi {{.first_name}},</p>
            <p>Thank you for signing up! Please verify your email address by clicking the button below:</p>
            <p style="text-align: center;">
                <a href="{{.verification_link}}" class="button">Verify Email Address</a>
            </p>
            <p>This verification link will expire at {{.expires_at}}.</p>
            <p>If you didn't create an account, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>The E-commerce Team</p>
        </div>
    </div>
</body>
</html>`
}

func (s *EmailTemplateService) getPasswordResetTextTemplate() string {
	return `Hi {{.first_name}},

You requested to reset your password. Click the link below to reset it:

{{.reset_link}}

This reset link will expire at {{.expires_at}}.

If you didn't request a password reset, please ignore this email.

Best regards,
The E-commerce Team`
}

func (s *EmailTemplateService) getPasswordResetHTMLTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Reset Your Password</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #dc3545; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .button { display: inline-block; padding: 12px 24px; background: #dc3545; color: white; text-decoration: none; border-radius: 4px; }
        .footer { padding: 20px; text-align: center; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Reset Your Password</h1>
        </div>
        <div class="content">
            <p>Hi {{.first_name}},</p>
            <p>You requested to reset your password. Click the button below to reset it:</p>
            <p style="text-align: center;">
                <a href="{{.reset_link}}" class="button">Reset Password</a>
            </p>
            <p>This reset link will expire at {{.expires_at}}.</p>
            <p>If you didn't request a password reset, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>The E-commerce Team</p>
        </div>
    </div>
</body>
</html>`
}

func (s *EmailTemplateService) getWelcomeTextTemplate() string {
	return `Hi {{.first_name}},

Welcome to our store! We're excited to have you as a customer.

You can now:
- Browse our products
- Add items to your cart
- Track your orders
- Manage your account

If you have any questions, feel free to contact our support team.

Happy shopping!

Best regards,
The E-commerce Team`
}

func (s *EmailTemplateService) getWelcomeHTMLTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to Our Store</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #28a745; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .features { list-style: none; padding: 0; }
        .features li { padding: 8px 0; }
        .footer { padding: 20px; text-align: center; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Our Store!</h1>
        </div>
        <div class="content">
            <p>Hi {{.first_name}},</p>
            <p>Welcome to our store! We're excited to have you as a customer.</p>
            <p>You can now:</p>
            <ul class="features">
                <li>✓ Browse our products</li>
                <li>✓ Add items to your cart</li>
                <li>✓ Track your orders</li>
                <li>✓ Manage your account</li>
            </ul>
            <p>If you have any questions, feel free to contact our support team.</p>
            <p>Happy shopping!</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>The E-commerce Team</p>
        </div>
    </div>
</body>
</html>`
}

func (s *EmailTemplateService) getOrderConfirmationTextTemplate() string {
	return `Hi {{.first_name}},

Thank you for your order! We've received your order and are processing it.

Order Details:
- Order Number: {{.order_number}}
- Total Items: {{.items_count}}
- Total Amount: ${{.total}}

We'll send you another email when your order ships.

Best regards,
The E-commerce Team`
}

func (s *EmailTemplateService) getOrderConfirmationHTMLTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Order Confirmation</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #007bff; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .order-details { background: white; padding: 15px; border-radius: 4px; margin: 15px 0; }
        .footer { padding: 20px; text-align: center; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Order Confirmation</h1>
        </div>
        <div class="content">
            <p>Hi {{.first_name}},</p>
            <p>Thank you for your order! We've received your order and are processing it.</p>
            <div class="order-details">
                <h3>Order Details:</h3>
                <p><strong>Order Number:</strong> {{.order_number}}</p>
                <p><strong>Total Items:</strong> {{.items_count}}</p>
                <p><strong>Total Amount:</strong> ${{.total}}</p>
            </div>
            <p>We'll send you another email when your order ships.</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>The E-commerce Team</p>
        </div>
    </div>
</body>
</html>`
}

func (s *EmailTemplateService) getOrderShippedTextTemplate() string {
	return `Hi {{.first_name}},

Great news! Your order has been shipped.

Order Number: {{.order_number}}

Your package is on its way and should arrive soon.

Best regards,
The E-commerce Team`
}

func (s *EmailTemplateService) getOrderShippedHTMLTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Order Shipped</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #28a745; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .footer { padding: 20px; text-align: center; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Order Shipped!</h1>
        </div>
        <div class="content">
            <p>Hi {{.first_name}},</p>
            <p>Great news! Your order has been shipped.</p>
            <p><strong>Order Number:</strong> {{.order_number}}</p>
            <p>Your package is on its way and should arrive soon.</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>The E-commerce Team</p>
        </div>
    </div>
</body>
</html>`
}
