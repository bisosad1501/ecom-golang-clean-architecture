package database

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type shippingRepository struct {
	db *gorm.DB
}

// NewShippingRepository creates a new shipping repository
func NewShippingRepository(db *gorm.DB) repositories.ShippingRepository {
	return &shippingRepository{db: db}
}

// CreateShipment creates a new shipment
func (r *shippingRepository) CreateShipment(ctx context.Context, shipment *entities.Shipment) error {
	return r.db.WithContext(ctx).Create(shipment).Error
}

// GetShipmentByID gets a shipment by ID
func (r *shippingRepository) GetShipmentByID(ctx context.Context, id uuid.UUID) (*entities.Shipment, error) {
	var shipment entities.Shipment
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("ShippingMethod").
		First(&shipment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

// GetShipmentByTrackingNumber gets a shipment by tracking number
func (r *shippingRepository) GetShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*entities.Shipment, error) {
	var shipment entities.Shipment
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("ShippingMethod").
		First(&shipment, "tracking_number = ?", trackingNumber).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

// GetShipmentsByOrder gets shipments for an order
func (r *shippingRepository) GetShipmentsByOrder(ctx context.Context, orderID uuid.UUID) ([]*entities.Shipment, error) {
	var shipments []*entities.Shipment
	err := r.db.WithContext(ctx).
		Preload("ShippingMethod").
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&shipments).Error
	return shipments, err
}

// UpdateShipment updates a shipment
func (r *shippingRepository) UpdateShipment(ctx context.Context, shipment *entities.Shipment) error {
	shipment.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(shipment).Error
}

// UpdateShipmentStatus updates shipment status
func (r *shippingRepository) UpdateShipmentStatus(ctx context.Context, shipmentID uuid.UUID, status entities.ShipmentStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	// Set delivered date if status is delivered
	if status == entities.ShipmentStatusDelivered {
		updates["delivered_at"] = time.Now()
	}

	return r.db.WithContext(ctx).
		Model(&entities.Shipment{}).
		Where("id = ?", shipmentID).
		Updates(updates).Error
}

// CreateShippingMethod creates a new shipping method
func (r *shippingRepository) CreateShippingMethod(ctx context.Context, method *entities.ShippingMethod) error {
	return r.db.WithContext(ctx).Create(method).Error
}

// GetShippingMethodByID gets a shipping method by ID
func (r *shippingRepository) GetShippingMethodByID(ctx context.Context, id uuid.UUID) (*entities.ShippingMethod, error) {
	var method entities.ShippingMethod
	err := r.db.WithContext(ctx).First(&method, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &method, nil
}

// GetActiveShippingMethods gets all active shipping methods
func (r *shippingRepository) GetActiveShippingMethods(ctx context.Context) ([]*entities.ShippingMethod, error) {
	var methods []*entities.ShippingMethod
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("display_order ASC, name ASC").
		Find(&methods).Error
	return methods, err
}

// UpdateShippingMethod updates a shipping method
func (r *shippingRepository) UpdateShippingMethod(ctx context.Context, method *entities.ShippingMethod) error {
	method.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(method).Error
}

// DeleteShippingMethod deletes a shipping method
func (r *shippingRepository) DeleteShippingMethod(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.ShippingMethod{}, "id = ?", id).Error
}

// CreateTrackingEvent creates a tracking event
func (r *shippingRepository) CreateTrackingEvent(ctx context.Context, event *entities.ShipmentTracking) error {
	return r.db.WithContext(ctx).Create(event).Error
}

// GetTrackingEvents gets tracking events for a shipment
func (r *shippingRepository) GetTrackingEvents(ctx context.Context, shipmentID uuid.UUID) ([]*entities.ShipmentTracking, error) {
	var events []*entities.ShipmentTracking
	err := r.db.WithContext(ctx).
		Where("shipment_id = ?", shipmentID).
		Order("created_at ASC").
		Find(&events).Error
	return events, err
}

// ListShipments lists shipments with filters
func (r *shippingRepository) ListShipments(ctx context.Context, filters repositories.ShipmentFilters) ([]*entities.Shipment, error) {
	var shipments []*entities.Shipment
	query := r.db.WithContext(ctx).
		Preload("Order").
		Preload("ShippingAddress")

	if filters.OrderID != nil {
		query = query.Where("order_id = ?", *filters.OrderID)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.Carrier != "" {
		query = query.Where("carrier = ?", filters.Carrier)
	}

	if filters.TrackingNumber != "" {
		query = query.Where("tracking_number LIKE ?", "%"+filters.TrackingNumber+"%")
	}

	if filters.CreatedAfter != nil {
		query = query.Where("created_at >= ?", *filters.CreatedAfter)
	}

	if filters.CreatedBefore != nil {
		query = query.Where("created_at <= ?", *filters.CreatedBefore)
	}

	// Apply sorting
	switch filters.SortBy {
	case "created_at":
		if filters.SortOrder == "desc" {
			query = query.Order("created_at DESC")
		} else {
			query = query.Order("created_at ASC")
		}
	case "status":
		if filters.SortOrder == "desc" {
			query = query.Order("status DESC")
		} else {
			query = query.Order("status ASC")
		}
	case "carrier":
		if filters.SortOrder == "desc" {
			query = query.Order("carrier DESC")
		} else {
			query = query.Order("carrier ASC")
		}
	default:
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&shipments).Error
	return shipments, err
}

// CountShipments counts shipments with filters
func (r *shippingRepository) CountShipments(ctx context.Context, filters repositories.ShipmentFilters) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Shipment{})

	if filters.OrderID != nil {
		query = query.Where("order_id = ?", *filters.OrderID)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.Carrier != "" {
		query = query.Where("carrier = ?", filters.Carrier)
	}

	if filters.TrackingNumber != "" {
		query = query.Where("tracking_number LIKE ?", "%"+filters.TrackingNumber+"%")
	}

	if filters.CreatedAfter != nil {
		query = query.Where("created_at >= ?", *filters.CreatedAfter)
	}

	if filters.CreatedBefore != nil {
		query = query.Where("created_at <= ?", *filters.CreatedBefore)
	}

	err := query.Count(&count).Error
	return count, err
}

// GetShipmentsByStatus gets shipments by status
func (r *shippingRepository) GetShipmentsByStatus(ctx context.Context, status entities.ShipmentStatus, limit, offset int) ([]*entities.Shipment, error) {
	var shipments []*entities.Shipment
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("ShippingAddress").
		Where("status = ?", status).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&shipments).Error
	return shipments, err
}

