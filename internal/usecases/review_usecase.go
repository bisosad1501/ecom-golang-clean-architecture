package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
)

// ReviewUseCase defines review use cases
type ReviewUseCase interface {
	CreateReview(ctx context.Context, userID uuid.UUID, req CreateReviewRequest) (*ReviewResponse, error)
	GetReview(ctx context.Context, reviewID uuid.UUID) (*ReviewResponse, error)
	UpdateReview(ctx context.Context, userID, reviewID uuid.UUID, req UpdateReviewRequest) (*ReviewResponse, error)
	DeleteReview(ctx context.Context, userID, reviewID uuid.UUID) error
	GetProductReviews(ctx context.Context, productID uuid.UUID, req GetReviewsRequest) (*ReviewsResponse, error)
	GetUserReviews(ctx context.Context, userID uuid.UUID, req GetReviewsRequest) (*ReviewsResponse, error)
	VoteReview(ctx context.Context, userID, reviewID uuid.UUID, voteType entities.ReviewVoteType) error
	RemoveVote(ctx context.Context, userID, reviewID uuid.UUID) error
	GetProductRatingSummary(ctx context.Context, productID uuid.UUID) (*ProductRatingSummaryResponse, error)

	// Admin operations
	ApproveReview(ctx context.Context, reviewID uuid.UUID) error
	HideReview(ctx context.Context, reviewID uuid.UUID) error
	RejectReview(ctx context.Context, reviewID uuid.UUID) error
	GetPendingReviews(ctx context.Context, req GetReviewsRequest) (*ReviewsResponse, error)
}

type reviewUseCase struct {
	reviewRepo        repositories.ReviewRepository
	reviewVoteRepo    repositories.ReviewVoteRepository
	productRatingRepo repositories.ProductRatingRepository
	productRepo       repositories.ProductRepository
	orderRepo         repositories.OrderRepository
	userRepo          repositories.UserRepository
}

// NewReviewUseCase creates a new review use case
func NewReviewUseCase(
	reviewRepo repositories.ReviewRepository,
	reviewVoteRepo repositories.ReviewVoteRepository,
	productRatingRepo repositories.ProductRatingRepository,
	productRepo repositories.ProductRepository,
	orderRepo repositories.OrderRepository,
	userRepo repositories.UserRepository,
) ReviewUseCase {
	return &reviewUseCase{
		reviewRepo:        reviewRepo,
		reviewVoteRepo:    reviewVoteRepo,
		productRatingRepo: productRatingRepo,
		productRepo:       productRepo,
		orderRepo:         orderRepo,
		userRepo:          userRepo,
	}
}

// CreateReviewRequest represents create review request
type CreateReviewRequest struct {
	ProductID uuid.UUID  `json:"product_id" validate:"required"`
	OrderID   *uuid.UUID `json:"order_id"`
	Rating    int        `json:"rating" validate:"required,min=1,max=5"`
	Title     string     `json:"title" validate:"max=200"`    // Optional title
	Comment   string     `json:"comment" validate:"max=2000"` // Optional comment
	Images    []string   `json:"images"`
}

// UpdateReviewRequest represents update review request
type UpdateReviewRequest struct {
	Rating  *int     `json:"rating" validate:"omitempty,min=1,max=5"`
	Title   *string  `json:"title" validate:"omitempty,max=200"`
	Comment *string  `json:"comment" validate:"omitempty,max=2000"`
	Images  []string `json:"images"`
}

// GetReviewsRequest represents get reviews request
type GetReviewsRequest struct {
	Rating     *int   `json:"rating"`
	IsVerified *bool  `json:"is_verified"`
	SortBy     string `json:"sort_by"`    // created_at, rating, helpful_count
	SortOrder  string `json:"sort_order"` // asc, desc
	Limit      int    `json:"limit" validate:"min=1,max=100"`
	Offset     int    `json:"offset" validate:"min=0"`
}

