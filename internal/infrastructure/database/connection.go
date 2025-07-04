package database

import (
	"fmt"
	"log"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewConnection creates a new database connection
func NewConnection(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	// Configure GORM logger
	var gormLogger logger.Interface
	if cfg.Host == "localhost" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		// Core entities
		&entities.User{},
		&entities.UserProfile{},
		&entities.Category{},
		&entities.Product{},
		&entities.ProductImage{},
		&entities.ProductTag{},
		// &entities.ProductProductTag{}, // Remove custom join table, let GORM handle it
		&entities.Cart{},
		&entities.CartItem{},
		&entities.Order{},
		&entities.OrderItem{},
		&entities.Payment{},

		// File uploads
		&entities.FileUpload{},

		// User management
		&entities.Address{},
		&entities.Wishlist{},
		&entities.UserPreference{},
		&entities.AccountVerification{},
		&entities.PasswordReset{},

		// Reviews & Ratings
		&entities.Review{},
		&entities.ReviewImage{},
		&entities.ReviewVote{},
		&entities.ProductRating{},

		// Coupons & Promotions
		&entities.Coupon{},
		&entities.CouponUsage{},
		&entities.Promotion{},
		&entities.LoyaltyProgram{},
		&entities.UserLoyaltyPoints{},

		// Inventory Management
		&entities.Inventory{},
		&entities.InventoryMovement{},
		&entities.Warehouse{},
		&entities.StockAlert{},
		&entities.Supplier{},

		// Shipping & Delivery
		&entities.ShippingMethod{},
		&entities.ShippingZone{},
		&entities.ShippingRate{},
		&entities.Shipment{},
		&entities.ShipmentTracking{},
		&entities.Return{},
		&entities.ReturnItem{},

		// Notifications
		&entities.Notification{},
		&entities.NotificationTemplate{},
		&entities.NotificationPreferences{},
		&entities.NotificationQueue{},

		// Analytics
		&entities.AnalyticsEvent{},
		&entities.SalesReport{},
		&entities.ProductAnalytics{},
		&entities.UserAnalytics{},
		&entities.CategoryAnalytics{},
		&entities.SearchAnalytics{},

		// Customer Support
		&entities.SupportTicket{},
		&entities.TicketMessage{},
		&entities.TicketAttachment{},
		&entities.FAQ{},
		&entities.KnowledgeBase{},
		&entities.LiveChatSession{},
		&entities.ChatMessage{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// CreateIndexes creates additional database indexes
func CreateIndexes(db *gorm.DB) error {
	log.Println("Creating database indexes...")

	// User indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active)")

	// Product indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_products_name ON products(name)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_products_status ON products(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_products_price ON products(price)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_products_stock ON products(stock)")

	// Category indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_categories_is_active ON categories(is_active)")

	// Cart indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_carts_user_id ON carts(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_cart_items_cart_id ON cart_items(cart_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_cart_items_product_id ON cart_items(product_id)")

	// Order indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_order_number ON orders(order_number)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_payment_status ON orders(payment_status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at)")

	// Order items indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id)")

	// Payment indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_payments_transaction_id ON payments(transaction_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_payments_method ON payments(method)")

	log.Println("Database indexes created successfully")
	return nil
}

// SeedData seeds initial data into the database
func SeedData(db *gorm.DB) error {
	log.Println("Seeding initial data...")

	// Create root categories
	categories := []entities.Category{
		{
			Name:        "Electronics",
			Description: "Electronic devices and accessories",
			Slug:        "electronics",
			IsActive:    true,
			SortOrder:   1,
		},
		{
			Name:        "Clothing",
			Description: "Fashion and apparel",
			Slug:        "clothing",
			IsActive:    true,
			SortOrder:   2,
		},
		{
			Name:        "Books",
			Description: "Books and literature",
			Slug:        "books",
			IsActive:    true,
			SortOrder:   3,
		},
		{
			Name:        "Home & Garden",
			Description: "Home improvement and gardening",
			Slug:        "home-garden",
			IsActive:    true,
			SortOrder:   4,
		},
	}

	for _, category := range categories {
		var existingCategory entities.Category
		if err := db.Where("slug = ?", category.Slug).First(&existingCategory).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&category).Error; err != nil {
					return fmt.Errorf("failed to create category %s: %w", category.Name, err)
				}
				log.Printf("Created category: %s", category.Name)
			}
		}
	}

	// Create product tags
	tags := []entities.ProductTag{
		{Name: "New", Slug: "new"},
		{Name: "Featured", Slug: "featured"},
		{Name: "Sale", Slug: "sale"},
		{Name: "Popular", Slug: "popular"},
		{Name: "Limited Edition", Slug: "limited-edition"},
	}

	for _, tag := range tags {
		var existingTag entities.ProductTag
		if err := db.Where("slug = ?", tag.Slug).First(&existingTag).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&tag).Error; err != nil {
					return fmt.Errorf("failed to create tag %s: %w", tag.Name, err)
				}
				log.Printf("Created tag: %s", tag.Name)
			}
		}
	}

	log.Println("Initial data seeded successfully")
	return nil
}
