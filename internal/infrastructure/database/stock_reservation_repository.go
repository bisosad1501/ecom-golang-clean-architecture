package database

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type stockReservationRepository struct {
	db *gorm.DB
}

// NewStockReservationRepository creates a new stock reservation repository
func NewStockReservationRepository(db *gorm.DB) repositories.StockReservationRepository {
	return &stockReservationRepository{db: db}
}

// Create creates a new stock reservation
func (r *stockReservationRepository) Create(ctx context.Context, reservation *entities.StockReservation) error {
	return r.db.WithContext(ctx).Create(reservation).Error
}

// GetByID gets a stock reservation by ID
func (r *stockReservationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.StockReservation, error) {
	var reservation entities.StockReservation
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Order").
		Preload("User").
		Where("id = ?", id).
		First(&reservation).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

// Update updates a stock reservation
func (r *stockReservationRepository) Update(ctx context.Context, reservation *entities.StockReservation) error {
	return r.db.WithContext(ctx).Save(reservation).Error
}

// Delete deletes a stock reservation
func (r *stockReservationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.StockReservation{}, id).Error
}

// GetByOrderID gets stock reservations by order ID
func (r *stockReservationRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.StockReservation, error) {
	var reservations []*entities.StockReservation
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("order_id = ?", orderID).
		Find(&reservations).Error
	return reservations, err
}

// GetByProductID gets stock reservations by product ID
func (r *stockReservationRepository) GetByProductID(ctx context.Context, productID uuid.UUID) ([]*entities.StockReservation, error) {
	var reservations []*entities.StockReservation
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Order").
		Where("product_id = ?", productID).
		Find(&reservations).Error
	return reservations, err
}

// GetByUserID gets stock reservations by user ID
func (r *stockReservationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.StockReservation, error) {
	var reservations []*entities.StockReservation
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Order").
		Where("user_id = ?", userID).
		Find(&reservations).Error
	return reservations, err
}

// GetActiveReservations gets active reservations for a product
func (r *stockReservationRepository) GetActiveReservations(ctx context.Context, productID uuid.UUID) ([]*entities.StockReservation, error) {
	var reservations []*entities.StockReservation
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND status = ? AND expires_at > ?",
			productID, entities.ReservationStatusActive, time.Now()).
		Find(&reservations).Error
	return reservations, err
}

// GetActiveReservationsByProduct gets active reservations for a product (alias for consistency)
func (r *stockReservationRepository) GetActiveReservationsByProduct(ctx context.Context, productID uuid.UUID) ([]*entities.StockReservation, error) {
	return r.GetActiveReservations(ctx, productID)
}

// GetExpiredReservations gets expired reservations
func (r *stockReservationRepository) GetExpiredReservations(ctx context.Context) ([]*entities.StockReservation, error) {
	var reservations []*entities.StockReservation
	err := r.db.WithContext(ctx).
		Where("status = ? AND expires_at <= ?", 
			entities.ReservationStatusActive, time.Now()).
		Find(&reservations).Error
	return reservations, err
}

// GetTotalReservedStock gets total reserved stock for a product
func (r *stockReservationRepository) GetTotalReservedStock(ctx context.Context, productID uuid.UUID) (int, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&entities.StockReservation{}).
		Select("COALESCE(SUM(quantity), 0)").
		Where("product_id = ? AND status = ? AND expires_at > ?", 
			productID, entities.ReservationStatusActive, time.Now()).
		Scan(&total).Error
	return int(total), err
}

// GetUserReservedStock gets user's reserved stock for a product
func (r *stockReservationRepository) GetUserReservedStock(ctx context.Context, userID, productID uuid.UUID) (int, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&entities.StockReservation{}).
		Select("COALESCE(SUM(quantity), 0)").
		Where("user_id = ? AND product_id = ? AND status = ? AND expires_at > ?", 
			userID, productID, entities.ReservationStatusActive, time.Now()).
		Scan(&total).Error
	return int(total), err
}

// CreateBatch creates multiple stock reservations
func (r *stockReservationRepository) CreateBatch(ctx context.Context, reservations []*entities.StockReservation) error {
	return r.db.WithContext(ctx).CreateInBatches(reservations, 100).Error
}

// ReleaseByOrderID releases all reservations for an order
func (r *stockReservationRepository) ReleaseByOrderID(ctx context.Context, orderID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entities.StockReservation{}).
		Where("order_id = ? AND status IN (?)", orderID, []string{
			string(entities.ReservationStatusActive),
			string(entities.ReservationStatusConfirmed),
		}).
		Updates(map[string]interface{}{
			"status":      entities.ReservationStatusReleased,
			"released_at": &now,
			"updated_at":  now,
		}).Error
}

// ReleaseExpiredReservations releases all expired reservations
func (r *stockReservationRepository) ReleaseExpiredReservations(ctx context.Context) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entities.StockReservation{}).
		Where("status = ? AND expires_at <= ?", entities.ReservationStatusActive, now).
		Updates(map[string]interface{}{
			"status":      entities.ReservationStatusExpired,
			"released_at": &now,
			"updated_at":  now,
		}).Error
}

// ConfirmByOrderID confirms all reservations for an order
func (r *stockReservationRepository) ConfirmByOrderID(ctx context.Context, orderID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entities.StockReservation{}).
		Where("order_id = ? AND status = ?", orderID, entities.ReservationStatusActive).
		Updates(map[string]interface{}{
			"status":       entities.ReservationStatusConfirmed,
			"confirmed_at": &now,
			"updated_at":   now,
		}).Error
}

// List lists stock reservations with filters
func (r *stockReservationRepository) List(ctx context.Context, filters repositories.StockReservationFilters) ([]*entities.StockReservation, error) {
	var reservations []*entities.StockReservation
	query := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Order").
		Preload("User")

	// Apply filters
	if filters.ProductID != nil {
		query = query.Where("product_id = ?", *filters.ProductID)
	}
	if filters.OrderID != nil {
		query = query.Where("order_id = ?", *filters.OrderID)
	}
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.ExpiresAt != nil {
		query = query.Where("expires_at <= ?", *filters.ExpiresAt)
	}
	if filters.CreatedAt != nil {
		query = query.Where("created_at >= ?", *filters.CreatedAt)
	}

	// Apply sorting
	if filters.SortBy != "" {
		order := filters.SortBy
		if filters.SortOrder != "" {
			order += " " + filters.SortOrder
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&reservations).Error
	return reservations, err
}

// Count counts stock reservations with filters
func (r *stockReservationRepository) Count(ctx context.Context, filters repositories.StockReservationFilters) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.StockReservation{})

	// Apply filters (same as List method)
	if filters.ProductID != nil {
		query = query.Where("product_id = ?", *filters.ProductID)
	}
	if filters.OrderID != nil {
		query = query.Where("order_id = ?", *filters.OrderID)
	}
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.ExpiresAt != nil {
		query = query.Where("expires_at <= ?", *filters.ExpiresAt)
	}
	if filters.CreatedAt != nil {
		query = query.Where("created_at >= ?", *filters.CreatedAt)
	}

	err := query.Count(&count).Error
	return count, err
}

// WithTransaction executes a function within a database transaction
func (r *stockReservationRepository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new context with the transaction
		txCtx := context.WithValue(ctx, "tx_repo", &stockReservationRepository{db: tx})
		return fn(txCtx)
	})
}
