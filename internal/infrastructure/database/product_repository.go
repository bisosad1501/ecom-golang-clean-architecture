package database

import (
	"context"
	"strings"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) repositories.ProductRepository {
	return &productRepository{db: db}
}

// Create creates a new product
func (r *productRepository) Create(ctx context.Context, product *entities.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// GetByID retrieves a product by ID
func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images").
		Preload("Tags").
		Where("id = ?", id).
		First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

// GetBySKU retrieves a product by SKU
func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images").
		Preload("Tags").
		Where("sku = ?", sku).
		First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

// Update updates an existing product
func (r *productRepository) Update(ctx context.Context, product *entities.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// Delete deletes a product by ID
func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrProductNotFound
	}
	return nil
}

// List retrieves products with pagination
func (r *productRepository) List(ctx context.Context, limit, offset int) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images").
		Preload("Tags").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

// Search searches products based on criteria
func (r *productRepository) Search(ctx context.Context, params repositories.ProductSearchParams) ([]*entities.Product, error) {
	query := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images").
		Preload("Tags")

	// Apply filters
	if params.Query != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+params.Query+"%", "%"+params.Query+"%")
	}

	if params.CategoryID != nil {
		query = query.Where("category_id = ?", *params.CategoryID)
	}

	if params.MinPrice != nil {
		query = query.Where("price >= ?", *params.MinPrice)
	}

	if params.MaxPrice != nil {
		query = query.Where("price <= ?", *params.MaxPrice)
	}

	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	// Apply sorting
	orderBy := "created_at DESC"
	if params.SortBy != "" {
		direction := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			direction = "DESC"
		}
		orderBy = params.SortBy + " " + direction
	}
	query = query.Order(orderBy)

	// Apply pagination
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	var products []*entities.Product
	err := query.Find(&products).Error
	return products, err
}

// Count returns the total number of products
func (r *productRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Product{}).Count(&count).Error
	return count, err
}

// CountByCategory returns the number of products in a category
func (r *productRepository) CountByCategory(ctx context.Context, categoryID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("category_id = ?", categoryID).
		Count(&count).Error
	return count, err
}

// GetByCategory retrieves products by category
func (r *productRepository) GetByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images").
		Preload("Tags").
		Where("category_id = ?", categoryID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

// UpdateStock updates product stock
func (r *productRepository) UpdateStock(ctx context.Context, productID uuid.UUID, stock int) error {
	result := r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("id = ?", productID).
		Update("stock", stock)
	
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrProductNotFound
	}
	return nil
}

// ExistsBySKU checks if a product exists with the given SKU
func (r *productRepository) ExistsBySKU(ctx context.Context, sku string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("sku = ?", sku).
		Count(&count).Error
	return count > 0, err
}

// GetFeatured retrieves featured products
func (r *productRepository) GetFeatured(ctx context.Context, limit int) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images").
		Preload("Tags").
		Joins("JOIN product_tags ON products.id = product_tags.id").
		Where("product_tags.slug = ?", "featured").
		Limit(limit).
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

// GetRelated retrieves related products
func (r *productRepository) GetRelated(ctx context.Context, productID uuid.UUID, limit int) ([]*entities.Product, error) {
	// Get the product to find its category
	product, err := r.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	var products []*entities.Product
	err = r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images").
		Preload("Tags").
		Where("category_id = ? AND id != ?", product.CategoryID, productID).
		Limit(limit).
		Order("RANDOM()").
		Find(&products).Error
	return products, err
}
