package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

// Request structs
type CreateProductRequest struct {
	Name         string                   `json:"name" validate:"required"`
	Description  string                   `json:"description" validate:"required"`
	SKU          string                   `json:"sku" validate:"required"`
	Price        float64                  `json:"price" validate:"required,gt=0"`
	ComparePrice *float64                 `json:"compare_price" validate:"omitempty,gt=0"`
	CostPrice    *float64                 `json:"cost_price" validate:"omitempty,gt=0"`
	Stock        int                      `json:"stock" validate:"required,min=0"`
	Weight       *float64                 `json:"weight" validate:"omitempty,gt=0"`
	Dimensions   *DimensionsRequest       `json:"dimensions"`
	CategoryID   uuid.UUID                `json:"category_id" validate:"required"`
	Images       []ProductImageRequest    `json:"images"`
	Tags         []string                 `json:"tags"`
	Status       entities.ProductStatus   `json:"status"`
	IsDigital    bool                     `json:"is_digital"`
}

type GetProductsRequest struct {
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}

type SearchProductsRequest struct {
	Query      string                  `json:"query"`
	CategoryID *uuid.UUID              `json:"category_id"`
	MinPrice   *float64                `json:"min_price" validate:"omitempty,gt=0"`
	MaxPrice   *float64                `json:"max_price" validate:"omitempty,gt=0"`
	Status     *entities.ProductStatus `json:"status"`
	Tags       []string                `json:"tags"`
	SortBy     string                  `json:"sort_by"`
	SortOrder  string                  `json:"sort_order"`
	Limit      int                     `json:"limit" validate:"min=1,max=100"`
	Offset     int                     `json:"offset" validate:"min=0"`
}

type DimensionsRequest struct {
	Length float64 `json:"length" validate:"required,gt=0"`
	Width  float64 `json:"width" validate:"required,gt=0"`
	Height float64 `json:"height" validate:"required,gt=0"`
}

type ProductImageRequest struct {
	URL     string `json:"url" validate:"required,url"`
	AltText string `json:"alt_text"`
}

// Response structs are defined in types.go

// ProductUseCase defines product use cases
type ProductUseCase interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*ProductResponse, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*ProductResponse, error)
	PatchProduct(ctx context.Context, id uuid.UUID, req PatchProductRequest) (*ProductResponse, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	GetProducts(ctx context.Context, req GetProductsRequest) ([]*ProductResponse, error)
	SearchProducts(ctx context.Context, req SearchProductsRequest) ([]*ProductResponse, error)
	GetProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*ProductResponse, error)
	UpdateStock(ctx context.Context, productID uuid.UUID, stock int) error
}

type productUseCase struct {
	productRepo  repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
	tagRepo      repositories.TagRepository
	imageRepo    repositories.ImageRepository
	cartRepo     repositories.CartRepository
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(
	productRepo repositories.ProductRepository,
	categoryRepo repositories.CategoryRepository,
	tagRepo repositories.TagRepository,
	imageRepo repositories.ImageRepository,
	cartRepo repositories.CartRepository,
) ProductUseCase {
	return &productUseCase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
		imageRepo:    imageRepo,
		cartRepo:     cartRepo,
	}
}

type UpdateProductRequest struct {
	Name         *string                  `json:"name"`
	Description  *string                  `json:"description"`
	Price        *float64                 `json:"price" validate:"omitempty,gt=0"`
	ComparePrice *float64                 `json:"compare_price" validate:"omitempty,gt=0"`
	CostPrice    *float64                 `json:"cost_price" validate:"omitempty,gt=0"`
	Stock        *int                     `json:"stock" validate:"omitempty,min=0"`
	Weight       *float64                 `json:"weight" validate:"omitempty,gt=0"`
	Dimensions   *DimensionsRequest       `json:"dimensions"`
	CategoryID   *uuid.UUID               `json:"category_id"`
	Images       []ProductImageRequest    `json:"images"`        // For PUT: replace all images
	Tags         []string                 `json:"tags"`          // For PUT: replace all tags
	Status       *entities.ProductStatus  `json:"status"`
	IsDigital    *bool                    `json:"is_digital"`
}

