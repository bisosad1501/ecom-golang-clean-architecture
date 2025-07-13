package handlers

import (
	"net/http"
	"strconv"
	"time"

	"ecom-golang-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CategoryHandler handles category-related HTTP requests
type CategoryHandler struct {
	categoryUseCase usecases.CategoryUseCase
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryUseCase usecases.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		categoryUseCase: categoryUseCase,
	}
}

// CreateCategory handles creating a new category
// @Summary Create a new category
// @Description Create a new category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.CreateCategoryRequest true "Create category request"
// @Success 201 {object} usecases.CategoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req usecases.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	category, err := h.categoryUseCase.CreateCategory(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Category created successfully",
		Data:    category,
	})
}

// GetCategory handles getting a category by ID
// @Summary Get category by ID
// @Description Get a single category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} usecases.CategoryResponse
// @Failure 404 {object} ErrorResponse
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	category, err := h.categoryUseCase.GetCategory(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: category,
	})
}

// GetCategories handles getting list of categories
// @Summary Get categories list
// @Description Get list of categories with pagination
// @Tags categories
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} usecases.CategoryResponse
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	req := usecases.GetCategoriesRequest{
		Limit:  limit,
		Offset: offset,
	}

	categories, err := h.categoryUseCase.GetCategories(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: categories,
	})
}

// GetCategoryTree handles getting category tree
// @Summary Get category tree
// @Description Get hierarchical category tree
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} usecases.CategoryResponse
// @Router /categories/tree [get]
func (h *CategoryHandler) GetCategoryTree(c *gin.Context) {
	categories, err := h.categoryUseCase.GetCategoryTree(c.Request.Context())
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: categories,
	})
}

// GetRootCategories handles getting root categories
// @Summary Get root categories
// @Description Get categories that have no parent
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} usecases.CategoryResponse
// @Router /categories/root [get]
func (h *CategoryHandler) GetRootCategories(c *gin.Context) {
	categories, err := h.categoryUseCase.GetRootCategories(c.Request.Context())
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: categories,
	})
}

// GetCategoryChildren handles getting category children
// @Summary Get category children
// @Description Get child categories of a specific category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Parent Category ID"
// @Success 200 {array} usecases.CategoryResponse
// @Failure 400 {object} ErrorResponse
// @Router /categories/{id}/children [get]
func (h *CategoryHandler) GetCategoryChildren(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	categories, err := h.categoryUseCase.GetCategoryChildren(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: categories,
	})
}

// UpdateCategory handles updating a category
// @Summary Update category
// @Description Update an existing category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param request body usecases.UpdateCategoryRequest true "Update category request"
// @Success 200 {object} usecases.CategoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	var req usecases.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	category, err := h.categoryUseCase.UpdateCategory(c.Request.Context(), categoryID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Category updated successfully",
		Data:    category,
	})
}

// DeleteCategory handles deleting a category
// @Summary Delete category
// @Description Delete a category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	err = h.categoryUseCase.DeleteCategory(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Category deleted successfully",
	})
}

// GetCategoryPath handles getting category path from root
// @Summary Get category path
// @Description Get full path from root to specified category (breadcrumbs)
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {array} usecases.CategoryResponse
// @Failure 400 {object} ErrorResponse
// @Router /categories/{id}/path [get]
func (h *CategoryHandler) GetCategoryPath(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	path, err := h.categoryUseCase.GetCategoryPath(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: path,
	})
}

// GetCategoryProductCount handles getting product count for category (including subcategories)
// @Summary Get category product count
// @Description Get total product count for category including all subcategories
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]int64
// @Failure 400 {object} ErrorResponse
// @Router /categories/{id}/product-count [get]
func (h *CategoryHandler) GetCategoryProductCount(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	count, err := h.categoryUseCase.GetCategoryProductCount(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: map[string]int64{
			"product_count": count,
		},
	})
}

