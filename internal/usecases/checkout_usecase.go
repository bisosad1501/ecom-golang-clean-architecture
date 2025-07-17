package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/domain/services"
	"ecom-golang-clean-architecture/internal/infrastructure/database"
	pkgErrors "ecom-golang-clean-architecture/pkg/errors"
)

// PaymentUseCaseInterface interface for payment operations (to avoid conflict)
type PaymentUseCaseInterface interface {
	CreateCheckoutSession(ctx context.Context, req CreateCheckoutSessionRequest) (*CreateCheckoutSessionResponse, error)
}

// CheckoutUseCase defines checkout use cases
type CheckoutUseCase interface {
	// Create checkout session for online payments
	CreateCheckoutSession(ctx context.Context, userID uuid.UUID, req CreateNewCheckoutSessionRequest) (*NewCheckoutSessionResponse, error)

	// Complete checkout session (after payment success)
	CompleteCheckoutSession(ctx context.Context, sessionID string) (*OrderResponse, error)

	// Create order directly for COD
	CreateCODOrder(ctx context.Context, userID uuid.UUID, req CreateOrderRequest) (*OrderResponse, error)

	// Get checkout session
	GetCheckoutSession(ctx context.Context, sessionID string) (*NewCheckoutSessionResponse, error)

	// Cancel checkout session
	CancelCheckoutSession(ctx context.Context, sessionID string) error
}

// CreateNewCheckoutSessionRequest represents create checkout session request
type CreateNewCheckoutSessionRequest struct {
	ShippingAddress AddressRequest         `json:"shipping_address" validate:"required"`
	BillingAddress  *AddressRequest        `json:"billing_address"`
	PaymentMethod   entities.PaymentMethod `json:"payment_method" validate:"required"`
	Notes           string                 `json:"notes"`
	TaxRate         float64                `json:"tax_rate" validate:"min=0,max=1"`
	ShippingCost    float64                `json:"shipping_cost" validate:"min=0"`
	DiscountAmount  float64                `json:"discount_amount" validate:"min=0"`
}

// NewCheckoutSessionResponse represents checkout session response
type NewCheckoutSessionResponse struct {
	ID              uuid.UUID                     `json:"id"`
	SessionID       string                        `json:"session_id"`
	Status          entities.CheckoutSessionStatus `json:"status"`
	PaymentMethod   entities.PaymentMethod        `json:"payment_method"`
	PaymentIntentID string                        `json:"payment_intent_id,omitempty"`
	StripeURL       string                        `json:"stripe_url,omitempty"`
	Subtotal        float64                       `json:"subtotal"`
	TaxAmount       float64                       `json:"tax_amount"`
	ShippingAmount  float64                       `json:"shipping_amount"`
	DiscountAmount  float64                       `json:"discount_amount"`
	Total           float64                       `json:"total"`
	Currency        string                        `json:"currency"`
	ExpiresAt       *time.Time                    `json:"expires_at"`
	CreatedAt       time.Time                     `json:"created_at"`
}

type checkoutUseCase struct {
	checkoutRepo    repositories.CheckoutSessionRepository
	cartRepo        repositories.CartRepository
	orderRepo       repositories.OrderRepository
	productRepo     repositories.ProductRepository
	stockService    services.SimpleStockService
	orderService    services.OrderService
	paymentUseCase  PaymentUseCaseInterface
	txManager       *database.TransactionManager
}

// NewCheckoutUseCase creates a new checkout use case
func NewCheckoutUseCase(
	checkoutRepo repositories.CheckoutSessionRepository,
	cartRepo repositories.CartRepository,
	orderRepo repositories.OrderRepository,
	productRepo repositories.ProductRepository,
	stockService services.SimpleStockService,
	orderService services.OrderService,
	paymentUseCase PaymentUseCaseInterface,
	txManager *database.TransactionManager,
) CheckoutUseCase {
	return &checkoutUseCase{
		checkoutRepo:   checkoutRepo,
		cartRepo:       cartRepo,
		orderRepo:      orderRepo,
		productRepo:    productRepo,
		stockService:   stockService,
		orderService:   orderService,
		paymentUseCase: paymentUseCase,
		txManager:      txManager,
	}
}

