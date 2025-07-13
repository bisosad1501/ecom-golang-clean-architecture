package repositories

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"

	"github.com/google/uuid"
)

// ProductSearchParams represents search parameters for products
type ProductSearchParams struct {
	Query      string
	CategoryID *uuid.UUID
	MinPrice   *float64
	MaxPrice   *float64
	Status     *entities.ProductStatus
	Tags       []string
	SortBy     string // name, price, created_at
	SortOrder  string // asc, desc
	Limit      int
	Offset     int
}

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	// Create creates a new product
	Create(ctx context.Context, product *entities.Product) error

	// GetByID retrieves a product by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)

	// GetByIDForUpdate retrieves a product by ID with row-level locking (SELECT FOR UPDATE)
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*entities.Product, error)

	// GetByIDs retrieves multiple products by IDs (bulk operation)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*entities.Product, error)

	// GetByIDsWithFullDetails retrieves multiple products by IDs with all relations (optimized for bulk operations)
	GetByIDsWithFullDetails(ctx context.Context, ids []uuid.UUID) ([]*entities.Product, error)

	// GetBySKU retrieves a product by SKU
	GetBySKU(ctx context.Context, sku string) (*entities.Product, error)

	// Update updates an existing product
	Update(ctx context.Context, product *entities.Product) error

	// Delete deletes a product by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves products with pagination
	List(ctx context.Context, limit, offset int) ([]*entities.Product, error)

	// Search searches products based on criteria
	Search(ctx context.Context, params ProductSearchParams) ([]*entities.Product, error)

	// Count returns the total number of products
	Count(ctx context.Context) (int64, error)

	// CountByCategory returns the number of products in a category
	CountByCategory(ctx context.Context, categoryID uuid.UUID) (int64, error)

	// GetByCategory retrieves products by category
	GetByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*entities.Product, error)

	// UpdateStock updates product stock
	UpdateStock(ctx context.Context, productID uuid.UUID, stock int) error

	// ExistsBySKU checks if a product exists with the given SKU
	ExistsBySKU(ctx context.Context, sku string) (bool, error)

	// GetFeatured retrieves featured products
	GetFeatured(ctx context.Context, limit int) ([]*entities.Product, error)

	// GetRelated retrieves related products
	GetRelated(ctx context.Context, productID uuid.UUID, limit int) ([]*entities.Product, error)

	// ClearTags removes all tag associations for a product
	ClearTags(ctx context.Context, productID uuid.UUID) error

	// AddTag adds a tag association to a product
	AddTag(ctx context.Context, productID, tagID uuid.UUID) error

	// ReplaceTags replaces all tag associations for a product with new ones
	ReplaceTags(ctx context.Context, productID uuid.UUID, tagIDs []uuid.UUID) error

	// Additional methods for admin dashboard
	CountProducts(ctx context.Context) (int64, error)

	// Brand-related methods
	GetByBrand(ctx context.Context, brandID uuid.UUID, limit, offset int) ([]*entities.Product, error)

	// Slug-related methods
	GetBySlug(ctx context.Context, slug string) (*entities.Product, error)
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
	ExistsBySlugExcludingID(ctx context.Context, slug string, excludeID uuid.UUID) (bool, error)
	GetExistingSlugs(ctx context.Context, prefix string) ([]string, error)

	// Advanced search methods
	SearchAdvanced(ctx context.Context, params AdvancedSearchParams) ([]*entities.Product, error)

	// Search autocomplete and suggestions
	GetSearchSuggestions(ctx context.Context, query string, limit int) (*SearchSuggestions, error)
	GetPopularSearches(ctx context.Context, limit int) ([]string, error)
	RecordSearchQuery(ctx context.Context, query string, userID *uuid.UUID, resultCount int) error
	GetSearchHistory(ctx context.Context, userID uuid.UUID, limit int) ([]string, error)
}

// AdvancedSearchParams represents advanced search parameters
type AdvancedSearchParams struct {
	Query       string
	CategoryID  *uuid.UUID
	BrandID     *uuid.UUID
	MinPrice    *float64
	MaxPrice    *float64
	InStock     *bool
	Featured    *bool
	Visibility  *entities.ProductVisibility
	ProductType *entities.ProductType
	Status      *entities.ProductStatus
	Tags        []string
	Attributes  map[uuid.UUID][]uuid.UUID // AttributeID -> TermIDs
	SortBy      string                    // price, name, created_at, etc.
	SortOrder   string                    // asc, desc
	Limit       int
	Offset      int
}

