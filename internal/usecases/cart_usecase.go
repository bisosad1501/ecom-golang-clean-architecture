package usecases

import (
	"context"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/domain/services"
	pkgErrors "ecom-golang-clean-architecture/pkg/errors"

	"github.com/google/uuid"
)

// MergeStrategy defines how to merge guest cart with user cart
type MergeStrategy string

const (
	MergeStrategyAuto     MergeStrategy = "auto"      // Auto merge (current behavior)
	MergeStrategyReplace  MergeStrategy = "replace"   // Replace user cart with guest cart
	MergeStrategyKeepUser MergeStrategy = "keep_user" // Keep user cart, discard guest cart
	MergeStrategyMerge    MergeStrategy = "merge"     // Merge items (add quantities)
)

// CartConflictInfo represents information about cart merge conflicts
type CartConflictInfo struct {
	HasConflict      bool              `json:"has_conflict"`
	UserCartExists   bool              `json:"user_cart_exists"`
	GuestCartExists  bool              `json:"guest_cart_exists"`
	ConflictingItems []ConflictingItem `json:"conflicting_items,omitempty"`
	UserCart         *CartResponse     `json:"user_cart,omitempty"`
	GuestCart        *CartResponse     `json:"guest_cart,omitempty"`
	Recommendations  []string          `json:"recommendations,omitempty"`
}

// ConflictingItem represents an item that exists in both carts
type ConflictingItem struct {
	ProductID       string  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	UserQuantity    int     `json:"user_quantity"`
	GuestQuantity   int     `json:"guest_quantity"`
	UserPrice       float64 `json:"user_price"`
	GuestPrice      float64 `json:"guest_price"`
	PriceDifference float64 `json:"price_difference"`
}

