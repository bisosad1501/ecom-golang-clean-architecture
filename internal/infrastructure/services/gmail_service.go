package services

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/infrastructure/config"
)

// GmailService handles email sending via Gmail SMTP
type GmailService struct {
	config *config.EmailConfig
	auth   smtp.Auth
}

// NewGmailService creates a new Gmail service
func NewGmailService(config *config.EmailConfig) *GmailService {
	var auth smtp.Auth
	if config.SMTPUsername != "" && config.SMTPPassword != "" {
		auth = smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)
	}

	return &GmailService{
		config: config,
		auth:   auth,
	}
}

// SendEmail sends an email via Gmail SMTP
func (g *GmailService) SendEmail(ctx context.Context, to, subject, body string) error {
	return g.SendEmailWithTemplate(ctx, to, subject, body, "")
}

// SendEmailWithTemplate sends an email with HTML template
func (g *GmailService) SendEmailWithTemplate(ctx context.Context, to, subject, bodyText, bodyHTML string) error {
	// Build email message
	message, err := g.buildEmailMessage(to, subject, bodyText, bodyHTML)
	if err != nil {
		return fmt.Errorf("failed to build email message: %w", err)
	}

	// Send email
	if err := g.sendSMTP(to, message); err != nil {
		return fmt.Errorf("failed to send email via Gmail SMTP: %w", err)
	}

	return nil
}

