package database

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRatingRepository struct {
	db *gorm.DB
}

// NewProductRatingRepository creates a new product rating repository
func NewProductRatingRepository(db *gorm.DB) repositories.ProductRatingRepository {
	return &productRatingRepository{db: db}
}

// Create creates a new product rating
func (r *productRatingRepository) Create(ctx context.Context, rating *entities.ProductRating) error {
	return r.db.WithContext(ctx).Create(rating).Error
}

// GetByID gets a product rating by ID
func (r *productRatingRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ProductRating, error) {
	var rating entities.ProductRating
	err := r.db.WithContext(ctx).First(&rating, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

// GetByProduct gets product rating by product ID
func (r *productRatingRepository) GetByProduct(ctx context.Context, productID uuid.UUID) (*entities.ProductRating, error) {
	var rating entities.ProductRating
	err := r.db.WithContext(ctx).First(&rating, "product_id = ?", productID).Error
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

// Update updates a product rating
func (r *productRatingRepository) Update(ctx context.Context, rating *entities.ProductRating) error {
	rating.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(rating).Error
}

// UpdateRating updates product rating based on reviews
func (r *productRatingRepository) UpdateRating(ctx context.Context, productID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Calculate rating statistics from reviews
		var stats struct {
			TotalReviews   int64   `json:"total_reviews"`
			AverageRating  float64 `json:"average_rating"`
			Rating1Count   int64   `json:"rating_1_count"`
			Rating2Count   int64   `json:"rating_2_count"`
			Rating3Count   int64   `json:"rating_3_count"`
			Rating4Count   int64   `json:"rating_4_count"`
			Rating5Count   int64   `json:"rating_5_count"`
		}

		// Get total reviews and average rating
		err := tx.Model(&entities.Review{}).
			Select("COUNT(*) as total_reviews, COALESCE(AVG(rating), 0) as average_rating").
			Where("product_id = ? AND status = ?", productID, entities.ReviewStatusApproved).
			Scan(&stats).Error
		if err != nil {
			return err
		}

		// Get rating distribution
		for rating := 1; rating <= 5; rating++ {
			var count int64
			err := tx.Model(&entities.Review{}).
				Where("product_id = ? AND rating = ? AND status = ?", productID, rating, entities.ReviewStatusApproved).
				Count(&count).Error
			if err != nil {
				return err
			}

			switch rating {
			case 1:
				stats.Rating1Count = count
			case 2:
				stats.Rating2Count = count
			case 3:
				stats.Rating3Count = count
			case 4:
				stats.Rating4Count = count
			case 5:
				stats.Rating5Count = count
			}
		}

		// Update or create product rating
		productRating := &entities.ProductRating{
			ProductID:      productID,
			TotalReviews:   int(stats.TotalReviews),
			AverageRating:  stats.AverageRating,
			Rating1Count:   int(stats.Rating1Count),
			Rating2Count:   int(stats.Rating2Count),
			Rating3Count:   int(stats.Rating3Count),
			Rating4Count:   int(stats.Rating4Count),
			Rating5Count:   int(stats.Rating5Count),
			UpdatedAt:      time.Now(),
		}

		// Try to update existing record
		result := tx.Model(&entities.ProductRating{}).
			Where("product_id = ?", productID).
			Updates(productRating)

		if result.Error != nil {
			return result.Error
		}

		// If no record was updated, create a new one
		if result.RowsAffected == 0 {
			productRating.ID = uuid.New()
			productRating.CreatedAt = time.Now()
			return tx.Create(productRating).Error
		}

		return nil
	})
}

// Delete deletes a product rating
func (r *productRatingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.ProductRating{}, "id = ?", id).Error
}

// DeleteByProduct deletes product rating by product ID
func (r *productRatingRepository) DeleteByProduct(ctx context.Context, productID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.ProductRating{}, "product_id = ?", productID).Error
}

