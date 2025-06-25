package entities

import (
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of analytics event
type EventType string

const (
	EventTypePageView      EventType = "page_view"
	EventTypeProductView   EventType = "product_view"
	EventTypeAddToCart     EventType = "add_to_cart"
	EventTypeRemoveFromCart EventType = "remove_from_cart"
	EventTypeCheckout      EventType = "checkout"
	EventTypePurchase      EventType = "purchase"
	EventTypeSearch        EventType = "search"
	EventTypeLogin         EventType = "login"
	EventTypeLogout        EventType = "logout"
	EventTypeRegister      EventType = "register"
	EventTypeWishlistAdd   EventType = "wishlist_add"
	EventTypeWishlistRemove EventType = "wishlist_remove"
	EventTypeReview        EventType = "review"
	EventTypeShare         EventType = "share"
	EventTypeClick         EventType = "click"
	EventTypeCustom        EventType = "custom"
)

// AnalyticsEvent represents user behavior tracking events
type AnalyticsEvent struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      *uuid.UUID `json:"user_id" gorm:"type:uuid;index"`        // null for anonymous users
	User        *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	SessionID   string    `json:"session_id" gorm:"index;not null"`
	
	// Event details
	EventType   EventType `json:"event_type" gorm:"not null;index"`
	EventName   string    `json:"event_name" gorm:"not null"`
	Category    string    `json:"category"`
	Action      string    `json:"action"`
	Label       string    `json:"label"`
	Value       float64   `json:"value" gorm:"default:0"`
	
	// Context information
	Page        string    `json:"page"`
	Referrer    string    `json:"referrer"`
	UserAgent   string    `json:"user_agent"`
	IPAddress   string    `json:"ip_address"`
	Country     string    `json:"country"`
	City        string    `json:"city"`
	Device      string    `json:"device"`
	Browser     string    `json:"browser"`
	OS          string    `json:"os"`
	
	// E-commerce specific
	ProductID   *uuid.UUID `json:"product_id" gorm:"type:uuid;index"`
	CategoryID  *uuid.UUID `json:"category_id" gorm:"type:uuid;index"`
	OrderID     *uuid.UUID `json:"order_id" gorm:"type:uuid;index"`
	
	// Additional data
	Properties  string    `json:"properties" gorm:"type:text"`           // JSON object with custom properties
	Revenue     float64   `json:"revenue" gorm:"default:0"`              // For purchase events
	Quantity    int       `json:"quantity" gorm:"default:0"`             // For cart/purchase events
	
	// Metadata
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;index"`
}

// TableName returns the table name for AnalyticsEvent entity
func (AnalyticsEvent) TableName() string {
	return "analytics_events"
}

// IsEcommerceEvent checks if event is e-commerce related
func (ae *AnalyticsEvent) IsEcommerceEvent() bool {
	return ae.EventType == EventTypeAddToCart ||
		   ae.EventType == EventTypeRemoveFromCart ||
		   ae.EventType == EventTypeCheckout ||
		   ae.EventType == EventTypePurchase ||
		   ae.EventType == EventTypeProductView
}

// IsConversionEvent checks if event represents a conversion
func (ae *AnalyticsEvent) IsConversionEvent() bool {
	return ae.EventType == EventTypePurchase ||
		   ae.EventType == EventTypeRegister ||
		   ae.EventType == EventTypeCheckout
}

