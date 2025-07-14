package handlers

import (
	"net/http"
	"regexp"

	"ecom-golang-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CartHandler handles cart-related HTTP requests
type CartHandler struct {
	cartUseCase usecases.CartUseCase
}

// NewCartHandler creates a new cart handler
func NewCartHandler(cartUseCase usecases.CartUseCase) *CartHandler {
	return &CartHandler{
		cartUseCase: cartUseCase,
	}
}

// GetCart handles getting user's cart or guest cart
// @Summary Get user's cart or guest cart
// @Description Get current user's shopping cart or guest cart by session ID
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Session-ID header string false "Session ID for guest cart"
// @Success 200 {object} usecases.CartResponse
// @Failure 401 {object} ErrorResponse
// @Router /cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	// Check if user is authenticated
	userIDInterface, exists := c.Get("user_id")
	if exists {
		// Authenticated user - userID is already a UUID from middleware
		userID, ok := userIDInterface.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid user ID format",
			})
			return
		}

		cart, err := h.cartUseCase.GetCart(c.Request.Context(), userID)
		if err != nil {
			c.JSON(getErrorStatusCode(err), ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse{
			Data: cart,
		})
		return
	}

	// Guest user - check for session ID
	sessionID := c.GetHeader("X-Session-ID")
	if !validateSessionID(sessionID) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Valid session ID is required for guest cart",
		})
		return
	}

	cart, err := h.cartUseCase.GetGuestCart(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: cart,
	})
}

// AddToCart handles adding item to cart for authenticated users or guests
// @Summary Add item to cart
// @Description Add a product to the shopping cart (user or guest)
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Session-ID header string false "Session ID for guest cart"
// @Param request body usecases.AddToCartRequest true "Add to cart request"
// @Success 200 {object} usecases.CartResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /cart/items [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	var req usecases.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Check if user is authenticated
	userIDInterface, exists := c.Get("user_id")
	if exists {
		// Authenticated user - userID is already a UUID from middleware
		userID, ok := userIDInterface.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid user ID format",
			})
			return
		}

		cart, err := h.cartUseCase.AddToCart(c.Request.Context(), userID, req)
		if err != nil {
			c.JSON(getErrorStatusCode(err), ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse{
			Message: "Item added to cart successfully",
			Data:    cart,
		})
		return
	}

	// Guest user - check for session ID
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Session ID is required for guest cart",
		})
		return
	}

	cart, err := h.cartUseCase.AddToGuestCart(c.Request.Context(), sessionID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Item added to cart successfully",
		Data:    cart,
	})
}

// UpdateCartItem handles updating cart item quantity
// @Summary Update cart item
// @Description Update quantity of an item in the cart
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productId path string true "Product ID"
// @Param request body usecases.UpdateCartItemRequest true "Update cart item request"
// @Success 200 {object} usecases.CartResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /cart/items/{productId} [put]
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID format",
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

	var req usecases.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Set the product ID from URL parameter
	req.ProductID = productID

	cart, err := h.cartUseCase.UpdateCartItem(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Cart item updated successfully",
		Data:    cart,
	})
}

// RemoveFromCart handles removing item from cart
// @Summary Remove item from cart
// @Description Remove a product from the shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productId path string true "Product ID"
// @Success 200 {object} usecases.CartResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /cart/items/{productId} [delete]
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID format",
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

	cart, err := h.cartUseCase.RemoveFromCart(c.Request.Context(), userID, productID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Item removed from cart successfully",
		Data:    cart,
	})
}

// ClearCart handles clearing all items from cart
// @Summary Clear cart
// @Description Remove all items from the shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /cart [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID format",
		})
		return
	}

	err := h.cartUseCase.ClearCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Cart cleared successfully",
	})
}

// MergeGuestCart handles merging guest cart with user cart when user logs in
// @Summary Merge guest cart with user cart
// @Description Merge guest cart items into user cart when user logs in
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body MergeCartRequest true "Merge cart request"
// @Success 200 {object} usecases.CartResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /cart/merge [post]
func (h *CartHandler) MergeGuestCart(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID format",
		})
		return
	}

	var req MergeCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Determine merge strategy
	strategy := usecases.MergeStrategyAuto // default
	if req.Strategy != "" {
		switch req.Strategy {
		case "auto":
			strategy = usecases.MergeStrategyAuto
		case "replace":
			strategy = usecases.MergeStrategyReplace
		case "keep_user":
			strategy = usecases.MergeStrategyKeepUser
		case "merge":
			strategy = usecases.MergeStrategyMerge
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid merge strategy. Valid options: auto, replace, keep_user, merge",
			})
			return
		}
	}

	cart, err := h.cartUseCase.MergeGuestCartWithStrategy(c.Request.Context(), userID, req.SessionID, strategy)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Guest cart merged successfully",
		Data:    cart,
	})
}

// CheckCartConflict checks if merging guest cart will cause conflicts
// @Summary Check cart merge conflicts
// @Description Check if guest cart merge will cause conflicts with existing user cart
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CheckConflictRequest true "Check conflict request"
// @Success 200 {object} CartConflictResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /cart/check-conflict [post]
func (h *CartHandler) CheckCartConflict(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID format",
		})
		return
	}

	var req CheckConflictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	conflict, err := h.cartUseCase.CheckMergeConflict(c.Request.Context(), userID, req.SessionID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: conflict,
	})
}

// MergeCartRequest represents the request to merge guest cart
type MergeCartRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Strategy  string `json:"strategy,omitempty"` // auto, replace, keep_user, merge
}

// CheckConflictRequest represents the request to check merge conflicts
type CheckConflictRequest struct {
	SessionID string `json:"session_id" binding:"required"`
}

// CartConflictResponse represents cart merge conflict information
type CartConflictResponse struct {
	HasConflict      bool                    `json:"has_conflict"`
	UserCartExists   bool                    `json:"user_cart_exists"`
	GuestCartExists  bool                    `json:"guest_cart_exists"`
	ConflictingItems []ConflictingItem       `json:"conflicting_items,omitempty"`
	UserCart         *usecases.CartResponse  `json:"user_cart,omitempty"`
	GuestCart        *usecases.CartResponse  `json:"guest_cart,omitempty"`
	Recommendations  []string                `json:"recommendations,omitempty"`
}

// ConflictingItem represents an item that exists in both carts
type ConflictingItem struct {
	ProductID       string  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	UserQuantity    int     `json:"user_quantity"`
	GuestQuantity   int     `json:"guest_quantity"`
	UserPrice       float64 `json:"user_price"`
	GuestPrice      float64 `json:"guest_price"`
	PriceDifference float64 `json:"price_difference"`
}

// validateSessionID validates the format of session ID
func validateSessionID(sessionID string) bool {
	if sessionID == "" {
		return false
	}

	// Check length (should be reasonable)
	if len(sessionID) < 10 || len(sessionID) > 100 {
		return false
	}

	// Check for basic alphanumeric format (allow hyphens for UUIDs)
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\-_]+$`, sessionID)
	return matched
}
