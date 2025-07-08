package database

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type analyticsRepository struct {
	db *gorm.DB
}

// NewAnalyticsRepository creates a new analytics repository
func NewAnalyticsRepository(db *gorm.DB) repositories.AnalyticsRepository {
	return &analyticsRepository{db: db}
}

// RecordEvent records an analytics event
func (r *analyticsRepository) RecordEvent(ctx context.Context, event *entities.AnalyticsEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

// GetSalesMetrics gets sales metrics with filters
func (r *analyticsRepository) GetSalesMetrics(ctx context.Context, filters repositories.SalesMetricsFilters) (*repositories.SalesMetrics, error) {
	var metrics repositories.SalesMetrics

	query := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Select("COALESCE(SUM(total), 0) as total_sales, COUNT(*) as total_orders").
		Where("status = ? AND payment_status = ?", entities.OrderStatusDelivered, entities.PaymentStatusPaid)

	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	err := query.Scan(&metrics).Error
	if err != nil {
		return nil, err
	}

	// Get average order value
	if metrics.TotalOrders > 0 {
		metrics.AverageOrderValue = metrics.TotalSales / float64(metrics.TotalOrders)
	}

	return &metrics, nil
}

// GetProductMetrics gets product performance metrics with filters
func (r *analyticsRepository) GetProductMetrics(ctx context.Context, filters repositories.ProductMetricsFilters) (*repositories.ProductMetrics, error) {
	var metrics repositories.ProductMetrics

	// Get sales data
	query := r.db.WithContext(ctx).
		Table("order_items").
		Select("COALESCE(SUM(quantity), 0) as units_sold, COALESCE(SUM(price * quantity), 0) as revenue").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.status = ?", entities.OrderStatusDelivered)

	if filters.ProductID != nil {
		query = query.Where("order_items.product_id = ?", *filters.ProductID)
	}

	if filters.DateFrom != nil {
		query = query.Where("orders.created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("orders.created_at <= ?", *filters.DateTo)
	}

	err := query.Scan(&metrics).Error
	if err != nil {
		return nil, err
	}

	// Get view count if product ID is specified
	if filters.ProductID != nil {
		viewQuery := r.db.WithContext(ctx).
			Model(&entities.AnalyticsEvent{}).
			Where("event_type = ? AND product_id = ?", "product_view", *filters.ProductID)

		if filters.DateFrom != nil {
			viewQuery = viewQuery.Where("created_at >= ?", *filters.DateFrom)
		}

		if filters.DateTo != nil {
			viewQuery = viewQuery.Where("created_at <= ?", *filters.DateTo)
		}

		err = viewQuery.Count(&metrics.ViewCount).Error
		if err != nil {
			return nil, err
		}
	}

	return &metrics, nil
}

// GetUserMetrics gets user analytics metrics with filters
func (r *analyticsRepository) GetUserMetrics(ctx context.Context, filters repositories.UserMetricsFilters) (*repositories.UserMetrics, error) {
	var metrics repositories.UserMetrics

	query := r.db.WithContext(ctx).Model(&entities.User{})

	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	// Get new users count
	err := query.Count(&metrics.NewUsers).Error
	if err != nil {
		return nil, err
	}

	// Get active users count (users who placed orders)
	activeQuery := r.db.WithContext(ctx).Model(&entities.Order{}).Select("COUNT(DISTINCT user_id)")
	if filters.DateFrom != nil {
		activeQuery = activeQuery.Where("created_at >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		activeQuery = activeQuery.Where("created_at <= ?", *filters.DateTo)
	}

	err = activeQuery.Scan(&metrics.ActiveUsers).Error
	if err != nil {
		return nil, err
	}

	// Get total users
	err = r.db.WithContext(ctx).Model(&entities.User{}).Count(&metrics.TotalUsers).Error
	if err != nil {
		return nil, err
	}

	return &metrics, nil
}

// GetTrafficMetrics gets website traffic metrics with filters
func (r *analyticsRepository) GetTrafficMetrics(ctx context.Context, filters repositories.TrafficMetricsFilters) (*repositories.TrafficMetrics, error) {
	var metrics repositories.TrafficMetrics

	query := r.db.WithContext(ctx).Model(&entities.AnalyticsEvent{}).Where("event_type = ?", "page_view")

	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	// Get page views
	err := query.Count(&metrics.PageViews).Error
	if err != nil {
		return nil, err
	}

	// Get unique visitors
	err = query.Select("COUNT(DISTINCT session_id)").Scan(&metrics.UniqueVisitors).Error
	if err != nil {
		return nil, err
	}

	// Get bounce rate (sessions with only one page view)
	var singlePageSessions int64
	err = query.Select("COUNT(DISTINCT session_id)").
		Group("session_id").
		Having("COUNT(*) = 1").
		Scan(&singlePageSessions).Error
	if err != nil {
		return nil, err
	}

	if metrics.UniqueVisitors > 0 {
		metrics.BounceRate = float64(singlePageSessions) / float64(metrics.UniqueVisitors) * 100
	}

	return &metrics, nil
}

// GetTopProducts gets top selling products
func (r *analyticsRepository) GetTopProducts(ctx context.Context, period string, limit int) ([]*repositories.TopProduct, error) {
	var topProducts []*repositories.TopProduct

	// Calculate date range based on period
	var from, to time.Time
	now := time.Now()
	switch period {
	case "week":
		from = now.AddDate(0, 0, -7)
		to = now
	case "month":
		from = now.AddDate(0, -1, 0)
		to = now
	case "year":
		from = now.AddDate(-1, 0, 0)
		to = now
	default: // today
		from = now.Truncate(24 * time.Hour)
		to = now
	}

	err := r.db.WithContext(ctx).
		Table("order_items").
		Select("products.id, products.name, products.price, SUM(order_items.quantity) as units_sold, SUM(order_items.price * order_items.quantity) as revenue").
		Joins("JOIN products ON order_items.product_id = products.id").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.created_at BETWEEN ? AND ? AND orders.status = ?", from, to, entities.OrderStatusDelivered).
		Group("products.id, products.name, products.price").
		Order("units_sold DESC").
		Limit(limit).
		Scan(&topProducts).Error

	return topProducts, err
}

// GetTopCategories gets top performing categories
func (r *analyticsRepository) GetTopCategories(ctx context.Context, period string, limit int) ([]*repositories.TopCategory, error) {
	var topCategories []*repositories.TopCategory

	// Calculate date range based on period
	var from, to time.Time
	now := time.Now()
	switch period {
	case "week":
		from = now.AddDate(0, 0, -7)
		to = now
	case "month":
		from = now.AddDate(0, -1, 0)
		to = now
	case "year":
		from = now.AddDate(-1, 0, 0)
		to = now
	default: // today
		from = now.Truncate(24 * time.Hour)
		to = now
	}

	err := r.db.WithContext(ctx).
		Table("order_items").
		Select("categories.id, categories.name, SUM(order_items.quantity) as units_sold, SUM(order_items.price * order_items.quantity) as revenue").
		Joins("JOIN products ON order_items.product_id = products.id").
		Joins("JOIN categories ON products.category_id = categories.id").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.created_at BETWEEN ? AND ? AND orders.status = ?", from, to, entities.OrderStatusDelivered).
		Group("categories.id, categories.name").
		Order("revenue DESC").
		Limit(limit).
		Scan(&topCategories).Error

	return topCategories, err
}

// GetRevenueByPeriod gets revenue data grouped by period
func (r *analyticsRepository) GetRevenueByPeriod(ctx context.Context, from, to time.Time, period string) ([]*entities.RevenueData, error) {
	var revenueData []*entities.RevenueData

	var dateFormat string
	switch period {
	case "day":
		dateFormat = "DATE(created_at)"
	case "week":
		dateFormat = "YEARWEEK(created_at)"
	case "month":
		dateFormat = "DATE_FORMAT(created_at, '%Y-%m')"
	default:
		dateFormat = "DATE(created_at)"
	}

	err := r.db.WithContext(ctx).
		Table("orders").
		Select(dateFormat+" as period, SUM(total) as revenue, COUNT(*) as order_count").
		Where("created_at BETWEEN ? AND ? AND status = ? AND payment_status = ?", from, to, entities.OrderStatusDelivered, entities.PaymentStatusPaid).
		Group("period").
		Order("period ASC").
		Scan(&revenueData).Error

	return revenueData, err
}

// GetConversionMetrics gets conversion rate metrics
func (r *analyticsRepository) GetConversionMetrics(ctx context.Context, from, to time.Time) (*entities.ConversionMetrics, error) {
	var metrics entities.ConversionMetrics

	// Get total sessions
	err := r.db.WithContext(ctx).
		Model(&entities.AnalyticsEvent{}).
		Select("COUNT(DISTINCT session_id)").
		Where("created_at BETWEEN ? AND ?", from, to).
		Scan(&metrics.TotalSessions).Error
	if err != nil {
		return nil, err
	}

	// Get sessions with orders
	err = r.db.WithContext(ctx).
		Table("analytics_events").
		Select("COUNT(DISTINCT analytics_events.session_id)").
		Joins("JOIN orders ON analytics_events.user_id = orders.user_id").
		Where("analytics_events.created_at BETWEEN ? AND ? AND orders.created_at BETWEEN ? AND ?",
			from, to, from, to).
		Scan(&metrics.ConvertedSessions).Error
	if err != nil {
		return nil, err
	}

	// Calculate conversion rate
	if metrics.TotalSessions > 0 {
		metrics.ConversionRate = float64(metrics.ConvertedSessions) / float64(metrics.TotalSessions) * 100
	}

	return &metrics, nil
}

// GetRealTimeMetrics gets real-time metrics
func (r *analyticsRepository) GetRealTimeMetrics(ctx context.Context) (*entities.RealTimeMetrics, error) {
	var metrics entities.RealTimeMetrics
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)

	// Get active users in last hour
	err := r.db.WithContext(ctx).
		Model(&entities.AnalyticsEvent{}).
		Select("COUNT(DISTINCT user_id)").
		Where("created_at >= ?", oneHourAgo).
		Scan(&metrics.ActiveUsers).Error
	if err != nil {
		return nil, err
	}

	// Get page views in last hour
	err = r.db.WithContext(ctx).
		Model(&entities.AnalyticsEvent{}).
		Where("event_type = ? AND created_at >= ?", "page_view", oneHourAgo).
		Count(&metrics.PageViews).Error
	if err != nil {
		return nil, err
	}

	// Get orders in last hour
	err = r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("created_at >= ?", oneHourAgo).
		Count(&metrics.Orders).Error
	if err != nil {
		return nil, err
	}

	// Get revenue in last hour
	err = r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Select("COALESCE(SUM(total), 0)").
		Where("created_at >= ? AND status = ? AND payment_status = ?", oneHourAgo, entities.OrderStatusDelivered, entities.PaymentStatusPaid).
		Scan(&metrics.Revenue).Error
	if err != nil {
		return nil, err
	}

	return &metrics, nil
}

// GetCustomerLifetimeValue gets customer lifetime value metrics
func (r *analyticsRepository) GetCustomerLifetimeValue(ctx context.Context, userID uuid.UUID) (*entities.CustomerLifetimeValue, error) {
	var clv entities.CustomerLifetimeValue

	// Get total spent by customer
	err := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Select("COALESCE(SUM(total), 0) as total_spent, COUNT(*) as order_count").
		Where("user_id = ? AND status = ? AND payment_status = ?", userID, entities.OrderStatusDelivered, entities.PaymentStatusPaid).
		Scan(&clv).Error
	if err != nil {
		return nil, err
	}

	// Get first order date
	err = r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Select("MIN(created_at)").
		Where("user_id = ?", userID).
		Scan(&clv.FirstOrderDate).Error
	if err != nil {
		return nil, err
	}

	// Get last order date
	err = r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Select("MAX(created_at)").
		Where("user_id = ?", userID).
		Scan(&clv.LastOrderDate).Error
	if err != nil {
		return nil, err
	}

	// Calculate average order value
	if clv.OrderCount > 0 {
		clv.AverageOrderValue = clv.TotalSpent / float64(clv.OrderCount)
	}

	return &clv, nil
}

// GetInventoryMetrics gets inventory analytics
func (r *analyticsRepository) GetInventoryMetrics(ctx context.Context) (*entities.InventoryMetrics, error) {
	var metrics entities.InventoryMetrics

	// Get total products
	err := r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Count(&metrics.TotalProducts).Error
	if err != nil {
		return nil, err
	}

	// Get low stock items
	err = r.db.WithContext(ctx).
		Table("inventories").
		Where("quantity_available <= reorder_level AND quantity_available > 0").
		Count(&metrics.LowStockItems).Error
	if err != nil {
		return nil, err
	}

	// Get out of stock items
	err = r.db.WithContext(ctx).
		Table("inventories").
		Where("quantity_available = 0").
		Count(&metrics.OutOfStockItems).Error
	if err != nil {
		return nil, err
	}

	// Get total inventory value
	err = r.db.WithContext(ctx).
		Table("inventories").
		Select("COALESCE(SUM(inventories.quantity_on_hand * products.price), 0)").
		Joins("JOIN products ON inventories.product_id = products.id").
		Scan(&metrics.TotalInventoryValue).Error
	if err != nil {
		return nil, err
	}

	return &metrics, nil
}

// CountEvents counts analytics events with filters
func (r *analyticsRepository) CountEvents(ctx context.Context, filters repositories.EventFilters) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.AnalyticsEvent{})

	if filters.EventType != "" {
		query = query.Where("event_type = ?", filters.EventType)
	}

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.ProductID != nil {
		query = query.Where("product_id = ?", *filters.ProductID)
	}

	if filters.SessionID != "" {
		query = query.Where("session_id = ?", filters.SessionID)
	}

	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	err := query.Count(&count).Error
	return count, err
}

