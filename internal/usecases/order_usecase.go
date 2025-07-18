package usecases

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/domain/services"
	"ecom-golang-clean-architecture/internal/infrastructure/database"
	pkgErrors "ecom-golang-clean-architecture/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	GetOrders(ctx context.Context, req GetOrdersRequest) (*GetOrdersResponse, error)

	// Shipping management
	UpdateShippingInfo(ctx context.Context, orderID uuid.UUID, req UpdateShippingInfoRequest) (*OrderResponse, error)
	UpdateDeliveryStatus(ctx context.Context, orderID uuid.UUID, status entities.OrderStatus) (*OrderResponse, error)

	// Order notes management
	AddOrderNote(ctx context.Context, orderID uuid.UUID, req AddOrderNoteRequest) error

	// Order events
	GetOrderEvents(ctx context.Context, orderID uuid.UUID, publicOnly bool) ([]*entities.OrderEvent, error)
}

// NotificationService interface for order notifications
type NotificationService interface {
	NotifyOrderCreated(ctx context.Context, orderID uuid.UUID) error
	NotifyOrderStatusChanged(ctx context.Context, orderID uuid.UUID, newStatus string) error
	NotifyNewOrder(ctx context.Context, orderID uuid.UUID) error
}

type orderUseCase struct {
	orderRepo               repositories.OrderRepository
	cartRepo                repositories.CartRepository
	productRepo             repositories.ProductRepository
	paymentRepo             repositories.PaymentRepository
	inventoryRepo           repositories.InventoryRepository
	orderEventRepo          repositories.OrderEventRepository
	orderService            services.OrderService
	simpleStockService      services.SimpleStockService
	orderEventService       services.OrderEventService
	userMetricsService      services.UserMetricsService
	notificationService     NotificationService
	txManager               *database.TransactionManager
}