// GetCategoryLandingPage handles getting category landing page data
// @Summary Get category landing page
// @Description Get category with products, subcategories, featured products, and navigation data
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Products per page" default(20)
// @Param sort_by query string false "Sort by field" Enums(name,price,created_at,popularity)
// @Param sort_order query string false "Sort order" Enums(asc,desc) default(asc)
// @Param include_subcategory_products query bool false "Include products from subcategories" default(false)
// @Param include_featured query bool false "Include featured products in category" default(false)
// @Param featured_limit query int false "Featured products limit" default(6)
// @Success 200 {object} usecases.CategoryLandingPageResponse
// @Failure 400 {object} ErrorResponse
// @Router /categories/{id}/landing [get]
func (h *CategoryHandler) GetCategoryLandingPage(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	sortBy := c.DefaultQuery("sort_by", "name")
	sortOrder := c.DefaultQuery("sort_order", "asc")
	includeSubcategoryProducts := c.DefaultQuery("include_subcategory_products", "false") == "true"
	includeFeatured := c.DefaultQuery("include_featured", "false") == "true"
	featuredLimit, _ := strconv.Atoi(c.DefaultQuery("featured_limit", "6"))

	req := usecases.GetCategoryLandingPageRequest{
		CategoryID:                 categoryID,
		Page:                      page,
		Limit:                     limit,
		SortBy:                    sortBy,
		SortOrder:                 sortOrder,
		IncludeSubcategoryProducts: includeSubcategoryProducts,
		IncludeFeatured:           includeFeatured,
		FeaturedLimit:             featuredLimit,
	}

	response, err := h.categoryUseCase.GetCategoryLandingPage(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: response,
	})
}

// BulkCreateCategories handles bulk creating categories
// @Summary Bulk create categories
// @Description Create multiple categories at once
// @Tags categories
// @Accept json
// @Produce json
// @Param categories body []usecases.CreateCategoryRequest true "Categories to create"
// @Success 201 {array} usecases.CategoryResponse
// @Router /admin/categories/bulk [post]
func (h *CategoryHandler) BulkCreateCategories(c *gin.Context) {
	var req []usecases.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	categories, err := h.categoryUseCase.BulkCreateCategories(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Data: categories,
	})
}

// BulkUpdateCategories handles bulk updating categories
// @Summary Bulk update categories
// @Description Update multiple categories at once
// @Tags categories
// @Accept json
// @Produce json
// @Param categories body []usecases.BulkUpdateCategoryRequest true "Categories to update"
// @Success 200 {array} usecases.CategoryResponse
// @Router /admin/categories/bulk [put]
func (h *CategoryHandler) BulkUpdateCategories(c *gin.Context) {
	var req []usecases.BulkUpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	categories, err := h.categoryUseCase.BulkUpdateCategories(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: categories,
	})
}

// BulkDeleteCategories handles bulk deleting categories
// @Summary Bulk delete categories
// @Description Delete multiple categories at once
// @Tags categories
// @Accept json
// @Produce json
// @Param request body map[string][]string true "Category IDs to delete"
// @Success 200 {object} SuccessResponse
// @Router /admin/categories/bulk [delete]
func (h *CategoryHandler) BulkDeleteCategories(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var ids []uuid.UUID
	for _, idStr := range req.IDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid category ID: " + idStr,
			})
			return
		}
		ids = append(ids, id)
	}

	err := h.categoryUseCase.BulkDeleteCategories(c.Request.Context(), ids)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Categories deleted successfully",
	})
}

// SearchCategories handles searching categories
// @Summary Search categories
// @Description Search categories by name and description
// @Tags categories
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} usecases.CategoriesListResponse
// @Router /categories/search [get]
func (h *CategoryHandler) SearchCategories(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Search query is required",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	req := usecases.SearchCategoriesRequest{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	}

	result, err := h.categoryUseCase.SearchCategories(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: result,
	})
}

