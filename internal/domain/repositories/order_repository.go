package repositories

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"

	"github.com/google/uuid"
)

// OrderSearchParams represents search parameters for orders
type OrderSearchParams struct {
	UserID        *uuid.UUID
	Status        *entities.OrderStatus
	PaymentStatus *entities.PaymentStatus
	StartDate     *time.Time
	EndDate       *time.Time
	MinTotal      *float64
	MaxTotal      *float64
	SortBy        string // created_at, total, status
	SortOrder     string // asc, desc
	Limit         int
	Offset        int
}

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	// Create creates a new order
	Create(ctx context.Context, order *entities.Order) error

	// GetByID retrieves an order by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Order, error)

	// GetByOrderNumber retrieves an order by order number
	GetByOrderNumber(ctx context.Context, orderNumber string) (*entities.Order, error)

	// Update updates an existing order
	Update(ctx context.Context, order *entities.Order) error

	// Delete deletes an order by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves orders with pagination
	List(ctx context.Context, limit, offset int) ([]*entities.Order, error)

	// Search searches orders based on criteria
	Search(ctx context.Context, params OrderSearchParams) ([]*entities.Order, error)

	// GetByUserID retrieves orders by user ID
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Order, error)

	// Count returns the total number of orders
	Count(ctx context.Context) (int64, error)

	// CountByUser returns the number of orders for a user
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)

	// UpdateStatus updates order status
	UpdateStatus(ctx context.Context, orderID uuid.UUID, status entities.OrderStatus) error

	// UpdatePaymentStatus updates payment status
	UpdatePaymentStatus(ctx context.Context, orderID uuid.UUID, status entities.PaymentStatus) error

	// GetRecentOrders retrieves recent orders
	GetRecentOrders(ctx context.Context, limit int) ([]*entities.Order, error)

	// GetOrdersByDateRange retrieves orders within a date range
	GetOrdersByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entities.Order, error)

	// GetTotalSales calculates total sales within a date range
	GetTotalSales(ctx context.Context, startDate, endDate time.Time) (float64, error)

	// Additional methods for admin dashboard
	GetTotalRevenue(ctx context.Context) (float64, error)    // Net revenue (total)
	GetGrossRevenue(ctx context.Context) (float64, error)    // Before discounts
	GetProductRevenue(ctx context.Context) (float64, error)  // Only subtotal
	GetTaxCollected(ctx context.Context) (float64, error)    // Total tax amount
	GetShippingRevenue(ctx context.Context) (float64, error) // Total shipping fees
	GetDiscountsGiven(ctx context.Context) (float64, error)  // Total discounts
	CountOrders(ctx context.Context) (int64, error)
	CountOrdersByStatus(ctx context.Context, status entities.OrderStatus) (int64, error)
}

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	// Create creates a new payment
	Create(ctx context.Context, payment *entities.Payment) error

	// GetByID retrieves a payment by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error)

	// GetByOrderID retrieves a payment by order ID
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*entities.Payment, error)

	// GetByTransactionID retrieves a payment by transaction ID
	GetByTransactionID(ctx context.Context, transactionID string) (*entities.Payment, error)

	// GetByExternalID retrieves a payment by external ID (e.g., Stripe session ID)
	GetByExternalID(ctx context.Context, externalID string) (*entities.Payment, error)

	// Update updates an existing payment
	Update(ctx context.Context, payment *entities.Payment) error

	// Delete deletes a payment by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves payments with pagination
	List(ctx context.Context, limit, offset int) ([]*entities.Payment, error)

	// GetByStatus retrieves payments by status
	GetByStatus(ctx context.Context, status entities.PaymentStatus, limit, offset int) ([]*entities.Payment, error)

	// Count returns the total number of payments
	Count(ctx context.Context) (int64, error)

	// GetFailedPayments retrieves failed payments
	GetFailedPayments(ctx context.Context, limit, offset int) ([]*entities.Payment, error)

	// GetRefundablePayments retrieves payments that can be refunded
	GetRefundablePayments(ctx context.Context, limit, offset int) ([]*entities.Payment, error)

	// Refund-related methods
	CreateRefund(ctx context.Context, refund *entities.Refund) error
	GetRefund(ctx context.Context, refundID uuid.UUID) (*entities.Refund, error)
	GetRefundsByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]*entities.Refund, error)
	UpdateRefund(ctx context.Context, refund *entities.Refund) error
	ListRefunds(ctx context.Context, limit, offset int) ([]*entities.Refund, error)
}
