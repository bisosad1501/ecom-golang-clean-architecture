package database

import (
	"ecom-golang-clean-architecture/internal/domain/entities"
	"gorm.io/gorm"
)

// migration001Up creates the initial database schema
func migration001Up(db *gorm.DB) error {
	// Create all core entities using AutoMigrate
	return db.AutoMigrate(
		// Core entities
		&entities.User{},
		&entities.UserProfile{},
		&entities.Category{},
		&entities.Product{},
		&entities.ProductImage{},
		&entities.ProductTag{},

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

		// Recommendation system
		&entities.UserProductInteraction{},
		&entities.ProductRecommendation{},
		&entities.ProductSimilarity{},
		&entities.FrequentlyBoughtTogether{},
		&entities.TrendingProduct{},

		// Product comparison
		&entities.ProductComparison{},
		&entities.ProductComparisonItem{},

		// Advanced filtering
		&entities.FilterSet{},
		&entities.FilterUsage{},
		&entities.ProductFilterOption{},
	)
}

// migration001Down drops the initial schema (dangerous - use with caution)
func migration001Down(db *gorm.DB) error {
	// Note: This is a destructive operation
	// In production, you might want to backup data first
	tables := []string{
		"chat_messages", "live_chat_sessions", "knowledge_bases", "faqs",
		"ticket_attachments", "ticket_messages", "support_tickets",
		"search_analytics", "category_analytics", "user_analytics",
		"product_analytics", "sales_reports", "analytics_events",
		"notification_queues", "notification_preferences", "notification_templates", "notifications",
		"return_items", "returns", "shipment_trackings", "shipments",
		"shipping_rates", "shipping_zones", "shipping_methods",
		"suppliers", "stock_alerts", "warehouses", "inventory_movements", "inventories",
		"user_loyalty_points", "loyalty_programs", "promotions", "coupon_usages", "coupons",
		"product_ratings", "review_votes", "review_images", "reviews",
		"password_resets", "account_verifications", "user_preferences", "wishlists", "addresses",
		"file_uploads", "stock_reservations", "payments", "order_events", "order_items", "orders",
		"cart_items", "carts", "product_variant_attributes", "product_attribute_values",
		"product_attribute_terms", "product_attributes", "product_variants", "brands",
		"product_images", "products", "categories", "user_profiles", "users",
	}

	for _, table := range tables {
		if err := db.Exec("DROP TABLE IF EXISTS " + table + " CASCADE").Error; err != nil {
			return err
		}
	}

	return nil
}

// migration002Up adds cart enhancements
func migration002Up(db *gorm.DB) error {
	// Add new fields to carts table
	sqls := []string{
		"ALTER TABLE carts ADD COLUMN IF NOT EXISTS session_id TEXT",
		"ALTER TABLE carts ADD COLUMN IF NOT EXISTS status TEXT DEFAULT 'active'",
		"ALTER TABLE carts ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP WITH TIME ZONE",
		"ALTER TABLE carts ADD COLUMN IF NOT EXISTS subtotal NUMERIC DEFAULT 0",
		"ALTER TABLE carts ADD COLUMN IF NOT EXISTS total NUMERIC DEFAULT 0",
		"ALTER TABLE carts ADD COLUMN IF NOT EXISTS item_count INTEGER DEFAULT 0",
		"ALTER TABLE carts ADD COLUMN IF NOT EXISTS currency TEXT DEFAULT 'USD'",
		"ALTER TABLE carts ADD COLUMN IF NOT EXISTS notes TEXT",
		"CREATE INDEX IF NOT EXISTS idx_carts_session_id ON carts(session_id)",
		"CREATE INDEX IF NOT EXISTS idx_carts_status ON carts(status)",
		"CREATE INDEX IF NOT EXISTS idx_carts_expires_at ON carts(expires_at)",
	}

	for _, sql := range sqls {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}

	return nil
}