// CreateCheckoutSession creates a checkout session for online payments
func (uc *checkoutUseCase) CreateCheckoutSession(ctx context.Context, userID uuid.UUID, req CreateNewCheckoutSessionRequest) (*NewCheckoutSessionResponse, error) {
	// Validate request
	if err := uc.validateCheckoutRequest(req); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInvalidInput, "Invalid checkout request")
	}

	// Only allow online payment methods for checkout sessions
	if req.PaymentMethod == entities.PaymentMethodCash {
		return nil, pkgErrors.InvalidInput("COD orders should use direct order creation")
	}

	// Get user's cart
	cart, err := uc.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, pkgErrors.CartNotFound()
	}

	if cart.IsEmpty() {
		return nil, pkgErrors.InvalidInput("Cart is empty")
	}

	// Check stock availability
	if err := uc.stockService.CheckStockAvailability(ctx, cart.Items); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInsufficientStock, "Stock not available")
	}

	// Calculate totals
	subtotal, taxAmount, total := uc.orderService.CalculateOrderTotal(
		cart.Items, req.TaxRate, req.ShippingCost, req.DiscountAmount,
	)

	// Create checkout session
	session := &entities.CheckoutSession{
		ID:              uuid.New(),
		UserID:          userID,
		CartID:          cart.ID,
		CartItems:       cart.Items, // Snapshot
		PaymentMethod:   req.PaymentMethod,
		Subtotal:        subtotal,
		TaxAmount:       taxAmount,
		ShippingAmount:  req.ShippingCost,
		DiscountAmount:  req.DiscountAmount,
		Total:           total,
		Currency:        "USD",
		TaxRate:         req.TaxRate,
		ShippingCost:    req.ShippingCost,
		Notes:           req.Notes,
		Status:          entities.CheckoutSessionStatusActive,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Set addresses
	session.ShippingAddress = &entities.OrderAddress{
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
		session.BillingAddress = &entities.OrderAddress{
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
		// Copy shipping address to billing address
		session.BillingAddress = session.ShippingAddress
	}

	// Generate session ID and set expiration
	session.GenerateSessionID()
	session.SetExpiration(15) // 15 minutes for online payments

	// For Stripe payment method, create Stripe checkout session
	if req.PaymentMethod == entities.PaymentMethodStripe {
		fmt.Printf("🔍 Processing Stripe payment method\n")

		// Create a simple temporary order for Stripe
		tempOrderID := uuid.New()
		tempOrder := &entities.Order{
			ID:             tempOrderID,
			OrderNumber:    fmt.Sprintf("STRIPE-TEMP-%s", session.SessionID),
			UserID:         userID,
			Status:         entities.OrderStatusPending,
			PaymentStatus:  entities.PaymentStatusPending,
			PaymentMethod:  entities.PaymentMethodStripe,
			Subtotal:       subtotal,
			TaxAmount:      taxAmount,
			ShippingAmount: req.ShippingCost,
			DiscountAmount: req.DiscountAmount,
			Total:          total,
			Currency:       "USD",
			Source:         entities.OrderSourceWeb,
			CustomerType:   entities.CustomerTypeRegistered,
			Priority:       entities.OrderPriorityNormal,
			Version:        1,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Set addresses
		tempOrder.ShippingAddress = session.ShippingAddress
		tempOrder.BillingAddress = session.BillingAddress

		// Add items to temp order
		for _, cartItem := range cart.Items {
			orderItem := entities.OrderItem{
				ID:          uuid.New(),
				OrderID:     tempOrder.ID,
				ProductID:   cartItem.ProductID,
				ProductName: cartItem.Product.Name,
				ProductSKU:  cartItem.Product.SKU,
				Quantity:    cartItem.Quantity,
				Price:       cartItem.Price,
				Total:       cartItem.Total,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			tempOrder.Items = append(tempOrder.Items, orderItem)
		}

		fmt.Printf("🔍 Creating temporary order with ID: %s\n", tempOrder.ID)
		// Save temp order
		if err := uc.orderRepo.Create(ctx, tempOrder); err != nil {
			fmt.Printf("❌ Failed to create temporary order: %v\n", err)
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to create temporary order")
		}
		fmt.Printf("✅ Temporary order created successfully\n")

		// Create Stripe checkout session
		stripeReq := CreateCheckoutSessionRequest{
			OrderID:     tempOrder.ID,
			Amount:      total,
			Currency:    "usd",
			Description: fmt.Sprintf("Payment for checkout session %s", session.SessionID),
			SuccessURL:  fmt.Sprintf("%s/checkout/success?session_id=%s&order_id=%s", "http://localhost:3000", session.SessionID, tempOrder.ID.String()),
			CancelURL:   fmt.Sprintf("%s/checkout/cancel?session_id=%s&order_id=%s", "http://localhost:3000", session.SessionID, tempOrder.ID.String()),
			Metadata: map[string]interface{}{
				"checkout_session_id": session.SessionID,
				"user_id":             userID.String(),
				"order_id":            tempOrder.ID.String(),
			},
		}

		fmt.Printf("🔍 Creating Stripe checkout session with request: %+v\n", stripeReq)
		stripeResp, err := uc.paymentUseCase.CreateCheckoutSession(ctx, stripeReq)
		if err != nil {
			fmt.Printf("❌ Stripe checkout session error: %v\n", err)
			return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to create Stripe checkout session")
		}

		fmt.Printf("✅ Stripe checkout session response: %+v\n", stripeResp)
		if !stripeResp.Success {
			fmt.Printf("❌ Stripe checkout session failed: %s\n", stripeResp.Message)
			return nil, pkgErrors.InvalidInput(stripeResp.Message)
		}

		// Store Stripe session ID and URL in our checkout session
		session.PaymentIntentID = stripeResp.SessionID
		// Store Stripe URL in notes for now (can be returned in response)
		session.Notes = fmt.Sprintf("Stripe URL: %s", stripeResp.SessionURL)
		fmt.Printf("✅ Stripe checkout session created: %s\n", stripeResp.SessionID)
	}

	// Validate and save
	if err := session.Validate(); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInvalidInput, "Invalid session data")
	}

	if err := uc.checkoutRepo.Create(ctx, session); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to create checkout session")
	}

	response := uc.toCheckoutSessionResponse(session)

	// Log Stripe URL if available
	if req.PaymentMethod == entities.PaymentMethodStripe && response.StripeURL != "" {
		fmt.Printf("✅ Stripe checkout URL available: %s\n", response.StripeURL)
	}

	return response, nil
}

