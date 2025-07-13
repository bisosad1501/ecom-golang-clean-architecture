package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SearchHandler handles search-related HTTP requests
type SearchHandler struct {
	searchUseCase usecases.SearchUseCase
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchUseCase usecases.SearchUseCase) *SearchHandler {
	return &SearchHandler{
		searchUseCase: searchUseCase,
	}
}

// FullTextSearch handles full-text search requests
// @Summary Perform full-text search
// @Description Search products using advanced full-text search with filters and facets
// @Tags search
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param category_ids query string false "Category IDs (comma-separated)"
// @Param brand_ids query string false "Brand IDs (comma-separated)"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param in_stock query boolean false "In stock only"
// @Param featured query boolean false "Featured products only"
// @Param on_sale query boolean false "On sale products only"
// @Param tags query string false "Tags (comma-separated)"
// @Param min_rating query number false "Minimum rating"
// @Param max_rating query number false "Maximum rating"
// @Param visibility query string false "Product visibility (public, private, hidden)"
// @Param product_type query string false "Product type (simple, variable, grouped)"
// @Param status query string false "Product status (active, inactive, draft)"
// @Param availability_status query string false "Availability status (in_stock, out_of_stock, low_stock)"
// @Param created_after query string false "Created after date (RFC3339)"
// @Param created_before query string false "Created before date (RFC3339)"
// @Param updated_after query string false "Updated after date (RFC3339)"
// @Param updated_before query string false "Updated before date (RFC3339)"
// @Param min_weight query number false "Minimum weight"
// @Param max_weight query number false "Maximum weight"
// @Param shipping_class query string false "Shipping class"
// @Param tax_class query string false "Tax class"
// @Param min_discount_percent query number false "Minimum discount percentage"
// @Param max_discount_percent query number false "Maximum discount percentage"
// @Param is_digital query boolean false "Digital products only"
// @Param requires_shipping query boolean false "Requires shipping"
// @Param allow_backorder query boolean false "Allow backorder"
// @Param track_quantity query boolean false "Track quantity"
// @Param sort_by query string false "Sort by (relevance, price, name, created_at, rating)" default(relevance)
// @Param sort_order query string false "Sort order (asc, desc)" default(desc)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} usecases.SearchResponse
// @Router /search [get]
func (h *SearchHandler) FullTextSearch(c *gin.Context) {
	req := usecases.FullTextSearchRequest{
		Query:     c.Query("q"),
		SortBy:    c.DefaultQuery("sort_by", "relevance"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
		Page:      1,
		Limit:     20,
		SessionID: c.GetString("session_id"),
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}

	// Parse user ID if authenticated
	if userIDStr, exists := c.Get("user_id"); exists {
		if userID, err := uuid.Parse(userIDStr.(string)); err == nil {
			req.UserID = &userID
		}
	}

	// Parse category IDs (support both category_id and category_ids)
	categoryIDsStr := c.Query("category_ids")
	if categoryIDsStr == "" {
		categoryIDsStr = c.Query("category_id") // Support singular form
	}
	if categoryIDsStr != "" {
		categoryIDStrs := strings.Split(categoryIDsStr, ",")
		for _, idStr := range categoryIDStrs {
			if id, err := uuid.Parse(strings.TrimSpace(idStr)); err == nil {
				req.CategoryIDs = append(req.CategoryIDs, id)
			}
		}
	}

	// Parse brand IDs (support both brand_id and brand_ids)
	brandIDsStr := c.Query("brand_ids")
	if brandIDsStr == "" {
		brandIDsStr = c.Query("brand_id") // Support singular form
	}
	if brandIDsStr != "" {
		brandIDStrs := strings.Split(brandIDsStr, ",")
		for _, idStr := range brandIDStrs {
			if id, err := uuid.Parse(strings.TrimSpace(idStr)); err == nil {
				req.BrandIDs = append(req.BrandIDs, id)
			}
		}
	}

	// Parse price filters
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

	// Parse boolean filters
	if inStockStr := c.Query("in_stock"); inStockStr != "" {
		if inStock, err := strconv.ParseBool(inStockStr); err == nil {
			req.InStock = &inStock
		}
	}

	if featuredStr := c.Query("featured"); featuredStr != "" {
		if featured, err := strconv.ParseBool(featuredStr); err == nil {
			req.Featured = &featured
		}
	}

	if onSaleStr := c.Query("on_sale"); onSaleStr != "" {
		if onSale, err := strconv.ParseBool(onSaleStr); err == nil {
			req.OnSale = &onSale
		}
	}

	// Parse tags
	if tagsStr := c.Query("tags"); tagsStr != "" {
		req.Tags = strings.Split(tagsStr, ",")
		for i, tag := range req.Tags {
			req.Tags[i] = strings.TrimSpace(tag)
		}
	}

	// Parse advanced filters
	// Rating filters
	if minRatingStr := c.Query("min_rating"); minRatingStr != "" {
		if minRating, err := strconv.ParseFloat(minRatingStr, 64); err == nil {
			req.MinRating = &minRating
		}
	}
	if maxRatingStr := c.Query("max_rating"); maxRatingStr != "" {
		if maxRating, err := strconv.ParseFloat(maxRatingStr, 64); err == nil {
			req.MaxRating = &maxRating
		}
	}

	// Enum filters
	if visibilityStr := c.Query("visibility"); visibilityStr != "" {
		visibility := entities.ProductVisibility(visibilityStr)
		req.Visibility = &visibility
	}
	if productTypeStr := c.Query("product_type"); productTypeStr != "" {
		productType := entities.ProductType(productTypeStr)
		req.ProductType = &productType
	}
	if statusStr := c.Query("status"); statusStr != "" {
		status := entities.ProductStatus(statusStr)
		req.Status = &status
	}
	if availabilityStr := c.Query("availability_status"); availabilityStr != "" {
		req.AvailabilityStatus = &availabilityStr
	}

	// Date filters
	if createdAfterStr := c.Query("created_after"); createdAfterStr != "" {
		if createdAfter, err := time.Parse(time.RFC3339, createdAfterStr); err == nil {
			req.CreatedAfter = &createdAfter
		}
	}
	if createdBeforeStr := c.Query("created_before"); createdBeforeStr != "" {
		if createdBefore, err := time.Parse(time.RFC3339, createdBeforeStr); err == nil {
			req.CreatedBefore = &createdBefore
		}
	}
	if updatedAfterStr := c.Query("updated_after"); updatedAfterStr != "" {
		if updatedAfter, err := time.Parse(time.RFC3339, updatedAfterStr); err == nil {
			req.UpdatedAfter = &updatedAfter
		}
	}
	if updatedBeforeStr := c.Query("updated_before"); updatedBeforeStr != "" {
		if updatedBefore, err := time.Parse(time.RFC3339, updatedBeforeStr); err == nil {
			req.UpdatedBefore = &updatedBefore
		}
	}

	// Weight filters
	if minWeightStr := c.Query("min_weight"); minWeightStr != "" {
		if minWeight, err := strconv.ParseFloat(minWeightStr, 64); err == nil {
			req.MinWeight = &minWeight
		}
	}
	if maxWeightStr := c.Query("max_weight"); maxWeightStr != "" {
		if maxWeight, err := strconv.ParseFloat(maxWeightStr, 64); err == nil {
			req.MaxWeight = &maxWeight
		}
	}

	// String filters
	if shippingClassStr := c.Query("shipping_class"); shippingClassStr != "" {
		req.ShippingClass = &shippingClassStr
	}
	if taxClassStr := c.Query("tax_class"); taxClassStr != "" {
		req.TaxClass = &taxClassStr
	}

	// Discount filters
	if minDiscountStr := c.Query("min_discount_percent"); minDiscountStr != "" {
		if minDiscount, err := strconv.ParseFloat(minDiscountStr, 64); err == nil {
			req.MinDiscountPercent = &minDiscount
		}
	}
	if maxDiscountStr := c.Query("max_discount_percent"); maxDiscountStr != "" {
		if maxDiscount, err := strconv.ParseFloat(maxDiscountStr, 64); err == nil {
			req.MaxDiscountPercent = &maxDiscount
		}
	}

	// Boolean filters
	if isDigitalStr := c.Query("is_digital"); isDigitalStr != "" {
		if isDigital, err := strconv.ParseBool(isDigitalStr); err == nil {
			req.IsDigital = &isDigital
		}
	}
	if requiresShippingStr := c.Query("requires_shipping"); requiresShippingStr != "" {
		if requiresShipping, err := strconv.ParseBool(requiresShippingStr); err == nil {
			req.RequiresShipping = &requiresShipping
		}
	}
	if allowBackorderStr := c.Query("allow_backorder"); allowBackorderStr != "" {
		if allowBackorder, err := strconv.ParseBool(allowBackorderStr); err == nil {
			req.AllowBackorder = &allowBackorder
		}
	}
	if trackQuantityStr := c.Query("track_quantity"); trackQuantityStr != "" {
		if trackQuantity, err := strconv.ParseBool(trackQuantityStr); err == nil {
			req.TrackQuantity = &trackQuantity
		}
	}

	// Parse pagination
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			req.Limit = limit
		}
	}

	// Perform search
	response, err := h.searchUseCase.FullTextSearch(c.Request.Context(), req)
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

