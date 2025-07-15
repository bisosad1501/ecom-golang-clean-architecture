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
	// Parse and validate pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))

	// Validate and normalize pagination for admin users
	page, limit, err := usecases.ValidateAndNormalizePaginationForEntity(page, limit, "admin_users")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Convert to offset for repository
	offset := (page - 1) * limit

	var req usecases.AdminUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	// Override pagination parameters
	req.Limit = limit
	req.Offset = offset

	// Debug logging
	fmt.Printf("DEBUG GetUsers - Page: %d, Limit: %d, Offset: %d, Status: %v, Role: %v, Search: %s\n",
		page, req.Limit, req.Offset, req.Status, req.Role, req.Search)

	response, err := h.adminUseCase.GetUsersPaginated(c.Request.Context(), req, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get users",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       response.Users,
		Pagination: response.Pagination,
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

// BulkUpdateUsers handles bulk user updates
func (h *AdminHandler) BulkUpdateUsers(c *gin.Context) {
	var req usecases.BulkUserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	response, err := h.adminUseCase.BulkUpdateUsers(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to bulk update users",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Bulk user update completed",
		Data:    response,
	})
}

// BulkDeleteUsers handles bulk user deletion
func (h *AdminHandler) BulkDeleteUsers(c *gin.Context) {
	var req usecases.BulkUserDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	response, err := h.adminUseCase.BulkDeleteUsers(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to bulk delete users",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Bulk user deletion completed",
		Data:    response,
	})
}

// BulkActivateUsers handles bulk user activation
func (h *AdminHandler) BulkActivateUsers(c *gin.Context) {
	var req usecases.BulkUserActivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	response, err := h.adminUseCase.BulkActivateUsers(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to bulk activate users",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Bulk user activation completed",
		Data:    response,
	})
}

// BulkDeactivateUsers handles bulk user deactivation
func (h *AdminHandler) BulkDeactivateUsers(c *gin.Context) {
	var req usecases.BulkUserDeactivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	response, err := h.adminUseCase.BulkDeactivateUsers(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to bulk deactivate users",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Bulk user deactivation completed",
		Data:    response,
	})
}

// BulkUpdateUserRoles handles bulk user role updates
func (h *AdminHandler) BulkUpdateUserRoles(c *gin.Context) {
	var req usecases.BulkUserRoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	response, err := h.adminUseCase.BulkUpdateUserRoles(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to bulk update user roles",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Bulk user role update completed",
		Data:    response,
	})
}

// SendUserNotification handles sending notification to a user
func (h *AdminHandler) SendUserNotification(c *gin.Context) {
	var req usecases.UserNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	response, err := h.adminUseCase.SendUserNotification(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to send notification",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Notification sent successfully",
		Data:    response,
	})
}

// SendBulkNotification handles sending notifications to multiple users
func (h *AdminHandler) SendBulkNotification(c *gin.Context) {
	var req usecases.BulkNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	response, err := h.adminUseCase.SendBulkNotification(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to send bulk notifications",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Bulk notifications sent",
		Data:    response,
	})
}

// SendUserEmail handles sending email to a user
func (h *AdminHandler) SendUserEmail(c *gin.Context) {
	var req usecases.UserEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	response, err := h.adminUseCase.SendUserEmail(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to send email",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Email sent successfully",
		Data:    response,
	})
}

// SendBulkEmail handles sending emails to multiple users
func (h *AdminHandler) SendBulkEmail(c *gin.Context) {
	var req usecases.BulkEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	response, err := h.adminUseCase.SendBulkEmail(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to send bulk emails",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Bulk emails sent",
		Data:    response,
	})
}

// CreateAnnouncement handles creating announcements
func (h *AdminHandler) CreateAnnouncement(c *gin.Context) {
	var req usecases.AnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	response, err := h.adminUseCase.CreateAnnouncement(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create announcement",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Announcement created successfully",
		Data:    response,
	})
}

