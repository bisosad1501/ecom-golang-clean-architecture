package usecases

import (
	"context"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/domain/services"

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

	// Distance-based methods
	CalculateDistanceBasedShipping(ctx context.Context, req DistanceBasedShippingRequest) (*DistanceBasedShippingResponse, error)
	GetShippingZones(ctx context.Context) ([]services.ShippingZoneInfo, error)

	// Address validation
	ValidateShippingAddress(ctx context.Context, req ValidateShippingAddressRequest) (*ValidateShippingAddressResponse, error)
}

type shippingUseCase struct {
	shippingRepo         repositories.ShippingRepository
	orderRepo            repositories.OrderRepository
	distanceService      services.DistanceService
	compatibilityService services.ShippingCompatibilityService
}

// NewShippingUseCase creates a new shipping use case
func NewShippingUseCase(
	shippingRepo repositories.ShippingRepository,
	orderRepo repositories.OrderRepository,
	distanceService services.DistanceService,
	compatibilityService services.ShippingCompatibilityService,
) ShippingUseCase {
	return &shippingUseCase{
		shippingRepo:         shippingRepo,
		orderRepo:            orderRepo,
		distanceService:      distanceService,
		compatibilityService: compatibilityService,
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

type DistanceBasedShippingRequest struct {
	FromLatitude  *float64 `json:"from_latitude"`
	FromLongitude *float64 `json:"from_longitude"`
	FromAddress   string   `json:"from_address"`
	ToLatitude    *float64 `json:"to_latitude"`
	ToLongitude   *float64 `json:"to_longitude"`
	ToAddress     string   `json:"to_address"`
	Destination   string   `json:"destination"` // Alternative field name for compatibility
	Weight        float64  `json:"weight" validate:"required,gt=0"`
	OrderValue    float64  `json:"order_value" validate:"required,gt=0"`
	MethodID      string   `json:"method_id"`
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

type DistanceBasedShippingResponse struct {
	Distance      float64                     `json:"distance_km"`
	Zone          string                      `json:"shipping_zone"`
	IsShippable   bool                        `json:"is_shippable"`
	Options       []DistanceShippingOption    `json:"shipping_options"`
	Recommendations []string                  `json:"recommendations"`
}

type DistanceShippingOption struct {
	MethodID      string  `json:"method_id"`
	MethodName    string  `json:"method_name"`
	Cost          float64 `json:"cost"`
	EstimatedDays int     `json:"estimated_days"`
	IsAvailable   bool    `json:"is_available"`
	Reason        string  `json:"reason,omitempty"`
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

// CalculateDistanceBasedShipping calculates shipping options based on distance
func (uc *shippingUseCase) CalculateDistanceBasedShipping(ctx context.Context, req DistanceBasedShippingRequest) (*DistanceBasedShippingResponse, error) {
	var distance float64
	var err error

	// Set default from address if not provided
	fromAddress := req.FromAddress
	if fromAddress == "" {
		fromAddress = "New York, NY, USA" // Default warehouse location
	}

	// Determine destination address
	toAddress := req.ToAddress
	if toAddress == "" {
		toAddress = req.Destination
	}

	// Calculate distance based on provided data
	if req.FromLatitude != nil && req.FromLongitude != nil && req.ToLatitude != nil && req.ToLongitude != nil {
		distance, err = uc.distanceService.CalculateDistance(ctx, *req.FromLatitude, *req.FromLongitude, *req.ToLatitude, *req.ToLongitude)
	} else if fromAddress != "" && toAddress != "" {
		distance, err = uc.distanceService.CalculateDistanceByAddress(ctx, fromAddress, toAddress)
	} else {
		return nil, fmt.Errorf("insufficient location data provided: need either coordinates or addresses")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to calculate distance: %w", err)
	}

	// Get shipping zone
	zone, err := uc.distanceService.GetShippingZoneByDistance(ctx, distance)
	if err != nil {
		return &DistanceBasedShippingResponse{
			Distance:    distance,
			Zone:        "unavailable",
			IsShippable: false,
			Options:     []DistanceShippingOption{},
			Recommendations: []string{"Shipping not available for this distance"},
		}, nil
	}

	// Get available shipping methods
	methods, err := uc.shippingRepo.GetShippingMethods(ctx, nil, &req.Weight)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipping methods: %w", err)
	}

	// Calculate shipping options
	var options []DistanceShippingOption
	var recommendations []string

	for _, method := range methods {
		// Validate if method supports this distance
		isValid, err := uc.distanceService.ValidateShippingDistance(ctx, distance, method.Name)
		if err != nil {
			continue
		}

		// Calculate cost using the method's CalculateCost function
		cost := method.CalculateCost(req.Weight, distance, req.OrderValue)

		option := DistanceShippingOption{
			MethodID:      method.ID.String(),
			MethodName:    method.Name,
			Cost:          cost,
			EstimatedDays: method.MaxDeliveryDays,
			IsAvailable:   isValid,
		}

		if !isValid {
			option.Reason = fmt.Sprintf("Distance %.1f km exceeds maximum for %s", distance, method.Name)
		}

		options = append(options, option)
	}

	// Add recommendations
	if distance > 500 {
		recommendations = append(recommendations, "Consider splitting large orders for faster delivery")
	}
	if req.OrderValue > 100 {
		recommendations = append(recommendations, "You may be eligible for free shipping on some methods")
	}

	return &DistanceBasedShippingResponse{
		Distance:        distance,
		Zone:           zone,
		IsShippable:    len(options) > 0,
		Options:        options,
		Recommendations: recommendations,
	}, nil
}

// GetShippingZones returns available shipping zones
func (uc *shippingUseCase) GetShippingZones(ctx context.Context) ([]services.ShippingZoneInfo, error) {
	return uc.distanceService.GetShippingZones(), nil
}

// SimpleAddress represents a simplified address for validation
type SimpleAddress struct {
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	Address1   string `json:"address1" validate:"required"`
	Address2   string `json:"address2"`
	City       string `json:"city" validate:"required"`
	State      string `json:"state" validate:"required"`
	ZipCode    string `json:"zip_code" validate:"required"`
	Country    string `json:"country" validate:"required"`
	Phone      string `json:"phone"`
}

// Validate validates the simple address
func (sa *SimpleAddress) Validate() error {
	if sa.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if sa.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if sa.Address1 == "" {
		return fmt.Errorf("address line 1 is required")
	}
	if sa.City == "" {
		return fmt.Errorf("city is required")
	}
	if sa.State == "" {
		return fmt.Errorf("state is required")
	}
	if sa.ZipCode == "" {
		return fmt.Errorf("zip code is required")
	}
	if sa.Country == "" {
		return fmt.Errorf("country is required")
	}
	return nil
}

// ToAddress converts SimpleAddress to entities.Address
func (sa *SimpleAddress) ToAddress() *entities.Address {
	return &entities.Address{
		FirstName: sa.FirstName,
		LastName:  sa.LastName,
		Address1:  sa.Address1,
		Address2:  sa.Address2,
		City:      sa.City,
		State:     sa.State,
		ZipCode:   sa.ZipCode,
		Country:   sa.Country,
		Phone:     sa.Phone,
		Type:      entities.AddressTypeShipping,
		IsActive:  true,
	}
}

// IsInternational checks if the address is international (non-US)
func (sa *SimpleAddress) IsInternational() bool {
	return sa.Country != "US" && sa.Country != "USA"
}

// ValidateShippingAddressRequest represents a request to validate shipping address
type ValidateShippingAddressRequest struct {
	Address    SimpleAddress `json:"address" validate:"required"`
	MethodID   *uuid.UUID    `json:"method_id,omitempty"`
	Weight     float64       `json:"weight" validate:"min=0"`
	OrderValue float64       `json:"order_value" validate:"min=0"`
}

// ValidateShippingAddressResponse represents the response for address validation
type ValidateShippingAddressResponse struct {
	IsValid              bool                    `json:"is_valid"`
	ValidationErrors     []string                `json:"validation_errors,omitempty"`
	CompatibleMethods    []ShippingMethodSummary `json:"compatible_methods"`
	IncompatibleMethods  []IncompatibleMethod    `json:"incompatible_methods"`
	Recommendations      []string                `json:"recommendations,omitempty"`
}

// ShippingMethodSummary represents a summary of shipping method
type ShippingMethodSummary struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Carrier     string  `json:"carrier"`
	EstimatedCost float64 `json:"estimated_cost"`
	DeliveryDays int     `json:"delivery_days"`
}

// IncompatibleMethod represents an incompatible shipping method with reason
type IncompatibleMethod struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

// ValidateShippingAddress validates shipping address and returns compatible methods
func (uc *shippingUseCase) ValidateShippingAddress(ctx context.Context, req ValidateShippingAddressRequest) (*ValidateShippingAddressResponse, error) {
	response := &ValidateShippingAddressResponse{
		IsValid:             true,
		ValidationErrors:    []string{},
		CompatibleMethods:   []ShippingMethodSummary{},
		IncompatibleMethods: []IncompatibleMethod{},
		Recommendations:     []string{},
	}

	// Validate address basic fields
	if err := req.Address.Validate(); err != nil {
		response.IsValid = false
		response.ValidationErrors = append(response.ValidationErrors, err.Error())
	}

	// Get all shipping methods
	methods, err := uc.shippingRepo.GetShippingMethods(ctx, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipping methods: %w", err)
	}

	// Convert SimpleAddress to entities.Address for compatibility checks
	address := req.Address.ToAddress()

	// Check compatibility with each method
	for _, method := range methods {
		// Validate method compatibility with address
		if err := uc.compatibilityService.ValidateShippingMethodForAddress(ctx, method, address); err != nil {
			response.IncompatibleMethods = append(response.IncompatibleMethods, IncompatibleMethod{
				ID:     method.ID.String(),
				Name:   method.Name,
				Reason: err.Error(),
			})
			continue
		}

		// Validate weight constraints
		if err := uc.compatibilityService.ValidateShippingConstraints(ctx, method, req.Weight, nil); err != nil {
			response.IncompatibleMethods = append(response.IncompatibleMethods, IncompatibleMethod{
				ID:     method.ID.String(),
				Name:   method.Name,
				Reason: err.Error(),
			})
			continue
		}

		// Method is compatible - calculate estimated cost
		estimatedCost := method.CalculateCost(req.Weight, 100, req.OrderValue) // Using 100km as default distance

		response.CompatibleMethods = append(response.CompatibleMethods, ShippingMethodSummary{
			ID:            method.ID.String(),
			Name:          method.Name,
			Type:          string(method.Type),
			Carrier:       method.Carrier,
			EstimatedCost: estimatedCost,
			DeliveryDays:  method.MaxDeliveryDays,
		})
	}

	// Add recommendations
	if req.Address.IsInternational() {
		response.Recommendations = append(response.Recommendations, "International shipping may require additional documentation")
		response.Recommendations = append(response.Recommendations, "Consider customs duties and taxes for international orders")
	}

	if req.Weight > 10 {
		response.Recommendations = append(response.Recommendations, "Heavy packages may have limited shipping options")
	}

	if len(response.CompatibleMethods) == 0 {
		response.IsValid = false
		response.ValidationErrors = append(response.ValidationErrors, "No compatible shipping methods available for this address")
	}

	return response, nil
}
