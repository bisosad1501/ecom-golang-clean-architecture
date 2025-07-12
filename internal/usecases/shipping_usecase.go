package usecases

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
)

// ShippingUseCase interface
type ShippingUseCase interface {
	// Shipping Methods
	GetShippingMethods(ctx context.Context, req GetShippingMethodsRequest) ([]*ShippingMethodResponse, error)
	CalculateShippingCost(ctx context.Context, req CalculateShippingRequest) (*ShippingCostResponse, error)

	// Shipments
	CreateShipment(ctx context.Context, req CreateShipmentRequest) (*ShipmentResponse, error)
	GetShipment(ctx context.Context, shipmentID uuid.UUID) (*ShipmentResponse, error)
	UpdateShipmentStatus(ctx context.Context, shipmentID uuid.UUID, status entities.ShipmentStatus) (*ShipmentResponse, error)
	TrackShipment(ctx context.Context, trackingNumber string) (*ShipmentTrackingResponse, error)

	// Returns
	CreateReturn(ctx context.Context, req CreateReturnRequest) (*ReturnResponse, error)
	GetReturn(ctx context.Context, returnID uuid.UUID) (*ReturnResponse, error)
	ProcessReturn(ctx context.Context, returnID uuid.UUID, status entities.ReturnStatus) (*ReturnResponse, error)
}

type shippingUseCase struct {
	shippingRepo repositories.ShippingRepository
	orderRepo    repositories.OrderRepository
}

// NewShippingUseCase creates a new shipping use case
func NewShippingUseCase(
	shippingRepo repositories.ShippingRepository,
	orderRepo repositories.OrderRepository,
) ShippingUseCase {
	return &shippingUseCase{
		shippingRepo: shippingRepo,
		orderRepo:    orderRepo,
	}
}

// Request/Response types
type GetShippingMethodsRequest struct {
	ZoneID      *uuid.UUID `json:"zone_id"`
	Weight      *float64   `json:"weight"`
	Destination string     `json:"destination"`
}

type CalculateShippingRequest struct {
	OrderID     uuid.UUID `json:"order_id" validate:"required"`
	MethodID    uuid.UUID `json:"method_id" validate:"required"`
	Destination string    `json:"destination" validate:"required"`
}

type CreateShipmentRequest struct {
	OrderID           uuid.UUID  `json:"order_id" validate:"required"`
	ShippingMethod    uuid.UUID  `json:"shipping_method_id" validate:"required"`
	TrackingNumber    string     `json:"tracking_number"`
	Carrier           string     `json:"carrier" validate:"required"`
	EstimatedDelivery *time.Time `json:"estimated_delivery"`
}

type CreateReturnRequest struct {
	OrderID     uuid.UUID             `json:"order_id" validate:"required"`
	Items       []ReturnItemRequest   `json:"items" validate:"required,dive"`
	Reason      entities.ReturnReason `json:"reason" validate:"required"`
	Description string                `json:"description"`
}

type ReturnItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
}

// Response types
type ShippingMethodResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cost        float64   `json:"cost"`
	MinWeight   float64   `json:"min_weight"`
	MaxWeight   float64   `json:"max_weight"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type ShippingCostResponse struct {
	MethodID      uuid.UUID `json:"method_id"`
	MethodName    string    `json:"method_name"`
	Cost          float64   `json:"cost"`
	EstimatedDays int       `json:"estimated_days"`
}

type ShipmentResponse struct {
	ID                uuid.UUID               `json:"id"`
	OrderID           uuid.UUID               `json:"order_id"`
	TrackingNumber    string                  `json:"tracking_number"`
	Carrier           string                  `json:"carrier"`
	Status            entities.ShipmentStatus `json:"status"`
	ShippedAt         *time.Time              `json:"shipped_at"`
	ActualDelivery    *time.Time              `json:"actual_delivery"`
	EstimatedDelivery *time.Time              `json:"estimated_delivery"`
	TrackingEvents    []ShipmentTrackingEvent `json:"tracking_events"`
	CreatedAt         time.Time               `json:"created_at"`
	UpdatedAt         time.Time               `json:"updated_at"`
}

type ShipmentTrackingEvent struct {
	ID          uuid.UUID `json:"id"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Timestamp   time.Time `json:"timestamp"`
}

type ShipmentTrackingResponse struct {
	TrackingNumber    string                  `json:"tracking_number"`
	Status            entities.ShipmentStatus `json:"status"`
	Events            []ShipmentTrackingEvent `json:"events"`
	EstimatedDelivery *time.Time              `json:"estimated_delivery"`
}

