package handlers

import (
	"net/http"
	"strconv"

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