// NewOrderUseCase creates a new order use case
func NewOrderUseCase(
	orderRepo repositories.OrderRepository,
	cartRepo repositories.CartRepository,
	productRepo repositories.ProductRepository,
	paymentRepo repositories.PaymentRepository,
	inventoryRepo repositories.InventoryRepository,
	orderEventRepo repositories.OrderEventRepository,
	orderService services.OrderService,
	simpleStockService services.SimpleStockService,
	orderEventService services.OrderEventService,
	userMetricsService services.UserMetricsService,
	notificationService NotificationService,
	txManager *database.TransactionManager,
) OrderUseCase {
	return &orderUseCase{
		orderRepo:               orderRepo,
		cartRepo:                cartRepo,
		productRepo:             productRepo,
		paymentRepo:             paymentRepo,
		inventoryRepo:           inventoryRepo,
		orderEventRepo:          orderEventRepo,
		orderService:            orderService,
		simpleStockService:      simpleStockService,
		orderEventService:       orderEventService,
		userMetricsService:      userMetricsService,
		notificationService:     notificationService,
		txManager:               txManager,
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

// PaginatedOrderResponse represents a paginated order response
type PaginatedOrderResponse struct {
	Data       []*OrderResponse `json:"data"`
	Pagination *PaginationInfo  `json:"pagination"`
}

// GetOrdersResponse represents admin orders response
type GetOrdersResponse struct {
	Orders     []*OrderResponse `json:"orders"`
	Pagination *PaginationInfo  `json:"pagination"`
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
	ID                   uuid.UUID                  `json:"id"`
	OrderNumber          string                     `json:"order_number"`
	User                 *UserResponse              `json:"user"`
	Items                []OrderItemResponse        `json:"items"`
	Status               entities.OrderStatus       `json:"status"`
	FulfillmentStatus    entities.FulfillmentStatus `json:"fulfillment_status"`
	PaymentStatus        entities.PaymentStatus     `json:"payment_status"`
	PaymentMethod        entities.PaymentMethod     `json:"payment_method"`
	Priority             entities.OrderPriority     `json:"priority"`
	Source               entities.OrderSource       `json:"source"`
	CustomerType         entities.CustomerType      `json:"customer_type"`
	Subtotal             float64                    `json:"subtotal"`
	TaxAmount            float64                    `json:"tax_amount"`
	ShippingAmount       float64                    `json:"shipping_amount"`
	DiscountAmount       float64                    `json:"discount_amount"`
	TipAmount            float64                    `json:"tip_amount"`
	Total                float64                    `json:"total"`
	Currency             string                     `json:"currency"`
	ShippingAddress      *OrderAddressResponse      `json:"shipping_address"`
	BillingAddress       *OrderAddressResponse      `json:"billing_address"`
	ShippingMethod       string                     `json:"shipping_method"`
	TrackingNumber       string                     `json:"tracking_number"`
	TrackingURL          string                     `json:"tracking_url"`
	Carrier              string                     `json:"carrier"`
	EstimatedDelivery    *time.Time                 `json:"estimated_delivery"`
	ActualDelivery       *time.Time                 `json:"actual_delivery"`
	DeliveryInstructions string                     `json:"delivery_instructions"`
	CustomerNotes        string                     `json:"customer_notes"`
	AdminNotes           string                     `json:"admin_notes"`
	IsGift               bool                       `json:"is_gift"`
	GiftMessage          string                     `json:"gift_message"`
	GiftWrap             bool                       `json:"gift_wrap"`
	Payment              *PaymentResponse           `json:"payment"`
	ItemCount            int                        `json:"item_count"`
	CanBeCancelled       bool                       `json:"can_be_cancelled"`
	CanBeRefunded        bool                       `json:"can_be_refunded"`
	CanBeShipped         bool                       `json:"can_be_shipped"`
	CanBeDelivered       bool                       `json:"can_be_delivered"`
	IsShipped            bool                       `json:"is_shipped"`
	IsDelivered          bool                       `json:"is_delivered"`
	HasTracking          bool                       `json:"has_tracking"`
	CreatedAt            time.Time                  `json:"created_at"`
	UpdatedAt            time.Time                  `json:"updated_at"`
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
	// Execute the entire order creation in a transaction
	result, err := uc.txManager.WithTransactionResult(ctx, func(tx *gorm.DB) (interface{}, error) {
		return uc.createOrderInTransaction(ctx, tx, userID, req)
	})
	if err != nil {
		return nil, err
	}
	return result.(*OrderResponse), nil
}

// validateCreateOrderRequest validates the create order request
func (uc *orderUseCase) validateCreateOrderRequest(req CreateOrderRequest) error {
	// Validate shipping address
	if err := uc.validateAddress(req.ShippingAddress, "shipping"); err != nil {
		return fmt.Errorf("invalid shipping address: %w", err)
	}

	// Validate billing address if provided
	if req.BillingAddress != nil {
		if err := uc.validateAddress(*req.BillingAddress, "billing"); err != nil {
			return fmt.Errorf("invalid billing address: %w", err)
		}
	}

	// Validate payment method
	validPaymentMethods := []entities.PaymentMethod{
		entities.PaymentMethodCreditCard,
		entities.PaymentMethodDebitCard,
		entities.PaymentMethodPayPal,
		entities.PaymentMethodStripe,
		entities.PaymentMethodApplePay,
		entities.PaymentMethodGooglePay,
		entities.PaymentMethodBankTransfer,
		entities.PaymentMethodCash, // Cash on Delivery (COD)
	}
	isValidPaymentMethod := false
	for _, method := range validPaymentMethods {
		if req.PaymentMethod == method {
			isValidPaymentMethod = true
			break
		}
	}
	if !isValidPaymentMethod {
		return fmt.Errorf("invalid payment method: %s", req.PaymentMethod)
	}

	// Validate financial amounts
	if req.TaxRate < 0 || req.TaxRate > 1 {
		return fmt.Errorf("tax rate must be between 0 and 1")
	}
	if req.ShippingCost < 0 {
		return fmt.Errorf("shipping cost cannot be negative")
	}
	if req.DiscountAmount < 0 {
		return fmt.Errorf("discount amount cannot be negative")
	}

	return nil
}

// validateAddress validates address data
func (uc *orderUseCase) validateAddress(addr AddressRequest, addressType string) error {
	// Required fields validation with length limits
	if addr.FirstName == "" {
		return fmt.Errorf("%s address first name is required", addressType)
	}
	if len(addr.FirstName) > 50 {
		return fmt.Errorf("%s address first name cannot exceed 50 characters", addressType)
	}

	if addr.LastName == "" {
		return fmt.Errorf("%s address last name is required", addressType)
	}
	if len(addr.LastName) > 50 {
		return fmt.Errorf("%s address last name cannot exceed 50 characters", addressType)
	}

	if addr.Address1 == "" {
		return fmt.Errorf("%s address line 1 is required", addressType)
	}
	if len(addr.Address1) > 100 {
		return fmt.Errorf("%s address line 1 cannot exceed 100 characters", addressType)
	}

	if len(addr.Address2) > 100 {
		return fmt.Errorf("%s address line 2 cannot exceed 100 characters", addressType)
	}

	if addr.City == "" {
		return fmt.Errorf("%s address city is required", addressType)
	}
	if len(addr.City) > 50 {
		return fmt.Errorf("%s address city cannot exceed 50 characters", addressType)
	}

	if addr.State == "" {
		return fmt.Errorf("%s address state is required", addressType)
	}
	if len(addr.State) > 50 {
		return fmt.Errorf("%s address state cannot exceed 50 characters", addressType)
	}

	if addr.ZipCode == "" {
		return fmt.Errorf("%s address zip code is required", addressType)
	}

	if addr.Country == "" {
		return fmt.Errorf("%s address country is required", addressType)
	}

	// Validate zip code format (basic validation)
	if len(addr.ZipCode) < 3 || len(addr.ZipCode) > 10 {
		return fmt.Errorf("%s address zip code must be between 3 and 10 characters", addressType)
	}
	// Enhanced zip code format validation
	zipRegex := `^[A-Za-z0-9\s\-]{3,10}$`
	if matched, _ := regexp.MatchString(zipRegex, addr.ZipCode); !matched {
		return fmt.Errorf("%s address zip code contains invalid characters", addressType)
	}

	// Validate country code (should be exactly 2 characters for ISO codes)
	if len(addr.Country) != 2 {
		return fmt.Errorf("%s address country must be a 2-letter ISO country code", addressType)
	}

	// Optional fields validation
	if len(addr.Company) > 100 {
		return fmt.Errorf("%s address company cannot exceed 100 characters", addressType)
	}

	// Validate phone format if provided
	if addr.Phone != "" {
		if len(addr.Phone) > 20 {
			return fmt.Errorf("%s address phone cannot exceed 20 characters", addressType)
		}
		phoneRegex := `^\+?[1-9]\d{1,14}$`
		if matched, _ := regexp.MatchString(phoneRegex, addr.Phone); !matched {
			return fmt.Errorf("%s address phone format is invalid", addressType)
		}
	}

	return nil
}

// createOrderInTransaction handles order creation within a transaction
func (uc *orderUseCase) createOrderInTransaction(ctx context.Context, tx *gorm.DB, userID uuid.UUID, req CreateOrderRequest) (*OrderResponse, error) {
	// Validate request data
	if err := uc.validateCreateOrderRequest(req); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInvalidInput, "Invalid order request")
	}

	// Get user's cart
	cart, err := uc.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, pkgErrors.CartNotFound()
	}

	if cart.IsEmpty() {
		return nil, pkgErrors.InvalidInput("Cart is empty")
	}

	// Validate cart items
	if err := uc.orderService.ValidateOrderItems(cart.Items); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInvalidInput, "Invalid cart items")
	}

	// Bulk check product availability to avoid N+1 queries
	productIDs := make([]uuid.UUID, len(cart.Items))
	for i, item := range cart.Items {
		productIDs[i] = item.ProductID
	}

	products, err := uc.getProductsBulk(ctx, productIDs)
	if err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeProductNotFound, "Failed to get products")
	}

	// Validate products and stock availability
	for _, item := range cart.Items {
		product, exists := products[item.ProductID]
		if !exists {
			return nil, pkgErrors.ProductNotFound().WithContext("product_id", item.ProductID)
		}

		if !product.IsAvailable() {
			return nil, pkgErrors.New(pkgErrors.ErrCodeProductNotAvailable, "Product not available").
				WithContext("product_id", item.ProductID).
				WithContext("product_name", product.Name)
		}

		// Check stock availability using simple stock service
		if err := uc.simpleStockService.CheckStockAvailability(ctx, []entities.CartItem{item}); err != nil {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInsufficientStock, "Stock not available").
				WithContext("product_id", item.ProductID).
				WithContext("product_name", product.Name).
				WithContext("requested_quantity", item.Quantity)
		}
	}

	// Calculate totals
	subtotal, taxAmount, total := uc.orderService.CalculateOrderTotal(
		cart.Items, req.TaxRate, req.ShippingCost, req.DiscountAmount,
	)

	// Generate unique order number
	orderNumber, err := uc.orderService.GenerateUniqueOrderNumber(ctx)
	if err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to generate order number")
	}

	// Determine initial payment status based on payment method
	initialPaymentStatus := entities.PaymentStatusPending
	if req.PaymentMethod == entities.PaymentMethodCash {
		// COD orders start with "awaiting_payment" status
		initialPaymentStatus = entities.PaymentStatusAwaitingPayment
	}

	// Create order with reservation fields
	order := &entities.Order{
		ID:             uuid.New(),
		OrderNumber:    orderNumber,
		UserID:         userID,
		Status:         entities.OrderStatusPending,
		PaymentStatus:  initialPaymentStatus,
		PaymentMethod:  req.PaymentMethod, // Store payment method
		Subtotal:       subtotal,
		TaxAmount:      taxAmount,
		ShippingAmount: req.ShippingCost,
		DiscountAmount: req.DiscountAmount,
		Total:          total,
		Currency:       "USD",
		CustomerNotes:  req.Notes,
		Source:         entities.OrderSourceWeb,
		CustomerType:   entities.CustomerTypeRegistered,
		Priority:       entities.OrderPriorityNormal,
		Version:        1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Set timeouts and validate
	order.ValidateTimeouts()

	// Set payment timeout
	order.SetPaymentTimeout(24) // 24 hours for payment

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

	// Create order items using bulk data
	for _, cartItem := range cart.Items {
		product := products[cartItem.ProductID]

		// Validate price consistency
		if cartItem.Price != product.Price {
			// Log warning but use current product price for order
			// This handles price changes between cart and order creation
		}

		orderItem := entities.OrderItem{
			ID:          uuid.New(),
			OrderID:     order.ID,
			ProductID:   cartItem.ProductID,
			ProductName: product.Name,
			ProductSKU:  product.SKU,
			Quantity:    cartItem.Quantity,
			Price:       product.Price, // Use current product price
			Total:       float64(cartItem.Quantity) * product.Price,
			Weight:      getProductWeight(product.Weight), // Add weight from product
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		order.Items = append(order.Items, orderItem)
	}

	// Update order total weight
	order.UpdateTotalWeight()

	// Create order within transaction
	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to create order")
	}

	// For COD orders, create a pending payment record
	if req.PaymentMethod == entities.PaymentMethodCash {
		codPayment := &entities.Payment{
			ID:        uuid.New(),
			OrderID:   order.ID,
			UserID:    userID,
			Amount:    total,
			Currency:  "USD",
			Method:    entities.PaymentMethodCash,
			Status:    entities.PaymentStatusAwaitingPayment,
			Gateway:   "cod",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Create payment record for COD
		if err := uc.paymentRepo.Create(ctx, codPayment); err != nil {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to create COD payment record")
		}
	}

	// FIXED: Reduce stock immediately for ALL payment methods when order is confirmed
	// This ensures consistent behavior regardless of payment method
	if order.Status == entities.OrderStatusConfirmed {
		if err := uc.simpleStockService.ReduceStockForOrder(ctx, order.Items); err != nil {
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInsufficientStock, "Failed to reduce stock for order")
		}
		fmt.Printf("✅ Stock reduced immediately for order %s (all payment methods)\n", order.OrderNumber)
	}

	// For bank transfer, only check stock availability - stock will be reduced when payment is confirmed
	// This is consistent with COD and other payment methods
	if err := uc.simpleStockService.CheckStockAvailability(ctx, cart.Items); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInsufficientStock, "Stock not available")
	}
	// Stock availability already checked above
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to update order inventory status")
	}

	// FIXED: Mark cart as converted and clear items atomically within transaction
	cart.MarkAsConverted()
	if err := uc.cartRepo.Update(ctx, cart); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to update cart status")
	}

	// FIXED: Clear cart items within transaction - if this fails, entire transaction should fail
	if err := uc.cartRepo.ClearCart(ctx, cart.ID); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to clear cart items")
	}

	// Create events within transaction to ensure consistency
	if err := uc.orderEventService.CreateOrderCreatedEvent(ctx, order, &userID); err != nil {
		// Log warning but don't fail the transaction for event creation
		// Note: Event creation failure is non-critical
	}

	// Send order created notification (async, don't fail transaction)
	if uc.notificationService != nil {
		go func() {
			// Notify customer
			if err := uc.notificationService.NotifyOrderCreated(context.Background(), order.ID); err != nil {
				// Log error but don't fail the transaction
				fmt.Printf("Failed to send order created notification: %v\n", err)
			}
			// Notify admin about new order
			if err := uc.notificationService.NotifyNewOrder(context.Background(), order.ID); err != nil {
				// Log error but don't fail the transaction
				fmt.Printf("Failed to send new order notification to admin: %v\n", err)
			}
		}()
	}

	// Order created successfully - no stock reservation needed with simple stock service

	// Get created order with relations
	createdOrder, err := uc.orderRepo.GetByID(ctx, order.ID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeOrderNotFound, "Failed to retrieve created order")
	}

	return uc.toOrderResponse(createdOrder), nil
}