// CreateEvent creates an analytics event (alias for RecordEvent)
func (r *analyticsRepository) CreateEvent(ctx context.Context, event *entities.AnalyticsEvent) error {
	return r.RecordEvent(ctx, event)
}

// GetEvents gets analytics events with filters
func (r *analyticsRepository) GetEvents(ctx context.Context, filters repositories.EventFilters) ([]*entities.AnalyticsEvent, error) {
	var events []*entities.AnalyticsEvent
	query := r.db.WithContext(ctx).Model(&entities.AnalyticsEvent{})

	if filters.EventType != "" {
		query = query.Where("event_type = ?", filters.EventType)
	}

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.ProductID != nil {
		query = query.Where("product_id = ?", *filters.ProductID)
	}

	if filters.SessionID != "" {
		query = query.Where("session_id = ?", filters.SessionID)
	}

	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	// Apply sorting
	sortBy := "created_at"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}
	sortOrder := "DESC"
	if filters.SortOrder != "" {
		sortOrder = filters.SortOrder
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&events).Error
	return events, err
}

// ExecuteCustomQuery executes a custom analytics query
func (r *analyticsRepository) ExecuteCustomQuery(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// Convert params map to slice for Raw query
	var args []interface{}
	for _, v := range params {
		args = append(args, v)
	}

	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		result := make(map[string]interface{})
		for i, col := range columns {
			result[col] = values[i]
		}
		results = append(results, result)
	}

	return results, nil
}