// GetCategoriesWithFilters handles getting categories with advanced filtering
// @Summary Get categories with filters
// @Description Get categories with advanced filtering options
// @Tags categories
// @Accept json
// @Produce json
// @Param search query string false "Search query"
// @Param parent_id query string false "Parent category ID"
// @Param is_active query bool false "Active status"
// @Param has_parent query bool false "Has parent filter"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Param sort_by query string false "Sort by field" default(name)
// @Param sort_order query string false "Sort order" default(asc)
// @Success 200 {object} usecases.CategoriesListResponse
// @Router /categories/filter [get]
func (h *CategoryHandler) GetCategoriesWithFilters(c *gin.Context) {
	req := usecases.GetCategoriesWithFiltersRequest{
		Search:    c.Query("search"),
		SortBy:    c.DefaultQuery("sort_by", "name"),
		SortOrder: c.DefaultQuery("sort_order", "asc"),
	}

	// Parse parent_id if provided
	if parentIDStr := c.Query("parent_id"); parentIDStr != "" {
		parentID, err := uuid.Parse(parentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid parent_id",
			})
			return
		}
		req.ParentID = &parentID
	}

	// Parse is_active if provided
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		req.IsActive = &isActive
	}

	// Parse has_parent if provided
	if hasParentStr := c.Query("has_parent"); hasParentStr != "" {
		hasParent := hasParentStr == "true"
		req.HasParent = &hasParent
	}

	// Parse pagination
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = offset
		}
	}

	result, err := h.categoryUseCase.GetCategoriesWithFilters(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: result,
	})
}

// MoveCategory handles moving a category to a new parent
// @Summary Move category to new parent
// @Description Move a category to a new parent (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.MoveCategoryRequest true "Move category request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/categories/move [post]
func (h *CategoryHandler) MoveCategory(c *gin.Context) {
	var req usecases.MoveCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	err := h.categoryUseCase.MoveCategory(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Category moved successfully",
	})
}

// ReorderCategories handles reordering multiple categories
// @Summary Reorder categories
// @Description Reorder multiple categories (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.ReorderCategoriesRequest true "Reorder categories request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/categories/reorder [post]
func (h *CategoryHandler) ReorderCategories(c *gin.Context) {
	var req usecases.ReorderCategoriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	err := h.categoryUseCase.ReorderCategories(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Categories reordered successfully",
	})
}

// GetCategoryTreeStats handles getting category tree statistics
// @Summary Get category tree statistics
// @Description Get statistics about the category tree (admin only)
// @Tags categories
// @Produce json
// @Security BearerAuth
// @Success 200 {object} usecases.CategoryTreeStatsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/categories/tree/stats [get]
func (h *CategoryHandler) GetCategoryTreeStats(c *gin.Context) {
	stats, err := h.categoryUseCase.GetCategoryTreeStats(c.Request.Context())
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: stats,
	})
}

// ValidateAndRepairTree handles validating and repairing the category tree
// @Summary Validate and repair category tree
// @Description Validate the category tree integrity and perform repairs if needed (admin only)
// @Tags categories
// @Produce json
// @Security BearerAuth
// @Success 200 {object} usecases.TreeValidationResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/categories/tree/validate [post]
func (h *CategoryHandler) ValidateAndRepairTree(c *gin.Context) {
	result, err := h.categoryUseCase.ValidateAndRepairTree(c.Request.Context())
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: result,
	})
}

// GetCategoryAnalytics handles getting comprehensive category analytics
// @Summary Get category analytics
// @Description Get comprehensive analytics for a category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param time_range query string false "Time range" default(30d)
// @Success 200 {object} usecases.CategoryAnalyticsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/categories/{id}/analytics [get]
func (h *CategoryHandler) GetCategoryAnalytics(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	req := usecases.GetCategoryAnalyticsRequest{
		CategoryID: categoryID,
		TimeRange:  c.DefaultQuery("time_range", "30d"),
	}

	analytics, err := h.categoryUseCase.GetCategoryAnalytics(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: analytics,
	})
}

// GetTopCategories handles getting top performing categories
// @Summary Get top categories
// @Description Get top performing categories (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param sort_by query string false "Sort by" default(sales)
// @Success 200 {object} usecases.TopCategoriesResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/categories/top [get]
func (h *CategoryHandler) GetTopCategories(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sort_by", "sales")

	req := usecases.GetTopCategoriesRequest{
		Limit:  limit,
		SortBy: sortBy,
	}

	result, err := h.categoryUseCase.GetTopCategories(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: result,
	})
}

// GetCategoryPerformanceMetrics handles getting detailed category performance metrics
// @Summary Get category performance metrics
// @Description Get detailed performance metrics for a category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} usecases.CategoryPerformanceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/categories/{id}/performance [get]
func (h *CategoryHandler) GetCategoryPerformanceMetrics(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	metrics, err := h.categoryUseCase.GetCategoryPerformanceMetrics(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: metrics,
	})
}

