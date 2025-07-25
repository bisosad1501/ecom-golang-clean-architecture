package repositories

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"

	"github.com/google/uuid"
)

// SettingsRepository defines the interface for settings data access
type SettingsRepository interface {
	// Settings CRUD
	GetSetting(ctx context.Context, key string) (*entities.Setting, error)
	GetSettingsByCategory(ctx context.Context, category string) ([]*entities.Setting, error)
	GetAllSettings(ctx context.Context) ([]*entities.Setting, error)
	SetSetting(ctx context.Context, setting *entities.Setting) error
	SetSettings(ctx context.Context, settings []*entities.Setting) error
	DeleteSetting(ctx context.Context, key string) error

	// Tax Management
	GetTaxRates(ctx context.Context, filters TaxRateFilters) ([]*entities.TaxRate, error)
	GetTaxRate(ctx context.Context, id uuid.UUID) (*entities.TaxRate, error)
	CreateTaxRate(ctx context.Context, taxRate *entities.TaxRate) error
	UpdateTaxRate(ctx context.Context, taxRate *entities.TaxRate) error
	DeleteTaxRate(ctx context.Context, id uuid.UUID) error
	GetTaxClasses(ctx context.Context) ([]*entities.TaxClass, error)
	GetTaxClass(ctx context.Context, id uuid.UUID) (*entities.TaxClass, error)
	CreateTaxClass(ctx context.Context, taxClass *entities.TaxClass) error
	UpdateTaxClass(ctx context.Context, taxClass *entities.TaxClass) error
	DeleteTaxClass(ctx context.Context, id uuid.UUID) error

	// Shipping Management
	GetShippingZones(ctx context.Context) ([]*entities.ShippingZone, error)
	GetShippingZone(ctx context.Context, id uuid.UUID) (*entities.ShippingZone, error)
	CreateShippingZone(ctx context.Context, zone *entities.ShippingZone) error
	UpdateShippingZone(ctx context.Context, zone *entities.ShippingZone) error
	DeleteShippingZone(ctx context.Context, id uuid.UUID) error
	GetShippingMethods(ctx context.Context, zoneID *uuid.UUID) ([]*entities.ShippingMethod, error)
	GetShippingMethod(ctx context.Context, id uuid.UUID) (*entities.ShippingMethod, error)
	CreateShippingMethod(ctx context.Context, method *entities.ShippingMethod) error
	UpdateShippingMethod(ctx context.Context, method *entities.ShippingMethod) error
	DeleteShippingMethod(ctx context.Context, id uuid.UUID) error
	GetShippingClasses(ctx context.Context) ([]*entities.ShippingClass, error)
	GetShippingClass(ctx context.Context, id uuid.UUID) (*entities.ShippingClass, error)
	CreateShippingClass(ctx context.Context, class *entities.ShippingClass) error
	UpdateShippingClass(ctx context.Context, class *entities.ShippingClass) error
	DeleteShippingClass(ctx context.Context, id uuid.UUID) error
	GetPackagingOptions(ctx context.Context) ([]*entities.PackagingOption, error)
	GetPackagingOption(ctx context.Context, id uuid.UUID) (*entities.PackagingOption, error)
	CreatePackagingOption(ctx context.Context, option *entities.PackagingOption) error
	UpdatePackagingOption(ctx context.Context, option *entities.PackagingOption) error
	DeletePackagingOption(ctx context.Context, id uuid.UUID) error
}

// TaxRateFilters represents filters for tax rates
type TaxRateFilters struct {
	Country  string
	State    string
	City     string
	ZipCode  string
	TaxClass string
	IsActive *bool
	Limit    int
	Offset   int
}

