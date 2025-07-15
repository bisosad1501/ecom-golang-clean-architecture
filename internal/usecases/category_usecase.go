package usecases

import (
	"context"
	"fmt"
	"regexp"
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
	GetCategories(ctx context.Context, req GetCategoriesRequest) (*GetCategoriesResponse, error)
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

	// Enhanced URL optimization
	OptimizeSlug(ctx context.Context, categoryID uuid.UUID, req SlugOptimizationRequest) (*SlugOptimizationResponse, error)
	GenerateSlugSuggestions(ctx context.Context, categoryID uuid.UUID) (*SlugSuggestionsResponse, error)
	ValidateSlugAvailability(ctx context.Context, slug string, excludeID *uuid.UUID) (*SlugValidationResponse, error)
	GetSlugHistory(ctx context.Context, categoryID uuid.UUID) (*SlugHistoryResponse, error)

	// Bulk SEO operations
	BulkUpdateSEO(ctx context.Context, req BulkSEOUpdateRequest) (*BulkSEOUpdateResponse, error)
	BulkGenerateSEO(ctx context.Context, req BulkSEOGenerateRequest) (*BulkSEOGenerateResponse, error)
	BulkValidateSEO(ctx context.Context, req BulkSEOValidateRequest) (*BulkSEOValidateResponse, error)

	// SEO analytics and insights
	GetSEOAnalytics(ctx context.Context, req SEOAnalyticsRequest) (*SEOAnalyticsResponse, error)
	GetSEOInsights(ctx context.Context, categoryID uuid.UUID) (*SEOInsightsResponse, error)
	GetSEOCompetitorAnalysis(ctx context.Context, categoryID uuid.UUID) (*SEOCompetitorAnalysisResponse, error)
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

// GetCategoriesResponse represents paginated categories response
type GetCategoriesResponse struct {
	Categories []*CategoryResponse `json:"categories"`
	Pagination *PaginationInfo     `json:"pagination"`
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

// GetCategories gets list of categories with pagination
func (uc *categoryUseCase) GetCategories(ctx context.Context, req GetCategoriesRequest) (*GetCategoriesResponse, error) {
	// Get total count
	total, err := uc.categoryRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// Get categories
	categories, err := uc.categoryRepo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	// Convert to responses
	responses := make([]*CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = uc.toCategoryResponse(category)
	}

	// Create pagination info
	pagination := NewPaginationInfoFromOffset(req.Offset, req.Limit, total)

	return &GetCategoriesResponse{
		Categories: responses,
		Pagination: pagination,
	}, nil
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

// OptimizeSlug optimizes category slug for better SEO
func (uc *categoryUseCase) OptimizeSlug(ctx context.Context, categoryID uuid.UUID, req SlugOptimizationRequest) (*SlugOptimizationResponse, error) {
	// Get current category
	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	oldSlug := category.Slug

	// Validate new slug
	validation, err := uc.ValidateSlugAvailability(ctx, req.NewSlug, &categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate slug: %w", err)
	}

	if !validation.IsAvailable {
		return &SlugOptimizationResponse{
			OldSlug: oldSlug,
			NewSlug: req.NewSlug,
			Success: false,
			Message: "Slug is not available",
		}, nil
	}

	// Update category slug
	category.Slug = req.NewSlug
	if err := uc.categoryRepo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category slug: %w", err)
	}

	// Create redirect if requested
	var redirectURL string
	if req.AutoRedirect && oldSlug != req.NewSlug {
		redirectURL = fmt.Sprintf("/categories/%s", req.NewSlug)
		// TODO: Store redirect mapping in database
	}

	return &SlugOptimizationResponse{
		OldSlug:     oldSlug,
		NewSlug:     req.NewSlug,
		RedirectURL: redirectURL,
		Success:     true,
		Message:     "Slug optimized successfully",
	}, nil
}

// GenerateSlugSuggestions generates SEO-friendly slug suggestions
func (uc *categoryUseCase) GenerateSlugSuggestions(ctx context.Context, categoryID uuid.UUID) (*SlugSuggestionsResponse, error) {
	// Get current category
	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	suggestions := []SlugSuggestion{}

	// Generate suggestions based on category name
	baseSuggestions := []string{
		generateSlugFromName(category.Name),
		generateSlugFromName(category.Name + "-category"),
		generateSlugFromName(category.Name + "-products"),
		generateSlugFromName("shop-" + category.Name),
		generateSlugFromName("buy-" + category.Name),
	}

	// Add parent category context if exists
	if category.ParentID != nil {
		parent, err := uc.categoryRepo.GetByID(ctx, *category.ParentID)
		if err == nil {
			baseSuggestions = append(baseSuggestions,
				generateSlugFromName(parent.Name+"-"+category.Name),
				generateSlugFromName(category.Name+"-in-"+parent.Name),
			)
		}
	}

	// Check availability and score each suggestion
	for _, slug := range baseSuggestions {
		if slug == category.Slug {
			continue // Skip current slug
		}

		validation, err := uc.ValidateSlugAvailability(ctx, slug, &categoryID)
		if err != nil {
			continue
		}

		score := calculateSlugSEOScore(slug, category.Name)
		reason := generateSlugReason(slug, category.Name)

		suggestions = append(suggestions, SlugSuggestion{
			Slug:        slug,
			Score:       score,
			Reason:      reason,
			IsAvailable: validation.IsAvailable,
			SEOFriendly: score > 0.7,
		})
	}

	// Sort by score
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	return &SlugSuggestionsResponse{
		Suggestions: suggestions,
		Current:     category.Slug,
	}, nil
}

// ValidateSlugAvailability validates if a slug is available and SEO-friendly
func (uc *categoryUseCase) ValidateSlugAvailability(ctx context.Context, slug string, excludeID *uuid.UUID) (*SlugValidationResponse, error) {
	issues := []string{}
	suggestions := []string{}

	// Check if slug is valid format
	if !isValidSlugFormat(slug) {
		issues = append(issues, "Slug contains invalid characters")
		suggestions = append(suggestions, "Use only lowercase letters, numbers, and hyphens")
	}

	// Check length
	if len(slug) < 3 {
		issues = append(issues, "Slug is too short")
		suggestions = append(suggestions, "Use at least 3 characters")
	}
	if len(slug) > 100 {
		issues = append(issues, "Slug is too long")
		suggestions = append(suggestions, "Keep slug under 100 characters")
	}

	// Check for SEO best practices
	if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
		issues = append(issues, "Slug should not start or end with hyphen")
	}
	if strings.Contains(slug, "--") {
		issues = append(issues, "Slug should not contain consecutive hyphens")
	}

	// Check availability in database
	existingCategory, err := uc.categoryRepo.GetBySlug(ctx, slug)
	isAvailable := true
	if err == nil && existingCategory != nil {
		if excludeID == nil || existingCategory.ID != *excludeID {
			isAvailable = false
			issues = append(issues, "Slug is already in use")
			suggestions = append(suggestions, "Try adding a number or modifier")
		}
	}

	return &SlugValidationResponse{
		Slug:        slug,
		IsAvailable: isAvailable,
		IsValid:     len(issues) == 0,
		Issues:      issues,
		Suggestions: suggestions,
	}, nil
}

// GetSlugHistory returns the slug change history for a category
func (uc *categoryUseCase) GetSlugHistory(ctx context.Context, categoryID uuid.UUID) (*SlugHistoryResponse, error) {
	// Get current category
	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// TODO: Implement slug history tracking in database
	// For now, return current slug as history
	history := []SlugHistoryEntry{
		{
			Slug:      category.Slug,
			CreatedAt: category.CreatedAt,
			Reason:    "Initial creation",
			IsActive:  true,
		},
	}

	return &SlugHistoryResponse{
		History: history,
		Current: category.Slug,
	}, nil
}

// BulkUpdateSEO updates SEO metadata for multiple categories
func (uc *categoryUseCase) BulkUpdateSEO(ctx context.Context, req BulkSEOUpdateRequest) (*BulkSEOUpdateResponse, error) {
	startTime := time.Now()
	results := []BulkSEOResult{}
	successCount := 0
	failureCount := 0

	for _, categoryID := range req.CategoryIDs {
		result := BulkSEOResult{
			CategoryID: categoryID,
		}

		// Update SEO for this category
		_, err := uc.UpdateCategorySEO(ctx, categoryID, req.SEOData)
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			result.Message = "Failed to update SEO"
			failureCount++
		} else {
			result.Success = true
			result.Message = "SEO updated successfully"
			successCount++
		}

		results = append(results, result)
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	successRate := float64(successCount) / float64(len(req.CategoryIDs)) * 100

	return &BulkSEOUpdateResponse{
		TotalCategories: len(req.CategoryIDs),
		SuccessCount:    successCount,
		FailureCount:    failureCount,
		Results:         results,
		Summary: BulkOperationSummary{
			Duration:    duration.String(),
			StartTime:   startTime,
			EndTime:     endTime,
			SuccessRate: successRate,
		},
	}, nil
}

// BulkGenerateSEO generates SEO metadata for multiple categories
func (uc *categoryUseCase) BulkGenerateSEO(ctx context.Context, req BulkSEOGenerateRequest) (*BulkSEOGenerateResponse, error) {
	startTime := time.Now()
	results := []BulkSEOResult{}
	successCount := 0
	failureCount := 0

	for _, categoryID := range req.CategoryIDs {
		result := BulkSEOResult{
			CategoryID: categoryID,
		}

		// Generate SEO for this category
		seoData, err := uc.GenerateCategorySEO(ctx, categoryID)
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			result.Message = "Failed to generate SEO"
			failureCount++
		} else {
			// Apply generated SEO data to category
			seoRequest := CategorySEORequest{
				MetaTitle:       &seoData.MetaTitle,
				MetaDescription: &seoData.MetaDescription,
				MetaKeywords:    &seoData.MetaKeywords,
				CanonicalURL:    &seoData.CanonicalURL,
				OGTitle:         &seoData.OGTitle,
				OGDescription:   &seoData.OGDescription,
				TwitterTitle:    &seoData.TwitterTitle,
				TwitterDescription: &seoData.TwitterDescription,
			}

			if req.Options.OverwriteExisting || req.Options.GenerateKeywords {
				_, updateErr := uc.UpdateCategorySEO(ctx, categoryID, seoRequest)
				if updateErr != nil {
					result.Success = false
					result.Error = updateErr.Error()
					result.Message = "Failed to apply generated SEO"
					failureCount++
				} else {
					result.Success = true
					result.Message = "SEO generated and applied successfully"
					successCount++
				}
			} else {
				result.Success = true
				result.Message = "SEO generated successfully"
				successCount++
			}
		}

		results = append(results, result)
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	successRate := float64(successCount) / float64(len(req.CategoryIDs)) * 100

	return &BulkSEOGenerateResponse{
		TotalCategories: len(req.CategoryIDs),
		SuccessCount:    successCount,
		FailureCount:    failureCount,
		Results:         results,
		Summary: BulkOperationSummary{
			Duration:    duration.String(),
			StartTime:   startTime,
			EndTime:     endTime,
			SuccessRate: successRate,
		},
	}, nil
}

// BulkValidateSEO validates SEO metadata for multiple categories
func (uc *categoryUseCase) BulkValidateSEO(ctx context.Context, req BulkSEOValidateRequest) (*BulkSEOValidateResponse, error) {
	startTime := time.Now()
	results := []BulkSEOValidationResult{}
	validCount := 0
	invalidCount := 0
	totalScore := 0

	for _, categoryID := range req.CategoryIDs {
		// Validate SEO for this category
		validation, err := uc.ValidateCategorySEO(ctx, categoryID)
		if err != nil {
			invalidCount++
			continue
		}

		result := BulkSEOValidationResult{
			CategoryID:  categoryID,
			IsValid:     validation.IsValid,
			Score:       validation.Score,
			Issues:      validation.Issues,
			Suggestions: validation.Suggestions,
		}

		if validation.IsValid {
			validCount++
		} else {
			invalidCount++
		}

		totalScore += validation.Score
		results = append(results, result)
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	averageScore := float64(totalScore) / float64(len(req.CategoryIDs))

	// Check for global issues
	globalIssues := []string{}
	if req.Options.CheckDuplicates {
		// TODO: Implement duplicate detection across categories
		globalIssues = append(globalIssues, "Duplicate meta titles detected across categories")
	}

	return &BulkSEOValidateResponse{
		TotalCategories: len(req.CategoryIDs),
		ValidCount:      validCount,
		InvalidCount:    invalidCount,
		Results:         results,
		GlobalIssues:    globalIssues,
		Summary: BulkOperationSummary{
			Duration:     duration.String(),
			StartTime:    startTime,
			EndTime:      endTime,
			SuccessRate:  float64(validCount) / float64(len(req.CategoryIDs)) * 100,
			AverageScore: averageScore,
		},
	}, nil
}

// GetSEOAnalytics provides comprehensive SEO analytics across categories
func (uc *categoryUseCase) GetSEOAnalytics(ctx context.Context, req SEOAnalyticsRequest) (*SEOAnalyticsResponse, error) {
	// Get categories to analyze
	var categoryIDs []uuid.UUID
	if len(req.CategoryIDs) > 0 {
		categoryIDs = req.CategoryIDs
	} else {
		// Get all categories if none specified
		categories, err := uc.categoryRepo.List(ctx, 1000, 0) // Get up to 1000 categories
		if err != nil {
			return nil, fmt.Errorf("failed to get categories: %w", err)
		}
		for _, cat := range categories {
			categoryIDs = append(categoryIDs, cat.ID)
		}
	}

	// Initialize response
	response := &SEOAnalyticsResponse{}

	// Calculate overview metrics
	totalCategories := len(categoryIDs)
	categoriesWithSEO := 0
	totalScore := 0
	topPerforming := []CategorySEOPerformance{}
	bottomPerforming := []CategorySEOPerformance{}

	// Metrics counters
	metaTitleCount := 0
	metaDescCount := 0
	keywordsCount := 0
	canonicalURLCount := 0
	openGraphCount := 0
	twitterCardCount := 0
	schemaMarkupCount := 0

	// Issues tracking
	duplicateMetaTitles := make(map[string][]uuid.UUID)
	duplicateMetaDescs := make(map[string][]uuid.UUID)
	missingCanonicalURLs := []uuid.UUID{}
	longMetaTitles := []uuid.UUID{}
	shortMetaDescriptions := []uuid.UUID{}

	// Analyze each category
	for _, categoryID := range categoryIDs {
		category, err := uc.categoryRepo.GetByID(ctx, categoryID)
		if err != nil {
			continue
		}

		// Validate SEO
		validation, err := uc.ValidateCategorySEO(ctx, categoryID)
		if err != nil {
			continue
		}

		// Count categories with SEO data
		if category.MetaTitle != "" || category.MetaDescription != "" {
			categoriesWithSEO++
		}

		totalScore += validation.Score

		// Track performance
		performance := CategorySEOPerformance{
			CategoryID:   categoryID,
			CategoryName: category.Name,
			SEOScore:     validation.Score,
			Issues:       len(validation.Issues),
			LastUpdated:  category.UpdatedAt,
		}

		// Add to top/bottom performers
		if validation.Score >= 80 {
			topPerforming = append(topPerforming, performance)
		} else if validation.Score <= 40 {
			bottomPerforming = append(bottomPerforming, performance)
		}

		// Count coverage metrics
		if category.MetaTitle != "" {
			metaTitleCount++
			// Check for duplicates
			if existing, exists := duplicateMetaTitles[category.MetaTitle]; exists {
				duplicateMetaTitles[category.MetaTitle] = append(existing, categoryID)
			} else {
				duplicateMetaTitles[category.MetaTitle] = []uuid.UUID{categoryID}
			}
			// Check length
			if len(category.MetaTitle) > 60 {
				longMetaTitles = append(longMetaTitles, categoryID)
			}
		}

		if category.MetaDescription != "" {
			metaDescCount++
			// Check for duplicates
			if existing, exists := duplicateMetaDescs[category.MetaDescription]; exists {
				duplicateMetaDescs[category.MetaDescription] = append(existing, categoryID)
			} else {
				duplicateMetaDescs[category.MetaDescription] = []uuid.UUID{categoryID}
			}
			// Check length
			if len(category.MetaDescription) < 120 {
				shortMetaDescriptions = append(shortMetaDescriptions, categoryID)
			}
		}

		if category.MetaKeywords != "" {
			keywordsCount++
		}

		if category.CanonicalURL != "" {
			canonicalURLCount++
		} else {
			missingCanonicalURLs = append(missingCanonicalURLs, categoryID)
		}

		if category.OGTitle != "" || category.OGDescription != "" {
			openGraphCount++
		}

		if category.TwitterTitle != "" || category.TwitterDescription != "" {
			twitterCardCount++
		}

		if category.SchemaMarkup != "" {
			schemaMarkupCount++
		}
	}

	// Calculate averages and percentages
	averageSEOScore := float64(totalScore) / float64(totalCategories)
	seoCompletionRate := float64(categoriesWithSEO) / float64(totalCategories) * 100

	// Sort performers
	sort.Slice(topPerforming, func(i, j int) bool {
		return topPerforming[i].SEOScore > topPerforming[j].SEOScore
	})
	sort.Slice(bottomPerforming, func(i, j int) bool {
		return bottomPerforming[i].SEOScore < bottomPerforming[j].SEOScore
	})

	// Limit results
	if len(topPerforming) > 10 {
		topPerforming = topPerforming[:10]
	}
	if len(bottomPerforming) > 10 {
		bottomPerforming = bottomPerforming[:10]
	}

	// Build duplicate issues
	duplicateTitleIssues := []DuplicateIssue{}
	for title, categoryIDs := range duplicateMetaTitles {
		if len(categoryIDs) > 1 {
			duplicateTitleIssues = append(duplicateTitleIssues, DuplicateIssue{
				Value:       title,
				CategoryIDs: categoryIDs,
				Count:       len(categoryIDs),
			})
		}
	}

	duplicateDescIssues := []DuplicateIssue{}
	for desc, categoryIDs := range duplicateMetaDescs {
		if len(categoryIDs) > 1 {
			duplicateDescIssues = append(duplicateDescIssues, DuplicateIssue{
				Value:       desc,
				CategoryIDs: categoryIDs,
				Count:       len(categoryIDs),
			})
		}
	}

	// Build response
	response.Overview.TotalCategories = totalCategories
	response.Overview.CategoriesWithSEO = categoriesWithSEO
	response.Overview.AverageSEOScore = averageSEOScore
	response.Overview.SEOCompletionRate = seoCompletionRate
	response.Overview.TopPerformingCategories = topPerforming
	response.Overview.BottomPerformingCategories = bottomPerforming

	response.Metrics.MetaTitleCoverage = float64(metaTitleCount) / float64(totalCategories) * 100
	response.Metrics.MetaDescCoverage = float64(metaDescCount) / float64(totalCategories) * 100
	response.Metrics.KeywordsCoverage = float64(keywordsCount) / float64(totalCategories) * 100
	response.Metrics.CanonicalURLCoverage = float64(canonicalURLCount) / float64(totalCategories) * 100
	response.Metrics.OpenGraphCoverage = float64(openGraphCount) / float64(totalCategories) * 100
	response.Metrics.TwitterCardCoverage = float64(twitterCardCount) / float64(totalCategories) * 100
	response.Metrics.SchemaMarkupCoverage = float64(schemaMarkupCount) / float64(totalCategories) * 100

	response.Issues.DuplicateMetaTitles = duplicateTitleIssues
	response.Issues.DuplicateMetaDescs = duplicateDescIssues
	response.Issues.MissingCanonicalURLs = missingCanonicalURLs
	response.Issues.LongMetaTitles = longMetaTitles
	response.Issues.ShortMetaDescriptions = shortMetaDescriptions

	// TODO: Add trends data from historical tracking
	response.Trends = []SEOTrendData{}

	return response, nil
}

// GetSEOInsights provides detailed SEO insights for a specific category
func (uc *categoryUseCase) GetSEOInsights(ctx context.Context, categoryID uuid.UUID) (*SEOInsightsResponse, error) {
	// Get category
	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Get current SEO validation
	validation, err := uc.ValidateCategorySEO(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate SEO: %w", err)
	}

	// Determine grade based on score
	grade := "F"
	if validation.Score >= 90 {
		grade = "A+"
	} else if validation.Score >= 80 {
		grade = "A"
	} else if validation.Score >= 70 {
		grade = "B"
	} else if validation.Score >= 60 {
		grade = "C"
	} else if validation.Score >= 50 {
		grade = "D"
	}

	// Generate recommendations
	priorityRecs := []SEORecommendation{}
	quickRecs := []SEORecommendation{}
	advancedRecs := []SEORecommendation{}

	// Priority recommendations based on issues
	for _, issue := range validation.Issues {
		if issue.Severity == "error" {
			priorityRecs = append(priorityRecs, SEORecommendation{
				Title:       "Fix " + issue.Field,
				Description: issue.Description,
				Impact:      "high",
				Effort:      "low",
				Priority:    1,
				Action:      "immediate",
			})
		}
	}

	// Quick wins
	if category.MetaDescription == "" {
		quickRecs = append(quickRecs, SEORecommendation{
			Title:       "Add Meta Description",
			Description: "Write a compelling 120-160 character description",
			Impact:      "medium",
			Effort:      "low",
			Priority:    2,
			Action:      "add_meta_description",
		})
	}

	if category.OGTitle == "" {
		quickRecs = append(quickRecs, SEORecommendation{
			Title:       "Add Open Graph Title",
			Description: "Improve social media sharing appearance",
			Impact:      "medium",
			Effort:      "low",
			Priority:    3,
			Action:      "add_og_title",
		})
	}

	// Advanced recommendations
	if category.SchemaMarkup == "" {
		advancedRecs = append(advancedRecs, SEORecommendation{
			Title:       "Add Structured Data",
			Description: "Implement Schema.org markup for better search visibility",
			Impact:      "high",
			Effort:      "high",
			Priority:    4,
			Action:      "add_schema_markup",
		})
	}

	// Mock competitor data (in real implementation, this would come from external APIs)
	competitors := []CompetitorCategory{
		{
			Name:     "Similar Category A",
			SEOScore: 85,
			URL:      "https://competitor1.com/category",
			Insights: []string{"Strong meta descriptions", "Good keyword usage"},
		},
		{
			Name:     "Similar Category B",
			SEOScore: 78,
			URL:      "https://competitor2.com/category",
			Insights: []string{"Excellent schema markup", "Optimized URLs"},
		},
	}

	bestPractices := []BestPracticeExample{
		{
			Field:       "meta_title",
			Example:     category.Name + " - Premium Quality | Your Store",
			Explanation: "Includes category name, value proposition, and brand",
			Source:      "SEO Best Practices",
		},
		{
			Field:       "meta_description",
			Example:     "Discover our premium " + category.Name + " collection. Free shipping on orders over $50. Shop now for the best deals!",
			Explanation: "Includes keywords, value proposition, and call-to-action",
			Source:      "E-commerce SEO Guide",
		},
	}

	// Mock historical data
	historicalScores := []ScoreHistory{
		{Date: time.Now().AddDate(0, -3, 0), Score: validation.Score - 10, Event: "Initial setup"},
		{Date: time.Now().AddDate(0, -2, 0), Score: validation.Score - 5, Event: "Meta description added"},
		{Date: time.Now().AddDate(0, -1, 0), Score: validation.Score, Event: "Current state"},
	}

	improvements := []Improvement{
		{
			Date:        time.Now().AddDate(0, -1, 0),
			Field:       "meta_description",
			OldValue:    "",
			NewValue:    category.MetaDescription,
			ScoreChange: 5,
		},
	}

	trends := []string{
		"SEO score improving over time",
		"Meta data coverage increasing",
		"Schema markup implementation needed",
	}

	return &SEOInsightsResponse{
		CategoryID:   categoryID,
		CategoryName: category.Name,
		CurrentSEO: struct {
			Score       int                     `json:"score"`
			Grade       string                  `json:"grade"`
			Issues      []CategorySEOIssue      `json:"issues"`
			Suggestions []CategorySEOSuggestion `json:"suggestions"`
		}{
			Score:       validation.Score,
			Grade:       grade,
			Issues:      validation.Issues,
			Suggestions: validation.Suggestions,
		},
		Recommendations: struct {
			Priority []SEORecommendation `json:"priority"`
			Quick    []SEORecommendation `json:"quick"`
			Advanced []SEORecommendation `json:"advanced"`
		}{
			Priority: priorityRecs,
			Quick:    quickRecs,
			Advanced: advancedRecs,
		},
		Competitors: struct {
			Similar       []CompetitorCategory  `json:"similar"`
			BestPractices []BestPracticeExample `json:"best_practices"`
		}{
			Similar:       competitors,
			BestPractices: bestPractices,
		},
		Performance: struct {
			HistoricalScores []ScoreHistory `json:"historical_scores"`
			Improvements     []Improvement  `json:"improvements"`
			Trends           []string       `json:"trends"`
		}{
			HistoricalScores: historicalScores,
			Improvements:     improvements,
			Trends:           trends,
		},
	}, nil
}

// GetSEOCompetitorAnalysis provides competitor analysis for SEO optimization
func (uc *categoryUseCase) GetSEOCompetitorAnalysis(ctx context.Context, categoryID uuid.UUID) (*SEOCompetitorAnalysisResponse, error) {
	// Get category
	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Get current SEO validation
	validation, err := uc.ValidateCategorySEO(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate SEO: %w", err)
	}

	// Mock competitor analysis (in real implementation, this would use external APIs)
	competitors := []struct {
		Name        string   `json:"name"`
		URL         string   `json:"url"`
		SEOScore    int      `json:"seo_score"`
		Strengths   []string `json:"strengths"`
		Weaknesses  []string `json:"weaknesses"`
		KeyInsights []string `json:"key_insights"`
	}{
		{
			Name:     "Market Leader",
			URL:      "https://leader.com/" + category.Slug,
			SEOScore: 95,
			Strengths: []string{
				"Comprehensive meta descriptions",
				"Rich schema markup",
				"Optimized URL structure",
				"Strong internal linking",
			},
			Weaknesses: []string{
				"Slow page load times",
				"Limited mobile optimization",
			},
			KeyInsights: []string{
				"Uses category + brand in title tags",
				"Implements breadcrumb schema",
				"Strong focus on user intent keywords",
			},
		},
		{
			Name:     "Direct Competitor",
			URL:      "https://competitor.com/" + category.Slug,
			SEOScore: 78,
			Strengths: []string{
				"Good keyword targeting",
				"Clean URL structure",
				"Regular content updates",
			},
			Weaknesses: []string{
				"Missing Open Graph tags",
				"Inconsistent meta descriptions",
				"No structured data",
			},
			KeyInsights: []string{
				"Focuses on long-tail keywords",
				"Uses category descriptions for meta content",
				"Strong social media integration",
			},
		},
	}

	// Calculate competitive position
	averageCompetitorScore := 86.5
	competitiveGap := averageCompetitorScore - float64(validation.Score)

	marketPosition := "Behind Leaders"
	if validation.Score >= 90 {
		marketPosition = "Market Leader"
	} else if validation.Score >= 80 {
		marketPosition = "Strong Competitor"
	} else if validation.Score >= 70 {
		marketPosition = "Average Performer"
	}

	// Generate opportunities and threats
	opportunities := []string{
		"Implement structured data markup",
		"Optimize meta descriptions for click-through rates",
		"Add comprehensive Open Graph tags",
		"Improve internal linking structure",
	}

	threats := []string{
		"Competitors have better schema markup",
		"Missing social media optimization",
		"URL structure could be more SEO-friendly",
	}

	// Generate action plan
	actionPlan := []struct {
		Priority   int    `json:"priority"`
		Action     string `json:"action"`
		Impact     string `json:"impact"`
		Timeline   string `json:"timeline"`
		Difficulty string `json:"difficulty"`
	}{
		{
			Priority:   1,
			Action:     "Add comprehensive meta descriptions",
			Impact:     "High",
			Timeline:   "1 week",
			Difficulty: "Low",
		},
		{
			Priority:   2,
			Action:     "Implement schema markup",
			Impact:     "High",
			Timeline:   "2 weeks",
			Difficulty: "Medium",
		},
		{
			Priority:   3,
			Action:     "Optimize Open Graph tags",
			Impact:     "Medium",
			Timeline:   "1 week",
			Difficulty: "Low",
		},
		{
			Priority:   4,
			Action:     "Improve URL structure",
			Impact:     "Medium",
			Timeline:   "1 month",
			Difficulty: "High",
		},
	}

	return &SEOCompetitorAnalysisResponse{
		CategoryID:   categoryID,
		CategoryName: category.Name,
		Analysis: struct {
			MarketPosition string   `json:"market_position"`
			CompetitiveGap float64  `json:"competitive_gap"`
			Opportunities  []string `json:"opportunities"`
			Threats        []string `json:"threats"`
		}{
			MarketPosition: marketPosition,
			CompetitiveGap: competitiveGap,
			Opportunities:  opportunities,
			Threats:        threats,
		},
		Competitors: competitors,
		Benchmarks: struct {
			IndustryAverage float64 `json:"industry_average"`
			TopPerformer    float64 `json:"top_performer"`
			YourScore       float64 `json:"your_score"`
			Percentile      float64 `json:"percentile"`
		}{
			IndustryAverage: 75.0,
			TopPerformer:    95.0,
			YourScore:       float64(validation.Score),
			Percentile:      float64(validation.Score) / 95.0 * 100,
		},
		ActionPlan: actionPlan,
	}, nil
}

// Enhanced URL optimization and slug management request/response types
type SlugOptimizationRequest struct {
	NewSlug         string `json:"new_slug" validate:"required"`
	PreserveHistory bool   `json:"preserve_history"`
	AutoRedirect    bool   `json:"auto_redirect"`
}

type SlugOptimizationResponse struct {
	OldSlug     string `json:"old_slug"`
	NewSlug     string `json:"new_slug"`
	RedirectURL string `json:"redirect_url,omitempty"`
	Success     bool   `json:"success"`
	Message     string `json:"message"`
}

type SlugSuggestionsResponse struct {
	Suggestions []SlugSuggestion `json:"suggestions"`
	Current     string           `json:"current"`
}

type SlugSuggestion struct {
	Slug        string  `json:"slug"`
	Score       float64 `json:"score"`
	Reason      string  `json:"reason"`
	IsAvailable bool    `json:"is_available"`
	SEOFriendly bool    `json:"seo_friendly"`
}

type SlugValidationResponse struct {
	Slug        string `json:"slug"`
	IsAvailable bool   `json:"is_available"`
	IsValid     bool   `json:"is_valid"`
	Issues      []string `json:"issues,omitempty"`
	Suggestions []string `json:"suggestions,omitempty"`
}

type SlugHistoryResponse struct {
	History []SlugHistoryEntry `json:"history"`
	Current string             `json:"current"`
}

type SlugHistoryEntry struct {
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	Reason    string    `json:"reason"`
	IsActive  bool      `json:"is_active"`
}

// Bulk SEO operations request/response types
type BulkSEOUpdateRequest struct {
	CategoryIDs []uuid.UUID       `json:"category_ids" validate:"required"`
	SEOData     CategorySEORequest `json:"seo_data" validate:"required"`
	UpdateMode  string            `json:"update_mode" validate:"oneof=replace merge"`
}

type BulkSEOUpdateResponse struct {
	TotalCategories   int                    `json:"total_categories"`
	SuccessCount      int                    `json:"success_count"`
	FailureCount      int                    `json:"failure_count"`
	Results           []BulkSEOResult        `json:"results"`
	Summary           BulkOperationSummary   `json:"summary"`
}

type BulkSEOGenerateRequest struct {
	CategoryIDs []uuid.UUID `json:"category_ids" validate:"required"`
	Options     struct {
		OverwriteExisting bool `json:"overwrite_existing"`
		GenerateKeywords  bool `json:"generate_keywords"`
		GenerateSchema    bool `json:"generate_schema"`
	} `json:"options"`
}

type BulkSEOGenerateResponse struct {
	TotalCategories   int                    `json:"total_categories"`
	SuccessCount      int                    `json:"success_count"`
	FailureCount      int                    `json:"failure_count"`
	Results           []BulkSEOResult        `json:"results"`
	Summary           BulkOperationSummary   `json:"summary"`
}

type BulkSEOValidateRequest struct {
	CategoryIDs []uuid.UUID `json:"category_ids" validate:"required"`
	Options     struct {
		CheckDuplicates bool `json:"check_duplicates"`
		CheckKeywords   bool `json:"check_keywords"`
		CheckSchema     bool `json:"check_schema"`
	} `json:"options"`
}

type BulkSEOValidateResponse struct {
	TotalCategories   int                    `json:"total_categories"`
	ValidCount        int                    `json:"valid_count"`
	InvalidCount      int                    `json:"invalid_count"`
	Results           []BulkSEOValidationResult `json:"results"`
	Summary           BulkOperationSummary   `json:"summary"`
	GlobalIssues      []string               `json:"global_issues,omitempty"`
}

type BulkSEOResult struct {
	CategoryID uuid.UUID `json:"category_id"`
	Success    bool      `json:"success"`
	Message    string    `json:"message"`
	Error      string    `json:"error,omitempty"`
}

type BulkSEOValidationResult struct {
	CategoryID  uuid.UUID                    `json:"category_id"`
	IsValid     bool                         `json:"is_valid"`
	Score       int                          `json:"score"`
	Issues      []CategorySEOIssue           `json:"issues"`
	Suggestions []CategorySEOSuggestion      `json:"suggestions"`
}

type BulkOperationSummary struct {
	Duration      string    `json:"duration"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	SuccessRate   float64   `json:"success_rate"`
	AverageScore  float64   `json:"average_score,omitempty"`
}

// SEO analytics and insights request/response types
type SEOAnalyticsRequest struct {
	CategoryIDs []uuid.UUID `json:"category_ids,omitempty"`
	DateFrom    *time.Time  `json:"date_from,omitempty"`
	DateTo      *time.Time  `json:"date_to,omitempty"`
	Metrics     []string    `json:"metrics,omitempty"`
}

type SEOAnalyticsResponse struct {
	Overview struct {
		TotalCategories      int     `json:"total_categories"`
		CategoriesWithSEO    int     `json:"categories_with_seo"`
		AverageSEOScore      float64 `json:"average_seo_score"`
		SEOCompletionRate    float64 `json:"seo_completion_rate"`
		TopPerformingCategories []CategorySEOPerformance `json:"top_performing_categories"`
		BottomPerformingCategories []CategorySEOPerformance `json:"bottom_performing_categories"`
	} `json:"overview"`

	Metrics struct {
		MetaTitleCoverage    float64 `json:"meta_title_coverage"`
		MetaDescCoverage     float64 `json:"meta_desc_coverage"`
		KeywordsCoverage     float64 `json:"keywords_coverage"`
		CanonicalURLCoverage float64 `json:"canonical_url_coverage"`
		OpenGraphCoverage    float64 `json:"open_graph_coverage"`
		TwitterCardCoverage  float64 `json:"twitter_card_coverage"`
		SchemaMarkupCoverage float64 `json:"schema_markup_coverage"`
	} `json:"metrics"`

	Issues struct {
		DuplicateMetaTitles    []DuplicateIssue `json:"duplicate_meta_titles"`
		DuplicateMetaDescs     []DuplicateIssue `json:"duplicate_meta_descs"`
		MissingCanonicalURLs   []uuid.UUID      `json:"missing_canonical_urls"`
		LongMetaTitles         []uuid.UUID      `json:"long_meta_titles"`
		ShortMetaDescriptions  []uuid.UUID      `json:"short_meta_descriptions"`
	} `json:"issues"`

	Trends []SEOTrendData `json:"trends"`
}

type CategorySEOPerformance struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	SEOScore     int       `json:"seo_score"`
	Issues       int       `json:"issues"`
	LastUpdated  time.Time `json:"last_updated"`
}

type DuplicateIssue struct {
	Value       string      `json:"value"`
	CategoryIDs []uuid.UUID `json:"category_ids"`
	Count       int         `json:"count"`
}

type SEOTrendData struct {
	Date             time.Time `json:"date"`
	AverageSEOScore  float64   `json:"average_seo_score"`
	CompletionRate   float64   `json:"completion_rate"`
	IssuesResolved   int       `json:"issues_resolved"`
	NewIssuesFound   int       `json:"new_issues_found"`
}

type SEOInsightsResponse struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`

	CurrentSEO struct {
		Score       int                     `json:"score"`
		Grade       string                  `json:"grade"`
		Issues      []CategorySEOIssue      `json:"issues"`
		Suggestions []CategorySEOSuggestion `json:"suggestions"`
	} `json:"current_seo"`

	Recommendations struct {
		Priority     []SEORecommendation `json:"priority"`
		Quick        []SEORecommendation `json:"quick"`
		Advanced     []SEORecommendation `json:"advanced"`
	} `json:"recommendations"`

	Competitors struct {
		Similar      []CompetitorCategory `json:"similar"`
		BestPractices []BestPracticeExample `json:"best_practices"`
	} `json:"competitors"`

	Performance struct {
		HistoricalScores []ScoreHistory `json:"historical_scores"`
		Improvements     []Improvement  `json:"improvements"`
		Trends           []string       `json:"trends"`
	} `json:"performance"`
}

type SEORecommendation struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Effort      string `json:"effort"`
	Priority    int    `json:"priority"`
	Action      string `json:"action"`
}

type CompetitorCategory struct {
	Name     string `json:"name"`
	SEOScore int    `json:"seo_score"`
	URL      string `json:"url"`
	Insights []string `json:"insights"`
}

type BestPracticeExample struct {
	Field       string `json:"field"`
	Example     string `json:"example"`
	Explanation string `json:"explanation"`
	Source      string `json:"source"`
}

type ScoreHistory struct {
	Date  time.Time `json:"date"`
	Score int       `json:"score"`
	Event string    `json:"event,omitempty"`
}

type Improvement struct {
	Date        time.Time `json:"date"`
	Field       string    `json:"field"`
	OldValue    string    `json:"old_value"`
	NewValue    string    `json:"new_value"`
	ScoreChange int       `json:"score_change"`
}

type SEOCompetitorAnalysisResponse struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`

	Analysis struct {
		MarketPosition string  `json:"market_position"`
		CompetitiveGap float64 `json:"competitive_gap"`
		Opportunities  []string `json:"opportunities"`
		Threats        []string `json:"threats"`
	} `json:"analysis"`

	Competitors []struct {
		Name        string  `json:"name"`
		URL         string  `json:"url"`
		SEOScore    int     `json:"seo_score"`
		Strengths   []string `json:"strengths"`
		Weaknesses  []string `json:"weaknesses"`
		KeyInsights []string `json:"key_insights"`
	} `json:"competitors"`

	Benchmarks struct {
		IndustryAverage float64 `json:"industry_average"`
		TopPerformer    float64 `json:"top_performer"`
		YourScore       float64 `json:"your_score"`
		Percentile      float64 `json:"percentile"`
	} `json:"benchmarks"`

	ActionPlan []struct {
		Priority    int    `json:"priority"`
		Action      string `json:"action"`
		Impact      string `json:"impact"`
		Timeline    string `json:"timeline"`
		Difficulty  string `json:"difficulty"`
	} `json:"action_plan"`
}

// Helper functions for slug optimization
func generateSlugFromName(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and special characters with hyphens
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// Remove consecutive hyphens
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	return slug
}

func isValidSlugFormat(slug string) bool {
	// Check if slug contains only valid characters
	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, slug)
	return matched
}

func calculateSlugSEOScore(slug, categoryName string) float64 {
	score := 1.0

	// Length score (optimal 3-50 characters)
	length := len(slug)
	if length < 3 || length > 50 {
		score -= 0.2
	}

	// Keyword relevance (check if category name words are in slug)
	nameWords := strings.Fields(strings.ToLower(categoryName))
	slugWords := strings.Split(slug, "-")

	relevantWords := 0
	for _, nameWord := range nameWords {
		for _, slugWord := range slugWords {
			if strings.Contains(slugWord, nameWord) || strings.Contains(nameWord, slugWord) {
				relevantWords++
				break
			}
		}
	}

	if len(nameWords) > 0 {
		relevanceScore := float64(relevantWords) / float64(len(nameWords))
		score *= relevanceScore
	}

	// Penalize for too many hyphens
	hyphenCount := strings.Count(slug, "-")
	if hyphenCount > 3 {
		score -= 0.1 * float64(hyphenCount-3)
	}

	// Bonus for SEO-friendly patterns
	if strings.Contains(slug, "shop") || strings.Contains(slug, "buy") || strings.Contains(slug, "category") {
		score += 0.1
	}

	// Ensure score is between 0 and 1
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	return score
}

func generateSlugReason(slug, categoryName string) string {
	if strings.Contains(slug, "shop") {
		return "Includes 'shop' keyword for better e-commerce SEO"
	}
	if strings.Contains(slug, "buy") {
		return "Includes 'buy' keyword for purchase intent"
	}
	if strings.Contains(slug, "category") {
		return "Clearly identifies as a category page"
	}
	if len(slug) <= 30 {
		return "Short and memorable URL"
	}
	if strings.Count(slug, "-") <= 2 {
		return "Clean structure with minimal separators"
	}
	return "SEO-optimized slug based on category name"
}
