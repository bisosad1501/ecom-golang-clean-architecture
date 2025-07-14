package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *gorm.DB) repositories.OrderRepository {
	return &orderRepository{db: db}
}

// Create creates a new order
func (r *orderRepository) Create(ctx context.Context, order *entities.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// GetByID retrieves an order by ID
func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	var order entities.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Images").
		Preload("Items.Product.Category").
		Preload("Payments").
		Where("id = ?", id).
		First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

// GetByOrderNumber retrieves an order by order number
func (r *orderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*entities.Order, error) {
	var order entities.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Images").
		Preload("Items.Product.Category").
		Preload("Payments").
		Where("order_number = ?", orderNumber).
		First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

// ExistsByOrderNumber checks if an order exists with the given order number
func (r *orderRepository) ExistsByOrderNumber(ctx context.Context, orderNumber string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("order_number = ?", orderNumber).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update updates an existing order
func (r *orderRepository) Update(ctx context.Context, order *entities.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// Delete deletes an order by ID
func (r *orderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.Order{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrOrderNotFound
	}
	return nil
}

// List retrieves orders with pagination
func (r *orderRepository) List(ctx context.Context, limit, offset int) ([]*entities.Order, error) {
	var orders []*entities.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Payments").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// Search searches orders based on criteria
func (r *orderRepository) Search(ctx context.Context, params repositories.OrderSearchParams) ([]*entities.Order, error) {
	query := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Payments")

	// Apply filters
	if params.UserID != nil {
		query = query.Where("user_id = ?", *params.UserID)
	}

	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	if params.PaymentStatus != nil {
		query = query.Where("payment_status = ?", *params.PaymentStatus)
	}

	if params.StartDate != nil {
		query = query.Where("created_at >= ?", *params.StartDate)
	}

	if params.EndDate != nil {
		query = query.Where("created_at <= ?", *params.EndDate)
	}

	if params.MinTotal != nil {
		query = query.Where("total >= ?", *params.MinTotal)
	}

	if params.MaxTotal != nil {
		query = query.Where("total <= ?", *params.MaxTotal)
	}

	// Apply sorting
	orderBy := "created_at DESC"
	if params.SortBy != "" {
		direction := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			direction = "DESC"
		}
		orderBy = params.SortBy + " " + direction
	}
	query = query.Order(orderBy)

	// Apply pagination
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	var orders []*entities.Order
	err := query.Find(&orders).Error
	return orders, err
}

// CountSearch counts orders based on search criteria
func (r *orderRepository) CountSearch(ctx context.Context, params repositories.OrderSearchParams) (int64, error) {
	query := r.db.WithContext(ctx).Model(&entities.Order{})

	// Apply the same filters as Search method
	if params.UserID != nil {
		query = query.Where("user_id = ?", *params.UserID)
	}

	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	if params.PaymentStatus != nil {
		query = query.Where("payment_status = ?", *params.PaymentStatus)
	}

	if params.StartDate != nil {
		query = query.Where("created_at >= ?", *params.StartDate)
	}

	if params.EndDate != nil {
		query = query.Where("created_at <= ?", *params.EndDate)
	}

	if params.MinTotal != nil {
		query = query.Where("total >= ?", *params.MinTotal)
	}

	if params.MaxTotal != nil {
		query = query.Where("total <= ?", *params.MaxTotal)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// GetByUserID retrieves orders by user ID
func (r *orderRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Order, error) {
	var orders []*entities.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("position >= 0").Order("position ASC")
		}).
		Preload("Payments").
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// Count returns the total number of orders
func (r *orderRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Order{}).Count(&count).Error
	return count, err
}

// CountByUser returns the number of orders for a user
func (r *orderRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// UpdateStatus updates order status
func (r *orderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status entities.OrderStatus) error {
	result := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("id = ?", orderID).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrOrderNotFound
	}
	return nil
}

// UpdatePaymentStatus updates payment status
func (r *orderRepository) UpdatePaymentStatus(ctx context.Context, orderID uuid.UUID, status entities.PaymentStatus) error {
	result := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("id = ?", orderID).
		Update("payment_status", status)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrOrderNotFound
	}
	return nil
}

// GetRecentOrders retrieves recent orders
func (r *orderRepository) GetRecentOrders(ctx context.Context, limit int) ([]*entities.Order, error) {
	var orders []*entities.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Payments").
		Limit(limit).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// GetOrdersByDateRange retrieves orders within a date range
func (r *orderRepository) GetOrdersByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entities.Order, error) {
	var orders []*entities.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Payments").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// GetTotalSales calculates total sales within a date range
func (r *orderRepository) GetTotalSales(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("created_at BETWEEN ? AND ? AND payment_status = ?", startDate, endDate, entities.PaymentStatusPaid).
		Select("COALESCE(SUM(total), 0)").
		Scan(&total).Error
	return total, err
}

// GetTotalRevenue gets total revenue from all orders
func (r *orderRepository) GetTotalRevenue(ctx context.Context) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("payment_status = ? AND status NOT IN ?",
			entities.PaymentStatusPaid,
			[]entities.OrderStatus{entities.OrderStatusCancelled, entities.OrderStatusRefunded}).
		Select("COALESCE(SUM(total), 0)").
		Scan(&total).Error
	return total, err
}

// CountOrders counts total number of orders
func (r *orderRepository) CountOrders(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Count(&count).Error
	return count, err
}

// CountOrdersByStatus counts orders by status
func (r *orderRepository) CountOrdersByStatus(ctx context.Context, status entities.OrderStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

// GetGrossRevenue gets gross revenue (before discounts)
func (r *orderRepository) GetGrossRevenue(ctx context.Context) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("status IN ? AND payment_status = ?",
			[]entities.OrderStatus{entities.OrderStatusDelivered, entities.OrderStatusShipped},
			entities.PaymentStatusPaid).
		Select("COALESCE(SUM(subtotal + tax_amount + shipping_amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetProductRevenue gets product revenue (only subtotal)
func (r *orderRepository) GetProductRevenue(ctx context.Context) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("status IN ? AND payment_status = ?",
			[]entities.OrderStatus{entities.OrderStatusDelivered, entities.OrderStatusShipped},
			entities.PaymentStatusPaid).
		Select("COALESCE(SUM(subtotal), 0)").
		Scan(&total).Error
	return total, err
}

// GetTaxCollected gets total tax collected
func (r *orderRepository) GetTaxCollected(ctx context.Context) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("status IN ? AND payment_status = ?",
			[]entities.OrderStatus{entities.OrderStatusDelivered, entities.OrderStatusShipped},
			entities.PaymentStatusPaid).
		Select("COALESCE(SUM(tax_amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetShippingRevenue gets total shipping revenue
func (r *orderRepository) GetShippingRevenue(ctx context.Context) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("status IN ? AND payment_status = ?",
			[]entities.OrderStatus{entities.OrderStatusDelivered, entities.OrderStatusShipped},
			entities.PaymentStatusPaid).
		Select("COALESCE(SUM(shipping_amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetDiscountsGiven gets total discounts given
func (r *orderRepository) GetDiscountsGiven(ctx context.Context) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&entities.Order{}).
		Where("status IN ? AND payment_status = ?",
			[]entities.OrderStatus{entities.OrderStatusDelivered, entities.OrderStatusShipped},
			entities.PaymentStatusPaid).
		Select("COALESCE(SUM(discount_amount), 0)").
		Scan(&total).Error
	return total, err
}

type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *gorm.DB) repositories.PaymentRepository {
	return &paymentRepository{db: db}
}

// Create creates a new payment
func (r *paymentRepository) Create(ctx context.Context, payment *entities.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// GetByID retrieves a payment by ID
func (r *paymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrPaymentNotFound
		}
		return nil, err
	}
	return &payment, nil
}

// GetByOrderID retrieves the latest payment by order ID (for backward compatibility)
func (r *paymentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrPaymentNotFound
		}
		return nil, err
	}
	return &payment, nil
}

