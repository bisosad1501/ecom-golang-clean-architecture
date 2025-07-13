package entities

import (
	"time"

	"github.com/google/uuid"
)

// SearchSuggestion represents search suggestions
type SearchSuggestion struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Query        string    `json:"query" gorm:"not null;uniqueIndex" validate:"required"`
	SearchCount  int       `json:"search_count" gorm:"default:1"`
	Frequency    int       `json:"frequency" gorm:"default:1"`
	ResultCount  int       `json:"result_count" gorm:"default:0"`
	LastSearched time.Time `json:"last_searched" gorm:"default:CURRENT_TIMESTAMP"`
	IsTrending   bool      `json:"is_trending" gorm:"default:false"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for SearchSuggestion entity
func (SearchSuggestion) TableName() string {
	return "search_suggestions"
}

// PopularSearch represents popular search terms
type PopularSearch struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Query       string    `json:"query" gorm:"not null;uniqueIndex" validate:"required"`
	SearchCount int       `json:"search_count" gorm:"default:1"`
	Period      string    `json:"period" gorm:"default:'daily'"` // daily, weekly, monthly
	Date        time.Time `json:"date" gorm:"index"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for PopularSearch entity
func (PopularSearch) TableName() string {
	return "popular_searches"
}

// SearchFilter represents saved search filters
type SearchFilter struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index" validate:"required"`
	Name        string    `json:"name" gorm:"not null" validate:"required"`
	Query       string    `json:"query"`
	Filters     string    `json:"filters" gorm:"type:jsonb"` // JSON string of filter parameters
	IsDefault   bool      `json:"is_default" gorm:"default:false"`
	IsPublic    bool      `json:"is_public" gorm:"default:false"`
	UsageCount  int       `json:"usage_count" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for SearchFilter entity
func (SearchFilter) TableName() string {
	return "search_filters"
}

// SearchHistory represents user search history
type SearchHistory struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index" validate:"required"`
	Query     string    `json:"query" gorm:"not null" validate:"required"`
	Filters   string    `json:"filters" gorm:"type:jsonb"` // JSON string of filter parameters
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;index"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for SearchHistory entity
func (SearchHistory) TableName() string {
	return "search_history"
}

// SearchEvent represents individual search events for detailed analytics
type SearchEvent struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Query            string     `json:"query" gorm:"not null;index" validate:"required"`
	UserID           *uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	ResultsCount     int        `json:"results_count" gorm:"default:0"`
	ClickedProductID *uuid.UUID `json:"clicked_product_id" gorm:"type:uuid"`
	SessionID        string     `json:"session_id" gorm:"index"`
	IPAddress        string     `json:"ip_address"`
	UserAgent        string     `json:"user_agent" gorm:"type:text"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime;index"`

	// Relationships
	User           *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ClickedProduct *Product `json:"clicked_product,omitempty" gorm:"foreignKey:ClickedProductID"`
}

// TableName returns the table name for SearchEvent entity
func (SearchEvent) TableName() string {
	return "search_events"
}

// AutocompleteEntry represents an autocomplete entry
type AutocompleteEntry struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Type        string    `json:"type" gorm:"index"`        // product, category, brand, tag, query
	Value       string    `json:"value" gorm:"index"`       // the actual text value
	DisplayText string    `json:"display_text"`             // formatted display text
	EntityID    *uuid.UUID `json:"entity_id" gorm:"type:uuid;index"` // reference to actual entity
	Priority    int       `json:"priority" gorm:"default:0"` // higher priority = shown first
	SearchCount int       `json:"search_count" gorm:"default:0"` // how often this was searched
	ClickCount  int       `json:"click_count" gorm:"default:0"`  // how often this was clicked
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Metadata    string    `json:"metadata" gorm:"type:jsonb"` // additional data (price, image, etc.)
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for AutocompleteEntry entity
func (AutocompleteEntry) TableName() string {
	return "autocomplete_entries"
}

// SearchTrend represents search trend data
type SearchTrend struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Query       string    `json:"query" gorm:"index"`
	SearchCount int       `json:"search_count" gorm:"default:0"`
	Period      string    `json:"period" gorm:"index"` // daily, weekly, monthly
	Date        time.Time `json:"date" gorm:"index"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for SearchTrend entity
func (SearchTrend) TableName() string {
	return "search_trends"
}

// UserSearchPreference represents user search preferences
type UserSearchPreference struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID            uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	PreferredCategories []string `json:"preferred_categories" gorm:"type:jsonb"`
	PreferredBrands     []string `json:"preferred_brands" gorm:"type:jsonb"`
	SearchLanguage      string   `json:"search_language" gorm:"default:'en'"`
	AutocompleteEnabled bool     `json:"autocomplete_enabled" gorm:"default:true"`
	SearchHistoryEnabled bool    `json:"search_history_enabled" gorm:"default:true"`
	PersonalizedResults  bool    `json:"personalized_results" gorm:"default:true"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for UserSearchPreference entity
func (UserSearchPreference) TableName() string {
	return "user_search_preferences"
}

// SearchSession represents a search session
type SearchSession struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SessionID     string    `json:"session_id" gorm:"index"`
	UserID        *uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	StartTime     time.Time `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
	SearchCount   int       `json:"search_count" gorm:"default:0"`
	ClickCount    int       `json:"click_count" gorm:"default:0"`
	ConversionCount int     `json:"conversion_count" gorm:"default:0"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for SearchSession entity
func (SearchSession) TableName() string {
	return "search_sessions"
}
