package database

import (
	"context"
	"fmt"
	"strings"
	"time"

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

// GetByIDForUpdate retrieves a product by ID with row-level locking (SELECT FOR UPDATE)
func (r *productRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).Session(&gorm.Session{}).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags").
		Where("id = ?", id).
		Set("gorm:query_option", "FOR UPDATE").
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

	// Apply filters with enhanced full-text search
	if params.Query != "" {
		// Use PostgreSQL full-text search for better performance and relevance
		searchVector := "to_tsvector('english', coalesce(name, '') || ' ' || coalesce(description, '') || ' ' || coalesce(short_description, '') || ' ' || coalesce(sku, '') || ' ' || coalesce(keywords, ''))"
		searchQuery := "plainto_tsquery('english', ?)"

		// Combine full-text search with ILIKE for partial matches
		query = query.Where(
			fmt.Sprintf("(%s @@ %s) OR name ILIKE ? OR description ILIKE ? OR sku ILIKE ?", searchVector, searchQuery),
			params.Query, "%"+params.Query+"%", "%"+params.Query+"%", "%"+params.Query+"%",
		)
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

	// Apply sorting with relevance ranking
	orderBy := r.buildSortOrder(params.SortBy, params.SortOrder, params.Query)
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

// UpdateStock updates product stock and stock status
func (r *productRepository) UpdateStock(ctx context.Context, productID uuid.UUID, stock int) error {
	// Get the product first to calculate stock status
	product, err := r.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	// Update stock and calculate new stock status
	product.Stock = stock
	product.UpdateStockStatus()

	// Update both stock and stock_status in database
	result := r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"stock":        stock,
			"stock_status": product.StockStatus,
			"updated_at":   time.Now(),
		})

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
	// Get category_id directly without loading full product
	var categoryID uuid.UUID
	err := r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Select("category_id").
		Where("id = ?", productID).
		Scan(&categoryID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrProductNotFound
		}
		return nil, err
	}

	var products []*entities.Product
	err = r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags").
		Where("category_id = ? AND id != ?", categoryID, productID).
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

