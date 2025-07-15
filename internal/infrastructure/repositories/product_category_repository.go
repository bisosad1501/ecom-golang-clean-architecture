package repositories

import (
	"context"
	"fmt"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productCategoryRepository struct {
	db *gorm.DB
}

// NewProductCategoryRepository creates a new product category repository
func NewProductCategoryRepository(db *gorm.DB) repositories.ProductCategoryRepository {
	return &productCategoryRepository{db: db}
}

// Create creates a new product category relationship
func (r *productCategoryRepository) Create(ctx context.Context, productCategory *entities.ProductCategory) error {
	return r.db.WithContext(ctx).Create(productCategory).Error
}

// GetByID gets a product category by ID
func (r *productCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ProductCategory, error) {
	var productCategory entities.ProductCategory
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Category").
		First(&productCategory, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &productCategory, nil
}

// Update updates a product category
func (r *productCategoryRepository) Update(ctx context.Context, productCategory *entities.ProductCategory) error {
	return r.db.WithContext(ctx).Save(productCategory).Error
}

// Delete deletes a product category
func (r *productCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.ProductCategory{}, "id = ?", id).Error
}

// List lists product categories with filters
func (r *productCategoryRepository) List(ctx context.Context, filters entities.ProductCategoryFilters) ([]*entities.ProductCategory, error) {
	query := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Category")

	if filters.ProductID != nil {
		query = query.Where("product_id = ?", *filters.ProductID)
	}
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}
	if filters.IsPrimary != nil {
		query = query.Where("is_primary = ?", *filters.IsPrimary)
	}

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var productCategories []*entities.ProductCategory
	err := query.Find(&productCategories).Error
	return productCategories, err
}

// GetCategoriesByProductID gets all categories for a product
func (r *productCategoryRepository) GetCategoriesByProductID(ctx context.Context, productID uuid.UUID) ([]*entities.Category, error) {
	var categories []*entities.Category
	err := r.db.WithContext(ctx).
		Table("categories").
		Joins("JOIN product_categories ON categories.id = product_categories.category_id").
		Where("product_categories.product_id = ?", productID).
		Find(&categories).Error
	return categories, err
}

// GetProductsByCategoryID gets all products in a category
func (r *productCategoryRepository) GetProductsByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.db.WithContext(ctx).
		Table("products").
		Joins("JOIN product_categories ON products.id = product_categories.product_id").
		Where("product_categories.category_id = ?", categoryID).
		Find(&products).Error
	return products, err
}

// GetProductWithCategories gets a product with all its categories
func (r *productCategoryRepository) GetProductWithCategories(ctx context.Context, productID uuid.UUID) (*entities.ProductWithCategories, error) {
	// Get product
	var product entities.Product
	err := r.db.WithContext(ctx).First(&product, "id = ?", productID).Error
	if err != nil {
		return nil, err
	}

	// Get categories
	categories, err := r.GetCategoriesByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Get primary category
	primaryCategory, _ := r.GetPrimaryCategory(ctx, productID)

	// Build result
	result := &entities.ProductWithCategories{
		Product:    &product,
		Categories: categories,
		PrimaryCategory: primaryCategory,
	}

	// Extract category IDs
	for _, cat := range categories {
		result.CategoryIDs = append(result.CategoryIDs, cat.ID)
	}

	if primaryCategory != nil {
		result.PrimaryCategoryID = &primaryCategory.ID
	}

	return result, nil
}

// GetCategoryWithProducts gets a category with all its products
func (r *productCategoryRepository) GetCategoryWithProducts(ctx context.Context, categoryID uuid.UUID) (*entities.CategoryWithProducts, error) {
	// Get category
	var category entities.Category
	err := r.db.WithContext(ctx).First(&category, "id = ?", categoryID).Error
	if err != nil {
		return nil, err
	}

	// Get products
	products, err := r.GetProductsByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// Build result
	result := &entities.CategoryWithProducts{
		Category: &category,
		Products: products,
	}

	// Extract product IDs
	for _, prod := range products {
		result.ProductIDs = append(result.ProductIDs, prod.ID)
	}

	return result, nil
}

// AssignProductToCategory assigns a product to a category
func (r *productCategoryRepository) AssignProductToCategory(ctx context.Context, productID, categoryID uuid.UUID, isPrimary bool) error {
	// Check if relationship already exists
	exists, err := r.ExistsProductCategory(ctx, productID, categoryID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("product %s is already assigned to category %s", productID, categoryID)
	}

	// If this is primary, unset other primary categories
	if isPrimary {
		err := r.db.WithContext(ctx).
			Model(&entities.ProductCategory{}).
			Where("product_id = ? AND is_primary = true", productID).
			Update("is_primary", false).Error
		if err != nil {
			return err
		}
	}

	// Create new relationship
	productCategory := &entities.ProductCategory{
		ID:         uuid.New(),
		ProductID:  productID,
		CategoryID: categoryID,
		IsPrimary:  isPrimary,
	}

	return r.Create(ctx, productCategory)
}

// RemoveProductFromCategory removes a product from a category
func (r *productCategoryRepository) RemoveProductFromCategory(ctx context.Context, productID, categoryID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&entities.ProductCategory{}, "product_id = ? AND category_id = ?", productID, categoryID).Error
}

