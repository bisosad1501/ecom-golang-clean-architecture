package services

import (
	"context"
	"fmt"
	"strings"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

// ProductCategoryService handles product categorization business logic
type ProductCategoryService interface {
	// Product categorization
	AssignProductToCategory(ctx context.Context, productID, categoryID uuid.UUID, isPrimary bool) error
	AssignProductToMultipleCategories(ctx context.Context, productID uuid.UUID, categoryIDs []uuid.UUID, primaryCategoryID *uuid.UUID) error
	RemoveProductFromCategory(ctx context.Context, productID, categoryID uuid.UUID) error
	SetPrimaryCategory(ctx context.Context, productID, categoryID uuid.UUID) error
	
	// Product queries
	GetProductWithCategories(ctx context.Context, productID uuid.UUID) (*entities.ProductWithCategories, error)
	GetProductsByCategory(ctx context.Context, categoryID uuid.UUID, includeSubcategories bool) ([]*entities.Product, error)
	GetProductsByMultipleCategories(ctx context.Context, categoryIDs []uuid.UUID) ([]*entities.Product, error)
	
	// Category queries
	GetCategoriesForProduct(ctx context.Context, productID uuid.UUID) ([]*entities.Category, error)
	GetPrimaryCategory(ctx context.Context, productID uuid.UUID) (*entities.Category, error)
	GetCategoryWithProducts(ctx context.Context, categoryID uuid.UUID) (*entities.CategoryWithProducts, error)
	
	// Search and filtering
	SearchProductsInCategories(ctx context.Context, categoryIDs []uuid.UUID, searchTerm string) ([]*entities.Product, error)
	GetFeaturedProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit int) ([]*entities.Product, error)
	
	// Migration and maintenance
	MigrateExistingProductCategories(ctx context.Context) error
	ValidateProductCategorization(ctx context.Context, productID uuid.UUID) error
}

type productCategoryService struct {
	productCategoryRepo repositories.ProductCategoryRepository
	productRepo         repositories.ProductRepository
	categoryRepo        repositories.CategoryRepository
}

// NewProductCategoryService creates a new product category service
func NewProductCategoryService(
	productCategoryRepo repositories.ProductCategoryRepository,
	productRepo repositories.ProductRepository,
	categoryRepo repositories.CategoryRepository,
) ProductCategoryService {
	return &productCategoryService{
		productCategoryRepo: productCategoryRepo,
		productRepo:         productRepo,
		categoryRepo:        categoryRepo,
	}
}

// AssignProductToCategory assigns a product to a category
func (s *productCategoryService) AssignProductToCategory(ctx context.Context, productID, categoryID uuid.UUID, isPrimary bool) error {
	// Validate the assignment
	if err := s.productCategoryRepo.ValidateProductCategoryAssignment(ctx, productID, categoryID); err != nil {
		return err
	}

	// Assign the product to category
	return s.productCategoryRepo.AssignProductToCategory(ctx, productID, categoryID, isPrimary)
}

// AssignProductToMultipleCategories assigns a product to multiple categories
func (s *productCategoryService) AssignProductToMultipleCategories(ctx context.Context, productID uuid.UUID, categoryIDs []uuid.UUID, primaryCategoryID *uuid.UUID) error {
	// Validate all assignments
	for _, categoryID := range categoryIDs {
		if err := s.productCategoryRepo.ValidateProductCategoryAssignment(ctx, productID, categoryID); err != nil {
			return fmt.Errorf("invalid category %s: %w", categoryID, err)
		}
	}

	// Validate primary category is in the list
	if primaryCategoryID != nil {
		found := false
		for _, categoryID := range categoryIDs {
			if categoryID == *primaryCategoryID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("primary category %s must be in the category list", *primaryCategoryID)
		}
	}

	// Assign to all categories
	return s.productCategoryRepo.AssignProductToCategories(ctx, productID, categoryIDs, primaryCategoryID)
}

// RemoveProductFromCategory removes a product from a category
func (s *productCategoryService) RemoveProductFromCategory(ctx context.Context, productID, categoryID uuid.UUID) error {
	return s.productCategoryRepo.RemoveProductFromCategory(ctx, productID, categoryID)
}

// SetPrimaryCategory sets a category as primary for a product
func (s *productCategoryService) SetPrimaryCategory(ctx context.Context, productID, categoryID uuid.UUID) error {
	return s.productCategoryRepo.SetPrimaryCategory(ctx, productID, categoryID)
}

// GetProductWithCategories gets a product with all its categories
func (s *productCategoryService) GetProductWithCategories(ctx context.Context, productID uuid.UUID) (*entities.ProductWithCategories, error) {
	return s.productCategoryRepo.GetProductWithCategories(ctx, productID)
}