// MarketingRepository defines the interface for marketing data access
type MarketingRepository interface {
	// Email Templates
	GetEmailTemplates(ctx context.Context, filters EmailTemplateFilters) ([]*entities.EmailTemplate, error)
	GetEmailTemplate(ctx context.Context, id uuid.UUID) (*entities.EmailTemplate, error)
	CreateEmailTemplate(ctx context.Context, template *entities.EmailTemplate) error
	UpdateEmailTemplate(ctx context.Context, template *entities.EmailTemplate) error
	DeleteEmailTemplate(ctx context.Context, id uuid.UUID) error
	CountEmailTemplates(ctx context.Context, filters EmailTemplateFilters) (int64, error)

	// Campaigns
	GetCampaigns(ctx context.Context, filters CampaignFilters) ([]*entities.Campaign, error)
	GetCampaign(ctx context.Context, id uuid.UUID) (*entities.Campaign, error)
	CreateCampaign(ctx context.Context, campaign *entities.Campaign) error
	UpdateCampaign(ctx context.Context, campaign *entities.Campaign) error
	DeleteCampaign(ctx context.Context, id uuid.UUID) error
	CountCampaigns(ctx context.Context, filters CampaignFilters) (int64, error)

	// Promotions
	GetPromotions(ctx context.Context, filters PromotionFilters) ([]*entities.Promotion, error)
	GetPromotion(ctx context.Context, id uuid.UUID) (*entities.Promotion, error)
	CreatePromotion(ctx context.Context, promotion *entities.Promotion) error
	UpdatePromotion(ctx context.Context, promotion *entities.Promotion) error
	DeletePromotion(ctx context.Context, id uuid.UUID) error
	CountPromotions(ctx context.Context, filters PromotionFilters) (int64, error)

	// Email Campaigns
	GetEmailCampaigns(ctx context.Context, filters EmailCampaignFilters) ([]*entities.EmailCampaign, error)
	GetEmailCampaign(ctx context.Context, id uuid.UUID) (*entities.EmailCampaign, error)
	CreateEmailCampaign(ctx context.Context, campaign *entities.EmailCampaign) error
	UpdateEmailCampaign(ctx context.Context, campaign *entities.EmailCampaign) error
	DeleteEmailCampaign(ctx context.Context, id uuid.UUID) error
	CountEmailCampaigns(ctx context.Context, filters EmailCampaignFilters) (int64, error)

	// Newsletter Subscribers
	GetNewsletterSubscribers(ctx context.Context, filters NewsletterSubscriberFilters) ([]*entities.NewsletterSubscriber, error)
	GetNewsletterSubscriber(ctx context.Context, id uuid.UUID) (*entities.NewsletterSubscriber, error)
	GetNewsletterSubscriberByEmail(ctx context.Context, email string) (*entities.NewsletterSubscriber, error)
	CreateNewsletterSubscriber(ctx context.Context, subscriber *entities.NewsletterSubscriber) error
	UpdateNewsletterSubscriber(ctx context.Context, subscriber *entities.NewsletterSubscriber) error
	DeleteNewsletterSubscriber(ctx context.Context, id uuid.UUID) error
	CountNewsletterSubscribers(ctx context.Context, filters NewsletterSubscriberFilters) (int64, error)

	// Loyalty Programs
	GetLoyaltyPrograms(ctx context.Context) ([]*entities.LoyaltyProgram, error)
	GetLoyaltyProgram(ctx context.Context, id uuid.UUID) (*entities.LoyaltyProgram, error)
	CreateLoyaltyProgram(ctx context.Context, program *entities.LoyaltyProgram) error
	UpdateLoyaltyProgram(ctx context.Context, program *entities.LoyaltyProgram) error
	DeleteLoyaltyProgram(ctx context.Context, id uuid.UUID) error

	// Referral Programs
	GetReferralPrograms(ctx context.Context) ([]*entities.ReferralProgram, error)
	GetReferralProgram(ctx context.Context, id uuid.UUID) (*entities.ReferralProgram, error)
	CreateReferralProgram(ctx context.Context, program *entities.ReferralProgram) error
	UpdateReferralProgram(ctx context.Context, program *entities.ReferralProgram) error
	DeleteReferralProgram(ctx context.Context, id uuid.UUID) error
}

// Filter types
type EmailTemplateFilters struct {
	Search   string
	Type     string
	IsActive *bool
	SortBy   string
	SortOrder string
	Limit    int
	Offset   int
}

type CampaignFilters struct {
	Search   string
	Status   string
	Type     string
	SortBy   string
	SortOrder string
	Limit    int
	Offset   int
}

type PromotionFilters struct {
	Search   string
	Type     string
	Status   string
	SortBy   string
	SortOrder string
	Limit    int
	Offset   int
}

type EmailCampaignFilters struct {
	Search   string
	Status   string
	SortBy   string
	SortOrder string
	Limit    int
	Offset   int
}

type NewsletterSubscriberFilters struct {
	Search   string
	Status   string
	Segment  string
	SortBy   string
	SortOrder string
	Limit    int
	Offset   int
}
