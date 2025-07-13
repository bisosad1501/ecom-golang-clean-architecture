package repositories

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// SearchRepository defines the interface for search-related operations
type SearchRepository interface {
	// Search Events
	RecordSearchEvent(ctx context.Context, event *entities.SearchEvent) error
	GetSearchEvents(ctx context.Context, filters SearchEventFilters) ([]*entities.SearchEvent, error)
	GetPopularSearchTerms(ctx context.Context, limit int, period string) ([]*entities.PopularSearch, error)
	
	// Search Suggestions
	GetSearchSuggestions(ctx context.Context, query string, limit int) ([]*entities.SearchSuggestion, error)
	UpdateSearchSuggestion(ctx context.Context, query string) error
	
	// Search History
	SaveSearchHistory(ctx context.Context, history *entities.SearchHistory) error
	GetUserSearchHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*entities.SearchHistory, error)
	ClearUserSearchHistory(ctx context.Context, userID uuid.UUID) error
	
	// Search Filters
	SaveSearchFilter(ctx context.Context, filter *entities.SearchFilter) error
	GetUserSearchFilters(ctx context.Context, userID uuid.UUID) ([]*entities.SearchFilter, error)
	GetSearchFilter(ctx context.Context, id uuid.UUID) (*entities.SearchFilter, error)
	UpdateSearchFilter(ctx context.Context, filter *entities.SearchFilter) error
	DeleteSearchFilter(ctx context.Context, id uuid.UUID) error
	
	// Advanced Search
	FullTextSearch(ctx context.Context, params FullTextSearchParams) ([]*entities.Product, int64, error)
	GetSearchFacets(ctx context.Context, query string) (*SearchFacets, error)

	// Enhanced Faceted Search
	EnhancedSearch(ctx context.Context, params EnhancedSearchParams) ([]*entities.Product, int64, *DynamicSearchFacets, error)
	GetDynamicFacets(ctx context.Context, params EnhancedSearchParams) (*DynamicSearchFacets, error)
	GetFacetCounts(ctx context.Context, params EnhancedSearchParams, facetType string) (map[string]int64, error)

	// Search Analytics
	RecordSearchAnalytics(ctx context.Context, query string, resultCount int) error
	GetSearchAnalytics(ctx context.Context, startDate, endDate time.Time, limit int) ([]map[string]interface{}, error)

	// Enhanced Autocomplete
	GetAutocompleteEntries(ctx context.Context, query string, types []string, limit int) ([]*entities.AutocompleteEntry, error)
	CreateAutocompleteEntry(ctx context.Context, entry *entities.AutocompleteEntry) error
	UpdateAutocompleteEntry(ctx context.Context, entry *entities.AutocompleteEntry) error
	DeleteAutocompleteEntry(ctx context.Context, id uuid.UUID) error
	IncrementAutocompleteUsage(ctx context.Context, id uuid.UUID, isClick bool) error

	// Search Trends
	GetSearchTrends(ctx context.Context, period string, limit int) ([]*entities.SearchTrend, error)
	UpdateSearchTrend(ctx context.Context, query string, period string) error

	// User Search Preferences
	GetUserSearchPreferences(ctx context.Context, userID uuid.UUID) (*entities.UserSearchPreference, error)
	SaveUserSearchPreferences(ctx context.Context, prefs *entities.UserSearchPreference) error

	// Search Sessions
	CreateSearchSession(ctx context.Context, session *entities.SearchSession) error
	UpdateSearchSession(ctx context.Context, session *entities.SearchSession) error
	GetSearchSession(ctx context.Context, sessionID string) (*entities.SearchSession, error)

	// Smart Suggestions
	GetPersonalizedSuggestions(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*entities.AutocompleteEntry, error)
	GetTrendingSuggestions(ctx context.Context, limit int) ([]*entities.AutocompleteEntry, error)
	GetCategorySuggestions(ctx context.Context, query string, limit int) ([]*entities.AutocompleteEntry, error)
	GetBrandSuggestions(ctx context.Context, query string, limit int) ([]*entities.AutocompleteEntry, error)
	GetProductSuggestions(ctx context.Context, query string, limit int) ([]*entities.AutocompleteEntry, error)

	// Enhanced Smart Autocomplete
	GetSmartAutocomplete(ctx context.Context, req entities.SmartAutocompleteRequest) (*entities.SmartAutocompleteResponse, error)
	GetFuzzyMatches(ctx context.Context, query string, types []string, limit int) ([]*entities.AutocompleteEntry, error)
	GetSynonymSuggestions(ctx context.Context, query string, limit int) ([]*entities.AutocompleteEntry, error)
	GetPopularSuggestions(ctx context.Context, limit int, timeframe string) ([]*entities.AutocompleteEntry, error)
	GetUserAutocompleteHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*entities.AutocompleteEntry, error)

	// Autocomplete Analytics
	TrackAutocompleteClick(ctx context.Context, entryID uuid.UUID, userID *uuid.UUID, sessionID string) error
	TrackAutocompleteImpression(ctx context.Context, entryID uuid.UUID, userID *uuid.UUID, sessionID string) error
	UpdateAutocompleteTrending(ctx context.Context) error
	CalculateAutocompleteScores(ctx context.Context) error

	// Autocomplete Management
	RebuildAutocompleteIndex(ctx context.Context) error
	CleanupOldAutocompleteEntries(ctx context.Context, days int) error
}

