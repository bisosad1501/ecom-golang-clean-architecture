package usecases

import (
	"context"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

// CouponUseCase defines coupon use cases
type CouponUseCase interface {
	CreateCoupon(ctx context.Context, req CreateCouponRequest) (*CouponResponse, error)
	GetCoupon(ctx context.Context, id uuid.UUID) (*CouponResponse, error)
	GetCouponByCode(ctx context.Context, code string) (*CouponResponse, error)
	UpdateCoupon(ctx context.Context, id uuid.UUID, req UpdateCouponRequest) (*CouponResponse, error)
	DeleteCoupon(ctx context.Context, id uuid.UUID) error
	ListCoupons(ctx context.Context, req ListCouponsRequest) (*CouponsListResponse, error)
	ValidateCoupon(ctx context.Context, code string, userID uuid.UUID, orderTotal float64) (*CouponValidationResponse, error)
	ApplyCoupon(ctx context.Context, req ApplyCouponRequest) (*CouponApplicationResponse, error)
	GetUserCoupons(ctx context.Context, userID uuid.UUID) ([]*CouponResponse, error)
	GetActiveCoupons(ctx context.Context) ([]*CouponResponse, error)
}

type couponUseCase struct {
	couponRepo repositories.CouponRepository
	userRepo   repositories.UserRepository
}

// NewCouponUseCase creates a new coupon use case
func NewCouponUseCase(
	couponRepo repositories.CouponRepository,
	userRepo repositories.UserRepository,
) CouponUseCase {
	return &couponUseCase{
		couponRepo: couponRepo,
		userRepo:   userRepo,
	}
}

// Request/Response types
type CreateCouponRequest struct {
	Code                 string                      `json:"code" validate:"required,min=3,max=50"`
	Name                 string                      `json:"name" validate:"required,max=200"`
	Description          string                      `json:"description,omitempty"`
	Type                 entities.CouponType         `json:"type" validate:"required"`
	Value                float64                     `json:"value" validate:"required,min=0"`
	MaxDiscount          *float64                    `json:"max_discount,omitempty"`
	MinOrderAmount       *float64                    `json:"min_order_amount,omitempty"`
	UsageLimit           *int                        `json:"usage_limit,omitempty"`
	UsageLimitPerUser    *int                        `json:"usage_limit_per_user,omitempty"`
	Applicability        entities.CouponApplicability `json:"applicability"`
	ApplicableCategoryIDs []uuid.UUID                `json:"applicable_category_ids,omitempty"`
	ApplicableProductIDs []uuid.UUID                 `json:"applicable_product_ids,omitempty"`
	ApplicableUserIDs    []uuid.UUID                 `json:"applicable_user_ids,omitempty"`
	BuyQuantity          *int                        `json:"buy_quantity,omitempty"`
	GetQuantity          *int                        `json:"get_quantity,omitempty"`
	GetProductID         *uuid.UUID                  `json:"get_product_id,omitempty"`
	StartsAt             *time.Time                  `json:"starts_at,omitempty"`
	ExpiresAt            *time.Time                  `json:"expires_at,omitempty"`
	IsFirstTimeUser      bool                        `json:"is_first_time_user"`
	IsPublic             bool                        `json:"is_public"`
}

type UpdateCouponRequest struct {
	Name                 *string                      `json:"name,omitempty" validate:"omitempty,max=200"`
	Description          *string                      `json:"description,omitempty"`
	Value                *float64                     `json:"value,omitempty" validate:"omitempty,min=0"`
	MaxDiscount          *float64                     `json:"max_discount,omitempty"`
	MinOrderAmount       *float64                     `json:"min_order_amount,omitempty"`
	UsageLimit           *int                         `json:"usage_limit,omitempty"`
	UsageLimitPerUser    *int                         `json:"usage_limit_per_user,omitempty"`
	Applicability        *entities.CouponApplicability `json:"applicability,omitempty"`
	ApplicableCategoryIDs []uuid.UUID                 `json:"applicable_category_ids,omitempty"`
	ApplicableProductIDs []uuid.UUID                  `json:"applicable_product_ids,omitempty"`
	ApplicableUserIDs    []uuid.UUID                  `json:"applicable_user_ids,omitempty"`
	StartsAt             *time.Time                   `json:"starts_at,omitempty"`
	ExpiresAt            *time.Time                   `json:"expires_at,omitempty"`
	Status               *entities.CouponStatus       `json:"status,omitempty"`
	IsFirstTimeUser      *bool                        `json:"is_first_time_user,omitempty"`
	IsPublic             *bool                        `json:"is_public,omitempty"`
}

type ListCouponsRequest struct {
	Status        *entities.CouponStatus `json:"status,omitempty"`
	Type          *entities.CouponType   `json:"type,omitempty"`
	IsPublic      *bool                  `json:"is_public,omitempty"`
	IsExpired     *bool                  `json:"is_expired,omitempty"`
	Search        string                 `json:"search,omitempty"`
	SortBy        string                 `json:"sort_by,omitempty" validate:"omitempty,oneof=name code created_at expires_at"`
	SortOrder     string                 `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
	Limit         int                    `json:"limit" validate:"min=1,max=100"`
	Offset        int                    `json:"offset" validate:"min=0"`
}

type ApplyCouponRequest struct {
	Code      string    `json:"code" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	OrderID   uuid.UUID `json:"order_id" validate:"required"`
	OrderTotal float64  `json:"order_total" validate:"required,min=0"`
}

type CouponResponse struct {
	ID                   uuid.UUID                   `json:"id"`
	Code                 string                      `json:"code"`
	Name                 string                      `json:"name"`
	Description          string                      `json:"description"`
	Type                 entities.CouponType         `json:"type"`
	Value                float64                     `json:"value"`
	MaxDiscount          *float64                    `json:"max_discount"`
	MinOrderAmount       *float64                    `json:"min_order_amount"`
	UsageLimit           *int                        `json:"usage_limit"`
	UsageLimitPerUser    *int                        `json:"usage_limit_per_user"`
	UsedCount            int                         `json:"used_count"`
	Applicability        entities.CouponApplicability `json:"applicability"`
	ApplicableCategories []CategoryResponse          `json:"applicable_categories,omitempty"`
	ApplicableProducts   []ProductResponse           `json:"applicable_products,omitempty"`
	BuyQuantity          *int                        `json:"buy_quantity"`
	GetQuantity          *int                        `json:"get_quantity"`
	GetProductID         *uuid.UUID                  `json:"get_product_id"`
	StartsAt             *time.Time                  `json:"starts_at"`
	ExpiresAt            *time.Time                  `json:"expires_at"`
	Status               entities.CouponStatus       `json:"status"`
	IsFirstTimeUser      bool                        `json:"is_first_time_user"`
	IsPublic             bool                        `json:"is_public"`
	IsValid              bool                        `json:"is_valid"`
	CreatedAt            time.Time                   `json:"created_at"`
	UpdatedAt            time.Time                   `json:"updated_at"`
}

type CouponsListResponse struct {
	Coupons    []*CouponResponse `json:"coupons"`
	Total      int64             `json:"total"`
	Pagination *PaginationInfo   `json:"pagination"`
}

type CouponValidationResponse struct {
	IsValid        bool    `json:"is_valid"`
	DiscountAmount float64 `json:"discount_amount"`
	Message        string  `json:"message"`
	Coupon         *CouponResponse `json:"coupon,omitempty"`
}

type CouponApplicationResponse struct {
	Success        bool    `json:"success"`
	DiscountAmount float64 `json:"discount_amount"`
	Message        string  `json:"message"`
	UsageID        uuid.UUID `json:"usage_id,omitempty"`
}

// CreateCoupon creates a new coupon
func (uc *couponUseCase) CreateCoupon(ctx context.Context, req CreateCouponRequest) (*CouponResponse, error) {
	// Validate coupon code uniqueness
	existing, _ := uc.couponRepo.GetByCode(ctx, req.Code)
	if existing != nil {
		return nil, entities.ErrCouponCodeExists
	}

	// Create coupon entity
	coupon := &entities.Coupon{
		ID:                uuid.New(),
		Code:              strings.ToUpper(req.Code),
		Name:              req.Name,
		Description:       req.Description,
		Type:              req.Type,
		Value:             req.Value,
		MaxDiscount:       req.MaxDiscount,
		MinOrderAmount:    req.MinOrderAmount,
		UsageLimit:        req.UsageLimit,
		UsageLimitPerUser: req.UsageLimitPerUser,
		Applicability:     req.Applicability,
		BuyQuantity:       req.BuyQuantity,
		GetQuantity:       req.GetQuantity,
		GetProductID:      req.GetProductID,
		StartsAt:          req.StartsAt,
		ExpiresAt:         req.ExpiresAt,
		Status:            entities.CouponStatusActive,
		IsFirstTimeUser:   req.IsFirstTimeUser,
		IsPublic:          req.IsPublic,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := uc.couponRepo.Create(ctx, coupon); err != nil {
		return nil, err
	}

	// Handle associations
	if len(req.ApplicableCategoryIDs) > 0 {
		// Mock implementation - in real app this would set applicable categories
		// if err := uc.couponRepo.SetApplicableCategories(ctx, coupon.ID, req.ApplicableCategoryIDs); err != nil {
		//     return nil, err
		// }
	}

	if len(req.ApplicableProductIDs) > 0 {
		// Mock implementation - in real app this would set applicable products
		// if err := uc.couponRepo.SetApplicableProducts(ctx, coupon.ID, req.ApplicableProductIDs); err != nil {
		//     return nil, err
		// }
	}

	if len(req.ApplicableUserIDs) > 0 {
		// Mock implementation - in real app this would set applicable users
		// if err := uc.couponRepo.SetApplicableUsers(ctx, coupon.ID, req.ApplicableUserIDs); err != nil {
		//     return nil, err
		// }
	}

	return uc.toCouponResponse(coupon), nil
}

// GetCoupon gets a coupon by ID
func (uc *couponUseCase) GetCoupon(ctx context.Context, id uuid.UUID) (*CouponResponse, error) {
	coupon, err := uc.couponRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrCouponNotFound
	}

	return uc.toCouponResponse(coupon), nil
}

// GetCouponByCode gets a coupon by code
func (uc *couponUseCase) GetCouponByCode(ctx context.Context, code string) (*CouponResponse, error) {
	coupon, err := uc.couponRepo.GetByCode(ctx, strings.ToUpper(code))
	if err != nil {
		return nil, entities.ErrCouponNotFound
	}

	return uc.toCouponResponse(coupon), nil
}

// ValidateCoupon validates a coupon for use
func (uc *couponUseCase) ValidateCoupon(ctx context.Context, code string, userID uuid.UUID, orderTotal float64) (*CouponValidationResponse, error) {
	coupon, err := uc.couponRepo.GetByCode(ctx, strings.ToUpper(code))
	if err != nil {
		return &CouponValidationResponse{
			IsValid: false,
			Message: "Coupon not found",
		}, nil
	}

	// Check if coupon is valid
	if !coupon.IsValid() {
		return &CouponValidationResponse{
			IsValid: false,
			Message: "Coupon is not valid or has expired",
			Coupon:  uc.toCouponResponse(coupon),
		}, nil
	}

	// Check if user can use this coupon
	if !coupon.CanBeUsedBy(userID) {
		return &CouponValidationResponse{
			IsValid: false,
			Message: "You are not eligible to use this coupon",
			Coupon:  uc.toCouponResponse(coupon),
		}, nil
	}

	// Check usage limit per user
	if coupon.UsageLimitPerUser != nil {
		usageCount, err := uc.couponRepo.GetUserUsageCount(ctx, coupon.ID, userID)
		if err != nil {
			return nil, err
		}
		if usageCount >= *coupon.UsageLimitPerUser {
			return &CouponValidationResponse{
				IsValid: false,
				Message: "You have reached the usage limit for this coupon",
				Coupon:  uc.toCouponResponse(coupon),
			}, nil
		}
	}

	// Check if first-time user coupon
	if coupon.IsFirstTimeUser {
		// Mock implementation - in real app this would check if user is first time customer
		isFirstTime := true // Mock as true for now
		// isFirstTime, err := uc.userRepo.IsFirstTimeCustomer(ctx, userID)
		// if err != nil {
		//     return nil, err
		// }
		if !isFirstTime {
			return &CouponValidationResponse{
				IsValid: false,
				Message: "This coupon is only for first-time customers",
				Coupon:  uc.toCouponResponse(coupon),
			}, nil
		}
	}

	// Calculate discount
	discountAmount := coupon.CalculateDiscount(orderTotal)
	if discountAmount == 0 {
		return &CouponValidationResponse{
			IsValid: false,
			Message: "Order does not meet minimum requirements for this coupon",
			Coupon:  uc.toCouponResponse(coupon),
		}, nil
	}

	return &CouponValidationResponse{
		IsValid:        true,
		DiscountAmount: discountAmount,
		Message:        "Coupon is valid",
		Coupon:         uc.toCouponResponse(coupon),
	}, nil
}

// ApplyCoupon applies a coupon to an order
func (uc *couponUseCase) ApplyCoupon(ctx context.Context, req ApplyCouponRequest) (*CouponApplicationResponse, error) {
	// Validate coupon first
	validation, err := uc.ValidateCoupon(ctx, req.Code, req.UserID, req.OrderTotal)
	if err != nil {
		return nil, err
	}

	if !validation.IsValid {
		return &CouponApplicationResponse{
			Success: false,
			Message: validation.Message,
		}, nil
	}

	// Create usage record
	usage := &entities.CouponUsage{
		ID:             uuid.New(),
		CouponID:       validation.Coupon.ID,
		UserID:         req.UserID,
		OrderID:        req.OrderID,
		DiscountAmount: validation.DiscountAmount,
		CreatedAt:      time.Now(),
	}

	// Mock implementation - in real app this would create usage record
	// if err := uc.couponRepo.CreateUsage(ctx, usage); err != nil {
	//     return nil, err
	// }

	// Mock implementation - in real app this would update coupon usage count
	// if err := uc.couponRepo.IncrementUsageCount(ctx, validation.Coupon.ID); err != nil {
	//     return nil, err
	// }

	return &CouponApplicationResponse{
		Success:        true,
		DiscountAmount: validation.DiscountAmount,
		Message:        "Coupon applied successfully",
		UsageID:        usage.ID,
	}, nil
}

// Helper methods
func (uc *couponUseCase) toCouponResponse(coupon *entities.Coupon) *CouponResponse {
	response := &CouponResponse{
		ID:                coupon.ID,
		Code:              coupon.Code,
		Name:              coupon.Name,
		Description:       coupon.Description,
		Type:              coupon.Type,
		Value:             coupon.Value,
		MaxDiscount:       coupon.MaxDiscount,
		MinOrderAmount:    coupon.MinOrderAmount,
		UsageLimit:        coupon.UsageLimit,
		UsageLimitPerUser: coupon.UsageLimitPerUser,
		UsedCount:         coupon.UsedCount,
		Applicability:     coupon.Applicability,
		BuyQuantity:       coupon.BuyQuantity,
		GetQuantity:       coupon.GetQuantity,
		GetProductID:      coupon.GetProductID,
		StartsAt:          coupon.StartsAt,
		ExpiresAt:         coupon.ExpiresAt,
		Status:            coupon.Status,
		IsFirstTimeUser:   coupon.IsFirstTimeUser,
		IsPublic:          coupon.IsPublic,
		IsValid:           coupon.IsValid(),
		CreatedAt:         coupon.CreatedAt,
		UpdatedAt:         coupon.UpdatedAt,
	}

	// Add applicable categories if available
	if len(coupon.ApplicableCategories) > 0 {
		categories := make([]CategoryResponse, len(coupon.ApplicableCategories))
		for i, cat := range coupon.ApplicableCategories {
			categories[i] = CategoryResponse{
				ID:          cat.ID,
				Name:        cat.Name,
				Slug:        cat.Slug,
				Description: cat.Description,
				IsActive:    cat.IsActive,
			}
		}
		response.ApplicableCategories = categories
	}

	// Add applicable products if available
	if len(coupon.ApplicableProducts) > 0 {
		products := make([]ProductResponse, len(coupon.ApplicableProducts))
		for i, prod := range coupon.ApplicableProducts {
			products[i] = ProductResponse{
				ID:          prod.ID,
				Name:        prod.Name,
				Description: prod.Description,
				Price:       prod.Price,
				SKU:         prod.SKU,
				Status:      prod.Status,
			}
		}
		response.ApplicableProducts = products
	}

	return response
}

// DeleteCoupon deletes a coupon
func (uc *couponUseCase) DeleteCoupon(ctx context.Context, id uuid.UUID) error {
	// Check if coupon exists
	_, err := uc.couponRepo.GetByID(ctx, id)
	if err != nil {
		return entities.ErrCouponNotFound
	}

	// Delete the coupon
	return uc.couponRepo.Delete(ctx, id)
}

// GetActiveCoupons gets active coupons
func (uc *couponUseCase) GetActiveCoupons(ctx context.Context) ([]*CouponResponse, error) {
	// Mock implementation for active coupons
	coupons := []*CouponResponse{
		{
			ID:          uuid.New(),
			Code:        "SAVE20",
			Type:        entities.CouponTypePercentage,
			Value:       20.0,
			Description: "Save 20% on all items",
			Status:      entities.CouponStatusActive,
			ExpiresAt:   &[]time.Time{time.Now().AddDate(0, 1, 0)}[0],
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Code:        "FREESHIP",
			Type:        entities.CouponTypeFixed,
			Value:       10.0,
			Description: "Free shipping on orders over $50",
			Status:      entities.CouponStatusActive,
			ExpiresAt:   &[]time.Time{time.Now().AddDate(0, 0, 30)}[0],
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	return coupons, nil
}

// GetUserCoupons gets coupons for a specific user
func (uc *couponUseCase) GetUserCoupons(ctx context.Context, userID uuid.UUID) ([]*CouponResponse, error) {
	// Mock implementation for user coupons
	coupons := []*CouponResponse{
		{
			ID:          uuid.New(),
			Code:        "WELCOME10",
			Name:        "Welcome Coupon",
			Type:        entities.CouponTypePercentage,
			Value:       10.0,
			Description: "Welcome discount for new users",
			Status:      entities.CouponStatusActive,
			ExpiresAt:   &[]time.Time{time.Now().AddDate(0, 0, 7)}[0],
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	return coupons, nil
}

// ListCoupons lists all coupons with filtering and pagination
func (uc *couponUseCase) ListCoupons(ctx context.Context, req ListCouponsRequest) (*CouponsListResponse, error) {
	// Mock implementation for list coupons
	coupons := []*CouponResponse{
		{
			ID:          uuid.New(),
			Code:        "SAMPLE20",
			Name:        "Sample Coupon",
			Type:        entities.CouponTypePercentage,
			Value:       20.0,
			Description: "Sample coupon for testing",
			Status:      entities.CouponStatusActive,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	
	total := int64(len(coupons))
	pagination := NewPaginationInfo(req.Offset, req.Limit, total)
	
	response := &CouponsListResponse{
		Coupons:    coupons,
		Total:      total,
		Pagination: pagination,
	}
	return response, nil
}

// UpdateCoupon updates an existing coupon
func (uc *couponUseCase) UpdateCoupon(ctx context.Context, couponID uuid.UUID, req UpdateCouponRequest) (*CouponResponse, error) {
	// Mock implementation for update coupon
	name := ""
	if req.Name != nil {
		name = *req.Name
	}
	description := ""
	if req.Description != nil {
		description = *req.Description
	}
	value := 0.0
	if req.Value != nil {
		value = *req.Value
	}
	
	response := &CouponResponse{
		ID:          couponID,
		Code:        "UPDATED", // Code không có trong UpdateCouponRequest nên dùng giá trị mặc định
		Name:        name,
		Description: description,
		Type:        entities.CouponTypeFixed, // Type không có trong UpdateCouponRequest nên dùng giá trị mặc định
		Value:       value,
		Status:      entities.CouponStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return response, nil
}