// GetAllByOrderID retrieves all payments for an order
func (r *paymentRepository) GetAllByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.Payment, error) {
	var payments []*entities.Payment
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}

// GetSuccessfulPaymentsByOrderID retrieves all successful payments for an order
func (r *paymentRepository) GetSuccessfulPaymentsByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.Payment, error) {
	var payments []*entities.Payment
	err := r.db.WithContext(ctx).
		Where("order_id = ? AND status = ?", orderID, entities.PaymentStatusPaid).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}

// GetByTransactionID retrieves a payment by transaction ID
func (r *paymentRepository) GetByTransactionID(ctx context.Context, transactionID string) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrPaymentNotFound
		}
		return nil, err
	}
	return &payment, nil
}

// GetByExternalID retrieves a payment by external ID (e.g., Stripe session ID)
func (r *paymentRepository) GetByExternalID(ctx context.Context, externalID string) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.WithContext(ctx).Where("external_id = ?", externalID).First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrPaymentNotFound
		}
		return nil, err
	}
	return &payment, nil
}

// Update updates an existing payment
func (r *paymentRepository) Update(ctx context.Context, payment *entities.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

// Delete deletes a payment by ID
func (r *paymentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.Payment{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrPaymentNotFound
	}
	return nil
}

// List retrieves payments with pagination
func (r *paymentRepository) List(ctx context.Context, limit, offset int) ([]*entities.Payment, error) {
	var payments []*entities.Payment
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}

// GetByStatus retrieves payments by status
func (r *paymentRepository) GetByStatus(ctx context.Context, status entities.PaymentStatus, limit, offset int) ([]*entities.Payment, error) {
	var payments []*entities.Payment
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}

// Count returns the total number of payments
func (r *paymentRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Payment{}).Count(&count).Error
	return count, err
}

