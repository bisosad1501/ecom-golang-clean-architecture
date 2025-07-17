package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
)

type checkoutSessionRepository struct {
	db *gorm.DB
}

// NewCheckoutSessionRepository creates a new checkout session repository
func NewCheckoutSessionRepository(db *gorm.DB) repositories.CheckoutSessionRepository {
	return &checkoutSessionRepository{db: db}
}

// Create creates a new checkout session
func (r *checkoutSessionRepository) Create(ctx context.Context, session *entities.CheckoutSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetByID retrieves a checkout session by ID
func (r *checkoutSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.CheckoutSession, error) {
	var session entities.CheckoutSession
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&session, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetBySessionID retrieves a checkout session by session ID
func (r *checkoutSessionRepository) GetBySessionID(ctx context.Context, sessionID string) (*entities.CheckoutSession, error) {
	var session entities.CheckoutSession
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&session, "session_id = ?", sessionID).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetByUserID retrieves active checkout sessions for a user
func (r *checkoutSessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.CheckoutSession, error) {
	var sessions []*entities.CheckoutSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, entities.CheckoutSessionStatusActive).
		Order("created_at DESC").
		Find(&sessions).Error
	return sessions, err
}

// Update updates a checkout session
func (r *checkoutSessionRepository) Update(ctx context.Context, session *entities.CheckoutSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// Delete deletes a checkout session
func (r *checkoutSessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.CheckoutSession{}, "id = ?", id).Error
}

// GetExpiredSessions retrieves expired checkout sessions for cleanup
func (r *checkoutSessionRepository) GetExpiredSessions(ctx context.Context, limit int) ([]*entities.CheckoutSession, error) {
	var sessions []*entities.CheckoutSession
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("status = ? AND expires_at < ?", entities.CheckoutSessionStatusActive, now).
		Limit(limit).
		Find(&sessions).Error
	return sessions, err
}

// MarkAsExpired marks checkout sessions as expired
func (r *checkoutSessionRepository) MarkAsExpired(ctx context.Context, ids []uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.CheckoutSession{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":     entities.CheckoutSessionStatusExpired,
			"updated_at": time.Now(),
		}).Error
}