// GetCategorySalesStats handles getting category sales statistics
// @Summary Get category sales statistics
// @Description Get sales statistics for a category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param time_range query string false "Time range" default(30d)
// @Success 200 {object} usecases.CategorySalesStatsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/categories/{id}/sales [get]
func (h *CategoryHandler) GetCategorySalesStats(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	req := usecases.GetCategorySalesStatsRequest{
		CategoryID: categoryID,
		TimeRange:  c.DefaultQuery("time_range", "30d"),
	}

	stats, err := h.categoryUseCase.GetCategorySalesStats(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: stats,
	})
}

// UpdateCategorySEO handles updating SEO metadata for a category
// @Summary Update category SEO
// @Description Update SEO metadata for a category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param request body usecases.CategorySEORequest true "SEO update request"
// @Success 200 {object} usecases.CategoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/categories/{id}/seo [put]
func (h *CategoryHandler) UpdateCategorySEO(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	var req usecases.CategorySEORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	category, err := h.categoryUseCase.UpdateCategorySEO(c.Request.Context(), categoryID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: category,
	})
}

// GetCategorySEO handles getting SEO metadata for a category
// @Summary Get category SEO
// @Description Get SEO metadata for a category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} usecases.CategorySEOResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /categories/{id}/seo [get]
func (h *CategoryHandler) GetCategorySEO(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	seo, err := h.categoryUseCase.GetCategorySEO(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: seo,
	})
}

// GenerateCategorySEO handles generating SEO metadata for a category
// @Summary Generate category SEO
// @Description Auto-generate SEO metadata for a category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} usecases.CategorySEOResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/categories/{id}/seo/generate [post]
func (h *CategoryHandler) GenerateCategorySEO(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	seo, err := h.categoryUseCase.GenerateCategorySEO(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: seo,
	})
}

// ValidateCategorySEO handles validating SEO metadata for a category
// @Summary Validate category SEO
// @Description Validate SEO metadata for a category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} usecases.CategorySEOValidationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/categories/{id}/seo/validate [get]
func (h *CategoryHandler) ValidateCategorySEO(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	validation, err := h.categoryUseCase.ValidateCategorySEO(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: validation,
	})
}

// OptimizeSlug handles optimizing category slug for better SEO
// @Summary Optimize category slug
// @Description Optimize category slug for better SEO (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param request body usecases.SlugOptimizationRequest true "Slug optimization request"
// @Success 200 {object} usecases.SlugOptimizationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/categories/{id}/slug/optimize [post]
func (h *CategoryHandler) OptimizeSlug(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	var req usecases.SlugOptimizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	result, err := h.categoryUseCase.OptimizeSlug(c.Request.Context(), categoryID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: result,
	})
}

// GenerateSlugSuggestions handles generating slug suggestions
// @Summary Generate slug suggestions
// @Description Generate SEO-friendly slug suggestions for a category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} usecases.SlugSuggestionsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/categories/{id}/slug/suggestions [get]
func (h *CategoryHandler) GenerateSlugSuggestions(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	suggestions, err := h.categoryUseCase.GenerateSlugSuggestions(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: suggestions,
	})
}

// ValidateSlugAvailability handles validating slug availability
// @Summary Validate slug availability
// @Description Check if a slug is available and SEO-friendly
// @Tags categories
// @Accept json
// @Produce json
// @Param slug query string true "Slug to validate"
// @Param exclude_id query string false "Category ID to exclude from check"
// @Success 200 {object} usecases.SlugValidationResponse
// @Failure 400 {object} ErrorResponse
// @Router /categories/slug/validate [get]
func (h *CategoryHandler) ValidateSlugAvailability(c *gin.Context) {
	slug := c.Query("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Slug parameter is required",
		})
		return
	}

	var excludeID *uuid.UUID
	if excludeIDStr := c.Query("exclude_id"); excludeIDStr != "" {
		if id, err := uuid.Parse(excludeIDStr); err == nil {
			excludeID = &id
		}
	}

	validation, err := h.categoryUseCase.ValidateSlugAvailability(c.Request.Context(), slug, excludeID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: validation,
	})
}