// PatchProductRequest for PATCH operations - only updates provided fields
type PatchProductRequest struct {
	Name         *string                  `json:"name"`
	Description  *string                  `json:"description"`
	Price        *float64                 `json:"price" validate:"omitempty,gt=0"`
	ComparePrice *float64                 `json:"compare_price" validate:"omitempty,gt=0"`
	CostPrice    *float64                 `json:"cost_price" validate:"omitempty,gt=0"`
	Stock        *int                     `json:"stock" validate:"omitempty,min=0"`
	Weight       *float64                 `json:"weight" validate:"omitempty,gt=0"`
	Dimensions   *DimensionsRequest       `json:"dimensions"`
	CategoryID   *uuid.UUID               `json:"category_id"`
	Images       *[]ProductImageRequest   `json:"images"`        // For PATCH: nil = no change, empty = clear all, values = replace
	Tags         *[]string                `json:"tags"`          // For PATCH: nil = no change, empty = clear all, values = replace  
	Status       *entities.ProductStatus  `json:"status"`
	IsDigital    *bool                    `json:"is_digital"`
}

// CreateProduct creates a new product
func (uc *productUseCase) CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error) {
	// Check if SKU already exists
	exists, err := uc.productRepo.ExistsBySKU(ctx, req.SKU)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entities.ErrConflict
	}

	// Verify category exists
	_, err = uc.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, entities.ErrCategoryNotFound
	}

	// Create product
	product := &entities.Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Price:       req.Price,
		ComparePrice: req.ComparePrice,
		CostPrice:   req.CostPrice,
		Stock:       req.Stock,
		Weight:      req.Weight,
		CategoryID:  req.CategoryID,
		Status:      req.Status,
		IsDigital:   req.IsDigital,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.Dimensions != nil {
		product.Dimensions = &entities.Dimensions{
			Length: req.Dimensions.Length,
			Width:  req.Dimensions.Width,
			Height: req.Dimensions.Height,
		}
	}

	// Create product first
	if err := uc.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	// Handle tags if provided
	if len(req.Tags) > 0 {
		if err := uc.replaceProductTags(ctx, product.ID, req.Tags); err != nil {
			return nil, err
		}
	}

	// Handle images if provided
	if len(req.Images) > 0 {
		if err := uc.replaceProductImages(ctx, product.ID, req.Images); err != nil {
			return nil, err
		}
	}

	// Reload and return
	updatedProduct, err := uc.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	return uc.toProductResponse(updatedProduct), nil
}

// GetProduct gets a product by ID
func (uc *productUseCase) GetProduct(ctx context.Context, id uuid.UUID) (*ProductResponse, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrProductNotFound
	}

	return uc.toProductResponse(product), nil
}

// UpdateProduct updates a product with improved business logic
func (uc *productUseCase) UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*ProductResponse, error) {
	// Get existing product
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrProductNotFound
	}

	// Track what needs to be updated
	hasChanges := false
	
	// Update basic fields only if they are provided
	if req.Name != nil {
		if *req.Name == "" {
			return nil, fmt.Errorf("name cannot be empty")
		}
		product.Name = *req.Name
		hasChanges = true
	}
	
	if req.Description != nil {
		product.Description = *req.Description
		hasChanges = true
	}
	
	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, fmt.Errorf("price must be greater than 0")
		}
		product.Price = *req.Price
		hasChanges = true
	}
	
	if req.ComparePrice != nil {
		if *req.ComparePrice <= 0 {
			return nil, fmt.Errorf("compare price must be greater than 0")
		}
		product.ComparePrice = req.ComparePrice
		hasChanges = true
	}
	
	if req.CostPrice != nil {
		if *req.CostPrice < 0 {
			return nil, fmt.Errorf("cost price cannot be negative")
		}
		product.CostPrice = req.CostPrice
		hasChanges = true
	}
	
	if req.Stock != nil {
		if *req.Stock < 0 {
			return nil, fmt.Errorf("stock cannot be negative")
		}
		product.Stock = *req.Stock
		hasChanges = true
	}
	
	if req.Weight != nil {
		if *req.Weight <= 0 {
			return nil, fmt.Errorf("weight must be greater than 0")
		}
		product.Weight = req.Weight
		hasChanges = true
	}
	
	if req.CategoryID != nil {
		// Verify category exists
		_, err := uc.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, entities.ErrCategoryNotFound
		}
		product.CategoryID = *req.CategoryID
		hasChanges = true
	}
	
	if req.Status != nil {
		product.Status = *req.Status
		hasChanges = true
	}
	
	if req.IsDigital != nil {
		product.IsDigital = *req.IsDigital
		hasChanges = true
	}
	
	if req.Dimensions != nil {
		if req.Dimensions.Length <= 0 || req.Dimensions.Width <= 0 || req.Dimensions.Height <= 0 {
			return nil, fmt.Errorf("all dimensions must be greater than 0")
		}
		product.Dimensions = &entities.Dimensions{
			Length: req.Dimensions.Length,
			Width:  req.Dimensions.Width,
			Height: req.Dimensions.Height,
		}
		hasChanges = true
	}

	// Handle Images - Complete replacement if provided
	if req.Images != nil {
		if err := uc.replaceProductImages(ctx, product.ID, req.Images); err != nil {
			return nil, fmt.Errorf("failed to update images: %w", err)
		}
		hasChanges = true
	}

	// Handle Tags - Complete replacement if provided
	if req.Tags != nil {
		if err := uc.replaceProductTags(ctx, product.ID, req.Tags); err != nil {
			return nil, fmt.Errorf("failed to update tags: %w", err)
		}
		hasChanges = true
	}

	// Only update product if there were actual changes to basic fields
	if hasChanges {
		product.UpdatedAt = time.Now()
		if err := uc.productRepo.Update(ctx, product); err != nil {
			return nil, fmt.Errorf("failed to update product: %w", err)
		}
	}

	// Return updated product with fresh data - force fresh reload from database
	// Clear any potential cache by creating a fresh query
	updatedProduct, err := uc.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated product: %w", err)
	}

	return uc.toProductResponse(updatedProduct), nil
}