// GetActiveUsers gets active users count within a duration
func (r *analyticsRepository) GetActiveUsers(ctx context.Context, duration time.Duration) (int64, error) {
	since := time.Now().Add(-duration)
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.AnalyticsEvent{}).
		Where("created_at >= ?", since).
		Select("COUNT(DISTINCT user_id)").
		Scan(&count).Error
	return count, err
}

// GetConversionRate gets conversion rate as percentage
func (r *analyticsRepository) GetConversionRate(ctx context.Context, from, to time.Time) (float64, error) {
	var totalSessions, convertedSessions int64

	// Get total sessions (unique visitors)
	err := r.db.WithContext(ctx).
		Model(&entities.AnalyticsEvent{}).
		Where("event_type = ? AND created_at BETWEEN ? AND ?", "page_view", from, to).
		Select("COUNT(DISTINCT session_id)").
		Scan(&totalSessions).Error
	if err != nil {
		return 0, err
	}

	// Get converted sessions (orders placed)
	err = r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("created_at BETWEEN ? AND ?", from, to).
		Select("COUNT(DISTINCT user_id)").
		Scan(&convertedSessions).Error
	if err != nil {
		return 0, err
	}

	// Calculate conversion rate
	if totalSessions > 0 {
		return float64(convertedSessions) / float64(totalSessions) * 100, nil
	}

	return 0, nil
}