// GetPendingShipments gets shipments that are pending
func (r *shippingRepository) GetPendingShipments(ctx context.Context, limit, offset int) ([]*entities.Shipment, error) {
	var shipments []*entities.Shipment
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("ShippingAddress").
		Where("status IN (?)", []entities.ShipmentStatus{
			entities.ShipmentStatusPending,
			entities.ShipmentStatusProcessing,
			entities.ShipmentStatusShipped,
		}).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&shipments).Error
	return shipments, err
}

// GetOverdueShipments gets shipments that are overdue
func (r *shippingRepository) GetOverdueShipments(ctx context.Context, limit, offset int) ([]*entities.Shipment, error) {
	var shipments []*entities.Shipment
	overdueDate := time.Now().AddDate(0, 0, -7) // 7 days ago
	
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("ShippingAddress").
		Where("status IN (?) AND created_at < ?", 
			[]entities.ShipmentStatus{
				entities.ShipmentStatusPending,
				entities.ShipmentStatusProcessing,
				entities.ShipmentStatusShipped,
			}, overdueDate).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&shipments).Error
	return shipments, err
}

// DeleteShipment deletes a shipment
func (r *shippingRepository) DeleteShipment(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Shipment{}, "id = ?", id).Error
}

// CreateReturn creates a return request
func (r *shippingRepository) CreateReturn(ctx context.Context, returnRequest *entities.Return) error {
	// Set return-specific properties
	returnRequest.ID = uuid.New()
	returnRequest.Status = entities.ReturnStatusRequested
	returnRequest.CreatedAt = time.Now()
	returnRequest.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).Create(returnRequest).Error
}

// GetReturnByID gets a return by ID
func (r *shippingRepository) GetReturnByID(ctx context.Context, id uuid.UUID) (*entities.Return, error) {
	var returnRequest entities.Return
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("User").
		First(&returnRequest, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &returnRequest, nil
}

// GetShippingMethods gets available shipping methods for location and weight
func (r *shippingRepository) GetShippingMethods(ctx context.Context, locationID *uuid.UUID, weight *float64) ([]*entities.ShippingMethod, error) {
	var methods []*entities.ShippingMethod
	query := r.db.WithContext(ctx).Where("is_active = ?", true)

	// Filter by weight limits if weight is provided
	if weight != nil {
		query = query.Where("max_weight IS NULL OR max_weight >= ?", *weight)
	}

	// Additional location-based filtering could be added here if needed

	err := query.Order("name ASC").Find(&methods).Error
	return methods, err
}

// UpdateReturn updates a return request
func (r *shippingRepository) UpdateReturn(ctx context.Context, returnRequest *entities.Return) error {
	returnRequest.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(returnRequest).Error
}