// SearchEventFilters represents filters for search events
type SearchEventFilters struct {
	UserID    *uuid.UUID `json:"user_id"`
	Query     string     `json:"query"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	SortBy    string     `json:"sort_by"`    // created_at, query, results_count
	SortOrder string     `json:"sort_order"` // asc, desc
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

// FullTextSearchParams represents parameters for full-text search
type FullTextSearchParams struct {
	Query       string                  `json:"query"`
	CategoryIDs []uuid.UUID             `json:"category_ids"`
	BrandIDs    []uuid.UUID             `json:"brand_ids"`
	MinPrice    *float64                `json:"min_price"`
	MaxPrice    *float64                `json:"max_price"`
	InStock     *bool                   `json:"in_stock"`
	Featured    *bool                   `json:"featured"`
	OnSale      *bool                   `json:"on_sale"`
	Tags        []string                `json:"tags"`
	Attributes  map[uuid.UUID][]string  `json:"attributes"` // AttributeID -> Values

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

	SortBy      string                  `json:"sort_by"`    // relevance, price, name, created_at, rating
	SortOrder   string                  `json:"sort_order"` // asc, desc
	Limit       int                     `json:"limit"`
	Offset      int                     `json:"offset"`
}

// EnhancedSearchParams represents enhanced search parameters with dynamic faceting
type EnhancedSearchParams struct {
	FullTextSearchParams
	IncludeFacets    bool                `json:"include_facets"`    // Whether to include facets in response
	FacetFilters     map[string][]string `json:"facet_filters"`     // Current filter selections for facet calculation
	ExcludeFilters   map[string][]string `json:"exclude_filters"`   // Filters to exclude from facet calculation
	DynamicFacets    bool                `json:"dynamic_facets"`    // Whether to calculate dynamic facet counts
}

// SearchFacets represents search facets for filtering
type SearchFacets struct {
	Categories   []CategoryFacet   `json:"categories"`
	Brands       []BrandFacet      `json:"brands"`
	PriceRange   PriceRangeFacet   `json:"price_range"`
	Tags         []TagFacet        `json:"tags"`
	Attributes   []AttributeFacet  `json:"attributes"`
	Status       []StatusFacet     `json:"status"`
	ProductTypes []ProductTypeFacet `json:"product_types"`
	Availability []AvailabilityFacet `json:"availability"`
	Ratings      []RatingFacet     `json:"ratings"`
	Shipping     []ShippingFacet   `json:"shipping"`
}

// DynamicSearchFacets represents dynamic facets with real-time counts
type DynamicSearchFacets struct {
	Categories   []DynamicCategoryFacet   `json:"categories"`
	Brands       []DynamicBrandFacet      `json:"brands"`
	PriceRange   DynamicPriceRangeFacet   `json:"price_range"`
	Tags         []DynamicTagFacet        `json:"tags"`
	Attributes   []DynamicAttributeFacet  `json:"attributes"`
	Status       []DynamicStatusFacet     `json:"status"`
	ProductTypes []DynamicProductTypeFacet `json:"product_types"`
	Availability []DynamicAvailabilityFacet `json:"availability"`
	Ratings      []DynamicRatingFacet     `json:"ratings"`
	Shipping     []DynamicShippingFacet   `json:"shipping"`
	TotalCount   int64                    `json:"total_count"`
}

// Dynamic facet types with selection state
type DynamicCategoryFacet struct {
	CategoryFacet
	IsSelected bool `json:"is_selected"`
	IsDisabled bool `json:"is_disabled"`
}

type DynamicBrandFacet struct {
	BrandFacet
	IsSelected bool `json:"is_selected"`
	IsDisabled bool `json:"is_disabled"`
}

type DynamicTagFacet struct {
	TagFacet
	IsSelected bool `json:"is_selected"`
	IsDisabled bool `json:"is_disabled"`
}

type DynamicAttributeFacet struct {
	AttributeFacet
	IsSelected bool `json:"is_selected"`
	IsDisabled bool `json:"is_disabled"`
}

type DynamicStatusFacet struct {
	StatusFacet
	IsSelected bool `json:"is_selected"`
	IsDisabled bool `json:"is_disabled"`
}

type DynamicProductTypeFacet struct {
	ProductTypeFacet
	IsSelected bool `json:"is_selected"`
	IsDisabled bool `json:"is_disabled"`
}

type DynamicAvailabilityFacet struct {
	AvailabilityFacet
	IsSelected bool `json:"is_selected"`
	IsDisabled bool `json:"is_disabled"`
}

type DynamicRatingFacet struct {
	RatingFacet
	IsSelected bool `json:"is_selected"`
	IsDisabled bool `json:"is_disabled"`
}

type DynamicShippingFacet struct {
	ShippingFacet
	IsSelected bool `json:"is_selected"`
	IsDisabled bool `json:"is_disabled"`
}

type DynamicPriceRangeFacet struct {
	PriceRangeFacet
	SelectedMin *float64 `json:"selected_min"`
	SelectedMax *float64 `json:"selected_max"`
}

// CategoryFacet represents category facet
type CategoryFacet struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ProductCount int64     `json:"product_count"`
}

// BrandFacet represents brand facet
type BrandFacet struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ProductCount int64     `json:"product_count"`
}

// PriceRangeFacet represents price range facet
type PriceRangeFacet struct {
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
	Ranges   []PriceRange `json:"ranges"`
}

// PriceRange represents a price range
type PriceRange struct {
	Min          *float64 `json:"min"`
	Max          *float64 `json:"max"`
	Label        string   `json:"label"`
	ProductCount int64    `json:"product_count"`
}

// TagFacet represents tag facet
type TagFacet struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ProductCount int64     `json:"product_count"`
}

// AttributeFacet represents attribute facet
type AttributeFacet struct {
	ID     uuid.UUID       `json:"id"`
	Name   string          `json:"name"`
	Values []AttributeValue `json:"values"`
}

// AttributeValue represents attribute value
type AttributeValue struct {
	Value        string `json:"value"`
	ProductCount int64  `json:"product_count"`
}

// StatusFacet represents status facet
type StatusFacet struct {
	Status       entities.ProductStatus `json:"status"`
	Label        string                 `json:"label"`
	ProductCount int64                  `json:"product_count"`
}

// ProductTypeFacet represents product type facet
type ProductTypeFacet struct {
	Type         entities.ProductType `json:"type"`
	Label        string               `json:"label"`
	ProductCount int64                `json:"product_count"`
}

// AvailabilityFacet represents availability facet
type AvailabilityFacet struct {
	Status       string `json:"status"`
	Label        string `json:"label"`
	ProductCount int64  `json:"product_count"`
}

// RatingFacet represents rating facet
type RatingFacet struct {
	Rating       int   `json:"rating"`
	Label        string `json:"label"`
	ProductCount int64  `json:"product_count"`
}

// ShippingFacet represents shipping facet
type ShippingFacet struct {
	Type         string `json:"type"`
	Label        string `json:"label"`
	ProductCount int64  `json:"product_count"`
}

// SearchResult represents search result with metadata
type SearchResult struct {
	Products     []*entities.Product `json:"products"`
	Total        int64               `json:"total"`
	Facets       *SearchFacets       `json:"facets"`
	Query        string              `json:"query"`
	SearchTime   time.Duration       `json:"search_time"`
	Suggestions  []string            `json:"suggestions"`
}

// AutocompleteResult represents autocomplete result
type AutocompleteResult struct {
	Products    []*ProductSuggestion `json:"products"`
	Categories  []*CategorySuggestion `json:"categories"`
	Brands      []*BrandSuggestion   `json:"brands"`
	Suggestions []string             `json:"suggestions"`
}


