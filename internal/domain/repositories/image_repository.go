package repositories

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// ImageRepository defines the interface for product image data access
type ImageRepository interface {
	// Create creates a new product image
	Create(ctx context.Context, image *entities.ProductImage) error
	
	// GetByID retrieves an image by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ProductImage, error)
	
	// GetByProductID retrieves images by product ID
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]*entities.ProductImage, error)
	
	// Update updates an existing image
	Update(ctx context.Context, image *entities.ProductImage) error
	
	// Delete deletes an image by ID
	Delete(ctx context.Context, id uuid.UUID) error
	
	// DeleteByProductID deletes all images for a product
	DeleteByProductID(ctx context.Context, productID uuid.UUID) error
	
	// CreateBatch creates multiple images for a product
	CreateBatch(ctx context.Context, images []*entities.ProductImage) error
	
	// MarkAsInactive marks all images for a product as inactive (position = -1)
	MarkAsInactive(ctx context.Context, productID uuid.UUID) error
}
