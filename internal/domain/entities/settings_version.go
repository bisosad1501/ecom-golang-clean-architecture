package entities

import (
	"time"

	"github.com/google/uuid"
)

// SettingsVersion represents a version of settings
type SettingsVersion struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Version     int       `json:"version" gorm:"not null"`
	Category    string    `json:"category" gorm:"not null;index"`
	Settings    string    `json:"settings" gorm:"type:text;not null"` // JSON of settings
	Checksum    string    `json:"checksum" gorm:"not null"`
	Description string    `json:"description"`
	CreatedBy   uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt   time.Time `json:"created_at"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
}

// TableName returns the table name for SettingsVersion
func (SettingsVersion) TableName() string {
	return "settings_versions"
}

// SettingsSnapshot represents a complete snapshot of all settings at a point in time
type SettingsSnapshot struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Version     string    `json:"version" gorm:"not null"`
	Settings    string    `json:"settings" gorm:"type:text;not null"` // JSON of all settings
	Checksum    string    `json:"checksum" gorm:"not null"`
	CreatedBy   uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt   time.Time `json:"created_at"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Tags        string    `json:"tags" gorm:"type:text"` // JSON array of tags
}

// TableName returns the table name for SettingsSnapshot
func (SettingsSnapshot) TableName() string {
	return "settings_snapshots"
}

// SettingsChangeLog represents a log of settings changes
type SettingsChangeLog struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Category    string    `json:"category" gorm:"not null;index"`
	SettingKey  string    `json:"setting_key" gorm:"not null;index"`
	OldValue    string    `json:"old_value" gorm:"type:text"`
	NewValue    string    `json:"new_value" gorm:"type:text"`
	ChangeType  string    `json:"change_type" gorm:"not null"` // CREATE, UPDATE, DELETE
	Reason      string    `json:"reason"`
	ChangedBy   uuid.UUID `json:"changed_by" gorm:"type:uuid;not null"`
	ChangedAt   time.Time `json:"changed_at"`
	VersionFrom int       `json:"version_from"`
	VersionTo   int       `json:"version_to"`
}

// TableName returns the table name for SettingsChangeLog
func (SettingsChangeLog) TableName() string {
	return "settings_change_logs"
}

// SettingsTemplate represents a template for settings
type SettingsTemplate struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null;unique"`
	Description string    `json:"description"`
	Category    string    `json:"category" gorm:"not null;index"`
	Template    string    `json:"template" gorm:"type:text;not null"` // JSON template
	Variables   string    `json:"variables" gorm:"type:text"`         // JSON of template variables
	IsPublic    bool      `json:"is_public" gorm:"default:false"`
	CreatedBy   uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UsageCount  int       `json:"usage_count" gorm:"default:0"`
}

// TableName returns the table name for SettingsTemplate
func (SettingsTemplate) TableName() string {
	return "settings_templates"
}

// SettingsEnvironment represents different environments for settings
type SettingsEnvironment struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null;unique"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	IsDefault   bool      `json:"is_default" gorm:"default:false"`
	Settings    string    `json:"settings" gorm:"type:text"` // JSON of environment-specific settings
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the table name for SettingsEnvironment
func (SettingsEnvironment) TableName() string {
	return "settings_environments"
}

// SettingsValidationRule represents validation rules for settings
type SettingsValidationRule struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SettingKey  string    `json:"setting_key" gorm:"not null;index"`
	Category    string    `json:"category" gorm:"not null;index"`
	RuleType    string    `json:"rule_type" gorm:"not null"` // required, regex, range, enum, etc.
	RuleValue   string    `json:"rule_value" gorm:"not null"`
	ErrorMessage string   `json:"error_message"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Priority    int       `json:"priority" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the table name for SettingsValidationRule
func (SettingsValidationRule) TableName() string {
	return "settings_validation_rules"
}

// SettingsChangeType represents the type of change
type SettingsChangeType string