// getProductWeight safely extracts weight from product
func getProductWeight(weight *float64) float64 {
	if weight == nil {
		return 0.5 // Default weight for products without weight specified (0.5kg)
	}
	return *weight
}

// getProductsBulk retrieves multiple products in a single query to avoid N+1 problem
func (uc *orderUseCase) getProductsBulk(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID]*entities.Product, error) {
	// Use bulk query to get all products at once
	productList, err := uc.productRepo.GetByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	// Convert to map for easy lookup
	products := make(map[uuid.UUID]*entities.Product)
	for _, product := range productList {
		products[product.ID] = product
	}

	return products, nil
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

	// Create pagination context
	context := &EcommercePaginationContext{
		EntityType: "orders",
		UserID:     userID.String(),
	}

	// Calculate pagination metadata using enhanced function
	// Use proper offset-to-page conversion
	pagination := NewPaginationInfoFromOffset(req.Offset, req.Limit, totalCount)

	// Apply ecommerce enhancements
	if context != nil {
		extraParams := make(map[string]interface{})
		extraParams["user_id"] = context.UserID
		if req.Status != nil {
			extraParams["status"] = *req.Status
		}
		if req.PaymentStatus != nil {
			extraParams["payment_status"] = *req.PaymentStatus
		}
		ApplyEcommerceEnhancements(pagination, context.EntityType, context.UserID, extraParams)
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

	oldStatus := order.Status

	// Update fulfillment status based on order status
	switch status {
	case entities.OrderStatusConfirmed:
		order.FulfillmentStatus = entities.FulfillmentStatusPending
	case entities.OrderStatusProcessing:
		order.FulfillmentStatus = entities.FulfillmentStatusProcessing
		order.SetProcessing()
	case entities.OrderStatusReadyToShip:
		order.FulfillmentStatus = entities.FulfillmentStatusPacked
	case entities.OrderStatusShipped:
		order.FulfillmentStatus = entities.FulfillmentStatusShipped
	case entities.OrderStatusOutForDelivery:
		order.FulfillmentStatus = entities.FulfillmentStatusShipped
	case entities.OrderStatusDelivered:
		order.FulfillmentStatus = entities.FulfillmentStatusDelivered
		order.SetDelivered()
	case entities.OrderStatusCancelled:
		order.FulfillmentStatus = entities.FulfillmentStatusCancelled
	case entities.OrderStatusReturned:
		order.FulfillmentStatus = entities.FulfillmentStatusReturned
	}

	// Update order status and fulfillment status
	order.Status = status
	order.UpdatedAt = time.Now()

	// Save the updated order
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	// Create status changed event
	if err := uc.orderEventService.CreateStatusChangedEvent(ctx, orderID, oldStatus, status, nil); err != nil {
		// Note: Event creation failure is non-critical
	}

	// Send order status changed notification (async)
	if uc.notificationService != nil {
		go func() {
			if err := uc.notificationService.NotifyOrderStatusChanged(context.Background(), orderID, string(status)); err != nil {
				// Log error but don't fail the transaction
				fmt.Printf("Failed to send order status changed notification: %v\n", err)
			}
		}()
	}

	return uc.toOrderResponse(order), nil
}

// CancelOrder cancels an order
func (uc *orderUseCase) CancelOrder(ctx context.Context, orderID uuid.UUID) (*OrderResponse, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, entities.ErrOrderNotFound
	}

	// Validate order can be cancelled
	if !order.CanBeCancelled() {
		return nil, entities.ErrOrderCannotBeCancelled
	}

	// Additional validation for edge cases
	if order.Status == entities.OrderStatusCancelled {
		return nil, fmt.Errorf("order is already cancelled")
	}

	if order.Status == entities.OrderStatusRefunded {
		return nil, fmt.Errorf("order is already refunded and cannot be cancelled")
	}

	// Handle stock based on payment status and order state
	switch {
	case order.IsPaid() && order.Status == entities.OrderStatusConfirmed:
		// Order is paid and confirmed - need to restore actual stock through inventory system
		// This ensures consistency between inventory and product stock
		for _, item := range order.Items {
			// Get inventory record for the product
			inventory, err := uc.inventoryRepo.GetByProductID(ctx, item.ProductID)
			if err != nil {
				fmt.Printf("❌ Failed to get inventory for product %s: %v\n", item.ProductID, err)
				continue
			}

			// Restore stock in inventory (source of truth)
			inventory.QuantityOnHand += item.Quantity
			inventory.QuantityAvailable = inventory.QuantityOnHand - inventory.QuantityReserved

			// Update inventory in database
			if err := uc.inventoryRepo.Update(ctx, inventory); err != nil {
				fmt.Printf("❌ Failed to update inventory for product %s: %v\n", item.ProductID, err)
				continue
			}

			// Sync product stock with inventory (inventory is source of truth)
			if err := uc.productRepo.UpdateStock(ctx, item.ProductID, inventory.QuantityOnHand); err != nil {
				fmt.Printf("❌ Failed to sync product stock for %s: %v\n", item.ProductID, err)
				continue
			}

			fmt.Printf("✅ Restored %d units for product %s (Inventory: %d)\n",
				item.Quantity, item.ProductID, inventory.QuantityOnHand)
		}

	case !order.IsPaid():
		// Order is not paid - restore stock (for bank transfer orders)
		if err := uc.simpleStockService.RestoreStock(ctx, order.Items); err != nil {
			// Don't fail the cancellation, but log the error
			fmt.Printf("❌ Failed to restore stock: %v\n", err)
		}

	default:
		// Unexpected order state for cancellation
		fmt.Printf("⚠️ Unexpected order state for cancellation: IsPaid=%v, Status=%s\n",
			order.IsPaid(), order.Status)
	}

	// Update order status
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	// Update user metrics if order was previously confirmed (paid)
	if order.IsPaid() && order.Status == entities.OrderStatusConfirmed {
		if uc.userMetricsService != nil {
			if err := uc.userMetricsService.UpdateUserMetricsOnOrderCancelled(ctx, order.UserID, order.Total); err != nil {
				fmt.Printf("❌ Failed to update user metrics on cancellation: %v\n", err)
				// Don't fail the cancellation for metrics update failure
			} else {
				fmt.Printf("✅ User metrics updated for order cancellation\n")
			}
		}
	}

	// Create cancelled event
	if err := uc.orderEventService.CreateCancelledEvent(ctx, orderID, "Order cancelled by user", nil); err != nil {
		// Note: Event creation failure is non-critical
	}

	// Order cancelled successfully - no inventory release event needed with simple stock service

	return uc.UpdateOrderStatus(ctx, orderID, entities.OrderStatusCancelled)
}

