package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"ecom-golang-clean-architecture/internal/domain/entities"
)

// UserSearchHistoryRepository defines the interface for user search history operations
type UserSearchHistoryRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, history *entities.UserSearchHistory) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.UserSearchHistory, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.UserSearchHistory, error)
	Update(ctx context.Context, history *entities.UserSearchHistory) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error

	// Search and filtering
	SearchHistory(ctx context.Context, req SearchHistoryRequest) (*SearchHistoryResponse, error)
	GetPopularSearches(ctx context.Context, limit int, timeRange time.Duration) ([]*PopularSearch, error)
	GetUserSearchStats(ctx context.Context, userID uuid.UUID) (*UserSearchStats, error)

	// Analytics
	GetSearchTrends(ctx context.Context, req SearchTrendsRequest) (*SearchTrendsResponse, error)
	GetSearchAnalytics(ctx context.Context, req SearchAnalyticsRequest) (*SearchAnalyticsResponse, error)

	// Cleanup
	CleanupOldHistory(ctx context.Context, olderThan time.Duration) (int64, error)
}

// SearchHistoryRequest represents search history query parameters
type SearchHistoryRequest struct {
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	Query     *string    `json:"query,omitempty"`
	DateFrom  *time.Time `json:"date_from,omitempty"`
	DateTo    *time.Time `json:"date_to,omitempty"`
	HasClicks *bool      `json:"has_clicks,omitempty"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

// SearchHistoryResponse represents search history query results
type SearchHistoryResponse struct {
	History []*entities.UserSearchHistory `json:"history"`
	Total   int64                         `json:"total"`
}

// PopularSearch represents popular search data
type PopularSearch struct {
	Query       string    `json:"query"`
	SearchCount int64     `json:"search_count"`
	ClickRate   float64   `json:"click_rate"`
	LastUsed    time.Time `json:"last_used"`
}

// UserSearchStats represents user search statistics
type UserSearchStats struct {
	TotalSearches    int64     `json:"total_searches"`
	UniqueQueries    int64     `json:"unique_queries"`
	ClickThroughRate float64   `json:"click_through_rate"`
	AverageResults   float64   `json:"average_results"`
	TopQueries       []string  `json:"top_queries"`
	LastSearch       *time.Time `json:"last_search,omitempty"`
}

// SearchTrendsRequest represents search trends query parameters
type SearchTrendsRequest struct {
	TimeRange string     `json:"time_range"` // day, week, month, year
	DateFrom  *time.Time `json:"date_from,omitempty"`
	DateTo    *time.Time `json:"date_to,omitempty"`
	Limit     int        `json:"limit"`
}

// SearchTrendsResponse represents search trends data
type SearchTrendsResponse struct {
	Trends []SearchTrendData `json:"trends"`
	Period string            `json:"period"`
}

// SearchTrendData represents individual trend data point
type SearchTrendData struct {
	Date        time.Time `json:"date"`
	Query       string    `json:"query"`
	SearchCount int64     `json:"search_count"`
	ClickCount  int64     `json:"click_count"`
	ClickRate   float64   `json:"click_rate"`
}

// SearchAnalyticsRequest represents search analytics query parameters
type SearchAnalyticsRequest struct {
	UserID   *uuid.UUID `json:"user_id,omitempty"`
	DateFrom *time.Time `json:"date_from,omitempty"`
	DateTo   *time.Time `json:"date_to,omitempty"`
	GroupBy  string     `json:"group_by"` // hour, day, week, month
}

// SearchAnalyticsResponse represents search analytics data
type SearchAnalyticsResponse struct {
	Overview struct {
		TotalSearches    int64   `json:"total_searches"`
		UniqueUsers      int64   `json:"unique_users"`
		UniqueQueries    int64   `json:"unique_queries"`
		ClickThroughRate float64 `json:"click_through_rate"`
		AverageResults   float64 `json:"average_results"`
	} `json:"overview"`
	
	TimeSeriesData []SearchTimeSeriesPoint `json:"time_series_data"`
	TopQueries     []PopularSearch   `json:"top_queries"`
	QueryAnalysis  struct {
		ShortQueries  int64 `json:"short_queries"`  // 1-2 words
		MediumQueries int64 `json:"medium_queries"` // 3-5 words
		LongQueries   int64 `json:"long_queries"`   // 6+ words
	} `json:"query_analysis"`
}

// SearchTimeSeriesPoint represents a data point in time series for search analytics
type SearchTimeSeriesPoint struct {
	Timestamp    time.Time `json:"timestamp"`
	SearchCount  int64     `json:"search_count"`
	ClickCount   int64     `json:"click_count"`
	UniqueUsers  int64     `json:"unique_users"`
	ClickRate    float64   `json:"click_rate"`
}