// PatchProduct partially updates a product - only updates provided fields
func (uc *productUseCase) PatchProduct(ctx context.Context, id uuid.UUID, req PatchProductRequest) (*ProductResponse, error) {
	// Get existing product
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrProductNotFound
	}

	var hasChanges bool

	// Basic field updates - only if provided
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return nil, fmt.Errorf("name cannot be empty")
		}
		product.Name = *req.Name
		hasChanges = true
	}

	if req.Description != nil {
		product.Description = *req.Description
		hasChanges = true
	}

	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, fmt.Errorf("price must be greater than 0")
		}
		product.Price = *req.Price
		hasChanges = true
	}

	if req.ComparePrice != nil {
		if *req.ComparePrice <= 0 {
			return nil, fmt.Errorf("compare price must be greater than 0")
		}
		product.ComparePrice = req.ComparePrice
		hasChanges = true
	}

	if req.CostPrice != nil {
		if *req.CostPrice < 0 {
			return nil, fmt.Errorf("cost price cannot be negative")
		}
		product.CostPrice = req.CostPrice
		hasChanges = true
	}

	if req.Stock != nil {
		if *req.Stock < 0 {
			return nil, fmt.Errorf("stock cannot be negative")
		}
		product.Stock = *req.Stock
		hasChanges = true
	}

	if req.Weight != nil {
		if *req.Weight <= 0 {
			return nil, fmt.Errorf("weight must be greater than 0")
		}
		product.Weight = req.Weight
		hasChanges = true
	}

	if req.CategoryID != nil {
		// Verify category exists
		_, err := uc.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, entities.ErrCategoryNotFound
		}
		product.CategoryID = *req.CategoryID
		hasChanges = true
	}

	if req.Status != nil {
		product.Status = *req.Status
		hasChanges = true
	}

	if req.IsDigital != nil {
		product.IsDigital = *req.IsDigital
		hasChanges = true
	}

	if req.Dimensions != nil {
		if req.Dimensions.Length <= 0 || req.Dimensions.Width <= 0 || req.Dimensions.Height <= 0 {
			return nil, fmt.Errorf("dimensions must be positive values")
		}
		product.Dimensions = &entities.Dimensions{
			Length: req.Dimensions.Length,
			Width:  req.Dimensions.Width,
			Height: req.Dimensions.Height,
		}
		hasChanges = true
	}

	// Handle Images - check if field is provided
	if req.Images != nil {
		// If empty slice, clear all images
		if len(*req.Images) == 0 {
			if err := uc.imageRepo.MarkAsInactive(ctx, product.ID); err != nil {
				return nil, fmt.Errorf("failed to clear images: %w", err)
			}
		} else {
			// Replace with new images
			if err := uc.replaceProductImages(ctx, product.ID, *req.Images); err != nil {
				return nil, fmt.Errorf("failed to update images: %w", err)
			}
		}
		hasChanges = true
	}

	// Handle Tags - check if field is provided
	if req.Tags != nil {
		// Convert to slice and process
		tagSlice := *req.Tags
		if err := uc.replaceProductTags(ctx, product.ID, tagSlice); err != nil {
			return nil, fmt.Errorf("failed to update tags: %w", err)
		}
		hasChanges = true
	}

	// Only update product if there were actual changes
	if hasChanges {
		product.UpdatedAt = time.Now()
		if err := uc.productRepo.Update(ctx, product); err != nil {
			return nil, fmt.Errorf("failed to update product: %w", err)
		}
	}

	// Return updated product with fresh data
	updatedProduct, err := uc.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated product: %w", err)
	}

	return uc.toProductResponse(updatedProduct), nil
}

