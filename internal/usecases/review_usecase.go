package usecases

import (
	"context"
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
	RejectReview(ctx context.Context, reviewID uuid.UUID) error
	GetPendingReviews(ctx context.Context, req GetReviewsRequest) (*ReviewsResponse, error)
}

type reviewUseCase struct {
	reviewRepo        repositories.ReviewRepository
	reviewVoteRepo    repositories.ReviewVoteRepository
	productRatingRepo repositories.ProductRatingRepository
	productRepo       repositories.ProductRepository
	orderRepo         repositories.OrderRepository
}

// NewReviewUseCase creates a new review use case
func NewReviewUseCase(
	reviewRepo repositories.ReviewRepository,
	reviewVoteRepo repositories.ReviewVoteRepository,
	productRatingRepo repositories.ProductRatingRepository,
	productRepo repositories.ProductRepository,
	orderRepo repositories.OrderRepository,
) ReviewUseCase {
	return &reviewUseCase{
		reviewRepo:        reviewRepo,
		reviewVoteRepo:    reviewVoteRepo,
		productRatingRepo: productRatingRepo,
		productRepo:       productRepo,
		orderRepo:         orderRepo,
	}
}

// CreateReviewRequest represents create review request
type CreateReviewRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	OrderID   *uuid.UUID `json:"order_id"`
	Rating    int       `json:"rating" validate:"required,min=1,max=5"`
	Title     string    `json:"title" validate:"required,max=200"`
	Comment   string    `json:"comment" validate:"max=2000"`
	Images    []string  `json:"images"`
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
	Rating    *int    `json:"rating"`
	SortBy    string  `json:"sort_by"`    // created_at, rating, helpful_count
	SortOrder string  `json:"sort_order"` // asc, desc
	Limit     int     `json:"limit" validate:"min=1,max=100"`
	Offset    int     `json:"offset" validate:"min=0"`
}

// ReviewResponse represents review response
type ReviewResponse struct {
	ID              uuid.UUID            `json:"id"`
	User            ReviewUserResponse   `json:"user"`
	Product         ReviewProductResponse `json:"product"`
	Rating          int                  `json:"rating"`
	Title           string               `json:"title"`
	Comment         string               `json:"comment"`
	Status          entities.ReviewStatus `json:"status"`
	IsVerified      bool                 `json:"is_verified"`
	HelpfulCount    int                  `json:"helpful_count"`
	NotHelpfulCount int                  `json:"not_helpful_count"`
	HelpfulPercentage float64            `json:"helpful_percentage"`
	Images          []ReviewImageResponse `json:"images"`
	UserVote        *entities.ReviewVoteType `json:"user_vote,omitempty"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
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
	ID       uuid.UUID `json:"id"`
	URL      string    `json:"url"`
	AltText  string    `json:"alt_text"`
	SortOrder int      `json:"sort_order"`
}

// ReviewsResponse represents reviews list response
type ReviewsResponse struct {
	Reviews    []*ReviewResponse `json:"reviews"`
	TotalCount int64             `json:"total_count"`
	Limit      int               `json:"limit"`
	Offset     int               `json:"offset"`
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

	// Check if user has already reviewed this product
	hasReviewed, err := uc.reviewRepo.HasUserReviewedProduct(ctx, userID, req.ProductID)
	if err != nil {
		return nil, err
	}
	if hasReviewed {
		return nil, entities.ErrConflict // User has already reviewed this product
	}

	// Verify order if provided
	var isVerified bool
	if req.OrderID != nil {
		order, err := uc.orderRepo.GetByID(ctx, *req.OrderID)
		if err == nil && order.UserID == userID {
			// Check if order contains this product
			for _, item := range order.Items {
				if item.ProductID == req.ProductID {
					isVerified = true
					break
				}
			}
		}
	}

	// Create review
	review := &entities.Review{
		ID:         uuid.New(),
		UserID:     userID,
		ProductID:  req.ProductID,
		OrderID:    req.OrderID,
		Rating:     req.Rating,
		Title:      req.Title,
		Comment:    req.Comment,
		Status:     entities.ReviewStatusPending,
		IsVerified: isVerified,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := uc.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	// Update product rating
	if err := uc.productRatingRepo.RecalculateRating(ctx, req.ProductID); err != nil {
		// Log error but don't fail the review creation
	}

	// Get the created review with relationships
	createdReview, err := uc.reviewRepo.GetByID(ctx, review.ID)
	if err != nil {
		return nil, err
	}

	return uc.toReviewResponse(createdReview, nil), nil
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
		ProductID: &productID,
		Rating:    req.Rating,
		Status:    &approvedStatus,
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

	return &ReviewsResponse{
		Reviews:    responses,
		TotalCount: totalCount,
		Limit:      req.Limit,
		Offset:     req.Offset,
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
	// Implementation would go here
	return nil, entities.ErrNotImplemented
}

// DeleteReview deletes a review
func (uc *reviewUseCase) DeleteReview(ctx context.Context, userID, reviewID uuid.UUID) error {
	// Implementation would go here
	return entities.ErrNotImplemented
}

// GetUserReviews gets reviews by user
func (uc *reviewUseCase) GetUserReviews(ctx context.Context, userID uuid.UUID, req GetReviewsRequest) (*ReviewsResponse, error) {
	// Implementation would go here
	return nil, entities.ErrNotImplemented
}

// RemoveVote removes a vote from review
func (uc *reviewUseCase) RemoveVote(ctx context.Context, userID, reviewID uuid.UUID) error {
	// Implementation would go here
	return entities.ErrNotImplemented
}

// ApproveReview approves a review (admin)
func (uc *reviewUseCase) ApproveReview(ctx context.Context, reviewID uuid.UUID) error {
	// Implementation would go here
	return entities.ErrNotImplemented
}

// RejectReview rejects a review (admin)
func (uc *reviewUseCase) RejectReview(ctx context.Context, reviewID uuid.UUID) error {
	// Implementation would go here
	return entities.ErrNotImplemented
}

// GetPendingReviews gets pending reviews (admin)
func (uc *reviewUseCase) GetPendingReviews(ctx context.Context, req GetReviewsRequest) (*ReviewsResponse, error) {
	// Implementation would go here
	return nil, entities.ErrNotImplemented
}
