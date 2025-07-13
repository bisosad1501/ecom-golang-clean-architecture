package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/usecases"
)

// RecommendationHandler handles recommendation-related HTTP requests
type RecommendationHandler struct {
	recommendationUseCase *usecases.RecommendationUseCase
}

// NewRecommendationHandler creates a new recommendation handler
func NewRecommendationHandler(recommendationUseCase *usecases.RecommendationUseCase) *RecommendationHandler {
	return &RecommendationHandler{
		recommendationUseCase: recommendationUseCase,
	}
}

// GetRecommendations gets product recommendations
// @Summary Get product recommendations
// @Description Get product recommendations based on type and context
// @Tags recommendations
// @Accept json
// @Produce json
// @Param type query string true "Recommendation type" Enums(related,similar,frequently_bought,trending,personalized,based_on_category,based_on_brand)
// @Param product_id query string false "Product ID for product-based recommendations"
// @Param user_id query string false "User ID for personalized recommendations"
// @Param category_id query string false "Category ID for category-based recommendations"
// @Param brand_id query string false "Brand ID for brand-based recommendations"
// @Param limit query int false "Number of recommendations to return" default(10)
// @Param period query string false "Period for trending recommendations" Enums(daily,weekly,monthly) default(weekly)
// @Success 200 {object} APIResponse{data=entities.RecommendationResponse}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/v1/recommendations [get]
func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {
	// Parse query parameters
	recType := c.Query("type")
	if recType == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "recommendation type is required",
		})
		return
	}

	// Parse limit
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Build recommendation request
	req := &entities.RecommendationRequest{
		Type:    entities.RecommendationType(recType),
		Limit:   limit,
		Context: make(map[string]interface{}),
	}

	// Parse optional parameters
	if productIDStr := c.Query("product_id"); productIDStr != "" {
		if productID, err := uuid.Parse(productIDStr); err == nil {
			req.ProductID = &productID
		}
	}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			req.UserID = &userID
		}
	}

	// Get user ID from JWT if available and not provided in query
	if req.UserID == nil {
		if userID, exists := c.Get("user_id"); exists {
			if uid, ok := userID.(uuid.UUID); ok {
				req.UserID = &uid
			}
		}
	}

	// Parse context parameters
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		req.Context["category_id"] = categoryIDStr
	}

	if brandIDStr := c.Query("brand_id"); brandIDStr != "" {
		req.Context["brand_id"] = brandIDStr
	}

	if period := c.Query("period"); period != "" {
		req.Context["period"] = period
	} else {
		req.Context["period"] = "weekly"
	}

	// Get recommendations
	response, err := h.recommendationUseCase.GetRecommendations(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get recommendations",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Recommendations retrieved successfully",
		Data:    response,
	})
}

// GetRelatedProducts gets related products for a specific product
// @Summary Get related products
// @Description Get products related to a specific product
// @Tags recommendations
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param limit query int false "Number of products to return" default(10)
// @Success 200 {object} APIResponse{data=entities.RecommendationResponse}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/v1/products/{product_id}/related [get]
func (h *RecommendationHandler) GetRelatedProducts(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid product ID",
		})
		return
	}

	// Parse limit
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	req := &entities.RecommendationRequest{
		ProductID: &productID,
		Type:      entities.RecommendationTypeRelated,
		Limit:     limit,
	}

	response, err := h.recommendationUseCase.GetRecommendations(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get related products",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Related products retrieved successfully",
		Data:    response,
	})
}

// GetFrequentlyBoughtTogether gets products frequently bought together
// @Summary Get frequently bought together products
// @Description Get products that are frequently bought together with a specific product
// @Tags recommendations
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param limit query int false "Number of products to return" default(5)
// @Success 200 {object} APIResponse{data=entities.RecommendationResponse}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/v1/products/{product_id}/frequently-bought-together [get]
func (h *RecommendationHandler) GetFrequentlyBoughtTogether(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid product ID",
		})
		return
	}

	// Parse limit
	limit := 5
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	req := &entities.RecommendationRequest{
		ProductID: &productID,
		Type:      entities.RecommendationTypeFrequentlyBought,
		Limit:     limit,
	}

	response, err := h.recommendationUseCase.GetRecommendations(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get frequently bought together products",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Frequently bought together products retrieved successfully",
		Data:    response,
	})
}

