package repositories

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"

	"github.com/google/uuid"
)

// CartRepository defines the interface for cart data access
type CartRepository interface {
	// Create creates a new cart
	Create(ctx context.Context, cart *entities.Cart) error

	// GetByID retrieves a cart by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Cart, error)

	// GetByUserID retrieves a cart by user ID
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.Cart, error)

	// GetBySessionID retrieves a cart by session ID (for guest users)
	GetBySessionID(ctx context.Context, sessionID string) (*entities.Cart, error)

	// GetBySessionIDForUpdate retrieves a cart by session ID with row-level locking
	GetBySessionIDForUpdate(ctx context.Context, sessionID string) (*entities.Cart, error)

	// GetByUserIDForUpdate retrieves a cart by user ID with row-level locking
	GetByUserIDForUpdate(ctx context.Context, userID uuid.UUID) (*entities.Cart, error)

	// WithTransaction executes a function within a database transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error)

	// Update updates an existing cart
	Update(ctx context.Context, cart *entities.Cart) error

	// Delete deletes a cart by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// AddItem adds an item to the cart
	AddItem(ctx context.Context, cartID uuid.UUID, item *entities.CartItem) error

	// UpdateItem updates a cart item
	UpdateItem(ctx context.Context, item *entities.CartItem) error

	// RemoveItem removes an item from the cart
	RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error

	// GetItem retrieves a cart item
	GetItem(ctx context.Context, cartID, productID uuid.UUID) (*entities.CartItem, error)

	// ClearCart removes all items from the cart
	ClearCart(ctx context.Context, cartID uuid.UUID) error

	// GetItems retrieves all items in a cart
	GetItems(ctx context.Context, cartID uuid.UUID) ([]*entities.CartItem, error)

	// RemoveItemsByProductID removes all cart items with the specified product ID
	RemoveItemsByProductID(ctx context.Context, productID uuid.UUID) error

	// GetExpiredCarts retrieves all expired carts
	GetExpiredCarts(ctx context.Context) ([]*entities.Cart, error)
}
