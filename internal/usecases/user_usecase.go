package usecases

import (
	"context"
	"fmt"
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
	TrackUserActivity(ctx context.Context, userID uuid.UUID, activityType entities.ActivityType, description string, entityType string, entityID *uuid.UUID, metadata map[string]interface{}) error
	GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStatsResponse, error)

	// User preferences methods
	GetUserPreferences(ctx context.Context, userID uuid.UUID) (*UserPreferencesResponse, error)
	UpdateUserPreferences(ctx context.Context, userID uuid.UUID, req UpdateUserPreferencesRequest) (*UserPreferencesResponse, error)
	UpdateTheme(ctx context.Context, userID uuid.UUID, theme string) error
	UpdateLanguage(ctx context.Context, userID uuid.UUID, language string) error

	// User verification methods
	SendEmailVerification(ctx context.Context, userID uuid.UUID) error
	VerifyEmail(ctx context.Context, token string) error
	SendPhoneVerification(ctx context.Context, userID uuid.UUID, phone string) error
	VerifyPhone(ctx context.Context, userID uuid.UUID, code string) error
	GetVerificationStatus(ctx context.Context, userID uuid.UUID) (*VerificationStatusResponse, error)
}

type userUseCase struct {
	userRepo             repositories.UserRepository
	userProfileRepo      repositories.UserProfileRepository
	userSessionRepo      repositories.UserSessionRepository
	userLoginHistoryRepo repositories.UserLoginHistoryRepository
	userActivityRepo     repositories.UserActivityRepository
	userPreferencesRepo  repositories.UserPreferencesRepository
	userVerificationRepo repositories.UserVerificationRepository
	passwordService      services.PasswordService
	jwtSecret            string
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
	passwordService services.PasswordService,
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
		passwordService:      passwordService,
		jwtSecret:            jwtSecret,
	}
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Phone     string `json:"phone"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
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

// UserResponse represents user response
type UserResponse struct {
	ID        uuid.UUID            `json:"id"`
	Email     string               `json:"email"`
	FirstName string               `json:"first_name"`
	LastName  string               `json:"last_name"`
	Phone     string               `json:"phone"`
	Role      entities.UserRole    `json:"role"`
	IsActive  bool                 `json:"is_active"`
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
	Pagination Pagination             `json:"pagination"`
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
	Pagination Pagination          `json:"pagination"`
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

	return uc.toUserResponse(user), nil
}

// Login authenticates a user
func (uc *userUseCase) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, entities.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, entities.ErrUserNotActive
	}

	// Check password
	if err := uc.passwordService.CheckPassword(req.Password, user.Password); err != nil {
		return nil, entities.ErrInvalidCredentials
	}

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
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

// toUserResponse converts user entity to response
func (uc *userUseCase) toUserResponse(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Role:      user.Role,
		IsActive:  user.IsActive,
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

	return &UserSessionsResponse{
		Sessions: sessionResponses,
		Total:    total,
		Pagination: Pagination{
			Page:       (offset / limit) + 1,
			Limit:      limit,
			Total:      total,
			TotalPages: int((total + int64(limit) - 1) / int64(limit)),
			HasNext:    offset+limit < int(total),
			HasPrev:    offset > 0,
		},
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

	return &UserActivityResponse{
		Activities: activityItems,
		Total:      total,
		Pagination: Pagination{
			Page:       (offset / limit) + 1,
			Limit:      limit,
			Total:      total,
			TotalPages: int((total + int64(limit) - 1) / int64(limit)),
			HasNext:    offset+limit < int(total),
			HasPrev:    offset > 0,
		},
	}, nil
}

// TrackUserActivity tracks user activity
func (uc *userUseCase) TrackUserActivity(ctx context.Context, userID uuid.UUID, activityType entities.ActivityType, description string, entityType string, entityID *uuid.UUID, metadata map[string]interface{}) error {
	activity := &entities.UserActivity{
		ID:          uuid.New(),
		UserID:      userID,
		Type:        activityType,
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
	_ = uc.TrackUserActivity(ctx, userID, entities.ActivityTypeProfileUpdate, "User preferences updated", "user_preferences", &preferences.ID, nil)

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

	// TODO: Store verification token in proper table
	// For now, just log the verification token
	fmt.Printf("Email verification token for %s: %s\n", user.Email, token)
	fmt.Printf("Verification link: http://localhost:3000/verify-email?token=%s\n", token)

	// Track activity
	_ = uc.TrackUserActivity(ctx, userID, entities.ActivityTypeProfileUpdate, "Email verification sent", "user", &user.ID, nil)

	return nil
}

// VerifyEmail verifies email with token
func (uc *userUseCase) VerifyEmail(ctx context.Context, token string) error {
	// For now, we'll implement a simple validation
	// In production, you would validate the token against a proper storage
	if token == "" {
		return fmt.Errorf("invalid verification token")
	}

	// For demo purposes, we'll accept any non-empty token
	// In production, you would:
	// 1. Validate token exists in verification table
	// 2. Check if token is not expired
	// 3. Get user ID from token
	// 4. Mark user email as verified

	// For now, we'll just return success
	// TODO: Implement proper token validation and email verification
	fmt.Printf("Email verification requested with token: %s\n", token)

	return nil
}

// SendPhoneVerification sends phone verification
func (uc *userUseCase) SendPhoneVerification(ctx context.Context, userID uuid.UUID, phone string) error {
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return entities.ErrUserNotFound
	}

	// Generate 6-digit OTP code
	code := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)

	// TODO: Store verification code in proper table
	// For now, just log the verification code
	fmt.Printf("Phone verification code for user %s (phone: %s): %s\n", userID, phone, code)

	// Track activity
	_ = uc.TrackUserActivity(ctx, userID, entities.ActivityTypeProfileUpdate, "Phone verification sent", "user", &userID, nil)

	return nil
}

// VerifyPhone verifies phone with code
func (uc *userUseCase) VerifyPhone(ctx context.Context, userID uuid.UUID, code string) error {
	// For now, we'll implement a simple validation
	// In production, you would validate the code against a proper storage
	if code == "" {
		return fmt.Errorf("invalid verification code")
	}

	// For demo purposes, we'll accept any non-empty code
	// In production, you would:
	// 1. Validate code exists in verification table
	// 2. Check if code is not expired
	// 3. Check if code belongs to this user
	// 4. Mark user phone as verified

	// For now, we'll just return success
	// TODO: Implement proper code validation and phone verification
	fmt.Printf("Phone verification requested for user %s with code: %s\n", userID, code)

	return nil
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

	// For now, we'll use a simple approach - store the reset token in user table
	// In production, you would create a proper password_reset_tokens table
	// TODO: Send password reset email
	// In production, you would send an email with the reset link
	// For now, we'll just log it
	fmt.Printf("Password reset token for %s: %s\n", user.Email, resetToken)
	fmt.Printf("Reset link: http://localhost:3000/reset-password?token=%s\n", resetToken)

	return nil
}

// ResetPassword resets user password using reset token
func (uc *userUseCase) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	// For now, we'll implement a simple validation
	// In production, you would validate the token against a proper storage
	if req.Token == "" {
		return fmt.Errorf("invalid reset token")
	}

	// For demo purposes, we'll accept any non-empty token
	// In production, you would:
	// 1. Validate token exists in password_reset_tokens table
	// 2. Check if token is not expired
	// 3. Get user ID from token

	// For now, we'll just return success
	// TODO: Implement proper token validation and password reset
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

	// For now, just log the verification request
	// TODO: Implement proper email verification sending
	fmt.Printf("Email verification resent for: %s\n", user.Email)

	return nil
}
