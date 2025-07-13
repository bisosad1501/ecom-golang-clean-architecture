package usecases

import (
	"context"
	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// ProductComparisonRequest represents a request to create/update comparison
type ProductComparisonRequest struct {
	Name      string      `json:"name"`
	ProductIDs []uuid.UUID `json:"product_ids" validate:"required,min=2,max=5"`
}

// ProductComparisonResponse represents a comparison response
type ProductComparisonResponse struct {
	ID        uuid.UUID                      `json:"id"`
	UserID    *uuid.UUID                     `json:"user_id,omitempty"`
	SessionID string                         `json:"session_id,omitempty"`
	Name      string                         `json:"name"`
	Products  []ProductComparisonItemResponse `json:"products"`
	CreatedAt time.Time                      `json:"created_at"`
	UpdatedAt time.Time                      `json:"updated_at"`
}

// ProductComparisonItemResponse represents a comparison item response
type ProductComparisonItemResponse struct {
	Position int              `json:"position"`
	Product  *ProductResponse `json:"product"`
}

// ComparisonMatrixResponse represents a comparison matrix
type ComparisonMatrixResponse struct {
	Comparison *ProductComparisonResponse `json:"comparison"`
	Matrix     map[string]interface{}     `json:"matrix"`
	Attributes []string                   `json:"attributes"`
}

// ProductComparisonUseCase defines the interface for product comparison operations
type ProductComparisonUseCase interface {
	// Comparison management
	CreateComparison(ctx context.Context, userID *uuid.UUID, sessionID string, req ProductComparisonRequest) (*ProductComparisonResponse, error)
	GetComparison(ctx context.Context, id uuid.UUID) (*ProductComparisonResponse, error)
	GetUserComparison(ctx context.Context, userID uuid.UUID) (*ProductComparisonResponse, error)
	GetSessionComparison(ctx context.Context, sessionID string) (*ProductComparisonResponse, error)
	UpdateComparison(ctx context.Context, id uuid.UUID, req ProductComparisonRequest) (*ProductComparisonResponse, error)
	DeleteComparison(ctx context.Context, id uuid.UUID) error

	// Comparison items management
	AddProductToComparison(ctx context.Context, comparisonID, productID uuid.UUID) (*ProductComparisonResponse, error)
	RemoveProductFromComparison(ctx context.Context, comparisonID, productID uuid.UUID) (*ProductComparisonResponse, error)
	ClearComparison(ctx context.Context, comparisonID uuid.UUID) (*ProductComparisonResponse, error)

	// Comparison queries
	CompareProducts(ctx context.Context, productIDs []uuid.UUID) (*ComparisonMatrixResponse, error)
	GetComparisonMatrix(ctx context.Context, comparisonID uuid.UUID) (*ComparisonMatrixResponse, error)
	GetPopularComparedProducts(ctx context.Context, limit int) ([]*ProductResponse, error)
}

type productComparisonUseCase struct {
	comparisonRepo repositories.ProductComparisonRepository
	productRepo    repositories.ProductRepository
}

// NewProductComparisonUseCase creates a new product comparison use case
func NewProductComparisonUseCase(
	comparisonRepo repositories.ProductComparisonRepository,
	productRepo repositories.ProductRepository,
) ProductComparisonUseCase {
	return &productComparisonUseCase{
		comparisonRepo: comparisonRepo,
		productRepo:    productRepo,
	}
}

// CreateComparison creates a new product comparison
func (uc *productComparisonUseCase) CreateComparison(ctx context.Context, userID *uuid.UUID, sessionID string, req ProductComparisonRequest) (*ProductComparisonResponse, error) {
	if len(req.ProductIDs) < 2 {
		return nil, fmt.Errorf("at least 2 products are required for comparison")
	}
	if len(req.ProductIDs) > 5 {
		return nil, fmt.Errorf("maximum 5 products can be compared at once")
	}

	// Validate all products exist
	for _, productID := range req.ProductIDs {
		_, err := uc.productRepo.GetByID(ctx, productID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found", productID)
		}
	}

	// Check if user/session already has a comparison
	var existingComparison *entities.ProductComparison
	var err error
	
	if userID != nil {
		existingComparison, err = uc.comparisonRepo.GetComparisonByUserID(ctx, *userID)
	} else {
		existingComparison, err = uc.comparisonRepo.GetComparisonBySessionID(ctx, sessionID)
	}

	if err == nil && existingComparison != nil {
		// Clear existing comparison and add new products
		if err := uc.comparisonRepo.ClearComparison(ctx, existingComparison.ID); err != nil {
			return nil, fmt.Errorf("failed to clear existing comparison: %w", err)
		}
		
		// Add new products
		for i, productID := range req.ProductIDs {
			if err := uc.comparisonRepo.AddProductToComparison(ctx, existingComparison.ID, productID, i); err != nil {
				return nil, fmt.Errorf("failed to add product to comparison: %w", err)
			}
		}

		// Update comparison name
		existingComparison.Name = req.Name
		if err := uc.comparisonRepo.UpdateComparison(ctx, existingComparison); err != nil {
			return nil, fmt.Errorf("failed to update comparison: %w", err)
		}

		return uc.GetComparison(ctx, existingComparison.ID)
	}

	// Create new comparison
	comparison := &entities.ProductComparison{
		UserID:    userID,
		SessionID: sessionID,
		Name:      req.Name,
	}

	if err := uc.comparisonRepo.CreateComparison(ctx, comparison); err != nil {
		return nil, fmt.Errorf("failed to create comparison: %w", err)
	}

	// Add products to comparison
	for i, productID := range req.ProductIDs {
		if err := uc.comparisonRepo.AddProductToComparison(ctx, comparison.ID, productID, i); err != nil {
			return nil, fmt.Errorf("failed to add product to comparison: %w", err)
		}
	}

	return uc.GetComparison(ctx, comparison.ID)
}

// GetComparison gets a comparison by ID
func (uc *productComparisonUseCase) GetComparison(ctx context.Context, id uuid.UUID) (*ProductComparisonResponse, error) {
	comparison, err := uc.comparisonRepo.GetComparison(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("comparison not found: %w", err)
	}

	return uc.mapComparisonToResponse(comparison), nil
}

// GetUserComparison gets user's comparison
func (uc *productComparisonUseCase) GetUserComparison(ctx context.Context, userID uuid.UUID) (*ProductComparisonResponse, error) {
	comparison, err := uc.comparisonRepo.GetComparisonByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user comparison not found: %w", err)
	}

	return uc.mapComparisonToResponse(comparison), nil
}

