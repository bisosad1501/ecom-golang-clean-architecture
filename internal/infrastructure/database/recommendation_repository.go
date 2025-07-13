package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
)

// recommendationRepository implements the RecommendationRepository interface
type recommendationRepository struct {
	db *gorm.DB
}

// NewRecommendationRepository creates a new recommendation repository
func NewRecommendationRepository(db *gorm.DB) repositories.RecommendationRepository {
	return &recommendationRepository{db: db}
}

// CreateRecommendation creates a new product recommendation
func (r *recommendationRepository) CreateRecommendation(ctx context.Context, recommendation *entities.ProductRecommendation) error {
	return r.db.WithContext(ctx).Create(recommendation).Error
}

// GetRecommendationsByProduct gets recommendations for a specific product
func (r *recommendationRepository) GetRecommendationsByProduct(ctx context.Context, productID uuid.UUID, recType entities.RecommendationType, limit int) ([]entities.ProductRecommendation, error) {
	var recommendations []entities.ProductRecommendation
	query := r.db.WithContext(ctx).
		Where("product_id = ? AND type = ? AND is_active = ?", productID, recType, true).
		Order("score DESC").
		Preload("Recommended").
		Preload("Recommended.Category").
		Preload("Recommended.Brand")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&recommendations).Error
	return recommendations, err
}

// GetRecommendationsByType gets recommendations by type
func (r *recommendationRepository) GetRecommendationsByType(ctx context.Context, recType entities.RecommendationType, limit int) ([]entities.ProductRecommendation, error) {
	var recommendations []entities.ProductRecommendation
	query := r.db.WithContext(ctx).
		Where("type = ? AND is_active = ?", recType, true).
		Order("score DESC").
		Preload("Product").
		Preload("Recommended")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&recommendations).Error
	return recommendations, err
}

// UpdateRecommendation updates a recommendation
func (r *recommendationRepository) UpdateRecommendation(ctx context.Context, recommendation *entities.ProductRecommendation) error {
	return r.db.WithContext(ctx).Save(recommendation).Error
}

// DeleteRecommendation deletes a recommendation
func (r *recommendationRepository) DeleteRecommendation(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.ProductRecommendation{}, id).Error
}

// BulkCreateRecommendations creates multiple recommendations
func (r *recommendationRepository) BulkCreateRecommendations(ctx context.Context, recommendations []entities.ProductRecommendation) error {
	if len(recommendations) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(recommendations, 100).Error
}

// CreateInteraction creates a user product interaction
func (r *recommendationRepository) CreateInteraction(ctx context.Context, interaction *entities.UserProductInteraction) error {
	return r.db.WithContext(ctx).Create(interaction).Error
}

// GetUserInteractions gets user interactions
func (r *recommendationRepository) GetUserInteractions(ctx context.Context, userID uuid.UUID, limit int) ([]entities.UserProductInteraction, error) {
	var interactions []entities.UserProductInteraction
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Preload("Product")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&interactions).Error
	return interactions, err
}

// GetSessionInteractions gets session interactions for guest users
func (r *recommendationRepository) GetSessionInteractions(ctx context.Context, sessionID string, limit int) ([]entities.UserProductInteraction, error) {
	var interactions []entities.UserProductInteraction
	query := r.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at DESC").
		Preload("Product")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&interactions).Error
	return interactions, err
}

// GetProductInteractions gets interactions for a specific product
func (r *recommendationRepository) GetProductInteractions(ctx context.Context, productID uuid.UUID, interactionType entities.InteractionType, limit int) ([]entities.UserProductInteraction, error) {
	var interactions []entities.UserProductInteraction
	query := r.db.WithContext(ctx).
		Where("product_id = ? AND interaction_type = ?", productID, interactionType).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&interactions).Error
	return interactions, err
}

// GetUserInteractionsByType gets user interactions by type
func (r *recommendationRepository) GetUserInteractionsByType(ctx context.Context, userID uuid.UUID, interactionType entities.InteractionType, limit int) ([]entities.UserProductInteraction, error) {
	var interactions []entities.UserProductInteraction
	query := r.db.WithContext(ctx).
		Where("user_id = ? AND interaction_type = ?", userID, interactionType).
		Order("created_at DESC").
		Preload("Product")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&interactions).Error
	return interactions, err
}

// CreateSimilarity creates a product similarity record
func (r *recommendationRepository) CreateSimilarity(ctx context.Context, similarity *entities.ProductSimilarity) error {
	return r.db.WithContext(ctx).Create(similarity).Error
}

