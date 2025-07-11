package usecases

import (
	"context"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/domain/services"
)

// StockCleanupUseCase handles cleanup of expired stock reservations and orders
type StockCleanupUseCase interface {
	// Cleanup expired stock reservations
	CleanupExpiredReservations(ctx context.Context) error
	
	// Cleanup expired unpaid orders
	CleanupExpiredOrders(ctx context.Context) error
	
	// Run full cleanup process
	RunCleanup(ctx context.Context) error
}

type stockCleanupUseCase struct {
	stockReservationService services.StockReservationService
	orderRepo               repositories.OrderRepository
	stockReservationRepo    repositories.StockReservationRepository
}

// NewStockCleanupUseCase creates a new stock cleanup use case
func NewStockCleanupUseCase(
	stockReservationService services.StockReservationService,
	orderRepo repositories.OrderRepository,
	stockReservationRepo repositories.StockReservationRepository,
) StockCleanupUseCase {
	return &stockCleanupUseCase{
		stockReservationService: stockReservationService,
		orderRepo:               orderRepo,
		stockReservationRepo:    stockReservationRepo,
	}
}

// CleanupExpiredReservations cleans up expired stock reservations
func (uc *stockCleanupUseCase) CleanupExpiredReservations(ctx context.Context) error {
	fmt.Printf("🧹 Starting cleanup of expired stock reservations...\n")
	
	// Get expired reservations
	expiredReservations, err := uc.stockReservationRepo.GetExpiredReservations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expired reservations: %w", err)
	}
	
	if len(expiredReservations) == 0 {
		fmt.Printf("✅ No expired reservations found\n")
		return nil
	}
	
	fmt.Printf("🔍 Found %d expired reservations\n", len(expiredReservations))
	
	// Release expired reservations
	if err := uc.stockReservationService.CleanupExpiredReservations(ctx); err != nil {
		return fmt.Errorf("failed to cleanup expired reservations: %w", err)
	}
	
	fmt.Printf("✅ Successfully cleaned up %d expired reservations\n", len(expiredReservations))
	return nil
}

// CleanupExpiredOrders cleans up expired unpaid orders
func (uc *stockCleanupUseCase) CleanupExpiredOrders(ctx context.Context) error {
	fmt.Printf("🧹 Starting cleanup of expired unpaid orders...\n")
	
	// Get orders that are pending payment and have expired payment timeout
	filters := repositories.OrderSearchParams{
		Status:        &[]entities.OrderStatus{entities.OrderStatusPending}[0],
		PaymentStatus: &[]entities.PaymentStatus{entities.PaymentStatusPending}[0],
		Limit:         100, // Process in batches
	}
	
	orders, err := uc.orderRepo.Search(ctx, filters)
	if err != nil {
		return fmt.Errorf("failed to get pending orders: %w", err)
	}
	
	expiredCount := 0
	for _, order := range orders {
		if order.IsPaymentExpired() {
			fmt.Printf("🕐 Order %s payment expired, cancelling...\n", order.OrderNumber)
			
			// Release stock reservations
			if err := uc.stockReservationService.ReleaseReservations(ctx, order.ID); err != nil {
				fmt.Printf("❌ Failed to release reservations for order %s: %v\n", order.OrderNumber, err)
				continue
			}
			
			// Update order status to cancelled
			order.Status = entities.OrderStatusCancelled
			order.ReleaseReservation()
			order.UpdatedAt = time.Now()
			
			if err := uc.orderRepo.Update(ctx, order); err != nil {
				fmt.Printf("❌ Failed to cancel expired order %s: %v\n", order.OrderNumber, err)
				continue
			}
			
			expiredCount++
			fmt.Printf("✅ Cancelled expired order %s\n", order.OrderNumber)
		}
	}
	
	if expiredCount == 0 {
		fmt.Printf("✅ No expired orders found\n")
	} else {
		fmt.Printf("✅ Successfully cancelled %d expired orders\n", expiredCount)
	}
	
	return nil
}

// RunCleanup runs the full cleanup process
func (uc *stockCleanupUseCase) RunCleanup(ctx context.Context) error {
	fmt.Printf("🚀 Starting full stock cleanup process...\n")
	
	// Cleanup expired reservations first
	if err := uc.CleanupExpiredReservations(ctx); err != nil {
		fmt.Printf("❌ Failed to cleanup expired reservations: %v\n", err)
		// Continue with order cleanup even if reservation cleanup fails
	}
	
	// Cleanup expired orders
	if err := uc.CleanupExpiredOrders(ctx); err != nil {
		fmt.Printf("❌ Failed to cleanup expired orders: %v\n", err)
		return err
	}
	
	fmt.Printf("🎉 Full cleanup process completed successfully\n")
	return nil
}

// StartCleanupScheduler starts a background scheduler for cleanup tasks
func StartCleanupScheduler(ctx context.Context, cleanupUseCase StockCleanupUseCase) {
	ticker := time.NewTicker(5 * time.Minute) // Run every 5 minutes
	defer ticker.Stop()
	
	fmt.Printf("📅 Starting cleanup scheduler (every 5 minutes)...\n")
	
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("🛑 Cleanup scheduler stopped\n")
			return
		case <-ticker.C:
			if err := cleanupUseCase.RunCleanup(ctx); err != nil {
				fmt.Printf("❌ Scheduled cleanup failed: %v\n", err)
			}
		}
	}
}
