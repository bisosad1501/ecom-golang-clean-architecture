package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/pkg/utils"

	"github.com/google/uuid"
)

// Request structs
type CreateProductRequest struct {
	Name             string `json:"name" validate:"required"`
	Description      string `json:"description" validate:"required"`
	ShortDescription string `json:"short_description"`
	SKU              string `json:"sku" validate:"required"`

	// SEO and Metadata
	Slug            string                     `json:"slug" validate:"required"`
	MetaTitle       string                     `json:"meta_title"`
	MetaDescription string                     `json:"meta_description"`
	Keywords        string                     `json:"keywords"`
	Featured        bool                       `json:"featured"`
	Visibility      entities.ProductVisibility `json:"visibility"`

	// Pricing
	Price        float64  `json:"price" validate:"required,gt=0"`
	ComparePrice *float64 `json:"compare_price" validate:"omitempty,gt=0"`
	CostPrice    *float64 `json:"cost_price" validate:"omitempty,gt=0"`

	// Sale Pricing
	SalePrice     *float64   `json:"sale_price" validate:"omitempty,gt=0"`
	SaleStartDate *time.Time `json:"sale_start_date"`
	SaleEndDate   *time.Time `json:"sale_end_date"`

	// Inventory
	Stock             int  `json:"stock" validate:"required,min=0"`
	LowStockThreshold int  `json:"low_stock_threshold"`
	TrackQuantity     bool `json:"track_quantity"`
	AllowBackorder    bool `json:"allow_backorder"`

	// Physical Properties
	Weight     *float64           `json:"weight" validate:"omitempty,gt=0"`
	Dimensions *DimensionsRequest `json:"dimensions"`

	// Shipping and Tax
	RequiresShipping bool   `json:"requires_shipping"`
	ShippingClass    string `json:"shipping_class"`
	TaxClass         string `json:"tax_class"`
	CountryOfOrigin  string `json:"country_of_origin"`

	// Categorization
	CategoryID uuid.UUID  `json:"category_id" validate:"required"`
	BrandID    *uuid.UUID `json:"brand_id"`

	// Content
	Images     []ProductImageRequest     `json:"images"`
	Tags       []string                  `json:"tags"`
	Attributes []ProductAttributeRequest `json:"attributes"`
	Variants   []ProductVariantRequest   `json:"variants"`

	// Status and Type
	Status      entities.ProductStatus `json:"status"`
	ProductType entities.ProductType   `json:"product_type"`
	IsDigital   bool                   `json:"is_digital"`
}

type GetProductsRequest struct {
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}

// GetProductsResponse represents paginated products response
type GetProductsResponse struct {
	Products   []*ProductResponse `json:"products"`
	Pagination *PaginationInfo    `json:"pagination"`
}

// FeaturedProductsPaginatedResponse represents paginated featured products
type FeaturedProductsPaginatedResponse struct {
	Products   []*ProductResponse `json:"products"`
	Pagination *PaginationInfo    `json:"pagination"`
}

// TrendingProductsPaginatedResponse represents paginated trending products
type TrendingProductsPaginatedResponse struct {
	Products   []*ProductResponse `json:"products"`
	Pagination *PaginationInfo    `json:"pagination"`
}

// RelatedProductsPaginatedResponse represents paginated related products
type RelatedProductsPaginatedResponse struct {
	Products   []*ProductResponse `json:"products"`
	Pagination *PaginationInfo    `json:"pagination"`
	ProductID  uuid.UUID          `json:"product_id"`
}

type SearchProductsRequest struct {
	Query      string                  `json:"query"`
	CategoryID *uuid.UUID              `json:"category_id"`
	MinPrice   *float64                `json:"min_price" validate:"omitempty,gt=0"`
	MaxPrice   *float64                `json:"max_price" validate:"omitempty,gt=0"`
	Status     *entities.ProductStatus `json:"status"`
	Tags       []string                `json:"tags"`
	SortBy     string                  `json:"sort_by"`
	SortOrder  string                  `json:"sort_order"`
	Limit      int                     `json:"limit" validate:"min=1,max=100"`
	Offset     int                     `json:"offset" validate:"min=0"`
}

type DimensionsRequest struct {
	Length float64 `json:"length" validate:"required,gt=0"`
	Width  float64 `json:"width" validate:"required,gt=0"`
	Height float64 `json:"height" validate:"required,gt=0"`
}

type ProductImageRequest struct {
	URL      string `json:"url" validate:"required,url"`
	AltText  string `json:"alt_text"`
	Position int    `json:"position"`
}

type ProductAttributeRequest struct {
	AttributeID uuid.UUID  `json:"attribute_id" validate:"required"`
	TermID      *uuid.UUID `json:"term_id"`
	Value       string     `json:"value"`
	Position    int        `json:"position"`
}

type ProductVariantRequest struct {
	SKU          string                           `json:"sku" validate:"required"`
	Price        float64                          `json:"price" validate:"required,gt=0"`
	ComparePrice *float64                         `json:"compare_price" validate:"omitempty,gt=0"`
	CostPrice    *float64                         `json:"cost_price" validate:"omitempty,gt=0"`
	Stock        int                              `json:"stock" validate:"min=0"`
	Weight       *float64                         `json:"weight" validate:"omitempty,gt=0"`
	Dimensions   *DimensionsRequest               `json:"dimensions"`
	Image        string                           `json:"image"`
	Position     int                              `json:"position"`
	IsActive     bool                             `json:"is_active"`
	Attributes   []ProductVariantAttributeRequest `json:"attributes"`
}

type ProductVariantAttributeRequest struct {
	AttributeID uuid.UUID `json:"attribute_id" validate:"required"`
	TermID      uuid.UUID `json:"term_id" validate:"required"`
}

// SearchSuggestionsRequest represents search suggestions request
type SearchSuggestionsRequest struct {
	Query string `json:"query" validate:"required,min=1"`
	Limit int    `json:"limit" validate:"min=1,max=50"`
}

// SearchSuggestionsResponse represents search suggestions response
type SearchSuggestionsResponse struct {
	Suggestions *repositories.SearchSuggestions `json:"suggestions"`
}

// PopularSearchesResponse represents popular searches response
type PopularSearchesResponse struct {
	Searches []string `json:"searches"`
}

// SearchHistoryResponse represents search history response
type SearchHistoryResponse struct {
	History []string `json:"history"`
}

// Response structs are defined in types.go

// ProductUseCase defines product use cases
type ProductUseCase interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*ProductResponse, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*ProductResponse, error)
	PatchProduct(ctx context.Context, id uuid.UUID, req PatchProductRequest) (*ProductResponse, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	GetProducts(ctx context.Context, req GetProductsRequest) (*GetProductsResponse, error)
	SearchProducts(ctx context.Context, req SearchProductsRequest) ([]*ProductResponse, error)
	SearchProductsPaginated(ctx context.Context, req SearchProductsRequest) (*GetProductsResponse, error)
	GetProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) (*GetProductsResponse, error)
	UpdateStock(ctx context.Context, productID uuid.UUID, stock int) error

	// Search autocomplete and suggestions
	GetSearchSuggestions(ctx context.Context, req SearchSuggestionsRequest) (*SearchSuggestionsResponse, error)
	GetPopularSearches(ctx context.Context, limit int) (*PopularSearchesResponse, error)
	GetSearchHistory(ctx context.Context, userID uuid.UUID, limit int) (*SearchHistoryResponse, error)

	// Paginated product methods
	GetFeaturedProductsPaginated(ctx context.Context, page, limit int) (*FeaturedProductsPaginatedResponse, error)
	GetTrendingProductsPaginated(ctx context.Context, page, limit int) (*TrendingProductsPaginatedResponse, error)
	GetRelatedProductsPaginated(ctx context.Context, productID uuid.UUID, page, limit int) (*RelatedProductsPaginatedResponse, error)
}

