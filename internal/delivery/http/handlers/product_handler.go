package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	// Parse product ID
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid product ID format",
		})
		return
	}

	// Get user info from middleware
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	
	fmt.Printf("UpdateProduct: ProductID=%s, UserID=%v, Role=%v\n", productID.String(), userID, role)

	// Parse the request
	var req usecases.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("UpdateProduct: JSON binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Validate the request
	if err := h.validateUpdateProductRequest(&req); err != nil {
		fmt.Printf("UpdateProduct: Validation error: %v\n", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Validation failed: " + err.Error(),
		})
		return
	}

	// Log the request data
	fmt.Printf("UpdateProduct: Request data: %+v\n", req)

	// Use the clean implementation
	product, err := h.productUseCase.UpdateProduct(c.Request.Context(), productID, req)
	if err != nil {
		fmt.Printf("UpdateProduct: UseCase error: %v\n", err)
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	fmt.Printf("UpdateProduct: Success for ProductID=%s\n", productID.String())
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product updated successfully",
		Data:    product,
	})
}

// PatchProduct handles partially updating a product
// @Summary Partially update product
// @Description Partially update an existing product - only updates provided fields (admin/moderator only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body usecases.PatchProductRequest true "Patch product request"
// @Success 200 {object} usecases.ProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/products/{id} [patch]
func (h *ProductHandler) PatchProduct(c *gin.Context) {
	// Parse product ID
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid product ID format",
		})
		return
	}

	// Get user info from middleware
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	
	fmt.Printf("PatchProduct: ProductID=%s, UserID=%v, Role=%v\n", productID.String(), userID, role)

	// Parse the request
	var req usecases.PatchProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("PatchProduct: JSON binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Validate the request
	if err := h.validatePatchProductRequest(&req); err != nil {
		fmt.Printf("PatchProduct: Validation error: %v\n", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Validation failed: " + err.Error(),
		})
		return
	}

	// Log the request data
	fmt.Printf("PatchProduct: Request data: %+v\n", req)

	// Use the patch implementation
	product, err := h.productUseCase.PatchProduct(c.Request.Context(), productID, req)
	if err != nil {
		fmt.Printf("PatchProduct: UseCase error: %v\n", err)
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	fmt.Printf("PatchProduct: Success for ProductID=%s\n", productID.String())
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

// validateUpdateProductRequest validates the update product request
func (h *ProductHandler) validateUpdateProductRequest(req *usecases.UpdateProductRequest) error {
	// Validate name
	if req.Name != nil {
		if len(strings.TrimSpace(*req.Name)) == 0 {
			return fmt.Errorf("name cannot be empty")
		}
		if len(*req.Name) > 255 {
			return fmt.Errorf("name cannot exceed 255 characters")
		}
	}
	
	// Validate description
	if req.Description != nil && len(*req.Description) > 2000 {
		return fmt.Errorf("description cannot exceed 2000 characters")
	}
	
	// Validate price fields (validation moved to usecase for better business logic)
	// Keep basic validation here for early feedback
	if req.Price != nil && *req.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	if req.ComparePrice != nil && *req.ComparePrice <= 0 {
		return fmt.Errorf("compare price must be greater than 0")
	}
	if req.CostPrice != nil && *req.CostPrice < 0 {
		return fmt.Errorf("cost price cannot be negative")
	}
	
	// Validate stock
	if req.Stock != nil && *req.Stock < 0 {
		return fmt.Errorf("stock cannot be negative")
	}
	
	// Validate weight
	if req.Weight != nil && *req.Weight <= 0 {
		return fmt.Errorf("weight must be greater than 0")
	}
	
	// Validate dimensions
	if req.Dimensions != nil {
		if req.Dimensions.Length <= 0 {
			return fmt.Errorf("dimension length must be greater than 0")
		}
		if req.Dimensions.Width <= 0 {
			return fmt.Errorf("dimension width must be greater than 0")
		}
		if req.Dimensions.Height <= 0 {
			return fmt.Errorf("dimension height must be greater than 0")
		}
	}
	
	// Validate images
	if req.Images != nil {
		if len(req.Images) > 10 { // Reasonable limit
			return fmt.Errorf("cannot have more than 10 images per product")
		}
		for i, img := range req.Images {
			if strings.TrimSpace(img.URL) == "" {
				return fmt.Errorf("image URL cannot be empty at position %d", i+1)
			}
			if len(img.URL) > 500 {
				return fmt.Errorf("image URL too long at position %d", i+1)
			}
			if len(img.AltText) > 255 {
				return fmt.Errorf("image alt text too long at position %d", i+1)
			}
		}
	}
	
	// Validate tags
	if req.Tags != nil {
		if len(req.Tags) > 20 { // Reasonable limit
			return fmt.Errorf("cannot have more than 20 tags per product")
		}
		for i, tag := range req.Tags {
			if len(strings.TrimSpace(tag)) == 0 {
				continue // Skip empty tags, they'll be filtered out
			}
			if len(tag) > 50 {
				return fmt.Errorf("tag too long at position %d (max 50 characters)", i+1)
			}
		}
	}
	
	// Validate that at least one field is being updated
	if req.Name == nil && req.Description == nil && req.Price == nil && 
	   req.ComparePrice == nil && req.CostPrice == nil && req.Stock == nil &&
	   req.Weight == nil && req.CategoryID == nil && req.Status == nil &&
	   req.IsDigital == nil && req.Dimensions == nil && req.Images == nil &&
	   req.Tags == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}
	
	return nil
}

