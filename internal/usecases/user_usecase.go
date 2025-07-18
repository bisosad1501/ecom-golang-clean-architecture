package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/domain/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UserUseCase defines user use cases
type UserUseCase interface {
	Register(ctx context.Context, req RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error)
	ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req ResetPasswordRequest) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req UpdateProfileRequest) (*UserResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req ChangePasswordRequest) error
	GetUsers(ctx context.Context, limit, offset int) (*UsersListResponse, error)
	DeactivateUser(ctx context.Context, userID uuid.UUID) error
	ActivateUser(ctx context.Context, userID uuid.UUID) error

	// Enhanced user methods
	GetUsersWithFilters(ctx context.Context, filters repositories.UserFilters) (*UsersListResponse, error)
	GetUserSessions(ctx context.Context, userID uuid.UUID, limit, offset int) (*UserSessionsResponse, error)
	InvalidateSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error
	InvalidateAllSessions(ctx context.Context, userID uuid.UUID) error
	GetUserActivity(ctx context.Context, userID uuid.UUID, limit, offset int) (*UserActivityResponse, error)
	TrackUserActivity(ctx context.Context, userID uuid.UUID, activityType string, description string, entityType string, entityID *uuid.UUID, metadata map[string]interface{}) error
	GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStatsResponse, error)

	// User preferences methods
	GetUserPreferences(ctx context.Context, userID uuid.UUID) (*UserPreferencesResponse, error)
	UpdateUserPreferences(ctx context.Context, userID uuid.UUID, req UpdateUserPreferencesRequest) (*UserPreferencesResponse, error)
	UpdateTheme(ctx context.Context, userID uuid.UUID, theme string) error

	// Search history
	TrackSearch(ctx context.Context, req TrackSearchRequest) error
	GetSearchHistory(ctx context.Context, userID uuid.UUID, req SearchHistoryRequest) (*UserSearchHistoryListResponse, error)
	ClearSearchHistory(ctx context.Context, userID uuid.UUID) error
	DeleteSearchHistoryItem(ctx context.Context, userID, historyID uuid.UUID) error
	GetPopularSearches(ctx context.Context, userID uuid.UUID, limit int) (*UserPopularSearchesResponse, error)

	// Saved searches
	CreateSavedSearch(ctx context.Context, req CreateSavedSearchRequest) (*SavedSearchResponse, error)
	GetSavedSearches(ctx context.Context, userID uuid.UUID, req GetSavedSearchesRequest) (*GetSavedSearchesResponse, error)
	UpdateSavedSearch(ctx context.Context, req UpdateSavedSearchRequest) (*SavedSearchResponse, error)
	DeleteSavedSearch(ctx context.Context, userID, savedSearchID uuid.UUID) error
	ExecuteSavedSearch(ctx context.Context, userID, savedSearchID uuid.UUID) (*SavedSearchExecutionResponse, error)

	// Browsing history
	TrackProductView(ctx context.Context, req TrackProductViewRequest) error
	GetBrowsingHistory(ctx context.Context, userID uuid.UUID, req BrowsingHistoryRequest) (*BrowsingHistoryResponse, error)
	ClearBrowsingHistory(ctx context.Context, userID uuid.UUID) error

	// Personalization
	GetPersonalization(ctx context.Context, userID uuid.UUID) (*PersonalizationResponse, error)
	UpdatePersonalization(ctx context.Context, req UpdatePersonalizationRequest) (*PersonalizationResponse, error)
	GetPersonalizedRecommendations(ctx context.Context, userID uuid.UUID, req PersonalizedRecommendationsRequest) (*PersonalizedRecommendationsResponse, error)
	AnalyzeUserBehavior(ctx context.Context, userID uuid.UUID) (*UserBehaviorAnalysisResponse, error)

	// Login History
	GetLoginHistory(ctx context.Context, userID uuid.UUID, req LoginHistoryRequest) (*LoginHistoryResponse, error)
	GetLoginStats(ctx context.Context, userID uuid.UUID) (*LoginStatsResponse, error)

	// Profile analytics
	GetProfileAnalytics(ctx context.Context, userID uuid.UUID) (*ProfileAnalyticsResponse, error)
	GetActivitySummary(ctx context.Context, userID uuid.UUID, timeRange string) (*ActivitySummaryResponse, error)
	UpdateLanguage(ctx context.Context, userID uuid.UUID, language string) error

	// User verification methods
	SendEmailVerification(ctx context.Context, userID uuid.UUID) error
	VerifyEmail(ctx context.Context, token string) error
	VerifyEmailByToken(ctx context.Context, token string) (*UserResponse, error)
	GetVerificationStatus(ctx context.Context, userID uuid.UUID) (*VerificationStatusResponse, error)
}

// UserNotificationService interface for user notifications
type UserNotificationService interface {
	NotifyNewUser(ctx context.Context, userID uuid.UUID) error
}

type userUseCase struct {
	userRepo             repositories.UserRepository
	userProfileRepo      repositories.UserProfileRepository
	userSessionRepo      repositories.UserSessionRepository
	userLoginHistoryRepo repositories.UserLoginHistoryRepository
	userActivityRepo     repositories.UserActivityRepository
	userPreferencesRepo  repositories.UserPreferencesRepository
	userVerificationRepo repositories.UserVerificationRepository
	passwordResetRepo    repositories.PasswordResetRepository
	passwordService      services.PasswordService
	gmailService         GmailService
	notificationService  UserNotificationService
	jwtSecret            string
}

// GmailService interface for email operations
type GmailService interface {
	SendVerificationEmail(ctx context.Context, to, firstName, verificationLink string) error
	SendPasswordResetEmail(ctx context.Context, to, firstName, resetLink string) error
	SendWelcomeEmail(ctx context.Context, to, firstName string) error
	ValidateConfiguration() error
}