// replaceProductImages completely replaces all product images with new ones
func (uc *productUseCase) replaceProductImages(ctx context.Context, productID uuid.UUID, images []ProductImageRequest) error {
	fmt.Printf("DEBUG: replaceProductImages called for productID: %s with %d new images\n", productID.String(), len(images))
	
	// Validate images
	for i, img := range images {
		if img.URL == "" {
			return fmt.Errorf("image URL cannot be empty at position %d", i)
		}
	}
	
	// Step 1: Get existing images
	fmt.Printf("DEBUG: Getting existing images for productID: %s\n", productID.String())
	existingImages, err := uc.imageRepo.GetByProductID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get existing images: %w", err)
	}
	fmt.Printf("DEBUG: Found %d existing images\n", len(existingImages))
	
	// Step 2: Update/Replace strategy instead of delete
	// Mark all existing images as "hidden" by setting position to -1
	if len(existingImages) > 0 {
		fmt.Printf("DEBUG: Marking %d existing images as hidden\n", len(existingImages))
		for _, img := range existingImages {
			img.Position = -1 // Mark as hidden
			if err := uc.imageRepo.Update(ctx, img); err != nil {
				fmt.Printf("DEBUG: Error hiding image %s: %v\n", img.ID.String(), err)
				return fmt.Errorf("failed to hide existing image: %w", err)
			}
		}
		fmt.Printf("DEBUG: Successfully marked existing images as hidden\n")
	}

	// Step 3: Create new images with positive positions
	if len(images) > 0 {
		fmt.Printf("DEBUG: Creating %d new images\n", len(images))
		var newImages []*entities.ProductImage
		for i, imgReq := range images {
			image := &entities.ProductImage{
				ID:        uuid.New(),
				ProductID: productID,
				URL:       imgReq.URL,
				AltText:   imgReq.AltText,
				Position:  i, // Positive position (0, 1, 2, ...)
				CreatedAt: time.Now(),
			}
			newImages = append(newImages, image)
			fmt.Printf("DEBUG: Prepared new image %d: %s at position %d\n", i, imgReq.URL, i)
		}

		if err := uc.imageRepo.CreateBatch(ctx, newImages); err != nil {
			fmt.Printf("DEBUG: Error creating new images: %v\n", err)
			return fmt.Errorf("failed to create new images: %w", err)
		}
		fmt.Printf("DEBUG: Successfully created %d new images\n", len(newImages))
	} else {
		fmt.Printf("DEBUG: No new images to create\n")
	}

	// Step 4: Verify by counting active images (position >= 0)
	activeImages, err := uc.getActiveImagesByProductID(ctx, productID)
	if err != nil {
		fmt.Printf("DEBUG: Error counting active images: %v\n", err)
	} else {
		fmt.Printf("DEBUG: After replacement, product has %d active images\n", len(activeImages))
	}

	fmt.Printf("DEBUG: replaceProductImages completed successfully\n")
	return nil
}

// Helper function to get active images (position >= 0)
func (uc *productUseCase) getActiveImagesByProductID(ctx context.Context, productID uuid.UUID) ([]*entities.ProductImage, error) {
	allImages, err := uc.imageRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}
	
	var activeImages []*entities.ProductImage
	for _, img := range allImages {
		if img.Position >= 0 {
			activeImages = append(activeImages, img)
		}
	}
	return activeImages, nil
}

