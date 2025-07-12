package handlers

import (
	"net/http"
	"strconv"

	"ecom-golang-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BrandHandler handles brand-related HTTP requests
type BrandHandler struct {
	brandUseCase usecases.BrandUseCase
}

// NewBrandHandler creates a new brand handler
func NewBrandHandler(brandUseCase usecases.BrandUseCase) *BrandHandler {
	return &BrandHandler{
		brandUseCase: brandUseCase,
	}
}

// CreateBrand handles brand creation
// @Summary Create brand
// @Description Create a new brand
// @Tags brands
// @Accept json
// @Produce json
// @Param request body usecases.CreateBrandRequest true "Create brand request"
// @Success 201 {object} usecases.BrandResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /brands [post]
func (h *BrandHandler) CreateBrand(c *gin.Context) {
	var req usecases.CreateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	brand, err := h.brandUseCase.CreateBrand(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "resource conflict" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Brand with this slug already exists",
			})
			return
		}
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Brand created successfully",
		Data:    brand,
	})
}

// GetBrand handles getting a single brand
// @Summary Get brand
// @Description Get brand by ID
// @Tags brands
// @Accept json
// @Produce json
// @Param id path string true "Brand ID"
// @Success 200 {object} usecases.BrandResponse
// @Failure 404 {object} ErrorResponse
// @Router /brands/{id} [get]
func (h *BrandHandler) GetBrand(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid brand ID format",
		})
		return
	}

	brand, err := h.brandUseCase.GetBrand(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Brand not found",
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: brand,
	})
}

// GetBrandBySlug handles getting a brand by slug
// @Summary Get brand by slug
// @Description Get brand by slug
// @Tags brands
// @Accept json
// @Produce json
// @Param slug path string true "Brand slug"
// @Success 200 {object} usecases.BrandResponse
// @Failure 404 {object} ErrorResponse
// @Router /brands/slug/{slug} [get]
func (h *BrandHandler) GetBrandBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Brand slug is required",
		})
		return
	}

	brand, err := h.brandUseCase.GetBrandBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Brand not found",
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: brand,
	})
}

// GetBrands handles getting list of brands
// @Summary Get brands list
// @Description Get list of brands with pagination
// @Tags brands
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} usecases.BrandsListResponse
// @Router /brands [get]
func (h *BrandHandler) GetBrands(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	var isActive *bool
	if activeStr := c.Query("is_active"); activeStr != "" {
		if active, err := strconv.ParseBool(activeStr); err == nil {
			isActive = &active
		}
	}

	req := usecases.GetBrandsRequest{
		Limit:    limit,
		Offset:   offset,
		IsActive: isActive,
	}

	brands, err := h.brandUseCase.GetBrands(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: brands,
	})
}

// SearchBrands handles brand search
// @Summary Search brands
// @Description Search brands by name or description
// @Tags brands
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} usecases.BrandsListResponse
// @Router /brands/search [get]
func (h *BrandHandler) SearchBrands(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Search query is required",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	req := usecases.SearchBrandsRequest{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	}

	brands, err := h.brandUseCase.SearchBrands(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: brands,
	})
}

// GetActiveBrands handles getting active brands
// @Summary Get active brands
// @Description Get list of active brands
// @Tags brands
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} usecases.BrandsListResponse
// @Router /brands/active [get]
func (h *BrandHandler) GetActiveBrands(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	brands, err := h.brandUseCase.GetActiveBrands(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: brands,
	})
}

// GetPopularBrands handles getting popular brands
// @Summary Get popular brands
// @Description Get list of popular brands by product count
// @Tags brands
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} usecases.BrandResponse
// @Router /brands/popular [get]
func (h *BrandHandler) GetPopularBrands(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	brands, err := h.brandUseCase.GetPopularBrands(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: brands,
	})
}

// UpdateBrand handles brand update
// @Summary Update brand
// @Description Update an existing brand
// @Tags brands
// @Accept json
// @Produce json
// @Param id path string true "Brand ID"
// @Param request body usecases.UpdateBrandRequest true "Update brand request"
// @Success 200 {object} usecases.BrandResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /brands/{id} [put]
func (h *BrandHandler) UpdateBrand(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid brand ID format",
		})
		return
	}

	var req usecases.UpdateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	brand, err := h.brandUseCase.UpdateBrand(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "brand not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Brand not found",
			})
			return
		}
		if err.Error() == "resource conflict" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Brand with this slug already exists",
			})
			return
		}
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Brand updated successfully",
		Data:    brand,
	})
}

// DeleteBrand handles brand deletion
// @Summary Delete brand
// @Description Delete a brand
// @Tags brands
// @Accept json
// @Produce json
// @Param id path string true "Brand ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /brands/{id} [delete]
func (h *BrandHandler) DeleteBrand(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid brand ID format",
		})
		return
	}

	err = h.brandUseCase.DeleteBrand(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "brand not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Brand not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Brand deleted successfully",
	})
}
