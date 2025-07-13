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
