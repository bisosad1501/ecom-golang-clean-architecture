package entities

import (
	"time"

	"github.com/google/uuid"
)

// RecommendationType represents the type of recommendation
type RecommendationType string

const (
	RecommendationTypeRelated           RecommendationType = "related"
	RecommendationTypeSimilar           RecommendationType = "similar"
	RecommendationTypeFrequentlyBought  RecommendationType = "frequently_bought"
	RecommendationTypeTrending          RecommendationType = "trending"
	RecommendationTypePersonalized      RecommendationType = "personalized"
	RecommendationTypeCrossSell         RecommendationType = "cross_sell"
	RecommendationTypeUpSell            RecommendationType = "up_sell"
	RecommendationTypeRecentlyViewed    RecommendationType = "recently_viewed"
	RecommendationTypeBasedOnCategory   RecommendationType = "based_on_category"
	RecommendationTypeBasedOnBrand      RecommendationType = "based_on_brand"
)

// ProductRecommendation represents a product recommendation
type ProductRecommendation struct {
	ID                uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID         uuid.UUID          `json:"product_id" gorm:"type:uuid;not null;index"`
	RecommendedID     uuid.UUID          `json:"recommended_id" gorm:"type:uuid;not null;index"`
	Type              RecommendationType `json:"type" gorm:"not null;index"`
	Score             float64            `json:"score" gorm:"default:0"`
	Reason            string             `json:"reason" gorm:"type:text"`
	IsActive          bool               `json:"is_active" gorm:"default:true"`
	CreatedAt         time.Time          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time          `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Product     Product `json:"product" gorm:"foreignKey:ProductID"`
	Recommended Product `json:"recommended" gorm:"foreignKey:RecommendedID"`
}

// TableName returns the table name for ProductRecommendation entity
func (ProductRecommendation) TableName() string {
	return "product_recommendations"
}

// UserProductInteraction represents user interactions with products
type UserProductInteraction struct {
	ID            uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        *uuid.UUID      `json:"user_id" gorm:"type:uuid;index"` // Nullable for guest users
	SessionID     *string         `json:"session_id" gorm:"index"`        // For guest users
	ProductID     uuid.UUID       `json:"product_id" gorm:"type:uuid;not null;index"`
	InteractionType InteractionType `json:"interaction_type" gorm:"not null;index"`
	Value         float64         `json:"value" gorm:"default:1"` // Weight/score for the interaction
	Metadata      string          `json:"metadata" gorm:"type:text"` // JSON metadata
	CreatedAt     time.Time       `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	User    *User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Product Product `json:"product" gorm:"foreignKey:ProductID"`
}

// TableName returns the table name for UserProductInteraction entity
func (UserProductInteraction) TableName() string {
	return "user_product_interactions"
}

// InteractionType represents the type of user interaction
type InteractionType string

const (
	InteractionTypeView         InteractionType = "view"
	InteractionTypeAddToCart    InteractionType = "add_to_cart"
	InteractionTypeRemoveFromCart InteractionType = "remove_from_cart"
	InteractionTypePurchase     InteractionType = "purchase"
	InteractionTypeWishlist     InteractionType = "wishlist"
	InteractionTypeReview       InteractionType = "review"
	InteractionTypeShare        InteractionType = "share"
	InteractionTypeCompare      InteractionType = "compare"
	InteractionTypeSearch       InteractionType = "search"
	InteractionTypeClick        InteractionType = "click"
)

// ProductSimilarity represents similarity between products
type ProductSimilarity struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID     uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	SimilarID     uuid.UUID `json:"similar_id" gorm:"type:uuid;not null;index"`
	SimilarityScore float64 `json:"similarity_score" gorm:"not null"`
	Algorithm     string    `json:"algorithm" gorm:"not null"` // cosine, jaccard, etc.
	Features      string    `json:"features" gorm:"type:text"` // JSON of features used
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Product Product `json:"product" gorm:"foreignKey:ProductID"`
	Similar Product `json:"similar" gorm:"foreignKey:SimilarID"`
}

// TableName returns the table name for ProductSimilarity entity
func (ProductSimilarity) TableName() string {
	return "product_similarities"
}

// FrequentlyBoughtTogether represents products frequently bought together
type FrequentlyBoughtTogether struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID   uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	WithID      uuid.UUID `json:"with_id" gorm:"type:uuid;not null;index"`
	Frequency   int       `json:"frequency" gorm:"default:1"`
	Confidence  float64   `json:"confidence" gorm:"default:0"` // Support/confidence from market basket analysis
	Support     float64   `json:"support" gorm:"default:0"`
	Lift        float64   `json:"lift" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Product Product `json:"product" gorm:"foreignKey:ProductID"`
	With    Product `json:"with" gorm:"foreignKey:WithID"`
}

// TableName returns the table name for FrequentlyBoughtTogether entity
func (FrequentlyBoughtTogether) TableName() string {
	return "frequently_bought_together"
}

// TrendingProduct represents trending products
type TrendingProduct struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID     uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	TrendScore    float64   `json:"trend_score" gorm:"not null"`
	ViewCount     int       `json:"view_count" gorm:"default:0"`
	SalesCount    int       `json:"sales_count" gorm:"default:0"`
	SearchCount   int       `json:"search_count" gorm:"default:0"`
	Period        string    `json:"period" gorm:"not null"` // daily, weekly, monthly
	Date          time.Time `json:"date" gorm:"not null;index"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Product Product `json:"product" gorm:"foreignKey:ProductID"`
}

// TableName returns the table name for TrendingProduct entity
func (TrendingProduct) TableName() string {
	return "trending_products"
}

// RecommendationRequest represents a request for recommendations
type RecommendationRequest struct {
	UserID      *uuid.UUID           `json:"user_id"`
	SessionID   *string              `json:"session_id"`
	ProductID   *uuid.UUID           `json:"product_id"`
	Type        RecommendationType   `json:"type"`
	Limit       int                  `json:"limit"`
	Filters     *RecommendationFilters `json:"filters"`
	Context     map[string]interface{} `json:"context"`
}

// RecommendationFilters represents filters for recommendations
type RecommendationFilters struct {
	CategoryIDs []uuid.UUID `json:"category_ids"`
	BrandIDs    []uuid.UUID `json:"brand_ids"`
	MinPrice    *float64    `json:"min_price"`
	MaxPrice    *float64    `json:"max_price"`
	InStock     *bool       `json:"in_stock"`
	ExcludeIDs  []uuid.UUID `json:"exclude_ids"`
}

// RecommendationResponse represents a recommendation response
type RecommendationResponse struct {
	Type             RecommendationType `json:"type"`
	Products         []ProductListItem  `json:"products"`
	Reason           string             `json:"reason"`
	ConfidenceScore  float64            `json:"confidence_score"`
	Algorithm        string             `json:"algorithm"`
	TotalCount       int                `json:"total_count"`
}

// ProductListItem represents a simplified product for lists
type ProductListItem struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	Price            float64   `json:"price"`
	CurrentPrice     float64   `json:"current_price"`
	IsOnSale         bool      `json:"is_on_sale"`
	SaleDiscountPercentage float64 `json:"sale_discount_percentage"`
	MainImage        string    `json:"main_image"`
	Stock            int       `json:"stock"`
	StockStatus      string    `json:"stock_status"`
	IsAvailable      bool      `json:"is_available"`
	RatingAverage    float64   `json:"rating_average"`
	RatingCount      int       `json:"rating_count"`
	Category         *Category `json:"category,omitempty"`
	Brand            *Brand    `json:"brand,omitempty"`
}