// GetProductsByCategory gets all products in a category
func (s *productCategoryService) GetProductsByCategory(ctx context.Context, categoryID uuid.UUID, includeSubcategories bool) ([]*entities.Product, error) {
	if includeSubcategories {
		return s.productCategoryRepo.GetProductsInCategoryHierarchy(ctx, categoryID)
	}
	return s.productCategoryRepo.GetProductsByCategoryID(ctx, categoryID)
}

// GetProductsByMultipleCategories gets products that belong to multiple categories
func (s *productCategoryService) GetProductsByMultipleCategories(ctx context.Context, categoryIDs []uuid.UUID) ([]*entities.Product, error) {
	return s.productCategoryRepo.GetProductsInMultipleCategories(ctx, categoryIDs)
}

// GetCategoriesForProduct gets all categories for a product
func (s *productCategoryService) GetCategoriesForProduct(ctx context.Context, productID uuid.UUID) ([]*entities.Category, error) {
	return s.productCategoryRepo.GetCategoriesByProductID(ctx, productID)
}

// GetPrimaryCategory gets the primary category for a product
func (s *productCategoryService) GetPrimaryCategory(ctx context.Context, productID uuid.UUID) (*entities.Category, error) {
	return s.productCategoryRepo.GetPrimaryCategory(ctx, productID)
}

// GetCategoryWithProducts gets a category with all its products
func (s *productCategoryService) GetCategoryWithProducts(ctx context.Context, categoryID uuid.UUID) (*entities.CategoryWithProducts, error) {
	return s.productCategoryRepo.GetCategoryWithProducts(ctx, categoryID)
}

// SearchProductsInCategories searches products within specific categories
func (s *productCategoryService) SearchProductsInCategories(ctx context.Context, categoryIDs []uuid.UUID, searchTerm string) ([]*entities.Product, error) {
	// Get products in categories first
	products, err := s.productCategoryRepo.SearchProductsByCategories(ctx, categoryIDs, true)
	if err != nil {
		return nil, err
	}

	// Filter by search term if provided
	if searchTerm == "" {
		return products, nil
	}

	var filteredProducts []*entities.Product
	for _, product := range products {
		// Simple search in name and description
		if strings.Contains(strings.ToLower(product.Name), strings.ToLower(searchTerm)) ||
			strings.Contains(strings.ToLower(product.Description), strings.ToLower(searchTerm)) {
			filteredProducts = append(filteredProducts, product)
		}
	}

	return filteredProducts, nil
}

// GetFeaturedProductsByCategory gets featured products in a category
func (s *productCategoryService) GetFeaturedProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit int) ([]*entities.Product, error) {
	// Get all products in category
	products, err := s.productCategoryRepo.GetProductsByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// Filter featured products (simple logic: take first N products)
	var featuredProducts []*entities.Product
	for _, product := range products {
		if len(featuredProducts) < limit {
			featuredProducts = append(featuredProducts, product)
		} else {
			break
		}
	}

	return featuredProducts, nil
}

// MigrateExistingProductCategories migrates existing product categories from products.category_id
func (s *productCategoryService) MigrateExistingProductCategories(ctx context.Context) error {
	// Get all products with category_id
	products, err := s.productRepo.List(ctx, 0, 1000) // Get first 1000 products
	if err != nil {
		return fmt.Errorf("failed to get products: %w", err)
	}

	migratedCount := 0
	for _, product := range products {
		// Check if product has a valid category ID (not zero UUID)
		if product.CategoryID != uuid.Nil {
			// Check if already migrated
			exists, err := s.productCategoryRepo.ExistsProductCategory(ctx, product.ID, product.CategoryID)
			if err != nil {
				fmt.Printf("Warning: Failed to check existing assignment for product %s: %v\n", product.ID, err)
				continue
			}

			if !exists {
				// Create the assignment
				err = s.productCategoryRepo.AssignProductToCategory(ctx, product.ID, product.CategoryID, true)
				if err != nil {
					fmt.Printf("Warning: Failed to migrate product %s to category %s: %v\n", product.ID, product.CategoryID, err)
					continue
				}
				migratedCount++
			}
		}
	}

	fmt.Printf("✅ Migrated %d products to new category system\n", migratedCount)
	return nil
}

// ValidateProductCategorization validates product categorization
func (s *productCategoryService) ValidateProductCategorization(ctx context.Context, productID uuid.UUID) error {
	// Get product categories
	categories, err := s.productCategoryRepo.GetCategoriesByProductID(ctx, productID)
	if err != nil {
		return err
	}

	if len(categories) == 0 {
		return fmt.Errorf("product %s has no categories assigned", productID)
	}

	// Check if primary category exists
	primaryCategory, err := s.productCategoryRepo.GetPrimaryCategory(ctx, productID)
	if err != nil {
		return fmt.Errorf("product %s has no primary category: %w", productID, err)
	}

	if primaryCategory == nil {
		return fmt.Errorf("product %s has no primary category", productID)
	}

	return nil
}