// List lists product ratings with filters
func (r *productRatingRepository) List(ctx context.Context, filters repositories.ProductRatingFilters) ([]*entities.ProductRating, error) {
	var ratings []*entities.ProductRating
	query := r.db.WithContext(ctx)

	if filters.ProductID != nil {
		query = query.Where("product_id = ?", *filters.ProductID)
	}

	if filters.MinRating != nil {
		query = query.Where("average_rating >= ?", *filters.MinRating)
	}

	if filters.MaxRating != nil {
		query = query.Where("average_rating <= ?", *filters.MaxRating)
	}

	if filters.MinReviews != nil {
		query = query.Where("total_reviews >= ?", *filters.MinReviews)
	}

	// Apply sorting
	switch filters.SortBy {
	case "average_rating":
		if filters.SortOrder == "desc" {
			query = query.Order("average_rating DESC")
		} else {
			query = query.Order("average_rating ASC")
		}
	case "total_reviews":
		if filters.SortOrder == "desc" {
			query = query.Order("total_reviews DESC")
		} else {
			query = query.Order("total_reviews ASC")
		}
	case "updated_at":
		if filters.SortOrder == "desc" {
			query = query.Order("updated_at DESC")
		} else {
			query = query.Order("updated_at ASC")
		}
	default:
		query = query.Order("average_rating DESC")
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&ratings).Error
	return ratings, err
}

// Count counts product ratings with filters
func (r *productRatingRepository) Count(ctx context.Context, filters repositories.ProductRatingFilters) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.ProductRating{})

	if filters.ProductID != nil {
		query = query.Where("product_id = ?", *filters.ProductID)
	}

	if filters.MinRating != nil {
		query = query.Where("average_rating >= ?", *filters.MinRating)
	}

	if filters.MaxRating != nil {
		query = query.Where("average_rating <= ?", *filters.MaxRating)
	}

	if filters.MinReviews != nil {
		query = query.Where("total_reviews >= ?", *filters.MinReviews)
	}

	err := query.Count(&count).Error
	return count, err
}

// GetTopRatedProducts gets products with highest ratings
func (r *productRatingRepository) GetTopRatedProducts(ctx context.Context, limit int, minReviews int64) ([]*entities.ProductRating, error) {
	var ratings []*entities.ProductRating
	err := r.db.WithContext(ctx).
		Where("total_reviews >= ?", minReviews).
		Order("average_rating DESC, total_reviews DESC").
		Limit(limit).
		Find(&ratings).Error
	return ratings, err
}

// GetMostReviewedProducts gets products with most reviews
func (r *productRatingRepository) GetMostReviewedProducts(ctx context.Context, limit int) ([]*entities.ProductRating, error) {
	var ratings []*entities.ProductRating
	err := r.db.WithContext(ctx).
		Order("total_reviews DESC, average_rating DESC").
		Limit(limit).
		Find(&ratings).Error
	return ratings, err
}

// GetRatingDistribution gets rating distribution for all products
func (r *productRatingRepository) GetRatingDistribution(ctx context.Context) (*entities.RatingDistribution, error) {
	var distribution entities.RatingDistribution
	
	err := r.db.WithContext(ctx).
		Model(&entities.ProductRating{}).
		Select(`
			SUM(rating_1_count) as rating_1_count,
			SUM(rating_2_count) as rating_2_count,
			SUM(rating_3_count) as rating_3_count,
			SUM(rating_4_count) as rating_4_count,
			SUM(rating_5_count) as rating_5_count,
			SUM(total_reviews) as total_reviews,
			AVG(average_rating) as average_rating
		`).
		Scan(&distribution).Error
	
	return &distribution, err
}

// GetProductsWithoutRating gets products that don't have rating records
func (r *productRatingRepository) GetProductsWithoutRating(ctx context.Context, limit, offset int) ([]uuid.UUID, error) {
	var productIDs []uuid.UUID
	err := r.db.WithContext(ctx).
		Table("products").
		Select("products.id").
		Joins("LEFT JOIN product_ratings ON products.id = product_ratings.product_id").
		Where("product_ratings.id IS NULL").
		Limit(limit).
		Offset(offset).
		Pluck("products.id", &productIDs).Error
	return productIDs, err
}

