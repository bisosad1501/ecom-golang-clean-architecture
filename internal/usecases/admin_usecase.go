package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
)

// AdminUseCase defines admin use cases
type AdminUseCase interface {
	// Dashboard
	GetDashboard(ctx context.Context, req AdminDashboardRequest) (*AdminDashboardResponse, error)
	GetSystemStats(ctx context.Context) (*SystemStatsResponse, error)

	// User management
	GetUsers(ctx context.Context, req AdminUsersRequest) (*AdminUsersResponse, error)
	UpdateUserStatus(ctx context.Context, userID uuid.UUID, status entities.UserStatus) error
	UpdateUserRole(ctx context.Context, userID uuid.UUID, role entities.UserRole) error
	GetUserActivity(ctx context.Context, userID uuid.UUID, req ActivityRequest) (*ActivityResponse, error)

	// Order management
	GetOrders(ctx context.Context, req AdminOrdersRequest) (*AdminOrdersResponse, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status entities.OrderStatus) error
	GetOrderDetails(ctx context.Context, orderID uuid.UUID) (*AdminOrderDetailsResponse, error)
	ProcessRefund(ctx context.Context, orderID uuid.UUID, amount float64, reason string) error

	// Product management
	GetProducts(ctx context.Context, req AdminProductsRequest) (*AdminProductsResponse, error)
	BulkUpdateProducts(ctx context.Context, req BulkUpdateProductsRequest) error
	GetProductAnalytics(ctx context.Context, productID uuid.UUID, period string) (*ProductAnalyticsResponse, error)

	// Content management
	ManageReviews(ctx context.Context, req ManageReviewsRequest) (*ManageReviewsResponse, error)
	UpdateReviewStatus(ctx context.Context, reviewID uuid.UUID, status entities.ReviewStatus) error

	// System management
	GetSystemLogs(ctx context.Context, req SystemLogsRequest) (*SystemLogsResponse, error)
	GetAuditLogs(ctx context.Context, req AuditLogsRequest) (*AuditLogsResponse, error)
	BackupDatabase(ctx context.Context) (*BackupResponse, error)

	// Reports
	GenerateReport(ctx context.Context, req GenerateReportRequest) (*ReportResponse, error)
	GetReports(ctx context.Context, req GetReportsRequest) (*ReportsListResponse, error)
	DownloadReport(ctx context.Context, reportID uuid.UUID) (*DownloadResponse, error)
}

type adminUseCase struct {
	userRepo      repositories.UserRepository
	orderRepo     repositories.OrderRepository
	productRepo   repositories.ProductRepository
	reviewRepo    repositories.ReviewRepository
	analyticsRepo repositories.AnalyticsRepository
	inventoryRepo repositories.InventoryRepository
	paymentRepo   repositories.PaymentRepository
	auditRepo     repositories.AuditRepository
}

// NewAdminUseCase creates a new admin use case
func NewAdminUseCase(
	userRepo repositories.UserRepository,
	orderRepo repositories.OrderRepository,
	productRepo repositories.ProductRepository,
	reviewRepo repositories.ReviewRepository,
	analyticsRepo repositories.AnalyticsRepository,
	inventoryRepo repositories.InventoryRepository,
	paymentRepo repositories.PaymentRepository,
	auditRepo repositories.AuditRepository,
) AdminUseCase {
	return &adminUseCase{
		userRepo:      userRepo,
		orderRepo:     orderRepo,
		productRepo:   productRepo,
		reviewRepo:    reviewRepo,
		analyticsRepo: analyticsRepo,
		inventoryRepo: inventoryRepo,
		paymentRepo:   paymentRepo,
		auditRepo:     auditRepo,
	}
}

// Request types
type AdminDashboardRequest struct {
	Period   string     `json:"period,omitempty" validate:"omitempty,oneof=today week month year"`
	DateFrom *time.Time `json:"date_from,omitempty"`
	DateTo   *time.Time `json:"date_to,omitempty"`
}

type AdminUsersRequest struct {
	Status    *entities.UserStatus `json:"status,omitempty"`
	Role      *entities.UserRole   `json:"role,omitempty"`
	Search    string               `json:"search,omitempty"`
	SortBy    string               `json:"sort_by,omitempty" validate:"omitempty,oneof=name email created_at last_login"`
	SortOrder string               `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
	Limit     int                  `json:"limit" validate:"min=1,max=100"`
	Offset    int                  `json:"offset" validate:"min=0"`
}

type AdminOrdersRequest struct {
	Status        *entities.OrderStatus   `json:"status,omitempty"`
	PaymentStatus *entities.PaymentStatus `json:"payment_status,omitempty"`
	UserID        *uuid.UUID              `json:"user_id,omitempty"`
	DateFrom      *time.Time              `json:"date_from,omitempty"`
	DateTo        *time.Time              `json:"date_to,omitempty"`
	Search        string                  `json:"search,omitempty"`
	SortBy        string                  `json:"sort_by,omitempty" validate:"omitempty,oneof=created_at total status"`
	SortOrder     string                  `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
	Limit         int                     `json:"limit" validate:"min=1,max=100"`
	Offset        int                     `json:"offset" validate:"min=0"`
}

type AdminProductsRequest struct {
	Status     *entities.ProductStatus `json:"status,omitempty"`
	CategoryID *uuid.UUID              `json:"category_id,omitempty"`
	Search     string                  `json:"search,omitempty"`
	LowStock   *bool                   `json:"low_stock,omitempty"`
	SortBy     string                  `json:"sort_by,omitempty" validate:"omitempty,oneof=name price stock created_at"`
	SortOrder  string                  `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
	Limit      int                     `json:"limit" validate:"min=1,max=100"`
	Offset     int                     `json:"offset" validate:"min=0"`
}

type BulkUpdateProductsRequest struct {
	ProductIDs []uuid.UUID `json:"product_ids" validate:"required,min=1"`
	Updates    struct {
		Status       *entities.ProductStatus `json:"status,omitempty"`
		CategoryID   *uuid.UUID              `json:"category_id,omitempty"`
		Price        *float64                `json:"price,omitempty"`
		ComparePrice *float64                `json:"compare_price,omitempty"`
		IsActive     *bool                   `json:"is_active,omitempty"`
	} `json:"updates" validate:"required"`
}

type ManageReviewsRequest struct {
	Status    *entities.ReviewStatus `json:"status,omitempty"`
	ProductID *uuid.UUID             `json:"product_id,omitempty"`
	UserID    *uuid.UUID             `json:"user_id,omitempty"`
	Rating    *int                   `json:"rating,omitempty"`
	Flagged   *bool                  `json:"flagged,omitempty"`
	SortBy    string                 `json:"sort_by,omitempty" validate:"omitempty,oneof=created_at rating helpful_votes"`
	SortOrder string                 `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
	Limit     int                    `json:"limit" validate:"min=1,max=100"`
	Offset    int                    `json:"offset" validate:"min=0"`
}

