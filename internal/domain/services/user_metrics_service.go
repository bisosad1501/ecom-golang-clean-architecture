package services

import (
	"context"
	"fmt"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

// UserMetricsService handles user metrics calculations and updates
type UserMetricsService interface {
	UpdateUserMetricsOnOrderConfirmed(ctx context.Context, userID uuid.UUID, orderTotal float64) error
	UpdateUserMetricsOnOrderCancelled(ctx context.Context, userID uuid.UUID, orderTotal float64) error
	RecalculateUserMetrics(ctx context.Context, userID uuid.UUID) error
	UpdateLoyaltyPoints(ctx context.Context, userID uuid.UUID, points int) error
	UpdateMembershipTier(ctx context.Context, userID uuid.UUID) error
}

type userMetricsService struct {
	userRepo  repositories.UserRepository
	orderRepo repositories.OrderRepository
}

// NewUserMetricsService creates a new user metrics service
func NewUserMetricsService(
	userRepo repositories.UserRepository,
	orderRepo repositories.OrderRepository,
) UserMetricsService {
	return &userMetricsService{
		userRepo:  userRepo,
		orderRepo: orderRepo,
	}
}

// UpdateUserMetricsOnOrderConfirmed updates user metrics when order is confirmed
func (s *userMetricsService) UpdateUserMetricsOnOrderConfirmed(ctx context.Context, userID uuid.UUID, orderTotal float64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Update metrics
	user.TotalOrders++
	user.TotalSpent += orderTotal

	// Calculate loyalty points (1 point per $1 spent)
	loyaltyPointsEarned := int(orderTotal)
	user.LoyaltyPoints += loyaltyPointsEarned

	// Update user in database
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user metrics: %w", err)
	}

	// Update membership tier based on new metrics
	if err := s.UpdateMembershipTier(ctx, userID); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: Failed to update membership tier for user %s: %v\n", userID, err)
	}

	return nil
}

// UpdateUserMetricsOnOrderCancelled updates user metrics when order is cancelled
func (s *userMetricsService) UpdateUserMetricsOnOrderCancelled(ctx context.Context, userID uuid.UUID, orderTotal float64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Only decrease if metrics are positive
	if user.TotalOrders > 0 {
		user.TotalOrders--
	}
	if user.TotalSpent >= orderTotal {
		user.TotalSpent -= orderTotal
	}

	// Remove loyalty points (1 point per $1)
	loyaltyPointsToRemove := int(orderTotal)
	if user.LoyaltyPoints >= loyaltyPointsToRemove {
		user.LoyaltyPoints -= loyaltyPointsToRemove
	}

	// Update user in database
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user metrics: %w", err)
	}

	// Update membership tier based on new metrics
	if err := s.UpdateMembershipTier(ctx, userID); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: Failed to update membership tier for user %s: %v\n", userID, err)
	}

	return nil
}

// RecalculateUserMetrics recalculates user metrics from actual order data
func (s *userMetricsService) RecalculateUserMetrics(ctx context.Context, userID uuid.UUID) error {
	// Get all confirmed orders for user
	orders, err := s.orderRepo.GetByUserID(ctx, userID, 0, 1000) // Get first 1000 orders
	if err != nil {
		return fmt.Errorf("failed to get user orders: %w", err)
	}

	// Calculate metrics from actual orders
	var totalOrders int
	var totalSpent float64

	for _, order := range orders {
		if order.Status == entities.OrderStatusConfirmed || order.Status == entities.OrderStatusDelivered {
			totalOrders++
			totalSpent += order.Total
		}
	}

	// Get user and update metrics
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.TotalOrders = totalOrders
	user.TotalSpent = totalSpent

	// Recalculate loyalty points based on total spent
	user.LoyaltyPoints = int(totalSpent)

	// Update user in database
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user metrics: %w", err)
	}

	// Update membership tier
	return s.UpdateMembershipTier(ctx, userID)
}

// UpdateLoyaltyPoints updates user loyalty points
func (s *userMetricsService) UpdateLoyaltyPoints(ctx context.Context, userID uuid.UUID, points int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.LoyaltyPoints += points

	return s.userRepo.Update(ctx, user)
}

// UpdateMembershipTier updates user membership tier based on total spent
func (s *userMetricsService) UpdateMembershipTier(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Determine membership tier based on total spent
	var newTier string
	switch {
	case user.TotalSpent >= 10000: // $10,000+
		newTier = "platinum"
	case user.TotalSpent >= 5000: // $5,000+
		newTier = "gold"
	case user.TotalSpent >= 1000: // $1,000+
		newTier = "silver"
	default:
		newTier = "bronze"
	}

	// Only update if tier changed
	if user.MembershipTier != newTier {
		user.MembershipTier = newTier
		return s.userRepo.Update(ctx, user)
	}

	return nil
}