// GetSimilarProducts gets similar products
func (r *recommendationRepository) GetSimilarProducts(ctx context.Context, productID uuid.UUID, limit int) ([]entities.ProductSimilarity, error) {
	var similarities []entities.ProductSimilarity
	query := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("similarity_score DESC").
		Preload("Similar").
		Preload("Similar.Category").
		Preload("Similar.Brand")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&similarities).Error
	return similarities, err
}

// UpdateSimilarity updates a similarity record
func (r *recommendationRepository) UpdateSimilarity(ctx context.Context, similarity *entities.ProductSimilarity) error {
	return r.db.WithContext(ctx).Save(similarity).Error
}

// BulkCreateSimilarities creates multiple similarity records
func (r *recommendationRepository) BulkCreateSimilarities(ctx context.Context, similarities []entities.ProductSimilarity) error {
	if len(similarities) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(similarities, 100).Error
}

// CreateFrequentlyBought creates a frequently bought together record
func (r *recommendationRepository) CreateFrequentlyBought(ctx context.Context, fbt *entities.FrequentlyBoughtTogether) error {
	return r.db.WithContext(ctx).Create(fbt).Error
}

// GetFrequentlyBoughtTogether gets frequently bought together products
func (r *recommendationRepository) GetFrequentlyBoughtTogether(ctx context.Context, productID uuid.UUID, limit int) ([]entities.FrequentlyBoughtTogether, error) {
	var fbts []entities.FrequentlyBoughtTogether
	query := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("confidence DESC, frequency DESC").
		Preload("With").
		Preload("With.Category").
		Preload("With.Brand")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&fbts).Error
	return fbts, err
}

// UpdateFrequentlyBought updates a frequently bought together record
func (r *recommendationRepository) UpdateFrequentlyBought(ctx context.Context, fbt *entities.FrequentlyBoughtTogether) error {
	return r.db.WithContext(ctx).Save(fbt).Error
}

// BulkCreateFrequentlyBought creates multiple frequently bought together records
func (r *recommendationRepository) BulkCreateFrequentlyBought(ctx context.Context, fbts []entities.FrequentlyBoughtTogether) error {
	if len(fbts) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(fbts, 100).Error
}

// CreateTrendingProduct creates a trending product record
func (r *recommendationRepository) CreateTrendingProduct(ctx context.Context, trending *entities.TrendingProduct) error {
	return r.db.WithContext(ctx).Create(trending).Error
}

// GetTrendingProducts gets trending products
func (r *recommendationRepository) GetTrendingProducts(ctx context.Context, period string, limit int) ([]entities.TrendingProduct, error) {
	var trending []entities.TrendingProduct
	query := r.db.WithContext(ctx).
		Where("period = ?", period).
		Order("trending_score DESC, updated_at DESC").
		Preload("Product").
		Preload("Product.Category").
		Preload("Product.Brand")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&trending).Error
	return trending, err
}

// UpdateTrendingProduct updates a trending product record
func (r *recommendationRepository) UpdateTrendingProduct(ctx context.Context, trending *entities.TrendingProduct) error {
	return r.db.WithContext(ctx).Save(trending).Error
}

// BulkCreateTrendingProducts creates multiple trending product records
func (r *recommendationRepository) BulkCreateTrendingProducts(ctx context.Context, trendings []entities.TrendingProduct) error {
	if len(trendings) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(trendings, 100).Error
}

// GetMostViewedProducts gets most viewed products
func (r *recommendationRepository) GetMostViewedProducts(ctx context.Context, days int, limit int) ([]entities.ProductListItem, error) {
	var queryResults []ProductQueryResult

	query := `
		SELECT p.id, p.name, p.slug, p.price,
			COALESCE(p.sale_price, p.price) as current_price,
			CASE WHEN p.sale_price IS NOT NULL AND p.sale_price < p.price THEN true ELSE false END as is_on_sale,
			CASE WHEN p.sale_price IS NOT NULL AND p.sale_price < p.price
				THEN ROUND(((p.price - p.sale_price) / p.price * 100)::numeric, 2)
				ELSE 0 END as sale_discount_percentage,
			COALESCE((SELECT url FROM product_images WHERE product_id = p.id ORDER BY position ASC LIMIT 1), '') as main_image,
			p.stock, p.stock_status,
			CASE WHEN p.stock > 0 OR p.allow_backorder = true THEN true ELSE false END as is_available,
			COALESCE(AVG(r.rating), 0) as rating_average,
			COUNT(r.id) as rating_count,
			(SELECT COUNT(*) FROM user_product_interactions upi2
			 WHERE upi2.product_id = p.id
				AND upi2.interaction_type = 'view'
				AND upi2.created_at >= NOW() - INTERVAL '%d days') as view_count
		FROM products p
		LEFT JOIN reviews r ON p.id = r.product_id AND r.status = 'approved'
		WHERE p.status = 'active'
			AND EXISTS (
				SELECT 1 FROM user_product_interactions upi
				WHERE upi.product_id = p.id
					AND upi.interaction_type = 'view'
					AND upi.created_at >= NOW() - INTERVAL '%d days'
			)
		GROUP BY p.id, p.name, p.slug, p.price, p.sale_price, p.stock, p.stock_status, p.allow_backorder
		ORDER BY view_count DESC
		LIMIT %d
	`

	err := r.db.WithContext(ctx).Raw(fmt.Sprintf(query, days, days, limit)).Scan(&queryResults).Error
	if err != nil {
		return nil, err
	}

	// Convert query results to ProductListItem
	products := make([]entities.ProductListItem, len(queryResults))
	for i, result := range queryResults {
		products[i] = result.ToProductListItem()
	}

	return products, nil
}

