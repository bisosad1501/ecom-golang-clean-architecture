package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"ecom-golang-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProductFilterHandler handles advanced product filtering requests
type ProductFilterHandler struct {
	filterUseCase usecases.ProductFilterUseCase
}

// NewProductFilterHandler creates a new product filter handler
func NewProductFilterHandler(filterUseCase usecases.ProductFilterUseCase) *ProductFilterHandler {
	return &ProductFilterHandler{
		filterUseCase: filterUseCase,
	}
}

// FilterProducts handles advanced product filtering
func (h *ProductFilterHandler) FilterProducts(c *gin.Context) {
	var req usecases.AdvancedFilterRequest

	// Parse query parameters
	req.Query = c.Query("query")
	req.CategoryIDs = parseStringSlice(c.Query("category_ids"))

	// Support both brand_id and brand_ids
	brandIDs := parseStringSlice(c.Query("brand_ids"))
	if brandID := c.Query("brand_id"); brandID != "" {
		brandIDs = append(brandIDs, brandID)
	}
	req.BrandIDs = brandIDs

	req.Tags = parseStringSlice(c.Query("tags"))
	req.ProductTypes = parseStringSlice(c.Query("product_types"))
	req.StockStatus = parseStringSlice(c.Query("stock_status"))
	req.Visibility = parseStringSlice(c.Query("visibility"))

	// Parse price filters
	if minPrice := c.Query("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			req.MinPrice = &price
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			req.MaxPrice = &price
		}
	}

	// Parse rating filters
	if minRating := c.Query("min_rating"); minRating != "" {
		if rating, err := strconv.ParseFloat(minRating, 64); err == nil {
			req.MinRating = &rating
		}
	}
	if maxRating := c.Query("max_rating"); maxRating != "" {
		if rating, err := strconv.ParseFloat(maxRating, 64); err == nil {
			req.MaxRating = &rating
		}
	}

	// Parse boolean filters
	if inStock := c.Query("in_stock"); inStock != "" {
		if val, err := strconv.ParseBool(inStock); err == nil {
			req.InStock = &val
		}
	}
	if lowStock := c.Query("low_stock"); lowStock != "" {
		if val, err := strconv.ParseBool(lowStock); err == nil {
			req.LowStock = &val
		}
	}
	if onSale := c.Query("on_sale"); onSale != "" {
		if val, err := strconv.ParseBool(onSale); err == nil {
			req.OnSale = &val
		}
	}
	if featured := c.Query("featured"); featured != "" {
		if val, err := strconv.ParseBool(featured); err == nil {
			req.Featured = &val
		}
	}
	if hasImages := c.Query("has_images"); hasImages != "" {
		if val, err := strconv.ParseBool(hasImages); err == nil {
			req.HasImages = &val
		}
	}
	if hasVariants := c.Query("has_variants"); hasVariants != "" {
		if val, err := strconv.ParseBool(hasVariants); err == nil {
			req.HasVariants = &val
		}
	}
	if hasReviews := c.Query("has_reviews"); hasReviews != "" {
		if val, err := strconv.ParseBool(hasReviews); err == nil {
			req.HasReviews = &val
		}
	}

	// Parse date filters
	if createdAfter := c.Query("created_after"); createdAfter != "" {
		req.CreatedAfter = &createdAfter
	}
	if createdBefore := c.Query("created_before"); createdBefore != "" {
		req.CreatedBefore = &createdBefore
	}
	if updatedAfter := c.Query("updated_after"); updatedAfter != "" {
		req.UpdatedAfter = &updatedAfter
	}
	if updatedBefore := c.Query("updated_before"); updatedBefore != "" {
		req.UpdatedBefore = &updatedBefore
	}

	// Parse attributes
	req.Attributes = make(map[string][]string)
	for key, values := range c.Request.URL.Query() {
		if strings.HasPrefix(key, "attr_") {
			attrID := strings.TrimPrefix(key, "attr_")
			req.Attributes[attrID] = values
		}
	}

	// Parse sorting and pagination
	req.SortBy = c.DefaultQuery("sort_by", "created_at")
	req.SortOrder = c.DefaultQuery("sort_order", "desc")

	// Parse and validate pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Validate and normalize pagination
	page, limit, err := usecases.ValidateAndNormalizePagination(page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	req.Page = page
	req.Limit = limit

	// Parse facet options
	if includeFacets := c.Query("include_facets"); includeFacets != "" {
		if val, err := strconv.ParseBool(includeFacets); err == nil {
			req.IncludeFacets = val
		}
	}
	
	if facetLimit, err := strconv.Atoi(c.DefaultQuery("facet_limit", "10")); err == nil {
		req.FacetLimit = facetLimit
	} else {
		req.FacetLimit = 10
	}

	// Execute filtering
	result, err := h.filterUseCase.FilterProducts(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to filter products: " + err.Error(),
		})
		return
	}

	// Track filter usage
	userID := getUserIDFromContext(c)
	sessionID := getSessionIDFromContext(c)
	
	// Track major filters
	totalCount := result.Pagination.Total
	if req.Query != "" {
		h.filterUseCase.TrackFilterUsage(c.Request.Context(), userID, sessionID, "query", "search", req.Query, int(totalCount))
	}
	if len(req.CategoryIDs) > 0 {
		h.filterUseCase.TrackFilterUsage(c.Request.Context(), userID, sessionID, "category", "category_ids", strings.Join(req.CategoryIDs, ","), int(totalCount))
	}
	if len(req.BrandIDs) > 0 {
		h.filterUseCase.TrackFilterUsage(c.Request.Context(), userID, sessionID, "brand", "brand_ids", strings.Join(req.BrandIDs, ","), int(totalCount))
	}

	// Create response with enhanced pagination format
	response := map[string]interface{}{
		"data":       result.Products,
		"pagination": result.Pagination,
	}

	// Include facets if available
	if result.Facets != nil {
		response["facets"] = result.Facets
	}

	c.JSON(http.StatusOK, response)
}

