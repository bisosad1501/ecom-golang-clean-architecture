package usecases

import (
	"context"
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
	GetUsers(ctx context.Context, limit, offset int) ([]*UserResponse, error)
	DeactivateUser(ctx context.Context, userID uuid.UUID) error
	ActivateUser(ctx context.Context, userID uuid.UUID) error
}

type userUseCase struct {
	userRepo        repositories.UserRepository
	userProfileRepo repositories.UserProfileRepository
	passwordService services.PasswordService
	jwtSecret       string
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(
	userRepo repositories.UserRepository,
	userProfileRepo repositories.UserProfileRepository,
	passwordService services.PasswordService,
	jwtSecret string,
) UserUseCase {
	return &userUseCase{
		userRepo:        userRepo,
		userProfileRepo: userProfileRepo,
		passwordService: passwordService,
		jwtSecret:       jwtSecret,
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
	ID        uuid.UUID           `json:"id"`
	Email     string              `json:"email"`
	FirstName string              `json:"first_name"`
	LastName  string              `json:"last_name"`
	Phone     string              `json:"phone"`
	Role      entities.UserRole   `json:"role"`
	IsActive  bool                `json:"is_active"`
	Profile   *UserProfileResponse `json:"profile,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
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
func (uc *userUseCase) GetUsers(ctx context.Context, limit, offset int) ([]*UserResponse, error) {
	users, err := uc.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = uc.toUserResponse(user)
	}

	return responses, nil
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