// GetMostPurchasedProducts gets most purchased products
func (r *recommendationRepository) GetMostPurchasedProducts(ctx context.Context, days int, limit int) ([]entities.ProductListItem, error) {
	var queryResults []ProductQueryResult

	query := `
		SELECT DISTINCT p.id, p.name, p.slug, p.price,
			COALESCE(p.sale_price, p.price) as current_price,
			CASE WHEN p.sale_price IS NOT NULL AND p.sale_price < p.price THEN true ELSE false END as is_on_sale,
			CASE WHEN p.sale_price IS NOT NULL AND p.sale_price < p.price
				THEN ROUND(((p.price - p.sale_price) / p.price * 100)::numeric, 2)
				ELSE 0 END as sale_discount_percentage,
			COALESCE((SELECT url FROM product_images WHERE product_id = p.id ORDER BY position ASC LIMIT 1), '') as main_image,
			p.stock, p.stock_status,
			CASE WHEN p.stock > 0 OR p.allow_backorder = true THEN true ELSE false END as is_available,
			COALESCE(AVG(r.rating), 0) as rating_average,
			COUNT(r.id) as rating_count
		FROM products p
		LEFT JOIN reviews r ON p.id = r.product_id AND r.status = 'approved'
		WHERE p.status = 'active'
			AND EXISTS (
				SELECT 1 FROM user_product_interactions upi
				WHERE upi.product_id = p.id
					AND upi.interaction_type = 'purchase'
					AND upi.created_at >= NOW() - INTERVAL '%d days'
			)
		GROUP BY p.id, p.name, p.slug, p.price, p.sale_price, p.stock, p.stock_status, p.allow_backorder
		ORDER BY (
			SELECT COUNT(*) FROM user_product_interactions upi2
			WHERE upi2.product_id = p.id
				AND upi2.interaction_type = 'purchase'
				AND upi2.created_at >= NOW() - INTERVAL '%d days'
		) DESC
		LIMIT %d
	`

	err := r.db.WithContext(ctx).Raw(fmt.Sprintf(query, days, days, limit)).Scan(&queryResults).Error
	if err != nil {
		return nil, err
	}

	// Convert query results to ProductListItem
	products := make([]entities.ProductListItem, len(queryResults))
	for i, result := range queryResults {
		products[i] = result.ToProductListItem()
	}

	return products, nil
}

// GetUserProductAffinities gets user product affinities based on interactions
func (r *recommendationRepository) GetUserProductAffinities(ctx context.Context, userID uuid.UUID, limit int) ([]entities.ProductListItem, error) {
	var queryResults []ProductQueryResult

	query := `
		SELECT DISTINCT p.id, p.name, p.slug, p.price,
			COALESCE(p.sale_price, p.price) as current_price,
			CASE WHEN p.sale_price IS NOT NULL AND p.sale_price < p.price THEN true ELSE false END as is_on_sale,
			CASE WHEN p.sale_price IS NOT NULL AND p.sale_price < p.price
				THEN ROUND(((p.price - p.sale_price) / p.price * 100)::numeric, 2)
				ELSE 0 END as sale_discount_percentage,
			COALESCE((SELECT url FROM product_images WHERE product_id = p.id ORDER BY position ASC LIMIT 1), '') as main_image,
			p.stock, p.stock_status,
			CASE WHEN p.stock > 0 OR p.allow_backorder = true THEN true ELSE false END as is_available,
			COALESCE(AVG(r.rating), 0) as rating_average,
			COUNT(r.id) as rating_count
		FROM products p
		LEFT JOIN reviews r ON p.id = r.product_id AND r.status = 'approved'
		WHERE p.status = 'active'
			AND EXISTS (
				SELECT 1 FROM user_product_interactions upi
				WHERE upi.product_id = p.id AND upi.user_id = $1
			)
		GROUP BY p.id, p.name, p.slug, p.price, p.sale_price, p.stock, p.stock_status, p.allow_backorder
		ORDER BY (
			SELECT SUM(upi2.value) FROM user_product_interactions upi2
			WHERE upi2.product_id = p.id AND upi2.user_id = $1
		) DESC
		LIMIT $2
	`

	err := r.db.WithContext(ctx).Raw(query, userID, limit).Scan(&queryResults).Error
	if err != nil {
		return nil, err
	}

	// Convert query results to ProductListItem
	products := make([]entities.ProductListItem, len(queryResults))
	for i, result := range queryResults {
		products[i] = result.ToProductListItem()
	}

	return products, nil
}

