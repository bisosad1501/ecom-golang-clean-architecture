package database

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type couponRepository struct {
	db *gorm.DB
}

// NewCouponRepository creates a new coupon repository
func NewCouponRepository(db *gorm.DB) repositories.CouponRepository {
	return &couponRepository{db: db}
}

// Create creates a new coupon
func (r *couponRepository) Create(ctx context.Context, coupon *entities.Coupon) error {
	return r.db.WithContext(ctx).Create(coupon).Error
}

// GetByID retrieves a coupon by ID
func (r *couponRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Coupon, error) {
	var coupon entities.Coupon
	err := r.db.WithContext(ctx).
		Preload("ApplicableCategories").
		Preload("ApplicableProducts").
		Preload("ApplicableUsers").
		Where("id = ?", id).
		First(&coupon).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrCouponNotFound
		}
		return nil, err
	}
	return &coupon, nil
}

// GetByCode retrieves a coupon by code
func (r *couponRepository) GetByCode(ctx context.Context, code string) (*entities.Coupon, error) {
	var coupon entities.Coupon
	err := r.db.WithContext(ctx).
		Preload("ApplicableCategories").
		Preload("ApplicableProducts").
		Preload("ApplicableUsers").
		Where("code = ?", code).
		First(&coupon).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrCouponNotFound
		}
		return nil, err
	}
	return &coupon, nil
}

// Update updates an existing coupon
func (r *couponRepository) Update(ctx context.Context, coupon *entities.Coupon) error {
	return r.db.WithContext(ctx).Save(coupon).Error
}

// Delete deletes a coupon by ID
func (r *couponRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.Coupon{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrCouponNotFound
	}
	return nil
}

// List retrieves coupons with pagination
func (r *couponRepository) List(ctx context.Context, limit, offset int) ([]*entities.Coupon, error) {
	var coupons []*entities.Coupon
	err := r.db.WithContext(ctx).
		Preload("ApplicableCategories").
		Preload("ApplicableProducts").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&coupons).Error
	return coupons, err
}

// GetActiveCoupons retrieves active coupons
func (r *couponRepository) GetActiveCoupons(ctx context.Context) ([]*entities.Coupon, error) {
	var coupons []*entities.Coupon
	now := time.Now()
	err := r.db.WithContext(ctx).
		Preload("ApplicableCategories").
		Preload("ApplicableProducts").
		Where("status = ? AND (starts_at IS NULL OR starts_at <= ?) AND (expires_at IS NULL OR expires_at > ?)",
			entities.CouponStatusActive, now, now).
		Find(&coupons).Error
	return coupons, err
}

// GetUserCoupons retrieves coupons applicable to a user
func (r *couponRepository) GetUserCoupons(ctx context.Context, userID uuid.UUID) ([]*entities.Coupon, error) {
	var coupons []*entities.Coupon
	now := time.Now()
	
	// Get public coupons and user-specific coupons
	err := r.db.WithContext(ctx).
		Preload("ApplicableCategories").
		Preload("ApplicableProducts").
		Where(`status = ? AND (starts_at IS NULL OR starts_at <= ?) AND (expires_at IS NULL OR expires_at > ?) 
			   AND (applicability != ? OR id IN (
				   SELECT coupon_id FROM coupon_users WHERE user_id = ?
			   ))`,
			entities.CouponStatusActive, now, now, entities.CouponApplicabilityUsers, userID).
		Find(&coupons).Error
	return coupons, err
}

// ValidateCoupon validates if a coupon can be used
func (r *couponRepository) ValidateCoupon(ctx context.Context, code string, userID uuid.UUID) (*entities.Coupon, error) {
	coupon, err := r.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if !coupon.IsValid() {
		return nil, entities.ErrCouponInvalid
	}

	if !coupon.CanBeUsedBy(userID) {
		return nil, entities.ErrCouponNotApplicable
	}

	// Check user usage limit
	if coupon.UsageLimitPerUser != nil {
		var usageCount int64
		r.db.WithContext(ctx).
			Model(&entities.CouponUsage{}).
			Where("coupon_id = ? AND user_id = ?", coupon.ID, userID).
			Count(&usageCount)
		
		if int(usageCount) >= *coupon.UsageLimitPerUser {
			return nil, entities.ErrCouponUsageLimitExceeded
		}
	}

	return coupon, nil
}

// IncrementUsage increments the usage count of a coupon
func (r *couponRepository) IncrementUsage(ctx context.Context, couponID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.Coupon{}).
		Where("id = ?", couponID).
		UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}

// RecordUsage records coupon usage
func (r *couponRepository) RecordUsage(ctx context.Context, usage *entities.CouponUsage) error {
	return r.db.WithContext(ctx).Create(usage).Error
}

// GetUsageHistory gets coupon usage history
func (r *couponRepository) GetUsageHistory(ctx context.Context, couponID uuid.UUID, limit, offset int) ([]*entities.CouponUsage, error) {
	var usage []*entities.CouponUsage
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Order").
		Where("coupon_id = ?", couponID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&usage).Error
	return usage, err
}

// GetUserUsageCount gets user's usage count for a coupon
func (r *couponRepository) GetUserUsageCount(ctx context.Context, couponID, userID uuid.UUID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.CouponUsage{}).
		Where("coupon_id = ? AND user_id = ?", couponID, userID).
		Count(&count).Error
	return int(count), err
}

// ExpireCoupons marks expired coupons as expired
func (r *couponRepository) ExpireCoupons(ctx context.Context) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entities.Coupon{}).
		Where("expires_at < ? AND status = ?", now, entities.CouponStatusActive).
		Update("status", entities.CouponStatusExpired).Error
}

type promotionRepository struct {
	db *gorm.DB
}

