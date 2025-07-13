package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

// SearchUseCase defines the interface for search use cases
type SearchUseCase interface {
	// Full-text search
	FullTextSearch(ctx context.Context, req FullTextSearchRequest) (*SearchResponse, error)
	GetSearchSuggestions(ctx context.Context, query string, limit int) ([]string, error)
	GetSearchFacets(ctx context.Context, query string) (*SearchFacetsResponse, error)

	// Enhanced search with dynamic faceting
	EnhancedSearch(ctx context.Context, req *EnhancedSearchRequest) (*EnhancedSearchResponse, error)

	// Search events
	RecordSearchEvent(ctx context.Context, req RecordSearchEventRequest) error
	GetPopularSearchTerms(ctx context.Context, limit int, period string) ([]PopularSearchResponse, error)

	// Search analytics
	GetSearchAnalytics(ctx context.Context, req SearchAnalyticsRequest) (*SearchAnalyticsResponse, error)

	// Search history
	SaveSearchHistory(ctx context.Context, userID uuid.UUID, req SaveSearchHistoryRequest) error
	GetUserSearchHistory(ctx context.Context, userID uuid.UUID, limit int) ([]UserSearchHistoryResponse, error)
	ClearUserSearchHistory(ctx context.Context, userID uuid.UUID) error

	// Search filters
	SaveSearchFilter(ctx context.Context, userID uuid.UUID, req SaveSearchFilterRequest) (*SearchFilterResponse, error)
	GetUserSearchFilters(ctx context.Context, userID uuid.UUID) ([]SearchFilterResponse, error)
	UpdateSearchFilter(ctx context.Context, userID uuid.UUID, filterID uuid.UUID, req UpdateSearchFilterRequest) error
	DeleteSearchFilter(ctx context.Context, userID uuid.UUID, filterID uuid.UUID) error

	// Autocomplete
	GetAutocomplete(ctx context.Context, query string, limit int) (*AutocompleteResponse, error)

	// Enhanced Autocomplete
	GetEnhancedAutocomplete(ctx context.Context, req EnhancedAutocompleteRequest) (*EnhancedAutocompleteResponse, error)
	GetPersonalizedAutocomplete(ctx context.Context, userID uuid.UUID, query string, limit int) (*EnhancedAutocompleteResponse, error)
	GetTrendingSearches(ctx context.Context, limit int) ([]TrendingSearchResponse, error)

	// Search Preferences
	GetUserSearchPreferences(ctx context.Context, userID uuid.UUID) (*UserSearchPreferencesResponse, error)
	UpdateUserSearchPreferences(ctx context.Context, userID uuid.UUID, req UpdateSearchPreferencesRequest) error

	// Search Analytics
	RecordAutocompleteClick(ctx context.Context, req AutocompleteClickRequest) error
	GetSearchTrends(ctx context.Context, period string, limit int) ([]SearchTrendResponse, error)

	// Admin Functions
	RebuildAutocompleteIndex(ctx context.Context) error
	CleanupSearchData(ctx context.Context, days int) error
}

type searchUseCase struct {
	searchRepo  repositories.SearchRepository
	productRepo repositories.ProductRepository
}

// NewSearchUseCase creates a new search use case
func NewSearchUseCase(searchRepo repositories.SearchRepository, productRepo repositories.ProductRepository) SearchUseCase {
	return &searchUseCase{
		searchRepo:  searchRepo,
		productRepo: productRepo,
	}
}

// FullTextSearchRequest represents a full-text search request
type FullTextSearchRequest struct {
	Query       string                  `json:"query"`
	CategoryIDs []uuid.UUID             `json:"category_ids"`
	BrandIDs    []uuid.UUID             `json:"brand_ids"`
	MinPrice    *float64                `json:"min_price"`
	MaxPrice    *float64                `json:"max_price"`
	InStock     *bool                   `json:"in_stock"`
	Featured    *bool                   `json:"featured"`
	OnSale      *bool                   `json:"on_sale"`
	Tags        []string                `json:"tags"`
	Attributes  map[uuid.UUID][]string  `json:"attributes"`

	// Advanced filters
	MinRating           *float64                    `json:"min_rating"`
	MaxRating           *float64                    `json:"max_rating"`
	Visibility          *entities.ProductVisibility `json:"visibility"`
	ProductType         *entities.ProductType       `json:"product_type"`
	Status              *entities.ProductStatus     `json:"status"`
	AvailabilityStatus  *string                     `json:"availability_status"` // in_stock, out_of_stock, low_stock
	CreatedAfter        *time.Time                  `json:"created_after"`
	CreatedBefore       *time.Time                  `json:"created_before"`
	UpdatedAfter        *time.Time                  `json:"updated_after"`
	UpdatedBefore       *time.Time                  `json:"updated_before"`
	MinWeight           *float64                    `json:"min_weight"`
	MaxWeight           *float64                    `json:"max_weight"`
	ShippingClass       *string                     `json:"shipping_class"`
	TaxClass            *string                     `json:"tax_class"`
	MinDiscountPercent  *float64                    `json:"min_discount_percent"`
	MaxDiscountPercent  *float64                    `json:"max_discount_percent"`
	IsDigital           *bool                       `json:"is_digital"`
	RequiresShipping    *bool                       `json:"requires_shipping"`
	AllowBackorder      *bool                       `json:"allow_backorder"`
	TrackQuantity       *bool                       `json:"track_quantity"`

	SortBy      string                  `json:"sort_by"`
	SortOrder   string                  `json:"sort_order"`
	Page        int                     `json:"page"`
	Limit       int                     `json:"limit"`
	UserID      *uuid.UUID              `json:"user_id"`
	SessionID   string                  `json:"session_id"`
	IPAddress   string                  `json:"ip_address"`
	UserAgent   string                  `json:"user_agent"`
}

// SearchResponse represents a search response
type SearchResponse struct {
	Products     []*ProductResponse      `json:"products"`
	Total        int64                   `json:"total"`
	Page         int                     `json:"page"`
	Limit        int                     `json:"limit"`
	TotalPages   int                     `json:"total_pages"`
	HasNext      bool                    `json:"has_next"`
	HasPrev      bool                    `json:"has_prev"`
	Facets       *SearchFacetsResponse   `json:"facets"`
	Query        string                  `json:"query"`
	SearchTime   string                  `json:"search_time"`
	Suggestions  []string                `json:"suggestions"`
}

// SearchFacetsResponse represents search facets response
type SearchFacetsResponse struct {
	Categories   []CategoryFacetResponse   `json:"categories"`
	Brands       []BrandFacetResponse      `json:"brands"`
	PriceRange   PriceRangeFacetResponse   `json:"price_range"`
	Tags         []TagFacetResponse        `json:"tags"`
	Status       []StatusFacetResponse     `json:"status"`
	ProductTypes []ProductTypeFacetResponse `json:"product_types"`
	Availability []AvailabilityFacetResponse `json:"availability"`
	Ratings      []RatingFacetResponse     `json:"ratings"`
	Shipping     []ShippingFacetResponse   `json:"shipping"`
}

