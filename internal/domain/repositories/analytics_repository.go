package repositories

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// AnalyticsRepository defines analytics repository interface
type AnalyticsRepository interface {
	// Event tracking
	CreateEvent(ctx context.Context, event *entities.AnalyticsEvent) error
	GetEvents(ctx context.Context, filters EventFilters) ([]*entities.AnalyticsEvent, error)
	CountEvents(ctx context.Context, filters EventFilters) (int64, error)

	// Dashboard metrics
	GetDashboardMetrics(ctx context.Context, dateFrom, dateTo time.Time) (*DashboardMetrics, error)
	GetSalesMetrics(ctx context.Context, filters SalesMetricsFilters) (*SalesMetrics, error)
	GetProductMetrics(ctx context.Context, filters ProductMetricsFilters) (*ProductMetrics, error)
	GetUserMetrics(ctx context.Context, filters UserMetricsFilters) (*UserMetrics, error)
	GetTrafficMetrics(ctx context.Context, filters TrafficMetricsFilters) (*TrafficMetrics, error)

	// Real-time metrics
	GetActiveUsers(ctx context.Context, duration time.Duration) (int64, error)
	GetOnlineVisitors(ctx context.Context) (int64, error)
	GetTodayOrders(ctx context.Context) (int64, error)
	GetTodayRevenue(ctx context.Context) (float64, error)

	// Top performers
	GetTopProducts(ctx context.Context, period string, limit int) ([]*TopProduct, error)
	GetTopCategories(ctx context.Context, period string, limit int) ([]*TopCategory, error)
	GetTopPages(ctx context.Context, period string, limit int) ([]*TopPage, error)

	// Conversion tracking
	GetConversionRate(ctx context.Context, dateFrom, dateTo time.Time) (float64, error)
	GetFunnelAnalysis(ctx context.Context, steps []string, dateFrom, dateTo time.Time) (*FunnelAnalysis, error)

	// Cohort analysis
	GetUserCohorts(ctx context.Context, period string) (*CohortAnalysis, error)
	GetRetentionRate(ctx context.Context, period string) (float64, error)

	// Custom reports
	ExecuteCustomQuery(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error)
}



























// TimeSeriesPoint represents a point in time series data
type TimeSeriesPoint struct {
	Period  string  `json:"period"`
	Revenue float64 `json:"revenue"`
	Orders  int64   `json:"orders"`
	Growth  float64 `json:"growth"`
}

// BreakdownItem represents breakdown data
type BreakdownItem struct {
	Category string  `json:"category"`
	Revenue  float64 `json:"revenue"`
	Orders   int64   `json:"orders"`
	Share    float64 `json:"share"`
}

// ProductMetric represents individual product metrics
type ProductMetric struct {
	ProductID      uuid.UUID `json:"product_id"`
	ProductName    string    `json:"product_name"`
	Views          int64     `json:"views"`
	Sales          int64     `json:"sales"`
	Revenue        float64   `json:"revenue"`
	ConversionRate float64   `json:"conversion_rate"`
	Stock          int       `json:"stock"`
}

// SegmentItem represents user segment data
type SegmentItem struct {
	Segment string  `json:"segment"`
	Count   int64   `json:"count"`
	Share   float64 `json:"share"`
}

// TrafficSource represents traffic source data
type TrafficSource struct {
	Source   string  `json:"source"`
	Visitors int64   `json:"visitors"`
	Share    float64 `json:"share"`
}

// PopularPage represents popular page data
type PopularPage struct {
	Page        string `json:"page"`
	Views       int64  `json:"views"`
	UniqueViews int64  `json:"unique_views"`
}





// Cohort represents a user cohort
type Cohort struct {
	Period        string    `json:"period"`
	Users         int64     `json:"users"`
	RetentionRate []float64 `json:"retention_rate"`
}