// migration002Down removes cart enhancements
func migration002Down(db *gorm.DB) error {
	sqls := []string{
		"DROP INDEX IF EXISTS idx_carts_expires_at",
		"DROP INDEX IF EXISTS idx_carts_status",
		"DROP INDEX IF EXISTS idx_carts_session_id",
		"ALTER TABLE carts DROP COLUMN IF EXISTS notes",
		"ALTER TABLE carts DROP COLUMN IF EXISTS currency",
		"ALTER TABLE carts DROP COLUMN IF EXISTS item_count",
		"ALTER TABLE carts DROP COLUMN IF EXISTS total",
		"ALTER TABLE carts DROP COLUMN IF EXISTS subtotal",
		"ALTER TABLE carts DROP COLUMN IF EXISTS expires_at",
		"ALTER TABLE carts DROP COLUMN IF EXISTS status",
		"ALTER TABLE carts DROP COLUMN IF EXISTS session_id",
	}

	for _, sql := range sqls {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}

	return nil
}

// migration003Up adds user enhancements
func migration003Up(db *gorm.DB) error {
	sqls := []string{
		// OAuth fields
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS google_id TEXT",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS facebook_id TEXT",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar TEXT",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS is_oauth_user BOOLEAN DEFAULT false",

		// Enhanced user fields
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS username TEXT",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS date_of_birth TIMESTAMP WITH TIME ZONE",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS gender TEXT",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS language TEXT DEFAULT 'en'",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS timezone TEXT DEFAULT 'UTC'",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS currency TEXT DEFAULT 'USD'",

		// Security fields
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS two_factor_enabled BOOLEAN DEFAULT false",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS two_factor_secret TEXT",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS backup_codes TEXT[]",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS security_score INTEGER DEFAULT 0",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS last_password_change TIMESTAMP WITH TIME ZONE",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS failed_login_attempts INTEGER DEFAULT 0",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS locked_until TIMESTAMP WITH TIME ZONE",

		// Activity tracking
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMP WITH TIME ZONE",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS last_activity_at TIMESTAMP WITH TIME ZONE",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS login_count INTEGER DEFAULT 0",

		// Customer metrics
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS total_orders INTEGER DEFAULT 0",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS total_spent NUMERIC DEFAULT 0",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS loyalty_points INTEGER DEFAULT 0",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS membership_tier TEXT DEFAULT 'bronze'",

		// Create indexes
		"CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id)",
		"CREATE INDEX IF NOT EXISTS idx_users_facebook_id ON users(facebook_id)",
		"CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)",
		"CREATE INDEX IF NOT EXISTS idx_users_two_factor_enabled ON users(two_factor_enabled)",
		"CREATE INDEX IF NOT EXISTS idx_users_last_login_at ON users(last_login_at)",
		"CREATE INDEX IF NOT EXISTS idx_users_membership_tier ON users(membership_tier)",
	}

	for _, sql := range sqls {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}

	return nil
}

