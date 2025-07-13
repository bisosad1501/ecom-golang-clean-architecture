package usecases

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/domain/services"
	"ecom-golang-clean-architecture/pkg/utils"
	"github.com/google/uuid"
)

// CategoryUseCase defines category use cases
type CategoryUseCase interface {
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (*CategoryResponse, error)
	GetCategory(ctx context.Context, id uuid.UUID) (*CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, req UpdateCategoryRequest) (*CategoryResponse, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	GetCategories(ctx context.Context, req GetCategoriesRequest) ([]*CategoryResponse, error)
	GetCategoryTree(ctx context.Context) ([]*CategoryResponse, error)
	GetRootCategories(ctx context.Context) ([]*CategoryResponse, error)
	GetCategoryChildren(ctx context.Context, parentID uuid.UUID) ([]*CategoryResponse, error)
	GetCategoryPath(ctx context.Context, categoryID uuid.UUID) ([]*CategoryResponse, error)
	GetCategoryProductCount(ctx context.Context, categoryID uuid.UUID) (int64, error)

	// GetCategoryLandingPage gets category landing page data
	GetCategoryLandingPage(ctx context.Context, req GetCategoryLandingPageRequest) (*CategoryLandingPageResponse, error)

	// Bulk operations
	BulkCreateCategories(ctx context.Context, req []CreateCategoryRequest) ([]*CategoryResponse, error)
	BulkUpdateCategories(ctx context.Context, req []BulkUpdateCategoryRequest) ([]*CategoryResponse, error)
	BulkDeleteCategories(ctx context.Context, ids []uuid.UUID) error

	// Advanced filtering
	SearchCategories(ctx context.Context, req SearchCategoriesRequest) (*CategoriesListResponse, error)
	GetCategoriesWithFilters(ctx context.Context, req GetCategoriesWithFiltersRequest) (*CategoriesListResponse, error)

	// Tree operations
	MoveCategory(ctx context.Context, req MoveCategoryRequest) error
	ReorderCategories(ctx context.Context, req ReorderCategoriesRequest) error
	GetCategoryTreeStats(ctx context.Context) (*CategoryTreeStatsResponse, error)
	ValidateAndRepairTree(ctx context.Context) (*TreeValidationResponse, error)

	// Analytics and statistics
	GetCategoryAnalytics(ctx context.Context, req GetCategoryAnalyticsRequest) (*CategoryAnalyticsResponse, error)
	GetTopCategories(ctx context.Context, req GetTopCategoriesRequest) (*TopCategoriesResponse, error)
	GetCategoryPerformanceMetrics(ctx context.Context, categoryID uuid.UUID) (*CategoryPerformanceResponse, error)
	GetCategorySalesStats(ctx context.Context, req GetCategorySalesStatsRequest) (*CategorySalesStatsResponse, error)

	// SEO operations
	UpdateCategorySEO(ctx context.Context, categoryID uuid.UUID, req CategorySEORequest) (*CategoryResponse, error)
	GetCategorySEO(ctx context.Context, categoryID uuid.UUID) (*CategorySEOResponse, error)
	GenerateCategorySEO(ctx context.Context, categoryID uuid.UUID) (*CategorySEOResponse, error)
	ValidateCategorySEO(ctx context.Context, categoryID uuid.UUID) (*CategorySEOValidationResponse, error)
}

type categoryUseCase struct {
	categoryRepo repositories.CategoryRepository
	productRepo  repositories.ProductRepository
	fileService  services.FileService
}

// NewCategoryUseCase creates a new category use case
func NewCategoryUseCase(categoryRepo repositories.CategoryRepository, productRepo repositories.ProductRepository, fileService services.FileService) CategoryUseCase {
	return &categoryUseCase{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
		fileService:  fileService,
	}
}

// CreateCategoryRequest represents create category request
type CreateCategoryRequest struct {
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description"`
	Slug        string     `json:"slug" validate:"required"`
	Image       string     `json:"image"`
	ParentID    *uuid.UUID `json:"parent_id"`
	IsActive    bool       `json:"is_active"`
	SortOrder   int        `json:"sort_order"`

	// SEO fields
	SEO *CategorySEORequest `json:"seo,omitempty"`
}

// UpdateCategoryRequest represents update category request
type UpdateCategoryRequest struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Slug        *string    `json:"slug"`
	Image       *string    `json:"image"`
	ParentID    *uuid.UUID `json:"parent_id"`
	IsActive    *bool      `json:"is_active"`
	SortOrder   *int       `json:"sort_order"`

	// SEO fields
	SEO *CategorySEORequest `json:"seo,omitempty"`
}

// CategorySEORequest represents category SEO metadata request
type CategorySEORequest struct {
	MetaTitle          *string `json:"meta_title,omitempty"`
	MetaDescription    *string `json:"meta_description,omitempty"`
	MetaKeywords       *string `json:"meta_keywords,omitempty"`
	CanonicalURL       *string `json:"canonical_url,omitempty"`
	OGTitle            *string `json:"og_title,omitempty"`
	OGDescription      *string `json:"og_description,omitempty"`
	OGImage            *string `json:"og_image,omitempty"`
	TwitterTitle       *string `json:"twitter_title,omitempty"`
	TwitterDescription *string `json:"twitter_description,omitempty"`
	TwitterImage       *string `json:"twitter_image,omitempty"`
	SchemaMarkup       *string `json:"schema_markup,omitempty"`
}

// GetCategoriesRequest represents get categories request
type GetCategoriesRequest struct {
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}

// GetCategoryLandingPageRequest represents category landing page request
type GetCategoryLandingPageRequest struct {
	CategoryID                 uuid.UUID `json:"category_id"`
	Page                      int       `json:"page"`
	Limit                     int       `json:"limit"`
	SortBy                    string    `json:"sort_by"`
	SortOrder                 string    `json:"sort_order"`
	IncludeSubcategoryProducts bool      `json:"include_subcategory_products"`
	IncludeFeatured           bool      `json:"include_featured"`
	FeaturedLimit             int       `json:"featured_limit"`
}

// CategoryResponse represents category response
type CategoryResponse struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Slug        string             `json:"slug"`
	Image       string             `json:"image"`
	ParentID    *uuid.UUID         `json:"parent_id"`
	Parent      *CategoryResponse  `json:"parent,omitempty"`
	Children    []CategoryResponse `json:"children,omitempty"`
	IsActive    bool               `json:"is_active"`
	SortOrder   int                `json:"sort_order"`
	Level       int                `json:"level"`
	Path        string             `json:"path"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`

	// SEO fields
	SEO *CategorySEOResponse `json:"seo,omitempty"`
}

