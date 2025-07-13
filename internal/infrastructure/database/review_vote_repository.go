package database

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type reviewVoteRepository struct {
	db *gorm.DB
}

// NewReviewVoteRepository creates a new review vote repository
func NewReviewVoteRepository(db *gorm.DB) repositories.ReviewVoteRepository {
	return &reviewVoteRepository{db: db}
}

// Create creates a new review vote
func (r *reviewVoteRepository) Create(ctx context.Context, vote *entities.ReviewVote) error {
	return r.db.WithContext(ctx).Create(vote).Error
}

// GetByID gets a review vote by ID
func (r *reviewVoteRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ReviewVote, error) {
	var vote entities.ReviewVote
	err := r.db.WithContext(ctx).First(&vote, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

// GetByUserAndReview gets a vote by user and review
func (r *reviewVoteRepository) GetByUserAndReview(ctx context.Context, userID, reviewID uuid.UUID) (*entities.ReviewVote, error) {
	var vote entities.ReviewVote
	err := r.db.WithContext(ctx).
		First(&vote, "user_id = ? AND review_id = ?", userID, reviewID).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

// Update updates a review vote
func (r *reviewVoteRepository) Update(ctx context.Context, vote *entities.ReviewVote) error {
	vote.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(vote).Error
}

// Delete deletes a review vote
func (r *reviewVoteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.ReviewVote{}, "id = ?", id).Error
}

// DeleteByUserAndReview deletes a vote by user and review
func (r *reviewVoteRepository) DeleteByUserAndReview(ctx context.Context, userID, reviewID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&entities.ReviewVote{}, "user_id = ? AND review_id = ?", userID, reviewID).Error
}

// GetVotesByReview gets all votes for a review
func (r *reviewVoteRepository) GetVotesByReview(ctx context.Context, reviewID uuid.UUID) ([]*entities.ReviewVote, error) {
	var votes []*entities.ReviewVote
	err := r.db.WithContext(ctx).
		Where("review_id = ?", reviewID).
		Find(&votes).Error
	return votes, err
}

// GetVotesByUser gets all votes by a user
func (r *reviewVoteRepository) GetVotesByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.ReviewVote, error) {
	var votes []*entities.ReviewVote
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&votes).Error
	return votes, err
}

// CountVotesByReview counts votes for a review by type
func (r *reviewVoteRepository) CountVotesByReview(ctx context.Context, reviewID uuid.UUID, voteType entities.ReviewVoteType) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.ReviewVote{}).
		Where("review_id = ? AND vote_type = ?", reviewID, voteType).
		Count(&count).Error
	return count, err
}

// GetVoteCounts gets vote counts for a review
func (r *reviewVoteRepository) GetVoteCounts(ctx context.Context, reviewID uuid.UUID) (helpful int, notHelpful int, err error) {
	var helpfulCount, notHelpfulCount int64

	// Get helpful votes
	err = r.db.WithContext(ctx).
		Model(&entities.ReviewVote{}).
		Where("review_id = ? AND vote_type = ?", reviewID, entities.ReviewVoteHelpful).
		Count(&helpfulCount).Error
	if err != nil {
		return 0, 0, err
	}

	// Get not helpful votes
	err = r.db.WithContext(ctx).
		Model(&entities.ReviewVote{}).
		Where("review_id = ? AND vote_type = ?", reviewID, entities.ReviewVoteNotHelpful).
		Count(&notHelpfulCount).Error
	if err != nil {
		return 0, 0, err
	}

	return int(helpfulCount), int(notHelpfulCount), nil
}

// HasUserVoted checks if a user has voted on a review
func (r *reviewVoteRepository) HasUserVoted(ctx context.Context, userID, reviewID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.ReviewVote{}).
		Where("user_id = ? AND review_id = ?", userID, reviewID).
		Count(&count).Error
	return count > 0, err
}

// GetUserVoteType gets the vote type for a user on a review
func (r *reviewVoteRepository) GetUserVoteType(ctx context.Context, userID, reviewID uuid.UUID) (*entities.ReviewVoteType, error) {
	var vote entities.ReviewVote
	err := r.db.WithContext(ctx).
		Select("vote_type").
		First(&vote, "user_id = ? AND review_id = ?", userID, reviewID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &vote.VoteType, nil
}

// List lists review votes with filters
func (r *reviewVoteRepository) List(ctx context.Context, filters repositories.ReviewVoteFilters) ([]*entities.ReviewVote, error) {
	var votes []*entities.ReviewVote
	query := r.db.WithContext(ctx)

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.ReviewID != nil {
		query = query.Where("review_id = ?", *filters.ReviewID)
	}

	if filters.VoteType != nil {
		query = query.Where("vote_type = ?", *filters.VoteType)
	}

	if filters.CreatedAfter != nil {
		query = query.Where("created_at >= ?", *filters.CreatedAfter)
	}

	if filters.CreatedBefore != nil {
		query = query.Where("created_at <= ?", *filters.CreatedBefore)
	}

	// Apply sorting
	switch filters.SortBy {
	case "created_at":
		if filters.SortOrder == "desc" {
			query = query.Order("created_at DESC")
		} else {
			query = query.Order("created_at ASC")
		}
	case "vote_type":
		if filters.SortOrder == "desc" {
			query = query.Order("vote_type DESC")
		} else {
			query = query.Order("vote_type ASC")
		}
	default:
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&votes).Error
	return votes, err
}

// Count counts review votes with filters
func (r *reviewVoteRepository) Count(ctx context.Context, filters repositories.ReviewVoteFilters) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.ReviewVote{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.ReviewID != nil {
		query = query.Where("review_id = ?", *filters.ReviewID)
	}

	if filters.VoteType != nil {
		query = query.Where("vote_type = ?", *filters.VoteType)
	}

	if filters.CreatedAfter != nil {
		query = query.Where("created_at >= ?", *filters.CreatedAfter)
	}

	if filters.CreatedBefore != nil {
		query = query.Where("created_at <= ?", *filters.CreatedBefore)
	}

	err := query.Count(&count).Error
	return count, err
}

// GetMostHelpfulReviews gets reviews with the most helpful votes
func (r *reviewVoteRepository) GetMostHelpfulReviews(ctx context.Context, productID *uuid.UUID, limit int) ([]uuid.UUID, error) {
	var reviewIDs []uuid.UUID
	query := r.db.WithContext(ctx).
		Table("review_votes").
		Select("review_id").
		Where("vote_type = ?", entities.ReviewVoteHelpful).
		Group("review_id").
		Order("COUNT(*) DESC").
		Limit(limit)

	if productID != nil {
		query = query.Joins("JOIN reviews ON review_votes.review_id = reviews.id").
			Where("reviews.product_id = ?", *productID)
	}

	err := query.Pluck("review_id", &reviewIDs).Error
	return reviewIDs, err
}

// DeleteByReview deletes all votes for a review
func (r *reviewVoteRepository) DeleteByReview(ctx context.Context, reviewID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.ReviewVote{}, "review_id = ?", reviewID).Error
}

