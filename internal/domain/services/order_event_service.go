package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

// OrderEventService handles order event business logic
type OrderEventService interface {
	// Create order events
	CreateEvent(ctx context.Context, orderID uuid.UUID, eventType entities.OrderEventType, title, description string, data interface{}, userID *uuid.UUID, isPublic bool) error
	
	// Predefined event creators
	CreateOrderCreatedEvent(ctx context.Context, order *entities.Order, userID *uuid.UUID) error
	CreateStatusChangedEvent(ctx context.Context, orderID uuid.UUID, oldStatus, newStatus entities.OrderStatus, userID *uuid.UUID) error
	CreatePaymentReceivedEvent(ctx context.Context, orderID uuid.UUID, amount float64, paymentMethod string, userID *uuid.UUID) error
	CreatePaymentFailedEvent(ctx context.Context, orderID uuid.UUID, reason string, userID *uuid.UUID) error
	CreateShippedEvent(ctx context.Context, orderID uuid.UUID, trackingNumber, carrier string, userID *uuid.UUID) error
	CreateDeliveredEvent(ctx context.Context, orderID uuid.UUID, userID *uuid.UUID) error
	CreateCancelledEvent(ctx context.Context, orderID uuid.UUID, reason string, userID *uuid.UUID) error
	CreateRefundedEvent(ctx context.Context, orderID uuid.UUID, amount float64, reason string, userID *uuid.UUID) error
	CreateNoteAddedEvent(ctx context.Context, orderID uuid.UUID, note string, userID *uuid.UUID, isPublic bool) error
	CreateTrackingUpdatedEvent(ctx context.Context, orderID uuid.UUID, trackingNumber, status string, userID *uuid.UUID) error
	CreateInventoryReservedEvent(ctx context.Context, orderID uuid.UUID, items []entities.CartItem, userID *uuid.UUID) error
	CreateInventoryReleasedEvent(ctx context.Context, orderID uuid.UUID, reason string, userID *uuid.UUID) error
	
	// Get events
	GetOrderEvents(ctx context.Context, orderID uuid.UUID, publicOnly bool) ([]*entities.OrderEvent, error)
	GetOrderTimeline(ctx context.Context, orderID uuid.UUID) ([]*entities.OrderEvent, error)
}

type orderEventService struct {
	eventRepo repositories.OrderEventRepository
}

// NewOrderEventService creates a new order event service
func NewOrderEventService(eventRepo repositories.OrderEventRepository) OrderEventService {
	return &orderEventService{
		eventRepo: eventRepo,
	}
}

// CreateEvent creates a generic order event
func (s *orderEventService) CreateEvent(ctx context.Context, orderID uuid.UUID, eventType entities.OrderEventType, title, description string, data interface{}, userID *uuid.UUID, isPublic bool) error {
	event := &entities.OrderEvent{
		ID:          uuid.New(),
		OrderID:     orderID,
		EventType:   eventType,
		Title:       title,
		Description: description,
		UserID:      userID,
		IsPublic:    isPublic,
		CreatedAt:   time.Now(),
	}
	
	// Serialize data to JSON if provided
	if data != nil {
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to serialize event data: %w", err)
		}
		event.Data = string(dataBytes)
	}
	
	return s.eventRepo.Create(ctx, event)
}

// CreateOrderCreatedEvent creates an order created event
func (s *orderEventService) CreateOrderCreatedEvent(ctx context.Context, order *entities.Order, userID *uuid.UUID) error {
	data := map[string]interface{}{
		"order_number": order.OrderNumber,
		"total":        order.Total,
		"currency":     order.Currency,
		"item_count":   order.GetItemCount(),
	}
	
	return s.CreateEvent(
		ctx,
		order.ID,
		entities.OrderEventTypeCreated,
		"Order Created",
		fmt.Sprintf("Order %s has been created with total %s %.2f", order.OrderNumber, order.Currency, order.Total),
		data,
		userID,
		true,
	)
}

// CreateStatusChangedEvent creates a status changed event
func (s *orderEventService) CreateStatusChangedEvent(ctx context.Context, orderID uuid.UUID, oldStatus, newStatus entities.OrderStatus, userID *uuid.UUID) error {
	data := map[string]interface{}{
		"old_status": oldStatus,
		"new_status": newStatus,
	}
	
	return s.CreateEvent(
		ctx,
		orderID,
		entities.OrderEventTypeStatusChanged,
		"Status Changed",
		fmt.Sprintf("Order status changed from %s to %s", oldStatus, newStatus),
		data,
		userID,
		true,
	)
}

// CreatePaymentReceivedEvent creates a payment received event
func (s *orderEventService) CreatePaymentReceivedEvent(ctx context.Context, orderID uuid.UUID, amount float64, paymentMethod string, userID *uuid.UUID) error {
	data := map[string]interface{}{
		"amount":         amount,
		"payment_method": paymentMethod,
	}
	
	return s.CreateEvent(
		ctx,
		orderID,
		entities.OrderEventTypePaymentReceived,
		"Payment Received",
		fmt.Sprintf("Payment of $%.2f received via %s", amount, paymentMethod),
		data,
		userID,
		true,
	)
}