// CategoryFacetResponse represents category facet response
type CategoryFacetResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ProductCount int64     `json:"product_count"`
}

// BrandFacetResponse represents brand facet response
type BrandFacetResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ProductCount int64     `json:"product_count"`
}

// PriceRangeFacetResponse represents price range facet response
type PriceRangeFacetResponse struct {
	MinPrice float64                `json:"min_price"`
	MaxPrice float64                `json:"max_price"`
	Ranges   []PriceRangeResponse   `json:"ranges"`
}

// PriceRangeResponse represents a price range response
type PriceRangeResponse struct {
	Min          *float64 `json:"min"`
	Max          *float64 `json:"max"`
	Label        string   `json:"label"`
	ProductCount int64    `json:"product_count"`
}

// TagFacetResponse represents tag facet response
type TagFacetResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ProductCount int64     `json:"product_count"`
}

// StatusFacetResponse represents status facet response
type StatusFacetResponse struct {
	Status       entities.ProductStatus `json:"status"`
	Label        string                 `json:"label"`
	ProductCount int64                  `json:"product_count"`
}

// ProductTypeFacetResponse represents product type facet response
type ProductTypeFacetResponse struct {
	Type         entities.ProductType `json:"type"`
	Label        string               `json:"label"`
	ProductCount int64                `json:"product_count"`
}

// AvailabilityFacetResponse represents availability facet response
type AvailabilityFacetResponse struct {
	Status       string `json:"status"`
	Label        string `json:"label"`
	ProductCount int64  `json:"product_count"`
}

// RatingFacetResponse represents rating facet response
type RatingFacetResponse struct {
	Rating       int    `json:"rating"`
	Label        string `json:"label"`
	ProductCount int64  `json:"product_count"`
}

// ShippingFacetResponse represents shipping facet response
type ShippingFacetResponse struct {
	Type         string `json:"type"`
	Label        string `json:"label"`
	ProductCount int64  `json:"product_count"`
}

// RecordSearchEventRequest represents a record search event request
type RecordSearchEventRequest struct {
	Query            string     `json:"query"`
	UserID           *uuid.UUID `json:"user_id"`
	ResultsCount     int        `json:"results_count"`
	ClickedProductID *uuid.UUID `json:"clicked_product_id"`
	SessionID        string     `json:"session_id"`
	IPAddress        string     `json:"ip_address"`
	UserAgent        string     `json:"user_agent"`
}

// PopularSearchResponse represents popular search response
type PopularSearchResponse struct {
	Query       string `json:"query"`
	SearchCount int    `json:"search_count"`
}

// SearchAnalyticsRequest represents request for search analytics
type SearchAnalyticsRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Limit     int       `json:"limit"`
}

// SearchAnalyticsResponse represents search analytics response
type SearchAnalyticsResponse struct {
	Analytics []SearchAnalyticItem `json:"analytics"`
	Summary   SearchAnalyticsSummary `json:"summary"`
}

type SearchAnalyticItem struct {
	Query              string  `json:"query"`
	TotalSearches      int     `json:"total_searches"`
	AvgResultCount     float64 `json:"avg_result_count"`
	AvgClickThroughRate float64 `json:"avg_click_through_rate"`
	AvgConversionRate  float64 `json:"avg_conversion_rate"`
	LastSearched       string  `json:"last_searched"`
}

type SearchAnalyticsSummary struct {
	TotalQueries       int     `json:"total_queries"`
	TotalSearches      int     `json:"total_searches"`
	AvgResultsPerQuery float64 `json:"avg_results_per_query"`
	TopQuery           string  `json:"top_query"`
}

// SaveSearchHistoryRequest represents save search history request
type SaveSearchHistoryRequest struct {
	Query   string `json:"query"`
	Filters string `json:"filters"`
}

// UserSearchHistoryResponse represents user search history response
type UserSearchHistoryResponse struct {
	ID        uuid.UUID `json:"id"`
	Query     string    `json:"query"`
	Filters   string    `json:"filters"`
	CreatedAt time.Time `json:"created_at"`
}

// SaveSearchFilterRequest represents save search filter request
type SaveSearchFilterRequest struct {
	Name      string `json:"name"`
	Query     string `json:"query"`
	Filters   string `json:"filters"`
	IsDefault bool   `json:"is_default"`
	IsPublic  bool   `json:"is_public"`
}

// UpdateSearchFilterRequest represents update search filter request
type UpdateSearchFilterRequest struct {
	Name      string `json:"name"`
	Query     string `json:"query"`
	Filters   string `json:"filters"`
	IsDefault bool   `json:"is_default"`
	IsPublic  bool   `json:"is_public"`
}