// validatePatchProductRequest validates the patch product request
func (h *ProductHandler) validatePatchProductRequest(req *usecases.PatchProductRequest) error {
	// Validate name if provided
	if req.Name != nil {
		if len(strings.TrimSpace(*req.Name)) == 0 {
			return fmt.Errorf("name cannot be empty")
		}
		if len(*req.Name) > 255 {
			return fmt.Errorf("name cannot exceed 255 characters")
		}
	}
	
	// Validate description if provided
	if req.Description != nil && len(*req.Description) > 2000 {
		return fmt.Errorf("description cannot exceed 2000 characters")
	}
	
	// Validate price fields if provided
	if req.Price != nil && *req.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	if req.ComparePrice != nil && *req.ComparePrice <= 0 {
		return fmt.Errorf("compare price must be greater than 0")
	}
	if req.CostPrice != nil && *req.CostPrice < 0 {
		return fmt.Errorf("cost price cannot be negative")
	}
	
	// Validate stock if provided
	if req.Stock != nil && *req.Stock < 0 {
		return fmt.Errorf("stock cannot be negative")
	}
	
	// Validate weight if provided
	if req.Weight != nil && *req.Weight <= 0 {
		return fmt.Errorf("weight must be greater than 0")
	}
	
	// Validate dimensions if provided
	if req.Dimensions != nil {
		if req.Dimensions.Length <= 0 {
			return fmt.Errorf("dimension length must be greater than 0")
		}
		if req.Dimensions.Width <= 0 {
			return fmt.Errorf("dimension width must be greater than 0")
		}
		if req.Dimensions.Height <= 0 {
			return fmt.Errorf("dimension height must be greater than 0")
		}
	}
	
	// Validate images if provided
	if req.Images != nil {
		if len(*req.Images) > 10 { // Reasonable limit
			return fmt.Errorf("cannot have more than 10 images per product")
		}
		for i, img := range *req.Images {
			if strings.TrimSpace(img.URL) == "" {
				return fmt.Errorf("image URL cannot be empty at position %d", i+1)
			}
			if len(img.URL) > 500 {
				return fmt.Errorf("image URL too long at position %d", i+1)
			}
			if len(img.AltText) > 255 {
				return fmt.Errorf("image alt text too long at position %d", i+1)
			}
		}
	}
	
	// Validate tags if provided
	if req.Tags != nil {
		if len(*req.Tags) > 20 { // Reasonable limit
			return fmt.Errorf("cannot have more than 20 tags per product")
		}
		for i, tag := range *req.Tags {
			if len(strings.TrimSpace(tag)) == 0 {
				continue // Skip empty tags, they'll be filtered out
			}
			if len(tag) > 50 {
				return fmt.Errorf("tag too long at position %d (max 50 characters)", i+1)
			}
		}
	}
	
	// For PATCH, we don't require at least one field (unlike PUT)
	// This allows for edge cases where user might want to trigger validation only
	
	return nil
}

// GetFeaturedProducts handles getting featured products
// @Summary Get featured products
// @Description Get list of featured products
// @Tags products
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} usecases.ProductResponse
// @Router /products/featured [get]
func (h *ProductHandler) GetFeaturedProducts(c *gin.Context) {
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
		Message: "Featured products retrieved successfully",
		Data:    products,
	})
}

// GetTrendingProducts handles getting trending products
// @Summary Get trending products
// @Description Get list of trending products
// @Tags products
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} usecases.ProductResponse
// @Router /products/trending [get]
func (h *ProductHandler) GetTrendingProducts(c *gin.Context) {
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
		Message: "Trending products retrieved successfully",
		Data:    products,
	})
}

// GetRelatedProducts handles getting products related to a specific product
// @Summary Get related products
// @Description Get products related to a specific product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} usecases.ProductResponse
// @Router /products/{id}/related [get]
func (h *ProductHandler) GetRelatedProducts(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid product ID",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Get the product first to find its category
	product, err := h.productUseCase.GetProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Use category to find related products
	products, err := h.productUseCase.GetProductsByCategory(c.Request.Context(), product.Category.ID, limit, 0)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Related products retrieved successfully",
		Data:    products,
	})
}

// GetProductFilters handles getting product filters for faceted search
// @Summary Get product filters
// @Description Get available filters for products including brands, price range, and attributes
// @Tags products
// @Accept json
// @Produce json
// @Param category_id query string false "Category ID to filter brands"
// @Success 200 {object} map[string]interface{}
// @Router /products/filters [get]
func (h *ProductHandler) GetProductFilters(c *gin.Context) {
	categoryIDStr := c.Query("category_id")
	var categoryID *uuid.UUID

	if categoryIDStr != "" {
		if id, err := uuid.Parse(categoryIDStr); err == nil {
			categoryID = &id
		}
	}

	// Get brand filters (this would need to be implemented in product use case)
	// For now, we'll return a basic structure
	// TODO: Use categoryID to filter brands by category
	_ = categoryID // Suppress unused variable warning

	filters := map[string]interface{}{
		"price_range": map[string]interface{}{
			"min": 0,
			"max": 10000,
		},
		"brands": []map[string]interface{}{},
		"attributes": []map[string]interface{}{},
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product filters retrieved successfully",
		Data:    filters,
	})
}

