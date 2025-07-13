package handlers

import (
	"ecom-golang-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
)

// ProductComparisonHandler handles product comparison HTTP requests
type ProductComparisonHandler struct {
	comparisonUseCase usecases.ProductComparisonUseCase
}

// NewProductComparisonHandler creates a new product comparison handler
func NewProductComparisonHandler(comparisonUseCase usecases.ProductComparisonUseCase) *ProductComparisonHandler {
	return &ProductComparisonHandler{
		comparisonUseCase: comparisonUseCase,
	}
}

// CreateComparison creates a new product comparison
// @Summary Create product comparison
// @Description Create a new product comparison for user or session
// @Tags product-comparison
// @Accept json
// @Produce json
// @Param request body usecases.ProductComparisonRequest true "Comparison request"
// @Success 201 {object} usecases.ProductComparisonResponse
// @Failure 400 {object} ErrorResponse
// @Router /products/compare [post]
func (h *ProductComparisonHandler) CreateComparison(c *gin.Context) {
	var req usecases.ProductComparisonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Get user ID from context (optional)
	var userID *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			userID = &uid
		}
	}

	// Get session ID from header or generate one
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	comparison, err := h.comparisonUseCase.CreateComparison(c.Request.Context(), userID, sessionID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Comparison created successfully",
		Data:    comparison,
	})
}

// GetComparison gets a comparison by ID
// @Summary Get product comparison
// @Description Get a product comparison by ID
// @Tags product-comparison
// @Produce json
// @Param id path string true "Comparison ID"
// @Success 200 {object} usecases.ProductComparisonResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/compare/{id} [get]
func (h *ProductComparisonHandler) GetComparison(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid comparison ID",
		})
		return
	}

	comparison, err := h.comparisonUseCase.GetComparison(c.Request.Context(), id)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Comparison retrieved successfully",
		Data:    comparison,
	})
}

// GetUserComparison gets user's comparison
// @Summary Get user's product comparison
// @Description Get the current user's product comparison
// @Tags product-comparison
// @Produce json
// @Security BearerAuth
// @Success 200 {object} usecases.ProductComparisonResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/compare/user [get]
func (h *ProductComparisonHandler) GetUserComparison(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	comparison, err := h.comparisonUseCase.GetUserComparison(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "User comparison retrieved successfully",
		Data:    comparison,
	})
}

// GetSessionComparison gets session's comparison
// @Summary Get session's product comparison
// @Description Get the current session's product comparison
// @Tags product-comparison
// @Produce json
// @Param X-Session-ID header string true "Session ID"
// @Success 200 {object} usecases.ProductComparisonResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/compare/session [get]
func (h *ProductComparisonHandler) GetSessionComparison(c *gin.Context) {
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Session ID is required",
		})
		return
	}

	comparison, err := h.comparisonUseCase.GetSessionComparison(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Session comparison retrieved successfully",
		Data:    comparison,
	})
}

// UpdateComparison updates a comparison
// @Summary Update product comparison
// @Description Update a product comparison
// @Tags product-comparison
// @Accept json
// @Produce json
// @Param id path string true "Comparison ID"
// @Param request body usecases.ProductComparisonRequest true "Comparison request"
// @Success 200 {object} usecases.ProductComparisonResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/compare/{id} [put]
func (h *ProductComparisonHandler) UpdateComparison(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid comparison ID",
		})
		return
	}

	var req usecases.ProductComparisonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	comparison, err := h.comparisonUseCase.UpdateComparison(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Comparison updated successfully",
		Data:    comparison,
	})
}

// DeleteComparison deletes a comparison
// @Summary Delete product comparison
// @Description Delete a product comparison
// @Tags product-comparison
// @Produce json
// @Param id path string true "Comparison ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/compare/{id} [delete]
func (h *ProductComparisonHandler) DeleteComparison(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid comparison ID",
		})
		return
	}

	err = h.comparisonUseCase.DeleteComparison(c.Request.Context(), id)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Comparison deleted successfully",
	})
}