type ActivityRequest struct {
	Type     string     `json:"type,omitempty"`
	DateFrom *time.Time `json:"date_from,omitempty"`
	DateTo   *time.Time `json:"date_to,omitempty"`
	Limit    int        `json:"limit" validate:"min=1,max=100"`
	Offset   int        `json:"offset" validate:"min=0"`
}

type SystemLogsRequest struct {
	Level    string     `json:"level,omitempty" validate:"omitempty,oneof=debug info warn error"`
	Service  string     `json:"service,omitempty"`
	DateFrom *time.Time `json:"date_from,omitempty"`
	DateTo   *time.Time `json:"date_to,omitempty"`
	Search   string     `json:"search,omitempty"`
	Limit    int        `json:"limit" validate:"min=1,max=1000"`
	Offset   int        `json:"offset" validate:"min=0"`
}

type AuditLogsRequest struct {
	UserID   *uuid.UUID `json:"user_id,omitempty"`
	Action   string     `json:"action,omitempty"`
	Resource string     `json:"resource,omitempty"`
	DateFrom *time.Time `json:"date_from,omitempty"`
	DateTo   *time.Time `json:"date_to,omitempty"`
	Limit    int        `json:"limit" validate:"min=1,max=1000"`
	Offset   int        `json:"offset" validate:"min=0"`
}

type GenerateReportRequest struct {
	Type      string                 `json:"type" validate:"required,oneof=sales products users inventory payments"`
	Format    string                 `json:"format" validate:"required,oneof=csv excel pdf"`
	DateFrom  time.Time              `json:"date_from" validate:"required"`
	DateTo    time.Time              `json:"date_to" validate:"required"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
	GroupBy   string                 `json:"group_by,omitempty"`
	CreatedBy uuid.UUID              `json:"created_by" validate:"required"`
}

type GetReportsRequest struct {
	Type      string     `json:"type,omitempty"`
	Status    string     `json:"status,omitempty"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	DateFrom  *time.Time `json:"date_from,omitempty"`
	DateTo    *time.Time `json:"date_to,omitempty"`
	Limit     int        `json:"limit" validate:"min=1,max=100"`
	Offset    int        `json:"offset" validate:"min=0"`
}

// Response types
type AdminDashboardResponse struct {
	Overview struct {
		TotalRevenue    float64 `json:"total_revenue"`    // Net revenue (current)
		GrossRevenue    float64 `json:"gross_revenue"`    // Before discounts
		ProductRevenue  float64 `json:"product_revenue"`  // Only product sales
		TaxCollected    float64 `json:"tax_collected"`    // Total tax amount
		ShippingRevenue float64 `json:"shipping_revenue"` // Shipping fees
		DiscountsGiven  float64 `json:"discounts_given"`  // Total discounts
		TotalOrders     int64   `json:"total_orders"`
		TotalCustomers  int64   `json:"total_customers"`
		TotalProducts   int64   `json:"total_products"`
		PendingOrders   int64   `json:"pending_orders"`
		LowStockItems   int64   `json:"low_stock_items"`
		PendingReviews  int64   `json:"pending_reviews"`
		ActiveUsers     int64   `json:"active_users"`
	} `json:"overview"`

	Charts struct {
		RevenueChart []struct {
			Date    string  `json:"date"`
			Revenue float64 `json:"revenue"`
			Orders  int64   `json:"orders"`
		} `json:"revenue_chart"`

		OrdersChart []struct {
			Date   string `json:"date"`
			Orders int64  `json:"orders"`
		} `json:"orders_chart"`

		TopProducts []struct {
			ProductID   uuid.UUID `json:"product_id"`
			ProductName string    `json:"product_name"`
			Revenue     float64   `json:"revenue"`
			Quantity    int64     `json:"quantity"`
		} `json:"top_products"`

		TopCategories []struct {
			CategoryID   uuid.UUID `json:"category_id"`
			CategoryName string    `json:"category_name"`
			Revenue      float64   `json:"revenue"`
			Orders       int64     `json:"orders"`
		} `json:"top_categories"`
	} `json:"charts"`

	RecentActivity []struct {
		Type        string    `json:"type"`
		Description string    `json:"description"`
		UserID      uuid.UUID `json:"user_id"`
		UserName    string    `json:"user_name"`
		Timestamp   time.Time `json:"timestamp"`
	} `json:"recent_activity"`

	RecentOrders []struct {
		ID           uuid.UUID `json:"id"`
		OrderNumber  string    `json:"order_number"`
		Status       string    `json:"status"`
		Total        float64   `json:"total"`
		TotalAmount  float64   `json:"total_amount"`
		CustomerName string    `json:"customer_name"`
		CreatedAt    time.Time `json:"created_at"`
		User         *struct {
			ID        uuid.UUID `json:"id"`
			FirstName string    `json:"first_name"`
			LastName  string    `json:"last_name"`
		} `json:"user,omitempty"`
	} `json:"recent_orders"`
}

type SystemStatsResponse struct {
	Database struct {
		TotalSize       string `json:"total_size"`
		TableCount      int    `json:"table_count"`
		ConnectionCount int    `json:"connection_count"`
		QueryCount      int64  `json:"query_count"`
	} `json:"database"`

	Server struct {
		Uptime       string  `json:"uptime"`
		CPUUsage     float64 `json:"cpu_usage"`
		MemoryUsage  float64 `json:"memory_usage"`
		DiskUsage    float64 `json:"disk_usage"`
		RequestCount int64   `json:"request_count"`
		ErrorRate    float64 `json:"error_rate"`
	} `json:"server"`

	Cache struct {
		HitRate     float64 `json:"hit_rate"`
		MissRate    float64 `json:"miss_rate"`
		KeyCount    int64   `json:"key_count"`
		MemoryUsage string  `json:"memory_usage"`
	} `json:"cache"`
}

type AdminUsersResponse struct {
	Users []struct {
		ID         uuid.UUID           `json:"id"`
		Email      string              `json:"email"`
		FirstName  string              `json:"first_name"`
		LastName   string              `json:"last_name"`
		Role       entities.UserRole   `json:"role"`
		Status     entities.UserStatus `json:"status"`
		LastLogin  *time.Time          `json:"last_login"`
		OrderCount int64               `json:"order_count"`
		TotalSpent float64             `json:"total_spent"`
		CreatedAt  time.Time           `json:"created_at"`
	} `json:"users"`
	Total      int64           `json:"total"`
	Pagination *PaginationInfo `json:"pagination"`
}