type productUseCase struct {
	productRepo         repositories.ProductRepository
	categoryRepo        repositories.CategoryRepository
	productCategoryRepo repositories.ProductCategoryRepository
	tagRepo             repositories.TagRepository
	imageRepo           repositories.ImageRepository
	cartRepo            repositories.CartRepository
	inventoryRepo       repositories.InventoryRepository
	warehouseRepo       repositories.WarehouseRepository
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(
	productRepo repositories.ProductRepository,
	categoryRepo repositories.CategoryRepository,
	productCategoryRepo repositories.ProductCategoryRepository,
	tagRepo repositories.TagRepository,
	imageRepo repositories.ImageRepository,
	cartRepo repositories.CartRepository,
	inventoryRepo repositories.InventoryRepository,
	warehouseRepo repositories.WarehouseRepository,
) ProductUseCase {
	return &productUseCase{
		productRepo:         productRepo,
		categoryRepo:        categoryRepo,
		productCategoryRepo: productCategoryRepo,
		tagRepo:             tagRepo,
		imageRepo:           imageRepo,
		cartRepo:            cartRepo,
		inventoryRepo:       inventoryRepo,
		warehouseRepo:       warehouseRepo,
	}
}

type UpdateProductRequest struct {
	Name             *string `json:"name"`
	Description      *string `json:"description"`
	ShortDescription *string `json:"short_description"`

	// SEO and Metadata
	Slug            *string                     `json:"slug"`
	MetaTitle       *string                     `json:"meta_title"`
	MetaDescription *string                     `json:"meta_description"`
	Keywords        *string                     `json:"keywords"`
	Featured        *bool                       `json:"featured"`
	Visibility      *entities.ProductVisibility `json:"visibility"`

	// Pricing
	Price        *float64 `json:"price" validate:"omitempty,gt=0"`
	ComparePrice *float64 `json:"compare_price" validate:"omitempty,gt=0"`
	CostPrice    *float64 `json:"cost_price" validate:"omitempty,gt=0"`

	// Sale Pricing
	SalePrice     *float64   `json:"sale_price" validate:"omitempty,gt=0"`
	SaleStartDate *time.Time `json:"sale_start_date"`
	SaleEndDate   *time.Time `json:"sale_end_date"`

	// Inventory
	Stock             *int  `json:"stock" validate:"omitempty,min=0"`
	LowStockThreshold *int  `json:"low_stock_threshold"`
	TrackQuantity     *bool `json:"track_quantity"`
	AllowBackorder    *bool `json:"allow_backorder"`

	// Physical Properties
	Weight     *float64           `json:"weight" validate:"omitempty,gt=0"`
	Dimensions *DimensionsRequest `json:"dimensions"`

	// Shipping and Tax
	RequiresShipping *bool   `json:"requires_shipping"`
	ShippingClass    *string `json:"shipping_class"`
	TaxClass         *string `json:"tax_class"`
	CountryOfOrigin  *string `json:"country_of_origin"`

	// Categorization
	CategoryID *uuid.UUID `json:"category_id"`
	BrandID    *uuid.UUID `json:"brand_id"`

	// Content
	Images     []ProductImageRequest     `json:"images"`     // For PUT: replace all images
	Tags       []string                  `json:"tags"`       // For PUT: replace all tags
	Attributes []ProductAttributeRequest `json:"attributes"` // For PUT: replace all attributes
	Variants   []ProductVariantRequest   `json:"variants"`   // For PUT: replace all variants

	// Status and Type
	Status      *entities.ProductStatus `json:"status"`
	ProductType *entities.ProductType   `json:"product_type"`
	IsDigital   *bool                   `json:"is_digital"`
}

// PatchProductRequest for PATCH operations - only updates provided fields
type PatchProductRequest struct {
	Name             *string `json:"name"`
	Description      *string `json:"description"`
	ShortDescription *string `json:"short_description"`

	// SEO and Metadata
	Slug            *string                     `json:"slug"`
	MetaTitle       *string                     `json:"meta_title"`
	MetaDescription *string                     `json:"meta_description"`
	Keywords        *string                     `json:"keywords"`
	Featured        *bool                       `json:"featured"`
	Visibility      *entities.ProductVisibility `json:"visibility"`

	// Pricing
	Price        *float64 `json:"price" validate:"omitempty,gt=0"`
	ComparePrice *float64 `json:"compare_price" validate:"omitempty,gt=0"`
	CostPrice    *float64 `json:"cost_price" validate:"omitempty,gt=0"`

	// Sale Pricing
	SalePrice     *float64   `json:"sale_price" validate:"omitempty,gt=0"`
	SaleStartDate *time.Time `json:"sale_start_date"`
	SaleEndDate   *time.Time `json:"sale_end_date"`

	// Inventory
	Stock             *int  `json:"stock" validate:"omitempty,min=0"`
	LowStockThreshold *int  `json:"low_stock_threshold"`
	TrackQuantity     *bool `json:"track_quantity"`
	AllowBackorder    *bool `json:"allow_backorder"`

	// Physical Properties
	Weight     *float64           `json:"weight" validate:"omitempty,gt=0"`
	Dimensions *DimensionsRequest `json:"dimensions"`

	// Shipping and Tax
	RequiresShipping *bool   `json:"requires_shipping"`
	ShippingClass    *string `json:"shipping_class"`
	TaxClass         *string `json:"tax_class"`
	CountryOfOrigin  *string `json:"country_of_origin"`

	// Categorization
	CategoryID *uuid.UUID `json:"category_id"`
	BrandID    *uuid.UUID `json:"brand_id"`

	// Content
	Images     *[]ProductImageRequest     `json:"images"`     // For PATCH: nil = no change, empty = clear all, values = replace
	Tags       *[]string                  `json:"tags"`       // For PATCH: nil = no change, empty = clear all, values = replace
	Attributes *[]ProductAttributeRequest `json:"attributes"` // For PATCH: nil = no change, empty = clear all, values = replace
	Variants   *[]ProductVariantRequest   `json:"variants"`   // For PATCH: nil = no change, empty = clear all, values = replace

	// Status and Type
	Status      *entities.ProductStatus `json:"status"`
	ProductType *entities.ProductType   `json:"product_type"`
	IsDigital   *bool                   `json:"is_digital"`
}

// CreateProduct creates a new product
func (uc *productUseCase) CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error) {
	// Check if SKU already exists
	exists, err := uc.productRepo.ExistsBySKU(ctx, req.SKU)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entities.ErrConflict
	}

	// Verify category exists
	_, err = uc.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, entities.ErrCategoryNotFound
	}

	// Generate unique slug
	slug := req.Slug
	if slug == "" {
		slug = utils.GenerateSlug(req.Name)
	}

	// Validate slug format
	if err := utils.ValidateSlug(slug); err != nil {
		return nil, fmt.Errorf("invalid slug: %w", err)
	}

	// Ensure slug is unique
	baseSlug := slug
	existingSlugs, err := uc.productRepo.GetExistingSlugs(ctx, baseSlug)
	if err != nil {
		return nil, err
	}
	slug = utils.GenerateUniqueSlug(baseSlug, existingSlugs)

	// Create product
	product := &entities.Product{
		ID:               uuid.New(),
		Name:             req.Name,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		SKU:              req.SKU,

		// SEO and Metadata
		Slug:            slug,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		Keywords:        req.Keywords,
		Featured:        req.Featured,
		Visibility:      req.Visibility,

		// Pricing
		Price:        req.Price,
		ComparePrice: req.ComparePrice,
		CostPrice:    req.CostPrice,

		// Sale Pricing
		SalePrice:     req.SalePrice,
		SaleStartDate: req.SaleStartDate,
		SaleEndDate:   req.SaleEndDate,

		// Inventory
		Stock:             req.Stock,
		LowStockThreshold: req.LowStockThreshold,
		TrackQuantity:     req.TrackQuantity,
		AllowBackorder:    req.AllowBackorder,

		// Physical Properties
		Weight: req.Weight,

		// Shipping and Tax
		RequiresShipping: req.RequiresShipping,
		ShippingClass:    req.ShippingClass,
		TaxClass:         req.TaxClass,
		CountryOfOrigin:  req.CountryOfOrigin,

		// Categorization (CategoryID removed - using ProductCategory many-to-many)
		BrandID:    req.BrandID,

		// Status and Type
		Status:      req.Status,
		ProductType: req.ProductType,
		IsDigital:   req.IsDigital,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set default values if not provided
	if product.Visibility == "" {
		product.Visibility = entities.ProductVisibilityVisible
	}
	if product.ProductType == "" {
		product.ProductType = entities.ProductTypeSimple
	}
	if product.Status == "" {
		product.Status = entities.ProductStatusDraft
	}
	if product.LowStockThreshold == 0 {
		product.LowStockThreshold = 5
	}
	if product.TaxClass == "" {
		product.TaxClass = "standard"
	}

	if req.Dimensions != nil {
		product.Dimensions = &entities.Dimensions{
			Length: req.Dimensions.Length,
			Width:  req.Dimensions.Width,
			Height: req.Dimensions.Height,
		}
	}

	// Update stock status based on current stock
	product.UpdateStockStatus()

	// Create product first
	if err := uc.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	// Assign category using ProductCategory many-to-many (as primary category)
	if req.CategoryID != uuid.Nil {
		if err := uc.productCategoryRepo.AssignProductToCategory(ctx, product.ID, req.CategoryID, true); err != nil {
			return nil, fmt.Errorf("failed to assign category: %w", err)
		}
	}

	// Create initial inventory record for default warehouse
	if err := uc.createInitialInventory(ctx, product); err != nil {
		// Log error but don't fail product creation
		// In production, you might want to handle this differently
		fmt.Printf("Warning: Failed to create initial inventory for product %s: %v\n", product.SKU, err)
	}

	// Handle tags if provided
	if len(req.Tags) > 0 {
		if err := uc.replaceProductTags(ctx, product.ID, req.Tags); err != nil {
			return nil, err
		}
	}

	// Handle images if provided
	if len(req.Images) > 0 {
		if err := uc.replaceProductImages(ctx, product.ID, req.Images); err != nil {
			return nil, err
		}
	}

	// Handle attributes if provided
	if len(req.Attributes) > 0 {
		if err := uc.replaceProductAttributes(ctx, product.ID, req.Attributes); err != nil {
			return nil, err
		}
	}

	// Handle variants if provided
	if len(req.Variants) > 0 {
		if err := uc.replaceProductVariants(ctx, product.ID, req.Variants); err != nil {
			return nil, err
		}
	}

	// Reload and return
	updatedProduct, err := uc.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	return uc.toProductResponse(updatedProduct), nil
}

