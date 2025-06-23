package services

import (
	"fmt"
	"math/rand"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
)

// OrderService handles order-related business logic
type OrderService interface {
	GenerateOrderNumber() string
	CalculateOrderTotal(items []entities.CartItem, taxRate, shippingCost, discountAmount float64) (subtotal, taxAmount, total float64)
	ValidateOrderItems(items []entities.CartItem) error
}

type orderService struct{}

// NewOrderService creates a new order service
func NewOrderService() OrderService {
	return &orderService{}
}

// GenerateOrderNumber generates a unique order number
func (s *orderService) GenerateOrderNumber() string {
	// Generate order number with format: ORD-YYYYMMDD-XXXXXX
	now := time.Now()
	dateStr := now.Format("20060102")
	
	// Generate random 6-digit number
	rand.Seed(now.UnixNano())
	randomNum := rand.Intn(999999-100000) + 100000
	
	return fmt.Sprintf("ORD-%s-%d", dateStr, randomNum)
}

// CalculateOrderTotal calculates the order totals
func (s *orderService) CalculateOrderTotal(items []entities.CartItem, taxRate, shippingCost, discountAmount float64) (subtotal, taxAmount, total float64) {
	// Calculate subtotal
	for _, item := range items {
		subtotal += item.GetSubtotal()
	}
	
	// Calculate tax amount
	taxAmount = subtotal * taxRate
	
	// Calculate total
	total = subtotal + taxAmount + shippingCost - discountAmount
	
	// Ensure total is not negative
	if total < 0 {
		total = 0
	}
	
	return subtotal, taxAmount, total
}

// ValidateOrderItems validates order items
func (s *orderService) ValidateOrderItems(items []entities.CartItem) error {
	if len(items) == 0 {
		return entities.ErrInvalidInput
	}
	
	for _, item := range items {
		if item.Quantity <= 0 {
			return entities.ErrInvalidQuantity
		}
		
		if item.Price <= 0 {
			return entities.ErrInvalidInput
		}
	}
	
	return nil
}
