package database

import (
	"context"
	"fmt"
	"strings"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) repositories.ProductRepository {
	return &productRepository{db: db}
}

// Create creates a new product
func (r *productRepository) Create(ctx context.Context, product *entities.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// GetByID retrieves a product by ID
func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).Session(&gorm.Session{}).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags").
		Where("id = ?", id).
		First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

// GetByIDs retrieves multiple products by IDs (bulk operation)
func (r *productRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*entities.Product, error) {
	if len(ids) == 0 {
		return []*entities.Product{}, nil
	}

	var products []*entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags").
		Where("id IN ?", ids).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// GetBySKU retrieves a product by SKU
func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags").
		Where("sku = ?", sku).
		First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

// Update updates an existing product
func (r *productRepository) Update(ctx context.Context, product *entities.Product) error {
	// Use Updates instead of Save to ensure all fields are updated properly
	// Select specific fields to avoid issues with relationships
	result := r.db.WithContext(ctx).Model(product).Select(
		// Basic fields
		"name", "description", "short_description", "sku", "updated_at",

		// SEO and Metadata
		"slug", "meta_title", "meta_description", "keywords", "featured", "visibility",

		// Pricing
		"price", "compare_price", "cost_price",

		// Sale Pricing
		"sale_price", "sale_start_date", "sale_end_date",

		// Inventory
		"stock", "low_stock_threshold", "track_quantity", "allow_backorder", "stock_status",

		// Physical Properties
		"weight", "length", "width", "height", // dimensions fields

		// Shipping and Tax
		"requires_shipping", "shipping_class", "tax_class", "country_of_origin",

		// Categorization
		"category_id", "brand_id",

		// Status and Type
		"status", "product_type", "is_digital",
	).Updates(product)

	return result.Error
}

// Delete deletes a product by ID
func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get the product first to check if it exists
	var product entities.Product
	err := tx.Where("id = ?", id).First(&product).Error
	if err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return entities.ErrProductNotFound
		}
		return err
	}

	// Remove all associations from product_tag_associations table
	err = tx.Exec("DELETE FROM product_tag_associations WHERE product_id = ?", id).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete product images (if any)
	err = tx.Where("product_id = ?", id).Delete(&entities.ProductImage{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Finally delete the product
	result := tx.Delete(&entities.Product{}, id)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}

// List retrieves products with pagination
func (r *productRepository) List(ctx context.Context, limit, offset int) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

// Search searches products based on criteria
func (r *productRepository) Search(ctx context.Context, params repositories.ProductSearchParams) ([]*entities.Product, error) {
	query := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags")

	// Apply filters
	if params.Query != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+params.Query+"%", "%"+params.Query+"%")
	}

	// Enhanced category filter with recursive search (includes subcategories)
	if params.CategoryID != nil {
		// Get all descendant categories using recursive CTE
		categoryQuery := `
			WITH RECURSIVE category_tree AS (
				-- Base case: start with the given category
				SELECT id FROM categories WHERE id = $1 AND is_active = true
				
				UNION ALL
				
				-- Recursive case: find all children
				SELECT c.id FROM categories c
				INNER JOIN category_tree ct ON c.parent_id = ct.id
				WHERE c.is_active = true
			)
			SELECT id FROM category_tree
		`

		var categoryIDs []uuid.UUID
		rows, err := r.db.WithContext(ctx).Raw(categoryQuery, *params.CategoryID).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var id uuid.UUID
			if err := rows.Scan(&id); err != nil {
				return nil, err
			}
			categoryIDs = append(categoryIDs, id)
		}

		if len(categoryIDs) > 0 {
			query = query.Where("category_id IN ?", categoryIDs)
		} else {
			// If no categories found, still filter by the original category
			query = query.Where("category_id = ?", *params.CategoryID)
		}
	}

	if params.MinPrice != nil {
		query = query.Where("price >= ?", *params.MinPrice)
	}

	if params.MaxPrice != nil {
		query = query.Where("price <= ?", *params.MaxPrice)
	}

	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	// Apply sorting
	orderBy := "created_at DESC"
	if params.SortBy != "" {
		direction := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			direction = "DESC"
		}
		orderBy = params.SortBy + " " + direction
	}
	query = query.Order(orderBy)

	// Apply pagination
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	var products []*entities.Product
	err := query.Find(&products).Error
	return products, err
}

// Count returns the total number of products
func (r *productRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Product{}).Count(&count).Error
	return count, err
}

// CountByCategory returns the number of products in a category
func (r *productRepository) CountByCategory(ctx context.Context, categoryID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("category_id = ?", categoryID).
		Count(&count).Error
	return count, err
}

// GetByCategory retrieves products by category
func (r *productRepository) GetByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags").
		Where("category_id = ?", categoryID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

// UpdateStock updates product stock
func (r *productRepository) UpdateStock(ctx context.Context, productID uuid.UUID, stock int) error {
	result := r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("id = ?", productID).
		Update("stock", stock)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrProductNotFound
	}
	return nil
}

// ExistsBySKU checks if a product exists with the given SKU
func (r *productRepository) ExistsBySKU(ctx context.Context, sku string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("sku = ?", sku).
		Count(&count).Error
	return count > 0, err
}

// ExistsBySlug checks if a product exists with the given slug
func (r *productRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("slug = ?", slug).
		Count(&count).Error
	return count > 0, err
}