// GetProduct gets a product by ID
func (uc *productUseCase) GetProduct(ctx context.Context, id uuid.UUID) (*ProductResponse, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrProductNotFound
	}

	return uc.toProductResponse(product), nil
}

// UpdateProduct updates a product with improved business logic
func (uc *productUseCase) UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*ProductResponse, error) {
	// Get existing product
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrProductNotFound
	}

	// Track what needs to be updated
	hasChanges := false

	// Update basic fields only if they are provided
	if req.Name != nil {
		if *req.Name == "" {
			return nil, fmt.Errorf("name cannot be empty")
		}
		product.Name = *req.Name
		hasChanges = true

		// If name changed and no explicit slug provided, regenerate slug
		if req.Slug == nil && product.Slug == "" {
			newSlug := utils.GenerateSlug(*req.Name)
			if err := utils.ValidateSlug(newSlug); err == nil {
				// Ensure slug is unique
				baseSlug := newSlug
				existingSlugs, err := uc.productRepo.GetExistingSlugs(ctx, baseSlug)
				if err == nil {
					product.Slug = utils.GenerateUniqueSlug(baseSlug, existingSlugs)
				}
			}
		}
	}

	if req.Description != nil {
		product.Description = *req.Description
		hasChanges = true
	}

	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, fmt.Errorf("price must be greater than 0")
		}
		product.Price = *req.Price
		hasChanges = true
	}

	if req.ComparePrice != nil {
		if *req.ComparePrice <= 0 {
			return nil, fmt.Errorf("compare price must be greater than 0")
		}
		product.ComparePrice = req.ComparePrice
		hasChanges = true
	}

	if req.CostPrice != nil {
		if *req.CostPrice < 0 {
			return nil, fmt.Errorf("cost price cannot be negative")
		}
		product.CostPrice = req.CostPrice
		hasChanges = true
	}

	if req.Stock != nil {
		if *req.Stock < 0 {
			return nil, fmt.Errorf("stock cannot be negative")
		}
		product.Stock = *req.Stock
		hasChanges = true
	}

	if req.Weight != nil {
		if *req.Weight <= 0 {
			return nil, fmt.Errorf("weight must be greater than 0")
		}
		product.Weight = req.Weight
		hasChanges = true
	}

	if req.CategoryID != nil {
		// Verify category exists
		_, err := uc.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, entities.ErrCategoryNotFound
		}
		// Use ProductCategory many-to-many as single source of truth
		// Remove all existing categories and assign new primary category
		if err := uc.productCategoryRepo.RemoveProductFromAllCategories(ctx, product.ID); err != nil {
			return nil, fmt.Errorf("failed to remove existing categories: %w", err)
		}
		if err := uc.productCategoryRepo.AssignProductToCategory(ctx, product.ID, *req.CategoryID, true); err != nil {
			return nil, fmt.Errorf("failed to assign primary category: %w", err)
		}
		hasChanges = true
	}

	if req.Status != nil {
		product.Status = *req.Status
		hasChanges = true
	}

	if req.IsDigital != nil {
		product.IsDigital = *req.IsDigital
		hasChanges = true
	}

	if req.Dimensions != nil {
		if req.Dimensions.Length <= 0 || req.Dimensions.Width <= 0 || req.Dimensions.Height <= 0 {
			return nil, fmt.Errorf("all dimensions must be greater than 0")
		}
		product.Dimensions = &entities.Dimensions{
			Length: req.Dimensions.Length,
			Width:  req.Dimensions.Width,
			Height: req.Dimensions.Height,
		}
		hasChanges = true
	}

	// Handle new SEO and Metadata fields
	if req.ShortDescription != nil {
		product.ShortDescription = *req.ShortDescription
		hasChanges = true
	}

	if req.Slug != nil {
		slug := *req.Slug
		if slug == "" {
			// Generate slug from name if not provided
			slug = utils.GenerateSlug(product.Name)
		}

		// Validate slug format
		if err := utils.ValidateSlug(slug); err != nil {
			return nil, fmt.Errorf("invalid slug: %w", err)
		}

		// Ensure slug is unique (excluding current product)
		exists, err := uc.productRepo.ExistsBySlugExcludingID(ctx, slug, product.ID)
		if err != nil {
			return nil, err
		}
		if exists {
			// Generate unique slug if conflicts
			baseSlug := slug
			existingSlugs, err := uc.productRepo.GetExistingSlugs(ctx, baseSlug)
			if err != nil {
				return nil, err
			}
			slug = utils.GenerateUniqueSlug(baseSlug, existingSlugs)
		}

		product.Slug = slug
		hasChanges = true
	}

	if req.MetaTitle != nil {
		product.MetaTitle = *req.MetaTitle
		hasChanges = true
	}

	if req.MetaDescription != nil {
		product.MetaDescription = *req.MetaDescription
		hasChanges = true
	}

	if req.Keywords != nil {
		product.Keywords = *req.Keywords
		hasChanges = true
	}

	if req.Featured != nil {
		product.Featured = *req.Featured
		hasChanges = true
	}

	if req.Visibility != nil {
		product.Visibility = *req.Visibility
		hasChanges = true
	}

	// Handle Sale Pricing
	if req.SalePrice != nil {
		product.SalePrice = req.SalePrice
		hasChanges = true
	}

	if req.SaleStartDate != nil {
		product.SaleStartDate = req.SaleStartDate
		hasChanges = true
	}

	if req.SaleEndDate != nil {
		product.SaleEndDate = req.SaleEndDate
		hasChanges = true
	}

	// Validate sale pricing business rules after all sale fields are updated
	if hasChanges && (req.SalePrice != nil || req.SaleStartDate != nil || req.SaleEndDate != nil) {
		if err := product.ValidateSalePricing(); err != nil {
			return nil, fmt.Errorf("sale pricing validation failed: %w", err)
		}
	}

	// Handle Inventory Management
	if req.LowStockThreshold != nil {
		if *req.LowStockThreshold < 0 {
			return nil, fmt.Errorf("low stock threshold cannot be negative")
		}
		product.LowStockThreshold = *req.LowStockThreshold
		hasChanges = true
	}

	if req.TrackQuantity != nil {
		product.TrackQuantity = *req.TrackQuantity
		hasChanges = true
	}

	if req.AllowBackorder != nil {
		product.AllowBackorder = *req.AllowBackorder
		hasChanges = true
	}

	// Handle Shipping and Tax
	if req.RequiresShipping != nil {
		product.RequiresShipping = *req.RequiresShipping
		hasChanges = true
	}

	if req.ShippingClass != nil {
		product.ShippingClass = *req.ShippingClass
		hasChanges = true
	}

	if req.TaxClass != nil {
		product.TaxClass = *req.TaxClass
		hasChanges = true
	}

	if req.CountryOfOrigin != nil {
		product.CountryOfOrigin = *req.CountryOfOrigin
		hasChanges = true
	}

	// Handle Brand
	if req.BrandID != nil {
		product.BrandID = req.BrandID
		hasChanges = true
	}

	// Handle Product Type
	if req.ProductType != nil {
		product.ProductType = *req.ProductType
		hasChanges = true
	}

	// Update stock status if stock-related fields changed
	if req.Stock != nil || req.LowStockThreshold != nil || req.TrackQuantity != nil || req.AllowBackorder != nil {
		product.UpdateStockStatus()
		hasChanges = true
	}

	// Handle Images - Complete replacement if provided
	if req.Images != nil {
		if err := uc.replaceProductImages(ctx, product.ID, req.Images); err != nil {
			return nil, fmt.Errorf("failed to update images: %w", err)
		}
		hasChanges = true
	}

	// Handle Tags - Complete replacement if provided
	if req.Tags != nil {
		if err := uc.replaceProductTags(ctx, product.ID, req.Tags); err != nil {
			return nil, fmt.Errorf("failed to update tags: %w", err)
		}
		hasChanges = true
	}

	// Only update product if there were actual changes to basic fields
	if hasChanges {
		product.UpdatedAt = time.Now()
		if err := uc.productRepo.Update(ctx, product); err != nil {
			return nil, fmt.Errorf("failed to update product: %w", err)
		}
	}

	// Return updated product with fresh data - force fresh reload from database
	// Clear any potential cache by creating a fresh query
	updatedProduct, err := uc.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated product: %w", err)
	}

	return uc.toProductResponse(updatedProduct), nil
}