// SalesReport represents sales analytics data
type SalesReport struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ReportDate      time.Time `json:"report_date" gorm:"not null;uniqueIndex:idx_sales_report_date"`
	Period          string    `json:"period" gorm:"not null;index"`        // daily, weekly, monthly, yearly
	
	// Sales metrics
	TotalOrders     int       `json:"total_orders" gorm:"default:0"`
	TotalRevenue    float64   `json:"total_revenue" gorm:"default:0"`
	TotalItems      int       `json:"total_items" gorm:"default:0"`
	AverageOrderValue float64 `json:"average_order_value" gorm:"default:0"`
	
	// Order status breakdown
	PendingOrders   int       `json:"pending_orders" gorm:"default:0"`
	ConfirmedOrders int       `json:"confirmed_orders" gorm:"default:0"`
	ShippedOrders   int       `json:"shipped_orders" gorm:"default:0"`
	DeliveredOrders int       `json:"delivered_orders" gorm:"default:0"`
	CancelledOrders int       `json:"cancelled_orders" gorm:"default:0"`
	RefundedOrders  int       `json:"refunded_orders" gorm:"default:0"`
	
	// Payment metrics
	PaidRevenue     float64   `json:"paid_revenue" gorm:"default:0"`
	PendingRevenue  float64   `json:"pending_revenue" gorm:"default:0"`
	RefundedRevenue float64   `json:"refunded_revenue" gorm:"default:0"`
	
	// Customer metrics
	NewCustomers    int       `json:"new_customers" gorm:"default:0"`
	ReturningCustomers int    `json:"returning_customers" gorm:"default:0"`
	
	// Product metrics
	TopSellingProducts string `json:"top_selling_products" gorm:"type:text"` // JSON array
	TopCategories     string  `json:"top_categories" gorm:"type:text"`       // JSON array
	
	// Metadata
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for SalesReport entity
func (SalesReport) TableName() string {
	return "sales_reports"
}

// GetConversionRate calculates conversion rate
func (sr *SalesReport) GetConversionRate(totalVisitors int) float64 {
	if totalVisitors == 0 {
		return 0
	}
	return float64(sr.TotalOrders) / float64(totalVisitors) * 100
}

// GetCancellationRate calculates order cancellation rate
func (sr *SalesReport) GetCancellationRate() float64 {
	if sr.TotalOrders == 0 {
		return 0
	}
	return float64(sr.CancelledOrders) / float64(sr.TotalOrders) * 100
}

// ProductAnalytics represents product performance analytics
type ProductAnalytics struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID       uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	ReportDate      time.Time `json:"report_date" gorm:"not null;index"`
	Period          string    `json:"period" gorm:"not null;index"`        // daily, weekly, monthly
	
	// View metrics
	PageViews       int       `json:"page_views" gorm:"default:0"`
	UniqueViews     int       `json:"unique_views" gorm:"default:0"`
	ViewDuration    float64   `json:"view_duration" gorm:"default:0"`      // Average time in seconds
	
	// Engagement metrics
	AddToCartCount  int       `json:"add_to_cart_count" gorm:"default:0"`
	WishlistCount   int       `json:"wishlist_count" gorm:"default:0"`
	ShareCount      int       `json:"share_count" gorm:"default:0"`
	ReviewCount     int       `json:"review_count" gorm:"default:0"`
	
	// Sales metrics
	OrderCount      int       `json:"order_count" gorm:"default:0"`
	QuantitySold    int       `json:"quantity_sold" gorm:"default:0"`
	Revenue         float64   `json:"revenue" gorm:"default:0"`
	
	// Conversion metrics
	ViewToCartRate  float64   `json:"view_to_cart_rate" gorm:"default:0"`  // Percentage
	CartToOrderRate float64   `json:"cart_to_order_rate" gorm:"default:0"` // Percentage
	ConversionRate  float64   `json:"conversion_rate" gorm:"default:0"`    // View to purchase
	
	// Search metrics
	SearchImpressions int     `json:"search_impressions" gorm:"default:0"`
	SearchClicks     int      `json:"search_clicks" gorm:"default:0"`
	SearchPosition   float64  `json:"search_position" gorm:"default:0"`    // Average position in search results
	
	// Metadata
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for ProductAnalytics entity
func (ProductAnalytics) TableName() string {
	return "product_analytics"
}

// CalculateConversionRate calculates view to purchase conversion rate
func (pa *ProductAnalytics) CalculateConversionRate() {
	if pa.UniqueViews > 0 {
		pa.ConversionRate = float64(pa.OrderCount) / float64(pa.UniqueViews) * 100
	}
}

// CalculateViewToCartRate calculates view to cart conversion rate
func (pa *ProductAnalytics) CalculateViewToCartRate() {
	if pa.UniqueViews > 0 {
		pa.ViewToCartRate = float64(pa.AddToCartCount) / float64(pa.UniqueViews) * 100
	}
}