// GetSlugHistory handles getting slug change history
// @Summary Get slug history
// @Description Get slug change history for a category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} usecases.SlugHistoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/categories/{id}/slug/history [get]
func (h *CategoryHandler) GetSlugHistory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	history, err := h.categoryUseCase.GetSlugHistory(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: history,
	})
}

// BulkUpdateSEO handles bulk SEO updates for multiple categories
// @Summary Bulk update SEO
// @Description Update SEO metadata for multiple categories (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.BulkSEOUpdateRequest true "Bulk SEO update request"
// @Success 200 {object} usecases.BulkSEOUpdateResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/categories/seo/bulk-update [post]
func (h *CategoryHandler) BulkUpdateSEO(c *gin.Context) {
	var req usecases.BulkSEOUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	result, err := h.categoryUseCase.BulkUpdateSEO(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: result,
	})
}

// BulkGenerateSEO handles bulk SEO generation for multiple categories
// @Summary Bulk generate SEO
// @Description Generate SEO metadata for multiple categories (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.BulkSEOGenerateRequest true "Bulk SEO generate request"
// @Success 200 {object} usecases.BulkSEOGenerateResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/categories/seo/bulk-generate [post]
func (h *CategoryHandler) BulkGenerateSEO(c *gin.Context) {
	var req usecases.BulkSEOGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	result, err := h.categoryUseCase.BulkGenerateSEO(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: result,
	})
}

// BulkValidateSEO handles bulk SEO validation for multiple categories
// @Summary Bulk validate SEO
// @Description Validate SEO metadata for multiple categories (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.BulkSEOValidateRequest true "Bulk SEO validate request"
// @Success 200 {object} usecases.BulkSEOValidateResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/categories/seo/bulk-validate [post]
func (h *CategoryHandler) BulkValidateSEO(c *gin.Context) {
	var req usecases.BulkSEOValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	result, err := h.categoryUseCase.BulkValidateSEO(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: result,
	})
}

// GetSEOAnalytics handles getting comprehensive SEO analytics
// @Summary Get SEO analytics
// @Description Get comprehensive SEO analytics across categories (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category_ids query []string false "Category IDs to analyze"
// @Param date_from query string false "Start date for analysis (YYYY-MM-DD)"
// @Param date_to query string false "End date for analysis (YYYY-MM-DD)"
// @Param metrics query []string false "Specific metrics to include"
// @Success 200 {object} usecases.SEOAnalyticsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/categories/seo/analytics [get]
func (h *CategoryHandler) GetSEOAnalytics(c *gin.Context) {
	req := usecases.SEOAnalyticsRequest{}

	// Parse category IDs
	if categoryIDsStr := c.QueryArray("category_ids"); len(categoryIDsStr) > 0 {
		for _, idStr := range categoryIDsStr {
			if id, err := uuid.Parse(idStr); err == nil {
				req.CategoryIDs = append(req.CategoryIDs, id)
			}
		}
	}

	// Parse date range
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if dateFrom, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			req.DateFrom = &dateFrom
		}
	}

	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if dateTo, err := time.Parse("2006-01-02", dateToStr); err == nil {
			req.DateTo = &dateTo
		}
	}

	// Parse metrics
	req.Metrics = c.QueryArray("metrics")

	analytics, err := h.categoryUseCase.GetSEOAnalytics(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: analytics,
	})
}

// GetSEOInsights handles getting detailed SEO insights for a category
// @Summary Get SEO insights
// @Description Get detailed SEO insights and recommendations for a category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} usecases.SEOInsightsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/categories/{id}/seo/insights [get]
func (h *CategoryHandler) GetSEOInsights(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	insights, err := h.categoryUseCase.GetSEOInsights(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: insights,
	})
}

// GetSEOCompetitorAnalysis handles getting competitor analysis for SEO
// @Summary Get SEO competitor analysis
// @Description Get competitor analysis and benchmarks for category SEO (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} usecases.SEOCompetitorAnalysisResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/categories/{id}/seo/competitor-analysis [get]
func (h *CategoryHandler) GetSEOCompetitorAnalysis(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	analysis, err := h.categoryUseCase.GetSEOCompetitorAnalysis(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: analysis,
	})
}