// PatchProduct partially updates a product - only updates provided fields
func (uc *productUseCase) PatchProduct(ctx context.Context, id uuid.UUID, req PatchProductRequest) (*ProductResponse, error) {
	// Get existing product
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrProductNotFound
	}

	var hasChanges bool

	// Basic field updates - only if provided
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return nil, fmt.Errorf("name cannot be empty")
		}
		product.Name = *req.Name
		hasChanges = true
	}

	if req.Description != nil {
		product.Description = *req.Description
		hasChanges = true
	}

	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, fmt.Errorf("price must be greater than 0")
		}
		product.Price = *req.Price
		hasChanges = true
	}

	if req.ComparePrice != nil {
		if *req.ComparePrice <= 0 {
			return nil, fmt.Errorf("compare price must be greater than 0")
		}
		product.ComparePrice = req.ComparePrice
		hasChanges = true
	}

	if req.CostPrice != nil {
		if *req.CostPrice < 0 {
			return nil, fmt.Errorf("cost price cannot be negative")
		}
		product.CostPrice = req.CostPrice
		hasChanges = true
	}

	if req.Stock != nil {
		if *req.Stock < 0 {
			return nil, fmt.Errorf("stock cannot be negative")
		}
		product.Stock = *req.Stock
		hasChanges = true
	}

	if req.Weight != nil {
		if *req.Weight <= 0 {
			return nil, fmt.Errorf("weight must be greater than 0")
		}
		product.Weight = req.Weight
		hasChanges = true
	}

	if req.CategoryID != nil {
		// Verify category exists
		_, err := uc.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, entities.ErrCategoryNotFound
		}
		// Use ProductCategory many-to-many as single source of truth
		// Remove all existing categories and assign new primary category
		if err := uc.productCategoryRepo.RemoveProductFromAllCategories(ctx, product.ID); err != nil {
			return nil, fmt.Errorf("failed to remove existing categories: %w", err)
		}
		if err := uc.productCategoryRepo.AssignProductToCategory(ctx, product.ID, *req.CategoryID, true); err != nil {
			return nil, fmt.Errorf("failed to assign primary category: %w", err)
		}
		hasChanges = true
	}

	if req.Status != nil {
		product.Status = *req.Status
		hasChanges = true
	}

	if req.IsDigital != nil {
		product.IsDigital = *req.IsDigital
		hasChanges = true
	}

	if req.Dimensions != nil {
		if req.Dimensions.Length <= 0 || req.Dimensions.Width <= 0 || req.Dimensions.Height <= 0 {
			return nil, fmt.Errorf("dimensions must be positive values")
		}
		product.Dimensions = &entities.Dimensions{
			Length: req.Dimensions.Length,
			Width:  req.Dimensions.Width,
			Height: req.Dimensions.Height,
		}
		hasChanges = true
	}

	// Handle new SEO and Metadata fields
	if req.ShortDescription != nil {
		product.ShortDescription = *req.ShortDescription
		hasChanges = true
	}

	if req.Slug != nil {
		if strings.TrimSpace(*req.Slug) == "" {
			return nil, fmt.Errorf("slug cannot be empty")
		}
		product.Slug = *req.Slug
		hasChanges = true
	}

	if req.MetaTitle != nil {
		product.MetaTitle = *req.MetaTitle
		hasChanges = true
	}

	if req.MetaDescription != nil {
		product.MetaDescription = *req.MetaDescription
		hasChanges = true
	}

	if req.Keywords != nil {
		product.Keywords = *req.Keywords
		hasChanges = true
	}

	if req.Featured != nil {
		product.Featured = *req.Featured
		hasChanges = true
	}

	if req.Visibility != nil {
		product.Visibility = *req.Visibility
		hasChanges = true
	}

	// Handle Sale Pricing
	if req.SalePrice != nil {
		product.SalePrice = req.SalePrice
		hasChanges = true
	}

	if req.SaleStartDate != nil {
		product.SaleStartDate = req.SaleStartDate
		hasChanges = true
	}

	if req.SaleEndDate != nil {
		product.SaleEndDate = req.SaleEndDate
		hasChanges = true
	}

	// Validate sale pricing business rules after all sale fields are updated
	if hasChanges && (req.SalePrice != nil || req.SaleStartDate != nil || req.SaleEndDate != nil) {
		if err := product.ValidateSalePricing(); err != nil {
			return nil, fmt.Errorf("sale pricing validation failed: %w", err)
		}
	}

	// Handle Inventory Management
	if req.LowStockThreshold != nil {
		if *req.LowStockThreshold < 0 {
			return nil, fmt.Errorf("low stock threshold cannot be negative")
		}
		product.LowStockThreshold = *req.LowStockThreshold
		hasChanges = true
	}

	if req.TrackQuantity != nil {
		product.TrackQuantity = *req.TrackQuantity
		hasChanges = true
	}

	if req.AllowBackorder != nil {
		product.AllowBackorder = *req.AllowBackorder
		hasChanges = true
	}

	// Handle Shipping and Tax
	if req.RequiresShipping != nil {
		product.RequiresShipping = *req.RequiresShipping
		hasChanges = true
	}

	if req.ShippingClass != nil {
		product.ShippingClass = *req.ShippingClass
		hasChanges = true
	}

	if req.TaxClass != nil {
		product.TaxClass = *req.TaxClass
		hasChanges = true
	}

	if req.CountryOfOrigin != nil {
		product.CountryOfOrigin = *req.CountryOfOrigin
		hasChanges = true
	}

	// Handle Brand
	if req.BrandID != nil {
		product.BrandID = req.BrandID
		hasChanges = true
	}

	// Handle Product Type
	if req.ProductType != nil {
		product.ProductType = *req.ProductType
		hasChanges = true
	}

	// Update stock status if stock-related fields changed
	if req.Stock != nil || req.LowStockThreshold != nil || req.TrackQuantity != nil || req.AllowBackorder != nil {
		product.UpdateStockStatus()
		hasChanges = true
	}

	// Handle Images - check if field is provided
	if req.Images != nil {
		// If empty slice, clear all images
		if len(*req.Images) == 0 {
			if err := uc.imageRepo.MarkAsInactive(ctx, product.ID); err != nil {
				return nil, fmt.Errorf("failed to clear images: %w", err)
			}
		} else {
			// Replace with new images
			if err := uc.replaceProductImages(ctx, product.ID, *req.Images); err != nil {
				return nil, fmt.Errorf("failed to update images: %w", err)
			}
		}
		hasChanges = true
	}

	// Handle Tags - check if field is provided
	if req.Tags != nil {
		// Convert to slice and process
		tagSlice := *req.Tags
		if err := uc.replaceProductTags(ctx, product.ID, tagSlice); err != nil {
			return nil, fmt.Errorf("failed to update tags: %w", err)
		}
		hasChanges = true
	}

	// Only update product if there were actual changes
	if hasChanges {
		product.UpdatedAt = time.Now()
		if err := uc.productRepo.Update(ctx, product); err != nil {
			return nil, fmt.Errorf("failed to update product: %w", err)
		}
	}

	// Return updated product with fresh data
	updatedProduct, err := uc.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated product: %w", err)
	}

	return uc.toProductResponse(updatedProduct), nil
}