// validateCheckoutRequest validates checkout request
func (uc *checkoutUseCase) validateCheckoutRequest(req CreateNewCheckoutSessionRequest) error {
	// Validate payment method
	validPaymentMethods := []entities.PaymentMethod{
		entities.PaymentMethodCreditCard,
		entities.PaymentMethodDebitCard,
		entities.PaymentMethodPayPal,
		entities.PaymentMethodStripe,
		entities.PaymentMethodApplePay,
		entities.PaymentMethodGooglePay,
		entities.PaymentMethodBankTransfer,
	}

	isValidPaymentMethod := false
	for _, method := range validPaymentMethods {
		if req.PaymentMethod == method {
			isValidPaymentMethod = true
			break
		}
	}
	if !isValidPaymentMethod {
		return fmt.Errorf("invalid payment method for checkout session: %s", req.PaymentMethod)
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

// CompleteCheckoutSession completes checkout session after payment success
func (uc *checkoutUseCase) CompleteCheckoutSession(ctx context.Context, sessionID string) (*OrderResponse, error) {
	// Execute in transaction
	result, err := uc.txManager.WithTransactionResult(ctx, func(tx *gorm.DB) (interface{}, error) {
		return uc.completeCheckoutSessionInTransaction(ctx, sessionID)
	})
	if err != nil {
		return nil, err
	}
	return result.(*OrderResponse), nil
}

// completeCheckoutSessionInTransaction handles checkout completion in transaction
func (uc *checkoutUseCase) completeCheckoutSessionInTransaction(ctx context.Context, sessionID string) (*OrderResponse, error) {
	// Get checkout session
	session, err := uc.checkoutRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeNotFound, "Checkout session not found")
	}

	// Validate session can be completed
	if !session.CanBeCompleted() {
		return nil, pkgErrors.InvalidInput("Checkout session cannot be completed")
	}

	// Check stock availability again
	if err := uc.stockService.CheckStockAvailability(ctx, session.CartItems); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInsufficientStock, "Stock not available")
	}

	// Generate order number
	orderNumber, err := uc.orderService.GenerateUniqueOrderNumber(ctx)
	if err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to generate order number")
	}

	// Create order from session
	order := &entities.Order{
		ID:             uuid.New(),
		OrderNumber:    orderNumber,
		UserID:         session.UserID,
		Status:         entities.OrderStatusConfirmed, // Confirmed because payment is already successful
		PaymentStatus:  entities.PaymentStatusPaid,
		PaymentMethod:  session.PaymentMethod,
		Subtotal:       session.Subtotal,
		TaxAmount:      session.TaxAmount,
		ShippingAmount: session.ShippingAmount,
		DiscountAmount: session.DiscountAmount,
		Total:          session.Total,
		Currency:       session.Currency,
		CustomerNotes:  session.Notes,
		Source:         entities.OrderSourceWeb,
		CustomerType:   entities.CustomerTypeRegistered,
		Priority:       entities.OrderPriorityNormal,
		Version:        1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Set addresses
	order.ShippingAddress = session.ShippingAddress
	order.BillingAddress = session.BillingAddress

	// Create order items
	for _, cartItem := range session.CartItems {
		orderItem := entities.OrderItem{
			ID:          uuid.New(),
			OrderID:     order.ID,
			ProductID:   cartItem.ProductID,
			ProductName: cartItem.Product.Name,
			ProductSKU:  cartItem.Product.SKU,
			Quantity:    cartItem.Quantity,
			Price:       cartItem.Price,
			Total:       cartItem.Total,
		}
		order.Items = append(order.Items, orderItem)
	}

	// Validate order
	if err := order.Validate(); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInvalidInput, "Invalid order data")
	}

	// Save order
	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to create order")
	}

	// Reduce stock
	if err := uc.stockService.ReduceStock(ctx, session.CartItems); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInsufficientStock, "Failed to reduce stock")
	}

	// Mark session as completed
	session.MarkAsCompleted(order.ID)
	if err := uc.checkoutRepo.Update(ctx, session); err != nil {
		// Log warning but don't fail
		fmt.Printf("Warning: Failed to update checkout session: %v\n", err)
	}

	// Clear cart
	if err := uc.cartRepo.ClearCart(ctx, session.CartID); err != nil {
		// Log warning but don't fail
		fmt.Printf("Warning: Failed to clear cart: %v\n", err)
	}

	// Get created order with relations
	createdOrder, err := uc.orderRepo.GetByID(ctx, order.ID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeOrderNotFound, "Failed to retrieve created order")
	}

	return toOrderResponse(createdOrder), nil
}