type ReturnResponse struct {
	ID           uuid.UUID             `json:"id"`
	OrderID      uuid.UUID             `json:"order_id"`
	Status       entities.ReturnStatus `json:"status"`
	Reason       entities.ReturnReason `json:"reason"`
	Description  string                `json:"description"`
	Items        []ReturnItemResponse  `json:"items"`
	RefundAmount float64               `json:"refund_amount"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

type ReturnItemResponse struct {
	ID           uuid.UUID `json:"id"`
	ProductID    uuid.UUID `json:"product_id"`
	ProductName  string    `json:"product_name"`
	Quantity     int       `json:"quantity"`
	RefundAmount float64   `json:"refund_amount"`
}

// GetShippingMethods gets available shipping methods
func (uc *shippingUseCase) GetShippingMethods(ctx context.Context, req GetShippingMethodsRequest) ([]*ShippingMethodResponse, error) {
	methods, err := uc.shippingRepo.GetShippingMethods(ctx, req.ZoneID, req.Weight)
	if err != nil {
		return nil, err
	}

	responses := make([]*ShippingMethodResponse, len(methods))
	for i, method := range methods {
		responses[i] = &ShippingMethodResponse{
			ID:          method.ID,
			Name:        method.Name,
			Description: method.Description,
			Cost:        method.BaseCost,
			MinWeight:   0, // Would come from shipping rates
			MaxWeight:   method.MaxWeight,
			IsActive:    method.IsActive,
			CreatedAt:   method.CreatedAt,
		}
	}

	return responses, nil
}

// CalculateShippingCost calculates shipping cost for an order
func (uc *shippingUseCase) CalculateShippingCost(ctx context.Context, req CalculateShippingRequest) (*ShippingCostResponse, error) {
	// Get order details
	order, err := uc.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, entities.ErrOrderNotFound
	}

	// Get shipping method
	method, err := uc.shippingRepo.GetShippingMethodByID(ctx, req.MethodID)
	if err != nil {
		return nil, entities.ErrShippingMethodNotFound
	}

	// Calculate total weight (simplified)
	totalWeight := 0.0
	for _, item := range order.Items {
		if item.Product.Weight != nil {
			totalWeight += *item.Product.Weight * float64(item.Quantity)
		}
	}

	// Calculate cost based on weight and distance (simplified)
	cost := method.BaseCost
	if totalWeight > 0 {
		cost += totalWeight * method.CostPerKg
	}

	return &ShippingCostResponse{
		MethodID:      method.ID,
		MethodName:    method.Name,
		Cost:          cost,
		EstimatedDays: method.MaxDeliveryDays,
	}, nil
}

// CreateShipment creates a new shipment
func (uc *shippingUseCase) CreateShipment(ctx context.Context, req CreateShipmentRequest) (*ShipmentResponse, error) {
	// Verify order exists
	order, err := uc.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, entities.ErrOrderNotFound
	}

	// Create shipment
	shipment := &entities.Shipment{
		ID:                uuid.New(),
		OrderID:           req.OrderID,
		TrackingNumber:    req.TrackingNumber,
		Carrier:           req.Carrier,
		Status:            entities.ShipmentStatusPending,
		EstimatedDelivery: req.EstimatedDelivery,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := uc.shippingRepo.CreateShipment(ctx, shipment); err != nil {
		return nil, err
	}

	// Update order status to shipped
	if err := uc.orderRepo.UpdateStatus(ctx, order.ID, entities.OrderStatusShipped); err != nil {
		// handle hoặc log lỗi nếu cần
	}

	return uc.toShipmentResponse(shipment), nil
}

// GetShipment gets shipment by ID
func (uc *shippingUseCase) GetShipment(ctx context.Context, shipmentID uuid.UUID) (*ShipmentResponse, error) {
	shipment, err := uc.shippingRepo.GetShipmentByID(ctx, shipmentID)
	if err != nil {
		return nil, entities.ErrShipmentNotFound
	}

	return uc.toShipmentResponse(shipment), nil
}

// UpdateShipmentStatus updates shipment status
func (uc *shippingUseCase) UpdateShipmentStatus(ctx context.Context, shipmentID uuid.UUID, status entities.ShipmentStatus) (*ShipmentResponse, error) {
	shipment, err := uc.shippingRepo.GetShipmentByID(ctx, shipmentID)
	if err != nil {
		return nil, entities.ErrShipmentNotFound
	}

	// Update status
	shipment.Status = status
	shipment.UpdatedAt = time.Now()

	if status == entities.ShipmentStatusDelivered {
		now := time.Now()
		shipment.ActualDelivery = &now
		// Update order status
		if err := uc.orderRepo.UpdateStatus(ctx, shipment.OrderID, entities.OrderStatusDelivered); err != nil {
			// handle hoặc log lỗi nếu cần
		}
	}

	if err := uc.shippingRepo.UpdateShipment(ctx, shipment); err != nil {
		return nil, err
	}

	return uc.toShipmentResponse(shipment), nil
}

// TrackShipment tracks shipment by tracking number
func (uc *shippingUseCase) TrackShipment(ctx context.Context, trackingNumber string) (*ShipmentTrackingResponse, error) {
	shipment, err := uc.shippingRepo.GetShipmentByTrackingNumber(ctx, trackingNumber)
	if err != nil {
		return nil, entities.ErrShipmentNotFound
	}

	// Get tracking events
	events, err := uc.shippingRepo.GetTrackingEvents(ctx, shipment.ID)
	if err != nil {
		return nil, err
	}

	trackingEvents := make([]ShipmentTrackingEvent, len(events))
	for i, event := range events {
		trackingEvents[i] = ShipmentTrackingEvent{
			ID:          event.ID,
			Status:      string(event.Status),
			Description: event.Description,
			Location:    event.Location,
			Timestamp:   event.EventTime,
		}
	}

	return &ShipmentTrackingResponse{
		TrackingNumber:    trackingNumber,
		Status:            shipment.Status,
		Events:            trackingEvents,
		EstimatedDelivery: shipment.EstimatedDelivery,
	}, nil
}

// CreateReturn creates a return request
func (uc *shippingUseCase) CreateReturn(ctx context.Context, req CreateReturnRequest) (*ReturnResponse, error) {
	// Verify order exists and is eligible for return
	order, err := uc.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, entities.ErrOrderNotFound
	}

	if !order.CanBeRefunded() {
		return nil, entities.ErrOrderCannotBeReturned
	}

	// Create return
	returnEntity := &entities.Return{
		ID:          uuid.New(),
		OrderID:     req.OrderID,
		Status:      entities.ReturnStatusRequested,
		Reason:      req.Reason,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create return items
	for _, item := range req.Items {
		returnItem := entities.ReturnItem{
			ID:        uuid.New(),
			ReturnID:  returnEntity.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
		returnEntity.Items = append(returnEntity.Items, returnItem)
	}

	if err := uc.shippingRepo.CreateReturn(ctx, returnEntity); err != nil {
		return nil, err
	}

	return uc.toReturnResponse(returnEntity), nil
}

// GetReturn gets return by ID
func (uc *shippingUseCase) GetReturn(ctx context.Context, returnID uuid.UUID) (*ReturnResponse, error) {
	returnEntity, err := uc.shippingRepo.GetReturnByID(ctx, returnID)
	if err != nil {
		return nil, entities.ErrReturnNotFound
	}

	return uc.toReturnResponse(returnEntity), nil
}

// ProcessReturn processes a return request
func (uc *shippingUseCase) ProcessReturn(ctx context.Context, returnID uuid.UUID, status entities.ReturnStatus) (*ReturnResponse, error) {
	returnEntity, err := uc.shippingRepo.GetReturnByID(ctx, returnID)
	if err != nil {
		return nil, entities.ErrReturnNotFound
	}

	returnEntity.Status = status
	returnEntity.UpdatedAt = time.Now()

	if err := uc.shippingRepo.UpdateReturn(ctx, returnEntity); err != nil {
		return nil, err
	}

	return uc.toReturnResponse(returnEntity), nil
}

// Helper methods
func (uc *shippingUseCase) toShipmentResponse(shipment *entities.Shipment) *ShipmentResponse {
	return &ShipmentResponse{
		ID:                shipment.ID,
		OrderID:           shipment.OrderID,
		TrackingNumber:    shipment.TrackingNumber,
		Carrier:           shipment.Carrier,
		Status:            shipment.Status,
		ShippedAt:         shipment.ShippedAt,
		ActualDelivery:    shipment.ActualDelivery,
		EstimatedDelivery: shipment.EstimatedDelivery,
		CreatedAt:         shipment.CreatedAt,
		UpdatedAt:         shipment.UpdatedAt,
	}
}

func (uc *shippingUseCase) toReturnResponse(returnEntity *entities.Return) *ReturnResponse {
	items := make([]ReturnItemResponse, len(returnEntity.Items))
	for i, item := range returnEntity.Items {
		items[i] = ReturnItemResponse{
			ID:           item.ID,
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			RefundAmount: item.RefundAmount,
		}
	}

	return &ReturnResponse{
		ID:           returnEntity.ID,
		OrderID:      returnEntity.OrderID,
		Status:       returnEntity.Status,
		Reason:       returnEntity.Reason,
		Description:  returnEntity.Description,
		Items:        items,
		RefundAmount: returnEntity.RefundAmount,
		CreatedAt:    returnEntity.CreatedAt,
		UpdatedAt:    returnEntity.UpdatedAt,
	}
}
