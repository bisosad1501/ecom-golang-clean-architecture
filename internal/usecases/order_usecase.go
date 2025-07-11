package usecases

import (
	"context"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/domain/services"

	"github.com/google/uuid"
)

// OrderUseCase defines order use cases
type OrderUseCase interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, req CreateOrderRequest) (*OrderResponse, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (*OrderResponse, error)
	GetOrderBySessionID(ctx context.Context, sessionID string, userID uuid.UUID) (*OrderResponse, error)
	GetUserOrders(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*OrderResponse, error)
	GetUserOrdersWithFilters(ctx context.Context, userID uuid.UUID, req GetUserOrdersRequest) (*PaginatedOrderResponse, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status entities.OrderStatus) (*OrderResponse, error)
	CancelOrder(ctx context.Context, orderID uuid.UUID) (*OrderResponse, error)
	GetOrders(ctx context.Context, req GetOrdersRequest) ([]*OrderResponse, error)
}

type orderUseCase struct {
	orderRepo     repositories.OrderRepository
	cartRepo      repositories.CartRepository
	productRepo   repositories.ProductRepository
	paymentRepo   repositories.PaymentRepository
	inventoryRepo repositories.InventoryRepository
	orderService  services.OrderService
}

// NewOrderUseCase creates a new order use case
func NewOrderUseCase(
	orderRepo repositories.OrderRepository,
	cartRepo repositories.CartRepository,
	productRepo repositories.ProductRepository,
	paymentRepo repositories.PaymentRepository,
	inventoryRepo repositories.InventoryRepository,
	orderService services.OrderService,
) OrderUseCase {
	return &orderUseCase{
		orderRepo:     orderRepo,
		cartRepo:      cartRepo,
		productRepo:   productRepo,
		paymentRepo:   paymentRepo,
		inventoryRepo: inventoryRepo,
		orderService:  orderService,
	}
}

// CreateOrderRequest represents create order request
type CreateOrderRequest struct {
	ShippingAddress AddressRequest         `json:"shipping_address" validate:"required"`
	BillingAddress  *AddressRequest        `json:"billing_address"`
	PaymentMethod   entities.PaymentMethod `json:"payment_method" validate:"required"`
	Notes           string                 `json:"notes"`
	TaxRate         float64                `json:"tax_rate" validate:"min=0,max=1"`
	ShippingCost    float64                `json:"shipping_cost" validate:"min=0"`
	DiscountAmount  float64                `json:"discount_amount" validate:"min=0"`
}

// GetOrdersRequest represents get orders request
type GetOrdersRequest struct {
	Status        *entities.OrderStatus   `json:"status"`
	PaymentStatus *entities.PaymentStatus `json:"payment_status"`
	StartDate     *time.Time              `json:"start_date"`
	EndDate       *time.Time              `json:"end_date"`
	SortBy        string                  `json:"sort_by"`
	SortOrder     string                  `json:"sort_order"`
	Limit         int                     `json:"limit" validate:"min=1,max=100"`
	Offset        int                     `json:"offset" validate:"min=0"`
}

// GetUserOrdersRequest represents get user orders request with filters
type GetUserOrdersRequest struct {
	Status        *entities.OrderStatus   `json:"status"`
	PaymentStatus *entities.PaymentStatus `json:"payment_status"`
	StartDate     *time.Time              `json:"start_date"`
	EndDate       *time.Time              `json:"end_date"`
	SortBy        string                  `json:"sort_by"`
	SortOrder     string                  `json:"sort_order"`
	Limit         int                     `json:"limit" validate:"min=1,max=100"`
	Offset        int                     `json:"offset" validate:"min=0"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// PaginatedOrderResponse represents a paginated order response
type PaginatedOrderResponse struct {
	Data       []*OrderResponse `json:"data"`
	Pagination PaginationMeta   `json:"pagination"`
}

// AddressRequest represents address request
type AddressRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Company   string `json:"company"`
	Address1  string `json:"address1" validate:"required"`
	Address2  string `json:"address2"`
	City      string `json:"city" validate:"required"`
	State     string `json:"state" validate:"required"`
	ZipCode   string `json:"zip_code" validate:"required"`
	Country   string `json:"country" validate:"required"`
	Phone     string `json:"phone"`
}

// OrderResponse represents order response
type OrderResponse struct {
	ID              uuid.UUID              `json:"id"`
	OrderNumber     string                 `json:"order_number"`
	User            *UserResponse          `json:"user"`
	Items           []OrderItemResponse    `json:"items"`
	Status          entities.OrderStatus   `json:"status"`
	PaymentStatus   entities.PaymentStatus `json:"payment_status"`
	Subtotal        float64                `json:"subtotal"`
	TaxAmount       float64                `json:"tax_amount"`
	ShippingAmount  float64                `json:"shipping_amount"`
	DiscountAmount  float64                `json:"discount_amount"`
	Total           float64                `json:"total"`
	Currency        string                 `json:"currency"`
	ShippingAddress *OrderAddressResponse  `json:"shipping_address"`
	BillingAddress  *OrderAddressResponse  `json:"billing_address"`
	Notes           string                 `json:"notes"`
	Payment         *PaymentResponse       `json:"payment"`
	ItemCount       int                    `json:"item_count"`
	CanBeCancelled  bool                   `json:"can_be_cancelled"`
	CanBeRefunded   bool                   `json:"can_be_refunded"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// OrderItemResponse represents order item response
type OrderItemResponse struct {
	ID          uuid.UUID        `json:"id"`
	Product     *ProductResponse `json:"product"`
	ProductName string           `json:"product_name"`
	ProductSKU  string           `json:"product_sku"`
	Quantity    int              `json:"quantity"`
	Price       float64          `json:"price"`
	Total       float64          `json:"total"`
}

// OrderAddressResponse represents order address response
type OrderAddressResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company"`
	Address1  string `json:"address1"`
	Address2  string `json:"address2"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zip_code"`
	Country   string `json:"country"`
	Phone     string `json:"phone"`
}