// GetUserAuditLogs handles retrieving user audit logs
func (h *AdminHandler) GetUserAuditLogs(c *gin.Context) {
	var req usecases.UserAuditLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	// Set default values
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	response, err := h.adminUseCase.GetUserAuditLogs(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get audit logs",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Audit logs retrieved successfully",
		Data:    response,
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

// SearchCustomers performs advanced customer search with filtering and segmentation
// @Summary Search customers with advanced filters
// @Description Search customers with advanced filtering, segmentation, and analytics
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param query query string false "Search query"
// @Param role query string false "User role filter"
// @Param status query string false "User status filter"
// @Param is_active query bool false "Active status filter"
// @Param email_verified query bool false "Email verified filter"
// @Param phone_verified query bool false "Phone verified filter"
// @Param two_factor_enabled query bool false "Two factor enabled filter"
// @Param membership_tier query string false "Membership tier filter"
// @Param customer_segment query string false "Customer segment filter"
// @Param min_total_spent query number false "Minimum total spent filter"
// @Param max_total_spent query number false "Maximum total spent filter"
// @Param min_total_orders query int false "Minimum total orders filter"
// @Param max_total_orders query int false "Maximum total orders filter"
// @Param min_loyalty_points query int false "Minimum loyalty points filter"
// @Param max_loyalty_points query int false "Maximum loyalty points filter"
// @Param created_from query string false "Created from date filter (RFC3339)"
// @Param created_to query string false "Created to date filter (RFC3339)"
// @Param last_login_from query string false "Last login from date filter (RFC3339)"
// @Param last_login_to query string false "Last login to date filter (RFC3339)"
// @Param last_activity_from query string false "Last activity from date filter (RFC3339)"
// @Param last_activity_to query string false "Last activity to date filter (RFC3339)"
// @Param include_inactive query bool false "Include inactive customers"
// @Param include_unverified query bool false "Include unverified customers"
// @Param sort_by query string false "Sort by field" Enums(name,email,created_at,last_login,total_spent,total_orders,loyalty_points)
// @Param sort_order query string false "Sort order" Enums(asc,desc)
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} usecases.CustomerSearchResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/customers/search [get]
func (h *AdminHandler) SearchCustomers(c *gin.Context) {
	// Parse and validate pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))

	// Validate and normalize pagination for admin users
	page, limit, err := usecases.ValidateAndNormalizePaginationForEntity(page, limit, "admin_users")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Convert to offset for repository
	offset := (page - 1) * limit

	var req usecases.CustomerSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	// Override pagination parameters
	req.Limit = limit
	req.Offset = offset

	response, err := h.adminUseCase.SearchCustomersPaginated(c.Request.Context(), req, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to search customers",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       response.Customers,
		Pagination: response.Pagination,
	})
}

// GetCustomerSegments returns customer segmentation analysis
// @Summary Get customer segments
// @Description Get customer segmentation analysis with statistics
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} usecases.CustomerSegmentsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/customers/segments [get]
func (h *AdminHandler) GetCustomerSegments(c *gin.Context) {
	result, err := h.adminUseCase.GetCustomerSegments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get customer segments",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Customer segments retrieved successfully",
		Data:    result,
	})
}

// GetCustomerAnalytics returns comprehensive customer analytics
// @Summary Get customer analytics
// @Description Get comprehensive customer analytics and insights
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date_from query string false "Date from filter (RFC3339)"
// @Param date_to query string false "Date to filter (RFC3339)"
// @Param segment query string false "Customer segment filter"
// @Success 200 {object} usecases.CustomerAnalyticsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/customers/analytics [get]
func (h *AdminHandler) GetCustomerAnalytics(c *gin.Context) {
	var req usecases.CustomerAnalyticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	result, err := h.adminUseCase.GetCustomerAnalytics(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get customer analytics",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Customer analytics retrieved successfully",
		Data:    result,
	})
}

// GetHighValueCustomers returns high value customers
// @Summary Get high value customers
// @Description Get list of high value customers based on spending and order criteria
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(50)
// @Success 200 {object} usecases.HighValueCustomersResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/customers/high-value [get]
func (h *AdminHandler) GetHighValueCustomers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 1000 {
		limit = 50
	}

	result, err := h.adminUseCase.GetHighValueCustomers(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get high value customers",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "High value customers retrieved successfully",
		Data:    result,
	})
}

// GetCustomersBySegment returns customers filtered by segment
// @Summary Get customers by segment
// @Description Get customers filtered by specific segment with pagination
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param segment path string true "Customer segment" Enums(new,occasional,regular,loyal)
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} usecases.CustomersBySegmentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/customers/segments/{segment} [get]
func (h *AdminHandler) GetCustomersBySegment(c *gin.Context) {
	segment := c.Param("segment")
	if segment == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Segment parameter is required",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	result, err := h.adminUseCase.GetCustomersBySegment(c.Request.Context(), segment, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get customers by segment",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Customers by segment retrieved successfully",
		Data:    result,
	})
}

// GetCustomerLifetimeValue calculates and returns customer lifetime value
// @Summary Get customer lifetime value
// @Description Calculate and return customer lifetime value with analytics
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param customer_id path string true "Customer ID"
// @Success 200 {object} usecases.CustomerLifetimeValueResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/customers/{customer_id}/lifetime-value [get]
func (h *AdminHandler) GetCustomerLifetimeValue(c *gin.Context) {
	customerIDStr := c.Param("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid customer ID",
			Details: err.Error(),
		})
		return
	}

	result, err := h.adminUseCase.GetCustomerLifetimeValue(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get customer lifetime value",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Customer lifetime value retrieved successfully",
		Data:    result,
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

	// Handle pagination: if page is provided, calculate offset
	if req.Page > 0 {
		// Validate and normalize pagination for admin orders
		page, limit, err := usecases.ValidateAndNormalizePaginationForEntity(req.Page, req.Limit, "admin_orders")
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
			})
			return
		}
		req.Page = page
		req.Limit = limit
		req.Offset = (page - 1) * limit
	} else if req.Offset == 0 && req.Limit == 0 {
		// Set defaults if no pagination provided
		req.Page = 1
		req.Limit = 20 // AdminOrdersPerPage default
		req.Offset = 0
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
