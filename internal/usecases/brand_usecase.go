package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/pkg/utils"
	"github.com/google/uuid"
)

// BrandUseCase defines brand use cases
type BrandUseCase interface {
	CreateBrand(ctx context.Context, req CreateBrandRequest) (*BrandResponse, error)
	GetBrand(ctx context.Context, id uuid.UUID) (*BrandResponse, error)
	GetBrandBySlug(ctx context.Context, slug string) (*BrandResponse, error)
	UpdateBrand(ctx context.Context, id uuid.UUID, req UpdateBrandRequest) (*BrandResponse, error)
	DeleteBrand(ctx context.Context, id uuid.UUID) error
	GetBrands(ctx context.Context, req GetBrandsRequest) (*BrandsListResponse, error)
	SearchBrands(ctx context.Context, req SearchBrandsRequest) (*BrandsListResponse, error)
	GetActiveBrands(ctx context.Context, limit, offset int) (*BrandsListResponse, error)
	GetPopularBrands(ctx context.Context, limit int) ([]*BrandResponse, error)
	GetBrandsForFiltering(ctx context.Context, categoryID *uuid.UUID) ([]BrandFilterOption, error)
}

type brandUseCase struct {
	brandRepo repositories.BrandRepository
}

// NewBrandUseCase creates a new brand use case
func NewBrandUseCase(brandRepo repositories.BrandRepository) BrandUseCase {
	return &brandUseCase{
		brandRepo: brandRepo,
	}
}

// CreateBrandRequest represents create brand request
type CreateBrandRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Slug        string `json:"slug" validate:"omitempty,min=2,max=100"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	Logo        string `json:"logo" validate:"omitempty,url"`
	Website     string `json:"website" validate:"omitempty,url"`
	IsActive    bool   `json:"is_active"`
}

// UpdateBrandRequest represents update brand request
type UpdateBrandRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Slug        string `json:"slug" validate:"omitempty,min=2,max=100"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	Logo        string `json:"logo" validate:"omitempty,url"`
	Website     string `json:"website" validate:"omitempty,url"`
	IsActive    bool   `json:"is_active"`
}

