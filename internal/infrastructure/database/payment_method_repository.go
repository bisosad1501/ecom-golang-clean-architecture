package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
)

type paymentMethodRepository struct {
	db *gorm.DB
}

// NewPaymentMethodRepository creates a new payment method repository
func NewPaymentMethodRepository(db *gorm.DB) repositories.PaymentMethodRepository {
	return &paymentMethodRepository{db: db}
}

// Create creates a new payment method
func (r *paymentMethodRepository) Create(ctx context.Context, paymentMethod *entities.PaymentMethodEntity) error {
	return r.db.WithContext(ctx).Create(paymentMethod).Error
}

// GetByID retrieves a payment method by ID
func (r *paymentMethodRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.PaymentMethodEntity, error) {
	var paymentMethod entities.PaymentMethodEntity
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&paymentMethod).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrPaymentMethodNotFound
		}
		return nil, err
	}
	return &paymentMethod, nil
}

// GetByUserID retrieves all payment methods for a user
func (r *paymentMethodRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.PaymentMethodEntity, error) {
	var paymentMethods []*entities.PaymentMethodEntity
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&paymentMethods).Error
	return paymentMethods, err
}

// GetActiveByUserID retrieves all active payment methods for a user
func (r *paymentMethodRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.PaymentMethodEntity, error) {
	var paymentMethods []*entities.PaymentMethodEntity
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("is_default DESC, created_at DESC").
		Find(&paymentMethods).Error
	return paymentMethods, err
}

// GetDefaultByUserID retrieves the default payment method for a user
func (r *paymentMethodRepository) GetDefaultByUserID(ctx context.Context, userID uuid.UUID) (*entities.PaymentMethodEntity, error) {
	var paymentMethod entities.PaymentMethodEntity
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_default = ? AND is_active = ?", userID, true, true).
		First(&paymentMethod).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrPaymentMethodNotFound
		}
		return nil, err
	}
	return &paymentMethod, nil
}

// GetByGatewayToken retrieves a payment method by gateway token
func (r *paymentMethodRepository) GetByGatewayToken(ctx context.Context, gatewayToken string) (*entities.PaymentMethodEntity, error) {
	var paymentMethod entities.PaymentMethodEntity
	err := r.db.WithContext(ctx).Where("gateway_token = ?", gatewayToken).First(&paymentMethod).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrPaymentMethodNotFound
		}
		return nil, err
	}
	return &paymentMethod, nil
}

// GetByFingerprint retrieves payment methods by fingerprint (to prevent duplicates)
func (r *paymentMethodRepository) GetByFingerprint(ctx context.Context, userID uuid.UUID, fingerprint string) (*entities.PaymentMethodEntity, error) {
	var paymentMethod entities.PaymentMethodEntity
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND fingerprint = ? AND is_active = ?", userID, fingerprint, true).
		First(&paymentMethod).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrPaymentMethodNotFound
		}
		return nil, err
	}
	return &paymentMethod, nil
}

// Update updates an existing payment method
func (r *paymentMethodRepository) Update(ctx context.Context, paymentMethod *entities.PaymentMethodEntity) error {
	return r.db.WithContext(ctx).Save(paymentMethod).Error
}

// Delete deletes a payment method by ID
func (r *paymentMethodRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.PaymentMethodEntity{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrPaymentMethodNotFound
	}
	return nil
}

// SetAsDefault sets a payment method as default and unsets others
func (r *paymentMethodRepository) SetAsDefault(ctx context.Context, userID, paymentMethodID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First, unset all default payment methods for the user
		if err := tx.Model(&entities.PaymentMethodEntity{}).
			Where("user_id = ?", userID).
			Update("is_default", false).Error; err != nil {
			return err
		}

		// Then set the specified payment method as default
		result := tx.Model(&entities.PaymentMethodEntity{}).
			Where("id = ? AND user_id = ?", paymentMethodID, userID).
			Update("is_default", true)
		
		if result.Error != nil {
			return result.Error
		}
		
		if result.RowsAffected == 0 {
			return entities.ErrPaymentMethodNotFound
		}
		
		return nil
	})
}

// UnsetDefault removes default status from all payment methods for a user
func (r *paymentMethodRepository) UnsetDefault(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.PaymentMethodEntity{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
}

// Deactivate deactivates a payment method
func (r *paymentMethodRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Model(&entities.PaymentMethodEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active":  false,
			"is_default": false,
			"updated_at": time.Now(),
		})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return entities.ErrPaymentMethodNotFound
	}
	
	return nil
}

// Count returns the total number of payment methods for a user
func (r *paymentMethodRepository) Count(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.PaymentMethodEntity{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Count(&count).Error
	return count, err
}

// GetExpiredCards retrieves expired card payment methods
func (r *paymentMethodRepository) GetExpiredCards(ctx context.Context, limit, offset int) ([]*entities.PaymentMethodEntity, error) {
	var paymentMethods []*entities.PaymentMethodEntity
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())
	
	err := r.db.WithContext(ctx).
		Where("(type = ? OR type = ?) AND is_active = ? AND ((expiry_year < ?) OR (expiry_year = ? AND expiry_month < ?))",
			entities.PaymentMethodCreditCard, entities.PaymentMethodDebitCard, true,
			currentYear, currentYear, currentMonth).
		Limit(limit).
		Offset(offset).
		Find(&paymentMethods).Error
	
	return paymentMethods, err
}

// CleanupInactive removes inactive payment methods older than specified days
func (r *paymentMethodRepository) CleanupInactive(ctx context.Context, daysOld int) error {
	cutoffDate := time.Now().AddDate(0, 0, -daysOld)
	
	result := r.db.WithContext(ctx).
		Where("is_active = ? AND updated_at < ?", false, cutoffDate).
		Delete(&entities.PaymentMethodEntity{})
	
	return result.Error
}
