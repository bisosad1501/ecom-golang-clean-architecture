package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus represents user account status
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusPending   UserStatus = "pending"
)

// UserRole represents the role of a user
type UserRole string

const (
	UserRoleCustomer  UserRole = "customer"
	UserRoleAdmin     UserRole = "admin"
	UserRoleModerator UserRole = "moderator"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string     `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password  string     `json:"-" gorm:"" validate:"omitempty,min=6"` // Made optional for OAuth users
	FirstName string     `json:"first_name" gorm:"not null" validate:"required"`
	LastName  string     `json:"last_name" gorm:"not null" validate:"required"`
	Phone     string     `json:"phone" gorm:"index"`
	Role      UserRole   `json:"role" gorm:"default:'customer'" validate:"required"`
	Status    UserStatus `json:"status" gorm:"default:'active'" validate:"required"`
	IsActive  bool       `json:"is_active" gorm:"default:true"`

	// OAuth fields
	GoogleID    string `json:"google_id,omitempty" gorm:"index"`
	FacebookID  string `json:"facebook_id,omitempty" gorm:"index"`
	Avatar      string `json:"avatar,omitempty"`
	IsOAuthUser bool   `json:"is_oauth_user" gorm:"default:false"`

	// Enhanced user fields
	Username    *string    `json:"username,omitempty" gorm:"index"` // Optional, non-unique display name
	DateOfBirth *time.Time `json:"date_of_birth"`
	Gender      string     `json:"gender"`
	Language    string     `json:"language" gorm:"default:'en'"`
	Timezone    string     `json:"timezone" gorm:"default:'UTC'"`
	Currency    string     `json:"currency" gorm:"default:'USD'"`

	// Account verification and status
	EmailVerified  bool       `json:"email_verified" gorm:"default:false"`
	PhoneVerified  bool       `json:"phone_verified" gorm:"default:false"`
	LastLoginAt    *time.Time `json:"last_login_at"`
	LastActivityAt *time.Time `json:"last_activity_at"`

	// Marketing preferences
	MarketingOptIn  bool `json:"marketing_opt_in" gorm:"default:false"`
	NewsletterOptIn bool `json:"newsletter_opt_in" gorm:"default:false"`

	// Security settings
	TwoFactorEnabled bool `json:"two_factor_enabled" gorm:"default:false"`
	SecurityScore    int  `json:"security_score" gorm:"default:0"`

	// Customer metrics
	TotalOrders    int     `json:"total_orders" gorm:"default:0"`
	TotalSpent     float64 `json:"total_spent" gorm:"default:0"`
	LoyaltyPoints  int     `json:"loyalty_points" gorm:"default:0"`
	MembershipTier string  `json:"membership_tier" gorm:"default:'bronze'"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Profile   *UserProfile `json:"profile,omitempty" gorm:"foreignKey:UserID"`
	Addresses []Address    `json:"addresses,omitempty" gorm:"foreignKey:UserID"`
	Wishlist  []Product    `json:"wishlist,omitempty" gorm:"many2many:user_wishlists;"`
}

// TableName returns the table name for User entity
func (User) TableName() string {
	return "users"
}

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

// IsModerator checks if the user is a moderator
func (u *User) IsModerator() bool {
	return u.Role == UserRoleModerator
}

// CanManageProducts checks if the user can manage products
func (u *User) CanManageProducts() bool {
	return u.Role == UserRoleAdmin || u.Role == UserRoleModerator
}

// IsVerified checks if the user is fully verified
func (u *User) IsVerified() bool {
	return u.EmailVerified && u.PhoneVerified
}

// IsVIP checks if the user is a VIP member
func (u *User) IsVIP() bool {
	return u.MembershipTier == "gold" || u.MembershipTier == "platinum" || u.MembershipTier == "diamond"
}

// GetSecurityLevel returns the security level based on security score
func (u *User) GetSecurityLevel() string {
	if u.SecurityScore >= 80 {
		return "high"
	} else if u.SecurityScore >= 50 {
		return "medium"
	}
	return "low"
}