// GetFailedPayments retrieves failed payments
func (r *paymentRepository) GetFailedPayments(ctx context.Context, limit, offset int) ([]*entities.Payment, error) {
	var payments []*entities.Payment
	err := r.db.WithContext(ctx).
		Where("status = ?", entities.PaymentStatusFailed).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}

// GetRefundablePayments retrieves payments that can be refunded
func (r *paymentRepository) GetRefundablePayments(ctx context.Context, limit, offset int) ([]*entities.Payment, error) {
	var payments []*entities.Payment
	err := r.db.WithContext(ctx).
		Where("status = ? AND refund_amount < amount", entities.PaymentStatusPaid).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}

// Refund-related methods
func (r *paymentRepository) CreateRefund(ctx context.Context, refund *entities.Refund) error {
	return r.db.WithContext(ctx).Create(refund).Error
}

func (r *paymentRepository) GetRefund(ctx context.Context, refundID uuid.UUID) (*entities.Refund, error) {
	var refund entities.Refund
	err := r.db.WithContext(ctx).Where("id = ?", refundID).First(&refund).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("refund not found")
		}
		return nil, err
	}
	return &refund, nil
}

func (r *paymentRepository) GetRefundsByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]*entities.Refund, error) {
	var refunds []*entities.Refund
	err := r.db.WithContext(ctx).Where("payment_id = ?", paymentID).Find(&refunds).Error
	return refunds, err
}

func (r *paymentRepository) UpdateRefund(ctx context.Context, refund *entities.Refund) error {
	return r.db.WithContext(ctx).Save(refund).Error
}

func (r *paymentRepository) ListRefunds(ctx context.Context, limit, offset int) ([]*entities.Refund, error) {
	var refunds []*entities.Refund
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&refunds).Error
	return refunds, err
}

func (r *paymentRepository) GetPendingRefunds(ctx context.Context, limit, offset int) ([]*entities.Refund, error) {
	var refunds []*entities.Refund
	err := r.db.WithContext(ctx).
		Where("status IN (?)", []entities.RefundStatus{
			entities.RefundStatusPending,
			entities.RefundStatusAwaitingApproval,
		}).
		Limit(limit).
		Offset(offset).
		Order("created_at ASC").
		Find(&refunds).Error
	return refunds, err
}
