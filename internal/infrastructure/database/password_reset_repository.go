package database

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type passwordResetRepository struct {
	db *gorm.DB
}

// NewPasswordResetRepository creates a new password reset repository
func NewPasswordResetRepository(db *gorm.DB) repositories.PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

// Create creates a new password reset record
func (r *passwordResetRepository) Create(ctx context.Context, reset *entities.PasswordReset) error {
	return r.db.WithContext(ctx).Create(reset).Error
}

// GetByToken retrieves a password reset record by token
func (r *passwordResetRepository) GetByToken(ctx context.Context, token string) (*entities.PasswordReset, error) {
	var reset entities.PasswordReset
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&reset).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return &reset, nil
}

// MarkAsUsed marks a password reset token as used
func (r *passwordResetRepository) MarkAsUsed(ctx context.Context, token string) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&entities.PasswordReset{}).
		Where("token = ?", token).
		Update("used_at", now)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrUserNotFound
	}
	return nil
}

// Delete deletes a password reset record by ID
func (r *passwordResetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.PasswordReset{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrUserNotFound
	}
	return nil
}

// DeleteExpired deletes expired password reset records
func (r *passwordResetRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Delete(&entities.PasswordReset{}, "expires_at < ?", time.Now()).Error
}

// DeleteByUserID deletes all password reset records for a user
func (r *passwordResetRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&entities.PasswordReset{}, "user_id = ?", userID).Error
}

// IsTokenValid checks if a token is valid (exists, not expired, not used)
func (r *passwordResetRepository) IsTokenValid(ctx context.Context, token string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.PasswordReset{}).
		Where("token = ? AND expires_at > ? AND used_at IS NULL", token, time.Now()).
		Count(&count).Error
	
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}
