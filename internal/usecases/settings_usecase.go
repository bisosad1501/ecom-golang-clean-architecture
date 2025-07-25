package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// SettingsUseCase defines settings management use cases
type SettingsUseCase interface {
	// General Settings
	GetGeneralSettings(ctx context.Context) (*GeneralSettingsResponse, error)
	UpdateGeneralSettings(ctx context.Context, req UpdateGeneralSettingsRequest) error

	// Store Configuration
	GetStoreConfig(ctx context.Context) (*StoreConfigResponse, error)
	UpdateStoreConfig(ctx context.Context, req UpdateStoreConfigRequest) error

	// Payment Settings
	GetPaymentSettings(ctx context.Context) (*PaymentSettingsResponse, error)
	UpdatePaymentSettings(ctx context.Context, req UpdatePaymentSettingsRequest) error

	// Email Settings
	GetEmailSettings(ctx context.Context) (*EmailSettingsResponse, error)
	UpdateEmailSettings(ctx context.Context, req UpdateEmailSettingsRequest) error
	TestEmailSettings(ctx context.Context, req TestEmailRequest) error

	// Tax Settings
	GetTaxSettings(ctx context.Context) (*TaxSettingsResponse, error)
	UpdateTaxSettings(ctx context.Context, req UpdateTaxSettingsRequest) error

	// Shipping Settings
	GetShippingSettings(ctx context.Context) (*ShippingSettingsResponse, error)
	UpdateShippingSettings(ctx context.Context, req UpdateShippingSettingsRequest) error

	// SEO Settings
	GetSEOSettings(ctx context.Context) (*SEOSettingsResponse, error)
	UpdateSEOSettings(ctx context.Context, req UpdateSEOSettingsRequest) error

	// Security Settings
	GetSecuritySettings(ctx context.Context) (*SecuritySettingsResponse, error)
	UpdateSecuritySettings(ctx context.Context, req UpdateSecuritySettingsRequest) error

	// Notification Settings
	GetNotificationSettings(ctx context.Context) (*NotificationSettingsResponse, error)
	UpdateNotificationSettings(ctx context.Context, req UpdateNotificationSettingsRequest) error

	// Integration Settings
	GetIntegrationSettings(ctx context.Context) (*IntegrationSettingsResponse, error)
	UpdateIntegrationSettings(ctx context.Context, req UpdateIntegrationSettingsRequest) error
}

// Request types
type UpdateGeneralSettingsRequest struct {
	StoreName        string `json:"store_name" validate:"required,min=1,max=100"`
	StoreDescription string `json:"store_description" validate:"max=500"`
	StoreEmail       string `json:"store_email" validate:"required,email"`
	StorePhone       string `json:"store_phone" validate:"max=20"`
	StoreAddress     string `json:"store_address" validate:"max=500"`
	Currency         string `json:"currency" validate:"required,len=3"`
	Timezone         string `json:"timezone" validate:"required"`
	Language         string `json:"language" validate:"required,len=2"`
	DateFormat       string `json:"date_format" validate:"required"`
	TimeFormat       string `json:"time_format" validate:"required"`
}

type UpdateStoreConfigRequest struct {
	MaintenanceMode    bool   `json:"maintenance_mode"`
	MaintenanceMessage string `json:"maintenance_message"`
	AllowRegistration  bool   `json:"allow_registration"`
	RequireEmailVerify bool   `json:"require_email_verify"`
	AutoApproveReviews bool   `json:"auto_approve_reviews"`
	EnableWishlist     bool   `json:"enable_wishlist"`
	EnableCompare      bool   `json:"enable_compare"`
	EnableReviews      bool   `json:"enable_reviews"`
	EnableCoupons      bool   `json:"enable_coupons"`
	EnableInventory    bool   `json:"enable_inventory"`
	EnableMultiVendor  bool   `json:"enable_multi_vendor"`
}