// migration003Down removes user enhancements
func migration003Down(db *gorm.DB) error {
	sqls := []string{
		// Drop indexes
		"DROP INDEX IF EXISTS idx_users_membership_tier",
		"DROP INDEX IF EXISTS idx_users_last_login_at",
		"DROP INDEX IF EXISTS idx_users_two_factor_enabled",
		"DROP INDEX IF EXISTS idx_users_username",
		"DROP INDEX IF EXISTS idx_users_facebook_id",
		"DROP INDEX IF EXISTS idx_users_google_id",

		// Drop columns
		"ALTER TABLE users DROP COLUMN IF EXISTS membership_tier",
		"ALTER TABLE users DROP COLUMN IF EXISTS loyalty_points",
		"ALTER TABLE users DROP COLUMN IF EXISTS total_spent",
		"ALTER TABLE users DROP COLUMN IF EXISTS total_orders",
		"ALTER TABLE users DROP COLUMN IF EXISTS login_count",
		"ALTER TABLE users DROP COLUMN IF EXISTS last_activity_at",
		"ALTER TABLE users DROP COLUMN IF EXISTS last_login_at",
		"ALTER TABLE users DROP COLUMN IF EXISTS locked_until",
		"ALTER TABLE users DROP COLUMN IF EXISTS failed_login_attempts",
		"ALTER TABLE users DROP COLUMN IF EXISTS last_password_change",
		"ALTER TABLE users DROP COLUMN IF EXISTS security_score",
		"ALTER TABLE users DROP COLUMN IF EXISTS backup_codes",
		"ALTER TABLE users DROP COLUMN IF EXISTS two_factor_secret",
		"ALTER TABLE users DROP COLUMN IF EXISTS two_factor_enabled",
		"ALTER TABLE users DROP COLUMN IF EXISTS currency",
		"ALTER TABLE users DROP COLUMN IF EXISTS timezone",
		"ALTER TABLE users DROP COLUMN IF EXISTS language",
		"ALTER TABLE users DROP COLUMN IF EXISTS gender",
		"ALTER TABLE users DROP COLUMN IF EXISTS date_of_birth",
		"ALTER TABLE users DROP COLUMN IF EXISTS username",
		"ALTER TABLE users DROP COLUMN IF EXISTS is_oauth_user",
		"ALTER TABLE users DROP COLUMN IF EXISTS avatar",
		"ALTER TABLE users DROP COLUMN IF EXISTS facebook_id",
		"ALTER TABLE users DROP COLUMN IF EXISTS google_id",
	}

	for _, sql := range sqls {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}

	return nil
}

// migration004Up adds performance indexes
func migration004Up(db *gorm.DB) error {
	sqls := []string{
		// User indexes
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)",
		"CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active)",
		"CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)",
		"CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at)",

		// Product indexes
		"CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id)",
		"CREATE INDEX IF NOT EXISTS idx_products_brand_id ON products(brand_id)",
		"CREATE INDEX IF NOT EXISTS idx_products_status ON products(status)",
		"CREATE INDEX IF NOT EXISTS idx_products_price ON products(price)",
		"CREATE INDEX IF NOT EXISTS idx_products_stock ON products(stock)",
		"CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_products_name_gin ON products USING gin(to_tsvector('english', name))",
		"CREATE INDEX IF NOT EXISTS idx_products_description_gin ON products USING gin(to_tsvector('english', description))",

		// Order indexes
		"CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status)",
		"CREATE INDEX IF NOT EXISTS idx_orders_payment_status ON orders(payment_status)",
		"CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_orders_total ON orders(total)",

		// Cart indexes
		"CREATE INDEX IF NOT EXISTS idx_carts_user_id ON carts(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_cart_items_cart_id ON cart_items(cart_id)",
		"CREATE INDEX IF NOT EXISTS idx_cart_items_product_id ON cart_items(product_id)",

		// Review indexes
		"CREATE INDEX IF NOT EXISTS idx_reviews_product_id ON reviews(product_id)",
		"CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_reviews_status ON reviews(status)",
		"CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews(rating)",
		"CREATE INDEX IF NOT EXISTS idx_reviews_created_at ON reviews(created_at)",

		// Category indexes
		"CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id)",
		"CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug)",
		"CREATE INDEX IF NOT EXISTS idx_categories_is_active ON categories(is_active)",

		// Payment indexes
		"CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id)",
		"CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status)",
		"CREATE INDEX IF NOT EXISTS idx_payments_method ON payments(method)",
		"CREATE INDEX IF NOT EXISTS idx_payments_created_at ON payments(created_at)",

		// Stock reservation indexes
		"CREATE INDEX IF NOT EXISTS idx_stock_reservations_product_id ON stock_reservations(product_id)",
		"CREATE INDEX IF NOT EXISTS idx_stock_reservations_user_id ON stock_reservations(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_stock_reservations_expires_at ON stock_reservations(expires_at)",
		"CREATE INDEX IF NOT EXISTS idx_stock_reservations_status ON stock_reservations(status)",

		// Analytics indexes
		"CREATE INDEX IF NOT EXISTS idx_analytics_events_event_type ON analytics_events(event_type)",
		"CREATE INDEX IF NOT EXISTS idx_analytics_events_user_id ON analytics_events(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_analytics_events_created_at ON analytics_events(created_at)",

		// Notification indexes
		"CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status)",
		"CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at)",
	}

	for _, sql := range sqls {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}

	return nil
}