// CartUseCase defines cart use cases
type CartUseCase interface {
	GetCart(ctx context.Context, userID uuid.UUID) (*CartResponse, error)
	GetGuestCart(ctx context.Context, sessionID string) (*CartResponse, error)
	AddToCart(ctx context.Context, userID uuid.UUID, req AddToCartRequest) (*CartResponse, error)
	AddToGuestCart(ctx context.Context, sessionID string, req AddToCartRequest) (*CartResponse, error)
	UpdateCartItem(ctx context.Context, userID uuid.UUID, req UpdateCartItemRequest) (*CartResponse, error)
	RemoveFromCart(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (*CartResponse, error)
	ClearCart(ctx context.Context, userID uuid.UUID) error
	MergeGuestCart(ctx context.Context, userID uuid.UUID, sessionID string) (*CartResponse, error)
	MergeGuestCartWithStrategy(ctx context.Context, userID uuid.UUID, sessionID string, strategy MergeStrategy) (*CartResponse, error)
	CheckMergeConflict(ctx context.Context, userID uuid.UUID, sessionID string) (*CartConflictInfo, error)

	// Cleanup methods
	CleanupExpiredCarts(ctx context.Context) error
	CleanupExpiredStockReservations(ctx context.Context) error
}

type cartUseCase struct {
	cartRepo                repositories.CartRepository
	productRepo             repositories.ProductRepository
	simpleStockService      services.SimpleStockService
}

// NewCartUseCase creates a new cart use case
func NewCartUseCase(
	cartRepo repositories.CartRepository,
	productRepo repositories.ProductRepository,
	simpleStockService services.SimpleStockService,
) CartUseCase {
	return &cartUseCase{
		cartRepo:                cartRepo,
		productRepo:             productRepo,
		simpleStockService:      simpleStockService,
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
	ID             uuid.UUID          `json:"id"`
	UserID         *uuid.UUID         `json:"user_id,omitempty"` // Nullable for guest carts
	SessionID      *string            `json:"session_id,omitempty"`
	Items          []CartItemResponse `json:"items"`
	ItemCount      int                `json:"item_count"`
	Subtotal       float64            `json:"subtotal"`
	TaxAmount      float64            `json:"tax_amount"`      // Added missing field
	ShippingAmount float64            `json:"shipping_amount"` // Added missing field
	Total          float64            `json:"total"`
	Status         string             `json:"status"`
	Currency       string             `json:"currency"`
	Notes          string             `json:"notes,omitempty"`
	ExpiresAt      *time.Time         `json:"expires_at,omitempty"`
	IsGuest        bool               `json:"is_guest"` // Added helper field
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
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
			UserID:    &userID,
			Items:     []entities.CartItem{},
			Status:    "active",
			Currency:  "USD",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Set expiration and update calculated fields
		cart.SetExpiration()
		cart.UpdateCalculatedFields()

		if err := uc.cartRepo.Create(ctx, cart); err != nil {
			return nil, err
		}
	}

	return uc.toCartResponse(cart), nil
}

// GetGuestCart gets guest cart by session ID
func (uc *cartUseCase) GetGuestCart(ctx context.Context, sessionID string) (*CartResponse, error) {
	cart, err := uc.cartRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		// Create new guest cart if not exists
		cart = &entities.Cart{
			ID:        uuid.New(),
			SessionID: &sessionID,
			Items:     []entities.CartItem{},
			Status:    "active",
			Currency:  "USD",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Set expiration and update calculated fields
		cart.SetExpiration()
		cart.UpdateCalculatedFields()

		if err := uc.cartRepo.Create(ctx, cart); err != nil {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to create guest cart")
		}
	}

	return uc.toCartResponse(cart), nil
}

// AddToGuestCart adds item to guest cart
func (uc *cartUseCase) AddToGuestCart(ctx context.Context, sessionID string, req AddToCartRequest) (*CartResponse, error) {
	return uc.addToGuestCartInTransaction(ctx, sessionID, req)
}

// AddToCart adds item to cart
func (uc *cartUseCase) AddToCart(ctx context.Context, userID uuid.UUID, req AddToCartRequest) (*CartResponse, error) {
	// Validate input
	if req.Quantity <= 0 {
		return nil, pkgErrors.InvalidInput("Quantity must be greater than 0")
	}

	if req.Quantity > 100 { // Max quantity per item
		return nil, pkgErrors.InvalidInput("Quantity cannot exceed 100")
	}

	// Execute add to cart
	return uc.addToCartInTransaction(ctx, userID, req)
}

// addToCartInTransaction handles adding item to cart
func (uc *cartUseCase) addToCartInTransaction(ctx context.Context, userID uuid.UUID, req AddToCartRequest) (*CartResponse, error) {
	// Validate product exists and is available
	product, err := uc.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeProductNotFound, "Product not found")
	}

	// Check if product is active and available
	if product.Status != "active" {
		return nil, pkgErrors.InvalidInput("Product is not available for purchase")
	}

	// Use current product price (will be used when adding/updating cart items)
	_ = product.Price // Suppress unused variable warning

	// Get or create cart
	cart, err := uc.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Only create new cart if cart not found, not for other errors
		if err != entities.ErrCartNotFound {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to get user cart")
		}

		// Create new cart only if not found
		cart = &entities.Cart{
			ID:        uuid.New(),
			UserID:    &userID,
			Items:     []entities.CartItem{},
			Status:    "active",
			Currency:  "USD",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Set expiration and update calculated fields
		cart.SetExpiration()
		cart.UpdateCalculatedFields()

		if err := uc.cartRepo.Create(ctx, cart); err != nil {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to create cart")
		}
	} else {
		// Check if cart is expired
		if cart.IsExpired() {
			// Mark cart as expired and create new one
			cart.MarkAsAbandoned()
			if err := uc.cartRepo.Update(ctx, cart); err != nil {
				// Log error but continue with new cart creation
				fmt.Printf("Warning: Failed to mark expired cart as abandoned: %v\n", err)
			}

			// Create new cart
			cart = &entities.Cart{
				ID:        uuid.New(),
				UserID:    &userID,
				Items:     []entities.CartItem{},
				Status:    "active",
				Currency:  "USD",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			cart.SetExpiration()
			cart.UpdateCalculatedFields()

			if err := uc.cartRepo.Create(ctx, cart); err != nil {
				return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to create new cart after expiration")
			}
		}
	}

	// Check if cart is expired
	if cart.IsExpired() {
		cart.MarkAsAbandoned()
		// Save abandoned cart
		if err := uc.cartRepo.Update(ctx, cart); err != nil {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to mark cart as abandoned")
		}

		// Create new cart
		cart = &entities.Cart{
			ID:        uuid.New(),
			UserID:    &userID,
			Items:     []entities.CartItem{},
			Status:    "active",
			Currency:  "USD",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		cart.SetExpiration()
		cart.UpdateCalculatedFields()

		if err := uc.cartRepo.Create(ctx, cart); err != nil {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to create new cart")
		}
	}

	// Get product with current price
	product, err = uc.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, pkgErrors.ProductNotFound().WithContext("product_id", req.ProductID)
	}

	// Check if product is available
	if !product.IsAvailable() {
		return nil, pkgErrors.New(pkgErrors.ErrCodeProductNotAvailable, "Product is not available").
			WithContext("product_id", req.ProductID).
			WithContext("product_name", product.Name)
	}

	// Check if item already exists in cart
	existingItem := cart.GetItem(req.ProductID)

	// Check stock availability (no reservation needed for cart)
	var totalQuantity int
	if existingItem != nil {
		totalQuantity = existingItem.Quantity + req.Quantity
	} else {
		totalQuantity = req.Quantity
	}

	// Create a temporary cart item to check stock availability
	tempCartItem := entities.CartItem{
		ProductID: req.ProductID,
		Product:   *product,
		Quantity:  totalQuantity,
	}

	if err := uc.simpleStockService.CheckStockAvailability(ctx, []entities.CartItem{tempCartItem}); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInsufficientStock, "Stock not available").
			WithContext("product_id", req.ProductID).
			WithContext("product_name", product.Name).
			WithContext("requested_quantity", totalQuantity)
	}

	// Now update cart item (stock already checked)
	if existingItem != nil {
		// Check total quantity after adding new quantity
		if totalQuantity > 100 { // Max quantity per item
			return nil, pkgErrors.InvalidInput(fmt.Sprintf("Total quantity %d exceeds maximum allowed (100)", totalQuantity))
		}

		// Update existing item with current price and new quantity
		existingItem.Quantity = totalQuantity
		existingItem.Price = product.Price // Update to current price
		existingItem.CalculateTotal()      // Recalculate total
		existingItem.UpdatedAt = time.Now()

		if err := uc.cartRepo.UpdateItem(ctx, existingItem); err != nil {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to update cart item")
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
		// Calculate and set total
		cartItem.CalculateTotal()

		if err := uc.cartRepo.AddItem(ctx, cart.ID, cartItem); err != nil {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to add item to cart")
		}
	}

	// Get updated cart
	updatedCart, err := uc.cartRepo.GetByID(ctx, cart.ID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeCartNotFound, "Failed to get updated cart")
	}

	return uc.toCartResponse(updatedCart), nil
}

