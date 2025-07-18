package database

import (
	"context"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type emailRepository struct {
	db *gorm.DB
}

// NewEmailRepository creates a new email repository
func NewEmailRepository(db *gorm.DB) repositories.EmailRepository {
	return &emailRepository{db: db}
}

// Create creates a new email
func (r *emailRepository) Create(ctx context.Context, email *entities.Email) error {
	return r.db.WithContext(ctx).Create(email).Error
}

// GetByID gets an email by ID
func (r *emailRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Email, error) {
	var email entities.Email
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&email).Error
	if err != nil {
		return nil, err
	}
	return &email, nil
}

// GetByExternalID gets an email by external ID
func (r *emailRepository) GetByExternalID(ctx context.Context, externalID string) (*entities.Email, error) {
	var email entities.Email
	err := r.db.WithContext(ctx).Where("external_id = ?", externalID).First(&email).Error
	if err != nil {
		return nil, err
	}
	return &email, nil
}

// Update updates an email
func (r *emailRepository) Update(ctx context.Context, email *entities.Email) error {
	return r.db.WithContext(ctx).Save(email).Error
}

// Delete deletes an email
func (r *emailRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Email{}, id).Error
}

// List lists emails with pagination
func (r *emailRepository) List(ctx context.Context, offset, limit int) ([]*entities.Email, error) {
	var emails []*entities.Email
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&emails).Error
	return emails, err
}

// GetByUserID gets emails by user ID
func (r *emailRepository) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entities.Email, error) {
	var emails []*entities.Email
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&emails).Error
	return emails, err
}

// GetByOrderID gets emails by order ID
func (r *emailRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entities.Email, error) {
	var emails []*entities.Email
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&emails).Error
	return emails, err
}

// GetByType gets emails by type
func (r *emailRepository) GetByType(ctx context.Context, emailType entities.EmailType, offset, limit int) ([]*entities.Email, error) {
	var emails []*entities.Email
	err := r.db.WithContext(ctx).
		Where("type = ?", emailType).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&emails).Error
	return emails, err
}

// GetByStatus gets emails by status
func (r *emailRepository) GetByStatus(ctx context.Context, status entities.EmailStatus, offset, limit int) ([]*entities.Email, error) {
	var emails []*entities.Email
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&emails).Error
	return emails, err
}

// GetRetryableEmails gets emails that can be retried
func (r *emailRepository) GetRetryableEmails(ctx context.Context) ([]*entities.Email, error) {
	var emails []*entities.Email
	err := r.db.WithContext(ctx).
		Where("status = ? AND retry_count < max_retries AND next_retry_at <= ?", 
			entities.EmailStatusFailed, time.Now()).
		Order("created_at ASC").
		Find(&emails).Error
	return emails, err
}

// GetFailedEmails gets failed emails since a specific time
func (r *emailRepository) GetFailedEmails(ctx context.Context, since time.Time) ([]*entities.Email, error) {
	var emails []*entities.Email
	err := r.db.WithContext(ctx).
		Where("status = ? AND created_at >= ?", entities.EmailStatusFailed, since).
		Order("created_at DESC").
		Find(&emails).Error
	return emails, err
}

// GetEmailStats gets email statistics
func (r *emailRepository) GetEmailStats(ctx context.Context, since time.Time) (*repositories.EmailStats, error) {
	stats := &repositories.EmailStats{
		Since: since,
		Until: time.Now(),
		TypeStats: make(map[entities.EmailType]repositories.TypeStats),
	}

	// Get overall stats
	err := r.db.WithContext(ctx).
		Model(&entities.Email{}).
		Where("created_at >= ?", since).
		Select(`
			COUNT(*) as total_sent,
			SUM(CASE WHEN status = 'delivered' THEN 1 ELSE 0 END) as total_delivered,
			SUM(CASE WHEN is_opened = true THEN 1 ELSE 0 END) as total_opened,
			SUM(CASE WHEN is_clicked = true THEN 1 ELSE 0 END) as total_clicked,
			SUM(CASE WHEN is_bounced = true THEN 1 ELSE 0 END) as total_bounced,
			SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as total_failed
		`).
		Scan(stats).Error

	if err != nil {
		return nil, err
	}

	// Calculate rates
	if stats.TotalSent > 0 {
		stats.DeliveryRate = float64(stats.TotalDelivered) / float64(stats.TotalSent)
		stats.BounceRate = float64(stats.TotalBounced) / float64(stats.TotalSent)
		stats.FailureRate = float64(stats.TotalFailed) / float64(stats.TotalSent)
	}

	if stats.TotalDelivered > 0 {
		stats.OpenRate = float64(stats.TotalOpened) / float64(stats.TotalDelivered)
		stats.ClickRate = float64(stats.TotalClicked) / float64(stats.TotalDelivered)
	}

	return stats, nil
}

