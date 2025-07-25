package usecases

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

// SettingsValidator provides validation for settings
type SettingsValidator struct{}

// NewSettingsValidator creates a new settings validator
func NewSettingsValidator() *SettingsValidator {
	return &SettingsValidator{}
}

// ValidateGeneralSettings validates general settings
func (v *SettingsValidator) ValidateGeneralSettings(req UpdateGeneralSettingsRequest) error {
	// Validate store name
	if strings.TrimSpace(req.StoreName) == "" {
		return fmt.Errorf("store name cannot be empty")
	}
	if len(req.StoreName) > 255 {
		return fmt.Errorf("store name cannot exceed 255 characters")
	}

	// Validate email
	if req.StoreEmail != "" {
		if _, err := mail.ParseAddress(req.StoreEmail); err != nil {
			return fmt.Errorf("invalid email format: %s", req.StoreEmail)
		}
	}

	// Validate phone
	if req.StorePhone != "" {
		phoneRegex := regexp.MustCompile(`^[\+]?[1-9][\d]{0,15}$`)
		cleanPhone := strings.ReplaceAll(strings.ReplaceAll(req.StorePhone, "-", ""), " ", "")
		if !phoneRegex.MatchString(cleanPhone) {
			return fmt.Errorf("invalid phone format: %s", req.StorePhone)
		}
	}

	// Validate currency
	validCurrencies := []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "CNY", "SEK", "NZD"}
	if !contains(validCurrencies, req.Currency) {
		return fmt.Errorf("invalid currency: %s. Valid currencies: %s", req.Currency, strings.Join(validCurrencies, ", "))
	}

	// Validate timezone
	if req.Timezone != "" {
		if _, err := time.LoadLocation(req.Timezone); err != nil {
			return fmt.Errorf("invalid timezone: %s", req.Timezone)
		}
	}

	// Validate language
	validLanguages := []string{"en", "es", "fr", "de", "it", "pt", "ru", "ja", "ko", "zh"}
	if !contains(validLanguages, req.Language) {
		return fmt.Errorf("invalid language: %s. Valid languages: %s", req.Language, strings.Join(validLanguages, ", "))
	}

	return nil
}

// ValidatePaymentSettings validates payment settings
func (v *SettingsValidator) ValidatePaymentSettings(req UpdatePaymentSettingsRequest) error {
	// Validate default currency
	validCurrencies := []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "CNY", "SEK", "NZD"}
	if !contains(validCurrencies, req.DefaultCurrency) {
		return fmt.Errorf("invalid default currency: %s", req.DefaultCurrency)
	}

	// Validate accepted currencies
	for _, currency := range req.AcceptedCurrencies {
		if !contains(validCurrencies, currency) {
			return fmt.Errorf("invalid accepted currency: %s", currency)
		}
	}

	// Validate Stripe settings
	if req.StripeEnabled {
		if req.StripePublicKey == "" {
			return fmt.Errorf("stripe public key is required when Stripe is enabled")
		}
		if req.StripeSecretKey == "" {
			return fmt.Errorf("stripe secret key is required when Stripe is enabled")
		}
		if !strings.HasPrefix(req.StripePublicKey, "pk_") {
			return fmt.Errorf("invalid Stripe public key format")
		}
		if !strings.HasPrefix(req.StripeSecretKey, "sk_") {
			return fmt.Errorf("invalid Stripe secret key format")
		}
	}

	// Validate PayPal settings
	if req.PayPalEnabled {
		if req.PayPalClientID == "" {
			return fmt.Errorf("PayPal client ID is required when PayPal is enabled")
		}
		if req.PayPalClientSecret == "" {
			return fmt.Errorf("PayPal client secret is required when PayPal is enabled")
		}
	}

	// Validate bank details if bank transfer is enabled
	if req.BankTransfer {
		if bankName, ok := req.BankDetails["bank_name"].(string); !ok || bankName == "" {
			return fmt.Errorf("bank name is required when bank transfer is enabled")
		}
		if accountNumber, ok := req.BankDetails["account_number"].(string); !ok || accountNumber == "" {
			return fmt.Errorf("account number is required when bank transfer is enabled")
		}
	}

	return nil
}

// ValidateEmailSettings validates email settings
func (v *SettingsValidator) ValidateEmailSettings(req UpdateEmailSettingsRequest) error {
	// Validate SMTP host
	if req.SMTPHost == "" {
		return fmt.Errorf("SMTP host is required")
	}

	// Validate SMTP port
	if req.SMTPPort <= 0 || req.SMTPPort > 65535 {
		return fmt.Errorf("SMTP port must be between 1 and 65535")
	}

	// Validate from email
	if req.FromEmail != "" {
		if _, err := mail.ParseAddress(req.FromEmail); err != nil {
			return fmt.Errorf("invalid from email format: %s", req.FromEmail)
		}
	}

	// Validate reply-to email
	if req.ReplyToEmail != "" {
		if _, err := mail.ParseAddress(req.ReplyToEmail); err != nil {
			return fmt.Errorf("invalid reply-to email format: %s", req.ReplyToEmail)
		}
	}

	// Validate SMTP security
	validSecurity := []string{"none", "tls", "ssl"}
	if !contains(validSecurity, req.SMTPSecurity) {
		return fmt.Errorf("invalid SMTP security: %s. Valid options: %s", req.SMTPSecurity, strings.Join(validSecurity, ", "))
	}

	return nil
}

