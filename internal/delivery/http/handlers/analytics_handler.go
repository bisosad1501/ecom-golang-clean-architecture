package handlers

import (
	"net/http"
	"strconv"
	"time"

	"ecom-golang-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AnalyticsHandler handles analytics-related HTTP requests
type AnalyticsHandler struct {
	analyticsUseCase usecases.AnalyticsUseCase
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(analyticsUseCase usecases.AnalyticsUseCase) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsUseCase: analyticsUseCase,
	}
}

// TrackEvent tracks a custom event
func (h *AnalyticsHandler) TrackEvent(c *gin.Context) {
	var req usecases.TrackEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.analyticsUseCase.TrackEvent(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to track event",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Event tracked successfully",
	})
}

// TrackPageView tracks a page view
func (h *AnalyticsHandler) TrackPageView(c *gin.Context) {
	var req usecases.TrackPageViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.analyticsUseCase.TrackPageView(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to track page view",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Page view tracked successfully",
	})
}

// TrackProductView tracks a product view
func (h *AnalyticsHandler) TrackProductView(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid product ID",
			Details: err.Error(),
		})
		return
	}

	var userID *uuid.UUID
	if userIDStr := c.GetHeader("X-User-ID"); userIDStr != "" {
		if id, err := uuid.Parse(userIDStr); err == nil {
			userID = &id
		}
	}

	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		sessionID = "anonymous"
	}

	if err := h.analyticsUseCase.TrackProductView(c.Request.Context(), productID, userID, sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to track product view",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product view tracked successfully",
	})
}

// GetDashboardMetrics returns dashboard metrics
func (h *AnalyticsHandler) GetDashboardMetrics(c *gin.Context) {
	var req usecases.DashboardMetricsRequest

	// Parse query parameters
	if period := c.Query("period"); period != "" {
		req.Period = period
	}
	if startDate := c.Query("date_from"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			req.DateFrom = &t
		}
	}
	if endDate := c.Query("date_to"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			req.DateTo = &t
		}
	}

	metrics, err := h.analyticsUseCase.GetDashboardMetrics(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get dashboard metrics",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Dashboard metrics retrieved successfully",
		Data: metrics,
	})
}

// GetSalesMetrics returns sales metrics
func (h *AnalyticsHandler) GetSalesMetrics(c *gin.Context) {
	var req usecases.SalesMetricsRequest

	// Parse query parameters
	if period := c.Query("period"); period != "" {
		req.Period = period
	}
	if startDate := c.Query("date_from"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			req.DateFrom = &t
		}
	}
	if endDate := c.Query("date_to"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			req.DateTo = &t
		}
	}

	metrics, err := h.analyticsUseCase.GetSalesMetrics(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get sales metrics",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Sales metrics retrieved successfully",
		Data: metrics,
	})
}

// GetProductMetrics returns product metrics
func (h *AnalyticsHandler) GetProductMetrics(c *gin.Context) {
	var req usecases.ProductMetricsRequest

	// Parse query parameters
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}

	metrics, err := h.analyticsUseCase.GetProductMetrics(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get product metrics",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product metrics retrieved successfully",
		Data: metrics,
	})
}

// GetUserMetrics returns user metrics
func (h *AnalyticsHandler) GetUserMetrics(c *gin.Context) {
	var req usecases.UserMetricsRequest

	// Parse query parameters
	if period := c.Query("period"); period != "" {
		req.Period = period
	}

	metrics, err := h.analyticsUseCase.GetUserMetrics(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get user metrics",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "User metrics retrieved successfully",
		Data: metrics,
	})
}

// GetTrafficMetrics returns traffic metrics
func (h *AnalyticsHandler) GetTrafficMetrics(c *gin.Context) {
	var req usecases.TrafficMetricsRequest

	// Parse query parameters
	if period := c.Query("period"); period != "" {
		req.Period = period
	}

	metrics, err := h.analyticsUseCase.GetTrafficMetrics(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get traffic metrics",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Traffic metrics retrieved successfully",
		Data: metrics,
	})
}

// GetRealTimeMetrics returns real-time metrics
func (h *AnalyticsHandler) GetRealTimeMetrics(c *gin.Context) {
	metrics, err := h.analyticsUseCase.GetRealTimeMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get real-time metrics",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Real-time metrics retrieved successfully",
		Data: metrics,
	})
}

// GetTopProducts returns top products
func (h *AnalyticsHandler) GetTopProducts(c *gin.Context) {
	period := c.DefaultQuery("period", "30d")

	// Parse and validate pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate and normalize pagination for analytics
	page, limit, err := usecases.ValidateAndNormalizePagination(page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Get top products with pagination
	response, err := h.analyticsUseCase.GetTopProductsPaginated(c.Request.Context(), period, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get top products",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       response.Products,
		Pagination: response.Pagination,
	})
}

// GetTopCategories returns top categories
func (h *AnalyticsHandler) GetTopCategories(c *gin.Context) {
	period := c.DefaultQuery("period", "30d")

	// Parse and validate pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate and normalize pagination for analytics
	page, limit, err := usecases.ValidateAndNormalizePagination(page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Get top categories with pagination
	response, err := h.analyticsUseCase.GetTopCategoriesPaginated(c.Request.Context(), period, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get top categories",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       response.Categories,
		Pagination: response.Pagination,
	})
}
