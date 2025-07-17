package repositories

import (
	"context"

	"github.com/google/uuid"
	"ecom-golang-clean-architecture/internal/domain/entities"
)

// CheckoutSessionRepository defines the interface for checkout session data access
type CheckoutSessionRepository interface {
	// Create creates a new checkout session
	Create(ctx context.Context, session *entities.CheckoutSession) error

	// GetByID retrieves a checkout session by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.CheckoutSession, error)

	// GetBySessionID retrieves a checkout session by session ID
	GetBySessionID(ctx context.Context, sessionID string) (*entities.CheckoutSession, error)

	// GetByUserID retrieves active checkout sessions for a user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.CheckoutSession, error)

	// Update updates a checkout session
	Update(ctx context.Context, session *entities.CheckoutSession) error

	// Delete deletes a checkout session
	Delete(ctx context.Context, id uuid.UUID) error

	// GetExpiredSessions retrieves expired checkout sessions for cleanup
	GetExpiredSessions(ctx context.Context, limit int) ([]*entities.CheckoutSession, error)

	// MarkAsExpired marks checkout sessions as expired
	MarkAsExpired(ctx context.Context, ids []uuid.UUID) error
}