// GetCategoryAffinities gets user category affinities
func (r *recommendationRepository) GetCategoryAffinities(ctx context.Context, userID uuid.UUID, limit int) ([]entities.Category, error) {
	var categories []entities.Category

	query := `
		SELECT c.*, COUNT(upi.id) as interaction_count
		FROM categories c
		JOIN products p ON c.id = p.category_id
		JOIN user_product_interactions upi ON p.id = upi.product_id
		WHERE upi.user_id = $1
		GROUP BY c.id, c.name, c.slug, c.description, c.parent_id, c.is_active, c.sort_order, c.created_at, c.updated_at
		ORDER BY interaction_count DESC
		LIMIT $2
	`

	err := r.db.WithContext(ctx).Raw(query, userID, limit).Scan(&categories).Error
	return categories, err
}

// GetBrandAffinities gets user brand affinities
func (r *recommendationRepository) GetBrandAffinities(ctx context.Context, userID uuid.UUID, limit int) ([]entities.Brand, error) {
	var brands []entities.Brand

	query := `
		SELECT b.*, COUNT(upi.id) as interaction_count
		FROM brands b
		JOIN products p ON b.id = p.brand_id
		JOIN user_product_interactions upi ON p.id = upi.product_id
		WHERE upi.user_id = $1
		GROUP BY b.id, b.name, b.slug, b.description, b.logo_url, b.website_url, b.is_active, b.created_at, b.updated_at
		ORDER BY interaction_count DESC
		LIMIT $2
	`

	err := r.db.WithContext(ctx).Raw(query, userID, limit).Scan(&brands).Error
	return brands, err
}

// ProductQueryResult represents the result from raw SQL queries
type ProductQueryResult struct {
	ID                     uuid.UUID `json:"id"`
	Name                   string    `json:"name"`
	Slug                   string    `json:"slug"`
	Price                  float64   `json:"price"`
	CurrentPrice           float64   `json:"current_price"`
	IsOnSale               bool      `json:"is_on_sale"`
	SaleDiscountPercentage float64   `json:"sale_discount_percentage"`
	MainImage              string    `json:"main_image"`
	Stock                  int       `json:"stock"`
	StockStatus            string    `json:"stock_status"`
	IsAvailable            bool      `json:"is_available"`
	RatingAverage          float64   `json:"rating_average"`
	RatingCount            int       `json:"rating_count"`
}

// ToProductListItem converts ProductQueryResult to ProductListItem
func (p ProductQueryResult) ToProductListItem() entities.ProductListItem {
	return entities.ProductListItem{
		ID:                     p.ID,
		Name:                   p.Name,
		Slug:                   p.Slug,
		Price:                  p.Price,
		CurrentPrice:           p.CurrentPrice,
		IsOnSale:               p.IsOnSale,
		SaleDiscountPercentage: p.SaleDiscountPercentage,
		MainImage:              p.MainImage,
		Stock:                  p.Stock,
		StockStatus:            p.StockStatus,
		IsAvailable:            p.IsAvailable,
		RatingAverage:          p.RatingAverage,
		RatingCount:            p.RatingCount,
	}
}

