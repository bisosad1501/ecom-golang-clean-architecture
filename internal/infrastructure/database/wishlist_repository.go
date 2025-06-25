package database

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type wishlistRepository struct {
	db *gorm.DB
}

// NewWishlistRepository creates a new wishlist repository
func NewWishlistRepository(db *gorm.DB) repositories.WishlistRepository {
	return &wishlistRepository{db: db}
}

// Create creates a new wishlist item
func (r *wishlistRepository) Create(ctx context.Context, wishlist *entities.Wishlist) error {
	return r.db.WithContext(ctx).Create(wishlist).Error
}

// GetByID gets a wishlist item by ID
func (r *wishlistRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Wishlist, error) {
	var wishlist entities.Wishlist
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Images").
		Preload("Product.Category").
		First(&wishlist, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &wishlist, nil
}

// GetByUserAndProduct gets a wishlist item by user and product
func (r *wishlistRepository) GetByUserAndProduct(ctx context.Context, userID, productID uuid.UUID) (*entities.Wishlist, error) {
	var wishlist entities.Wishlist
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Images").
		Preload("Product.Category").
		First(&wishlist, "user_id = ? AND product_id = ?", userID, productID).Error
	if err != nil {
		return nil, err
	}
	return &wishlist, nil
}

// GetByUser gets all wishlist items for a user
func (r *wishlistRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Wishlist, error) {
	var wishlists []*entities.Wishlist
	query := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Images").
		Preload("Product.Category").
		Where("user_id = ?", userID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&wishlists).Error
	return wishlists, err
}

// CountByUser counts wishlist items for a user
func (r *wishlistRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Wishlist{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// Delete deletes a wishlist item
func (r *wishlistRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Wishlist{}, "id = ?", id).Error
}

// DeleteByUserAndProduct deletes a wishlist item by user and product
func (r *wishlistRepository) DeleteByUserAndProduct(ctx context.Context, userID, productID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&entities.Wishlist{}, "user_id = ? AND product_id = ?", userID, productID).Error
}

// ClearByUser clears all wishlist items for a user
func (r *wishlistRepository) ClearByUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&entities.Wishlist{}, "user_id = ?", userID).Error
}

// Exists checks if a wishlist item exists
func (r *wishlistRepository) Exists(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Wishlist{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	return count > 0, err
}

// GetPopularProducts gets most wishlisted products
func (r *wishlistRepository) GetPopularProducts(ctx context.Context, limit int) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.db.WithContext(ctx).
		Table("products").
		Select("products.*, COUNT(user_wishlists.product_id) as wishlist_count").
		Joins("JOIN user_wishlists ON products.id = user_wishlists.product_id").
		Group("products.id").
		Order("wishlist_count DESC").
		Limit(limit).
		Find(&products).Error
	return products, err
}

// Update updates a wishlist item
func (r *wishlistRepository) Update(ctx context.Context, wishlist *entities.Wishlist) error {
	wishlist.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(wishlist).Error
}

// List lists wishlist items with filters
func (r *wishlistRepository) List(ctx context.Context, filters repositories.WishlistFilters) ([]*entities.Wishlist, error) {
	var wishlists []*entities.Wishlist
	query := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Images").
		Preload("Product.Category")

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.ProductID != nil {
		query = query.Where("product_id = ?", *filters.ProductID)
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
	case "product_name":
		query = query.Joins("JOIN products ON user_wishlists.product_id = products.id")
		if filters.SortOrder == "desc" {
			query = query.Order("products.name DESC")
		} else {
			query = query.Order("products.name ASC")
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

	err := query.Find(&wishlists).Error
	return wishlists, err
}

// Count counts wishlist items with filters
func (r *wishlistRepository) Count(ctx context.Context, filters repositories.WishlistFilters) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Wishlist{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.ProductID != nil {
		query = query.Where("product_id = ?", *filters.ProductID)
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

// AddToWishlist adds a product to user's wishlist
func (r *wishlistRepository) AddToWishlist(ctx context.Context, userID, productID uuid.UUID) error {
	wishlist := &entities.Wishlist{
		ID:        uuid.New(),
		UserID:    userID,
		ProductID: productID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return r.Create(ctx, wishlist)
}

// ClearWishlist clears all items from user's wishlist (alias for ClearByUser)
func (r *wishlistRepository) ClearWishlist(ctx context.Context, userID uuid.UUID) error {
	return r.ClearByUser(ctx, userID)
}

// CountByUserID counts wishlist items for a user (alias for CountByUser)
func (r *wishlistRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.CountByUser(ctx, userID)
}

// GetByUserID gets wishlist items for a user (alias for GetByUser)
func (r *wishlistRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Wishlist, error) {
	return r.GetByUser(ctx, userID, limit, offset)
}

// GetWishlistProductIDs gets product IDs in user's wishlist
func (r *wishlistRepository) GetWishlistProductIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var productIDs []uuid.UUID
	err := r.db.WithContext(ctx).
		Model(&entities.Wishlist{}).
		Where("user_id = ?", userID).
		Pluck("product_id", &productIDs).Error
	return productIDs, err
}

// GetWishlistProducts gets products in user's wishlist with details
func (r *wishlistRepository) GetWishlistProducts(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.db.WithContext(ctx).
		Table("products").
		Select("products.*").
		Joins("JOIN user_wishlists ON products.id = user_wishlists.product_id").
		Where("user_wishlists.user_id = ?", userID).
		Preload("Images").
		Preload("Category").
		Order("user_wishlists.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, err
}

// IsInWishlist checks if a product is in user's wishlist
func (r *wishlistRepository) IsInWishlist(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	return r.Exists(ctx, userID, productID)
}

// RemoveFromWishlist removes a product from user's wishlist
func (r *wishlistRepository) RemoveFromWishlist(ctx context.Context, userID, productID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&entities.Wishlist{}, "user_id = ? AND product_id = ?", userID, productID).Error
}
