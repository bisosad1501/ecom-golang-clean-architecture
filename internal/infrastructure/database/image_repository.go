package database

import (
	"context"
	"fmt"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type imageRepository struct {
	db *gorm.DB
}

// NewImageRepository creates a new image repository
func NewImageRepository(db *gorm.DB) repositories.ImageRepository {
	return &imageRepository{db: db}
}

// Create creates a new product image
func (r *imageRepository) Create(ctx context.Context, image *entities.ProductImage) error {
	return r.db.WithContext(ctx).Create(image).Error
}

// GetByID retrieves an image by ID
func (r *imageRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ProductImage, error) {
	var image entities.ProductImage
	if err := r.db.WithContext(ctx).First(&image, id).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

// GetByProductID retrieves active images by product ID (position >= 0)
func (r *imageRepository) GetByProductID(ctx context.Context, productID uuid.UUID) ([]*entities.ProductImage, error) {
	var images []*entities.ProductImage
	if err := r.db.WithContext(ctx).Where("product_id = ? AND position >= 0", productID).Order("position").Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

// Update updates an existing image
func (r *imageRepository) Update(ctx context.Context, image *entities.ProductImage) error {
	return r.db.WithContext(ctx).Save(image).Error
}

// Delete deletes an image by ID
func (r *imageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Use Unscoped() to force hard delete
	return r.db.WithContext(ctx).Unscoped().Delete(&entities.ProductImage{}, id).Error
}

// DeleteByProductID marks all images for a product as inactive (position = -1)
func (r *imageRepository) DeleteByProductID(ctx context.Context, productID uuid.UUID) error {
	// Instead of deleting, mark all images as inactive by setting position = -1
	result := r.db.WithContext(ctx).Model(&entities.ProductImage{}).
		Where("product_id = ?", productID).
		Update("position", -1)
	
	if result.Error != nil {
		return fmt.Errorf("failed to mark images as inactive: %w", result.Error)
	}

	return nil
}

// MarkAsInactive marks all images for a product as inactive (position = -1)
func (r *imageRepository) MarkAsInactive(ctx context.Context, productID uuid.UUID) error {
	return r.DeleteByProductID(ctx, productID)
}

// CreateBatch creates multiple images for a product
func (r *imageRepository) CreateBatch(ctx context.Context, images []*entities.ProductImage) error {
	if len(images) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(images, 100).Error
}