// AddProductToComparison adds a product to comparison
// @Summary Add product to comparison
// @Description Add a product to an existing comparison
// @Tags product-comparison
// @Produce json
// @Param id path string true "Comparison ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} usecases.ProductComparisonResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/compare/{id}/products/{product_id} [post]
func (h *ProductComparisonHandler) AddProductToComparison(c *gin.Context) {
	comparisonIDStr := c.Param("id")
	comparisonID, err := uuid.Parse(comparisonIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid comparison ID",
		})
		return
	}

	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid product ID",
		})
		return
	}

	comparison, err := h.comparisonUseCase.AddProductToComparison(c.Request.Context(), comparisonID, productID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product added to comparison successfully",
		Data:    comparison,
	})
}

// RemoveProductFromComparison removes a product from comparison
// @Summary Remove product from comparison
// @Description Remove a product from an existing comparison
// @Tags product-comparison
// @Produce json
// @Param id path string true "Comparison ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} usecases.ProductComparisonResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/compare/{id}/products/{product_id} [delete]
func (h *ProductComparisonHandler) RemoveProductFromComparison(c *gin.Context) {
	comparisonIDStr := c.Param("id")
	comparisonID, err := uuid.Parse(comparisonIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid comparison ID",
		})
		return
	}

	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid product ID",
		})
		return
	}

	comparison, err := h.comparisonUseCase.RemoveProductFromComparison(c.Request.Context(), comparisonID, productID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product removed from comparison successfully",
		Data:    comparison,
	})
}

// ClearComparison clears all products from comparison
// @Summary Clear comparison
// @Description Remove all products from a comparison
// @Tags product-comparison
// @Produce json
// @Param id path string true "Comparison ID"
// @Success 200 {object} usecases.ProductComparisonResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/compare/{id}/clear [post]
func (h *ProductComparisonHandler) ClearComparison(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid comparison ID",
		})
		return
	}

	comparison, err := h.comparisonUseCase.ClearComparison(c.Request.Context(), id)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Comparison cleared successfully",
		Data:    comparison,
	})
}

// CompareProducts compares products directly without creating a comparison
// @Summary Compare products
// @Description Compare multiple products and get comparison matrix
// @Tags product-comparison
// @Produce json
// @Param product_ids query string true "Comma-separated product IDs"
// @Success 200 {object} usecases.ComparisonMatrixResponse
// @Failure 400 {object} ErrorResponse
// @Router /products/compare/matrix [get]
func (h *ProductComparisonHandler) CompareProducts(c *gin.Context) {
	productIDsStr := c.Query("product_ids")
	if productIDsStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "product_ids parameter is required",
		})
		return
	}

	// Parse product IDs
	productIDStrings := strings.Split(productIDsStr, ",")
	productIDs := make([]uuid.UUID, len(productIDStrings))

	for i, idStr := range productIDStrings {
		id, err := uuid.Parse(strings.TrimSpace(idStr))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid product ID: " + idStr,
			})
			return
		}
		productIDs[i] = id
	}

	matrix, err := h.comparisonUseCase.CompareProducts(c.Request.Context(), productIDs)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Products compared successfully",
		Data:    matrix,
	})
}

// GetComparisonMatrix gets comparison matrix for existing comparison
// @Summary Get comparison matrix
// @Description Get comparison matrix for an existing comparison
// @Tags product-comparison
// @Produce json
// @Param id path string true "Comparison ID"
// @Success 200 {object} usecases.ComparisonMatrixResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/compare/{id}/matrix [get]
func (h *ProductComparisonHandler) GetComparisonMatrix(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid comparison ID",
		})
		return
	}

	matrix, err := h.comparisonUseCase.GetComparisonMatrix(c.Request.Context(), id)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Comparison matrix retrieved successfully",
		Data:    matrix,
	})
}

// GetPopularComparedProducts gets most compared products
// @Summary Get popular compared products
// @Description Get the most frequently compared products
// @Tags product-comparison
// @Produce json
// @Param limit query int false "Number of products to return" default(10)
// @Success 200 {object} []usecases.ProductResponse
// @Failure 400 {object} ErrorResponse
// @Router /products/compare/popular [get]
func (h *ProductComparisonHandler) GetPopularComparedProducts(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	products, err := h.comparisonUseCase.GetPopularComparedProducts(c.Request.Context(), limit)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Popular compared products retrieved successfully",
		Data:    products,
	})
}
