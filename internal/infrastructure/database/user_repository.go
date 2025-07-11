package database

import (
	"context"
	"errors"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByGoogleID retrieves a user by Google ID
func (r *userRepository) GetByGoogleID(ctx context.Context, googleID string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("google_id = ?", googleID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByFacebookID retrieves a user by Facebook ID
func (r *userRepository) GetByFacebookID(ctx context.Context, facebookID string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("facebook_id = ?", facebookID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user by ID
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrUserNotFound
	}
	return nil
}

// List retrieves users with pagination
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error
	return users, err
}

// Count returns the total number of users
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.User{}).Count(&count).Error
	return count, err
}

// ExistsByEmail checks if a user exists with the given email
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("email = ?", email).
		Count(&count).Error
	return count > 0, err
}

// UpdatePassword updates user password
func (r *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	result := r.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("id = ?", userID).
		Update("password", hashedPassword)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrUserNotFound
	}
	return nil
}

// SetActive sets user active status
func (r *userRepository) SetActive(ctx context.Context, userID uuid.UUID, isActive bool) error {
	result := r.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("id = ?", userID).
		Update("is_active", isActive)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrUserNotFound
	}
	return nil
}

// CountUsers counts total number of users
func (r *userRepository) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.User{}).
		Count(&count).Error
	return count, err
}

// CountActiveUsers counts active users
func (r *userRepository) CountActiveUsers(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("is_active = ? AND status = ?", true, entities.UserStatusActive).
		Count(&count).Error
	return count, err
}

// GetUsersWithFilters gets users with filters
func (r *userRepository) GetUsersWithFilters(ctx context.Context, filters repositories.UserFilters) ([]*entities.User, error) {
	query := r.db.WithContext(ctx).Model(&entities.User{})

	// Apply filters
	query = r.applyUserFilters(query, filters)

	// Apply sorting
	if filters.SortBy != "" {
		order := filters.SortBy
		if filters.SortOrder == "desc" {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var users []*entities.User
	err := query.Find(&users).Error
	return users, err
}

// CountUsersWithFilters counts users with filters
func (r *userRepository) CountUsersWithFilters(ctx context.Context, filters repositories.UserFilters) (int64, error) {
	query := r.db.WithContext(ctx).Model(&entities.User{})
	query = r.applyUserFilters(query, filters)

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// UpdateLastLogin updates user's last login timestamp
func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entities.User{}).
		Where("id = ?", userID).
		Update("last_login_at", now).Error
}

// UpdateLastActivity updates user's last activity timestamp
func (r *userRepository) UpdateLastActivity(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entities.User{}).
		Where("id = ?", userID).
		Update("last_activity_at", now).Error
}

// GetUsersByRole gets users by role
func (r *userRepository) GetUsersByRole(ctx context.Context, role entities.UserRole, limit, offset int) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.WithContext(ctx).
		Where("role = ?", role).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	return users, err
}

// GetUsersByStatus gets users by status
func (r *userRepository) GetUsersByStatus(ctx context.Context, status entities.UserStatus, limit, offset int) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	return users, err
}

// GetHighValueCustomers gets high value customers
func (r *userRepository) GetHighValueCustomers(ctx context.Context, limit int) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.WithContext(ctx).
		Where("total_spent > ? OR total_orders > ?", 1000, 10).
		Order("total_spent DESC").
		Limit(limit).
		Find(&users).Error
	return users, err
}

// GetRecentlyRegistered gets recently registered users
func (r *userRepository) GetRecentlyRegistered(ctx context.Context, limit int) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&users).Error
	return users, err
}

// applyUserFilters applies filters to the query
func (r *userRepository) applyUserFilters(query *gorm.DB, filters repositories.UserFilters) *gorm.DB {
	if filters.Role != nil {
		query = query.Where("role = ?", *filters.Role)
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}
	if filters.EmailVerified != nil {
		query = query.Where("email_verified = ?", *filters.EmailVerified)
	}
	if filters.PhoneVerified != nil {
		query = query.Where("phone_verified = ?", *filters.PhoneVerified)
	}
	if filters.TwoFactorEnabled != nil {
		query = query.Where("two_factor_enabled = ?", *filters.TwoFactorEnabled)
	}
	if filters.MembershipTier != "" {
		query = query.Where("membership_tier = ?", filters.MembershipTier)
	}
	if filters.MinTotalSpent != nil {
		query = query.Where("total_spent >= ?", *filters.MinTotalSpent)
	}
	if filters.MaxTotalSpent != nil {
		query = query.Where("total_spent <= ?", *filters.MaxTotalSpent)
	}
	if filters.MinTotalOrders != nil {
		query = query.Where("total_orders >= ?", *filters.MinTotalOrders)
	}
	if filters.MaxTotalOrders != nil {
		query = query.Where("total_orders <= ?", *filters.MaxTotalOrders)
	}
	if filters.CreatedFrom != nil {
		query = query.Where("created_at >= ?", *filters.CreatedFrom)
	}
	if filters.CreatedTo != nil {
		query = query.Where("created_at <= ?", *filters.CreatedTo)
	}
	if filters.LastLoginFrom != nil {
		query = query.Where("last_login_at >= ?", *filters.LastLoginFrom)
	}
	if filters.LastLoginTo != nil {
		query = query.Where("last_login_at <= ?", *filters.LastLoginTo)
	}
	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	return query
}

type userProfileRepository struct {
	db *gorm.DB
}

// NewUserProfileRepository creates a new user profile repository
func NewUserProfileRepository(db *gorm.DB) repositories.UserProfileRepository {
	return &userProfileRepository{db: db}
}

// Create creates a new user profile
func (r *userProfileRepository) Create(ctx context.Context, profile *entities.UserProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

// GetByUserID retrieves a user profile by user ID
func (r *userProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.UserProfile, error) {
	var profile entities.UserProfile
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}
	return &profile, nil
}

// Update updates an existing user profile
func (r *userProfileRepository) Update(ctx context.Context, profile *entities.UserProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

// Delete deletes a user profile by user ID
func (r *userProfileRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&entities.UserProfile{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrNotFound
	}
	return nil
}
