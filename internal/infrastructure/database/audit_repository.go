package database

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type auditRepository struct {
	db *gorm.DB
}

// NewAuditRepository creates a new audit repository
func NewAuditRepository(db *gorm.DB) repositories.AuditRepository {
	return &auditRepository{db: db}
}

// Create creates a new audit log entry
func (r *auditRepository) Create(ctx context.Context, log *entities.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetByID gets an audit log by ID
func (r *auditRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.AuditLog, error) {
	var log entities.AuditLog
	err := r.db.WithContext(ctx).First(&log, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// List lists audit logs with filters
func (r *auditRepository) List(ctx context.Context, filters repositories.AuditFilters) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	query := r.db.WithContext(ctx)

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Action != "" {
		query = query.Where("action = ?", filters.Action)
	}

	if filters.Resource != "" {
		query = query.Where("resource = ?", filters.Resource)
	}

	if filters.ResourceID != nil {
		query = query.Where("resource_id = ?", *filters.ResourceID)
	}

	if filters.IPAddress != "" {
		query = query.Where("ip_address = ?", filters.IPAddress)
	}

	if filters.UserAgent != "" {
		query = query.Where("user_agent LIKE ?", "%"+filters.UserAgent+"%")
	}

	if filters.CreatedAfter != nil {
		query = query.Where("created_at >= ?", *filters.CreatedAfter)
	}

	if filters.CreatedBefore != nil {
		query = query.Where("created_at <= ?", *filters.CreatedBefore)
	}

	// Apply sorting
	switch filters.SortBy {
	case "created_at":
		if filters.SortOrder == "desc" {
			query = query.Order("created_at DESC")
		} else {
			query = query.Order("created_at ASC")
		}
	case "action":
		if filters.SortOrder == "desc" {
			query = query.Order("action DESC")
		} else {
			query = query.Order("action ASC")
		}
	case "resource":
		if filters.SortOrder == "desc" {
			query = query.Order("resource DESC")
		} else {
			query = query.Order("resource ASC")
		}
	default:
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&logs).Error
	return logs, err
}

// Count counts audit logs with filters
func (r *auditRepository) Count(ctx context.Context, filters repositories.AuditFilters) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.AuditLog{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Action != "" {
		query = query.Where("action = ?", filters.Action)
	}

	if filters.Resource != "" {
		query = query.Where("resource = ?", filters.Resource)
	}

	if filters.ResourceID != nil {
		query = query.Where("resource_id = ?", *filters.ResourceID)
	}

	if filters.IPAddress != "" {
		query = query.Where("ip_address = ?", filters.IPAddress)
	}

	if filters.UserAgent != "" {
		query = query.Where("user_agent LIKE ?", "%"+filters.UserAgent+"%")
	}

	if filters.CreatedAfter != nil {
		query = query.Where("created_at >= ?", *filters.CreatedAfter)
	}

	if filters.CreatedBefore != nil {
		query = query.Where("created_at <= ?", *filters.CreatedBefore)
	}

	err := query.Count(&count).Error
	return count, err
}

// GetByUser gets audit logs for a specific user
func (r *auditRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetByResource gets audit logs for a specific resource
func (r *auditRepository) GetByResource(ctx context.Context, resource string, resourceID uuid.UUID, limit, offset int) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	err := r.db.WithContext(ctx).
		Where("resource = ? AND resource_id = ?", resource, resourceID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetByAction gets audit logs for a specific action
func (r *auditRepository) GetByAction(ctx context.Context, action string, limit, offset int) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	err := r.db.WithContext(ctx).
		Where("action = ?", action).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetByDateRange gets audit logs within a date range
func (r *auditRepository) GetByDateRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	err := r.db.WithContext(ctx).
		Where("created_at BETWEEN ? AND ?", from, to).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetSecurityEvents gets security-related audit logs
func (r *auditRepository) GetSecurityEvents(ctx context.Context, limit, offset int) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	securityActions := []string{
		"login",
		"logout",
		"login_failed",
		"password_change",
		"password_reset",
		"account_locked",
		"account_unlocked",
		"permission_denied",
		"unauthorized_access",
	}
	
	err := r.db.WithContext(ctx).
		Where("action IN (?)", securityActions).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetFailedLogins gets failed login attempts
func (r *auditRepository) GetFailedLogins(ctx context.Context, ipAddress string, since time.Time) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	query := r.db.WithContext(ctx).
		Where("action = ? AND created_at >= ?", "login_failed", since)
	
	if ipAddress != "" {
		query = query.Where("ip_address = ?", ipAddress)
	}
	
	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// GetUserActivity gets recent activity for a user
func (r *auditRepository) GetUserActivity(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// DeleteOld deletes old audit logs
func (r *auditRepository) DeleteOld(ctx context.Context, olderThan time.Time) error {
	return r.db.WithContext(ctx).
		Delete(&entities.AuditLog{}, "created_at < ?", olderThan).Error
}

// CreateBulk creates multiple audit logs in batch
func (r *auditRepository) CreateBulk(ctx context.Context, logs []*entities.AuditLog) error {
	return r.db.WithContext(ctx).CreateInBatches(logs, 100).Error
}

// GetStats gets audit statistics
func (r *auditRepository) GetStats(ctx context.Context, from, to time.Time) (*entities.AuditStats, error) {
	var stats entities.AuditStats
	
	// Get total events
	err := r.db.WithContext(ctx).
		Model(&entities.AuditLog{}).
		Where("created_at BETWEEN ? AND ?", from, to).
		Count(&stats.TotalEvents).Error
	if err != nil {
		return nil, err
	}

	// Get unique users
	err = r.db.WithContext(ctx).
		Model(&entities.AuditLog{}).
		Select("COUNT(DISTINCT user_id)").
		Where("created_at BETWEEN ? AND ? AND user_id IS NOT NULL", from, to).
		Scan(&stats.UniqueUsers).Error
	if err != nil {
		return nil, err
	}

	// Get failed logins
	err = r.db.WithContext(ctx).
		Model(&entities.AuditLog{}).
		Where("action = ? AND created_at BETWEEN ? AND ?", "login_failed", from, to).
		Count(&stats.FailedLogins).Error
	if err != nil {
		return nil, err
	}

	// Get successful logins
	err = r.db.WithContext(ctx).
		Model(&entities.AuditLog{}).
		Where("action = ? AND created_at BETWEEN ? AND ?", "login", from, to).
		Count(&stats.SuccessfulLogins).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// ArchiveLogs archives old audit logs
func (r *auditRepository) ArchiveLogs(ctx context.Context, olderThan time.Time) error {
	// In a real implementation, this would move logs to an archive table or storage
	// For now, we'll just mark them as archived or delete them
	return r.db.WithContext(ctx).
		Delete(&entities.AuditLog{}, "created_at < ?", olderThan).Error
}

// DeleteOldLogs deletes old audit logs (alias for ArchiveLogs)
func (r *auditRepository) DeleteOldLogs(ctx context.Context, olderThan time.Time) error {
	return r.ArchiveLogs(ctx, olderThan)
}

// GetActivitySummary gets activity summary for a time period
func (r *auditRepository) GetActivitySummary(ctx context.Context, from, to time.Time) (*repositories.ActivitySummary, error) {
	var stats repositories.ActivitySummary
	query := r.db.WithContext(ctx).Model(&entities.AuditLog{})

	if !from.IsZero() {
		query = query.Where("created_at >= ?", from)
	}

	if !to.IsZero() {
		query = query.Where("created_at <= ?", to)
	}

	// Get total events
	err := query.Count(&stats.TotalEvents).Error
	if err != nil {
		return nil, err
	}

	// Get unique users
	err = query.Select("COUNT(DISTINCT user_id)").Scan(&stats.UniqueUsers).Error
	if err != nil {
		return nil, err
	}

	// Get failed logins
	failedQuery := r.db.WithContext(ctx).Model(&entities.AuditLog{}).
		Where("action = ? AND created_at BETWEEN ? AND ?", "login_failed", from, to)
	err = failedQuery.Count(&stats.FailedLogins).Error
	if err != nil {
		return nil, err
	}

	// Get successful logins
	successQuery := r.db.WithContext(ctx).Model(&entities.AuditLog{}).
		Where("action = ? AND created_at BETWEEN ? AND ?", "login", from, to)
	err = successQuery.Count(&stats.SuccessfulLogins).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetAdminActions gets admin actions from audit logs with filters
func (r *auditRepository) GetAdminActions(ctx context.Context, filters repositories.AdminActionFilters) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	query := r.db.WithContext(ctx).Model(&entities.AuditLog{})

	if filters.AdminID != nil {
		query = query.Where("user_id = ?", *filters.AdminID)
	}

	if filters.Action != "" {
		query = query.Where("action = ?", filters.Action)
	}

	if filters.Resource != "" {
		query = query.Where("resource = ?", filters.Resource)
	}

	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	// Filter for admin-type actions
	query = query.Where("action IN (?)", []string{"admin_login", "user_create", "user_update", "user_delete", "role_change", "system_config"})

	// Apply sorting
	sortBy := "created_at"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}
	sortOrder := "DESC"
	if filters.SortOrder != "" {
		sortOrder = filters.SortOrder
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&logs).Error
	return logs, err
}

// GetComplianceReport gets compliance report for audit
func (r *auditRepository) GetComplianceReport(ctx context.Context, from, to time.Time) (*repositories.ComplianceReport, error) {
	var report repositories.ComplianceReport

	// Get total events in period
	err := r.db.WithContext(ctx).
		Model(&entities.AuditLog{}).
		Where("created_at BETWEEN ? AND ?", from, to).
		Count(&report.TotalEvents).Error
	if err != nil {
		return nil, err
	}

	// Get security events
	err = r.db.WithContext(ctx).
		Model(&entities.AuditLog{}).
		Where("created_at BETWEEN ? AND ? AND category = ?", from, to, "security").
		Count(&report.SecurityEvents).Error
	if err != nil {
		return nil, err
	}

	// Get failed login attempts
	err = r.db.WithContext(ctx).
		Model(&entities.AuditLog{}).
		Where("created_at BETWEEN ? AND ? AND action = ?", from, to, "login_failed").
		Count(&report.FailedLogins).Error
	if err != nil {
		return nil, err
	}

	// Get data access events
	err = r.db.WithContext(ctx).
		Model(&entities.AuditLog{}).
		Where("created_at BETWEEN ? AND ? AND action IN (?)", from, to, []string{"data_access", "data_export", "data_view"}).
		Count(&report.DataAccessEvents).Error
	if err != nil {
		return nil, err
	}

	report.Period = from.Format("2006-01-02") + " to " + to.Format("2006-01-02")

	return &report, nil
}

// GetCriticalEvents gets critical events since a specific time
func (r *auditRepository) GetCriticalEvents(ctx context.Context, since time.Time) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	err := r.db.WithContext(ctx).
		Where("created_at >= ? AND level IN (?)", since, []string{"critical", "error"}).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// GetDataChanges gets data changes from audit logs (placeholder)
func (r *auditRepository) GetDataChanges(ctx context.Context, table, recordID string, limit, offset int) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	query := r.db.WithContext(ctx).Model(&entities.AuditLog{})

	if table != "" {
		query = query.Where("resource = ?", table)
	}

	if recordID != "" {
		query = query.Where("resource_id = ?", recordID)
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, err
}

// GetDataChangesByUser gets data changes by user from audit logs
func (r *auditRepository) GetDataChangesByUser(ctx context.Context, userID uuid.UUID, table string, limit, offset int) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	query := r.db.WithContext(ctx).Model(&entities.AuditLog{}).
		Where("user_id = ?", userID)

	if table != "" {
		query = query.Where("resource = ?", table)
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, err
}

// GetFailedLoginAttempts gets failed login attempts for a user
func (r *auditRepository) GetFailedLoginAttempts(ctx context.Context, userID uuid.UUID, since time.Time) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	err := r.db.WithContext(ctx).Model(&entities.AuditLog{}).
		Where("user_id = ? AND action = ? AND created_at >= ?", userID, "login_failed", since).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, err
}

// GetLogRetentionStats gets log retention statistics
func (r *auditRepository) GetLogRetentionStats(ctx context.Context) (*repositories.LogRetentionStats, error) {
	var stats repositories.LogRetentionStats

	// Get total logs
	err := r.db.WithContext(ctx).
		Model(&entities.AuditLog{}).
		Count(&stats.TotalLogs).Error
	if err != nil {
		return nil, err
	}

	// Get logs older than 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	err = r.db.WithContext(ctx).
		Model(&entities.AuditLog{}).
		Where("created_at < ?", thirtyDaysAgo).
		Count(&stats.LogsOlderThan30Days).Error
	if err != nil {
		return nil, err
	}

	// Get logs older than 90 days
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90)
	err = r.db.WithContext(ctx).
		Model(&entities.AuditLog{}).
		Where("created_at < ?", ninetyDaysAgo).
		Count(&stats.LogsOlderThan90Days).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetSecurityLogs gets security-related audit logs with filters
func (r *auditRepository) GetSecurityLogs(ctx context.Context, filters repositories.SecurityLogFilters) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	query := r.db.WithContext(ctx).Where("category = ?", "security")

	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	err := query.Order("created_at DESC").
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&logs).Error
	return logs, err
}

// GetSuspiciousActivity gets suspicious activity logs
func (r *auditRepository) GetSuspiciousActivity(ctx context.Context, limit, offset int) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	err := r.db.WithContext(ctx).
		Where("level = ? OR action IN (?)", "critical", []string{"login_failed", "unauthorized_access", "data_breach"}).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetSystemLogs gets system-related audit logs with filters
func (r *auditRepository) GetSystemLogs(ctx context.Context, filters repositories.SystemLogFilters) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	query := r.db.WithContext(ctx).Where("category = ?", "system")

	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	err := query.Order("created_at DESC").
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&logs).Error
	return logs, err
}

