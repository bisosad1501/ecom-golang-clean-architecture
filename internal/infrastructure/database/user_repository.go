package database

import (
	"context"

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