// GetSearchSuggestions handles search suggestions requests
// @Summary Get search suggestions
// @Description Get search suggestions based on query
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Number of suggestions" default(10)
// @Success 200 {array} string
// @Router /search/suggestions [get]
func (h *SearchHandler) GetSearchSuggestions(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Query parameter 'q' is required",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	suggestions, err := h.searchUseCase.GetSearchSuggestions(c.Request.Context(), query, limit)
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

// GetSearchFacets handles search facets requests
// @Summary Get search facets
// @Description Get search facets for filtering
// @Tags search
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Success 200 {object} usecases.SearchFacetsResponse
// @Router /search/facets [get]
func (h *SearchHandler) GetSearchFacets(c *gin.Context) {
	query := c.Query("q")

	facets, err := h.searchUseCase.GetSearchFacets(c.Request.Context(), query)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: facets,
	})
}

// RecordSearchEvent handles recording search events
// @Summary Record search event
// @Description Record a search event for analytics
// @Tags search
// @Accept json
// @Produce json
// @Param request body usecases.RecordSearchEventRequest true "Record search event request"
// @Success 200 {object} SuccessResponse
// @Router /search/record [post]
func (h *SearchHandler) RecordSearchEvent(c *gin.Context) {
	var req usecases.RecordSearchEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Set additional context
	req.SessionID = c.GetString("session_id")
	req.IPAddress = c.ClientIP()
	req.UserAgent = c.GetHeader("User-Agent")

	// Parse user ID if authenticated
	if userIDStr, exists := c.Get("user_id"); exists {
		if userID, err := uuid.Parse(userIDStr.(string)); err == nil {
			req.UserID = &userID
		}
	}

	err := h.searchUseCase.RecordSearchEvent(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Search recorded successfully",
	})
}