// addToGuestCartInTransaction handles adding item to guest cart
func (uc *cartUseCase) addToGuestCartInTransaction(ctx context.Context, sessionID string, req AddToCartRequest) (*CartResponse, error) {
	// Similar implementation to addToCartInTransaction but for guest carts
	// Get or create guest cart
	cart, err := uc.cartRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		// Only create new cart if cart not found, not for other errors
		if err != entities.ErrCartNotFound {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to get guest cart")
		}

		// Create new cart only if not found
		cart = &entities.Cart{
			ID:        uuid.New(),
			SessionID: &sessionID,
			Items:     []entities.CartItem{},
			Status:    "active",
			Currency:  "USD",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		cart.SetExpiration()
		cart.UpdateCalculatedFields()

		if err := uc.cartRepo.Create(ctx, cart); err != nil {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to create guest cart")
		}
	}

	// Rest of implementation similar to user cart...
	// Get product and validate
	product, err := uc.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, pkgErrors.ProductNotFound().WithContext("product_id", req.ProductID)
	}

	// Validate product is active and available
	if product.Status != entities.ProductStatusActive {
		return nil, pkgErrors.InvalidInput("Product is not available for purchase")
	}

	// Check stock availability
	if product.Stock < req.Quantity {
		return nil, pkgErrors.InsufficientStock().
			WithContext("product_id", req.ProductID).
			WithContext("available_stock", product.Stock).
			WithContext("requested_quantity", req.Quantity)
	}

	// Check stock availability using simple stock service
	tempCartItem := entities.CartItem{
		ProductID: req.ProductID,
		Product:   *product,
		Quantity:  req.Quantity,
	}
	if err := uc.simpleStockService.CheckStockAvailability(ctx, []entities.CartItem{tempCartItem}); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInsufficientStock, "Stock not available").
			WithContext("product_id", req.ProductID).
			WithContext("requested_quantity", req.Quantity)
	}

	if !product.IsAvailable() {
		return nil, pkgErrors.New(pkgErrors.ErrCodeProductNotAvailable, "Product is not available")
	}

	// Handle existing item or add new item (same logic as user cart)
	existingItem := cart.GetItem(req.ProductID)

	// Calculate total quantity and check stock availability
	var totalQuantity int
	if existingItem != nil {
		totalQuantity = existingItem.Quantity + req.Quantity
		if totalQuantity > 100 { // Max quantity per item
			return nil, pkgErrors.InvalidInput(fmt.Sprintf("Total quantity %d exceeds maximum allowed (100)", totalQuantity))
		}
	} else {
		totalQuantity = req.Quantity
	}

	// Check stock availability for guest cart (no reservation needed)
	guestCartItem := entities.CartItem{
		ProductID: req.ProductID,
		Product:   *product,
		Quantity:  totalQuantity,
	}
	if err := uc.simpleStockService.CheckStockAvailability(ctx, []entities.CartItem{guestCartItem}); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInsufficientStock, "Stock not available").
			WithContext("product_id", req.ProductID).
			WithContext("product_name", product.Name).
			WithContext("requested_quantity", totalQuantity)
	}

	if existingItem != nil {
		existingItem.Quantity += req.Quantity
		existingItem.Price = product.Price
		existingItem.CalculateTotal() // Recalculate total
		existingItem.UpdatedAt = time.Now()
		if err := uc.cartRepo.UpdateItem(ctx, existingItem); err != nil {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to update guest cart item")
		}
	} else {
		cartItem := &entities.CartItem{
			ID:        uuid.New(),
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			Price:     product.Price,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		// Calculate and set total
		cartItem.CalculateTotal()
		if err := uc.cartRepo.AddItem(ctx, cart.ID, cartItem); err != nil {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to add item to guest cart")
		}
	}

	// Stock reservation was already created atomically above

	updatedCart, err := uc.cartRepo.GetByID(ctx, cart.ID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeCartNotFound, "Failed to get updated guest cart")
	}

	return uc.toCartResponse(updatedCart), nil
}

