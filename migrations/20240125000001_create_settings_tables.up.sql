-- Create settings table
CREATE TABLE IF NOT EXISTS settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(255) UNIQUE NOT NULL,
    value TEXT,
    type VARCHAR(50) NOT NULL DEFAULT 'string',
    category VARCHAR(100) NOT NULL,
    description TEXT,
    is_public BOOLEAN DEFAULT FALSE,
    is_encrypted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on category for faster queries
CREATE INDEX IF NOT EXISTS idx_settings_category ON settings(category);
CREATE INDEX IF NOT EXISTS idx_settings_key ON settings(key);

-- Create tax_rates table
CREATE TABLE IF NOT EXISTS tax_rates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    rate DECIMAL(5,4) NOT NULL,
    country CHAR(2) NOT NULL,
    state VARCHAR(50),
    city VARCHAR(100),
    zip_code VARCHAR(20),
    tax_class VARCHAR(50),
    priority INTEGER DEFAULT 0,
    compound BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for tax rate lookups
CREATE INDEX IF NOT EXISTS idx_tax_rates_location ON tax_rates(country, state, city, zip_code);
CREATE INDEX IF NOT EXISTS idx_tax_rates_active ON tax_rates(is_active);

-- Create tax_classes table
CREATE TABLE IF NOT EXISTS tax_classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create shipping_zones table
CREATE TABLE IF NOT EXISTS shipping_zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    countries TEXT, -- JSON array
    states TEXT,    -- JSON array
    zip_codes TEXT, -- JSON array
    is_active BOOLEAN DEFAULT TRUE,
    priority INTEGER DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create shipping_methods table
