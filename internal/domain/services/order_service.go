package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
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
		// Generate order number with format: ORD-YYYYMMDD-HHMMSS-XXXX
		now := time.Now()
		dateStr := now.Format("20060102")
		timeStr := now.Format("150405") // HHMMSS format

		// Generate cryptographically secure random 4-digit number
		// Using smaller random part since we have time component for uniqueness
		randomBig, err := rand.Int(rand.Reader, big.NewInt(9000))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		randomNum := randomBig.Int64() + 1000

		orderNumber := fmt.Sprintf("ORD-%s-%s-%d", dateStr, timeStr, randomNum)

		// Check if order number already exists
		exists, err := s.orderRepo.ExistsByOrderNumber(ctx, orderNumber)
		if err != nil {
			return "", fmt.Errorf("failed to check order number existence: %w", err)
		}

		if !exists {
			return orderNumber, nil
		}

		// Add small delay between attempts to reduce collision probability
		if attempt < maxAttempts-1 {
			time.Sleep(time.Millisecond * 10)
		}
	}

	return "", fmt.Errorf("failed to generate unique order number after %d attempts", maxAttempts)
}

// CalculateOrderTotal calculates the order totals
func (s *orderService) CalculateOrderTotal(items []entities.CartItem, taxRate, shippingCost, discountAmount float64) (subtotal, taxAmount, total float64) {
	// Validate inputs
	if taxRate < 0 {
		taxRate = 0
	}
	if taxRate > 1 { // Assume tax rate is percentage (0.1 = 10%)
		taxRate = taxRate / 100
	}
	if shippingCost < 0 {
		shippingCost = 0
	}
	if discountAmount < 0 {
		discountAmount = 0
	}

	// Calculate subtotal
	for _, item := range items {
		subtotal += item.GetSubtotal()
	}

	// Calculate tax amount (round to 2 decimal places)
	taxAmount = subtotal * taxRate
	taxAmount = float64(int(taxAmount*100+0.5)) / 100

	// Calculate total
	total = subtotal + taxAmount + shippingCost - discountAmount

	// Ensure discount doesn't exceed subtotal + tax + shipping
	maxDiscount := subtotal + taxAmount + shippingCost
	if discountAmount > maxDiscount {
		discountAmount = maxDiscount
		total = 0
	} else {
		// Round total to 2 decimal places
		total = float64(int(total*100+0.5)) / 100
	}

	// Ensure total is not negative
	if total < 0 {
		total = 0
	}

	return subtotal, taxAmount, total
}

// ValidateOrderItems validates order items
func (s *orderService) ValidateOrderItems(items []entities.CartItem) error {
	if len(items) == 0 {
		return fmt.Errorf("order must contain at least one item")
	}

	// Track product IDs to check for duplicates
	productIDs := make(map[string]bool)
	totalItems := 0

	for i, item := range items {
		// Validate basic fields
		if item.Quantity <= 0 {
			return fmt.Errorf("item %d: quantity must be greater than 0", i+1)
		}

		if item.Quantity > 100 {
			return fmt.Errorf("item %d: quantity cannot exceed 100", i+1)
		}

		if item.Price <= 0 {
			return fmt.Errorf("item %d: price must be greater than 0", i+1)
		}

		if item.Price > 999999.99 {
			return fmt.Errorf("item %d: price cannot exceed $999,999.99", i+1)
		}

		// Validate ProductID
		if item.ProductID.String() == "00000000-0000-0000-0000-000000000000" {
			return fmt.Errorf("item %d: invalid product ID", i+1)
		}

		// Check for duplicate products
		productIDStr := item.ProductID.String()
		if productIDs[productIDStr] {
			return fmt.Errorf("item %d: duplicate product in order", i+1)
		}
		productIDs[productIDStr] = true

		// Validate calculated total with floating point tolerance
		expectedTotal := float64(item.Quantity) * item.Price
		const epsilon = 0.01
		if math.Abs(item.GetSubtotal() - expectedTotal) > epsilon {
			return fmt.Errorf("item %d: subtotal %.2f does not match calculated subtotal %.2f", i+1, item.GetSubtotal(), expectedTotal)
		}

		totalItems += item.Quantity
	}

	// Validate total items limit
	if totalItems > 1000 {
		return fmt.Errorf("total items in order (%d) cannot exceed 1000", totalItems)
	}

	return nil
}