// GenerateRelatedProducts generates related products based on category and brand similarity
func (r *recommendationRepository) GenerateRelatedProducts(ctx context.Context, productID uuid.UUID, limit int) ([]entities.ProductListItem, error) {
	var queryResults []ProductQueryResult

	query := `
		SELECT p2.id, p2.name, p2.slug, p2.price,
			COALESCE(p2.sale_price, p2.price) as current_price,
			CASE WHEN p2.sale_price IS NOT NULL AND p2.sale_price < p2.price THEN true ELSE false END as is_on_sale,
			CASE WHEN p2.sale_price IS NOT NULL AND p2.sale_price < p2.price
				THEN ROUND(((p2.price - p2.sale_price) / p2.price * 100)::numeric, 2)
				ELSE 0 END as sale_discount_percentage,
			COALESCE((SELECT url FROM product_images WHERE product_id = p2.id ORDER BY position ASC LIMIT 1), '') as main_image,
			p2.stock, p2.stock_status,
			CASE WHEN p2.stock > 0 OR p2.allow_backorder = true THEN true ELSE false END as is_available,
			COALESCE(AVG(r.rating), 0) as rating_average,
			COUNT(r.id) as rating_count,
			CASE WHEN p1.category_id = p2.category_id AND p1.brand_id = p2.brand_id THEN 3
				 WHEN p1.category_id = p2.category_id THEN 2
				 WHEN p1.brand_id = p2.brand_id THEN 1
				 ELSE 0 END as relevance_score
		FROM products p1
		JOIN products p2 ON (p1.category_id = p2.category_id OR p1.brand_id = p2.brand_id)
		LEFT JOIN reviews r ON p2.id = r.product_id AND r.status = 'approved'
		WHERE p1.id = $1 AND p2.id != $1 AND p2.status = 'active'
		GROUP BY p2.id, p2.name, p2.slug, p2.price, p2.sale_price, p2.stock, p2.stock_status, p2.allow_backorder, p1.category_id, p1.brand_id, p2.category_id, p2.brand_id
		ORDER BY relevance_score DESC, RANDOM()
		LIMIT $2
	`

	err := r.db.WithContext(ctx).Raw(query, productID, limit).Scan(&queryResults).Error
	if err != nil {
		return nil, err
	}

	// Convert query results to ProductListItem
	products := make([]entities.ProductListItem, len(queryResults))
	for i, result := range queryResults {
		products[i] = result.ToProductListItem()
	}

	return products, nil
}

// GenerateSimilarProducts generates similar products using content-based similarity
func (r *recommendationRepository) GenerateSimilarProducts(ctx context.Context, productID uuid.UUID, limit int) ([]entities.ProductListItem, error) {
	// First try to get from similarity table
	similarities, err := r.GetSimilarProducts(ctx, productID, limit)
	if err == nil && len(similarities) > 0 {
		products := make([]entities.ProductListItem, len(similarities))
		for i, sim := range similarities {
			products[i] = entities.ProductListItem{
				ID:                     sim.Similar.ID,
				Name:                   sim.Similar.Name,
				Slug:                   sim.Similar.Slug,
				Price:                  sim.Similar.Price,
				CurrentPrice:           sim.Similar.Price, // TODO: Calculate current price
				IsOnSale:               false,             // TODO: Calculate sale status
				SaleDiscountPercentage: 0,                 // TODO: Calculate discount
				Stock:                  sim.Similar.Stock,
				StockStatus:            string(sim.Similar.StockStatus),
				IsAvailable:            sim.Similar.Stock > 0 || sim.Similar.AllowBackorder,
			}
		}
		return products, nil
	}

	// Fallback to related products if no similarities found
	return r.GenerateRelatedProducts(ctx, productID, limit)
}

