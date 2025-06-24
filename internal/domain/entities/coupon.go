package entities

import (
	"time"

	"github.com/google/uuid"
)

// CouponType represents the type of coupon
type CouponType string

const (
	CouponTypePercentage CouponType = "percentage"
	CouponTypeFixed      CouponType = "fixed"
	CouponTypeFreeShipping CouponType = "free_shipping"
	CouponTypeBuyXGetY   CouponType = "buy_x_get_y"
)

// CouponStatus represents the status of a coupon
type CouponStatus string

const (
	CouponStatusActive   CouponStatus = "active"
	CouponStatusInactive CouponStatus = "inactive"
	CouponStatusExpired  CouponStatus = "expired"
	CouponStatusUsedUp   CouponStatus = "used_up"
)

// CouponApplicability represents what the coupon applies to
type CouponApplicability string

const (
	CouponApplicabilityAll        CouponApplicability = "all"
	CouponApplicabilityCategories CouponApplicability = "categories"
	CouponApplicabilityProducts   CouponApplicability = "products"
	CouponApplicabilityUsers      CouponApplicability = "users"
)

// Coupon represents a discount coupon
type Coupon struct {
	ID          uuid.UUID           `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code        string              `json:"code" gorm:"uniqueIndex;not null" validate:"required"`
	Name        string              `json:"name" gorm:"not null" validate:"required"`
	Description string              `json:"description"`
	Type        CouponType          `json:"type" gorm:"not null" validate:"required"`
	Value       float64             `json:"value" gorm:"not null" validate:"required,min=0"`
	MaxDiscount *float64            `json:"max_discount"` // For percentage coupons
	MinOrderAmount *float64         `json:"min_order_amount"`
	
	// Usage limits
	UsageLimit      *int `json:"usage_limit"`      // Total usage limit
	UsageLimitPerUser *int `json:"usage_limit_per_user"` // Per user limit
	UsedCount       int  `json:"used_count" gorm:"default:0"`
	
	// Applicability
	Applicability   CouponApplicability `json:"applicability" gorm:"default:'all'"`
	ApplicableCategories []Category     `json:"applicable_categories,omitempty" gorm:"many2many:coupon_categories;"`
	ApplicableProducts   []Product      `json:"applicable_products,omitempty" gorm:"many2many:coupon_products;"`
	ApplicableUsers      []User         `json:"applicable_users,omitempty" gorm:"many2many:coupon_users;"`
	
	// Buy X Get Y specific fields
	BuyQuantity *int     `json:"buy_quantity"`  // For buy_x_get_y type
	GetQuantity *int     `json:"get_quantity"`  // For buy_x_get_y type
	GetProductID *uuid.UUID `json:"get_product_id"` // Specific product to get free
	
	// Validity
	StartsAt  *time.Time    `json:"starts_at"`
	ExpiresAt *time.Time    `json:"expires_at"`
	Status    CouponStatus  `json:"status" gorm:"default:'active'"`
	
	// Metadata
	IsFirstTimeUser bool      `json:"is_first_time_user" gorm:"default:false"`
	IsPublic        bool      `json:"is_public" gorm:"default:true"`
	CreatedBy       uuid.UUID `json:"created_by" gorm:"type:uuid"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	Usage []CouponUsage `json:"usage,omitempty" gorm:"foreignKey:CouponID"`
}

// TableName returns the table name for Coupon entity
func (Coupon) TableName() string {
	return "coupons"
}

// IsValid checks if the coupon is valid for use
func (c *Coupon) IsValid() bool {
	now := time.Now()
	
	// Check status
	if c.Status != CouponStatusActive {
		return false
	}
	
	// Check start date
	if c.StartsAt != nil && now.Before(*c.StartsAt) {
		return false
	}
	
	// Check expiry date
	if c.ExpiresAt != nil && now.After(*c.ExpiresAt) {
		return false
	}
	
	// Check usage limit
	if c.UsageLimit != nil && c.UsedCount >= *c.UsageLimit {
		return false
	}
	
	return true
}

// CanBeUsedBy checks if the coupon can be used by a specific user
func (c *Coupon) CanBeUsedBy(userID uuid.UUID) bool {
	if !c.IsValid() {
		return false
	}
	
	// Check if coupon is restricted to specific users
	if c.Applicability == CouponApplicabilityUsers {
		for _, user := range c.ApplicableUsers {
			if user.ID == userID {
				return true
			}
		}
		return false
	}
	
	return true
}

// CalculateDiscount calculates the discount amount for given order total
func (c *Coupon) CalculateDiscount(orderTotal float64) float64 {
	if !c.IsValid() {
		return 0
	}
	
	// Check minimum order amount
	if c.MinOrderAmount != nil && orderTotal < *c.MinOrderAmount {
		return 0
	}
	
	switch c.Type {
	case CouponTypePercentage:
		discount := orderTotal * (c.Value / 100)
		if c.MaxDiscount != nil && discount > *c.MaxDiscount {
			return *c.MaxDiscount
		}
		return discount
		
	case CouponTypeFixed:
		if c.Value > orderTotal {
			return orderTotal
		}
		return c.Value
		
	case CouponTypeFreeShipping:
		// This should be handled separately in shipping calculation
		return 0
		
	default:
		return 0
	}
}