// replaceProductTags completely replaces all product tags with new ones
func (uc *productUseCase) replaceProductTags(ctx context.Context, productID uuid.UUID, tagNames []string) error {
	// Validate and clean tag names
	var validTags []string
	for _, tagName := range tagNames {
		cleanTag := strings.TrimSpace(tagName)
		if cleanTag != "" && len(cleanTag) <= 50 { // Reasonable limit for tag length
			validTags = append(validTags, cleanTag)
		}
	}
	
	// If no valid tags, clear all tags
	if len(validTags) == 0 {
		return uc.productRepo.ClearTags(ctx, productID)
	}
	
	// Find or create all tags and collect their IDs
	var tagIDs []uuid.UUID
	for _, tagName := range validTags {
		tag, err := uc.tagRepo.FindOrCreate(ctx, tagName)
		if err != nil {
			return fmt.Errorf("failed to find/create tag '%s': %w", tagName, err)
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	// Replace all tags at once using the new ReplaceTags method
	if err := uc.productRepo.ReplaceTags(ctx, productID, tagIDs); err != nil {
		return fmt.Errorf("failed to replace product tags: %w", err)
	}

	return nil
}

// DeleteProduct deletes a product (same as original)
func (uc *productUseCase) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return entities.ErrProductNotFound
	}

	// First, remove all cart items that reference this product
	err = uc.cartRepo.RemoveItemsByProductID(ctx, id)
	if err != nil {
		return err
	}

	// Then delete the product
	return uc.productRepo.Delete(ctx, id)
}

// GetProducts gets list of products (same as original)
func (uc *productUseCase) GetProducts(ctx context.Context, req GetProductsRequest) ([]*ProductResponse, error) {
	products, err := uc.productRepo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*ProductResponse, len(products))
	for i, product := range products {
		responses[i] = uc.toProductResponse(product)
	}

	return responses, nil
}

// SearchProducts searches products (same as original)
func (uc *productUseCase) SearchProducts(ctx context.Context, req SearchProductsRequest) ([]*ProductResponse, error) {
	params := repositories.ProductSearchParams{
		Query:      req.Query,
		CategoryID: req.CategoryID,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		Status:     req.Status,
		Tags:       req.Tags,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}

	products, err := uc.productRepo.Search(ctx, params)
	if err != nil {
		return nil, err
	}

	responses := make([]*ProductResponse, len(products))
	for i, product := range products {
		responses[i] = uc.toProductResponse(product)
	}

	return responses, nil
}

// GetProductsByCategory gets products by category (same as original)
func (uc *productUseCase) GetProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*ProductResponse, error) {
	products, err := uc.productRepo.GetByCategory(ctx, categoryID, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*ProductResponse, len(products))
	for i, product := range products {
		responses[i] = uc.toProductResponse(product)
	}

	return responses, nil
}

// UpdateStock updates product stock (same as original)
func (uc *productUseCase) UpdateStock(ctx context.Context, productID uuid.UUID, stock int) error {
	_, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return entities.ErrProductNotFound
	}

	return uc.productRepo.UpdateStock(ctx, productID, stock)
}

// toProductResponse converts product entity to response (same as original)
func (uc *productUseCase) toProductResponse(product *entities.Product) *ProductResponse {
	response := &ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		SKU:         product.SKU,
		Price:       product.Price,
		ComparePrice: product.ComparePrice,
		CostPrice:   product.CostPrice,
		Stock:       product.Stock,
		Weight:      product.Weight,
		Status:      product.Status,
		IsDigital:   product.IsDigital,
		IsAvailable: product.IsAvailable(),
		HasDiscount: product.HasDiscount(),
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	if product.Dimensions != nil {
		response.Dimensions = &DimensionsResponse{
			Length: product.Dimensions.Length,
			Width:  product.Dimensions.Width,
			Height: product.Dimensions.Height,
		}
	}

	if product.Category.ID != uuid.Nil {
		response.Category = &ProductCategoryResponse{
			ID:          product.Category.ID,
			Name:        product.Category.Name,
			Description: product.Category.Description,
			Slug:        product.Category.Slug,
			Image:       product.Category.Image,
		}
	}

	// Convert images - only include active images (position >= 0)
	var activeImages []ProductImageResponse
	for _, img := range product.Images {
		if img.Position >= 0 { // Only include active images
			activeImages = append(activeImages, ProductImageResponse{
				ID:       img.ID,
				URL:      img.URL,
				AltText:  img.AltText,
				Position: img.Position,
			})
		}
	}
	response.Images = activeImages

	// Convert tags
	response.Tags = make([]ProductTagResponse, len(product.Tags))
	for i, tag := range product.Tags {
		response.Tags[i] = ProductTagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		}
	}

	return response
}