// replaceProductImages completely replaces all product images with new ones
func (uc *productUseCase) replaceProductImages(ctx context.Context, productID uuid.UUID, images []ProductImageRequest) error {
	fmt.Printf("DEBUG: replaceProductImages called for productID: %s with %d new images\n", productID.String(), len(images))

	// Validate images
	for i, img := range images {
		if img.URL == "" {
			return fmt.Errorf("image URL cannot be empty at position %d", i)
		}
	}

	// Step 1: Get existing images
	fmt.Printf("DEBUG: Getting existing images for productID: %s\n", productID.String())
	existingImages, err := uc.imageRepo.GetByProductID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get existing images: %w", err)
	}
	fmt.Printf("DEBUG: Found %d existing images\n", len(existingImages))

	// Step 2: Update/Replace strategy instead of delete
	// Mark all existing images as "hidden" by setting position to -1
	if len(existingImages) > 0 {
		fmt.Printf("DEBUG: Marking %d existing images as hidden\n", len(existingImages))
		for _, img := range existingImages {
			img.Position = -1 // Mark as hidden
			if err := uc.imageRepo.Update(ctx, img); err != nil {
				fmt.Printf("DEBUG: Error hiding image %s: %v\n", img.ID.String(), err)
				return fmt.Errorf("failed to hide existing image: %w", err)
			}
		}
		fmt.Printf("DEBUG: Successfully marked existing images as hidden\n")
	}

	// Step 3: Create new images with positive positions
	if len(images) > 0 {
		fmt.Printf("DEBUG: Creating %d new images\n", len(images))
		var newImages []*entities.ProductImage
		for i, imgReq := range images {
			image := &entities.ProductImage{
				ID:        uuid.New(),
				ProductID: productID,
				URL:       imgReq.URL,
				AltText:   imgReq.AltText,
				Position:  i, // Positive position (0, 1, 2, ...)
				CreatedAt: time.Now(),
			}
			newImages = append(newImages, image)
			fmt.Printf("DEBUG: Prepared new image %d: %s at position %d\n", i, imgReq.URL, i)
		}

		if err := uc.imageRepo.CreateBatch(ctx, newImages); err != nil {
			fmt.Printf("DEBUG: Error creating new images: %v\n", err)
			return fmt.Errorf("failed to create new images: %w", err)
		}
		fmt.Printf("DEBUG: Successfully created %d new images\n", len(newImages))
	} else {
		fmt.Printf("DEBUG: No new images to create\n")
	}

	// Step 4: Verify by counting active images (position >= 0)
	activeImages, err := uc.getActiveImagesByProductID(ctx, productID)
	if err != nil {
		fmt.Printf("DEBUG: Error counting active images: %v\n", err)
	} else {
		fmt.Printf("DEBUG: After replacement, product has %d active images\n", len(activeImages))
	}

	fmt.Printf("DEBUG: replaceProductImages completed successfully\n")
	return nil
}