const (
	SettingsChangeTypeCreate SettingsChangeType = "CREATE"
	SettingsChangeTypeUpdate SettingsChangeType = "UPDATE"
	SettingsChangeTypeDelete SettingsChangeType = "DELETE"
)

// SettingsValidationRuleType represents the type of validation rule
type SettingsValidationRuleType string

const (
	ValidationRuleRequired SettingsValidationRuleType = "required"
	ValidationRuleRegex    SettingsValidationRuleType = "regex"
	ValidationRuleRange    SettingsValidationRuleType = "range"
	ValidationRuleEnum     SettingsValidationRuleType = "enum"
	ValidationRuleEmail    SettingsValidationRuleType = "email"
	ValidationRuleURL      SettingsValidationRuleType = "url"
	ValidationRuleNumeric  SettingsValidationRuleType = "numeric"
	ValidationRuleBoolean  SettingsValidationRuleType = "boolean"
	ValidationRuleJSON     SettingsValidationRuleType = "json"
	ValidationRuleCustom   SettingsValidationRuleType = "custom"
)

// VersionedSetting represents a setting with version information
type VersionedSetting struct {
	Setting
	Version     int       `json:"version"`
	PreviousValue string  `json:"previous_value,omitempty"`
	ChangeType  string    `json:"change_type"`
	ChangedBy   uuid.UUID `json:"changed_by"`
	ChangedAt   time.Time `json:"changed_at"`
}

// SettingsDiff represents the difference between two versions
type SettingsDiff struct {
	Category    string                    `json:"category"`
	FromVersion int                       `json:"from_version"`
	ToVersion   int                       `json:"to_version"`
	Changes     []SettingsChangeDetail    `json:"changes"`
	Summary     SettingsDiffSummary       `json:"summary"`
}

// SettingsDiffSummary provides a summary of changes
type SettingsDiffSummary struct {
	TotalChanges int `json:"total_changes"`
	Created      int `json:"created"`
	Updated      int `json:"updated"`
	Deleted      int `json:"deleted"`
}

// SettingsRollbackInfo contains information needed for rollback
type SettingsRollbackInfo struct {
	TargetVersion int                    `json:"target_version"`
	Category      string                 `json:"category"`
	Changes       []SettingsChangeDetail `json:"changes"`
	CanRollback   bool                   `json:"can_rollback"`
	Warnings      []string               `json:"warnings,omitempty"`
}

// SettingsExportFormat represents export format options
type SettingsExportFormat string

const (
	ExportFormatJSON SettingsExportFormat = "json"
	ExportFormatYAML SettingsExportFormat = "yaml"
	ExportFormatTOML SettingsExportFormat = "toml"
	ExportFormatENV  SettingsExportFormat = "env"
	ExportFormatSQL  SettingsExportFormat = "sql"
)

// SettingsImportResult represents the result of an import operation
type SettingsImportResult struct {
	Success      bool                   `json:"success"`
	TotalItems   int                    `json:"total_items"`
	Imported     int                    `json:"imported"`
	Skipped      int                    `json:"skipped"`
	Errors       []string               `json:"errors,omitempty"`
	Warnings     []string               `json:"warnings,omitempty"`
	Changes      []SettingsChangeDetail `json:"changes,omitempty"`
	BackupID     string                 `json:"backup_id,omitempty"`
}

// SettingsExportOptions represents options for exporting settings
type SettingsExportOptions struct {
	Format      SettingsExportFormat `json:"format"`
	Categories  []string             `json:"categories,omitempty"`
	IncludeAll  bool                 `json:"include_all"`
	Encrypted   bool                 `json:"encrypted"`
	Compressed  bool                 `json:"compressed"`
	Metadata    bool                 `json:"metadata"`
}

// SettingsImportOptions represents options for importing settings
type SettingsImportOptions struct {
	Format       SettingsExportFormat `json:"format"`
	Categories   []string             `json:"categories,omitempty"`
	Overwrite    bool                 `json:"overwrite"`
	CreateBackup bool                 `json:"create_backup"`
	Validate     bool                 `json:"validate"`
	DryRun       bool                 `json:"dry_run"`
}
