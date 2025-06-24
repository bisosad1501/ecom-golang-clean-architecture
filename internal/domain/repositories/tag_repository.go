package repositories

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// TagRepository defines the interface for tag data access
type TagRepository interface {
	// Create creates a new tag
	Create(ctx context.Context, tag *entities.ProductTag) error
	
	// GetByID retrieves a tag by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ProductTag, error)
	
	// GetByName retrieves a tag by name
	GetByName(ctx context.Context, name string) (*entities.ProductTag, error)
	
	// GetBySlug retrieves a tag by slug
	GetBySlug(ctx context.Context, slug string) (*entities.ProductTag, error)
	
	// Update updates an existing tag
	Update(ctx context.Context, tag *entities.ProductTag) error
	
	// Delete deletes a tag by ID
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List retrieves tags with pagination
	List(ctx context.Context, limit, offset int) ([]*entities.ProductTag, error)
	
	// ExistsByName checks if a tag exists with the given name
	ExistsByName(ctx context.Context, name string) (bool, error)
	
	// ExistsBySlug checks if a tag exists with the given slug
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
	
	// FindOrCreate finds existing tag by name or creates new one
	FindOrCreate(ctx context.Context, name string) (*entities.ProductTag, error)
}
