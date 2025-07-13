package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"ecom-golang-clean-architecture/internal/domain/entities"
)

// UserPersonalizationRepository defines the interface for user personalization operations
type UserPersonalizationRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, personalization *entities.UserPersonalization) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.UserPersonalization, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.UserPersonalization, error)
	Update(ctx context.Context, personalization *entities.UserPersonalization) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error

	// Preference management
	UpdateCategoryPreferences(ctx context.Context, userID uuid.UUID, preferences map[string]float64) error
	UpdateBrandPreferences(ctx context.Context, userID uuid.UUID, preferences map[string]float64) error
	UpdatePriceRangePreference(ctx context.Context, userID uuid.UUID, minPrice, maxPrice float64) error
	
	// Analytics updates
	IncrementViews(ctx context.Context, userID uuid.UUID, count int) error
	IncrementSearches(ctx context.Context, userID uuid.UUID, count int) error
	UpdateUniqueProductsViewed(ctx context.Context, userID uuid.UUID, count int) error
	UpdateLastAnalyzed(ctx context.Context, userID uuid.UUID, timestamp time.Time) error

	// Recommendation data
	GetUsersForRecommendation(ctx context.Context, engine string, limit int) ([]*entities.UserPersonalization, error)
	GetSimilarUsers(ctx context.Context, userID uuid.UUID, limit int) ([]*entities.UserPersonalization, error)
	
	// Bulk operations
	BulkUpdatePersonalization(ctx context.Context, updates []PersonalizationUpdate) error
	AnalyzeUserBehavior(ctx context.Context, userID uuid.UUID) (*UserBehaviorAnalysis, error)
}

// PersonalizationUpdate represents a bulk personalization update
type PersonalizationUpdate struct {
	UserID               uuid.UUID              `json:"user_id"`
	CategoryPreferences  map[string]float64     `json:"category_preferences,omitempty"`
	BrandPreferences     map[string]float64     `json:"brand_preferences,omitempty"`
	PriceRangePreference *PriceRangePreference  `json:"price_range_preference,omitempty"`
	BehavioralData       *BehavioralData        `json:"behavioral_data,omitempty"`
}

// PriceRangePreference represents price range preferences
type PriceRangePreference struct {
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
	Currency string  `json:"currency"`
}

// BehavioralData represents user behavioral data
type BehavioralData struct {
	AverageOrderValue     float64 `json:"average_order_value"`
	PreferredShoppingTime string  `json:"preferred_shopping_time"`
	ShoppingFrequency     string  `json:"shopping_frequency"`
}

// UserBehaviorAnalysis represents analyzed user behavior
type UserBehaviorAnalysis struct {
	UserID uuid.UUID `json:"user_id"`
	
	// Category analysis
	TopCategories []CategoryPreference `json:"top_categories"`
	
	// Brand analysis
	TopBrands []BrandPreference `json:"top_brands"`
	
	// Price analysis
	PriceRange struct {
		MinPrice    float64 `json:"min_price"`
		MaxPrice    float64 `json:"max_price"`
		AveragePrice float64 `json:"average_price"`
		PriceVariance float64 `json:"price_variance"`
	} `json:"price_range"`
	
	// Behavioral patterns
	ShoppingPatterns struct {
		PreferredDays  []string `json:"preferred_days"`
		PreferredHours []int    `json:"preferred_hours"`
		SessionLength  float64  `json:"average_session_length"`
		PagesPerSession float64 `json:"average_pages_per_session"`
	} `json:"shopping_patterns"`
	
	// Engagement metrics
	EngagementScore float64 `json:"engagement_score"`
	LoyaltyScore    float64 `json:"loyalty_score"`
	
	AnalyzedAt time.Time `json:"analyzed_at"`
}

// CategoryPreference represents category preference data
type CategoryPreference struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Score        float64   `json:"score"`
	ViewCount    int       `json:"view_count"`
	PurchaseCount int      `json:"purchase_count"`
}

// BrandPreference represents brand preference data
type BrandPreference struct {
	BrandID       uuid.UUID `json:"brand_id"`
	BrandName     string    `json:"brand_name"`
	Score         float64   `json:"score"`
	ViewCount     int       `json:"view_count"`
	PurchaseCount int       `json:"purchase_count"`
}
