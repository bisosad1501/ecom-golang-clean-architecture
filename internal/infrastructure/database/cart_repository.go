package database

import (
	"context"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a new cart repository
func NewCartRepository(db *gorm.DB) repositories.CartRepository {
	return &cartRepository{db: db}
}

// updateCartCalculatedFields efficiently updates cart calculated fields using SQL
func (r *cartRepository) updateCartCalculatedFields(ctx context.Context, tx *gorm.DB, cartID uuid.UUID) error {
	// Calculate values using SQL aggregation for better performance
	var result struct {
		Subtotal  float64
		ItemCount int
	}

	err := tx.WithContext(ctx).
		Table("cart_items ci").
		Select("COALESCE(SUM(ci.total), 0) as subtotal, COALESCE(SUM(ci.quantity), 0) as item_count").
		Where("ci.cart_id = ?", cartID).
		Scan(&result).Error

	if err != nil {
		return fmt.Errorf("failed to calculate cart totals: %w", err)
	}

	// Get current cart to preserve tax and shipping amounts
	var currentCart entities.Cart
	if err := tx.WithContext(ctx).Select("tax_amount", "shipping_amount").Where("id = ?", cartID).First(&currentCart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("cart not found: %s", cartID)
		}
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// Calculate total including tax and shipping
	total := result.Subtotal + currentCart.TaxAmount + currentCart.ShippingAmount

	// Validate calculated values
	if result.Subtotal < 0 {
		return fmt.Errorf("calculated subtotal cannot be negative: %.2f", result.Subtotal)
	}
	if result.ItemCount < 0 {
		return fmt.Errorf("calculated item count cannot be negative: %d", result.ItemCount)
	}

	// Update cart with calculated values only if they changed
	updateResult := tx.WithContext(ctx).
		Model(&entities.Cart{}).
		Where("id = ?", cartID).
		Where("subtotal != ? OR total != ? OR item_count != ?", result.Subtotal, total, result.ItemCount).
		Updates(map[string]interface{}{
			"subtotal":   result.Subtotal,
			"total":      total,
			"item_count": result.ItemCount,
			"updated_at": time.Now(),
		})

	if updateResult.Error != nil {
		return fmt.Errorf("failed to update cart calculated fields: %w", updateResult.Error)
	}

	return nil
}

// Create creates a new cart
func (r *cartRepository) Create(ctx context.Context, cart *entities.Cart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

// GetByID retrieves a cart by ID
func (r *cartRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Cart, error) {
	var cart entities.Cart
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Images").
		Where("id = ?", id).
		First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrCartNotFound
		}
		return nil, err
	}
	return &cart, nil
}

// GetByUserID retrieves a cart by user ID
func (r *cartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.Cart, error) {
	var cart entities.Cart
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Images").
		Where("user_id = ? AND status = ?", userID, "active").
		First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrCartNotFound
		}
		return nil, err
	}

	// Calculated fields are updated by updateCartCalculatedFields after item operations
	// No need to call UpdateCalculatedFields here as it's done by the SQL aggregation

	return &cart, nil
}

// GetBySessionID retrieves a cart by session ID (for guest users)
func (r *cartRepository) GetBySessionID(ctx context.Context, sessionID string) (*entities.Cart, error) {
	var cart entities.Cart
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Images").
		Where("session_id = ? AND status = ?", sessionID, "active").
		First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrCartNotFound
		}
		return nil, err
	}

	// Calculated fields are updated by updateCartCalculatedFields after item operations
	// No need to call UpdateCalculatedFields here as it's done by the SQL aggregation

	return &cart, nil
}

// GetBySessionIDForUpdate retrieves a cart by session ID with row-level locking
func (r *cartRepository) GetBySessionIDForUpdate(ctx context.Context, sessionID string) (*entities.Cart, error) {
	var cart entities.Cart
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Images").
		Where("session_id = ? AND status = ?", sessionID, "active").
		Set("gorm:query_option", "FOR UPDATE").
		First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrCartNotFound
		}
		return nil, err
	}
	return &cart, nil
}