// GenerateFrequentlyBoughtTogether generates frequently bought together products
func (r *recommendationRepository) GenerateFrequentlyBoughtTogether(ctx context.Context, productID uuid.UUID, limit int) ([]entities.ProductListItem, error) {
	// First try to get from frequently bought together table
	fbts, err := r.GetFrequentlyBoughtTogether(ctx, productID, limit)
	if err == nil && len(fbts) > 0 {
		products := make([]entities.ProductListItem, len(fbts))
		for i, fbt := range fbts {
			products[i] = entities.ProductListItem{
				ID:                     fbt.With.ID,
				Name:                   fbt.With.Name,
				Slug:                   fbt.With.Slug,
				Price:                  fbt.With.Price,
				CurrentPrice:           fbt.With.Price, // TODO: Calculate current price
				IsOnSale:               false,          // TODO: Calculate sale status
				SaleDiscountPercentage: 0,              // TODO: Calculate discount
				Stock:                  fbt.With.Stock,
				StockStatus:            string(fbt.With.StockStatus),
				IsAvailable:            fbt.With.Stock > 0 || fbt.With.AllowBackorder,
			}
		}
		return products, nil
	}

	// Fallback: analyze order items to find frequently bought together
	var queryResults []ProductQueryResult

	query := `
		SELECT p2.id, p2.name, p2.slug, p2.price,
			COALESCE(p2.sale_price, p2.price) as current_price,
			CASE WHEN p2.sale_price IS NOT NULL AND p2.sale_price < p2.price THEN true ELSE false END as is_on_sale,
			CASE WHEN p2.sale_price IS NOT NULL AND p2.sale_price < p2.price
				THEN ROUND(((p2.price - p2.sale_price) / p2.price * 100)::numeric, 2)
				ELSE 0 END as sale_discount_percentage,
			COALESCE((SELECT url FROM product_images WHERE product_id = p2.id ORDER BY position ASC LIMIT 1), '') as main_image,
			p2.stock, p2.stock_status,
			CASE WHEN p2.stock > 0 OR p2.allow_backorder = true THEN true ELSE false END as is_available,
			COALESCE(AVG(r.rating), 0) as rating_average,
			COUNT(r.id) as rating_count,
			COUNT(DISTINCT oi1.order_id) as frequency
		FROM order_items oi1
		JOIN order_items oi2 ON oi1.order_id = oi2.order_id AND oi1.product_id != oi2.product_id
		JOIN products p2 ON oi2.product_id = p2.id
		LEFT JOIN reviews r ON p2.id = r.product_id AND r.status = 'approved'
		WHERE oi1.product_id = $1 AND p2.status = 'active'
		GROUP BY p2.id, p2.name, p2.slug, p2.price, p2.sale_price, p2.stock, p2.stock_status, p2.allow_backorder
		HAVING COUNT(DISTINCT oi1.order_id) >= 2
		ORDER BY frequency DESC, RANDOM()
		LIMIT $2
	`

	err = r.db.WithContext(ctx).Raw(query, productID, limit).Scan(&queryResults).Error
	if err != nil {
		return nil, err
	}

	// Convert query results to ProductListItem
	products := make([]entities.ProductListItem, len(queryResults))
	for i, result := range queryResults {
		products[i] = result.ToProductListItem()
	}

	return products, nil
}

// GeneratePersonalizedRecommendations generates personalized recommendations for a user
func (r *recommendationRepository) GeneratePersonalizedRecommendations(ctx context.Context, userID uuid.UUID, limit int) ([]entities.ProductListItem, error) {
	var queryResults []ProductQueryResult

	query := `
		WITH user_preferences AS (
			SELECT
				p.category_id,
				p.brand_id,
				SUM(upi.value) as affinity_score
			FROM user_product_interactions upi
			JOIN products p ON upi.product_id = p.id
			WHERE upi.user_id = $1
			GROUP BY p.category_id, p.brand_id
		),
		recommended_products AS (
			SELECT p.id, p.name, p.slug, p.price,
				COALESCE(p.sale_price, p.price) as current_price,
				CASE WHEN p.sale_price IS NOT NULL AND p.sale_price < p.price THEN true ELSE false END as is_on_sale,
				CASE WHEN p.sale_price IS NOT NULL AND p.sale_price < p.price
					THEN ROUND(((p.price - p.sale_price) / p.price * 100)::numeric, 2)
					ELSE 0 END as sale_discount_percentage,
				COALESCE((SELECT url FROM product_images WHERE product_id = p.id ORDER BY position ASC LIMIT 1), '') as main_image,
				p.stock, p.stock_status,
				CASE WHEN p.stock > 0 OR p.allow_backorder = true THEN true ELSE false END as is_available,
				COALESCE(AVG(r.rating), 0) as rating_average,
				COUNT(r.id) as rating_count,
				COALESCE(up.affinity_score, 0) as recommendation_score
			FROM products p
			LEFT JOIN user_preferences up ON (p.category_id = up.category_id OR p.brand_id = up.brand_id)
			LEFT JOIN reviews r ON p.id = r.product_id AND r.status = 'approved'
			WHERE p.status = 'active'
				AND p.id NOT IN (
					SELECT product_id FROM user_product_interactions
					WHERE user_id = $1 AND interaction_type IN ('purchase', 'view')
				)
			GROUP BY p.id, p.name, p.slug, p.price, p.sale_price, p.stock, p.stock_status, p.allow_backorder, up.affinity_score
		)
		SELECT * FROM recommended_products
		WHERE recommendation_score > 0
		ORDER BY recommendation_score DESC, RANDOM()
		LIMIT $2
	`

	err := r.db.WithContext(ctx).Raw(query, userID, limit).Scan(&queryResults).Error
	if err != nil {
		return nil, err
	}

	// Convert query results to ProductListItem
	products := make([]entities.ProductListItem, len(queryResults))
	for i, result := range queryResults {
		products[i] = result.ToProductListItem()
	}

	return products, nil
}