// CreateCODOrder creates order directly for COD payments
func (uc *checkoutUseCase) CreateCODOrder(ctx context.Context, userID uuid.UUID, req CreateOrderRequest) (*OrderResponse, error) {
	// Execute in transaction
	result, err := uc.txManager.WithTransactionResult(ctx, func(tx *gorm.DB) (interface{}, error) {
		return uc.createCODOrderInTransaction(ctx, userID, req)
	})
	if err != nil {
		return nil, err
	}
	return result.(*OrderResponse), nil
}

// createCODOrderInTransaction handles COD order creation in transaction
func (uc *checkoutUseCase) createCODOrderInTransaction(ctx context.Context, userID uuid.UUID, req CreateOrderRequest) (*OrderResponse, error) {
	// Validate request
	if req.PaymentMethod != entities.PaymentMethodCash {
		return nil, pkgErrors.InvalidInput("This method is only for COD orders")
	}

	// Get user's cart
	cart, err := uc.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, pkgErrors.CartNotFound()
	}

	if cart.IsEmpty() {
		return nil, pkgErrors.InvalidInput("Cart is empty")
	}

	// Check stock availability and reduce immediately for COD
	if err := uc.stockService.CheckStockAvailability(ctx, cart.Items); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInsufficientStock, "Stock not available")
	}

	// Calculate totals
	subtotal, taxAmount, total := uc.orderService.CalculateOrderTotal(
		cart.Items, req.TaxRate, req.ShippingCost, req.DiscountAmount,
	)

	// Generate order number
	orderNumber, err := uc.orderService.GenerateUniqueOrderNumber(ctx)
	if err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to generate order number")
	}

	// Create order
	order := &entities.Order{
		ID:             uuid.New(),
		OrderNumber:    orderNumber,
		UserID:         userID,
		Status:         entities.OrderStatusPending, // Pending for COD
		PaymentStatus:  entities.PaymentStatusAwaitingPayment,
		PaymentMethod:  entities.PaymentMethodCash,
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

	// Set addresses (same logic as before)
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
		orderItem := entities.OrderItem{
			ID:          uuid.New(),
			OrderID:     order.ID,
			ProductID:   cartItem.ProductID,
			ProductName: cartItem.Product.Name,
			ProductSKU:  cartItem.Product.SKU,
			Quantity:    cartItem.Quantity,
			Price:       cartItem.Price,
			Total:       cartItem.Total,
		}
		order.Items = append(order.Items, orderItem)
	}

	// Validate order
	if err := order.Validate(); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInvalidInput, "Invalid order data")
	}

	// Save order
	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInternalError, "Failed to create order")
	}

	// For COD, only check stock availability - don't reduce until delivery confirmed
	if err := uc.stockService.CheckStockAvailability(ctx, cart.Items); err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeInsufficientStock, "Stock not available")
	}

	// Clear cart
	if err := uc.cartRepo.ClearCart(ctx, cart.ID); err != nil {
		// Log warning but don't fail
		fmt.Printf("Warning: Failed to clear cart: %v\n", err)
	}

	// Get created order with relations
	createdOrder, err := uc.orderRepo.GetByID(ctx, order.ID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeOrderNotFound, "Failed to retrieve created order")
	}

	return toOrderResponse(createdOrder), nil
}