// GetSessionComparison gets session's comparison
func (uc *productComparisonUseCase) GetSessionComparison(ctx context.Context, sessionID string) (*ProductComparisonResponse, error) {
	comparison, err := uc.comparisonRepo.GetComparisonBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session comparison not found: %w", err)
	}

	return uc.mapComparisonToResponse(comparison), nil
}

// UpdateComparison updates a comparison
func (uc *productComparisonUseCase) UpdateComparison(ctx context.Context, id uuid.UUID, req ProductComparisonRequest) (*ProductComparisonResponse, error) {
	comparison, err := uc.comparisonRepo.GetComparison(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("comparison not found: %w", err)
	}

	// Clear existing items
	if err := uc.comparisonRepo.ClearComparison(ctx, id); err != nil {
		return nil, fmt.Errorf("failed to clear comparison: %w", err)
	}

	// Add new products
	for i, productID := range req.ProductIDs {
		if err := uc.comparisonRepo.AddProductToComparison(ctx, id, productID, i); err != nil {
			return nil, fmt.Errorf("failed to add product to comparison: %w", err)
		}
	}

	// Update comparison
	comparison.Name = req.Name
	if err := uc.comparisonRepo.UpdateComparison(ctx, comparison); err != nil {
		return nil, fmt.Errorf("failed to update comparison: %w", err)
	}

	return uc.GetComparison(ctx, id)
}

// DeleteComparison deletes a comparison
func (uc *productComparisonUseCase) DeleteComparison(ctx context.Context, id uuid.UUID) error {
	return uc.comparisonRepo.DeleteComparison(ctx, id)
}

// AddProductToComparison adds a product to comparison
func (uc *productComparisonUseCase) AddProductToComparison(ctx context.Context, comparisonID, productID uuid.UUID) (*ProductComparisonResponse, error) {
	// Check if comparison exists
	_, err := uc.comparisonRepo.GetComparison(ctx, comparisonID)
	if err != nil {
		return nil, fmt.Errorf("comparison not found: %w", err)
	}

	// Check if product exists
	_, err = uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Check comparison limit
	count, err := uc.comparisonRepo.CountComparisonItems(ctx, comparisonID)
	if err != nil {
		return nil, fmt.Errorf("failed to count comparison items: %w", err)
	}
	if count >= 5 {
		return nil, fmt.Errorf("maximum 5 products can be compared at once")
	}

	// Add product
	if err := uc.comparisonRepo.AddProductToComparison(ctx, comparisonID, productID, int(count)); err != nil {
		return nil, fmt.Errorf("failed to add product to comparison: %w", err)
	}

	return uc.GetComparison(ctx, comparisonID)
}