// DeleteByUser deletes all votes by a user
func (r *reviewVoteRepository) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.ReviewVote{}, "user_id = ?", userID).Error
}

// CountHelpfulVotes counts helpful votes for a review
func (r *reviewVoteRepository) CountHelpfulVotes(ctx context.Context, reviewID uuid.UUID) (int, error) {
	count, err := r.CountVotesByReview(ctx, reviewID, entities.ReviewVoteHelpful)
	return int(count), err
}

// CountNotHelpfulVotes counts not helpful votes for a review
func (r *reviewVoteRepository) CountNotHelpfulVotes(ctx context.Context, reviewID uuid.UUID) (int, error) {
	count, err := r.CountVotesByReview(ctx, reviewID, entities.ReviewVoteNotHelpful)
	return int(count), err
}

// GetUserVote gets a user's vote for a specific review
func (r *reviewVoteRepository) GetUserVote(ctx context.Context, userID, reviewID uuid.UUID) (*entities.ReviewVote, error) {
	var vote entities.ReviewVote
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND review_id = ?", userID, reviewID).
		First(&vote).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

// GetUserVotesForReviews gets user votes for multiple reviews as a map
func (r *reviewVoteRepository) GetUserVotesForReviews(ctx context.Context, userID uuid.UUID, reviewIDs []uuid.UUID) (map[uuid.UUID]*entities.ReviewVote, error) {
	var votes []*entities.ReviewVote
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND review_id IN (?)", userID, reviewIDs).
		Find(&votes).Error
	if err != nil {
		return nil, err
	}

	// Convert to map
	voteMap := make(map[uuid.UUID]*entities.ReviewVote)
	for _, vote := range votes {
		voteMap[vote.ReviewID] = vote
	}

	return voteMap, nil
}

// GetVotesByReviewIDs gets votes for multiple reviews grouped by review ID
func (r *reviewVoteRepository) GetVotesByReviewIDs(ctx context.Context, reviewIDs []uuid.UUID) (map[uuid.UUID][]*entities.ReviewVote, error) {
	var votes []*entities.ReviewVote
	err := r.db.WithContext(ctx).
		Where("review_id IN (?)", reviewIDs).
		Find(&votes).Error
	if err != nil {
		return nil, err
	}

	// Group votes by review ID
	voteMap := make(map[uuid.UUID][]*entities.ReviewVote)
	for _, vote := range votes {
		voteMap[vote.ReviewID] = append(voteMap[vote.ReviewID], vote)
	}

	return voteMap, nil
}

// RemoveVote removes a user's vote for a review
func (r *reviewVoteRepository) RemoveVote(ctx context.Context, userID, reviewID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&entities.ReviewVote{}, "user_id = ? AND review_id = ?", userID, reviewID).Error
}

// UpdateReviewVoteCounts updates vote counts for a review
func (r *reviewVoteRepository) UpdateReviewVoteCounts(ctx context.Context, reviewID uuid.UUID) error {
	// This would typically update cached vote counts in the review table
	// For now, it's a placeholder
	return nil
}

// VoteReview adds or updates a vote for a review
func (r *reviewVoteRepository) VoteReview(ctx context.Context, reviewID, userID uuid.UUID, voteType entities.ReviewVoteType) error {
	// Check if vote already exists
	var existingVote entities.ReviewVote
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND review_id = ?", userID, reviewID).
		First(&existingVote).Error

	if err == gorm.ErrRecordNotFound {
		// Create new vote
		vote := &entities.ReviewVote{
			ID:       uuid.New(),
			UserID:   userID,
			ReviewID: reviewID,
			VoteType: voteType,
			CreatedAt: time.Now(),
		}
		return r.db.WithContext(ctx).Create(vote).Error
	} else if err != nil {
		return err
	} else {
		// Update existing vote
		existingVote.VoteType = voteType
		existingVote.UpdatedAt = time.Now()
		return r.db.WithContext(ctx).Save(&existingVote).Error
	}
}
