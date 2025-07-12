package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"ecom-golang-clean-architecture/internal/config"
	"ecom-golang-clean-architecture/internal/infrastructure/database"
)

func main() {
	var (
		action = flag.String("action", "up", "Migration action: up, down, status")
		configPath = flag.String("config", "configs/config.yaml", "Path to config file")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database connection
	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create migration manager
	migrationManager := database.NewMigrationManager(db)
	ctx := context.Background()

	switch *action {
	case "up":
		fmt.Println("🔄 Running migrations...")
		if err := migrationManager.RunMigrations(ctx); err != nil {
			log.Fatal("Migration failed:", err)
		}
		fmt.Println("✅ Migrations completed successfully")

	case "down":
		fmt.Println("🔄 Rolling back last migration...")
		if err := migrationManager.RollbackMigration(ctx); err != nil {
			log.Fatal("Rollback failed:", err)
		}
		fmt.Println("✅ Rollback completed successfully")

	case "status":
		fmt.Println("📊 Migration Status:")
		status, err := migrationManager.GetMigrationStatus()
		if err != nil {
			log.Fatal("Failed to get migration status:", err)
		}

		fmt.Printf("%-25s %-50s %-10s %s\n", "Version", "Name", "Applied", "Applied At")
		fmt.Println(strings.Repeat("-", 100))
		
		for _, migration := range status {
			appliedStatus := "❌ No"
			appliedAt := ""
			if migration.Applied {
				appliedStatus = "✅ Yes"
				if migration.AppliedAt != nil {
					appliedAt = migration.AppliedAt.Format("2006-01-02 15:04:05")
				}
			}
			fmt.Printf("%-25s %-50s %-10s %s\n", 
				migration.Version, 
				migration.Name, 
				appliedStatus, 
				appliedAt)
		}

	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: up, down, status")
		os.Exit(1)
	}
}
