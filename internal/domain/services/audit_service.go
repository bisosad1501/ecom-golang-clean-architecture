package services

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"

	"github.com/google/uuid"
)

// AuditService interface for audit logging
type AuditService interface {
	// LogSettingsChange logs a settings change
	LogSettingsChange(ctx context.Context, userID uuid.UUID, userEmail, category string, changes []entities.SettingsChangeDetail, ipAddress, userAgent, sessionID string) error
	
	// LogSecurityEvent logs a security event
	LogSecurityEvent(ctx context.Context, userID *uuid.UUID, eventType entities.SecurityEventType, severity entities.SecuritySeverity, description, ipAddress, userAgent, location string, metadata map[string]interface{}) error
	
	// LogSystemEvent logs a system event
	LogSystemEvent(ctx context.Context, eventType entities.SystemEventType, category entities.SystemEventCategory, description string, status entities.EventStatus, duration int64, metadata map[string]interface{}) error
	
	// LogUserAction logs a general user action
	LogUserAction(ctx context.Context, userID uuid.UUID, action entities.AuditAction, resource entities.AuditResource, resourceID string, oldValues, newValues map[string]interface{}, ipAddress, userAgent, sessionID string) error
	
	// GetAuditLogs retrieves audit logs with filtering
	GetAuditLogs(ctx context.Context, filters AuditLogFilters) ([]*entities.AuditLog, int64, error)
	
	// GetSecurityEvents retrieves security events with filtering
	GetSecurityEvents(ctx context.Context, filters SecurityEventFilters) ([]*entities.SecurityEvent, int64, error)
	
	// ResolveSecurityEvent marks a security event as resolved
	ResolveSecurityEvent(ctx context.Context, eventID, resolvedBy uuid.UUID) error
	
	// GetAuditStats returns audit statistics
	GetAuditStats(ctx context.Context, days int) (map[string]interface{}, error)
}

// AuditLogFilters represents filters for audit log queries
type AuditLogFilters struct {
	UserID     *uuid.UUID `json:"user_id"`
	Action     string     `json:"action"`
	Resource   string     `json:"resource"`
	ResourceID string     `json:"resource_id"`
	IPAddress  string     `json:"ip_address"`
	StartDate  time.Time  `json:"start_date"`
	EndDate    time.Time  `json:"end_date"`
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset"`
}

// SecurityEventFilters represents filters for security event queries
type SecurityEventFilters struct {
	UserID    *uuid.UUID `json:"user_id"`
	EventType string     `json:"event_type"`
	Severity  string     `json:"severity"`
	Resolved  *bool      `json:"resolved"`
	StartDate time.Time  `json:"start_date"`
	EndDate   time.Time  `json:"end_date"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}