// GetDeliveryRate gets delivery rate for a specific email type
func (r *emailRepository) GetDeliveryRate(ctx context.Context, emailType entities.EmailType, since time.Time) (float64, error) {
	var result struct {
		TotalSent      int64
		TotalDelivered int64
	}

	err := r.db.WithContext(ctx).
		Model(&entities.Email{}).
		Where("type = ? AND created_at >= ?", emailType, since).
		Select(`
			COUNT(*) as total_sent,
			SUM(CASE WHEN status = 'delivered' THEN 1 ELSE 0 END) as total_delivered
		`).
		Scan(&result).Error

	if err != nil || result.TotalSent == 0 {
		return 0, err
	}

	return float64(result.TotalDelivered) / float64(result.TotalSent), nil
}

// GetOpenRate gets open rate for a specific email type
func (r *emailRepository) GetOpenRate(ctx context.Context, emailType entities.EmailType, since time.Time) (float64, error) {
	var result struct {
		TotalDelivered int64
		TotalOpened    int64
	}

	err := r.db.WithContext(ctx).
		Model(&entities.Email{}).
		Where("type = ? AND created_at >= ? AND status = 'delivered'", emailType, since).
		Select(`
			COUNT(*) as total_delivered,
			SUM(CASE WHEN is_opened = true THEN 1 ELSE 0 END) as total_opened
		`).
		Scan(&result).Error

	if err != nil || result.TotalDelivered == 0 {
		return 0, err
	}

	return float64(result.TotalOpened) / float64(result.TotalDelivered), nil
}

// GetClickRate gets click rate for a specific email type
func (r *emailRepository) GetClickRate(ctx context.Context, emailType entities.EmailType, since time.Time) (float64, error) {
	var result struct {
		TotalDelivered int64
		TotalClicked   int64
	}

	err := r.db.WithContext(ctx).
		Model(&entities.Email{}).
		Where("type = ? AND created_at >= ? AND status = 'delivered'", emailType, since).
		Select(`
			COUNT(*) as total_delivered,
			SUM(CASE WHEN is_clicked = true THEN 1 ELSE 0 END) as total_clicked
		`).
		Scan(&result).Error

	if err != nil || result.TotalDelivered == 0 {
		return 0, err
	}

	return float64(result.TotalClicked) / float64(result.TotalDelivered), nil
}

// CreateBatch creates multiple emails in a batch
func (r *emailRepository) CreateBatch(ctx context.Context, emails []*entities.Email) error {
	return r.db.WithContext(ctx).CreateInBatches(emails, 100).Error
}

// UpdateBatch updates multiple emails in a batch
func (r *emailRepository) UpdateBatch(ctx context.Context, emails []*entities.Email) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, email := range emails {
			if err := tx.Save(email).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Search searches emails based on query parameters
func (r *emailRepository) Search(ctx context.Context, query repositories.EmailSearchQuery) ([]*entities.Email, int, error) {
	var emails []*entities.Email
	var total int64

	db := r.db.WithContext(ctx).Model(&entities.Email{})

	// Apply filters
	if query.UserID != nil {
		db = db.Where("user_id = ?", *query.UserID)
	}
	if query.OrderID != nil {
		db = db.Where("order_id = ?", *query.OrderID)
	}
	if query.Type != nil {
		db = db.Where("type = ?", *query.Type)
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if query.Priority != nil {
		db = db.Where("priority = ?", *query.Priority)
	}
	if query.ToEmail != "" {
		db = db.Where("to_email ILIKE ?", "%"+query.ToEmail+"%")
	}
	if query.Subject != "" {
		db = db.Where("subject ILIKE ?", "%"+query.Subject+"%")
	}
	if query.CreatedAfter != nil {
		db = db.Where("created_at >= ?", *query.CreatedAfter)
	}
	if query.CreatedBefore != nil {
		db = db.Where("created_at <= ?", *query.CreatedBefore)
	}
	if query.SentAfter != nil {
		db = db.Where("sent_at >= ?", *query.SentAfter)
	}
	if query.SentBefore != nil {
		db = db.Where("sent_at <= ?", *query.SentBefore)
	}
	if query.IsDelivered {
		db = db.Where("status = 'delivered'")
	}
	if query.IsOpened {
		db = db.Where("is_opened = true")
	}
	if query.IsClicked {
		db = db.Where("is_clicked = true")
	}
	if query.IsBounced {
		db = db.Where("is_bounced = true")
	}
	if query.HasFailed {
		db = db.Where("status = 'failed'")
	}
	if query.CanRetry {
		db = db.Where("status = 'failed' AND retry_count < max_retries")
	}

	// Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := query.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := query.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}
	db = db.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	// Apply pagination
	if query.Limit > 0 {
		db = db.Limit(query.Limit)
	}
	if query.Offset > 0 {
		db = db.Offset(query.Offset)
	}

	// Execute query
	err := db.Find(&emails).Error
	return emails, int(total), err
}
