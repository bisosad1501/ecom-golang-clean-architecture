package repositories

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// ReviewRepository defines the interface for review data access
type ReviewRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, review *entities.Review) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Review, error)
	Update(ctx context.Context, review *entities.Review) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetByProductID(ctx context.Context, productID uuid.UUID, filter entities.ReviewFilter) ([]*entities.Review, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, filter entities.ReviewFilter) ([]*entities.Review, error)
	Search(ctx context.Context, filter entities.ReviewFilter) ([]*entities.Review, error)
	Count(ctx context.Context, filter entities.ReviewFilter) (int64, error)

	// Product-specific operations
	GetProductReviews(ctx context.Context, productID uuid.UUID, limit, offset int) ([]*entities.Review, error)
	GetProductReviewsWithRating(ctx context.Context, productID uuid.UUID, rating int, limit, offset int) ([]*entities.Review, error)
	CountProductReviews(ctx context.Context, productID uuid.UUID) (int64, error)
	CountProductReviewsByRating(ctx context.Context, productID uuid.UUID, rating int) (int64, error)

	// User-specific operations
	GetUserReviews(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Review, error)
	HasUserReviewedProduct(ctx context.Context, userID, productID uuid.UUID) (bool, error)
	GetUserReviewForProduct(ctx context.Context, userID, productID uuid.UUID) (*entities.Review, error)

	// Status operations
	GetPendingReviews(ctx context.Context, limit, offset int) ([]*entities.Review, error)
	ApproveReview(ctx context.Context, reviewID uuid.UUID) error
	RejectReview(ctx context.Context, reviewID uuid.UUID) error
	BulkUpdateStatus(ctx context.Context, reviewIDs []uuid.UUID, status entities.ReviewStatus) error

	// Verification operations
	MarkAsVerified(ctx context.Context, reviewID uuid.UUID) error
	GetVerifiedReviews(ctx context.Context, productID uuid.UUID, limit, offset int) ([]*entities.Review, error)

	// Statistics
	GetAverageRating(ctx context.Context, productID uuid.UUID) (float64, error)
	GetRatingDistribution(ctx context.Context, productID uuid.UUID) (map[int]int, error)
	GetReviewStats(ctx context.Context, productID uuid.UUID) (*entities.ReviewSummary, error)
	CountReviewsByStatus(ctx context.Context, status entities.ReviewStatus) (int64, error)
}

// ReviewVoteRepository defines the interface for review vote data access
type ReviewVoteRepository interface {
	// Basic operations
	Create(ctx context.Context, vote *entities.ReviewVote) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ReviewVote, error)
	Update(ctx context.Context, vote *entities.ReviewVote) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Vote operations
	VoteReview(ctx context.Context, reviewID, userID uuid.UUID, voteType entities.ReviewVoteType) error
	RemoveVote(ctx context.Context, reviewID, userID uuid.UUID) error
	GetUserVote(ctx context.Context, reviewID, userID uuid.UUID) (*entities.ReviewVote, error)
	HasUserVoted(ctx context.Context, reviewID, userID uuid.UUID) (bool, error)

	// Statistics
	CountHelpfulVotes(ctx context.Context, reviewID uuid.UUID) (int, error)
	CountNotHelpfulVotes(ctx context.Context, reviewID uuid.UUID) (int, error)
	GetVoteCounts(ctx context.Context, reviewID uuid.UUID) (helpful int, notHelpful int, err error)
	UpdateReviewVoteCounts(ctx context.Context, reviewID uuid.UUID) error

	// Bulk operations
	GetVotesByReviewIDs(ctx context.Context, reviewIDs []uuid.UUID) (map[uuid.UUID][]*entities.ReviewVote, error)
	GetUserVotesForReviews(ctx context.Context, userID uuid.UUID, reviewIDs []uuid.UUID) (map[uuid.UUID]*entities.ReviewVote, error)
}

// ProductRatingRepository defines the interface for product rating data access
type ProductRatingRepository interface {
	// Basic operations
	Create(ctx context.Context, rating *entities.ProductRating) error
	GetByProductID(ctx context.Context, productID uuid.UUID) (*entities.ProductRating, error)
	Update(ctx context.Context, rating *entities.ProductRating) error
	Delete(ctx context.Context, productID uuid.UUID) error

	// Rating calculations
	RecalculateRating(ctx context.Context, productID uuid.UUID) error
	UpdateRatingFromReview(ctx context.Context, productID uuid.UUID, oldRating, newRating int) error
	RemoveRatingFromReview(ctx context.Context, productID uuid.UUID, rating int) error

	// Bulk operations
	RecalculateAllRatings(ctx context.Context) error
	GetProductsWithHighestRating(ctx context.Context, limit int) ([]*entities.ProductRating, error)
	GetProductsWithMostReviews(ctx context.Context, limit int) ([]*entities.ProductRating, error)

	// Statistics
	GetAverageRatingAcrossProducts(ctx context.Context) (float64, error)
	GetTotalReviewsCount(ctx context.Context) (int64, error)
}

// ReviewImageRepository defines the interface for review image data access
type ReviewImageRepository interface {
	// Basic operations
	Create(ctx context.Context, image *entities.ReviewImage) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ReviewImage, error)
	Update(ctx context.Context, image *entities.ReviewImage) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Review-specific operations
	GetByReviewID(ctx context.Context, reviewID uuid.UUID) ([]*entities.ReviewImage, error)
	DeleteByReviewID(ctx context.Context, reviewID uuid.UUID) error
	CountByReviewID(ctx context.Context, reviewID uuid.UUID) (int, error)

	// Bulk operations
	CreateBatch(ctx context.Context, images []*entities.ReviewImage) error
	UpdateSortOrder(ctx context.Context, reviewID uuid.UUID, imageOrders map[uuid.UUID]int) error
}