// IsHighValue checks if the user is a high-value customer
func (u *User) IsHighValue() bool {
	return u.TotalSpent > 1000 || u.TotalOrders > 10
}

// GetCustomerSegment returns the customer segment
func (u *User) GetCustomerSegment() string {
	if u.TotalOrders == 0 {
		return "new"
	} else if u.TotalOrders < 5 {
		return "occasional"
	} else if u.TotalOrders < 20 {
		return "regular"
	}
	return "loyal"
}

// UpdateLastActivity updates the last activity timestamp
func (u *User) UpdateLastActivity() {
	now := time.Now()
	u.LastActivityAt = &now
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}

// UserProfile represents additional user profile information
type UserProfile struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	User        User       `json:"user" gorm:"foreignKey:UserID"`
	Avatar      string     `json:"avatar"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Gender      string     `json:"gender"`
	Address     string     `json:"address"`
	City        string     `json:"city"`
	State       string     `json:"state"`
	Country     string     `json:"country"`
	ZipCode     string     `json:"zip_code"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for UserProfile entity
func (UserProfile) TableName() string {
	return "user_profiles"
}

// UserSession represents an active user session
type UserSession struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	User         User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	SessionToken string    `json:"-" gorm:"uniqueIndex;not null"`
	DeviceInfo   string    `json:"device_info"`
	IPAddress    string    `json:"ip_address" gorm:"index"`
	UserAgent    string    `json:"user_agent"`
	Location     string    `json:"location"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	LastActivity time.Time `json:"last_activity" gorm:"index"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"index"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for UserSession entity
func (UserSession) TableName() string {
	return "user_sessions"
}

// IsExpired checks if the session is expired
func (us *UserSession) IsExpired() bool {
	return time.Now().After(us.ExpiresAt)
}

// IsValid checks if the session is valid and active
func (us *UserSession) IsValid() bool {
	return us.IsActive && !us.IsExpired()
}

// UserLoginHistory represents user login history
type UserLoginHistory struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	User       User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	IPAddress  string    `json:"ip_address" gorm:"index"`
	UserAgent  string    `json:"user_agent"`
	DeviceInfo string    `json:"device_info"`
	Location   string    `json:"location"`
	LoginType  string    `json:"login_type" gorm:"default:'password'"` // password, google, facebook, etc.
	Success    bool      `json:"success" gorm:"index"`
	FailReason string    `json:"fail_reason"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime;index"`
}

// TableName returns the table name for UserLoginHistory entity
func (UserLoginHistory) TableName() string {
	return "user_login_history"
}

// ActivityType represents the type of user activity
type ActivityType string

const (
	ActivityTypeLogin          ActivityType = "login"
	ActivityTypeLogout         ActivityType = "logout"
	ActivityTypeRegister       ActivityType = "register"
	ActivityTypeProfileUpdate  ActivityType = "profile_update"
	ActivityTypePasswordChange ActivityType = "password_change"
	ActivityTypeProductView    ActivityType = "product_view"
	ActivityTypeProductSearch  ActivityType = "product_search"
	ActivityTypeCartAdd        ActivityType = "cart_add"
	ActivityTypeCartRemove     ActivityType = "cart_remove"
	ActivityTypeWishlistAdd    ActivityType = "wishlist_add"
	ActivityTypeWishlistRemove ActivityType = "wishlist_remove"
	ActivityTypeOrderPlace     ActivityType = "order_place"
	ActivityTypeOrderCancel    ActivityType = "order_cancel"
	ActivityTypeReviewCreate   ActivityType = "review_create"
	ActivityTypeAddressAdd     ActivityType = "address_add"
)

// UserActivity represents user activity log
type UserActivity struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID    `json:"user_id" gorm:"type:uuid;not null;index"`
	User        User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Type        ActivityType `json:"type" gorm:"not null;index"`
	Description string       `json:"description"`
	EntityType  string       `json:"entity_type"` // product, order, review, etc.
	EntityID    *uuid.UUID   `json:"entity_id" gorm:"type:uuid;index"`
	Metadata    string       `json:"metadata" gorm:"type:jsonb"` // Additional data as JSON
	IPAddress   string       `json:"ip_address" gorm:"index"`
	UserAgent   string       `json:"user_agent"`
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime;index"`
}

