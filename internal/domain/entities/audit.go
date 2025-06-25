package entities

import (
	"time"

	"github.com/google/uuid"
)

// LogLevel represents log severity level
type LogLevel string

const (
	LogLevelDebug   LogLevel = "debug"
	LogLevelInfo    LogLevel = "info"
	LogLevelWarning LogLevel = "warning"
	LogLevelError   LogLevel = "error"
	LogLevelCritical LogLevel = "critical"
)

// LogCategory represents log category
type LogCategory string

const (
	LogCategoryUser     LogCategory = "user"
	LogCategorySystem   LogCategory = "system"
	LogCategorySecurity LogCategory = "security"
	LogCategoryData     LogCategory = "data"
	LogCategoryAdmin    LogCategory = "admin"
	LogCategoryAPI      LogCategory = "api"
	LogCategoryAuth     LogCategory = "auth"
	LogCategoryOrder    LogCategory = "order"
	LogCategoryPayment  LogCategory = "payment"
	LogCategoryInventory LogCategory = "inventory"
)

// SecuritySeverity represents security event severity
type SecuritySeverity string

const (
	SecuritySeverityLow      SecuritySeverity = "low"
	SecuritySeverityMedium   SecuritySeverity = "medium"
	SecuritySeverityHigh     SecuritySeverity = "high"
	SecuritySeverityCritical SecuritySeverity = "critical"
)

// DataAction represents data change action
type DataAction string