// ValidateTaxSettings validates tax settings
func (v *SettingsValidator) ValidateTaxSettings(req UpdateTaxSettingsRequest) error {
	// Validate default tax rate
	if req.DefaultTaxRate < 0 || req.DefaultTaxRate > 100 {
		return fmt.Errorf("default tax rate must be between 0 and 100")
	}

	// Validate digital tax rate
	if req.DigitalTaxRate < 0 || req.DigitalTaxRate > 100 {
		return fmt.Errorf("digital tax rate must be between 0 and 100")
	}

	// Validate tax calculation method
	validCalculations := []string{"exclusive", "inclusive"}
	if !contains(validCalculations, req.TaxCalculation) {
		return fmt.Errorf("invalid tax calculation: %s. Valid options: %s", req.TaxCalculation, strings.Join(validCalculations, ", "))
	}

	// Validate tax display format
	validFormats := []string{"including", "excluding", "both"}
	if !contains(validFormats, req.TaxDisplayFormat) {
		return fmt.Errorf("invalid tax display format: %s. Valid options: %s", req.TaxDisplayFormat, strings.Join(validFormats, ", "))
	}

	// Validate tax rounding
	validRounding := []string{"round_up", "round_down", "round_nearest"}
	if !contains(validRounding, req.TaxRounding) {
		return fmt.Errorf("invalid tax rounding: %s. Valid options: %s", req.TaxRounding, strings.Join(validRounding, ", "))
	}

	// Validate tax rates
	for _, rate := range req.TaxRates {
		if rate.Rate < 0 || rate.Rate > 100 {
			return fmt.Errorf("tax rate %s: rate must be between 0 and 100", rate.Name)
		}
		if strings.TrimSpace(rate.Name) == "" {
			return fmt.Errorf("tax rate name cannot be empty")
		}
		if rate.Country == "" {
			return fmt.Errorf("tax rate %s: country is required", rate.Name)
		}
	}

	return nil
}

// ValidateSecuritySettings validates security settings
func (v *SettingsValidator) ValidateSecuritySettings(req UpdateSecuritySettingsRequest) error {
	// Validate session timeout
	if req.SessionTimeout < 5 || req.SessionTimeout > 1440 { // 5 minutes to 24 hours
		return fmt.Errorf("session timeout must be between 5 and 1440 minutes")
	}

	// Validate password policy
	validPolicies := []string{"weak", "medium", "strong", "very_strong"}
	if !contains(validPolicies, req.PasswordPolicy) {
		return fmt.Errorf("invalid password policy: %s. Valid options: %s", req.PasswordPolicy, strings.Join(validPolicies, ", "))
	}

	// Validate max login attempts
	if req.MaxLoginAttempts < 1 || req.MaxLoginAttempts > 20 {
		return fmt.Errorf("max login attempts must be between 1 and 20")
	}

	// Validate lockout duration
	if req.LockoutDuration < 1 || req.LockoutDuration > 1440 { // 1 minute to 24 hours
		return fmt.Errorf("lockout duration must be between 1 and 1440 minutes")
	}

	// Validate data retention days
	if req.DataRetentionDays < 30 || req.DataRetentionDays > 2555 { // 30 days to 7 years
		return fmt.Errorf("data retention days must be between 30 and 2555 days")
	}

	// Validate captcha provider
	if req.EnableCaptcha {
		validProviders := []string{"recaptcha", "hcaptcha", "turnstile"}
		if !contains(validProviders, req.CaptchaProvider) {
			return fmt.Errorf("invalid captcha provider: %s. Valid options: %s", req.CaptchaProvider, strings.Join(validProviders, ", "))
		}
	}

	return nil
}

// ValidateSEOSettings validates SEO settings
func (v *SettingsValidator) ValidateSEOSettings(req UpdateSEOSettingsRequest) error {
	// Validate site title
	if strings.TrimSpace(req.SiteTitle) == "" {
		return fmt.Errorf("site title cannot be empty")
	}
	if len(req.SiteTitle) > 60 {
		return fmt.Errorf("site title should not exceed 60 characters for optimal SEO")
	}

	// Validate site description
	if len(req.SiteDescription) > 160 {
		return fmt.Errorf("site description should not exceed 160 characters for optimal SEO")
	}

	// Validate meta robots
	validRobots := []string{"index,follow", "index,nofollow", "noindex,follow", "noindex,nofollow"}
	if !contains(validRobots, req.MetaRobots) {
		return fmt.Errorf("invalid meta robots: %s. Valid options: %s", req.MetaRobots, strings.Join(validRobots, ", "))
	}

	// Validate Twitter card
	validCards := []string{"summary", "summary_large_image", "app", "player"}
	if req.TwitterCard != "" && !contains(validCards, req.TwitterCard) {
		return fmt.Errorf("invalid Twitter card: %s. Valid options: %s", req.TwitterCard, strings.Join(validCards, ", "))
	}

	return nil
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
