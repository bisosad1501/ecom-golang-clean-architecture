package handlers

import (
	"net/http"
	"strconv"

	"ecom-golang-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProductHandler handles product-related HTTP requests
type ProductHandler struct {
	productUseCase usecases.ProductUseCase
}

// NewProductHandler creates a new product handler
func NewProductHandler(productUseCase usecases.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
	}
}

// CreateProduct handles creating a new product
// @Summary Create a new product
// @Description Create a new product (admin/moderator only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.CreateProductRequest true "Create product request"
// @Success 201 {object} usecases.ProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req usecases.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	product, err := h.productUseCase.CreateProduct(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Product created successfully",
		Data:    product,
	})
}

// GetProduct handles getting a product by ID
// @Summary Get product by ID
// @Description Get a single product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} usecases.ProductResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid product ID",
		})
		return
	}

	product, err := h.productUseCase.GetProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: product,
	})
}

// GetProducts handles getting list of products
// @Summary Get products list
// @Description Get list of products with pagination
// @Tags products
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} usecases.ProductResponse
// @Router /products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	req := usecases.GetProductsRequest{
		Limit:  limit,
		Offset: offset,
	}

	products, err := h.productUseCase.GetProducts(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: products,
	})
}

// SearchProducts handles searching products
// @Summary Search products
// @Description Search products with various filters
// @Tags products
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param category_id query string false "Category ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param status query string false "Product status"
// @Param sort_by query string false "Sort by field" default(created_at)
// @Param sort_order query string false "Sort order" default(desc)
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} usecases.ProductResponse
// @Router /products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	req := usecases.SearchProductsRequest{
		Query:     c.Query("q"),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
		Limit:     10,
		Offset:    0,
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := uuid.Parse(categoryIDStr); err == nil {
			req.CategoryID = &categoryID
		}
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			req.MinPrice = &minPrice
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			req.MaxPrice = &maxPrice
		}
	}

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

	products, err := h.productUseCase.SearchProducts(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: products,
	})
}

// UpdateProduct handles updating a product
// @Summary Update product
// @Description Update an existing product (admin/moderator only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body usecases.UpdateProductRequest true "Update product request"
// @Success 200 {object} usecases.ProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid product ID",
		})
		return
	}

	var req usecases.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	product, err := h.productUseCase.UpdateProduct(c.Request.Context(), productID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product updated successfully",
		Data:    product,
	})
}

// DeleteProduct handles deleting a product
// @Summary Delete product
// @Description Delete a product (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid product ID",
		})
		return
	}

	err = h.productUseCase.DeleteProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product deleted successfully",
	})
}

// GetProductsByCategory handles getting products by category
// @Summary Get products by category
// @Description Get products belonging to a specific category
// @Tags products
// @Accept json
// @Produce json
// @Param categoryId path string true "Category ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} usecases.ProductResponse
// @Failure 400 {object} ErrorResponse
// @Router /products/category/{categoryId} [get]
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("categoryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid category ID",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	products, err := h.productUseCase.GetProductsByCategory(c.Request.Context(), categoryID, limit, offset)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: products,
	})
}

// UpdateStock handles updating product stock
// @Summary Update product stock
// @Description Update product stock quantity (admin/moderator only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body map[string]int true "Stock update request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/products/{id}/stock [put]
func (h *ProductHandler) UpdateStock(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid product ID",
		})
		return
	}

	var req struct {
		Stock int `json:"stock" validate:"min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	err = h.productUseCase.UpdateStock(c.Request.Context(), productID, req.Stock)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product stock updated successfully",
	})
}
