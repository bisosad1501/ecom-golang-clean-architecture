package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	domainServices "ecom-golang-clean-architecture/internal/domain/services"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditService handles audit logging
type AuditService struct {
	db *gorm.DB
}

// NewAuditService creates a new audit service
func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

// LogSettingsChange logs a settings change
func (s *AuditService) LogSettingsChange(ctx context.Context, userID uuid.UUID, userEmail, category string, changes []entities.SettingsChangeDetail, ipAddress, userAgent, sessionID string) error {
	// Create audit log entry
	auditLog := &entities.AuditLog{
		UserID:     &userID,
		Action:     string(entities.AuditActionUpdate),
		Resource:   string(entities.AuditResourceSettings),
		ResourceID: &category,
		Level:      entities.LogLevelInfo,
		Category:   entities.LogCategoryAdmin,
		Message:    fmt.Sprintf("Settings updated in category: %s", category),
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		SessionID:  &sessionID,
		Success:    true,
		CreatedAt:  time.Now(),
	}

	// Serialize changes as details
	if len(changes) > 0 {
		changesMap := make(map[string]interface{})
		for _, change := range changes {
			changesMap[change.Key] = map[string]interface{}{
				"old_value": change.OldValue,
				"new_value": change.NewValue,
				"type":      change.Type,
			}
		}
		auditLog.Details = changesMap
	}

	// Add metadata
	auditLog.Metadata = map[string]interface{}{
		"category":      category,
		"changes_count": len(changes),
		"user_email":    userEmail,
	}

	return s.db.WithContext(ctx).Create(auditLog).Error
}

// LogSecurityEvent logs a security event
func (s *AuditService) LogSecurityEvent(ctx context.Context, userID *uuid.UUID, eventType entities.SecurityEventType, severity entities.SecuritySeverity, description, ipAddress, userAgent, location string, metadata map[string]interface{}) error {
	event := &entities.SecurityEvent{
		UserID:      userID,
		EventType:   string(eventType),
		Severity:    string(severity),
		Description: description,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Location:    location,
		Resolved:    false,
		CreatedAt:   time.Now(),
	}

	if metadata != nil {
		if metadataJSON, err := json.Marshal(metadata); err == nil {
			event.Metadata = string(metadataJSON)
		}
	}

	return s.db.WithContext(ctx).Create(event).Error
}

// LogSystemEvent logs a system event
func (s *AuditService) LogSystemEvent(ctx context.Context, eventType entities.SystemEventType, category entities.SystemEventCategory, description string, status entities.EventStatus, duration int64, metadata map[string]interface{}) error {
	event := &entities.SystemEvent{
		EventType:   string(eventType),
		Category:    string(category),
		Description: description,
		Status:      string(status),
		Duration:    duration,
		CreatedAt:   time.Now(),
	}

	if metadata != nil {
		if metadataJSON, err := json.Marshal(metadata); err == nil {
			event.Metadata = string(metadataJSON)
		}
	}

	return s.db.WithContext(ctx).Create(event).Error
}

// LogUserAction logs a general user action
func (s *AuditService) LogUserAction(ctx context.Context, userID uuid.UUID, action entities.AuditAction, resource entities.AuditResource, resourceID string, oldValues, newValues map[string]interface{}, ipAddress, userAgent, sessionID string) error {
	auditLog := &entities.AuditLog{
		UserID:     &userID,
		Action:     string(action),
		Resource:   string(resource),
		ResourceID: &resourceID,
		Level:      entities.LogLevelInfo,
		Category:   entities.LogCategoryAdmin,
		Message:    fmt.Sprintf("User action: %s on %s", action, resource),
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		SessionID:  &sessionID,
		Success:    true,
		CreatedAt:  time.Now(),
	}

	// Combine old and new values in details
	if oldValues != nil || newValues != nil {
		details := make(map[string]interface{})
		if oldValues != nil {
			details["old_values"] = oldValues
		}
		if newValues != nil {
			details["new_values"] = newValues
		}
		auditLog.Details = details
	}

	return s.db.WithContext(ctx).Create(auditLog).Error
}

