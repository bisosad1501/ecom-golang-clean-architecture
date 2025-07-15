package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

// AdvancedFilterRequest represents an advanced filter request
type AdvancedFilterRequest struct {
	// Basic filters
	Query       string     `json:"query"`
	CategoryIDs []string   `json:"category_ids"`
	BrandIDs    []string   `json:"brand_ids"`
	MinPrice    *float64   `json:"min_price"`
	MaxPrice    *float64   `json:"max_price"`
	MinRating   *float64   `json:"min_rating"`
	MaxRating   *float64   `json:"max_rating"`
	
	// Stock and availability
	InStock     *bool `json:"in_stock"`
	LowStock    *bool `json:"low_stock"`
	OnSale      *bool `json:"on_sale"`
	Featured    *bool `json:"featured"`
	
	// Product properties
	ProductTypes []string `json:"product_types"`
	StockStatus  []string `json:"stock_status"`
	Visibility   []string `json:"visibility"`
	
	// Custom attributes
	Attributes map[string][]string `json:"attributes"` // AttributeID -> Values
	
	// Date filters
	CreatedAfter  *string `json:"created_after"`
	CreatedBefore *string `json:"created_before"`
	UpdatedAfter  *string `json:"updated_after"`
	UpdatedBefore *string `json:"updated_before"`
	
	// Advanced options
	Tags         []string `json:"tags"`
	HasImages    *bool    `json:"has_images"`
	HasVariants  *bool    `json:"has_variants"`
	HasReviews   *bool    `json:"has_reviews"`
	
	// Sorting and pagination
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	
	// Filter options
	IncludeFacets bool `json:"include_facets"`
	FacetLimit    int  `json:"facet_limit"`
}

// FilteredProductResponse represents filtered product response
type FilteredProductResponse struct {
	Products   []*ProductResponse         `json:"products"`
	Facets     *repositories.FilterFacets `json:"facets,omitempty"`
	Pagination *PaginationInfo            `json:"pagination"`
}

// FilterSetRequest represents a filter set request
type FilterSetRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description"`
	Filters     AdvancedFilterRequest  `json:"filters" validate:"required"`
	IsPublic    bool                   `json:"is_public"`
}

