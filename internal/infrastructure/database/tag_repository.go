package database

import (
	"context"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *gorm.DB) repositories.TagRepository {
	return &tagRepository{db: db}
}

// Create creates a new tag
func (r *tagRepository) Create(ctx context.Context, tag *entities.ProductTag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

// GetByID retrieves a tag by ID
func (r *tagRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ProductTag, error) {
	var tag entities.ProductTag
	if err := r.db.WithContext(ctx).First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetByName retrieves a tag by name
func (r *tagRepository) GetByName(ctx context.Context, name string) (*entities.ProductTag, error) {
	var tag entities.ProductTag
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetBySlug retrieves a tag by slug
func (r *tagRepository) GetBySlug(ctx context.Context, slug string) (*entities.ProductTag, error) {
	var tag entities.ProductTag
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// Update updates an existing tag
func (r *tagRepository) Update(ctx context.Context, tag *entities.ProductTag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

// Delete deletes a tag by ID
func (r *tagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.ProductTag{}, id).Error
}

// List retrieves tags with pagination
func (r *tagRepository) List(ctx context.Context, limit, offset int) ([]*entities.ProductTag, error) {
	var tags []*entities.ProductTag
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// ExistsByName checks if a tag exists with the given name
func (r *tagRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.ProductTag{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsBySlug checks if a tag exists with the given slug
func (r *tagRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.ProductTag{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindOrCreate finds existing tag by name or creates new one
func (r *tagRepository) FindOrCreate(ctx context.Context, name string) (*entities.ProductTag, error) {
	// Try to find existing tag by name first
	tag, err := r.GetByName(ctx, name)
	if err == nil {
		return tag, nil
	}
	
	// If not found, create new tag
	if err == gorm.ErrRecordNotFound {
		slug := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "-"))
		
		newTag := &entities.ProductTag{
			ID:        uuid.New(),
			Name:      strings.TrimSpace(name),
			Slug:      slug,
			CreatedAt: time.Now(),
		}
		
		if err := r.Create(ctx, newTag); err != nil {
			return nil, err
		}
		
		return newTag, nil
	}
	
	return nil, err
}