// CreatePaymentFailedEvent creates a payment failed event
func (s *orderEventService) CreatePaymentFailedEvent(ctx context.Context, orderID uuid.UUID, reason string, userID *uuid.UUID) error {
	data := map[string]interface{}{
		"reason": reason,
	}
	
	return s.CreateEvent(
		ctx,
		orderID,
		entities.OrderEventTypePaymentFailed,
		"Payment Failed",
		fmt.Sprintf("Payment failed: %s", reason),
		data,
		userID,
		true,
	)
}

// CreateShippedEvent creates a shipped event
func (s *orderEventService) CreateShippedEvent(ctx context.Context, orderID uuid.UUID, trackingNumber, carrier string, userID *uuid.UUID) error {
	data := map[string]interface{}{
		"tracking_number": trackingNumber,
		"carrier":         carrier,
	}
	
	return s.CreateEvent(
		ctx,
		orderID,
		entities.OrderEventTypeShipped,
		"Order Shipped",
		fmt.Sprintf("Order shipped via %s with tracking number %s", carrier, trackingNumber),
		data,
		userID,
		true,
	)
}

// CreateDeliveredEvent creates a delivered event
func (s *orderEventService) CreateDeliveredEvent(ctx context.Context, orderID uuid.UUID, userID *uuid.UUID) error {
	return s.CreateEvent(
		ctx,
		orderID,
		entities.OrderEventTypeDelivered,
		"Order Delivered",
		"Order has been successfully delivered",
		nil,
		userID,
		true,
	)
}

// CreateCancelledEvent creates a cancelled event
func (s *orderEventService) CreateCancelledEvent(ctx context.Context, orderID uuid.UUID, reason string, userID *uuid.UUID) error {
	data := map[string]interface{}{
		"reason": reason,
	}
	
	return s.CreateEvent(
		ctx,
		orderID,
		entities.OrderEventTypeCancelled,
		"Order Cancelled",
		fmt.Sprintf("Order cancelled: %s", reason),
		data,
		userID,
		true,
	)
}

// CreateRefundedEvent creates a refunded event
func (s *orderEventService) CreateRefundedEvent(ctx context.Context, orderID uuid.UUID, amount float64, reason string, userID *uuid.UUID) error {
	data := map[string]interface{}{
		"amount": amount,
		"reason": reason,
	}
	
	return s.CreateEvent(
		ctx,
		orderID,
		entities.OrderEventTypeRefunded,
		"Order Refunded",
		fmt.Sprintf("Refund of $%.2f processed: %s", amount, reason),
		data,
		userID,
		true,
	)
}

// CreateNoteAddedEvent creates a note added event
func (s *orderEventService) CreateNoteAddedEvent(ctx context.Context, orderID uuid.UUID, note string, userID *uuid.UUID, isPublic bool) error {
	data := map[string]interface{}{
		"note": note,
	}
	
	visibility := "internal"
	if isPublic {
		visibility = "public"
	}
	
	return s.CreateEvent(
		ctx,
		orderID,
		entities.OrderEventTypeNoteAdded,
		"Note Added",
		fmt.Sprintf("A %s note has been added to the order", visibility),
		data,
		userID,
		isPublic,
	)
}

// CreateTrackingUpdatedEvent creates a tracking updated event
func (s *orderEventService) CreateTrackingUpdatedEvent(ctx context.Context, orderID uuid.UUID, trackingNumber, status string, userID *uuid.UUID) error {
	data := map[string]interface{}{
		"tracking_number": trackingNumber,
		"status":          status,
	}
	
	return s.CreateEvent(
		ctx,
		orderID,
		entities.OrderEventTypeTrackingUpdated,
		"Tracking Updated",
		fmt.Sprintf("Tracking status updated: %s", status),
		data,
		userID,
		true,
	)
}

// CreateInventoryReservedEvent creates an inventory reserved event
func (s *orderEventService) CreateInventoryReservedEvent(ctx context.Context, orderID uuid.UUID, items []entities.CartItem, userID *uuid.UUID) error {
	data := map[string]interface{}{
		"items": items,
	}
	
	return s.CreateEvent(
		ctx,
		orderID,
		entities.OrderEventTypeInventoryReserved,
		"Inventory Reserved",
		"Inventory has been reserved for this order",
		data,
		userID,
		false, // Internal event
	)
}

// CreateInventoryReleasedEvent creates an inventory released event
func (s *orderEventService) CreateInventoryReleasedEvent(ctx context.Context, orderID uuid.UUID, reason string, userID *uuid.UUID) error {
	data := map[string]interface{}{
		"reason": reason,
	}
	
	return s.CreateEvent(
		ctx,
		orderID,
		entities.OrderEventTypeInventoryReleased,
		"Inventory Released",
		fmt.Sprintf("Inventory reservation released: %s", reason),
		data,
		userID,
		false, // Internal event
	)
}

// GetOrderEvents gets order events
func (s *orderEventService) GetOrderEvents(ctx context.Context, orderID uuid.UUID, publicOnly bool) ([]*entities.OrderEvent, error) {
	if publicOnly {
		return s.eventRepo.GetPublicByOrderID(ctx, orderID)
	}
	return s.eventRepo.GetByOrderID(ctx, orderID)
}

// GetOrderTimeline gets order timeline (public events only)
func (s *orderEventService) GetOrderTimeline(ctx context.Context, orderID uuid.UUID) ([]*entities.OrderEvent, error) {
	return s.eventRepo.GetPublicByOrderID(ctx, orderID)
}