// BrandRepository defines the interface for brand data access
type BrandRepository interface {
	// Create creates a new brand
	Create(ctx context.Context, brand *entities.Brand) error

	// GetByID retrieves a brand by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Brand, error)

	// GetBySlug retrieves a brand by slug
	GetBySlug(ctx context.Context, slug string) (*entities.Brand, error)

	// Update updates an existing brand
	Update(ctx context.Context, brand *entities.Brand) error

	// Delete deletes a brand by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves brands with pagination
	List(ctx context.Context, limit, offset int) ([]*entities.Brand, error)

	// Search searches brands
	Search(ctx context.Context, query string, limit, offset int) ([]*entities.Brand, error)

	// ExistsBySlug checks if a brand exists with the given slug
	ExistsBySlug(ctx context.Context, slug string) (bool, error)

	// GetActive retrieves active brands
	GetActive(ctx context.Context, limit, offset int) ([]*entities.Brand, error)

	// GetPopularBrands retrieves brands ordered by product count
	GetPopularBrands(ctx context.Context, limit int) ([]*entities.Brand, error)

	// GetBrandsForFiltering retrieves brands for product filtering with counts
	GetBrandsForFiltering(ctx context.Context, categoryID *uuid.UUID) ([]map[string]interface{}, error)

	// CountByStatus counts brands by status
	CountByStatus(ctx context.Context, isActive bool) (int64, error)

	// GetTotal gets total number of brands
	GetTotal(ctx context.Context) (int64, error)
}

// ProductAttributeRepository defines the interface for product attribute data access
type ProductAttributeRepository interface {
	// Create creates a new product attribute
	Create(ctx context.Context, attribute *entities.ProductAttribute) error

	// GetByID retrieves a product attribute by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ProductAttribute, error)

	// GetBySlug retrieves a product attribute by slug
	GetBySlug(ctx context.Context, slug string) (*entities.ProductAttribute, error)

	// Update updates an existing product attribute
	Update(ctx context.Context, attribute *entities.ProductAttribute) error

	// Delete deletes a product attribute by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves product attributes with pagination
	List(ctx context.Context, limit, offset int) ([]*entities.ProductAttribute, error)

	// GetVisible retrieves visible product attributes
	GetVisible(ctx context.Context) ([]*entities.ProductAttribute, error)

	// CreateTerm creates a new attribute term
	CreateTerm(ctx context.Context, term *entities.ProductAttributeTerm) error

	// GetTermsByAttribute retrieves terms for an attribute
	GetTermsByAttribute(ctx context.Context, attributeID uuid.UUID) ([]*entities.ProductAttributeTerm, error)

	// DeleteTerm deletes an attribute term
	DeleteTerm(ctx context.Context, termID uuid.UUID) error
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	// Create creates a new category
	Create(ctx context.Context, category *entities.Category) error

	// GetByID retrieves a category by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Category, error)

	// GetBySlug retrieves a category by slug
	GetBySlug(ctx context.Context, slug string) (*entities.Category, error)

	// Update updates an existing category
	Update(ctx context.Context, category *entities.Category) error

	// Delete deletes a category by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves categories with pagination
	List(ctx context.Context, limit, offset int) ([]*entities.Category, error)

	// GetRootCategories retrieves root categories
	GetRootCategories(ctx context.Context) ([]*entities.Category, error)

	// GetChildren retrieves child categories
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entities.Category, error)

	// Count returns the total number of categories
	Count(ctx context.Context) (int64, error)

	// ExistsBySlug checks if a category exists with the given slug
	ExistsBySlug(ctx context.Context, slug string) (bool, error)

	// GetTree retrieves the category tree
	GetTree(ctx context.Context) ([]*entities.Category, error)

	// GetCategoryTree returns all descendant category IDs for a given category (including itself)
	GetCategoryTree(ctx context.Context, categoryID uuid.UUID) ([]uuid.UUID, error)

	// Bulk operations
	BulkCreate(ctx context.Context, categories []*entities.Category) error
	BulkUpdate(ctx context.Context, categories []*entities.Category) error
	BulkDelete(ctx context.Context, ids []uuid.UUID) error

	// Advanced filtering
	ListWithFilters(ctx context.Context, filters CategoryFilters) ([]*entities.Category, error)
	CountWithFilters(ctx context.Context, filters CategoryFilters) (int64, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*entities.Category, error)

	// Validation helpers
	ValidateHierarchy(ctx context.Context, categoryID, parentID uuid.UUID) error
	GetCategoryPath(ctx context.Context, categoryID uuid.UUID) ([]*entities.Category, error)
	GetProductCount(ctx context.Context, categoryID uuid.UUID, includeSubcategories bool) (int64, error)

	// Tree operations
	MoveCategory(ctx context.Context, categoryID, newParentID uuid.UUID) error
	ReorderCategories(ctx context.Context, reorderRequests []CategoryReorderRequest) error
	GetCategoryDepth(ctx context.Context, categoryID uuid.UUID) (int, error)
	GetMaxDepth(ctx context.Context) (int, error)
	ValidateTreeIntegrity(ctx context.Context) error
	RebuildCategoryPaths(ctx context.Context) error

	// Analytics and statistics
	GetCategoryAnalytics(ctx context.Context, categoryID uuid.UUID, timeRange string) (*CategoryAnalytics, error)
	GetTopCategories(ctx context.Context, limit int, sortBy string) ([]*CategoryStats, error)
	GetCategoryPerformanceMetrics(ctx context.Context, categoryID uuid.UUID) (*CategoryPerformanceMetrics, error)
	GetCategorySalesStats(ctx context.Context, categoryID uuid.UUID, timeRange string) (*CategorySalesStats, error)



	// GetProductCountByCategory returns product count for each category (including descendants)
	GetProductCountByCategory(ctx context.Context, categoryID uuid.UUID) (int64, error)

	// Optimized bulk operations
	GetWithProductsOptimized(ctx context.Context, id uuid.UUID, limit, offset int) (*entities.Category, []*entities.Product, error)
	GetCategoriesWithProductCount(ctx context.Context) ([]*entities.Category, map[uuid.UUID]int64, error)
}

