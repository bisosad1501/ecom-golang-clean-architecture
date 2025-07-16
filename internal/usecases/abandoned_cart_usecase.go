package usecases

import (
	"context"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
)

// AbandonedCartUseCase defines the interface for abandoned cart recovery
type AbandonedCartUseCase interface {
	// Detection and recovery
	DetectAbandonedCarts(ctx context.Context) error
	SendAbandonedCartEmails(ctx context.Context) error

	// Analytics
	GetAbandonedCartStats(ctx context.Context, since time.Time) (*AbandonedCartStats, error)
	GetAbandonedCarts(ctx context.Context, offset, limit int) ([]*AbandonedCartResponse, error)

	// Recovery tracking
	TrackCartRecovery(ctx context.Context, cartID uuid.UUID) error
	GetRecoveryRate(ctx context.Context, since time.Time) (float64, error)
}

type abandonedCartUseCase struct {
	cartRepo     repositories.CartRepository
	userRepo     repositories.UserRepository
	emailUseCase EmailUseCase
	productRepo  repositories.ProductRepository
	orderRepo    repositories.OrderRepository
}

// NewAbandonedCartUseCase creates a new abandoned cart use case
func NewAbandonedCartUseCase(
	cartRepo repositories.CartRepository,
	userRepo repositories.UserRepository,
	emailUseCase EmailUseCase,
	productRepo repositories.ProductRepository,
	orderRepo repositories.OrderRepository,
) AbandonedCartUseCase {
	return &abandonedCartUseCase{
		cartRepo:     cartRepo,
		userRepo:     userRepo,
		emailUseCase: emailUseCase,
		productRepo:  productRepo,
		orderRepo:    orderRepo,
	}
}

// DetectAbandonedCarts detects carts that have been abandoned
func (uc *abandonedCartUseCase) DetectAbandonedCarts(ctx context.Context) error {
	// Define abandonment criteria
	abandonmentThreshold := time.Now().Add(-24 * time.Hour) // 24 hours ago

	// Get carts that haven't been updated recently
	carts, err := uc.cartRepo.GetAbandonedCarts(ctx, abandonmentThreshold)
	if err != nil {
		return fmt.Errorf("failed to get abandoned carts: %w", err)
	}

	fmt.Printf("🔍 Found %d potentially abandoned carts\n", len(carts))

	abandonedCount := 0
	for _, cart := range carts {
		// Skip if cart is empty or has no user
		if len(cart.Items) == 0 || cart.UserID == nil {
			continue
		}

		// Skip if user has already completed an order recently
		hasRecentOrder, err := uc.hasRecentOrder(ctx, *cart.UserID, cart.UpdatedAt)
		if err != nil {
			fmt.Printf("❌ Failed to check recent orders for user %s: %v\n", *cart.UserID, err)
			continue
		}
		if hasRecentOrder {
			continue
		}

		// Mark as abandoned if not already marked
		if !cart.IsAbandoned {
			cart.IsAbandoned = true
			cart.AbandonedAt = &cart.UpdatedAt

			if err := uc.cartRepo.Update(ctx, cart); err != nil {
				fmt.Printf("❌ Failed to mark cart as abandoned: %v\n", err)
				continue
			}

			abandonedCount++
		}

		// Check if we should send reminder emails
		if cart.AbandonedAt != nil {
			timeSinceAbandoned := time.Since(*cart.AbandonedAt)

			// Send first reminder after 1 hour
			if timeSinceAbandoned >= time.Hour && cart.FirstReminderSent == nil {
				if err := uc.sendFirstReminder(ctx, cart); err != nil {
					fmt.Printf("❌ Failed to send first reminder for cart %s: %v\n", cart.ID, err)
				} else {
					now := time.Now()
					cart.FirstReminderSent = &now
					_ = uc.cartRepo.Update(ctx, cart)
				}
			}

			// Send second reminder after 24 hours
			if timeSinceAbandoned >= 24*time.Hour && cart.SecondReminderSent == nil {
				if err := uc.sendSecondReminder(ctx, cart); err != nil {
					fmt.Printf("❌ Failed to send second reminder for cart %s: %v\n", cart.ID, err)
				} else {
					now := time.Now()
					cart.SecondReminderSent = &now
					_ = uc.cartRepo.Update(ctx, cart)
				}
			}

			// Send final reminder after 72 hours
			if timeSinceAbandoned >= 72*time.Hour && cart.FinalReminderSent == nil {
				if err := uc.sendFinalReminder(ctx, cart); err != nil {
					fmt.Printf("❌ Failed to send final reminder for cart %s: %v\n", cart.ID, err)
				} else {
					now := time.Now()
					cart.FinalReminderSent = &now
					_ = uc.cartRepo.Update(ctx, cart)
				}
			}
		}
	}

	fmt.Printf("✅ Marked %d carts as abandoned\n", abandonedCount)
	return nil
}

// hasRecentOrder checks if user has completed an order recently
func (uc *abandonedCartUseCase) hasRecentOrder(ctx context.Context, userID uuid.UUID, since time.Time) (bool, error) {
	orders, err := uc.orderRepo.GetByUserID(ctx, userID, 0, 5) // Check last 5 orders
	if err != nil {
		return false, err
	}

	for _, order := range orders {
		if order.CreatedAt.After(since) && order.Status != entities.OrderStatusCancelled {
			return true, nil
		}
	}

	return false, nil
}

// sendFirstReminder sends the first abandonment reminder
func (uc *abandonedCartUseCase) sendFirstReminder(ctx context.Context, cart *entities.Cart) error {
	if cart.UserID == nil {
		return fmt.Errorf("cart has no user ID")
	}

	return uc.emailUseCase.SendAbandonedCartEmail(ctx, *cart.UserID)
}

