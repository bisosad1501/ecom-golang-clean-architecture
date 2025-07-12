package database

import (
	"context"
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

	// Configure optimized connection pool
	sqlDB.SetMaxIdleConns(25)                // Increased from 10
	sqlDB.SetMaxOpenConns(200)               // Increased from 100
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Reduced from 1 hour
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // Added idle timeout

	// Enable query logging for slow queries in development
	// Note: Add environment check when config supports it
	// if cfg.Environment == "development" {
	//     db = db.Debug()
	// }

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}

// TransactionManager provides transaction management utilities
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// WithTransaction executes a function within a database transaction
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

// WithTransactionResult executes a function within a database transaction and returns a result
func (tm *TransactionManager) WithTransactionResult(ctx context.Context, fn func(*gorm.DB) (interface{}, error)) (interface{}, error) {
	var result interface{}
	err := tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = fn(tx)
		return err
	})
	return result, err
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

		// Brand and Product Extensions
		&entities.Brand{},
		&entities.ProductVariant{},
		&entities.ProductAttribute{},
		&entities.ProductAttributeTerm{},
		&entities.ProductAttributeValue{},
		&entities.ProductVariantAttribute{},

		&entities.Cart{},
		&entities.CartItem{},
		&entities.Order{},
		&entities.OrderItem{},
		&entities.OrderEvent{},
		&entities.Payment{},
		&entities.StockReservation{},

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

	// New indexes for existing entities with new fields
	// Cart indexes (if new fields exist)
	db.Exec("CREATE INDEX IF NOT EXISTS idx_carts_status ON carts(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_carts_expires_at ON carts(expires_at)")

	// Payment indexes (if new fields exist)
	db.Exec("CREATE INDEX IF NOT EXISTS idx_payments_payment_intent_id ON payments(payment_intent_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_payments_gateway ON payments(gateway)")

	// Order indexes for new fields (if they exist)
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_reserved_until ON orders(reserved_until)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_payment_timeout ON orders(payment_timeout)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_version ON orders(version)")

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

	// Create brands
	brands := []entities.Brand{
		{Name: "Apple", Slug: "apple", Description: "Technology company", IsActive: true},
		{Name: "Samsung", Slug: "samsung", Description: "Electronics manufacturer", IsActive: true},
		{Name: "Nike", Slug: "nike", Description: "Sportswear brand", IsActive: true},
		{Name: "Adidas", Slug: "adidas", Description: "Athletic apparel", IsActive: true},
		{Name: "Sony", Slug: "sony", Description: "Electronics and entertainment", IsActive: true},
	}

	for _, brand := range brands {
		var existingBrand entities.Brand
		if err := db.Where("slug = ?", brand.Slug).First(&existingBrand).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&brand).Error; err != nil {
					return fmt.Errorf("failed to create brand %s: %w", brand.Name, err)
				}
				log.Printf("Created brand: %s", brand.Name)
			}
		}
	}

	// Create product attributes
	attributes := []entities.ProductAttribute{
		{Name: "Color", Slug: "color", Type: "color", IsVisible: true},
		{Name: "Size", Slug: "size", Type: "select", IsVisible: true},
		{Name: "Material", Slug: "material", Type: "text", IsVisible: true},
		{Name: "Brand", Slug: "brand", Type: "select", IsVisible: true},
	}

	for _, attr := range attributes {
		var existingAttr entities.ProductAttribute
		if err := db.Where("slug = ?", attr.Slug).First(&existingAttr).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&attr).Error; err != nil {
					return fmt.Errorf("failed to create attribute %s: %w", attr.Name, err)
				}
				log.Printf("Created attribute: %s", attr.Name)
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