CREATE TABLE IF NOT EXISTS shipping_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- flat_rate, free_shipping, local_pickup
    cost DECIMAL(10,2) DEFAULT 0,
    min_amount DECIMAL(10,2) DEFAULT 0,
    max_amount DECIMAL(10,2) DEFAULT 0,
    zone_id UUID REFERENCES shipping_zones(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT TRUE,
    settings TEXT, -- JSON
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create shipping_classes table
CREATE TABLE IF NOT EXISTS shipping_classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    cost DECIMAL(10,2) DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create packaging_options table
CREATE TABLE IF NOT EXISTS packaging_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    length DECIMAL(8,2) DEFAULT 0,
    width DECIMAL(8,2) DEFAULT 0,
    height DECIMAL(8,2) DEFAULT 0,
    weight DECIMAL(8,2) DEFAULT 0,
    cost DECIMAL(10,2) DEFAULT 0,
    is_default BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create email_templates table
CREATE TABLE IF NOT EXISTS email_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- newsletter, promotional, transactional
    content TEXT NOT NULL,
    plain_text TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    settings TEXT, -- JSON
    variables TEXT, -- JSON
    usage_count BIGINT DEFAULT 0,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create campaigns table
CREATE TABLE IF NOT EXISTS campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL, -- email, social, display, affiliate
    status VARCHAR(50) NOT NULL DEFAULT 'draft', -- draft, active, paused, completed
    budget DECIMAL(12,2) DEFAULT 0,
    spent DECIMAL(12,2) DEFAULT 0,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE,
    target_audience TEXT, -- JSON
    goals TEXT, -- JSON
    settings TEXT, -- JSON
    metrics TEXT, -- JSON
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create promotions table
CREATE TABLE IF NOT EXISTS promotions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL, -- percentage, fixed, buy_x_get_y, free_shipping
    value DECIMAL(10,2) NOT NULL,
    min_order_amount DECIMAL(10,2) DEFAULT 0,
    max_discount DECIMAL(10,2) DEFAULT 0,
    usage_limit INTEGER DEFAULT 0,
    usage_count INTEGER DEFAULT 0,
    usage_per_customer INTEGER DEFAULT 0,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(50) DEFAULT 'active', -- active, inactive, scheduled, expired
    applicable_products TEXT, -- JSON array of UUIDs
    applicable_categories TEXT, -- JSON array of UUIDs
    excluded_products TEXT, -- JSON array of UUIDs
    excluded_categories TEXT, -- JSON array of UUIDs
    customer_segments TEXT, -- JSON array
    settings TEXT, -- JSON
    metrics TEXT, -- JSON
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create newsletter_subscribers table
CREATE TABLE IF NOT EXISTS newsletter_subscribers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- active, inactive, pending
    segments TEXT, -- JSON array
    tags TEXT, -- JSON array
    custom_fields TEXT, -- JSON
    subscribed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    unsubscribed_at TIMESTAMP WITH TIME ZONE,
    last_email_sent TIMESTAMP WITH TIME ZONE,
    emails_sent BIGINT DEFAULT 0,
    emails_opened BIGINT DEFAULT 0,
    emails_clicked BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create email_campaigns table
CREATE TABLE IF NOT EXISTS email_campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    template_id UUID REFERENCES email_templates(id) ON DELETE SET NULL,
    recipients TEXT, -- JSON
    status VARCHAR(50) NOT NULL DEFAULT 'draft', -- draft, scheduled, sending, sent, failed
    scheduled_at TIMESTAMP WITH TIME ZONE,
    sent_at TIMESTAMP WITH TIME ZONE,
    variables TEXT, -- JSON
    settings TEXT, -- JSON
    metrics TEXT, -- JSON
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create loyalty_programs table
CREATE TABLE IF NOT EXISTS loyalty_programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL, -- points, cashback, tier
    is_active BOOLEAN DEFAULT TRUE,
    rules TEXT, -- JSON
    rewards TEXT, -- JSON
    settings TEXT, -- JSON
    member_count BIGINT DEFAULT 0,
    total_points_issued BIGINT DEFAULT 0,
    total_points_redeemed BIGINT DEFAULT 0,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create referral_programs table
CREATE TABLE IF NOT EXISTS referral_programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    referrer_reward TEXT, -- JSON
    referee_reward TEXT, -- JSON
    rules TEXT, -- JSON
    settings TEXT, -- JSON
    stats TEXT, -- JSON
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_email_templates_type ON email_templates(type);
CREATE INDEX IF NOT EXISTS idx_email_templates_active ON email_templates(is_active);
CREATE INDEX IF NOT EXISTS idx_campaigns_status ON campaigns(status);
CREATE INDEX IF NOT EXISTS idx_campaigns_type ON campaigns(type);
CREATE INDEX IF NOT EXISTS idx_promotions_status ON promotions(status);
CREATE INDEX IF NOT EXISTS idx_promotions_dates ON promotions(start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_newsletter_subscribers_email ON newsletter_subscribers(email);
CREATE INDEX IF NOT EXISTS idx_newsletter_subscribers_status ON newsletter_subscribers(status);
CREATE INDEX IF NOT EXISTS idx_email_campaigns_status ON email_campaigns(status);
CREATE INDEX IF NOT EXISTS idx_loyalty_programs_active ON loyalty_programs(is_active);
CREATE INDEX IF NOT EXISTS idx_referral_programs_active ON referral_programs(is_active);

-- Insert default settings
INSERT INTO settings (key, value, type, category, description) VALUES
('store_name', 'My Store', 'string', 'general', 'Store name'),
('store_description', 'Welcome to our store', 'string', 'general', 'Store description'),
('store_email', 'admin@store.com', 'string', 'general', 'Store email'),
('currency', 'USD', 'string', 'general', 'Default currency'),
('timezone', 'UTC', 'string', 'general', 'Store timezone'),
('language', 'en', 'string', 'general', 'Default language'),
('maintenance_mode', 'false', 'boolean', 'store', 'Maintenance mode'),
('allow_registration', 'true', 'boolean', 'store', 'Allow user registration'),
('tax_enabled', 'true', 'boolean', 'tax', 'Enable tax calculation'),
('shipping_enabled', 'true', 'boolean', 'shipping', 'Enable shipping'),
('email_notifications', 'true', 'boolean', 'notifications', 'Enable email notifications')
ON CONFLICT (key) DO NOTHING;
