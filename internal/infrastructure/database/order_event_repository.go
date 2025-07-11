package database

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderEventRepository struct {
	db *gorm.DB
}

// NewOrderEventRepository creates a new order event repository
func NewOrderEventRepository(db *gorm.DB) repositories.OrderEventRepository {
	return &orderEventRepository{db: db}
}

// Create creates a new order event
func (r *orderEventRepository) Create(ctx context.Context, event *entities.OrderEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

// GetByID gets an order event by ID
func (r *orderEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.OrderEvent, error) {
	var event entities.OrderEvent
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", id).
		First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// Update updates an order event
func (r *orderEventRepository) Update(ctx context.Context, event *entities.OrderEvent) error {
	return r.db.WithContext(ctx).Save(event).Error
}

// Delete deletes an order event
func (r *orderEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.OrderEvent{}, id).Error
}

// GetByOrderID gets order events by order ID
func (r *orderEventRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.OrderEvent, error) {
	var events []*entities.OrderEvent
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Find(&events).Error
	return events, err
}

// GetPublicByOrderID gets public order events by order ID
func (r *orderEventRepository) GetPublicByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.OrderEvent, error) {
	var events []*entities.OrderEvent
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("order_id = ? AND is_public = ?", orderID, true).
		Order("created_at ASC").
		Find(&events).Error
	return events, err
}

// GetByEventType gets order events by event type
func (r *orderEventRepository) GetByEventType(ctx context.Context, eventType entities.OrderEventType) ([]*entities.OrderEvent, error) {
	var events []*entities.OrderEvent
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("event_type = ?", eventType).
		Order("created_at DESC").
		Find(&events).Error
	return events, err
}

// GetByUserID gets order events by user ID
func (r *orderEventRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.OrderEvent, error) {
	var events []*entities.OrderEvent
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&events).Error
	return events, err
}

// List lists order events with filters
func (r *orderEventRepository) List(ctx context.Context, filters repositories.OrderEventFilters) ([]*entities.OrderEvent, error) {
	var events []*entities.OrderEvent
	query := r.db.WithContext(ctx).
		Preload("User")

	// Apply filters
	if filters.OrderID != nil {
		query = query.Where("order_id = ?", *filters.OrderID)
	}
	if filters.EventType != nil {
		query = query.Where("event_type = ?", *filters.EventType)
	}
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if filters.IsPublic != nil {
		query = query.Where("is_public = ?", *filters.IsPublic)
	}
	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
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

	err := query.Find(&events).Error
	return events, err
}

// Count counts order events with filters
func (r *orderEventRepository) Count(ctx context.Context, filters repositories.OrderEventFilters) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.OrderEvent{})

	// Apply filters (same as List method)
	if filters.OrderID != nil {
		query = query.Where("order_id = ?", *filters.OrderID)
	}
	if filters.EventType != nil {
		query = query.Where("event_type = ?", *filters.EventType)
	}
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if filters.IsPublic != nil {
		query = query.Where("is_public = ?", *filters.IsPublic)
	}
	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	err := query.Count(&count).Error
	return count, err
}

// CreateBatch creates multiple order events
func (r *orderEventRepository) CreateBatch(ctx context.Context, events []*entities.OrderEvent) error {
	return r.db.WithContext(ctx).CreateInBatches(events, 100).Error
}

// DeleteByOrderID deletes all events for an order
func (r *orderEventRepository) DeleteByOrderID(ctx context.Context, orderID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Delete(&entities.OrderEvent{}).Error
}