type AdminOrdersResponse struct {
	Orders []struct {
		ID            uuid.UUID              `json:"id"`
		OrderNumber   string                 `json:"order_number"`
		UserID        uuid.UUID              `json:"user_id"`
		UserName      string                 `json:"user_name"`
		UserEmail     string                 `json:"user_email"`
		Status        entities.OrderStatus   `json:"status"`
		PaymentStatus entities.PaymentStatus `json:"payment_status"`
		Total         float64                `json:"total"`
		ItemCount     int                    `json:"item_count"`
		CreatedAt     time.Time              `json:"created_at"`
		UpdatedAt     time.Time              `json:"updated_at"`
	} `json:"orders"`
	Total      int64           `json:"total"`
	Pagination *PaginationInfo `json:"pagination"`
}

type AdminOrderDetailsResponse struct {
	Order struct {
		ID             uuid.UUID              `json:"id"`
		OrderNumber    string                 `json:"order_number"`
		Status         entities.OrderStatus   `json:"status"`
		PaymentStatus  entities.PaymentStatus `json:"payment_status"`
		Subtotal       float64                `json:"subtotal"`
		TaxAmount      float64                `json:"tax_amount"`
		ShippingAmount float64                `json:"shipping_amount"`
		DiscountAmount float64                `json:"discount_amount"`
		Total          float64                `json:"total"`
		CreatedAt      time.Time              `json:"created_at"`
		UpdatedAt      time.Time              `json:"updated_at"`
	} `json:"order"`

	Customer struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Phone     string    `json:"phone"`
	} `json:"customer"`

	Items []struct {
		ID          uuid.UUID `json:"id"`
		ProductID   uuid.UUID `json:"product_id"`
		ProductName string    `json:"product_name"`
		ProductSKU  string    `json:"product_sku"`
		Quantity    int       `json:"quantity"`
		Price       float64   `json:"price"`
		Total       float64   `json:"total"`
	} `json:"items"`

	ShippingAddress *struct {
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Company      string `json:"company"`
		AddressLine1 string `json:"address_line_1"`
		AddressLine2 string `json:"address_line_2"`
		City         string `json:"city"`
		State        string `json:"state"`
		PostalCode   string `json:"postal_code"`
		Country      string `json:"country"`
		Phone        string `json:"phone"`
	} `json:"shipping_address,omitempty"`

	BillingAddress *struct {
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Company      string `json:"company"`
		AddressLine1 string `json:"address_line_1"`
		AddressLine2 string `json:"address_line_2"`
		City         string `json:"city"`
		State        string `json:"state"`
		PostalCode   string `json:"postal_code"`
		Country      string `json:"country"`
		Phone        string `json:"phone"`
	} `json:"billing_address,omitempty"`

	Payments []struct {
		ID            uuid.UUID              `json:"id"`
		Amount        float64                `json:"amount"`
		Method        entities.PaymentMethod `json:"method"`
		Status        entities.PaymentStatus `json:"status"`
		TransactionID string                 `json:"transaction_id"`
		ProcessedAt   *time.Time             `json:"processed_at"`
	} `json:"payments"`

	Timeline []struct {
		Event       string     `json:"event"`
		Description string     `json:"description"`
		Timestamp   time.Time  `json:"timestamp"`
		UserID      *uuid.UUID `json:"user_id,omitempty"`
		UserName    string     `json:"user_name,omitempty"`
	} `json:"timeline"`
}

type AdminProductsResponse struct {
	Products []struct {
		ID            uuid.UUID              `json:"id"`
		Name          string                 `json:"name"`
		SKU           string                 `json:"sku"`
		Price         float64                `json:"price"`
		ComparePrice  float64                `json:"compare_price"`
		Status        entities.ProductStatus `json:"status"`
		StockQuantity int                    `json:"stock_quantity"`
		CategoryID    uuid.UUID              `json:"category_id"`
		CategoryName  string                 `json:"category_name"`
		ViewCount     int64                  `json:"view_count"`
		SalesCount    int64                  `json:"sales_count"`
		Revenue       float64                `json:"revenue"`
		CreatedAt     time.Time              `json:"created_at"`
		UpdatedAt     time.Time              `json:"updated_at"`
	} `json:"products"`
	Total      int64           `json:"total"`
	Pagination *PaginationInfo `json:"pagination"`
}

type ProductAnalyticsResponse struct {
	ProductID uuid.UUID `json:"product_id"`
	Period    string    `json:"period"`

	Metrics struct {
		Views          int64   `json:"views"`
		Sales          int64   `json:"sales"`
		Revenue        float64 `json:"revenue"`
		ConversionRate float64 `json:"conversion_rate"`
		AverageRating  float64 `json:"average_rating"`
		ReviewCount    int64   `json:"review_count"`
	} `json:"metrics"`

	Charts struct {
		ViewsChart []struct {
			Date  string `json:"date"`
			Views int64  `json:"views"`
		} `json:"views_chart"`

		SalesChart []struct {
			Date  string `json:"date"`
			Sales int64  `json:"sales"`
		} `json:"sales_chart"`

		RevenueChart []struct {
			Date    string  `json:"date"`
			Revenue float64 `json:"revenue"`
		} `json:"revenue_chart"`
	} `json:"charts"`
}

type ManageReviewsResponse struct {
	Reviews []struct {
		ID           uuid.UUID             `json:"id"`
		ProductID    uuid.UUID             `json:"product_id"`
		ProductName  string                `json:"product_name"`
		UserID       uuid.UUID             `json:"user_id"`
		UserName     string                `json:"user_name"`
		Rating       int                   `json:"rating"`
		Title        string                `json:"title"`
		Content      string                `json:"content"`
		Status       entities.ReviewStatus `json:"status"`
		HelpfulVotes int                   `json:"helpful_votes"`
		TotalVotes   int                   `json:"total_votes"`
		IsFlagged    bool                  `json:"is_flagged"`
		CreatedAt    time.Time             `json:"created_at"`
	} `json:"reviews"`
	Total      int64           `json:"total"`
	Pagination *PaginationInfo `json:"pagination"`
}

type ActivityResponse struct {
	Activities []struct {
		ID          uuid.UUID `json:"id"`
		Type        string    `json:"type"`
		Description string    `json:"description"`
		IPAddress   string    `json:"ip_address"`
		UserAgent   string    `json:"user_agent"`
		Metadata    string    `json:"metadata"`
		CreatedAt   time.Time `json:"created_at"`
	} `json:"activities"`
	Total      int64           `json:"total"`
	Pagination *PaginationInfo `json:"pagination"`
}

