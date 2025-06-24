package repositories

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// CouponRepository defines the interface for coupon data access
type CouponRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, coupon *entities.Coupon) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Coupon, error)
	GetByCode(ctx context.Context, code string) (*entities.Coupon, error)
	Update(ctx context.Context, coupon *entities.Coupon) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	List(ctx context.Context, limit, offset int) ([]*entities.Coupon, error)
	GetActiveCoupons(ctx context.Context) ([]*entities.Coupon, error)
	GetUserCoupons(ctx context.Context, userID uuid.UUID) ([]*entities.Coupon, error)

	// Validation and usage
	ValidateCoupon(ctx context.Context, code string, userID uuid.UUID) (*entities.Coupon, error)
	IncrementUsage(ctx context.Context, couponID uuid.UUID) error
	RecordUsage(ctx context.Context, usage *entities.CouponUsage) error

	// Usage tracking
	GetUsageHistory(ctx context.Context, couponID uuid.UUID, limit, offset int) ([]*entities.CouponUsage, error)
	GetUserUsageCount(ctx context.Context, couponID, userID uuid.UUID) (int, error)

	// Maintenance
	ExpireCoupons(ctx context.Context) error
}

// PromotionRepository defines the interface for promotion data access
type PromotionRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, promotion *entities.Promotion) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Promotion, error)
	Update(ctx context.Context, promotion *entities.Promotion) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetActivePromotions(ctx context.Context) ([]*entities.Promotion, error)
	GetFeaturedPromotions(ctx context.Context, limit int) ([]*entities.Promotion, error)
	GetPromotionsForProduct(ctx context.Context, productID uuid.UUID) ([]*entities.Promotion, error)
}

// LoyaltyRepository defines the interface for loyalty program data access
type LoyaltyRepository interface {
	// Points management
	GetUserPoints(ctx context.Context, userID uuid.UUID) (*entities.UserLoyaltyPoints, error)
	AddPoints(ctx context.Context, userID uuid.UUID, points int, reason string) error
	RedeemPoints(ctx context.Context, userID uuid.UUID, points int) error

	// Program management
	GetLoyaltyProgram(ctx context.Context) (*entities.LoyaltyProgram, error)
}