// CategoryFilters represents filters for category queries
type CategoryFilters struct {
	Search    string     `json:"search"`
	ParentID  *uuid.UUID `json:"parent_id"`
	IsActive  *bool      `json:"is_active"`
	Level     *int       `json:"level"`
	HasParent *bool      `json:"has_parent"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
	SortBy    string     `json:"sort_by"`    // name, created_at, sort_order
	SortOrder string     `json:"sort_order"` // asc, desc
}

// CategoryReorderRequest represents a category reorder request
type CategoryReorderRequest struct {
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
	SortOrder  int       `json:"sort_order" validate:"required"`
}

// CategoryAnalytics represents comprehensive category analytics
type CategoryAnalytics struct {
	CategoryID       uuid.UUID                 `json:"category_id"`
	CategoryName     string                    `json:"category_name"`
	ProductCount     int64                     `json:"product_count"`
	ActiveProducts   int64                     `json:"active_products"`
	InactiveProducts int64                     `json:"inactive_products"`
	TotalViews       int64                     `json:"total_views"`
	TotalSales       int64                     `json:"total_sales"`
	Revenue          float64                   `json:"revenue"`
	AveragePrice     float64                   `json:"average_price"`
	ConversionRate   float64                   `json:"conversion_rate"`
	TopProducts      []ProductPerformance      `json:"top_products"`
	SalesHistory     []SalesDataPoint          `json:"sales_history"`
	ViewsHistory     []ViewsDataPoint          `json:"views_history"`
}

// CategoryStats represents category statistics for ranking
type CategoryStats struct {
	CategoryID     uuid.UUID `json:"category_id"`
	CategoryName   string    `json:"category_name"`
	ProductCount   int64     `json:"product_count"`
	TotalSales     int64     `json:"total_sales"`
	Revenue        float64   `json:"revenue"`
	AverageRating  float64   `json:"average_rating"`
	ConversionRate float64   `json:"conversion_rate"`
	GrowthRate     float64   `json:"growth_rate"`
}

// CategoryPerformanceMetrics represents detailed performance metrics
type CategoryPerformanceMetrics struct {
	CategoryID          uuid.UUID `json:"category_id"`
	CategoryName        string    `json:"category_name"`
	ProductCount        int64     `json:"product_count"`
	ActiveProductCount  int64     `json:"active_product_count"`
	AverageProductPrice float64   `json:"average_product_price"`
	TotalInventoryValue float64   `json:"total_inventory_value"`
	LowStockProducts    int64     `json:"low_stock_products"`
	OutOfStockProducts  int64     `json:"out_of_stock_products"`
	AverageRating       float64   `json:"average_rating"`
	TotalReviews        int64     `json:"total_reviews"`
	PopularityScore     float64   `json:"popularity_score"`
}

// CategorySalesStats represents sales statistics for a category
type CategorySalesStats struct {
	CategoryID      uuid.UUID        `json:"category_id"`
	CategoryName    string           `json:"category_name"`
	TimeRange       string           `json:"time_range"`
	TotalSales      int64            `json:"total_sales"`
	TotalRevenue    float64          `json:"total_revenue"`
	AverageOrderValue float64        `json:"average_order_value"`
	TopSellingProducts []ProductSales `json:"top_selling_products"`
	SalesByPeriod   []PeriodSales    `json:"sales_by_period"`
	GrowthMetrics   GrowthMetrics    `json:"growth_metrics"`
}

// ProductPerformance represents product performance data
type ProductPerformance struct {
	ProductID    uuid.UUID `json:"product_id"`
	ProductName  string    `json:"product_name"`
	SKU          string    `json:"sku"`
	Sales        int64     `json:"sales"`
	Revenue      float64   `json:"revenue"`
	Views        int64     `json:"views"`
	Rating       float64   `json:"rating"`
	ReviewCount  int64     `json:"review_count"`
}

// ProductSales represents product sales data
type ProductSales struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	SKU         string    `json:"sku"`
	Quantity    int64     `json:"quantity"`
	Revenue     float64   `json:"revenue"`
}

// SalesDataPoint represents a sales data point over time
type SalesDataPoint struct {
	Date     string  `json:"date"`
	Sales    int64   `json:"sales"`
	Revenue  float64 `json:"revenue"`
}

// ViewsDataPoint represents a views data point over time
type ViewsDataPoint struct {
	Date  string `json:"date"`
	Views int64  `json:"views"`
}

// PeriodSales represents sales data for a specific period
type PeriodSales struct {
	Period   string  `json:"period"`
	Sales    int64   `json:"sales"`
	Revenue  float64 `json:"revenue"`
	Orders   int64   `json:"orders"`
}

// GrowthMetrics represents growth metrics
type GrowthMetrics struct {
	SalesGrowth   float64 `json:"sales_growth"`
	RevenueGrowth float64 `json:"revenue_growth"`
	OrderGrowth   float64 `json:"order_growth"`
}

// SearchSuggestions represents search suggestions response
type SearchSuggestions struct {
	Products    []ProductSuggestion  `json:"products"`
	Categories  []CategorySuggestion `json:"categories"`
	Brands      []BrandSuggestion    `json:"brands"`
	Popular     []string             `json:"popular"`
	Corrections []string             `json:"corrections"`
}

// ProductSuggestion represents a product suggestion
type ProductSuggestion struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	SKU         string    `json:"sku"`
	Price       float64   `json:"price"`
	Image       string    `json:"image"`
	CategoryID  uuid.UUID `json:"category_id"`
	Category    string    `json:"category"`
	Relevance   float64   `json:"relevance"`
}

// CategorySuggestion represents a category suggestion
type CategorySuggestion struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ProductCount int64     `json:"product_count"`
	Relevance    float64   `json:"relevance"`
}

// BrandSuggestion represents a brand suggestion
type BrandSuggestion struct {
	Name         string  `json:"name"`
	ProductCount int64   `json:"product_count"`
	Relevance    float64 `json:"relevance"`
}

// SearchQuery represents a search query record
type SearchQuery struct {
	ID          uuid.UUID  `json:"id"`
	Query       string     `json:"query"`
	UserID      *uuid.UUID `json:"user_id"`
	ResultCount int        `json:"result_count"`
	CreatedAt   time.Time  `json:"created_at"`
}