// FilterSetResponse represents a filter set response
type FilterSetResponse struct {
	ID          uuid.UUID             `json:"id"`
	UserID      *uuid.UUID            `json:"user_id,omitempty"`
	SessionID   string                `json:"session_id,omitempty"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Filters     AdvancedFilterRequest `json:"filters"`
	IsPublic    bool                  `json:"is_public"`
	UsageCount  int                   `json:"usage_count"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// ProductFilterUseCase defines the interface for advanced product filtering
type ProductFilterUseCase interface {
	// Advanced filtering
	FilterProducts(ctx context.Context, req AdvancedFilterRequest) (*FilteredProductResponse, error)
	GetFilterFacets(ctx context.Context, categoryID *string) (*repositories.FilterFacets, error)
	GetDynamicFilters(ctx context.Context, req AdvancedFilterRequest) (*repositories.FilterFacets, error)
	
	// Filter sets management
	SaveFilterSet(ctx context.Context, userID *uuid.UUID, sessionID string, req FilterSetRequest) (*FilterSetResponse, error)
	GetFilterSet(ctx context.Context, id uuid.UUID) (*FilterSetResponse, error)
	GetUserFilterSets(ctx context.Context, userID uuid.UUID) ([]*FilterSetResponse, error)
	GetSessionFilterSets(ctx context.Context, sessionID string) ([]*FilterSetResponse, error)
	UpdateFilterSet(ctx context.Context, id uuid.UUID, req FilterSetRequest) (*FilterSetResponse, error)
	DeleteFilterSet(ctx context.Context, id uuid.UUID) error
	
	// Filter analytics
	TrackFilterUsage(ctx context.Context, userID *uuid.UUID, sessionID string, filterType, filterKey, filterValue string, resultCount int) error
	GetFilterAnalytics(ctx context.Context, days int) (map[string]interface{}, error)
	GetPopularFilters(ctx context.Context, limit int) ([]*entities.FilterUsage, error)
	
	// Filter suggestions
	GetFilterSuggestions(ctx context.Context, query string, limit int) ([]string, error)
	GetRelatedFilters(ctx context.Context, req AdvancedFilterRequest) ([]string, error)
	
	// Attribute management
	GetAttributeFilters(ctx context.Context, categoryID *string) ([]*entities.ProductAttribute, error)
	GetAttributeTerms(ctx context.Context, attributeID uuid.UUID, categoryID *string) ([]*entities.ProductAttributeTerm, error)
}

type productFilterUseCase struct {
	filterRepo  repositories.ProductFilterRepository
	productRepo repositories.ProductRepository
}

// NewProductFilterUseCase creates a new product filter use case
func NewProductFilterUseCase(
	filterRepo repositories.ProductFilterRepository,
	productRepo repositories.ProductRepository,
) ProductFilterUseCase {
	return &productFilterUseCase{
		filterRepo:  filterRepo,
		productRepo: productRepo,
	}
}

// FilterProducts performs advanced product filtering
func (uc *productFilterUseCase) FilterProducts(ctx context.Context, req AdvancedFilterRequest) (*FilteredProductResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (req.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// Convert request to repository parameters
	params := uc.convertToRepositoryParams(req)
	params.Limit = limit
	params.Offset = offset

	result, err := uc.filterRepo.FilterProducts(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	products := result.Products

	// Convert products to response format using the proper mapping function
	convertedProducts := make([]*ProductResponse, len(products))
	for i, product := range products {
		convertedProducts[i] = uc.mapProductToResponse(product)
	}

	// Create pagination context
	context := &EcommercePaginationContext{
		EntityType:  "products",
		SearchQuery: req.Query,
	}

	// Create enhanced pagination info
	pagination := NewEcommercePaginationInfo(req.Page, req.Limit, result.Total, context)

	response := &FilteredProductResponse{
		Products:   convertedProducts,
		Facets:     result.Facets,
		Pagination: pagination,
	}
	
	return response, nil
}

// GetFilterFacets gets available filter facets
func (uc *productFilterUseCase) GetFilterFacets(ctx context.Context, categoryID *string) (*repositories.FilterFacets, error) {
	var categoryUUID *uuid.UUID
	if categoryID != nil && *categoryID != "" {
		if id, err := uuid.Parse(*categoryID); err == nil {
			categoryUUID = &id
		}
	}
	
	return uc.filterRepo.GetFilterFacets(ctx, categoryUUID)
}

// GetDynamicFilters gets dynamic filters based on current state
func (uc *productFilterUseCase) GetDynamicFilters(ctx context.Context, req AdvancedFilterRequest) (*repositories.FilterFacets, error) {
	params := uc.convertToRepositoryParams(req)
	return uc.filterRepo.GetDynamicFilters(ctx, params)
}

// SaveFilterSet saves a filter set
func (uc *productFilterUseCase) SaveFilterSet(ctx context.Context, userID *uuid.UUID, sessionID string, req FilterSetRequest) (*FilterSetResponse, error) {
	// Convert filters to JSON
	filtersJSON, err := json.Marshal(req.Filters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filters: %w", err)
	}
	
	filterSet := &entities.FilterSet{
		UserID:      userID,
		SessionID:   sessionID,
		Name:        req.Name,
		Description: req.Description,
		Filters:     string(filtersJSON),
		IsPublic:    req.IsPublic,
	}
	
	if err := uc.filterRepo.SaveFilterSet(ctx, filterSet); err != nil {
		return nil, fmt.Errorf("failed to save filter set: %w", err)
	}
	
	return uc.mapFilterSetToResponse(filterSet), nil
}

// GetFilterSet gets a filter set by ID
func (uc *productFilterUseCase) GetFilterSet(ctx context.Context, id uuid.UUID) (*FilterSetResponse, error) {
	filterSet, err := uc.filterRepo.GetFilterSet(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("filter set not found: %w", err)
	}
	
	return uc.mapFilterSetToResponse(filterSet), nil
}

// GetUserFilterSets gets filter sets for a user
func (uc *productFilterUseCase) GetUserFilterSets(ctx context.Context, userID uuid.UUID) ([]*FilterSetResponse, error) {
	filterSets, err := uc.filterRepo.GetUserFilterSets(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user filter sets: %w", err)
	}
	
	responses := make([]*FilterSetResponse, len(filterSets))
	for i, filterSet := range filterSets {
		responses[i] = uc.mapFilterSetToResponse(filterSet)
	}
	
	return responses, nil
}

// GetSessionFilterSets gets filter sets for a session
func (uc *productFilterUseCase) GetSessionFilterSets(ctx context.Context, sessionID string) ([]*FilterSetResponse, error) {
	filterSets, err := uc.filterRepo.GetSessionFilterSets(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session filter sets: %w", err)
	}
	
	responses := make([]*FilterSetResponse, len(filterSets))
	for i, filterSet := range filterSets {
		responses[i] = uc.mapFilterSetToResponse(filterSet)
	}
	
	return responses, nil
}

// UpdateFilterSet updates a filter set
func (uc *productFilterUseCase) UpdateFilterSet(ctx context.Context, id uuid.UUID, req FilterSetRequest) (*FilterSetResponse, error) {
	filterSet, err := uc.filterRepo.GetFilterSet(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("filter set not found: %w", err)
	}
	
	// Convert filters to JSON
	filtersJSON, err := json.Marshal(req.Filters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filters: %w", err)
	}
	
	filterSet.Name = req.Name
	filterSet.Description = req.Description
	filterSet.Filters = string(filtersJSON)
	filterSet.IsPublic = req.IsPublic
	
	if err := uc.filterRepo.UpdateFilterSet(ctx, filterSet); err != nil {
		return nil, fmt.Errorf("failed to update filter set: %w", err)
	}
	
	return uc.mapFilterSetToResponse(filterSet), nil
}

// DeleteFilterSet deletes a filter set
func (uc *productFilterUseCase) DeleteFilterSet(ctx context.Context, id uuid.UUID) error {
	return uc.filterRepo.DeleteFilterSet(ctx, id)
}

// TrackFilterUsage tracks filter usage
func (uc *productFilterUseCase) TrackFilterUsage(ctx context.Context, userID *uuid.UUID, sessionID string, filterType, filterKey, filterValue string, resultCount int) error {
	usage := &entities.FilterUsage{
		UserID:      userID,
		SessionID:   sessionID,
		FilterType:  filterType,
		FilterKey:   filterKey,
		FilterValue: filterValue,
		ResultCount: resultCount,
	}

	return uc.filterRepo.TrackFilterUsage(ctx, usage)
}

// GetFilterAnalytics gets filter analytics
func (uc *productFilterUseCase) GetFilterAnalytics(ctx context.Context, days int) (map[string]interface{}, error) {
	return uc.filterRepo.GetFilterAnalytics(ctx, days)
}

// GetPopularFilters gets popular filters
func (uc *productFilterUseCase) GetPopularFilters(ctx context.Context, limit int) ([]*entities.FilterUsage, error) {
	return uc.filterRepo.GetPopularFilters(ctx, limit)
}

// GetFilterSuggestions gets filter suggestions
func (uc *productFilterUseCase) GetFilterSuggestions(ctx context.Context, query string, limit int) ([]string, error) {
	return uc.filterRepo.GetFilterSuggestions(ctx, query, limit)
}

// GetRelatedFilters gets related filters
func (uc *productFilterUseCase) GetRelatedFilters(ctx context.Context, req AdvancedFilterRequest) ([]string, error) {
	params := uc.convertToRepositoryParams(req)
	return uc.filterRepo.GetRelatedFilters(ctx, params)
}

// GetAttributeFilters gets attribute filters
func (uc *productFilterUseCase) GetAttributeFilters(ctx context.Context, categoryID *string) ([]*entities.ProductAttribute, error) {
	var categoryUUID *uuid.UUID
	if categoryID != nil && *categoryID != "" {
		if id, err := uuid.Parse(*categoryID); err == nil {
			categoryUUID = &id
		}
	}

	return uc.filterRepo.GetAttributeFilters(ctx, categoryUUID)
}

// GetAttributeTerms gets attribute terms
func (uc *productFilterUseCase) GetAttributeTerms(ctx context.Context, attributeID uuid.UUID, categoryID *string) ([]*entities.ProductAttributeTerm, error) {
	var categoryUUID *uuid.UUID
	if categoryID != nil && *categoryID != "" {
		if id, err := uuid.Parse(*categoryID); err == nil {
			categoryUUID = &id
		}
	}

	return uc.filterRepo.GetAttributeTerms(ctx, attributeID, categoryUUID)
}

// Helper method to convert request to repository parameters
func (uc *productFilterUseCase) convertToRepositoryParams(req AdvancedFilterRequest) repositories.AdvancedFilterParams {
	params := repositories.AdvancedFilterParams{
		Query:         req.Query,
		MinPrice:      req.MinPrice,
		MaxPrice:      req.MaxPrice,
		MinRating:     req.MinRating,
		MaxRating:     req.MaxRating,
		InStock:       req.InStock,
		LowStock:      req.LowStock,
		OnSale:        req.OnSale,
		Featured:      req.Featured,
		CreatedAfter:  req.CreatedAfter,
		CreatedBefore: req.CreatedBefore,
		UpdatedAfter:  req.UpdatedAfter,
		UpdatedBefore: req.UpdatedBefore,
		Tags:          req.Tags,
		HasImages:     req.HasImages,
		HasVariants:   req.HasVariants,
		HasReviews:    req.HasReviews,
		SortBy:        req.SortBy,
		SortOrder:     req.SortOrder,
		IncludeFacets: req.IncludeFacets,
		FacetLimit:    req.FacetLimit,
	}

	// Convert string IDs to UUIDs
	for _, idStr := range req.CategoryIDs {
		if id, err := uuid.Parse(idStr); err == nil {
			params.CategoryIDs = append(params.CategoryIDs, id)
		}
	}

	for _, idStr := range req.BrandIDs {
		if id, err := uuid.Parse(idStr); err == nil {
			params.BrandIDs = append(params.BrandIDs, id)
		}
	}

	// Convert product types
	for _, typeStr := range req.ProductTypes {
		params.ProductTypes = append(params.ProductTypes, entities.ProductType(typeStr))
	}

	// Convert stock status
	for _, statusStr := range req.StockStatus {
		params.StockStatus = append(params.StockStatus, entities.StockStatus(statusStr))
	}

	// Convert visibility
	for _, visStr := range req.Visibility {
		params.Visibility = append(params.Visibility, entities.ProductVisibility(visStr))
	}

	// Convert attributes
	params.Attributes = make(map[uuid.UUID][]string)
	for attrIDStr, values := range req.Attributes {
		if attrID, err := uuid.Parse(attrIDStr); err == nil {
			params.Attributes[attrID] = values
		}
	}

	// Set pagination
	if req.Page > 0 && req.Limit > 0 {
		params.Offset = (req.Page - 1) * req.Limit
		params.Limit = req.Limit
	} else {
		params.Limit = 20
		params.Offset = 0
	}

	return params
}

// Helper method to map product to response
func (uc *productFilterUseCase) mapProductToResponse(product *entities.Product) *ProductResponse {
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

// Helper method to map filter set to response
func (uc *productFilterUseCase) mapFilterSetToResponse(filterSet *entities.FilterSet) *FilterSetResponse {
	var filters AdvancedFilterRequest
	json.Unmarshal([]byte(filterSet.Filters), &filters)

	return &FilterSetResponse{
		ID:          filterSet.ID,
		UserID:      filterSet.UserID,
		SessionID:   filterSet.SessionID,
		Name:        filterSet.Name,
		Description: filterSet.Description,
		Filters:     filters,
		IsPublic:    filterSet.IsPublic,
		UsageCount:  filterSet.UsageCount,
		CreatedAt:   filterSet.CreatedAt,
		UpdatedAt:   filterSet.UpdatedAt,
	}
}
