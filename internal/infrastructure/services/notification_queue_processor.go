package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/usecases"
)

// NotificationQueueProcessor handles background processing of notification queue
type NotificationQueueProcessor struct {
	notificationRepo repositories.NotificationRepository
	notificationUC   usecases.NotificationUseCase
	workers          int
	batchSize        int
	pollInterval     time.Duration
	retryInterval    time.Duration
	maxRetries       int
	stopChan         chan struct{}
	wg               sync.WaitGroup
	running          bool
	mu               sync.RWMutex
}

// NewNotificationQueueProcessor creates a new notification queue processor
func NewNotificationQueueProcessor(
	notificationRepo repositories.NotificationRepository,
	notificationUC usecases.NotificationUseCase,
	workers int,
	batchSize int,
	pollInterval time.Duration,
	retryInterval time.Duration,
	maxRetries int,
) *NotificationQueueProcessor {
	if workers <= 0 {
		workers = 3
	}
	if batchSize <= 0 {
		batchSize = 10
	}
	if pollInterval <= 0 {
		pollInterval = 30 * time.Second
	}
	if retryInterval <= 0 {
		retryInterval = 5 * time.Minute
	}
	if maxRetries <= 0 {
		maxRetries = 3
	}

	return &NotificationQueueProcessor{
		notificationRepo: notificationRepo,
		notificationUC:   notificationUC,
		workers:          workers,
		batchSize:        batchSize,
		pollInterval:     pollInterval,
		retryInterval:    retryInterval,
		maxRetries:       maxRetries,
		stopChan:         make(chan struct{}),
	}
}

// Start starts the notification queue processor
func (p *NotificationQueueProcessor) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return fmt.Errorf("notification queue processor is already running")
	}

	p.running = true
	log.Printf("Starting notification queue processor with %d workers", p.workers)

	// Start worker goroutines
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}

	// Start retry processor
	p.wg.Add(1)
	go p.retryProcessor(ctx)

	return nil
}

// Stop stops the notification queue processor
func (p *NotificationQueueProcessor) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return fmt.Errorf("notification queue processor is not running")
	}

	log.Println("Stopping notification queue processor...")
	close(p.stopChan)
	p.wg.Wait()
	p.running = false
	log.Println("Notification queue processor stopped")

	return nil
}

// IsRunning returns whether the processor is running
func (p *NotificationQueueProcessor) IsRunning() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.running
}

// worker processes notifications from the queue
func (p *NotificationQueueProcessor) worker(ctx context.Context, workerID int) {
	defer p.wg.Done()
	log.Printf("Worker %d started", workerID)

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopped due to context cancellation", workerID)
			return
		case <-p.stopChan:
			log.Printf("Worker %d stopped", workerID)
			return
		case <-ticker.C:
			p.processBatch(ctx, workerID)
		}
	}
}

// processBatch processes a batch of pending notifications
func (p *NotificationQueueProcessor) processBatch(ctx context.Context, workerID int) {
	// Get pending notifications
	notifications, err := p.notificationRepo.GetPendingNotificationsForQueue(ctx, p.batchSize)
	if err != nil {
		log.Printf("Worker %d: Failed to get pending notifications: %v", workerID, err)
		return
	}

	if len(notifications) == 0 {
		return
	}

	log.Printf("Worker %d: Processing %d notifications", workerID, len(notifications))

	for _, notification := range notifications {
		select {
		case <-ctx.Done():
			return
		case <-p.stopChan:
			return
		default:
			p.processNotification(ctx, workerID, notification)
		}
	}
}

// processNotification processes a single notification
func (p *NotificationQueueProcessor) processNotification(ctx context.Context, workerID int, notification *entities.Notification) {
	log.Printf("Worker %d: Processing notification %s (type: %s)", workerID, notification.ID, notification.Type)

	// Try to claim the notification atomically (prevent race condition between workers)
	if !p.claimNotification(ctx, notification.ID, workerID) {
		log.Printf("Worker %d: Notification %s already claimed by another worker", workerID, notification.ID)
		return
	}

	// Send notification
	err := p.notificationUC.SendNotification(ctx, notification)
	if err != nil {
		log.Printf("Worker %d: Failed to send notification %s: %v", workerID, notification.ID, err)
		p.handleFailedNotification(ctx, notification, err)
		return
	}

	log.Printf("Worker %d: Successfully sent notification %s", workerID, notification.ID)
}

