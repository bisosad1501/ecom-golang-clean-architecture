package repositories

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// AuditRepository defines audit repository interface
type AuditRepository interface {
	// Basic operations
	Create(ctx context.Context, log *entities.AuditLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.AuditLog, error)
	List(ctx context.Context, filters AuditFilters) ([]*entities.AuditLog, error)
	Count(ctx context.Context, filters AuditFilters) (int64, error)

	// User activity tracking
	LogUserAction(ctx context.Context, userID uuid.UUID, action, resource string, details map[string]interface{}) error
	GetUserActivity(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.AuditLog, error)
	GetUserActivityByDateRange(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]*entities.AuditLog, error)

	// System activity tracking
	LogSystemEvent(ctx context.Context, event, description string, metadata map[string]interface{}) error
	GetSystemLogs(ctx context.Context, filters SystemLogFilters) ([]*entities.AuditLog, error)
	GetSystemLogsByLevel(ctx context.Context, level entities.LogLevel, limit, offset int) ([]*entities.AuditLog, error)

	// Security tracking
	LogSecurityEvent(ctx context.Context, userID *uuid.UUID, event, description string, severity entities.SecuritySeverity, metadata map[string]interface{}) error
	GetSecurityLogs(ctx context.Context, filters SecurityLogFilters) ([]*entities.AuditLog, error)
	GetFailedLoginAttempts(ctx context.Context, userID uuid.UUID, since time.Time) ([]*entities.AuditLog, error)
	GetSuspiciousActivity(ctx context.Context, limit, offset int) ([]*entities.AuditLog, error)

	// Data change tracking
	LogDataChange(ctx context.Context, userID uuid.UUID, table, recordID string, action entities.DataAction, oldData, newData map[string]interface{}) error
	GetDataChanges(ctx context.Context, table, recordID string, limit, offset int) ([]*entities.AuditLog, error)
	GetDataChangesByUser(ctx context.Context, userID uuid.UUID, table string, limit, offset int) ([]*entities.AuditLog, error)

	// Admin operations
	GetAdminActions(ctx context.Context, filters AdminActionFilters) ([]*entities.AuditLog, error)
	GetCriticalEvents(ctx context.Context, since time.Time) ([]*entities.AuditLog, error)
	GetComplianceReport(ctx context.Context, from, to time.Time) (*ComplianceReport, error)

	// Cleanup operations
	DeleteOldLogs(ctx context.Context, olderThan time.Time) error
	ArchiveLogs(ctx context.Context, olderThan time.Time) error
	GetLogRetentionStats(ctx context.Context) (*LogRetentionStats, error)

	// Search and analytics
	SearchLogs(ctx context.Context, query string, filters SearchFilters) ([]*entities.AuditLog, error)
	GetActivitySummary(ctx context.Context, from, to time.Time) (*ActivitySummary, error)
	GetUserActivitySummary(ctx context.Context, userID uuid.UUID, from, to time.Time) (*UserActivitySummary, error)
}



















// UserActivityCount represents user activity count
type UserActivityCount struct {
	UserID    uuid.UUID `json:"user_id"`
	UserEmail string    `json:"user_email"`
	Count     int64     `json:"count"`
}

// ActionCount represents action count
type ActionCount struct {
	Action string `json:"action"`
	Count  int64  `json:"count"`
}

// SecurityIncident represents security incident
type SecurityIncident struct {
	Timestamp   time.Time                `json:"timestamp"`
	UserID      *uuid.UUID               `json:"user_id"`
	EventType   string                   `json:"event_type"`
	Severity    entities.SecuritySeverity `json:"severity"`
	Description string                   `json:"description"`
	IPAddress   string                   `json:"ip_address"`
	UserAgent   string                   `json:"user_agent"`
}

// DataAccessPattern represents data access pattern
type DataAccessPattern struct {
	Table       string    `json:"table"`
	Action      string    `json:"action"`
	Count       int64     `json:"count"`
	UniqueUsers int64     `json:"unique_users"`
	LastAccess  time.Time `json:"last_access"`
}

// HourlyActivityCount represents hourly activity count
type HourlyActivityCount struct {
	Hour  int   `json:"hour"`
	Count int64 `json:"count"`
}

// DailyActivityCount represents daily activity count
type DailyActivityCount struct {
	Date  time.Time `json:"date"`
	Count int64     `json:"count"`
}
