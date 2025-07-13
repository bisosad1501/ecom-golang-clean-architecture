package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"ecom-golang-clean-architecture/internal/domain/entities"
)

// SavedSearchRepository defines the interface for saved search operations
type SavedSearchRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, savedSearch *entities.SavedSearch) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.SavedSearch, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.SavedSearch, error)
	Update(ctx context.Context, savedSearch *entities.SavedSearch) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error

	// Search and filtering
	SearchSavedSearches(ctx context.Context, req SavedSearchRequest) (*SavedSearchResponse, error)
	GetActiveSavedSearches(ctx context.Context, userID uuid.UUID) ([]*entities.SavedSearch, error)
	GetSavedSearchesByAlert(ctx context.Context, alertType string) ([]*entities.SavedSearch, error)

	// Alert management
	GetSearchesForPriceAlert(ctx context.Context) ([]*entities.SavedSearch, error)
	GetSearchesForStockAlert(ctx context.Context) ([]*entities.SavedSearch, error)
	GetSearchesForNewItemAlert(ctx context.Context) ([]*entities.SavedSearch, error)
	UpdateLastChecked(ctx context.Context, id uuid.UUID, timestamp time.Time) error
	UpdateLastNotified(ctx context.Context, id uuid.UUID, timestamp time.Time) error

	// Analytics
	GetSavedSearchStats(ctx context.Context, userID uuid.UUID) (*SavedSearchStats, error)
	GetPopularSavedSearches(ctx context.Context, limit int) ([]*PopularSavedSearch, error)
}

// SavedSearchRequest represents saved search query parameters
type SavedSearchRequest struct {
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	Name       *string    `json:"name,omitempty"`
	Query      *string    `json:"query,omitempty"`
	IsActive   *bool      `json:"is_active,omitempty"`
	HasAlerts  *bool      `json:"has_alerts,omitempty"`
	DateFrom   *time.Time `json:"date_from,omitempty"`
	DateTo     *time.Time `json:"date_to,omitempty"`
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset"`
}

// SavedSearchResponse represents saved search query results
type SavedSearchResponse struct {
	SavedSearches []*entities.SavedSearch `json:"saved_searches"`
	Total         int64                   `json:"total"`
}

// SavedSearchStats represents saved search statistics
type SavedSearchStats struct {
	TotalSavedSearches int64     `json:"total_saved_searches"`
	ActiveSearches     int64     `json:"active_searches"`
	SearchesWithAlerts int64     `json:"searches_with_alerts"`
	LastCreated        *time.Time `json:"last_created,omitempty"`
	MostUsedQuery      string    `json:"most_used_query"`
}

// PopularSavedSearch represents popular saved search data
type PopularSavedSearch struct {
	Query     string `json:"query"`
	UserCount int64  `json:"user_count"`
	AlertRate float64 `json:"alert_rate"`
}

// SavedSearchAlert represents an alert to be sent
type SavedSearchAlert struct {
	SavedSearchID uuid.UUID `json:"saved_search_id"`
	UserID        uuid.UUID `json:"user_id"`
	AlertType     string    `json:"alert_type"` // price, stock, new_item
	Message       string    `json:"message"`
	Data          map[string]interface{} `json:"data"`
	CreatedAt     time.Time `json:"created_at"`
}
