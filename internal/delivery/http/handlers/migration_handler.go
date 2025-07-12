package handlers

import (
	"net/http"

	"ecom-golang-clean-architecture/internal/infrastructure/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MigrationHandler handles migration-related HTTP requests
type MigrationHandler struct {
	migrationManager *database.MigrationManager
}

// NewMigrationHandler creates a new migration handler
func NewMigrationHandler(db *gorm.DB) *MigrationHandler {
	return &MigrationHandler{
		migrationManager: database.NewMigrationManager(db),
	}
}

// GetMigrationStatus returns the status of all migrations
// @Summary Get migration status
// @Description Get the status of all database migrations
// @Tags migrations
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/migrations/status [get]
func (h *MigrationHandler) GetMigrationStatus(c *gin.Context) {
	status, err := h.migrationManager.GetMigrationStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get migration status",
			"details": err.Error(),
		})
		return
	}

	// Calculate summary
	var applied, pending int
	for _, migration := range status {
		if migration.Applied {
			applied++
		} else {
			pending++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"migrations": status,
			"summary": gin.H{
				"total":   len(status),
				"applied": applied,
				"pending": pending,
			},
		},
	})
}

// RunMigrations runs all pending migrations
// @Summary Run migrations
// @Description Run all pending database migrations
// @Tags migrations
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/migrations/run [post]
func (h *MigrationHandler) RunMigrations(c *gin.Context) {
	if err := h.migrationManager.RunMigrations(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to run migrations",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "All migrations completed successfully",
	})
}

// RollbackMigration rolls back the last migration
// @Summary Rollback migration
// @Description Rollback the last applied migration
// @Tags migrations
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/migrations/rollback [post]
func (h *MigrationHandler) RollbackMigration(c *gin.Context) {
	if err := h.migrationManager.RollbackMigration(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to rollback migration",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Migration rolled back successfully",
	})
}
