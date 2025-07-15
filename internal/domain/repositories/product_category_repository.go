package repositories

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// ProductCategoryRepository defines the interface for product category operations
type ProductCategoryRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, productCategory *entities.ProductCategory) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ProductCategory, error)
	Update(ctx context.Context, productCategory *entities.ProductCategory) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters entities.ProductCategoryFilters) ([]*entities.ProductCategory, error)

	// Product-Category relationship operations
	GetCategoriesByProductID(ctx context.Context, productID uuid.UUID) ([]*entities.Category, error)
	GetProductsByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]*entities.Product, error)
	GetProductWithCategories(ctx context.Context, productID uuid.UUID) (*entities.ProductWithCategories, error)
	GetCategoryWithProducts(ctx context.Context, categoryID uuid.UUID) (*entities.CategoryWithProducts, error)

	// Assignment operations
	AssignProductToCategory(ctx context.Context, productID, categoryID uuid.UUID, isPrimary bool) error
	RemoveProductFromCategory(ctx context.Context, productID, categoryID uuid.UUID) error
	SetPrimaryCategory(ctx context.Context, productID, categoryID uuid.UUID) error
	GetPrimaryCategory(ctx context.Context, productID uuid.UUID) (*entities.Category, error)

	// Bulk operations
	AssignProductToCategories(ctx context.Context, productID uuid.UUID, categoryIDs []uuid.UUID, primaryCategoryID *uuid.UUID) error
	RemoveProductFromAllCategories(ctx context.Context, productID uuid.UUID) error
	GetProductsInMultipleCategories(ctx context.Context, categoryIDs []uuid.UUID) ([]*entities.Product, error)

	// Search and filtering
	SearchProductsByCategories(ctx context.Context, categoryIDs []uuid.UUID, includeSubcategories bool) ([]*entities.Product, error)
	GetProductsInCategoryHierarchy(ctx context.Context, categoryID uuid.UUID) ([]*entities.Product, error)

	// Statistics
	CountProductsInCategory(ctx context.Context, categoryID uuid.UUID) (int64, error)
	CountCategoriesForProduct(ctx context.Context, productID uuid.UUID) (int64, error)

	// Validation
	ExistsProductCategory(ctx context.Context, productID, categoryID uuid.UUID) (bool, error)
	ValidateProductCategoryAssignment(ctx context.Context, productID, categoryID uuid.UUID) error
}
