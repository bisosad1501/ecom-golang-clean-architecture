package repositories

import (
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// WishlistFilters represents filters for wishlist queries
type WishlistFilters struct {
	UserID        *uuid.UUID `json:"user_id"`
	ProductID     *uuid.UUID `json:"product_id"`
	CreatedAfter  *time.Time `json:"created_after"`
	CreatedBefore *time.Time `json:"created_before"`
	SortBy        string     `json:"sort_by"`    // created_at, product_name
	SortOrder     string     `json:"sort_order"` // asc, desc
	Limit         int        `json:"limit"`
	Offset        int        `json:"offset"`
}

// NotificationFilters represents filters for notification queries
type NotificationFilters struct {
	UserID        *uuid.UUID                    `json:"user_id"`
	Type          *entities.NotificationType    `json:"type"`
	IsRead        *bool                         `json:"is_read"`
	Priority      *entities.NotificationPriority `json:"priority"`
	DateFrom      *time.Time                    `json:"date_from"`
	DateTo        *time.Time                    `json:"date_to"`
	CreatedAfter  *time.Time                    `json:"created_after"`
	CreatedBefore *time.Time                    `json:"created_before"`
	SortBy        string                        `json:"sort_by"`    // created_at, type
	SortOrder     string                        `json:"sort_order"` // asc, desc
	Limit         int                           `json:"limit"`
	Offset        int                           `json:"offset"`
}

// ReviewVoteFilters represents filters for review vote queries
type ReviewVoteFilters struct {
	UserID        *uuid.UUID                   `json:"user_id"`
	ReviewID      *uuid.UUID                   `json:"review_id"`
	VoteType      *entities.ReviewVoteType     `json:"vote_type"`
	CreatedAfter  *time.Time                   `json:"created_after"`
	CreatedBefore *time.Time                   `json:"created_before"`
	SortBy        string                       `json:"sort_by"`    // created_at, vote_type
	SortOrder     string                       `json:"sort_order"` // asc, desc
	Limit         int                          `json:"limit"`
	Offset        int                          `json:"offset"`
}

// ProductRatingFilters represents filters for product rating queries
type ProductRatingFilters struct {
	ProductID  *uuid.UUID `json:"product_id"`
	MinRating  *float64   `json:"min_rating"`
	MaxRating  *float64   `json:"max_rating"`
	MinReviews *int64     `json:"min_reviews"`
	SortBy     string     `json:"sort_by"`    // average_rating, total_reviews, updated_at
	SortOrder  string     `json:"sort_order"` // asc, desc
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset"`
}

// ShipmentFilters represents filters for shipment queries
type ShipmentFilters struct {
	OrderID        *uuid.UUID                 `json:"order_id"`
	Status         *entities.ShipmentStatus   `json:"status"`
	Carrier        string                     `json:"carrier"`
	TrackingNumber string                     `json:"tracking_number"`
	CreatedAfter   *time.Time                 `json:"created_after"`
	CreatedBefore  *time.Time                 `json:"created_before"`
	SortBy         string                     `json:"sort_by"`    // created_at, status, carrier
	SortOrder      string                     `json:"sort_order"` // asc, desc
	Limit          int                        `json:"limit"`
	Offset         int                        `json:"offset"`
}

// AuditFilters represents filters for audit log queries
type AuditFilters struct {
	UserID        *uuid.UUID `json:"user_id"`
	Action        string     `json:"action"`
	Resource      string     `json:"resource"`
	ResourceID    *uuid.UUID `json:"resource_id"`
	IPAddress     string     `json:"ip_address"`
	UserAgent     string     `json:"user_agent"`
	CreatedAfter  *time.Time `json:"created_after"`
	CreatedBefore *time.Time `json:"created_before"`
	SortBy        string     `json:"sort_by"`    // created_at, action, resource
	SortOrder     string     `json:"sort_order"` // asc, desc
	Limit         int        `json:"limit"`
	Offset        int        `json:"offset"`
}



// EventFilters represents filters for analytics event queries
type EventFilters struct {
	EventType string     `json:"event_type"`
	UserID    *uuid.UUID `json:"user_id"`
	ProductID *uuid.UUID `json:"product_id"`
	SessionID string     `json:"session_id"`
	DateFrom  *time.Time `json:"date_from"`
	DateTo    *time.Time `json:"date_to"`
	SortBy    string     `json:"sort_by"`    // created_at, event_type
	SortOrder string     `json:"sort_order"` // asc, desc
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

// ActivitySummary represents activity summary data
type ActivitySummary struct {
	TotalEvents      int64 `json:"total_events"`
	UniqueUsers      int64 `json:"unique_users"`
	FailedLogins     int64 `json:"failed_logins"`
	SuccessfulLogins int64 `json:"successful_logins"`
}

// AdminActionFilters represents filters for admin action queries
type AdminActionFilters struct {
	AdminID   *uuid.UUID `json:"admin_id"`
	Action    string     `json:"action"`
	Resource  string     `json:"resource"`
	DateFrom  *time.Time `json:"date_from"`
	DateTo    *time.Time `json:"date_to"`
	SortBy    string     `json:"sort_by"`    // created_at, action
	SortOrder string     `json:"sort_order"` // asc, desc
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

// ComplianceReport represents compliance audit report
type ComplianceReport struct {
	Period            string `json:"period"`
	TotalEvents       int64  `json:"total_events"`
	SecurityEvents    int64  `json:"security_events"`
	FailedLogins      int64  `json:"failed_logins"`
	DataAccessEvents  int64  `json:"data_access_events"`
}

// DashboardMetrics represents dashboard metrics
type DashboardMetrics struct {
	TotalUsers     int64   `json:"total_users"`
	TotalOrders    int64   `json:"total_orders"`
	TotalRevenue   float64 `json:"total_revenue"`
	ConversionRate float64 `json:"conversion_rate"`
}

// FunnelAnalysis represents funnel analysis data
type FunnelAnalysis struct {
	Steps          []string `json:"steps"`
	TotalUsers     int64    `json:"total_users"`
	ConversionRate float64  `json:"conversion_rate"`
}

// DeliveryStats represents notification delivery statistics
type DeliveryStats struct {
	TotalSent    int64   `json:"total_sent"`
	Delivered    int64   `json:"delivered"`
	Failed       int64   `json:"failed"`
	DeliveryRate float64 `json:"delivery_rate"`
}

// ProductMetricsFilters represents filters for product metrics
type ProductMetricsFilters struct {
	ProductID *uuid.UUID `json:"product_id"`
	DateFrom  *time.Time `json:"date_from"`
	DateTo    *time.Time `json:"date_to"`
}

// ProductMetrics represents product performance metrics
type ProductMetrics struct {
	UnitsSold int64   `json:"units_sold"`
	Revenue   float64 `json:"revenue"`
	ViewCount int64   `json:"view_count"`
}

// LogRetentionStats represents log retention statistics
type LogRetentionStats struct {
	TotalLogs            int64 `json:"total_logs"`
	LogsOlderThan30Days  int64 `json:"logs_older_than_30_days"`
	LogsOlderThan90Days  int64 `json:"logs_older_than_90_days"`
}

// EngagementStats represents engagement statistics
type EngagementStats struct {
	TotalNotifications  int64   `json:"total_notifications"`
	OpenedNotifications int64   `json:"opened_notifications"`
	OpenRate           float64 `json:"open_rate"`
}

// SalesMetricsFilters represents filters for sales metrics
type SalesMetricsFilters struct {
	DateFrom *time.Time `json:"date_from"`
	DateTo   *time.Time `json:"date_to"`
}

// SalesMetrics represents sales metrics
type SalesMetrics struct {
	TotalSales        float64 `json:"total_sales"`
	TotalOrders       int64   `json:"total_orders"`
	AverageOrderValue float64 `json:"average_order_value"`
}

// SecurityLogFilters represents filters for security logs
type SecurityLogFilters struct {
	DateFrom *time.Time `json:"date_from"`
	DateTo   *time.Time `json:"date_to"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
}

// NotificationStats represents notification statistics
type NotificationStats struct {
	TotalSent int64 `json:"total_sent"`
	TotalRead int64 `json:"total_read"`
}

// TopCategory represents top performing category
type TopCategory struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	UnitsSold int64     `json:"units_sold"`
	Revenue   float64   `json:"revenue"`
}

// SystemLogFilters represents filters for system logs
type SystemLogFilters struct {
	DateFrom *time.Time `json:"date_from"`
	DateTo   *time.Time `json:"date_to"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
}

// TopPage represents top performing page
type TopPage struct {
	Page        string `json:"page"`
	Views       int64  `json:"views"`
	UniqueViews int64  `json:"unique_views"`
}

// TopProduct represents top performing product
type TopProduct struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	UnitsSold int64     `json:"units_sold"`
	Revenue   float64   `json:"revenue"`
}

// WarehousePerformanceReport represents warehouse performance report
type WarehousePerformanceReport struct {
	WarehouseID   uuid.UUID `json:"warehouse_id"`
	WarehouseName string    `json:"warehouse_name"`
	TotalOrders   int64     `json:"total_orders"`
	TotalRevenue  float64   `json:"total_revenue"`
}

// TrafficMetricsFilters represents filters for traffic metrics
type TrafficMetricsFilters struct {
	DateFrom *time.Time `json:"date_from"`
	DateTo   *time.Time `json:"date_to"`
}

// TrafficMetrics represents traffic metrics
type TrafficMetrics struct {
	PageViews      int64   `json:"page_views"`
	UniqueVisitors int64   `json:"unique_visitors"`
	BounceRate     float64 `json:"bounce_rate"`
}

// ReportFilters represents filters for reports
type ReportFilters struct {
	DateFrom *time.Time `json:"date_from"`
	DateTo   *time.Time `json:"date_to"`
}

// PerformanceReport represents performance report
type PerformanceReport struct {
	TotalOrders  int64   `json:"total_orders"`
	TotalRevenue float64 `json:"total_revenue"`
}

// CohortAnalysis represents cohort analysis data
type CohortAnalysis struct {
	Period        string  `json:"period"`
	TotalUsers    int64   `json:"total_users"`
	RetentionRate float64 `json:"retention_rate"`
}

// UserActivitySummary represents user activity summary
type UserActivitySummary struct {
	UserID          uuid.UUID `json:"user_id"`
	TotalActivities int64     `json:"total_activities"`
}

// UserMetricsFilters represents filters for user metrics
type UserMetricsFilters struct {
	DateFrom *time.Time `json:"date_from"`
	DateTo   *time.Time `json:"date_to"`
}

// UserMetrics represents user metrics
type UserMetrics struct {
	NewUsers    int64 `json:"new_users"`
	ActiveUsers int64 `json:"active_users"`
	TotalUsers  int64 `json:"total_users"`
}



// WarehouseCapacity represents warehouse capacity information
type WarehouseCapacity struct {
	TotalCapacity     float64 `json:"total_capacity"`
	UsedCapacity      float64 `json:"used_capacity"`
	AvailableCapacity float64 `json:"available_capacity"`
	CapacityUnit      string  `json:"capacity_unit"`
}

// WarehouseInventoryFilters represents filters for warehouse inventory
type WarehouseInventoryFilters struct {
	ProductID *uuid.UUID `json:"product_id"`
	LowStock  bool       `json:"low_stock"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

// SearchFilters represents filters for search operations
type SearchFilters struct {
	DateFrom *time.Time `json:"date_from"`
	DateTo   *time.Time `json:"date_to"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
}

// WarehouseMetrics represents warehouse metrics
type WarehouseMetrics struct {
	TotalProducts   int64 `json:"total_products"`
	LowStockCount   int64 `json:"low_stock_count"`
	OutOfStockCount int64 `json:"out_of_stock_count"`
}

// MetricsFilters represents filters for metrics
type MetricsFilters struct {
	DateFrom *time.Time `json:"date_from"`
	DateTo   *time.Time `json:"date_to"`
	Period   string     `json:"period"`
}

// WarehouseStaff represents warehouse staff information
type WarehouseStaff struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	AssignedAt time.Time `json:"assigned_at"`
	IsActive   bool      `json:"is_active"`
}

// WarehouseFilters represents filters for warehouse search
type WarehouseFilters struct {
	Name      string `json:"name"`
	Code      string `json:"code"`
	IsActive  *bool  `json:"is_active"`
	Type      string `json:"type"`
	Country   string `json:"country"`
	State     string `json:"state"`
	City      string `json:"city"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

// InventoryFilters represents filters for inventory queries
type InventoryFilters struct {
	ProductID   *uuid.UUID `json:"product_id"`
	WarehouseID *uuid.UUID `json:"warehouse_id"`
	LowStock    bool       `json:"low_stock"`
	OutOfStock  bool       `json:"out_of_stock"`
	SortBy      string     `json:"sort_by"`    // created_at, updated_at, quantity
	SortOrder   string     `json:"sort_order"` // asc, desc
	Limit       int        `json:"limit"`
	Offset      int        `json:"offset"`
}
