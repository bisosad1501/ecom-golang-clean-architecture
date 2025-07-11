package repositories

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// StockReservationRepository defines the interface for stock reservation operations
type StockReservationRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, reservation *entities.StockReservation) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.StockReservation, error)
	Update(ctx context.Context, reservation *entities.StockReservation) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Reservation-specific operations
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.StockReservation, error)
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]*entities.StockReservation, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.StockReservation, error)
	GetActiveReservations(ctx context.Context, productID uuid.UUID) ([]*entities.StockReservation, error)
	GetExpiredReservations(ctx context.Context) ([]*entities.StockReservation, error)
	
	// Stock calculation operations
	GetTotalReservedStock(ctx context.Context, productID uuid.UUID) (int, error)
	GetUserReservedStock(ctx context.Context, userID, productID uuid.UUID) (int, error)
	
	// Batch operations
	CreateBatch(ctx context.Context, reservations []*entities.StockReservation) error
	ReleaseByOrderID(ctx context.Context, orderID uuid.UUID) error
	ReleaseExpiredReservations(ctx context.Context) error
	ConfirmByOrderID(ctx context.Context, orderID uuid.UUID) error
	
	// Query operations
	List(ctx context.Context, filters StockReservationFilters) ([]*entities.StockReservation, error)
	Count(ctx context.Context, filters StockReservationFilters) (int64, error)
}

// StockReservationFilters represents filters for stock reservation queries
type StockReservationFilters struct {
	ProductID  *uuid.UUID                        `json:"product_id"`
	OrderID    *uuid.UUID                        `json:"order_id"`
	UserID     *uuid.UUID                        `json:"user_id"`
	Type       *entities.StockReservationType    `json:"type"`
	Status     *entities.StockReservationStatus  `json:"status"`
	ExpiresAt  *time.Time                        `json:"expires_at"`
	CreatedAt  *time.Time                        `json:"created_at"`
	Limit      int                               `json:"limit"`
	Offset     int                               `json:"offset"`
	SortBy     string                            `json:"sort_by"`
	SortOrder  string                            `json:"sort_order"`
}
