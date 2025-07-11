package repositories

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// OrderEventRepository defines the interface for order event operations
type OrderEventRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, event *entities.OrderEvent) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.OrderEvent, error)
	Update(ctx context.Context, event *entities.OrderEvent) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Order-specific operations
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.OrderEvent, error)
	GetPublicByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.OrderEvent, error)
	GetByEventType(ctx context.Context, eventType entities.OrderEventType) ([]*entities.OrderEvent, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.OrderEvent, error)
	
	// Query operations
	List(ctx context.Context, filters OrderEventFilters) ([]*entities.OrderEvent, error)
	Count(ctx context.Context, filters OrderEventFilters) (int64, error)
	
	// Batch operations
	CreateBatch(ctx context.Context, events []*entities.OrderEvent) error
	DeleteByOrderID(ctx context.Context, orderID uuid.UUID) error
}

// OrderEventFilters represents filters for order event queries
type OrderEventFilters struct {
	OrderID    *uuid.UUID                  `json:"order_id"`
	EventType  *entities.OrderEventType    `json:"event_type"`
	UserID     *uuid.UUID                  `json:"user_id"`
	IsPublic   *bool                       `json:"is_public"`
	CreatedAt  *time.Time                  `json:"created_at"`
	DateFrom   *time.Time                  `json:"date_from"`
	DateTo     *time.Time                  `json:"date_to"`
	Limit      int                         `json:"limit"`
	Offset     int                         `json:"offset"`
	SortBy     string                      `json:"sort_by"`
	SortOrder  string                      `json:"sort_order"`
}
