package handlers

import (
	"net/http"

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

// GetCart handles getting user's cart
// @Summary Get user's cart
// @Description Get current user's shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} usecases.CartResponse
// @Failure 401 {object} ErrorResponse
// @Router /cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
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
}

// AddToCart handles adding item to cart
// @Summary Add item to cart
// @Description Add a product to the shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.AddToCartRequest true "Add to cart request"
// @Success 200 {object} usecases.CartResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /cart/items [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
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

	var req usecases.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
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
}

// UpdateCartItem handles updating cart item quantity
// @Summary Update cart item
// @Description Update quantity of an item in the cart
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.UpdateCartItemRequest true "Update cart item request"
// @Success 200 {object} usecases.CartResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /cart/items [put]
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
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

	var req usecases.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

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

	err = h.cartUseCase.ClearCart(c.Request.Context(), userID)
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