// GetSearchAnalytics handles getting search analytics for admin
// @Summary Get search analytics
// @Description Get search analytics for admin dashboard
// @Tags search
// @Accept json
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param limit query int false "Limit" default(50)
// @Success 200 {object} usecases.SearchAnalyticsResponse
// @Failure 400 {object} ErrorResponse
// @Router /search/analytics [get]
func (h *SearchHandler) GetSearchAnalytics(c *gin.Context) {
	req := usecases.SearchAnalyticsRequest{
		StartDate: time.Now().AddDate(0, 0, -30), // Default to last 30 days
		EndDate:   time.Now(),
		Limit:     50,
	}

	// Parse start date
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			req.StartDate = startDate
		}
	}

	// Parse end date
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			req.EndDate = endDate
		}
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			req.Limit = limit
		}
	}

	analytics, err := h.searchUseCase.GetSearchAnalytics(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetPopularSearchTerms handles popular search terms requests
// @Summary Get popular search terms
// @Description Get popular search terms for a given period
// @Tags search
// @Accept json
// @Produce json
// @Param period query string false "Period (daily, weekly, monthly)" default(daily)
// @Param limit query int false "Number of terms" default(10)
// @Success 200 {array} usecases.PopularSearchResponse
// @Router /search/popular [get]
func (h *SearchHandler) GetPopularSearchTerms(c *gin.Context) {
	period := c.DefaultQuery("period", "daily")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	terms, err := h.searchUseCase.GetPopularSearchTerms(c.Request.Context(), limit, period)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: terms,
	})
}

