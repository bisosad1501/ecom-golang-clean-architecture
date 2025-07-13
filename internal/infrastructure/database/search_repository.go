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

type searchRepository struct {
	db *gorm.DB
}

// NewSearchRepository creates a new search repository
func NewSearchRepository(db *gorm.DB) repositories.SearchRepository {
	return &searchRepository{db: db}
}

// RecordSearchEvent records a search event
func (r *searchRepository) RecordSearchEvent(ctx context.Context, event *entities.SearchEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

// GetSearchEvents retrieves search events with filters
func (r *searchRepository) GetSearchEvents(ctx context.Context, filters repositories.SearchEventFilters) ([]*entities.SearchEvent, error) {
	query := r.db.WithContext(ctx).Model(&entities.SearchEvent{})

	// Apply filters
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if filters.Query != "" {
		query = query.Where("query ILIKE ?", "%"+filters.Query+"%")
	}
	if filters.StartDate != nil {
		query = query.Where("created_at >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", *filters.EndDate)
	}

	// Apply sorting
	orderBy := "created_at DESC"
	if filters.SortBy != "" {
		direction := "ASC"
		if strings.ToUpper(filters.SortOrder) == "DESC" {
			direction = "DESC"
		}
		orderBy = filters.SortBy + " " + direction
	}
	query = query.Order(orderBy)

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var events []*entities.SearchEvent
	err := query.Find(&events).Error
	return events, err
}

// GetPopularSearchTerms retrieves popular search terms
func (r *searchRepository) GetPopularSearchTerms(ctx context.Context, limit int, period string) ([]*entities.PopularSearch, error) {
	query := r.db.WithContext(ctx).Model(&entities.PopularSearch{}).
		Where("period = ?", period).
		Order("search_count DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	var popularSearches []*entities.PopularSearch
	err := query.Find(&popularSearches).Error
	return popularSearches, err
}

// GetSearchSuggestions retrieves enhanced search suggestions with fuzzy matching
func (r *searchRepository) GetSearchSuggestions(ctx context.Context, query string, limit int) ([]*entities.SearchSuggestion, error) {
	if limit <= 0 {
		limit = 10
	}

	var suggestions []*entities.SearchSuggestion

	// Use the enhanced function for better suggestions
	type SuggestionResult struct {
		Suggestion   string `json:"suggestion"`
		Frequency    int    `json:"frequency"`
		ResultCount  int    `json:"result_count"`
	}

	var results []SuggestionResult
	err := r.db.WithContext(ctx).Raw(
		"SELECT * FROM get_search_suggestions_with_synonyms(?, ?)",
		query, limit,
	).Scan(&results).Error

	if err != nil {
		// Fallback to basic search if function fails
		dbQuery := r.db.WithContext(ctx).Model(&entities.SearchSuggestion{}).
			Where("is_active = true")

		if query != "" {
			dbQuery = dbQuery.Where("query ILIKE ? OR query % ?", query+"%", query)
		}

		dbQuery = dbQuery.Order("search_count DESC")

		if limit > 0 {
			dbQuery = dbQuery.Limit(limit)
		}

		err = dbQuery.Find(&suggestions).Error
		return suggestions, err
	}

	// Convert results to entities
	for _, result := range results {
		suggestion := &entities.SearchSuggestion{
			Query:       result.Suggestion,
			Frequency:   result.Frequency,
			ResultCount: result.ResultCount,
		}
		suggestions = append(suggestions, suggestion)
	}

	return suggestions, nil
}

// UpdateSearchSuggestion updates or creates a search suggestion using enhanced function
func (r *searchRepository) UpdateSearchSuggestion(ctx context.Context, query string) error {
	// Use the enhanced function for better suggestion management
	err := r.db.WithContext(ctx).Exec(
		"SELECT update_search_suggestion(?, ?)",
		query, 0, // result_count will be updated separately
	).Error

	if err != nil {
		// Fallback to manual update if function fails
		result := r.db.WithContext(ctx).Model(&entities.SearchSuggestion{}).
			Where("query = ?", query).
			Update("search_count", gorm.Expr("search_count + 1"))

		if result.Error != nil {
			return result.Error
		}

		// If no rows affected, create new suggestion
		if result.RowsAffected == 0 {
			suggestion := &entities.SearchSuggestion{
				Query:       query,
				SearchCount: 1,
				IsActive:    true,
			}
			return r.db.WithContext(ctx).Create(suggestion).Error
		}
	}

	return nil
}

// SaveSearchHistory saves user search history
func (r *searchRepository) SaveSearchHistory(ctx context.Context, history *entities.SearchHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetUserSearchHistory retrieves user search history
func (r *searchRepository) GetUserSearchHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*entities.SearchHistory, error) {
	query := r.db.WithContext(ctx).Model(&entities.SearchHistory{}).
		Where("user_id = ?", userID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	var history []*entities.SearchHistory
	err := query.Find(&history).Error
	return history, err
}

// ClearUserSearchHistory clears user search history
func (r *searchRepository) ClearUserSearchHistory(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&entities.SearchHistory{}).Error
}

// SaveSearchFilter saves a search filter
func (r *searchRepository) SaveSearchFilter(ctx context.Context, filter *entities.SearchFilter) error {
	return r.db.WithContext(ctx).Create(filter).Error
}

// GetUserSearchFilters retrieves user search filters
func (r *searchRepository) GetUserSearchFilters(ctx context.Context, userID uuid.UUID) ([]*entities.SearchFilter, error) {
	var filters []*entities.SearchFilter
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("is_default DESC, usage_count DESC, created_at DESC").
		Find(&filters).Error
	return filters, err
}

// GetSearchFilter retrieves a search filter by ID
func (r *searchRepository) GetSearchFilter(ctx context.Context, id uuid.UUID) (*entities.SearchFilter, error) {
	var filter entities.SearchFilter
	err := r.db.WithContext(ctx).First(&filter, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &filter, nil
}

// UpdateSearchFilter updates a search filter
func (r *searchRepository) UpdateSearchFilter(ctx context.Context, filter *entities.SearchFilter) error {
	return r.db.WithContext(ctx).Save(filter).Error
}

// DeleteSearchFilter deletes a search filter
func (r *searchRepository) DeleteSearchFilter(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.SearchFilter{}, "id = ?", id).Error
}

// FullTextSearch performs advanced full-text search with enhanced ranking and fuzzy matching
func (r *searchRepository) FullTextSearch(ctx context.Context, params repositories.FullTextSearchParams) ([]*entities.Product, int64, error) {
	query := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Tags")

	// Enhanced full-text search with multiple search strategies
	if params.Query != "" {
		// Strategy 1: Use pre-computed search vector for best performance
		searchQuery := "plainto_tsquery('english', ?)"

		// Strategy 2: Add fuzzy matching for typos and partial matches
		fuzzyCondition := "(name % ? OR sku % ?)"

		// Strategy 3: Add exact phrase matching for high relevance
		exactCondition := "(name ILIKE ? OR sku ILIKE ?)"

		// Strategy 4: Add synonym expansion
		synonymCondition := r.buildSynonymCondition(params.Query)

		// Combine all search strategies
		searchCondition := fmt.Sprintf(
			"(search_vector @@ %s) OR %s OR %s",
			searchQuery, fuzzyCondition, exactCondition,
		)

		if synonymCondition != "" {
			searchCondition += " OR " + synonymCondition
		}

		query = query.Where(
			searchCondition,
			params.Query, // for search_vector
			params.Query, params.Query, // for fuzzy matching
			"%"+params.Query+"%", "%"+params.Query+"%", // for exact matching
		)
	}

	// Apply filters
	if len(params.CategoryIDs) > 0 {
		query = query.Where("category_id IN ?", params.CategoryIDs)
	}
	if len(params.BrandIDs) > 0 {
		query = query.Where("brand_id IN ?", params.BrandIDs)
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
	if params.OnSale != nil && *params.OnSale {
		query = query.Where("sale_price IS NOT NULL AND sale_price > 0")
	}

	// Advanced filters
	if params.MinRating != nil {
		query = query.Where("rating >= ?", *params.MinRating)
	}
	if params.MaxRating != nil {
		query = query.Where("rating <= ?", *params.MaxRating)
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
	if params.AvailabilityStatus != nil {
		switch *params.AvailabilityStatus {
		case "in_stock":
			query = query.Where("stock > 0")
		case "out_of_stock":
			query = query.Where("stock = 0")
		case "low_stock":
			query = query.Where("stock > 0 AND stock <= low_stock_threshold")
		}
	}
	if params.CreatedAfter != nil {
		query = query.Where("created_at >= ?", *params.CreatedAfter)
	}
	if params.CreatedBefore != nil {
		query = query.Where("created_at <= ?", *params.CreatedBefore)
	}
	if params.UpdatedAfter != nil {
		query = query.Where("updated_at >= ?", *params.UpdatedAfter)
	}
	if params.UpdatedBefore != nil {
		query = query.Where("updated_at <= ?", *params.UpdatedBefore)
	}
	if params.MinWeight != nil {
		query = query.Where("weight >= ?", *params.MinWeight)
	}
	if params.MaxWeight != nil {
		query = query.Where("weight <= ?", *params.MaxWeight)
	}
	if params.ShippingClass != nil {
		query = query.Where("shipping_class = ?", *params.ShippingClass)
	}
	if params.TaxClass != nil {
		query = query.Where("tax_class = ?", *params.TaxClass)
	}
	if params.MinDiscountPercent != nil {
		query = query.Where("CASE WHEN sale_price IS NOT NULL AND sale_price > 0 THEN ((price - sale_price) / price * 100) ELSE 0 END >= ?", *params.MinDiscountPercent)
	}
	if params.MaxDiscountPercent != nil {
		query = query.Where("CASE WHEN sale_price IS NOT NULL AND sale_price > 0 THEN ((price - sale_price) / price * 100) ELSE 0 END <= ?", *params.MaxDiscountPercent)
	}
	if params.IsDigital != nil {
		query = query.Where("is_digital = ?", *params.IsDigital)
	}
	if params.RequiresShipping != nil {
		query = query.Where("requires_shipping = ?", *params.RequiresShipping)
	}
	if params.AllowBackorder != nil {
		query = query.Where("allow_backorder = ?", *params.AllowBackorder)
	}
	if params.TrackQuantity != nil {
		query = query.Where("track_quantity = ?", *params.TrackQuantity)
	}

	// Tags filter
	if len(params.Tags) > 0 {
		query = query.Joins("JOIN product_tag_associations pta ON products.id = pta.product_id").
			Joins("JOIN tags t ON pta.product_tag_id = t.id").
			Where("t.name IN ?", params.Tags).
			Group("products.id")
	}

	// Count total results
	var total int64
	countQuery := query
	if err := countQuery.Model(&entities.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
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
	return products, total, err
}

// buildSynonymCondition builds search condition with synonym expansion
func (r *searchRepository) buildSynonymCondition(query string) string {
	// Get synonyms for the query - check if query matches any synonym first
	var synonyms []string
	err := r.db.Raw(`
		SELECT unnest(synonyms) as synonym
		FROM search_synonyms
		WHERE ? = ANY(synonyms) AND is_active = true
		UNION
		SELECT unnest(synonyms) as synonym
		FROM search_synonyms
		WHERE term ILIKE ? AND is_active = true
	`, query, "%"+query+"%").Pluck("synonym", &synonyms).Error

	if err != nil || len(synonyms) == 0 {
		return ""
	}

	// Build condition for synonyms
	conditions := make([]string, len(synonyms))
	for i, synonym := range synonyms {
		conditions[i] = fmt.Sprintf("(name ILIKE '%%%s%%' OR description ILIKE '%%%s%%')", synonym, synonym)
	}

	return "(" + strings.Join(conditions, " OR ") + ")"
}

// RecordSearchAnalytics records search analytics for performance tracking
func (r *searchRepository) RecordSearchAnalytics(ctx context.Context, query string, resultCount int) error {
	// Use UPSERT to update existing or create new analytics record
	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO search_analytics (query, result_count, search_date, total_searches)
		VALUES (?, ?, CURRENT_DATE, 1)
		ON CONFLICT (query, search_date) DO UPDATE SET
			total_searches = search_analytics.total_searches + 1,
			result_count = EXCLUDED.result_count,
			updated_at = CURRENT_TIMESTAMP
	`, query, resultCount).Error

	return err
}

// GetSearchAnalytics retrieves search analytics for admin dashboard
func (r *searchRepository) GetSearchAnalytics(ctx context.Context, startDate, endDate time.Time, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			query,
			SUM(total_searches) as total_searches,
			AVG(result_count) as avg_result_count,
			AVG(click_through_rate) as avg_ctr,
			AVG(conversion_rate) as avg_conversion_rate,
			MAX(search_date) as last_searched
		FROM search_analytics
		WHERE search_date BETWEEN ? AND ?
		GROUP BY query
		ORDER BY total_searches DESC
		LIMIT ?
	`, startDate, endDate, limit).Scan(&results).Error

	return results, err
}

// buildSortOrder builds the enhanced sort order clause with advanced ranking
func (r *searchRepository) buildSortOrder(sortBy, sortOrder, searchQuery string) string {
	direction := "ASC"
	if strings.ToUpper(sortOrder) == "DESC" {
		direction = "DESC"
	}

	switch sortBy {
	case "relevance":
		if searchQuery != "" {
			// Enhanced relevance ranking with multiple factors
			return fmt.Sprintf(`
				(
					-- Full-text search ranking (highest weight)
					ts_rank(search_vector, plainto_tsquery('english', '%s')) * 4.0 +

					-- Exact name match bonus
					CASE WHEN name ILIKE '%%%s%%' THEN 3.0 ELSE 0 END +

					-- SKU match bonus
					CASE WHEN sku ILIKE '%%%s%%' THEN 2.0 ELSE 0 END +

					-- Fuzzy similarity bonus
					similarity(name, '%s') * 2.0 +

					-- Featured product bonus
					CASE WHEN featured = true THEN 1.5 ELSE 0 END +

					-- Stock availability bonus
					CASE WHEN stock > 0 THEN 1.0 ELSE 0 END +

					-- Recent products bonus (within 30 days)
					CASE WHEN created_at > NOW() - INTERVAL '30 days' THEN 0.5 ELSE 0 END
				) %s, created_at DESC`,
				searchQuery, searchQuery, searchQuery, searchQuery, direction)
		}
		return "featured DESC, stock DESC, created_at DESC"
	case "price":
		return "price " + direction
	case "name":
		return "name " + direction
	case "created_at":
		return "created_at " + direction
	case "popularity":
		return "view_count DESC, created_at DESC"
	case "rating":
		return "average_rating DESC, review_count DESC"
	default:
		return "created_at DESC"
	}
}

// GetSearchFacets retrieves search facets
func (r *searchRepository) GetSearchFacets(ctx context.Context, query string) (*repositories.SearchFacets, error) {
	// This is a simplified implementation
	// In a real-world scenario, you might want to use more sophisticated aggregation
	
	facets := &repositories.SearchFacets{}
	
	// Get category facets
	var categoryFacets []repositories.CategoryFacet
	err := r.db.WithContext(ctx).Raw(`
		SELECT c.id, c.name, COUNT(p.id) as product_count
		FROM categories c
		LEFT JOIN products p ON c.id = p.category_id
		WHERE p.status = 'active'
		GROUP BY c.id, c.name
		HAVING COUNT(p.id) > 0
		ORDER BY product_count DESC
		LIMIT 10
	`).Scan(&categoryFacets).Error
	if err != nil {
		return nil, err
	}
	facets.Categories = categoryFacets

	// Get brand facets
	var brandFacets []repositories.BrandFacet
	err = r.db.WithContext(ctx).Raw(`
		SELECT b.id, b.name, COUNT(p.id) as product_count
		FROM brands b
		LEFT JOIN products p ON b.id = p.brand_id
		WHERE p.status = 'active'
		GROUP BY b.id, b.name
		HAVING COUNT(p.id) > 0
		ORDER BY product_count DESC
		LIMIT 10
	`).Scan(&brandFacets).Error
	if err != nil {
		return nil, err
	}
	facets.Brands = brandFacets

	// Get price range facets
	var priceStats struct {
		MinPrice float64 `json:"min_price"`
		MaxPrice float64 `json:"max_price"`
	}
	err = r.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(MIN(price), 0) as min_price,
			COALESCE(MAX(price), 0) as max_price
		FROM products
		WHERE status = 'active' AND price > 0
	`).Scan(&priceStats).Error
	if err != nil {
		return nil, err
	}

	// Create price ranges
	priceRanges := []repositories.PriceRange{}
	if priceStats.MaxPrice > 0 {
		ranges := []struct {
			min   *float64
			max   *float64
			label string
		}{
			{nil, float64Ptr(50), "Under $50"},
			{float64Ptr(50), float64Ptr(100), "$50 - $100"},
			{float64Ptr(100), float64Ptr(250), "$100 - $250"},
			{float64Ptr(250), float64Ptr(500), "$250 - $500"},
			{float64Ptr(500), float64Ptr(1000), "$500 - $1000"},
			{float64Ptr(1000), nil, "Over $1000"},
		}

		for _, rng := range ranges {
			var count int64
			query := r.db.WithContext(ctx).Model(&entities.Product{}).Where("status = 'active'")
			if rng.min != nil {
				query = query.Where("price >= ?", *rng.min)
			}
			if rng.max != nil {
				query = query.Where("price <= ?", *rng.max)
			}
			query.Count(&count)

			if count > 0 {
				priceRanges = append(priceRanges, repositories.PriceRange{
					Min:          rng.min,
					Max:          rng.max,
					Label:        rng.label,
					ProductCount: count,
				})
			}
		}
	}

	facets.PriceRange = repositories.PriceRangeFacet{
		MinPrice: priceStats.MinPrice,
		MaxPrice: priceStats.MaxPrice,
		Ranges:   priceRanges,
	}

	// Get tag facets
	var tagFacets []repositories.TagFacet
	err = r.db.WithContext(ctx).Raw(`
		SELECT t.id, t.name, COUNT(pta.product_id) as product_count
		FROM tags t
		LEFT JOIN product_tag_associations pta ON t.id = pta.product_tag_id
		LEFT JOIN products p ON pta.product_id = p.id
		WHERE p.status = 'active'
		GROUP BY t.id, t.name
		HAVING COUNT(pta.product_id) > 0
		ORDER BY product_count DESC
		LIMIT 10
	`).Scan(&tagFacets).Error
	if err != nil {
		return nil, err
	}
	facets.Tags = tagFacets

	// Get status facets
	var statusFacets []repositories.StatusFacet
	err = r.db.WithContext(ctx).Raw(`
		SELECT
			status,
			CASE
				WHEN status = 'active' THEN 'Active'
				WHEN status = 'inactive' THEN 'Inactive'
				WHEN status = 'draft' THEN 'Draft'
				ELSE status
			END as label,
			COUNT(*) as product_count
		FROM products
		GROUP BY status
		HAVING COUNT(*) > 0
		ORDER BY product_count DESC
	`).Scan(&statusFacets).Error
	if err != nil {
		return nil, err
	}
	facets.Status = statusFacets

	// Get product type facets
	var productTypeFacets []repositories.ProductTypeFacet
	err = r.db.WithContext(ctx).Raw(`
		SELECT
			product_type as type,
			CASE
				WHEN product_type = 'simple' THEN 'Simple Product'
				WHEN product_type = 'variable' THEN 'Variable Product'
				WHEN product_type = 'grouped' THEN 'Grouped Product'
				ELSE product_type
			END as label,
			COUNT(*) as product_count
		FROM products
		WHERE status = 'active'
		GROUP BY product_type
		HAVING COUNT(*) > 0
		ORDER BY product_count DESC
	`).Scan(&productTypeFacets).Error
	if err != nil {
		return nil, err
	}
	facets.ProductTypes = productTypeFacets

	// Get availability facets
	var availabilityFacets []repositories.AvailabilityFacet
	err = r.db.WithContext(ctx).Raw(`
		SELECT
			availability_status as status,
			availability_label as label,
			COUNT(*) as product_count
		FROM (
			SELECT
				CASE
					WHEN stock > 0 THEN 'in_stock'
					WHEN stock = 0 THEN 'out_of_stock'
					ELSE 'unknown'
				END as availability_status,
				CASE
					WHEN stock > 0 THEN 'In Stock'
					WHEN stock = 0 THEN 'Out of Stock'
					ELSE 'Unknown'
				END as availability_label
			FROM products
			WHERE status = 'active'
		) availability_subquery
		GROUP BY availability_status, availability_label
		HAVING COUNT(*) > 0
		ORDER BY product_count DESC
	`).Scan(&availabilityFacets).Error
	if err != nil {
		return nil, err
	}
	facets.Availability = availabilityFacets

	// Get rating facets (skip for now since rating column doesn't exist)
	// TODO: Add rating column to products table and implement rating facets
	facets.Ratings = []repositories.RatingFacet{}

	// Get shipping facets
	var shippingFacets []repositories.ShippingFacet
	err = r.db.WithContext(ctx).Raw(`
		SELECT
			shipping_type as type,
			shipping_label as label,
			COUNT(*) as product_count
		FROM (
			SELECT
				CASE
					WHEN requires_shipping = true THEN 'physical'
					WHEN is_digital = true THEN 'digital'
					ELSE 'unknown'
				END as shipping_type,
				CASE
					WHEN requires_shipping = true THEN 'Physical Products'
					WHEN is_digital = true THEN 'Digital Products'
					ELSE 'Unknown'
				END as shipping_label
			FROM products
			WHERE status = 'active'
		) shipping_subquery
		GROUP BY shipping_type, shipping_label
		HAVING COUNT(*) > 0
		ORDER BY product_count DESC
	`).Scan(&shippingFacets).Error
	if err != nil {
		return nil, err
	}
	facets.Shipping = shippingFacets

	return facets, nil
}

// EnhancedSearch performs enhanced search with dynamic faceting
func (r *searchRepository) EnhancedSearch(ctx context.Context, params repositories.EnhancedSearchParams) ([]*entities.Product, int64, *repositories.DynamicSearchFacets, error) {
	// Convert enhanced params to full-text search params
	searchParams := params.FullTextSearchParams

	// Perform the search
	products, total, err := r.FullTextSearch(ctx, searchParams)
	if err != nil {
		return nil, 0, nil, err
	}

	var facets *repositories.DynamicSearchFacets
	if params.IncludeFacets {
		facets, err = r.GetDynamicFacets(ctx, params)
		if err != nil {
			// Log error but don't fail the search
			fmt.Printf("Error getting dynamic facets: %v\n", err)
		}
	}

	return products, total, facets, nil
}

// GetDynamicFacets retrieves dynamic facets with real-time counts based on current filters
func (r *searchRepository) GetDynamicFacets(ctx context.Context, params repositories.EnhancedSearchParams) (*repositories.DynamicSearchFacets, error) {
	facets := &repositories.DynamicSearchFacets{}
	baseQuery := r.buildFacetBaseQuery(params)

	// Get category facets
	categoryFacets, err := r.getDynamicCategoryFacets(ctx, params, baseQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get category facets: %w", err)
	}
	facets.Categories = categoryFacets

	// Get brand facets
	brandFacets, err := r.getDynamicBrandFacets(ctx, params, baseQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get brand facets: %w", err)
	}
	facets.Brands = brandFacets

	// Get tag facets
	tagFacets, err := r.getDynamicTagFacets(ctx, params, baseQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag facets: %w", err)
	}
	facets.Tags = tagFacets

	// Get price range facets
	priceRangeFacets, err := r.getDynamicPriceRangeFacet(ctx, params, baseQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get price range facets: %w", err)
	}
	facets.PriceRange = *priceRangeFacets

	// Get status facets
	statusFacets, err := r.getDynamicStatusFacets(ctx, params, baseQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get status facets: %w", err)
	}
	facets.Status = statusFacets

	// Calculate total count
	var totalCount int64
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}
	facets.TotalCount = totalCount

	return facets, nil
}

// buildFacetBaseQuery builds the base query for facet calculation
func (r *searchRepository) buildFacetBaseQuery(params repositories.EnhancedSearchParams) *gorm.DB {
	// Simple query without complex search conditions to avoid SQL errors
	query := r.db.Model(&entities.Product{}).Where("status = ?", "active")
	return query
}

// buildSearchCondition builds the search condition for facet queries
func (r *searchRepository) buildSearchCondition(searchQuery string) string {
	return `(
		name ILIKE ? OR
		sku ILIKE ? OR
		description ILIKE ?
	)`
}

// applyFacetFilters applies filters to the facet query
func (r *searchRepository) applyFacetFilters(query *gorm.DB, params repositories.EnhancedSearchParams) *gorm.DB {
	// Apply price filters
	if params.MinPrice != nil {
		query = query.Where("price >= ?", *params.MinPrice)
	}
	if params.MaxPrice != nil {
		query = query.Where("price <= ?", *params.MaxPrice)
	}

	// Apply category filters (if not calculating category facets)
	if len(params.CategoryIDs) > 0 {
		query = query.Where("category_id IN ?", params.CategoryIDs)
	}

	// Apply brand filters (if not calculating brand facets)
	if len(params.BrandIDs) > 0 {
		query = query.Where("brand_id IN ?", params.BrandIDs)
	}

	// Apply other filters...
	if params.Featured != nil {
		query = query.Where("featured = ?", *params.Featured)
	}

	if params.InStock != nil {
		if *params.InStock {
			query = query.Where("stock > 0")
		} else {
			query = query.Where("stock = 0")
		}
	}

	if params.OnSale != nil && *params.OnSale {
		query = query.Where("sale_price IS NOT NULL AND sale_price > 0")
	}

	return query
}

// getDynamicCategoryFacets gets category facets with dynamic counts
func (r *searchRepository) getDynamicCategoryFacets(ctx context.Context, params repositories.EnhancedSearchParams, baseQuery *gorm.DB) ([]repositories.DynamicCategoryFacet, error) {
	// Create query excluding category filters for accurate counts
	// facetQuery := r.buildFacetBaseQuery(params)
	// Remove category filter for this facet calculation
	// if len(params.BrandIDs) > 0 {
	//	facetQuery = facetQuery.Where("brand_id IN ?", params.BrandIDs)
	// }

	type CategoryCount struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		ProductCount int64  `json:"product_count"`
	}

	var categoryCounts []CategoryCount
	err := r.db.Table("products p").
		Select("c.id, c.name, COUNT(p.id) as product_count").
		Joins("JOIN categories c ON p.category_id = c.id").
		Where("p.status = ?", "active").
		Group("c.id, c.name").
		Having("COUNT(p.id) > 0").
		Order("product_count DESC").
		Limit(20).
		Scan(&categoryCounts).Error

	if err != nil {
		return nil, err
	}

	var dynamicFacets []repositories.DynamicCategoryFacet
	for _, cat := range categoryCounts {
		categoryID, _ := uuid.Parse(cat.ID)
		isSelected := r.isIDSelected(categoryID.String(), params.CategoryIDs)

		dynamicFacets = append(dynamicFacets, repositories.DynamicCategoryFacet{
			CategoryFacet: repositories.CategoryFacet{
				ID:           categoryID,
				Name:         cat.Name,
				ProductCount: cat.ProductCount,
			},
			IsSelected: isSelected,
			IsDisabled: cat.ProductCount == 0,
		})
	}

	return dynamicFacets, nil
}

// getDynamicBrandFacets gets brand facets with dynamic counts
func (r *searchRepository) getDynamicBrandFacets(ctx context.Context, params repositories.EnhancedSearchParams, baseQuery *gorm.DB) ([]repositories.DynamicBrandFacet, error) {
	// Create query excluding brand filters for accurate counts
	facetQuery := r.buildFacetBaseQuery(params)
	// Remove brand filter for this facet calculation
	if len(params.CategoryIDs) > 0 {
		facetQuery = facetQuery.Where("category_id IN ?", params.CategoryIDs)
	}

	type BrandCount struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		ProductCount int64  `json:"product_count"`
	}

	var brandCounts []BrandCount
	err := r.db.Table("products p").
		Select("b.id, b.name, COUNT(p.id) as product_count").
		Joins("JOIN brands b ON p.brand_id = b.id").
		Where("p.status = ?", "active").
		Group("b.id, b.name").
		Having("COUNT(p.id) > 0").
		Order("product_count DESC").
		Limit(20).
		Scan(&brandCounts).Error

	if err != nil {
		return nil, err
	}

	var dynamicFacets []repositories.DynamicBrandFacet
	for _, brand := range brandCounts {
		brandID, _ := uuid.Parse(brand.ID)
		isSelected := r.isIDSelected(brandID.String(), params.BrandIDs)

		dynamicFacets = append(dynamicFacets, repositories.DynamicBrandFacet{
			BrandFacet: repositories.BrandFacet{
				ID:           brandID,
				Name:         brand.Name,
				ProductCount: brand.ProductCount,
			},
			IsSelected: isSelected,
			IsDisabled: brand.ProductCount == 0,
		})
	}

	return dynamicFacets, nil
}

// getDynamicTagFacets gets tag facets with dynamic counts
func (r *searchRepository) getDynamicTagFacets(ctx context.Context, params repositories.EnhancedSearchParams, baseQuery *gorm.DB) ([]repositories.DynamicTagFacet, error) {
	// facetQuery := r.buildFacetBaseQuery(params)

	type TagCount struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		ProductCount int64  `json:"product_count"`
	}

	var tagCounts []TagCount
	err := r.db.Table("products p").
		Select("t.id, t.name, COUNT(DISTINCT p.id) as product_count").
		Joins("JOIN product_tag_associations pt ON p.id = pt.product_id").
		Joins("JOIN tags t ON pt.product_tag_id = t.id").
		Where("p.status = ?", "active").
		Group("t.id, t.name").
		Having("COUNT(DISTINCT p.id) > 0").
		Order("product_count DESC").
		Limit(20).
		Scan(&tagCounts).Error

	if err != nil {
		return nil, err
	}

	var dynamicFacets []repositories.DynamicTagFacet
	for _, tag := range tagCounts {
		tagID, _ := uuid.Parse(tag.ID)
		isSelected := r.isStringSelected(tag.ID, params.Tags)

		dynamicFacets = append(dynamicFacets, repositories.DynamicTagFacet{
			TagFacet: repositories.TagFacet{
				ID:           tagID,
				Name:         tag.Name,
				ProductCount: tag.ProductCount,
			},
			IsSelected: isSelected,
			IsDisabled: tag.ProductCount == 0,
		})
	}

	return dynamicFacets, nil
}

// getDynamicPriceRangeFacet gets price range facet with dynamic data
func (r *searchRepository) getDynamicPriceRangeFacet(ctx context.Context, params repositories.EnhancedSearchParams, baseQuery *gorm.DB) (*repositories.DynamicPriceRangeFacet, error) {
	facetQuery := r.buildFacetBaseQuery(params)

	type PriceStats struct {
		MinPrice float64 `json:"min_price"`
		MaxPrice float64 `json:"max_price"`
	}

	var priceStats PriceStats
	err := r.db.Table("products").
		Select("MIN(price) as min_price, MAX(price) as max_price").
		Where("status = ?", "active").
		Scan(&priceStats).Error

	if err != nil {
		return nil, err
	}

	// Define price ranges
	ranges := []repositories.PriceRange{
		{Min: nil, Max: &[]float64{50}[0], Label: "Under $50"},
		{Min: &[]float64{50}[0], Max: &[]float64{100}[0], Label: "$50 - $100"},
		{Min: &[]float64{100}[0], Max: &[]float64{250}[0], Label: "$100 - $250"},
		{Min: &[]float64{250}[0], Max: &[]float64{500}[0], Label: "$250 - $500"},
		{Min: &[]float64{500}[0], Max: &[]float64{1000}[0], Label: "$500 - $1000"},
		{Min: &[]float64{1000}[0], Max: nil, Label: "Over $1000"},
	}

	// Calculate counts for each range
	for i := range ranges {
		var count int64
		rangeQuery := facetQuery

		if ranges[i].Min != nil {
			rangeQuery = rangeQuery.Where("current_price >= ?", *ranges[i].Min)
		}
		if ranges[i].Max != nil {
			rangeQuery = rangeQuery.Where("current_price <= ?", *ranges[i].Max)
		}

		rangeQuery.Count(&count)
		ranges[i].ProductCount = count
	}

	return &repositories.DynamicPriceRangeFacet{
		PriceRangeFacet: repositories.PriceRangeFacet{
			MinPrice: priceStats.MinPrice,
			MaxPrice: priceStats.MaxPrice,
			Ranges:   ranges,
		},
		SelectedMin: params.MinPrice,
		SelectedMax: params.MaxPrice,
	}, nil
}

// getDynamicStatusFacets gets status facets with dynamic counts
func (r *searchRepository) getDynamicStatusFacets(ctx context.Context, params repositories.EnhancedSearchParams, baseQuery *gorm.DB) ([]repositories.DynamicStatusFacet, error) {
	// facetQuery := r.buildFacetBaseQuery(params)

	type StatusCount struct {
		Status       string `json:"status"`
		Label        string `json:"label"`
		ProductCount int64  `json:"product_count"`
	}

	var statusCounts []StatusCount
	err := r.db.Table("products").
		Select(`
			status,
			CASE
				WHEN status = 'active' THEN 'Active'
				WHEN status = 'inactive' THEN 'Inactive'
				WHEN status = 'draft' THEN 'Draft'
				ELSE status
			END as label,
			COUNT(*) as product_count
		`).
		Group("status").
		Having("COUNT(*) > 0").
		Order("product_count DESC").
		Scan(&statusCounts).Error

	if err != nil {
		return nil, err
	}

	var dynamicFacets []repositories.DynamicStatusFacet
	for _, status := range statusCounts {
		isSelected := r.isStringSelected(status.Status, []string{}) // Add status selection logic

		dynamicFacets = append(dynamicFacets, repositories.DynamicStatusFacet{
			StatusFacet: repositories.StatusFacet{
				Status:       entities.ProductStatus(status.Status),
				Label:        status.Label,
				ProductCount: status.ProductCount,
			},
			IsSelected: isSelected,
			IsDisabled: status.ProductCount == 0,
		})
	}

	return dynamicFacets, nil
}

// Helper methods
func (r *searchRepository) isIDSelected(id string, selectedIDs []uuid.UUID) bool {
	for _, selectedID := range selectedIDs {
		if selectedID.String() == id {
			return true
		}
	}
	return false
}

func (r *searchRepository) isStringSelected(value string, selectedValues []string) bool {
	for _, selectedValue := range selectedValues {
		if selectedValue == value {
			return true
		}
	}
	return false
}

// GetFacetCounts gets counts for a specific facet type
func (r *searchRepository) GetFacetCounts(ctx context.Context, params repositories.EnhancedSearchParams, facetType string) (map[string]int64, error) {
	counts := make(map[string]int64)

	switch facetType {
	case "categories":
		var results []struct {
			ID    string `json:"id"`
			Count int64  `json:"count"`
		}
		err := r.db.Table("products").
			Select("category_id as id, COUNT(*) as count").
			Where("status = ?", "active").
			Group("category_id").
			Scan(&results).Error
		if err != nil {
			return nil, err
		}
		for _, result := range results {
			counts[result.ID] = result.Count
		}
	case "brands":
		var results []struct {
			ID    string `json:"id"`
			Count int64  `json:"count"`
		}
		err := r.db.Table("products").
			Select("brand_id as id, COUNT(*) as count").
			Where("status = ?", "active").
			Group("brand_id").
			Scan(&results).Error
		if err != nil {
			return nil, err
		}
		for _, result := range results {
			counts[result.ID] = result.Count
		}
	}

	return counts, nil
}

// Helper function to create float64 pointer
func float64Ptr(f float64) *float64 {
	return &f
}