// SetPrimaryCategory sets a category as primary for a product
func (r *productCategoryRepository) SetPrimaryCategory(ctx context.Context, productID, categoryID uuid.UUID) error {
	// Verify the relationship exists
	exists, err := r.ExistsProductCategory(ctx, productID, categoryID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("product %s is not assigned to category %s", productID, categoryID)
	}

	// Unset all primary categories for this product
	err = r.db.WithContext(ctx).
		Model(&entities.ProductCategory{}).
		Where("product_id = ?", productID).
		Update("is_primary", false).Error
	if err != nil {
		return err
	}

	// Set the specified category as primary
	return r.db.WithContext(ctx).
		Model(&entities.ProductCategory{}).
		Where("product_id = ? AND category_id = ?", productID, categoryID).
		Update("is_primary", true).Error
}

// GetPrimaryCategory gets the primary category for a product
func (r *productCategoryRepository) GetPrimaryCategory(ctx context.Context, productID uuid.UUID) (*entities.Category, error) {
	var category entities.Category
	err := r.db.WithContext(ctx).
		Table("categories").
		Joins("JOIN product_categories ON categories.id = product_categories.category_id").
		Where("product_categories.product_id = ? AND product_categories.is_primary = true", productID).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// ExistsProductCategory checks if a product-category relationship exists
func (r *productCategoryRepository) ExistsProductCategory(ctx context.Context, productID, categoryID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.ProductCategory{}).
		Where("product_id = ? AND category_id = ?", productID, categoryID).
		Count(&count).Error
	return count > 0, err
}

// ValidateProductCategoryAssignment validates a product-category assignment
func (r *productCategoryRepository) ValidateProductCategoryAssignment(ctx context.Context, productID, categoryID uuid.UUID) error {
	// Check if product exists
	var productCount int64
	err := r.db.WithContext(ctx).Model(&entities.Product{}).Where("id = ?", productID).Count(&productCount).Error
	if err != nil {
		return err
	}
	if productCount == 0 {
		return fmt.Errorf("product %s does not exist", productID)
	}

	// Check if category exists
	var categoryCount int64
	err = r.db.WithContext(ctx).Model(&entities.Category{}).Where("id = ?", categoryID).Count(&categoryCount).Error
	if err != nil {
		return err
	}
	if categoryCount == 0 {
		return fmt.Errorf("category %s does not exist", categoryID)
	}

	return nil
}

// AssignProductToCategories assigns a product to multiple categories
func (r *productCategoryRepository) AssignProductToCategories(ctx context.Context, productID uuid.UUID, categoryIDs []uuid.UUID, primaryCategoryID *uuid.UUID) error {
	// Remove existing assignments
	err := r.RemoveProductFromAllCategories(ctx, productID)
	if err != nil {
		return err
	}

	// Add new assignments
	for _, categoryID := range categoryIDs {
		isPrimary := primaryCategoryID != nil && *primaryCategoryID == categoryID
		err := r.AssignProductToCategory(ctx, productID, categoryID, isPrimary)
		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveProductFromAllCategories removes a product from all categories
func (r *productCategoryRepository) RemoveProductFromAllCategories(ctx context.Context, productID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&entities.ProductCategory{}, "product_id = ?", productID).Error
}

// GetProductsInMultipleCategories gets products that belong to multiple categories
func (r *productCategoryRepository) GetProductsInMultipleCategories(ctx context.Context, categoryIDs []uuid.UUID) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.db.WithContext(ctx).
		Table("products").
		Joins("JOIN product_categories ON products.id = product_categories.product_id").
		Where("product_categories.category_id IN ?", categoryIDs).
		Group("products.id").
		Find(&products).Error
	return products, err
}

// SearchProductsByCategories searches products by categories
func (r *productCategoryRepository) SearchProductsByCategories(ctx context.Context, categoryIDs []uuid.UUID, includeSubcategories bool) ([]*entities.Product, error) {
	query := r.db.WithContext(ctx).
		Table("products").
		Joins("JOIN product_categories ON products.id = product_categories.product_id")

	if includeSubcategories {
		// Include subcategories using recursive CTE
		query = query.Where(`product_categories.category_id IN (
			WITH RECURSIVE category_tree AS (
				SELECT id FROM categories WHERE id IN ?
				UNION ALL
				SELECT c.id FROM categories c
				INNER JOIN category_tree ct ON c.parent_id = ct.id
			)
			SELECT id FROM category_tree
		)`, categoryIDs)
	} else {
		query = query.Where("product_categories.category_id IN ?", categoryIDs)
	}

	var products []*entities.Product
	err := query.Group("products.id").Find(&products).Error
	return products, err
}

// GetProductsInCategoryHierarchy gets all products in a category and its subcategories
func (r *productCategoryRepository) GetProductsInCategoryHierarchy(ctx context.Context, categoryID uuid.UUID) ([]*entities.Product, error) {
	return r.SearchProductsByCategories(ctx, []uuid.UUID{categoryID}, true)
}

// CountProductsInCategory counts products in a category
func (r *productCategoryRepository) CountProductsInCategory(ctx context.Context, categoryID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.ProductCategory{}).
		Where("category_id = ?", categoryID).
		Count(&count).Error
	return count, err
}

// CountCategoriesForProduct counts categories for a product
func (r *productCategoryRepository) CountCategoriesForProduct(ctx context.Context, productID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.ProductCategory{}).
		Where("product_id = ?", productID).
		Count(&count).Error
	return count, err
}
