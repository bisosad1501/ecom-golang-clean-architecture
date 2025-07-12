package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/usecases"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AdminHandler handles admin-related HTTP requests
type AdminHandler struct {
	adminUseCase        usecases.AdminUseCase
	stockCleanupUseCase usecases.StockCleanupUseCase
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(adminUseCase usecases.AdminUseCase, stockCleanupUseCase usecases.StockCleanupUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase:        adminUseCase,
		stockCleanupUseCase: stockCleanupUseCase,
	}
}

// GetDashboard returns admin dashboard data
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	var req usecases.AdminDashboardRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	dashboard, err := h.adminUseCase.GetDashboard(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get dashboard",
			Details: err.Error(),
		})
		return
	}

	// Debug log
	fmt.Printf("Dashboard response - Total Revenue: %f\n", dashboard.Overview.TotalRevenue)

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Dashboard retrieved successfully",
		Data:    dashboard,
	})
}

// GetSystemStats returns system statistics
func (h *AdminHandler) GetSystemStats(c *gin.Context) {
	stats, err := h.adminUseCase.GetSystemStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get system stats",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "System stats retrieved successfully",
		Data:    stats,
	})
}

// GetUsers returns paginated list of users
func (h *AdminHandler) GetUsers(c *gin.Context) {
	var req usecases.AdminUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	// Set default values if not provided
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Debug logging
	fmt.Printf("DEBUG GetUsers - Limit: %d, Offset: %d, Status: %v, Role: %v, Search: %s\n",
		req.Limit, req.Offset, req.Status, req.Role, req.Search)

	users, err := h.adminUseCase.GetUsers(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get users",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

// UpdateUserStatus updates a user's status
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Status entities.UserStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.adminUseCase.UpdateUserStatus(c.Request.Context(), userID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update user status",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "User status updated successfully",
	})
}

// UpdateUserRole updates a user's role
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Role entities.UserRole `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.adminUseCase.UpdateUserRole(c.Request.Context(), userID, req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update user role",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "User role updated successfully",
	})
}

// GetUserActivity returns user activity
func (h *AdminHandler) GetUserActivity(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Details: err.Error(),
		})
		return
	}

	var req usecases.ActivityRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	activity, err := h.adminUseCase.GetUserActivity(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get user activity",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "User activity retrieved successfully",
		Data:    activity,
	})
}

// GetOrders returns paginated list of orders
func (h *AdminHandler) GetOrders(c *gin.Context) {
	var req usecases.AdminOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	orders, err := h.adminUseCase.GetOrders(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get orders",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Orders retrieved successfully",
		Data:    orders,
	})
}

// UpdateOrderStatus updates an order's status
func (h *AdminHandler) UpdateOrderStatus(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid order ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Status entities.OrderStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.adminUseCase.UpdateOrderStatus(c.Request.Context(), orderID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update order status",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Order status updated successfully",
	})
}

// GetOrderDetails returns detailed order information
func (h *AdminHandler) GetOrderDetails(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid order ID",
			Details: err.Error(),
		})
		return
	}

	details, err := h.adminUseCase.GetOrderDetails(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Order not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Order details retrieved successfully",
		Data:    details,
	})
}

// ProcessRefund processes a refund for an order
func (h *AdminHandler) ProcessRefund(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid order ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
		Reason string  `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.adminUseCase.ProcessRefund(c.Request.Context(), orderID, req.Amount, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to process refund",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Refund processed successfully",
	})
}

// GetProducts returns paginated list of products for admin
func (h *AdminHandler) GetProducts(c *gin.Context) {
	var req usecases.AdminProductsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	products, err := h.adminUseCase.GetProducts(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get products",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Products retrieved successfully",
		Data:    products,
	})
}

// BulkUpdateProducts updates multiple products
func (h *AdminHandler) BulkUpdateProducts(c *gin.Context) {
	var req usecases.BulkUpdateProductsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.adminUseCase.BulkUpdateProducts(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to bulk update products",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Products updated successfully",
	})
}

// GetAuditLogs returns audit logs
func (h *AdminHandler) GetAuditLogs(c *gin.Context) {
	var req usecases.AuditLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	logs, err := h.adminUseCase.GetAuditLogs(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get audit logs",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Audit logs retrieved successfully",
		Data:    logs,
	})
}

// ManageReviews returns paginated list of reviews for admin management
func (h *AdminHandler) ManageReviews(c *gin.Context) {
	var req usecases.ManageReviewsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	reviews, err := h.adminUseCase.ManageReviews(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get reviews",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Reviews retrieved successfully",
		Data:    reviews,
	})
}

