package database

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) repositories.CategoryRepository {
	return &categoryRepository{db: db}
}

// Create creates a new category
func (r *categoryRepository) Create(ctx context.Context, category *entities.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// GetByID retrieves a category by ID
func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Category, error) {
	var category entities.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("id = ?", id).
		First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

// GetBySlug retrieves a category by slug
func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*entities.Category, error) {
	var category entities.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("slug = ?", slug).
		First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

// Update updates an existing category
func (r *categoryRepository) Update(ctx context.Context, category *entities.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// Delete deletes a category by ID
func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrCategoryNotFound
	}
	return nil
}

// List retrieves categories with pagination
func (r *categoryRepository) List(ctx context.Context, limit, offset int) ([]*entities.Category, error) {
	var categories []*entities.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Limit(limit).
		Offset(offset).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// GetRootCategories retrieves root categories
func (r *categoryRepository) GetRootCategories(ctx context.Context) ([]*entities.Category, error) {
	var categories []*entities.Category
	err := r.db.WithContext(ctx).
		Preload("Children").
		Where("parent_id IS NULL AND is_active = ?", true).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// GetChildren retrieves child categories
func (r *categoryRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entities.Category, error) {
	var categories []*entities.Category
	err := r.db.WithContext(ctx).
		Where("parent_id = ? AND is_active = ?", parentID, true).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// Count returns the total number of categories
func (r *categoryRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Category{}).Count(&count).Error
	return count, err
}

// ExistsBySlug checks if a category exists with the given slug
func (r *categoryRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Category{}).
		Where("slug = ?", slug).
		Count(&count).Error
	return count > 0, err
}

// GetTree retrieves the category tree
func (r *categoryRepository) GetTree(ctx context.Context) ([]*entities.Category, error) {
	var categories []*entities.Category
	err := r.db.WithContext(ctx).
		Preload("Children").
		Where("parent_id IS NULL AND is_active = ?", true).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// GetCategoryTree returns all descendant category IDs for a given category (including itself)
func (r *categoryRepository) GetCategoryTree(ctx context.Context, categoryID uuid.UUID) ([]uuid.UUID, error) {
	var categoryIDs []uuid.UUID

	// Using recursive CTE query to get all descendants
	query := `
		WITH RECURSIVE category_tree AS (
			-- Base case: start with the given category
			SELECT id, parent_id, name
			FROM categories 
			WHERE id = $1 AND is_active = true

			UNION ALL

			-- Recursive case: find all children
			SELECT c.id, c.parent_id, c.name
			FROM categories c
			INNER JOIN category_tree ct ON c.parent_id = ct.id
			WHERE c.is_active = true
		)
		SELECT id FROM category_tree
	`

	rows, err := r.db.WithContext(ctx).Raw(query, categoryID).Rows()
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

	return categoryIDs, nil
}

// GetCategoryPath returns the full path from root to the given category
func (r *categoryRepository) GetCategoryPath(ctx context.Context, categoryID uuid.UUID) ([]*entities.Category, error) {
	var categories []*entities.Category

	// Using recursive CTE query to get path from root to category
	query := `
		WITH RECURSIVE category_path AS (
			-- Start with the target category
			SELECT id, parent_id, name, slug, sort_order, 0 as level
			FROM categories 
			WHERE id = $1 AND is_active = true

			UNION ALL

			-- Recursively find parent categories
			SELECT c.id, c.parent_id, c.name, c.slug, c.sort_order, cp.level + 1 as level
			FROM categories c
			INNER JOIN category_path cp ON c.id = cp.parent_id
			WHERE c.is_active = true
		)
		SELECT id, parent_id, name, slug, sort_order, level FROM category_path
		ORDER BY level DESC
	`

	rows, err := r.db.WithContext(ctx).Raw(query, categoryID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category entities.Category
		var level int
		if err := rows.Scan(&category.ID, &category.ParentID, &category.Name,
			&category.Slug, &category.SortOrder, &level); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

// GetProductCountByCategory returns product count for each category (including descendants)
func (r *categoryRepository) GetProductCountByCategory(ctx context.Context, categoryID uuid.UUID) (int64, error) {
	// Get all descendant categories
	categoryIDs, err := r.GetCategoryTree(ctx, categoryID)
	if err != nil {
		return 0, err
	}

	if len(categoryIDs) == 0 {
		return 0, nil
	}

	var count int64
	err = r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("category_id IN ? AND status = ?", categoryIDs, "active").
		Count(&count).Error

	return count, err
}

// GetWithProductsOptimized retrieves a category with its products and all relations (optimized)
func (r *categoryRepository) GetWithProductsOptimized(ctx context.Context, id uuid.UUID, limit, offset int) (*entities.Category, []*entities.Product, error) {
	// Get category first
	category, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	// Get products with all relations in one query
	var products []*entities.Product
	err = r.db.WithContext(ctx).
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags").
		Where("category_id = ? AND status = ?", id, "active").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&products).Error

	return category, products, err
}

// GetCategoriesWithProductCount retrieves categories with product count (optimized)
func (r *categoryRepository) GetCategoriesWithProductCount(ctx context.Context) ([]*entities.Category, map[uuid.UUID]int64, error) {
	// Get all categories using List method
	categories, err := r.List(ctx, 1000, 0) // Get up to 1000 categories
	if err != nil {
		return nil, nil, err
	}

	// Get product counts for all categories in one query
	type CategoryCount struct {
		CategoryID uuid.UUID `json:"category_id"`
		Count      int64     `json:"count"`
	}

	var counts []CategoryCount
	err = r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Select("category_id, COUNT(*) as count").
		Where("status = ?", "active").
		Group("category_id").
		Scan(&counts).Error
	if err != nil {
		return nil, nil, err
	}

	// Build count map
	countMap := make(map[uuid.UUID]int64)
	for _, count := range counts {
		countMap[count.CategoryID] = count.Count
	}

	return categories, countMap, nil
}

// BulkCreate creates multiple categories
func (r *categoryRepository) BulkCreate(ctx context.Context, categories []*entities.Category) error {
	if len(categories) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(categories).Error
}

// BulkUpdate updates multiple categories
func (r *categoryRepository) BulkUpdate(ctx context.Context, categories []*entities.Category) error {
	if len(categories) == 0 {
		return nil
	}

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, category := range categories {
		if err := tx.Save(category).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// BulkDelete deletes multiple categories
func (r *categoryRepository) BulkDelete(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&entities.Category{}).Error
}

// ListWithFilters retrieves categories with advanced filtering
func (r *categoryRepository) ListWithFilters(ctx context.Context, filters repositories.CategoryFilters) ([]*entities.Category, error) {
	var categories []*entities.Category

	query := r.db.WithContext(ctx).Model(&entities.Category{})

	// Apply filters
	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	if filters.ParentID != nil {
		query = query.Where("parent_id = ?", *filters.ParentID)
	}

	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}

	if filters.HasParent != nil {
		if *filters.HasParent {
			query = query.Where("parent_id IS NOT NULL")
		} else {
			query = query.Where("parent_id IS NULL")
		}
	}

	// Apply sorting
	sortBy := "name"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}

	sortOrder := "ASC"
	if filters.SortOrder == "desc" {
		sortOrder = "DESC"
	}

	query = query.Order(sortBy + " " + sortOrder)

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset >= 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&categories).Error
	return categories, err
}

// CountWithFilters counts categories with advanced filtering
func (r *categoryRepository) CountWithFilters(ctx context.Context, filters repositories.CategoryFilters) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&entities.Category{})

	// Apply filters (same as ListWithFilters)
	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	if filters.ParentID != nil {
		query = query.Where("parent_id = ?", *filters.ParentID)
	}

	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}

	if filters.HasParent != nil {
		if *filters.HasParent {
			query = query.Where("parent_id IS NOT NULL")
		} else {
			query = query.Where("parent_id IS NULL")
		}
	}

	err := query.Count(&count).Error
	return count, err
}