// GetAutocomplete handles autocomplete requests
// @Summary Get autocomplete suggestions
// @Description Get autocomplete suggestions including products, categories, brands
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Number of suggestions" default(10)
// @Success 200 {object} usecases.AutocompleteResponse
// @Router /search/autocomplete [get]
func (h *SearchHandler) GetAutocomplete(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Query parameter 'q' is required",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	autocomplete, err := h.searchUseCase.GetAutocomplete(c.Request.Context(), query, limit)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: autocomplete,
	})
}

// EnhancedSearch performs enhanced search with dynamic faceting
// @Summary Enhanced search with dynamic faceting
// @Description Perform enhanced product search with multi-select filters and dynamic facets
// @Tags search
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param category_ids query string false "Category IDs (comma-separated)"
// @Param brand_ids query string false "Brand IDs (comma-separated)"
// @Param tag_ids query string false "Tag IDs (comma-separated)"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param featured query boolean false "Featured products only"
// @Param in_stock query boolean false "In stock products only"
// @Param on_sale query boolean false "On sale products only"
// @Param sort_by query string false "Sort by (relevance, price, name, created_at, rating)" default(relevance)
// @Param sort_order query string false "Sort order (asc, desc)" default(desc)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param include_facets query boolean false "Include facets in response" default(true)
// @Param dynamic_facets query boolean false "Use dynamic facets" default(true)
// @Success 200 {object} usecases.EnhancedSearchResponse
// @Router /search/enhanced [get]
func (h *SearchHandler) EnhancedSearch(c *gin.Context) {
	req := usecases.EnhancedSearchRequest{
		Query:         c.Query("q"),
		SortBy:        c.DefaultQuery("sort_by", "relevance"),
		SortOrder:     c.DefaultQuery("sort_order", "desc"),
		Page:          1,
		Limit:         20,
		IncludeFacets: true,
		DynamicFacets: true,
		SessionID:     c.GetString("session_id"),
		IPAddress:     c.ClientIP(),
		UserAgent:     c.GetHeader("User-Agent"),
	}

	// Parse user ID if authenticated
	if userID := c.GetString("user_id"); userID != "" {
		req.UserID = &userID
	}

	// Parse multi-select filters
	if categoryIDs := c.Query("category_ids"); categoryIDs != "" {
		req.CategoryIDs = strings.Split(categoryIDs, ",")
	}

	if brandIDs := c.Query("brand_ids"); brandIDs != "" {
		req.BrandIDs = strings.Split(brandIDs, ",")
	}

	if tagIDs := c.Query("tag_ids"); tagIDs != "" {
		req.TagIDs = strings.Split(tagIDs, ",")
	}

	// Parse price filters
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

	// Parse boolean filters
	if featuredStr := c.Query("featured"); featuredStr != "" {
		if featured, err := strconv.ParseBool(featuredStr); err == nil {
			req.Featured = &featured
		}
	}

	if inStockStr := c.Query("in_stock"); inStockStr != "" {
		if inStock, err := strconv.ParseBool(inStockStr); err == nil {
			req.InStock = &inStock
		}
	}

	if onSaleStr := c.Query("on_sale"); onSaleStr != "" {
		if onSale, err := strconv.ParseBool(onSaleStr); err == nil {
			req.OnSale = &onSale
		}
	}

	// Parse facet options
	if includeFacetsStr := c.Query("include_facets"); includeFacetsStr != "" {
		if includeFacets, err := strconv.ParseBool(includeFacetsStr); err == nil {
			req.IncludeFacets = includeFacets
		}
	}

	if dynamicFacetsStr := c.Query("dynamic_facets"); dynamicFacetsStr != "" {
		if dynamicFacets, err := strconv.ParseBool(dynamicFacetsStr); err == nil {
			req.DynamicFacets = dynamicFacets
		}
	}

	// Parse pagination
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			req.Limit = limit
		}
	}

	// Perform enhanced search
	response, err := h.searchUseCase.EnhancedSearch(c.Request.Context(), &req)
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