// ExistsBySlugExcludingID checks if a product exists with the given slug, excluding a specific ID
func (r *productRepository) ExistsBySlugExcludingID(ctx context.Context, slug string, excludeID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("slug = ? AND id != ?", slug, excludeID).
		Count(&count).Error
	return count > 0, err
}

// GetExistingSlugs gets all existing slugs that start with the given prefix
func (r *productRepository) GetExistingSlugs(ctx context.Context, prefix string) ([]string, error) {
	var slugs []string
	err := r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("slug LIKE ?", prefix+"%").
		Pluck("slug", &slugs).Error
	return slugs, err
}

// GetFeatured retrieves featured products
func (r *productRepository) GetFeatured(ctx context.Context, limit int) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags").
		// Joins("JOIN product_tags ON products.id = product_tags.id"). // Temporarily disabled
		// Where("product_tags.slug = ?", "featured"). // Temporarily disabled
		Limit(limit).
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

// GetRelated retrieves related products
func (r *productRepository) GetRelated(ctx context.Context, productID uuid.UUID, limit int) ([]*entities.Product, error) {
	// Get the product to find its category
	product, err := r.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	var products []*entities.Product
	err = r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags").
		Where("category_id = ? AND id != ?", product.CategoryID, productID).
		Limit(limit).
		Order("RANDOM()").
		Find(&products).Error
	return products, err
}

// ClearTags removes all tag associations for a product using GORM Association
func (r *productRepository) ClearTags(ctx context.Context, productID uuid.UUID) error {
	// Get the product first
	var product entities.Product
	if err := r.db.WithContext(ctx).First(&product, productID).Error; err != nil {
		return err
	}

	// Use GORM Association to clear all tags
	if err := r.db.WithContext(ctx).Model(&product).Association("Tags").Clear(); err != nil {
		return fmt.Errorf("failed to clear tags: %w", err)
	}

	return nil
}

// AddTag adds a tag association to a product using GORM Association
func (r *productRepository) AddTag(ctx context.Context, productID, tagID uuid.UUID) error {
	// Get the product and tag
	var product entities.Product
	if err := r.db.WithContext(ctx).First(&product, productID).Error; err != nil {
		return err
	}

	var tag entities.ProductTag
	if err := r.db.WithContext(ctx).First(&tag, tagID).Error; err != nil {
		return err
	}

	// Use GORM Association to append tag
	if err := r.db.WithContext(ctx).Model(&product).Association("Tags").Append(&tag); err != nil {
		return fmt.Errorf("failed to add tag: %w", err)
	}

	return nil
}

// ReplaceTags replaces all tag associations for a product with new ones
func (r *productRepository) ReplaceTags(ctx context.Context, productID uuid.UUID, tagIDs []uuid.UUID) error {
	if len(tagIDs) == 0 {
		// If no tags provided, just clear all
		return r.ClearTags(ctx, productID)
	}

	// Get the product
	var product entities.Product
	if err := r.db.WithContext(ctx).First(&product, productID).Error; err != nil {
		return err
	}

	// Get all tags
	var tags []entities.ProductTag
	if err := r.db.WithContext(ctx).Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return err
	}

	// Validate that all tag IDs exist
	if len(tags) != len(tagIDs) {
		return fmt.Errorf("some tag IDs do not exist")
	}

	// Use GORM Association to replace all tags
	if err := r.db.WithContext(ctx).Model(&product).Association("Tags").Replace(tags); err != nil {
		return fmt.Errorf("failed to replace tags: %w", err)
	}

	return nil
}

// CountProducts counts total number of products
func (r *productRepository) CountProducts(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Count(&count).Error
	return count, err
}

// GetByBrand retrieves products by brand
func (r *productRepository) GetByBrand(ctx context.Context, brandID uuid.UUID, limit, offset int) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags").
		Where("brand_id = ?", brandID).
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, err
}

// GetBySlug retrieves a product by slug
func (r *productRepository) GetBySlug(ctx context.Context, slug string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags").
		Where("slug = ?", slug).
		First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

// SearchAdvanced performs advanced search with multiple filters
func (r *productRepository) SearchAdvanced(ctx context.Context, params repositories.AdvancedSearchParams) ([]*entities.Product, error) {
	query := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags")

	// Apply filters
	if params.Query != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ? OR short_description ILIKE ?",
			"%"+params.Query+"%", "%"+params.Query+"%", "%"+params.Query+"%")
	}

	if params.CategoryID != nil {
		query = query.Where("category_id = ?", *params.CategoryID)
	}

	if params.BrandID != nil {
		query = query.Where("brand_id = ?", *params.BrandID)
	}

	if params.MinPrice != nil {
		query = query.Where("price >= ?", *params.MinPrice)
	}

	if params.MaxPrice != nil {
		query = query.Where("price <= ?", *params.MaxPrice)
	}

	if params.InStock != nil && *params.InStock {
		query = query.Where("stock > 0")
	}

	if params.Featured != nil {
		query = query.Where("featured = ?", *params.Featured)
	}

	if params.Visibility != nil {
		query = query.Where("visibility = ?", *params.Visibility)
	}

	if params.ProductType != nil {
		query = query.Where("product_type = ?", *params.ProductType)
	}

	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	// Apply sorting
	if params.SortBy != "" {
		order := params.SortBy
		if params.SortOrder == "desc" {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	var products []*entities.Product
	err := query.Find(&products).Error
	return products, err
}