// GetByUserIDForUpdate retrieves a cart by user ID with row-level locking
func (r *cartRepository) GetByUserIDForUpdate(ctx context.Context, userID uuid.UUID) (*entities.Cart, error) {
	var cart entities.Cart
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Images").
		Where("user_id = ? AND status = ?", userID, "active").
		Set("gorm:query_option", "FOR UPDATE").
		First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrCartNotFound
		}
		return nil, err
	}
	return &cart, nil
}

// WithTransaction executes a function within a database transaction
func (r *cartRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	var result interface{}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new repository instance with the transaction
		txRepo := &cartRepository{db: tx}

		// Create a new context with the transaction repository
		txCtx := context.WithValue(ctx, "tx_repo", txRepo)

		var err error
		result, err = fn(txCtx)
		return err
	})
	return result, err
}

// Update updates an existing cart
func (r *cartRepository) Update(ctx context.Context, cart *entities.Cart) error {
	// Calculated fields are updated by updateCartCalculatedFields after item operations
	// No need to call UpdateCalculatedFields here as it's done by the SQL aggregation
	return r.db.WithContext(ctx).Save(cart).Error
}

// Delete deletes a cart by ID
func (r *cartRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.Cart{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrCartNotFound
	}
	return nil
}

// AddItem adds an item to the cart
func (r *cartRepository) AddItem(ctx context.Context, cartID uuid.UUID, item *entities.CartItem) error {
	item.CartID = cartID

	// Use transaction to ensure consistency
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the item
		if err := tx.Create(item).Error; err != nil {
			return err
		}

		// Update cart calculated fields efficiently
		return r.updateCartCalculatedFields(ctx, tx, cartID)
	})
}

// UpdateItem updates a cart item
func (r *cartRepository) UpdateItem(ctx context.Context, item *entities.CartItem) error {
	// Use transaction to ensure consistency
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update the item
		if err := tx.Save(item).Error; err != nil {
			return err
		}

		// Update cart calculated fields efficiently
		return r.updateCartCalculatedFields(ctx, tx, item.CartID)
	})
}

// RemoveItem removes an item from the cart
func (r *cartRepository) RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error {
	// Use transaction to ensure consistency
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("cart_id = ? AND product_id = ?", cartID, productID).
			Delete(&entities.CartItem{})

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return entities.ErrCartItemNotFound
		}

		// Update cart calculated fields efficiently
		return r.updateCartCalculatedFields(ctx, tx, cartID)
	})
}

// GetItem retrieves a cart item
func (r *cartRepository) GetItem(ctx context.Context, cartID, productID uuid.UUID) (*entities.CartItem, error) {
	var item entities.CartItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrCartItemNotFound
		}
		return nil, err
	}
	return &item, nil
}

// ClearCart removes all items from the cart
func (r *cartRepository) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("cart_id = ?", cartID).
		Delete(&entities.CartItem{}).Error
}

// GetItems retrieves all items in a cart
func (r *cartRepository) GetItems(ctx context.Context, cartID uuid.UUID) ([]*entities.CartItem, error) {
	var items []*entities.CartItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Category").
		Preload("Product.Images").
		Where("cart_id = ?", cartID).
		Find(&items).Error
	return items, err
}

// RemoveItemsByProductID removes all cart items with the specified product ID
func (r *cartRepository) RemoveItemsByProductID(ctx context.Context, productID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Delete(&entities.CartItem{}).Error
}

// GetExpiredCarts retrieves all expired carts
func (r *cartRepository) GetExpiredCarts(ctx context.Context) ([]*entities.Cart, error) {
	var carts []*entities.Cart
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Where("status = ? AND expires_at < ?", "active", time.Now()).
		Find(&carts).Error
	if err != nil {
		return nil, err
	}
	return carts, nil
}
