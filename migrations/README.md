# Migration System

## Overview

This project uses **programmatic migrations** managed by Go code, not SQL files.

## Migration System Location

- **Migration Manager**: `internal/infrastructure/database/migration_manager.go`
- **Migration Functions**: `internal/infrastructure/database/migrations.go`
- **Execution**: Automatically runs on application startup

## How It Works

1. **Programmatic Migrations**: All migrations are written in Go using GORM
2. **Version Tracking**: Migrations are tracked in `schema_migrations` table
3. **Auto-execution**: Migrations run automatically when the application starts
4. **Rollback Support**: Each migration has an Up and Down function

## Adding New Migrations

1. Add migration function to `migrations.go`:
```go
func migration008Up(db *gorm.DB) error {
    // Your migration logic here
    return db.AutoMigrate(&entities.NewEntity{})
}

func migration008Down(db *gorm.DB) error {
    // Rollback logic here
    return db.Migrator().DropTable(&entities.NewEntity{})
}
```

2. Register migration in `migration_manager.go`:
```go
{
    Version: "008_your_migration_name",
    Name:    "Description of your migration",
    Up:      migration008Up,
    Down:    migration008Down,
},
```

## Current Migrations

- `001_initial_schema`: Creates all core entities
- `002_add_indexes`: Adds database indexes for performance
- `003_add_constraints`: Adds foreign key constraints
- `004_add_triggers`: Adds database triggers
- `005_add_cleanup_fields`: Adds cleanup and expiration fields
- `006_add_recommendation_tables`: Adds recommendation system tables
- `007_sync_inventory_data`: Creates inventory records for existing products

## Important Notes

- **DO NOT** add SQL files to this directory - they will be ignored
- **DO NOT** use `db.AutoMigrate()` in application code - use migrations instead
- All schema changes must go through the migration system
- Test migrations thoroughly before deploying to production

## Checking Migration Status

```bash
# Connect to database and check applied migrations
docker exec -it ecom_postgres psql -U postgres -d ecommerce_db -c "SELECT version, name, applied_at FROM schema_migrations ORDER BY applied_at;"
```