// migration004Down removes performance indexes
func migration004Down(db *gorm.DB) error {
	indexes := []string{
		"idx_notifications_created_at", "idx_notifications_status", "idx_notifications_user_id",
		"idx_analytics_events_created_at", "idx_analytics_events_user_id", "idx_analytics_events_event_type",
		"idx_stock_reservations_status", "idx_stock_reservations_expires_at", "idx_stock_reservations_user_id", "idx_stock_reservations_product_id",
		"idx_payments_created_at", "idx_payments_method", "idx_payments_status", "idx_payments_order_id",
		"idx_categories_is_active", "idx_categories_slug", "idx_categories_parent_id",
		"idx_reviews_created_at", "idx_reviews_rating", "idx_reviews_status", "idx_reviews_user_id", "idx_reviews_product_id",
		"idx_cart_items_product_id", "idx_cart_items_cart_id", "idx_carts_user_id",
		"idx_orders_total", "idx_orders_created_at", "idx_orders_payment_status", "idx_orders_status", "idx_orders_user_id",
		"idx_products_description_gin", "idx_products_name_gin", "idx_products_created_at", "idx_products_stock", "idx_products_price", "idx_products_status", "idx_products_brand_id", "idx_products_category_id",
		"idx_users_created_at", "idx_users_status", "idx_users_is_active", "idx_users_role", "idx_users_email",
	}

	for _, index := range indexes {
		if err := db.Exec("DROP INDEX IF EXISTS " + index).Error; err != nil {
			return err
		}
	}

	return nil
}

// migration005Up adds cleanup and expiration fields
func migration005Up(db *gorm.DB) error {
	sqls := []string{
		// Add expiration fields to various tables
		"ALTER TABLE orders ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP WITH TIME ZONE",
		"ALTER TABLE stock_reservations ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP WITH TIME ZONE",

		// Add cleanup tracking fields
		"ALTER TABLE carts ADD COLUMN IF NOT EXISTS last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()",
		"ALTER TABLE orders ADD COLUMN IF NOT EXISTS last_status_change_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()",
		"ALTER TABLE stock_reservations ADD COLUMN IF NOT EXISTS last_updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()",

		// Create indexes for cleanup operations
		"CREATE INDEX IF NOT EXISTS idx_orders_expires_at ON orders(expires_at)",
		"CREATE INDEX IF NOT EXISTS idx_stock_reservations_expires_at ON stock_reservations(expires_at)",
		"CREATE INDEX IF NOT EXISTS idx_carts_last_activity_at ON carts(last_activity_at)",
		"CREATE INDEX IF NOT EXISTS idx_orders_last_status_change_at ON orders(last_status_change_at)",
		"CREATE INDEX IF NOT EXISTS idx_stock_reservations_last_updated_at ON stock_reservations(last_updated_at)",
	}

	for _, sql := range sqls {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}

	return nil
}