type SystemLogsResponse struct {
	Logs []struct {
		ID        uuid.UUID `json:"id"`
		Level     string    `json:"level"`
		Service   string    `json:"service"`
		Message   string    `json:"message"`
		Context   string    `json:"context"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"logs"`
	Total      int64           `json:"total"`
	Pagination *PaginationInfo `json:"pagination"`
}

type AuditLogsResponse struct {
	Logs []struct {
		ID         uuid.UUID `json:"id"`
		UserID     uuid.UUID `json:"user_id"`
		UserName   string    `json:"user_name"`
		Action     string    `json:"action"`
		Resource   string    `json:"resource"`
		ResourceID string    `json:"resource_id"`
		OldValues  string    `json:"old_values"`
		NewValues  string    `json:"new_values"`
		IPAddress  string    `json:"ip_address"`
		UserAgent  string    `json:"user_agent"`
		CreatedAt  time.Time `json:"created_at"`
	} `json:"logs"`
	Total      int64           `json:"total"`
	Pagination *PaginationInfo `json:"pagination"`
}

type BackupResponse struct {
	BackupID    uuid.UUID `json:"backup_id"`
	Filename    string    `json:"filename"`
	Size        int64     `json:"size"`
	Status      string    `json:"status"`
	DownloadURL string    `json:"download_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type ReportResponse struct {
	ReportID    uuid.UUID  `json:"report_id"`
	Type        string     `json:"type"`
	Format      string     `json:"format"`
	Status      string     `json:"status"`
	Progress    int        `json:"progress"`
	DownloadURL string     `json:"download_url,omitempty"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type ReportsListResponse struct {
	Reports    []*ReportResponse `json:"reports"`
	Total      int64             `json:"total"`
	Pagination *PaginationInfo   `json:"pagination"`
}

type DownloadResponse struct {
	URL       string    `json:"url"`
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	ExpiresAt time.Time `json:"expires_at"`
}

type AuditLogFilters struct {
	UserID   *uuid.UUID `json:"user_id,omitempty"`
	Action   string     `json:"action,omitempty"`
	Resource string     `json:"resource,omitempty"`
	DateFrom *time.Time `json:"date_from,omitempty"`
	DateTo   *time.Time `json:"date_to,omitempty"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
}

type AuditLogSummary struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

type ProductAnalyticsRequest struct {
	ProductID *uuid.UUID `json:"product_id,omitempty"`
	Period    string     `json:"period,omitempty"`
	DateFrom  *time.Time `json:"date_from,omitempty"`
	DateTo    *time.Time `json:"date_to,omitempty"`
}

type OrderDetailsResponse struct {
	OrderID      uuid.UUID     `json:"order_id"`
	OrderNumber  string        `json:"order_number"`
	CustomerID   uuid.UUID     `json:"customer_id"`
	CustomerName string        `json:"customer_name"`
	Status       string        `json:"status"`
	Total        float64       `json:"total"`
	Items        []interface{} `json:"items"`
	CreatedAt    time.Time     `json:"created_at"`
}

type GetReportsResponse struct {
	Reports []interface{} `json:"reports"`
	Total   int64         `json:"total"`
}

// GetDashboard gets admin dashboard data
func (uc *adminUseCase) GetDashboard(ctx context.Context, req AdminDashboardRequest) (*AdminDashboardResponse, error) {
	// Set default period if not provided
	if req.Period == "" {
		req.Period = "month"
	}

	// Calculate date range based on period
	now := time.Now()
	var dateFrom, dateTo time.Time

	switch req.Period {
	case "today":
		dateFrom = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		dateTo = dateFrom.Add(24 * time.Hour)
	case "week":
		dateFrom = now.AddDate(0, 0, -7)
		dateTo = now
	case "month":
		dateFrom = now.AddDate(0, -1, 0)
		dateTo = now
	case "year":
		dateFrom = now.AddDate(-1, 0, 0)
		dateTo = now
	default:
		if req.DateFrom != nil && req.DateTo != nil {
			dateFrom = *req.DateFrom
			dateTo = *req.DateTo
		} else {
			dateFrom = now.AddDate(0, -1, 0)
			dateTo = now
		}
	}

	// Use date range for filtering (in real implementation)
	_ = dateFrom
	_ = dateTo

	// Get overview metrics
	totalRevenue, _ := uc.orderRepo.GetTotalRevenue(ctx)       // Net revenue (current)
	grossRevenue, _ := uc.orderRepo.GetGrossRevenue(ctx)       // Before discounts
	productRevenue, _ := uc.orderRepo.GetProductRevenue(ctx)   // Only products
	taxCollected, _ := uc.orderRepo.GetTaxCollected(ctx)       // Tax amount
	shippingRevenue, _ := uc.orderRepo.GetShippingRevenue(ctx) // Shipping fees
	discountsGiven, _ := uc.orderRepo.GetDiscountsGiven(ctx)   // Discounts
	totalOrders, _ := uc.orderRepo.CountOrders(ctx)
	totalCustomers, _ := uc.userRepo.CountUsers(ctx)
	totalProducts, _ := uc.productRepo.CountProducts(ctx)
	pendingOrders, _ := uc.orderRepo.CountOrdersByStatus(ctx, entities.OrderStatusPending)
	lowStockItems, _ := uc.inventoryRepo.CountLowStockItems(ctx)
	pendingReviews, _ := uc.reviewRepo.CountReviewsByStatus(ctx, entities.ReviewStatusPending)
	activeUsers, _ := uc.userRepo.CountActiveUsers(ctx)

	response := &AdminDashboardResponse{
		Overview: struct {
			TotalRevenue    float64 `json:"total_revenue"`    // Net revenue (current)
			GrossRevenue    float64 `json:"gross_revenue"`    // Before discounts
			ProductRevenue  float64 `json:"product_revenue"`  // Only product sales
			TaxCollected    float64 `json:"tax_collected"`    // Total tax amount
			ShippingRevenue float64 `json:"shipping_revenue"` // Shipping fees
			DiscountsGiven  float64 `json:"discounts_given"`  // Total discounts
			TotalOrders     int64   `json:"total_orders"`
			TotalCustomers  int64   `json:"total_customers"`
			TotalProducts   int64   `json:"total_products"`
			PendingOrders   int64   `json:"pending_orders"`
			LowStockItems   int64   `json:"low_stock_items"`
			PendingReviews  int64   `json:"pending_reviews"`
			ActiveUsers     int64   `json:"active_users"`
		}{
			TotalRevenue:    totalRevenue,
			GrossRevenue:    grossRevenue,
			ProductRevenue:  productRevenue,
			TaxCollected:    taxCollected,
			ShippingRevenue: shippingRevenue,
			DiscountsGiven:  discountsGiven,
			TotalOrders:     totalOrders,
			TotalCustomers:  totalCustomers,
			TotalProducts:   totalProducts,
			PendingOrders:   pendingOrders,
			LowStockItems:   lowStockItems,
			PendingReviews:  pendingReviews,
			ActiveUsers:     activeUsers,
		},
	}
	// Get recent orders (limit to 5 for dashboard)
	recentOrdersReq := AdminOrdersRequest{
		Limit:     5,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
	recentOrdersResp, err := uc.GetOrders(ctx, recentOrdersReq)
	if err == nil && recentOrdersResp != nil {
		for _, order := range recentOrdersResp.Orders {
			recentOrder := struct {
				ID           uuid.UUID `json:"id"`
				OrderNumber  string    `json:"order_number"`
				Status       string    `json:"status"`
				Total        float64   `json:"total"`
				TotalAmount  float64   `json:"total_amount"`
				CustomerName string    `json:"customer_name"`
				CreatedAt    time.Time `json:"created_at"`
				User         *struct {
					ID        uuid.UUID `json:"id"`
					FirstName string    `json:"first_name"`
					LastName  string    `json:"last_name"`
				} `json:"user,omitempty"`
			}{
				ID:           order.ID,
				OrderNumber:  order.OrderNumber,
				Status:       string(order.Status),
				Total:        order.Total,
				TotalAmount:  order.Total,
				CustomerName: order.UserName,
				CreatedAt:    order.CreatedAt,
			}

			// Add user info if available (use UserName for now)
			if order.UserName != "" {
				names := strings.Split(order.UserName, " ")
				firstName := names[0]
				lastName := ""
				if len(names) > 1 {
					lastName = strings.Join(names[1:], " ")
				}

				recentOrder.User = &struct {
					ID        uuid.UUID `json:"id"`
					FirstName string    `json:"first_name"`
					LastName  string    `json:"last_name"`
				}{
					ID:        order.UserID,
					FirstName: firstName,
					LastName:  lastName,
				}
			}

			response.RecentOrders = append(response.RecentOrders, recentOrder)
		}
	}

	// Get chart data (simplified implementation)
	// In a real implementation, you would fetch actual chart data from repositories
	response.Charts.RevenueChart = []struct {
		Date    string  `json:"date"`
		Revenue float64 `json:"revenue"`
		Orders  int64   `json:"orders"`
	}{
		{Date: "2024-01-01", Revenue: 10000, Orders: 50},
		{Date: "2024-01-02", Revenue: 12000, Orders: 60},
		// Add more data points...
	}

	return response, nil
}

// BackupDatabase creates a database backup
func (uc *adminUseCase) BackupDatabase(ctx context.Context) (*BackupResponse, error) {
	// In a real implementation, this would trigger a database backup
	// For now, return a mock response
	response := &BackupResponse{
		BackupID:    uuid.New(),
		Status:      "completed",
		Filename:    "db_backup_" + time.Now().Format("20060102_150405") + ".sql",
		Size:        131072000, // 125.5 MB in bytes
		DownloadURL: "/api/v1/admin/backups/" + uuid.New().String(),
		CreatedAt:   time.Now(),
	}

	return response, nil
}

// BulkUpdateProducts updates multiple products
func (uc *adminUseCase) BulkUpdateProducts(ctx context.Context, req BulkUpdateProductsRequest) error {
	// Mock implementation for bulk update
	// In real implementation, this would update multiple products
	return nil
}

// GenerateReport generates a report
func (uc *adminUseCase) GenerateReport(ctx context.Context, req GenerateReportRequest) (*ReportResponse, error) {
	// Mock implementation for generate report
	response := &ReportResponse{
		ReportID:    uuid.New(),
		Type:        req.Type,
		Format:      req.Format,
		Status:      "completed",
		Progress:    100,
		DownloadURL: "/api/v1/admin/reports/" + uuid.New().String() + "/download",
		CreatedBy:   uuid.New(), // Should be current user ID
		CreatedAt:   time.Now(),
		CompletedAt: &time.Time{},
	}
	*response.CompletedAt = time.Now()
	return response, nil
}

// DownloadReport downloads a report
func (uc *adminUseCase) DownloadReport(ctx context.Context, reportID uuid.UUID) (*DownloadResponse, error) {
	// Mock implementation for download report
	response := &DownloadResponse{
		Filename: "report_" + reportID.String() + ".pdf",
		Size:     1024000, // 1MB
	}
	return response, nil
}

// GetAuditLogs gets audit logs
func (uc *adminUseCase) GetAuditLogs(ctx context.Context, req AuditLogsRequest) (*AuditLogsResponse, error) {
	// Mock implementation for audit logs
	response := &AuditLogsResponse{
		Logs: []struct {
			ID         uuid.UUID `json:"id"`
			UserID     uuid.UUID `json:"user_id"`
			UserName   string    `json:"user_name"`
			Action     string    `json:"action"`
			Resource   string    `json:"resource"`
			ResourceID string    `json:"resource_id"`
			OldValues  string    `json:"old_values"`
			NewValues  string    `json:"new_values"`
			IPAddress  string    `json:"ip_address"`
			UserAgent  string    `json:"user_agent"`
			CreatedAt  time.Time `json:"created_at"`
		}{
			{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				UserName:   "John Doe",
				Action:     "login",
				Resource:   "user",
				ResourceID: uuid.New().String(),
				IPAddress:  "192.168.1.1",
				UserAgent:  "Mozilla/5.0",
				CreatedAt:  time.Now().Add(-1 * time.Hour),
			},
			{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				UserName:   "Jane Smith",
				Action:     "create_product",
				Resource:   "product",
				ResourceID: uuid.New().String(),
				IPAddress:  "192.168.1.2",
				UserAgent:  "Mozilla/5.0",
				CreatedAt:  time.Now().Add(-2 * time.Hour),
			},
		},
		Total: 100,
	}

	return response, nil
}

// GetOrderDetails gets order details
func (uc *adminUseCase) GetOrderDetails(ctx context.Context, orderID uuid.UUID) (*AdminOrderDetailsResponse, error) {
	// Get order from repository with preloaded relationships
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Get user information
	user, err := uc.userRepo.GetByID(ctx, order.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Convert order items
	items := make([]struct {
		ID          uuid.UUID `json:"id"`
		ProductID   uuid.UUID `json:"product_id"`
		ProductName string    `json:"product_name"`
		ProductSKU  string    `json:"product_sku"`
		Quantity    int       `json:"quantity"`
		Price       float64   `json:"price"`
		Total       float64   `json:"total"`
	}, len(order.Items))

	for i, item := range order.Items {
		items[i] = struct {
			ID          uuid.UUID `json:"id"`
			ProductID   uuid.UUID `json:"product_id"`
			ProductName string    `json:"product_name"`
			ProductSKU  string    `json:"product_sku"`
			Quantity    int       `json:"quantity"`
			Price       float64   `json:"price"`
			Total       float64   `json:"total"`
		}{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			ProductSKU:  item.ProductSKU,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Total:       item.Total,
		}
	}

	// Build response
	response := &AdminOrderDetailsResponse{
		Order: struct {
			ID             uuid.UUID              `json:"id"`
			OrderNumber    string                 `json:"order_number"`
			Status         entities.OrderStatus   `json:"status"`
			PaymentStatus  entities.PaymentStatus `json:"payment_status"`
			Subtotal       float64                `json:"subtotal"`
			TaxAmount      float64                `json:"tax_amount"`
			ShippingAmount float64                `json:"shipping_amount"`
			DiscountAmount float64                `json:"discount_amount"`
			Total          float64                `json:"total"`
			CreatedAt      time.Time              `json:"created_at"`
			UpdatedAt      time.Time              `json:"updated_at"`
		}{
			ID:             order.ID,
			OrderNumber:    order.OrderNumber,
			Status:         order.Status,
			PaymentStatus:  order.PaymentStatus,
			Subtotal:       order.Subtotal,
			TaxAmount:      order.TaxAmount,
			ShippingAmount: order.ShippingAmount,
			DiscountAmount: order.DiscountAmount,
			Total:          order.Total,
			CreatedAt:      order.CreatedAt,
			UpdatedAt:      order.UpdatedAt,
		},
		Customer: struct {
			ID        uuid.UUID `json:"id"`
			Email     string    `json:"email"`
			FirstName string    `json:"first_name"`
			LastName  string    `json:"last_name"`
			Phone     string    `json:"phone"`
		}{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
		},
		Items: items,
	}

	// Add shipping address if exists
	if order.ShippingAddress != nil {
		response.ShippingAddress = &struct {
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Company      string `json:"company"`
			AddressLine1 string `json:"address_line_1"`
			AddressLine2 string `json:"address_line_2"`
			City         string `json:"city"`
			State        string `json:"state"`
			PostalCode   string `json:"postal_code"`
			Country      string `json:"country"`
			Phone        string `json:"phone"`
		}{
			FirstName:    order.ShippingAddress.FirstName,
			LastName:     order.ShippingAddress.LastName,
			Company:      order.ShippingAddress.Company,
			AddressLine1: order.ShippingAddress.Address1,
			AddressLine2: order.ShippingAddress.Address2,
			City:         order.ShippingAddress.City,
			State:        order.ShippingAddress.State,
			PostalCode:   order.ShippingAddress.ZipCode,
			Country:      order.ShippingAddress.Country,
			Phone:        order.ShippingAddress.Phone,
		}
	}

	// Add billing address if exists
	if order.BillingAddress != nil {
		response.BillingAddress = &struct {
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Company      string `json:"company"`
			AddressLine1 string `json:"address_line_1"`
			AddressLine2 string `json:"address_line_2"`
			City         string `json:"city"`
			State        string `json:"state"`
			PostalCode   string `json:"postal_code"`
			Country      string `json:"country"`
			Phone        string `json:"phone"`
		}{
			FirstName:    order.BillingAddress.FirstName,
			LastName:     order.BillingAddress.LastName,
			Company:      order.BillingAddress.Company,
			AddressLine1: order.BillingAddress.Address1,
			AddressLine2: order.BillingAddress.Address2,
			City:         order.BillingAddress.City,
			State:        order.BillingAddress.State,
			PostalCode:   order.BillingAddress.ZipCode,
			Country:      order.BillingAddress.Country,
			Phone:        order.BillingAddress.Phone,
		}
	}

	return response, nil
}

// GetOrders gets orders
func (uc *adminUseCase) GetOrders(ctx context.Context, req AdminOrdersRequest) (*AdminOrdersResponse, error) {
	// Build search parameters for order repository
	searchParams := repositories.OrderSearchParams{
		SortBy:    "created_at",
		SortOrder: "desc",
		Limit:     req.Limit,
		Offset:    req.Offset,
	}

	// Add filters if provided
	if req.Status != nil {
		searchParams.Status = req.Status
	}

	if req.PaymentStatus != nil {
		searchParams.PaymentStatus = req.PaymentStatus
	}

	if req.UserID != nil {
		searchParams.UserID = req.UserID
	}

	if req.DateFrom != nil {
		searchParams.StartDate = req.DateFrom
	}

	if req.DateTo != nil {
		searchParams.EndDate = req.DateTo
	}

	// Get orders from repository
	orders, err := uc.orderRepo.Search(ctx, searchParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search orders: %w", err)
	}

	// Preload items for each order to get accurate item count
	for i, order := range orders {
		// Load order items if not already loaded
		if len(order.Items) == 0 {
			fullOrder, err := uc.orderRepo.GetByID(ctx, order.ID)
			if err == nil && fullOrder != nil {
				orders[i].Items = fullOrder.Items
			}
		}
	}

	// Get total count for pagination
	totalCount, err := uc.orderRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count orders: %w", err)
	}

	// Convert to response format
	orderResponses := make([]struct {
		ID            uuid.UUID              `json:"id"`
		OrderNumber   string                 `json:"order_number"`
		UserID        uuid.UUID              `json:"user_id"`
		UserName      string                 `json:"user_name"`
		UserEmail     string                 `json:"user_email"`
		Status        entities.OrderStatus   `json:"status"`
		PaymentStatus entities.PaymentStatus `json:"payment_status"`
		Total         float64                `json:"total"`
		ItemCount     int                    `json:"item_count"`
		CreatedAt     time.Time              `json:"created_at"`
		UpdatedAt     time.Time              `json:"updated_at"`
	}, len(orders))

	for i, order := range orders {
		// Get user information
		user, err := uc.userRepo.GetByID(ctx, order.UserID)
		userName := "Unknown User"
		userEmail := "unknown@example.com"
		if err == nil && user != nil {
			userName = user.GetFullName()
			userEmail = user.Email
		}

		orderResponses[i] = struct {
			ID            uuid.UUID              `json:"id"`
			OrderNumber   string                 `json:"order_number"`
			UserID        uuid.UUID              `json:"user_id"`
			UserName      string                 `json:"user_name"`
			UserEmail     string                 `json:"user_email"`
			Status        entities.OrderStatus   `json:"status"`
			PaymentStatus entities.PaymentStatus `json:"payment_status"`
			Total         float64                `json:"total"`
			ItemCount     int                    `json:"item_count"`
			CreatedAt     time.Time              `json:"created_at"`
			UpdatedAt     time.Time              `json:"updated_at"`
		}{
			ID:            order.ID,
			OrderNumber:   order.OrderNumber,
			UserID:        order.UserID,
			UserName:      userName,
			UserEmail:     userEmail,
			Status:        order.Status,
			PaymentStatus: order.PaymentStatus,
			Total:         order.Total,
			ItemCount:     len(order.Items),
			CreatedAt:     order.CreatedAt,
			UpdatedAt:     order.UpdatedAt,
		}
	}

	response := &AdminOrdersResponse{
		Orders:     orderResponses,
		Total:      int64(totalCount),
		Pagination: NewPaginationInfo(req.Offset, req.Limit, int64(totalCount)),
	}

	return response, nil
}

// GetSystemStats gets system statistics
func (uc *adminUseCase) GetSystemStats(ctx context.Context) (*SystemStatsResponse, error) {
	// Mock implementation for system stats
	response := &SystemStatsResponse{
		Database: struct {
			TotalSize       string `json:"total_size"`
			TableCount      int    `json:"table_count"`
			ConnectionCount int    `json:"connection_count"`
			QueryCount      int64  `json:"query_count"`
		}{
			TotalSize:       "2.5 GB",
			TableCount:      25,
			ConnectionCount: 10,
			QueryCount:      1250000,
		},
		Server: struct {
			Uptime       string  `json:"uptime"`
			CPUUsage     float64 `json:"cpu_usage"`
			MemoryUsage  float64 `json:"memory_usage"`
			DiskUsage    float64 `json:"disk_usage"`
			RequestCount int64   `json:"request_count"`
			ErrorRate    float64 `json:"error_rate"`
		}{
			Uptime:       "15 days, 8 hours",
			CPUUsage:     45.2,
			MemoryUsage:  65.5,
			DiskUsage:    78.9,
			RequestCount: 1250000,
			ErrorRate:    0.02,
		},
		Cache: struct {
			HitRate     float64 `json:"hit_rate"`
			MissRate    float64 `json:"miss_rate"`
			KeyCount    int64   `json:"key_count"`
			MemoryUsage string  `json:"memory_usage"`
		}{
			HitRate:     92.5,
			MissRate:    7.5,
			KeyCount:    50000,
			MemoryUsage: "256 MB",
		},
	}

	return response, nil
}

// ManageReviews manages reviews
func (uc *adminUseCase) ManageReviews(ctx context.Context, req ManageReviewsRequest) (*ManageReviewsResponse, error) {
	// Mock implementation for manage reviews
	response := &ManageReviewsResponse{
		Reviews: []struct {
			ID           uuid.UUID             `json:"id"`
			ProductID    uuid.UUID             `json:"product_id"`
			ProductName  string                `json:"product_name"`
			UserID       uuid.UUID             `json:"user_id"`
			UserName     string                `json:"user_name"`
			Rating       int                   `json:"rating"`
			Title        string                `json:"title"`
			Content      string                `json:"content"`
			Status       entities.ReviewStatus `json:"status"`
			HelpfulVotes int                   `json:"helpful_votes"`
			TotalVotes   int                   `json:"total_votes"`
			IsFlagged    bool                  `json:"is_flagged"`
			CreatedAt    time.Time             `json:"created_at"`
		}{
			{
				ID:           uuid.New(),
				ProductID:    uuid.New(),
				ProductName:  "iPhone 15",
				UserID:       uuid.New(),
				UserName:     "John Doe",
				Rating:       5,
				Title:        "Great product!",
				Content:      "Really satisfied with this purchase",
				Status:       entities.ReviewStatusPending,
				HelpfulVotes: 10,
				TotalVotes:   12,
				IsFlagged:    false,
				CreatedAt:    time.Now().Add(-2 * time.Hour),
			},
		},
		Total: 25,
	}

	return response, nil
}

// UpdateReviewStatus updates review status
func (uc *adminUseCase) UpdateReviewStatus(ctx context.Context, reviewID uuid.UUID, status entities.ReviewStatus) error {
	// Mock implementation for update review status
	// In real implementation, this would update the review status in database
	return nil
}

// ProcessRefund processes a refund
func (uc *adminUseCase) ProcessRefund(ctx context.Context, orderID uuid.UUID, amount float64, reason string) error {
	// Mock implementation for process refund
	// In real implementation, this would process the refund through payment service
	return nil
}

// GetReports gets reports
func (uc *adminUseCase) GetReports(ctx context.Context, req GetReportsRequest) (*ReportsListResponse, error) {
	// Mock implementation for get reports
	now := time.Now()
	response := &ReportsListResponse{
		Reports: []*ReportResponse{
			{
				ReportID:    uuid.New(),
				Type:        "sales",
				Format:      "pdf",
				Status:      "completed",
				Progress:    100,
				DownloadURL: "/api/v1/admin/reports/" + uuid.New().String() + "/download",
				CreatedBy:   uuid.New(),
				CreatedAt:   now.Add(-1 * time.Hour),
				CompletedAt: &now,
			},
			{
				ReportID:    uuid.New(),
				Type:        "inventory",
				Format:      "excel",
				Status:      "pending",
				Progress:    45,
				CreatedBy:   uuid.New(),
				CreatedAt:   now.Add(-30 * time.Minute),
				CompletedAt: nil,
			},
		},
		Total: 10,
	}

	return response, nil
}

// GetSystemLogs gets system logs
func (uc *adminUseCase) GetSystemLogs(ctx context.Context, req SystemLogsRequest) (*SystemLogsResponse, error) {
	// Mock implementation for system logs
	response := &SystemLogsResponse{
		Logs: []struct {
			ID        uuid.UUID `json:"id"`
			Level     string    `json:"level"`
			Service   string    `json:"service"`
			Message   string    `json:"message"`
			Context   string    `json:"context"`
			Timestamp time.Time `json:"timestamp"`
		}{
			{
				ID:        uuid.New(),
				Level:     "info",
				Service:   "api",
				Message:   "User login successful",
				Context:   "User ID: " + uuid.New().String(),
				Timestamp: time.Now().Add(-30 * time.Minute),
			},
			{
				ID:        uuid.New(),
				Level:     "error",
				Service:   "database",
				Message:   "Connection timeout",
				Context:   "Timeout after 30 seconds",
				Timestamp: time.Now().Add(-1 * time.Hour),
			},
		},
		Total: 500,
	}

	return response, nil
}

// GetProductAnalytics gets product analytics
func (uc *adminUseCase) GetProductAnalytics(ctx context.Context, productID uuid.UUID, period string) (*ProductAnalyticsResponse, error) {
	// Mock implementation for product analytics
	response := &ProductAnalyticsResponse{
		ProductID: productID,
		Period:    period,
		Metrics: struct {
			Views          int64   `json:"views"`
			Sales          int64   `json:"sales"`
			Revenue        float64 `json:"revenue"`
			ConversionRate float64 `json:"conversion_rate"`
			AverageRating  float64 `json:"average_rating"`
			ReviewCount    int64   `json:"review_count"`
		}{
			Views:          15000,
			Sales:          500,
			Revenue:        250000,
			ConversionRate: 3.3,
			AverageRating:  4.5,
			ReviewCount:    125,
		},
		Charts: struct {
			ViewsChart []struct {
				Date  string `json:"date"`
				Views int64  `json:"views"`
			} `json:"views_chart"`
			SalesChart []struct {
				Date  string `json:"date"`
				Sales int64  `json:"sales"`
			} `json:"sales_chart"`
			RevenueChart []struct {
				Date    string  `json:"date"`
				Revenue float64 `json:"revenue"`
			} `json:"revenue_chart"`
		}{
			ViewsChart: []struct {
				Date  string `json:"date"`
				Views int64  `json:"views"`
			}{
				{Date: "2024-01-01", Views: 1500},
				{Date: "2024-01-02", Views: 1800},
			},
			SalesChart: []struct {
				Date  string `json:"date"`
				Sales int64  `json:"sales"`
			}{
				{Date: "2024-01-01", Sales: 50},
				{Date: "2024-01-02", Sales: 65},
			},
			RevenueChart: []struct {
				Date    string  `json:"date"`
				Revenue float64 `json:"revenue"`
			}{
				{Date: "2024-01-01", Revenue: 25000},
				{Date: "2024-01-02", Revenue: 32500},
			},
		},
	}

	return response, nil
}

// GetUsers gets users for admin
func (uc *adminUseCase) GetUsers(ctx context.Context, req AdminUsersRequest) (*AdminUsersResponse, error) {
	// Get total count first
	total, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users with pagination
	userEntities, err := uc.userRepo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	// Transform entities to response format
	users := make([]struct {
		ID         uuid.UUID           `json:"id"`
		Email      string              `json:"email"`
		FirstName  string              `json:"first_name"`
		LastName   string              `json:"last_name"`
		Role       entities.UserRole   `json:"role"`
		Status     entities.UserStatus `json:"status"`
		LastLogin  *time.Time          `json:"last_login"`
		OrderCount int64               `json:"order_count"`
		TotalSpent float64             `json:"total_spent"`
		CreatedAt  time.Time           `json:"created_at"`
	}, len(userEntities))

	for i, user := range userEntities {
		// Convert user status to UserStatus enum
		var status entities.UserStatus
		if user.IsActive {
			status = entities.UserStatusActive
		} else {
			status = entities.UserStatusInactive
		}

		users[i] = struct {
			ID         uuid.UUID           `json:"id"`
			Email      string              `json:"email"`
			FirstName  string              `json:"first_name"`
			LastName   string              `json:"last_name"`
			Role       entities.UserRole   `json:"role"`
			Status     entities.UserStatus `json:"status"`
			LastLogin  *time.Time          `json:"last_login"`
			OrderCount int64               `json:"order_count"`
			TotalSpent float64             `json:"total_spent"`
			CreatedAt  time.Time           `json:"created_at"`
		}{
			ID:         user.ID,
			Email:      user.Email,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Role:       user.Role,
			Status:     status,
			LastLogin:  nil, // TODO: Implement last login tracking
			OrderCount: 0,   // TODO: Get actual order count from order repository
			TotalSpent: 0,   // TODO: Get actual total spent from order repository
			CreatedAt:  user.CreatedAt,
		}
	}

	pagination := NewPaginationInfo(req.Offset, req.Limit, total)

	response := &AdminUsersResponse{
		Users:      users,
		Total:      total,
		Pagination: pagination,
	}

	return response, nil
}

// UpdateUserStatus updates user status
func (uc *adminUseCase) UpdateUserStatus(ctx context.Context, userID uuid.UUID, status entities.UserStatus) error {
	// Mock implementation for update user status
	// In real implementation, this would update the user status in database
	return nil
}

// UpdateUserRole updates user role
func (uc *adminUseCase) UpdateUserRole(ctx context.Context, userID uuid.UUID, role entities.UserRole) error {
	// Mock implementation for update user role
	// In real implementation, this would update the user role in database
	return nil
}

// GetUserActivity gets user activity
func (uc *adminUseCase) GetUserActivity(ctx context.Context, userID uuid.UUID, req ActivityRequest) (*ActivityResponse, error) {
	// Mock implementation for user activity
	activities := []struct {
		ID          uuid.UUID `json:"id"`
		Type        string    `json:"type"`
		Description string    `json:"description"`
		IPAddress   string    `json:"ip_address"`
		UserAgent   string    `json:"user_agent"`
		Metadata    string    `json:"metadata"`
		CreatedAt   time.Time `json:"created_at"`
	}{
		{
			ID:          uuid.New(),
			Type:        "login",
			Description: "User logged in",
			IPAddress:   "192.168.1.1",
			UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			Metadata:    `{"browser": "Chrome", "os": "Windows"}`,
			CreatedAt:   time.Now().AddDate(0, 0, -1),
		},
		{
			ID:          uuid.New(),
			Type:        "order",
			Description: "Order placed: #ORD-001",
			IPAddress:   "192.168.1.1",
			UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			Metadata:    `{"order_id": "ORD-001", "total": 99.99}`,
			CreatedAt:   time.Now().AddDate(0, 0, -2),
		},
	}

	total := int64(len(activities))
	pagination := NewPaginationInfo(req.Offset, req.Limit, total)

	response := &ActivityResponse{
		Activities: activities,
		Total:      total,
		Pagination: pagination,
	}

	return response, nil
}

// UpdateOrderStatus updates order status
func (uc *adminUseCase) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status entities.OrderStatus) error {
	// Check if order exists
	_, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Update order status in database
	return uc.orderRepo.UpdateStatus(ctx, orderID, status)
}

// GetProducts gets products for admin
func (uc *adminUseCase) GetProducts(ctx context.Context, req AdminProductsRequest) (*AdminProductsResponse, error) {
	// Mock implementation for admin products
	products := []struct {
		ID            uuid.UUID              `json:"id"`
		Name          string                 `json:"name"`
		SKU           string                 `json:"sku"`
		Price         float64                `json:"price"`
		ComparePrice  float64                `json:"compare_price"`
		Status        entities.ProductStatus `json:"status"`
		StockQuantity int                    `json:"stock_quantity"`
		CategoryID    uuid.UUID              `json:"category_id"`
		CategoryName  string                 `json:"category_name"`
		ViewCount     int64                  `json:"view_count"`
		SalesCount    int64                  `json:"sales_count"`
		Revenue       float64                `json:"revenue"`
		CreatedAt     time.Time              `json:"created_at"`
		UpdatedAt     time.Time              `json:"updated_at"`
	}{
		{
			ID:            uuid.New(),
			Name:          "iPhone 15",
			SKU:           "IPHONE15-001",
			Price:         999.99,
			ComparePrice:  1099.99,
			Status:        entities.ProductStatusActive,
			StockQuantity: 50,
			CategoryID:    uuid.New(),
			CategoryName:  "Electronics",
			ViewCount:     15000,
			SalesCount:    500,
			Revenue:       499950,
			CreatedAt:     time.Now().AddDate(0, -1, 0),
			UpdatedAt:     time.Now(),
		},
		{
			ID:            uuid.New(),
			Name:          "MacBook Pro",
			SKU:           "MBP-001",
			Price:         1999.99,
			ComparePrice:  2199.99,
			Status:        entities.ProductStatusActive,
			StockQuantity: 25,
			CategoryID:    uuid.New(),
			CategoryName:  "Computers",
			ViewCount:     12000,
			SalesCount:    300,
			Revenue:       599997,
			CreatedAt:     time.Now().AddDate(0, -2, 0),
			UpdatedAt:     time.Now(),
		},
	}

	total := int64(len(products))
	pagination := NewPaginationInfo(req.Offset, req.Limit, total)

	response := &AdminProductsResponse{
		Products:   products,
		Total:      total,
		Pagination: pagination,
	}

	return response, nil
}
