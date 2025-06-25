package repositories

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entities.User) error
	
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	
	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)

	// GetByGoogleID retrieves a user by Google ID
	GetByGoogleID(ctx context.Context, googleID string) (*entities.User, error)

	// GetByFacebookID retrieves a user by Facebook ID
	GetByFacebookID(ctx context.Context, facebookID string) (*entities.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entities.User) error
	
	// Delete deletes a user by ID
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List retrieves users with pagination
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)
	
	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)
	
	// ExistsByEmail checks if a user exists with the given email
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	
	// UpdatePassword updates user password
	UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error
	
	// SetActive sets user active status
	SetActive(ctx context.Context, userID uuid.UUID, isActive bool) error

	// Additional methods for admin dashboard
	CountUsers(ctx context.Context) (int64, error)
	CountActiveUsers(ctx context.Context) (int64, error)
}

// UserProfileRepository defines the interface for user profile data access
type UserProfileRepository interface {
	// Create creates a new user profile
	Create(ctx context.Context, profile *entities.UserProfile) error
	
	// GetByUserID retrieves a user profile by user ID
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.UserProfile, error)
	
	// Update updates an existing user profile
	Update(ctx context.Context, profile *entities.UserProfile) error
	
	// Delete deletes a user profile by user ID
	Delete(ctx context.Context, userID uuid.UUID) error
}