type UpdatePaymentSettingsRequest struct {
	DefaultCurrency    string                 `json:"default_currency" validate:"required,len=3"`
	AcceptedCurrencies []string               `json:"accepted_currencies"`
	StripeEnabled      bool                   `json:"stripe_enabled"`
	StripePublicKey    string                 `json:"stripe_public_key"`
	StripeSecretKey    string                 `json:"stripe_secret_key"`
	PayPalEnabled      bool                   `json:"paypal_enabled"`
	PayPalClientID     string                 `json:"paypal_client_id"`
	PayPalClientSecret string                 `json:"paypal_client_secret"`
	PayPalSandbox      bool                   `json:"paypal_sandbox"`
	CashOnDelivery     bool                   `json:"cash_on_delivery"`
	BankTransfer       bool                   `json:"bank_transfer"`
	BankDetails        map[string]interface{} `json:"bank_details"`
}

type UpdateEmailSettingsRequest struct {
	SMTPHost     string `json:"smtp_host" validate:"required"`
	SMTPPort     int    `json:"smtp_port" validate:"required,min=1,max=65535"`
	SMTPUsername string `json:"smtp_username" validate:"required"`
	SMTPPassword string `json:"smtp_password" validate:"required"`
	SMTPSecurity string `json:"smtp_security" validate:"oneof=none tls ssl"`
	FromEmail    string `json:"from_email" validate:"required,email"`
	FromName     string `json:"from_name" validate:"required"`
	ReplyToEmail string `json:"reply_to_email" validate:"email"`
}

type TestEmailRequest struct {
	ToEmail string `json:"to_email" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type UpdateTaxSettingsRequest struct {
	TaxEnabled       bool                   `json:"tax_enabled"`
	TaxIncluded      bool                   `json:"tax_included"`
	DefaultTaxRate   float64                `json:"default_tax_rate" validate:"min=0,max=100"`
	TaxCalculation   string                 `json:"tax_calculation" validate:"oneof=exclusive inclusive"`
	TaxRates         []TaxRateRequest       `json:"tax_rates"`
	TaxClasses       []TaxClassRequest      `json:"tax_classes"`
	TaxExemptions    []string               `json:"tax_exemptions"`
	DigitalTaxRate   float64                `json:"digital_tax_rate" validate:"min=0,max=100"`
	ShippingTaxable  bool                   `json:"shipping_taxable"`
	TaxDisplayFormat string                 `json:"tax_display_format" validate:"oneof=excluding including both"`
	TaxRounding      string                 `json:"tax_rounding" validate:"oneof=round_up round_down round_nearest"`
	CustomTaxRules   map[string]interface{} `json:"custom_tax_rules"`
}

type TaxRateRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Rate        float64   `json:"rate" validate:"required,min=0,max=100"`
	Country     string    `json:"country" validate:"required,len=2"`
	State       string    `json:"state"`
	City        string    `json:"city"`
	ZipCode     string    `json:"zip_code"`
	TaxClass    string    `json:"tax_class"`
	Priority    int       `json:"priority"`
	Compound    bool      `json:"compound"`
	IsActive    bool      `json:"is_active"`
	Description string    `json:"description"`
}

type TaxClassRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
}

type UpdateShippingSettingsRequest struct {
	ShippingEnabled     bool                     `json:"shipping_enabled"`
	FreeShippingEnabled bool                     `json:"free_shipping_enabled"`
	FreeShippingAmount  float64                  `json:"free_shipping_amount" validate:"min=0"`
	DefaultShippingCost float64                  `json:"default_shipping_cost" validate:"min=0"`
	ShippingCalculation string                   `json:"shipping_calculation" validate:"oneof=flat_rate weight_based distance_based"`
	ShippingZones       []ShippingZoneRequest    `json:"shipping_zones"`
	ShippingMethods     []ShippingMethodRequest  `json:"shipping_methods"`
	ShippingClasses     []ShippingClassRequest   `json:"shipping_classes"`
	PackagingOptions    []PackagingOptionRequest `json:"packaging_options"`
	DimensionUnit       string                   `json:"dimension_unit" validate:"oneof=cm in"`
	WeightUnit          string                   `json:"weight_unit" validate:"oneof=kg lb"`
}

type ShippingZoneRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Countries   []string  `json:"countries"`
	States      []string  `json:"states"`
	ZipCodes    []string  `json:"zip_codes"`
	IsActive    bool      `json:"is_active"`
	Priority    int       `json:"priority"`
	Description string    `json:"description"`
}

type ShippingMethodRequest struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name" validate:"required"`
	Type        string                 `json:"type" validate:"oneof=flat_rate free_shipping local_pickup"`
	Cost        float64                `json:"cost" validate:"min=0"`
	MinAmount   float64                `json:"min_amount" validate:"min=0"`
	MaxAmount   float64                `json:"max_amount" validate:"min=0"`
	ZoneID      uuid.UUID              `json:"zone_id"`
	IsActive    bool                   `json:"is_active"`
	Settings    map[string]interface{} `json:"settings"`
	Description string                 `json:"description"`
}

type ShippingClassRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Slug        string    `json:"slug" validate:"required"`
	Description string    `json:"description"`
	Cost        float64   `json:"cost" validate:"min=0"`
	IsActive    bool      `json:"is_active"`
}

type PackagingOptionRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Length      float64   `json:"length" validate:"min=0"`
	Width       float64   `json:"width" validate:"min=0"`
	Height      float64   `json:"height" validate:"min=0"`
	Weight      float64   `json:"weight" validate:"min=0"`
	Cost        float64   `json:"cost" validate:"min=0"`
	IsDefault   bool      `json:"is_default"`
	IsActive    bool      `json:"is_active"`
	Description string    `json:"description"`
}

type UpdateSEOSettingsRequest struct {
	SiteTitle        string            `json:"site_title" validate:"required,max=60"`
	SiteDescription  string            `json:"site_description" validate:"required,max=160"`
	SiteKeywords     []string          `json:"site_keywords"`
	MetaRobots       string            `json:"meta_robots" validate:"oneof=index,follow noindex,nofollow index,nofollow noindex,follow"`
	CanonicalURL     string            `json:"canonical_url" validate:"url"`
	OGTitle          string            `json:"og_title" validate:"max=60"`
	OGDescription    string            `json:"og_description" validate:"max=160"`
	OGImage          string            `json:"og_image" validate:"url"`
	TwitterCard      string            `json:"twitter_card" validate:"oneof=summary summary_large_image app player"`
	TwitterSite      string            `json:"twitter_site"`
	TwitterCreator   string            `json:"twitter_creator"`
	GoogleAnalytics  string            `json:"google_analytics"`
	GoogleTagManager string            `json:"google_tag_manager"`
	FacebookPixel    string            `json:"facebook_pixel"`
	CustomMeta       map[string]string `json:"custom_meta"`
	Sitemap          SitemapSettings   `json:"sitemap"`
	Robots           RobotsSettings    `json:"robots"`
}

type SitemapSettings struct {
	Enabled           bool     `json:"enabled"`
	IncludeProducts   bool     `json:"include_products"`
	IncludeCategories bool     `json:"include_categories"`
	IncludePages      bool     `json:"include_pages"`
	ExcludeURLs       []string `json:"exclude_urls"`
	ChangeFreq        string   `json:"change_freq" validate:"oneof=always hourly daily weekly monthly yearly never"`
	Priority          float64  `json:"priority" validate:"min=0,max=1"`
}

type RobotsSettings struct {
	Content     string   `json:"content"`
	Disallow    []string `json:"disallow"`
	Allow       []string `json:"allow"`
	Crawldelay  int      `json:"crawldelay" validate:"min=0"`
	SitemapURLs []string `json:"sitemap_urls"`
}