// GenerateTrendingRecommendations generates trending product recommendations
func (r *recommendationRepository) GenerateTrendingRecommendations(ctx context.Context, period string, limit int) ([]entities.ProductListItem, error) {
	// First try to get from trending table
	trending, err := r.GetTrendingProducts(ctx, period, limit)
	if err == nil && len(trending) > 0 {
		products := make([]entities.ProductListItem, len(trending))
		for i, t := range trending {
			products[i] = entities.ProductListItem{
				ID:                     t.Product.ID,
				Name:                   t.Product.Name,
				Slug:                   t.Product.Slug,
				Price:                  t.Product.Price,
				CurrentPrice:           t.Product.Price, // TODO: Calculate current price
				IsOnSale:               false,           // TODO: Calculate sale status
				SaleDiscountPercentage: 0,               // TODO: Calculate discount
				Stock:                  t.Product.Stock,
				StockStatus:            string(t.Product.StockStatus),
				IsAvailable:            t.Product.Stock > 0 || t.Product.AllowBackorder,
			}
		}
		return products, nil
	}

	// Fallback: calculate trending based on recent interactions
	var queryResults []ProductQueryResult

	days := 7 // Default to weekly
	if period == "daily" {
		days = 1
	} else if period == "monthly" {
		days = 30
	}

	query := `
		SELECT DISTINCT p.id, p.name, p.slug, p.price,
			COALESCE(p.sale_price, p.price) as current_price,
			CASE WHEN p.sale_price IS NOT NULL AND p.sale_price < p.price THEN true ELSE false END as is_on_sale,
			CASE WHEN p.sale_price IS NOT NULL AND p.sale_price < p.price
				THEN ROUND(((p.price - p.sale_price) / p.price * 100)::numeric, 2)
				ELSE 0 END as sale_discount_percentage,
			COALESCE((SELECT url FROM product_images WHERE product_id = p.id ORDER BY position ASC LIMIT 1), '') as main_image,
			p.stock, p.stock_status,
			CASE WHEN p.stock > 0 OR p.allow_backorder = true THEN true ELSE false END as is_available,
			COALESCE(AVG(r.rating), 0) as rating_average,
			COUNT(r.id) as rating_count,
			(
				SELECT SUM(upi.value) FROM user_product_interactions upi
				WHERE upi.product_id = p.id
					AND upi.created_at >= NOW() - INTERVAL '%d days'
			) as trend_score
		FROM products p
		LEFT JOIN reviews r ON p.id = r.product_id AND r.status = 'approved'
		WHERE p.status = 'active'
		GROUP BY p.id, p.name, p.slug, p.price, p.sale_price, p.stock, p.stock_status, p.allow_backorder
		HAVING (
			SELECT SUM(upi.value) FROM user_product_interactions upi
			WHERE upi.product_id = p.id
				AND upi.created_at >= NOW() - INTERVAL '%d days'
		) > 0
		ORDER BY trend_score DESC
		LIMIT %d
	`

	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(query, days, days, limit)).Scan(&queryResults).Error
	if err != nil {
		return nil, err
	}

	// Convert query results to ProductListItem
	products := make([]entities.ProductListItem, len(queryResults))
	for i, result := range queryResults {
		products[i] = result.ToProductListItem()
	}

	return products, nil
}

// GenerateCategoryBasedRecommendations generates recommendations based on category
func (r *recommendationRepository) GenerateCategoryBasedRecommendations(ctx context.Context, categoryID uuid.UUID, excludeProductID *uuid.UUID, limit int) ([]entities.ProductListItem, error) {
	var queryResults []ProductQueryResult

	query := `
		SELECT p.id, p.name, p.slug, p.price,
			COALESCE(p.sale_price, p.price) as current_price,
			CASE WHEN p.sale_price IS NOT NULL AND p.sale_price < p.price THEN true ELSE false END as is_on_sale,
			CASE WHEN p.sale_price IS NOT NULL AND p.sale_price < p.price
				THEN ROUND(((p.price - p.sale_price) / p.price * 100)::numeric, 2)
				ELSE 0 END as sale_discount_percentage,
			COALESCE((SELECT url FROM product_images WHERE product_id = p.id ORDER BY position ASC LIMIT 1), '') as main_image,
			p.stock, p.stock_status,
			CASE WHEN p.stock > 0 OR p.allow_backorder = true THEN true ELSE false END as is_available,
			COALESCE(AVG(r.rating), 0) as rating_average,
			COUNT(r.id) as rating_count
		FROM products p
		LEFT JOIN reviews r ON p.id = r.product_id AND r.status = 'approved'
		WHERE p.category_id = $1 AND p.status = 'active'
	`

	args := []interface{}{categoryID}
	if excludeProductID != nil {
		query += " AND p.id != $2"
		args = append(args, *excludeProductID)
	}

	query += `
		GROUP BY p.id, p.name, p.slug, p.price, p.sale_price, p.stock, p.stock_status, p.allow_backorder
		ORDER BY RANDOM()
		LIMIT $` + fmt.Sprintf("%d", len(args)+1)

	args = append(args, limit)

	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&queryResults).Error
	if err != nil {
		return nil, err
	}

	// Convert query results to ProductListItem
	products := make([]entities.ProductListItem, len(queryResults))
	for i, result := range queryResults {
		products[i] = result.ToProductListItem()
	}

	return products, nil
}