// NewUserUseCase creates a new user use case
func NewUserUseCase(
	userRepo repositories.UserRepository,
	userProfileRepo repositories.UserProfileRepository,
	userSessionRepo repositories.UserSessionRepository,
	userLoginHistoryRepo repositories.UserLoginHistoryRepository,
	userActivityRepo repositories.UserActivityRepository,
	userPreferencesRepo repositories.UserPreferencesRepository,
	userVerificationRepo repositories.UserVerificationRepository,
	passwordResetRepo repositories.PasswordResetRepository,
	passwordService services.PasswordService,
	gmailService GmailService,
	notificationService UserNotificationService,
	jwtSecret string,
) UserUseCase {
	return &userUseCase{
		userRepo:             userRepo,
		userProfileRepo:      userProfileRepo,
		userSessionRepo:      userSessionRepo,
		userLoginHistoryRepo: userLoginHistoryRepo,
		userActivityRepo:     userActivityRepo,
		userPreferencesRepo:  userPreferencesRepo,
		userVerificationRepo: userVerificationRepo,
		passwordResetRepo:    passwordResetRepo,
		passwordService:      passwordService,
		gmailService:         gmailService,
		notificationService:  notificationService,
		jwtSecret:            jwtSecret,
	}
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Phone     string `json:"phone" validate:"omitempty,min=10,max=15"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	IPAddress  string `json:"ip_address,omitempty"`  // Client IP address
	UserAgent  string `json:"user_agent,omitempty"`  // Browser/client user agent
	DeviceInfo string `json:"device_info,omitempty"` // Device information
}

// ForgotPasswordRequest represents forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// VerifyEmailRequest represents email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// ResendVerificationRequest represents resend verification request
type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// RefreshTokenResponse represents refresh token response
type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// UpdateProfileRequest represents update profile request
type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

// ChangePasswordRequest represents change password request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

// UserResponse represents user response (matches frontend User interface)
type UserResponse struct {
	ID        uuid.UUID            `json:"id"`
	Email     string               `json:"email"`
	FirstName string               `json:"first_name"`
	LastName  string               `json:"last_name"`
	Phone     string               `json:"phone"`
	Role      entities.UserRole    `json:"role"`
	IsActive  bool                 `json:"is_active"`

	// User metrics (from User entity - matches frontend expectations)
	TotalOrders    int     `json:"total_orders"`
	TotalSpent     float64 `json:"total_spent"`
	LoyaltyPoints  int     `json:"loyalty_points"`
	MembershipTier string  `json:"membership_tier"`

	Profile   *UserProfileResponse `json:"profile,omitempty"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

// UsersListResponse represents paginated users response
type UsersListResponse struct {
	Users []*UserResponse `json:"users"`
	Total int64           `json:"total"`
}

// UserSessionsResponse represents user sessions response
type UserSessionsResponse struct {
	Sessions   []*UserSessionResponse `json:"sessions"`
	Total      int64                  `json:"total"`
	Pagination PaginationInfo         `json:"pagination"`
}

// UserSessionResponse represents user session response
type UserSessionResponse struct {
	ID           uuid.UUID `json:"id"`
	DeviceInfo   string    `json:"device_info"`
	IPAddress    string    `json:"ip_address"`
	Location     string    `json:"location"`
	IsActive     bool      `json:"is_active"`
	IsCurrent    bool      `json:"is_current"`
	LastActivity time.Time `json:"last_activity"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserActivityResponse represents user activity response
type UserActivityResponse struct {
	Activities []*UserActivityItem `json:"activities"`
	Total      int64               `json:"total"`
	Pagination PaginationInfo      `json:"pagination"`
}

// UserActivityItem represents user activity item
type UserActivityItem struct {
	ID          uuid.UUID             `json:"id"`
	Type        entities.ActivityType `json:"type"`
	Description string                `json:"description"`
	EntityType  string                `json:"entity_type"`
	EntityID    *uuid.UUID            `json:"entity_id"`
	IPAddress   string                `json:"ip_address"`
	CreatedAt   time.Time             `json:"created_at"`
}

// UserStatsResponse represents user statistics response
type UserStatsResponse struct {
	TotalSessions     int                             `json:"total_sessions"`
	TotalPageViews    int                             `json:"total_page_views"`
	TotalProductViews int                             `json:"total_product_views"`
	TotalSearches     int                             `json:"total_searches"`
	TotalOrders       int                             `json:"total_orders"`
	TotalSpent        float64                         `json:"total_spent"`
	AverageOrderValue float64                         `json:"average_order_value"`
	ActivityBreakdown map[entities.ActivityType]int64 `json:"activity_breakdown"`
	RecentActivities  []*UserActivityItem             `json:"recent_activities"`
}

// UserPreferencesResponse represents user preferences response
type UserPreferencesResponse struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`

	// Display preferences
	Theme    string `json:"theme"`
	Language string `json:"language"`
	Currency string `json:"currency"`
	Timezone string `json:"timezone"`

	// Notification preferences
	EmailNotifications     bool `json:"email_notifications"`
	SMSNotifications       bool `json:"sms_notifications"`
	PushNotifications      bool `json:"push_notifications"`
	MarketingEmails        bool `json:"marketing_emails"`
	OrderUpdates           bool `json:"order_updates"`
	ProductRecommendations bool `json:"product_recommendations"`
	NewsletterSubscription bool `json:"newsletter_subscription"`
	SecurityAlerts         bool `json:"security_alerts"`

	// Privacy preferences
	ProfileVisibility      string `json:"profile_visibility"`
	ShowOnlineStatus       bool   `json:"show_online_status"`
	AllowDataCollection    bool   `json:"allow_data_collection"`
	AllowPersonalization   bool   `json:"allow_personalization"`
	AllowThirdPartySharing bool   `json:"allow_third_party_sharing"`

	// Shopping preferences
	DefaultShippingMethod string `json:"default_shipping_method"`
	DefaultPaymentMethod  string `json:"default_payment_method"`
	SavePaymentMethods    bool   `json:"save_payment_methods"`
	AutoApplyCoupons      bool   `json:"auto_apply_coupons"`
	WishlistVisibility    string `json:"wishlist_visibility"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateUserPreferencesRequest represents update user preferences request
type UpdateUserPreferencesRequest struct {
	// Display preferences
	Theme    *string `json:"theme,omitempty"`
	Language *string `json:"language,omitempty"`
	Currency *string `json:"currency,omitempty"`
	Timezone *string `json:"timezone,omitempty"`

	// Notification preferences
	EmailNotifications     *bool `json:"email_notifications,omitempty"`
	SMSNotifications       *bool `json:"sms_notifications,omitempty"`
	PushNotifications      *bool `json:"push_notifications,omitempty"`
	MarketingEmails        *bool `json:"marketing_emails,omitempty"`
	OrderUpdates           *bool `json:"order_updates,omitempty"`
	ProductRecommendations *bool `json:"product_recommendations,omitempty"`
	NewsletterSubscription *bool `json:"newsletter_subscription,omitempty"`
	SecurityAlerts         *bool `json:"security_alerts,omitempty"`

	// Privacy preferences
	ProfileVisibility      *string `json:"profile_visibility,omitempty"`
	ShowOnlineStatus       *bool   `json:"show_online_status,omitempty"`
	AllowDataCollection    *bool   `json:"allow_data_collection,omitempty"`
	AllowPersonalization   *bool   `json:"allow_personalization,omitempty"`
	AllowThirdPartySharing *bool   `json:"allow_third_party_sharing,omitempty"`

	// Shopping preferences
	DefaultShippingMethod *string `json:"default_shipping_method,omitempty"`
	DefaultPaymentMethod  *string `json:"default_payment_method,omitempty"`
	SavePaymentMethods    *bool   `json:"save_payment_methods,omitempty"`
	AutoApplyCoupons      *bool   `json:"auto_apply_coupons,omitempty"`
	WishlistVisibility    *string `json:"wishlist_visibility,omitempty"`
}

// VerificationStatusResponse represents verification status response
type VerificationStatusResponse struct {
	UserID        uuid.UUID `json:"user_id"`
	EmailVerified bool      `json:"email_verified"`
	PhoneVerified bool      `json:"phone_verified"`

	// Active verifications
	PendingEmailVerification bool `json:"pending_email_verification"`
	PendingPhoneVerification bool `json:"pending_phone_verification"`

	// Verification history
	LastEmailVerificationSent *time.Time `json:"last_email_verification_sent"`
	LastPhoneVerificationSent *time.Time `json:"last_phone_verification_sent"`
	EmailVerifiedAt           *time.Time `json:"email_verified_at"`
	PhoneVerifiedAt           *time.Time `json:"phone_verified_at"`
}

// UserProfileResponse represents user profile response
type UserProfileResponse struct {
	Avatar      string     `json:"avatar"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Gender      string     `json:"gender"`
	Address     string     `json:"address"`
	City        string     `json:"city"`
	State       string     `json:"state"`
	Country     string     `json:"country"`
	ZipCode     string     `json:"zip_code"`
}

// LoginResponse represents login response
type LoginResponse struct {
	User         *UserResponse `json:"user"`
	Token        string        `json:"token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresAt    int64         `json:"expires_at"`
}

// Register registers a new user
func (uc *userUseCase) Register(ctx context.Context, req RegisterRequest) (*UserResponse, error) {
	// Validate password complexity
	if err := uc.validatePasswordComplexity(req.Password); err != nil {
		return nil, err
	}

	// Validate email format more strictly
	if err := uc.validateEmailFormat(req.Email); err != nil {
		return nil, err
	}

	// Validate phone format if provided
	if req.Phone != "" {
		if err := uc.validatePhoneFormat(req.Phone); err != nil {
			return nil, err
		}
	}

	// Check if user already exists
	exists, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entities.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := uc.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entities.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      entities.UserRoleCustomer,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Send email verification automatically after registration
	go func() {
		if err := uc.SendEmailVerification(context.Background(), user.ID); err != nil {
			fmt.Printf("❌ Failed to send email verification to %s: %v\n", user.Email, err)
		}
	}()

	// Send new user notification to admin (async)
	if uc.notificationService != nil {
		go func() {
			if err := uc.notificationService.NotifyNewUser(context.Background(), user.ID); err != nil {
				fmt.Printf("❌ Failed to send new user notification: %v\n", err)
			} else {
				fmt.Printf("✅ New user notification sent to admin\n")
			}
		}()
	}

	return uc.toUserResponse(user), nil
}

// validatePasswordComplexity validates password complexity requirements
func (uc *userUseCase) validatePasswordComplexity(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return fmt.Errorf("password must be less than 128 characters long")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char >= 32 && char <= 126: // Printable ASCII special characters
			if !((char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')) {
				hasSpecial = true
			}
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// validateEmailFormat validates email format more strictly
func (uc *userUseCase) validateEmailFormat(email string) error {
	if len(email) == 0 {
		return fmt.Errorf("email is required")
	}

	if len(email) > 254 {
		return fmt.Errorf("email is too long")
	}

	// Basic email format validation
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email format")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email format")
	}

	localPart := parts[0]
	domainPart := parts[1]

	if len(localPart) == 0 || len(localPart) > 64 {
		return fmt.Errorf("invalid email local part")
	}

	if len(domainPart) == 0 || len(domainPart) > 253 {
		return fmt.Errorf("invalid email domain part")
	}

	if !strings.Contains(domainPart, ".") {
		return fmt.Errorf("invalid email domain format")
	}

	return nil
}

// validatePhoneFormat validates phone number format
func (uc *userUseCase) validatePhoneFormat(phone string) error {
	if len(phone) == 0 {
		return nil // Phone is optional
	}

	// Remove common phone number separators
	cleanPhone := strings.ReplaceAll(phone, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "(", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, ")", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "+", "")

	if len(cleanPhone) < 10 || len(cleanPhone) > 15 {
		return fmt.Errorf("phone number must be between 10 and 15 digits")
	}

	// Check if all characters are digits
	for _, char := range cleanPhone {
		if char < '0' || char > '9' {
			return fmt.Errorf("phone number must contain only digits and common separators")
		}
	}

	return nil
}

// checkLoginRateLimit checks if user has exceeded login attempt limits
func (uc *userUseCase) checkLoginRateLimit(ctx context.Context, email string) error {
	// Simple in-memory rate limiting (in production, use Redis or database)
	// For now, we'll implement a basic check
	// TODO: Implement proper rate limiting with Redis/database storage

	// This is a simplified implementation
	// In production, you would store failed attempts in cache/database
	return nil
}

// incrementFailedLoginAttempts increments failed login attempts for an email
func (uc *userUseCase) incrementFailedLoginAttempts(ctx context.Context, email string) error {
	// TODO: Implement proper failed attempt tracking
	// This would typically store in Redis or database with expiration
	return nil
}

// resetFailedLoginAttempts resets failed login attempts for an email
func (uc *userUseCase) resetFailedLoginAttempts(ctx context.Context, email string) error {
	// TODO: Implement proper failed attempt reset
	// This would typically clear the counter in Redis or database
	return nil
}

// Login authenticates a user
func (uc *userUseCase) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Check rate limiting for this email
	if err := uc.checkLoginRateLimit(ctx, req.Email); err != nil {
		return nil, err
	}

	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Log failed login attempt
		_ = uc.logLoginAttemptEnhanced(ctx, req.Email, false, "user not found", req.IPAddress, req.UserAgent, req.DeviceInfo)
		_ = uc.incrementFailedLoginAttempts(ctx, req.Email)
		return nil, entities.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		// Log failed login attempt
		_ = uc.logLoginAttemptEnhanced(ctx, req.Email, false, "user not active", req.IPAddress, req.UserAgent, req.DeviceInfo)
		_ = uc.incrementFailedLoginAttempts(ctx, req.Email)
		return nil, entities.ErrUserNotActive
	}

	// Check if email is verified
	if !user.EmailVerified {
		// Log failed login attempt
		_ = uc.logLoginAttemptEnhanced(ctx, req.Email, false, "email not verified", req.IPAddress, req.UserAgent, req.DeviceInfo)
		_ = uc.incrementFailedLoginAttempts(ctx, req.Email)
		return nil, fmt.Errorf("email not verified. Please check your email and verify your account before logging in")
	}

	// Check password
	if err := uc.passwordService.CheckPassword(req.Password, user.Password); err != nil {
		// Log failed login attempt
		_ = uc.logLoginAttemptEnhanced(ctx, req.Email, false, "invalid password", req.IPAddress, req.UserAgent, req.DeviceInfo)
		_ = uc.incrementFailedLoginAttempts(ctx, req.Email)
		return nil, entities.ErrInvalidCredentials
	}

	// Reset failed login attempts on successful login
	_ = uc.resetFailedLoginAttempts(ctx, req.Email)

	// Generate JWT token
	token, err := uc.generateJWTToken(user)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := uc.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Create user session with enhanced tracking
	session := &entities.UserSession{
		ID:           uuid.New(),
		UserID:       user.ID,
		SessionToken: token,
		DeviceInfo:   req.DeviceInfo,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
		Location:     uc.getLocationFromIP(req.IPAddress), // TODO: Implement IP geolocation
		IsActive:     true,
		LastActivity: time.Now(),
		ExpiresAt:    time.Now().Add(time.Hour * 24),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save session
	if err := uc.userSessionRepo.Create(ctx, session); err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to create user session: %v\n", err)
	}

	// Update user last login
	now := time.Now()
	user.LastLoginAt = &now
	user.LastActivityAt = &now
	user.UpdatedAt = now
	_ = uc.userRepo.Update(ctx, user)

	// Log successful login attempt with enhanced tracking
	_ = uc.logLoginAttemptEnhanced(ctx, req.Email, true, "", req.IPAddress, req.UserAgent, req.DeviceInfo)

	return &LoginResponse{
		User:         uc.toUserResponse(user),
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Hour * 24).Unix(),
	}, nil
}

// GetProfile gets user profile
func (uc *userUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, entities.ErrUserNotFound
	}

	userResponse := uc.toUserResponse(user)

	// Get user profile
	profile, err := uc.userProfileRepo.GetByUserID(ctx, userID)
	if err == nil {
		userResponse.Profile = uc.toUserProfileResponse(profile)
	}

	return userResponse, nil
}

// UpdateProfile updates user profile
func (uc *userUseCase) UpdateProfile(ctx context.Context, userID uuid.UUID, req UpdateProfileRequest) (*UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, entities.ErrUserNotFound
	}

	// Update user fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return uc.toUserResponse(user), nil
}

// ChangePassword changes user password
func (uc *userUseCase) ChangePassword(ctx context.Context, userID uuid.UUID, req ChangePasswordRequest) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return entities.ErrUserNotFound
	}

	// Check current password
	if err := uc.passwordService.CheckPassword(req.CurrentPassword, user.Password); err != nil {
		return entities.ErrInvalidCredentials
	}

	// Validate new password complexity
	if err := uc.validatePasswordComplexity(req.NewPassword); err != nil {
		return err
	}

	// Check if new password is different from current password
	if err := uc.passwordService.CheckPassword(req.NewPassword, user.Password); err == nil {
		return fmt.Errorf("new password must be different from current password")
	}

	// Hash new password
	hashedPassword, err := uc.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return uc.userRepo.UpdatePassword(ctx, userID, hashedPassword)
}

// GetUsers gets list of users
func (uc *userUseCase) GetUsers(ctx context.Context, limit, offset int) (*UsersListResponse, error) {
	users, err := uc.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = uc.toUserResponse(user)
	}

	total, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &UsersListResponse{
		Users: responses,
		Total: total,
	}, nil
}

// DeactivateUser deactivates a user
func (uc *userUseCase) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	return uc.userRepo.SetActive(ctx, userID, false)
}

// ActivateUser activates a user
func (uc *userUseCase) ActivateUser(ctx context.Context, userID uuid.UUID) error {
	return uc.userRepo.SetActive(ctx, userID, true)
}

// generateJWTToken generates a JWT token for the user
func (uc *userUseCase) generateJWTToken(user *entities.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"iat":     now.Unix(),
		"exp":     now.Add(time.Hour * 24 * 7).Unix(), // 7 days for better UX
		"jti":     uuid.New().String(),                // JWT ID for token tracking
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

// toUserResponse converts user entity to response (includes user metrics)
func (uc *userUseCase) toUserResponse(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Role:      user.Role,
		IsActive:  user.IsActive,

		// User metrics (matches frontend expectations)
		TotalOrders:    user.TotalOrders,
		TotalSpent:     user.TotalSpent,
		LoyaltyPoints:  user.LoyaltyPoints,
		MembershipTier: user.MembershipTier,

		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// toUserProfileResponse converts user profile entity to response
func (uc *userUseCase) toUserProfileResponse(profile *entities.UserProfile) *UserProfileResponse {
	return &UserProfileResponse{
		Avatar:      profile.Avatar,
		DateOfBirth: profile.DateOfBirth,
		Gender:      profile.Gender,
		Address:     profile.Address,
		City:        profile.City,
		State:       profile.State,
		Country:     profile.Country,
		ZipCode:     profile.ZipCode,
	}
}

// GetUsersWithFilters gets users with filters
func (uc *userUseCase) GetUsersWithFilters(ctx context.Context, filters repositories.UserFilters) (*UsersListResponse, error) {
	users, err := uc.userRepo.GetUsersWithFilters(ctx, filters)
	if err != nil {
		return nil, err
	}

	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = uc.toUserResponse(user)
	}

	total, err := uc.userRepo.CountUsersWithFilters(ctx, filters)
	if err != nil {
		return nil, err
	}

	return &UsersListResponse{
		Users: responses,
		Total: total,
	}, nil
}

// GetUserSessions gets user sessions
func (uc *userUseCase) GetUserSessions(ctx context.Context, userID uuid.UUID, limit, offset int) (*UserSessionsResponse, error) {
	sessions, err := uc.userSessionRepo.GetSessionsByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	sessionResponses := make([]*UserSessionResponse, len(sessions))
	for i, session := range sessions {
		sessionResponses[i] = &UserSessionResponse{
			ID:           session.ID,
			DeviceInfo:   session.DeviceInfo,
			IPAddress:    session.IPAddress,
			Location:     session.Location,
			IsActive:     session.IsActive,
			IsCurrent:    false, // TODO: Determine current session
			LastActivity: session.LastActivity,
			CreatedAt:    session.CreatedAt,
		}
	}

	// Get total count (simplified)
	total := int64(len(sessions))

	// Create pagination info using the standardized function
	pagination := NewPaginationInfoFromOffset(offset, limit, total)

	return &UserSessionsResponse{
		Sessions:   sessionResponses,
		Total:      total,
		Pagination: *pagination,
	}, nil
}

// InvalidateSession invalidates a specific session
func (uc *userUseCase) InvalidateSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	session, err := uc.userSessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}

	if session.UserID != userID {
		return entities.ErrUserNotFound
	}

	session.IsActive = false
	return uc.userSessionRepo.Update(ctx, session)
}

// InvalidateAllSessions invalidates all user sessions
func (uc *userUseCase) InvalidateAllSessions(ctx context.Context, userID uuid.UUID) error {
	return uc.userSessionRepo.InvalidateUserSessions(ctx, userID)
}

// GetUserActivity gets user activity
func (uc *userUseCase) GetUserActivity(ctx context.Context, userID uuid.UUID, limit, offset int) (*UserActivityResponse, error) {
	activities, err := uc.userActivityRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	activityItems := make([]*UserActivityItem, len(activities))
	for i, activity := range activities {
		activityItems[i] = &UserActivityItem{
			ID:          activity.ID,
			Type:        activity.Type,
			Description: activity.Description,
			EntityType:  activity.EntityType,
			EntityID:    activity.EntityID,
			IPAddress:   activity.IPAddress,
			CreatedAt:   activity.CreatedAt,
		}
	}

	// Get total count (simplified)
	total := int64(len(activities))

	// Create pagination info using standardized function
	pagination := NewPaginationInfoFromOffset(offset, limit, total)

	return &UserActivityResponse{
		Activities: activityItems,
		Total:      total,
		Pagination: *pagination,
	}, nil
}

// TrackUserActivity tracks user activity
func (uc *userUseCase) TrackUserActivity(ctx context.Context, userID uuid.UUID, activityType string, description string, entityType string, entityID *uuid.UUID, metadata map[string]interface{}) error {
	activity := &entities.UserActivity{
		ID:          uuid.New(),
		UserID:      userID,
		Type:        entities.ActivityType(activityType),
		Description: description,
		EntityType:  entityType,
		EntityID:    entityID,
		CreatedAt:   time.Now(),
	}

	// Convert metadata to JSON string if provided
	if metadata != nil {
		// In a real implementation, you'd marshal this to JSON
		// For now, we'll leave it empty
		activity.Metadata = ""
	}

	return uc.userActivityRepo.Create(ctx, activity)
}

// GetUserStats gets user statistics
func (uc *userUseCase) GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStatsResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, entities.ErrUserNotFound
	}

	// Get activity breakdown for the last 30 days
	dateFrom := time.Now().AddDate(0, 0, -30)
	dateTo := time.Now()

	activityBreakdown, err := uc.userActivityRepo.GetActivityStats(ctx, userID, dateFrom, dateTo)
	if err != nil {
		activityBreakdown = make(map[entities.ActivityType]int64)
	}

	// Get recent activities
	recentActivities, err := uc.userActivityRepo.GetRecentActivity(ctx, userID, time.Now().AddDate(0, 0, -7))
	if err != nil {
		recentActivities = []*entities.UserActivity{}
	}

	recentActivityItems := make([]*UserActivityItem, len(recentActivities))
	for i, activity := range recentActivities {
		recentActivityItems[i] = &UserActivityItem{
			ID:          activity.ID,
			Type:        activity.Type,
			Description: activity.Description,
			EntityType:  activity.EntityType,
			EntityID:    activity.EntityID,
			IPAddress:   activity.IPAddress,
			CreatedAt:   activity.CreatedAt,
		}
	}

	// Calculate average order value
	averageOrderValue := float64(0)
	if user.TotalOrders > 0 {
		averageOrderValue = user.TotalSpent / float64(user.TotalOrders)
	}

	return &UserStatsResponse{
		TotalSessions:     0, // TODO: Implement session counting
		TotalPageViews:    0, // TODO: Implement page view counting
		TotalProductViews: 0, // TODO: Implement product view counting
		TotalSearches:     0, // TODO: Implement search counting
		TotalOrders:       user.TotalOrders,
		TotalSpent:        user.TotalSpent,
		AverageOrderValue: averageOrderValue,
		ActivityBreakdown: activityBreakdown,
		RecentActivities:  recentActivityItems,
	}, nil
}

// GetUserPreferences gets user preferences
func (uc *userUseCase) GetUserPreferences(ctx context.Context, userID uuid.UUID) (*UserPreferencesResponse, error) {
	preferences, err := uc.userPreferencesRepo.GetByUserID(ctx, userID)
	if err != nil {
		// If preferences don't exist, create default ones
		if err == entities.ErrUserNotFound {
			defaultPreferences := &entities.UserPreferences{
				ID:                 uuid.New(),
				UserID:             userID,
				Theme:              "system",
				Language:           "en",
				Currency:           "USD",
				Timezone:           "UTC",
				EmailNotifications: true,
				SecurityAlerts:     true,
				ProfileVisibility:  "private",
				SavePaymentMethods: true,
				AutoApplyCoupons:   true,
				WishlistVisibility: "private",
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			}

			if createErr := uc.userPreferencesRepo.Create(ctx, defaultPreferences); createErr != nil {
				return nil, createErr
			}
			preferences = defaultPreferences
		} else {
			return nil, err
		}
	}

	return uc.toUserPreferencesResponse(preferences), nil
}

// UpdateUserPreferences updates user preferences
func (uc *userUseCase) UpdateUserPreferences(ctx context.Context, userID uuid.UUID, req UpdateUserPreferencesRequest) (*UserPreferencesResponse, error) {
	preferences, err := uc.userPreferencesRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Theme != nil {
		preferences.Theme = *req.Theme
	}
	if req.Language != nil {
		preferences.Language = *req.Language
	}
	if req.Currency != nil {
		preferences.Currency = *req.Currency
	}
	if req.Timezone != nil {
		preferences.Timezone = *req.Timezone
	}
	if req.EmailNotifications != nil {
		preferences.EmailNotifications = *req.EmailNotifications
	}
	if req.SMSNotifications != nil {
		preferences.SMSNotifications = *req.SMSNotifications
	}
	if req.PushNotifications != nil {
		preferences.PushNotifications = *req.PushNotifications
	}
	if req.MarketingEmails != nil {
		preferences.MarketingEmails = *req.MarketingEmails
	}
	if req.OrderUpdates != nil {
		preferences.OrderUpdates = *req.OrderUpdates
	}
	if req.ProductRecommendations != nil {
		preferences.ProductRecommendations = *req.ProductRecommendations
	}
	if req.NewsletterSubscription != nil {
		preferences.NewsletterSubscription = *req.NewsletterSubscription
	}
	if req.SecurityAlerts != nil {
		preferences.SecurityAlerts = *req.SecurityAlerts
	}
	if req.ProfileVisibility != nil {
		preferences.ProfileVisibility = *req.ProfileVisibility
	}
	if req.ShowOnlineStatus != nil {
		preferences.ShowOnlineStatus = *req.ShowOnlineStatus
	}
	if req.AllowDataCollection != nil {
		preferences.AllowDataCollection = *req.AllowDataCollection
	}
	if req.AllowPersonalization != nil {
		preferences.AllowPersonalization = *req.AllowPersonalization
	}
	if req.AllowThirdPartySharing != nil {
		preferences.AllowThirdPartySharing = *req.AllowThirdPartySharing
	}
	if req.DefaultShippingMethod != nil {
		preferences.DefaultShippingMethod = *req.DefaultShippingMethod
	}
	if req.DefaultPaymentMethod != nil {
		preferences.DefaultPaymentMethod = *req.DefaultPaymentMethod
	}
	if req.SavePaymentMethods != nil {
		preferences.SavePaymentMethods = *req.SavePaymentMethods
	}
	if req.AutoApplyCoupons != nil {
		preferences.AutoApplyCoupons = *req.AutoApplyCoupons
	}
	if req.WishlistVisibility != nil {
		preferences.WishlistVisibility = *req.WishlistVisibility
	}

	preferences.UpdatedAt = time.Now()

	if err := uc.userPreferencesRepo.Update(ctx, preferences); err != nil {
		return nil, err
	}

	// Track activity
	_ = uc.TrackUserActivity(ctx, userID, "profile_update", "User preferences updated", "user_preferences", &preferences.ID, nil)

	return uc.toUserPreferencesResponse(preferences), nil
}

// UpdateTheme updates user theme preference
func (uc *userUseCase) UpdateTheme(ctx context.Context, userID uuid.UUID, theme string) error {
	return uc.userPreferencesRepo.UpdateTheme(ctx, userID, theme)
}

// UpdateLanguage updates user language preference
func (uc *userUseCase) UpdateLanguage(ctx context.Context, userID uuid.UUID, language string) error {
	return uc.userPreferencesRepo.UpdateLanguage(ctx, userID, language)
}

// toUserPreferencesResponse converts user preferences entity to response
func (uc *userUseCase) toUserPreferencesResponse(preferences *entities.UserPreferences) *UserPreferencesResponse {
	return &UserPreferencesResponse{
		ID:                     preferences.ID,
		UserID:                 preferences.UserID,
		Theme:                  preferences.Theme,
		Language:               preferences.Language,
		Currency:               preferences.Currency,
		Timezone:               preferences.Timezone,
		EmailNotifications:     preferences.EmailNotifications,
		SMSNotifications:       preferences.SMSNotifications,
		PushNotifications:      preferences.PushNotifications,
		MarketingEmails:        preferences.MarketingEmails,
		OrderUpdates:           preferences.OrderUpdates,
		ProductRecommendations: preferences.ProductRecommendations,
		NewsletterSubscription: preferences.NewsletterSubscription,
		SecurityAlerts:         preferences.SecurityAlerts,
		ProfileVisibility:      preferences.ProfileVisibility,
		ShowOnlineStatus:       preferences.ShowOnlineStatus,
		AllowDataCollection:    preferences.AllowDataCollection,
		AllowPersonalization:   preferences.AllowPersonalization,
		AllowThirdPartySharing: preferences.AllowThirdPartySharing,
		DefaultShippingMethod:  preferences.DefaultShippingMethod,
		DefaultPaymentMethod:   preferences.DefaultPaymentMethod,
		SavePaymentMethods:     preferences.SavePaymentMethods,
		AutoApplyCoupons:       preferences.AutoApplyCoupons,
		WishlistVisibility:     preferences.WishlistVisibility,
		CreatedAt:              preferences.CreatedAt,
		UpdatedAt:              preferences.UpdatedAt,
	}
}

// SendEmailVerification sends email verification
func (uc *userUseCase) SendEmailVerification(ctx context.Context, userID uuid.UUID) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return entities.ErrUserNotFound
	}

	if user.EmailVerified {
		return fmt.Errorf("email already verified")
	}

	// Generate verification token
	token := uuid.New().String()

	// Check if verification record already exists for this user (any type)
	existingVerification, err := uc.userVerificationRepo.GetByUserID(ctx, userID)
	if err != nil && err != entities.ErrUserNotFound {
		return fmt.Errorf("failed to check existing verification: %w", err)
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	if existingVerification != nil {
		// Update existing verification record for email verification
		existingVerification.VerificationCode = token
		existingVerification.CodeExpiresAt = &expiresAt
		existingVerification.VerificationType = "email"
		existingVerification.IsUsed = false
		existingVerification.VerifiedAt = nil
		existingVerification.UpdatedAt = time.Now()

		if err := uc.userVerificationRepo.Update(ctx, existingVerification); err != nil {
			return fmt.Errorf("failed to update verification record: %w", err)
		}
	} else {
		// Create new verification record
		verification := &entities.UserVerification{
			ID:               uuid.New(),
			UserID:           userID,
			VerificationType: "email",
			VerificationCode: token,
			CodeExpiresAt:    &expiresAt, // 24 hours expiry
			IsUsed:           false,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		if err := uc.userVerificationRepo.Create(ctx, verification); err != nil {
			return fmt.Errorf("failed to create verification record: %w", err)
		}
	}

	// Send email verification via Gmail service
	verificationLink := fmt.Sprintf("http://localhost:8080/api/v1/auth/verify-email?token=%s", token)

	if err := uc.gmailService.SendVerificationEmail(ctx, user.Email, user.FirstName, verificationLink); err != nil {
		// Log error but don't fail the operation - fallback to console logging
		fmt.Printf("⚠️ Failed to send email verification to %s: %v\n", user.Email, err)
		fmt.Printf("📧 FALLBACK - Email verification token for %s: %s\n", user.Email, token)
		fmt.Printf("📧 FALLBACK - Verification link: %s\n", verificationLink)
	} else {
		fmt.Printf("✅ Email verification sent to %s\n", user.Email)
	}

	// Track activity
	_ = uc.TrackUserActivity(ctx, userID, "profile_update", "Email verification sent", "user", &user.ID, nil)

	return nil
}

// VerifyEmail verifies email with token
func (uc *userUseCase) VerifyEmail(ctx context.Context, token string) error {
	if token == "" {
		return entities.ErrInvalidVerificationCode
	}

	// Find verification record by token
	verification, err := uc.userVerificationRepo.GetByCode(ctx, token, "email")
	if err != nil {
		return entities.ErrAccountVerificationNotFound
	}

	// Check if verification is expired
	if verification.IsExpired() {
		return entities.ErrVerificationCodeExpired
	}

	// Check if verification is already used
	if verification.IsUsed {
		return fmt.Errorf("verification code already used")
	}

	// Mark verification as used
	verification.IsUsed = true
	verifiedAt := time.Now()
	verification.VerifiedAt = &verifiedAt
	verification.UpdatedAt = time.Now()

	if err := uc.userVerificationRepo.Update(ctx, verification); err != nil {
		return fmt.Errorf("failed to update verification record: %w", err)
	}

	// Mark user email as verified
	user, err := uc.userRepo.GetByID(ctx, verification.UserID)
	if err != nil {
		return entities.ErrUserNotFound
	}

	user.EmailVerified = true
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Track activity
	_ = uc.TrackUserActivity(ctx, user.ID, "profile_update", "Email verified", "user", &user.ID, nil)

	fmt.Printf("Email verification successful for user: %s\n", user.Email)

	return nil
}

// VerifyEmailByToken verifies email using verification token from link
func (uc *userUseCase) VerifyEmailByToken(ctx context.Context, token string) (*UserResponse, error) {
	if token == "" {
		return nil, fmt.Errorf("verification token is required")
	}

	// Find verification record by token
	verification, err := uc.userVerificationRepo.GetByCode(ctx, token, "email")
	if err != nil {
		return nil, fmt.Errorf("invalid or expired verification token")
	}

	// Check if token is expired
	if verification.CodeExpiresAt != nil && time.Now().After(*verification.CodeExpiresAt) {
		return nil, fmt.Errorf("verification token has expired")
	}

	// Check if token is already used
	if verification.IsUsed {
		return nil, fmt.Errorf("verification token has already been used")
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, verification.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if email is already verified
	if user.EmailVerified {
		return nil, fmt.Errorf("email is already verified")
	}

	// Mark verification as used
	now := time.Now()
	verification.IsUsed = true
	verification.VerifiedAt = &now
	verification.UpdatedAt = now

	if err := uc.userVerificationRepo.Update(ctx, verification); err != nil {
		return nil, fmt.Errorf("failed to update verification record: %w", err)
	}

	// Mark user email as verified
	user.EmailVerified = true
	user.UpdatedAt = now

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Track activity
	_ = uc.TrackUserActivity(ctx, user.ID, "profile_update", "Email verified", "user", &user.ID, nil)

	fmt.Printf("✅ Email verification successful for user: %s (%s)\n", user.Email, user.ID)

	return uc.toUserResponse(user), nil
}

// GetVerificationStatus gets verification status
func (uc *userUseCase) GetVerificationStatus(ctx context.Context, userID uuid.UUID) (*VerificationStatusResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, entities.ErrUserNotFound
	}

	response := &VerificationStatusResponse{
		UserID:        userID,
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
		// For now, we'll set pending verifications to false
		// In production, you would check against verification table
		PendingEmailVerification: false,
		PendingPhoneVerification: false,
	}

	return response, nil
}

// Logout invalidates a user token
func (uc *userUseCase) Logout(ctx context.Context, token string) error {
	// TODO: Implement token blacklisting
	// For now, we'll just return success since JWT tokens are stateless
	// In production, you should store blacklisted tokens in Redis or database
	return nil
}

// RefreshToken generates a new access token using refresh token
func (uc *userUseCase) RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	// Parse and validate refresh token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(uc.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Extract user ID from claims
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	// Get user to ensure they still exist and are active
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, entities.ErrUserNotFound
	}

	if !user.IsActive {
		return nil, entities.ErrUserNotActive
	}

	// Generate new tokens
	newToken, err := uc.generateJWTToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := uc.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(time.Hour * 24).Unix(),
	}, nil
}

// ForgotPassword initiates password reset process
func (uc *userUseCase) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	// Check if user exists
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		return nil
	}

	// Generate password reset token
	resetToken := uuid.New().String()

	// Create password reset record
	expiresAt := time.Now().Add(1 * time.Hour) // 1 hour expiry
	passwordReset := &entities.PasswordReset{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     resetToken,
		ExpiresAt: expiresAt,
		UsedAt:    nil,
		CreatedAt: time.Now(),
	}

	if err := uc.passwordResetRepo.Create(ctx, passwordReset); err != nil {
		return fmt.Errorf("failed to create password reset record: %w", err)
	}

	// Send password reset email via Gmail service
	resetLink := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", resetToken)

	if err := uc.gmailService.SendPasswordResetEmail(ctx, user.Email, user.FirstName, resetLink); err != nil {
		// Log error but don't fail the operation - fallback to console logging
		fmt.Printf("⚠️ Failed to send password reset email to %s: %v\n", user.Email, err)
		fmt.Printf("📧 FALLBACK - Password reset token for %s: %s\n", user.Email, resetToken)
		fmt.Printf("📧 FALLBACK - Reset link: %s\n", resetLink)
	} else {
		fmt.Printf("✅ Password reset email sent to %s\n", user.Email)
	}

	return nil
}

// ResetPassword resets user password using reset token
func (uc *userUseCase) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	if req.Token == "" {
		return fmt.Errorf("invalid reset token")
	}

	// Get password reset record by token
	passwordReset, err := uc.passwordResetRepo.GetByToken(ctx, req.Token)
	if err != nil {
		return fmt.Errorf("invalid reset token")
	}

	// Check if token is expired
	if time.Now().After(passwordReset.ExpiresAt) {
		return fmt.Errorf("reset token has expired")
	}

	// Check if token is already used
	if passwordReset.UsedAt != nil {
		return fmt.Errorf("reset token already used")
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, passwordReset.UserID)
	if err != nil {
		return entities.ErrUserNotFound
	}

	// Hash new password
	hashedPassword, err := uc.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	if err := uc.userRepo.UpdatePassword(ctx, user.ID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	if err := uc.passwordResetRepo.MarkAsUsed(ctx, req.Token); err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Log for testing
	fmt.Printf("Password reset requested with token: %s\n", req.Token)
	fmt.Printf("New password would be set to: %s\n", req.NewPassword)

	return nil
}

// generateRefreshToken generates a refresh token for the user
func (uc *userUseCase) generateRefreshToken(user *entities.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

// ResendVerification resends email verification
func (uc *userUseCase) ResendVerification(ctx context.Context, req ResendVerificationRequest) error {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		return nil
	}

	// Check if already verified
	if user.EmailVerified {
		return fmt.Errorf("email already verified")
	}

	// Send email verification using the existing SendEmailVerification method
	if err := uc.SendEmailVerification(ctx, user.ID); err != nil {
		fmt.Printf("❌ Failed to resend email verification to %s: %v\n", user.Email, err)
		return fmt.Errorf("failed to send verification email")
	}

	fmt.Printf("✅ Email verification resent for: %s\n", user.Email)
	return nil
}

// TrackSearch tracks user search activity
func (uc *userUseCase) TrackSearch(ctx context.Context, req TrackSearchRequest) error {
	// Create search history entry
	searchHistory := &entities.UserSearchHistory{
		ID:        uuid.New(),
		UserID:    req.UserID,
		Query:     req.Query,
		Results:   req.Results,
		Clicked:   req.Clicked,
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
		CreatedAt: time.Now(),
	}

	// Convert filters to JSON string
	if req.Filters != nil {
		filtersJSON, err := json.Marshal(req.Filters)
		if err == nil {
			searchHistory.Filters = string(filtersJSON)
		}
	}

	// TODO: Save to search history repository
	// For now, just track the activity
	_ = uc.TrackUserActivity(ctx, req.UserID, "search",
		fmt.Sprintf("Searched for: %s", req.Query), "search", nil, nil)

	return nil
}

// GetSearchHistory retrieves user's search history
func (uc *userUseCase) GetSearchHistory(ctx context.Context, userID uuid.UUID, req SearchHistoryRequest) (*UserSearchHistoryListResponse, error) {
	// TODO: Implement search history retrieval from repository
	// For now, return empty results
	return &UserSearchHistoryListResponse{
		History: []*SearchHistoryItem{},
		Total:   0,
	}, nil
}

// ClearSearchHistory clears user's search history
func (uc *userUseCase) ClearSearchHistory(ctx context.Context, userID uuid.UUID) error {
	// TODO: Implement search history clearing
	// For now, just track the activity
	_ = uc.TrackUserActivity(ctx, userID, "profile_update",
		"Cleared search history", "search_history", nil, nil)
	return nil
}

// DeleteSearchHistoryItem deletes a specific search history item
func (uc *userUseCase) DeleteSearchHistoryItem(ctx context.Context, userID, historyID uuid.UUID) error {
	// TODO: Implement specific search history item deletion
	// For now, just track the activity
	_ = uc.TrackUserActivity(ctx, userID, "profile_update",
		"Deleted search history item", "search_history", &historyID, nil)
	return nil
}

// GetPopularSearches gets user's popular searches
func (uc *userUseCase) GetPopularSearches(ctx context.Context, userID uuid.UUID, limit int) (*UserPopularSearchesResponse, error) {
	// TODO: Implement popular searches retrieval
	// For now, return empty results
	return &UserPopularSearchesResponse{
		Searches: []PopularSearchItem{},
	}, nil
}

// CreateSavedSearch creates a new saved search
func (uc *userUseCase) CreateSavedSearch(ctx context.Context, req CreateSavedSearchRequest) (*SavedSearchResponse, error) {
	// Create saved search entity
	savedSearch := &entities.SavedSearch{
		ID:           uuid.New(),
		UserID:       req.UserID,
		Name:         req.Name,
		Query:        req.Query,
		IsActive:     true,
		PriceAlert:   req.PriceAlert,
		StockAlert:   req.StockAlert,
		NewItemAlert: req.NewItemAlert,
		MaxPrice:     req.MaxPrice,
		MinPrice:     req.MinPrice,
		EmailNotify:  req.EmailNotify,
		PushNotify:   req.PushNotify,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Convert filters to JSON string
	if req.Filters != nil {
		filtersJSON, err := json.Marshal(req.Filters)
		if err == nil {
			savedSearch.Filters = string(filtersJSON)
		}
	}

	// TODO: Save to saved search repository
	// For now, just track the activity
	_ = uc.TrackUserActivity(ctx, req.UserID, "profile_update",
		fmt.Sprintf("Created saved search: %s", req.Name), "saved_search", &savedSearch.ID, nil)

	return &SavedSearchResponse{
		ID:           savedSearch.ID,
		Name:         savedSearch.Name,
		Query:        savedSearch.Query,
		Filters:      req.Filters,
		IsActive:     savedSearch.IsActive,
		PriceAlert:   savedSearch.PriceAlert,
		StockAlert:   savedSearch.StockAlert,
		NewItemAlert: savedSearch.NewItemAlert,
		MaxPrice:     savedSearch.MaxPrice,
		MinPrice:     savedSearch.MinPrice,
		EmailNotify:  savedSearch.EmailNotify,
		PushNotify:   savedSearch.PushNotify,
		CreatedAt:    savedSearch.CreatedAt,
		UpdatedAt:    savedSearch.UpdatedAt,
	}, nil
}

// GetSavedSearches retrieves user's saved searches
func (uc *userUseCase) GetSavedSearches(ctx context.Context, userID uuid.UUID, req GetSavedSearchesRequest) (*GetSavedSearchesResponse, error) {
	// TODO: Implement saved searches retrieval from repository
	// For now, return empty results
	return &GetSavedSearchesResponse{
		SavedSearches: []*SavedSearchResponse{},
		Total:         0,
	}, nil
}

// TrackProductView tracks user product viewing activity
func (uc *userUseCase) TrackProductView(ctx context.Context, req TrackProductViewRequest) error {
	// TODO: Save to browsing history repository
	// For now, just track the activity
	_ = uc.TrackUserActivity(ctx, req.UserID, "product_view",
		"Viewed product", "product", &req.ProductID, nil)

	// Update personalization data
	// TODO: Update user personalization based on product view

	return nil
}

// GetPersonalization retrieves user personalization data
func (uc *userUseCase) GetPersonalization(ctx context.Context, userID uuid.UUID) (*PersonalizationResponse, error) {
	// TODO: Implement personalization retrieval from repository
	// For now, return default personalization
	return &PersonalizationResponse{
		ID:                  uuid.New(),
		UserID:              userID,
		CategoryPreferences: make(map[string]float64),
		BrandPreferences:    make(map[string]float64),
		PriceRangePreference: PriceRangePreference{
			MinPrice: 0,
			MaxPrice: 1000,
			Currency: "USD",
		},
		AverageOrderValue:     0,
		PreferredShoppingTime: "evening",
		ShoppingFrequency:     "weekly",
		RecommendationEngine:  "collaborative",
		PersonalizationLevel:  "medium",
		TotalViews:            0,
		TotalSearches:         0,
		UniqueProductsViewed:  0,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}, nil
}

// UpdateSavedSearch updates an existing saved search
func (uc *userUseCase) UpdateSavedSearch(ctx context.Context, req UpdateSavedSearchRequest) (*SavedSearchResponse, error) {
	// TODO: Implement saved search update
	// For now, just track the activity
	_ = uc.TrackUserActivity(ctx, req.UserID, "profile_update",
		"Updated saved search", "saved_search", &req.SavedSearchID, nil)

	return &SavedSearchResponse{
		ID:        req.SavedSearchID,
		Name:      *req.Name,
		Query:     *req.Query,
		IsActive:  *req.IsActive,
		UpdatedAt: time.Now(),
	}, nil
}

// DeleteSavedSearch deletes a saved search
func (uc *userUseCase) DeleteSavedSearch(ctx context.Context, userID, savedSearchID uuid.UUID) error {
	// TODO: Implement saved search deletion
	// For now, just track the activity
	_ = uc.TrackUserActivity(ctx, userID, "profile_update",
		"Deleted saved search", "saved_search", &savedSearchID, nil)
	return nil
}

// ExecuteSavedSearch executes a saved search
func (uc *userUseCase) ExecuteSavedSearch(ctx context.Context, userID, savedSearchID uuid.UUID) (*SavedSearchExecutionResponse, error) {
	// TODO: Implement saved search execution
	// For now, return placeholder response
	return &SavedSearchExecutionResponse{
		SavedSearchID: savedSearchID,
		Query:         "placeholder query",
		Results:       0,
		ExecutedAt:    time.Now(),
		SearchURL:     "/search?q=placeholder",
	}, nil
}

// GetBrowsingHistory retrieves user's browsing history
func (uc *userUseCase) GetBrowsingHistory(ctx context.Context, userID uuid.UUID, req BrowsingHistoryRequest) (*BrowsingHistoryResponse, error) {
	// TODO: Implement browsing history retrieval from repository
	// For now, return empty results
	return &BrowsingHistoryResponse{
		History: []*BrowsingHistoryItem{},
		Total:   0,
	}, nil
}

// ClearBrowsingHistory clears user's browsing history
func (uc *userUseCase) ClearBrowsingHistory(ctx context.Context, userID uuid.UUID) error {
	// TODO: Implement browsing history clearing
	// For now, just track the activity
	_ = uc.TrackUserActivity(ctx, userID, "profile_update",
		"Cleared browsing history", "browsing_history", nil, nil)
	return nil
}

// UpdatePersonalization updates user personalization data
func (uc *userUseCase) UpdatePersonalization(ctx context.Context, req UpdatePersonalizationRequest) (*PersonalizationResponse, error) {
	// TODO: Implement personalization update
	// For now, just track the activity and return current data
	_ = uc.TrackUserActivity(ctx, req.UserID, "profile_update",
		"Updated personalization settings", "personalization", nil, nil)

	return uc.GetPersonalization(ctx, req.UserID)
}

// GetPersonalizedRecommendations gets personalized recommendations for user
func (uc *userUseCase) GetPersonalizedRecommendations(ctx context.Context, userID uuid.UUID, req PersonalizedRecommendationsRequest) (*PersonalizedRecommendationsResponse, error) {
	// TODO: Implement personalized recommendations
	// For now, return empty recommendations
	return &PersonalizedRecommendationsResponse{
		Type:            req.Type,
		Recommendations: []PersonalizedRecommendation{},
		Algorithm:       "collaborative",
		GeneratedAt:     time.Now(),
	}, nil
}

// AnalyzeUserBehavior analyzes user behavior and returns insights
func (uc *userUseCase) AnalyzeUserBehavior(ctx context.Context, userID uuid.UUID) (*UserBehaviorAnalysisResponse, error) {
	// TODO: Implement user behavior analysis
	// For now, return placeholder analysis
	return &UserBehaviorAnalysisResponse{
		UserID:        userID,
		TopCategories: []CategoryPreference{},
		TopBrands:     []BrandPreference{},
		PriceAnalysis: PriceAnalysis{
			MinPrice:     0,
			MaxPrice:     1000,
			AveragePrice: 100,
			PriceSegment: "mid-range",
		},
		ShoppingPatterns: ShoppingPatterns{
			PreferredDays:        []string{"Saturday", "Sunday"},
			PreferredHours:       []int{19, 20, 21},
			AverageSessionLength: 15.5,
			PagesPerSession:      8.2,
			ConversionRate:       0.05,
		},
		EngagementScore: 75.0,
		LoyaltyScore:    60.0,
		Insights: []BehaviorInsight{
			{
				Type:        "trend",
				Title:       "Weekend Shopping Preference",
				Description: "You tend to shop more on weekends",
				Confidence:  0.85,
				ActionItems: []string{"Check weekend deals", "Set weekend shopping reminders"},
			},
		},
		AnalyzedAt: time.Now(),
	}, nil
}

// GetProfileAnalytics gets user profile analytics
func (uc *userUseCase) GetProfileAnalytics(ctx context.Context, userID uuid.UUID) (*ProfileAnalyticsResponse, error) {
	// TODO: Implement profile analytics
	// For now, return placeholder analytics
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &ProfileAnalyticsResponse{
		UserID: userID,
		Overview: struct {
			TotalViews        int       `json:"total_views"`
			TotalSearches     int       `json:"total_searches"`
			TotalOrders       int       `json:"total_orders"`
			TotalSpent        float64   `json:"total_spent"`
			AverageOrderValue float64   `json:"average_order_value"`
			LastActivity      time.Time `json:"last_activity"`
			MemberSince       time.Time `json:"member_since"`
			EngagementScore   float64   `json:"engagement_score"`
			LoyaltyScore      float64   `json:"loyalty_score"`
		}{
			TotalViews:        0,
			TotalSearches:     0,
			TotalOrders:       user.TotalOrders,
			TotalSpent:        user.TotalSpent,
			AverageOrderValue: 0,
			LastActivity:      time.Now(),
			MemberSince:       user.CreatedAt,
			EngagementScore:   75.0,
			LoyaltyScore:      60.0,
		},
		ActivityTrends: []ActivityTrendData{},
		TopCategories:  []CategoryStats{},
		TopBrands:      []BrandStats{},
		Preferences: struct {
			Theme                string `json:"theme"`
			Language             string `json:"language"`
			Currency             string `json:"currency"`
			NotificationsEnabled bool   `json:"notifications_enabled"`
			PersonalizationLevel string `json:"personalization_level"`
		}{
			Theme:                user.Language,
			Language:             user.Language,
			Currency:             user.Currency,
			NotificationsEnabled: true,
			PersonalizationLevel: "medium",
		},
	}, nil
}

// GetActivitySummary gets user activity summary for a time range
func (uc *userUseCase) GetActivitySummary(ctx context.Context, userID uuid.UUID, timeRange string) (*ActivitySummaryResponse, error) {
	// TODO: Implement activity summary
	// For now, return placeholder summary
	now := time.Now()
	var startDate time.Time

	switch timeRange {
	case "day":
		startDate = now.AddDate(0, 0, -1)
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	default:
		startDate = now.AddDate(0, 0, -7) // Default to week
	}

	return &ActivitySummaryResponse{
		UserID:    userID,
		TimeRange: timeRange,
		Period: struct {
			StartDate time.Time `json:"start_date"`
			EndDate   time.Time `json:"end_date"`
		}{
			StartDate: startDate,
			EndDate:   now,
		},
		Summary: struct {
			Views          int     `json:"views"`
			Searches       int     `json:"searches"`
			Orders         int     `json:"orders"`
			AmountSpent    float64 `json:"amount_spent"`
			TimeSpent      int     `json:"time_spent"`
			PagesVisited   int     `json:"pages_visited"`
			UniqueProducts int     `json:"unique_products"`
		}{
			Views:          0,
			Searches:       0,
			Orders:         0,
			AmountSpent:    0,
			TimeSpent:      0,
			PagesVisited:   0,
			UniqueProducts: 0,
		},
		DailyActivity: []DailyActivityData{},
		TopActions:    []ActionData{},
	}, nil
}

// Search history request/response types
type TrackSearchRequest struct {
	UserID    uuid.UUID              `json:"user_id" validate:"required"`
	Query     string                 `json:"query" validate:"required"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
	Results   int                    `json:"results"`
	Clicked   bool                   `json:"clicked"`
	IPAddress string                 `json:"ip_address,omitempty"`
	UserAgent string                 `json:"user_agent,omitempty"`
}

type SearchHistoryRequest struct {
	Query    *string    `json:"query,omitempty"`
	DateFrom *time.Time `json:"date_from,omitempty"`
	DateTo   *time.Time `json:"date_to,omitempty"`
	Limit    int        `json:"limit,omitempty"`
	Offset   int        `json:"offset,omitempty"`
}

type UserSearchHistoryListResponse struct {
	History []*SearchHistoryItem `json:"history"`
	Total   int64                `json:"total"`
}

type SearchHistoryItem struct {
	ID        uuid.UUID              `json:"id"`
	Query     string                 `json:"query"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
	Results   int                    `json:"results"`
	Clicked   bool                   `json:"clicked"`
	CreatedAt time.Time              `json:"created_at"`
}

type UserPopularSearchesResponse struct {
	Searches []PopularSearchItem `json:"searches"`
}

type PopularSearchItem struct {
	Query       string    `json:"query"`
	SearchCount int       `json:"search_count"`
	LastUsed    time.Time `json:"last_used"`
}

// Saved searches request/response types
type CreateSavedSearchRequest struct {
	UserID       uuid.UUID              `json:"user_id" validate:"required"`
	Name         string                 `json:"name" validate:"required"`
	Query        string                 `json:"query" validate:"required"`
	Filters      map[string]interface{} `json:"filters,omitempty"`
	PriceAlert   bool                   `json:"price_alert"`
	StockAlert   bool                   `json:"stock_alert"`
	NewItemAlert bool                   `json:"new_item_alert"`
	MaxPrice     *float64               `json:"max_price,omitempty"`
	MinPrice     *float64               `json:"min_price,omitempty"`
	EmailNotify  bool                   `json:"email_notify"`
	PushNotify   bool                   `json:"push_notify"`
}

type SavedSearchResponse struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	Query        string                 `json:"query"`
	Filters      map[string]interface{} `json:"filters,omitempty"`
	IsActive     bool                   `json:"is_active"`
	PriceAlert   bool                   `json:"price_alert"`
	StockAlert   bool                   `json:"stock_alert"`
	NewItemAlert bool                   `json:"new_item_alert"`
	MaxPrice     *float64               `json:"max_price,omitempty"`
	MinPrice     *float64               `json:"min_price,omitempty"`
	EmailNotify  bool                   `json:"email_notify"`
	PushNotify   bool                   `json:"push_notify"`
	LastChecked  *time.Time             `json:"last_checked,omitempty"`
	LastNotified *time.Time             `json:"last_notified,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type GetSavedSearchesRequest struct {
	IsActive *bool `json:"is_active,omitempty"`
	Limit    int   `json:"limit,omitempty"`
	Offset   int   `json:"offset,omitempty"`
}

type GetSavedSearchesResponse struct {
	SavedSearches []*SavedSearchResponse `json:"saved_searches"`
	Total         int64                  `json:"total"`
}

type UpdateSavedSearchRequest struct {
	UserID        uuid.UUID              `json:"user_id" validate:"required"`
	SavedSearchID uuid.UUID              `json:"saved_search_id" validate:"required"`
	Name          *string                `json:"name,omitempty"`
	Query         *string                `json:"query,omitempty"`
	Filters       map[string]interface{} `json:"filters,omitempty"`
	IsActive      *bool                  `json:"is_active,omitempty"`
	PriceAlert    *bool                  `json:"price_alert,omitempty"`
	StockAlert    *bool                  `json:"stock_alert,omitempty"`
	NewItemAlert  *bool                  `json:"new_item_alert,omitempty"`
	MaxPrice      *float64               `json:"max_price,omitempty"`
	MinPrice      *float64               `json:"min_price,omitempty"`
	EmailNotify   *bool                  `json:"email_notify,omitempty"`
	PushNotify    *bool                  `json:"push_notify,omitempty"`
}

type SavedSearchExecutionResponse struct {
	SavedSearchID uuid.UUID `json:"saved_search_id"`
	Query         string    `json:"query"`
	Results       int       `json:"results"`
	ExecutedAt    time.Time `json:"executed_at"`
	SearchURL     string    `json:"search_url"`
}

// Browsing history request/response types
type TrackProductViewRequest struct {
	UserID       uuid.UUID  `json:"user_id" validate:"required"`
	ProductID    uuid.UUID  `json:"product_id" validate:"required"`
	CategoryID   *uuid.UUID `json:"category_id,omitempty"`
	ViewDuration int        `json:"view_duration"` // in seconds
	Source       string     `json:"source"`        // search, category, recommendation, etc.
	IPAddress    string     `json:"ip_address,omitempty"`
	UserAgent    string     `json:"user_agent,omitempty"`
}

type BrowsingHistoryRequest struct {
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	Source     *string    `json:"source,omitempty"`
	DateFrom   *time.Time `json:"date_from,omitempty"`
	DateTo     *time.Time `json:"date_to,omitempty"`
	Limit      int        `json:"limit,omitempty"`
	Offset     int        `json:"offset,omitempty"`
}

type BrowsingHistoryResponse struct {
	History []*BrowsingHistoryItem `json:"history"`
	Total   int64                  `json:"total"`
}

type BrowsingHistoryItem struct {
	ID           uuid.UUID  `json:"id"`
	ProductID    uuid.UUID  `json:"product_id"`
	ProductName  string     `json:"product_name"`
	ProductImage string     `json:"product_image,omitempty"`
	CategoryID   *uuid.UUID `json:"category_id,omitempty"`
	CategoryName string     `json:"category_name,omitempty"`
	ViewDuration int        `json:"view_duration"`
	Source       string     `json:"source"`
	CreatedAt    time.Time  `json:"created_at"`
}

// Personalization request/response types
type PersonalizationResponse struct {
	ID                    uuid.UUID            `json:"id"`
	UserID                uuid.UUID            `json:"user_id"`
	CategoryPreferences   map[string]float64   `json:"category_preferences"`
	BrandPreferences      map[string]float64   `json:"brand_preferences"`
	PriceRangePreference  PriceRangePreference `json:"price_range_preference"`
	AverageOrderValue     float64              `json:"average_order_value"`
	PreferredShoppingTime string               `json:"preferred_shopping_time"`
	ShoppingFrequency     string               `json:"shopping_frequency"`
	RecommendationEngine  string               `json:"recommendation_engine"`
	PersonalizationLevel  string               `json:"personalization_level"`
	TotalViews            int                  `json:"total_views"`
	TotalSearches         int                  `json:"total_searches"`
	UniqueProductsViewed  int                  `json:"unique_products_viewed"`
	LastAnalyzed          *time.Time           `json:"last_analyzed,omitempty"`
	CreatedAt             time.Time            `json:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at"`
}

type UpdatePersonalizationRequest struct {
	UserID                uuid.UUID             `json:"user_id" validate:"required"`
	CategoryPreferences   map[string]float64    `json:"category_preferences,omitempty"`
	BrandPreferences      map[string]float64    `json:"brand_preferences,omitempty"`
	PriceRangePreference  *PriceRangePreference `json:"price_range_preference,omitempty"`
	PreferredShoppingTime *string               `json:"preferred_shopping_time,omitempty"`
	ShoppingFrequency     *string               `json:"shopping_frequency,omitempty"`
	RecommendationEngine  *string               `json:"recommendation_engine,omitempty"`
	PersonalizationLevel  *string               `json:"personalization_level,omitempty"`
}

type PriceRangePreference struct {
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
	Currency string  `json:"currency"`
}

type PersonalizedRecommendationsRequest struct {
	UserID     uuid.UUID   `json:"user_id" validate:"required"`
	Type       string      `json:"type"` // products, categories, brands
	Limit      int         `json:"limit,omitempty"`
	CategoryID *uuid.UUID  `json:"category_id,omitempty"`
	Exclude    []uuid.UUID `json:"exclude,omitempty"` // Exclude specific items
}

type PersonalizedRecommendationsResponse struct {
	Type            string                       `json:"type"`
	Recommendations []PersonalizedRecommendation `json:"recommendations"`
	Algorithm       string                       `json:"algorithm"`
	GeneratedAt     time.Time                    `json:"generated_at"`
}

type PersonalizedRecommendation struct {
	ID          uuid.UUID              `json:"id"`
	Type        string                 `json:"type"` // product, category, brand
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Image       string                 `json:"image,omitempty"`
	Price       *float64               `json:"price,omitempty"`
	Score       float64                `json:"score"`
	Reason      string                 `json:"reason"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type UserBehaviorAnalysisResponse struct {
	UserID           uuid.UUID            `json:"user_id"`
	TopCategories    []CategoryPreference `json:"top_categories"`
	TopBrands        []BrandPreference    `json:"top_brands"`
	PriceAnalysis    PriceAnalysis        `json:"price_analysis"`
	ShoppingPatterns ShoppingPatterns     `json:"shopping_patterns"`
	EngagementScore  float64              `json:"engagement_score"`
	LoyaltyScore     float64              `json:"loyalty_score"`
	Insights         []BehaviorInsight    `json:"insights"`
	AnalyzedAt       time.Time            `json:"analyzed_at"`
}

type CategoryPreference struct {
	CategoryID    uuid.UUID `json:"category_id"`
	CategoryName  string    `json:"category_name"`
	Score         float64   `json:"score"`
	ViewCount     int       `json:"view_count"`
	PurchaseCount int       `json:"purchase_count"`
}

type BrandPreference struct {
	BrandID       uuid.UUID `json:"brand_id"`
	BrandName     string    `json:"brand_name"`
	Score         float64   `json:"score"`
	ViewCount     int       `json:"view_count"`
	PurchaseCount int       `json:"purchase_count"`
}

type PriceAnalysis struct {
	MinPrice      float64 `json:"min_price"`
	MaxPrice      float64 `json:"max_price"`
	AveragePrice  float64 `json:"average_price"`
	PriceVariance float64 `json:"price_variance"`
	PriceSegment  string  `json:"price_segment"` // budget, mid-range, premium
}

type ShoppingPatterns struct {
	PreferredDays        []string `json:"preferred_days"`
	PreferredHours       []int    `json:"preferred_hours"`
	AverageSessionLength float64  `json:"average_session_length"`
	PagesPerSession      float64  `json:"average_pages_per_session"`
	ConversionRate       float64  `json:"conversion_rate"`
}

type BehaviorInsight struct {
	Type        string   `json:"type"` // trend, preference, opportunity
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Confidence  float64  `json:"confidence"`
	ActionItems []string `json:"action_items,omitempty"`
}

// Profile analytics request/response types
type ProfileAnalyticsResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Overview struct {
		TotalViews        int       `json:"total_views"`
		TotalSearches     int       `json:"total_searches"`
		TotalOrders       int       `json:"total_orders"`
		TotalSpent        float64   `json:"total_spent"`
		AverageOrderValue float64   `json:"average_order_value"`
		LastActivity      time.Time `json:"last_activity"`
		MemberSince       time.Time `json:"member_since"`
		EngagementScore   float64   `json:"engagement_score"`
		LoyaltyScore      float64   `json:"loyalty_score"`
	} `json:"overview"`

	ActivityTrends []ActivityTrendData `json:"activity_trends"`
	TopCategories  []CategoryStats     `json:"top_categories"`
	TopBrands      []BrandStats        `json:"top_brands"`

	Preferences struct {
		Theme                string `json:"theme"`
		Language             string `json:"language"`
		Currency             string `json:"currency"`
		NotificationsEnabled bool   `json:"notifications_enabled"`
		PersonalizationLevel string `json:"personalization_level"`
	} `json:"preferences"`
}

type ActivitySummaryResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	TimeRange string    `json:"time_range"`
	Period    struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	} `json:"period"`

	Summary struct {
		Views          int     `json:"views"`
		Searches       int     `json:"searches"`
		Orders         int     `json:"orders"`
		AmountSpent    float64 `json:"amount_spent"`
		TimeSpent      int     `json:"time_spent"` // in minutes
		PagesVisited   int     `json:"pages_visited"`
		UniqueProducts int     `json:"unique_products"`
	} `json:"summary"`

	DailyActivity []DailyActivityData `json:"daily_activity"`
	TopActions    []ActionData        `json:"top_actions"`
}

type ActivityTrendData struct {
	Date     time.Time `json:"date"`
	Views    int       `json:"views"`
	Searches int       `json:"searches"`
	Orders   int       `json:"orders"`
	Spent    float64   `json:"spent"`
}

type CategoryStats struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Views        int       `json:"views"`
	Orders       int       `json:"orders"`
	AmountSpent  float64   `json:"amount_spent"`
}

type BrandStats struct {
	BrandID     uuid.UUID `json:"brand_id"`
	BrandName   string    `json:"brand_name"`
	Views       int       `json:"views"`
	Orders      int       `json:"orders"`
	AmountSpent float64   `json:"amount_spent"`
}

type DailyActivityData struct {
	Date         time.Time `json:"date"`
	Views        int       `json:"views"`
	Searches     int       `json:"searches"`
	TimeSpent    int       `json:"time_spent"`
	PagesVisited int       `json:"pages_visited"`
}

type ActionData struct {
	Action string `json:"action"`
	Count  int    `json:"count"`
}

// GetLoginHistory retrieves user's login history
func (uc *userUseCase) GetLoginHistory(ctx context.Context, userID uuid.UUID, req LoginHistoryRequest) (*LoginHistoryResponse, error) {
	// Set default values
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Get login history from repository
	loginHistoryEntities, err := uc.userLoginHistoryRepo.GetByUserID(ctx, userID, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get login history: %w", err)
	}

	// Convert to response format
	loginHistory := make([]LoginHistoryItem, len(loginHistoryEntities))
	for i, entity := range loginHistoryEntities {
		loginHistory[i] = LoginHistoryItem{
			ID:         entity.ID,
			IPAddress:  entity.IPAddress,
			UserAgent:  entity.UserAgent,
			DeviceInfo: entity.DeviceInfo,
			Location:   entity.Location,
			LoginType:  entity.LoginType,
			Success:    entity.Success,
			FailReason: entity.FailReason,
			CreatedAt:  entity.CreatedAt,
		}
	}

	// Get total count for pagination
	totalCount, err := uc.userLoginHistoryRepo.CountLoginAttempts(ctx, userID, time.Time{})
	if err != nil {
		totalCount = int64(len(loginHistory))
	}

	// Calculate stats
	stats, err := uc.calculateLoginStats(ctx, userID, req.DateFrom, req.DateTo)
	if err != nil {
		// Log error but don't fail the request
		stats = nil
	}

	// Create pagination info
	pagination := &PaginationInfo{
		Page:       (req.Offset / req.Limit) + 1,
		Limit:      req.Limit,
		Total:      totalCount,
		TotalPages: int((totalCount + int64(req.Limit) - 1) / int64(req.Limit)),
	}

	return &LoginHistoryResponse{
		LoginHistory: loginHistory,
		Total:        totalCount,
		Pagination:   pagination,
		Stats:        stats,
	}, nil
}

// GetLoginStats retrieves user's login statistics
func (uc *userUseCase) GetLoginStats(ctx context.Context, userID uuid.UUID) (*LoginStatsResponse, error) {
	// Get all login attempts
	allLogins, err := uc.userLoginHistoryRepo.GetByUserID(ctx, userID, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get login history: %w", err)
	}

	// Calculate statistics
	var totalLogins, successfulLogins, failedLogins int64
	var lastLogin, lastFailedLogin *time.Time
	uniqueIPs := make(map[string]bool)
	uniqueDevices := make(map[string]bool)
	deviceCount := make(map[string]int)
	locationCount := make(map[string]int)

	// Recent failed attempts (last 24 hours)
	recentFailedAttempts, _ := uc.userLoginHistoryRepo.CountFailedAttempts(ctx, userID, time.Now().Add(-24*time.Hour))

	for _, login := range allLogins {
		totalLogins++

		if login.Success {
			successfulLogins++
			if lastLogin == nil || login.CreatedAt.After(*lastLogin) {
				lastLogin = &login.CreatedAt
			}
		} else {
			failedLogins++
			if lastFailedLogin == nil || login.CreatedAt.After(*lastFailedLogin) {
				lastFailedLogin = &login.CreatedAt
			}
		}

		// Track unique IPs
		if login.IPAddress != "" {
			uniqueIPs[login.IPAddress] = true
		}

		// Track unique devices
		if login.DeviceInfo != "" {
			uniqueDevices[login.DeviceInfo] = true
			deviceCount[login.DeviceInfo]++
		}

		// Track locations
		if login.Location != "" {
			locationCount[login.Location]++
		}
	}

	// Calculate success rate
	var successRate float64
	if totalLogins > 0 {
		successRate = float64(successfulLogins) / float64(totalLogins) * 100
	}

	// Find most used device and location
	var mostUsedDevice, mostUsedLocation string
	maxDeviceCount, maxLocationCount := 0, 0

	for device, count := range deviceCount {
		if count > maxDeviceCount {
			maxDeviceCount = count
			mostUsedDevice = device
		}
	}

	for location, count := range locationCount {
		if count > maxLocationCount {
			maxLocationCount = count
			mostUsedLocation = location
		}
	}

	return &LoginStatsResponse{
		UserID:               userID,
		TotalLogins:          totalLogins,
		SuccessfulLogins:     successfulLogins,
		FailedLogins:         failedLogins,
		SuccessRate:          successRate,
		LastLogin:            lastLogin,
		LastFailedLogin:      lastFailedLogin,
		UniqueIPs:            len(uniqueIPs),
		UniqueDevices:        len(uniqueDevices),
		MostUsedDevice:       mostUsedDevice,
		MostUsedLocation:     mostUsedLocation,
		RecentFailedAttempts: recentFailedAttempts,
	}, nil
}

// calculateLoginStats calculates login statistics for a date range
func (uc *userUseCase) calculateLoginStats(ctx context.Context, userID uuid.UUID, dateFrom, dateTo *time.Time) (*LoginStatsInfo, error) {
	// Use default date range if not provided
	if dateFrom == nil {
		defaultFrom := time.Now().AddDate(0, -1, 0) // Last month
		dateFrom = &defaultFrom
	}
	if dateTo == nil {
		defaultTo := time.Now()
		dateTo = &defaultTo
	}

	// Get login attempts in date range
	totalCount, err := uc.userLoginHistoryRepo.CountLoginAttempts(ctx, userID, *dateFrom)
	if err != nil {
		return nil, err
	}

	failedCount, err := uc.userLoginHistoryRepo.CountFailedAttempts(ctx, userID, *dateFrom)
	if err != nil {
		return nil, err
	}

	successfulCount := totalCount - failedCount

	var successRate float64
	if totalCount > 0 {
		successRate = float64(successfulCount) / float64(totalCount) * 100
	}

	// Get unique IPs count (simplified - would need additional repository method for exact count)
	logins, err := uc.userLoginHistoryRepo.GetByUserID(ctx, userID, 100, 0)
	if err != nil {
		return nil, err
	}

	uniqueIPs := make(map[string]bool)
	for _, login := range logins {
		if login.IPAddress != "" && login.CreatedAt.After(*dateFrom) && login.CreatedAt.Before(*dateTo) {
			uniqueIPs[login.IPAddress] = true
		}
	}

	return &LoginStatsInfo{
		TotalLogins:      totalCount,
		SuccessfulLogins: successfulCount,
		FailedLogins:     failedCount,
		SuccessRate:      successRate,
		UniqueIPs:        len(uniqueIPs),
	}, nil
}

// logLoginAttemptEnhanced logs a login attempt with enhanced tracking information
func (uc *userUseCase) logLoginAttemptEnhanced(ctx context.Context, email string, success bool, failReason string, ipAddress, userAgent, deviceInfo string) error {
	// Try to get user ID by email
	var userID uuid.UUID
	if user, err := uc.userRepo.GetByEmail(ctx, email); err == nil {
		userID = user.ID
	}

	// Set default values if not provided
	if ipAddress == "" {
		ipAddress = "unknown"
	}
	if userAgent == "" {
		userAgent = "unknown"
	}
	if deviceInfo == "" {
		deviceInfo = uc.extractDeviceInfoFromUserAgent(userAgent)
	}

	loginHistory := &entities.UserLoginHistory{
		ID:         uuid.New(),
		UserID:     userID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		DeviceInfo: deviceInfo,
		Location:   uc.getLocationFromIP(ipAddress),
		LoginType:  "password",
		Success:    success,
		FailReason: failReason,
		CreatedAt:  time.Now(),
	}

	return uc.userLoginHistoryRepo.Create(ctx, loginHistory)
}

// logLoginAttempt logs a login attempt (legacy method for backward compatibility)
func (uc *userUseCase) logLoginAttempt(ctx context.Context, email string, success bool, failReason string, ipAddress string) error {
	return uc.logLoginAttemptEnhanced(ctx, email, success, failReason, ipAddress, "", "")
}

// getLocationFromIP gets location information from IP address
func (uc *userUseCase) getLocationFromIP(ipAddress string) string {
	// TODO: Implement IP geolocation service integration
	// For now, return empty string or basic classification
	if ipAddress == "" || ipAddress == "unknown" || ipAddress == "127.0.0.1" || ipAddress == "::1" {
		return "Local"
	}

	// This would integrate with a geolocation service like MaxMind, IPStack, etc.
	return "Unknown Location"
}

// extractDeviceInfoFromUserAgent extracts device information from user agent string
func (uc *userUseCase) extractDeviceInfoFromUserAgent(userAgent string) string {
	if userAgent == "" || userAgent == "unknown" {
		return "Unknown Device"
	}

	// Simple device detection - in production you'd use a proper user agent parser
	userAgent = strings.ToLower(userAgent)

	if strings.Contains(userAgent, "mobile") || strings.Contains(userAgent, "android") || strings.Contains(userAgent, "iphone") {
		if strings.Contains(userAgent, "android") {
			return "Android Mobile"
		} else if strings.Contains(userAgent, "iphone") {
			return "iPhone"
		} else if strings.Contains(userAgent, "ipad") {
			return "iPad"
		}
		return "Mobile Device"
	} else if strings.Contains(userAgent, "tablet") || strings.Contains(userAgent, "ipad") {
		return "Tablet"
	} else if strings.Contains(userAgent, "windows") {
		return "Windows Desktop"
	} else if strings.Contains(userAgent, "mac") || strings.Contains(userAgent, "macintosh") {
		return "Mac Desktop"
	} else if strings.Contains(userAgent, "linux") {
		return "Linux Desktop"
	}

	return "Desktop Browser"
}

// Login History request/response types
type LoginHistoryRequest struct {
	Limit    int        `json:"limit" validate:"min=1,max=100"`
	Offset   int        `json:"offset" validate:"min=0"`
	DateFrom *time.Time `json:"date_from,omitempty"`
	DateTo   *time.Time `json:"date_to,omitempty"`
	Success  *bool      `json:"success,omitempty"` // Filter by success/failure
}

type LoginHistoryResponse struct {
	LoginHistory []LoginHistoryItem `json:"login_history"`
	Total        int64              `json:"total"`
	Pagination   *PaginationInfo    `json:"pagination"`
	Stats        *LoginStatsInfo    `json:"stats,omitempty"`
}

type LoginHistoryItem struct {
	ID         uuid.UUID `json:"id"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	DeviceInfo string    `json:"device_info"`
	Location   string    `json:"location"`
	LoginType  string    `json:"login_type"`
	Success    bool      `json:"success"`
	FailReason string    `json:"fail_reason,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type LoginStatsRequest struct {
	DateFrom *time.Time `json:"date_from,omitempty"`
	DateTo   *time.Time `json:"date_to,omitempty"`
}

type LoginStatsResponse struct {
	UserID              uuid.UUID           `json:"user_id"`
	TotalLogins         int64               `json:"total_logins"`
	SuccessfulLogins    int64               `json:"successful_logins"`
	FailedLogins        int64               `json:"failed_logins"`
	SuccessRate         float64             `json:"success_rate"`
	LastLogin           *time.Time          `json:"last_login,omitempty"`
	LastFailedLogin     *time.Time          `json:"last_failed_login,omitempty"`
	UniqueIPs           int                 `json:"unique_ips"`
	UniqueDevices       int                 `json:"unique_devices"`
	MostUsedDevice      string              `json:"most_used_device,omitempty"`
	MostUsedLocation    string              `json:"most_used_location,omitempty"`
	RecentFailedAttempts int64              `json:"recent_failed_attempts"` // Last 24h
	LoginsByDay         []LoginDayStats     `json:"logins_by_day,omitempty"`
}

type LoginStatsInfo struct {
	TotalLogins      int64   `json:"total_logins"`
	SuccessfulLogins int64   `json:"successful_logins"`
	FailedLogins     int64   `json:"failed_logins"`
	SuccessRate      float64 `json:"success_rate"`
	UniqueIPs        int     `json:"unique_ips"`
}

type LoginDayStats struct {
	Date             time.Time `json:"date"`
	TotalLogins      int64     `json:"total_logins"`
	SuccessfulLogins int64     `json:"successful_logins"`
	FailedLogins     int64     `json:"failed_logins"`
}
