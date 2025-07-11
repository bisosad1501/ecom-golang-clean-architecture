package repositories

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entities.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)

	// GetByGoogleID retrieves a user by Google ID
	GetByGoogleID(ctx context.Context, googleID string) (*entities.User, error)

	// GetByFacebookID retrieves a user by Facebook ID
	GetByFacebookID(ctx context.Context, facebookID string) (*entities.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entities.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves users with pagination
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)

	// ExistsByEmail checks if a user exists with the given email
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// UpdatePassword updates user password
	UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error

	// SetActive sets user active status
	SetActive(ctx context.Context, userID uuid.UUID, isActive bool) error

	// Additional methods for admin dashboard
	CountUsers(ctx context.Context) (int64, error)
	CountActiveUsers(ctx context.Context) (int64, error)

	// Enhanced user methods
	GetUsersWithFilters(ctx context.Context, filters UserFilters) ([]*entities.User, error)
	CountUsersWithFilters(ctx context.Context, filters UserFilters) (int64, error)
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	UpdateLastActivity(ctx context.Context, userID uuid.UUID) error
	GetUsersByRole(ctx context.Context, role entities.UserRole, limit, offset int) ([]*entities.User, error)
	GetUsersByStatus(ctx context.Context, status entities.UserStatus, limit, offset int) ([]*entities.User, error)
	GetHighValueCustomers(ctx context.Context, limit int) ([]*entities.User, error)
	GetRecentlyRegistered(ctx context.Context, limit int) ([]*entities.User, error)
}

// UserFilters represents filters for user queries
type UserFilters struct {
	Role             *entities.UserRole   `json:"role"`
	Status           *entities.UserStatus `json:"status"`
	IsActive         *bool                `json:"is_active"`
	EmailVerified    *bool                `json:"email_verified"`
	PhoneVerified    *bool                `json:"phone_verified"`
	TwoFactorEnabled *bool                `json:"two_factor_enabled"`
	MembershipTier   string               `json:"membership_tier"`
	MinTotalSpent    *float64             `json:"min_total_spent"`
	MaxTotalSpent    *float64             `json:"max_total_spent"`
	MinTotalOrders   *int                 `json:"min_total_orders"`
	MaxTotalOrders   *int                 `json:"max_total_orders"`
	CreatedFrom      *time.Time           `json:"created_from"`
	CreatedTo        *time.Time           `json:"created_to"`
	LastLoginFrom    *time.Time           `json:"last_login_from"`
	LastLoginTo      *time.Time           `json:"last_login_to"`
	Search           string               `json:"search"`
	SortBy           string               `json:"sort_by"`
	SortOrder        string               `json:"sort_order"`
	Limit            int                  `json:"limit"`
	Offset           int                  `json:"offset"`
}

// UserProfileRepository defines the interface for user profile data access
type UserProfileRepository interface {
	// Create creates a new user profile
	Create(ctx context.Context, profile *entities.UserProfile) error

	// GetByUserID retrieves a user profile by user ID
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.UserProfile, error)

	// Update updates an existing user profile
	Update(ctx context.Context, profile *entities.UserProfile) error

	// Delete deletes a user profile by user ID
	Delete(ctx context.Context, userID uuid.UUID) error
}

// UserSessionRepository defines the interface for user session data access
type UserSessionRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, session *entities.UserSession) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.UserSession, error)
	GetByToken(ctx context.Context, token string) (*entities.UserSession, error)
	Update(ctx context.Context, session *entities.UserSession) error
	Delete(ctx context.Context, id uuid.UUID) error

	// User-specific operations
	GetActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.UserSession, error)
	GetSessionsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.UserSession, error)
	InvalidateUserSessions(ctx context.Context, userID uuid.UUID) error
	InvalidateSessionByToken(ctx context.Context, token string) error

	// Cleanup operations
	DeleteExpiredSessions(ctx context.Context) error
	DeleteInactiveSessions(ctx context.Context, inactiveThreshold time.Duration) error

	// Analytics
	CountActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}