// BulkUpdateRatings updates ratings for multiple products
func (r *productRatingRepository) BulkUpdateRatings(ctx context.Context, productIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, productID := range productIDs {
			err := r.UpdateRating(ctx, productID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Exists checks if a product rating exists
func (r *productRatingRepository) Exists(ctx context.Context, productID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.ProductRating{}).
		Where("product_id = ?", productID).
		Count(&count).Error
	return count > 0, err
}

// GetAverageRatingAcrossProducts gets average rating across all products
func (r *productRatingRepository) GetAverageRatingAcrossProducts(ctx context.Context) (float64, error) {
	var avgRating float64
	err := r.db.WithContext(ctx).
		Model(&entities.ProductRating{}).
		Select("AVG(average_rating)").
		Where("total_reviews > 0").
		Scan(&avgRating).Error
	return avgRating, err
}

// GetByProductID gets product rating by product ID (alias for GetByProduct)
func (r *productRatingRepository) GetByProductID(ctx context.Context, productID uuid.UUID) (*entities.ProductRating, error) {
	return r.GetByProduct(ctx, productID)
}

// GetProductsWithHighestRating gets products with highest ratings
func (r *productRatingRepository) GetProductsWithHighestRating(ctx context.Context, limit int) ([]*entities.ProductRating, error) {
	var ratings []*entities.ProductRating
	err := r.db.WithContext(ctx).
		Where("total_reviews >= ?", 5). // Only products with at least 5 reviews
		Order("average_rating DESC, total_reviews DESC").
		Limit(limit).
		Find(&ratings).Error
	return ratings, err
}

// GetProductsWithMostReviews gets products with most reviews
func (r *productRatingRepository) GetProductsWithMostReviews(ctx context.Context, limit int) ([]*entities.ProductRating, error) {
	var ratings []*entities.ProductRating
	err := r.db.WithContext(ctx).
		Order("total_reviews DESC, average_rating DESC").
		Limit(limit).
		Find(&ratings).Error
	return ratings, err
}

// GetTotalReviewsCount gets total count of all reviews across all products
func (r *productRatingRepository) GetTotalReviewsCount(ctx context.Context) (int64, error) {
	var totalReviews int64
	err := r.db.WithContext(ctx).
		Model(&entities.ProductRating{}).
		Select("SUM(total_reviews)").
		Scan(&totalReviews).Error
	return totalReviews, err
}

// RecalculateAllRatings recalculates ratings for all products
func (r *productRatingRepository) RecalculateAllRatings(ctx context.Context) error {
	// Get all products that have reviews
	var productIDs []uuid.UUID
	err := r.db.WithContext(ctx).
		Model(&entities.Review{}).
		Select("DISTINCT product_id").
		Scan(&productIDs).Error
	if err != nil {
		return err
	}

	// Recalculate rating for each product
	for _, productID := range productIDs {
		err = r.UpdateRating(ctx, productID)
		if err != nil {
			return err
		}
	}

	return nil
}

// RecalculateRating recalculates rating for a specific product
func (r *productRatingRepository) RecalculateRating(ctx context.Context, productID uuid.UUID) error {
	return r.UpdateRating(ctx, productID)
}

// RemoveRatingFromReview removes rating contribution from a review
func (r *productRatingRepository) RemoveRatingFromReview(ctx context.Context, reviewID uuid.UUID, oldRating int) error {
	// Get the review to find the product ID
	var review entities.Review
	err := r.db.WithContext(ctx).First(&review, "id = ?", reviewID).Error
	if err != nil {
		return err
	}

	// Recalculate the product rating after removing this review's contribution
	return r.UpdateRating(ctx, review.ProductID)
}

// UpdateRatingFromReview updates rating when a review is added/updated
func (r *productRatingRepository) UpdateRatingFromReview(ctx context.Context, reviewID uuid.UUID, oldRating, newRating int) error {
	// Get the review to find the product ID
	var review entities.Review
	err := r.db.WithContext(ctx).First(&review, "id = ?", reviewID).Error
	if err != nil {
		return err
	}

	// Recalculate the product rating
	return r.UpdateRating(ctx, review.ProductID)
}