// Helper function to get active images (position >= 0)
func (uc *productUseCase) getActiveImagesByProductID(ctx context.Context, productID uuid.UUID) ([]*entities.ProductImage, error) {
	allImages, err := uc.imageRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	var activeImages []*entities.ProductImage
	for _, img := range allImages {
		if img.Position >= 0 {
			activeImages = append(activeImages, img)
		}
	}
	return activeImages, nil
}

// replaceProductTags completely replaces all product tags with new ones
func (uc *productUseCase) replaceProductTags(ctx context.Context, productID uuid.UUID, tagNames []string) error {
	// Validate and clean tag names
	var validTags []string
	for _, tagName := range tagNames {
		cleanTag := strings.TrimSpace(tagName)
		if cleanTag != "" && len(cleanTag) <= 50 { // Reasonable limit for tag length
			validTags = append(validTags, cleanTag)
		}
	}

	// If no valid tags, clear all tags
	if len(validTags) == 0 {
		return uc.productRepo.ClearTags(ctx, productID)
	}

	// Find or create all tags and collect their IDs
	var tagIDs []uuid.UUID
	for _, tagName := range validTags {
		tag, err := uc.tagRepo.FindOrCreate(ctx, tagName)
		if err != nil {
			return fmt.Errorf("failed to find/create tag '%s': %w", tagName, err)
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	// Replace all tags at once using the new ReplaceTags method
	if err := uc.productRepo.ReplaceTags(ctx, productID, tagIDs); err != nil {
		return fmt.Errorf("failed to replace product tags: %w", err)
	}

	return nil
}

// DeleteProduct deletes a product (same as original)
func (uc *productUseCase) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return entities.ErrProductNotFound
	}

	// First, remove all cart items that reference this product
	err = uc.cartRepo.RemoveItemsByProductID(ctx, id)
	if err != nil {
		return err
	}

	// Then delete the product
	return uc.productRepo.Delete(ctx, id)
}

// GetProducts gets list of products with pagination
func (uc *productUseCase) GetProducts(ctx context.Context, req GetProductsRequest) (*GetProductsResponse, error) {
	// Get total count
	total, err := uc.productRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// Get products
	products, err := uc.productRepo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	// Convert to responses
	responses := make([]*ProductResponse, len(products))
	for i, product := range products {
		responses[i] = uc.toProductResponse(product)
	}

	// Create pagination context
	context := &EcommercePaginationContext{
		EntityType: "products",
	}

	// Create pagination info with enhanced features
	pagination := NewPaginationInfoFromOffset(req.Offset, req.Limit, total)

	// Apply ecommerce enhancements
	if context != nil {
		// Adjust page sizes based on entity type
		pagination.PageSizes = []int{12, 24, 48, 96} // Grid-friendly sizes for products

		// Check if cursor pagination should be used
		pagination.UseCursor = ShouldUseCursorPagination(total, context.EntityType)

		// Generate cache key
		cacheParams := map[string]interface{}{
			"page":  pagination.Page,
			"limit": pagination.Limit,
		}
		pagination.CacheKey = GenerateCacheKey("products", "", cacheParams)
	}

	return &GetProductsResponse{
		Products:   responses,
		Pagination: pagination,
	}, nil
}

// SearchProducts searches products (same as original)
func (uc *productUseCase) SearchProducts(ctx context.Context, req SearchProductsRequest) ([]*ProductResponse, error) {
	params := repositories.ProductSearchParams{
		Query:      req.Query,
		CategoryID: req.CategoryID,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		Status:     req.Status,
		Tags:       req.Tags,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}

	products, err := uc.productRepo.Search(ctx, params)
	if err != nil {
		return nil, err
	}

	responses := make([]*ProductResponse, len(products))
	for i, product := range products {
		responses[i] = uc.toProductResponse(product)
	}

	return responses, nil
}

// calculateDiscountPercentage calculates the discount percentage for a product
// Uses the new unified discount logic from entity
func (uc *productUseCase) calculateDiscountPercentage(product *entities.Product) float64 {
	return product.GetDiscountPercentage()
}

// SearchProductsPaginated searches products with pagination info
func (uc *productUseCase) SearchProductsPaginated(ctx context.Context, req SearchProductsRequest) (*GetProductsResponse, error) {
	params := repositories.ProductSearchParams{
		Query:      req.Query,
		CategoryID: req.CategoryID,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		Status:     req.Status,
		Tags:       req.Tags,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}

	// Get total count using the new SearchCount method
	total, err := uc.productRepo.SearchCount(ctx, params)
	if err != nil {
		return nil, err
	}

	// Get products
	products, err := uc.productRepo.Search(ctx, params)
	if err != nil {
		return nil, err
	}

	responses := make([]*ProductResponse, len(products))
	for i, product := range products {
		responses[i] = uc.toProductResponse(product)
	}

	// Create pagination context
	context := &EcommercePaginationContext{
		EntityType: "products",
	}

	// Create pagination info with enhanced features
	pagination := NewPaginationInfoFromOffset(req.Offset, req.Limit, total)

	// Apply ecommerce enhancements
	if context != nil {
		// Adjust page sizes based on entity type
		pagination.PageSizes = []int{12, 24, 48, 96} // Grid-friendly sizes for products

		// Check if cursor pagination should be used
		pagination.UseCursor = ShouldUseCursorPagination(total, context.EntityType)

		// Generate cache key
		cacheParams := map[string]interface{}{
			"page":  pagination.Page,
			"limit": pagination.Limit,
			"query": req.Query,
		}
		if req.CategoryID != nil {
			cacheParams["category_id"] = *req.CategoryID
		}
		pagination.CacheKey = GenerateCacheKey("product_search", "", cacheParams)
	}

	return &GetProductsResponse{
		Products:   responses,
		Pagination: pagination,
	}, nil
}

// GetProductsByCategory gets products by category with pagination
func (uc *productUseCase) GetProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) (*GetProductsResponse, error) {
	// Get products
	products, err := uc.productRepo.GetByCategory(ctx, categoryID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Get total count for the category
	total, err := uc.productRepo.CountByCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// Convert to responses
	responses := make([]*ProductResponse, len(products))
	for i, product := range products {
		responses[i] = uc.toProductResponse(product)
	}

	// Create pagination context
	context := &EcommercePaginationContext{
		EntityType: "products",
	}

	// Create pagination info with enhanced features
	pagination := NewPaginationInfoFromOffset(offset, limit, total)

	// Apply ecommerce enhancements
	if context != nil {
		// Adjust page sizes based on entity type
		pagination.PageSizes = []int{12, 24, 48, 96} // Grid-friendly sizes for products

		// Check if cursor pagination should be used
		pagination.UseCursor = ShouldUseCursorPagination(total, context.EntityType)

		// Generate cache key
		cacheParams := map[string]interface{}{
			"page":        pagination.Page,
			"limit":       pagination.Limit,
			"category_id": categoryID.String(),
		}
		pagination.CacheKey = GenerateCacheKey("products_by_category", "", cacheParams)
	}

	return &GetProductsResponse{
		Products:   responses,
		Pagination: pagination,
	}, nil
}

// UpdateStock updates product stock (same as original)
func (uc *productUseCase) UpdateStock(ctx context.Context, productID uuid.UUID, stock int) error {
	_, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return entities.ErrProductNotFound
	}

	return uc.productRepo.UpdateStock(ctx, productID, stock)
}

