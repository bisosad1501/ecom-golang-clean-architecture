package entities

import (
	"time"

	"github.com/google/uuid"
)

// Setting represents a system setting
type Setting struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Key         string    `json:"key" gorm:"uniqueIndex;not null"`
	Value       string    `json:"value" gorm:"type:text"`
	Type        string    `json:"type" gorm:"not null"` // string, number, boolean, json
	Category    string    `json:"category" gorm:"not null;index"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public" gorm:"default:false"`
	IsEncrypted bool      `json:"is_encrypted" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the table name for Setting
func (Setting) TableName() string {
	return "settings"
}

// TaxRate represents a tax rate configuration
type TaxRate struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Rate        float64   `json:"rate" gorm:"not null"`
	Country     string    `json:"country" gorm:"not null;size:2"`
	State       string    `json:"state" gorm:"size:50"`
	City        string    `json:"city" gorm:"size:100"`
	ZipCode     string    `json:"zip_code" gorm:"size:20"`
	TaxClass    string    `json:"tax_class" gorm:"size:50"`
	Priority    int       `json:"priority" gorm:"default:0"`
	Compound    bool      `json:"compound" gorm:"default:false"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the table name for TaxRate
func (TaxRate) TableName() string {
	return "tax_rates"
}

// TaxClass represents a tax class
type TaxClass struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null;uniqueIndex"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the table name for TaxClass
func (TaxClass) TableName() string {
	return "tax_classes"
}

// Note: ShippingZone and ShippingMethod are already defined in shipping.go

// ShippingClass represents a shipping class
type ShippingClass struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"not null;uniqueIndex"`
	Description string    `json:"description"`
	Cost        float64   `json:"cost" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the table name for ShippingClass
func (ShippingClass) TableName() string {
	return "shipping_classes"
}

// PackagingOption represents a packaging option
type PackagingOption struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Length      float64   `json:"length" gorm:"default:0"`
	Width       float64   `json:"width" gorm:"default:0"`
	Height      float64   `json:"height" gorm:"default:0"`
	Weight      float64   `json:"weight" gorm:"default:0"`
	Cost        float64   `json:"cost" gorm:"default:0"`
	IsDefault   bool      `json:"is_default" gorm:"default:false"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the table name for PackagingOption
func (PackagingOption) TableName() string {
	return "packaging_options"
}

// Note: EmailTemplate is already defined in email.go

// Campaign represents a marketing campaign
type Campaign struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name           string     `json:"name" gorm:"not null"`
	Description    string     `json:"description"`
	Type           string     `json:"type" gorm:"not null"`                 // email, social, display, affiliate
	Status         string     `json:"status" gorm:"not null;default:draft"` // draft, active, paused, completed
	Budget         float64    `json:"budget" gorm:"default:0"`
	Spent          float64    `json:"spent" gorm:"default:0"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	TargetAudience string     `json:"target_audience" gorm:"type:text"` // JSON
	Goals          string     `json:"goals" gorm:"type:text"`           // JSON
	Settings       string     `json:"settings" gorm:"type:text"`        // JSON
	Metrics        string     `json:"metrics" gorm:"type:text"`         // JSON
	CreatedBy      uuid.UUID  `json:"created_by" gorm:"type:uuid"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TableName returns the table name for Campaign
func (Campaign) TableName() string {
	return "campaigns"
}

// Note: Promotion is already defined in coupon.go

// NewsletterSubscriber represents a newsletter subscriber
type NewsletterSubscriber struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email          string     `json:"email" gorm:"not null;uniqueIndex"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	Status         string     `json:"status" gorm:"not null;default:active"` // active, inactive, pending
	Segments       string     `json:"segments" gorm:"type:text"`             // JSON array
	Tags           string     `json:"tags" gorm:"type:text"`                 // JSON array
	CustomFields   string     `json:"custom_fields" gorm:"type:text"`        // JSON
	SubscribedAt   time.Time  `json:"subscribed_at"`
	UnsubscribedAt *time.Time `json:"unsubscribed_at"`
	LastEmailSent  *time.Time `json:"last_email_sent"`
	EmailsSent     int64      `json:"emails_sent" gorm:"default:0"`
	EmailsOpened   int64      `json:"emails_opened" gorm:"default:0"`
	EmailsClicked  int64      `json:"emails_clicked" gorm:"default:0"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TableName returns the table name for NewsletterSubscriber
func (NewsletterSubscriber) TableName() string {
	return "newsletter_subscribers"
}

// EmailCampaign represents an email campaign
type EmailCampaign struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null"`
	TemplateID  uuid.UUID      `json:"template_id" gorm:"type:uuid"`
	Template    *EmailTemplate `json:"template,omitempty" gorm:"foreignKey:TemplateID"`
	Recipients  string         `json:"recipients" gorm:"type:text"`          // JSON
	Status      string         `json:"status" gorm:"not null;default:draft"` // draft, scheduled, sending, sent, failed
	ScheduledAt *time.Time     `json:"scheduled_at"`
	SentAt      *time.Time     `json:"sent_at"`
	Variables   string         `json:"variables" gorm:"type:text"` // JSON
	Settings    string         `json:"settings" gorm:"type:text"`  // JSON
	Metrics     string         `json:"metrics" gorm:"type:text"`   // JSON
	CreatedBy   uuid.UUID      `json:"created_by" gorm:"type:uuid"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// TableName returns the table name for EmailCampaign
func (EmailCampaign) TableName() string {
	return "email_campaigns"
}

// Note: LoyaltyProgram is already defined in coupon.go

// ReferralProgram represents a referral program
type ReferralProgram struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name           string    `json:"name" gorm:"not null"`
	Description    string    `json:"description"`
	IsActive       bool      `json:"is_active" gorm:"default:true"`
	ReferrerReward string    `json:"referrer_reward" gorm:"type:text"` // JSON
	RefereeReward  string    `json:"referee_reward" gorm:"type:text"`  // JSON
	Rules          string    `json:"rules" gorm:"type:text"`           // JSON
	Settings       string    `json:"settings" gorm:"type:text"`        // JSON
	Stats          string    `json:"stats" gorm:"type:text"`           // JSON
	CreatedBy      uuid.UUID `json:"created_by" gorm:"type:uuid"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName returns the table name for ReferralProgram
func (ReferralProgram) TableName() string {
	return "referral_programs"
}
