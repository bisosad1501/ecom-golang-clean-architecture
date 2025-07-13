package entities

import (
	"fmt"
	"regexp"
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

// Validate validates user data
func (u *User) Validate() error {
	// Validate required fields
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	if u.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if u.LastName == "" {
		return fmt.Errorf("last name is required")
	}

	// Validate email format
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(emailRegex, u.Email); !matched {
		return fmt.Errorf("invalid email format")
	}

	// Validate password for non-OAuth users
	if !u.IsOAuthUser && u.Password == "" {
		return fmt.Errorf("password is required for non-OAuth users")
	}

	// Validate phone format if provided
	if u.Phone != "" {
		phoneRegex := `^\+?[1-9]\d{1,14}$` // Basic international phone format
		if matched, _ := regexp.MatchString(phoneRegex, u.Phone); !matched {
			return fmt.Errorf("invalid phone format")
		}
	}

	// Validate role
	validRoles := []UserRole{UserRoleCustomer, UserRoleModerator, UserRoleAdmin}
	isValidRole := false
	for _, role := range validRoles {
		if u.Role == role {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
		return fmt.Errorf("invalid user role: %s", u.Role)
	}

	// Validate currency
	if u.Currency != "" {
		validCurrencies := []string{"USD", "EUR", "GBP", "JPY", "VND"}
		isValidCurrency := false
		for _, currency := range validCurrencies {
			if u.Currency == currency {
				isValidCurrency = true
				break
			}
		}
		if !isValidCurrency {
			return fmt.Errorf("invalid currency: %s", u.Currency)
		}
	}

	// Validate membership tier
	if u.MembershipTier != "" {
		validTiers := []string{"bronze", "silver", "gold", "platinum", "diamond"}
		isValidTier := false
		for _, tier := range validTiers {
			if u.MembershipTier == tier {
				isValidTier = true
				break
			}
		}
		if !isValidTier {
			return fmt.Errorf("invalid membership tier: %s", u.MembershipTier)
		}
	}

	// Validate metrics are non-negative
	if u.TotalOrders < 0 {
		return fmt.Errorf("total orders cannot be negative")
	}
	if u.TotalSpent < 0 {
		return fmt.Errorf("total spent cannot be negative")
	}
	if u.LoyaltyPoints < 0 {
		return fmt.Errorf("loyalty points cannot be negative")
	}
	if u.SecurityScore < 0 || u.SecurityScore > 100 {
		return fmt.Errorf("security score must be between 0 and 100")
	}

	return nil
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
	ActivityTypeSecurityUpdate ActivityType = "security_update"
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

// UserSearchHistory represents user's search history
type UserSearchHistory struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Query     string    `json:"query" gorm:"not null"`
	Filters   string    `json:"filters" gorm:"type:jsonb"` // JSON string of applied filters
	Results   int       `json:"results" gorm:"default:0"`  // Number of results returned
	Clicked   bool      `json:"clicked" gorm:"default:false"` // Whether user clicked on any result
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for UserSearchHistory entity
func (UserSearchHistory) TableName() string {
	return "user_search_history"
}

// SavedSearch represents user's saved searches with alerts
type SavedSearch struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	User        User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Name        string     `json:"name" gorm:"not null"` // User-defined name for the search
	Query       string     `json:"query" gorm:"not null"`
	Filters     string     `json:"filters" gorm:"type:jsonb"` // JSON string of filters
	IsActive    bool       `json:"is_active" gorm:"default:true"`

	// Alert settings
	PriceAlert     bool    `json:"price_alert" gorm:"default:false"`
	StockAlert     bool    `json:"stock_alert" gorm:"default:false"`
	NewItemAlert   bool    `json:"new_item_alert" gorm:"default:false"`
	MaxPrice       *float64 `json:"max_price,omitempty"`
	MinPrice       *float64 `json:"min_price,omitempty"`

	// Notification settings
	EmailNotify    bool `json:"email_notify" gorm:"default:true"`
	PushNotify     bool `json:"push_notify" gorm:"default:false"`

	LastChecked    *time.Time `json:"last_checked,omitempty"`
	LastNotified   *time.Time `json:"last_notified,omitempty"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for SavedSearch entity
func (SavedSearch) TableName() string {
	return "saved_searches"
}

// UserBrowsingHistory represents user's product browsing history
type UserBrowsingHistory struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	User       User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ProductID  uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	Product    Product   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	CategoryID *uuid.UUID `json:"category_id,omitempty" gorm:"type:uuid;index"`
	Category   *Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`

	// Interaction details
	ViewDuration int    `json:"view_duration" gorm:"default:0"` // in seconds
	Source       string `json:"source"` // search, category, recommendation, etc.
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for UserBrowsingHistory entity
func (UserBrowsingHistory) TableName() string {
	return "user_browsing_history"
}

// UserPersonalization represents user's personalization data
type UserPersonalization struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	User   User      `json:"user,omitempty" gorm:"foreignKey:UserID"`

	// Preference scores (0-100)
	CategoryPreferences  string `json:"category_preferences" gorm:"type:jsonb"`  // JSON map of category_id -> score
	BrandPreferences     string `json:"brand_preferences" gorm:"type:jsonb"`     // JSON map of brand_id -> score
	PriceRangePreference string `json:"price_range_preference" gorm:"type:jsonb"` // JSON object with min/max preferences

	// Behavioral data
	AverageOrderValue    float64 `json:"average_order_value" gorm:"default:0"`
	PreferredShoppingTime string  `json:"preferred_shopping_time"` // morning, afternoon, evening, night
	ShoppingFrequency    string  `json:"shopping_frequency"` // daily, weekly, monthly, occasional

	// Recommendation settings
	RecommendationEngine string `json:"recommendation_engine" gorm:"default:'collaborative'"` // collaborative, content_based, hybrid
	PersonalizationLevel string `json:"personalization_level" gorm:"default:'medium'"` // low, medium, high

	// Analytics data
	TotalViews           int       `json:"total_views" gorm:"default:0"`
	TotalSearches        int       `json:"total_searches" gorm:"default:0"`
	UniqueProductsViewed int       `json:"unique_products_viewed" gorm:"default:0"`
	LastAnalyzed         *time.Time `json:"last_analyzed,omitempty"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for UserPersonalization entity
func (UserPersonalization) TableName() string {
	return "user_personalization"
}

// UserActivityLog represents detailed user activity logging
type UserActivityLog struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	User       User      `json:"user,omitempty" gorm:"foreignKey:UserID"`

	// Activity details
	ActivityType string `json:"activity_type" gorm:"not null;index"` // search, view, purchase, etc.
	EntityType   string `json:"entity_type"` // product, category, brand, etc.
	EntityID     *uuid.UUID `json:"entity_id,omitempty" gorm:"type:uuid;index"`

	// Context data
	SessionID    string `json:"session_id" gorm:"index"`
	Source       string `json:"source"` // web, mobile, api
	Page         string `json:"page"`
	Referrer     string `json:"referrer"`

	// Metadata
	Metadata  string `json:"metadata" gorm:"type:jsonb"` // Additional context as JSON
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for UserActivityLog entity
func (UserActivityLog) TableName() string {
	return "user_activity_logs"
}

// UserVerification represents user verification records
type UserVerification struct {
	ID                uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID            uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	User              User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	EmailVerified     bool       `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at"`
	PhoneVerified     bool       `json:"phone_verified" gorm:"default:false"`
	PhoneVerifiedAt   *time.Time `json:"phone_verified_at"`
	VerificationCode  string     `json:"verification_code" gorm:"index"`
	CodeExpiresAt     *time.Time `json:"code_expires_at"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
}

// Legacy UserVerification for backward compatibility
type LegacyUserVerification struct {
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
	return "account_verifications"
}

// IsExpired checks if the verification code is expired
func (uv *UserVerification) IsExpired() bool {
	if uv.CodeExpiresAt == nil {
		return false
	}
	return time.Now().After(*uv.CodeExpiresAt)
}

// CanAttempt checks if more attempts are allowed (always true for new structure)
func (uv *UserVerification) CanAttempt() bool {
	return true
}

// IncrementAttempt increments the attempt count (no-op for new structure)
func (uv *UserVerification) IncrementAttempt() {
	now := time.Now()
	uv.UpdatedAt = &now
}

// MarkAsVerified marks the verification as completed
func (uv *UserVerification) MarkAsVerified() {
	now := time.Now()
	uv.EmailVerified = true
	uv.EmailVerifiedAt = &now
	uv.UpdatedAt = &now
}

// UserOrderStats represents user order statistics (for optimization)
type UserOrderStats struct {
	TotalOrders int64   `json:"total_orders"`
	TotalSpent  float64 `json:"total_spent"`
}