// toProductResponse converts product entity to response (same as original)
func (uc *productUseCase) toProductResponse(product *entities.Product) *ProductResponse {
	response := &ProductResponse{
		ID:               product.ID,
		Name:             product.Name,
		Description:      product.Description,
		ShortDescription: product.ShortDescription,
		SKU:              product.SKU,

		// SEO and Metadata
		Slug:            product.Slug,
		MetaTitle:       product.MetaTitle,
		MetaDescription: product.MetaDescription,
		Keywords:        product.Keywords,
		Featured:        product.Featured,
		Visibility:      product.Visibility,

		// Pricing
		Price:        product.Price,
		ComparePrice: product.ComparePrice,
		CostPrice:    product.CostPrice,

		// Sale Pricing
		SalePrice:     product.SalePrice,
		SaleStartDate: product.SaleStartDate,
		SaleEndDate:   product.SaleEndDate,

		// Computed Price Fields
		CurrentPrice:           product.GetCurrentPrice(),
		OriginalPrice:          product.GetOriginalPrice(),
		IsOnSale:               product.IsOnSale(),
		HasDiscount:            product.HasDiscount(),
		SaleDiscountPercentage: product.GetSaleDiscountPercentage(),
		DiscountPercentage:     product.GetDiscountPercentage(),

		// Inventory
		Stock:             product.Stock,
		LowStockThreshold: product.LowStockThreshold,
		TrackQuantity:     product.TrackQuantity,
		AllowBackorder:    product.AllowBackorder,
		StockStatus:       product.StockStatus,
		IsLowStock:        product.IsLowStock(),

		// Physical Properties
		Weight: product.Weight,

		// Shipping and Tax
		RequiresShipping: product.RequiresShipping,
		ShippingClass:    product.ShippingClass,
		TaxClass:         product.TaxClass,
		CountryOfOrigin:  product.CountryOfOrigin,

		// Status and Type
		Status:      product.Status,
		ProductType: product.ProductType,
		IsDigital:   product.IsDigital,
		IsAvailable: product.IsAvailable(),
		HasVariants: product.HasVariants(),
		MainImage:   product.GetMainImage(),

		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}

	if product.Dimensions != nil {
		response.Dimensions = &DimensionsResponse{
			Length: product.Dimensions.Length,
			Width:  product.Dimensions.Width,
			Height: product.Dimensions.Height,
		}
	}

	// Get primary category from ProductCategory many-to-many (single source of truth)
	if primaryCategory, err := uc.productCategoryRepo.GetPrimaryCategory(context.Background(), product.ID); err == nil && primaryCategory != nil {
		response.Category = &ProductCategoryResponse{
			ID:          primaryCategory.ID,
			Name:        primaryCategory.Name,
			Description: primaryCategory.Description,
			Slug:        primaryCategory.Slug,
			Image:       primaryCategory.Image,
		}
	}

	// Convert brand
	if product.Brand != nil && product.Brand.ID != uuid.Nil {
		response.Brand = &ProductBrandResponse{
			ID:          product.Brand.ID,
			Name:        product.Brand.Name,
			Slug:        product.Brand.Slug,
			Description: product.Brand.Description,
			Logo:        product.Brand.Logo,
			Website:     product.Brand.Website,
			IsActive:    product.Brand.IsActive,
		}
	}

	// Convert images - only include active images (position >= 0)
	var activeImages []ProductImageResponse
	for _, img := range product.Images {
		if img.Position >= 0 { // Only include active images
			activeImages = append(activeImages, ProductImageResponse{
				ID:       img.ID,
				URL:      img.URL,
				AltText:  img.AltText,
				Position: img.Position,
			})
		}
	}
	response.Images = activeImages

	// Convert tags
	response.Tags = make([]ProductTagResponse, len(product.Tags))
	for i, tag := range product.Tags {
		response.Tags[i] = ProductTagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		}
	}

	// Convert attributes (placeholder for now)
	response.Attributes = make([]ProductAttributeResponse, len(product.Attributes))
	for i, attr := range product.Attributes {
		response.Attributes[i] = ProductAttributeResponse{
			ID:          attr.ID,
			AttributeID: attr.AttributeID,
			TermID:      attr.TermID,
			Value:       attr.Value,
			Position:    attr.Position,
		}
	}

	// Convert variants (placeholder for now)
	response.Variants = make([]ProductVariantResponse, len(product.Variants))
	for i, variant := range product.Variants {
		variantResponse := ProductVariantResponse{
			ID:           variant.ID,
			SKU:          variant.SKU,
			Price:        variant.Price,
			ComparePrice: variant.ComparePrice,
			CostPrice:    variant.CostPrice,
			Stock:        variant.Stock,
			Weight:       variant.Weight,
			Image:        variant.Image,
			Position:     variant.Position,
			IsActive:     variant.IsActive,
		}

		if variant.Dimensions != nil {
			variantResponse.Dimensions = &DimensionsResponse{
				Length: variant.Dimensions.Length,
				Width:  variant.Dimensions.Width,
				Height: variant.Dimensions.Height,
			}
		}

		// Convert variant attributes (placeholder)
		variantResponse.Attributes = make([]ProductVariantAttributeResponse, len(variant.Attributes))
		for j, varAttr := range variant.Attributes {
			variantResponse.Attributes[j] = ProductVariantAttributeResponse{
				ID:          varAttr.ID,
				AttributeID: varAttr.AttributeID,
				TermID:      varAttr.TermID,
			}
		}

		response.Variants[i] = variantResponse
	}

	return response
}

// replaceProductAttributes replaces all attributes for a product
func (uc *productUseCase) replaceProductAttributes(ctx context.Context, productID uuid.UUID, attributes []ProductAttributeRequest) error {
	// For now, we'll implement a basic version
	// In a full implementation, you would:
	// 1. Delete existing product attribute values
	// 2. Create new attribute values
	// 3. Validate that attributes and terms exist

	// TODO: Implement full attribute management
	// This is a placeholder for the attribute system
	return nil
}

// replaceProductVariants replaces all variants for a product
func (uc *productUseCase) replaceProductVariants(ctx context.Context, productID uuid.UUID, variants []ProductVariantRequest) error {
	// For now, we'll implement a basic version
	// In a full implementation, you would:
	// 1. Delete existing product variants
	// 2. Create new variants with their attributes
	// 3. Validate variant data

	// TODO: Implement full variant management
	// This is a placeholder for the variant system
	return nil
}