const (
	DataActionCreate DataAction = "create"
	DataActionUpdate DataAction = "update"
	DataActionDelete DataAction = "delete"
	DataActionView   DataAction = "view"
	DataActionExport DataAction = "export"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID          uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      *uuid.UUID             `json:"user_id,omitempty" gorm:"type:uuid;index"`
	Action      string                 `json:"action" gorm:"not null;index"`
	Resource    string                 `json:"resource" gorm:"not null;index"`
	ResourceID  *string                `json:"resource_id,omitempty" gorm:"index"`
	Level       LogLevel               `json:"level" gorm:"not null;index"`
	Category    LogCategory            `json:"category" gorm:"not null;index"`
	Message     string                 `json:"message" gorm:"not null"`
	Details     map[string]interface{} `json:"details,omitempty" gorm:"type:jsonb"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" gorm:"type:jsonb"`
	IPAddress   string                 `json:"ip_address,omitempty" gorm:"index"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	SessionID   *string                `json:"session_id,omitempty" gorm:"index"`
	RequestID   *string                `json:"request_id,omitempty" gorm:"index"`
	Success     bool                   `json:"success" gorm:"default:true;index"`
	ErrorCode   *string                `json:"error_code,omitempty"`
	ErrorMessage *string               `json:"error_message,omitempty"`
	Duration    *int64                 `json:"duration,omitempty"` // in milliseconds
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime;index"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for AuditLog entity
func (AuditLog) TableName() string {
	return "audit_logs"
}

// SecurityLog represents a security-specific audit log
type SecurityLog struct {
	ID          uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      *uuid.UUID       `json:"user_id,omitempty" gorm:"type:uuid;index"`
	EventType   string           `json:"event_type" gorm:"not null;index"`
	Severity    SecuritySeverity `json:"severity" gorm:"not null;index"`
	Description string           `json:"description" gorm:"not null"`
	IPAddress   string           `json:"ip_address" gorm:"not null;index"`
	UserAgent   string           `json:"user_agent"`
	Location    string           `json:"location,omitempty"`
	Successful  bool             `json:"successful" gorm:"default:false;index"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt   time.Time        `json:"created_at" gorm:"autoCreateTime;index"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for SecurityLog entity
func (SecurityLog) TableName() string {
	return "security_logs"
}

// DataChangeLog represents a data change audit log
type DataChangeLog struct {
	ID        uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID              `json:"user_id" gorm:"type:uuid;not null;index"`
	Table     string                 `json:"table" gorm:"not null;index"`
	RecordID  string                 `json:"record_id" gorm:"not null;index"`
	Action    DataAction             `json:"action" gorm:"not null;index"`
	OldData   map[string]interface{} `json:"old_data,omitempty" gorm:"type:jsonb"`
	NewData   map[string]interface{} `json:"new_data,omitempty" gorm:"type:jsonb"`
	Changes   map[string]interface{} `json:"changes,omitempty" gorm:"type:jsonb"`
	Reason    string                 `json:"reason,omitempty"`
	CreatedAt time.Time              `json:"created_at" gorm:"autoCreateTime;index"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for DataChangeLog entity
func (DataChangeLog) TableName() string {
	return "data_change_logs"
}

// SystemLog represents a system-level audit log
type SystemLog struct {
	ID          uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Component   string                 `json:"component" gorm:"not null;index"`
	Level       LogLevel               `json:"level" gorm:"not null;index"`
	Category    LogCategory            `json:"category" gorm:"not null;index"`
	Event       string                 `json:"event" gorm:"not null;index"`
	Message     string                 `json:"message" gorm:"not null"`
	Details     map[string]interface{} `json:"details,omitempty" gorm:"type:jsonb"`
	ServerID    string                 `json:"server_id,omitempty" gorm:"index"`
	ProcessID   *int                   `json:"process_id,omitempty"`
	ThreadID    *int                   `json:"thread_id,omitempty"`
	MemoryUsage *int64                 `json:"memory_usage,omitempty"`
	CPUUsage    *float64               `json:"cpu_usage,omitempty"`
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime;index"`
}

// TableName returns the table name for SystemLog entity
func (SystemLog) TableName() string {
	return "system_logs"
}

// AdminActionLog represents admin-specific actions
type AdminActionLog struct {
	ID          uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AdminID     uuid.UUID              `json:"admin_id" gorm:"type:uuid;not null;index"`
	Action      string                 `json:"action" gorm:"not null;index"`
	Resource    string                 `json:"resource" gorm:"not null;index"`
	ResourceID  *string                `json:"resource_id,omitempty" gorm:"index"`
	TargetUserID *uuid.UUID            `json:"target_user_id,omitempty" gorm:"type:uuid;index"`
	Description string                 `json:"description" gorm:"not null"`
	Details     map[string]interface{} `json:"details,omitempty" gorm:"type:jsonb"`
	IPAddress   string                 `json:"ip_address" gorm:"not null;index"`
	UserAgent   string                 `json:"user_agent"`
	Successful  bool                   `json:"successful" gorm:"default:true;index"`
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime;index"`

	// Relationships
	Admin      *User `json:"admin,omitempty" gorm:"foreignKey:AdminID"`
	TargetUser *User `json:"target_user,omitempty" gorm:"foreignKey:TargetUserID"`
}

// TableName returns the table name for AdminActionLog entity
func (AdminActionLog) TableName() string {
	return "admin_action_logs"
}

// IsSecurityEvent checks if the audit log represents a security event
func (al *AuditLog) IsSecurityEvent() bool {
	return al.Category == LogCategorySecurity || 
		   al.Category == LogCategoryAuth ||
		   al.Action == "login" ||
		   al.Action == "logout" ||
		   al.Action == "failed_login" ||
		   al.Action == "password_change" ||
		   al.Action == "permission_denied"
}

// IsCritical checks if the audit log is critical
func (al *AuditLog) IsCritical() bool {
	return al.Level == LogLevelCritical || al.Level == LogLevelError
}

// GetSeverityScore returns a numeric severity score
func (ss SecuritySeverity) GetSeverityScore() int {
	switch ss {
	case SecuritySeverityLow:
		return 1
	case SecuritySeverityMedium:
		return 2
	case SecuritySeverityHigh:
		return 3
	case SecuritySeverityCritical:
		return 4
	default:
		return 0
	}
}

// GetLevelScore returns a numeric level score
func (ll LogLevel) GetLevelScore() int {
	switch ll {
	case LogLevelDebug:
		return 1
	case LogLevelInfo:
		return 2
	case LogLevelWarning:
		return 3
	case LogLevelError:
		return 4
	case LogLevelCritical:
		return 5
	default:
		return 0
	}
}

// AuditStats represents audit statistics
type AuditStats struct {
	TotalEvents      int64 `json:"total_events"`
	UniqueUsers      int64 `json:"unique_users"`
	FailedLogins     int64 `json:"failed_logins"`
	SuccessfulLogins int64 `json:"successful_logins"`
}
