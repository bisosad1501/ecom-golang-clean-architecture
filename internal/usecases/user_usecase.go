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
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
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

	return &LoginResponse{
		User:  uc.toUserResponse(user),
		Token: token,
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

	verification := &entities.UserVerification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      "email",
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours expiry
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.userVerificationRepo.Create(ctx, verification); err != nil {
		return err
	}

	// TODO: Send email with verification link
	// In a real implementation, you would send an email here

	// Track activity
	_ = uc.TrackUserActivity(ctx, userID, entities.ActivityTypeProfileUpdate, "Email verification sent", "user_verification", &verification.ID, nil)

	return nil
}

// VerifyEmail verifies email with token
func (uc *userUseCase) VerifyEmail(ctx context.Context, token string) error {
	verification, err := uc.userVerificationRepo.GetByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid or expired verification token")
	}

	if verification.IsExpired() {
		return fmt.Errorf("verification token has expired")
	}

	// Mark verification as verified
	if err := uc.userVerificationRepo.MarkAsVerified(ctx, verification.ID); err != nil {
		return err
	}

	// Update user email verification status
	user, err := uc.userRepo.GetByID(ctx, verification.UserID)
	if err != nil {
		return err
	}

	user.EmailVerified = true
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Track activity
	_ = uc.TrackUserActivity(ctx, verification.UserID, entities.ActivityTypeProfileUpdate, "Email verified", "user", &user.ID, nil)

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

	verification := &entities.UserVerification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      "phone",
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute), // 10 minutes expiry
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.userVerificationRepo.Create(ctx, verification); err != nil {
		return err
	}

	// TODO: Send SMS with verification code
	// In a real implementation, you would send an SMS here

	// Track activity
	_ = uc.TrackUserActivity(ctx, userID, entities.ActivityTypeProfileUpdate, "Phone verification sent", "user_verification", &verification.ID, nil)

	return nil
}

// VerifyPhone verifies phone with code
func (uc *userUseCase) VerifyPhone(ctx context.Context, userID uuid.UUID, code string) error {
	verification, err := uc.userVerificationRepo.GetByCode(ctx, code, "phone")
	if err != nil {
		return fmt.Errorf("invalid verification code")
	}

	if verification.UserID != userID {
		return fmt.Errorf("verification code does not belong to this user")
	}

	if verification.IsExpired() {
		return fmt.Errorf("verification code has expired")
	}

	if !verification.CanAttempt() {
		return fmt.Errorf("maximum verification attempts exceeded")
	}

	// Mark verification as verified
	if err := uc.userVerificationRepo.MarkAsVerified(ctx, verification.ID); err != nil {
		return err
	}

	// Update user phone verification status
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.PhoneVerified = true
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Track activity
	_ = uc.TrackUserActivity(ctx, userID, entities.ActivityTypeProfileUpdate, "Phone verified", "user", &user.ID, nil)

	return nil
}

// GetVerificationStatus gets verification status
func (uc *userUseCase) GetVerificationStatus(ctx context.Context, userID uuid.UUID) (*VerificationStatusResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, entities.ErrUserNotFound
	}

	// Get active verifications
	activeVerifications, err := uc.userVerificationRepo.GetActiveVerifications(ctx, userID)
	if err != nil {
		activeVerifications = []*entities.UserVerification{}
	}

	response := &VerificationStatusResponse{
		UserID:        userID,
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
	}

	// Check for pending verifications
	for _, verification := range activeVerifications {
		switch verification.Type {
		case "email":
			response.PendingEmailVerification = true
			response.LastEmailVerificationSent = &verification.CreatedAt
		case "phone":
			response.PendingPhoneVerification = true
			response.LastPhoneVerificationSent = &verification.CreatedAt
		}
	}

	return response, nil
}