// GetSystemLogsByLevel gets system logs by level
func (r *auditRepository) GetSystemLogsByLevel(ctx context.Context, level entities.LogLevel, limit, offset int) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	err := r.db.WithContext(ctx).
		Where("category = ? AND level = ?", "system", level).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetUserActivityByDateRange gets user activity within date range
func (r *auditRepository) GetUserActivityByDateRange(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, from, to).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// GetUserActivitySummary gets user activity summary
func (r *auditRepository) GetUserActivitySummary(ctx context.Context, userID uuid.UUID, from, to time.Time) (*repositories.UserActivitySummary, error) {
	var summary repositories.UserActivitySummary

	// Get total activities
	err := r.db.WithContext(ctx).
		Model(&entities.AuditLog{}).
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, from, to).
		Count(&summary.TotalActivities).Error
	if err != nil {
		return nil, err
	}

	summary.UserID = userID
	return &summary, nil
}

// LogDataChange logs data change events
func (r *auditRepository) LogDataChange(ctx context.Context, userID uuid.UUID, table, recordID string, action entities.DataAction, oldData, newData map[string]interface{}) error {
	details := map[string]interface{}{
		"record_id": recordID,
		"old_data":  oldData,
		"new_data":  newData,
	}
	log := &entities.AuditLog{
		ID:       uuid.New(),
		UserID:   &userID,
		Action:   string(action),
		Resource: table,
		Details:  details,
		Level:    entities.LogLevelInfo,
		Category: entities.LogCategoryData,
		CreatedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Create(log).Error
}

// LogSecurityEvent logs security events
func (r *auditRepository) LogSecurityEvent(ctx context.Context, userID *uuid.UUID, event, description string, severity entities.SecuritySeverity, metadata map[string]interface{}) error {
	details := map[string]interface{}{
		"description": description,
		"severity":    severity,
		"metadata":    metadata,
	}
	log := &entities.AuditLog{
		ID:       uuid.New(),
		UserID:   userID,
		Action:   event,
		Resource: "security",
		Details:  details,
		Level:    entities.LogLevelWarning,
		Category: entities.LogCategorySecurity,
		CreatedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Create(log).Error
}

// LogSystemEvent logs system events
func (r *auditRepository) LogSystemEvent(ctx context.Context, event, description string, metadata map[string]interface{}) error {
	details := map[string]interface{}{
		"description": description,
		"metadata":    metadata,
	}
	log := &entities.AuditLog{
		ID:       uuid.New(),
		Action:   event,
		Resource: "system",
		Details:  details,
		Level:    entities.LogLevelInfo,
		Category: entities.LogCategorySystem,
		CreatedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Create(log).Error
}

// LogUserAction logs user actions
func (r *auditRepository) LogUserAction(ctx context.Context, userID uuid.UUID, action, resource string, details map[string]interface{}) error {
	log := &entities.AuditLog{
		ID:       uuid.New(),
		UserID:   &userID,
		Action:   action,
		Resource: resource,
		Details:  details,
		Level:    entities.LogLevelInfo,
		Category: entities.LogCategoryUser,
		CreatedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Create(log).Error
}

// SearchLogs searches audit logs
func (r *auditRepository) SearchLogs(ctx context.Context, query string, filters repositories.SearchFilters) ([]*entities.AuditLog, error) {
	var logs []*entities.AuditLog
	dbQuery := r.db.WithContext(ctx)

	if query != "" {
		dbQuery = dbQuery.Where("action LIKE ? OR resource LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if filters.DateFrom != nil {
		dbQuery = dbQuery.Where("created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		dbQuery = dbQuery.Where("created_at <= ?", *filters.DateTo)
	}

	err := dbQuery.Order("created_at DESC").
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&logs).Error
	return logs, err
}