// SendVerificationEmail sends email verification
func (g *GmailService) SendVerificationEmail(ctx context.Context, to, firstName, verificationLink string) error {
	subject := "Verify Your Email Address"
	
	bodyText := fmt.Sprintf(`Hi %s,

Thank you for signing up! Please verify your email address by clicking the link below:

%s

This verification link will expire in 24 hours.

If you didn't create an account, please ignore this email.

Best regards,
%s`, firstName, verificationLink, g.config.FromName)

	bodyHTML := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email - BiHub Store</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #ffffff;
            margin: 0;
            padding: 20px;
            background: linear-gradient(135deg, #0f0f0f 0%%, #1a1a1a 50%%, #2d2d2d 100%%);
            min-height: 100vh;
        }
        .email-wrapper {
            max-width: 600px;
            margin: 0 auto;
            background: linear-gradient(145deg, #1a1a1a 0%%, #2a2a2a 100%%);
            border-radius: 20px;
            overflow: hidden;
            box-shadow:
                0 20px 40px rgba(0,0,0,0.4),
                0 0 0 1px rgba(255,144,0,0.1),
                inset 0 1px 0 rgba(255,255,255,0.1);
            backdrop-filter: blur(10px);
        }
        .header {
            background: linear-gradient(135deg, #FF9000 0%%, #FF7A00 50%%, #e67e00 100%%);
            color: white;
            padding: 50px 30px;
            text-align: center;
            position: relative;
            overflow: hidden;
        }
        .header::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background:
                radial-gradient(circle at 20%% 80%%, rgba(255,255,255,0.1) 0%%, transparent 50%%),
                radial-gradient(circle at 80%% 20%%, rgba(255,255,255,0.1) 0%%, transparent 50%%);
        }
        .logo-container {
            display: flex;
            align-items: center;
            justify-content: center;
            margin-bottom: 25px;
            position: relative;
            z-index: 1;
        }
        .logo-text {
            font-size: 36px;
            font-weight: 800;
            display: flex;
            align-items: center;
            text-shadow: 0 2px 8px rgba(0,0,0,0.3);
        }
        .logo-bi {
            color: white;
        }
        .logo-hub {
            background: rgba(0,0,0,0.2);
            color: white;
            padding: 2px 8px;
            border-radius: 6px;
            margin-left: 2px;
            font-weight: 900;
            letter-spacing: 0.5px;
        }
        .header h1 {
            margin: 0;
            font-size: 32px;
            font-weight: 700;
            position: relative;
            z-index: 1;
            text-shadow: 0 2px 8px rgba(0,0,0,0.3);
            letter-spacing: -0.5px;
        }
        .content {
            padding: 50px 40px;
            background: linear-gradient(145deg, #2a2a2a 0%%, #1f1f1f 100%%);
            color: #e8e8e8;
            position: relative;
        }
        .content::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            height: 1px;
            background: linear-gradient(90deg, transparent, rgba(255,144,0,0.3), transparent);
        }
        .content p {
            margin: 0 0 24px 0;
            font-size: 17px;
            line-height: 1.7;
            color: #e8e8e8;
        }
        .content strong {
            color: #FF9000;
        }
        .button-container {
            text-align: center;
            margin: 40px 0;
        }
        .button {
            display: inline-block;
            padding: 18px 40px;
            background: linear-gradient(135deg, #FF9000 0%%, #FF7A00 50%%, #e67e00 100%%);
            color: white;
            text-decoration: none;
            border-radius: 12px;
            font-weight: 700;
            font-size: 16px;
            text-transform: uppercase;
            letter-spacing: 1px;
            box-shadow:
                0 8px 25px rgba(255, 144, 0, 0.4),
                0 0 0 1px rgba(255, 144, 0, 0.2),
                inset 0 1px 0 rgba(255, 255, 255, 0.2);
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
            position: relative;
            overflow: hidden;
        }
        .button::before {
            content: '';
            position: absolute;
            top: 0;
            left: -100%%;
            width: 100%%;
            height: 100%%;
            background: linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent);
            transition: left 0.5s;
        }
        .button:hover::before {
            left: 100%%;
        }
        .features-list {
            background: rgba(255, 144, 0, 0.05);
            border: 1px solid rgba(255, 144, 0, 0.2);
            border-radius: 12px;
            padding: 25px;
            margin: 30px 0;
        }
        .features-list ul {
            margin: 0;
            padding-left: 20px;
            color: #e8e8e8;
        }
        .features-list li {
            margin: 8px 0;
            font-size: 16px;
            line-height: 1.6;
        }
        .features-list li::marker {
            color: #FF9000;
        }
        .warning {
            background: linear-gradient(135deg, rgba(255, 144, 0, 0.1) 0%%, rgba(255, 144, 0, 0.05) 100%%);
            border: 1px solid rgba(255, 144, 0, 0.3);
            border-radius: 12px;
            padding: 20px;
            margin: 30px 0;
            font-size: 15px;
            color: #ffcc80;
            position: relative;
        }
        .warning::before {
            content: '⚡';
            position: absolute;
            top: 20px;
            left: 20px;
            font-size: 20px;
        }
        .warning-content {
            margin-left: 35px;
        }
        .footer {
            padding: 40px 30px;
            text-align: center;
            color: #999;
            font-size: 14px;
            background: linear-gradient(145deg, #1a1a1a 0%%, #0f0f0f 100%%);
            border-top: 1px solid rgba(255, 144, 0, 0.1);
        }
        .footer p {
            margin: 0 0 10px 0;
            line-height: 1.5;
        }
        .brand-name {
            color: #FF9000;
            font-weight: 700;
        }
    </style>
</head>
<body>
    <div class="email-wrapper">
        <div class="header">
            <div class="logo-container">
                <div class="logo-text">
                    <span class="logo-bi">Bi</span><span class="logo-hub">hub</span>
                </div>
            </div>
            <h1>🎉 Welcome to BiHub Store!</h1>
        </div>
        <div class="content">
            <p>Hi <strong>%s</strong>,</p>
            <p>Welcome to <span class="brand-name">BiHub Store</span>! 🛍️ Thank you for joining our amazing community of shoppers.</p>
            <p>To unlock your account and start your shopping journey, please verify your email address by clicking the button below:</p>

            <div class="button-container">
                <a href="%s" class="button">✨ Verify My Email ✨</a>
            </div>

            <div class="warning">
                <div class="warning-content">
                    <strong>⏰ Important:</strong> This verification link will expire in 24 hours for security reasons.
                </div>
            </div>

            <p>If you didn't create an account with us, please ignore this email and no further action is required.</p>

            <div class="features-list">
                <p><strong>🎁 Once verified, you'll unlock:</strong></p>
                <ul>
                    <li>🛒 Browse our premium product collection</li>
                    <li>💳 Seamless checkout experience</li>
                    <li>📦 Real-time order tracking</li>
                    <li>⭐ Exclusive member benefits & deals</li>
                    <li>🎯 Personalized recommendations</li>
                </ul>
            </div>
        </div>
        <div class="footer">
            <p>Best regards,<br><span class="brand-name">The BiHub Store Team</span> 💙</p>
            <p style="margin-top: 15px; font-size: 12px; color: #666;">
                This email was sent to <strong>%s</strong><br>
                Need help? Contact our support team anytime!
            </p>
        </div>
    </div>
</body>
</html>`, firstName, verificationLink, to)

	return g.SendEmailWithTemplate(ctx, to, subject, bodyText, bodyHTML)
}

// SendPasswordResetEmail sends password reset email
func (g *GmailService) SendPasswordResetEmail(ctx context.Context, to, firstName, resetLink string) error {
	subject := "Reset Your Password"
	
	bodyText := fmt.Sprintf(`Hi %s,

You requested to reset your password. Click the link below to reset it:

%s

This reset link will expire in 1 hour.

If you didn't request a password reset, please ignore this email.

Best regards,
%s`, firstName, resetLink, g.config.FromName)

	bodyHTML := fmt.Sprintf(`<!DOCTYPE html>
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
            <p>Hi %s,</p>
            <p>You requested to reset your password. Click the button below to reset it:</p>
            <p style="text-align: center;">
                <a href="%s" class="button">Reset Password</a>
            </p>
            <p>This reset link will expire in 1 hour.</p>
            <p>If you didn't request a password reset, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>%s</p>
        </div>
    </div>
</body>
</html>`, firstName, resetLink, g.config.FromName)

	return g.SendEmailWithTemplate(ctx, to, subject, bodyText, bodyHTML)
}

// SendWelcomeEmail sends welcome email
func (g *GmailService) SendWelcomeEmail(ctx context.Context, to, firstName string) error {
	subject := "Welcome to " + g.config.FromName + "!"
	
	bodyText := fmt.Sprintf(`Hi %s,

Welcome to %s! We're excited to have you as a customer.

You can now:
- Browse our products
- Add items to your cart
- Track your orders
- Manage your account

If you have any questions, feel free to contact our support team.

Happy shopping!

Best regards,
%s`, firstName, g.config.FromName, g.config.FromName)

	bodyHTML := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to %s</title>
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
            <h1>Welcome to %s!</h1>
        </div>
        <div class="content">
            <p>Hi %s,</p>
            <p>Welcome to %s! We're excited to have you as a customer.</p>
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
            <p>Best regards,<br>%s</p>
        </div>
    </div>
</body>
</html>`, g.config.FromName, g.config.FromName, firstName, g.config.FromName, g.config.FromName)

	return g.SendEmailWithTemplate(ctx, to, subject, bodyText, bodyHTML)
}

// ValidateConfiguration validates Gmail SMTP configuration
func (g *GmailService) ValidateConfiguration() error {
	if g.config.SMTPHost == "" {
		return fmt.Errorf("SMTP host is required")
	}
	if g.config.SMTPPort == "" {
		return fmt.Errorf("SMTP port is required")
	}
	if g.config.SMTPUsername == "" {
		return fmt.Errorf("SMTP username is required")
	}
	if g.config.SMTPPassword == "" {
		return fmt.Errorf("SMTP password is required")
	}
	if g.config.FromEmail == "" {
		return fmt.Errorf("from email is required")
	}

	// Test connection
	return g.testConnection()
}

// buildEmailMessage builds the email message in RFC 5322 format
func (g *GmailService) buildEmailMessage(to, subject, bodyText, bodyHTML string) ([]byte, error) {
	var message strings.Builder

	// Headers
	message.WriteString(fmt.Sprintf("From: %s <%s>\r\n", g.config.FromName, g.config.FromEmail))
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	message.WriteString("MIME-Version: 1.0\r\n")

	if g.config.ReplyToEmail != "" {
		message.WriteString(fmt.Sprintf("Reply-To: %s\r\n", g.config.ReplyToEmail))
	}

	// Content type
	if bodyHTML != "" {
		message.WriteString("Content-Type: multipart/alternative; boundary=\"boundary123\"\r\n")
		message.WriteString("\r\n")

		// Text part
		if bodyText != "" {
			message.WriteString("--boundary123\r\n")
			message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
			message.WriteString("\r\n")
			message.WriteString(bodyText)
			message.WriteString("\r\n")
		}

		// HTML part
		message.WriteString("--boundary123\r\n")
		message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		message.WriteString("\r\n")
		message.WriteString(bodyHTML)
		message.WriteString("\r\n")
		message.WriteString("--boundary123--\r\n")
	} else {
		message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		message.WriteString("\r\n")
		message.WriteString(bodyText)
	}

	return []byte(message.String()), nil
}

// sendSMTP sends the email via SMTP with TLS
func (g *GmailService) sendSMTP(to string, message []byte) error {
	addr := fmt.Sprintf("%s:%s", g.config.SMTPHost, g.config.SMTPPort)
	log.Printf("🔄 Connecting to SMTP server: %s", addr)
	log.Printf("🔄 From: %s, To: %s", g.config.FromEmail, to)

	client, err := smtp.Dial(addr)
	if err != nil {
		log.Printf("❌ Failed to connect to SMTP server: %v", err)
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()
	log.Printf("✅ Connected to SMTP server")

	// Start TLS
	tlsConfig := &tls.Config{
		ServerName: g.config.SMTPHost,
	}
	if err := client.StartTLS(tlsConfig); err != nil {
		log.Printf("❌ Failed to start TLS: %v", err)
		return fmt.Errorf("failed to start TLS: %w", err)
	}
	log.Printf("✅ TLS connection established")

	// Authenticate
	if g.auth != nil {
		if err := client.Auth(g.auth); err != nil {
			log.Printf("❌ Failed to authenticate with username %s: %v", g.config.SMTPUsername, err)
			return fmt.Errorf("failed to authenticate: %w", err)
		}
		log.Printf("✅ SMTP authentication successful")
	}

	// Send email
	if err := client.Mail(g.config.FromEmail); err != nil {
		log.Printf("❌ Failed to set sender %s: %v", g.config.FromEmail, err)
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		log.Printf("❌ Failed to set recipient %s: %v", to, err)
		return fmt.Errorf("failed to set recipient: %w", err)
	}
	log.Printf("✅ Sender and recipient set successfully")

	writer, err := client.Data()
	if err != nil {
		log.Printf("❌ Failed to get data writer: %v", err)
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer writer.Close()

	if _, err := writer.Write(message); err != nil {
		log.Printf("❌ Failed to write message: %v", err)
		return fmt.Errorf("failed to write message: %w", err)
	}

	log.Printf("✅ Email message sent successfully to %s", to)
	return nil
}

// testConnection tests the Gmail SMTP connection
func (g *GmailService) testConnection() error {
	addr := fmt.Sprintf("%s:%s", g.config.SMTPHost, g.config.SMTPPort)
	
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	tlsConfig := &tls.Config{
		ServerName: g.config.SMTPHost,
	}
	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	if g.auth != nil {
		if err := client.Auth(g.auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	return nil
}
