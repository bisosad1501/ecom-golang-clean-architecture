package usecases

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

// ProductUseCase defines product use cases
type ProductUseCase interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*ProductResponse, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*ProductResponse, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	GetProducts(ctx context.Context, req GetProductsRequest) ([]*ProductResponse, error)
	SearchProducts(ctx context.Context, req SearchProductsRequest) ([]*ProductResponse, error)
	GetProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*ProductResponse, error)
	UpdateStock(ctx context.Context, productID uuid.UUID, stock int) error
}

type productUseCase struct {
	productRepo  repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(
	productRepo repositories.ProductRepository,
	categoryRepo repositories.CategoryRepository,
) ProductUseCase {
	return &productUseCase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateProductRequest represents create product request
type CreateProductRequest struct {
	Name         string                   `json:"name" validate:"required"`
	Description  string                   `json:"description"`
	SKU          string                   `json:"sku" validate:"required"`
	Price        float64                  `json:"price" validate:"required,gt=0"`
	ComparePrice *float64                 `json:"compare_price" validate:"omitempty,gt=0"`
	CostPrice    *float64                 `json:"cost_price" validate:"omitempty,gt=0"`
	Stock        int                      `json:"stock" validate:"min=0"`
	Weight       *float64                 `json:"weight" validate:"omitempty,gt=0"`
	Dimensions   *DimensionsRequest       `json:"dimensions"`
	CategoryID   uuid.UUID                `json:"category_id" validate:"required"`
	Images       []ProductImageRequest    `json:"images"`
	Tags         []string                 `json:"tags"`
	Status       entities.ProductStatus   `json:"status" validate:"required"`
	IsDigital    bool                     `json:"is_digital"`
}

// UpdateProductRequest represents update product request
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
	Images       []ProductImageRequest    `json:"images"`
	Tags         []string                 `json:"tags"`
	Status       *entities.ProductStatus  `json:"status"`
	IsDigital    *bool                    `json:"is_digital"`
}

// GetProductsRequest represents get products request
type GetProductsRequest struct {
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}

// SearchProductsRequest represents search products request
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

// DimensionsRequest represents dimensions request
type DimensionsRequest struct {
	Length float64 `json:"length" validate:"gt=0"`
	Width  float64 `json:"width" validate:"gt=0"`
	Height float64 `json:"height" validate:"gt=0"`
}

// ProductImageRequest represents product image request
type ProductImageRequest struct {
	URL      string `json:"url" validate:"required,url"`
	AltText  string `json:"alt_text"`
	Position int    `json:"position"`
}

// ProductResponse represents product response
type ProductResponse struct {
	ID           uuid.UUID                `json:"id"`
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	SKU          string                   `json:"sku"`
	Price        float64                  `json:"price"`
	ComparePrice *float64                 `json:"compare_price"`
	CostPrice    *float64                 `json:"cost_price"`
	Stock        int                      `json:"stock"`
	Weight       *float64                 `json:"weight"`
	Dimensions   *DimensionsResponse      `json:"dimensions"`
	Category     *ProductCategoryResponse `json:"category"`
	Images       []ProductImageResponse   `json:"images"`
	Tags         []ProductTagResponse     `json:"tags"`
	Status       entities.ProductStatus   `json:"status"`
	IsDigital    bool                     `json:"is_digital"`
	IsAvailable  bool                     `json:"is_available"`
	HasDiscount  bool                     `json:"has_discount"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
}

// DimensionsResponse represents dimensions response
type DimensionsResponse struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// ProductImageResponse represents product image response
type ProductImageResponse struct {
	ID       uuid.UUID `json:"id"`
	URL      string    `json:"url"`
	AltText  string    `json:"alt_text"`
	Position int       `json:"position"`
}

// ProductTagResponse represents product tag response
type ProductTagResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

// ProductCategoryResponse represents category response in product context
type ProductCategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Slug        string    `json:"slug"`
	Image       string    `json:"image"`
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

	if err := uc.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return uc.toProductResponse(product), nil
}

// GetProduct gets a product by ID
func (uc *productUseCase) GetProduct(ctx context.Context, id uuid.UUID) (*ProductResponse, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrProductNotFound
	}

	return uc.toProductResponse(product), nil
}

// UpdateProduct updates a product
func (uc *productUseCase) UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*ProductResponse, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, entities.ErrProductNotFound
	}

	// Update fields
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.ComparePrice != nil {
		product.ComparePrice = req.ComparePrice
	}
	if req.CostPrice != nil {
		product.CostPrice = req.CostPrice
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Weight != nil {
		product.Weight = req.Weight
	}
	if req.CategoryID != nil {
		// Verify category exists
		_, err = uc.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, entities.ErrCategoryNotFound
		}
		product.CategoryID = *req.CategoryID
	}
	if req.Status != nil {
		product.Status = *req.Status
	}
	if req.IsDigital != nil {
		product.IsDigital = *req.IsDigital
	}
	if req.Dimensions != nil {
		product.Dimensions = &entities.Dimensions{
			Length: req.Dimensions.Length,
			Width:  req.Dimensions.Width,
			Height: req.Dimensions.Height,
		}
	}

	product.UpdatedAt = time.Now()

	if err := uc.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return uc.toProductResponse(product), nil
}

// DeleteProduct deletes a product
func (uc *productUseCase) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return entities.ErrProductNotFound
	}

	return uc.productRepo.Delete(ctx, id)
}

// GetProducts gets list of products
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

// SearchProducts searches products
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

// GetProductsByCategory gets products by category
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

// UpdateStock updates product stock
func (uc *productUseCase) UpdateStock(ctx context.Context, productID uuid.UUID, stock int) error {
	_, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return entities.ErrProductNotFound
	}

	return uc.productRepo.UpdateStock(ctx, productID, stock)
}

// toProductResponse converts product entity to response
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

	// Convert images
	response.Images = make([]ProductImageResponse, len(product.Images))
	for i, img := range product.Images {
		response.Images[i] = ProductImageResponse{
			ID:       img.ID,
			URL:      img.URL,
			AltText:  img.AltText,
			Position: img.Position,
		}
	}

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