// GetFilterFacets gets available filter facets
func (h *ProductFilterHandler) GetFilterFacets(c *gin.Context) {
	var req usecases.AdvancedFilterRequest

	// Parse query parameters for filtering facets
	req.CategoryIDs = parseStringSlice(c.Query("category_ids"))
	req.BrandIDs = parseStringSlice(c.Query("brand_ids"))
	req.Tags = parseStringSlice(c.Query("tags"))

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request: " + err.Error(),
		})
		return
	}

	facets, err := h.filterUseCase.GetDynamicFilters(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get dynamic filters: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Dynamic filters retrieved successfully",
		Data:    facets,
	})
}

// GetDynamicFilters gets dynamic filters based on current state
func (h *ProductFilterHandler) GetDynamicFilters(c *gin.Context) {
	var req usecases.AdvancedFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request: " + err.Error(),
		})
		return
	}

	facets, err := h.filterUseCase.GetDynamicFilters(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get dynamic filters: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Dynamic filters retrieved successfully",
		Data:    facets,
	})
}

// SaveFilterSet saves a filter set
func (h *ProductFilterHandler) SaveFilterSet(c *gin.Context) {
	var req usecases.FilterSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request: " + err.Error(),
		})
		return
	}

	userID := getUserIDFromContext(c)
	sessionID := getSessionIDFromContext(c)

	result, err := h.filterUseCase.SaveFilterSet(c.Request.Context(), userID, sessionID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to save filter set: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Filter set saved successfully",
		Data:    result,
	})
}

// GetFilterSet gets a filter set by ID
func (h *ProductFilterHandler) GetFilterSet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid filter set ID",
		})
		return
	}

	result, err := h.filterUseCase.GetFilterSet(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Filter set not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Filter set retrieved successfully",
		Data:    result,
	})
}

// GetUserFilterSets gets filter sets for current user
func (h *ProductFilterHandler) GetUserFilterSets(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Authentication required",
		})
		return
	}

	result, err := h.filterUseCase.GetUserFilterSets(c.Request.Context(), *userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get user filter sets: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "User filter sets retrieved successfully",
		Data:    result,
	})
}