// UserAnalytics represents user behavior analytics
type UserAnalytics struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	ReportDate      time.Time `json:"report_date" gorm:"not null;index"`
	Period          string    `json:"period" gorm:"not null;index"`        // daily, weekly, monthly
	
	// Session metrics
	SessionCount    int       `json:"session_count" gorm:"default:0"`
	TotalDuration   float64   `json:"total_duration" gorm:"default:0"`     // Total time in seconds
	AverageDuration float64   `json:"average_duration" gorm:"default:0"`   // Average session duration
	PageViews       int       `json:"page_views" gorm:"default:0"`
	
	// E-commerce metrics
	ProductViews    int       `json:"product_views" gorm:"default:0"`
	CartAdditions   int       `json:"cart_additions" gorm:"default:0"`
	OrdersPlaced    int       `json:"orders_placed" gorm:"default:0"`
	TotalSpent      float64   `json:"total_spent" gorm:"default:0"`
	
	// Engagement metrics
	ReviewsWritten  int       `json:"reviews_written" gorm:"default:0"`
	WishlistItems   int       `json:"wishlist_items" gorm:"default:0"`
	ShareActions    int       `json:"share_actions" gorm:"default:0"`
	
	// Device/Browser info
	PrimaryDevice   string    `json:"primary_device"`
	PrimaryBrowser  string    `json:"primary_browser"`
	PrimaryOS       string    `json:"primary_os"`
	
	// Metadata
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for UserAnalytics entity
func (UserAnalytics) TableName() string {
	return "user_analytics"
}

// GetEngagementScore calculates user engagement score
func (ua *UserAnalytics) GetEngagementScore() float64 {
	// Simple engagement score based on various activities
	score := float64(ua.PageViews)*0.1 +
		float64(ua.ProductViews)*0.2 +
		float64(ua.CartAdditions)*0.5 +
		float64(ua.OrdersPlaced)*2.0 +
		float64(ua.ReviewsWritten)*1.0 +
		float64(ua.WishlistItems)*0.3
	
	return score
}

// SalesMetrics represents sales analytics metrics
type SalesMetrics struct {
	TotalSales        float64 `json:"total_sales"`
	TotalOrders       int64   `json:"total_orders"`
	AverageOrderValue float64 `json:"average_order_value"`
}

// ProductMetrics represents product performance metrics
type ProductMetrics struct {
	UnitsSold   int64   `json:"units_sold"`
	Revenue     float64 `json:"revenue"`
	ViewCount   int64   `json:"view_count"`
}

// UserMetrics represents user analytics metrics
type UserMetrics struct {
	NewUsers    int64 `json:"new_users"`
	ActiveUsers int64 `json:"active_users"`
	TotalUsers  int64 `json:"total_users"`
}

// TrafficMetrics represents website traffic metrics
type TrafficMetrics struct {
	PageViews      int64   `json:"page_views"`
	UniqueVisitors int64   `json:"unique_visitors"`
	BounceRate     float64 `json:"bounce_rate"`
}

// TopProduct represents top selling product data
type TopProduct struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	UnitsSold int64     `json:"units_sold"`
	Revenue   float64   `json:"revenue"`
}

// TopCategory represents top performing category data
type TopCategory struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	UnitsSold int64     `json:"units_sold"`
	Revenue   float64   `json:"revenue"`
}

// RevenueData represents revenue data for a period
type RevenueData struct {
	Period     string  `json:"period"`
	Revenue    float64 `json:"revenue"`
	OrderCount int64   `json:"order_count"`
}

// ConversionMetrics represents conversion rate metrics
type ConversionMetrics struct {
	TotalSessions     int64   `json:"total_sessions"`
	ConvertedSessions int64   `json:"converted_sessions"`
	ConversionRate    float64 `json:"conversion_rate"`
}

// RealTimeMetrics represents real-time metrics
type RealTimeMetrics struct {
	ActiveUsers int64   `json:"active_users"`
	PageViews   int64   `json:"page_views"`
	Orders      int64   `json:"orders"`
	Revenue     float64 `json:"revenue"`
}

// CustomerLifetimeValue represents customer lifetime value metrics
type CustomerLifetimeValue struct {
	TotalSpent       float64    `json:"total_spent"`
	OrderCount       int64      `json:"order_count"`
	AverageOrderValue float64   `json:"average_order_value"`
	FirstOrderDate   *time.Time `json:"first_order_date"`
	LastOrderDate    *time.Time `json:"last_order_date"`
}

// InventoryMetrics represents inventory analytics
type InventoryMetrics struct {
	TotalProducts        int64   `json:"total_products"`
	LowStockItems        int64   `json:"low_stock_items"`
	OutOfStockItems      int64   `json:"out_of_stock_items"`
	TotalInventoryValue  float64 `json:"total_inventory_value"`
}