// SaveSearchHistory handles saving search history
// @Summary Save search history
// @Description Save user search history
// @Tags search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.SaveSearchHistoryRequest true "Save search history request"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /search/history [post]
func (h *SearchHandler) SaveSearchHistory(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
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

	var req usecases.SaveSearchHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	err = h.searchUseCase.SaveSearchHistory(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Search history saved successfully",
	})
}

// GetUserSearchHistory handles getting user search history
// @Summary Get user search history
// @Description Get user's search history
// @Tags search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of history items" default(20)
// @Success 200 {array} usecases.SearchHistoryResponse
// @Failure 401 {object} ErrorResponse
// @Router /search/history [get]
func (h *SearchHandler) GetUserSearchHistory(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
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

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	history, err := h.searchUseCase.GetUserSearchHistory(c.Request.Context(), userID, limit)
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

// ClearUserSearchHistory handles clearing user search history
// @Summary Clear user search history
// @Description Clear all user's search history
// @Tags search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /search/history [delete]
func (h *SearchHandler) ClearUserSearchHistory(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
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

	err = h.searchUseCase.ClearUserSearchHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Search history cleared successfully",
	})
}

// SaveSearchFilter handles saving search filters
// @Summary Save search filter
// @Description Save a search filter for reuse
// @Tags search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.SaveSearchFilterRequest true "Save search filter request"
// @Success 201 {object} usecases.SearchFilterResponse
// @Failure 401 {object} ErrorResponse
// @Router /search/filters [post]
func (h *SearchHandler) SaveSearchFilter(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
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

	var req usecases.SaveSearchFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	filter, err := h.searchUseCase.SaveSearchFilter(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Data: filter,
	})
}

// GetUserSearchFilters handles getting user search filters
// @Summary Get user search filters
// @Description Get user's saved search filters
// @Tags search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} usecases.SearchFilterResponse
// @Failure 401 {object} ErrorResponse
// @Router /search/filters [get]
func (h *SearchHandler) GetUserSearchFilters(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
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

	filters, err := h.searchUseCase.GetUserSearchFilters(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: filters,
	})
}

// UpdateSearchFilter handles updating search filters
// @Summary Update search filter
// @Description Update a saved search filter
// @Tags search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Filter ID"
// @Param request body usecases.UpdateSearchFilterRequest true "Update search filter request"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /search/filters/{id} [put]
func (h *SearchHandler) UpdateSearchFilter(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
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

	filterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid filter ID",
		})
		return
	}

	var req usecases.UpdateSearchFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	err = h.searchUseCase.UpdateSearchFilter(c.Request.Context(), userID, filterID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Search filter updated successfully",
	})
}

// DeleteSearchFilter handles deleting search filters
// @Summary Delete search filter
// @Description Delete a saved search filter
// @Tags search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Filter ID"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /search/filters/{id} [delete]
func (h *SearchHandler) DeleteSearchFilter(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
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

	filterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid filter ID",
		})
		return
	}

	err = h.searchUseCase.DeleteSearchFilter(c.Request.Context(), userID, filterID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Search filter deleted successfully",
	})
}
