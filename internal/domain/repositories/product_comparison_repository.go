package repositories

import (
	"context"
	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// ProductComparisonRepository defines the interface for product comparison operations
type ProductComparisonRepository interface {
	// Comparison management
	CreateComparison(ctx context.Context, comparison *entities.ProductComparison) error
	GetComparison(ctx context.Context, id uuid.UUID) (*entities.ProductComparison, error)
	GetComparisonByUserID(ctx context.Context, userID uuid.UUID) (*entities.ProductComparison, error)
	GetComparisonBySessionID(ctx context.Context, sessionID string) (*entities.ProductComparison, error)
	UpdateComparison(ctx context.Context, comparison *entities.ProductComparison) error
	DeleteComparison(ctx context.Context, id uuid.UUID) error

	// Comparison items management
	AddProductToComparison(ctx context.Context, comparisonID, productID uuid.UUID, position int) error
	RemoveProductFromComparison(ctx context.Context, comparisonID, productID uuid.UUID) error
	GetComparisonItems(ctx context.Context, comparisonID uuid.UUID) ([]entities.ProductComparisonItem, error)
	UpdateItemPosition(ctx context.Context, itemID uuid.UUID, position int) error
	ClearComparison(ctx context.Context, comparisonID uuid.UUID) error

	// Comparison queries
	GetComparisonWithProducts(ctx context.Context, id uuid.UUID) (*entities.ProductComparison, error)
	GetUserComparisons(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entities.ProductComparison, error)
	CountComparisonItems(ctx context.Context, comparisonID uuid.UUID) (int64, error)
	IsProductInComparison(ctx context.Context, comparisonID, productID uuid.UUID) (bool, error)

	// Comparison analytics
	GetPopularComparedProducts(ctx context.Context, limit int) ([]entities.Product, error)
	GetComparisonStats(ctx context.Context) (map[string]interface{}, error)
}
