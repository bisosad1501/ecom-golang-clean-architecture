package database

import (
	"context"

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
		Where("user_id = ?", userID).
		First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrCartNotFound
		}
		return nil, err
	}
	return &cart, nil
}

// Update updates an existing cart
func (r *cartRepository) Update(ctx context.Context, cart *entities.Cart) error {
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
	return r.db.WithContext(ctx).Create(item).Error
}

// UpdateItem updates a cart item
func (r *cartRepository) UpdateItem(ctx context.Context, item *entities.CartItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// RemoveItem removes an item from the cart
func (r *cartRepository) RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		Delete(&entities.CartItem{})
	
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrCartItemNotFound
	}
	return nil
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
