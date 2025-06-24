package repositories

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// AddressRepository defines the interface for address data access
type AddressRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, address *entities.Address) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Address, error)
	Update(ctx context.Context, address *entities.Address) error
	Delete(ctx context.Context, id uuid.UUID) error

	// User-specific operations
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Address, error)
	GetDefaultByUserID(ctx context.Context, userID uuid.UUID, addressType entities.AddressType) (*entities.Address, error)
	SetAsDefault(ctx context.Context, userID, addressID uuid.UUID, addressType entities.AddressType) error
	GetByUserIDAndType(ctx context.Context, userID uuid.UUID, addressType entities.AddressType) ([]*entities.Address, error)

	// Validation
	ExistsByUserIDAndID(ctx context.Context, userID, addressID uuid.UUID) (bool, error)
}

// WishlistRepository defines the interface for wishlist data access
type WishlistRepository interface {
	// Basic operations
	AddToWishlist(ctx context.Context, userID, productID uuid.UUID) error
	RemoveFromWishlist(ctx context.Context, userID, productID uuid.UUID) error
	IsInWishlist(ctx context.Context, userID, productID uuid.UUID) (bool, error)

	// Get operations
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Wishlist, error)
	GetWishlistProducts(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Product, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)

	// Bulk operations
	ClearWishlist(ctx context.Context, userID uuid.UUID) error
	GetWishlistProductIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

// UserPreferenceRepository defines the interface for user preferences data access
type UserPreferenceRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, preference *entities.UserPreference) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.UserPreference, error)
	Update(ctx context.Context, preference *entities.UserPreference) error
	Delete(ctx context.Context, userID uuid.UUID) error

	// Specific updates
	UpdateNotificationSettings(ctx context.Context, userID uuid.UUID, settings map[string]bool) error
	UpdateLanguageAndCurrency(ctx context.Context, userID uuid.UUID, language, currency string) error
}

// AccountVerificationRepository defines the interface for account verification data access
type AccountVerificationRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, verification *entities.AccountVerification) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.AccountVerification, error)
	Update(ctx context.Context, verification *entities.AccountVerification) error

	// Verification operations
	VerifyEmail(ctx context.Context, userID uuid.UUID) error
	VerifyPhone(ctx context.Context, userID uuid.UUID) error
	SetVerificationCode(ctx context.Context, userID uuid.UUID, code string, expiresAt time.Time) error
	GetByVerificationCode(ctx context.Context, code string) (*entities.AccountVerification, error)

	// Status checks
	IsEmailVerified(ctx context.Context, userID uuid.UUID) (bool, error)
	IsPhoneVerified(ctx context.Context, userID uuid.UUID) (bool, error)
}

// PasswordResetRepository defines the interface for password reset data access
type PasswordResetRepository interface {
	// Basic operations
	Create(ctx context.Context, reset *entities.PasswordReset) error
	GetByToken(ctx context.Context, token string) (*entities.PasswordReset, error)
	MarkAsUsed(ctx context.Context, token string) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Cleanup operations
	DeleteExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error

	// Validation
	IsTokenValid(ctx context.Context, token string) (bool, error)
}