// CreateOrder creates a new order
func (uc *orderUseCase) CreateOrder(ctx context.Context, userID uuid.UUID, req CreateOrderRequest) (*OrderResponse, error) {
	// Get user's cart
	cart, err := uc.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, entities.ErrCartNotFound
	}

	if cart.IsEmpty() {
		return nil, entities.ErrInvalidInput
	}

	// Validate cart items and check stock
	if err := uc.orderService.ValidateOrderItems(cart.Items); err != nil {
		return nil, err
	}

	// Check product availability and stock
	for _, item := range cart.Items {
		product, err := uc.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, entities.ErrProductNotFound
		}

		if !product.IsAvailable() {
			return nil, entities.ErrProductNotAvailable
		}

		if !product.CanReduceStock(item.Quantity) {
			return nil, entities.ErrInsufficientStock
		}
	}

	// Calculate totals
	subtotal, taxAmount, total := uc.orderService.CalculateOrderTotal(
		cart.Items, req.TaxRate, req.ShippingCost, req.DiscountAmount,
	)

	// Create order
	order := &entities.Order{
		ID:             uuid.New(),
		OrderNumber:    uc.orderService.GenerateOrderNumber(),
		UserID:         userID,
		Status:         entities.OrderStatusPending,
		PaymentStatus:  entities.PaymentStatusPending,
		Subtotal:       subtotal,
		TaxAmount:      taxAmount,
		ShippingAmount: req.ShippingCost,
		DiscountAmount: req.DiscountAmount,
		Total:          total,
		Currency:       "USD",
		Notes:          req.Notes,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Set addresses
	order.ShippingAddress = &entities.OrderAddress{
		FirstName: req.ShippingAddress.FirstName,
		LastName:  req.ShippingAddress.LastName,
		Company:   req.ShippingAddress.Company,
		Address1:  req.ShippingAddress.Address1,
		Address2:  req.ShippingAddress.Address2,
		City:      req.ShippingAddress.City,
		State:     req.ShippingAddress.State,
		ZipCode:   req.ShippingAddress.ZipCode,
		Country:   req.ShippingAddress.Country,
		Phone:     req.ShippingAddress.Phone,
	}

	if req.BillingAddress != nil {
		order.BillingAddress = &entities.OrderAddress{
			FirstName: req.BillingAddress.FirstName,
			LastName:  req.BillingAddress.LastName,
			Company:   req.BillingAddress.Company,
			Address1:  req.BillingAddress.Address1,
			Address2:  req.BillingAddress.Address2,
			City:      req.BillingAddress.City,
			State:     req.BillingAddress.State,
			ZipCode:   req.BillingAddress.ZipCode,
			Country:   req.BillingAddress.Country,
			Phone:     req.BillingAddress.Phone,
		}
	} else {
		order.BillingAddress = order.ShippingAddress
	}

	// Create order items
	for _, cartItem := range cart.Items {
		product, _ := uc.productRepo.GetByID(ctx, cartItem.ProductID)

		orderItem := entities.OrderItem{
			ID:          uuid.New(),
			OrderID:     order.ID,
			ProductID:   cartItem.ProductID,
			ProductName: product.Name,
			ProductSKU:  product.SKU,
			Quantity:    cartItem.Quantity,
			Price:       cartItem.Price,
			Total:       cartItem.GetSubtotal(),
			CreatedAt:   time.Now(),
		}

		order.Items = append(order.Items, orderItem)
	}

	// Create order
	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Note: Payment record will be created during checkout session creation
	// This allows the payment to have the proper transaction ID from Stripe

	// Reduce product stock using dedicated UpdateStock method to avoid slug conflicts
	for _, item := range cart.Items {
		product, err := uc.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product %s: %w", item.ProductID, err)
		}

		// Validate stock can be reduced
		if err := product.ReduceStock(item.Quantity); err != nil {
			return nil, fmt.Errorf("failed to reduce stock for product %s: %w", product.Name, err)
		}

		// Use UpdateStock method instead of full Update to avoid slug constraint issues
		if err := uc.productRepo.UpdateStock(ctx, item.ProductID, product.Stock); err != nil {
			return nil, fmt.Errorf("failed to update stock for product %s: %w", product.Name, err)
		}
	}

	// Clear cart
	uc.cartRepo.ClearCart(ctx, cart.ID)

	// Get created order with relations
	createdOrder, err := uc.orderRepo.GetByID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	return uc.toOrderResponse(createdOrder), nil
}