// GetDashboardMetrics gets dashboard metrics (placeholder implementation)
func (r *analyticsRepository) GetDashboardMetrics(ctx context.Context, dateFrom, dateTo time.Time) (*repositories.DashboardMetrics, error) {
	// This is a placeholder implementation
	// In a real system, you would aggregate various metrics
	return &repositories.DashboardMetrics{
		TotalUsers:     0,
		TotalOrders:    0,
		TotalRevenue:   0,
		ConversionRate: 0,
	}, nil
}

// GetFunnelAnalysis gets funnel analysis data (placeholder)
func (r *analyticsRepository) GetFunnelAnalysis(ctx context.Context, steps []string, from, to time.Time) (*repositories.FunnelAnalysis, error) {
	// Placeholder implementation
	return &repositories.FunnelAnalysis{
		Steps:          steps,
		TotalUsers:     0,
		ConversionRate: 0,
	}, nil
}

// GetOnlineVisitors gets current online visitors count
func (r *analyticsRepository) GetOnlineVisitors(ctx context.Context) (int64, error) {
	// Consider users active in the last 5 minutes as online
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.AnalyticsEvent{}).
		Where("created_at >= ?", fiveMinutesAgo).
		Select("COUNT(DISTINCT user_id)").
		Scan(&count).Error
	return count, err
}

