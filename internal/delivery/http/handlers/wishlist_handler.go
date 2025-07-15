package handlers

import (
	"net/http"
	"strconv"

	"ecom-golang-clean-architecture/internal/usecases"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WishlistHandler handles wishlist-related HTTP requests
type WishlistHandler struct {
	wishlistUseCase usecases.WishlistUseCase
}

// NewWishlistHandler creates a new wishlist handler
func NewWishlistHandler(wishlistUseCase usecases.WishlistUseCase) *WishlistHandler {
	return &WishlistHandler{
		wishlistUseCase: wishlistUseCase,
	}
}

// GetWishlist handles getting user's wishlist
// @Summary Get user's wishlist
// @Description Get current user's wishlist with pagination
// @Tags wishlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} usecases.WishlistResponse
// @Failure 401 {object} ErrorResponse
// @Router /wishlist [get]
func (h *WishlistHandler) GetWishlist(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	// Parse and validate pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0")) // 0 means use default

	// Validate and normalize pagination for wishlist
	page, limit, err = usecases.ValidateAndNormalizePaginationForEntity(page, limit, "wishlist")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Convert to offset for repository
	offset := (page - 1) * limit

	req := usecases.GetWishlistRequest{
		Limit:  limit,
		Offset: offset,
	}

	response, err := h.wishlistUseCase.GetWishlist(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       response.Items,
		Pagination: response.Pagination,
	})
}

// AddToWishlist handles adding a product to wishlist
// @Summary Add product to wishlist
// @Description Add a product to the user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Product ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /wishlist [post]
func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req struct {
		ProductID string `json:"product_id" validate:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid product ID",
		})
		return
	}

	err = h.wishlistUseCase.AddToWishlist(c.Request.Context(), userID, productID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product added to wishlist successfully",
	})
}

// RemoveFromWishlist handles removing a product from wishlist
// @Summary Remove product from wishlist
// @Description Remove a product from the user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productId path string true "Product ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /wishlist/{productId} [delete]
func (h *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid product ID",
		})
		return
	}

	err = h.wishlistUseCase.RemoveFromWishlist(c.Request.Context(), userID, productID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product removed from wishlist successfully",
	})
}

// CheckWishlistStatus handles checking if a product is in wishlist
// @Summary Check if product is in wishlist
// @Description Check if a specific product is in the user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productId path string true "Product ID"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /wishlist/{productId}/status [get]
func (h *WishlistHandler) CheckWishlistStatus(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid product ID",
		})
		return
	}

	isInWishlist, err := h.wishlistUseCase.IsInWishlist(c.Request.Context(), userID, productID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: map[string]bool{
			"in_wishlist": isInWishlist,
		},
	})
}

// ClearWishlist handles clearing all items from wishlist
// @Summary Clear wishlist
// @Description Remove all items from the user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /wishlist [delete]
func (h *WishlistHandler) ClearWishlist(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	err = h.wishlistUseCase.ClearWishlist(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Wishlist cleared successfully",
	})
}

// GetWishlistCount handles getting wishlist item count
// @Summary Get wishlist count
// @Description Get the total number of items in the user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]int64
// @Failure 401 {object} ErrorResponse
// @Router /wishlist/count [get]
func (h *WishlistHandler) GetWishlistCount(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	count, err := h.wishlistUseCase.GetWishlistCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: map[string]int64{
			"count": count,
		},
	})
}
