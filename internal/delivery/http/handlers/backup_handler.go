package handlers

import (
	"net/http"
	"time"

	"ecom-golang-clean-architecture/internal/usecases"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BackupHandler handles backup-related requests
type BackupHandler struct {
	settingsUseCase usecases.SettingsUseCase
}

// NewBackupHandler creates a new backup handler
func NewBackupHandler(settingsUseCase usecases.SettingsUseCase) *BackupHandler {
	return &BackupHandler{
		settingsUseCase: settingsUseCase,
	}
}

// BackupData represents backup data structure
type BackupData struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Description string                 `json:"description"`
	Settings    map[string]interface{} `json:"settings"`
	Metadata    BackupMetadata         `json:"metadata"`
}

// BackupMetadata contains backup metadata
type BackupMetadata struct {
	CreatedBy     string   `json:"created_by"`
	TotalSettings int      `json:"total_settings"`
	Categories    []string `json:"categories"`
	Version       string   `json:"version"`
}

// CreateBackupRequest represents the request to create a backup
type CreateBackupRequest struct {
	Description string `json:"description" binding:"required"`
}

// CreateBackup creates a new settings backup
func (h *BackupHandler) CreateBackup(c *gin.Context) {
	var req CreateBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Get all settings categories
	categories := []string{"general", "store", "payment", "email", "tax", "shipping", "seo", "security", "notifications", "integrations"}
	allSettings := make(map[string]interface{})
	totalSettings := 0

	for _, category := range categories {
		var settings interface{}
		var err error

		switch category {
		case "general":
			settings, err = h.settingsUseCase.GetGeneralSettings(c.Request.Context())
		case "store":
			settings, err = h.settingsUseCase.GetStoreConfig(c.Request.Context())
		case "payment":
			settings, err = h.settingsUseCase.GetPaymentSettings(c.Request.Context())
		case "email":
			settings, err = h.settingsUseCase.GetEmailSettings(c.Request.Context())
		case "tax":
			settings, err = h.settingsUseCase.GetTaxSettings(c.Request.Context())
		case "shipping":
			settings, err = h.settingsUseCase.GetShippingSettings(c.Request.Context())
		case "seo":
			settings, err = h.settingsUseCase.GetSEOSettings(c.Request.Context())
		case "security":
			settings, err = h.settingsUseCase.GetSecuritySettings(c.Request.Context())
		case "notifications":
			settings, err = h.settingsUseCase.GetNotificationSettings(c.Request.Context())
		case "integrations":
			settings, err = h.settingsUseCase.GetIntegrationSettings(c.Request.Context())
		}

		if err == nil && settings != nil {
			allSettings[category] = settings
			totalSettings++
		}
	}

	// Create backup
	backup := &BackupData{
		ID:          uuid.New().String(),
		Timestamp:   time.Now(),
		Description: req.Description,
		Settings:    allSettings,
		Metadata: BackupMetadata{
			CreatedBy:     "admin@example.com", // Should come from auth context
			TotalSettings: totalSettings,
			Categories:    categories,
			Version:       "1.0",
		},
	}

	c.JSON(http.StatusOK, gin.H{"message": "Backup created successfully", "data": backup})
}

// ListBackups lists all available backups (mock implementation)
func (h *BackupHandler) ListBackups(c *gin.Context) {
	// Mock backup list
	backups := []map[string]interface{}{
		{
			"id":          "backup-001",
			"timestamp":   time.Now().Add(-24 * time.Hour),
			"description": "Daily backup",
			"size":        "2.5 MB",
			"status":      "completed",
		},
		{
			"id":          "backup-002",
			"timestamp":   time.Now().Add(-48 * time.Hour),
			"description": "Pre-update backup",
			"size":        "2.3 MB",
			"status":      "completed",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Backups retrieved successfully",
		"data": map[string]interface{}{
			"backups": backups,
			"total":   len(backups),
		},
	})
}

// RestoreBackupRequest represents the request to restore a backup
type RestoreBackupRequest struct {
	BackupID string `json:"backup_id" binding:"required"`
}