// GetByIDsWithFullDetails retrieves multiple products by IDs with all relations (optimized for bulk operations)
func (r *productRepository) GetByIDsWithFullDetails(ctx context.Context, ids []uuid.UUID) ([]*entities.Product, error) {
	if len(ids) == 0 {
		return []*entities.Product{}, nil
	}

	var products []*entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Brand").
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

	// Apply filters with enhanced full-text search
	if params.Query != "" {
		// Use PostgreSQL full-text search for better performance and relevance
		searchVector := "to_tsvector('english', coalesce(name, '') || ' ' || coalesce(description, '') || ' ' || coalesce(short_description, '') || ' ' || coalesce(sku, '') || ' ' || coalesce(keywords, ''))"
		searchQuery := "plainto_tsquery('english', ?)"

		// Combine full-text search with ILIKE for partial matches and add relevance ranking
		query = query.Where(
			fmt.Sprintf("(%s @@ %s) OR name ILIKE ? OR description ILIKE ? OR short_description ILIKE ? OR sku ILIKE ?", searchVector, searchQuery),
			params.Query, "%"+params.Query+"%", "%"+params.Query+"%", "%"+params.Query+"%", "%"+params.Query+"%",
		)
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

// buildSortOrder builds the sort order clause with relevance ranking
func (r *productRepository) buildSortOrder(sortBy, sortOrder, searchQuery string) string {
	direction := "ASC"
	if strings.ToUpper(sortOrder) == "DESC" {
		direction = "DESC"
	}

	switch sortBy {
	case "relevance":
		if searchQuery != "" {
			// Use PostgreSQL's ts_rank for relevance scoring
			searchVector := "to_tsvector('english', coalesce(name, '') || ' ' || coalesce(description, '') || ' ' || coalesce(short_description, '') || ' ' || coalesce(sku, '') || ' ' || coalesce(keywords, ''))"
			return fmt.Sprintf("ts_rank(%s, plainto_tsquery('english', '%s')) %s, featured DESC, created_at DESC", searchVector, searchQuery, direction)
		}
		return "featured DESC, created_at DESC"
	case "price":
		return "price " + direction
	case "name":
		return "name " + direction
	case "created_at":
		return "created_at " + direction
	case "rating":
		// Assuming we have a rating field or calculate it from reviews
		return "rating " + direction + ", created_at DESC"
	case "popularity":
		// Assuming we track view counts or sales
		return "view_count " + direction + ", created_at DESC"
	default:
		return "created_at DESC"
	}
}

// GetSearchSuggestions returns search suggestions based on query
func (r *productRepository) GetSearchSuggestions(ctx context.Context, query string, limit int) (*repositories.SearchSuggestions, error) {
	suggestions := &repositories.SearchSuggestions{
		Products:    []repositories.ProductSuggestion{},
		Categories:  []repositories.CategorySuggestion{},
		Brands:      []repositories.BrandSuggestion{},
		Popular:     []string{},
		Corrections: []string{},
	}

	if limit <= 0 {
		limit = 10
	}

	// Get product suggestions
	var products []entities.Product
	err := r.db.WithContext(ctx).
		Where("name ILIKE ? OR description ILIKE ? OR sku ILIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Where("status = ?", "active").
		Limit(limit).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	for _, product := range products {
		// Get category name
		var category entities.Category
		r.db.WithContext(ctx).Where("id = ?", product.CategoryID).First(&category)

		// Get first image
		var image string
		if len(product.Images) > 0 {
			image = product.Images[0].URL
		}

		suggestion := repositories.ProductSuggestion{
			ID:         product.ID,
			Name:       product.Name,
			SKU:        product.SKU,
			Price:      product.Price,
			Image:      image,
			CategoryID: product.CategoryID,
			Category:   category.Name,
			Relevance:  calculateRelevance(product.Name, query),
		}
		suggestions.Products = append(suggestions.Products, suggestion)
	}

	// Get category suggestions
	var categories []entities.Category
	err = r.db.WithContext(ctx).
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Where("is_active = ?", true).
		Limit(limit/2).
		Find(&categories).Error
	if err == nil {
		for _, category := range categories {
			// Count products in category
			var productCount int64
			r.db.WithContext(ctx).Model(&entities.Product{}).
				Where("category_id = ? AND status = ?", category.ID, "active").
				Count(&productCount)

			suggestion := repositories.CategorySuggestion{
				ID:           category.ID,
				Name:         category.Name,
				ProductCount: productCount,
				Relevance:    calculateRelevance(category.Name, query),
			}
			suggestions.Categories = append(suggestions.Categories, suggestion)
		}
	}

	// Get popular searches
	popular, err := r.GetPopularSearches(ctx, 5)
	if err == nil {
		suggestions.Popular = popular
	}

	return suggestions, nil
}

// GetPopularSearches returns popular search queries
func (r *productRepository) GetPopularSearches(ctx context.Context, limit int) ([]string, error) {
	var results []struct {
		Query string
		Count int64
	}

	err := r.db.WithContext(ctx).
		Table("search_queries").
		Select("query, COUNT(*) as count").
		Where("created_at > ?", time.Now().AddDate(0, 0, -30)). // Last 30 days
		Group("query").
		Order("count DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		// Return default popular searches if no data
		return []string{"laptop", "phone", "headphones", "camera", "watch"}, nil
	}

	var queries []string
	for _, result := range results {
		queries = append(queries, result.Query)
	}

	return queries, nil
}

// RecordSearchQuery records a search query for analytics
func (r *productRepository) RecordSearchQuery(ctx context.Context, query string, userID *uuid.UUID, resultCount int) error {
	searchQuery := repositories.SearchQuery{
		ID:          uuid.New(),
		Query:       query,
		UserID:      userID,
		ResultCount: resultCount,
		CreatedAt:   time.Now(),
	}

	return r.db.WithContext(ctx).Table("search_queries").Create(&searchQuery).Error
}

// GetSearchHistory returns search history for a user
func (r *productRepository) GetSearchHistory(ctx context.Context, userID uuid.UUID, limit int) ([]string, error) {
	var queries []string
	err := r.db.WithContext(ctx).
		Table("search_history").
		Select("query").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Pluck("query", &queries).Error

	return queries, err
}

// calculateRelevance calculates relevance score between product name and search query
func calculateRelevance(productName, query string) float64 {
	productLower := strings.ToLower(productName)
	queryLower := strings.ToLower(query)

	// Exact match gets highest score
	if productLower == queryLower {
		return 1.0
	}

	// Starts with query gets high score
	if strings.HasPrefix(productLower, queryLower) {
		return 0.9
	}

	// Contains query gets medium score
	if strings.Contains(productLower, queryLower) {
		return 0.7
	}

	// Default relevance
	return 0.5
}
