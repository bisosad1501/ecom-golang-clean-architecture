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

// GetCategoryTree returns all descendant category IDs for a given category (including itself)
func (r *categoryRepository) GetCategoryTree(ctx context.Context, categoryID uuid.UUID) ([]uuid.UUID, error) {
	var categoryIDs []uuid.UUID

	// Using recursive CTE query to get all descendants
	query := `
		WITH RECURSIVE category_tree AS (
			-- Base case: start with the given category
			SELECT id, parent_id, name
			FROM categories 
			WHERE id = $1 AND is_active = true

			UNION ALL

			-- Recursive case: find all children
			SELECT c.id, c.parent_id, c.name
			FROM categories c
			INNER JOIN category_tree ct ON c.parent_id = ct.id
			WHERE c.is_active = true
		)
		SELECT id FROM category_tree
	`

	rows, err := r.db.WithContext(ctx).Raw(query, categoryID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		categoryIDs = append(categoryIDs, id)
	}

	return categoryIDs, nil
}

// GetCategoryPath returns the full path from root to the given category
func (r *categoryRepository) GetCategoryPath(ctx context.Context, categoryID uuid.UUID) ([]*entities.Category, error) {
	var categories []*entities.Category

	// Using recursive CTE query to get path from root to category
	query := `
		WITH RECURSIVE category_path AS (
			-- Start with the target category
			SELECT id, parent_id, name, slug, sort_order, 0 as level
			FROM categories 
			WHERE id = $1 AND is_active = true

			UNION ALL

			-- Recursively find parent categories
			SELECT c.id, c.parent_id, c.name, c.slug, c.sort_order, cp.level + 1 as level
			FROM categories c
			INNER JOIN category_path cp ON c.id = cp.parent_id
			WHERE c.is_active = true
		)
		SELECT id, parent_id, name, slug, sort_order, level FROM category_path
		ORDER BY level DESC
	`

	rows, err := r.db.WithContext(ctx).Raw(query, categoryID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category entities.Category
		var level int
		if err := rows.Scan(&category.ID, &category.ParentID, &category.Name,
			&category.Slug, &category.SortOrder, &level); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

// GetProductCountByCategory returns product count for each category (including descendants)
func (r *categoryRepository) GetProductCountByCategory(ctx context.Context, categoryID uuid.UUID) (int64, error) {
	// Get all descendant categories
	categoryIDs, err := r.GetCategoryTree(ctx, categoryID)
	if err != nil {
		return 0, err
	}

	if len(categoryIDs) == 0 {
		return 0, nil
	}

	var count int64
	err = r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("category_id IN ? AND status = ?", categoryIDs, "active").
		Count(&count).Error

	return count, err
}
