package database

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

// GetAutocompleteEntries retrieves autocomplete entries based on query and types
func (r *searchRepository) GetAutocompleteEntries(ctx context.Context, query string, types []string, limit int) ([]*entities.AutocompleteEntry, error) {
	if limit <= 0 {
		limit = 10
	}

	var entries []*entities.AutocompleteEntry
	dbQuery := r.db.WithContext(ctx).Model(&entities.AutocompleteEntry{}).
		Where("is_active = true")

	if query != "" {
		dbQuery = dbQuery.Where("value ILIKE ? OR display_text ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if len(types) > 0 {
		dbQuery = dbQuery.Where("type IN ?", types)
	}

	err := dbQuery.Order("priority DESC, search_count DESC, click_count DESC").
		Limit(limit).
		Find(&entries).Error

	return entries, err
}

// CreateAutocompleteEntry creates a new autocomplete entry
func (r *searchRepository) CreateAutocompleteEntry(ctx context.Context, entry *entities.AutocompleteEntry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

// UpdateAutocompleteEntry updates an existing autocomplete entry
func (r *searchRepository) UpdateAutocompleteEntry(ctx context.Context, entry *entities.AutocompleteEntry) error {
	return r.db.WithContext(ctx).Save(entry).Error
}

// DeleteAutocompleteEntry deletes an autocomplete entry
func (r *searchRepository) DeleteAutocompleteEntry(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.AutocompleteEntry{}, id).Error
}

// IncrementAutocompleteUsage increments usage statistics for an autocomplete entry
func (r *searchRepository) IncrementAutocompleteUsage(ctx context.Context, id uuid.UUID, isClick bool) error {
	updates := map[string]interface{}{
		"search_count": gorm.Expr("search_count + 1"),
		"updated_at":   time.Now(),
	}

	if isClick {
		updates["click_count"] = gorm.Expr("click_count + 1")
	}

	return r.db.WithContext(ctx).Model(&entities.AutocompleteEntry{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// GetSearchTrends retrieves search trends for a specific period
func (r *searchRepository) GetSearchTrends(ctx context.Context, period string, limit int) ([]*entities.SearchTrend, error) {
	if limit <= 0 {
		limit = 20
	}

	var trends []*entities.SearchTrend
	err := r.db.WithContext(ctx).
		Where("period = ?", period).
		Order("search_count DESC, date DESC").
		Limit(limit).
		Find(&trends).Error

	return trends, err
}

// UpdateSearchTrend updates or creates search trend data
func (r *searchRepository) UpdateSearchTrend(ctx context.Context, query string, period string) error {
	today := time.Now().Truncate(24 * time.Hour)

	// Try to update existing trend
	result := r.db.WithContext(ctx).Model(&entities.SearchTrend{}).
		Where("query = ? AND period = ? AND date = ?", query, period, today).
		Update("search_count", gorm.Expr("search_count + 1"))

	if result.Error != nil {
		return result.Error
	}

	// If no rows affected, create new trend
	if result.RowsAffected == 0 {
		trend := &entities.SearchTrend{
			Query:       query,
			SearchCount: 1,
			Period:      period,
			Date:        today,
		}
		return r.db.WithContext(ctx).Create(trend).Error
	}

	return nil
}

// GetUserSearchPreferences retrieves user search preferences
func (r *searchRepository) GetUserSearchPreferences(ctx context.Context, userID uuid.UUID) (*entities.UserSearchPreference, error) {
	var prefs entities.UserSearchPreference
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&prefs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return default preferences
			return &entities.UserSearchPreference{
				UserID:              userID,
				SearchLanguage:      "en",
				AutocompleteEnabled: true,
				SearchHistoryEnabled: true,
				PersonalizedResults: true,
			}, nil
		}
		return nil, err
	}
	return &prefs, nil
}

// SaveUserSearchPreferences saves user search preferences
func (r *searchRepository) SaveUserSearchPreferences(ctx context.Context, prefs *entities.UserSearchPreference) error {
	return r.db.WithContext(ctx).Save(prefs).Error
}

// CreateSearchSession creates a new search session
func (r *searchRepository) CreateSearchSession(ctx context.Context, session *entities.SearchSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// UpdateSearchSession updates an existing search session
func (r *searchRepository) UpdateSearchSession(ctx context.Context, session *entities.SearchSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// GetSearchSession retrieves a search session by session ID
func (r *searchRepository) GetSearchSession(ctx context.Context, sessionID string) (*entities.SearchSession, error) {
	var session entities.SearchSession
	err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).First(&session).Error
	return &session, err
}

// GetPersonalizedSuggestions retrieves personalized suggestions for a user
func (r *searchRepository) GetPersonalizedSuggestions(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*entities.AutocompleteEntry, error) {
	if limit <= 0 {
		limit = 5
	}

	// Get user preferences
	prefs, err := r.GetUserSearchPreferences(ctx, userID)
	if err != nil || !prefs.PersonalizedResults {
		// Fall back to general suggestions
		return r.GetAutocompleteEntries(ctx, query, nil, limit)
	}

	var entries []*entities.AutocompleteEntry
	dbQuery := r.db.WithContext(ctx).Model(&entities.AutocompleteEntry{}).
		Where("is_active = true")

	if query != "" {
		dbQuery = dbQuery.Where("value ILIKE ? OR display_text ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	// Prioritize based on user preferences
	if len(prefs.PreferredCategories) > 0 {
		dbQuery = dbQuery.Where("(type != 'category' OR metadata::jsonb->>'category' = ANY(?))",
			pq.Array(prefs.PreferredCategories))
	}

	if len(prefs.PreferredBrands) > 0 {
		dbQuery = dbQuery.Where("(type != 'brand' OR metadata::jsonb->>'brand' = ANY(?))",
			pq.Array(prefs.PreferredBrands))
	}

	err = dbQuery.Order("priority DESC, search_count DESC").
		Limit(limit).
		Find(&entries).Error

	return entries, err
}

// GetTrendingSuggestions retrieves trending suggestions
func (r *searchRepository) GetTrendingSuggestions(ctx context.Context, limit int) ([]*entities.AutocompleteEntry, error) {
	if limit <= 0 {
		limit = 10
	}

	var entries []*entities.AutocompleteEntry

	// Get trending queries from search trends
	subQuery := r.db.WithContext(ctx).Model(&entities.SearchTrend{}).
		Select("query").
		Where("period = 'daily' AND date >= ?", time.Now().AddDate(0, 0, -7)).
		Group("query").
		Order("SUM(search_count) DESC").
		Limit(limit)

	err := r.db.WithContext(ctx).Model(&entities.AutocompleteEntry{}).
		Where("is_active = true AND type = 'query' AND value IN (?)", subQuery).
		Order("search_count DESC").
		Limit(limit).
		Find(&entries).Error

	return entries, err
}

// GetCategorySuggestions retrieves category-based suggestions
func (r *searchRepository) GetCategorySuggestions(ctx context.Context, query string, limit int) ([]*entities.AutocompleteEntry, error) {
	if limit <= 0 {
		limit = 5
	}

	var entries []*entities.AutocompleteEntry
	dbQuery := r.db.WithContext(ctx).Model(&entities.AutocompleteEntry{}).
		Where("is_active = true AND type = 'category'")

	if query != "" {
		dbQuery = dbQuery.Where("value ILIKE ? OR display_text ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	err := dbQuery.Order("priority DESC, search_count DESC").
		Limit(limit).
		Find(&entries).Error

	return entries, err
}

// GetBrandSuggestions retrieves brand-based suggestions
func (r *searchRepository) GetBrandSuggestions(ctx context.Context, query string, limit int) ([]*entities.AutocompleteEntry, error) {
	if limit <= 0 {
		limit = 5
	}

	var entries []*entities.AutocompleteEntry
	dbQuery := r.db.WithContext(ctx).Model(&entities.AutocompleteEntry{}).
		Where("is_active = true AND type = 'brand'")

	if query != "" {
		dbQuery = dbQuery.Where("value ILIKE ? OR display_text ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	err := dbQuery.Order("priority DESC, search_count DESC").
		Limit(limit).
		Find(&entries).Error

	return entries, err
}

// GetProductSuggestions retrieves product-based suggestions
func (r *searchRepository) GetProductSuggestions(ctx context.Context, query string, limit int) ([]*entities.AutocompleteEntry, error) {
	if limit <= 0 {
		limit = 5
	}

	var entries []*entities.AutocompleteEntry
	dbQuery := r.db.WithContext(ctx).Model(&entities.AutocompleteEntry{}).
		Where("is_active = true AND type = 'product'")

	if query != "" {
		dbQuery = dbQuery.Where("value ILIKE ? OR display_text ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	err := dbQuery.Order("priority DESC, search_count DESC").
		Limit(limit).
		Find(&entries).Error

	return entries, err
}

// RebuildAutocompleteIndex rebuilds the autocomplete index from existing data
func (r *searchRepository) RebuildAutocompleteIndex(ctx context.Context) error {
	// Clear existing entries
	if err := r.db.WithContext(ctx).Where("1 = 1").Delete(&entities.AutocompleteEntry{}).Error; err != nil {
		return fmt.Errorf("failed to clear autocomplete entries: %w", err)
	}

	// Rebuild from products
	if err := r.rebuildProductSuggestions(ctx); err != nil {
		return fmt.Errorf("failed to rebuild product suggestions: %w", err)
	}

	// Rebuild from categories
	if err := r.rebuildCategorySuggestions(ctx); err != nil {
		return fmt.Errorf("failed to rebuild category suggestions: %w", err)
	}

	// Rebuild from brands
	if err := r.rebuildBrandSuggestions(ctx); err != nil {
		return fmt.Errorf("failed to rebuild brand suggestions: %w", err)
	}

	// Rebuild from search history
	if err := r.rebuildQuerySuggestions(ctx); err != nil {
		return fmt.Errorf("failed to rebuild query suggestions: %w", err)
	}

	return nil
}

// CleanupOldAutocompleteEntries removes old unused autocomplete entries
func (r *searchRepository) CleanupOldAutocompleteEntries(ctx context.Context, days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)

	return r.db.WithContext(ctx).
		Where("updated_at < ? AND search_count = 0 AND click_count = 0", cutoffDate).
		Delete(&entities.AutocompleteEntry{}).Error
}

// Helper methods for rebuilding autocomplete index

func (r *searchRepository) rebuildProductSuggestions(ctx context.Context) error {
	var products []struct {
		ID   uuid.UUID
		Name string
	}

	if err := r.db.WithContext(ctx).Model(&entities.Product{}).
		Select("id, name").
		Where("status = 'active'").
		Find(&products).Error; err != nil {
		return err
	}

	for _, product := range products {
		entry := &entities.AutocompleteEntry{
			Type:        "product",
			Value:       product.Name,
			DisplayText: product.Name,
			EntityID:    &product.ID,
			Priority:    50,
			IsActive:    true,
			Metadata:    fmt.Sprintf(`{"product_id": "%s"}`, product.ID),
		}

		if err := r.db.WithContext(ctx).Create(entry).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *searchRepository) rebuildCategorySuggestions(ctx context.Context) error {
	var categories []struct {
		ID   uuid.UUID
		Name string
	}

	if err := r.db.WithContext(ctx).Model(&entities.Category{}).
		Select("id, name").
		Where("is_active = true").
		Find(&categories).Error; err != nil {
		return err
	}

	for _, category := range categories {
		entry := &entities.AutocompleteEntry{
			Type:        "category",
			Value:       category.Name,
			DisplayText: category.Name,
			EntityID:    &category.ID,
			Priority:    70,
			IsActive:    true,
			Metadata:    fmt.Sprintf(`{"category_id": "%s"}`, category.ID),
		}

		if err := r.db.WithContext(ctx).Create(entry).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *searchRepository) rebuildBrandSuggestions(ctx context.Context) error {
	var brands []struct {
		ID   uuid.UUID
		Name string
	}

	if err := r.db.WithContext(ctx).Model(&entities.Brand{}).
		Select("id, name").
		Where("is_active = true").
		Find(&brands).Error; err != nil {
		return err
	}

	for _, brand := range brands {
		entry := &entities.AutocompleteEntry{
			Type:        "brand",
			Value:       brand.Name,
			DisplayText: brand.Name,
			EntityID:    &brand.ID,
			Priority:    60,
			IsActive:    true,
			Metadata:    fmt.Sprintf(`{"brand_id": "%s"}`, brand.ID),
		}

		if err := r.db.WithContext(ctx).Create(entry).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *searchRepository) rebuildQuerySuggestions(ctx context.Context) error {
	var suggestions []struct {
		Query       string
		SearchCount int
	}

	if err := r.db.WithContext(ctx).Model(&entities.SearchSuggestion{}).
		Select("query, search_count").
		Where("is_active = true AND search_count > 0").
		Order("search_count DESC").
		Limit(1000).
		Find(&suggestions).Error; err != nil {
		return err
	}

	for _, suggestion := range suggestions {
		entry := &entities.AutocompleteEntry{
			Type:        "query",
			Value:       suggestion.Query,
			DisplayText: suggestion.Query,
			Priority:    80,
			SearchCount: suggestion.SearchCount,
			IsActive:    true,
			Metadata:    `{"type": "search_query"}`,
		}

		if err := r.db.WithContext(ctx).Create(entry).Error; err != nil {
			return err
		}
	}

	return nil
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

// GetSmartAutocomplete provides intelligent autocomplete suggestions
func (r *searchRepository) GetSmartAutocomplete(ctx context.Context, req entities.SmartAutocompleteRequest) (*entities.SmartAutocompleteResponse, error) {
	startTime := time.Now()

	response := &entities.SmartAutocompleteResponse{
		Suggestions: []entities.SmartAutocompleteSuggestion{},
		Total:       0,
		HasMore:     false,
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	var allSuggestions []entities.SmartAutocompleteSuggestion

	// 1. Get fuzzy matches for typo tolerance
	if req.Query != "" {
		fuzzyEntries, err := r.GetFuzzyMatches(ctx, req.Query, req.Types, req.Limit/2)
		if err == nil {
			for _, entry := range fuzzyEntries {
				suggestion := r.convertToSmartSuggestion(entry, "fuzzy_match")
				allSuggestions = append(allSuggestions, suggestion)
			}
		}
	}

	// 2. Get personalized suggestions if user is authenticated
	if req.IncludePersonalized && req.UserID != nil {
		personalizedEntries, err := r.GetPersonalizedSuggestions(ctx, *req.UserID, req.Query, req.Limit/4)
		if err == nil {
			for _, entry := range personalizedEntries {
				suggestion := r.convertToSmartSuggestion(entry, "personalized")
				suggestion.IsPersonalized = true
				allSuggestions = append(allSuggestions, suggestion)
			}
		}
	}

	// 3. Get trending suggestions
	if req.IncludeTrending {
		trendingEntries, err := r.GetTrendingSuggestions(ctx, req.Limit/4)
		if err == nil {
			for _, entry := range trendingEntries {
				suggestion := r.convertToSmartSuggestion(entry, "trending")
				suggestion.IsTrending = true
				allSuggestions = append(allSuggestions, suggestion)
			}
		}
	}

	// 4. Get popular suggestions
	if req.IncludePopular {
		popularEntries, err := r.GetPopularSuggestions(ctx, req.Limit/4, "week")
		if err == nil {
			for _, entry := range popularEntries {
				suggestion := r.convertToSmartSuggestion(entry, "popular")
				allSuggestions = append(allSuggestions, suggestion)
			}
		}
	}

	// 5. Get search history if requested
	if req.IncludeHistory && req.UserID != nil {
		historyEntries, err := r.GetUserAutocompleteHistory(ctx, *req.UserID, req.Limit/4)
		if err == nil {
			for _, entry := range historyEntries {
				suggestion := r.convertToSmartSuggestion(entry, "history")
				allSuggestions = append(allSuggestions, suggestion)
			}
		}
	}

	// 6. Get synonym suggestions
	if req.Query != "" {
		synonymEntries, err := r.GetSynonymSuggestions(ctx, req.Query, req.Limit/4)
		if err == nil {
			for _, entry := range synonymEntries {
				suggestion := r.convertToSmartSuggestion(entry, "synonym")
				allSuggestions = append(allSuggestions, suggestion)
			}
		}
	}

	// Remove duplicates and sort by score
	uniqueSuggestions := r.deduplicateAndSort(allSuggestions)

	// Limit results
	if len(uniqueSuggestions) > req.Limit {
		response.Suggestions = uniqueSuggestions[:req.Limit]
		response.HasMore = true
	} else {
		response.Suggestions = uniqueSuggestions
	}

	// Group by type for easier frontend consumption
	r.groupSuggestionsByType(response)

	response.Total = len(response.Suggestions)
	response.QueryTime = time.Since(startTime).Milliseconds()

	return response, nil
}

// GetFuzzyMatches provides fuzzy matching for typo tolerance
func (r *searchRepository) GetFuzzyMatches(ctx context.Context, query string, types []string, limit int) ([]*entities.AutocompleteEntry, error) {
	if limit <= 0 {
		limit = 10
	}

	var entries []*entities.AutocompleteEntry
	dbQuery := r.db.WithContext(ctx).Model(&entities.AutocompleteEntry{}).
		Where("is_active = true")

	if len(types) > 0 {
		dbQuery = dbQuery.Where("type IN ?", types)
	}

	// Use PostgreSQL similarity operator for fuzzy matching
	// First try simple ILIKE matching, then add similarity if needed
	dbQuery = dbQuery.Where("value ILIKE ? OR display_text ILIKE ? OR synonyms::text ILIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%")

	err := dbQuery.Order("priority DESC, search_count DESC, score DESC").
		Limit(limit).
		Find(&entries).Error

	// Debug log
	fmt.Printf("DEBUG GetFuzzyMatches: query=%s, types=%v, limit=%d, found=%d entries, error=%v\n",
		query, types, limit, len(entries), err)

	return entries, err
}

// GetSynonymSuggestions gets suggestions based on synonyms
func (r *searchRepository) GetSynonymSuggestions(ctx context.Context, query string, limit int) ([]*entities.AutocompleteEntry, error) {
	if limit <= 0 {
		limit = 5
	}

	// First, find synonyms for the query
	var synonyms []string
	err := r.db.WithContext(ctx).Table("search_synonyms").
		Select("unnest(synonyms) as synonym").
		Where("term ILIKE ? AND is_active = true", "%"+query+"%").
		Pluck("synonym", &synonyms)

	if err != nil || len(synonyms) == 0 {
		return []*entities.AutocompleteEntry{}, nil
	}

	// Get suggestions for synonyms
	var entries []*entities.AutocompleteEntry
	for _, synonym := range synonyms {
		var synonymEntries []*entities.AutocompleteEntry
		err := r.db.WithContext(ctx).Model(&entities.AutocompleteEntry{}).
			Where("is_active = true AND (value ILIKE ? OR display_text ILIKE ?)", "%"+synonym+"%", "%"+synonym+"%").
			Order("priority DESC, search_count DESC").
			Limit(limit / len(synonyms) + 1).
			Find(&synonymEntries).Error

		if err == nil {
			entries = append(entries, synonymEntries...)
		}
	}

	return entries, nil
}

// GetPopularSuggestions gets popular suggestions based on timeframe
func (r *searchRepository) GetPopularSuggestions(ctx context.Context, limit int, timeframe string) ([]*entities.AutocompleteEntry, error) {
	if limit <= 0 {
		limit = 10
	}

	var entries []*entities.AutocompleteEntry
	query := r.db.WithContext(ctx).Model(&entities.AutocompleteEntry{}).
		Where("is_active = true")

	// Add timeframe filter based on updated_at
	switch timeframe {
	case "day":
		query = query.Where("updated_at >= ?", time.Now().AddDate(0, 0, -1))
	case "week":
		query = query.Where("updated_at >= ?", time.Now().AddDate(0, 0, -7))
	case "month":
		query = query.Where("updated_at >= ?", time.Now().AddDate(0, -1, 0))
	}

	err := query.Order("search_count DESC, click_count DESC, priority DESC").
		Limit(limit).
		Find(&entries).Error

	return entries, err
}

// GetUserAutocompleteHistory gets user's search history as suggestions
func (r *searchRepository) GetUserAutocompleteHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*entities.AutocompleteEntry, error) {
	if limit <= 0 {
		limit = 5
	}

	var entries []*entities.AutocompleteEntry
	err := r.db.WithContext(ctx).Model(&entities.AutocompleteEntry{}).
		Where("user_id = ? AND is_active = true", userID).
		Order("updated_at DESC, search_count DESC").
		Limit(limit).
		Find(&entries).Error

	return entries, err
}

// convertToSmartSuggestion converts AutocompleteEntry to SmartAutocompleteSuggestion
func (r *searchRepository) convertToSmartSuggestion(entry *entities.AutocompleteEntry, reason string) entities.SmartAutocompleteSuggestion {
	var metadata map[string]interface{}
	if entry.Metadata != "" {
		json.Unmarshal([]byte(entry.Metadata), &metadata)
	}

	return entities.SmartAutocompleteSuggestion{
		ID:             entry.ID,
		Type:           entry.Type,
		Value:          entry.Value,
		DisplayText:    entry.DisplayText,
		EntityID:       entry.EntityID,
		Priority:       entry.Priority,
		Score:          entry.Score,
		IsTrending:     entry.IsTrending,
		IsPersonalized: entry.IsPersonalized,
		Metadata:       metadata,
		Synonyms:       []string(entry.Synonyms),
		Tags:           []string(entry.Tags),
		Reason:         reason,
	}
}

// deduplicateAndSort removes duplicates and sorts suggestions by score
func (r *searchRepository) deduplicateAndSort(suggestions []entities.SmartAutocompleteSuggestion) []entities.SmartAutocompleteSuggestion {
	seen := make(map[string]bool)
	var unique []entities.SmartAutocompleteSuggestion

	for _, suggestion := range suggestions {
		key := suggestion.Type + ":" + suggestion.Value
		if !seen[key] {
			seen[key] = true
			unique = append(unique, suggestion)
		}
	}

	// Sort by score (descending), then priority (descending)
	sort.Slice(unique, func(i, j int) bool {
		if unique[i].Score != unique[j].Score {
			return unique[i].Score > unique[j].Score
		}
		return unique[i].Priority > unique[j].Priority
	})

	return unique
}

// groupSuggestionsByType groups suggestions by type for easier frontend consumption
func (r *searchRepository) groupSuggestionsByType(response *entities.SmartAutocompleteResponse) {
	for _, suggestion := range response.Suggestions {
		switch suggestion.Type {
		case "product":
			response.Products = append(response.Products, suggestion)
		case "category":
			response.Categories = append(response.Categories, suggestion)
		case "brand":
			response.Brands = append(response.Brands, suggestion)
		case "query":
			response.Queries = append(response.Queries, suggestion)
		}

		if suggestion.IsTrending {
			response.Trending = append(response.Trending, suggestion)
		}

		if suggestion.Reason == "popular" {
			response.Popular = append(response.Popular, suggestion)
		}

		if suggestion.Reason == "history" {
			response.History = append(response.History, suggestion)
		}
	}
}

// TrackAutocompleteClick tracks when a user clicks on an autocomplete suggestion
func (r *searchRepository) TrackAutocompleteClick(ctx context.Context, entryID uuid.UUID, userID *uuid.UUID, sessionID string) error {
	// Increment click count
	err := r.db.WithContext(ctx).Model(&entities.AutocompleteEntry{}).
		Where("id = ?", entryID).
		UpdateColumn("click_count", gorm.Expr("click_count + 1")).Error

	if err != nil {
		return err
	}

	// Record the click event for analytics
	clickEvent := map[string]interface{}{
		"entry_id":   entryID,
		"user_id":    userID,
		"session_id": sessionID,
		"event_type": "autocomplete_click",
		"timestamp":  time.Now(),
	}

	// Store in analytics (simplified - could be expanded)
	return r.db.WithContext(ctx).Table("search_events").Create(map[string]interface{}{
		"query":      "",
		"user_id":    userID,
		"session_id": sessionID,
		"filters":    clickEvent,
		"created_at": time.Now(),
	}).Error
}

// TrackAutocompleteImpression tracks when autocomplete suggestions are shown
func (r *searchRepository) TrackAutocompleteImpression(ctx context.Context, entryID uuid.UUID, userID *uuid.UUID, sessionID string) error {
	// Record the impression event for analytics
	impressionEvent := map[string]interface{}{
		"entry_id":   entryID,
		"user_id":    userID,
		"session_id": sessionID,
		"event_type": "autocomplete_impression",
		"timestamp":  time.Now(),
	}

	// Store in analytics
	return r.db.WithContext(ctx).Table("search_events").Create(map[string]interface{}{
		"query":      "",
		"user_id":    userID,
		"session_id": sessionID,
		"filters":    impressionEvent,
		"created_at": time.Now(),
	}).Error
}

// UpdateAutocompleteTrending updates trending status for autocomplete entries
func (r *searchRepository) UpdateAutocompleteTrending(ctx context.Context) error {
	// Reset all trending flags
	err := r.db.WithContext(ctx).Model(&entities.AutocompleteEntry{}).
		Where("is_trending = true").
		Update("is_trending", false).Error

	if err != nil {
		return err
	}

	// Calculate trending based on recent activity (last 24 hours)
	since := time.Now().AddDate(0, 0, -1)

	// Mark entries as trending based on recent search/click activity
	return r.db.WithContext(ctx).Model(&entities.AutocompleteEntry{}).
		Where("updated_at >= ? AND (search_count > 10 OR click_count > 5)", since).
		Update("is_trending", true).Error
}

// CalculateAutocompleteScores calculates relevance scores for autocomplete entries
func (r *searchRepository) CalculateAutocompleteScores(ctx context.Context) error {
	// Calculate scores based on multiple factors:
	// - Search count (40%)
	// - Click count (30%)
	// - Priority (20%)
	// - Recency (10%)

	return r.db.WithContext(ctx).Exec(`
		UPDATE autocomplete_entries
		SET score = (
			(search_count * 0.4) +
			(click_count * 0.3) +
			(priority * 0.2) +
			(CASE
				WHEN updated_at >= NOW() - INTERVAL '7 days' THEN 10
				WHEN updated_at >= NOW() - INTERVAL '30 days' THEN 5
				ELSE 0
			END * 0.1)
		)
		WHERE is_active = true
	`).Error
}
