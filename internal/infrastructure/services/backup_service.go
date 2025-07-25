package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
)

// BackupService handles settings backup and restore
type BackupService struct {
	settingsRepo repositories.SettingsRepository
	auditService *AuditService
	backupDir    string
}

// NewBackupService creates a new backup service
func NewBackupService(settingsRepo repositories.SettingsRepository, auditService *AuditService, backupDir string) *BackupService {
	return &BackupService{
		settingsRepo: settingsRepo,
		auditService: auditService,
		backupDir:    backupDir,
	}
}

// SettingsBackup represents a settings backup
type SettingsBackup struct {
	ID          string                         `json:"id"`
	Timestamp   time.Time                      `json:"timestamp"`
	Version     string                         `json:"version"`
	Description string                         `json:"description"`
	Settings    map[string][]*entities.Setting `json:"settings"` // Grouped by category
	Metadata    BackupMetadata                 `json:"metadata"`
}

// BackupMetadata contains backup metadata
type BackupMetadata struct {
	CreatedBy     string                 `json:"created_by"`
	CreatedByID   string                 `json:"created_by_id"`
	TotalSettings int                    `json:"total_settings"`
	Categories    []string               `json:"categories"`
	Size          int64                  `json:"size"`
	Checksum      string                 `json:"checksum"`
	Environment   string                 `json:"environment"`
	AppVersion    string                 `json:"app_version"`
	Extra         map[string]interface{} `json:"extra,omitempty"`
}

// CreateBackup creates a backup of all settings
func (s *BackupService) CreateBackup(ctx context.Context, description, createdBy, createdByID string) (*SettingsBackup, error) {
	// Get all settings categories
	categories := []string{"general", "store", "payment", "email", "tax", "shipping", "seo", "security", "notifications", "integrations"}

	backup := &SettingsBackup{
		ID:          uuid.New().String(),
		Timestamp:   time.Now(),
		Version:     "1.0",
		Description: description,
		Settings:    make(map[string][]*entities.Setting),
		Metadata: BackupMetadata{
			CreatedBy:   createdBy,
			CreatedByID: createdByID,
			Categories:  categories,
			Environment: os.Getenv("APP_ENV"),
			AppVersion:  "1.0.0", // This should come from build info
		},
	}

	totalSettings := 0

	// Backup each category
	for _, category := range categories {
		settings, err := s.settingsRepo.GetSettingsByCategory(ctx, category)
		if err != nil {
			return nil, fmt.Errorf("failed to get settings for category %s: %w", category, err)
		}

		backup.Settings[category] = settings
		totalSettings += len(settings)
	}

	backup.Metadata.TotalSettings = totalSettings

	// Save backup to file
	if err := s.saveBackupToFile(backup); err != nil {
		return nil, fmt.Errorf("failed to save backup to file: %w", err)
	}

	// Log audit event
	if s.auditService != nil {
		if userID, err := uuid.Parse(createdByID); err == nil {
			s.auditService.LogSystemEvent(
				ctx,
				entities.SystemEventBackupCreated,
				entities.SystemCategoryBackup,
				fmt.Sprintf("Settings backup created: %s", description),
				entities.EventStatusSuccess,
				0,
				map[string]interface{}{
					"backup_id":      backup.ID,
					"total_settings": totalSettings,
					"categories":     categories,
					"created_by":     createdBy,
					"user_id":        userID.String(),
				},
			)
		}
	}

	return backup, nil
}

// RestoreBackup restores settings from a backup
func (s *BackupService) RestoreBackup(ctx context.Context, backupID, restoredBy, restoredByID string) error {
	// Load backup from file
	backup, err := s.loadBackupFromFile(backupID)
	if err != nil {
		return fmt.Errorf("failed to load backup: %w", err)
	}

	// Create a backup of current settings before restore
	currentBackup, err := s.CreateBackup(ctx, fmt.Sprintf("Pre-restore backup before restoring %s", backupID), restoredBy, restoredByID)
	if err != nil {
		return fmt.Errorf("failed to create pre-restore backup: %w", err)
	}

	// Restore each category
	for category, settings := range backup.Settings {
		if err := s.settingsRepo.SetSettings(ctx, settings); err != nil {
			return fmt.Errorf("failed to restore settings for category %s: %w", category, err)
		}
	}

	// Log audit event
	if s.auditService != nil {
		if userID, err := uuid.Parse(restoredByID); err == nil {
			s.auditService.LogSystemEvent(
				ctx,
				entities.SystemEventBackupRestored,
				entities.SystemCategoryBackup,
				fmt.Sprintf("Settings restored from backup: %s", backup.Description),
				entities.EventStatusSuccess,
				0,
				map[string]interface{}{
					"backup_id":          backupID,
					"pre_restore_backup": currentBackup.ID,
					"restored_by":        restoredBy,
					"total_settings":     backup.Metadata.TotalSettings,
					"user_id":            userID.String(),
				},
			)
		}
	}

	return nil
}

