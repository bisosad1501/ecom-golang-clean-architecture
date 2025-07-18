package database

import (
	"context"
	"database/sql"
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
		Where("user_id = ? AND status = ?", userID, "active").
		Order("created_at DESC").
		First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrCartNotFound
		}
		return nil, err
	}

	return &cart, nil
}

// GetBySessionID retrieves a cart by session ID (for guest users)
func (r *cartRepository) GetBySessionID(ctx context.Context, sessionID string) (*entities.Cart, error) {
	var cart entities.Cart
	err := r.db.WithContext(ctx).
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Items.Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, description, sku, slug, price, current_price, stock, status, category_id, created_at, updated_at")
		}).
		Preload("Items.Product.Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, description, slug, image")
		}).
		Preload("Items.Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC").Select("id, product_id, url, alt_text, position")
		}).
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

// GetAbandonedCarts retrieves carts that haven't been updated since the given time
func (r *cartRepository) GetAbandonedCarts(ctx context.Context, since time.Time) ([]*entities.Cart, error) {
	var carts []*entities.Cart

	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Where("updated_at < ? AND is_abandoned = false", since).
		Find(&carts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get abandoned carts: %w", err)
	}

	return carts, nil
}

// GetAbandonedCartsList retrieves abandoned carts with pagination
func (r *cartRepository) GetAbandonedCartsList(ctx context.Context, offset, limit int) ([]*entities.Cart, error) {
	var carts []*entities.Cart

	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Where("is_abandoned = true").
		Order("abandoned_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&carts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get abandoned carts list: %w", err)
	}

	return carts, nil
}

// GetAbandonedCartStats retrieves statistics for abandoned carts
func (r *cartRepository) GetAbandonedCartStats(ctx context.Context, since time.Time) (*repositories.AbandonedCartStats, error) {
	var stats repositories.AbandonedCartStats

	// Get total abandoned carts
	err := r.db.WithContext(ctx).
		Model(&entities.Cart{}).
		Where("is_abandoned = true AND abandoned_at >= ?", since).
		Count(&stats.TotalAbandoned).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count abandoned carts: %w", err)
	}

	// Get total recovered carts
	err = r.db.WithContext(ctx).
		Model(&entities.Cart{}).
		Where("is_abandoned = false AND recovered_at >= ?", since).
		Count(&stats.TotalRecovered).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count recovered carts: %w", err)
	}

	// Calculate recovery rate
	if stats.TotalAbandoned > 0 {
		stats.RecoveryRate = float64(stats.TotalRecovered) / float64(stats.TotalAbandoned) * 100
	}

	// Get average cart value for abandoned carts
	var avgValue sql.NullFloat64
	err = r.db.WithContext(ctx).
		Model(&entities.Cart{}).
		Select("AVG(total)").
		Where("is_abandoned = true AND abandoned_at >= ?", since).
		Scan(&avgValue).Error
	if err != nil {
		return nil, fmt.Errorf("failed to calculate average cart value: %w", err)
	}
	if avgValue.Valid {
		stats.AverageCartValue = avgValue.Float64
	}

	// Calculate total lost revenue
	var totalLost sql.NullFloat64
	err = r.db.WithContext(ctx).
		Model(&entities.Cart{}).
		Select("SUM(total)").
		Where("is_abandoned = true AND abandoned_at >= ?", since).
		Scan(&totalLost).Error
	if err != nil {
		return nil, fmt.Errorf("failed to calculate total lost revenue: %w", err)
	}
	if totalLost.Valid {
		stats.TotalLostRevenue = totalLost.Float64
	}

	// Calculate recovered revenue
	var recoveredRevenue sql.NullFloat64
	err = r.db.WithContext(ctx).
		Model(&entities.Cart{}).
		Select("SUM(total)").
		Where("is_abandoned = false AND recovered_at >= ?", since).
		Scan(&recoveredRevenue).Error
	if err != nil {
		return nil, fmt.Errorf("failed to calculate recovered revenue: %w", err)
	}
	if recoveredRevenue.Valid {
		stats.RecoveredRevenue = recoveredRevenue.Float64
	}

	// Count reminder emails sent
	err = r.db.WithContext(ctx).
		Model(&entities.Cart{}).
		Where("first_reminder_sent IS NOT NULL AND abandoned_at >= ?", since).
		Count(&stats.FirstReminderSent).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count first reminders: %w", err)
	}

	err = r.db.WithContext(ctx).
		Model(&entities.Cart{}).
		Where("second_reminder_sent IS NOT NULL AND abandoned_at >= ?", since).
		Count(&stats.SecondReminderSent).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count second reminders: %w", err)
	}

	err = r.db.WithContext(ctx).
		Model(&entities.Cart{}).
		Where("final_reminder_sent IS NOT NULL AND abandoned_at >= ?", since).
		Count(&stats.FinalReminderSent).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count final reminders: %w", err)
	}

	return &stats, nil
}
