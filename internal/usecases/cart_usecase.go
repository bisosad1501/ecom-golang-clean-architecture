package usecases

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

// CartUseCase defines cart use cases
type CartUseCase interface {
	GetCart(ctx context.Context, userID uuid.UUID) (*CartResponse, error)
	AddToCart(ctx context.Context, userID uuid.UUID, req AddToCartRequest) (*CartResponse, error)
	UpdateCartItem(ctx context.Context, userID uuid.UUID, req UpdateCartItemRequest) (*CartResponse, error)
	RemoveFromCart(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (*CartResponse, error)
	ClearCart(ctx context.Context, userID uuid.UUID) error
}

type cartUseCase struct {
	cartRepo    repositories.CartRepository
	productRepo repositories.ProductRepository
}

// NewCartUseCase creates a new cart use case
func NewCartUseCase(
	cartRepo repositories.CartRepository,
	productRepo repositories.ProductRepository,
) CartUseCase {
	return &cartUseCase{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// AddToCartRequest represents add to cart request
type AddToCartRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
}

// UpdateCartItemRequest represents update cart item request
type UpdateCartItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
}

// CartResponse represents cart response
type CartResponse struct {
	ID         uuid.UUID          `json:"id"`
	UserID     uuid.UUID          `json:"user_id"`
	Items      []CartItemResponse `json:"items"`
	ItemCount  int                `json:"item_count"`
	Total      float64            `json:"total"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

// CartItemResponse represents cart item response
type CartItemResponse struct {
	ID        uuid.UUID        `json:"id"`
	Product   *ProductResponse `json:"product"`
	Quantity  int              `json:"quantity"`
	Price     float64          `json:"price"`
	Subtotal  float64          `json:"subtotal"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// GetCart gets user's cart
func (uc *cartUseCase) GetCart(ctx context.Context, userID uuid.UUID) (*CartResponse, error) {
	cart, err := uc.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Create new cart if not exists
		cart = &entities.Cart{
			ID:        uuid.New(),
			UserID:    userID,
			Items:     []entities.CartItem{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		if err := uc.cartRepo.Create(ctx, cart); err != nil {
			return nil, err
		}
	}

	return uc.toCartResponse(cart), nil
}

// AddToCart adds item to cart
func (uc *cartUseCase) AddToCart(ctx context.Context, userID uuid.UUID, req AddToCartRequest) (*CartResponse, error) {
	// Get or create cart
	cart, err := uc.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		cart = &entities.Cart{
			ID:        uuid.New(),
			UserID:    userID,
			Items:     []entities.CartItem{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		if err := uc.cartRepo.Create(ctx, cart); err != nil {
			return nil, err
		}
	}

	// Get product
	product, err := uc.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, entities.ErrProductNotFound
	}

	// Check if product is available
	if !product.IsAvailable() {
		return nil, entities.ErrProductNotAvailable
	}

	// Check stock
	if !product.CanReduceStock(req.Quantity) {
		return nil, entities.ErrInsufficientStock
	}

	// Check if item already exists in cart
	existingItem := cart.GetItem(req.ProductID)
	if existingItem != nil {
		// Update quantity
		newQuantity := existingItem.Quantity + req.Quantity
		if !product.CanReduceStock(newQuantity) {
			return nil, entities.ErrInsufficientStock
		}
		
		existingItem.Quantity = newQuantity
		existingItem.UpdatedAt = time.Now()
		
		if err := uc.cartRepo.UpdateItem(ctx, existingItem); err != nil {
			return nil, err
		}
	} else {
		// Add new item
		cartItem := &entities.CartItem{
			ID:        uuid.New(),
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			Price:     product.Price,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		if err := uc.cartRepo.AddItem(ctx, cart.ID, cartItem); err != nil {
			return nil, err
		}
	}

	// Get updated cart
	updatedCart, err := uc.cartRepo.GetByID(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	return uc.toCartResponse(updatedCart), nil
}

// UpdateCartItem updates cart item quantity
func (uc *cartUseCase) UpdateCartItem(ctx context.Context, userID uuid.UUID, req UpdateCartItemRequest) (*CartResponse, error) {
	// Get cart
	cart, err := uc.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, entities.ErrCartNotFound
	}

	// Get product
	product, err := uc.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, entities.ErrProductNotFound
	}

	// Check stock
	if !product.CanReduceStock(req.Quantity) {
		return nil, entities.ErrInsufficientStock
	}

	// Get cart item
	cartItem, err := uc.cartRepo.GetItem(ctx, cart.ID, req.ProductID)
	if err != nil {
		return nil, entities.ErrCartItemNotFound
	}

	// Update quantity
	cartItem.Quantity = req.Quantity
	cartItem.UpdatedAt = time.Now()

	if err := uc.cartRepo.UpdateItem(ctx, cartItem); err != nil {
		return nil, err
	}

	// Get updated cart
	updatedCart, err := uc.cartRepo.GetByID(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	return uc.toCartResponse(updatedCart), nil
}

// RemoveFromCart removes item from cart
func (uc *cartUseCase) RemoveFromCart(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (*CartResponse, error) {
	// Get cart
	cart, err := uc.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, entities.ErrCartNotFound
	}

	// Remove item
	if err := uc.cartRepo.RemoveItem(ctx, cart.ID, productID); err != nil {
		return nil, err
	}

	// Get updated cart
	updatedCart, err := uc.cartRepo.GetByID(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	return uc.toCartResponse(updatedCart), nil
}

// ClearCart clears all items from cart
func (uc *cartUseCase) ClearCart(ctx context.Context, userID uuid.UUID) error {
	// Get cart
	cart, err := uc.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return entities.ErrCartNotFound
	}

	return uc.cartRepo.ClearCart(ctx, cart.ID)
}

// toCartResponse converts cart entity to response
func (uc *cartUseCase) toCartResponse(cart *entities.Cart) *CartResponse {
	response := &CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		ItemCount: cart.GetItemCount(),
		Total:     cart.GetTotal(),
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}

	// Convert items
	response.Items = make([]CartItemResponse, len(cart.Items))
	for i, item := range cart.Items {
		response.Items[i] = CartItemResponse{
			ID:        item.ID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Subtotal:  item.GetSubtotal(),
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}

		// Add product info if available
		if item.Product.ID != uuid.Nil {
			productUseCase := &productUseCase{}
			response.Items[i].Product = productUseCase.toProductResponse(&item.Product)
		}
	}

	return response
}