// RemoveProductFromComparison removes a product from comparison
func (uc *productComparisonUseCase) RemoveProductFromComparison(ctx context.Context, comparisonID, productID uuid.UUID) (*ProductComparisonResponse, error) {
	if err := uc.comparisonRepo.RemoveProductFromComparison(ctx, comparisonID, productID); err != nil {
		return nil, fmt.Errorf("failed to remove product from comparison: %w", err)
	}

	return uc.GetComparison(ctx, comparisonID)
}

// ClearComparison clears all products from comparison
func (uc *productComparisonUseCase) ClearComparison(ctx context.Context, comparisonID uuid.UUID) (*ProductComparisonResponse, error) {
	if err := uc.comparisonRepo.ClearComparison(ctx, comparisonID); err != nil {
		return nil, fmt.Errorf("failed to clear comparison: %w", err)
	}

	return uc.GetComparison(ctx, comparisonID)
}

// Helper method to map comparison entity to response
func (uc *productComparisonUseCase) mapComparisonToResponse(comparison *entities.ProductComparison) *ProductComparisonResponse {
	response := &ProductComparisonResponse{
		ID:        comparison.ID,
		UserID:    comparison.UserID,
		SessionID: comparison.SessionID,
		Name:      comparison.Name,
		CreatedAt: comparison.CreatedAt,
		UpdatedAt: comparison.UpdatedAt,
		Products:  make([]ProductComparisonItemResponse, len(comparison.Items)),
	}

	for i, item := range comparison.Items {
		response.Products[i] = ProductComparisonItemResponse{
			Position: item.Position,
			Product:  mapProductToResponse(&item.Product),
		}
	}

	return response
}

