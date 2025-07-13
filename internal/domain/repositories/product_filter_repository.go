package repositories

import (
	"context"
	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// AdvancedFilterParams represents advanced filtering parameters
type AdvancedFilterParams struct {
	// Basic filters
	Query       string     `json:"query"`
	CategoryIDs []uuid.UUID `json:"category_ids"`
	BrandIDs    []uuid.UUID `json:"brand_ids"`
	MinPrice    *float64   `json:"min_price"`
	MaxPrice    *float64   `json:"max_price"`
	MinRating   *float64   `json:"min_rating"`
	MaxRating   *float64   `json:"max_rating"`
	
	// Stock and availability
	InStock     *bool `json:"in_stock"`
	LowStock    *bool `json:"low_stock"`
	OnSale      *bool `json:"on_sale"`
	Featured    *bool `json:"featured"`
	
	// Product properties
	ProductTypes []entities.ProductType `json:"product_types"`
	StockStatus  []entities.StockStatus `json:"stock_status"`
	Visibility   []entities.ProductVisibility `json:"visibility"`
	
	// Custom attributes
	Attributes map[uuid.UUID][]string `json:"attributes"` // AttributeID -> Values
	
	// Date filters
	CreatedAfter  *string `json:"created_after"`
	CreatedBefore *string `json:"created_before"`
	UpdatedAfter  *string `json:"updated_after"`
	UpdatedBefore *string `json:"updated_before"`
	
	// Advanced options
	Tags         []string `json:"tags"`
	HasImages    *bool    `json:"has_images"`
	HasVariants  *bool    `json:"has_variants"`
	HasReviews   *bool    `json:"has_reviews"`
	
	// Sorting and pagination
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	
	// Filter options
	IncludeFacets bool `json:"include_facets"`
	FacetLimit    int  `json:"facet_limit"`
}

// FilterFacets represents available filter facets (reusing existing types)
type FilterFacets struct {
	Categories []FilterCategoryFacet   `json:"categories"`
	Brands     []FilterBrandFacet      `json:"brands"`
	Attributes []FilterAttributeFacet  `json:"attributes"`
	PriceRange FilterPriceRangeFacet   `json:"price_range"`
	Rating     FilterRatingFacet       `json:"rating"`
	Stock      FilterStockFacet        `json:"stock"`
	Tags       []FilterTagFacet        `json:"tags"`
}

// FilterCategoryFacet represents category filter facet
type FilterCategoryFacet struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Count    int       `json:"count"`
	Children []FilterCategoryFacet `json:"children,omitempty"`
}

// FilterBrandFacet represents brand filter facet
type FilterBrandFacet struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Slug  string    `json:"slug"`
	Count int       `json:"count"`
	Logo  string    `json:"logo,omitempty"`
}

// FilterAttributeFacet represents attribute filter facet
type FilterAttributeFacet struct {
	ID      uuid.UUID           `json:"id"`
	Name    string              `json:"name"`
	Slug    string              `json:"slug"`
	Type    string              `json:"type"`
	Terms   []FilterAttributeTermFacet `json:"terms"`
}

// FilterAttributeTermFacet represents attribute term facet
type FilterAttributeTermFacet struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Value string    `json:"value"`
	Count int       `json:"count"`
	Color string    `json:"color,omitempty"`
	Image string    `json:"image,omitempty"`
}

// FilterPriceRangeFacet represents price range facet
type FilterPriceRangeFacet struct {
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Ranges []FilterPriceRange `json:"ranges"`
}

// FilterPriceRange represents a price range option
type FilterPriceRange struct {
	Min   *float64 `json:"min"`
	Max   *float64 `json:"max"`
	Label string   `json:"label"`
	Count int      `json:"count"`
}

// FilterRatingFacet represents rating filter facet
type FilterRatingFacet struct {
	Ranges []FilterRatingRange `json:"ranges"`
}

// FilterRatingRange represents a rating range option
type FilterRatingRange struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Label string  `json:"label"`
	Count int     `json:"count"`
}

// FilterStockFacet represents stock filter facet
type FilterStockFacet struct {
	InStock   int `json:"in_stock"`
	LowStock  int `json:"low_stock"`
	OutStock  int `json:"out_of_stock"`
	OnSale    int `json:"on_sale"`
	Featured  int `json:"featured"`
}

// FilterTagFacet represents tag filter facet
type FilterTagFacet struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// FilteredProductResult represents the result of filtered product search
type FilteredProductResult struct {
	Products []*entities.Product `json:"products"`
	Facets   *FilterFacets       `json:"facets,omitempty"`
	Total    int64               `json:"total"`
}

// ProductFilterRepository defines the interface for advanced product filtering
type ProductFilterRepository interface {
	// Advanced filtering
	FilterProducts(ctx context.Context, params AdvancedFilterParams) (*FilteredProductResult, error)
	GetFilterFacets(ctx context.Context, categoryID *uuid.UUID) (*FilterFacets, error)
	GetDynamicFilters(ctx context.Context, params AdvancedFilterParams) (*FilterFacets, error)
	
	// Filter sets management
	SaveFilterSet(ctx context.Context, filterSet *entities.FilterSet) error
	GetFilterSet(ctx context.Context, id uuid.UUID) (*entities.FilterSet, error)
	GetUserFilterSets(ctx context.Context, userID uuid.UUID) ([]*entities.FilterSet, error)
	GetSessionFilterSets(ctx context.Context, sessionID string) ([]*entities.FilterSet, error)
	UpdateFilterSet(ctx context.Context, filterSet *entities.FilterSet) error
	DeleteFilterSet(ctx context.Context, id uuid.UUID) error
	
	// Filter analytics
	TrackFilterUsage(ctx context.Context, usage *entities.FilterUsage) error
	GetFilterAnalytics(ctx context.Context, days int) (map[string]interface{}, error)
	GetPopularFilters(ctx context.Context, limit int) ([]*entities.FilterUsage, error)
	
	// Filter options management
	UpdateFilterOptions(ctx context.Context, categoryID *uuid.UUID) error
	GetFilterOptions(ctx context.Context, categoryID *uuid.UUID) ([]*entities.ProductFilterOption, error)
	
	// Attribute-based filtering
	GetAttributeFilters(ctx context.Context, categoryID *uuid.UUID) ([]*entities.ProductAttribute, error)
	GetAttributeTerms(ctx context.Context, attributeID uuid.UUID, categoryID *uuid.UUID) ([]*entities.ProductAttributeTerm, error)
	
	// Filter suggestions
	GetFilterSuggestions(ctx context.Context, query string, limit int) ([]string, error)
	GetRelatedFilters(ctx context.Context, currentFilters AdvancedFilterParams) ([]string, error)
}