// TableName returns the table name for UserActivity entity
func (UserActivity) TableName() string {
	return "user_activities"
}

// UserPreferences represents user preferences and settings
type UserPreferences struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	User   User      `json:"user,omitempty" gorm:"foreignKey:UserID"`

	// Display preferences
	Theme    string `json:"theme" gorm:"default:'system'"` // light, dark, system
	Language string `json:"language" gorm:"default:'en'"`
	Currency string `json:"currency" gorm:"default:'USD'"`
	Timezone string `json:"timezone" gorm:"default:'UTC'"`

	// Notification preferences
	EmailNotifications     bool `json:"email_notifications" gorm:"default:true"`
	SMSNotifications       bool `json:"sms_notifications" gorm:"default:false"`
	PushNotifications      bool `json:"push_notifications" gorm:"default:true"`
	MarketingEmails        bool `json:"marketing_emails" gorm:"default:false"`
	OrderUpdates           bool `json:"order_updates" gorm:"default:true"`
	ProductRecommendations bool `json:"product_recommendations" gorm:"default:true"`
	NewsletterSubscription bool `json:"newsletter_subscription" gorm:"default:false"`
	SecurityAlerts         bool `json:"security_alerts" gorm:"default:true"`

	// Privacy preferences
	ProfileVisibility      string `json:"profile_visibility" gorm:"default:'private'"` // public, private, friends
	ShowOnlineStatus       bool   `json:"show_online_status" gorm:"default:false"`
	AllowDataCollection    bool   `json:"allow_data_collection" gorm:"default:true"`
	AllowPersonalization   bool   `json:"allow_personalization" gorm:"default:true"`
	AllowThirdPartySharing bool   `json:"allow_third_party_sharing" gorm:"default:false"`

	// Shopping preferences
	DefaultShippingMethod string `json:"default_shipping_method"`
	DefaultPaymentMethod  string `json:"default_payment_method"`
	SavePaymentMethods    bool   `json:"save_payment_methods" gorm:"default:true"`
	AutoApplyCoupons      bool   `json:"auto_apply_coupons" gorm:"default:true"`
	WishlistVisibility    string `json:"wishlist_visibility" gorm:"default:'private'"` // public, private

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for UserPreferences entity
func (UserPreferences) TableName() string {
	return "user_preferences"
}

// UserVerification represents user verification records
type UserVerification struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	User         User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Type         string     `json:"type" gorm:"not null;index"` // email, phone, identity
	Token        string     `json:"token" gorm:"not null;index"`
	Code         string     `json:"code" gorm:"index"` // For OTP verification
	IsVerified   bool       `json:"is_verified" gorm:"default:false"`
	VerifiedAt   *time.Time `json:"verified_at"`
	ExpiresAt    time.Time  `json:"expires_at" gorm:"index"`
	AttemptCount int        `json:"attempt_count" gorm:"default:0"`
	MaxAttempts  int        `json:"max_attempts" gorm:"default:3"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for UserVerification entity
func (UserVerification) TableName() string {
	return "user_verifications"
}

// IsExpired checks if the verification token is expired
func (uv *UserVerification) IsExpired() bool {
	return time.Now().After(uv.ExpiresAt)
}

// CanAttempt checks if more attempts are allowed
func (uv *UserVerification) CanAttempt() bool {
	return uv.AttemptCount < uv.MaxAttempts
}

// IncrementAttempt increments the attempt count
func (uv *UserVerification) IncrementAttempt() {
	uv.AttemptCount++
	uv.UpdatedAt = time.Now()
}

// MarkAsVerified marks the verification as completed
func (uv *UserVerification) MarkAsVerified() {
	uv.IsVerified = true
	now := time.Now()
	uv.VerifiedAt = &now
	uv.UpdatedAt = now
}
