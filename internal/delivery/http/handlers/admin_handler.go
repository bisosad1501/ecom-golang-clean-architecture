package handlers

import (
	"net/http"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AdminHandler handles admin-related HTTP requests
type AdminHandler struct {
	adminUseCase usecases.AdminUseCase
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(adminUseCase usecases.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase: adminUseCase,
	}
}

// GetDashboard returns admin dashboard data
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	var req usecases.AdminDashboardRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	dashboard, err := h.adminUseCase.GetDashboard(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get dashboard",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Dashboard retrieved successfully",
		Data: dashboard,
	})
}

// GetSystemStats returns system statistics
func (h *AdminHandler) GetSystemStats(c *gin.Context) {
	stats, err := h.adminUseCase.GetSystemStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get system stats",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "System stats retrieved successfully",
		Data: stats,
	})
}

// GetUsers returns paginated list of users
func (h *AdminHandler) GetUsers(c *gin.Context) {
	var req usecases.AdminUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	users, err := h.adminUseCase.GetUsers(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get users",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Users retrieved successfully",
		Data: users,
	})
}

// UpdateUserStatus updates a user's status
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Status entities.UserStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.adminUseCase.UpdateUserStatus(c.Request.Context(), userID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to update user status",
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
			Error: "Invalid user ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Role entities.UserRole `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.adminUseCase.UpdateUserRole(c.Request.Context(), userID, req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to update user role",
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
			Error: "Invalid user ID",
			Details: err.Error(),
		})
		return
	}

	var req usecases.ActivityRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	activity, err := h.adminUseCase.GetUserActivity(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get user activity",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "User activity retrieved successfully",
		Data: activity,
	})
}

// GetOrders returns paginated list of orders
func (h *AdminHandler) GetOrders(c *gin.Context) {
	var req usecases.AdminOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	orders, err := h.adminUseCase.GetOrders(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get orders",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Orders retrieved successfully",
		Data: orders,
	})
}

// UpdateOrderStatus updates an order's status
func (h *AdminHandler) UpdateOrderStatus(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid order ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Status entities.OrderStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.adminUseCase.UpdateOrderStatus(c.Request.Context(), orderID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to update order status",
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
	orderIDStr := c.Param("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid order ID",
			Details: err.Error(),
		})
		return
	}

	details, err := h.adminUseCase.GetOrderDetails(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Order not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Order details retrieved successfully",
		Data: details,
	})
}

// ProcessRefund processes a refund for an order
func (h *AdminHandler) ProcessRefund(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid order ID",
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
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.adminUseCase.ProcessRefund(c.Request.Context(), orderID, req.Amount, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to process refund",
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
			Error: "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	products, err := h.adminUseCase.GetProducts(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get products",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Products retrieved successfully",
		Data: products,
	})
}

// BulkUpdateProducts updates multiple products
func (h *AdminHandler) BulkUpdateProducts(c *gin.Context) {
	var req usecases.BulkUpdateProductsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.adminUseCase.BulkUpdateProducts(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to bulk update products",
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
			Error: "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	logs, err := h.adminUseCase.GetAuditLogs(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get audit logs",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Audit logs retrieved successfully",
		Data: logs,
	})
}

// ManageReviews returns paginated list of reviews for admin management
func (h *AdminHandler) ManageReviews(c *gin.Context) {
	var req usecases.ManageReviewsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	reviews, err := h.adminUseCase.ManageReviews(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get reviews",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Reviews retrieved successfully",
		Data: reviews,
	})
}

// UpdateReviewStatus updates review status (approve/reject/hide)
func (h *AdminHandler) UpdateReviewStatus(c *gin.Context) {
	reviewIDStr := c.Param("id")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid review ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Status entities.ReviewStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	err = h.adminUseCase.UpdateReviewStatus(c.Request.Context(), reviewID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to update review status",
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
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	report, err := h.adminUseCase.GenerateReport(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to generate report",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Report generated successfully",
		Data: report,
	})
}

// GetReports returns paginated list of reports
func (h *AdminHandler) GetReports(c *gin.Context) {
	var req usecases.GetReportsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	reports, err := h.adminUseCase.GetReports(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get reports",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Reports retrieved successfully",
		Data: reports,
	})
}

// DownloadReport downloads a report
func (h *AdminHandler) DownloadReport(c *gin.Context) {
	reportIDStr := c.Param("id")
	reportID, err := uuid.Parse(reportIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid report ID",
			Details: err.Error(),
		})
		return
	}

	download, err := h.adminUseCase.DownloadReport(c.Request.Context(), reportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to download report",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Report download prepared successfully",
		Data: download,
	})
}

// GetSystemLogs returns system logs
func (h *AdminHandler) GetSystemLogs(c *gin.Context) {
	var req usecases.SystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	logs, err := h.adminUseCase.GetSystemLogs(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get system logs",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "System logs retrieved successfully",
		Data: logs,
	})
}

// BackupDatabase creates a database backup
func (h *AdminHandler) BackupDatabase(c *gin.Context) {
	backup, err := h.adminUseCase.BackupDatabase(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to backup database",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Database backup created successfully",
		Data: backup,
	})
}
