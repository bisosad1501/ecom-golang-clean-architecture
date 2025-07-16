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

	// Cleanup expired carts
	CleanupExpiredCarts(ctx context.Context) error

	// Cleanup expired payment timeouts
	CleanupExpiredPayments(ctx context.Context) error

	// Run full cleanup process
	RunCleanup(ctx context.Context) error

	// Get cleanup statistics
	GetCleanupStats(ctx context.Context) (map[string]interface{}, error)
}

type stockCleanupUseCase struct {
	stockReservationService services.StockReservationService
	orderRepo               repositories.OrderRepository
	stockReservationRepo    repositories.StockReservationRepository
	cartRepo                repositories.CartRepository // Add cart repository
}

// NewStockCleanupUseCase creates a new stock cleanup use case
func NewStockCleanupUseCase(
	stockReservationService services.StockReservationService,
	orderRepo repositories.OrderRepository,
	stockReservationRepo repositories.StockReservationRepository,
	cartRepo repositories.CartRepository, // Add cart repository
) StockCleanupUseCase {
	return &stockCleanupUseCase{
		stockReservationService: stockReservationService,
		orderRepo:               orderRepo,
		stockReservationRepo:    stockReservationRepo,
		cartRepo:                cartRepo,
	}
}

// CleanupExpiredReservations cleans up expired stock reservations
func (uc *stockCleanupUseCase) CleanupExpiredReservations(ctx context.Context) error {
	startTime := time.Now()
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

	duration := time.Since(startTime)
	fmt.Printf("✅ Successfully cleaned up %d expired reservations in %v\n", len(expiredReservations), duration)
	return nil
}

// CleanupExpiredOrders cleans up expired unpaid orders
func (uc *stockCleanupUseCase) CleanupExpiredOrders(ctx context.Context) error {
	startTime := time.Now()
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
	errorCount := 0
	for _, order := range orders {
		if order.IsPaymentExpired() {
			fmt.Printf("🕐 Order %s payment expired, cancelling...\n", order.OrderNumber)

			// Release stock reservations first
			if order.HasInventoryReserved() {
				if err := uc.stockReservationService.ReleaseReservations(ctx, order.ID); err != nil {
					fmt.Printf("❌ Failed to release reservations for order %s: %v\n", order.OrderNumber, err)
					errorCount++
					continue
				}
			}

			// Update order status to cancelled
			order.Status = entities.OrderStatusCancelled
			order.PaymentStatus = entities.PaymentStatusFailed
			order.ReleaseReservation()
			order.UpdatedAt = time.Now()

			if err := uc.orderRepo.Update(ctx, order); err != nil {
				fmt.Printf("❌ Failed to cancel expired order %s: %v\n", order.OrderNumber, err)
				errorCount++
				continue
			}

			expiredCount++
			fmt.Printf("✅ Cancelled expired order %s\n", order.OrderNumber)
		}
	}

	duration := time.Since(startTime)
	if expiredCount == 0 {
		fmt.Printf("✅ No expired orders found\n")
	} else if errorCount > 0 {
		fmt.Printf("⚠️ Cancelled %d expired orders with %d errors in %v\n", expiredCount, errorCount, duration)
	} else {
		fmt.Printf("✅ Successfully cancelled %d expired orders in %v\n", expiredCount, duration)
	}

	return nil
}

// RunCleanup runs the full cleanup process
func (uc *stockCleanupUseCase) RunCleanup(ctx context.Context) error {
	overallStart := time.Now()
	fmt.Printf("🚀 Starting full stock cleanup process...\n")

	var hasErrors bool

	// Cleanup expired reservations first
	if err := uc.CleanupExpiredReservations(ctx); err != nil {
		fmt.Printf("❌ Failed to cleanup expired reservations: %v\n", err)
		hasErrors = true
		// Continue with other cleanups even if reservation cleanup fails
	}

	// Cleanup expired orders
	if err := uc.CleanupExpiredOrders(ctx); err != nil {
		fmt.Printf("❌ Failed to cleanup expired orders: %v\n", err)
		hasErrors = true
		// Continue with cart cleanup even if order cleanup fails
	}

	// Cleanup expired carts
	if err := uc.CleanupExpiredCarts(ctx); err != nil {
		fmt.Printf("❌ Failed to cleanup expired carts: %v\n", err)
		hasErrors = true
	}

	// Cleanup expired payments
	if err := uc.CleanupExpiredPayments(ctx); err != nil {
		fmt.Printf("❌ Failed to cleanup expired payments: %v\n", err)
		hasErrors = true
	}

	overallDuration := time.Since(overallStart)
	if hasErrors {
		fmt.Printf("⚠️ Cleanup process completed with errors in %v\n", overallDuration)
		return fmt.Errorf("cleanup process completed with errors")
	}

	fmt.Printf("🎉 Full cleanup process completed successfully in %v\n", overallDuration)
	return nil
}

