package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"

	"gorm.io/gorm"
)

// MigrationRecord tracks applied migrations
type MigrationRecord struct {
	ID        uint      `gorm:"primaryKey"`
	Version   string    `gorm:"uniqueIndex;not null"`
	Name      string    `gorm:"not null"`
	AppliedAt time.Time `gorm:"autoCreateTime"`
}

// TableName returns the table name for MigrationRecord
func (MigrationRecord) TableName() string {
	return "schema_migrations"
}

// Migration represents a database migration
type Migration struct {
	Version string
	Name    string
	Up      func(*gorm.DB) error
	Down    func(*gorm.DB) error
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db         *gorm.DB
	migrations []Migration
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *gorm.DB) *MigrationManager {
	return &MigrationManager{
		db:         db,
		migrations: getMigrations(),
	}
}

// RunMigrations runs all pending migrations
func (m *MigrationManager) RunMigrations(ctx context.Context) error {
	log.Println("🔄 Starting database migrations...")

	// Create migration tracking table
	if err := m.db.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("failed to create migration tracking table: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Run pending migrations
	for _, migration := range m.migrations {
		if _, applied := appliedMigrations[migration.Version]; applied {
			log.Printf("⏭️  Skipping migration %s (already applied)", migration.Version)
			continue
		}

		log.Printf("🔧 Running migration %s: %s", migration.Version, migration.Name)

		// Run migration in transaction
		err := m.db.Transaction(func(tx *gorm.DB) error {
			// Run the migration
			if err := migration.Up(tx); err != nil {
				return fmt.Errorf("migration %s failed: %w", migration.Version, err)
			}

			// Record migration as applied
			record := MigrationRecord{
				Version: migration.Version,
				Name:    migration.Name,
			}
			if err := tx.Create(&record).Error; err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
			}

			return nil
		})

		if err != nil {
			return err
		}

		log.Printf("✅ Migration %s completed successfully", migration.Version)
	}

	log.Println("🎉 All migrations completed successfully")
	return nil
}

