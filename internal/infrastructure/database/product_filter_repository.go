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

type productFilterRepository struct {
	db *gorm.DB
}

// NewProductFilterRepository creates a new product filter repository
func NewProductFilterRepository(db *gorm.DB) repositories.ProductFilterRepository {
	return &productFilterRepository{db: db}
}

// FilterProducts performs advanced product filtering
func (r *productFilterRepository) FilterProducts(ctx context.Context, params repositories.AdvancedFilterParams) (*repositories.FilteredProductResult, error) {
	// Simple implementation for now to avoid GORM issues
	var products []*entities.Product
	query := r.db.WithContext(ctx).Model(&entities.Product{})

	// Apply text search (simplified for now)
	if params.Query != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ? OR sku ILIKE ?",
			"%"+params.Query+"%", "%"+params.Query+"%", "%"+params.Query+"%")
	}

	// Apply category filters
	if len(params.CategoryIDs) > 0 {
		query = query.Where("category_id IN ?", params.CategoryIDs)
	}

	// Apply brand filters
	if len(params.BrandIDs) > 0 {
		query = query.Where("brand_id IN ?", params.BrandIDs)
	}

	// Apply price filters
	if params.MinPrice != nil {
		query = query.Where("COALESCE(sale_price, price) >= ?", *params.MinPrice)
	}
	if params.MaxPrice != nil {
		query = query.Where("COALESCE(sale_price, price) <= ?", *params.MaxPrice)
	}

	// Apply stock filters
	if params.InStock != nil && *params.InStock {
		query = query.Where("stock > 0 AND stock_status = ?", entities.StockStatusInStock)
	}
	if params.LowStock != nil && *params.LowStock {
		query = query.Where("stock <= low_stock_threshold AND stock > 0")
	}

	// Apply sale filter
	if params.OnSale != nil && *params.OnSale {
		query = query.Where("sale_price IS NOT NULL AND sale_price > 0 AND (sale_start_date IS NULL OR sale_start_date <= NOW()) AND (sale_end_date IS NULL OR sale_end_date >= NOW())")
	}

	// Apply featured filter
	if params.Featured != nil && *params.Featured {
		query = query.Where("featured = ?", true)
	}

	// Apply product type filters
	if len(params.ProductTypes) > 0 {
		query = query.Where("product_type IN ?", params.ProductTypes)
	}

	// Apply stock status filters
	if len(params.StockStatus) > 0 {
		query = query.Where("stock_status IN ?", params.StockStatus)
	}

	// Apply visibility filters
	if len(params.Visibility) > 0 {
		query = query.Where("visibility IN ?", params.Visibility)
	}

	// Apply tag filters (temporarily disabled)
	/*
	if len(params.Tags) > 0 {
		query = query.Joins("JOIN product_tag_associations pt ON products.id = pt.product_id").
			Joins("JOIN tags t ON pt.product_tag_id = t.id").
			Where("t.name IN ?", params.Tags)
	}
	*/

	// Apply attribute filters (temporarily disabled for debugging)
	/*
	if len(params.Attributes) > 0 {
		for attributeID, values := range params.Attributes {
			if len(values) > 0 {
				subQuery := r.db.Table("product_attribute_values pav").
					Select("pav.product_id").
					Where("pav.attribute_id = ?", attributeID)

				// Handle both term-based and value-based attributes
				termConditions := r.db.Table("product_attribute_terms pat").
					Select("pat.id").
					Where("pat.attribute_id = ? AND (pat.name IN ? OR pat.value IN ?)", attributeID, values, values)

				subQuery = subQuery.Where("(pav.term_id IN (?) OR pav.value IN ?)", termConditions, values)

				query = query.Where("products.id IN (?)", subQuery)
			}
		}
	}
	*/

	// Apply date filters
	if params.CreatedAfter != nil {
		if createdAfter, err := time.Parse("2006-01-02", *params.CreatedAfter); err == nil {
			query = query.Where("created_at >= ?", createdAfter)
		}
	}
	if params.CreatedBefore != nil {
		if createdBefore, err := time.Parse("2006-01-02", *params.CreatedBefore); err == nil {
			query = query.Where("created_at <= ?", createdBefore.Add(24*time.Hour))
		}
	}

	// Apply advanced filters (temporarily disabled)
	/*
	if params.HasImages != nil && *params.HasImages {
		query = query.Joins("JOIN product_images pi ON products.id = pi.product_id")
	}
	if params.HasVariants != nil && *params.HasVariants {
		query = query.Joins("JOIN product_variants pv ON products.id = pv.product_id")
	}
	*/

	// Count total results
	var total int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Model(&entities.Product{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count products: %w", err)
	}

	// Apply sorting
	switch params.SortBy {
	case "price":
		if params.SortOrder == "desc" {
			query = query.Order("COALESCE(sale_price, price) DESC")
		} else {
			query = query.Order("COALESCE(sale_price, price) ASC")
		}
	case "name":
		if params.SortOrder == "desc" {
			query = query.Order("name DESC")
		} else {
			query = query.Order("name ASC")
		}
	case "created_at":
		if params.SortOrder == "desc" {
			query = query.Order("created_at DESC")
		} else {
			query = query.Order("created_at ASC")
		}
	case "popularity":
		// TODO: Implement popularity sorting based on views/sales
		query = query.Order("created_at DESC")
	case "rating":
		// TODO: Implement rating sorting
		query = query.Order("created_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	// Preload relationships
	query = query.
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags")

	// Execute query
	if err := query.Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to filter products: %w", err)
	}

	result := &repositories.FilteredProductResult{
		Products: products,
		Total:    total,
	}

	// Generate facets if requested
	if params.IncludeFacets {
		facets, err := r.generateFacets(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to generate facets: %w", err)
		}
		result.Facets = facets
	}

	return result, nil
}

// GetFilterFacets gets available filter facets for a category
func (r *productFilterRepository) GetFilterFacets(ctx context.Context, categoryID *uuid.UUID) (*repositories.FilterFacets, error) {
	facets := &repositories.FilterFacets{}

	// Get categories
	categories, err := r.getCategoryFacets(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	facets.Categories = categories

	// Get brands
	brands, err := r.getBrandFacets(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	facets.Brands = brands

	// Get attributes
	attributes, err := r.getAttributeFacets(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	facets.Attributes = attributes

	// Get price range
	priceRange, err := r.getPriceRangeFacets(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	facets.PriceRange = priceRange

	// Get stock facets
	stockFacet, err := r.getStockFacets(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	facets.Stock = stockFacet

	// Get tags
	tags, err := r.getTagFacets(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	facets.Tags = tags

	return facets, nil
}

// GetDynamicFilters gets dynamic filters based on current filter state
func (r *productFilterRepository) GetDynamicFilters(ctx context.Context, params repositories.AdvancedFilterParams) (*repositories.FilterFacets, error) {
	// This would generate facets based on the current filter state
	// For now, we'll use the same logic as GetFilterFacets
	var categoryID *uuid.UUID
	if len(params.CategoryIDs) > 0 {
		categoryID = &params.CategoryIDs[0]
	}
	return r.GetFilterFacets(ctx, categoryID)
}

// Helper method to generate facets
func (r *productFilterRepository) generateFacets(ctx context.Context, params repositories.AdvancedFilterParams) (*repositories.FilterFacets, error) {
	var categoryID *uuid.UUID
	if len(params.CategoryIDs) > 0 {
		categoryID = &params.CategoryIDs[0]
	}
	return r.GetFilterFacets(ctx, categoryID)
}

// getCategoryFacets gets category facets with hierarchy support
func (r *productFilterRepository) getCategoryFacets(ctx context.Context, categoryID *uuid.UUID) ([]repositories.FilterCategoryFacet, error) {
	var facets []repositories.FilterCategoryFacet

	// First get all categories
	var categories []struct {
		ID   uuid.UUID
		Name string
		Slug string
	}

	categoryQuery := `
		SELECT id, name, slug
		FROM categories
		WHERE is_active = true
	`

	var categoryArgs []interface{}
	if categoryID != nil {
		categoryQuery += " AND parent_id = ?"
		categoryArgs = append(categoryArgs, *categoryID)
	} else {
		categoryQuery += " AND parent_id IS NULL"
	}

	categoryQuery += " ORDER BY name"

	rows, err := r.db.WithContext(ctx).Raw(categoryQuery, categoryArgs...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cat struct {
			ID   uuid.UUID
			Name string
			Slug string
		}
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Slug); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	// Now count products for each category (including subcategories)
	for _, cat := range categories {
		count, err := r.countProductsInCategoryHierarchy(ctx, cat.ID)
		if err != nil {
			return nil, err
		}

		facets = append(facets, repositories.FilterCategoryFacet{
			ID:    cat.ID,
			Name:  cat.Name,
			Slug:  cat.Slug,
			Count: count,
		})
	}

	return facets, nil
}

// countProductsInCategoryHierarchy counts products in a category and all its subcategories
func (r *productFilterRepository) countProductsInCategoryHierarchy(ctx context.Context, categoryID uuid.UUID) (int, error) {
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
	rows, err := r.db.WithContext(ctx).Raw(categoryQuery, categoryID).Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
		categoryIDs = append(categoryIDs, id)
	}

	if len(categoryIDs) == 0 {
		return 0, nil
	}

	// Count products in all these categories
	var count int64
	err = r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("category_id IN ? AND status = ?", categoryIDs, entities.ProductStatusActive).
		Count(&count).Error

	return int(count), err
}

// getBrandFacets gets brand facets
func (r *productFilterRepository) getBrandFacets(ctx context.Context, categoryID *uuid.UUID) ([]repositories.FilterBrandFacet, error) {
	var facets []repositories.FilterBrandFacet

	query := `
		SELECT b.id, b.name, b.slug, COALESCE(b.logo, '') as logo, COUNT(p.id) as count
		FROM brands b
		LEFT JOIN products p ON b.id = p.brand_id AND p.status = 'active'
	`

	var args []interface{}
	if categoryID != nil {
		query += " AND p.category_id = ?"
		args = append(args, *categoryID)
	}

	query += " GROUP BY b.id, b.name, b.slug, b.logo HAVING COUNT(p.id) > 0 ORDER BY b.name"

	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var facet repositories.FilterBrandFacet
		if err := rows.Scan(&facet.ID, &facet.Name, &facet.Slug, &facet.Logo, &facet.Count); err != nil {
			return nil, err
		}
		facets = append(facets, facet)
	}

	return facets, nil
}

// getAttributeFacets gets attribute facets
func (r *productFilterRepository) getAttributeFacets(ctx context.Context, categoryID *uuid.UUID) ([]repositories.FilterAttributeFacet, error) {
	var facets []repositories.FilterAttributeFacet

	// Get attributes
	query := `
		SELECT DISTINCT pa.id, pa.name, pa.slug, pa.type, pa.position
		FROM product_attributes pa
		JOIN product_attribute_values pav ON pa.id = pav.attribute_id
		JOIN products p ON pav.product_id = p.id
		WHERE pa.is_visible = true AND p.status = 'active'
	`

	var args []interface{}
	if categoryID != nil {
		query += " AND p.category_id = ?"
		args = append(args, *categoryID)
	}

	query += " ORDER BY pa.position, pa.name"

	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var facet repositories.FilterAttributeFacet
		var position int
		if err := rows.Scan(&facet.ID, &facet.Name, &facet.Slug, &facet.Type, &position); err != nil {
			return nil, err
		}

		// Get terms for this attribute
		terms, err := r.getAttributeTermFacets(ctx, facet.ID, categoryID)
		if err != nil {
			return nil, err
		}
		facet.Terms = terms

		if len(terms) > 0 {
			facets = append(facets, facet)
		}
	}

	return facets, nil
}

// getAttributeTermFacets gets attribute term facets
func (r *productFilterRepository) getAttributeTermFacets(ctx context.Context, attributeID uuid.UUID, categoryID *uuid.UUID) ([]repositories.FilterAttributeTermFacet, error) {
	var facets []repositories.FilterAttributeTermFacet

	query := `
		SELECT pat.id, pat.name, COALESCE(pat.value, '') as value,
		       COALESCE(pat.color, '') as color, COALESCE(pat.image, '') as image,
		       COUNT(DISTINCT p.id) as count
		FROM product_attribute_terms pat
		JOIN product_attribute_values pav ON pat.id = pav.term_id
		JOIN products p ON pav.product_id = p.id
		WHERE pat.attribute_id = ? AND p.status = 'active'
	`

	args := []interface{}{attributeID}
	if categoryID != nil {
		query += " AND p.category_id = ?"
		args = append(args, *categoryID)
	}

	query += " GROUP BY pat.id, pat.name, pat.value, pat.color, pat.image ORDER BY pat.position, pat.name"

	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var facet repositories.FilterAttributeTermFacet
		if err := rows.Scan(&facet.ID, &facet.Name, &facet.Value, &facet.Color, &facet.Image, &facet.Count); err != nil {
			return nil, err
		}
		facets = append(facets, facet)
	}

	return facets, nil
}

// getPriceRangeFacets gets price range facets
func (r *productFilterRepository) getPriceRangeFacets(ctx context.Context, categoryID *uuid.UUID) (repositories.FilterPriceRangeFacet, error) {
	var facet repositories.FilterPriceRangeFacet

	query := `
		SELECT MIN(COALESCE(sale_price, price)) as min_price,
		       MAX(COALESCE(sale_price, price)) as max_price
		FROM products
		WHERE status = 'active'
	`

	var args []interface{}
	if categoryID != nil {
		query += " AND category_id = ?"
		args = append(args, *categoryID)
	}

	var row *sql.Row
	if len(args) > 0 {
		row = r.db.WithContext(ctx).Raw(query, args...).Row()
	} else {
		row = r.db.WithContext(ctx).Raw(query).Row()
	}
	if err := row.Scan(&facet.Min, &facet.Max); err != nil {
		return facet, err
	}

	// Generate price ranges
	ranges := []repositories.FilterPriceRange{
		{Min: nil, Max: &[]float64{50}[0], Label: "Under $50"},
		{Min: &[]float64{50}[0], Max: &[]float64{100}[0], Label: "$50 - $100"},
		{Min: &[]float64{100}[0], Max: &[]float64{200}[0], Label: "$100 - $200"},
		{Min: &[]float64{200}[0], Max: &[]float64{500}[0], Label: "$200 - $500"},
		{Min: &[]float64{500}[0], Max: nil, Label: "Over $500"},
	}

	// Count products in each range
	for i, priceRange := range ranges {
		countQuery := "SELECT COUNT(*) FROM products WHERE status = 'active'"
		countArgs := []interface{}{}

		if categoryID != nil {
			countQuery += " AND category_id = ?"
			countArgs = append(countArgs, *categoryID)
		}

		if priceRange.Min != nil {
			countQuery += " AND COALESCE(sale_price, price) >= ?"
			countArgs = append(countArgs, *priceRange.Min)
		}
		if priceRange.Max != nil {
			countQuery += " AND COALESCE(sale_price, price) <= ?"
			countArgs = append(countArgs, *priceRange.Max)
		}

		var count int
		r.db.WithContext(ctx).Raw(countQuery, countArgs...).Scan(&count)
		ranges[i].Count = count
	}

	facet.Ranges = ranges
	return facet, nil
}

// getStockFacets gets stock facets
func (r *productFilterRepository) getStockFacets(ctx context.Context, categoryID *uuid.UUID) (repositories.FilterStockFacet, error) {
	var facet repositories.FilterStockFacet

	baseQuery := "SELECT COUNT(*) FROM products WHERE status = 'active'"
	categoryFilter := ""
	stockArgs := []interface{}{}

	if categoryID != nil {
		categoryFilter = " AND category_id = ?"
		stockArgs = append(stockArgs, *categoryID)
	}

	// In stock
	r.db.WithContext(ctx).Raw(baseQuery+" AND stock > 0"+categoryFilter, stockArgs...).Scan(&facet.InStock)

	// Low stock
	r.db.WithContext(ctx).Raw(baseQuery+" AND stock <= low_stock_threshold AND stock > 0"+categoryFilter, stockArgs...).Scan(&facet.LowStock)

	// Out of stock
	r.db.WithContext(ctx).Raw(baseQuery+" AND stock = 0"+categoryFilter, stockArgs...).Scan(&facet.OutStock)

	// On sale
	r.db.WithContext(ctx).Raw(baseQuery+" AND sale_price IS NOT NULL AND sale_price > 0"+categoryFilter, stockArgs...).Scan(&facet.OnSale)

	// Featured
	r.db.WithContext(ctx).Raw(baseQuery+" AND featured = true"+categoryFilter, stockArgs...).Scan(&facet.Featured)

	return facet, nil
}

// getTagFacets gets tag facets
func (r *productFilterRepository) getTagFacets(ctx context.Context, categoryID *uuid.UUID) ([]repositories.FilterTagFacet, error) {
	var facets []repositories.FilterTagFacet

	query := `
		SELECT t.name, COUNT(DISTINCT p.id) as count
		FROM tags t
		JOIN product_tag_associations pt ON t.id = pt.product_tag_id
		JOIN products p ON pt.product_id = p.id
		WHERE p.status = 'active'
	`

	var tagArgs []interface{}
	if categoryID != nil {
		query += " AND p.category_id = ?"
		tagArgs = append(tagArgs, *categoryID)
	}

	query += " GROUP BY t.name ORDER BY count DESC, t.name LIMIT 20"

	rows, err := r.db.WithContext(ctx).Raw(query, tagArgs...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var facet repositories.FilterTagFacet
		if err := rows.Scan(&facet.Name, &facet.Count); err != nil {
			return nil, err
		}
		facets = append(facets, facet)
	}

	return facets, nil
}

// SaveFilterSet saves a filter set
func (r *productFilterRepository) SaveFilterSet(ctx context.Context, filterSet *entities.FilterSet) error {
	return r.db.WithContext(ctx).Create(filterSet).Error
}

// GetFilterSet gets a filter set by ID
func (r *productFilterRepository) GetFilterSet(ctx context.Context, id uuid.UUID) (*entities.FilterSet, error) {
	var filterSet entities.FilterSet
	err := r.db.WithContext(ctx).First(&filterSet, id).Error
	if err != nil {
		return nil, err
	}
	return &filterSet, nil
}

// GetUserFilterSets gets filter sets for a user
func (r *productFilterRepository) GetUserFilterSets(ctx context.Context, userID uuid.UUID) ([]*entities.FilterSet, error) {
	var filterSets []*entities.FilterSet
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Find(&filterSets).Error
	return filterSets, err
}

// GetSessionFilterSets gets filter sets for a session
func (r *productFilterRepository) GetSessionFilterSets(ctx context.Context, sessionID string) ([]*entities.FilterSet, error) {
	var filterSets []*entities.FilterSet
	err := r.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("updated_at DESC").
		Find(&filterSets).Error
	return filterSets, err
}

// UpdateFilterSet updates a filter set
func (r *productFilterRepository) UpdateFilterSet(ctx context.Context, filterSet *entities.FilterSet) error {
	return r.db.WithContext(ctx).Save(filterSet).Error
}

// DeleteFilterSet deletes a filter set
func (r *productFilterRepository) DeleteFilterSet(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.FilterSet{}, id).Error
}

// TrackFilterUsage tracks filter usage
func (r *productFilterRepository) TrackFilterUsage(ctx context.Context, usage *entities.FilterUsage) error {
	return r.db.WithContext(ctx).Create(usage).Error
}

// GetFilterAnalytics gets filter analytics
func (r *productFilterRepository) GetFilterAnalytics(ctx context.Context, days int) (map[string]interface{}, error) {
	analytics := make(map[string]interface{})

	// Most used filters
	var topFilters []struct {
		FilterType  string `json:"filter_type"`
		FilterKey   string `json:"filter_key"`
		UsageCount  int    `json:"usage_count"`
	}

	err := r.db.WithContext(ctx).
		Table("filter_usage").
		Select("filter_type, filter_key, COUNT(*) as usage_count").
		Where("created_at >= ?", time.Now().AddDate(0, 0, -days)).
		Group("filter_type, filter_key").
		Order("usage_count DESC").
		Limit(10).
		Find(&topFilters).Error

	if err != nil {
		return nil, err
	}

	analytics["top_filters"] = topFilters

	// Filter usage by day
	var dailyUsage []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	err = r.db.WithContext(ctx).
		Table("filter_usage").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", time.Now().AddDate(0, 0, -days)).
		Group("DATE(created_at)").
		Order("date").
		Find(&dailyUsage).Error

	if err != nil {
		return nil, err
	}

	analytics["daily_usage"] = dailyUsage

	return analytics, nil
}

// GetPopularFilters gets popular filters
func (r *productFilterRepository) GetPopularFilters(ctx context.Context, limit int) ([]*entities.FilterUsage, error) {
	var filters []*entities.FilterUsage
	err := r.db.WithContext(ctx).
		Select("filter_type, filter_key, filter_value, COUNT(*) as result_count").
		Group("filter_type, filter_key, filter_value").
		Order("result_count DESC").
		Limit(limit).
		Find(&filters).Error
	return filters, err
}

// UpdateFilterOptions updates filter options for a category
func (r *productFilterRepository) UpdateFilterOptions(ctx context.Context, categoryID *uuid.UUID) error {
	// This would update the product_filter_options table
	// For now, we'll implement a basic version
	return nil
}

// GetFilterOptions gets filter options for a category
func (r *productFilterRepository) GetFilterOptions(ctx context.Context, categoryID *uuid.UUID) ([]*entities.ProductFilterOption, error) {
	var options []*entities.ProductFilterOption
	query := r.db.WithContext(ctx)

	if categoryID != nil {
		query = query.Where("category_id = ? OR category_id IS NULL", *categoryID)
	} else {
		query = query.Where("category_id IS NULL")
	}

	err := query.Where("is_active = ?", true).
		Order("position, filter_type, label").
		Find(&options).Error

	return options, err
}

// GetAttributeFilters gets attribute filters for a category
func (r *productFilterRepository) GetAttributeFilters(ctx context.Context, categoryID *uuid.UUID) ([]*entities.ProductAttribute, error) {
	var attributes []*entities.ProductAttribute
	query := r.db.WithContext(ctx).
		Preload("Terms").
		Where("is_visible = ?", true)

	if categoryID != nil {
		// Get attributes that are used by products in this category
		query = query.Joins("JOIN product_attribute_values pav ON product_attributes.id = pav.attribute_id").
			Joins("JOIN products p ON pav.product_id = p.id").
			Where("p.category_id = ?", *categoryID).
			Distinct()
	}

	err := query.Order("position, name").Find(&attributes).Error
	return attributes, err
}

// GetAttributeTerms gets attribute terms for an attribute
func (r *productFilterRepository) GetAttributeTerms(ctx context.Context, attributeID uuid.UUID, categoryID *uuid.UUID) ([]*entities.ProductAttributeTerm, error) {
	var terms []*entities.ProductAttributeTerm
	query := r.db.WithContext(ctx).Where("attribute_id = ?", attributeID)

	if categoryID != nil {
		// Get terms that are used by products in this category
		query = query.Joins("JOIN product_attribute_values pav ON product_attribute_terms.id = pav.term_id").
			Joins("JOIN products p ON pav.product_id = p.id").
			Where("p.category_id = ?", *categoryID).
			Distinct()
	}

	err := query.Order("position, name").Find(&terms).Error
	return terms, err
}

// GetFilterSuggestions gets filter suggestions based on query
func (r *productFilterRepository) GetFilterSuggestions(ctx context.Context, query string, limit int) ([]string, error) {
	var suggestions []string

	// Get suggestions from product names, brands, categories, and attributes
	sqlQuery := `
		(SELECT DISTINCT name as suggestion FROM products WHERE name ILIKE ? LIMIT ?)
		UNION
		(SELECT DISTINCT name as suggestion FROM brands WHERE name ILIKE ? LIMIT ?)
		UNION
		(SELECT DISTINCT name as suggestion FROM categories WHERE name ILIKE ? LIMIT ?)
		UNION
		(SELECT DISTINCT name as suggestion FROM product_attribute_terms WHERE name ILIKE ? LIMIT ?)
		ORDER BY suggestion
		LIMIT ?
	`

	searchTerm := "%" + query + "%"
	rows, err := r.db.WithContext(ctx).Raw(sqlQuery, searchTerm, limit/4, searchTerm, limit/4, searchTerm, limit/4, searchTerm, limit/4, limit).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var suggestion string
		if err := rows.Scan(&suggestion); err != nil {
			return nil, err
		}
		suggestions = append(suggestions, suggestion)
	}

	return suggestions, nil
}

// GetRelatedFilters gets related filters based on current filters
func (r *productFilterRepository) GetRelatedFilters(ctx context.Context, currentFilters repositories.AdvancedFilterParams) ([]string, error) {
	var related []string

	// This would implement collaborative filtering logic
	// For now, return empty slice
	return related, nil
}