// GetSearchSuggestions returns search suggestions based on query
func (uc *productUseCase) GetSearchSuggestions(ctx context.Context, req SearchSuggestionsRequest) (*SearchSuggestionsResponse, error) {
	// Set default limit if not provided
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	suggestions, err := uc.productRepo.GetSearchSuggestions(ctx, req.Query, limit)
	if err != nil {
		return nil, err
	}

	return &SearchSuggestionsResponse{
		Suggestions: suggestions,
	}, nil
}

// GetPopularSearches returns popular search queries
func (uc *productUseCase) GetPopularSearches(ctx context.Context, limit int) (*PopularSearchesResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	searches, err := uc.productRepo.GetPopularSearches(ctx, limit)
	if err != nil {
		return nil, err
	}

	return &PopularSearchesResponse{
		Searches: searches,
	}, nil
}

// GetSearchHistory returns search history for a user
func (uc *productUseCase) GetSearchHistory(ctx context.Context, userID uuid.UUID, limit int) (*SearchHistoryResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	history, err := uc.productRepo.GetSearchHistory(ctx, userID, limit)
	if err != nil {
		return nil, err
	}

	return &SearchHistoryResponse{
		History: history,
	}, nil
}

// createInitialInventory creates initial inventory record for a new product
func (uc *productUseCase) createInitialInventory(ctx context.Context, product *entities.Product) error {
	// Get default warehouse
	warehouses, err := uc.warehouseRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get warehouses: %w", err)
	}

	var defaultWarehouse *entities.Warehouse
	for _, warehouse := range warehouses {
		if warehouse.IsDefault {
			defaultWarehouse = warehouse
			break
		}
	}

	// If no default warehouse, use the first one or create a default
	if defaultWarehouse == nil && len(warehouses) > 0 {
		defaultWarehouse = warehouses[0]
	}

	// If no warehouses exist, create a default one
	if defaultWarehouse == nil {
		defaultWarehouse = &entities.Warehouse{
			ID:          uuid.New(),
			Name:        "Main Warehouse",
			Code:        "MAIN",
			Description: "Default warehouse",
			Address:     "Default Address",
			City:        "Default City",
			State:       "Default State",
			ZipCode:     "00000",
			Country:     "USA",
			IsActive:    true,
			IsDefault:   true,
		}

		if err := uc.warehouseRepo.Create(ctx, defaultWarehouse); err != nil {
			return fmt.Errorf("failed to create default warehouse: %w", err)
		}
	}

	// Create initial inventory record
	inventory := &entities.Inventory{
		ID:                uuid.New(),
		ProductID:         product.ID,
		WarehouseID:       defaultWarehouse.ID,
		QuantityOnHand:    product.Stock,
		QuantityAvailable: product.Stock,
		QuantityReserved:  0,
		ReorderLevel:      product.LowStockThreshold,
		MaxStockLevel:     product.LowStockThreshold * 10, // Default to 10x reorder level
		MinStockLevel:     0,
		AverageCost:       0,
		LastCost:          0,
		IsActive:          true,
	}

	return uc.inventoryRepo.Create(ctx, inventory)
}

// GetFeaturedProductsPaginated gets featured products with pagination
func (uc *productUseCase) GetFeaturedProductsPaginated(ctx context.Context, page, limit int) (*FeaturedProductsPaginatedResponse, error) {
	// Get featured products using existing GetProducts method with featured filter
	req := GetProductsRequest{
		Limit:  limit * 10, // Get more to simulate featured products
		Offset: 0,
	}

	// Get all products and filter featured ones (in real implementation, this would be optimized)
	allProductsResponse, err := uc.GetProducts(ctx, req)
	if err != nil {
		return nil, err
	}

	// Filter featured products (mock implementation)
	var featuredProducts []*ProductResponse
	for _, product := range allProductsResponse.Products {
		if product.Featured { // Assuming Featured field exists
			featuredProducts = append(featuredProducts, product)
		}
	}

	// Calculate pagination
	total := int64(len(featuredProducts))
	offset := (page - 1) * limit

	// Get products for current page
	var products []*ProductResponse
	if offset < len(featuredProducts) {
		end := offset + limit
		if end > len(featuredProducts) {
			end = len(featuredProducts)
		}
		products = featuredProducts[offset:end]
	}

	// Create pagination context
	context := &EcommercePaginationContext{
		EntityType: "products",
	}

	// Create enhanced pagination info
	pagination := NewEcommercePaginationInfo(page, limit, total, context)

	return &FeaturedProductsPaginatedResponse{
		Products:   products,
		Pagination: pagination,
	}, nil
}

// GetTrendingProductsPaginated gets trending products with pagination
func (uc *productUseCase) GetTrendingProductsPaginated(ctx context.Context, page, limit int) (*TrendingProductsPaginatedResponse, error) {
	// Get trending products (in real implementation, this would be based on analytics)
	req := GetProductsRequest{
		Limit:  limit * 10, // Get more to simulate trending products
		Offset: 0,
	}

	// Get all products and sort by popularity (mock implementation)
	allProductsResponse, err := uc.GetProducts(ctx, req)
	if err != nil {
		return nil, err
	}

	// Sort products by popularity (mock implementation)
	trendingProducts := allProductsResponse.Products

	// Calculate pagination
	total := int64(len(trendingProducts))
	offset := (page - 1) * limit

	// Get products for current page
	var products []*ProductResponse
	if offset < len(trendingProducts) {
		end := offset + limit
		if end > len(trendingProducts) {
			end = len(trendingProducts)
		}
		products = trendingProducts[offset:end]
	}

	// Create pagination context
	context := &EcommercePaginationContext{
		EntityType: "products",
	}

	// Create enhanced pagination info
	pagination := NewEcommercePaginationInfo(page, limit, total, context)

	return &TrendingProductsPaginatedResponse{
		Products:   products,
		Pagination: pagination,
	}, nil
}

// GetRelatedProductsPaginated gets related products with pagination
func (uc *productUseCase) GetRelatedProductsPaginated(ctx context.Context, productID uuid.UUID, page, limit int) (*RelatedProductsPaginatedResponse, error) {
	// Get all products and filter related ones (in real implementation, this would be optimized)
	// Note: Product.CategoryID removed - related products logic simplified
	req := GetProductsRequest{
		Limit:  limit * 10, // Get more to simulate related products
		Offset: 0,
	}

	allProductsResponse, err := uc.GetProducts(ctx, req)
	if err != nil {
		return nil, err
	}

	// Filter related products (mock implementation)
	// In a real implementation, this would use category, tags, etc.
	var relatedProducts []*ProductResponse
	for _, p := range allProductsResponse.Products {
		if p.ID != productID { // Exclude current product
			relatedProducts = append(relatedProducts, p)
		}
	}

	// Calculate pagination
	total := int64(len(relatedProducts))
	offset := (page - 1) * limit

	// Get products for current page
	var products []*ProductResponse
	if offset < len(relatedProducts) {
		end := offset + limit
		if end > len(relatedProducts) {
			end = len(relatedProducts)
		}
		products = relatedProducts[offset:end]
	}

	// Create pagination context
	context := &EcommercePaginationContext{
		EntityType: "products",
		// CategoryID removed - using ProductCategory many-to-many
	}

	// Create enhanced pagination info
	pagination := NewEcommercePaginationInfo(page, limit, total, context)

	return &RelatedProductsPaginatedResponse{
		Products:   products,
		Pagination: pagination,
		ProductID:  productID,
	}, nil
}