// CouponUsage represents the usage of a coupon
type CouponUsage struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CouponID     uuid.UUID `json:"coupon_id" gorm:"type:uuid;not null;index"`
	Coupon       Coupon    `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	User         User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OrderID      uuid.UUID `json:"order_id" gorm:"type:uuid;not null;index"`
	Order        Order     `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	DiscountAmount float64 `json:"discount_amount" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for CouponUsage entity
func (CouponUsage) TableName() string {
	return "coupon_usage"
}

// Promotion represents a promotional campaign
type Promotion struct {
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string          `json:"name" gorm:"not null" validate:"required"`
	Description string          `json:"description"`
	Type        string          `json:"type" gorm:"not null"` // flash_sale, seasonal, clearance, etc.
	
	// Discount settings
	DiscountType       CouponType `json:"discount_type" gorm:"not null"`
	DiscountValue      float64    `json:"discount_value" gorm:"not null"`
	MaxDiscountAmount  *float64   `json:"max_discount_amount"`
	MinOrderAmount     *float64   `json:"min_order_amount"`
	
	// Applicability
	ApplicableCategories []Category `json:"applicable_categories,omitempty" gorm:"many2many:promotion_categories;"`
	ApplicableProducts   []Product  `json:"applicable_products,omitempty" gorm:"many2many:promotion_products;"`
	
	// Validity
	StartsAt  time.Time     `json:"starts_at" gorm:"not null"`
	EndsAt    time.Time     `json:"ends_at" gorm:"not null"`
	Status    CouponStatus  `json:"status" gorm:"default:'active'"`
	
	// Display settings
	BannerImage   string `json:"banner_image"`
	BannerText    string `json:"banner_text"`
	IsPublic      bool   `json:"is_public" gorm:"default:true"`
	IsFeatured    bool   `json:"is_featured" gorm:"default:false"`
	
	// Metadata
	CreatedBy uuid.UUID `json:"created_by" gorm:"type:uuid"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Promotion entity
func (Promotion) TableName() string {
	return "promotions"
}

// IsActive checks if the promotion is currently active
func (p *Promotion) IsActive() bool {
	now := time.Now()
	return p.Status == CouponStatusActive && 
		   now.After(p.StartsAt) && 
		   now.Before(p.EndsAt)
}

// CalculatePromotionDiscount calculates discount for a promotion
func (p *Promotion) CalculatePromotionDiscount(amount float64) float64 {
	if !p.IsActive() {
		return 0
	}
	
	// Check minimum order amount
	if p.MinOrderAmount != nil && amount < *p.MinOrderAmount {
		return 0
	}
	
	switch p.DiscountType {
	case CouponTypePercentage:
		discount := amount * (p.DiscountValue / 100)
		if p.MaxDiscountAmount != nil && discount > *p.MaxDiscountAmount {
			return *p.MaxDiscountAmount
		}
		return discount
		
	case CouponTypeFixed:
		if p.DiscountValue > amount {
			return amount
		}
		return p.DiscountValue
		
	default:
		return 0
	}
}

// LoyaltyProgram represents a customer loyalty program
type LoyaltyProgram struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name              string    `json:"name" gorm:"not null" validate:"required"`
	Description       string    `json:"description"`
	PointsPerDollar   float64   `json:"points_per_dollar" gorm:"default:1"`
	DollarsPerPoint   float64   `json:"dollars_per_point" gorm:"default:0.01"`
	MinPointsToRedeem int       `json:"min_points_to_redeem" gorm:"default:100"`
	MaxPointsPerOrder *int      `json:"max_points_per_order"`
	IsActive          bool      `json:"is_active" gorm:"default:true"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for LoyaltyProgram entity
func (LoyaltyProgram) TableName() string {
	return "loyalty_programs"
}

// UserLoyaltyPoints represents user's loyalty points
type UserLoyaltyPoints struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	User            User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	TotalPoints     int       `json:"total_points" gorm:"default:0"`
	AvailablePoints int       `json:"available_points" gorm:"default:0"`
	UsedPoints      int       `json:"used_points" gorm:"default:0"`
	ExpiredPoints   int       `json:"expired_points" gorm:"default:0"`
	TierLevel       string    `json:"tier_level" gorm:"default:'bronze'"` // bronze, silver, gold, platinum
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for UserLoyaltyPoints entity
func (UserLoyaltyPoints) TableName() string {
	return "user_loyalty_points"
}

// CanRedeem checks if user can redeem points
func (ulp *UserLoyaltyPoints) CanRedeem(points int, program *LoyaltyProgram) bool {
	return ulp.AvailablePoints >= points && 
		   points >= program.MinPointsToRedeem &&
		   (program.MaxPointsPerOrder == nil || points <= *program.MaxPointsPerOrder)
}

// CalculateRedemptionValue calculates dollar value of points
func (ulp *UserLoyaltyPoints) CalculateRedemptionValue(points int, program *LoyaltyProgram) float64 {
	if !ulp.CanRedeem(points, program) {
		return 0
	}
	return float64(points) * program.DollarsPerPoint
}