// UpdateCartItem updates cart item quantity
func (uc *cartUseCase) UpdateCartItem(ctx context.Context, userID uuid.UUID, req UpdateCartItemRequest) (*CartResponse, error) {
	// Validate input
	if req.Quantity <= 0 {
		return nil, pkgErrors.InvalidInput("Quantity must be greater than 0")
	}
	if req.Quantity > 100 { // Max quantity per item
		return nil, pkgErrors.InvalidInput("Quantity cannot exceed 100")
	}

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

	// Check if product is available
	if !product.IsAvailable() {
		return nil, pkgErrors.New(pkgErrors.ErrCodeProductNotAvailable, "Product is not available").
			WithContext("product_id", req.ProductID).
			WithContext("product_name", product.Name)
	}

	// Get cart item
	cartItem, err := uc.cartRepo.GetItem(ctx, cart.ID, req.ProductID)
	if err != nil {
		return nil, entities.ErrCartItemNotFound
	}

	// Check stock availability for the new quantity
	if req.Quantity > cartItem.Quantity {
		// Check if we have enough stock for the new total quantity
		tempCartItem := entities.CartItem{
			ProductID: req.ProductID,
			Product:   *product,
			Quantity:  req.Quantity,
		}
		if err := uc.simpleStockService.CheckStockAvailability(ctx, []entities.CartItem{tempCartItem}); err != nil {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInsufficientStock, "Stock not available").
				WithContext("product_id", req.ProductID).
				WithContext("product_name", product.Name).
				WithContext("requested_quantity", req.Quantity)
		}
	}

	// Update quantity and price
	cartItem.Quantity = req.Quantity
	cartItem.Price = product.Price // Update to current price
	cartItem.CalculateTotal()      // Recalculate total
	cartItem.UpdatedAt = time.Now()

	if err := uc.cartRepo.UpdateItem(ctx, cartItem); err != nil {
		return nil, err
	}

	// No stock reservation needed - stock will be reduced when order is placed

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
		ID:             cart.ID,
		UserID:         cart.UserID, // Now properly nullable
		SessionID:      cart.SessionID,
		ItemCount:      cart.ItemCount,
		Subtotal:       cart.Subtotal,
		TaxAmount:      cart.TaxAmount,      // Added missing field
		ShippingAmount: cart.ShippingAmount, // Added missing field
		Total:          cart.Total,
		Status:         cart.Status,
		Currency:       cart.Currency,
		Notes:          cart.Notes,
		ExpiresAt:      cart.ExpiresAt,
		IsGuest:        cart.IsGuest(), // Added helper field
		CreatedAt:      cart.CreatedAt,
		UpdatedAt:      cart.UpdatedAt,
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
			response.Items[i].Product = uc.toProductResponse(&item.Product)
		}
	}

	return response
}