// claimNotification atomically claims a notification for processing
func (p *NotificationQueueProcessor) claimNotification(ctx context.Context, notificationID uuid.UUID, workerID int) bool {
	// Get the notification first to update it
	notification, err := p.notificationRepo.GetByID(ctx, notificationID)
	if err != nil {
		log.Printf("Worker %d: Failed to get notification %s: %v", workerID, notificationID, err)
		return false
	}

	// Check if it's still pending
	if notification.Status != entities.NotificationStatusPending {
		return false
	}

	// Try to update the notification status from pending to processing atomically
	notification.Status = entities.NotificationStatusProcessing
	notification.UpdatedAt = time.Now()

	if err := p.notificationRepo.Update(ctx, notification); err != nil {
		log.Printf("Worker %d: Failed to claim notification %s: %v", workerID, notificationID, err)
		return false
	}

	log.Printf("Worker %d: Successfully claimed notification %s", workerID, notificationID)
	return true
}

// handleFailedNotification handles a failed notification
func (p *NotificationQueueProcessor) handleFailedNotification(ctx context.Context, notification *entities.Notification, err error) {
	notification.RetryCount++
	notification.ErrorMessage = err.Error()
	notification.UpdatedAt = time.Now()

	if notification.RetryCount >= p.maxRetries {
		notification.Status = entities.NotificationStatusFailed
		log.Printf("Notification %s failed permanently after %d retries", notification.ID, notification.RetryCount)
	} else {
		notification.Status = entities.NotificationStatusPending
		notification.NextRetryAt = &[]time.Time{time.Now().Add(p.retryInterval)}[0]
		log.Printf("Notification %s will be retried (attempt %d/%d)", notification.ID, notification.RetryCount, p.maxRetries)
	}

	if err := p.notificationRepo.Update(ctx, notification); err != nil {
		log.Printf("Failed to update failed notification %s: %v", notification.ID, err)
	}
}

// retryProcessor handles retrying failed notifications
func (p *NotificationQueueProcessor) retryProcessor(ctx context.Context) {
	defer p.wg.Done()
	log.Println("Retry processor started")

	ticker := time.NewTicker(p.retryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Retry processor stopped due to context cancellation")
			return
		case <-p.stopChan:
			log.Println("Retry processor stopped")
			return
		case <-ticker.C:
			p.processRetries(ctx)
		}
	}
}

// processRetries processes notifications that are ready for retry
func (p *NotificationQueueProcessor) processRetries(ctx context.Context) {
	notifications, err := p.notificationRepo.GetRetryableNotifications(ctx, p.batchSize)
	if err != nil {
		log.Printf("Retry processor: Failed to get retryable notifications: %v", err)
		return
	}

	if len(notifications) == 0 {
		return
	}

	log.Printf("Retry processor: Found %d notifications ready for retry", len(notifications))

	for _, notification := range notifications {
		select {
		case <-ctx.Done():
			return
		case <-p.stopChan:
			return
		default:
			// Reset status to pending so workers can pick it up
			notification.Status = entities.NotificationStatusPending
			notification.NextRetryAt = nil
			notification.UpdatedAt = time.Now()

			if err := p.notificationRepo.Update(ctx, notification); err != nil {
				log.Printf("Retry processor: Failed to reset notification %s status: %v", notification.ID, err)
			} else {
				log.Printf("Retry processor: Reset notification %s for retry", notification.ID)
			}
		}
	}
}

// GetStats returns processor statistics
func (p *NotificationQueueProcessor) GetStats(ctx context.Context) (map[string]interface{}, error) {
	pendingCount, err := p.notificationRepo.GetPendingCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending count: %w", err)
	}

	processingCount, err := p.notificationRepo.GetProcessingCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get processing count: %w", err)
	}

	failedCount, err := p.notificationRepo.GetFailedCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get failed count: %w", err)
	}

	return map[string]interface{}{
		"running":          p.IsRunning(),
		"workers":          p.workers,
		"batch_size":       p.batchSize,
		"poll_interval":    p.pollInterval.String(),
		"retry_interval":   p.retryInterval.String(),
		"max_retries":      p.maxRetries,
		"pending_count":    pendingCount,
		"processing_count": processingCount,
		"failed_count":     failedCount,
	}, nil
}