// GetSessionFilterSets gets filter sets for current session
func (h *ProductFilterHandler) GetSessionFilterSets(c *gin.Context) {
	sessionID := getSessionIDFromContext(c)

	result, err := h.filterUseCase.GetSessionFilterSets(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get session filter sets: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Session filter sets retrieved successfully",
		Data:    result,
	})
}

// UpdateFilterSet updates a filter set
func (h *ProductFilterHandler) UpdateFilterSet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid filter set ID",
		})
		return
	}

	var req usecases.FilterSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request: " + err.Error(),
		})
		return
	}

	result, err := h.filterUseCase.UpdateFilterSet(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to update filter set: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Filter set updated successfully",
		Data:    result,
	})
}

// DeleteFilterSet deletes a filter set
func (h *ProductFilterHandler) DeleteFilterSet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid filter set ID",
		})
		return
	}

	if err := h.filterUseCase.DeleteFilterSet(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to delete filter set: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Filter set deleted successfully",
	})
}

// GetFilterAnalytics gets filter analytics
func (h *ProductFilterHandler) GetFilterAnalytics(c *gin.Context) {
	days := 30
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	analytics, err := h.filterUseCase.GetFilterAnalytics(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get filter analytics: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Filter analytics retrieved successfully",
		Data:    analytics,
	})
}

// GetPopularFilters gets popular filters
func (h *ProductFilterHandler) GetPopularFilters(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	filters, err := h.filterUseCase.GetPopularFilters(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get popular filters: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Popular filters retrieved successfully",
		Data:    filters,
	})
}

// GetFilterSuggestions gets filter suggestions
func (h *ProductFilterHandler) GetFilterSuggestions(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Query parameter 'q' is required",
		})
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	suggestions, err := h.filterUseCase.GetFilterSuggestions(c.Request.Context(), query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get filter suggestions: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Filter suggestions retrieved successfully",
		Data:    suggestions,
	})
}

// GetRelatedFilters gets related filters
func (h *ProductFilterHandler) GetRelatedFilters(c *gin.Context) {
	var req usecases.AdvancedFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request: " + err.Error(),
		})
		return
	}

	related, err := h.filterUseCase.GetRelatedFilters(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get related filters: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Related filters retrieved successfully",
		Data:    related,
	})
}

// GetAttributeFilters gets attribute filters
func (h *ProductFilterHandler) GetAttributeFilters(c *gin.Context) {
	categoryID := c.Query("category_id")
	var categoryPtr *string
	if categoryID != "" {
		categoryPtr = &categoryID
	}

	attributes, err := h.filterUseCase.GetAttributeFilters(c.Request.Context(), categoryPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get attribute filters: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Attribute filters retrieved successfully",
		Data:    attributes,
	})
}

// GetAttributeTerms gets attribute terms
func (h *ProductFilterHandler) GetAttributeTerms(c *gin.Context) {
	attributeIDStr := c.Param("attribute_id")
	attributeID, err := uuid.Parse(attributeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid attribute ID",
		})
		return
	}

	categoryID := c.Query("category_id")
	var categoryPtr *string
	if categoryID != "" {
		categoryPtr = &categoryID
	}

	terms, err := h.filterUseCase.GetAttributeTerms(c.Request.Context(), attributeID, categoryPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get attribute terms: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Attribute terms retrieved successfully",
		Data:    terms,
	})
}

// Helper functions
func parseStringSlice(value string) []string {
	if value == "" {
		return nil
	}
	return strings.Split(value, ",")
}

func getUserIDFromContext(c *gin.Context) *uuid.UUID {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uuid.UUID); ok {
			return &id
		}
		if idStr, ok := userID.(string); ok {
			if id, err := uuid.Parse(idStr); err == nil {
				return &id
			}
		}
	}
	return nil
}

func getSessionIDFromContext(c *gin.Context) string {
	if sessionID := c.GetHeader("X-Session-ID"); sessionID != "" {
		return sessionID
	}
	if sessionID, exists := c.Get("session_id"); exists {
		if id, ok := sessionID.(string); ok {
			return id
		}
	}
	return "anonymous"
}
