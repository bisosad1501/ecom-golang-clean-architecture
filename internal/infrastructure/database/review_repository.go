package database

import (
	"context"
	"fmt"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type reviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository creates a new review repository
func NewReviewRepository(db *gorm.DB) repositories.ReviewRepository {
	return &reviewRepository{db: db}
}

// Create creates a new review
func (r *reviewRepository) Create(ctx context.Context, review *entities.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

// GetByID retrieves a review by ID
func (r *reviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Review, error) {
	var review entities.Review
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		Preload("Images").
		Preload("Votes").
		First(&review, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrReviewNotFound
		}
		return nil, err
	}

	return &review, nil
}

// Update updates an existing review
func (r *reviewRepository) Update(ctx context.Context, review *entities.Review) error {
	result := r.db.WithContext(ctx).Save(review)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrReviewNotFound
	}
	return nil
}

// Delete deletes a review by ID
func (r *reviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.Review{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrReviewNotFound
	}
	return nil
}

// GetByProductID retrieves reviews by product ID with filters
func (r *reviewRepository) GetByProductID(ctx context.Context, productID uuid.UUID, filter entities.ReviewFilter) ([]*entities.Review, error) {
	var reviews []*entities.Review

	query := r.db.WithContext(ctx).
		Preload("User").
		Preload("Images").
		Where("product_id = ?", productID)

	query = r.applyFilters(query, filter)

	err := query.Find(&reviews).Error
	return reviews, err
}

// GetByUserID retrieves reviews by user ID with filters
func (r *reviewRepository) GetByUserID(ctx context.Context, userID uuid.UUID, filter entities.ReviewFilter) ([]*entities.Review, error) {
	var reviews []*entities.Review

	query := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Images").
		Where("user_id = ?", userID)

	query = r.applyFilters(query, filter)

	err := query.Find(&reviews).Error
	return reviews, err
}

// Search searches reviews with filters
func (r *reviewRepository) Search(ctx context.Context, filter entities.ReviewFilter) ([]*entities.Review, error) {
	var reviews []*entities.Review

	query := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		Preload("Images")

	query = r.applyFilters(query, filter)

	err := query.Find(&reviews).Error
	return reviews, err
}

// Count counts reviews with filters
func (r *reviewRepository) Count(ctx context.Context, filter entities.ReviewFilter) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&entities.Review{})
	query = r.applyFilters(query, filter)

	err := query.Count(&count).Error
	return count, err
}

// GetProductRating gets aggregated rating for a product
func (r *reviewRepository) GetProductRating(ctx context.Context, productID uuid.UUID) (*entities.ProductRating, error) {
	var rating entities.ProductRating
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		First(&rating).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return zero rating if no reviews exist
			return &entities.ProductRating{
				ProductID:     productID,
				AverageRating: 0,
				TotalReviews:  0,
			}, nil
		}
		return nil, err
	}

	return &rating, nil
}

// GetRatingBreakdown gets rating breakdown for a product
func (r *reviewRepository) GetRatingBreakdown(ctx context.Context, productID uuid.UUID) (map[int]int64, error) {
	var results []struct {
		Rating int   `json:"rating"`
		Count  int64 `json:"count"`
	}

	err := r.db.WithContext(ctx).
		Model(&entities.Review{}).
		Select("rating, COUNT(*) as count").
		Where("product_id = ? AND status = ?", productID, entities.ReviewStatusApproved).
		Group("rating").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	breakdown := make(map[int]int64)
	for i := 1; i <= 5; i++ {
		breakdown[i] = 0
	}

	for _, result := range results {
		breakdown[result.Rating] = result.Count
	}

	return breakdown, nil
}