// SearchFilterResponse represents search filter response
type SearchFilterResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Query      string    `json:"query"`
	Filters    string    `json:"filters"`
	IsDefault  bool      `json:"is_default"`
	IsPublic   bool      `json:"is_public"`
	UsageCount int       `json:"usage_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AutocompleteResponse represents autocomplete response
type AutocompleteResponse struct {
	Products    []ProductSuggestionResponse    `json:"products"`
	Categories  []CategorySuggestionResponse   `json:"categories"`
	Brands      []BrandSuggestionResponse      `json:"brands"`
	Suggestions []string                       `json:"suggestions"`
}

// ProductSuggestionResponse represents product suggestion response
type ProductSuggestionResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
	Image string    `json:"image"`
}

// CategorySuggestionResponse represents category suggestion response
type CategorySuggestionResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// BrandSuggestionResponse represents brand suggestion response
type BrandSuggestionResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// FullTextSearch performs full-text search
func (uc *searchUseCase) FullTextSearch(ctx context.Context, req FullTextSearchRequest) (*SearchResponse, error) {
	startTime := time.Now()
	
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.SortBy == "" {
		req.SortBy = "relevance"
	}
	
	// Calculate offset
	offset := (req.Page - 1) * req.Limit
	
	// Build search parameters
	params := repositories.FullTextSearchParams{
		Query:       req.Query,
		CategoryIDs: req.CategoryIDs,
		BrandIDs:    req.BrandIDs,
		MinPrice:    req.MinPrice,
		MaxPrice:    req.MaxPrice,
		InStock:     req.InStock,
		Featured:    req.Featured,
		OnSale:      req.OnSale,
		Tags:        req.Tags,
		Attributes:  req.Attributes,

		// Advanced filters
		MinRating:           req.MinRating,
		MaxRating:           req.MaxRating,
		Visibility:          req.Visibility,
		ProductType:         req.ProductType,
		Status:              req.Status,
		AvailabilityStatus:  req.AvailabilityStatus,
		CreatedAfter:        req.CreatedAfter,
		CreatedBefore:       req.CreatedBefore,
		UpdatedAfter:        req.UpdatedAfter,
		UpdatedBefore:       req.UpdatedBefore,
		MinWeight:           req.MinWeight,
		MaxWeight:           req.MaxWeight,
		ShippingClass:       req.ShippingClass,
		TaxClass:            req.TaxClass,
		MinDiscountPercent:  req.MinDiscountPercent,
		MaxDiscountPercent:  req.MaxDiscountPercent,
		IsDigital:           req.IsDigital,
		RequiresShipping:    req.RequiresShipping,
		AllowBackorder:      req.AllowBackorder,
		TrackQuantity:       req.TrackQuantity,

		SortBy:      req.SortBy,
		SortOrder:   req.SortOrder,
		Limit:       req.Limit,
		Offset:      offset,
	}
	
	// Perform search
	products, total, err := uc.searchRepo.FullTextSearch(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to perform full-text search: %w", err)
	}
	
	// Convert to response format
	productResponses := make([]*ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = uc.toProductResponse(product)
	}
	
	// Calculate pagination
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))
	
	// Get search facets
	facets, err := uc.searchRepo.GetSearchFacets(ctx, req.Query)
	var facetsResponse *SearchFacetsResponse
	if err != nil {
		// Log error but don't fail the search
		fmt.Printf("Error getting search facets: %v\n", err)
	} else if facets != nil {
		facetsResponse = uc.toSearchFacetsResponse(facets)
	}
	
	// Get suggestions if query is provided
	var suggestions []string
	if req.Query != "" {
		suggestionEntities, _ := uc.searchRepo.GetSearchSuggestions(ctx, req.Query, 5)
		for _, s := range suggestionEntities {
			suggestions = append(suggestions, s.Query)
		}
	}
	
	// Record search event and analytics
	if req.Query != "" {
		event := &entities.SearchEvent{
			Query:        req.Query,
			UserID:       req.UserID,
			ResultsCount: int(total),
			SessionID:    req.SessionID,
			IPAddress:    req.IPAddress,
			UserAgent:    req.UserAgent,
		}
		uc.searchRepo.RecordSearchEvent(ctx, event)

		// Update search suggestions
		uc.searchRepo.UpdateSearchSuggestion(ctx, req.Query)

		// Record search analytics for performance tracking
		uc.searchRepo.RecordSearchAnalytics(ctx, req.Query, int(total))
	}
	
	searchTime := time.Since(startTime)
	
	return &SearchResponse{
		Products:    productResponses,
		Total:       total,
		Page:        req.Page,
		Limit:       req.Limit,
		TotalPages:  totalPages,
		HasNext:     req.Page < totalPages,
		HasPrev:     req.Page > 1,
		Facets:      facetsResponse,
		Query:       req.Query,
		SearchTime:  searchTime.String(),
		Suggestions: suggestions,
	}, nil
}

// Helper method to convert product to response
func (uc *searchUseCase) toProductResponse(product *entities.Product) *ProductResponse {
	// Use the existing product usecase conversion logic
	// For now, create a simplified response for search results
	response := &ProductResponse{
		ID:               product.ID,
		Name:             product.Name,
		Description:      product.Description,
		ShortDescription: product.ShortDescription,
		SKU:              product.SKU,
		Slug:             product.Slug,
		MetaTitle:        product.MetaTitle,
		MetaDescription:  product.MetaDescription,
		Keywords:         product.Keywords,
		Featured:         product.Featured,
		Visibility:       product.Visibility,
		Price:            product.Price,
		ComparePrice:     product.ComparePrice,
		CostPrice:        product.CostPrice,
		SalePrice:        product.SalePrice,
		SaleStartDate:    product.SaleStartDate,
		SaleEndDate:      product.SaleEndDate,
		Stock:            product.Stock,
		LowStockThreshold: product.LowStockThreshold,
		TrackQuantity:    product.TrackQuantity,
		AllowBackorder:   product.AllowBackorder,
		StockStatus:      product.StockStatus,
		Weight:           product.Weight,
		RequiresShipping: product.RequiresShipping,
		ShippingClass:    product.ShippingClass,
		TaxClass:         product.TaxClass,
		CountryOfOrigin:  product.CountryOfOrigin,
		Status:           product.Status,
		ProductType:      product.ProductType,
		IsDigital:        product.IsDigital,
		CreatedAt:        product.CreatedAt,
		UpdatedAt:        product.UpdatedAt,
	}

	// Calculate computed fields
	response.CurrentPrice = product.Price
	if product.SalePrice != nil && *product.SalePrice > 0 {
		response.CurrentPrice = *product.SalePrice
		response.IsOnSale = true
		if product.Price > 0 {
			response.SaleDiscountPercentage = ((product.Price - *product.SalePrice) / product.Price) * 100
		}
	}

	response.IsLowStock = product.Stock <= product.LowStockThreshold
	response.IsAvailable = product.Status == entities.ProductStatusActive && product.Stock > 0
	response.HasDiscount = product.ComparePrice != nil && *product.ComparePrice > product.Price
	response.HasVariants = len(product.Variants) > 0

	// Set main image
	if len(product.Images) > 0 {
		response.MainImage = product.Images[0].URL
	}

	// Convert related entities
	if product.Category.ID != uuid.Nil {
		response.Category = &ProductCategoryResponse{
			ID:   product.Category.ID,
			Name: product.Category.Name,
			Slug: product.Category.Slug,
		}
	}

	if product.Brand != nil {
		response.Brand = &ProductBrandResponse{
			ID:   product.Brand.ID,
			Name: product.Brand.Name,
			Slug: product.Brand.Slug,
		}
	}

	// Convert images
	for _, img := range product.Images {
		response.Images = append(response.Images, ProductImageResponse{
			ID:       img.ID,
			URL:      img.URL,
			AltText:  img.AltText,
			Position: img.Position,
		})
	}

	// Convert tags
	for _, tag := range product.Tags {
		response.Tags = append(response.Tags, ProductTagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	return response
}

// Helper method to convert search facets to response
func (uc *searchUseCase) toSearchFacetsResponse(facets *repositories.SearchFacets) *SearchFacetsResponse {
	response := &SearchFacetsResponse{}
	
	// Convert categories
	for _, cat := range facets.Categories {
		response.Categories = append(response.Categories, CategoryFacetResponse{
			ID:           cat.ID,
			Name:         cat.Name,
			ProductCount: cat.ProductCount,
		})
	}
	
	// Convert brands
	for _, brand := range facets.Brands {
		response.Brands = append(response.Brands, BrandFacetResponse{
			ID:           brand.ID,
			Name:         brand.Name,
			ProductCount: brand.ProductCount,
		})
	}

	// Convert price range
	response.PriceRange = PriceRangeFacetResponse{
		MinPrice: facets.PriceRange.MinPrice,
		MaxPrice: facets.PriceRange.MaxPrice,
		Ranges:   make([]PriceRangeResponse, len(facets.PriceRange.Ranges)),
	}
	for i, priceRange := range facets.PriceRange.Ranges {
		response.PriceRange.Ranges[i] = PriceRangeResponse{
			Min:          priceRange.Min,
			Max:          priceRange.Max,
			Label:        priceRange.Label,
			ProductCount: priceRange.ProductCount,
		}
	}

	// Convert tags
	for _, tag := range facets.Tags {
		response.Tags = append(response.Tags, TagFacetResponse{
			ID:           tag.ID,
			Name:         tag.Name,
			ProductCount: tag.ProductCount,
		})
	}

	// Convert status facets
	for _, status := range facets.Status {
		response.Status = append(response.Status, StatusFacetResponse{
			Status:       status.Status,
			Label:        status.Label,
			ProductCount: status.ProductCount,
		})
	}

	// Convert product type facets
	for _, productType := range facets.ProductTypes {
		response.ProductTypes = append(response.ProductTypes, ProductTypeFacetResponse{
			Type:         productType.Type,
			Label:        productType.Label,
			ProductCount: productType.ProductCount,
		})
	}

	// Convert availability facets
	for _, availability := range facets.Availability {
		response.Availability = append(response.Availability, AvailabilityFacetResponse{
			Status:       availability.Status,
			Label:        availability.Label,
			ProductCount: availability.ProductCount,
		})
	}

	// Convert rating facets
	for _, rating := range facets.Ratings {
		response.Ratings = append(response.Ratings, RatingFacetResponse{
			Rating:       rating.Rating,
			Label:        rating.Label,
			ProductCount: rating.ProductCount,
		})
	}

	// Convert shipping facets
	for _, shipping := range facets.Shipping {
		response.Shipping = append(response.Shipping, ShippingFacetResponse{
			Type:         shipping.Type,
			Label:        shipping.Label,
			ProductCount: shipping.ProductCount,
		})
	}

	return response
}

// GetSearchSuggestions gets search suggestions
func (uc *searchUseCase) GetSearchSuggestions(ctx context.Context, query string, limit int) ([]string, error) {
	suggestions, err := uc.searchRepo.GetSearchSuggestions(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get search suggestions: %w", err)
	}

	result := make([]string, len(suggestions))
	for i, s := range suggestions {
		result[i] = s.Query
	}

	return result, nil
}

// GetSearchFacets gets search facets
func (uc *searchUseCase) GetSearchFacets(ctx context.Context, query string) (*SearchFacetsResponse, error) {
	facets, err := uc.searchRepo.GetSearchFacets(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get search facets: %w", err)
	}

	return uc.toSearchFacetsResponse(facets), nil
}

// RecordSearchEvent records a search event
func (uc *searchUseCase) RecordSearchEvent(ctx context.Context, req RecordSearchEventRequest) error {
	event := &entities.SearchEvent{
		Query:            req.Query,
		UserID:           req.UserID,
		ResultsCount:     req.ResultsCount,
		ClickedProductID: req.ClickedProductID,
		SessionID:        req.SessionID,
		IPAddress:        req.IPAddress,
		UserAgent:        req.UserAgent,
	}

	if err := uc.searchRepo.RecordSearchEvent(ctx, event); err != nil {
		return fmt.Errorf("failed to record search event: %w", err)
	}

	// Update search suggestions
	if req.Query != "" {
		uc.searchRepo.UpdateSearchSuggestion(ctx, req.Query)
	}

	return nil
}

// GetPopularSearchTerms gets popular search terms
func (uc *searchUseCase) GetPopularSearchTerms(ctx context.Context, limit int, period string) ([]PopularSearchResponse, error) {
	popularSearches, err := uc.searchRepo.GetPopularSearchTerms(ctx, limit, period)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular search terms: %w", err)
	}

	result := make([]PopularSearchResponse, len(popularSearches))
	for i, ps := range popularSearches {
		result[i] = PopularSearchResponse{
			Query:       ps.Query,
			SearchCount: ps.SearchCount,
		}
	}

	return result, nil
}

// SaveSearchHistory saves user search history
func (uc *searchUseCase) SaveSearchHistory(ctx context.Context, userID uuid.UUID, req SaveSearchHistoryRequest) error {
	history := &entities.SearchHistory{
		UserID:  userID,
		Query:   req.Query,
		Filters: req.Filters,
	}

	if err := uc.searchRepo.SaveSearchHistory(ctx, history); err != nil {
		return fmt.Errorf("failed to save search history: %w", err)
	}

	return nil
}

// GetUserSearchHistory gets user search history
func (uc *searchUseCase) GetUserSearchHistory(ctx context.Context, userID uuid.UUID, limit int) ([]UserSearchHistoryResponse, error) {
	history, err := uc.searchRepo.GetUserSearchHistory(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user search history: %w", err)
	}

	result := make([]UserSearchHistoryResponse, len(history))
	for i, h := range history {
		result[i] = UserSearchHistoryResponse{
			ID:        h.ID,
			Query:     h.Query,
			Filters:   h.Filters,
			CreatedAt: h.CreatedAt,
		}
	}

	return result, nil
}

// ClearUserSearchHistory clears user search history
func (uc *searchUseCase) ClearUserSearchHistory(ctx context.Context, userID uuid.UUID) error {
	if err := uc.searchRepo.ClearUserSearchHistory(ctx, userID); err != nil {
		return fmt.Errorf("failed to clear user search history: %w", err)
	}

	return nil
}

// SaveSearchFilter saves a search filter
func (uc *searchUseCase) SaveSearchFilter(ctx context.Context, userID uuid.UUID, req SaveSearchFilterRequest) (*SearchFilterResponse, error) {
	filter := &entities.SearchFilter{
		UserID:    userID,
		Name:      req.Name,
		Query:     req.Query,
		Filters:   req.Filters,
		IsDefault: req.IsDefault,
		IsPublic:  req.IsPublic,
	}

	if err := uc.searchRepo.SaveSearchFilter(ctx, filter); err != nil {
		return nil, fmt.Errorf("failed to save search filter: %w", err)
	}

	return &SearchFilterResponse{
		ID:         filter.ID,
		Name:       filter.Name,
		Query:      filter.Query,
		Filters:    filter.Filters,
		IsDefault:  filter.IsDefault,
		IsPublic:   filter.IsPublic,
		UsageCount: filter.UsageCount,
		CreatedAt:  filter.CreatedAt,
		UpdatedAt:  filter.UpdatedAt,
	}, nil
}

// GetUserSearchFilters gets user search filters
func (uc *searchUseCase) GetUserSearchFilters(ctx context.Context, userID uuid.UUID) ([]SearchFilterResponse, error) {
	filters, err := uc.searchRepo.GetUserSearchFilters(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user search filters: %w", err)
	}

	result := make([]SearchFilterResponse, len(filters))
	for i, f := range filters {
		result[i] = SearchFilterResponse{
			ID:         f.ID,
			Name:       f.Name,
			Query:      f.Query,
			Filters:    f.Filters,
			IsDefault:  f.IsDefault,
			IsPublic:   f.IsPublic,
			UsageCount: f.UsageCount,
			CreatedAt:  f.CreatedAt,
			UpdatedAt:  f.UpdatedAt,
		}
	}

	return result, nil
}

// UpdateSearchFilter updates a search filter
func (uc *searchUseCase) UpdateSearchFilter(ctx context.Context, userID uuid.UUID, filterID uuid.UUID, req UpdateSearchFilterRequest) error {
	filter, err := uc.searchRepo.GetSearchFilter(ctx, filterID)
	if err != nil {
		return fmt.Errorf("failed to get search filter: %w", err)
	}

	// Check ownership
	if filter.UserID != userID {
		return fmt.Errorf("unauthorized to update this filter")
	}

	// Update fields
	filter.Name = req.Name
	filter.Query = req.Query
	filter.Filters = req.Filters
	filter.IsDefault = req.IsDefault
	filter.IsPublic = req.IsPublic

	if err := uc.searchRepo.UpdateSearchFilter(ctx, filter); err != nil {
		return fmt.Errorf("failed to update search filter: %w", err)
	}

	return nil
}

// DeleteSearchFilter deletes a search filter
func (uc *searchUseCase) DeleteSearchFilter(ctx context.Context, userID uuid.UUID, filterID uuid.UUID) error {
	filter, err := uc.searchRepo.GetSearchFilter(ctx, filterID)
	if err != nil {
		return fmt.Errorf("failed to get search filter: %w", err)
	}

	// Check ownership
	if filter.UserID != userID {
		return fmt.Errorf("unauthorized to delete this filter")
	}

	if err := uc.searchRepo.DeleteSearchFilter(ctx, filterID); err != nil {
		return fmt.Errorf("failed to delete search filter: %w", err)
	}

	return nil
}

// GetAutocomplete gets autocomplete suggestions
func (uc *searchUseCase) GetAutocomplete(ctx context.Context, query string, limit int) (*AutocompleteResponse, error) {
	response := &AutocompleteResponse{}

	// Get product suggestions
	productParams := repositories.ProductSearchParams{
		Query:  query,
		Limit:  limit / 3, // Divide limit among different types
		Offset: 0,
	}

	products, err := uc.productRepo.Search(ctx, productParams)
	if err == nil && len(products) > 0 {
		for _, p := range products {
			if len(response.Products) < limit/3 {
				suggestion := ProductSuggestionResponse{
					ID:    p.ID,
					Name:  p.Name,
					Price: p.Price,
				}
				if len(p.Images) > 0 {
					suggestion.Image = p.Images[0].URL
				}
				response.Products = append(response.Products, suggestion)
			}
		}
	}

	// Get search suggestions from history
	suggestions, err := uc.searchRepo.GetSearchSuggestions(ctx, query, limit)
	if err == nil {
		for _, s := range suggestions {
			if len(response.Suggestions) < limit {
				response.Suggestions = append(response.Suggestions, s.Query)
			}
		}
	}

	return response, nil
}

// GetEnhancedAutocomplete provides enhanced autocomplete with multiple sources
func (uc *searchUseCase) GetEnhancedAutocomplete(ctx context.Context, req EnhancedAutocompleteRequest) (*EnhancedAutocompleteResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}

	response := &EnhancedAutocompleteResponse{}

	// Get suggestions by type
	if len(req.Types) == 0 || contains(req.Types, "product") {
		products, _ := uc.searchRepo.GetProductSuggestions(ctx, req.Query, req.Limit/4)
		response.Products = uc.convertToAutocompleteSuggestions(products)
	}

	if len(req.Types) == 0 || contains(req.Types, "category") {
		categories, _ := uc.searchRepo.GetCategorySuggestions(ctx, req.Query, req.Limit/4)
		response.Categories = uc.convertToAutocompleteSuggestions(categories)
	}

	if len(req.Types) == 0 || contains(req.Types, "brand") {
		brands, _ := uc.searchRepo.GetBrandSuggestions(ctx, req.Query, req.Limit/4)
		response.Brands = uc.convertToAutocompleteSuggestions(brands)
	}

	if len(req.Types) == 0 || contains(req.Types, "query") {
		queries, _ := uc.searchRepo.GetAutocompleteEntries(ctx, req.Query, []string{"query"}, req.Limit/4)
		response.Queries = uc.convertToAutocompleteSuggestions(queries)
	}

	// Get trending suggestions if requested
	if req.IncludeTrending {
		trending, _ := uc.searchRepo.GetTrendingSuggestions(ctx, 5)
		response.Trending = uc.convertToAutocompleteSuggestions(trending)
	}

	// Get personalized suggestions if requested and user is provided
	if req.IncludePersonalized && req.UserID != nil {
		personalized, _ := uc.searchRepo.GetPersonalizedSuggestions(ctx, *req.UserID, req.Query, 5)
		response.Personalized = uc.convertToAutocompleteSuggestions(personalized)
	}

	return response, nil
}

// GetPersonalizedAutocomplete provides personalized autocomplete for a user
func (uc *searchUseCase) GetPersonalizedAutocomplete(ctx context.Context, userID uuid.UUID, query string, limit int) (*EnhancedAutocompleteResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	response := &EnhancedAutocompleteResponse{}

	// Get personalized suggestions
	personalized, err := uc.searchRepo.GetPersonalizedSuggestions(ctx, userID, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get personalized suggestions: %w", err)
	}

	response.Personalized = uc.convertToAutocompleteSuggestions(personalized)

	// Also get general suggestions as fallback
	general, _ := uc.searchRepo.GetAutocompleteEntries(ctx, query, nil, limit/2)
	response.Queries = uc.convertToAutocompleteSuggestions(general)

	return response, nil
}

// GetTrendingSearches retrieves trending search terms
func (uc *searchUseCase) GetTrendingSearches(ctx context.Context, limit int) ([]TrendingSearchResponse, error) {
	if limit <= 0 {
		limit = 20
	}

	trends, err := uc.searchRepo.GetSearchTrends(ctx, "daily", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get search trends: %w", err)
	}

	response := make([]TrendingSearchResponse, len(trends))
	for i, trend := range trends {
		response[i] = TrendingSearchResponse{
			Query:       trend.Query,
			SearchCount: trend.SearchCount,
			Period:      trend.Period,
			Trend:       "stable", // TODO: Calculate actual trend
		}
	}

	return response, nil
}

// GetUserSearchPreferences retrieves user search preferences
func (uc *searchUseCase) GetUserSearchPreferences(ctx context.Context, userID uuid.UUID) (*UserSearchPreferencesResponse, error) {
	prefs, err := uc.searchRepo.GetUserSearchPreferences(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user search preferences: %w", err)
	}

	return &UserSearchPreferencesResponse{
		UserID:              prefs.UserID,
		PreferredCategories: prefs.PreferredCategories,
		PreferredBrands:     prefs.PreferredBrands,
		SearchLanguage:      prefs.SearchLanguage,
		AutocompleteEnabled: prefs.AutocompleteEnabled,
		SearchHistoryEnabled: prefs.SearchHistoryEnabled,
		PersonalizedResults: prefs.PersonalizedResults,
	}, nil
}

// UpdateUserSearchPreferences updates user search preferences
func (uc *searchUseCase) UpdateUserSearchPreferences(ctx context.Context, userID uuid.UUID, req UpdateSearchPreferencesRequest) error {
	prefs, err := uc.searchRepo.GetUserSearchPreferences(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get current preferences: %w", err)
	}

	// Update only provided fields
	if req.PreferredCategories != nil {
		prefs.PreferredCategories = *req.PreferredCategories
	}
	if req.PreferredBrands != nil {
		prefs.PreferredBrands = *req.PreferredBrands
	}
	if req.SearchLanguage != nil {
		prefs.SearchLanguage = *req.SearchLanguage
	}
	if req.AutocompleteEnabled != nil {
		prefs.AutocompleteEnabled = *req.AutocompleteEnabled
	}
	if req.SearchHistoryEnabled != nil {
		prefs.SearchHistoryEnabled = *req.SearchHistoryEnabled
	}
	if req.PersonalizedResults != nil {
		prefs.PersonalizedResults = *req.PersonalizedResults
	}

	return uc.searchRepo.SaveUserSearchPreferences(ctx, prefs)
}

// RecordAutocompleteClick records autocomplete usage analytics
func (uc *searchUseCase) RecordAutocompleteClick(ctx context.Context, req AutocompleteClickRequest) error {
	// Record the click/search in autocomplete entry
	if err := uc.searchRepo.IncrementAutocompleteUsage(ctx, req.EntryID, req.IsClick); err != nil {
		return fmt.Errorf("failed to record autocomplete usage: %w", err)
	}

	// Also record as search event if it's a search
	if !req.IsClick && req.Query != "" {
		event := &entities.SearchEvent{
			Query:     req.Query,
			UserID:    req.UserID,
			SessionID: req.SessionID,
		}
		uc.searchRepo.RecordSearchEvent(ctx, event)
	}

	return nil
}

// GetSearchTrends retrieves search trends for analytics
func (uc *searchUseCase) GetSearchTrends(ctx context.Context, period string, limit int) ([]SearchTrendResponse, error) {
	if limit <= 0 {
		limit = 50
	}

	trends, err := uc.searchRepo.GetSearchTrends(ctx, period, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get search trends: %w", err)
	}

	response := make([]SearchTrendResponse, len(trends))
	for i, trend := range trends {
		response[i] = SearchTrendResponse{
			Query:       trend.Query,
			SearchCount: trend.SearchCount,
			Period:      trend.Period,
			Date:        trend.Date,
			Change:      0, // TODO: Calculate change from previous period
		}
	}

	return response, nil
}

// RebuildAutocompleteIndex rebuilds the autocomplete index
func (uc *searchUseCase) RebuildAutocompleteIndex(ctx context.Context) error {
	return uc.searchRepo.RebuildAutocompleteIndex(ctx)
}

// CleanupSearchData cleans up old search data
func (uc *searchUseCase) CleanupSearchData(ctx context.Context, days int) error {
	return uc.searchRepo.CleanupOldAutocompleteEntries(ctx, days)
}

// Helper methods

func (uc *searchUseCase) convertToAutocompleteSuggestions(entries []*entities.AutocompleteEntry) []AutocompleteSuggestion {
	suggestions := make([]AutocompleteSuggestion, len(entries))
	for i, entry := range entries {
		var metadata interface{}
		if entry.Metadata != "" {
			// Parse JSON metadata
			json.Unmarshal([]byte(entry.Metadata), &metadata)
		}

		suggestions[i] = AutocompleteSuggestion{
			ID:          entry.ID,
			Type:        entry.Type,
			Value:       entry.Value,
			DisplayText: entry.DisplayText,
			EntityID:    entry.EntityID,
			Priority:    entry.Priority,
			SearchCount: entry.SearchCount,
			ClickCount:  entry.ClickCount,
			Metadata:    metadata,
		}
	}
	return suggestions
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetSearchAnalytics gets search analytics for admin dashboard
func (uc *searchUseCase) GetSearchAnalytics(ctx context.Context, req SearchAnalyticsRequest) (*SearchAnalyticsResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 50
	}

	analytics, err := uc.searchRepo.GetSearchAnalytics(ctx, req.StartDate, req.EndDate, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get search analytics: %w", err)
	}

	response := &SearchAnalyticsResponse{
		Analytics: make([]SearchAnalyticItem, len(analytics)),
		Summary: SearchAnalyticsSummary{},
	}

	var totalSearches int
	var totalQueries int
	var totalResults float64
	var topQuery string
	var maxSearches int

	for i, item := range analytics {
		searches := int(item["total_searches"].(int64))
		avgResults := float64(item["avg_result_count"].(float64))

		response.Analytics[i] = SearchAnalyticItem{
			Query:              item["query"].(string),
			TotalSearches:      searches,
			AvgResultCount:     avgResults,
			AvgClickThroughRate: item["avg_ctr"].(float64),
			AvgConversionRate:  item["avg_conversion_rate"].(float64),
			LastSearched:       item["last_searched"].(time.Time).Format("2006-01-02"),
		}

		totalSearches += searches
		totalQueries++
		totalResults += avgResults

		if searches > maxSearches {
			maxSearches = searches
			topQuery = item["query"].(string)
		}
	}

	// Calculate summary
	response.Summary.TotalQueries = totalQueries
	response.Summary.TotalSearches = totalSearches
	response.Summary.TopQuery = topQuery
	if totalQueries > 0 {
		response.Summary.AvgResultsPerQuery = totalResults / float64(totalQueries)
	}

	return response, nil
}

// EnhancedSearch performs enhanced search with dynamic faceting
func (uc *searchUseCase) EnhancedSearch(ctx context.Context, req *EnhancedSearchRequest) (*EnhancedSearchResponse, error) {
	startTime := time.Now()

	// Convert request to repository parameters
	params := repositories.EnhancedSearchParams{
		FullTextSearchParams: repositories.FullTextSearchParams{
			Query:       req.Query,
			CategoryIDs: uc.convertStringIDsToUUIDs(req.CategoryIDs),
			BrandIDs:    uc.convertStringIDsToUUIDs(req.BrandIDs),
			MinPrice:    req.MinPrice,
			MaxPrice:    req.MaxPrice,
			Tags:        req.TagIDs,
			Featured:    req.Featured,
			InStock:     req.InStock,
			OnSale:      req.OnSale,
			SortBy:      req.SortBy,
			SortOrder:   req.SortOrder,
			Limit:       req.Limit,
			Offset:      (req.Page - 1) * req.Limit,
		},
		IncludeFacets: req.IncludeFacets,
		DynamicFacets: req.DynamicFacets,
	}

	// Perform enhanced search
	products, total, facets, err := uc.searchRepo.EnhancedSearch(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to perform enhanced search: %w", err)
	}

	// Convert products to response format
	var productResponses []ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, *uc.toProductResponse(product))
	}

	// Calculate pagination
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	// Convert facets to response format
	var facetsResponse *DynamicSearchFacetsResponse
	if facets != nil {
		facetsResponse = uc.toDynamicSearchFacetsResponse(facets)
	}

	// Get suggestions if query is provided
	var suggestions []string
	if req.Query != "" {
		suggestionEntities, _ := uc.searchRepo.GetSearchSuggestions(ctx, req.Query, 5)
		for _, s := range suggestionEntities {
			suggestions = append(suggestions, s.Query)
		}
	}

	// Record search event and analytics
	if req.Query != "" {
		var userID *uuid.UUID
		if req.UserID != nil {
			if parsedUUID, err := uuid.Parse(*req.UserID); err == nil {
				userID = &parsedUUID
			}
		}

		event := &entities.SearchEvent{
			Query:        req.Query,
			UserID:       userID,
			ResultsCount: int(total),
			SessionID:    req.SessionID,
			IPAddress:    req.IPAddress,
			UserAgent:    req.UserAgent,
		}
		uc.searchRepo.RecordSearchEvent(ctx, event)

		// Update search suggestions
		uc.searchRepo.UpdateSearchSuggestion(ctx, req.Query)

		// Record search analytics for performance tracking
		uc.searchRepo.RecordSearchAnalytics(ctx, req.Query, int(total))
	}

	searchTime := time.Since(startTime)

	return &EnhancedSearchResponse{
		Products:    productResponses,
		Total:       total,
		Page:        req.Page,
		Limit:       req.Limit,
		TotalPages:  totalPages,
		HasNext:     req.Page < totalPages,
		HasPrev:     req.Page > 1,
		Facets:      facetsResponse,
		Query:       req.Query,
		SearchTime:  searchTime.String(),
		Suggestions: suggestions,
	}, nil
}

// Helper methods
func (uc *searchUseCase) convertStringIDsToUUIDs(stringIDs []string) []uuid.UUID {
	var uuids []uuid.UUID
	for _, id := range stringIDs {
		if parsedUUID, err := uuid.Parse(id); err == nil {
			uuids = append(uuids, parsedUUID)
		}
	}
	return uuids
}

func (uc *searchUseCase) toDynamicSearchFacetsResponse(facets *repositories.DynamicSearchFacets) *DynamicSearchFacetsResponse {
	response := &DynamicSearchFacetsResponse{
		TotalCount: facets.TotalCount,
	}

	// Convert category facets
	for _, cat := range facets.Categories {
		response.Categories = append(response.Categories, DynamicCategoryFacetResponse{
			ID:           cat.ID.String(),
			Name:         cat.Name,
			ProductCount: cat.ProductCount,
			IsSelected:   cat.IsSelected,
			IsDisabled:   cat.IsDisabled,
		})
	}

	// Convert brand facets
	for _, brand := range facets.Brands {
		response.Brands = append(response.Brands, DynamicBrandFacetResponse{
			ID:           brand.ID.String(),
			Name:         brand.Name,
			ProductCount: brand.ProductCount,
			IsSelected:   brand.IsSelected,
			IsDisabled:   brand.IsDisabled,
		})
	}

	// Convert tag facets
	for _, tag := range facets.Tags {
		response.Tags = append(response.Tags, DynamicTagFacetResponse{
			ID:           tag.ID.String(),
			Name:         tag.Name,
			ProductCount: tag.ProductCount,
			IsSelected:   tag.IsSelected,
			IsDisabled:   tag.IsDisabled,
		})
	}

	// Convert status facets
	for _, status := range facets.Status {
		response.Status = append(response.Status, DynamicStatusFacetResponse{
			Status:       string(status.Status),
			Label:        status.Label,
			ProductCount: status.ProductCount,
			IsSelected:   status.IsSelected,
			IsDisabled:   status.IsDisabled,
		})
	}

	// Convert price range facet
	var priceRanges []PriceRangeResponse
	for _, pr := range facets.PriceRange.Ranges {
		priceRanges = append(priceRanges, PriceRangeResponse{
			Min:          pr.Min,
			Max:          pr.Max,
			Label:        pr.Label,
			ProductCount: pr.ProductCount,
		})
	}

	response.PriceRange = DynamicPriceRangeFacetResponse{
		MinPrice:    facets.PriceRange.MinPrice,
		MaxPrice:    facets.PriceRange.MaxPrice,
		Ranges:      priceRanges,
		SelectedMin: facets.PriceRange.SelectedMin,
		SelectedMax: facets.PriceRange.SelectedMax,
	}

	return response
}

// Enhanced Search Request and Response Types

// EnhancedSearchRequest represents enhanced search request with multi-select filters
type EnhancedSearchRequest struct {
	Query              string    `json:"query"`
	CategoryIDs        []string  `json:"category_ids"`
	BrandIDs           []string  `json:"brand_ids"`
	MinPrice           *float64  `json:"min_price"`
	MaxPrice           *float64  `json:"max_price"`
	TagIDs             []string  `json:"tag_ids"`
	StatusList         []string  `json:"status_list"`
	ProductTypes       []string  `json:"product_types"`
	AvailabilityTypes  []string  `json:"availability_types"`
	ShippingTypes      []string  `json:"shipping_types"`
	RatingMin          *float64  `json:"rating_min"`
	RatingMax          *float64  `json:"rating_max"`
	Featured           *bool     `json:"featured"`
	InStock            *bool     `json:"in_stock"`
	OnSale             *bool     `json:"on_sale"`
	SortBy             string    `json:"sort_by"`
	SortOrder          string    `json:"sort_order"`
	Page               int       `json:"page"`
	Limit              int       `json:"limit"`
	IncludeFacets      bool      `json:"include_facets"`
	DynamicFacets      bool      `json:"dynamic_facets"`
	SessionID          string    `json:"session_id"`
	IPAddress          string    `json:"ip_address"`
	UserAgent          string    `json:"user_agent"`
	UserID             *string   `json:"user_id"`
}

// EnhancedSearchResponse represents enhanced search response with dynamic facets
type EnhancedSearchResponse struct {
	Products     []ProductResponse                    `json:"products"`
	Total        int64                                `json:"total"`
	Page         int                                  `json:"page"`
	Limit        int                                  `json:"limit"`
	TotalPages   int                                  `json:"total_pages"`
	HasNext      bool                                 `json:"has_next"`
	HasPrev      bool                                 `json:"has_prev"`
	Facets       *DynamicSearchFacetsResponse         `json:"facets,omitempty"`
	Query        string                               `json:"query"`
	SearchTime   string                               `json:"search_time"`
	Suggestions  []string                             `json:"suggestions,omitempty"`
}

// DynamicSearchFacetsResponse represents dynamic facets response
type DynamicSearchFacetsResponse struct {
	Categories   []DynamicCategoryFacetResponse   `json:"categories"`
	Brands       []DynamicBrandFacetResponse      `json:"brands"`
	PriceRange   DynamicPriceRangeFacetResponse   `json:"price_range"`
	Tags         []DynamicTagFacetResponse        `json:"tags"`
	Status       []DynamicStatusFacetResponse     `json:"status"`
	ProductTypes []DynamicProductTypeFacetResponse `json:"product_types"`
	Availability []DynamicAvailabilityFacetResponse `json:"availability"`
	Shipping     []DynamicShippingFacetResponse   `json:"shipping"`
	TotalCount   int64                            `json:"total_count"`
}

// Dynamic facet response types
type DynamicCategoryFacetResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ProductCount int64  `json:"product_count"`
	IsSelected   bool   `json:"is_selected"`
	IsDisabled   bool   `json:"is_disabled"`
}

type DynamicBrandFacetResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ProductCount int64  `json:"product_count"`
	IsSelected   bool   `json:"is_selected"`
	IsDisabled   bool   `json:"is_disabled"`
}

type DynamicTagFacetResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ProductCount int64  `json:"product_count"`
	IsSelected   bool   `json:"is_selected"`
	IsDisabled   bool   `json:"is_disabled"`
}

type DynamicStatusFacetResponse struct {
	Status       string `json:"status"`
	Label        string `json:"label"`
	ProductCount int64  `json:"product_count"`
	IsSelected   bool   `json:"is_selected"`
	IsDisabled   bool   `json:"is_disabled"`
}

type DynamicProductTypeFacetResponse struct {
	Type         string `json:"type"`
	Label        string `json:"label"`
	ProductCount int64  `json:"product_count"`
	IsSelected   bool   `json:"is_selected"`
	IsDisabled   bool   `json:"is_disabled"`
}

type DynamicAvailabilityFacetResponse struct {
	Status       string `json:"status"`
	Label        string `json:"label"`
	ProductCount int64  `json:"product_count"`
	IsSelected   bool   `json:"is_selected"`
	IsDisabled   bool   `json:"is_disabled"`
}

type DynamicShippingFacetResponse struct {
	Type         string `json:"type"`
	Label        string `json:"label"`
	ProductCount int64  `json:"product_count"`
	IsSelected   bool   `json:"is_selected"`
	IsDisabled   bool   `json:"is_disabled"`
}

type DynamicPriceRangeFacetResponse struct {
	MinPrice    float64                      `json:"min_price"`
	MaxPrice    float64                      `json:"max_price"`
	Ranges      []PriceRangeResponse         `json:"ranges"`
	SelectedMin *float64                     `json:"selected_min"`
	SelectedMax *float64                     `json:"selected_max"`
}

// Enhanced Autocomplete Types

type EnhancedAutocompleteRequest struct {
	Query   string   `json:"query"`
	Types   []string `json:"types"` // product, category, brand, tag, query
	Limit   int      `json:"limit"`
	UserID  *uuid.UUID `json:"user_id"`
	IncludePersonalized bool `json:"include_personalized"`
	IncludeTrending     bool `json:"include_trending"`
}

type EnhancedAutocompleteResponse struct {
	Products    []AutocompleteSuggestion `json:"products"`
	Categories  []AutocompleteSuggestion `json:"categories"`
	Brands      []AutocompleteSuggestion `json:"brands"`
	Queries     []AutocompleteSuggestion `json:"queries"`
	Trending    []AutocompleteSuggestion `json:"trending"`
	Personalized []AutocompleteSuggestion `json:"personalized"`
}

type AutocompleteSuggestion struct {
	ID          uuid.UUID   `json:"id"`
	Type        string      `json:"type"`
	Value       string      `json:"value"`
	DisplayText string      `json:"display_text"`
	EntityID    *uuid.UUID  `json:"entity_id"`
	Priority    int         `json:"priority"`
	SearchCount int         `json:"search_count"`
	ClickCount  int         `json:"click_count"`
	Metadata    interface{} `json:"metadata"`
}

type TrendingSearchResponse struct {
	Query       string `json:"query"`
	SearchCount int    `json:"search_count"`
	Period      string `json:"period"`
	Trend       string `json:"trend"` // up, down, stable
}

type UserSearchPreferencesResponse struct {
	UserID              uuid.UUID `json:"user_id"`
	PreferredCategories []string  `json:"preferred_categories"`
	PreferredBrands     []string  `json:"preferred_brands"`
	SearchLanguage      string    `json:"search_language"`
	AutocompleteEnabled bool      `json:"autocomplete_enabled"`
	SearchHistoryEnabled bool     `json:"search_history_enabled"`
	PersonalizedResults  bool     `json:"personalized_results"`
}

type UpdateSearchPreferencesRequest struct {
	PreferredCategories  *[]string `json:"preferred_categories"`
	PreferredBrands      *[]string `json:"preferred_brands"`
	SearchLanguage       *string   `json:"search_language"`
	AutocompleteEnabled  *bool     `json:"autocomplete_enabled"`
	SearchHistoryEnabled *bool     `json:"search_history_enabled"`
	PersonalizedResults  *bool     `json:"personalized_results"`
}

type AutocompleteClickRequest struct {
	EntryID   uuid.UUID  `json:"entry_id"`
	UserID    *uuid.UUID `json:"user_id"`
	SessionID string     `json:"session_id"`
	Query     string     `json:"query"`
	IsClick   bool       `json:"is_click"` // true for click, false for search
}

type SearchTrendResponse struct {
	Query       string    `json:"query"`
	SearchCount int       `json:"search_count"`
	Period      string    `json:"period"`
	Date        time.Time `json:"date"`
	Change      float64   `json:"change"` // percentage change from previous period
}