// toProductResponse converts product entity to product response
// This is a simplified version for cart use case, a more complete version
// might exist in product use case.
func (uc *cartUseCase) toProductResponse(product *entities.Product) *ProductResponse {
	if product == nil {
		return nil
	}

	// Category information is now handled via ProductCategory many-to-many
	// For cart items, we don't need category info, so set to nil
	var categoryResponse *ProductCategoryResponse = nil

	var imageResponses []ProductImageResponse
	for _, img := range product.Images {
		imageResponses = append(imageResponses, ProductImageResponse{
			ID:       img.ID,
			URL:      img.URL,
			AltText:  img.AltText,
			Position: img.Position,
		})
	}

	return &ProductResponse{
		ID:                     product.ID,
		Name:                   product.Name,
		Description:            product.Description,
		ShortDescription:       product.ShortDescription,
		SKU:                    product.SKU,
		Slug:                   product.Slug,
		MetaTitle:              product.MetaTitle,
		MetaDescription:        product.MetaDescription,
		Keywords:               product.Keywords,
		Featured:               product.Featured,
		Visibility:             product.Visibility,
		Price:                  product.Price,
		ComparePrice:           product.ComparePrice,
		CostPrice:              product.CostPrice,
		SalePrice:              product.SalePrice,
		SaleStartDate:          product.SaleStartDate,
		SaleEndDate:            product.SaleEndDate,
		CurrentPrice:           product.GetCurrentPrice(),
		OriginalPrice:          product.GetOriginalPrice(),
		IsOnSale:               product.IsOnSale(),
		HasDiscount:            product.HasDiscount(),
		SaleDiscountPercentage: product.GetSaleDiscountPercentage(),
		DiscountPercentage:     product.GetDiscountPercentage(),
		Stock:                  product.Stock,
		LowStockThreshold:      product.LowStockThreshold,
		TrackQuantity:          product.TrackQuantity,
		AllowBackorder:         product.AllowBackorder,
		StockStatus:            product.StockStatus,
		IsLowStock:             product.IsLowStock(),
		Weight:                 product.Weight,
		Dimensions:             toDimensionsResponse(product.Dimensions),
		RequiresShipping:       product.RequiresShipping,
		ShippingClass:          product.ShippingClass,
		TaxClass:               product.TaxClass,
		CountryOfOrigin:        product.CountryOfOrigin,
		Category:               categoryResponse,
		// Brand:             toProductBrandResponse(product.Brand), // Assuming Brand conversion is needed elsewhere
		Images: imageResponses,
		// Tags:       toProductTagResponses(product.Tags), // Assuming Tag conversion is needed elsewhere
		// Attributes: toProductAttributeResponses(product.Attributes), // Assuming Attribute conversion is needed elsewhere
		// Variants:   toProductVariantResponses(product.Variants), // Assuming Variant conversion is needed elsewhere
		Status:      product.Status,
		ProductType: product.ProductType,
		IsDigital:   product.IsDigital,
		IsAvailable: product.IsAvailable(),
		HasVariants: product.HasVariants(),
		MainImage:   product.GetMainImage(),
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

// toDimensionsResponse converts Dimensions entity to DimensionsResponse
func toDimensionsResponse(d *entities.Dimensions) *DimensionsResponse {
	if d == nil {
		return nil
	}
	return &DimensionsResponse{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
	}
}

// MergeGuestCart merges guest cart with user cart when user logs in (uses auto strategy)
func (uc *cartUseCase) MergeGuestCart(ctx context.Context, userID uuid.UUID, sessionID string) (*CartResponse, error) {
	return uc.MergeGuestCartWithStrategy(ctx, userID, sessionID, MergeStrategyAuto)
}

// MergeGuestCartWithStrategy merges guest cart with user cart using specified strategy
func (uc *cartUseCase) MergeGuestCartWithStrategy(ctx context.Context, userID uuid.UUID, sessionID string, strategy MergeStrategy) (*CartResponse, error) {
	// Use transaction to prevent race conditions
	result, err := uc.cartRepo.WithTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
		// Get transaction repository from context
		txRepo, ok := txCtx.Value("tx_repo").(repositories.CartRepository)
		if !ok {
			txRepo = uc.cartRepo // fallback to original repo
		}

		// Get guest cart with row-level locking
		guestCart, err := txRepo.GetBySessionIDForUpdate(txCtx, sessionID)
		if err != nil {
			// No guest cart to merge, just return user cart
			return uc.getCartWithRepo(txCtx, txRepo, userID)
		}

		// Get user cart with row-level locking
		userCart, err := txRepo.GetByUserIDForUpdate(txCtx, userID)
		if err != nil {
			// No user cart exists, convert guest cart to user cart
			guestCart.UserID = &userID
			guestCart.SessionID = nil
			guestCart.UpdateCalculatedFields()

			// Also need to update any stock reservations to point to user instead of session
			if err := uc.transferGuestReservationsToUser(txCtx, sessionID, userID); err != nil {
				return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to transfer guest reservations to user")
			}

			if err := txRepo.Update(txCtx, guestCart); err != nil {
				return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to convert guest cart to user cart")
			}

			return uc.toCartResponse(guestCart), nil
		}

		// User cart exists, apply merge strategy
		switch strategy {
		case MergeStrategyKeepUser:
			// Keep user cart, mark guest cart as abandoned
			guestCart.MarkAsAbandoned()
			if err := txRepo.Update(txCtx, guestCart); err != nil {
				return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to mark guest cart as abandoned")
			}
			return uc.toCartResponse(userCart), nil

		case MergeStrategyReplace:
			// Replace user cart with guest cart
			// Clear user cart items first
			if err := txRepo.ClearCart(txCtx, userCart.ID); err != nil {
				return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to clear user cart")
			}

			// Move guest cart items to user cart
			for _, guestItem := range guestCart.Items {
				newItem := &entities.CartItem{
					ID:        uuid.New(),
					CartID:    userCart.ID,
					ProductID: guestItem.ProductID,
					Quantity:  guestItem.Quantity,
					Price:     guestItem.Price,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				// Calculate and set total
				newItem.CalculateTotal()

				if err := txRepo.AddItem(txCtx, userCart.ID, newItem); err != nil {
					return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to replace cart item")
				}
			}

			// Mark guest cart as converted
			guestCart.MarkAsConverted()
			if err := txRepo.Update(txCtx, guestCart); err != nil {
				return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to mark guest cart as converted")
			}

		case MergeStrategyMerge, MergeStrategyAuto:
			// Merge guest cart items into user cart (existing logic)
			return uc.mergeCartItemsWithRepo(txCtx, txRepo, userCart, guestCart)
		}

		// Get updated user cart
		updatedUserCart, err := txRepo.GetByID(txCtx, userCart.ID)
		if err != nil {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeCartNotFound, "Failed to get updated user cart")
		}

		return uc.toCartResponse(updatedUserCart), nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*CartResponse), nil
}

// mergeCartItems merges guest cart items into user cart
func (uc *cartUseCase) mergeCartItems(ctx context.Context, userCart, guestCart *entities.Cart) (*CartResponse, error) {
	// Merge guest cart items into user cart
	for _, guestItem := range guestCart.Items {
		// Validate product availability
		product, err := uc.productRepo.GetByID(ctx, guestItem.ProductID)
		if err != nil {
			// Skip unavailable products but log warning
			continue
		}
		if !product.IsAvailable() {
			// Skip unavailable products
			continue
		}

		existingItem := userCart.GetItem(guestItem.ProductID)
		if existingItem != nil {
			// Check if merged quantity exceeds maximum
			newQuantity := existingItem.Quantity + guestItem.Quantity
			if newQuantity > 100 {
				// Cap at maximum allowed quantity
				newQuantity = 100
			}

			// Update quantity and use current product price
			existingItem.Quantity = newQuantity
			existingItem.Price = product.Price // Use current price
			existingItem.CalculateTotal()
			existingItem.UpdatedAt = time.Now()

			if err := uc.cartRepo.UpdateItem(ctx, existingItem); err != nil {
				return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to merge cart item")
			}
		} else {
			// Add new item to user cart
			newItem := &entities.CartItem{
				ID:        uuid.New(),
				CartID:    userCart.ID,
				ProductID: guestItem.ProductID,
				Quantity:  guestItem.Quantity,
				Price:     product.Price, // Use current price
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			// Calculate and set total
			newItem.CalculateTotal()

			if err := uc.cartRepo.AddItem(ctx, userCart.ID, newItem); err != nil {
				return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to add merged cart item")
			}
		}
	}

	// Mark guest cart as converted
	guestCart.MarkAsConverted()
	if err := uc.cartRepo.Update(ctx, guestCart); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to mark guest cart as converted")
	}

	// Get updated user cart
	updatedUserCart, err := uc.cartRepo.GetByID(ctx, userCart.ID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeCartNotFound, "Failed to get updated user cart")
	}

	return uc.toCartResponse(updatedUserCart), nil
}