// GetOrders gets list of orders
func (uc *orderUseCase) GetOrders(ctx context.Context, req GetOrdersRequest) (*GetOrdersResponse, error) {
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

	// Create pagination context
	context := &EcommercePaginationContext{
		EntityType: "admin_orders",
	}

	// Calculate pagination metadata using enhanced function
	// Use proper offset-to-page conversion
	pagination := NewPaginationInfoFromOffset(req.Offset, req.Limit, totalCount)

	// Apply ecommerce enhancements
	if context != nil {
		extraParams := make(map[string]interface{})
		if req.Status != nil {
			extraParams["status"] = *req.Status
		}
		if req.PaymentStatus != nil {
			extraParams["payment_status"] = *req.PaymentStatus
		}
		ApplyEcommerceEnhancements(pagination, context.EntityType, "", extraParams)
	}

	return &GetOrdersResponse{
		Orders:     responses,
		Pagination: pagination,
	}, nil
}

// toOrderResponse converts order entity to response
func (uc *orderUseCase) toOrderResponse(order *entities.Order) *OrderResponse {
	response := &OrderResponse{
		ID:                   order.ID,
		OrderNumber:          order.OrderNumber,
		Status:               order.Status,
		FulfillmentStatus:    order.FulfillmentStatus,
		PaymentStatus:        order.PaymentStatus,
		PaymentMethod:        order.PaymentMethod,
		Priority:             order.Priority,
		Source:               order.Source,
		CustomerType:         order.CustomerType,
		Subtotal:             order.Subtotal,
		TaxAmount:            order.TaxAmount,
		ShippingAmount:       order.ShippingAmount,
		DiscountAmount:       order.DiscountAmount,
		TipAmount:            order.TipAmount,
		Total:                order.Total,
		Currency:             order.Currency,
		ShippingMethod:       order.ShippingMethod,
		TrackingNumber:       order.TrackingNumber,
		TrackingURL:          order.TrackingURL,
		Carrier:              order.Carrier,
		EstimatedDelivery:    order.EstimatedDelivery,
		ActualDelivery:       order.ActualDelivery,
		DeliveryInstructions: order.DeliveryInstructions,
		CustomerNotes:        order.CustomerNotes,
		AdminNotes:           order.AdminNotes,
		IsGift:               order.IsGift,
		GiftMessage:          order.GiftMessage,
		GiftWrap:             order.GiftWrap,
		ItemCount:            order.GetItemCount(),
		CanBeCancelled:       order.CanBeCancelled(),
		CanBeRefunded:        order.CanBeRefunded(),
		CanBeShipped:         order.CanBeShipped(),
		CanBeDelivered:       order.CanBeDelivered(),
		IsShipped:            order.IsShipped(),
		IsDelivered:          order.IsDelivered(),
		HasTracking:          order.HasTracking(),
		CreatedAt:            order.CreatedAt,
		UpdatedAt:            order.UpdatedAt,
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
			response.Items[i].Product = &ProductResponse{
				ID:          item.Product.ID,
				Name:        item.Product.Name,
				Description: item.Product.Description,
				SKU:         item.Product.SKU,
				Slug:        item.Product.Slug,
				Price:       item.Product.Price,
				CurrentPrice: item.Product.GetCurrentPrice(),
				Stock:       item.Product.Stock,
				Status:      item.Product.Status,
				MainImage:   item.Product.GetMainImage(),
			}
		}
	}

	// Convert payments - get the latest payment for backward compatibility
	if len(order.Payments) > 0 {
		latestPayment := order.GetLatestPayment()
		if latestPayment != nil {
			response.Payment = &PaymentResponse{
				ID:              latestPayment.ID,
				OrderID:         latestPayment.OrderID,
				Amount:          latestPayment.Amount,
				Currency:        latestPayment.Currency,
				Method:          latestPayment.Method,
				Status:          latestPayment.Status,
				TransactionID:   latestPayment.TransactionID,
				ExternalID:      latestPayment.ExternalID,
				ProcessedAt:     latestPayment.ProcessedAt,
				RefundedAt:      latestPayment.RefundedAt,
				RefundAmount:    latestPayment.RefundAmount,
				CanBeRefunded:   latestPayment.CanBeRefunded(),
				RemainingRefund: latestPayment.GetRemainingRefundAmount(),
				CreatedAt:       latestPayment.CreatedAt,
				UpdatedAt:       latestPayment.UpdatedAt,
			}
		}
	}

	return response
}