// UpdateReviewStatus updates review status (approve/reject/hide)
func (h *AdminHandler) UpdateReviewStatus(c *gin.Context) {
	reviewIDStr := c.Param("id")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid review ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Status entities.ReviewStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	err = h.adminUseCase.UpdateReviewStatus(c.Request.Context(), reviewID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update review status",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Review status updated successfully",
	})
}

// GenerateReport generates a new report
func (h *AdminHandler) GenerateReport(c *gin.Context) {
	var req usecases.GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	report, err := h.adminUseCase.GenerateReport(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to generate report",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Report generated successfully",
		Data:    report,
	})
}

// GetReports returns paginated list of reports
func (h *AdminHandler) GetReports(c *gin.Context) {
	var req usecases.GetReportsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	reports, err := h.adminUseCase.GetReports(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get reports",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Reports retrieved successfully",
		Data:    reports,
	})
}

// DownloadReport downloads a report
func (h *AdminHandler) DownloadReport(c *gin.Context) {
	reportIDStr := c.Param("id")
	reportID, err := uuid.Parse(reportIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid report ID",
			Details: err.Error(),
		})
		return
	}

	download, err := h.adminUseCase.DownloadReport(c.Request.Context(), reportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to download report",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Report download prepared successfully",
		Data:    download,
	})
}

// GetSystemLogs returns system logs
func (h *AdminHandler) GetSystemLogs(c *gin.Context) {
	var req usecases.SystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	logs, err := h.adminUseCase.GetSystemLogs(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get system logs",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "System logs retrieved successfully",
		Data:    logs,
	})
}

// BackupDatabase creates a database backup
func (h *AdminHandler) BackupDatabase(c *gin.Context) {
	backup, err := h.adminUseCase.BackupDatabase(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to backup database",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Database backup created successfully",
		Data:    backup,
	})
}

// GetRecentActivity returns recent admin activity
func (h *AdminHandler) GetRecentActivity(c *gin.Context) {
	// Parse query parameters
	limit := 5
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Create mock recent activity data
	// In a real implementation, this would come from audit logs or activity tracking
	activities := []map[string]interface{}{
		{
			"id":          uuid.New().String(),
			"type":        "order_created",
			"description": "New order #ORD-001 placed by John Doe",
			"timestamp":   time.Now().Add(-5 * time.Minute),
			"user_id":     uuid.New().String(),
			"user_name":   "John Doe",
		},
		{
			"id":          uuid.New().String(),
			"type":        "product_updated",
			"description": "Product 'iPhone 15' stock updated",
			"timestamp":   time.Now().Add(-15 * time.Minute),
			"user_id":     uuid.New().String(),
			"user_name":   "Admin User",
		},
		{
			"id":          uuid.New().String(),
			"type":        "user_registered",
			"description": "New user Jane Smith registered",
			"timestamp":   time.Now().Add(-30 * time.Minute),
			"user_id":     uuid.New().String(),
			"user_name":   "Jane Smith",
		},
		{
			"id":          uuid.New().String(),
			"type":        "payment_processed",
			"description": "Payment processed for order #ORD-002",
			"timestamp":   time.Now().Add(-45 * time.Minute),
			"user_id":     uuid.New().String(),
			"user_name":   "Mike Johnson",
		},
		{
			"id":          uuid.New().String(),
			"type":        "review_submitted",
			"description": "New review submitted for MacBook Pro",
			"timestamp":   time.Now().Add(-1 * time.Hour),
			"user_id":     uuid.New().String(),
			"user_name":   "Sarah Wilson",
		},
	}

	// Limit results
	if len(activities) > limit {
		activities = activities[:limit]
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Recent activity retrieved successfully",
		Data:    activities,
	})
}

// ReplyToReview allows admin to reply to a review
func (h *AdminHandler) ReplyToReview(c *gin.Context) {
	reviewIDStr := c.Param("id")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid review ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Reply string `json:"reply" validate:"required,max=1000"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.adminUseCase.AdminReplyToReview(c.Request.Context(), reviewID, req.Reply); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to reply to review",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Reply added successfully",
	})
}

// GetCleanupStats returns cleanup statistics
func (h *AdminHandler) GetCleanupStats(c *gin.Context) {
	stats, err := h.stockCleanupUseCase.GetCleanupStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get cleanup statistics",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Cleanup statistics retrieved successfully",
		Data:    stats,
	})
}

// TriggerCleanup manually triggers cleanup process
func (h *AdminHandler) TriggerCleanup(c *gin.Context) {
	err := h.stockCleanupUseCase.RunCleanup(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to run cleanup process",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Cleanup process completed successfully",
	})
}