// GenerateBrandBasedRecommendations generates recommendations based on brand
func (r *recommendationRepository) GenerateBrandBasedRecommendations(ctx context.Context, brandID uuid.UUID, excludeProductID *uuid.UUID, limit int) ([]entities.ProductListItem, error) {
	var queryResults []ProductQueryResult

	query := `
		SELECT p.id, p.name, p.slug, p.price,
			COALESCE(p.sale_price, p.price) as current_price,
			CASE WHEN p.sale_price IS NOT NULL AND p.sale_price < p.price THEN true ELSE false END as is_on_sale,
			CASE WHEN p.sale_price IS NOT NULL AND p.sale_price < p.price
				THEN ROUND(((p.price - p.sale_price) / p.price * 100)::numeric, 2)
				ELSE 0 END as sale_discount_percentage,
			COALESCE((SELECT url FROM product_images WHERE product_id = p.id ORDER BY position ASC LIMIT 1), '') as main_image,
			p.stock, p.stock_status,
			CASE WHEN p.stock > 0 OR p.allow_backorder = true THEN true ELSE false END as is_available,
			COALESCE(AVG(r.rating), 0) as rating_average,
			COUNT(r.id) as rating_count
		FROM products p
		LEFT JOIN reviews r ON p.id = r.product_id AND r.status = 'approved'
		WHERE p.brand_id = $1 AND p.status = 'active'
	`

	args := []interface{}{brandID}
	if excludeProductID != nil {
		query += " AND p.id != $2"
		args = append(args, *excludeProductID)
	}

	query += `
		GROUP BY p.id, p.name, p.slug, p.price, p.sale_price, p.stock, p.stock_status, p.allow_backorder
		ORDER BY RANDOM()
		LIMIT $` + fmt.Sprintf("%d", len(args)+1)

	args = append(args, limit)

	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&queryResults).Error
	if err != nil {
		return nil, err
	}

	// Convert query results to ProductListItem
	products := make([]entities.ProductListItem, len(queryResults))
	for i, result := range queryResults {
		products[i] = result.ToProductListItem()
	}

	return products, nil
}

// Batch operations
func (r *recommendationRepository) BatchUpdateRecommendations(ctx context.Context, productID uuid.UUID) error {
	// This would implement batch update logic for recommendations
	// For now, return nil as this would be implemented as background jobs
	return nil
}

func (r *recommendationRepository) BatchUpdateSimilarities(ctx context.Context, productID uuid.UUID) error {
	// This would implement batch update logic for similarities
	// For now, return nil as this would be implemented as background jobs
	return nil
}

func (r *recommendationRepository) BatchUpdateFrequentlyBought(ctx context.Context) error {
	// This would implement batch update logic for frequently bought together
	// For now, return nil as this would be implemented as background jobs
	return nil
}

func (r *recommendationRepository) BatchUpdateTrending(ctx context.Context, period string) error {
	// This would implement batch update logic for trending products
	// For now, return nil as this would be implemented as background jobs
	return nil
}

// Cleanup operations
func (r *recommendationRepository) CleanupOldInteractions(ctx context.Context, days int) error {
	return r.db.WithContext(ctx).
		Where("created_at < ?", time.Now().AddDate(0, 0, -days)).
		Delete(&entities.UserProductInteraction{}).Error
}

func (r *recommendationRepository) CleanupOldTrending(ctx context.Context, days int) error {
	return r.db.WithContext(ctx).
		Where("date < ?", time.Now().AddDate(0, 0, -days)).
		Delete(&entities.TrendingProduct{}).Error
}