// getCartWithRepo gets cart using specific repository (for transaction support)
func (uc *cartUseCase) getCartWithRepo(ctx context.Context, repo repositories.CartRepository, userID uuid.UUID) (*CartResponse, error) {
	cart, err := repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeCartNotFound, "Cart not found")
	}
	return uc.toCartResponse(cart), nil
}

// mergeCartItemsWithRepo merges guest cart items into user cart using specific repository
func (uc *cartUseCase) mergeCartItemsWithRepo(ctx context.Context, repo repositories.CartRepository, userCart, guestCart *entities.Cart) (*CartResponse, error) {
	// Merge guest cart items into user cart with proper validation
	for _, guestItem := range guestCart.Items {
		// Get current product to ensure we use current price and validate availability
		product, err := uc.productRepo.GetByID(ctx, guestItem.ProductID)
		if err != nil {
			// Skip unavailable products but log warning
			fmt.Printf("Warning: Skipping product %s during merge: %v\n", guestItem.ProductID, err)
			continue
		}
		if !product.IsAvailable() {
			// Skip unavailable products
			fmt.Printf("Warning: Skipping unavailable product %s during merge\n", product.Name)
			continue
		}

		// Check stock availability for the quantity being merged
		if err := uc.simpleStockService.CheckStockAvailability(ctx, []entities.CartItem{guestItem}); err != nil {
			// Skip items with insufficient stock but log warning
			fmt.Printf("Warning: Skipping product %s due to insufficient stock: %v\n", product.Name, err)
			continue
		}

		existingItem := userCart.GetItem(guestItem.ProductID)
		if existingItem != nil {
			// Check if merged quantity exceeds maximum
			newQuantity := existingItem.Quantity + guestItem.Quantity
			if newQuantity > 100 {
				// Cap at maximum allowed quantity
				newQuantity = 100
			}

			// Update quantity and use current product price (consistent with other merge method)
			existingItem.Quantity = newQuantity
			existingItem.Price = product.Price // Use current price instead of old guest price
			existingItem.CalculateTotal()
			existingItem.UpdatedAt = time.Now()

			if err := repo.UpdateItem(ctx, existingItem); err != nil {
				return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to merge cart item")
			}
		} else {
			// Add new item to user cart
			newItem := &entities.CartItem{
				ID:        uuid.New(),
				CartID:    userCart.ID,
				ProductID: guestItem.ProductID,
				Quantity:  guestItem.Quantity,
				Price:     product.Price, // Use current price instead of old guest price
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			// Calculate and set total
			newItem.CalculateTotal()

			if err := repo.AddItem(ctx, userCart.ID, newItem); err != nil {
				return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to add merged cart item")
			}
		}
	}

	// Mark guest cart as converted
	guestCart.MarkAsConverted()
	if err := repo.Update(ctx, guestCart); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to mark guest cart as converted")
	}

	// Get updated user cart
	updatedUserCart, err := repo.GetByID(ctx, userCart.ID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeCartNotFound, "Failed to get updated user cart")
	}

	return uc.toCartResponse(updatedUserCart), nil
}