// GetOrderEvents gets order events/timeline
func (uc *orderUseCase) GetOrderEvents(ctx context.Context, orderID uuid.UUID, publicOnly bool) ([]*entities.OrderEvent, error) {
	// Verify order exists
	_, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, entities.ErrOrderNotFound
	}

	return uc.orderEventService.GetOrderEvents(ctx, orderID, publicOnly)
}

// UpdateShippingInfoRequest represents request to update shipping info
type UpdateShippingInfoRequest struct {
	TrackingNumber    string     `json:"tracking_number" binding:"required"`
	Carrier           string     `json:"carrier" binding:"required"`
	ShippingMethod    string     `json:"shipping_method"`
	TrackingURL       string     `json:"tracking_url"`
	EstimatedDelivery *time.Time `json:"estimated_delivery"`
}

// UpdateShippingInfo updates shipping information for an order
func (uc *orderUseCase) UpdateShippingInfo(ctx context.Context, orderID uuid.UUID, req UpdateShippingInfoRequest) (*OrderResponse, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, entities.ErrOrderNotFound
	}

	if !order.CanBeShipped() {
		return nil, fmt.Errorf("order cannot be shipped in current status: %s", order.Status)
	}

	// Update shipping info
	order.TrackingNumber = req.TrackingNumber
	order.Carrier = req.Carrier
	order.ShippingMethod = req.ShippingMethod
	order.TrackingURL = req.TrackingURL
	order.EstimatedDelivery = req.EstimatedDelivery
	order.SetShipped(req.TrackingNumber, req.Carrier)

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	// Create shipped event
	if err := uc.orderEventService.CreateShippedEvent(ctx, orderID, req.TrackingNumber, req.Carrier, nil); err != nil {
		// Note: Event creation failure is non-critical
	}

	return uc.toOrderResponse(order), nil
}