// GetCheckoutSession gets checkout session by session ID
func (uc *checkoutUseCase) GetCheckoutSession(ctx context.Context, sessionID string) (*NewCheckoutSessionResponse, error) {
	session, err := uc.checkoutRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, pkgErrors.ErrCodeNotFound, "Checkout session not found")
	}

	return uc.toCheckoutSessionResponse(session), nil
}

// CancelCheckoutSession cancels a checkout session
func (uc *checkoutUseCase) CancelCheckoutSession(ctx context.Context, sessionID string) error {
	session, err := uc.checkoutRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		return pkgErrors.Wrap(err, pkgErrors.ErrCodeNotFound, "Checkout session not found")
	}

	session.MarkAsCancelled()
	return uc.checkoutRepo.Update(ctx, session)
}

// toCheckoutSessionResponse converts entity to response
func (uc *checkoutUseCase) toCheckoutSessionResponse(session *entities.CheckoutSession) *NewCheckoutSessionResponse {
	response := &NewCheckoutSessionResponse{
		ID:              session.ID,
		SessionID:       session.SessionID,
		Status:          session.Status,
		PaymentMethod:   session.PaymentMethod,
		PaymentIntentID: session.PaymentIntentID,
		Subtotal:        session.Subtotal,
		TaxAmount:       session.TaxAmount,
		ShippingAmount:  session.ShippingAmount,
		DiscountAmount:  session.DiscountAmount,
		Total:           session.Total,
		Currency:        session.Currency,
		ExpiresAt:       session.ExpiresAt,
		CreatedAt:       session.CreatedAt,
	}

	// Extract Stripe URL from notes if available
	if session.PaymentMethod == entities.PaymentMethodStripe && session.Notes != "" && strings.Contains(session.Notes, "Stripe URL: ") {
		response.StripeURL = strings.TrimPrefix(session.Notes, "Stripe URL: ")
	}

	return response
}

