package entities

import (
	"time"
)

// CleanupMetrics represents metrics for cleanup operations
type CleanupMetrics struct {
	// Reservation cleanup metrics
	ExpiredReservationsFound    int           `json:"expired_reservations_found"`
	ExpiredReservationsCleaned  int           `json:"expired_reservations_cleaned"`
	ReservationCleanupDuration  time.Duration `json:"reservation_cleanup_duration"`
	ReservationCleanupErrors    int           `json:"reservation_cleanup_errors"`

	// Order cleanup metrics
	ExpiredOrdersFound    int           `json:"expired_orders_found"`
	ExpiredOrdersCleaned  int           `json:"expired_orders_cleaned"`
	OrderCleanupDuration  time.Duration `json:"order_cleanup_duration"`
	OrderCleanupErrors    int           `json:"order_cleanup_errors"`

	// Cart cleanup metrics
	ExpiredCartsFound    int           `json:"expired_carts_found"`
	ExpiredCartsCleaned  int           `json:"expired_carts_cleaned"`
	CartCleanupDuration  time.Duration `json:"cart_cleanup_duration"`
	CartCleanupErrors    int           `json:"cart_cleanup_errors"`

	// Overall metrics
	TotalCleanupDuration time.Duration `json:"total_cleanup_duration"`
	CleanupStartTime     time.Time     `json:"cleanup_start_time"`
	CleanupEndTime       time.Time     `json:"cleanup_end_time"`
	Success              bool          `json:"success"`
	ErrorMessage         string        `json:"error_message,omitempty"`
}

// NewCleanupMetrics creates a new cleanup metrics instance
func NewCleanupMetrics() *CleanupMetrics {
	return &CleanupMetrics{
		CleanupStartTime: time.Now(),
		Success:          true,
	}
}

// FinishCleanup marks the cleanup as finished and calculates total duration
func (cm *CleanupMetrics) FinishCleanup() {
	cm.CleanupEndTime = time.Now()
	cm.TotalCleanupDuration = cm.CleanupEndTime.Sub(cm.CleanupStartTime)
}

// MarkError marks the cleanup as failed with an error message
func (cm *CleanupMetrics) MarkError(err error) {
	cm.Success = false
	if err != nil {
		cm.ErrorMessage = err.Error()
	}
}

// AddReservationMetrics adds reservation cleanup metrics
func (cm *CleanupMetrics) AddReservationMetrics(found, cleaned, errors int, duration time.Duration) {
	cm.ExpiredReservationsFound = found
	cm.ExpiredReservationsCleaned = cleaned
	cm.ReservationCleanupErrors = errors
	cm.ReservationCleanupDuration = duration
}

// AddOrderMetrics adds order cleanup metrics
func (cm *CleanupMetrics) AddOrderMetrics(found, cleaned, errors int, duration time.Duration) {
	cm.ExpiredOrdersFound = found
	cm.ExpiredOrdersCleaned = cleaned
	cm.OrderCleanupErrors = errors
	cm.OrderCleanupDuration = duration
}

// AddCartMetrics adds cart cleanup metrics
func (cm *CleanupMetrics) AddCartMetrics(found, cleaned, errors int, duration time.Duration) {
	cm.ExpiredCartsFound = found
	cm.ExpiredCartsCleaned = cleaned
	cm.CartCleanupErrors = errors
	cm.CartCleanupDuration = duration
}

// GetSummary returns a summary string of the cleanup metrics
func (cm *CleanupMetrics) GetSummary() string {
	if !cm.Success {
		return "❌ Cleanup failed: " + cm.ErrorMessage
	}

	totalFound := cm.ExpiredReservationsFound + cm.ExpiredOrdersFound + cm.ExpiredCartsFound
	totalErrors := cm.ReservationCleanupErrors + cm.OrderCleanupErrors + cm.CartCleanupErrors

	if totalFound == 0 {
		return "✅ No expired items found"
	}

	summary := "🧹 Cleanup completed: "
	if cm.ExpiredReservationsFound > 0 {
		summary += "reservations(" + string(rune(cm.ExpiredReservationsCleaned)) + "/" + string(rune(cm.ExpiredReservationsFound)) + ") "
	}
	if cm.ExpiredOrdersFound > 0 {
		summary += "orders(" + string(rune(cm.ExpiredOrdersCleaned)) + "/" + string(rune(cm.ExpiredOrdersFound)) + ") "
	}
	if cm.ExpiredCartsFound > 0 {
		summary += "carts(" + string(rune(cm.ExpiredCartsCleaned)) + "/" + string(rune(cm.ExpiredCartsFound)) + ") "
	}

	if totalErrors > 0 {
		summary += "⚠️ " + string(rune(totalErrors)) + " errors"
	}

	return summary
}

// CleanupJobConfig represents configuration for cleanup jobs
type CleanupJobConfig struct {
	// Cleanup intervals
	ReservationCleanupInterval time.Duration `json:"reservation_cleanup_interval"`
	OrderCleanupInterval       time.Duration `json:"order_cleanup_interval"`
	CartCleanupInterval        time.Duration `json:"cart_cleanup_interval"`

	// Batch sizes
	ReservationBatchSize int `json:"reservation_batch_size"`
	OrderBatchSize       int `json:"order_batch_size"`
	CartBatchSize        int `json:"cart_batch_size"`

	// Timeouts
	ReservationTimeout time.Duration `json:"reservation_timeout"`
	OrderTimeout       time.Duration `json:"order_timeout"`
	CartTimeout        time.Duration `json:"cart_timeout"`

	// Retry settings
	MaxRetries    int           `json:"max_retries"`
	RetryInterval time.Duration `json:"retry_interval"`

	// Monitoring
	EnableMetrics bool `json:"enable_metrics"`
	EnableLogging bool `json:"enable_logging"`
}

// DefaultCleanupJobConfig returns default configuration for cleanup jobs
func DefaultCleanupJobConfig() *CleanupJobConfig {
	return &CleanupJobConfig{
		ReservationCleanupInterval: 5 * time.Minute,
		OrderCleanupInterval:       10 * time.Minute,
		CartCleanupInterval:        15 * time.Minute,

		ReservationBatchSize: 100,
		OrderBatchSize:       50,
		CartBatchSize:        100,

		ReservationTimeout: 30 * time.Minute,
		OrderTimeout:       24 * time.Hour,
		CartTimeout:        7 * 24 * time.Hour, // 7 days

		MaxRetries:    3,
		RetryInterval: 1 * time.Minute,

		EnableMetrics: true,
		EnableLogging: true,
	}
}