// HasUserReviewedProduct checks if user has reviewed a product
func (r *reviewRepository) HasUserReviewedProduct(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Review{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error

	return count > 0, err
}

// GetProductReviews gets reviews for a product
func (r *reviewRepository) GetProductReviews(ctx context.Context, productID uuid.UUID, limit, offset int) ([]*entities.Review, error) {
	var reviews []*entities.Review
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND status = ?", productID, entities.ReviewStatusApproved).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Preload("User").
		Preload("Product").
		Find(&reviews).Error
	return reviews, err
}

// GetProductReviewsWithRating gets reviews for a product with specific rating
func (r *reviewRepository) GetProductReviewsWithRating(ctx context.Context, productID uuid.UUID, rating int, limit, offset int) ([]*entities.Review, error) {
	var reviews []*entities.Review
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND rating = ? AND status = ?", productID, rating, entities.ReviewStatusApproved).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

// CreateOrUpdateVote creates or updates a review vote
func (r *reviewRepository) CreateOrUpdateVote(ctx context.Context, vote *entities.ReviewVote) error {
	// Check if vote already exists
	var existingVote entities.ReviewVote
	err := r.db.WithContext(ctx).
		Where("review_id = ? AND user_id = ?", vote.ReviewID, vote.UserID).
		First(&existingVote).Error

	if err == gorm.ErrRecordNotFound {
		// Create new vote
		return r.db.WithContext(ctx).Create(vote).Error
	} else if err != nil {
		return err
	}

	// Update existing vote
	existingVote.VoteType = vote.VoteType
	existingVote.UpdatedAt = vote.UpdatedAt
	return r.db.WithContext(ctx).Save(&existingVote).Error
}

// UpdateProductRating updates the aggregated rating for a product
func (r *reviewRepository) UpdateProductRating(ctx context.Context, productID uuid.UUID) error {
	// Calculate new rating statistics
	var stats struct {
		AverageRating float64 `json:"average_rating"`
		TotalReviews  int     `json:"total_reviews"`
		Rating1Count  int     `json:"rating_1_count"`
		Rating2Count  int     `json:"rating_2_count"`
		Rating3Count  int     `json:"rating_3_count"`
		Rating4Count  int     `json:"rating_4_count"`
		Rating5Count  int     `json:"rating_5_count"`
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			COALESCE(AVG(rating::numeric), 0) as average_rating,
			COUNT(*) as total_reviews,
			COUNT(CASE WHEN rating = 1 THEN 1 END) as rating_1_count,
			COUNT(CASE WHEN rating = 2 THEN 1 END) as rating_2_count,
			COUNT(CASE WHEN rating = 3 THEN 1 END) as rating_3_count,
			COUNT(CASE WHEN rating = 4 THEN 1 END) as rating_4_count,
			COUNT(CASE WHEN rating = 5 THEN 1 END) as rating_5_count
		FROM reviews 
		WHERE product_id = ? AND status = ?
	`, productID, entities.ReviewStatusApproved).Scan(&stats).Error

	if err != nil {
		return err
	}

	// Upsert product rating
	productRating := &entities.ProductRating{
		ProductID:     productID,
		AverageRating: stats.AverageRating,
		TotalReviews:  stats.TotalReviews,
		Rating1Count:  stats.Rating1Count,
		Rating2Count:  stats.Rating2Count,
		Rating3Count:  stats.Rating3Count,
		Rating4Count:  stats.Rating4Count,
		Rating5Count:  stats.Rating5Count,
	}

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "product_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"average_rating", "total_reviews", "rating_1_count",
				"rating_2_count", "rating_3_count", "rating_4_count",
				"rating_5_count", "updated_at",
			}),
		}).
		Create(productRating).Error
}

// CountReviewsByStatus counts reviews by status
func (r *reviewRepository) CountReviewsByStatus(ctx context.Context, status entities.ReviewStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Review{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

// ApproveReview approves a review
func (r *reviewRepository) ApproveReview(ctx context.Context, reviewID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.Review{}).
		Where("id = ?", reviewID).
		Update("status", entities.ReviewStatusApproved).Error
}

// RejectReview rejects a review
func (r *reviewRepository) RejectReview(ctx context.Context, reviewID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.Review{}).
		Where("id = ?", reviewID).
		Update("status", entities.ReviewStatusRejected).Error
}

// BulkUpdateStatus updates status for multiple reviews
func (r *reviewRepository) BulkUpdateStatus(ctx context.Context, reviewIDs []uuid.UUID, status entities.ReviewStatus) error {
	return r.db.WithContext(ctx).
		Model(&entities.Review{}).
		Where("id IN ?", reviewIDs).
		Update("status", status).Error
}

// MarkAsVerified marks a review as verified
func (r *reviewRepository) MarkAsVerified(ctx context.Context, reviewID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.Review{}).
		Where("id = ?", reviewID).
		Update("is_verified", true).Error
}

// GetVerifiedReviews gets verified reviews for a product
func (r *reviewRepository) GetVerifiedReviews(ctx context.Context, productID uuid.UUID, limit, offset int) ([]*entities.Review, error) {
	var reviews []*entities.Review
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND is_verified = ? AND status = ?", productID, true, entities.ReviewStatusApproved).
		Limit(limit).
		Offset(offset).
		Find(&reviews).Error
	return reviews, err
}

// GetPendingReviews gets pending reviews
func (r *reviewRepository) GetPendingReviews(ctx context.Context, limit, offset int) ([]*entities.Review, error) {
	var reviews []*entities.Review
	err := r.db.WithContext(ctx).
		Where("status = ?", entities.ReviewStatusPending).
		Limit(limit).
		Offset(offset).
		Find(&reviews).Error
	return reviews, err
}

// CountProductReviews counts reviews for a product
func (r *reviewRepository) CountProductReviews(ctx context.Context, productID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Review{}).
		Where("product_id = ? AND status = ?", productID, entities.ReviewStatusApproved).
		Count(&count).Error
	return count, err
}

// CountProductReviewsByRating counts reviews for a product by rating
func (r *reviewRepository) CountProductReviewsByRating(ctx context.Context, productID uuid.UUID, rating int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Review{}).
		Where("product_id = ? AND rating = ? AND status = ?", productID, rating, entities.ReviewStatusApproved).
		Count(&count).Error
	return count, err
}

// GetAverageRating gets average rating for a product
func (r *reviewRepository) GetAverageRating(ctx context.Context, productID uuid.UUID) (float64, error) {
	var avg float64
	err := r.db.WithContext(ctx).
		Model(&entities.Review{}).
		Where("product_id = ? AND status = ?", productID, entities.ReviewStatusApproved).
		Select("AVG(rating)").
		Scan(&avg).Error
	return avg, err
}

// GetRatingDistribution gets rating distribution for a product
func (r *reviewRepository) GetRatingDistribution(ctx context.Context, productID uuid.UUID) (map[int]int, error) {
	type RatingCount struct {
		Rating int
		Count  int
	}

	var results []RatingCount
	err := r.db.WithContext(ctx).
		Model(&entities.Review{}).
		Where("product_id = ? AND status = ?", productID, entities.ReviewStatusApproved).
		Select("rating, COUNT(*) as count").
		Group("rating").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	distribution := make(map[int]int)
	for _, result := range results {
		distribution[result.Rating] = result.Count
	}

	return distribution, nil
}

// GetReviewStats gets review statistics for a product
func (r *reviewRepository) GetReviewStats(ctx context.Context, productID uuid.UUID) (*entities.ReviewSummary, error) {
	var stats struct {
		AverageRating float64
		TotalReviews  int
		Rating1Count  int
		Rating2Count  int
		Rating3Count  int
		Rating4Count  int
		Rating5Count  int
	}

	// Get aggregated statistics
	err := r.db.WithContext(ctx).
		Model(&entities.Review{}).
		Select(`
			COALESCE(AVG(rating), 0) as average_rating,
			COUNT(*) as total_reviews,
			SUM(CASE WHEN rating = 1 THEN 1 ELSE 0 END) as rating_1_count,
			SUM(CASE WHEN rating = 2 THEN 1 ELSE 0 END) as rating_2_count,
			SUM(CASE WHEN rating = 3 THEN 1 ELSE 0 END) as rating_3_count,
			SUM(CASE WHEN rating = 4 THEN 1 ELSE 0 END) as rating_4_count,
			SUM(CASE WHEN rating = 5 THEN 1 ELSE 0 END) as rating_5_count
		`).
		Where("product_id = ? AND status = ?", productID, entities.ReviewStatusApproved).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	// Calculate rating distribution
	ratingDistribution := make(map[int]float64)
	ratingCounts := make(map[int]int)

	if stats.TotalReviews > 0 {
		ratingCounts[1] = stats.Rating1Count
		ratingCounts[2] = stats.Rating2Count
		ratingCounts[3] = stats.Rating3Count
		ratingCounts[4] = stats.Rating4Count
		ratingCounts[5] = stats.Rating5Count

		ratingDistribution[1] = float64(stats.Rating1Count) / float64(stats.TotalReviews) * 100
		ratingDistribution[2] = float64(stats.Rating2Count) / float64(stats.TotalReviews) * 100
		ratingDistribution[3] = float64(stats.Rating3Count) / float64(stats.TotalReviews) * 100
		ratingDistribution[4] = float64(stats.Rating4Count) / float64(stats.TotalReviews) * 100
		ratingDistribution[5] = float64(stats.Rating5Count) / float64(stats.TotalReviews) * 100
	}

	// Get recent reviews
	var recentReviews []entities.Review
	err = r.db.WithContext(ctx).
		Where("product_id = ? AND status = ?", productID, entities.ReviewStatusApproved).
		Order("created_at DESC").
		Limit(5).
		Preload("User").
		Find(&recentReviews).Error

	if err != nil {
		return nil, err
	}

	return &entities.ReviewSummary{
		ProductID:          productID,
		AverageRating:      stats.AverageRating,
		TotalReviews:       stats.TotalReviews,
		RatingDistribution: ratingDistribution,
		RatingCounts:       ratingCounts,
		RecentReviews:      recentReviews,
	}, nil
}

// GetByProductIDsWithUser retrieves reviews for multiple products with user data (bulk operation)
func (r *reviewRepository) GetByProductIDsWithUser(ctx context.Context, productIDs []uuid.UUID, limit int) ([]*entities.Review, error) {
	if len(productIDs) == 0 {
		return []*entities.Review{}, nil
	}

	var reviews []*entities.Review
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		Where("product_id IN ? AND status = ?", productIDs, entities.ReviewStatusApproved).
		Order("created_at DESC").
		Limit(limit).
		Find(&reviews).Error
	return reviews, err
}

// GetRecentReviewsWithUser retrieves recent reviews with user data (optimized)
func (r *reviewRepository) GetRecentReviewsWithUser(ctx context.Context, limit int) ([]*entities.Review, error) {
	var reviews []*entities.Review
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		Preload("Product.Category").
		Preload("Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC").Limit(1)
		}).
		Where("status = ?", entities.ReviewStatusApproved).
		Order("created_at DESC").
		Limit(limit).
		Find(&reviews).Error
	return reviews, err
}

// GetUserReviewForProduct gets user's review for a product
func (r *reviewRepository) GetUserReviewForProduct(ctx context.Context, userID, productID uuid.UUID) (*entities.Review, error) {
	var review entities.Review
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID, productID).
		First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// GetUserReviews gets reviews by a user
func (r *reviewRepository) GetUserReviews(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Review, error) {
	var reviews []*entities.Review
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

// Helper method to apply filters
func (r *reviewRepository) applyFilters(query *gorm.DB, filter entities.ReviewFilter) *gorm.DB {
	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if filter.Rating != nil {
		query = query.Where("rating = ?", *filter.Rating)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.IsVerified != nil {
		query = query.Where("is_verified = ?", *filter.IsVerified)
	}

	if filter.MinRating != nil {
		query = query.Where("rating >= ?", *filter.MinRating)
	}

	if filter.MaxRating != nil {
		query = query.Where("rating <= ?", *filter.MaxRating)
	}

	// Apply sorting
	if filter.SortBy != "" {
		order := "DESC"
		if filter.SortOrder == "asc" {
			order = "ASC"
		}

		switch filter.SortBy {
		case "rating":
			query = query.Order(fmt.Sprintf("rating %s", order))
		case "helpful_count":
			query = query.Order(fmt.Sprintf("helpful_count %s", order))
		case "created_at":
			query = query.Order(fmt.Sprintf("created_at %s", order))
		default:
			query = query.Order("created_at DESC")
		}
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset >= 0 {
		query = query.Offset(filter.Offset)
	}

	return query
}
