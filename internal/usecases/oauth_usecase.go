package usecases

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/infrastructure/config"
	"ecom-golang-clean-architecture/internal/infrastructure/oauth"
)

// JWTService defines JWT service interface
type JWTService interface {
	GenerateToken(userID, role string) (string, error)
}

// OAuthUseCase defines OAuth use case interface
type OAuthUseCase interface {
	GetGoogleAuthURL(ctx context.Context) (*OAuthURLResponse, error)
	GetFacebookAuthURL(ctx context.Context) (*OAuthURLResponse, error)
	HandleGoogleCallback(ctx context.Context, req *OAuthCallbackRequest) (*LoginResponse, error)
	HandleFacebookCallback(ctx context.Context, req *OAuthCallbackRequest) (*LoginResponse, error)
}

// OAuthURLResponse represents OAuth URL response
type OAuthURLResponse struct {
	URL   string `json:"url"`
	State string `json:"state"`
}

// OAuthCallbackRequest represents OAuth callback request
type OAuthCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}

type oauthUseCase struct {
	userRepo     repositories.UserRepository
	oauthService *oauth.Service
	jwtService   JWTService
}

// NewOAuthUseCase creates a new OAuth use case
func NewOAuthUseCase(
	userRepo repositories.UserRepository,
	oauthService *oauth.Service,
	jwtService JWTService,
) OAuthUseCase {
	return &oauthUseCase{
		userRepo:     userRepo,
		oauthService: oauthService,
		jwtService:   jwtService,
	}
}

// GetGoogleAuthURL generates Google OAuth URL
func (uc *oauthUseCase) GetGoogleAuthURL(ctx context.Context) (*OAuthURLResponse, error) {
	state, err := generateRandomState()
	if err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}

	url := uc.oauthService.GetGoogleAuthURL(state)

	return &OAuthURLResponse{
		URL:   url,
		State: state,
	}, nil
}

// GetFacebookAuthURL generates Facebook OAuth URL
func (uc *oauthUseCase) GetFacebookAuthURL(ctx context.Context) (*OAuthURLResponse, error) {
	state, err := generateRandomState()
	if err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}

	url := uc.oauthService.GetFacebookAuthURL(state)

	return &OAuthURLResponse{
		URL:   url,
		State: state,
	}, nil
}

// HandleGoogleCallback handles Google OAuth callback
func (uc *oauthUseCase) HandleGoogleCallback(ctx context.Context, req *OAuthCallbackRequest) (*LoginResponse, error) {
	// Exchange code for user info
	userInfo, err := uc.oauthService.ExchangeGoogleCode(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange Google code: %w", err)
	}

	// Find or create user
	user, err := uc.findOrCreateOAuthUser(ctx, userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// Generate JWT token
	token, err := uc.jwtService.GenerateToken(user.ID.String(), string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  uc.toUserResponse(user),
	}, nil
}

// HandleFacebookCallback handles Facebook OAuth callback
func (uc *oauthUseCase) HandleFacebookCallback(ctx context.Context, req *OAuthCallbackRequest) (*LoginResponse, error) {
	// Exchange code for user info
	userInfo, err := uc.oauthService.ExchangeFacebookCode(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange Facebook code: %w", err)
	}

	// Find or create user
	user, err := uc.findOrCreateOAuthUser(ctx, userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// Generate JWT token
	token, err := uc.jwtService.GenerateToken(user.ID.String(), string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  uc.toUserResponse(user),
	}, nil
}

// findOrCreateOAuthUser finds existing user or creates new one from OAuth info
func (uc *oauthUseCase) findOrCreateOAuthUser(ctx context.Context, userInfo *config.OAuthUserInfo) (*entities.User, error) {
	// Try to find user by email first
	existingUser, err := uc.userRepo.GetByEmail(ctx, userInfo.Email)
	if err == nil {
		// User exists, update OAuth info
		return uc.updateUserOAuthInfo(ctx, existingUser, userInfo)
	}

	// Try to find by OAuth provider ID
	var user *entities.User
	switch userInfo.Provider {
	case config.ProviderGoogle:
		user, err = uc.userRepo.GetByGoogleID(ctx, userInfo.ProviderID)
	case config.ProviderFacebook:
		user, err = uc.userRepo.GetByFacebookID(ctx, userInfo.ProviderID)
	}

	if err == nil {
		// User found by OAuth ID, update info
		return uc.updateUserOAuthInfo(ctx, user, userInfo)
	}

	// Create new user
	return uc.createOAuthUser(ctx, userInfo)
}

// updateUserOAuthInfo updates existing user with OAuth information
func (uc *oauthUseCase) updateUserOAuthInfo(ctx context.Context, user *entities.User, userInfo *config.OAuthUserInfo) (*entities.User, error) {
	// Update OAuth fields
	switch userInfo.Provider {
	case config.ProviderGoogle:
		user.GoogleID = userInfo.ProviderID
	case config.ProviderFacebook:
		user.FacebookID = userInfo.ProviderID
	}

	user.Avatar = userInfo.Picture
	user.IsOAuthUser = true
	user.EmailVerified = userInfo.Verified

	// Update user in database
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// createOAuthUser creates a new user from OAuth information
func (uc *oauthUseCase) createOAuthUser(ctx context.Context, userInfo *config.OAuthUserInfo) (*entities.User, error) {
	user := &entities.User{
		ID:            uuid.New(),
		Email:         userInfo.Email,
		FirstName:     userInfo.FirstName,
		LastName:      userInfo.LastName,
		Role:          entities.UserRoleCustomer,
		Status:        entities.UserStatusActive,
		IsActive:      true,
		Avatar:        userInfo.Picture,
		IsOAuthUser:   true,
		EmailVerified: userInfo.Verified,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Set OAuth provider ID
	switch userInfo.Provider {
	case config.ProviderGoogle:
		user.GoogleID = userInfo.ProviderID
	case config.ProviderFacebook:
		user.FacebookID = userInfo.ProviderID
	}

	// If name is not split, try to split it
	if user.FirstName == "" && user.LastName == "" && userInfo.Name != "" {
		user.FirstName = userInfo.Name
		user.LastName = ""
	}

	// Create user in database
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// toUserResponse converts user entity to response
func (uc *oauthUseCase) toUserResponse(user *entities.User) *UserResponse {
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

// generateRandomState generates a random state for OAuth
func generateRandomState() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
