package usecases

import (
	"context"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

// CategoryUseCase defines category use cases
type CategoryUseCase interface {
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (*CategoryResponse, error)
	GetCategory(ctx context.Context, id uuid.UUID) (*CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, req UpdateCategoryRequest) (*CategoryResponse, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	GetCategories(ctx context.Context, req GetCategoriesRequest) ([]*CategoryResponse, error)
	GetCategoryTree(ctx context.Context) ([]*CategoryResponse, error)
	GetRootCategories(ctx context.Context) ([]*CategoryResponse, error)
	GetCategoryChildren(ctx context.Context, parentID uuid.UUID) ([]*CategoryResponse, error)
}

type categoryUseCase struct {
	categoryRepo repositories.CategoryRepository
}

// NewCategoryUseCase creates a new category use case
func NewCategoryUseCase(categoryRepo repositories.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{
		categoryRepo: categoryRepo,
	}
}

// CreateCategoryRequest represents create category request
type CreateCategoryRequest struct {
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description"`
	Slug        string     `json:"slug" validate:"required"`
	Image       string     `json:"image"`
	ParentID    *uuid.UUID `json:"parent_id"`
	IsActive    bool       `json:"is_active"`
	SortOrder   int        `json:"sort_order"`
}

// UpdateCategoryRequest represents update category request
type UpdateCategoryRequest struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Slug        *string    `json:"slug"`
	Image       *string    `json:"image"`
	ParentID    *uuid.UUID `json:"parent_id"`
	IsActive    *bool      `json:"is_active"`
	SortOrder   *int       `json:"sort_order"`
}

// GetCategoriesRequest represents get categories request
type GetCategoriesRequest struct {
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}

// CategoryResponse represents category response
type CategoryResponse struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Slug        string             `json:"slug"`
	Image       string             `json:"image"`
	ParentID    *uuid.UUID         `json:"parent_id"`
	Parent      *CategoryResponse  `json:"parent,omitempty"`
	Children    []CategoryResponse `json:"children,omitempty"`
	IsActive    bool               `json:"is_active"`
	SortOrder   int                `json:"sort_order"`
	Level       int                `json:"level"`
	Path        string             `json:"path"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// CreateCategory creates a new category
func (uc *categoryUseCase) CreateCategory(ctx context.Context, req CreateCategoryRequest) (*CategoryResponse, error) {
	// Check if slug already exists
	exists, err := uc.categoryRepo.ExistsBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entities.ErrConflict
	}

	// Verify parent category exists if provided
	if req.ParentID != nil {
		_, err = uc.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, entities.ErrCategoryNotFound
		}
	}

	// Generate slug if not provided
	if req.Slug == "" {
		req.Slug = generateSlug(req.Name)
	}

	// Create category
	category := &entities.Category{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Slug:        req.Slug,
		Image:       req.Image,
		ParentID:    req.ParentID,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return uc.toCategoryResponse(category), nil
}

// GetCategory gets a category by ID
func (uc *categoryUseCase) GetCategory(ctx context.Context, id uuid.UUID) (*CategoryResponse, error) {
	category, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrCategoryNotFound
	}

	return uc.toCategoryResponse(category), nil
}

// UpdateCategory updates a category
func (uc *categoryUseCase) UpdateCategory(ctx context.Context, id uuid.UUID, req UpdateCategoryRequest) (*CategoryResponse, error) {
	category, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrCategoryNotFound
	}

	// Update fields
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = *req.Description
	}
	if req.Slug != nil {
		// Check if new slug already exists
		if *req.Slug != category.Slug {
			exists, err := uc.categoryRepo.ExistsBySlug(ctx, *req.Slug)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, entities.ErrConflict
			}
		}
		category.Slug = *req.Slug
	}
	if req.Image != nil {
		category.Image = *req.Image
	}
	if req.ParentID != nil {
		// Verify parent category exists
		_, err = uc.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, entities.ErrCategoryNotFound
		}
		category.ParentID = req.ParentID
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}

	category.UpdatedAt = time.Now()

	if err := uc.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return uc.toCategoryResponse(category), nil
}

// DeleteCategory deletes a category
func (uc *categoryUseCase) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	_, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return entities.ErrCategoryNotFound
	}

	// Check if category has children
	children, err := uc.categoryRepo.GetChildren(ctx, id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return entities.ErrConflict // Cannot delete category with children
	}

	return uc.categoryRepo.Delete(ctx, id)
}

// GetCategories gets list of categories
func (uc *categoryUseCase) GetCategories(ctx context.Context, req GetCategoriesRequest) ([]*CategoryResponse, error) {
	categories, err := uc.categoryRepo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = uc.toCategoryResponse(category)
	}

	return responses, nil
}

// GetCategoryTree gets the category tree
func (uc *categoryUseCase) GetCategoryTree(ctx context.Context) ([]*CategoryResponse, error) {
	categories, err := uc.categoryRepo.GetTree(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = uc.toCategoryResponseWithChildren(category)
	}

	return responses, nil
}

// GetRootCategories gets root categories
func (uc *categoryUseCase) GetRootCategories(ctx context.Context) ([]*CategoryResponse, error) {
	categories, err := uc.categoryRepo.GetRootCategories(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = uc.toCategoryResponse(category)
	}

	return responses, nil
}

// GetCategoryChildren gets category children
func (uc *categoryUseCase) GetCategoryChildren(ctx context.Context, parentID uuid.UUID) ([]*CategoryResponse, error) {
	categories, err := uc.categoryRepo.GetChildren(ctx, parentID)
	if err != nil {
		return nil, err
	}

	responses := make([]*CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = uc.toCategoryResponse(category)
	}

	return responses, nil
}

// toCategoryResponse converts category entity to response
func (uc *categoryUseCase) toCategoryResponse(category *entities.Category) *CategoryResponse {
	response := &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Slug:        category.Slug,
		Image:       category.Image,
		ParentID:    category.ParentID,
		IsActive:    category.IsActive,
		SortOrder:   category.SortOrder,
		Level:       category.GetLevel(),
		Path:        category.GetPath(),
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	// Convert parent if available
	if category.Parent != nil {
		response.Parent = &CategoryResponse{
			ID:          category.Parent.ID,
			Name:        category.Parent.Name,
			Description: category.Parent.Description,
			Slug:        category.Parent.Slug,
			Image:       category.Parent.Image,
			ParentID:    category.Parent.ParentID,
			IsActive:    category.Parent.IsActive,
			SortOrder:   category.Parent.SortOrder,
			Level:       category.Parent.GetLevel(),
			Path:        category.Parent.GetPath(),
			CreatedAt:   category.Parent.CreatedAt,
			UpdatedAt:   category.Parent.UpdatedAt,
		}
	}

	return response
}

// toCategoryResponseWithChildren converts category entity to response with children
func (uc *categoryUseCase) toCategoryResponseWithChildren(category *entities.Category) *CategoryResponse {
	response := uc.toCategoryResponse(category)

	// Convert children
	if len(category.Children) > 0 {
		response.Children = make([]CategoryResponse, len(category.Children))
		for i, child := range category.Children {
			childResponse := uc.toCategoryResponseWithChildren(&child)
			response.Children[i] = *childResponse
		}
	}

	return response
}

// generateSlug generates a URL-friendly slug from a name
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "&", "and")
	// Remove special characters (basic implementation)
	return slug
}