// Search searches categories by name and description
func (r *categoryRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entities.Category, error) {
	var categories []*entities.Category

	err := r.db.WithContext(ctx).
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Where("is_active = ?", true).
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Find(&categories).Error

	return categories, err
}

// ValidateHierarchy validates that setting parentID for categoryID won't create circular reference
func (r *categoryRepository) ValidateHierarchy(ctx context.Context, categoryID, parentID uuid.UUID) error {
	if categoryID == parentID {
		return entities.ErrInvalidInput
	}

	// Check if parentID is a descendant of categoryID
	descendants, err := r.GetCategoryTree(ctx, categoryID)
	if err != nil {
		return err
	}

	for _, descendantID := range descendants {
		if descendantID == parentID {
			return entities.ErrInvalidInput // Would create circular reference
		}
	}

	return nil
}



// GetProductCount returns product count for a category (with option to include subcategories)
func (r *categoryRepository) GetProductCount(ctx context.Context, categoryID uuid.UUID, includeSubcategories bool) (int64, error) {
	var count int64

	if includeSubcategories {
		// Get all descendant category IDs
		categoryIDs, err := r.GetCategoryTree(ctx, categoryID)
		if err != nil {
			return 0, err
		}

		err = r.db.WithContext(ctx).
			Model(&entities.Product{}).
			Where("category_id IN ? AND status = ?", categoryIDs, "active").
			Count(&count).Error
		return count, err
	} else {
		// Count only direct products
		err := r.db.WithContext(ctx).
			Model(&entities.Product{}).
			Where("category_id = ? AND status = ?", categoryID, "active").
			Count(&count).Error
		return count, err
	}
}

