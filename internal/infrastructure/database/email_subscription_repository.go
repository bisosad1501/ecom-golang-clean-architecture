package database

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type emailSubscriptionRepository struct {
	db *gorm.DB
}

// NewEmailSubscriptionRepository creates a new email subscription repository
func NewEmailSubscriptionRepository(db *gorm.DB) repositories.EmailSubscriptionRepository {
	return &emailSubscriptionRepository{db: db}
}

// Create creates a new email subscription
func (r *emailSubscriptionRepository) Create(ctx context.Context, subscription *entities.EmailSubscription) error {
	return r.db.WithContext(ctx).Create(subscription).Error
}

// GetByID gets an email subscription by ID
func (r *emailSubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.EmailSubscription, error) {
	var subscription entities.EmailSubscription
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// GetByUserID gets an email subscription by user ID
func (r *emailSubscriptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.EmailSubscription, error) {
	var subscription entities.EmailSubscription
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// Update updates an email subscription
func (r *emailSubscriptionRepository) Update(ctx context.Context, subscription *entities.EmailSubscription) error {
	return r.db.WithContext(ctx).Save(subscription).Error
}

// Delete deletes an email subscription
func (r *emailSubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.EmailSubscription{}, id).Error
}

// List lists email subscriptions with pagination
func (r *emailSubscriptionRepository) List(ctx context.Context, offset, limit int) ([]*entities.EmailSubscription, error) {
	var subscriptions []*entities.EmailSubscription
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&subscriptions).Error
	return subscriptions, err
}

// GetSubscribedUsers gets users subscribed to a specific email type
func (r *emailSubscriptionRepository) GetSubscribedUsers(ctx context.Context, emailType entities.EmailType) ([]uuid.UUID, error) {
	var userIDs []uuid.UUID

	var whereClause string
	switch emailType {
	case entities.EmailTypeNewsletter:
		whereClause = "newsletter = true"
	case entities.EmailTypePromotion:
		whereClause = "promotions = true"
	case entities.EmailTypeOrderConfirmation, entities.EmailTypeOrderShipped,
		 entities.EmailTypeOrderDelivered, entities.EmailTypeOrderCancelled:
		whereClause = "order_updates = true"
	case entities.EmailTypeReviewRequest:
		whereClause = "review_requests = true"
	case entities.EmailTypeAbandonedCart:
		whereClause = "abandoned_cart = true"
	case entities.EmailTypeSupport:
		whereClause = "support = true"
	default:
		// For system emails, assume all users are subscribed
		whereClause = "1 = 1"
	}

	err := r.db.WithContext(ctx).
		Model(&entities.EmailSubscription{}).
		Where(whereClause).
		Pluck("user_id", &userIDs).Error

	return userIDs, err
}

// GetUnsubscribedUsers gets users unsubscribed from a specific email type
func (r *emailSubscriptionRepository) GetUnsubscribedUsers(ctx context.Context, emailType entities.EmailType) ([]uuid.UUID, error) {
	var userIDs []uuid.UUID

	var whereClause string
	switch emailType {
	case entities.EmailTypeNewsletter:
		whereClause = "newsletter = false"
	case entities.EmailTypePromotion:
		whereClause = "promotions = false"
	case entities.EmailTypeOrderConfirmation, entities.EmailTypeOrderShipped,
		 entities.EmailTypeOrderDelivered, entities.EmailTypeOrderCancelled:
		whereClause = "order_updates = false"
	case entities.EmailTypeReviewRequest:
		whereClause = "review_requests = false"
	case entities.EmailTypeAbandonedCart:
		whereClause = "abandoned_cart = false"
	case entities.EmailTypeSupport:
		whereClause = "support = false"
	default:
		// For system emails, return empty list (all users are subscribed)
		return userIDs, nil
	}

	err := r.db.WithContext(ctx).
		Model(&entities.EmailSubscription{}).
		Where(whereClause).
		Pluck("user_id", &userIDs).Error

	return userIDs, err
}

// UpdateSubscriptions updates multiple subscriptions for a user
func (r *emailSubscriptionRepository) UpdateSubscriptions(ctx context.Context, userID uuid.UUID, subscriptions map[entities.EmailType]bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Find or create subscription record for user
		var subscription entities.EmailSubscription
		err := tx.Where("user_id = ?", userID).First(&subscription).Error

		if err == gorm.ErrRecordNotFound {
			// Create new subscription with default values
			subscription = entities.EmailSubscription{
				ID:             uuid.New(),
				UserID:         userID,
				Newsletter:     true,
				Promotions:     true,
				OrderUpdates:   true,
				ReviewRequests: true,
				AbandonedCart:  true,
				Support:        true,
			}
		} else if err != nil {
			return err
		}

		// Update subscription preferences based on email types
		for emailType, isSubscribed := range subscriptions {
			switch emailType {
			case entities.EmailTypeNewsletter:
				subscription.Newsletter = isSubscribed
			case entities.EmailTypePromotion:
				subscription.Promotions = isSubscribed
			case entities.EmailTypeOrderConfirmation, entities.EmailTypeOrderShipped,
				 entities.EmailTypeOrderDelivered, entities.EmailTypeOrderCancelled:
				subscription.OrderUpdates = isSubscribed
			case entities.EmailTypeReviewRequest:
				subscription.ReviewRequests = isSubscribed
			case entities.EmailTypeAbandonedCart:
				subscription.AbandonedCart = isSubscribed
			case entities.EmailTypeSupport:
				subscription.Support = isSubscribed
			}
		}

		// Save the subscription
		if err == gorm.ErrRecordNotFound {
			return tx.Create(&subscription).Error
		} else {
			return tx.Save(&subscription).Error
		}
	})
}