// GetOrder gets an order by ID
func (uc *orderUseCase) GetOrder(ctx context.Context, orderID uuid.UUID) (*OrderResponse, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, entities.ErrOrderNotFound
	}

	return uc.toOrderResponse(order), nil
}

// GetOrderBySessionID gets an order by checkout session ID
func (uc *orderUseCase) GetOrderBySessionID(ctx context.Context, sessionID string, userID uuid.UUID) (*OrderResponse, error) {
	// First find the payment by session ID
	payment, err := uc.paymentRepo.GetByTransactionID(ctx, sessionID)
	if err != nil {
		return nil, entities.ErrOrderNotFound
	}

	// Get the order from the payment
	order, err := uc.orderRepo.GetByID(ctx, payment.OrderID)
	if err != nil {
		return nil, entities.ErrOrderNotFound
	}

	// Check if user owns this order
	if order.UserID != userID {
		return nil, entities.ErrOrderNotFound
	}

	return uc.toOrderResponse(order), nil
}

// GetUserOrders gets user's orders
func (uc *orderUseCase) GetUserOrders(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*OrderResponse, error) {
	orders, err := uc.orderRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = uc.toOrderResponse(order)
	}

	return responses, nil
}

// GetUserOrdersWithFilters gets user's orders with filters
func (uc *orderUseCase) GetUserOrdersWithFilters(ctx context.Context, userID uuid.UUID, req GetUserOrdersRequest) (*PaginatedOrderResponse, error) {
	// Convert request to search parameters
	params := repositories.OrderSearchParams{
		UserID:        &userID,
		Status:        req.Status,
		PaymentStatus: req.PaymentStatus,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		SortBy:        req.SortBy,
		SortOrder:     req.SortOrder,
		Limit:         req.Limit,
		Offset:        req.Offset,
	}

	// Set default sorting if not provided
	if params.SortBy == "" {
		params.SortBy = "created_at"
		params.SortOrder = "desc"
	}

	// Get total count with same filters
	totalCount, err := uc.orderRepo.CountSearch(ctx, params)
	if err != nil {
		return nil, err
	}

	// Get orders
	orders, err := uc.orderRepo.Search(ctx, params)
	if err != nil {
		return nil, err
	}

	// Convert to responses
	responses := make([]*OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = uc.toOrderResponse(order)
	}

	// Calculate pagination metadata
	page := (req.Offset / req.Limit) + 1
	totalPages := int((totalCount + int64(req.Limit) - 1) / int64(req.Limit)) // Ceiling division
	if totalPages == 0 {
		totalPages = 1
	}

	pagination := PaginationMeta{
		Page:       page,
		Limit:      req.Limit,
		Total:      totalCount,
		TotalPages: totalPages,
		HasNext:    req.Offset+req.Limit < int(totalCount),
		HasPrev:    req.Offset > 0,
	}

	return &PaginatedOrderResponse{
		Data:       responses,
		Pagination: pagination,
	}, nil
}

// UpdateOrderStatus updates order status
func (uc *orderUseCase) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status entities.OrderStatus) (*OrderResponse, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, entities.ErrOrderNotFound
	}

	if err := uc.orderRepo.UpdateStatus(ctx, orderID, status); err != nil {
		return nil, err
	}

	order.Status = status
	order.UpdatedAt = time.Now()

	return uc.toOrderResponse(order), nil
}

