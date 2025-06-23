package database

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) repositories.CategoryRepository {
	return &categoryRepository{db: db}
}

// Create creates a new category
func (r *categoryRepository) Create(ctx context.Context, category *entities.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// GetByID retrieves a category by ID
func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Category, error) {
	var category entities.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("id = ?", id).
		First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

// GetBySlug retrieves a category by slug
func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*entities.Category, error) {
	var category entities.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("slug = ?", slug).
		First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

// Update updates an existing category
func (r *categoryRepository) Update(ctx context.Context, category *entities.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// Delete deletes a category by ID
func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrCategoryNotFound
	}
	return nil
}

// List retrieves categories with pagination
func (r *categoryRepository) List(ctx context.Context, limit, offset int) ([]*entities.Category, error) {
	var categories []*entities.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Limit(limit).
		Offset(offset).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// GetRootCategories retrieves root categories
func (r *categoryRepository) GetRootCategories(ctx context.Context) ([]*entities.Category, error) {
	var categories []*entities.Category
	err := r.db.WithContext(ctx).
		Preload("Children").
		Where("parent_id IS NULL AND is_active = ?", true).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// GetChildren retrieves child categories
func (r *categoryRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entities.Category, error) {
	var categories []*entities.Category
	err := r.db.WithContext(ctx).
		Where("parent_id = ? AND is_active = ?", parentID, true).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// Count returns the total number of categories
func (r *categoryRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Category{}).Count(&count).Error
	return count, err
}

// ExistsBySlug checks if a category exists with the given slug
func (r *categoryRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Category{}).
		Where("slug = ?", slug).
		Count(&count).Error
	return count > 0, err
}

// GetTree retrieves the category tree
func (r *categoryRepository) GetTree(ctx context.Context) ([]*entities.Category, error) {
	var categories []*entities.Category
	err := r.db.WithContext(ctx).
		Preload("Children").
		Where("parent_id IS NULL AND is_active = ?", true).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}