// RollbackMigration rolls back the last migration
func (m *MigrationManager) RollbackMigration(ctx context.Context) error {
	log.Println("🔄 Rolling back last migration...")

	// Get the last applied migration
	var lastMigration MigrationRecord
	err := m.db.Order("applied_at DESC").First(&lastMigration).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("ℹ️  No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	// Find the migration definition
	var migrationDef *Migration
	for _, migration := range m.migrations {
		if migration.Version == lastMigration.Version {
			migrationDef = &migration
			break
		}
	}

	if migrationDef == nil {
		return fmt.Errorf("migration definition not found for version %s", lastMigration.Version)
	}

	log.Printf("🔧 Rolling back migration %s: %s", migrationDef.Version, migrationDef.Name)

	// Run rollback in transaction
	err = m.db.Transaction(func(tx *gorm.DB) error {
		// Run the rollback
		if err := migrationDef.Down(tx); err != nil {
			return fmt.Errorf("rollback %s failed: %w", migrationDef.Version, err)
		}

		// Remove migration record
		if err := tx.Delete(&lastMigration).Error; err != nil {
			return fmt.Errorf("failed to remove migration record %s: %w", migrationDef.Version, err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	log.Printf("✅ Migration %s rolled back successfully", migrationDef.Version)
	return nil
}

// GetMigrationStatus returns the status of all migrations
func (m *MigrationManager) GetMigrationStatus() ([]MigrationStatus, error) {
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return nil, err
	}

	var status []MigrationStatus
	for _, migration := range m.migrations {
		record, applied := appliedMigrations[migration.Version]
		migrationStatus := MigrationStatus{
			Version: migration.Version,
			Name:    migration.Name,
			Applied: applied,
		}
		if applied {
			migrationStatus.AppliedAt = &record.AppliedAt
		}
		status = append(status, migrationStatus)
	}

	return status, nil
}

// MigrationStatus represents the status of a migration
type MigrationStatus struct {
	Version   string     `json:"version"`
	Name      string     `json:"name"`
	Applied   bool       `json:"applied"`
	AppliedAt *time.Time `json:"applied_at,omitempty"`
}

// getAppliedMigrations returns a map of applied migrations
func (m *MigrationManager) getAppliedMigrations() (map[string]MigrationRecord, error) {
	var records []MigrationRecord
	if err := m.db.Find(&records).Error; err != nil {
		return nil, err
	}

	appliedMigrations := make(map[string]MigrationRecord)
	for _, record := range records {
		appliedMigrations[record.Version] = record
	}

	return appliedMigrations, nil
}

// getMigrations returns all available migrations in order
func getMigrations() []Migration {
	return []Migration{
		{
			Version: "001_initial_schema",
			Name:    "Create initial database schema",
			Up:      migration001Up,
			Down:    migration001Down,
		},
		{
			Version: "002_add_cart_enhancements",
			Name:    "Add cart session and calculated fields",
			Up:      migration002Up,
			Down:    migration002Down,
		},
		{
			Version: "003_add_user_enhancements",
			Name:    "Add user OAuth and enhanced fields",
			Up:      migration003Up,
			Down:    migration003Down,
		},
		{
			Version: "004_add_indexes",
			Name:    "Add performance indexes",
			Up:      migration004Up,
			Down:    migration004Down,
		},
		{
			Version: "005_add_cleanup_fields",
			Name:    "Add cleanup and expiration fields",
			Up:      migration005Up,
			Down:    migration005Down,
		},
		{
			Version: "007_sync_inventory_data",
			Name:    "Create inventory records for existing products",
			Up:      migration007Up,
			Down:    migration007Down,
		},
		{
			Version: "008_update_user_verification_and_sessions",
			Name:    "Update user verification structure and add user sessions",
			Up:      migration008Up,
			Down:    migration008Down,
		},
		{
			Version: "009_add_payment_method_to_orders",
			Name:    "Add payment method field to orders table",
			Up:      migration009Up,
			Down:    migration009Down,
		},
		{
			Version: "010_add_weight_to_order_items",
			Name:    "Add weight field to order_items table",
			Up:      migration010Up,
			Down:    migration010Down,
		},
		{
			Version: "011_create_product_categories",
			Name:    "Create product_categories table for many-to-many relationship",
			Up:      migration011Up,
			Down:    migration011Down,
		},
		{
			Version: "012_add_abandoned_cart_fields",
			Name:    "Add abandoned cart tracking fields to carts table",
			Up:      migration012Up,
			Down:    migration012Down,
		},
		{
			Version: "013_add_email_system",
			Name:    "Add email system tables",
			Up:      migration013Up,
			Down:    migration013Down,
		},
		// Temporarily disabled due to product_tags issue
		// {
		// 	Version: "006_enhance_search",
		// 	Name:    "Enhance full-text search capabilities",
		// 	Up:      migration006Up,
		// 	Down:    migration006Down,
		// },
	}
}

// migration013Up adds email system tables
func migration013Up(db *gorm.DB) error {
	// Create email system tables
	if err := db.AutoMigrate(
		&entities.Email{},
		&entities.EmailTemplate{},
		&entities.EmailSubscription{},
	); err != nil {
		return fmt.Errorf("failed to create email system tables: %w", err)
	}

	return nil
}

// migration013Down removes email system tables
func migration013Down(db *gorm.DB) error {
	// Drop email system tables
	if err := db.Migrator().DropTable(&entities.EmailSubscription{}); err != nil {
		return fmt.Errorf("failed to drop email_subscriptions table: %w", err)
	}
	if err := db.Migrator().DropTable(&entities.EmailTemplate{}); err != nil {
		return fmt.Errorf("failed to drop email_templates table: %w", err)
	}
	if err := db.Migrator().DropTable(&entities.Email{}); err != nil {
		return fmt.Errorf("failed to drop emails table: %w", err)
	}

	return nil
}