// RatingDistribution represents rating distribution data
type RatingDistribution struct {
	Rating1Count  int64   `json:"rating_1_count"`
	Rating2Count  int64   `json:"rating_2_count"`
	Rating3Count  int64   `json:"rating_3_count"`
	Rating4Count  int64   `json:"rating_4_count"`
	Rating5Count  int64   `json:"rating_5_count"`
	TotalReviews  int64   `json:"total_reviews"`
	AverageRating float64 `json:"average_rating"`
}

// CategoryAnalytics represents category performance analytics
type CategoryAnalytics struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CategoryID      uuid.UUID `json:"category_id" gorm:"type:uuid;not null;index"`
	ReportDate      time.Time `json:"report_date" gorm:"not null;index"`
	Period          string    `json:"period" gorm:"not null;index"`        // daily, weekly, monthly
	
	// View metrics
	PageViews       int       `json:"page_views" gorm:"default:0"`
	UniqueViews     int       `json:"unique_views" gorm:"default:0"`
	ProductViews    int       `json:"product_views" gorm:"default:0"`      // Views of products in this category
	
	// Sales metrics
	OrderCount      int       `json:"order_count" gorm:"default:0"`
	ItemsSold       int       `json:"items_sold" gorm:"default:0"`
	Revenue         float64   `json:"revenue" gorm:"default:0"`
	
	// Product metrics
	ActiveProducts  int       `json:"active_products" gorm:"default:0"`
	OutOfStockProducts int    `json:"out_of_stock_products" gorm:"default:0"`
	
	// Engagement metrics
	AddToCartCount  int       `json:"add_to_cart_count" gorm:"default:0"`
	WishlistCount   int       `json:"wishlist_count" gorm:"default:0"`
	
	// Metadata
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for CategoryAnalytics entity
func (CategoryAnalytics) TableName() string {
	return "category_analytics"
}

// GetConversionRate calculates category conversion rate
func (ca *CategoryAnalytics) GetConversionRate() float64 {
	if ca.UniqueViews == 0 {
		return 0
	}
	return float64(ca.OrderCount) / float64(ca.UniqueViews) * 100
}

// SearchAnalytics represents search behavior analytics
type SearchAnalytics struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SearchTerm      string    `json:"search_term" gorm:"not null;index"`
	ReportDate      time.Time `json:"report_date" gorm:"not null;index"`
	Period          string    `json:"period" gorm:"not null;index"`        // daily, weekly, monthly
	
	// Search metrics
	SearchCount     int       `json:"search_count" gorm:"default:0"`
	UniqueSearchers int       `json:"unique_searchers" gorm:"default:0"`
	ResultsCount    int       `json:"results_count" gorm:"default:0"`      // Average results returned
	
	// Engagement metrics
	ClickCount      int       `json:"click_count" gorm:"default:0"`        // Clicks on search results
	ClickThroughRate float64  `json:"click_through_rate" gorm:"default:0"` // CTR percentage
	
	// Conversion metrics
	ConversionsCount int      `json:"conversions_count" gorm:"default:0"`  // Purchases from search
	ConversionRate  float64   `json:"conversion_rate" gorm:"default:0"`    // Conversion percentage
	Revenue         float64   `json:"revenue" gorm:"default:0"`            // Revenue from search
	
	// Result quality
	NoResultsCount  int       `json:"no_results_count" gorm:"default:0"`   // Searches with no results
	RefinementCount int       `json:"refinement_count" gorm:"default:0"`   // Search refinements
	
	// Metadata
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for SearchAnalytics entity
func (SearchAnalytics) TableName() string {
	return "search_analytics"
}

// CalculateClickThroughRate calculates CTR
func (sa *SearchAnalytics) CalculateClickThroughRate() {
	if sa.SearchCount > 0 {
		sa.ClickThroughRate = float64(sa.ClickCount) / float64(sa.SearchCount) * 100
	}
}

// CalculateConversionRate calculates search conversion rate
func (sa *SearchAnalytics) CalculateConversionRate() {
	if sa.SearchCount > 0 {
		sa.ConversionRate = float64(sa.ConversionsCount) / float64(sa.SearchCount) * 100
	}
}