// GetBrandsRequest represents get brands request
type GetBrandsRequest struct {
	Limit    int  `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset   int  `json:"offset" validate:"omitempty,min=0"`
	IsActive *bool `json:"is_active"`
}

// SearchBrandsRequest represents search brands request
type SearchBrandsRequest struct {
	Query  string `json:"query" validate:"required,min=1"`
	Limit  int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset int    `json:"offset" validate:"omitempty,min=0"`
}

// BrandResponse represents brand response
type BrandResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	Logo         string    `json:"logo"`
	Website      string    `json:"website"`
	IsActive     bool      `json:"is_active"`
	ProductCount int       `json:"product_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BrandsListResponse represents brands list response
type BrandsListResponse struct {
	Brands []BrandResponse `json:"brands"`
	Total  int64           `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

// BrandFilterOption represents brand filter option
type BrandFilterOption struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// CreateBrand creates a new brand
func (uc *brandUseCase) CreateBrand(ctx context.Context, req CreateBrandRequest) (*BrandResponse, error) {
	// Generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = utils.GenerateSlug(req.Name)
	}

	// Validate slug format
	if err := utils.ValidateSlug(slug); err != nil {
		return nil, fmt.Errorf("invalid slug: %w", err)
	}

	// Check if slug already exists
	exists, err := uc.brandRepo.ExistsBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entities.ErrConflict
	}

	// Create brand
	brand := &entities.Brand{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(req.Name),
		Slug:        slug,
		Description: strings.TrimSpace(req.Description),
		Logo:        req.Logo,
		Website:     req.Website,
		IsActive:    req.IsActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.brandRepo.Create(ctx, brand); err != nil {
		return nil, err
	}

	return uc.toBrandResponse(brand), nil
}

// GetBrand gets a brand by ID
func (uc *brandUseCase) GetBrand(ctx context.Context, id uuid.UUID) (*BrandResponse, error) {
	brand, err := uc.brandRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrBrandNotFound
	}

	return uc.toBrandResponse(brand), nil
}

// GetBrandBySlug gets a brand by slug
func (uc *brandUseCase) GetBrandBySlug(ctx context.Context, slug string) (*BrandResponse, error) {
	brand, err := uc.brandRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, entities.ErrBrandNotFound
	}

	return uc.toBrandResponse(brand), nil
}

// UpdateBrand updates a brand
func (uc *brandUseCase) UpdateBrand(ctx context.Context, id uuid.UUID, req UpdateBrandRequest) (*BrandResponse, error) {
	// Get existing brand
	brand, err := uc.brandRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrBrandNotFound
	}

	// Generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = utils.GenerateSlug(req.Name)
	}

	// Validate slug format
	if err := utils.ValidateSlug(slug); err != nil {
		return nil, fmt.Errorf("invalid slug: %w", err)
	}

	// Check if slug already exists (excluding current brand)
	if slug != brand.Slug {
		exists, err := uc.brandRepo.ExistsBySlug(ctx, slug)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, entities.ErrConflict
		}
	}

	// Update brand
	brand.Name = strings.TrimSpace(req.Name)
	brand.Slug = slug
	brand.Description = strings.TrimSpace(req.Description)
	brand.Logo = req.Logo
	brand.Website = req.Website
	brand.IsActive = req.IsActive
	brand.UpdatedAt = time.Now()

	if err := uc.brandRepo.Update(ctx, brand); err != nil {
		return nil, err
	}

	return uc.toBrandResponse(brand), nil
}

// DeleteBrand deletes a brand
func (uc *brandUseCase) DeleteBrand(ctx context.Context, id uuid.UUID) error {
	// Check if brand exists
	_, err := uc.brandRepo.GetByID(ctx, id)
	if err != nil {
		return entities.ErrBrandNotFound
	}

	// TODO: Check if brand has products and handle accordingly
	// For now, we'll allow deletion (products will have null brand_id)

	return uc.brandRepo.Delete(ctx, id)
}

// GetBrands gets brands with pagination
func (uc *brandUseCase) GetBrands(ctx context.Context, req GetBrandsRequest) (*BrandsListResponse, error) {
	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	var brands []*entities.Brand
	var err error

	if req.IsActive != nil && *req.IsActive {
		brands, err = uc.brandRepo.GetActive(ctx, req.Limit, req.Offset)
	} else {
		brands, err = uc.brandRepo.List(ctx, req.Limit, req.Offset)
	}

	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := uc.brandRepo.GetTotal(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to response
	brandResponses := make([]BrandResponse, len(brands))
	for i, brand := range brands {
		brandResponses[i] = *uc.toBrandResponse(brand)
	}

	return &BrandsListResponse{
		Brands: brandResponses,
		Total:  total,
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}

// SearchBrands searches brands
func (uc *brandUseCase) SearchBrands(ctx context.Context, req SearchBrandsRequest) (*BrandsListResponse, error) {
	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	brands, err := uc.brandRepo.Search(ctx, req.Query, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	// Convert to response
	brandResponses := make([]BrandResponse, len(brands))
	for i, brand := range brands {
		brandResponses[i] = *uc.toBrandResponse(brand)
	}

	return &BrandsListResponse{
		Brands: brandResponses,
		Total:  int64(len(brands)), // For search, we return actual count
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}

// GetActiveBrands gets active brands
func (uc *brandUseCase) GetActiveBrands(ctx context.Context, limit, offset int) (*BrandsListResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	brands, err := uc.brandRepo.GetActive(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Get total active count
	total, err := uc.brandRepo.CountByStatus(ctx, true)
	if err != nil {
		return nil, err
	}

	// Convert to response
	brandResponses := make([]BrandResponse, len(brands))
	for i, brand := range brands {
		brandResponses[i] = *uc.toBrandResponse(brand)
	}

	return &BrandsListResponse{
		Brands: brandResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// GetPopularBrands gets popular brands by product count
func (uc *brandUseCase) GetPopularBrands(ctx context.Context, limit int) ([]*BrandResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	brands, err := uc.brandRepo.GetPopularBrands(ctx, limit)
	if err != nil {
		return nil, err
	}

	// Convert to response
	brandResponses := make([]*BrandResponse, len(brands))
	for i, brand := range brands {
		brandResponses[i] = uc.toBrandResponse(brand)
	}

	return brandResponses, nil
}

// GetBrandsForFiltering gets brands for product filtering
func (uc *brandUseCase) GetBrandsForFiltering(ctx context.Context, categoryID *uuid.UUID) ([]BrandFilterOption, error) {
	results, err := uc.brandRepo.GetBrandsForFiltering(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// Convert to filter options
	options := make([]BrandFilterOption, len(results))
	for i, result := range results {
		options[i] = BrandFilterOption{
			ID:    result["id"].(string),
			Name:  result["name"].(string),
			Count: int(result["count"].(int64)),
		}
	}

	return options, nil
}

// toBrandResponse converts brand entity to response
func (uc *brandUseCase) toBrandResponse(brand *entities.Brand) *BrandResponse {
	return &BrandResponse{
		ID:           brand.ID,
		Name:         brand.Name,
		Slug:         brand.Slug,
		Description:  brand.Description,
		Logo:         brand.Logo,
		Website:      brand.Website,
		IsActive:     brand.IsActive,
		ProductCount: len(brand.Products), // This will be populated by repository if needed
		CreatedAt:    brand.CreatedAt,
		UpdatedAt:    brand.UpdatedAt,
	}
}