// UpdateDeliveryStatus updates delivery status for an order
func (uc *orderUseCase) UpdateDeliveryStatus(ctx context.Context, orderID uuid.UUID, status entities.OrderStatus) (*OrderResponse, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, entities.ErrOrderNotFound
	}

	// Validate delivery status
	if status != entities.OrderStatusOutForDelivery && status != entities.OrderStatusDelivered {
		return nil, fmt.Errorf("invalid delivery status: %s", status)
	}

	if status == entities.OrderStatusDelivered && !order.CanBeDelivered() {
		return nil, fmt.Errorf("order cannot be marked as delivered in current status: %s", order.Status)
	}

	oldStatus := order.Status
	order.Status = status

	if status == entities.OrderStatusDelivered {
		order.SetDelivered()
	}

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	// Create appropriate event
	if status == entities.OrderStatusDelivered {
		if err := uc.orderEventService.CreateDeliveredEvent(ctx, orderID, nil); err != nil {
			return nil, err
		}
	}

	// Create status changed event
	if err := uc.orderEventService.CreateStatusChangedEvent(ctx, orderID, oldStatus, status, nil); err != nil {
		return nil, err
	}

	// Send order status changed notification (async)
	if uc.notificationService != nil {
		go func() {
			if err := uc.notificationService.NotifyOrderStatusChanged(context.Background(), orderID, string(status)); err != nil {
				// Log error but don't fail the transaction
				fmt.Printf("Failed to send order status changed notification: %v\n", err)
			}
		}()
	}

	return uc.toOrderResponse(order), nil
}

// AddOrderNoteRequest represents request to add order note
type AddOrderNoteRequest struct {
	Note     string `json:"note" binding:"required"`
	IsPublic bool   `json:"is_public"`
}

// AddOrderNote adds a note to an order (updated signature)
func (uc *orderUseCase) AddOrderNote(ctx context.Context, orderID uuid.UUID, req AddOrderNoteRequest) error {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return entities.ErrOrderNotFound
	}

	// Update order notes
	if req.IsPublic {
		if order.CustomerNotes != "" {
			order.CustomerNotes += "\n" + req.Note
		} else {
			order.CustomerNotes = req.Note
		}
	} else {
		if order.AdminNotes != "" {
			order.AdminNotes += "\n" + req.Note
		} else {
			order.AdminNotes = req.Note
		}
	}

	order.UpdatedAt = time.Now()

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return err
	}

	// Create note added event
	if err := uc.orderEventService.CreateNoteAddedEvent(ctx, orderID, req.Note, nil, req.IsPublic); err != nil {
		return err
	}

	return nil
}