// CheckMergeConflict checks for conflicts when merging guest cart with user cart
func (uc *cartUseCase) CheckMergeConflict(ctx context.Context, userID uuid.UUID, sessionID string) (*CartConflictInfo, error) {
	conflict := &CartConflictInfo{
		HasConflict:     false,
		UserCartExists:  false,
		GuestCartExists: false,
		Recommendations: []string{},
	}

	// Check if guest cart exists
	guestCart, err := uc.cartRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		conflict.GuestCartExists = false
		conflict.Recommendations = append(conflict.Recommendations, "No guest cart found - nothing to merge")
		return conflict, nil
	}

	conflict.GuestCartExists = true
	conflict.GuestCart = uc.toCartResponse(guestCart)

	// Check if user cart exists
	userCart, err := uc.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		conflict.UserCartExists = false
		conflict.Recommendations = append(conflict.Recommendations, "No user cart exists - guest cart will become user cart")
		return conflict, nil
	}

	conflict.UserCartExists = true
	conflict.UserCart = uc.toCartResponse(userCart)

	// Check for conflicting items
	conflictingItems := []ConflictingItem{}
	for _, guestItem := range guestCart.Items {
		if userItem := userCart.GetItem(guestItem.ProductID); userItem != nil {
			// Get product name
			product, err := uc.productRepo.GetByID(ctx, guestItem.ProductID)
			productName := "Unknown Product"
			if err == nil && product != nil {
				productName = product.Name
			}

			conflictItem := ConflictingItem{
				ProductID:       guestItem.ProductID.String(),
				ProductName:     productName,
				UserQuantity:    userItem.Quantity,
				GuestQuantity:   guestItem.Quantity,
				UserPrice:       userItem.Price,
				GuestPrice:      guestItem.Price,
				PriceDifference: guestItem.Price - userItem.Price,
			}
			conflictingItems = append(conflictingItems, conflictItem)
		}
	}

	if len(conflictingItems) > 0 {
		conflict.HasConflict = true
		conflict.ConflictingItems = conflictingItems

		// Generate recommendations
		conflict.Recommendations = append(conflict.Recommendations,
			"Items exist in both carts - choose merge strategy:")
		conflict.Recommendations = append(conflict.Recommendations,
			"• 'merge' - Add quantities together")
		conflict.Recommendations = append(conflict.Recommendations,
			"• 'replace' - Replace user cart with guest cart")
		conflict.Recommendations = append(conflict.Recommendations,
			"• 'keep_user' - Keep user cart, discard guest cart")
	} else {
		conflict.Recommendations = append(conflict.Recommendations,
			"No conflicts found - guest cart items will be added to user cart")
	}

	return conflict, nil
}

