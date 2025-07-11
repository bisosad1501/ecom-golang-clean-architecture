package repositories

import (
	"context"

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

	// GetByIDs retrieves multiple products by IDs (bulk operation)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*entities.Product, error)

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

	// GetCategoryPath returns the full path from root to the given category
	GetCategoryPath(ctx context.Context, categoryID uuid.UUID) ([]*entities.Category, error)

	// GetProductCountByCategory returns product count for each category (including descendants)
	GetProductCountByCategory(ctx context.Context, categoryID uuid.UUID) (int64, error)
}