// MoveCategory moves a category to a new parent
func (r *categoryRepository) MoveCategory(ctx context.Context, categoryID, newParentID uuid.UUID) error {
	// Validate hierarchy to prevent circular references
	if err := r.ValidateHierarchy(ctx, categoryID, newParentID); err != nil {
		return err
	}

	// Update the category's parent
	err := r.db.WithContext(ctx).
		Model(&entities.Category{}).
		Where("id = ?", categoryID).
		Update("parent_id", newParentID).Error

	if err != nil {
		return err
	}

	// Rebuild paths for the moved category and its descendants
	return r.RebuildCategoryPaths(ctx)
}

// ReorderCategories reorders multiple categories
func (r *categoryRepository) ReorderCategories(ctx context.Context, reorderRequests []repositories.CategoryReorderRequest) error {
	if len(reorderRequests) == 0 {
		return nil
	}

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, req := range reorderRequests {
		err := tx.Model(&entities.Category{}).
			Where("id = ?", req.CategoryID).
			Update("sort_order", req.SortOrder).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetCategoryDepth returns the depth of a category in the tree
func (r *categoryRepository) GetCategoryDepth(ctx context.Context, categoryID uuid.UUID) (int, error) {
	var depth int

	query := `
		WITH RECURSIVE category_depth AS (
			-- Base case: start with the given category
			SELECT id, parent_id, 0 as depth
			FROM categories
			WHERE id = $1

			UNION ALL

			-- Recursive case: find parent and increment depth
			SELECT c.id, c.parent_id, cd.depth + 1
			FROM categories c
			INNER JOIN category_depth cd ON c.id = cd.parent_id
		)
		SELECT MAX(depth) FROM category_depth
	`

	err := r.db.WithContext(ctx).Raw(query, categoryID).Scan(&depth).Error
	return depth, err
}

// GetMaxDepth returns the maximum depth in the category tree
func (r *categoryRepository) GetMaxDepth(ctx context.Context) (int, error) {
	var maxDepth int

	query := `
		WITH RECURSIVE category_tree AS (
			-- Base case: root categories
			SELECT id, parent_id, 0 as depth
			FROM categories
			WHERE parent_id IS NULL

			UNION ALL

			-- Recursive case: children
			SELECT c.id, c.parent_id, ct.depth + 1
			FROM categories c
			INNER JOIN category_tree ct ON c.parent_id = ct.id
		)
		SELECT COALESCE(MAX(depth), 0) FROM category_tree
	`

	err := r.db.WithContext(ctx).Raw(query).Scan(&maxDepth).Error
	return maxDepth, err
}

// ValidateTreeIntegrity validates the entire category tree for consistency
func (r *categoryRepository) ValidateTreeIntegrity(ctx context.Context) error {
	// Check for circular references
	query := `
		WITH RECURSIVE category_check AS (
			SELECT id, parent_id, ARRAY[id] as path
			FROM categories
			WHERE parent_id IS NOT NULL

			UNION ALL

			SELECT c.id, c.parent_id, cc.path || c.id
			FROM categories c
			INNER JOIN category_check cc ON c.parent_id = cc.id
			WHERE NOT (c.id = ANY(cc.path))
		)
		SELECT COUNT(*) FROM category_check WHERE array_length(path, 1) > 10
	`

	var circularCount int
	err := r.db.WithContext(ctx).Raw(query).Scan(&circularCount).Error
	if err != nil {
		return err
	}

	if circularCount > 0 {
		return entities.ErrCircularReference
	}

	return nil
}

// RebuildCategoryPaths rebuilds the path field for all categories
func (r *categoryRepository) RebuildCategoryPaths(ctx context.Context) error {
	// Update paths using recursive query
	query := `
		WITH RECURSIVE category_paths AS (
			-- Base case: root categories
			SELECT id, parent_id, name, name as path, 0 as level
			FROM categories
			WHERE parent_id IS NULL

			UNION ALL

			-- Recursive case: children
			SELECT c.id, c.parent_id, c.name, cp.path || ' > ' || c.name as path, cp.level + 1
			FROM categories c
			INNER JOIN category_paths cp ON c.parent_id = cp.id
		)
		UPDATE categories
		SET path = category_paths.path, level = category_paths.level
		FROM category_paths
		WHERE categories.id = category_paths.id
	`

	return r.db.WithContext(ctx).Exec(query).Error
}

// GetCategoryAnalytics returns comprehensive analytics for a category
func (r *categoryRepository) GetCategoryAnalytics(ctx context.Context, categoryID uuid.UUID, timeRange string) (*repositories.CategoryAnalytics, error) {
	analytics := &repositories.CategoryAnalytics{
		CategoryID: categoryID,
	}

	// Get category name
	var category entities.Category
	err := r.db.WithContext(ctx).Where("id = ?", categoryID).First(&category).Error
	if err != nil {
		return nil, err
	}
	analytics.CategoryName = category.Name

	// Get product counts including subcategories
	categoryIDs, err := r.GetCategoryTree(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	var productCount, activeProducts, inactiveProducts int64
	if len(categoryIDs) > 0 {
		r.db.WithContext(ctx).Model(&entities.Product{}).Where("category_id IN ?", categoryIDs).Count(&productCount)
		r.db.WithContext(ctx).Model(&entities.Product{}).Where("category_id IN ? AND status = ?", categoryIDs, "active").Count(&activeProducts)
		inactiveProducts = productCount - activeProducts
	}

	analytics.ProductCount = productCount
	analytics.ActiveProducts = activeProducts
	analytics.InactiveProducts = inactiveProducts

	// Calculate average price from all products in category hierarchy
	var avgPrice float64
	if len(categoryIDs) > 0 {
		r.db.WithContext(ctx).Model(&entities.Product{}).
			Where("category_id IN ? AND status = ?", categoryIDs, "active").
			Select("AVG(price)").Scan(&avgPrice)
	}
	analytics.AveragePrice = avgPrice

	// Get sales data (mock data for now - would integrate with actual order system)
	analytics.TotalSales = productCount * 10 // Mock calculation
	analytics.Revenue = avgPrice * float64(analytics.TotalSales)
	analytics.ConversionRate = 0.05 // Mock 5% conversion rate

	return analytics, nil
}

// GetTopCategories returns top performing categories with hierarchy support
func (r *categoryRepository) GetTopCategories(ctx context.Context, limit int, sortBy string) ([]*repositories.CategoryStats, error) {
	var categories []*entities.Category
	err := r.db.WithContext(ctx).Limit(limit).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	var stats []*repositories.CategoryStats
	for _, category := range categories {
		// Get product count including subcategories
		productCount, err := r.GetProductCount(ctx, category.ID, true)
		if err != nil {
			continue // Skip this category if error
		}

		// Calculate average price from all products in category hierarchy
		categoryIDs, err := r.GetCategoryTree(ctx, category.ID)
		if err != nil {
			continue // Skip this category if error
		}

		var avgPrice float64
		if len(categoryIDs) > 0 {
			r.db.WithContext(ctx).Model(&entities.Product{}).
				Where("category_id IN ? AND status = ?", categoryIDs, "active").
				Select("AVG(price)").Scan(&avgPrice)
		}

		// Mock sales and revenue data
		totalSales := productCount * 8
		revenue := avgPrice * float64(totalSales)

		stat := &repositories.CategoryStats{
			CategoryID:     category.ID,
			CategoryName:   category.Name,
			ProductCount:   productCount,
			TotalSales:     totalSales,
			Revenue:        revenue,
			AverageRating:  4.2, // Mock rating
			ConversionRate: 0.04, // Mock conversion rate
			GrowthRate:     0.15, // Mock 15% growth
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// GetCategoryPerformanceMetrics returns detailed performance metrics for a category
func (r *categoryRepository) GetCategoryPerformanceMetrics(ctx context.Context, categoryID uuid.UUID) (*repositories.CategoryPerformanceMetrics, error) {
	// Get category
	var category entities.Category
	err := r.db.WithContext(ctx).Where("id = ?", categoryID).First(&category).Error
	if err != nil {
		return nil, err
	}

	metrics := &repositories.CategoryPerformanceMetrics{
		CategoryID:   categoryID,
		CategoryName: category.Name,
	}

	// Get product counts including subcategories
	categoryIDs, err := r.GetCategoryTree(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	var productCount, activeProductCount int64
	if len(categoryIDs) > 0 {
		r.db.WithContext(ctx).Model(&entities.Product{}).Where("category_id IN ?", categoryIDs).Count(&productCount)
		r.db.WithContext(ctx).Model(&entities.Product{}).Where("category_id IN ? AND status = ?", categoryIDs, "active").Count(&activeProductCount)
	}

	metrics.ProductCount = productCount
	metrics.ActiveProductCount = activeProductCount

	// Calculate average price and inventory value from all products in category hierarchy
	var avgPrice, totalValue float64
	if len(categoryIDs) > 0 {
		r.db.WithContext(ctx).Model(&entities.Product{}).
			Where("category_id IN ? AND status = ?", categoryIDs, "active").
			Select("AVG(price)").Scan(&avgPrice)

		r.db.WithContext(ctx).Model(&entities.Product{}).
			Where("category_id IN ? AND status = ?", categoryIDs, "active").
			Select("SUM(price * stock)").Scan(&totalValue)
	}

	metrics.AverageProductPrice = avgPrice
	metrics.TotalInventoryValue = totalValue

	// Mock stock data (would integrate with inventory system)
	metrics.LowStockProducts = productCount / 10    // 10% low stock
	metrics.OutOfStockProducts = productCount / 20  // 5% out of stock

	// Mock review data (would integrate with review system)
	metrics.AverageRating = 4.3
	metrics.TotalReviews = productCount * 5
	metrics.PopularityScore = float64(activeProductCount) * metrics.AverageRating

	return metrics, nil
}

// GetCategorySalesStats returns sales statistics for a category
func (r *categoryRepository) GetCategorySalesStats(ctx context.Context, categoryID uuid.UUID, timeRange string) (*repositories.CategorySalesStats, error) {
	// Get category
	var category entities.Category
	err := r.db.WithContext(ctx).Where("id = ?", categoryID).First(&category).Error
	if err != nil {
		return nil, err
	}

	stats := &repositories.CategorySalesStats{
		CategoryID:   categoryID,
		CategoryName: category.Name,
		TimeRange:    timeRange,
	}

	// Get product count for calculations including subcategories
	categoryIDs, err := r.GetCategoryTree(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	var productCount int64
	if len(categoryIDs) > 0 {
		r.db.WithContext(ctx).Model(&entities.Product{}).Where("category_id IN ?", categoryIDs).Count(&productCount)
	}

	// Calculate average price from all products in category hierarchy
	var avgPrice float64
	if len(categoryIDs) > 0 {
		r.db.WithContext(ctx).Model(&entities.Product{}).
			Where("category_id IN ? AND status = ?", categoryIDs, "active").
			Select("AVG(price)").Scan(&avgPrice)
	}

	// Mock sales data based on time range
	multiplier := int64(1)
	switch timeRange {
	case "7d":
		multiplier = 1
	case "30d":
		multiplier = 4
	case "90d":
		multiplier = 12
	case "1y":
		multiplier = 52
	default:
		multiplier = 4 // Default to 30 days
	}

	stats.TotalSales = productCount * multiplier * 2
	stats.TotalRevenue = avgPrice * float64(stats.TotalSales)
	stats.AverageOrderValue = avgPrice * 1.5 // Mock AOV

	// Mock growth metrics
	stats.GrowthMetrics = repositories.GrowthMetrics{
		SalesGrowth:   0.12,  // 12% growth
		RevenueGrowth: 0.15,  // 15% growth
		OrderGrowth:   0.10,  // 10% growth
	}

	// Mock top selling products from category hierarchy
	var products []entities.Product
	if len(categoryIDs) > 0 {
		r.db.WithContext(ctx).Where("category_id IN ? AND status = ?", categoryIDs, "active").
			Limit(5).Find(&products)
	}

	for _, product := range products {
		productSales := repositories.ProductSales{
			ProductID:   product.ID,
			ProductName: product.Name,
			SKU:         product.SKU,
			Quantity:    multiplier * 3, // Mock quantity
			Revenue:     product.Price * float64(multiplier) * 3,
		}
		stats.TopSellingProducts = append(stats.TopSellingProducts, productSales)
	}

	return stats, nil
}