// GetRetentionRate gets user retention rate
func (r *analyticsRepository) GetRetentionRate(ctx context.Context, period string) (float64, error) {
	// Placeholder implementation - calculate retention rate
	return 75.5, nil
}

// GetTodayOrders gets today's order count
func (r *analyticsRepository) GetTodayOrders(ctx context.Context) (int64, error) {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("created_at >= ? AND created_at < ?", today, tomorrow).
		Count(&count).Error
	return count, err
}

// GetTodayRevenue gets today's revenue
func (r *analyticsRepository) GetTodayRevenue(ctx context.Context) (float64, error) {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	var revenue float64
	err := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Select("COALESCE(SUM(total), 0)").
		Where("created_at >= ? AND created_at < ? AND status = ? AND payment_status = ?", today, tomorrow, entities.OrderStatusDelivered, entities.PaymentStatusPaid).
		Scan(&revenue).Error
	return revenue, err
}

// GetTopPages gets top performing pages
func (r *analyticsRepository) GetTopPages(ctx context.Context, period string, limit int) ([]*repositories.TopPage, error) {
	var topPages []*repositories.TopPage

	// Calculate date range based on period
	var from, to time.Time
	now := time.Now()
	switch period {
	case "week":
		from = now.AddDate(0, 0, -7)
		to = now
	case "month":
		from = now.AddDate(0, -1, 0)
		to = now
	case "year":
		from = now.AddDate(-1, 0, 0)
		to = now
	default: // today
		from = now.Truncate(24 * time.Hour)
		to = now
	}

	err := r.db.WithContext(ctx).
		Model(&entities.AnalyticsEvent{}).
		Select("page, COUNT(*) as views, COUNT(DISTINCT user_id) as unique_views").
		Where("event_type = ? AND created_at BETWEEN ? AND ?", "page_view", from, to).
		Group("page").
		Order("views DESC").
		Limit(limit).
		Scan(&topPages).Error

	return topPages, err
}

// GetUserCohorts gets user cohort analysis
func (r *analyticsRepository) GetUserCohorts(ctx context.Context, period string) (*repositories.CohortAnalysis, error) {
	// Placeholder implementation
	return &repositories.CohortAnalysis{
		Period:        period,
		TotalUsers:    0,
		RetentionRate: 0,
	}, nil
}