// GetAuditLogs retrieves audit logs with filtering
func (s *AuditService) GetAuditLogs(ctx context.Context, filters domainServices.AuditLogFilters) ([]*entities.AuditLog, int64, error) {
	query := s.db.WithContext(ctx).Model(&entities.AuditLog{})

	// Apply filters
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if filters.Action != "" {
		query = query.Where("action = ?", filters.Action)
	}
	if filters.Resource != "" {
		query = query.Where("resource = ?", filters.Resource)
	}
	if filters.ResourceID != "" {
		query = query.Where("resource_id = ?", filters.ResourceID)
	}
	if filters.IPAddress != "" {
		query = query.Where("ip_address = ?", filters.IPAddress)
	}
	if !filters.StartDate.IsZero() {
		query = query.Where("created_at >= ?", filters.StartDate)
	}
	if !filters.EndDate.IsZero() {
		query = query.Where("created_at <= ?", filters.EndDate)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	// Order by created_at desc
	query = query.Order("created_at DESC")

	var logs []*entities.AuditLog
	if err := query.Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetSecurityEvents retrieves security events with filtering
func (s *AuditService) GetSecurityEvents(ctx context.Context, filters domainServices.SecurityEventFilters) ([]*entities.SecurityEvent, int64, error) {
	query := s.db.WithContext(ctx).Model(&entities.SecurityEvent{})

	// Apply filters
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if filters.EventType != "" {
		query = query.Where("event_type = ?", filters.EventType)
	}
	if filters.Severity != "" {
		query = query.Where("severity = ?", filters.Severity)
	}
	if filters.Resolved != nil {
		query = query.Where("resolved = ?", *filters.Resolved)
	}
	if !filters.StartDate.IsZero() {
		query = query.Where("created_at >= ?", filters.StartDate)
	}
	if !filters.EndDate.IsZero() {
		query = query.Where("created_at <= ?", filters.EndDate)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	// Order by created_at desc
	query = query.Order("created_at DESC")

	var events []*entities.SecurityEvent
	if err := query.Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// ResolveSecurityEvent marks a security event as resolved
func (s *AuditService) ResolveSecurityEvent(ctx context.Context, eventID, resolvedBy uuid.UUID) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&entities.SecurityEvent{}).
		Where("id = ?", eventID).
		Updates(map[string]interface{}{
			"resolved":    true,
			"resolved_by": resolvedBy,
			"resolved_at": &now,
		}).Error
}

// GetAuditStats returns audit statistics
func (s *AuditService) GetAuditStats(ctx context.Context, days int) (map[string]interface{}, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	stats := make(map[string]interface{})

	// Total audit logs
	var totalLogs int64
	s.db.WithContext(ctx).Model(&entities.AuditLog{}).
		Where("created_at >= ?", startDate).
		Count(&totalLogs)
	stats["total_logs"] = totalLogs

	// Logs by action
	var actionStats []struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}
	s.db.WithContext(ctx).Model(&entities.AuditLog{}).
		Select("action, COUNT(*) as count").
		Where("created_at >= ?", startDate).
		Group("action").
		Find(&actionStats)
	stats["by_action"] = actionStats

	// Logs by resource
	var resourceStats []struct {
		Resource string `json:"resource"`
		Count    int64  `json:"count"`
	}
	s.db.WithContext(ctx).Model(&entities.AuditLog{}).
		Select("resource, COUNT(*) as count").
		Where("created_at >= ?", startDate).
		Group("resource").
		Find(&resourceStats)
	stats["by_resource"] = resourceStats

	// Security events
	var securityEvents int64
	s.db.WithContext(ctx).Model(&entities.SecurityEvent{}).
		Where("created_at >= ?", startDate).
		Count(&securityEvents)
	stats["security_events"] = securityEvents

	// Unresolved security events
	var unresolvedEvents int64
	s.db.WithContext(ctx).Model(&entities.SecurityEvent{}).
		Where("created_at >= ? AND resolved = false", startDate).
		Count(&unresolvedEvents)
	stats["unresolved_security_events"] = unresolvedEvents

	return stats, nil
}