// CategorySEOResponse represents category SEO metadata
type CategorySEOResponse struct {
	MetaTitle          string `json:"meta_title,omitempty"`
	MetaDescription    string `json:"meta_description,omitempty"`
	MetaKeywords       string `json:"meta_keywords,omitempty"`
	CanonicalURL       string `json:"canonical_url,omitempty"`
	OGTitle            string `json:"og_title,omitempty"`
	OGDescription      string `json:"og_description,omitempty"`
	OGImage            string `json:"og_image,omitempty"`
	TwitterTitle       string `json:"twitter_title,omitempty"`
	TwitterDescription string `json:"twitter_description,omitempty"`
	TwitterImage       string `json:"twitter_image,omitempty"`
	SchemaMarkup       string `json:"schema_markup,omitempty"`
}

// CategorySEOValidationResponse represents category SEO validation response
type CategorySEOValidationResponse struct {
	IsValid    bool                    `json:"is_valid"`
	Score      int                     `json:"score"` // SEO score out of 100
	Issues     []CategorySEOIssue      `json:"issues"`
	Suggestions []CategorySEOSuggestion `json:"suggestions"`
}

// CategorySEOIssue represents an SEO issue
type CategorySEOIssue struct {
	Field       string `json:"field"`
	Issue       string `json:"issue"`
	Severity    string `json:"severity"` // "error", "warning", "info"
	Description string `json:"description"`
}

// CategorySEOSuggestion represents an SEO suggestion
type CategorySEOSuggestion struct {
	Field       string `json:"field"`
	Suggestion  string `json:"suggestion"`
	Impact      string `json:"impact"` // "high", "medium", "low"
	Description string `json:"description"`
}

// CategoryLandingPageResponse represents category landing page response
type CategoryLandingPageResponse struct {
	Category      *CategoryResponse `json:"category"`
	Breadcrumbs   []*CategoryResponse `json:"breadcrumbs"`
	Children      []*CategoryResponse `json:"children"`
	Products      []*ProductResponse `json:"products"`
	FeaturedProducts []*ProductResponse `json:"featured_products,omitempty"`
	TotalProducts int64 `json:"total_products"`
	Page          int   `json:"page"`
	Limit         int   `json:"limit"`
	TotalPages    int   `json:"total_pages"`
}

// BulkUpdateCategoryRequest represents bulk update category request
type BulkUpdateCategoryRequest struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Slug        string     `json:"slug"`
	Image       string     `json:"image"`
	ParentID    *uuid.UUID `json:"parent_id"`
	IsActive    *bool      `json:"is_active"`
	SortOrder   *int       `json:"sort_order"`
}

// SearchCategoriesRequest represents search categories request
type SearchCategoriesRequest struct {
	Query  string `json:"query" validate:"required"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// GetCategoriesWithFiltersRequest represents advanced filtering request
type GetCategoriesWithFiltersRequest struct {
	Search    string     `json:"search"`
	ParentID  *uuid.UUID `json:"parent_id"`
	IsActive  *bool      `json:"is_active"`
	HasParent *bool      `json:"has_parent"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
	SortBy    string     `json:"sort_by"`    // name, created_at, sort_order
	SortOrder string     `json:"sort_order"` // asc, desc
}

