package database

import (
	"context"
	"strings"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type brandRepository struct {
	db *gorm.DB
}

// NewBrandRepository creates a new brand repository
func NewBrandRepository(db *gorm.DB) repositories.BrandRepository {
	return &brandRepository{db: db}
}

// Create creates a new brand
func (r *brandRepository) Create(ctx context.Context, brand *entities.Brand) error {
	return r.db.WithContext(ctx).Create(brand).Error
}

// GetByID retrieves a brand by ID
func (r *brandRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Brand, error) {
	var brand entities.Brand
	err := r.db.WithContext(ctx).
		Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "published").Limit(10)
		}).
		Where("id = ?", id).
		First(&brand).Error
	if err != nil {
		return nil, err
	}
	return &brand, nil
}

// GetBySlug retrieves a brand by slug
func (r *brandRepository) GetBySlug(ctx context.Context, slug string) (*entities.Brand, error) {
	var brand entities.Brand
	err := r.db.WithContext(ctx).
		Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "published").Limit(10)
		}).
		Where("slug = ?", slug).
		First(&brand).Error
	if err != nil {
		return nil, err
	}
	return &brand, nil
}

// Update updates an existing brand
func (r *brandRepository) Update(ctx context.Context, brand *entities.Brand) error {
	return r.db.WithContext(ctx).Save(brand).Error
}

// Delete deletes a brand by ID
func (r *brandRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Brand{}, id).Error
}

// List retrieves brands with pagination
func (r *brandRepository) List(ctx context.Context, limit, offset int) ([]*entities.Brand, error) {
	var brands []*entities.Brand
	err := r.db.WithContext(ctx).
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Find(&brands).Error
	return brands, err
}

// Search searches brands
func (r *brandRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entities.Brand, error) {
	var brands []*entities.Brand
	searchQuery := "%" + strings.ToLower(query) + "%"

	err := r.db.WithContext(ctx).
		Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchQuery, searchQuery).
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Find(&brands).Error
	return brands, err
}

// CountSearch counts brands matching search query
func (r *brandRepository) CountSearch(ctx context.Context, query string) (int64, error) {
	var count int64
	searchQuery := "%" + strings.ToLower(query) + "%"

	err := r.db.WithContext(ctx).
		Model(&entities.Brand{}).
		Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchQuery, searchQuery).
		Count(&count).Error
	return count, err
}

// ExistsBySlug checks if a brand exists with the given slug
func (r *brandRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Brand{}).
		Where("slug = ?", slug).
		Count(&count).Error
	return count > 0, err
}

// GetActive retrieves active brands
func (r *brandRepository) GetActive(ctx context.Context, limit, offset int) ([]*entities.Brand, error) {
	var brands []*entities.Brand
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Find(&brands).Error
	return brands, err
}

// BrandWithCount represents a brand with product count
type BrandWithCount struct {
	entities.Brand
	ProductCount int `json:"product_count"`
}

// GetBrandWithProductCount retrieves brands with product count
func (r *brandRepository) GetBrandWithProductCount(ctx context.Context, limit, offset int) ([]*entities.Brand, error) {
	var results []struct {
		entities.Brand
		ProductCount int64 `gorm:"column:product_count"`
	}

	// Use a single query with LEFT JOIN to get brands and their product counts
	err := r.db.WithContext(ctx).Debug().
		Table("brands").
		Select("brands.*, COALESCE(COUNT(products.id), 0) as product_count").
		Joins("LEFT JOIN products ON brands.id = products.brand_id AND products.status = 'active'").
		Group("brands.id").
		Order("brands.name ASC").
		Limit(limit).
		Offset(offset).
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	// Convert results to brand entities
	brands := make([]*entities.Brand, len(results))
	for i, result := range results {
		brands[i] = &result.Brand
		brands[i].ProductCount = int(result.ProductCount)
	}

	return brands, nil
}

// GetPopularBrands retrieves brands ordered by product count
func (r *brandRepository) GetPopularBrands(ctx context.Context, limit int) ([]*entities.Brand, error) {
	var brandsWithCount []BrandWithCount
	err := r.db.WithContext(ctx).
		Select("brands.*, COUNT(products.id) as product_count").
		Joins("LEFT JOIN products ON brands.id = products.brand_id AND products.status = 'active'").
		Where("brands.is_active = ?", true).
		Group("brands.id").
		Having("COUNT(products.id) > 0").
		Order("COUNT(products.id) DESC, brands.name ASC").
		Limit(limit).
		Find(&brandsWithCount).Error

	if err != nil {
		return nil, err
	}

	// Convert to Brand entities with ProductCount set
	brands := make([]*entities.Brand, len(brandsWithCount))
	for i, bwc := range brandsWithCount {
		brands[i] = &bwc.Brand
		brands[i].ProductCount = bwc.ProductCount
	}

	return brands, nil
}

// GetBrandsForFiltering retrieves brands for product filtering with counts
func (r *brandRepository) GetBrandsForFiltering(ctx context.Context, categoryID *uuid.UUID) ([]map[string]interface{}, error) {
	query := r.db.WithContext(ctx).
		Select("brands.id, brands.name, COUNT(products.id) as count").
		Table("brands").
		Joins("INNER JOIN products ON brands.id = products.brand_id").
		Where("brands.is_active = ? AND products.status = ?", true, "active").
		Group("brands.id, brands.name").
		Having("COUNT(products.id) > 0").
		Order("brands.name ASC")

	if categoryID != nil {
		query = query.Where("products.category_id = ?", *categoryID)
	}

	var results []map[string]interface{}
	err := query.Find(&results).Error
	return results, err
}

// CountByStatus counts brands by status
func (r *brandRepository) CountByStatus(ctx context.Context, isActive bool) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Brand{}).
		Where("is_active = ?", isActive).
		Count(&count).Error
	return count, err
}

// GetTotal gets total number of brands
func (r *brandRepository) GetTotal(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Brand{}).
		Count(&count).Error
	return count, err
}