// NewPromotionRepository creates a new promotion repository
func NewPromotionRepository(db *gorm.DB) repositories.PromotionRepository {
	return &promotionRepository{db: db}
}

// Create creates a new promotion
func (r *promotionRepository) Create(ctx context.Context, promotion *entities.Promotion) error {
	return r.db.WithContext(ctx).Create(promotion).Error
}

// GetByID retrieves a promotion by ID
func (r *promotionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Promotion, error) {
	var promotion entities.Promotion
	err := r.db.WithContext(ctx).
		Preload("ApplicableCategories").
		Preload("ApplicableProducts").
		Where("id = ?", id).
		First(&promotion).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrPromotionNotFound
		}
		return nil, err
	}
	return &promotion, nil
}

// Update updates an existing promotion
func (r *promotionRepository) Update(ctx context.Context, promotion *entities.Promotion) error {
	return r.db.WithContext(ctx).Save(promotion).Error
}

// Delete deletes a promotion by ID
func (r *promotionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.Promotion{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrPromotionNotFound
	}
	return nil
}

// GetActivePromotions retrieves active promotions
func (r *promotionRepository) GetActivePromotions(ctx context.Context) ([]*entities.Promotion, error) {
	var promotions []*entities.Promotion
	now := time.Now()
	err := r.db.WithContext(ctx).
		Preload("ApplicableCategories").
		Preload("ApplicableProducts").
		Where("status = ? AND starts_at <= ? AND ends_at > ?",
			entities.CouponStatusActive, now, now).
		Find(&promotions).Error
	return promotions, err
}

// GetFeaturedPromotions retrieves featured promotions
func (r *promotionRepository) GetFeaturedPromotions(ctx context.Context, limit int) ([]*entities.Promotion, error) {
	var promotions []*entities.Promotion
	now := time.Now()
	err := r.db.WithContext(ctx).
		Preload("ApplicableCategories").
		Preload("ApplicableProducts").
		Where("status = ? AND is_featured = ? AND starts_at <= ? AND ends_at > ?",
			entities.CouponStatusActive, true, now, now).
		Limit(limit).
		Order("created_at DESC").
		Find(&promotions).Error
	return promotions, err
}

// GetPromotionsForProduct retrieves promotions applicable to a product
func (r *promotionRepository) GetPromotionsForProduct(ctx context.Context, productID uuid.UUID) ([]*entities.Promotion, error) {
	var promotions []*entities.Promotion
	now := time.Now()
	
	// Get promotions that apply to all products or specifically to this product
	err := r.db.WithContext(ctx).
		Preload("ApplicableCategories").
		Preload("ApplicableProducts").
		Where(`status = ? AND starts_at <= ? AND ends_at > ? 
			   AND (id NOT IN (SELECT promotion_id FROM promotion_products) 
			   OR id IN (SELECT promotion_id FROM promotion_products WHERE product_id = ?))`,
			entities.CouponStatusActive, now, now, productID).
		Find(&promotions).Error
	return promotions, err
}

type loyaltyRepository struct {
	db *gorm.DB
}

// NewLoyaltyRepository creates a new loyalty repository
func NewLoyaltyRepository(db *gorm.DB) repositories.LoyaltyRepository {
	return &loyaltyRepository{db: db}
}

// GetUserPoints retrieves user's loyalty points
func (r *loyaltyRepository) GetUserPoints(ctx context.Context, userID uuid.UUID) (*entities.UserLoyaltyPoints, error) {
	var points entities.UserLoyaltyPoints
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&points).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new points record
			points = entities.UserLoyaltyPoints{
				ID:              uuid.New(),
				UserID:          userID,
				TotalPoints:     0,
				AvailablePoints: 0,
				UsedPoints:      0,
				ExpiredPoints:   0,
				TierLevel:       "bronze",
				UpdatedAt:       time.Now(),
			}
			if err := r.db.WithContext(ctx).Create(&points).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &points, nil
}

// AddPoints adds points to user's account
func (r *loyaltyRepository) AddPoints(ctx context.Context, userID uuid.UUID, points int, reason string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update user points
		result := tx.Model(&entities.UserLoyaltyPoints{}).
			Where("user_id = ?", userID).
			Updates(map[string]interface{}{
				"total_points":     gorm.Expr("total_points + ?", points),
				"available_points": gorm.Expr("available_points + ?", points),
				"updated_at":       time.Now(),
			})
		
		if result.Error != nil {
			return result.Error
		}
		
		if result.RowsAffected == 0 {
			// Create new record if doesn't exist
			userPoints := &entities.UserLoyaltyPoints{
				ID:              uuid.New(),
				UserID:          userID,
				TotalPoints:     points,
				AvailablePoints: points,
				UsedPoints:      0,
				ExpiredPoints:   0,
				TierLevel:       "bronze",
				UpdatedAt:       time.Now(),
			}
			return tx.Create(userPoints).Error
		}
		
		return nil
	})
}

// RedeemPoints redeems points from user's account
func (r *loyaltyRepository) RedeemPoints(ctx context.Context, userID uuid.UUID, points int) error {
	return r.db.WithContext(ctx).
		Model(&entities.UserLoyaltyPoints{}).
		Where("user_id = ? AND available_points >= ?", userID, points).
		Updates(map[string]interface{}{
			"available_points": gorm.Expr("available_points - ?", points),
			"used_points":      gorm.Expr("used_points + ?", points),
			"updated_at":       time.Now(),
		}).Error
}

// GetLoyaltyProgram retrieves the active loyalty program
func (r *loyaltyRepository) GetLoyaltyProgram(ctx context.Context) (*entities.LoyaltyProgram, error) {
	var program entities.LoyaltyProgram
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		First(&program).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrLoyaltyProgramNotFound
		}
		return nil, err
	}
	return &program, nil
}