// CancelOrder cancels an order
func (uc *orderUseCase) CancelOrder(ctx context.Context, orderID uuid.UUID) (*OrderResponse, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, entities.ErrOrderNotFound
	}

	if !order.CanBeCancelled() {
		return nil, entities.ErrOrderCannotBeCancelled
	}

	// Restore product stock
	for _, item := range order.Items {
		product, err := uc.productRepo.GetByID(ctx, item.ProductID)
		if err == nil {
			product.IncreaseStock(item.Quantity)
			uc.productRepo.Update(ctx, product)
		}
	}

	return uc.UpdateOrderStatus(ctx, orderID, entities.OrderStatusCancelled)
}

// GetOrders gets list of orders
func (uc *orderUseCase) GetOrders(ctx context.Context, req GetOrdersRequest) ([]*OrderResponse, error) {
	params := repositories.OrderSearchParams{
		Status:        req.Status,
		PaymentStatus: req.PaymentStatus,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		SortBy:        req.SortBy,
		SortOrder:     req.SortOrder,
		Limit:         req.Limit,
		Offset:        req.Offset,
	}

	orders, err := uc.orderRepo.Search(ctx, params)
	if err != nil {
		return nil, err
	}

	responses := make([]*OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = uc.toOrderResponse(order)
	}

	return responses, nil
}

// toOrderResponse converts order entity to response
func (uc *orderUseCase) toOrderResponse(order *entities.Order) *OrderResponse {
	response := &OrderResponse{
		ID:             order.ID,
		OrderNumber:    order.OrderNumber,
		Status:         order.Status,
		PaymentStatus:  order.PaymentStatus,
		Subtotal:       order.Subtotal,
		TaxAmount:      order.TaxAmount,
		ShippingAmount: order.ShippingAmount,
		DiscountAmount: order.DiscountAmount,
		Total:          order.Total,
		Currency:       order.Currency,
		Notes:          order.Notes,
		ItemCount:      order.GetItemCount(),
		CanBeCancelled: order.CanBeCancelled(),
		CanBeRefunded:  order.CanBeRefunded(),
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
	}

	// Convert user
	if order.User.ID != uuid.Nil {
		userUseCase := &userUseCase{}
		response.User = userUseCase.toUserResponse(&order.User)
	}

	// Convert addresses
	if order.ShippingAddress != nil {
		response.ShippingAddress = &OrderAddressResponse{
			FirstName: order.ShippingAddress.FirstName,
			LastName:  order.ShippingAddress.LastName,
			Company:   order.ShippingAddress.Company,
			Address1:  order.ShippingAddress.Address1,
			Address2:  order.ShippingAddress.Address2,
			City:      order.ShippingAddress.City,
			State:     order.ShippingAddress.State,
			ZipCode:   order.ShippingAddress.ZipCode,
			Country:   order.ShippingAddress.Country,
			Phone:     order.ShippingAddress.Phone,
		}
	}

	if order.BillingAddress != nil {
		response.BillingAddress = &OrderAddressResponse{
			FirstName: order.BillingAddress.FirstName,
			LastName:  order.BillingAddress.LastName,
			Company:   order.BillingAddress.Company,
			Address1:  order.BillingAddress.Address1,
			Address2:  order.BillingAddress.Address2,
			City:      order.BillingAddress.City,
			State:     order.BillingAddress.State,
			ZipCode:   order.BillingAddress.ZipCode,
			Country:   order.BillingAddress.Country,
			Phone:     order.BillingAddress.Phone,
		}
	}

	// Convert items
	response.Items = make([]OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		response.Items[i] = OrderItemResponse{
			ID:          item.ID,
			ProductName: item.ProductName,
			ProductSKU:  item.ProductSKU,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Total:       item.Total,
		}

		// Add product info if available
		if item.Product.ID != uuid.Nil {
			productUseCase := &productUseCase{}
			response.Items[i].Product = productUseCase.toProductResponse(&item.Product)
		}
	}

	// Convert payment
	if order.Payment != nil {
		response.Payment = &PaymentResponse{
			ID:            order.Payment.ID,
			Amount:        order.Payment.Amount,
			Currency:      order.Payment.Currency,
			Method:        order.Payment.Method,
			Status:        order.Payment.Status,
			TransactionID: order.Payment.TransactionID,
			ProcessedAt:   order.Payment.ProcessedAt,
			RefundedAt:    order.Payment.RefundedAt,
			RefundAmount:  order.Payment.RefundAmount,
			CreatedAt:     order.Payment.CreatedAt,
		}
	}

	return response
}