// CleanupExpiredCarts cleans up expired carts
func (uc *stockCleanupUseCase) CleanupExpiredCarts(ctx context.Context) error {
	startTime := time.Now()
	fmt.Printf("🧹 Starting cleanup of expired carts...\n")

	expiredCarts, err := uc.cartRepo.GetExpiredCarts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expired carts: %w", err)
	}

	if len(expiredCarts) == 0 {
		fmt.Printf("✅ No expired carts found\n")
		return nil
	}

	fmt.Printf("🔍 Found %d expired carts\n", len(expiredCarts))

	cleanedCount := 0
	errorCount := 0
	for _, cart := range expiredCarts {
		cart.MarkAsAbandoned()
		if err := uc.cartRepo.Update(ctx, cart); err != nil {
			fmt.Printf("❌ Failed to mark cart %s as abandoned: %v\n", cart.ID, err)
			errorCount++
			continue
		}
		cleanedCount++
		fmt.Printf("✅ Marked expired cart %s as abandoned\n", cart.ID)
	}

	duration := time.Since(startTime)
	if errorCount > 0 {
		fmt.Printf("⚠️ Cleaned up %d expired carts with %d errors in %v\n", cleanedCount, errorCount, duration)
	} else {
		fmt.Printf("✅ Successfully cleaned up %d expired carts in %v\n", cleanedCount, duration)
	}
	return nil
}

// CleanupExpiredPayments cleans up orders with expired payment timeouts
func (uc *stockCleanupUseCase) CleanupExpiredPayments(ctx context.Context) error {
	fmt.Printf("🧹 Starting cleanup of expired payment timeouts...\n")
	startTime := time.Now()

	// Get orders with expired payment timeouts
	filters := repositories.OrderSearchParams{
		PaymentStatus: &[]entities.PaymentStatus{entities.PaymentStatusPending}[0],
		Limit:         1000,
	}
	orders, err := uc.orderRepo.Search(ctx, filters)
	if err != nil {
		return fmt.Errorf("failed to get pending orders: %w", err)
	}

	cleanedCount := 0
	errorCount := 0

	for _, order := range orders {
		if !order.IsPaymentExpired() {
			continue // Skip non-expired orders
		}

		fmt.Printf("🔄 Processing expired payment for order %s (timeout: %v)\n",
			order.OrderNumber, order.PaymentTimeout)

		// Cancel the order due to payment timeout
		order.Status = entities.OrderStatusCancelled
		order.PaymentStatus = entities.PaymentStatusFailed
		order.ReleaseReservation()
		order.UpdatedAt = time.Now()

		if err := uc.orderRepo.Update(ctx, order); err != nil {
			fmt.Printf("❌ Failed to cancel expired order %s: %v\n", order.OrderNumber, err)
			errorCount++
			continue
		}

		// Release stock reservations
		if err := uc.stockReservationService.ReleaseReservations(ctx, order.ID); err != nil {
			fmt.Printf("❌ Failed to release reservations for order %s: %v\n", order.OrderNumber, err)
			// Continue anyway, order is already cancelled
		}

		cleanedCount++
		fmt.Printf("✅ Cancelled expired order %s and released reservations\n", order.OrderNumber)
	}

	duration := time.Since(startTime)
	if errorCount > 0 {
		fmt.Printf("⚠️ Cleaned up %d expired payments with %d errors in %v\n", cleanedCount, errorCount, duration)
	} else {
		fmt.Printf("✅ Successfully cleaned up %d expired payments in %v\n", cleanedCount, duration)
	}
	return nil
}

// GetCleanupStats returns statistics about items that need cleanup
func (uc *stockCleanupUseCase) GetCleanupStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get expired reservations count
	expiredReservations, err := uc.stockReservationRepo.GetExpiredReservations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired reservations: %w", err)
	}
	stats["expired_reservations"] = len(expiredReservations)

	// Get expired carts count
	expiredCarts, err := uc.cartRepo.GetExpiredCarts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired carts: %w", err)
	}
	stats["expired_carts"] = len(expiredCarts)

	// Get pending orders that might be expired
	filters := repositories.OrderSearchParams{
		Status:        &[]entities.OrderStatus{entities.OrderStatusPending}[0],
		PaymentStatus: &[]entities.PaymentStatus{entities.PaymentStatusPending}[0],
		Limit:         1000,
	}
	orders, err := uc.orderRepo.Search(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending orders: %w", err)
	}

	expiredOrderCount := 0
	for _, order := range orders {
		if order.IsPaymentExpired() {
			expiredOrderCount++
		}
	}
	stats["expired_orders"] = expiredOrderCount
	stats["total_pending_orders"] = len(orders)

	return stats, nil
}

// StartCleanupScheduler starts a background scheduler for cleanup tasks
func StartCleanupScheduler(ctx context.Context, cleanupUseCase StockCleanupUseCase) {
	ticker := time.NewTicker(5 * time.Minute) // Run every 5 minutes
	defer ticker.Stop()

	fmt.Printf("📅 Starting cleanup scheduler (every 5 minutes)...\n")

	// Run initial cleanup
	if err := cleanupUseCase.RunCleanup(ctx); err != nil {
		fmt.Printf("❌ Initial cleanup failed: %v\n", err)
	}

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