// sendSecondReminder sends the second abandonment reminder
func (uc *abandonedCartUseCase) sendSecondReminder(ctx context.Context, cart *entities.Cart) error {
	if cart.UserID == nil {
		return fmt.Errorf("cart has no user ID")
	}

	return uc.emailUseCase.SendAbandonedCartEmail(ctx, *cart.UserID)
}

// sendFinalReminder sends the final abandonment reminder
func (uc *abandonedCartUseCase) sendFinalReminder(ctx context.Context, cart *entities.Cart) error {
	if cart.UserID == nil {
		return fmt.Errorf("cart has no user ID")
	}

	return uc.emailUseCase.SendAbandonedCartEmail(ctx, *cart.UserID)
}

// SendAbandonedCartEmails sends emails for abandoned carts
func (uc *abandonedCartUseCase) SendAbandonedCartEmails(ctx context.Context) error {
	return uc.DetectAbandonedCarts(ctx)
}

// TrackCartRecovery tracks when an abandoned cart is recovered
func (uc *abandonedCartUseCase) TrackCartRecovery(ctx context.Context, cartID uuid.UUID) error {
	cart, err := uc.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	if cart.IsAbandoned {
		cart.IsAbandoned = false
		now := time.Now()
		cart.RecoveredAt = &now

		return uc.cartRepo.Update(ctx, cart)
	}

	return nil
}

// GetAbandonedCartStats returns abandoned cart statistics
func (uc *abandonedCartUseCase) GetAbandonedCartStats(ctx context.Context, since time.Time) (*AbandonedCartStats, error) {
	stats, err := uc.cartRepo.GetAbandonedCartStats(ctx, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get abandoned cart stats: %w", err)
	}

	return &AbandonedCartStats{
		TotalAbandoned:     stats.TotalAbandoned,
		TotalRecovered:     stats.TotalRecovered,
		RecoveryRate:       stats.RecoveryRate,
		AverageCartValue:   stats.AverageCartValue,
		TotalLostRevenue:   stats.TotalLostRevenue,
		RecoveredRevenue:   stats.RecoveredRevenue,
		FirstReminderSent:  stats.FirstReminderSent,
		SecondReminderSent: stats.SecondReminderSent,
		FinalReminderSent:  stats.FinalReminderSent,
		Since:              since,
		Until:              time.Now(),
	}, nil
}

// GetAbandonedCarts returns a list of abandoned carts
func (uc *abandonedCartUseCase) GetAbandonedCarts(ctx context.Context, offset, limit int) ([]*AbandonedCartResponse, error) {
	carts, err := uc.cartRepo.GetAbandonedCartsList(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get abandoned carts: %w", err)
	}

	responses := make([]*AbandonedCartResponse, len(carts))
	for i, cart := range carts {
		var user *entities.User
		var userID uuid.UUID
		if cart.UserID != nil {
			userID = *cart.UserID
			user, _ = uc.userRepo.GetByID(ctx, userID)
		}

		total := 0.0
		for _, item := range cart.Items {
			total += item.Price * float64(item.Quantity)
		}

		userEmail := ""
		userName := ""
		if user != nil {
			userEmail = user.Email
			userName = user.GetFullName()
		}

		responses[i] = &AbandonedCartResponse{
			ID:                 cart.ID,
			UserID:             userID,
			UserEmail:          userEmail,
			UserName:           userName,
			ItemCount:          len(cart.Items),
			Total:              total,
			AbandonedAt:        cart.AbandonedAt,
			FirstReminderSent:  cart.FirstReminderSent,
			SecondReminderSent: cart.SecondReminderSent,
			FinalReminderSent:  cart.FinalReminderSent,
			RecoveredAt:        cart.RecoveredAt,
			CreatedAt:          cart.CreatedAt,
			UpdatedAt:          cart.UpdatedAt,
		}
	}

	return responses, nil
}

// GetRecoveryRate returns the cart recovery rate
func (uc *abandonedCartUseCase) GetRecoveryRate(ctx context.Context, since time.Time) (float64, error) {
	stats, err := uc.GetAbandonedCartStats(ctx, since)
	if err != nil {
		return 0, err
	}

	return stats.RecoveryRate, nil
}

// Response types
type AbandonedCartStats struct {
	TotalAbandoned     int64     `json:"total_abandoned"`
	TotalRecovered     int64     `json:"total_recovered"`
	RecoveryRate       float64   `json:"recovery_rate"`
	AverageCartValue   float64   `json:"average_cart_value"`
	TotalLostRevenue   float64   `json:"total_lost_revenue"`
	RecoveredRevenue   float64   `json:"recovered_revenue"`
	FirstReminderSent  int64     `json:"first_reminder_sent"`
	SecondReminderSent int64     `json:"second_reminder_sent"`
	FinalReminderSent  int64     `json:"final_reminder_sent"`
	Since              time.Time `json:"since"`
	Until              time.Time `json:"until"`
}

type AbandonedCartResponse struct {
	ID                 uuid.UUID  `json:"id"`
	UserID             uuid.UUID  `json:"user_id"`
	UserEmail          string     `json:"user_email"`
	UserName           string     `json:"user_name"`
	ItemCount          int        `json:"item_count"`
	Total              float64    `json:"total"`
	AbandonedAt        *time.Time `json:"abandoned_at"`
	FirstReminderSent  *time.Time `json:"first_reminder_sent"`
	SecondReminderSent *time.Time `json:"second_reminder_sent"`
	FinalReminderSent  *time.Time `json:"final_reminder_sent"`
	RecoveredAt        *time.Time `json:"recovered_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