// UserLoginHistoryRepository defines the interface for user login history data access
type UserLoginHistoryRepository interface {
	// Basic operations
	Create(ctx context.Context, history *entities.UserLoginHistory) error
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.UserLoginHistory, error)

	// Analytics
	GetFailedLoginAttempts(ctx context.Context, userID uuid.UUID, since time.Time) ([]*entities.UserLoginHistory, error)
	CountLoginAttempts(ctx context.Context, userID uuid.UUID, since time.Time) (int64, error)
	CountFailedAttempts(ctx context.Context, userID uuid.UUID, since time.Time) (int64, error)

	// Cleanup
	DeleteOldHistory(ctx context.Context, olderThan time.Time) error
}

// UserActivityRepository defines the interface for user activity data access
type UserActivityRepository interface {
	// Basic operations
	Create(ctx context.Context, activity *entities.UserActivity) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.UserActivity, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.UserActivity, error)
	GetByUserIDAndType(ctx context.Context, userID uuid.UUID, activityType entities.ActivityType, limit, offset int) ([]*entities.UserActivity, error)
	GetRecentActivity(ctx context.Context, userID uuid.UUID, since time.Time) ([]*entities.UserActivity, error)

	// Analytics
	GetActivityStats(ctx context.Context, userID uuid.UUID, dateFrom, dateTo time.Time) (map[entities.ActivityType]int64, error)
	GetMostActiveUsers(ctx context.Context, limit int, dateFrom, dateTo time.Time) ([]*entities.User, error)

	// Cleanup
	DeleteOldActivities(ctx context.Context, olderThan time.Time) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

// UserPreferencesRepository defines the interface for user preferences data access
type UserPreferencesRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, preferences *entities.UserPreferences) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.UserPreferences, error)
	Update(ctx context.Context, preferences *entities.UserPreferences) error
	Delete(ctx context.Context, userID uuid.UUID) error

	// Specific preference updates
	UpdateTheme(ctx context.Context, userID uuid.UUID, theme string) error
	UpdateLanguage(ctx context.Context, userID uuid.UUID, language string) error
	UpdateCurrency(ctx context.Context, userID uuid.UUID, currency string) error
	UpdateTimezone(ctx context.Context, userID uuid.UUID, timezone string) error
	UpdateNotificationSettings(ctx context.Context, userID uuid.UUID, settings map[string]bool) error
	UpdatePrivacySettings(ctx context.Context, userID uuid.UUID, settings map[string]interface{}) error
	UpdateShoppingSettings(ctx context.Context, userID uuid.UUID, settings map[string]interface{}) error

	// Bulk operations
	GetPreferencesByTheme(ctx context.Context, theme string, limit, offset int) ([]*entities.UserPreferences, error)
	GetPreferencesByLanguage(ctx context.Context, language string, limit, offset int) ([]*entities.UserPreferences, error)
}

// UserVerificationRepository defines the interface for user verification data access
type UserVerificationRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, verification *entities.UserVerification) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.UserVerification, error)
	GetByToken(ctx context.Context, token string) (*entities.UserVerification, error)
	GetByUserIDAndType(ctx context.Context, userID uuid.UUID, verificationType string) (*entities.UserVerification, error)
	Update(ctx context.Context, verification *entities.UserVerification) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Verification operations
	GetActiveVerifications(ctx context.Context, userID uuid.UUID) ([]*entities.UserVerification, error)
	GetByCode(ctx context.Context, code string, verificationType string) (*entities.UserVerification, error)
	MarkAsVerified(ctx context.Context, id uuid.UUID) error
	IncrementAttempt(ctx context.Context, id uuid.UUID) error

	// Cleanup operations
	DeleteExpiredVerifications(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error

	// Analytics
	CountVerificationsByType(ctx context.Context, verificationType string, dateFrom, dateTo time.Time) (int64, error)
	GetFailedVerifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.UserVerification, error)
}