// mapProductToResponse converts product entity to response
func mapProductToResponse(product *entities.Product) *ProductResponse {
	if product == nil {
		return nil
	}

	response := &ProductResponse{
		ID:               product.ID,
		Name:             product.Name,
		Description:      product.Description,
		ShortDescription: product.ShortDescription,
		SKU:              product.SKU,
		Slug:             product.Slug,
		MetaTitle:        product.MetaTitle,
		MetaDescription:  product.MetaDescription,
		Keywords:         product.Keywords,
		Featured:         product.Featured,
		Visibility:       product.Visibility,
		Price:            product.Price,
		ComparePrice:     product.ComparePrice,
		CostPrice:        product.CostPrice,
		SalePrice:        product.SalePrice,
		CurrentPrice:     product.GetCurrentPrice(),
		IsOnSale:         product.IsOnSale(),
		SaleDiscountPercentage: product.GetSaleDiscountPercentage(),
		Stock:            product.Stock,
		LowStockThreshold: product.LowStockThreshold,
		TrackQuantity:    product.TrackQuantity,
		AllowBackorder:   product.AllowBackorder,
		StockStatus:      product.StockStatus,
		IsLowStock:       product.IsLowStock(),
		Weight:           product.Weight,
		RequiresShipping: product.RequiresShipping,
		ShippingClass:    product.ShippingClass,
		TaxClass:         product.TaxClass,
		CountryOfOrigin:  product.CountryOfOrigin,
		Status:           product.Status,
		ProductType:      product.ProductType,
		IsDigital:        product.IsDigital,
		IsAvailable:      product.IsAvailable(),
		HasDiscount:      product.HasDiscount(),
		HasVariants:      product.HasVariants(),
		MainImage:        product.GetMainImage(),
		CreatedAt:        product.CreatedAt,
		UpdatedAt:        product.UpdatedAt,
	}

	// Convert category
	if product.Category.ID != uuid.Nil {
		response.Category = &ProductCategoryResponse{
			ID:          product.Category.ID,
			Name:        product.Category.Name,
			Description: product.Category.Description,
			Slug:        product.Category.Slug,
			Image:       product.Category.Image,
		}
	}

	// Convert brand
	if product.Brand != nil {
		response.Brand = &ProductBrandResponse{
			ID:          product.Brand.ID,
			Name:        product.Brand.Name,
			Description: product.Brand.Description,
			Slug:        product.Brand.Slug,
			Logo:        product.Brand.Logo,
		}
	}

	// Convert images
	for _, img := range product.Images {
		response.Images = append(response.Images, ProductImageResponse{
			ID:       img.ID,
			URL:      img.URL,
			AltText:  img.AltText,
			Position: img.Position,
		})
	}

	// Convert tags
	for _, tag := range product.Tags {
		response.Tags = append(response.Tags, ProductTagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	return response
}

// CompareProducts creates a temporary comparison for given product IDs
func (uc *productComparisonUseCase) CompareProducts(ctx context.Context, productIDs []uuid.UUID) (*ComparisonMatrixResponse, error) {
	if len(productIDs) < 2 {
		return nil, fmt.Errorf("at least 2 products are required for comparison")
	}
	if len(productIDs) > 5 {
		return nil, fmt.Errorf("maximum 5 products can be compared at once")
	}

	// Get all products
	products := make([]*entities.Product, len(productIDs))
	for i, productID := range productIDs {
		product, err := uc.productRepo.GetByID(ctx, productID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found", productID)
		}
		products[i] = product
	}

	// Create comparison response
	comparisonResponse := &ProductComparisonResponse{
		ID:       uuid.New(), // Temporary ID
		Name:     "Product Comparison",
		Products: make([]ProductComparisonItemResponse, len(products)),
	}

	for i, product := range products {
		comparisonResponse.Products[i] = ProductComparisonItemResponse{
			Position: i,
			Product:  mapProductToResponse(product),
		}
	}

	// Generate comparison matrix
	matrix := uc.generateComparisonMatrix(products)
	attributes := uc.getComparisonAttributes()

	return &ComparisonMatrixResponse{
		Comparison: comparisonResponse,
		Matrix:     matrix,
		Attributes: attributes,
	}, nil
}

// GetComparisonMatrix gets comparison matrix for existing comparison
func (uc *productComparisonUseCase) GetComparisonMatrix(ctx context.Context, comparisonID uuid.UUID) (*ComparisonMatrixResponse, error) {
	comparison, err := uc.comparisonRepo.GetComparison(ctx, comparisonID)
	if err != nil {
		return nil, fmt.Errorf("comparison not found: %w", err)
	}

	// Extract products from comparison items
	products := make([]*entities.Product, len(comparison.Items))
	for i, item := range comparison.Items {
		products[i] = &item.Product
	}

	// Generate comparison matrix
	matrix := uc.generateComparisonMatrix(products)
	attributes := uc.getComparisonAttributes()

	return &ComparisonMatrixResponse{
		Comparison: uc.mapComparisonToResponse(comparison),
		Matrix:     matrix,
		Attributes: attributes,
	}, nil
}

// GetPopularComparedProducts gets most compared products
func (uc *productComparisonUseCase) GetPopularComparedProducts(ctx context.Context, limit int) ([]*ProductResponse, error) {
	products, err := uc.comparisonRepo.GetPopularComparedProducts(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular compared products: %w", err)
	}

	responses := make([]*ProductResponse, len(products))
	for i, product := range products {
		responses[i] = mapProductToResponse(&product)
	}

	return responses, nil
}

// generateComparisonMatrix generates a comparison matrix for products
func (uc *productComparisonUseCase) generateComparisonMatrix(products []*entities.Product) map[string]interface{} {
	matrix := make(map[string]interface{})

	if len(products) == 0 {
		return matrix
	}

	// Basic product information
	matrix["names"] = make([]string, len(products))
	matrix["prices"] = make([]float64, len(products))
	matrix["current_prices"] = make([]float64, len(products))
	matrix["categories"] = make([]string, len(products))
	matrix["brands"] = make([]string, len(products))
	matrix["stock"] = make([]int, len(products))
	matrix["stock_status"] = make([]string, len(products))
	matrix["ratings"] = make([]float64, len(products))
	matrix["is_on_sale"] = make([]bool, len(products))
	matrix["sale_discount"] = make([]float64, len(products))

	for i, product := range products {
		matrix["names"].([]string)[i] = product.Name
		matrix["prices"].([]float64)[i] = product.Price
		matrix["current_prices"].([]float64)[i] = product.GetCurrentPrice()
		matrix["categories"].([]string)[i] = product.Category.Name
		if product.Brand != nil {
			matrix["brands"].([]string)[i] = product.Brand.Name
		} else {
			matrix["brands"].([]string)[i] = ""
		}
		matrix["stock"].([]int)[i] = product.Stock
		matrix["stock_status"].([]string)[i] = string(product.StockStatus)
		matrix["ratings"].([]float64)[i] = 0 // TODO: Calculate from reviews
		matrix["is_on_sale"].([]bool)[i] = product.IsOnSale()
		matrix["sale_discount"].([]float64)[i] = product.GetSaleDiscountPercentage()
	}

	return matrix
}

// getComparisonAttributes returns the list of attributes used in comparison
func (uc *productComparisonUseCase) getComparisonAttributes() []string {
	return []string{
		"name",
		"price",
		"current_price",
		"category",
		"brand",
		"stock",
		"stock_status",
		"rating",
		"is_on_sale",
		"sale_discount",
	}
}