// transferGuestReservationsToUser transfers stock reservations from guest session to user
func (uc *cartUseCase) transferGuestReservationsToUser(ctx context.Context, sessionID string, userID uuid.UUID) error {
	// For now, we'll implement a simple approach:
	// 1. Get all reservations for the session (this would need to be implemented in repository)
	// 2. Update them to point to the user instead of session
	// This is a simplified implementation - in production, you'd want proper session-based reservation tracking
	fmt.Printf("INFO: Transferring stock reservations from session %s to user %s (simplified implementation)\n", sessionID, userID)
	return nil
}

// CleanupExpiredCarts cleans up expired carts and their stock reservations
func (uc *cartUseCase) CleanupExpiredCarts(ctx context.Context) error {
	// Get all expired carts
	expiredCarts, err := uc.cartRepo.GetExpiredCarts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expired carts: %w", err)
	}

	for _, cart := range expiredCarts {
		// Mark cart as abandoned
		cart.MarkAsAbandoned()
		if err := uc.cartRepo.Update(ctx, cart); err != nil {
			fmt.Printf("Warning: Failed to mark cart %s as abandoned: %v\n", cart.ID, err)
			continue
		}

		// No stock reservations to release - using simple stock service
	}

	return nil
}

// CleanupExpiredStockReservations cleans up expired stock reservations - DEPRECATED
func (uc *cartUseCase) CleanupExpiredStockReservations(ctx context.Context) error {
	// No longer needed with simple stock service
	return nil
}