// ListBackups lists all available backups
func (s *BackupService) ListBackups() ([]*SettingsBackup, error) {
	if err := s.ensureBackupDir(); err != nil {
		return nil, err
	}

	files, err := filepath.Glob(filepath.Join(s.backupDir, "settings-backup-*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list backup files: %w", err)
	}

	var backups []*SettingsBackup
	for _, file := range files {
		backup, err := s.loadBackupFromFilePath(file)
		if err != nil {
			continue // Skip invalid backup files
		}
		backups = append(backups, backup)
	}

	return backups, nil
}

// DeleteBackup deletes a backup
func (s *BackupService) DeleteBackup(backupID string) error {
	filename := fmt.Sprintf("settings-backup-%s.json", backupID)
	filepath := filepath.Join(s.backupDir, filename)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("backup not found: %s", backupID)
	}

	return os.Remove(filepath)
}

// GetBackup gets a specific backup
func (s *BackupService) GetBackup(backupID string) (*SettingsBackup, error) {
	return s.loadBackupFromFile(backupID)
}

// saveBackupToFile saves backup to a JSON file
func (s *BackupService) saveBackupToFile(backup *SettingsBackup) error {
	if err := s.ensureBackupDir(); err != nil {
		return err
	}

	filename := fmt.Sprintf("settings-backup-%s.json", backup.ID)
	filepath := filepath.Join(s.backupDir, filename)

	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup: %w", err)
	}

	backup.Metadata.Size = int64(len(data))

	// Re-marshal with size info
	data, err = json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup with metadata: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// loadBackupFromFile loads backup from a JSON file
func (s *BackupService) loadBackupFromFile(backupID string) (*SettingsBackup, error) {
	filename := fmt.Sprintf("settings-backup-%s.json", backupID)
	filepath := filepath.Join(s.backupDir, filename)

	return s.loadBackupFromFilePath(filepath)
}

// loadBackupFromFilePath loads backup from a file path
func (s *BackupService) loadBackupFromFilePath(filepath string) (*SettingsBackup, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup file: %w", err)
	}

	var backup SettingsBackup
	if err := json.Unmarshal(data, &backup); err != nil {
		return nil, fmt.Errorf("failed to unmarshal backup: %w", err)
	}

	return &backup, nil
}

// ensureBackupDir ensures the backup directory exists
func (s *BackupService) ensureBackupDir() error {
	if s.backupDir == "" {
		s.backupDir = "./backups"
	}

	if _, err := os.Stat(s.backupDir); os.IsNotExist(err) {
		if err := os.MkdirAll(s.backupDir, 0755); err != nil {
			return fmt.Errorf("failed to create backup directory: %w", err)
		}
	}

	return nil
}

// ValidateBackup validates a backup file
func (s *BackupService) ValidateBackup(backupID string) error {
	backup, err := s.loadBackupFromFile(backupID)
	if err != nil {
		return fmt.Errorf("failed to load backup: %w", err)
	}

	// Basic validation
	if backup.ID == "" {
		return fmt.Errorf("backup ID is empty")
	}

	if backup.Settings == nil {
		return fmt.Errorf("backup settings is nil")
	}

	if backup.Metadata.TotalSettings == 0 {
		return fmt.Errorf("backup contains no settings")
	}

	// Validate each category has valid settings
	for category, settings := range backup.Settings {
		if len(settings) == 0 {
			continue // Empty category is OK
		}

		for _, setting := range settings {
			if setting.Key == "" {
				return fmt.Errorf("setting in category %s has empty key", category)
			}
			if setting.Category != category {
				return fmt.Errorf("setting %s has mismatched category: expected %s, got %s", setting.Key, category, setting.Category)
			}
		}
	}

	return nil
}