// GetPersonalizedRecommendations gets personalized recommendations for the current user
// @Summary Get personalized recommendations
// @Description Get personalized product recommendations for the authenticated user
// @Tags recommendations
// @Accept json
// @Produce json
// @Param limit query int false "Number of recommendations to return" default(20)
// @Success 200 {object} APIResponse{data=entities.RecommendationResponse}
// @Failure 401 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Security BearerAuth
// @Router /api/v1/recommendations/personalized [get]
func (h *RecommendationHandler) GetPersonalizedRecommendations(c *gin.Context) {
	// Get user ID from JWT
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Authentication required",
		})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	// Parse limit
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	req := &entities.RecommendationRequest{
		UserID: &uid,
		Type:   entities.RecommendationTypePersonalized,
		Limit:  limit,
	}

	response, err := h.recommendationUseCase.GetRecommendations(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get personalized recommendations",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Personalized recommendations retrieved successfully",
		Data:    response,
	})
}

// GetTrendingProducts gets trending products
// @Summary Get trending products
// @Description Get trending products for a specific period
// @Tags recommendations
// @Accept json
// @Produce json
// @Param period query string false "Period for trending" Enums(daily,weekly,monthly) default(weekly)
// @Param limit query int false "Number of products to return" default(20)
// @Success 200 {object} APIResponse{data=entities.RecommendationResponse}
// @Failure 500 {object} APIResponse
// @Router /api/v1/recommendations/trending [get]
func (h *RecommendationHandler) GetTrendingProducts(c *gin.Context) {
	// Parse period
	period := c.DefaultQuery("period", "weekly")

	// Parse limit
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	req := &entities.RecommendationRequest{
		Type:    entities.RecommendationTypeTrending,
		Limit:   limit,
		Context: map[string]interface{}{"period": period},
	}

	response, err := h.recommendationUseCase.GetRecommendations(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get trending products",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Trending products retrieved successfully",
		Data:    response,
	})
}

// TrackInteraction tracks user interaction with products
// @Summary Track user interaction
// @Description Track user interaction with a product for recommendation purposes
// @Tags recommendations
// @Accept json
// @Produce json
// @Param interaction body TrackInteractionRequest true "Interaction data"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/v1/recommendations/track [post]
func (h *RecommendationHandler) TrackInteraction(c *gin.Context) {
	var req TrackInteractionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Get user ID from JWT if available
	var userID *uuid.UUID
	if uid, exists := c.Get("user_id"); exists {
		if parsedUID, ok := uid.(uuid.UUID); ok {
			userID = &parsedUID
		}
	}

	// Get session ID for guest users
	var sessionID *string
	if userID == nil {
		if sid := c.GetHeader("X-Session-ID"); sid != "" {
			sessionID = &sid
		}
	}

	interaction := &entities.UserProductInteraction{
		UserID:          userID,
		SessionID:       sessionID,
		ProductID:       req.ProductID,
		InteractionType: req.InteractionType,
		Value:           req.Value,
		Metadata:        req.Metadata,
	}

	err := h.recommendationUseCase.TrackInteraction(c.Request.Context(), interaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to track interaction",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Interaction tracked successfully",
	})
}

// TrackInteractionRequest represents the request for tracking interactions
type TrackInteractionRequest struct {
	ProductID       uuid.UUID                   `json:"product_id" binding:"required"`
	InteractionType entities.InteractionType    `json:"interaction_type" binding:"required"`
	Value           float64                     `json:"value,omitempty"`
	Metadata        string                      `json:"metadata,omitempty"`
}
