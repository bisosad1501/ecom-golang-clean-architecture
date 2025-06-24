package repositories

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// ShippingRepository interface
type ShippingRepository interface {
	// Shipping Methods
	GetShippingMethods(ctx context.Context, zoneID *uuid.UUID, weight *float64) ([]*entities.ShippingMethod, error)
	GetShippingMethodByID(ctx context.Context, id uuid.UUID) (*entities.ShippingMethod, error)
	
	// Shipments
	CreateShipment(ctx context.Context, shipment *entities.Shipment) error
	GetShipmentByID(ctx context.Context, id uuid.UUID) (*entities.Shipment, error)
	GetShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*entities.Shipment, error)
	UpdateShipment(ctx context.Context, shipment *entities.Shipment) error
	GetTrackingEvents(ctx context.Context, shipmentID uuid.UUID) ([]*entities.ShipmentTracking, error)
	
	// Returns
	CreateReturn(ctx context.Context, returnEntity *entities.Return) error
	GetReturnByID(ctx context.Context, id uuid.UUID) (*entities.Return, error)
	UpdateReturn(ctx context.Context, returnEntity *entities.Return) error
}