// toOrderResponse converts order entity to response (simplified version)
func toOrderResponse(order *entities.Order) *OrderResponse {
	response := &OrderResponse{
		ID:                order.ID,
		OrderNumber:       order.OrderNumber,
		Status:            order.Status,
		FulfillmentStatus: order.FulfillmentStatus,
		PaymentStatus:     order.PaymentStatus,
		PaymentMethod:     order.PaymentMethod,
		Priority:          order.Priority,
		Source:            order.Source,
		CustomerType:      order.CustomerType,
		Subtotal:          order.Subtotal,
		TaxAmount:         order.TaxAmount,
		ShippingAmount:    order.ShippingAmount,
		DiscountAmount:    order.DiscountAmount,
		TipAmount:         order.TipAmount,
		Total:             order.Total,
		Currency:          order.Currency,
		CustomerNotes:     order.CustomerNotes,
		AdminNotes:        order.AdminNotes,
		IsGift:            order.IsGift,
		GiftMessage:       order.GiftMessage,
		GiftWrap:          order.GiftWrap,
		ItemCount:         len(order.Items),
		CanBeCancelled:    order.CanBeCancelled(),
		CanBeRefunded:     order.CanBeRefunded(),
		CanBeShipped:      order.CanBeShipped(),
		CanBeDelivered:    order.CanBeDelivered(),
		IsShipped:         order.IsShipped(),
		IsDelivered:       order.IsDelivered(),
		HasTracking:       order.HasTracking(),
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
	}

	// Convert user
	if order.User.ID != uuid.Nil {
		response.User = &UserResponse{
			ID:        order.User.ID,
			Email:     order.User.Email,
			FirstName: order.User.FirstName,
			LastName:  order.User.LastName,
			Phone:     order.User.Phone,
			Role:      order.User.Role,
			IsActive:  order.User.IsActive,
			CreatedAt: order.User.CreatedAt,
			UpdatedAt: order.User.UpdatedAt,
		}
	}

	// Convert items
	for _, item := range order.Items {
		orderItem := OrderItemResponse{
			ID:          item.ID,
			ProductName: item.ProductName,
			ProductSKU:  item.ProductSKU,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Total:       item.Total,
		}

		// Add product details if available
		if item.Product.ID != uuid.Nil {
			orderItem.Product = &ProductResponse{
				ID:          item.Product.ID,
				Name:        item.Product.Name,
				Description: item.Product.Description,
				SKU:         item.Product.SKU,
				Price:       item.Product.Price,
				Stock:       item.Product.Stock,
				Status:      item.Product.Status,
				CreatedAt:   item.Product.CreatedAt,
				UpdatedAt:   item.Product.UpdatedAt,
			}
		}

		response.Items = append(response.Items, orderItem)
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

	return response
}