// CategoriesListResponse represents paginated categories response
type CategoriesListResponse struct {
	Categories []*CategoryResponse `json:"categories"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}

// MoveCategoryRequest represents move category request
type MoveCategoryRequest struct {
	CategoryID    uuid.UUID  `json:"category_id" validate:"required"`
	NewParentID   *uuid.UUID `json:"new_parent_id"` // null for root level
	ValidateOnly  bool       `json:"validate_only"` // true to only validate without moving
}

// ReorderCategoriesRequest represents reorder categories request
type ReorderCategoriesRequest struct {
	Categories []CategoryReorderItem `json:"categories" validate:"required"`
}

// CategoryReorderItem represents a single category reorder item
type CategoryReorderItem struct {
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
	SortOrder  int       `json:"sort_order" validate:"required"`
}

// CategoryTreeStatsResponse represents category tree statistics
type CategoryTreeStatsResponse struct {
	TotalCategories   int                    `json:"total_categories"`
	RootCategories    int                    `json:"root_categories"`
	MaxDepth          int                    `json:"max_depth"`
	AverageDepth      float64                `json:"average_depth"`
	CategoriesByLevel map[int]int            `json:"categories_by_level"`
	LargestBranches   []CategoryBranchInfo   `json:"largest_branches"`
}

// CategoryBranchInfo represents information about a category branch
type CategoryBranchInfo struct {
	CategoryID       uuid.UUID `json:"category_id"`
	CategoryName     string    `json:"category_name"`
	DescendantCount  int       `json:"descendant_count"`
	DirectChildren   int       `json:"direct_children"`
	ProductCount     int64     `json:"product_count"`
}

// TreeValidationResponse represents tree validation results
type TreeValidationResponse struct {
	IsValid           bool                    `json:"is_valid"`
	Issues            []TreeValidationIssue   `json:"issues"`
	RepairsPerformed  []string                `json:"repairs_performed"`
	TotalIssuesFound  int                     `json:"total_issues_found"`
	TotalRepairs      int                     `json:"total_repairs"`
}

// TreeValidationIssue represents a tree validation issue
type TreeValidationIssue struct {
	Type        string    `json:"type"`        // circular_reference, orphaned_category, invalid_depth
	CategoryID  uuid.UUID `json:"category_id"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`    // critical, warning, info
}

// GetCategoryAnalyticsRequest represents get category analytics request
type GetCategoryAnalyticsRequest struct {
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
	TimeRange  string    `json:"time_range"` // 7d, 30d, 90d, 1y
}

// CategoryAnalyticsResponse represents category analytics response
type CategoryAnalyticsResponse struct {
	Analytics *repositories.CategoryAnalytics `json:"analytics"`
}

// GetTopCategoriesRequest represents get top categories request
type GetTopCategoriesRequest struct {
	Limit  int    `json:"limit" validate:"min=1,max=100"`
	SortBy string `json:"sort_by"` // sales, revenue, products, rating, growth
}

// TopCategoriesResponse represents top categories response
type TopCategoriesResponse struct {
	Categories []*repositories.CategoryStats `json:"categories"`
	Total      int                           `json:"total"`
}

// CategoryPerformanceResponse represents category performance response
type CategoryPerformanceResponse struct {
	Metrics *repositories.CategoryPerformanceMetrics `json:"metrics"`
}

// GetCategorySalesStatsRequest represents get category sales stats request
type GetCategorySalesStatsRequest struct {
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
	TimeRange  string    `json:"time_range"` // 7d, 30d, 90d, 1y
}

// CategorySalesStatsResponse represents category sales stats response
type CategorySalesStatsResponse struct {
	Stats *repositories.CategorySalesStats `json:"stats"`
}

// CreateCategory creates a new category
func (uc *categoryUseCase) CreateCategory(ctx context.Context, req CreateCategoryRequest) (*CategoryResponse, error) {
	// Check if slug already exists
	exists, err := uc.categoryRepo.ExistsBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entities.ErrConflict
	}

	// Verify parent category exists if provided
	if req.ParentID != nil {
		_, err = uc.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, entities.ErrCategoryNotFound
		}
	}

	// Generate slug if not provided
	if req.Slug == "" {
		req.Slug = generateSlug(req.Name)
	}

	// Create category
	category := &entities.Category{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Slug:        req.Slug,
		Image:       req.Image,
		ParentID:    req.ParentID,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return uc.toCategoryResponse(category), nil
}

// GetCategory gets a category by ID
func (uc *categoryUseCase) GetCategory(ctx context.Context, id uuid.UUID) (*CategoryResponse, error) {
	category, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrCategoryNotFound
	}

	return uc.toCategoryResponse(category), nil
}

// UpdateCategory updates a category
func (uc *categoryUseCase) UpdateCategory(ctx context.Context, id uuid.UUID, req UpdateCategoryRequest) (*CategoryResponse, error) {
	category, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrCategoryNotFound
	}

	// Store old image URL for cleanup
	oldImageURL := category.Image

	// Update fields
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = *req.Description
	}
	if req.Slug != nil {
		// Check if new slug already exists
		if *req.Slug != category.Slug {
			exists, err := uc.categoryRepo.ExistsBySlug(ctx, *req.Slug)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, entities.ErrConflict
			}
		}
		category.Slug = *req.Slug
	}
	if req.Image != nil {
		category.Image = *req.Image
	}
	if req.ParentID != nil {
		// Verify parent category exists
		_, err = uc.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, entities.ErrCategoryNotFound
		}
		category.ParentID = req.ParentID
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}

	category.UpdatedAt = time.Now()

	if err := uc.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	// Delete old image file if image was updated and it's different
	if req.Image != nil && oldImageURL != "" && oldImageURL != *req.Image {
		// Log for debugging
		println("DEBUG: Attempting to delete old image")
		println("DEBUG: Old image URL:", oldImageURL)
		println("DEBUG: New image URL:", *req.Image)
		
		// Extract object key from URL and delete using storage service
		if objectKey := utils.ExtractFilePathFromURL(oldImageURL); objectKey != "" {
			if err := uc.fileService.DeleteFile(ctx, objectKey); err != nil {
				// Log error but don't fail the update
				println("DEBUG: Failed to delete old image file:", err.Error())
			} else {
				println("DEBUG: Successfully deleted old image file")
			}
		}
	}

	return uc.toCategoryResponse(category), nil
}

// DeleteCategory deletes a category
func (uc *categoryUseCase) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	category, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return entities.ErrCategoryNotFound
	}

	// Store image URL for cleanup
	imageURL := category.Image

	// Check if category has children
	children, err := uc.categoryRepo.GetChildren(ctx, id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return entities.ErrConflict // Cannot delete category with children
	}

	// Delete category from database
	if err := uc.categoryRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Delete associated image file if exists
	if imageURL != "" {
		if err := utils.DeleteImageFile(imageURL); err != nil {
			// Log error but don't fail the deletion
			// The category is already deleted from database
		}
	}

	return nil
}

// GetCategories gets list of categories
func (uc *categoryUseCase) GetCategories(ctx context.Context, req GetCategoriesRequest) ([]*CategoryResponse, error) {
	categories, err := uc.categoryRepo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = uc.toCategoryResponse(category)
	}

	return responses, nil
}

// GetCategoryTree gets the category tree
func (uc *categoryUseCase) GetCategoryTree(ctx context.Context) ([]*CategoryResponse, error) {
	categories, err := uc.categoryRepo.GetTree(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = uc.toCategoryResponseWithChildren(category)
	}

	return responses, nil
}

// GetRootCategories gets root categories
func (uc *categoryUseCase) GetRootCategories(ctx context.Context) ([]*CategoryResponse, error) {
	categories, err := uc.categoryRepo.GetRootCategories(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = uc.toCategoryResponse(category)
	}

	return responses, nil
}

// GetCategoryChildren gets category children
func (uc *categoryUseCase) GetCategoryChildren(ctx context.Context, parentID uuid.UUID) ([]*CategoryResponse, error) {
	categories, err := uc.categoryRepo.GetChildren(ctx, parentID)
	if err != nil {
		return nil, err
	}

	responses := make([]*CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = uc.toCategoryResponse(category)
	}

	return responses, nil
}

// GetCategoryPath returns the path from root to a specific category
func (uc *categoryUseCase) GetCategoryPath(ctx context.Context, categoryID uuid.UUID) ([]*CategoryResponse, error) {
	categories, err := uc.categoryRepo.GetCategoryPath(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	responses := make([]*CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = uc.toCategoryResponse(category)
	}

	return responses, nil
}

// GetCategoryProductCount returns the total number of products in a category and its subcategories
func (uc *categoryUseCase) GetCategoryProductCount(ctx context.Context, categoryID uuid.UUID) (int64, error) {
	return uc.categoryRepo.GetProductCountByCategory(ctx, categoryID)
}

// GetCategoryLandingPage gets category landing page data
func (uc *categoryUseCase) GetCategoryLandingPage(ctx context.Context, req GetCategoryLandingPageRequest) (*CategoryLandingPageResponse, error) {
	// Get the category
	category, err := uc.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, entities.ErrCategoryNotFound
	}

	// Get breadcrumbs (path from root to current category)
	breadcrumbCategories, err := uc.categoryRepo.GetCategoryPath(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}

	breadcrumbs := make([]*CategoryResponse, len(breadcrumbCategories))
	for i, cat := range breadcrumbCategories {
		breadcrumbs[i] = uc.toCategoryResponse(cat)
	}

	// Get children categories
	childrenCategories, err := uc.categoryRepo.GetChildren(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}

	children := make([]*CategoryResponse, len(childrenCategories))
	for i, cat := range childrenCategories {
		children[i] = uc.toCategoryResponse(cat)
	}

	// Get products for this category
	offset := (req.Page - 1) * req.Limit
	if offset < 0 {
		offset = 0
	}

	var products []*entities.Product
	var totalProducts int64

	if req.IncludeSubcategoryProducts {
		// For now, just get products from the main category
		// TODO: Implement multi-category product search including subcategories
		products, err = uc.productRepo.GetByCategory(ctx, req.CategoryID, req.Limit, offset)
		if err != nil {
			return nil, err
		}

		totalProducts, err = uc.productRepo.CountByCategory(ctx, req.CategoryID)
		if err != nil {
			totalProducts = 0
		}
	} else {
		// Get products only from this category
		products, err = uc.productRepo.GetByCategory(ctx, req.CategoryID, req.Limit, offset)
		if err != nil {
			return nil, err
		}

		totalProducts, err = uc.productRepo.CountByCategory(ctx, req.CategoryID)
		if err != nil {
			totalProducts = 0
		}
	}

	// Convert products to response format
	productResponses := make([]*ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = uc.toProductResponse(product)
	}

	// Get featured products in this category if requested
	var featuredProductResponses []*ProductResponse
	if req.IncludeFeatured {
		featuredLimit := req.FeaturedLimit
		if featuredLimit <= 0 {
			featuredLimit = 6 // Default featured products limit
		}

		featuredProducts, err := uc.productRepo.GetFeaturedByCategory(ctx, req.CategoryID, featuredLimit)
		if err == nil && len(featuredProducts) > 0 {
			featuredProductResponses = make([]*ProductResponse, len(featuredProducts))
			for i, product := range featuredProducts {
				featuredProductResponses[i] = uc.toProductResponse(product)
			}
		}
	}

	// Calculate pagination
	totalPages := int((totalProducts + int64(req.Limit) - 1) / int64(req.Limit))

	response := &CategoryLandingPageResponse{
		Category:         uc.toCategoryResponse(category),
		Breadcrumbs:      breadcrumbs,
		Children:         children,
		Products:         productResponses,
		FeaturedProducts: featuredProductResponses,
		TotalProducts:    totalProducts,
		Page:             req.Page,
		Limit:            req.Limit,
		TotalPages:       totalPages,
	}

	return response, nil
}

// toCategoryResponse converts category entity to response
func (uc *categoryUseCase) toCategoryResponse(category *entities.Category) *CategoryResponse {
	response := &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Slug:        category.Slug,
		Image:       category.Image,
		ParentID:    category.ParentID,
		IsActive:    category.IsActive,
		SortOrder:   category.SortOrder,
		Level:       category.GetLevel(),
		Path:        category.GetPath(),
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	// Convert parent if available
	if category.Parent != nil {
		response.Parent = &CategoryResponse{
			ID:          category.Parent.ID,
			Name:        category.Parent.Name,
			Description: category.Parent.Description,
			Slug:        category.Parent.Slug,
			Image:       category.Parent.Image,
			ParentID:    category.Parent.ParentID,
			IsActive:    category.Parent.IsActive,
			SortOrder:   category.Parent.SortOrder,
			Level:       category.Parent.GetLevel(),
			Path:        category.Parent.GetPath(),
			CreatedAt:   category.Parent.CreatedAt,
			UpdatedAt:   category.Parent.UpdatedAt,
		}
	}

	// Add SEO data if available
	if category.MetaTitle != "" || category.MetaDescription != "" || category.MetaKeywords != "" ||
		category.CanonicalURL != "" || category.OGTitle != "" || category.OGDescription != "" ||
		category.OGImage != "" || category.TwitterTitle != "" || category.TwitterDescription != "" ||
		category.TwitterImage != "" || category.SchemaMarkup != "" {
		response.SEO = &CategorySEOResponse{
			MetaTitle:          category.MetaTitle,
			MetaDescription:    category.MetaDescription,
			MetaKeywords:       category.MetaKeywords,
			CanonicalURL:       category.CanonicalURL,
			OGTitle:            category.OGTitle,
			OGDescription:      category.OGDescription,
			OGImage:            category.OGImage,
			TwitterTitle:       category.TwitterTitle,
			TwitterDescription: category.TwitterDescription,
			TwitterImage:       category.TwitterImage,
			SchemaMarkup:       category.SchemaMarkup,
		}
	}

	return response
}

// toProductResponse converts product entity to response
func (uc *categoryUseCase) toProductResponse(product *entities.Product) *ProductResponse {
	response := &ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		SKU:         product.SKU,
		Price:       product.Price,
		SalePrice:   product.SalePrice,
		ComparePrice: product.ComparePrice,
		Stock:       product.Stock,
		Status:      product.Status,
		Weight:      product.Weight,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	// Set dimensions
	if product.Dimensions != nil {
		response.Dimensions = &DimensionsResponse{
			Length: product.Dimensions.Length,
			Width:  product.Dimensions.Width,
			Height: product.Dimensions.Height,
		}
	}

	// Set category
	if product.CategoryID != uuid.Nil && product.Category.ID != uuid.Nil {
		response.Category = &ProductCategoryResponse{
			ID:   product.Category.ID,
			Name: product.Category.Name,
			Slug: product.Category.Slug,
		}
	}

	// Set brand
	if product.BrandID != nil && product.Brand != nil {
		response.Brand = &ProductBrandResponse{
			ID:   product.Brand.ID,
			Name: product.Brand.Name,
			Slug: product.Brand.Slug,
		}
	}

	// Set images
	if len(product.Images) > 0 {
		response.Images = make([]ProductImageResponse, len(product.Images))
		for i, img := range product.Images {
			response.Images[i] = ProductImageResponse{
				ID:       img.ID,
				URL:      img.URL,
				AltText:  img.AltText,
				Position: img.Position,
			}
		}
		response.MainImage = product.Images[0].URL
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

	return response
}

// toCategoryResponseWithChildren converts category entity to response with children
func (uc *categoryUseCase) toCategoryResponseWithChildren(category *entities.Category) *CategoryResponse {
	response := uc.toCategoryResponse(category)

	// Convert children
	if len(category.Children) > 0 {
		response.Children = make([]CategoryResponse, len(category.Children))
		for i, child := range category.Children {
			childResponse := uc.toCategoryResponseWithChildren(&child)
			response.Children[i] = *childResponse
		}
	}

	return response
}

// generateSlug generates a URL-friendly slug from a name
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "&", "and")
	// Remove special characters (basic implementation)
	return slug
}

// BulkCreateCategories creates multiple categories
func (uc *categoryUseCase) BulkCreateCategories(ctx context.Context, req []CreateCategoryRequest) ([]*CategoryResponse, error) {
	if len(req) == 0 {
		return []*CategoryResponse{}, nil
	}

	var categories []*entities.Category

	for _, r := range req {
		// Check if slug already exists
		exists, err := uc.categoryRepo.ExistsBySlug(ctx, r.Slug)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, entities.ErrConflict
		}

		// Verify parent category exists if provided
		if r.ParentID != nil {
			_, err = uc.categoryRepo.GetByID(ctx, *r.ParentID)
			if err != nil {
				return nil, entities.ErrCategoryNotFound
			}
		}

		// Generate slug if not provided
		slug := r.Slug
		if slug == "" {
			slug = generateSlug(r.Name)
		}

		category := &entities.Category{
			ID:          uuid.New(),
			Name:        r.Name,
			Description: r.Description,
			Slug:        slug,
			Image:       r.Image,
			ParentID:    r.ParentID,
			IsActive:    r.IsActive,
			SortOrder:   r.SortOrder,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		categories = append(categories, category)
	}

	if err := uc.categoryRepo.BulkCreate(ctx, categories); err != nil {
		return nil, err
	}

	var responses []*CategoryResponse
	for _, category := range categories {
		responses = append(responses, uc.toCategoryResponse(category))
	}

	return responses, nil
}

// BulkUpdateCategories updates multiple categories
func (uc *categoryUseCase) BulkUpdateCategories(ctx context.Context, req []BulkUpdateCategoryRequest) ([]*CategoryResponse, error) {
	if len(req) == 0 {
		return []*CategoryResponse{}, nil
	}

	var categories []*entities.Category

	for _, r := range req {
		// Get existing category
		category, err := uc.categoryRepo.GetByID(ctx, r.ID)
		if err != nil {
			return nil, entities.ErrCategoryNotFound
		}

		// Update fields if provided
		if r.Name != "" {
			category.Name = r.Name
		}
		if r.Description != "" {
			category.Description = r.Description
		}
		if r.Slug != "" {
			// Check if new slug already exists (excluding current category)
			exists, err := uc.categoryRepo.ExistsBySlug(ctx, r.Slug)
			if err != nil {
				return nil, err
			}
			if exists {
				// Check if it's the same category
				existingCategory, err := uc.categoryRepo.GetBySlug(ctx, r.Slug)
				if err != nil {
					return nil, err
				}
				if existingCategory.ID != category.ID {
					return nil, entities.ErrConflict
				}
			}
			category.Slug = r.Slug
		}
		if r.Image != "" {
			category.Image = r.Image
		}
		if r.ParentID != nil {
			// Validate hierarchy to prevent circular references
			if err := uc.categoryRepo.ValidateHierarchy(ctx, category.ID, *r.ParentID); err != nil {
				return nil, err
			}
			category.ParentID = r.ParentID
		}
		if r.IsActive != nil {
			category.IsActive = *r.IsActive
		}
		if r.SortOrder != nil {
			category.SortOrder = *r.SortOrder
		}

		category.UpdatedAt = time.Now()
		categories = append(categories, category)
	}

	if err := uc.categoryRepo.BulkUpdate(ctx, categories); err != nil {
		return nil, err
	}

	var responses []*CategoryResponse
	for _, category := range categories {
		responses = append(responses, uc.toCategoryResponse(category))
	}

	return responses, nil
}

// BulkDeleteCategories deletes multiple categories
func (uc *categoryUseCase) BulkDeleteCategories(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	// Check if any category has children
	for _, id := range ids {
		children, err := uc.categoryRepo.GetChildren(ctx, id)
		if err != nil {
			return err
		}
		if len(children) > 0 {
			return entities.ErrCategoryHasChildren
		}

		// Check if category has products
		count, err := uc.categoryRepo.GetProductCount(ctx, id, false)
		if err != nil {
			return err
		}
		if count > 0 {
			return entities.ErrCategoryHasProducts
		}
	}

	return uc.categoryRepo.BulkDelete(ctx, ids)
}

// SearchCategories searches categories by query
func (uc *categoryUseCase) SearchCategories(ctx context.Context, req SearchCategoriesRequest) (*CategoriesListResponse, error) {
	// Set default pagination
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	categories, err := uc.categoryRepo.Search(ctx, req.Query, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	// Get total count for pagination
	filters := repositories.CategoryFilters{
		Search: req.Query,
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	total, err := uc.categoryRepo.CountWithFilters(ctx, filters)
	if err != nil {
		return nil, err
	}

	var responses []*CategoryResponse
	for _, category := range categories {
		responses = append(responses, uc.toCategoryResponse(category))
	}

	page := (req.Offset / req.Limit) + 1
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &CategoriesListResponse{
		Categories: responses,
		Total:      total,
		Page:       page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetCategoriesWithFilters gets categories with advanced filtering
func (uc *categoryUseCase) GetCategoriesWithFilters(ctx context.Context, req GetCategoriesWithFiltersRequest) (*CategoriesListResponse, error) {
	// Set default pagination
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Set default sorting
	if req.SortBy == "" {
		req.SortBy = "name"
	}
	if req.SortOrder == "" {
		req.SortOrder = "asc"
	}

	filters := repositories.CategoryFilters{
		Search:    req.Search,
		ParentID:  req.ParentID,
		IsActive:  req.IsActive,
		HasParent: req.HasParent,
		Limit:     req.Limit,
		Offset:    req.Offset,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	categories, err := uc.categoryRepo.ListWithFilters(ctx, filters)
	if err != nil {
		return nil, err
	}

	total, err := uc.categoryRepo.CountWithFilters(ctx, filters)
	if err != nil {
		return nil, err
	}

	var responses []*CategoryResponse
	for _, category := range categories {
		responses = append(responses, uc.toCategoryResponse(category))
	}

	page := (req.Offset / req.Limit) + 1
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &CategoriesListResponse{
		Categories: responses,
		Total:      total,
		Page:       page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// MoveCategory moves a category to a new parent
func (uc *categoryUseCase) MoveCategory(ctx context.Context, req MoveCategoryRequest) error {
	// Validate category exists
	_, err := uc.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return entities.ErrCategoryNotFound
	}

	// Validate new parent exists if provided
	if req.NewParentID != nil {
		_, err = uc.categoryRepo.GetByID(ctx, *req.NewParentID)
		if err != nil {
			return entities.ErrCategoryNotFound
		}

		// Validate hierarchy to prevent circular references
		if err := uc.categoryRepo.ValidateHierarchy(ctx, req.CategoryID, *req.NewParentID); err != nil {
			return err
		}
	}

	// If validate only, return without making changes
	if req.ValidateOnly {
		return nil
	}

	// Perform the move
	var newParentID uuid.UUID
	if req.NewParentID != nil {
		newParentID = *req.NewParentID
	}

	return uc.categoryRepo.MoveCategory(ctx, req.CategoryID, newParentID)
}

// ReorderCategories reorders multiple categories
func (uc *categoryUseCase) ReorderCategories(ctx context.Context, req ReorderCategoriesRequest) error {
	if len(req.Categories) == 0 {
		return nil
	}

	// Validate all categories exist
	for _, item := range req.Categories {
		_, err := uc.categoryRepo.GetByID(ctx, item.CategoryID)
		if err != nil {
			return entities.ErrCategoryNotFound
		}
	}

	// Convert to repository format
	var reorderRequests []repositories.CategoryReorderRequest
	for _, item := range req.Categories {
		reorderRequests = append(reorderRequests, repositories.CategoryReorderRequest{
			CategoryID: item.CategoryID,
			SortOrder:  item.SortOrder,
		})
	}

	return uc.categoryRepo.ReorderCategories(ctx, reorderRequests)
}

// GetCategoryTreeStats returns statistics about the category tree
func (uc *categoryUseCase) GetCategoryTreeStats(ctx context.Context) (*CategoryTreeStatsResponse, error) {
	// Get total categories
	totalCategories, err := uc.categoryRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// Get root categories count
	rootCategories, err := uc.categoryRepo.CountWithFilters(ctx, repositories.CategoryFilters{
		HasParent: &[]bool{false}[0],
	})
	if err != nil {
		return nil, err
	}

	// Get max depth
	maxDepth, err := uc.categoryRepo.GetMaxDepth(ctx)
	if err != nil {
		return nil, err
	}

	// Get all categories to calculate statistics
	allCategories, err := uc.categoryRepo.List(ctx, 10000, 0) // Large limit to get all
	if err != nil {
		return nil, err
	}

	// Calculate categories by level
	categoriesByLevel := make(map[int]int)
	totalDepth := 0

	for _, category := range allCategories {
		level := category.GetLevel()
		categoriesByLevel[level]++
		totalDepth += level
	}

	// Calculate average depth
	averageDepth := 0.0
	if len(allCategories) > 0 {
		averageDepth = float64(totalDepth) / float64(len(allCategories))
	}

	// Get largest branches (top 5 categories with most descendants)
	largestBranches := []CategoryBranchInfo{}
	for _, category := range allCategories {
		children, err := uc.categoryRepo.GetChildren(ctx, category.ID)
		if err != nil {
			continue
		}

		productCount, err := uc.categoryRepo.GetProductCount(ctx, category.ID, true)
		if err != nil {
			productCount = 0
		}

		// Count all descendants
		descendantCount := uc.countDescendants(ctx, category.ID)

		largestBranches = append(largestBranches, CategoryBranchInfo{
			CategoryID:      category.ID,
			CategoryName:    category.Name,
			DescendantCount: descendantCount,
			DirectChildren:  len(children),
			ProductCount:    productCount,
		})
	}

	// Sort by descendant count and take top 5
	sort.Slice(largestBranches, func(i, j int) bool {
		return largestBranches[i].DescendantCount > largestBranches[j].DescendantCount
	})

	if len(largestBranches) > 5 {
		largestBranches = largestBranches[:5]
	}

	return &CategoryTreeStatsResponse{
		TotalCategories:   int(totalCategories),
		RootCategories:    int(rootCategories),
		MaxDepth:          maxDepth,
		AverageDepth:      averageDepth,
		CategoriesByLevel: categoriesByLevel,
		LargestBranches:   largestBranches,
	}, nil
}

// ValidateAndRepairTree validates the category tree and performs repairs if needed
func (uc *categoryUseCase) ValidateAndRepairTree(ctx context.Context) (*TreeValidationResponse, error) {
	var issues []TreeValidationIssue
	var repairsPerformed []string

	// Validate tree integrity
	err := uc.categoryRepo.ValidateTreeIntegrity(ctx)
	if err != nil {
		if err == entities.ErrCircularReference {
			issues = append(issues, TreeValidationIssue{
				Type:        "circular_reference",
				Description: "Circular reference detected in category tree",
				Severity:    "critical",
			})
		}
	}

	// Check for orphaned categories (categories with non-existent parents)
	allCategories, err := uc.categoryRepo.List(ctx, 10000, 0) // Large limit to get all
	if err != nil {
		return nil, err
	}

	for _, category := range allCategories {
		if category.ParentID != nil {
			_, err := uc.categoryRepo.GetByID(ctx, *category.ParentID)
			if err != nil {
				issues = append(issues, TreeValidationIssue{
					Type:        "orphaned_category",
					CategoryID:  category.ID,
					Description: fmt.Sprintf("Category '%s' has non-existent parent", category.Name),
					Severity:    "warning",
				})
			}
		}
	}

	// Check for invalid depths
	maxDepth, err := uc.categoryRepo.GetMaxDepth(ctx)
	if err != nil {
		return nil, err
	}

	if maxDepth > 10 { // Assuming max depth of 10 levels
		issues = append(issues, TreeValidationIssue{
			Type:        "invalid_depth",
			Description: fmt.Sprintf("Category tree depth (%d) exceeds recommended maximum (10)", maxDepth),
			Severity:    "warning",
		})
	}

	// Perform repairs if needed
	if len(issues) > 0 {
		// Rebuild category paths
		err = uc.categoryRepo.RebuildCategoryPaths(ctx)
		if err == nil {
			repairsPerformed = append(repairsPerformed, "Rebuilt category paths")
		}
	}

	return &TreeValidationResponse{
		IsValid:          len(issues) == 0,
		Issues:           issues,
		RepairsPerformed: repairsPerformed,
		TotalIssuesFound: len(issues),
		TotalRepairs:     len(repairsPerformed),
	}, nil
}

// countDescendants counts all descendants of a category
func (uc *categoryUseCase) countDescendants(ctx context.Context, categoryID uuid.UUID) int {
	children, err := uc.categoryRepo.GetChildren(ctx, categoryID)
	if err != nil {
		return 0
	}

	count := len(children)
	for _, child := range children {
		count += uc.countDescendants(ctx, child.ID)
	}

	return count
}

// GetCategoryAnalytics returns comprehensive analytics for a category
func (uc *categoryUseCase) GetCategoryAnalytics(ctx context.Context, req GetCategoryAnalyticsRequest) (*CategoryAnalyticsResponse, error) {
	// Validate category exists
	_, err := uc.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, entities.ErrCategoryNotFound
	}

	// Set default time range if not provided
	timeRange := req.TimeRange
	if timeRange == "" {
		timeRange = "30d"
	}

	// Validate time range
	validRanges := map[string]bool{
		"7d": true, "30d": true, "90d": true, "1y": true,
	}
	if !validRanges[timeRange] {
		timeRange = "30d"
	}

	analytics, err := uc.categoryRepo.GetCategoryAnalytics(ctx, req.CategoryID, timeRange)
	if err != nil {
		return nil, err
	}

	return &CategoryAnalyticsResponse{
		Analytics: analytics,
	}, nil
}

// GetTopCategories returns top performing categories
func (uc *categoryUseCase) GetTopCategories(ctx context.Context, req GetTopCategoriesRequest) (*TopCategoriesResponse, error) {
	// Set default limit if not provided
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Set default sort by if not provided
	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "sales"
	}

	// Validate sort by
	validSortBy := map[string]bool{
		"sales": true, "revenue": true, "products": true, "rating": true, "growth": true,
	}
	if !validSortBy[sortBy] {
		sortBy = "sales"
	}

	categories, err := uc.categoryRepo.GetTopCategories(ctx, limit, sortBy)
	if err != nil {
		return nil, err
	}

	return &TopCategoriesResponse{
		Categories: categories,
		Total:      len(categories),
	}, nil
}

// GetCategoryPerformanceMetrics returns detailed performance metrics for a category
func (uc *categoryUseCase) GetCategoryPerformanceMetrics(ctx context.Context, categoryID uuid.UUID) (*CategoryPerformanceResponse, error) {
	// Validate category exists
	_, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, entities.ErrCategoryNotFound
	}

	metrics, err := uc.categoryRepo.GetCategoryPerformanceMetrics(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	return &CategoryPerformanceResponse{
		Metrics: metrics,
	}, nil
}

// GetCategorySalesStats returns sales statistics for a category
func (uc *categoryUseCase) GetCategorySalesStats(ctx context.Context, req GetCategorySalesStatsRequest) (*CategorySalesStatsResponse, error) {
	// Validate category exists
	_, err := uc.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, entities.ErrCategoryNotFound
	}

	// Set default time range if not provided
	timeRange := req.TimeRange
	if timeRange == "" {
		timeRange = "30d"
	}

	// Validate time range
	validRanges := map[string]bool{
		"7d": true, "30d": true, "90d": true, "1y": true,
	}
	if !validRanges[timeRange] {
		timeRange = "30d"
	}

	stats, err := uc.categoryRepo.GetCategorySalesStats(ctx, req.CategoryID, timeRange)
	if err != nil {
		return nil, err
	}

	return &CategorySalesStatsResponse{
		Stats: stats,
	}, nil
}

// UpdateCategorySEO updates SEO metadata for a category
func (uc *categoryUseCase) UpdateCategorySEO(ctx context.Context, categoryID uuid.UUID, req CategorySEORequest) (*CategoryResponse, error) {
	// Get existing category
	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, entities.ErrCategoryNotFound
	}

	// Update SEO fields
	if req.MetaTitle != nil {
		category.MetaTitle = *req.MetaTitle
	}
	if req.MetaDescription != nil {
		category.MetaDescription = *req.MetaDescription
	}
	if req.MetaKeywords != nil {
		category.MetaKeywords = *req.MetaKeywords
	}
	if req.CanonicalURL != nil {
		category.CanonicalURL = *req.CanonicalURL
	}
	if req.OGTitle != nil {
		category.OGTitle = *req.OGTitle
	}
	if req.OGDescription != nil {
		category.OGDescription = *req.OGDescription
	}
	if req.OGImage != nil {
		category.OGImage = *req.OGImage
	}
	if req.TwitterTitle != nil {
		category.TwitterTitle = *req.TwitterTitle
	}
	if req.TwitterDescription != nil {
		category.TwitterDescription = *req.TwitterDescription
	}
	if req.TwitterImage != nil {
		category.TwitterImage = *req.TwitterImage
	}
	if req.SchemaMarkup != nil {
		category.SchemaMarkup = *req.SchemaMarkup
	}

	// Update category
	err = uc.categoryRepo.Update(ctx, category)
	if err != nil {
		return nil, err
	}

	return uc.toCategoryResponse(category), nil
}

// GetCategorySEO gets SEO metadata for a category
func (uc *categoryUseCase) GetCategorySEO(ctx context.Context, categoryID uuid.UUID) (*CategorySEOResponse, error) {
	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, entities.ErrCategoryNotFound
	}

	return &CategorySEOResponse{
		MetaTitle:          category.MetaTitle,
		MetaDescription:    category.MetaDescription,
		MetaKeywords:       category.MetaKeywords,
		CanonicalURL:       category.CanonicalURL,
		OGTitle:            category.OGTitle,
		OGDescription:      category.OGDescription,
		OGImage:            category.OGImage,
		TwitterTitle:       category.TwitterTitle,
		TwitterDescription: category.TwitterDescription,
		TwitterImage:       category.TwitterImage,
		SchemaMarkup:       category.SchemaMarkup,
	}, nil
}

// GenerateCategorySEO automatically generates SEO metadata for a category
func (uc *categoryUseCase) GenerateCategorySEO(ctx context.Context, categoryID uuid.UUID) (*CategorySEOResponse, error) {
	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, entities.ErrCategoryNotFound
	}

	// Generate SEO metadata based on category data
	seo := &CategorySEOResponse{}

	// Generate meta title
	if category.MetaTitle == "" {
		seo.MetaTitle = category.Name + " - Shop Online"
		if len(seo.MetaTitle) > 60 {
			seo.MetaTitle = category.Name
		}
	} else {
		seo.MetaTitle = category.MetaTitle
	}

	// Generate meta description
	if category.MetaDescription == "" {
		if category.Description != "" {
			seo.MetaDescription = category.Description
			if len(seo.MetaDescription) > 160 {
				seo.MetaDescription = seo.MetaDescription[:157] + "..."
			}
		} else {
			seo.MetaDescription = "Shop " + category.Name + " products online. Find the best deals and latest products in " + category.Name + " category."
		}
	} else {
		seo.MetaDescription = category.MetaDescription
	}

	// Generate meta keywords
	if category.MetaKeywords == "" {
		seo.MetaKeywords = category.Name + ", shop " + category.Name + ", buy " + category.Name + " online"
	} else {
		seo.MetaKeywords = category.MetaKeywords
	}

	// Generate Open Graph data
	if category.OGTitle == "" {
		seo.OGTitle = seo.MetaTitle
	} else {
		seo.OGTitle = category.OGTitle
	}

	if category.OGDescription == "" {
		seo.OGDescription = seo.MetaDescription
	} else {
		seo.OGDescription = category.OGDescription
	}

	if category.OGImage == "" && category.Image != "" {
		seo.OGImage = category.Image
	} else {
		seo.OGImage = category.OGImage
	}

	// Generate Twitter Card data
	if category.TwitterTitle == "" {
		seo.TwitterTitle = seo.MetaTitle
	} else {
		seo.TwitterTitle = category.TwitterTitle
	}

	if category.TwitterDescription == "" {
		seo.TwitterDescription = seo.MetaDescription
	} else {
		seo.TwitterDescription = category.TwitterDescription
	}

	if category.TwitterImage == "" && category.Image != "" {
		seo.TwitterImage = category.Image
	} else {
		seo.TwitterImage = category.TwitterImage
	}

	// Generate canonical URL
	if category.CanonicalURL == "" {
		seo.CanonicalURL = "/categories/" + category.Slug
	} else {
		seo.CanonicalURL = category.CanonicalURL
	}

	return seo, nil
}

// ValidateCategorySEO validates SEO metadata for a category
func (uc *categoryUseCase) ValidateCategorySEO(ctx context.Context, categoryID uuid.UUID) (*CategorySEOValidationResponse, error) {
	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, entities.ErrCategoryNotFound
	}

	var issues []CategorySEOIssue
	var suggestions []CategorySEOSuggestion
	score := 100

	// Validate meta title
	if category.MetaTitle == "" {
		issues = append(issues, CategorySEOIssue{
			Field:       "meta_title",
			Issue:       "Missing meta title",
			Severity:    "error",
			Description: "Meta title is required for SEO",
		})
		suggestions = append(suggestions, CategorySEOSuggestion{
			Field:       "meta_title",
			Suggestion:  "Add a descriptive meta title (50-60 characters)",
			Impact:      "high",
			Description: "Meta title appears in search results and browser tabs",
		})
		score -= 20
	} else if len(category.MetaTitle) > 60 {
		issues = append(issues, CategorySEOIssue{
			Field:       "meta_title",
			Issue:       "Meta title too long",
			Severity:    "warning",
			Description: "Meta title should be under 60 characters",
		})
		score -= 10
	} else if len(category.MetaTitle) < 30 {
		issues = append(issues, CategorySEOIssue{
			Field:       "meta_title",
			Issue:       "Meta title too short",
			Severity:    "warning",
			Description: "Meta title should be at least 30 characters",
		})
		score -= 5
	}

	// Validate meta description
	if category.MetaDescription == "" {
		issues = append(issues, CategorySEOIssue{
			Field:       "meta_description",
			Issue:       "Missing meta description",
			Severity:    "error",
			Description: "Meta description is required for SEO",
		})
		suggestions = append(suggestions, CategorySEOSuggestion{
			Field:       "meta_description",
			Suggestion:  "Add a compelling meta description (150-160 characters)",
			Impact:      "high",
			Description: "Meta description appears in search results",
		})
		score -= 20
	} else if len(category.MetaDescription) > 160 {
		issues = append(issues, CategorySEOIssue{
			Field:       "meta_description",
			Issue:       "Meta description too long",
			Severity:    "warning",
			Description: "Meta description should be under 160 characters",
		})
		score -= 10
	} else if len(category.MetaDescription) < 120 {
		issues = append(issues, CategorySEOIssue{
			Field:       "meta_description",
			Issue:       "Meta description too short",
			Severity:    "info",
			Description: "Meta description could be longer for better SEO",
		})
		score -= 5
	}

	// Validate slug
	if category.Slug == "" {
		issues = append(issues, CategorySEOIssue{
			Field:       "slug",
			Issue:       "Missing URL slug",
			Severity:    "error",
			Description: "URL slug is required for SEO-friendly URLs",
		})
		score -= 15
	}

	// Validate Open Graph data
	if category.OGTitle == "" {
		suggestions = append(suggestions, CategorySEOSuggestion{
			Field:       "og_title",
			Suggestion:  "Add Open Graph title for social media sharing",
			Impact:      "medium",
			Description: "Improves appearance when shared on social media",
		})
		score -= 5
	}

	if category.OGDescription == "" {
		suggestions = append(suggestions, CategorySEOSuggestion{
			Field:       "og_description",
			Suggestion:  "Add Open Graph description for social media sharing",
			Impact:      "medium",
			Description: "Improves appearance when shared on social media",
		})
		score -= 5
	}

	if category.OGImage == "" {
		suggestions = append(suggestions, CategorySEOSuggestion{
			Field:       "og_image",
			Suggestion:  "Add Open Graph image for social media sharing",
			Impact:      "medium",
			Description: "Improves visual appeal when shared on social media",
		})
		score -= 5
	}

	// Validate canonical URL
	if category.CanonicalURL == "" {
		suggestions = append(suggestions, CategorySEOSuggestion{
			Field:       "canonical_url",
			Suggestion:  "Add canonical URL to prevent duplicate content issues",
			Impact:      "medium",
			Description: "Helps search engines understand the preferred URL",
		})
		score -= 5
	}

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}

	return &CategorySEOValidationResponse{
		IsValid:     len(issues) == 0 || (len(issues) > 0 && issues[0].Severity != "error"),
		Score:       score,
		Issues:      issues,
		Suggestions: suggestions,
	}, nil
}
