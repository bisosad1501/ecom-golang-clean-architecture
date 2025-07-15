package usecases

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

// WishlistUseCase defines wishlist use cases
type WishlistUseCase interface {
	AddToWishlist(ctx context.Context, userID, productID uuid.UUID) error
	RemoveFromWishlist(ctx context.Context, userID, productID uuid.UUID) error
	GetWishlist(ctx context.Context, userID uuid.UUID, req GetWishlistRequest) (*WishlistResponse, error)
	IsInWishlist(ctx context.Context, userID, productID uuid.UUID) (bool, error)
	ClearWishlist(ctx context.Context, userID uuid.UUID) error
	GetWishlistCount(ctx context.Context, userID uuid.UUID) (int64, error)
}

type wishlistUseCase struct {
	wishlistRepo repositories.WishlistRepository
	productRepo  repositories.ProductRepository
}

// NewWishlistUseCase creates a new wishlist use case
func NewWishlistUseCase(
	wishlistRepo repositories.WishlistRepository,
	productRepo repositories.ProductRepository,
) WishlistUseCase {
	return &wishlistUseCase{
		wishlistRepo: wishlistRepo,
		productRepo:  productRepo,
	}
}

// GetWishlistRequest represents get wishlist request
type GetWishlistRequest struct {
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}

// WishlistItemResponse represents wishlist item response
type WishlistItemResponse struct {
	ID        uuid.UUID       `json:"id"`
	Product   ProductResponse `json:"product"`
	AddedAt   time.Time       `json:"added_at"`
}

// WishlistResponse represents wishlist response
type WishlistResponse struct {
	Items      []*WishlistItemResponse `json:"items"`
	Pagination *PaginationInfo         `json:"pagination"`
}

// AddToWishlist adds a product to user's wishlist
func (uc *wishlistUseCase) AddToWishlist(ctx context.Context, userID, productID uuid.UUID) error {
	// Check if product exists
	_, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return entities.ErrProductNotFound
	}

	// Check if already in wishlist
	exists, err := uc.wishlistRepo.IsInWishlist(ctx, userID, productID)
	if err != nil {
		return err
	}
	if exists {
		return entities.ErrConflict // Already in wishlist
	}

	return uc.wishlistRepo.AddToWishlist(ctx, userID, productID)
}

// RemoveFromWishlist removes a product from user's wishlist
func (uc *wishlistUseCase) RemoveFromWishlist(ctx context.Context, userID, productID uuid.UUID) error {
	// Check if in wishlist
	exists, err := uc.wishlistRepo.IsInWishlist(ctx, userID, productID)
	if err != nil {
		return err
	}
	if !exists {
		return entities.ErrWishlistItemNotFound
	}

	return uc.wishlistRepo.RemoveFromWishlist(ctx, userID, productID)
}

// GetWishlist gets user's wishlist with pagination
func (uc *wishlistUseCase) GetWishlist(ctx context.Context, userID uuid.UUID, req GetWishlistRequest) (*WishlistResponse, error) {
	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Get wishlist items
	wishlistItems, err := uc.wishlistRepo.GetByUserID(ctx, userID, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	// Get total count
	totalCount, err := uc.wishlistRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	items := make([]*WishlistItemResponse, len(wishlistItems))
	for i, item := range wishlistItems {
		productResponse := &ProductResponse{
			ID:          item.Product.ID,
			Name:        item.Product.Name,
			Description: item.Product.Description,
			SKU:         item.Product.SKU,
			Price:       item.Product.Price,
			Stock:       item.Product.Stock,
			Status:      item.Product.Status,
			IsDigital:   item.Product.IsDigital,
			Weight:      item.Product.Weight,
			CreatedAt:   item.Product.CreatedAt,
			UpdatedAt:   item.Product.UpdatedAt,
		}

		// Add category if available
		if item.Product.Category.ID != uuid.Nil {
			productResponse.Category = &ProductCategoryResponse{
				ID:          item.Product.Category.ID,
				Name:        item.Product.Category.Name,
				Description: item.Product.Category.Description,
				Slug:        item.Product.Category.Slug,
				Image:       item.Product.Category.Image,
			}
		}

		// Add images
		if len(item.Product.Images) > 0 {
			images := make([]ProductImageResponse, len(item.Product.Images))
			for j, img := range item.Product.Images {
				images[j] = ProductImageResponse{
					ID:      img.ID,
					URL:     img.URL,
					AltText: img.AltText,
				}
			}
			productResponse.Images = images
		}

		// Add tags
		if len(item.Product.Tags) > 0 {
			tags := make([]ProductTagResponse, len(item.Product.Tags))
			for j, tag := range item.Product.Tags {
				tags[j] = ProductTagResponse{
					ID:   tag.ID,
					Name: tag.Name,
					Slug: tag.Slug,
				}
			}
			productResponse.Tags = tags
		}

		items[i] = &WishlistItemResponse{
			ID:      item.ID,
			Product: *productResponse,
			AddedAt: item.CreatedAt,
		}
	}

	// Create pagination info using enhanced function
	context := &EcommercePaginationContext{
		EntityType: "wishlist",
		UserID:     userID.String(),
	}
	pagination := NewEcommercePaginationInfo((req.Offset/req.Limit)+1, req.Limit, totalCount, context)

	return &WishlistResponse{
		Items:      items,
		Pagination: pagination,
	}, nil
}

// IsInWishlist checks if a product is in user's wishlist
func (uc *wishlistUseCase) IsInWishlist(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	return uc.wishlistRepo.IsInWishlist(ctx, userID, productID)
}

// ClearWishlist removes all items from user's wishlist
func (uc *wishlistUseCase) ClearWishlist(ctx context.Context, userID uuid.UUID) error {
	return uc.wishlistRepo.ClearWishlist(ctx, userID)
}

// GetWishlistCount gets the total count of items in user's wishlist
func (uc *wishlistUseCase) GetWishlistCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return uc.wishlistRepo.CountByUserID(ctx, userID)
}
