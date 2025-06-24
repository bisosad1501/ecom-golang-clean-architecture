package entities

import (
	"time"

	"github.com/google/uuid"
)

// AddressType represents the type of address
type AddressType string

const (
	AddressTypeShipping AddressType = "shipping"
	AddressTypeBilling  AddressType = "billing"
	AddressTypeBoth     AddressType = "both"
)

// Address represents a user address
type Address struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID   `json:"user_id" gorm:"type:uuid;not null;index"`
	User        User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Type        AddressType `json:"type" gorm:"not null;default:'shipping'"`
	FirstName   string      `json:"first_name" gorm:"not null" validate:"required"`
	LastName    string      `json:"last_name" gorm:"not null" validate:"required"`
	Company     string      `json:"company"`
	Address1    string      `json:"address1" gorm:"not null" validate:"required"`
	Address2    string      `json:"address2"`
	City        string      `json:"city" gorm:"not null" validate:"required"`
	State       string      `json:"state" gorm:"not null" validate:"required"`
	ZipCode     string      `json:"zip_code" gorm:"not null" validate:"required"`
	Country     string      `json:"country" gorm:"not null;default:'USA'" validate:"required"`
	Phone       string      `json:"phone"`
	IsDefault   bool        `json:"is_default" gorm:"default:false"`
	IsActive    bool        `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Address entity
func (Address) TableName() string {
	return "addresses"
}

// GetFullName returns the full name for the address
func (a *Address) GetFullName() string {
	return a.FirstName + " " + a.LastName
}

// GetFullAddress returns the complete address string
func (a *Address) GetFullAddress() string {
	address := a.Address1
	if a.Address2 != "" {
		address += ", " + a.Address2
	}
	address += ", " + a.City + ", " + a.State + " " + a.ZipCode
	if a.Country != "" {
		address += ", " + a.Country
	}
	return address
}

// IsShippingAddress checks if this is a shipping address
func (a *Address) IsShippingAddress() bool {
	return a.Type == AddressTypeShipping || a.Type == AddressTypeBoth
}

// IsBillingAddress checks if this is a billing address
func (a *Address) IsBillingAddress() bool {
	return a.Type == AddressTypeBilling || a.Type == AddressTypeBoth
}

// Wishlist represents a user's wishlist item
type Wishlist struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	Product   Product   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Wishlist entity
func (Wishlist) TableName() string {
	return "user_wishlists"
}

// UserPreference represents user preferences and settings
type UserPreference struct {
	ID                    uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID                uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	User                  User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Language              string    `json:"language" gorm:"default:'en'"`
	Currency              string    `json:"currency" gorm:"default:'USD'"`
	Timezone              string    `json:"timezone" gorm:"default:'UTC'"`
	EmailNotifications    bool      `json:"email_notifications" gorm:"default:true"`
	SMSNotifications      bool      `json:"sms_notifications" gorm:"default:false"`
	PushNotifications     bool      `json:"push_notifications" gorm:"default:true"`
	MarketingEmails       bool      `json:"marketing_emails" gorm:"default:true"`
	OrderUpdates          bool      `json:"order_updates" gorm:"default:true"`
	ProductRecommendations bool     `json:"product_recommendations" gorm:"default:true"`
	NewsletterSubscription bool     `json:"newsletter_subscription" gorm:"default:false"`
	CreatedAt             time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for UserPreference entity
func (UserPreference) TableName() string {
	return "user_preferences"
}

// AccountVerification represents account verification status
type AccountVerification struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID           uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	User             User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	EmailVerified    bool       `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt  *time.Time `json:"email_verified_at"`
	PhoneVerified    bool       `json:"phone_verified" gorm:"default:false"`
	PhoneVerifiedAt  *time.Time `json:"phone_verified_at"`
	VerificationCode string     `json:"-" gorm:"index"`
	CodeExpiresAt    *time.Time `json:"-"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for AccountVerification entity
func (AccountVerification) TableName() string {
	return "account_verifications"
}

// IsFullyVerified checks if both email and phone are verified
func (av *AccountVerification) IsFullyVerified() bool {
	return av.EmailVerified && av.PhoneVerified
}

// IsCodeValid checks if the verification code is still valid
func (av *AccountVerification) IsCodeValid() bool {
	if av.CodeExpiresAt == nil {
		return false
	}
	return time.Now().Before(*av.CodeExpiresAt)
}

// PasswordReset represents password reset requests
type PasswordReset struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	User      User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Token     string     `json:"-" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for PasswordReset entity
func (PasswordReset) TableName() string {
	return "password_resets"
}

// IsValid checks if the password reset token is still valid
func (pr *PasswordReset) IsValid() bool {
	return pr.UsedAt == nil && time.Now().Before(pr.ExpiresAt)
}

// MarkAsUsed marks the password reset token as used
func (pr *PasswordReset) MarkAsUsed() {
	now := time.Now()
	pr.UsedAt = &now
}