// ReviewResponse represents review response
type ReviewResponse struct {
	ID                uuid.UUID                `json:"id"`
	User              ReviewUserResponse       `json:"user"`
	Product           ReviewProductResponse    `json:"product"`
	Rating            int                      `json:"rating"`
	Title             string                   `json:"title"`
	Comment           string                   `json:"comment"`
	Status            entities.ReviewStatus    `json:"status"`
	IsVerified        bool                     `json:"is_verified"`
	AdminReply        string                   `json:"admin_reply,omitempty"`
	AdminReplyAt      *time.Time               `json:"admin_reply_at,omitempty"`
	HelpfulCount      int                      `json:"helpful_count"`
	NotHelpfulCount   int                      `json:"not_helpful_count"`
	HelpfulPercentage float64                  `json:"helpful_percentage"`
	Images            []ReviewImageResponse    `json:"images"`
	UserVote          *entities.ReviewVoteType `json:"user_vote,omitempty"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
}

// ReviewUserResponse represents user info in review response
type ReviewUserResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Avatar    string    `json:"avatar,omitempty"`
}

// ReviewProductResponse represents product info in review response
type ReviewProductResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Image string    `json:"image,omitempty"`
}

// ReviewImageResponse represents review image response
type ReviewImageResponse struct {
	ID        uuid.UUID `json:"id"`
	URL       string    `json:"url"`
	AltText   string    `json:"alt_text"`
	SortOrder int       `json:"sort_order"`
}

// ReviewsResponse represents reviews list response
type ReviewsResponse struct {
	Reviews    []*ReviewResponse `json:"reviews"`
	Pagination *PaginationInfo   `json:"pagination"`
}

// ProductRatingSummaryResponse represents product rating summary
type ProductRatingSummaryResponse struct {
	ProductID          uuid.UUID       `json:"product_id"`
	AverageRating      float64         `json:"average_rating"`
	TotalReviews       int             `json:"total_reviews"`
	RatingDistribution map[int]float64 `json:"rating_distribution"`
	RatingCounts       map[int]int     `json:"rating_counts"`
}

// CreateReview creates a new review
func (uc *reviewUseCase) CreateReview(ctx context.Context, userID uuid.UUID, req CreateReviewRequest) (*ReviewResponse, error) {
	// Check if product exists
	_, err := uc.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, entities.ErrProductNotFound
	}

	// Business rule: Allow multiple comments but only one rating per user per product
	// Check if user already has a rating for this product
	existingReview, err := uc.reviewRepo.GetUserReviewForProduct(ctx, userID, req.ProductID)
	if err == nil && existingReview != nil {
		// User already has a review - update the rating and add new comment
		return uc.updateExistingReview(ctx, userID, existingReview, req)
	}

	// Verify order if provided
	var isVerified bool
	if req.OrderID != nil {
		order, err := uc.orderRepo.GetByID(ctx, *req.OrderID)
		if err == nil && order.UserID == userID {
			// Check if order contains this product and is delivered
			for _, item := range order.Items {
				if item.ProductID == req.ProductID {
					// Only verify if order is delivered (customer actually received product)
					if order.Status == entities.OrderStatusDelivered {
						isVerified = true
					}
					break
				}
			}
		}
	}

	// Generate default title if not provided
	title := req.Title
	if title == "" {
		switch req.Rating {
		case 5:
			title = "Excellent!"
		case 4:
			title = "Very Good"
		case 3:
			title = "Good"
		case 2:
			title = "Fair"
		case 1:
			title = "Poor"
		default:
			title = "Review"
		}
	}

	// Smart auto-approval logic for optimal UX
	status := uc.determineReviewStatus(req.Rating, req.Comment, req.Title, isVerified)

	// Create review
	review := &entities.Review{
		ID:         uuid.New(),
		UserID:     userID,
		ProductID:  req.ProductID,
		OrderID:    req.OrderID,
		Rating:     req.Rating,
		Title:      title,
		Comment:    req.Comment,
		Status:     status,
		IsVerified: isVerified,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := uc.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	// Real-time rating update for approved reviews
	if review.Status == entities.ReviewStatusApproved {
		if err := uc.productRatingRepo.RecalculateRating(ctx, req.ProductID); err != nil {
			// Log error but don't fail the review creation
			fmt.Printf("❌ Failed to update product rating: %v\n", err)
		} else {
			fmt.Printf("✅ Product rating updated for product %s\n", req.ProductID)
		}

		// Award loyalty points for approved reviews
		uc.awardReviewLoyaltyPoints(ctx, userID, req.Rating, len(strings.TrimSpace(req.Comment)), isVerified)
	}

	// Get the created review with relationships
	createdReview, err := uc.reviewRepo.GetByID(ctx, review.ID)
	if err != nil {
		return nil, err
	}

	return uc.toReviewResponse(createdReview, nil), nil
}

// updateExistingReview updates existing review with new rating and adds comment
func (uc *reviewUseCase) updateExistingReview(ctx context.Context, userID uuid.UUID, existingReview *entities.Review, req CreateReviewRequest) (*ReviewResponse, error) {
	// Update rating (only one rating per user per product)
	existingReview.Rating = req.Rating

	// Update title if provided
	if req.Title != "" {
		existingReview.Title = req.Title
	}

	// Handle comments: Allow multiple comments by appending new ones
	if req.Comment != "" {
		if existingReview.Comment == "" {
			// First comment
			existingReview.Comment = req.Comment
		} else {
			// Append new comment with timestamp
			timestamp := time.Now().Format("2006-01-02 15:04")
			existingReview.Comment += fmt.Sprintf("\n\n[Update %s]: %s", timestamp, req.Comment)
		}
	}

	// Re-evaluate approval status with new content
	isVerified := existingReview.IsVerified
	newStatus := uc.determineReviewStatus(existingReview.Rating, existingReview.Comment, existingReview.Title, isVerified)
	existingReview.Status = newStatus
	existingReview.UpdatedAt = time.Now()

	// Update in database
	if err := uc.reviewRepo.Update(ctx, existingReview); err != nil {
		return nil, err
	}

	// Update product rating if approved
	if existingReview.Status == entities.ReviewStatusApproved {
		if err := uc.productRatingRepo.RecalculateRating(ctx, req.ProductID); err != nil {
			fmt.Printf("❌ Failed to update product rating: %v\n", err)
		} else {
			fmt.Printf("✅ Product rating updated after review update\n")
		}

		// Award loyalty points for the update (smaller amount)
		uc.awardReviewLoyaltyPoints(ctx, userID, req.Rating, len(strings.TrimSpace(req.Comment)), isVerified)
	}

	return uc.toReviewResponse(existingReview, nil), nil
}

// determineReviewStatus determines if a review should be auto-approved based on business rules
func (uc *reviewUseCase) determineReviewStatus(rating int, comment, title string, isVerified bool) entities.ReviewStatus {
	// Check for suspicious content first
	if uc.isSuspiciousContent(comment, title) {
		return entities.ReviewStatusPending
	}

	// Auto-approve verified purchases (best UX for real customers)
	if isVerified {
		return entities.ReviewStatusApproved
	}

	// Auto-approve all positive reviews (4-5 stars) regardless of comment
	if rating >= 4 {
		return entities.ReviewStatusApproved
	}

	// Auto-approve neutral reviews (3 stars) - customers are honest
	if rating == 3 {
		return entities.ReviewStatusApproved
	}

	// For negative reviews (1-2 stars), be more flexible:
	if rating <= 2 {
		// Auto-approve if has any comment (even short ones)
		if len(strings.TrimSpace(comment)) > 0 {
			return entities.ReviewStatusApproved
		}
		// Auto-approve rating-only negative reviews too (customer frustration is valid)
		// Only flag if completely empty and suspicious patterns
		return entities.ReviewStatusApproved
	}

	// Default to approved for maximum UX flexibility
	return entities.ReviewStatusApproved
}

// isSuspiciousContent checks for suspicious patterns in review content
func (uc *reviewUseCase) isSuspiciousContent(comment, title string) bool {
	suspiciousWords := []string{
		"fake", "spam", "bot", "paid", "advertisement", "promo",
		"discount code", "coupon", "free shipping", "click here",
		"visit my", "check out my", "follow me", "subscribe",
	}

	content := strings.ToLower(comment + " " + title)

	// Check for suspicious keywords
	for _, word := range suspiciousWords {
		if strings.Contains(content, word) {
			return true
		}
	}

	// Check for excessive repetition (spam pattern)
	words := strings.Fields(content)
	if len(words) > 5 {
		wordCount := make(map[string]int)
		for _, word := range words {
			if len(word) > 3 { // Only count meaningful words
				wordCount[word]++
				if wordCount[word] > 3 { // Same word repeated more than 3 times
					return true
				}
			}
		}
	}

	// Check for excessive capitalization (spam pattern)
	if len(comment) > 10 {
		upperCount := 0
		for _, char := range comment {
			if char >= 'A' && char <= 'Z' {
				upperCount++
			}
		}
		if float64(upperCount)/float64(len(comment)) > 0.5 { // More than 50% uppercase
			return true
		}
	}

	return false
}

// isSimilarContent checks if two content strings are similar (for edit detection)
func (uc *reviewUseCase) isSimilarContent(original, updated string) bool {
	// Normalize strings
	orig := strings.ToLower(strings.TrimSpace(original))
	upd := strings.ToLower(strings.TrimSpace(updated))

	// If both empty, they're similar
	if orig == "" && upd == "" {
		return true
	}

	// If one is empty and other isn't, not similar
	if (orig == "") != (upd == "") {
		return false
	}

	// If exactly the same, they're similar
	if orig == upd {
		return true
	}

	// Calculate similarity based on length difference
	lenDiff := len(upd) - len(orig)
	if lenDiff < 0 {
		lenDiff = -lenDiff
	}

	// If length difference is more than 50% of original, not similar
	if len(orig) > 0 && float64(lenDiff)/float64(len(orig)) > 0.5 {
		return false
	}

	// For short content, be more strict
	if len(orig) < 20 {
		return orig == upd
	}

	// For longer content, allow minor changes (typo fixes, etc.)
	// Simple similarity check: count common words
	origWords := strings.Fields(orig)
	updWords := strings.Fields(upd)

	if len(origWords) == 0 && len(updWords) == 0 {
		return true
	}

	commonWords := 0
	origWordMap := make(map[string]bool)
	for _, word := range origWords {
		origWordMap[word] = true
	}

	for _, word := range updWords {
		if origWordMap[word] {
			commonWords++
		}
	}

	// If more than 70% words are common, consider similar
	totalWords := len(origWords)
	if len(updWords) > totalWords {
		totalWords = len(updWords)
	}

	if totalWords == 0 {
		return true
	}

	similarity := float64(commonWords) / float64(totalWords)
	return similarity >= 0.7
}

// awardReviewLoyaltyPoints awards loyalty points for writing reviews
func (uc *reviewUseCase) awardReviewLoyaltyPoints(ctx context.Context, userID uuid.UUID, rating, commentLength int, isVerified bool) {
	// Calculate points based on review quality
	points := 0

	// Base points for any review
	points += 5

	// Bonus for detailed reviews
	if commentLength >= 50 {
		points += 5 // +5 for detailed review
	}
	if commentLength >= 100 {
		points += 5 // +5 more for very detailed review
	}

	// Bonus for verified purchase reviews
	if isVerified {
		points += 10
	}

	// Bonus for balanced reviews (not just 5 stars)
	if rating >= 3 && rating <= 4 {
		points += 3 // Encourage honest, balanced reviews
	}

	// Award points to user
	if points > 0 {
		user, err := uc.userRepo.GetByID(ctx, userID)
		if err == nil {
			user.LoyaltyPoints += points
			if err := uc.userRepo.Update(ctx, user); err == nil {
				fmt.Printf("✅ Awarded %d loyalty points for review\n", points)
			}
		}
	}
}

// GetReview gets a review by ID
func (uc *reviewUseCase) GetReview(ctx context.Context, reviewID uuid.UUID) (*ReviewResponse, error) {
	review, err := uc.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return nil, entities.ErrReviewNotFound
	}

	return uc.toReviewResponse(review, nil), nil
}

// GetProductReviews gets reviews for a product
func (uc *reviewUseCase) GetProductReviews(ctx context.Context, productID uuid.UUID, req GetReviewsRequest) (*ReviewsResponse, error) {
	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	approvedStatus := entities.ReviewStatusApproved
	filter := entities.ReviewFilter{
		ProductID:  &productID,
		Rating:     req.Rating,
		IsVerified: req.IsVerified,
		Status:     &approvedStatus,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}

	reviews, err := uc.reviewRepo.Search(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalCount, err := uc.reviewRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]*ReviewResponse, len(reviews))
	for i, review := range reviews {
		responses[i] = uc.toReviewResponse(review, nil)
	}

	// Create pagination info using enhanced function
	context := &EcommercePaginationContext{
		EntityType: "reviews",
	}

	// Use proper offset-to-page conversion
	pagination := NewPaginationInfoFromOffset(req.Offset, req.Limit, totalCount)

	// Apply ecommerce enhancements
	if context != nil {
		// Adjust page sizes based on entity type
		pagination.PageSizes = []int{5, 10, 20} // Smaller sizes for detailed content

		// Check if cursor pagination should be used
		pagination.UseCursor = ShouldUseCursorPagination(totalCount, context.EntityType)

		// Generate cache key
		cacheParams := map[string]interface{}{
			"page":  pagination.Page,
			"limit": pagination.Limit,
		}
		if req.Rating != nil {
			cacheParams["rating"] = *req.Rating
		}
		if req.IsVerified != nil {
			cacheParams["is_verified"] = *req.IsVerified
		}
		if req.SortBy != "" {
			cacheParams["sort_by"] = req.SortBy
		}
		pagination.CacheKey = GenerateCacheKey("reviews", "", cacheParams)
	}

	return &ReviewsResponse{
		Reviews:    responses,
		Pagination: pagination,
	}, nil
}

// VoteReview votes on a review
func (uc *reviewUseCase) VoteReview(ctx context.Context, userID, reviewID uuid.UUID, voteType entities.ReviewVoteType) error {
	// Check if review exists
	_, err := uc.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return entities.ErrReviewNotFound
	}

	// Vote on the review
	if err := uc.reviewVoteRepo.VoteReview(ctx, reviewID, userID, voteType); err != nil {
		return err
	}

	// Update review vote counts
	return uc.reviewVoteRepo.UpdateReviewVoteCounts(ctx, reviewID)
}

// GetProductRatingSummary gets rating summary for a product
func (uc *reviewUseCase) GetProductRatingSummary(ctx context.Context, productID uuid.UUID) (*ProductRatingSummaryResponse, error) {
	rating, err := uc.productRatingRepo.GetByProductID(ctx, productID)
	if err != nil {
		// If no rating exists, return empty summary
		return &ProductRatingSummaryResponse{
			ProductID:          productID,
			AverageRating:      0,
			TotalReviews:       0,
			RatingDistribution: map[int]float64{1: 0, 2: 0, 3: 0, 4: 0, 5: 0},
			RatingCounts:       map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0},
		}, nil
	}

	return &ProductRatingSummaryResponse{
		ProductID:          rating.ProductID,
		AverageRating:      rating.AverageRating,
		TotalReviews:       rating.TotalReviews,
		RatingDistribution: rating.GetRatingDistribution(),
		RatingCounts:       rating.GetRatingCounts(),
	}, nil
}

// toReviewResponse converts review entity to response
func (uc *reviewUseCase) toReviewResponse(review *entities.Review, userVote *entities.ReviewVoteType) *ReviewResponse {
	response := &ReviewResponse{
		ID:                review.ID,
		Rating:            review.Rating,
		Title:             review.Title,
		Comment:           review.Comment,
		Status:            review.Status,
		IsVerified:        review.IsVerified,
		AdminReply:        review.AdminReply,
		AdminReplyAt:      review.AdminReplyAt,
		HelpfulCount:      review.HelpfulCount,
		NotHelpfulCount:   review.NotHelpfulCount,
		HelpfulPercentage: review.GetHelpfulPercentage(),
		UserVote:          userVote,
		CreatedAt:         review.CreatedAt,
		UpdatedAt:         review.UpdatedAt,
	}

	// Add user info
	if review.User.ID != uuid.Nil {
		response.User = ReviewUserResponse{
			ID:        review.User.ID,
			FirstName: review.User.FirstName,
			LastName:  review.User.LastName,
		}
		if review.User.Profile != nil {
			response.User.Avatar = review.User.Profile.Avatar
		}
	}

	// Add product info
	if review.Product.ID != uuid.Nil {
		response.Product = ReviewProductResponse{
			ID:   review.Product.ID,
			Name: review.Product.Name,
		}
		if len(review.Product.Images) > 0 {
			response.Product.Image = review.Product.Images[0].URL
		}
	}

	// Add images
	if len(review.Images) > 0 {
		images := make([]ReviewImageResponse, len(review.Images))
		for i, img := range review.Images {
			images[i] = ReviewImageResponse{
				ID:        img.ID,
				URL:       img.URL,
				AltText:   img.AltText,
				SortOrder: img.SortOrder,
			}
		}
		response.Images = images
	}

	return response
}

// UpdateReview updates an existing review
func (uc *reviewUseCase) UpdateReview(ctx context.Context, userID, reviewID uuid.UUID, req UpdateReviewRequest) (*ReviewResponse, error) {
	// Get existing review
	review, err := uc.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return nil, entities.ErrReviewNotFound
	}

	// Check if user owns the review
	if review.UserID != userID {
		return nil, entities.ErrUnauthorized
	}

	// Business rule: Allow editing anytime (like real ecommerce platforms)
	// No time restriction for editing - customers should be able to update their reviews freely
	// This matches Amazon, Shopee, Lazada behavior

	// Store original values for comparison
	originalRating := review.Rating
	originalComment := review.Comment
	originalTitle := review.Title

	// Update review fields
	if req.Rating != nil {
		review.Rating = *req.Rating
	}
	if req.Title != nil {
		review.Title = *req.Title
	}
	if req.Comment != nil {
		review.Comment = *req.Comment
	}

	// Check if verified purchase status (for better approval logic)
	isVerified := review.IsVerified

	// Smart re-approval logic for edits
	newStatus := uc.determineReviewStatus(review.Rating, review.Comment, review.Title, isVerified)

	// If only minor changes (same rating, similar content), keep approved status
	if review.Status == entities.ReviewStatusApproved &&
		originalRating == review.Rating &&
		uc.isSimilarContent(originalComment, review.Comment) &&
		uc.isSimilarContent(originalTitle, review.Title) {
		// Keep approved status for minor edits
		newStatus = entities.ReviewStatusApproved
	}

	review.Status = newStatus
	review.UpdatedAt = time.Now()

	if err := uc.reviewRepo.Update(ctx, review); err != nil {
		return nil, err
	}

	// Update product rating if approved
	if review.Status == entities.ReviewStatusApproved {
		if err := uc.productRatingRepo.RecalculateRating(ctx, review.ProductID); err != nil {
			fmt.Printf("❌ Failed to update product rating after review edit: %v\n", err)
		} else {
			fmt.Printf("✅ Product rating updated after review edit\n")
		}
	}

	return uc.toReviewResponse(review, nil), nil
}

// DeleteReview deletes a review
func (uc *reviewUseCase) DeleteReview(ctx context.Context, userID, reviewID uuid.UUID) error {
	// Get existing review
	review, err := uc.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return entities.ErrReviewNotFound
	}

	// Check if user owns the review
	if review.UserID != userID {
		return entities.ErrUnauthorized
	}

	// Delete review
	if err := uc.reviewRepo.Delete(ctx, reviewID); err != nil {
		return err
	}

	// Update product rating
	if err := uc.productRatingRepo.RecalculateRating(ctx, review.ProductID); err != nil {
		// Log error but don't fail the request
		// TODO: Add proper logging
	}

	return nil
}

// GetUserReviews gets reviews by user
func (uc *reviewUseCase) GetUserReviews(ctx context.Context, userID uuid.UUID, req GetReviewsRequest) (*ReviewsResponse, error) {
	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	filter := entities.ReviewFilter{
		UserID:    &userID,
		Rating:    req.Rating,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
		Limit:     req.Limit,
		Offset:    req.Offset,
	}

	reviews, err := uc.reviewRepo.Search(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalCount, err := uc.reviewRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]*ReviewResponse, len(reviews))
	for i, review := range reviews {
		responses[i] = uc.toReviewResponse(review, nil)
	}

	// Create pagination info using enhanced function
	context := &EcommercePaginationContext{
		EntityType: "reviews",
		UserID:     userID.String(),
	}

	// Use proper offset-to-page conversion
	pagination := NewPaginationInfoFromOffset(req.Offset, req.Limit, totalCount)

	// Apply ecommerce enhancements
	if context != nil {
		// Adjust page sizes based on entity type
		pagination.PageSizes = []int{5, 10, 20} // Smaller sizes for detailed content

		// Check if cursor pagination should be used
		pagination.UseCursor = ShouldUseCursorPagination(totalCount, context.EntityType)

		// Generate cache key
		cacheParams := map[string]interface{}{
			"page":    pagination.Page,
			"limit":   pagination.Limit,
			"user_id": context.UserID,
		}
		if req.Rating != nil {
			cacheParams["rating"] = *req.Rating
		}
		pagination.CacheKey = GenerateCacheKey("user_reviews", context.UserID, cacheParams)
	}

	return &ReviewsResponse{
		Reviews:    responses,
		Pagination: pagination,
	}, nil
}

// RemoveVote removes a vote from review
func (uc *reviewUseCase) RemoveVote(ctx context.Context, userID, reviewID uuid.UUID) error {
	// Check if review exists
	_, err := uc.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return entities.ErrReviewNotFound
	}

	// Remove vote
	if err := uc.reviewVoteRepo.RemoveVote(ctx, reviewID, userID); err != nil {
		return err
	}

	// Update review vote counts
	return uc.reviewVoteRepo.UpdateReviewVoteCounts(ctx, reviewID)
}

// ApproveReview approves a review (admin)
func (uc *reviewUseCase) ApproveReview(ctx context.Context, reviewID uuid.UUID) error {
	review, err := uc.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return entities.ErrReviewNotFound
	}

	review.Status = entities.ReviewStatusApproved
	review.UpdatedAt = time.Now()

	if err := uc.reviewRepo.Update(ctx, review); err != nil {
		return err
	}

	// Real-time rating update
	if err := uc.productRatingRepo.RecalculateRating(ctx, review.ProductID); err != nil {
		fmt.Printf("❌ Failed to update product rating after approval: %v\n", err)
	} else {
		fmt.Printf("✅ Product rating updated after review approval\n")
	}

	return nil
}

// HideReview hides a review (admin) - keeps it in database but not visible to public
func (uc *reviewUseCase) HideReview(ctx context.Context, reviewID uuid.UUID) error {
	review, err := uc.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return entities.ErrReviewNotFound
	}

	review.Status = entities.ReviewStatusHidden
	review.UpdatedAt = time.Now()

	if err := uc.reviewRepo.Update(ctx, review); err != nil {
		return err
	}

	// Real-time rating update (hidden reviews don't count)
	if err := uc.productRatingRepo.RecalculateRating(ctx, review.ProductID); err != nil {
		fmt.Printf("❌ Failed to update product rating after hiding review: %v\n", err)
	} else {
		fmt.Printf("✅ Product rating updated after hiding review\n")
	}

	return nil
}

// RejectReview rejects a review (admin) - completely removes from consideration
func (uc *reviewUseCase) RejectReview(ctx context.Context, reviewID uuid.UUID) error {
	review, err := uc.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return entities.ErrReviewNotFound
	}

	review.Status = entities.ReviewStatusRejected
	review.UpdatedAt = time.Now()

	if err := uc.reviewRepo.Update(ctx, review); err != nil {
		return err
	}

	// Update product rating
	if err := uc.productRatingRepo.RecalculateRating(ctx, review.ProductID); err != nil {
		// Log error but don't fail the request
		// TODO: Add proper logging
	}

	return nil
}

// GetPendingReviews gets pending reviews (admin)
func (uc *reviewUseCase) GetPendingReviews(ctx context.Context, req GetReviewsRequest) (*ReviewsResponse, error) {
	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	pendingStatus := entities.ReviewStatusPending
	filter := entities.ReviewFilter{
		Status:    &pendingStatus,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
		Limit:     req.Limit,
		Offset:    req.Offset,
	}

	reviews, err := uc.reviewRepo.Search(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalCount, err := uc.reviewRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]*ReviewResponse, len(reviews))
	for i, review := range reviews {
		responses[i] = uc.toReviewResponse(review, nil)
	}

	// Create pagination info using enhanced function
	context := &EcommercePaginationContext{
		EntityType: "reviews",
	}

	// Use proper offset-to-page conversion
	pagination := NewPaginationInfoFromOffset(req.Offset, req.Limit, totalCount)

	// Apply ecommerce enhancements
	if context != nil {
		// Adjust page sizes based on entity type
		pagination.PageSizes = []int{5, 10, 20} // Smaller sizes for detailed content

		// Check if cursor pagination should be used
		pagination.UseCursor = ShouldUseCursorPagination(totalCount, context.EntityType)

		// Generate cache key
		cacheParams := map[string]interface{}{
			"page":  pagination.Page,
			"limit": pagination.Limit,
		}
		if req.Rating != nil {
			cacheParams["rating"] = *req.Rating
		}
		pagination.CacheKey = GenerateCacheKey("admin_reviews", "", cacheParams)
	}

	return &ReviewsResponse{
		Reviews:    responses,
		Pagination: pagination,
	}, nil
}
