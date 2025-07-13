package database

import (
	"context"
	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productComparisonRepository struct {
	db *gorm.DB
}

// NewProductComparisonRepository creates a new product comparison repository
func NewProductComparisonRepository(db *gorm.DB) repositories.ProductComparisonRepository {
	return &productComparisonRepository{db: db}
}

// CreateComparison creates a new product comparison
func (r *productComparisonRepository) CreateComparison(ctx context.Context, comparison *entities.ProductComparison) error {
	return r.db.WithContext(ctx).Create(comparison).Error
}

// GetComparison gets a comparison by ID
func (r *productComparisonRepository) GetComparison(ctx context.Context, id uuid.UUID) (*entities.ProductComparison, error) {
	var comparison entities.ProductComparison
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Brand").
		Preload("Items.Product.Images").
		First(&comparison, id).Error
	if err != nil {
		return nil, err
	}
	return &comparison, nil
}

// GetComparisonByUserID gets a comparison by user ID
func (r *productComparisonRepository) GetComparisonByUserID(ctx context.Context, userID uuid.UUID) (*entities.ProductComparison, error) {
	var comparison entities.ProductComparison
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Brand").
		Preload("Items.Product.Images").
		Where("user_id = ?", userID).
		First(&comparison).Error
	if err != nil {
		return nil, err
	}
	return &comparison, nil
}

// GetComparisonBySessionID gets a comparison by session ID
func (r *productComparisonRepository) GetComparisonBySessionID(ctx context.Context, sessionID string) (*entities.ProductComparison, error) {
	var comparison entities.ProductComparison
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Brand").
		Preload("Items.Product.Images").
		Where("session_id = ?", sessionID).
		First(&comparison).Error
	if err != nil {
		return nil, err
	}
	return &comparison, nil
}

// UpdateComparison updates a comparison
func (r *productComparisonRepository) UpdateComparison(ctx context.Context, comparison *entities.ProductComparison) error {
	return r.db.WithContext(ctx).Save(comparison).Error
}

// DeleteComparison deletes a comparison
func (r *productComparisonRepository) DeleteComparison(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete comparison items first
		if err := tx.Where("comparison_id = ?", id).Delete(&entities.ProductComparisonItem{}).Error; err != nil {
			return err
		}
		// Delete comparison
		return tx.Delete(&entities.ProductComparison{}, id).Error
	})
}

// AddProductToComparison adds a product to comparison
func (r *productComparisonRepository) AddProductToComparison(ctx context.Context, comparisonID, productID uuid.UUID, position int) error {
	// Check if product already exists in comparison
	var count int64
	r.db.WithContext(ctx).Model(&entities.ProductComparisonItem{}).
		Where("comparison_id = ? AND product_id = ?", comparisonID, productID).
		Count(&count)
	
	if count > 0 {
		return fmt.Errorf("product already exists in comparison")
	}

	item := &entities.ProductComparisonItem{
		ComparisonID: comparisonID,
		ProductID:    productID,
		Position:     position,
	}
	return r.db.WithContext(ctx).Create(item).Error
}

// RemoveProductFromComparison removes a product from comparison
func (r *productComparisonRepository) RemoveProductFromComparison(ctx context.Context, comparisonID, productID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("comparison_id = ? AND product_id = ?", comparisonID, productID).
		Delete(&entities.ProductComparisonItem{}).Error
}

// GetComparisonItems gets all items in a comparison
func (r *productComparisonRepository) GetComparisonItems(ctx context.Context, comparisonID uuid.UUID) ([]entities.ProductComparisonItem, error) {
	var items []entities.ProductComparisonItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Category").
		Preload("Product.Brand").
		Preload("Product.Images").
		Where("comparison_id = ?", comparisonID).
		Order("position ASC").
		Find(&items).Error
	return items, err
}

// UpdateItemPosition updates the position of an item in comparison
func (r *productComparisonRepository) UpdateItemPosition(ctx context.Context, itemID uuid.UUID, position int) error {
	return r.db.WithContext(ctx).
		Model(&entities.ProductComparisonItem{}).
		Where("id = ?", itemID).
		Update("position", position).Error
}

// ClearComparison removes all items from a comparison
func (r *productComparisonRepository) ClearComparison(ctx context.Context, comparisonID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("comparison_id = ?", comparisonID).
		Delete(&entities.ProductComparisonItem{}).Error
}

// GetComparisonWithProducts gets a comparison with all products
func (r *productComparisonRepository) GetComparisonWithProducts(ctx context.Context, id uuid.UUID) (*entities.ProductComparison, error) {
	return r.GetComparison(ctx, id)
}

// GetUserComparisons gets all comparisons for a user
func (r *productComparisonRepository) GetUserComparisons(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entities.ProductComparison, error) {
	var comparisons []entities.ProductComparison
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comparisons).Error
	return comparisons, err
}

// CountComparisonItems counts items in a comparison
func (r *productComparisonRepository) CountComparisonItems(ctx context.Context, comparisonID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.ProductComparisonItem{}).
		Where("comparison_id = ?", comparisonID).
		Count(&count).Error
	return count, err
}

// IsProductInComparison checks if a product is in comparison
func (r *productComparisonRepository) IsProductInComparison(ctx context.Context, comparisonID, productID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.ProductComparisonItem{}).
		Where("comparison_id = ? AND product_id = ?", comparisonID, productID).
		Count(&count).Error
	return count > 0, err
}

// GetPopularComparedProducts gets most compared products
func (r *productComparisonRepository) GetPopularComparedProducts(ctx context.Context, limit int) ([]entities.Product, error) {
	var products []entities.Product
	err := r.db.WithContext(ctx).
		Table("products").
		Select("products.*, COUNT(pci.product_id) as comparison_count").
		Joins("JOIN product_comparison_items pci ON products.id = pci.product_id").
		Group("products.id").
		Order("comparison_count DESC").
		Limit(limit).
		Find(&products).Error
	return products, err
}

// GetComparisonStats gets comparison statistics
func (r *productComparisonRepository) GetComparisonStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Total comparisons
	var totalComparisons int64
	r.db.WithContext(ctx).Model(&entities.ProductComparison{}).Count(&totalComparisons)
	stats["total_comparisons"] = totalComparisons
	
	// Total comparison items
	var totalItems int64
	r.db.WithContext(ctx).Model(&entities.ProductComparisonItem{}).Count(&totalItems)
	stats["total_items"] = totalItems
	
	// Average items per comparison
	if totalComparisons > 0 {
		stats["avg_items_per_comparison"] = float64(totalItems) / float64(totalComparisons)
	} else {
		stats["avg_items_per_comparison"] = 0
	}
	
	return stats, nil
}
