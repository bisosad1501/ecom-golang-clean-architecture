package entities

import (
	"time"

	"github.com/google/uuid"
)

// ReviewStatus represents the status of a review
type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "pending"
	ReviewStatusApproved ReviewStatus = "approved"
	ReviewStatusRejected ReviewStatus = "rejected"
)

// Review represents a product review
type Review struct {
	ID        uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID    `json:"user_id" gorm:"type:uuid;not null;index"`
	User      User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ProductID uuid.UUID    `json:"product_id" gorm:"type:uuid;not null;index"`
	Product   Product      `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	OrderID   *uuid.UUID   `json:"order_id" gorm:"type:uuid;index"` // Optional: link to order for verified purchases
	Order     *Order       `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Rating    int          `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5" validate:"required,min=1,max=5"`
	Title     string       `json:"title" gorm:"not null" validate:"required,max=200"`
	Comment   string       `json:"comment" gorm:"type:text" validate:"max=2000"`
	Status    ReviewStatus `json:"status" gorm:"default:'pending'"`
	IsVerified bool        `json:"is_verified" gorm:"default:false"` // Verified purchase
	HelpfulCount int       `json:"helpful_count" gorm:"default:0"`
	NotHelpfulCount int    `json:"not_helpful_count" gorm:"default:0"`
	Images      []ReviewImage `json:"images,omitempty" gorm:"foreignKey:ReviewID"`
	Votes       []ReviewVote  `json:"votes,omitempty" gorm:"foreignKey:ReviewID"`
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Review entity
func (Review) TableName() string {
	return "reviews"
}

// IsApproved checks if the review is approved
func (r *Review) IsApproved() bool {
	return r.Status == ReviewStatusApproved
}

// GetHelpfulPercentage calculates the helpful percentage
func (r *Review) GetHelpfulPercentage() float64 {
	totalVotes := r.HelpfulCount + r.NotHelpfulCount
	if totalVotes == 0 {
		return 0
	}
	return float64(r.HelpfulCount) / float64(totalVotes) * 100
}

// ReviewImage represents images attached to reviews
type ReviewImage struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ReviewID uuid.UUID `json:"review_id" gorm:"type:uuid;not null;index"`
	Review   Review    `json:"review,omitempty" gorm:"foreignKey:ReviewID"`
	URL      string    `json:"url" gorm:"not null" validate:"required,url"`
	AltText  string    `json:"alt_text"`
	SortOrder int      `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for ReviewImage entity
func (ReviewImage) TableName() string {
	return "review_images"
}

// ReviewVoteType represents the type of vote on a review
type ReviewVoteType string

const (
	ReviewVoteHelpful    ReviewVoteType = "helpful"
	ReviewVoteNotHelpful ReviewVoteType = "not_helpful"
)

// ReviewVote represents votes on reviews (helpful/not helpful)
type ReviewVote struct {
	ID       uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ReviewID uuid.UUID      `json:"review_id" gorm:"type:uuid;not null;index"`
	Review   Review         `json:"review,omitempty" gorm:"foreignKey:ReviewID"`
	UserID   uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	User     User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	VoteType ReviewVoteType `json:"vote_type" gorm:"not null"`
	CreatedAt time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for ReviewVote entity
func (ReviewVote) TableName() string {
	return "review_votes"
}

// IsHelpful checks if the vote is helpful
func (rv *ReviewVote) IsHelpful() bool {
	return rv.VoteType == ReviewVoteHelpful
}

// ProductRating represents aggregated rating data for a product
type ProductRating struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID     uuid.UUID `json:"product_id" gorm:"type:uuid;not null;uniqueIndex"`
	Product       Product   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	AverageRating float64   `json:"average_rating" gorm:"default:0"`
	TotalReviews  int       `json:"total_reviews" gorm:"default:0"`
	Rating1Count  int       `json:"rating_1_count" gorm:"default:0"`
	Rating2Count  int       `json:"rating_2_count" gorm:"default:0"`
	Rating3Count  int       `json:"rating_3_count" gorm:"default:0"`
	Rating4Count  int       `json:"rating_4_count" gorm:"default:0"`
	Rating5Count  int       `json:"rating_5_count" gorm:"default:0"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for ProductRating entity
func (ProductRating) TableName() string {
	return "product_ratings"
}

// GetRatingDistribution returns the rating distribution as percentages
func (pr *ProductRating) GetRatingDistribution() map[int]float64 {
	if pr.TotalReviews == 0 {
		return map[int]float64{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}
	}

	total := float64(pr.TotalReviews)
	return map[int]float64{
		1: float64(pr.Rating1Count) / total * 100,
		2: float64(pr.Rating2Count) / total * 100,
		3: float64(pr.Rating3Count) / total * 100,
		4: float64(pr.Rating4Count) / total * 100,
		5: float64(pr.Rating5Count) / total * 100,
	}
}

// GetRatingCounts returns the rating counts
func (pr *ProductRating) GetRatingCounts() map[int]int {
	return map[int]int{
		1: pr.Rating1Count,
		2: pr.Rating2Count,
		3: pr.Rating3Count,
		4: pr.Rating4Count,
		5: pr.Rating5Count,
	}
}

// ReviewSummary represents a summary of reviews for a product
type ReviewSummary struct {
	ProductID         uuid.UUID          `json:"product_id"`
	AverageRating     float64            `json:"average_rating"`
	TotalReviews      int                `json:"total_reviews"`
	RatingDistribution map[int]float64   `json:"rating_distribution"`
	RatingCounts      map[int]int        `json:"rating_counts"`
	RecentReviews     []Review           `json:"recent_reviews,omitempty"`
}

// ReviewFilter represents filters for review queries
type ReviewFilter struct {
	ProductID    *uuid.UUID    `json:"product_id"`
	UserID       *uuid.UUID    `json:"user_id"`
	Rating       *int          `json:"rating"`
	Status       *ReviewStatus `json:"status"`
	IsVerified   *bool         `json:"is_verified"`
	MinRating    *int          `json:"min_rating"`
	MaxRating    *int          `json:"max_rating"`
	SortBy       string        `json:"sort_by"`       // created_at, rating, helpful_count
	SortOrder    string        `json:"sort_order"`    // asc, desc
	Limit        int           `json:"limit"`
	Offset       int           `json:"offset"`
}
