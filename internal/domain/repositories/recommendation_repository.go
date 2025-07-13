package repositories

import (
	"context"

	"github.com/google/uuid"
	"ecom-golang-clean-architecture/internal/domain/entities"
)

// RecommendationRepository defines the interface for recommendation data operations
type RecommendationRepository interface {
	// Product Recommendations
	CreateRecommendation(ctx context.Context, recommendation *entities.ProductRecommendation) error
	GetRecommendationsByProduct(ctx context.Context, productID uuid.UUID, recType entities.RecommendationType, limit int) ([]entities.ProductRecommendation, error)
	GetRecommendationsByType(ctx context.Context, recType entities.RecommendationType, limit int) ([]entities.ProductRecommendation, error)
	UpdateRecommendation(ctx context.Context, recommendation *entities.ProductRecommendation) error
	DeleteRecommendation(ctx context.Context, id uuid.UUID) error
	BulkCreateRecommendations(ctx context.Context, recommendations []entities.ProductRecommendation) error
	
	// User Product Interactions
	CreateInteraction(ctx context.Context, interaction *entities.UserProductInteraction) error
	GetUserInteractions(ctx context.Context, userID uuid.UUID, limit int) ([]entities.UserProductInteraction, error)
	GetSessionInteractions(ctx context.Context, sessionID string, limit int) ([]entities.UserProductInteraction, error)
	GetProductInteractions(ctx context.Context, productID uuid.UUID, interactionType entities.InteractionType, limit int) ([]entities.UserProductInteraction, error)
	GetUserInteractionsByType(ctx context.Context, userID uuid.UUID, interactionType entities.InteractionType, limit int) ([]entities.UserProductInteraction, error)
	
	// Product Similarity
	CreateSimilarity(ctx context.Context, similarity *entities.ProductSimilarity) error
	GetSimilarProducts(ctx context.Context, productID uuid.UUID, limit int) ([]entities.ProductSimilarity, error)
	UpdateSimilarity(ctx context.Context, similarity *entities.ProductSimilarity) error
	BulkCreateSimilarities(ctx context.Context, similarities []entities.ProductSimilarity) error
	
	// Frequently Bought Together
	CreateFrequentlyBought(ctx context.Context, fbt *entities.FrequentlyBoughtTogether) error
	GetFrequentlyBoughtTogether(ctx context.Context, productID uuid.UUID, limit int) ([]entities.FrequentlyBoughtTogether, error)
	UpdateFrequentlyBought(ctx context.Context, fbt *entities.FrequentlyBoughtTogether) error
	BulkCreateFrequentlyBought(ctx context.Context, fbts []entities.FrequentlyBoughtTogether) error
	
	// Trending Products
	CreateTrendingProduct(ctx context.Context, trending *entities.TrendingProduct) error
	GetTrendingProducts(ctx context.Context, period string, limit int) ([]entities.TrendingProduct, error)
	UpdateTrendingProduct(ctx context.Context, trending *entities.TrendingProduct) error
	BulkCreateTrendingProducts(ctx context.Context, trendings []entities.TrendingProduct) error
	
	// Analytics and Insights
	GetMostViewedProducts(ctx context.Context, days int, limit int) ([]entities.ProductListItem, error)
	GetMostPurchasedProducts(ctx context.Context, days int, limit int) ([]entities.ProductListItem, error)
	GetUserProductAffinities(ctx context.Context, userID uuid.UUID, limit int) ([]entities.ProductListItem, error)
	GetCategoryAffinities(ctx context.Context, userID uuid.UUID, limit int) ([]entities.Category, error)
	GetBrandAffinities(ctx context.Context, userID uuid.UUID, limit int) ([]entities.Brand, error)
	
	// Recommendation Generation
	GenerateRelatedProducts(ctx context.Context, productID uuid.UUID, limit int) ([]entities.ProductListItem, error)
	GenerateSimilarProducts(ctx context.Context, productID uuid.UUID, limit int) ([]entities.ProductListItem, error)
	GenerateFrequentlyBoughtTogether(ctx context.Context, productID uuid.UUID, limit int) ([]entities.ProductListItem, error)
	GeneratePersonalizedRecommendations(ctx context.Context, userID uuid.UUID, limit int) ([]entities.ProductListItem, error)
	GenerateTrendingRecommendations(ctx context.Context, period string, limit int) ([]entities.ProductListItem, error)
	GenerateCategoryBasedRecommendations(ctx context.Context, categoryID uuid.UUID, excludeProductID *uuid.UUID, limit int) ([]entities.ProductListItem, error)
	GenerateBrandBasedRecommendations(ctx context.Context, brandID uuid.UUID, excludeProductID *uuid.UUID, limit int) ([]entities.ProductListItem, error)
	
	// Batch Operations
	BatchUpdateRecommendations(ctx context.Context, productID uuid.UUID) error
	BatchUpdateSimilarities(ctx context.Context, productID uuid.UUID) error
	BatchUpdateFrequentlyBought(ctx context.Context) error
	BatchUpdateTrending(ctx context.Context, period string) error
	
	// Cleanup
	CleanupOldInteractions(ctx context.Context, days int) error
	CleanupOldTrending(ctx context.Context, days int) error
}