// Response types
type GeneralSettingsResponse struct {
	StoreName        string    `json:"store_name"`
	StoreDescription string    `json:"store_description"`
	StoreEmail       string    `json:"store_email"`
	StorePhone       string    `json:"store_phone"`
	StoreAddress     string    `json:"store_address"`
	Currency         string    `json:"currency"`
	Timezone         string    `json:"timezone"`
	Language         string    `json:"language"`
	DateFormat       string    `json:"date_format"`
	TimeFormat       string    `json:"time_format"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type StoreConfigResponse struct {
	MaintenanceMode    bool      `json:"maintenance_mode"`
	MaintenanceMessage string    `json:"maintenance_message"`
	AllowRegistration  bool      `json:"allow_registration"`
	RequireEmailVerify bool      `json:"require_email_verify"`
	AutoApproveReviews bool      `json:"auto_approve_reviews"`
	EnableWishlist     bool      `json:"enable_wishlist"`
	EnableCompare      bool      `json:"enable_compare"`
	EnableReviews      bool      `json:"enable_reviews"`
	EnableCoupons      bool      `json:"enable_coupons"`
	EnableInventory    bool      `json:"enable_inventory"`
	EnableMultiVendor  bool      `json:"enable_multi_vendor"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type PaymentSettingsResponse struct {
	DefaultCurrency    string                 `json:"default_currency"`
	AcceptedCurrencies []string               `json:"accepted_currencies"`
	StripeEnabled      bool                   `json:"stripe_enabled"`
	StripePublicKey    string                 `json:"stripe_public_key"`
	PayPalEnabled      bool                   `json:"paypal_enabled"`
	PayPalClientID     string                 `json:"paypal_client_id"`
	PayPalSandbox      bool                   `json:"paypal_sandbox"`
	CashOnDelivery     bool                   `json:"cash_on_delivery"`
	BankTransfer       bool                   `json:"bank_transfer"`
	BankDetails        map[string]interface{} `json:"bank_details"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

type EmailSettingsResponse struct {
	SMTPHost     string    `json:"smtp_host"`
	SMTPPort     int       `json:"smtp_port"`
	SMTPUsername string    `json:"smtp_username"`
	SMTPSecurity string    `json:"smtp_security"`
	FromEmail    string    `json:"from_email"`
	FromName     string    `json:"from_name"`
	ReplyToEmail string    `json:"reply_to_email"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TaxSettingsResponse struct {
	TaxEnabled       bool                   `json:"tax_enabled"`
	TaxIncluded      bool                   `json:"tax_included"`
	DefaultTaxRate   float64                `json:"default_tax_rate"`
	TaxCalculation   string                 `json:"tax_calculation"`
	TaxRates         []TaxRateResponse      `json:"tax_rates"`
	TaxClasses       []TaxClassResponse     `json:"tax_classes"`
	TaxExemptions    []string               `json:"tax_exemptions"`
	DigitalTaxRate   float64                `json:"digital_tax_rate"`
	ShippingTaxable  bool                   `json:"shipping_taxable"`
	TaxDisplayFormat string                 `json:"tax_display_format"`
	TaxRounding      string                 `json:"tax_rounding"`
	CustomTaxRules   map[string]interface{} `json:"custom_tax_rules"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

type TaxRateResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Rate        float64   `json:"rate"`
	Country     string    `json:"country"`
	State       string    `json:"state"`
	City        string    `json:"city"`
	ZipCode     string    `json:"zip_code"`
	TaxClass    string    `json:"tax_class"`
	Priority    int       `json:"priority"`
	Compound    bool      `json:"compound"`
	IsActive    bool      `json:"is_active"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaxClassResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ShippingSettingsResponse struct {
	ShippingEnabled     bool                             `json:"shipping_enabled"`
	FreeShippingEnabled bool                             `json:"free_shipping_enabled"`
	FreeShippingAmount  float64                          `json:"free_shipping_amount"`
	DefaultShippingCost float64                          `json:"default_shipping_cost"`
	ShippingCalculation string                           `json:"shipping_calculation"`
	ShippingZones       []ShippingZoneResponse           `json:"shipping_zones"`
	ShippingMethods     []SettingsShippingMethodResponse `json:"shipping_methods"`
	ShippingClasses     []ShippingClassResponse          `json:"shipping_classes"`
	PackagingOptions    []PackagingOptionResponse        `json:"packaging_options"`
	DimensionUnit       string                           `json:"dimension_unit"`
	WeightUnit          string                           `json:"weight_unit"`
	UpdatedAt           time.Time                        `json:"updated_at"`
}

type ShippingZoneResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Countries   []string  `json:"countries"`
	States      []string  `json:"states"`
	ZipCodes    []string  `json:"zip_codes"`
	IsActive    bool      `json:"is_active"`
	Priority    int       `json:"priority"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SettingsShippingMethodResponse struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Cost        float64                `json:"cost"`
	MinAmount   float64                `json:"min_amount"`
	MaxAmount   float64                `json:"max_amount"`
	ZoneID      uuid.UUID              `json:"zone_id"`
	IsActive    bool                   `json:"is_active"`
	Settings    map[string]interface{} `json:"settings"`
	Description string                 `json:"description"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type ShippingClassResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Cost        float64   `json:"cost"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PackagingOptionResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Length      float64   `json:"length"`
	Width       float64   `json:"width"`
	Height      float64   `json:"height"`
	Weight      float64   `json:"weight"`
	Cost        float64   `json:"cost"`
	IsDefault   bool      `json:"is_default"`
	IsActive    bool      `json:"is_active"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SEOSettingsResponse struct {
	SiteTitle        string            `json:"site_title"`
	SiteDescription  string            `json:"site_description"`
	SiteKeywords     []string          `json:"site_keywords"`
	MetaRobots       string            `json:"meta_robots"`
	CanonicalURL     string            `json:"canonical_url"`
	OGTitle          string            `json:"og_title"`
	OGDescription    string            `json:"og_description"`
	OGImage          string            `json:"og_image"`
	TwitterCard      string            `json:"twitter_card"`
	TwitterSite      string            `json:"twitter_site"`
	TwitterCreator   string            `json:"twitter_creator"`
	GoogleAnalytics  string            `json:"google_analytics"`
	GoogleTagManager string            `json:"google_tag_manager"`
	FacebookPixel    string            `json:"facebook_pixel"`
	CustomMeta       map[string]string `json:"custom_meta"`
	Sitemap          SitemapSettings   `json:"sitemap"`
	Robots           RobotsSettings    `json:"robots"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

type SecuritySettingsResponse struct {
	TwoFactorEnabled  bool      `json:"two_factor_enabled"`
	SessionTimeout    int       `json:"session_timeout"`
	PasswordPolicy    string    `json:"password_policy"`
	MaxLoginAttempts  int       `json:"max_login_attempts"`
	LockoutDuration   int       `json:"lockout_duration"`
	RequireSSL        bool      `json:"require_ssl"`
	IPWhitelist       []string  `json:"ip_whitelist"`
	IPBlacklist       []string  `json:"ip_blacklist"`
	EnableCaptcha     bool      `json:"enable_captcha"`
	CaptchaProvider   string    `json:"captcha_provider"`
	EnableAuditLog    bool      `json:"enable_audit_log"`
	DataRetentionDays int       `json:"data_retention_days"`
	EncryptionEnabled bool      `json:"encryption_enabled"`
	BackupEncryption  bool      `json:"backup_encryption"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type NotificationSettingsResponse struct {
	EmailNotifications  bool      `json:"email_notifications"`
	SMSNotifications    bool      `json:"sms_notifications"`
	PushNotifications   bool      `json:"push_notifications"`
	OrderNotifications  bool      `json:"order_notifications"`
	StockNotifications  bool      `json:"stock_notifications"`
	ReviewNotifications bool      `json:"review_notifications"`
	SecurityAlerts      bool      `json:"security_alerts"`
	MarketingEmails     bool      `json:"marketing_emails"`
	NewsletterEnabled   bool      `json:"newsletter_enabled"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type IntegrationSettingsResponse struct {
	GoogleAnalytics GoogleAnalyticsSettings `json:"google_analytics"`
	FacebookPixel   FacebookPixelSettings   `json:"facebook_pixel"`
	MailChimp       MailChimpSettings       `json:"mailchimp"`
	Zapier          ZapierSettings          `json:"zapier"`
	Webhooks        []WebhookSettings       `json:"webhooks"`
	UpdatedAt       time.Time               `json:"updated_at"`
}

type GoogleAnalyticsSettings struct {
	Enabled      bool   `json:"enabled"`
	TrackingID   string `json:"tracking_id"`
	Enhanced     bool   `json:"enhanced"`
	Demographics bool   `json:"demographics"`
	Advertising  bool   `json:"advertising"`
	Ecommerce    bool   `json:"ecommerce"`
}

type FacebookPixelSettings struct {
	Enabled  bool   `json:"enabled"`
	PixelID  string `json:"pixel_id"`
	Advanced bool   `json:"advanced"`
}

type MailChimpSettings struct {
	Enabled     bool   `json:"enabled"`
	APIKey      string `json:"api_key"`
	ListID      string `json:"list_id"`
	DoubleOptin bool   `json:"double_optin"`
}

type ZapierSettings struct {
	Enabled bool   `json:"enabled"`
	APIKey  string `json:"api_key"`
}

type WebhookSettings struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	Secret    string    `json:"secret"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateSecuritySettingsRequest struct {
	TwoFactorEnabled  bool     `json:"two_factor_enabled"`
	SessionTimeout    int      `json:"session_timeout" validate:"min=5,max=1440"`
	PasswordPolicy    string   `json:"password_policy" validate:"oneof=weak medium strong"`
	MaxLoginAttempts  int      `json:"max_login_attempts" validate:"min=3,max=10"`
	LockoutDuration   int      `json:"lockout_duration" validate:"min=5,max=1440"`
	RequireSSL        bool     `json:"require_ssl"`
	IPWhitelist       []string `json:"ip_whitelist"`
	IPBlacklist       []string `json:"ip_blacklist"`
	EnableCaptcha     bool     `json:"enable_captcha"`
	CaptchaProvider   string   `json:"captcha_provider" validate:"oneof=recaptcha hcaptcha"`
	EnableAuditLog    bool     `json:"enable_audit_log"`
	DataRetentionDays int      `json:"data_retention_days" validate:"min=30,max=2555"`
	EncryptionEnabled bool     `json:"encryption_enabled"`
	BackupEncryption  bool     `json:"backup_encryption"`
}

type UpdateNotificationSettingsRequest struct {
	EmailNotifications  bool `json:"email_notifications"`
	SMSNotifications    bool `json:"sms_notifications"`
	PushNotifications   bool `json:"push_notifications"`
	OrderNotifications  bool `json:"order_notifications"`
	StockNotifications  bool `json:"stock_notifications"`
	ReviewNotifications bool `json:"review_notifications"`
	SecurityAlerts      bool `json:"security_alerts"`
	MarketingEmails     bool `json:"marketing_emails"`
	NewsletterEnabled   bool `json:"newsletter_enabled"`
}

type UpdateIntegrationSettingsRequest struct {
	GoogleAnalytics GoogleAnalyticsSettings `json:"google_analytics"`
	FacebookPixel   FacebookPixelSettings   `json:"facebook_pixel"`
	MailChimp       MailChimpSettings       `json:"mailchimp"`
	Zapier          ZapierSettings          `json:"zapier"`
	Webhooks        []WebhookSettings       `json:"webhooks"`
}
