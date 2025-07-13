package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
)

// RecommendationUseCase handles recommendation business logic
type RecommendationUseCase struct {
	recommendationRepo repositories.RecommendationRepository
	productRepo        repositories.ProductRepository
	userRepo           repositories.UserRepository
}

// NewRecommendationUseCase creates a new recommendation use case
func NewRecommendationUseCase(
	recommendationRepo repositories.RecommendationRepository,
	productRepo repositories.ProductRepository,
	userRepo repositories.UserRepository,
) *RecommendationUseCase {
	return &RecommendationUseCase{
		recommendationRepo: recommendationRepo,
		productRepo:        productRepo,
		userRepo:           userRepo,
	}
}

// GetRecommendations gets recommendations based on request
func (uc *RecommendationUseCase) GetRecommendations(ctx context.Context, req *entities.RecommendationRequest) (*entities.RecommendationResponse, error) {
	switch req.Type {
	case entities.RecommendationTypeRelated:
		return uc.getRelatedProducts(ctx, req)
	case entities.RecommendationTypeSimilar:
		return uc.getSimilarProducts(ctx, req)
	case entities.RecommendationTypeFrequentlyBought:
		return uc.getFrequentlyBoughtTogether(ctx, req)
	case entities.RecommendationTypePersonalized:
		return uc.getPersonalizedRecommendations(ctx, req)
	case entities.RecommendationTypeTrending:
		return uc.getTrendingRecommendations(ctx, req)
	case entities.RecommendationTypeBasedOnCategory:
		return uc.getCategoryBasedRecommendations(ctx, req)
	case entities.RecommendationTypeBasedOnBrand:
		return uc.getBrandBasedRecommendations(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported recommendation type: %s", req.Type)
	}
}

// TrackInteraction tracks user interaction with products
func (uc *RecommendationUseCase) TrackInteraction(ctx context.Context, interaction *entities.UserProductInteraction) error {
	// Set default values
	if interaction.Value == 0 {
		interaction.Value = uc.getInteractionWeight(interaction.InteractionType)
	}

	// Create interaction
	if err := uc.recommendationRepo.CreateInteraction(ctx, interaction); err != nil {
		return fmt.Errorf("failed to create interaction: %w", err)
	}

	// Trigger async recommendation updates if needed
	go uc.updateRecommendationsAsync(interaction.ProductID, interaction.UserID)

	return nil
}

// getRelatedProducts gets related products for a given product
func (uc *RecommendationUseCase) getRelatedProducts(ctx context.Context, req *entities.RecommendationRequest) (*entities.RecommendationResponse, error) {
	if req.ProductID == nil {
		return nil, fmt.Errorf("product_id is required for related products")
	}

	products, err := uc.recommendationRepo.GenerateRelatedProducts(ctx, *req.ProductID, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get related products: %w", err)
	}

	return &entities.RecommendationResponse{
		Type:            entities.RecommendationTypeRelated,
		Products:        products,
		Reason:          "Products related to your current selection",
		ConfidenceScore: 0.8,
		Algorithm:       "category_brand_similarity",
		TotalCount:      len(products),
	}, nil
}

// getSimilarProducts gets similar products using similarity algorithms
func (uc *RecommendationUseCase) getSimilarProducts(ctx context.Context, req *entities.RecommendationRequest) (*entities.RecommendationResponse, error) {
	if req.ProductID == nil {
		return nil, fmt.Errorf("product_id is required for similar products")
	}

	products, err := uc.recommendationRepo.GenerateSimilarProducts(ctx, *req.ProductID, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get similar products: %w", err)
	}

	return &entities.RecommendationResponse{
		Type:            entities.RecommendationTypeSimilar,
		Products:        products,
		Reason:          "Products similar to your current selection",
		ConfidenceScore: 0.75,
		Algorithm:       "content_based_similarity",
		TotalCount:      len(products),
	}, nil
}

// getFrequentlyBoughtTogether gets products frequently bought together
func (uc *RecommendationUseCase) getFrequentlyBoughtTogether(ctx context.Context, req *entities.RecommendationRequest) (*entities.RecommendationResponse, error) {
	if req.ProductID == nil {
		return nil, fmt.Errorf("product_id is required for frequently bought together")
	}

	products, err := uc.recommendationRepo.GenerateFrequentlyBoughtTogether(ctx, *req.ProductID, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get frequently bought together: %w", err)
	}

	return &entities.RecommendationResponse{
		Type:            entities.RecommendationTypeFrequentlyBought,
		Products:        products,
		Reason:          "Customers who bought this item also bought",
		ConfidenceScore: 0.85,
		Algorithm:       "market_basket_analysis",
		TotalCount:      len(products),
	}, nil
}

// getPersonalizedRecommendations gets personalized recommendations for a user
func (uc *RecommendationUseCase) getPersonalizedRecommendations(ctx context.Context, req *entities.RecommendationRequest) (*entities.RecommendationResponse, error) {
	if req.UserID == nil {
		return uc.getTrendingRecommendations(ctx, req) // Fallback to trending for anonymous users
	}

	products, err := uc.recommendationRepo.GeneratePersonalizedRecommendations(ctx, *req.UserID, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get personalized recommendations: %w", err)
	}

	return &entities.RecommendationResponse{
		Type:            entities.RecommendationTypePersonalized,
		Products:        products,
		Reason:          "Recommended for you based on your activity",
		ConfidenceScore: 0.9,
		Algorithm:       "collaborative_filtering",
		TotalCount:      len(products),
	}, nil
}

// getTrendingRecommendations gets trending products
func (uc *RecommendationUseCase) getTrendingRecommendations(ctx context.Context, req *entities.RecommendationRequest) (*entities.RecommendationResponse, error) {
	period := "weekly" // Default period
	if req.Context != nil {
		if p, ok := req.Context["period"].(string); ok {
			period = p
		}
	}

	products, err := uc.recommendationRepo.GenerateTrendingRecommendations(ctx, period, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending recommendations: %w", err)
	}

	return &entities.RecommendationResponse{
		Type:            entities.RecommendationTypeTrending,
		Products:        products,
		Reason:          "Trending products this week",
		ConfidenceScore: 0.7,
		Algorithm:       "trending_score",
		TotalCount:      len(products),
	}, nil
}

// getCategoryBasedRecommendations gets recommendations based on category
func (uc *RecommendationUseCase) getCategoryBasedRecommendations(ctx context.Context, req *entities.RecommendationRequest) (*entities.RecommendationResponse, error) {
	var categoryID uuid.UUID
	var excludeProductID *uuid.UUID

	if req.Context != nil {
		if cid, ok := req.Context["category_id"].(string); ok {
			if parsed, err := uuid.Parse(cid); err == nil {
				categoryID = parsed
			}
		}
	}

	if req.ProductID != nil {
		excludeProductID = req.ProductID
	}

	products, err := uc.recommendationRepo.GenerateCategoryBasedRecommendations(ctx, categoryID, excludeProductID, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get category-based recommendations: %w", err)
	}

	return &entities.RecommendationResponse{
		Type:            entities.RecommendationTypeBasedOnCategory,
		Products:        products,
		Reason:          "More products from this category",
		ConfidenceScore: 0.6,
		Algorithm:       "category_based",
		TotalCount:      len(products),
	}, nil
}

// getBrandBasedRecommendations gets recommendations based on brand
func (uc *RecommendationUseCase) getBrandBasedRecommendations(ctx context.Context, req *entities.RecommendationRequest) (*entities.RecommendationResponse, error) {
	var brandID uuid.UUID
	var excludeProductID *uuid.UUID

	if req.Context != nil {
		if bid, ok := req.Context["brand_id"].(string); ok {
			if parsed, err := uuid.Parse(bid); err == nil {
				brandID = parsed
			}
		}
	}

	if req.ProductID != nil {
		excludeProductID = req.ProductID
	}

	products, err := uc.recommendationRepo.GenerateBrandBasedRecommendations(ctx, brandID, excludeProductID, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get brand-based recommendations: %w", err)
	}

	return &entities.RecommendationResponse{
		Type:            entities.RecommendationTypeBasedOnBrand,
		Products:        products,
		Reason:          "More products from this brand",
		ConfidenceScore: 0.65,
		Algorithm:       "brand_based",
		TotalCount:      len(products),
	}, nil
}

// getInteractionWeight returns the weight for different interaction types
func (uc *RecommendationUseCase) getInteractionWeight(interactionType entities.InteractionType) float64 {
	weights := map[entities.InteractionType]float64{
		entities.InteractionTypeView:           1.0,
		entities.InteractionTypeClick:          1.5,
		entities.InteractionTypeAddToCart:      3.0,
		entities.InteractionTypeWishlist:       2.0,
		entities.InteractionTypePurchase:       5.0,
		entities.InteractionTypeReview:         4.0,
		entities.InteractionTypeShare:          2.5,
		entities.InteractionTypeCompare:        2.0,
		entities.InteractionTypeSearch:         1.2,
		entities.InteractionTypeRemoveFromCart: -1.0,
	}

	if weight, exists := weights[interactionType]; exists {
		return weight
	}
	return 1.0 // Default weight
}

// updateRecommendationsAsync triggers async recommendation updates
func (uc *RecommendationUseCase) updateRecommendationsAsync(productID uuid.UUID, userID *uuid.UUID) {
	// This would typically be handled by a background job queue
	// For now, we'll just log that an update should happen
	// In production, you'd use something like Redis Queue, Celery, etc.
}

// BatchUpdateRecommendations updates recommendations for all products
func (uc *RecommendationUseCase) BatchUpdateRecommendations(ctx context.Context) error {
	// This would typically be run as a scheduled job
	return uc.recommendationRepo.BatchUpdateFrequentlyBought(ctx)
}

// BatchUpdateTrending updates trending products
func (uc *RecommendationUseCase) BatchUpdateTrending(ctx context.Context, period string) error {
	return uc.recommendationRepo.BatchUpdateTrending(ctx, period)
}