// RestoreBackup restores settings from a backup (mock implementation)
func (h *BackupHandler) RestoreBackup(c *gin.Context) {
	var req RestoreBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Mock restore process
	result := map[string]interface{}{
		"backup_id":         req.BackupID,
		"restore_status":    "completed",
		"restored_at":       time.Now(),
		"settings_restored": 10,
		"categories":        []string{"general", "store", "payment", "email", "tax", "shipping", "seo", "security", "notifications", "integrations"},
	}

	c.JSON(http.StatusOK, gin.H{"message": "Backup restored successfully", "data": result})
}

// GetBackup gets a specific backup (mock implementation)
func (h *BackupHandler) GetBackup(c *gin.Context) {
	backupID := c.Param("id")
	if backupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Backup ID is required"})
		return
	}

	// Mock backup data
	backup := map[string]interface{}{
		"id":          backupID,
		"timestamp":   time.Now().Add(-24 * time.Hour),
		"description": "Mock backup for ID: " + backupID,
		"size":        "2.5 MB",
		"status":      "completed",
		"metadata": map[string]interface{}{
			"created_by":     "admin@example.com",
			"total_settings": 10,
			"categories":     []string{"general", "store", "payment", "email", "tax", "shipping", "seo", "security", "notifications", "integrations"},
			"version":        "1.0",
		},
	}

	c.JSON(http.StatusOK, gin.H{"message": "Backup retrieved successfully", "data": backup})
}

// DeleteBackup deletes a backup (mock implementation)
func (h *BackupHandler) DeleteBackup(c *gin.Context) {
	backupID := c.Param("id")
	if backupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Backup ID is required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Backup deleted successfully",
		"data": map[string]interface{}{
			"backup_id":  backupID,
			"deleted_at": time.Now(),
		},
	})
}

// GetBackupStats returns backup statistics (mock implementation)
func (h *BackupHandler) GetBackupStats(c *gin.Context) {
	stats := map[string]interface{}{
		"total_backups":    5,
		"total_size":       "12.5 MB",
		"latest_backup":    time.Now().Add(-6 * time.Hour),
		"oldest_backup":    time.Now().Add(-30 * 24 * time.Hour),
		"backup_frequency": "daily",
		"retention_days":   30,
		"auto_cleanup":     true,
		"storage_location": "./backups",
	}

	c.JSON(http.StatusOK, gin.H{"message": "Backup statistics retrieved successfully", "data": stats})
}

// ExportSettings exports current settings
func (h *BackupHandler) ExportSettings(c *gin.Context) {
	format := c.Query("format")
	if format == "" {
		format = "json"
	}

	// Get all settings (reuse backup logic)
	categories := []string{"general", "store", "payment", "email", "tax", "shipping", "seo", "security", "notifications", "integrations"}
	allSettings := make(map[string]interface{})

	for _, category := range categories {
		var settings interface{}
		var err error

		switch category {
		case "general":
			settings, err = h.settingsUseCase.GetGeneralSettings(c.Request.Context())
		case "store":
			settings, err = h.settingsUseCase.GetStoreConfig(c.Request.Context())
		case "payment":
			settings, err = h.settingsUseCase.GetPaymentSettings(c.Request.Context())
		case "email":
			settings, err = h.settingsUseCase.GetEmailSettings(c.Request.Context())
		case "tax":
			settings, err = h.settingsUseCase.GetTaxSettings(c.Request.Context())
		case "shipping":
			settings, err = h.settingsUseCase.GetShippingSettings(c.Request.Context())
		case "seo":
			settings, err = h.settingsUseCase.GetSEOSettings(c.Request.Context())
		case "security":
			settings, err = h.settingsUseCase.GetSecuritySettings(c.Request.Context())
		case "notifications":
			settings, err = h.settingsUseCase.GetNotificationSettings(c.Request.Context())
		case "integrations":
			settings, err = h.settingsUseCase.GetIntegrationSettings(c.Request.Context())
		}

		if err == nil && settings != nil {
			allSettings[category] = settings
		}
	}

	// Set headers for download
	filename := "settings_export_" + time.Now().Format("20060102_150405") + "." + format
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/json")

	exportData := map[string]interface{}{
		"export_id":   uuid.New().String(),
		"exported_at": time.Now(),
		"format":      format,
		"settings":    allSettings,
		"metadata": map[string]interface{}{
			"version":        "1.0",
			"total_settings": len(allSettings),
			"categories":     categories,
		},
	}

	c.JSON(http.StatusOK, exportData)
}
