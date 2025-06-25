package handlers

import (
	"net/http"
	"strconv"

	"ecom-golang-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CouponHandler handles coupon-related HTTP requests
type CouponHandler struct {
	couponUseCase usecases.CouponUseCase
}

// NewCouponHandler creates a new coupon handler
func NewCouponHandler(couponUseCase usecases.CouponUseCase) *CouponHandler {
	return &CouponHandler{
		couponUseCase: couponUseCase,
	}
}

// CreateCoupon creates a new coupon
func (h *CouponHandler) CreateCoupon(c *gin.Context) {
	var req usecases.CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	coupon, err := h.couponUseCase.CreateCoupon(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to create coupon",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Coupon created successfully",
		Data: coupon,
	})
}

// GetCoupon retrieves a coupon by ID
func (h *CouponHandler) GetCoupon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid coupon ID",
			Details: err.Error(),
		})
		return
	}

	coupon, err := h.couponUseCase.GetCoupon(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Coupon not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Coupon retrieved successfully",
		Data: coupon,
	})
}

// GetCoupons retrieves coupons with pagination
func (h *CouponHandler) GetCoupons(c *gin.Context) {
	var req usecases.ListCouponsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	coupons, err := h.couponUseCase.ListCoupons(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get coupons",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Coupons retrieved successfully",
		Data: coupons,
	})
}

// UpdateCoupon updates a coupon
func (h *CouponHandler) UpdateCoupon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid coupon ID",
			Details: err.Error(),
		})
		return
	}

	var req usecases.UpdateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	coupon, err := h.couponUseCase.UpdateCoupon(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to update coupon",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Coupon updated successfully",
		Data: coupon,
	})
}

// DeleteCoupon deletes a coupon
func (h *CouponHandler) DeleteCoupon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid coupon ID",
			Details: err.Error(),
		})
		return
	}

	if err := h.couponUseCase.DeleteCoupon(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to delete coupon",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Coupon deleted successfully",
	})
}

// ValidateCoupon validates a coupon
func (h *CouponHandler) ValidateCoupon(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Coupon code is required",
		})
		return
	}

	// Extract parameters from request
	var requestBody struct {
		UserID     uuid.UUID `json:"user_id" validate:"required"`
		OrderTotal float64   `json:"order_total" validate:"required,min=0"`
	}
	
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	validation, err := h.couponUseCase.ValidateCoupon(c.Request.Context(), code, requestBody.UserID, requestBody.OrderTotal)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Coupon validation failed",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Coupon validated successfully",
		Data: validation,
	})
}

// ApplyCoupon applies a coupon to an order
func (h *CouponHandler) ApplyCoupon(c *gin.Context) {
	var req usecases.ApplyCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	result, err := h.couponUseCase.ApplyCoupon(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Failed to apply coupon",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Coupon applied successfully",
		Data: result,
	})
}

// ListCoupons returns paginated list of coupons
func (h *CouponHandler) ListCoupons(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// Calculate offset from page
	offset := (page - 1) * limit

	req := usecases.ListCouponsRequest{
		Search:    search,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Limit:     limit,
		Offset:    offset,
	}

	coupons, err := h.couponUseCase.ListCoupons(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to list coupons",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Coupons retrieved successfully",
		Data: coupons,
	})
}