// migration005Down removes cleanup and expiration fields
func migration005Down(db *gorm.DB) error {
	sqls := []string{
		// Drop indexes
		"DROP INDEX IF EXISTS idx_stock_reservations_last_updated_at",
		"DROP INDEX IF EXISTS idx_orders_last_status_change_at",
		"DROP INDEX IF EXISTS idx_carts_last_activity_at",
		"DROP INDEX IF EXISTS idx_stock_reservations_expires_at",
		"DROP INDEX IF EXISTS idx_orders_expires_at",

		// Drop columns
		"ALTER TABLE stock_reservations DROP COLUMN IF EXISTS last_updated_at",
		"ALTER TABLE orders DROP COLUMN IF EXISTS last_status_change_at",
		"ALTER TABLE carts DROP COLUMN IF EXISTS last_activity_at",
		"ALTER TABLE stock_reservations DROP COLUMN IF EXISTS expires_at",
		"ALTER TABLE orders DROP COLUMN IF EXISTS expires_at",
	}

	for _, sql := range sqls {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}

	return nil
}

// migration006Up enhances full-text search capabilities
func migration006Up(db *gorm.DB) error {
	sqls := []string{
		// Enhanced product search indexes
		"CREATE INDEX IF NOT EXISTS idx_products_featured ON products(featured)",
		"CREATE INDEX IF NOT EXISTS idx_products_visibility ON products(visibility)",
		"CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku)",
		"CREATE INDEX IF NOT EXISTS idx_products_sale_price ON products(sale_price)",
		"CREATE INDEX IF NOT EXISTS idx_products_stock_status ON products(stock_status)",

		// Composite search vector for better full-text search
		"CREATE INDEX IF NOT EXISTS idx_products_search_vector ON products USING gin(to_tsvector('english', coalesce(name, '') || ' ' || coalesce(description, '') || ' ' || coalesce(short_description, '') || ' ' || coalesce(sku, '') || ' ' || coalesce(keywords, '')))",

		// SKU specific search
		"CREATE INDEX IF NOT EXISTS idx_products_sku_gin ON products USING gin(to_tsvector('english', sku))",

		// Brand and category search indexes
		"CREATE INDEX IF NOT EXISTS idx_brands_name_gin ON brands USING gin(to_tsvector('english', name))",
		"CREATE INDEX IF NOT EXISTS idx_categories_name_gin ON categories USING gin(to_tsvector('english', name))",

		// Tags search (using existing tags table)
		"CREATE INDEX IF NOT EXISTS idx_tags_name_gin ON tags USING gin(to_tsvector('english', name))",

		// Composite indexes for common filter combinations
		"CREATE INDEX IF NOT EXISTS idx_products_category_status ON products(category_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_products_brand_status ON products(brand_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_products_price_status ON products(price, status)",
		"CREATE INDEX IF NOT EXISTS idx_products_featured_status ON products(featured, status)",

		// Search events table
		`CREATE TABLE IF NOT EXISTS search_events (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			query TEXT NOT NULL,
			user_id UUID,
			results_count INTEGER DEFAULT 0,
			clicked_product_id UUID,
			session_id VARCHAR(255),
			ip_address INET,
			user_agent TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// Search suggestions table
		`CREATE TABLE IF NOT EXISTS search_suggestions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			query TEXT NOT NULL UNIQUE,
			search_count INTEGER DEFAULT 1,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// Popular searches table
		`CREATE TABLE IF NOT EXISTS popular_searches (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			query TEXT NOT NULL,
			search_count INTEGER DEFAULT 1,
			period VARCHAR(20) DEFAULT 'daily',
			date TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			UNIQUE(query, period, date)
		)`,

		// Search filters table
		`CREATE TABLE IF NOT EXISTS search_filters (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL,
			name TEXT NOT NULL,
			query TEXT,
			filters JSONB,
			is_default BOOLEAN DEFAULT false,
			is_public BOOLEAN DEFAULT false,
			usage_count INTEGER DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// Search history table
		`CREATE TABLE IF NOT EXISTS search_history (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL,
			query TEXT NOT NULL,
			filters JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// Search events indexes
		"CREATE INDEX IF NOT EXISTS idx_search_events_query ON search_events(query)",
		"CREATE INDEX IF NOT EXISTS idx_search_events_user_id ON search_events(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_search_events_created_at ON search_events(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_search_events_query_gin ON search_events USING gin(to_tsvector('english', query))",

		// Search suggestions indexes
		"CREATE INDEX IF NOT EXISTS idx_search_suggestions_query ON search_suggestions(query)",
		"CREATE INDEX IF NOT EXISTS idx_search_suggestions_search_count ON search_suggestions(search_count)",
		"CREATE INDEX IF NOT EXISTS idx_search_suggestions_is_active ON search_suggestions(is_active)",

		// Popular searches indexes
		"CREATE INDEX IF NOT EXISTS idx_popular_searches_query ON popular_searches(query)",
		"CREATE INDEX IF NOT EXISTS idx_popular_searches_period ON popular_searches(period)",
		"CREATE INDEX IF NOT EXISTS idx_popular_searches_date ON popular_searches(date)",
		"CREATE INDEX IF NOT EXISTS idx_popular_searches_search_count ON popular_searches(search_count)",

		// Search filters indexes
		"CREATE INDEX IF NOT EXISTS idx_search_filters_user_id ON search_filters(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_search_filters_is_default ON search_filters(is_default)",
		"CREATE INDEX IF NOT EXISTS idx_search_filters_is_public ON search_filters(is_public)",

		// Search history indexes
		"CREATE INDEX IF NOT EXISTS idx_search_history_user_id ON search_history(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_search_history_created_at ON search_history(created_at)",
	}

	for _, sql := range sqls {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}

	return nil
}

// migration006Down removes enhanced search indexes
func migration006Down(db *gorm.DB) error {
	sqls := []string{
		"DROP INDEX IF EXISTS idx_products_featured",
		"DROP INDEX IF EXISTS idx_products_visibility",
		"DROP INDEX IF EXISTS idx_products_sku",
		"DROP INDEX IF EXISTS idx_products_sale_price",
		"DROP INDEX IF EXISTS idx_products_stock_status",
		"DROP INDEX IF EXISTS idx_products_search_vector",
		"DROP INDEX IF EXISTS idx_products_sku_gin",
		"DROP INDEX IF EXISTS idx_brands_name_gin",
		"DROP INDEX IF EXISTS idx_categories_name_gin",
		"DROP INDEX IF EXISTS idx_tags_name_gin",
		"DROP INDEX IF EXISTS idx_products_category_status",
		"DROP INDEX IF EXISTS idx_products_brand_status",
		"DROP INDEX IF EXISTS idx_products_price_status",
		"DROP INDEX IF EXISTS idx_products_featured_status",
		"DROP INDEX IF EXISTS idx_search_events_query",
		"DROP INDEX IF EXISTS idx_search_events_user_id",
		"DROP INDEX IF EXISTS idx_search_events_created_at",
		"DROP INDEX IF EXISTS idx_search_events_query_gin",
		"DROP INDEX IF EXISTS idx_search_suggestions_query",
		"DROP INDEX IF EXISTS idx_search_suggestions_search_count",
		"DROP INDEX IF EXISTS idx_search_suggestions_is_active",
		"DROP INDEX IF EXISTS idx_popular_searches_query",
		"DROP INDEX IF EXISTS idx_popular_searches_period",
		"DROP INDEX IF EXISTS idx_popular_searches_date",
		"DROP INDEX IF EXISTS idx_popular_searches_search_count",
		"DROP INDEX IF EXISTS idx_search_filters_user_id",
		"DROP INDEX IF EXISTS idx_search_filters_is_default",
		"DROP INDEX IF EXISTS idx_search_filters_is_public",
		"DROP INDEX IF EXISTS idx_search_history_user_id",
		"DROP INDEX IF EXISTS idx_search_history_created_at",
		"DROP TABLE IF EXISTS search_history",
		"DROP TABLE IF EXISTS search_filters",
		"DROP TABLE IF EXISTS popular_searches",
		"DROP TABLE IF EXISTS search_suggestions",
		"DROP TABLE IF EXISTS search_events",
	}

	for _, sql := range sqls {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}

	return nil
}
