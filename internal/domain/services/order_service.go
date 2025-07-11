package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
)

// OrderService handles order-related business logic
type OrderService interface {
	GenerateUniqueOrderNumber(ctx context.Context) (string, error)
	CalculateOrderTotal(items []entities.CartItem, taxRate, shippingCost, discountAmount float64) (subtotal, taxAmount, total float64)
	ValidateOrderItems(items []entities.CartItem) error
}

type orderService struct {
	orderRepo repositories.OrderRepository
}

// NewOrderService creates a new order service
func NewOrderService(orderRepo repositories.OrderRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
	}
}

// GenerateUniqueOrderNumber generates a unique order number
func (s *orderService) GenerateUniqueOrderNumber(ctx context.Context) (string, error) {
	const maxAttempts = 10

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Generate order number with format: ORD-YYYYMMDD-XXXXXX
		now := time.Now()
		dateStr := now.Format("20060102")

		// Generate cryptographically secure random 6-digit number
		randomBig, err := rand.Int(rand.Reader, big.NewInt(900000))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		randomNum := randomBig.Int64() + 100000

		orderNumber := fmt.Sprintf("ORD-%s-%d", dateStr, randomNum)

		// Check if order number already exists
		exists, err := s.orderRepo.ExistsByOrderNumber(ctx, orderNumber)
		if err != nil {
			return "", fmt.Errorf("failed to check order number existence: %w", err)
		}

		if !exists {
			return orderNumber, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique order number after %d attempts", maxAttempts)
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
