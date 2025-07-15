package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/usecases"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReviewHandler struct {
	reviewUseCase usecases.ReviewUseCase
	fileUseCase   usecases.FileUseCase
}

// NewReviewHandler creates a new review handler
func NewReviewHandler(reviewUseCase usecases.ReviewUseCase, fileUseCase usecases.FileUseCase) *ReviewHandler {
	return &ReviewHandler{
		reviewUseCase: reviewUseCase,
		fileUseCase:   fileUseCase,
	}
}

// CreateReview creates a new review (supports both JSON and multipart form with images)
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req usecases.CreateReviewRequest

	// Check content type to determine how to parse request
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		// Handle multipart form (with images)
		if err := h.parseMultipartReviewRequest(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data", "details": err.Error()})
			return
		}
	} else {
		// Handle JSON request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}
	}

	review, err := h.reviewUseCase.CreateReview(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Review created successfully", "data": review})
}

// GetReview gets a review by ID
func (h *ReviewHandler) GetReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	review, err := h.reviewUseCase.GetReview(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": review})
}

// UpdateReview updates a review
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var req usecases.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user owns this review
	existingReview, err := h.reviewUseCase.GetReview(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	if existingReview.User.ID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own reviews"})
		return
	}

	review, err := h.reviewUseCase.UpdateReview(c.Request.Context(), userID, id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review updated successfully", "data": review})
}

// DeleteReview deletes a review
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user owns this review or is admin
	existingReview, err := h.reviewUseCase.GetReview(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	userRole, _ := c.Get("role")
	if existingReview.User.ID != userID && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own reviews"})
		return
	}

	if err := h.reviewUseCase.DeleteReview(c.Request.Context(), userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

// GetProductReviews gets reviews for a product
func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Parse and validate pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	// Validate and normalize pagination for reviews
	page, limit, err = usecases.ValidateAndNormalizePaginationForEntity(page, limit, "reviews")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to offset for repository
	offset := (page - 1) * limit

	// Parse query parameters
	req := usecases.GetReviewsRequest{
		Limit:  limit,
		Offset: offset,
	}

	if ratingStr := c.Query("rating"); ratingStr != "" {
		if rating, err := strconv.Atoi(ratingStr); err == nil && rating >= 1 && rating <= 5 {
			req.Rating = &rating
		}
	}

	// Filter by verified purchase
	if verifiedStr := c.Query("verified"); verifiedStr != "" {
		if verified, err := strconv.ParseBool(verifiedStr); err == nil {
			req.IsVerified = &verified
		}
	}

	req.SortBy = c.Query("sort_by")
	req.SortOrder = c.Query("sort_order")

	response, err := h.reviewUseCase.GetProductReviews(c.Request.Context(), productID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product reviews"})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       response.Reviews,
		Pagination: response.Pagination,
	})
}

// GetUserReviews gets reviews by a user
func (h *ReviewHandler) GetUserReviews(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse and validate pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	// Validate and normalize pagination for reviews
	page, limit, err = usecases.ValidateAndNormalizePaginationForEntity(page, limit, "reviews")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to offset for repository
	offset := (page - 1) * limit

	// Parse query parameters
	req := usecases.GetReviewsRequest{
		Limit:  limit,
		Offset: offset,
	}

	if ratingStr := c.Query("rating"); ratingStr != "" {
		if rating, err := strconv.Atoi(ratingStr); err == nil && rating >= 1 && rating <= 5 {
			req.Rating = &rating
		}
	}

	req.SortBy = c.Query("sort_by")
	req.SortOrder = c.Query("sort_order")

	response, err := h.reviewUseCase.GetUserReviews(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user reviews"})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       response.Reviews,
		Pagination: response.Pagination,
	})
}

// VoteReview votes on a review
func (h *ReviewHandler) VoteReview(c *gin.Context) {
	idStr := c.Param("id")
	reviewID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var req struct {
		IsHelpful bool `json:"is_helpful"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	voteType := entities.ReviewVoteHelpful
	if !req.IsHelpful {
		voteType = entities.ReviewVoteNotHelpful
	}

	if err := h.reviewUseCase.VoteReview(c.Request.Context(), userID, reviewID, voteType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to vote on review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote recorded successfully"})
}

// GetProductRating gets aggregated rating for a product
func (h *ReviewHandler) GetProductRating(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	rating, err := h.reviewUseCase.GetProductRatingSummary(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rating})
}

// parseMultipartReviewRequest parses multipart form data for review creation with images
func (h *ReviewHandler) parseMultipartReviewRequest(c *gin.Context, req *usecases.CreateReviewRequest) error {
	// Parse form data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		return err
	}

	// Parse product_id
	if productIDStr := c.PostForm("product_id"); productIDStr != "" {
		productID, err := uuid.Parse(productIDStr)
		if err != nil {
			return err
		}
		req.ProductID = productID
	}

	// Parse order_id (optional)
	if orderIDStr := c.PostForm("order_id"); orderIDStr != "" {
		orderID, err := uuid.Parse(orderIDStr)
		if err != nil {
			return err
		}
		req.OrderID = &orderID
	}

	// Parse rating
	if ratingStr := c.PostForm("rating"); ratingStr != "" {
		rating, err := strconv.Atoi(ratingStr)
		if err != nil || rating < 1 || rating > 5 {
			return err
		}
		req.Rating = rating
	}

	// Parse title and comment
	req.Title = c.PostForm("title")
	req.Comment = c.PostForm("comment")

	// Handle image files
	form := c.Request.MultipartForm
	if files := form.File["images"]; len(files) > 0 {
		// Upload images and get URLs
		var imageURLs []string
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				continue // Skip invalid files
			}
			defer file.Close()

			// Use file service to upload image
			if h.fileUseCase != nil {
				userIDStr := c.GetString("user_id")
				uploadResp, err := h.fileUseCase.UploadImage(c.Request.Context(), file, fileHeader, entities.FileUploadTypeUser, &userIDStr)
				if err == nil && uploadResp != nil {
					imageURLs = append(imageURLs, uploadResp.URL)
				}
			}
		}
		req.Images = imageURLs
	}

	return nil
}
