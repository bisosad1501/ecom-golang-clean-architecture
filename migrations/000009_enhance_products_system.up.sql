-- Add new columns to products table for enhanced functionality
-- First add slug column as nullable
ALTER TABLE products
ADD COLUMN short_description TEXT,
ADD COLUMN slug VARCHAR(255),
ADD COLUMN meta_title VARCHAR(255),
ADD COLUMN meta_description TEXT,
ADD COLUMN keywords VARCHAR(500),
ADD COLUMN featured BOOLEAN DEFAULT FALSE,
ADD COLUMN visibility ENUM('visible', 'hidden', 'private') DEFAULT 'visible',
ADD COLUMN sale_price DECIMAL(10,2),
ADD COLUMN sale_start_date TIMESTAMP NULL,
ADD COLUMN sale_end_date TIMESTAMP NULL,
ADD COLUMN low_stock_threshold INT DEFAULT 5,
ADD COLUMN track_quantity BOOLEAN DEFAULT TRUE,
ADD COLUMN allow_backorder BOOLEAN DEFAULT FALSE,
ADD COLUMN stock_status ENUM('in_stock', 'out_of_stock', 'on_backorder', 'low_stock') DEFAULT 'in_stock',
ADD COLUMN requires_shipping BOOLEAN DEFAULT TRUE,
ADD COLUMN shipping_class VARCHAR(100),
ADD COLUMN tax_class VARCHAR(100) DEFAULT 'standard',
ADD COLUMN country_of_origin VARCHAR(100),
ADD COLUMN brand_id VARCHAR(36),
ADD COLUMN product_type ENUM('simple', 'variable', 'grouped', 'external') DEFAULT 'simple';

-- Add indexes for performance
CREATE INDEX idx_products_slug ON products(slug);
CREATE INDEX idx_products_featured ON products(featured);
CREATE INDEX idx_products_visibility ON products(visibility);
CREATE INDEX idx_products_brand_id ON products(brand_id);
CREATE INDEX idx_products_product_type ON products(product_type);
CREATE INDEX idx_products_stock_status ON products(stock_status);
CREATE INDEX idx_products_sale_price ON products(sale_price);

-- Create brands table
CREATE TABLE brands (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    logo VARCHAR(500),
    website VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_brands_name (name),
    INDEX idx_brands_slug (slug),
    INDEX idx_brands_is_active (is_active)
);

-- Create product_attributes table
CREATE TABLE product_attributes (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    slug VARCHAR(255) UNIQUE NOT NULL,
    type VARCHAR(50) DEFAULT 'text',
    description TEXT,
    position INT DEFAULT 0,
    is_required BOOLEAN DEFAULT FALSE,
    is_visible BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_product_attributes_name (name),
    INDEX idx_product_attributes_slug (slug),
    INDEX idx_product_attributes_type (type),
    INDEX idx_product_attributes_position (position)
);

-- Create product_attribute_terms table
CREATE TABLE product_attribute_terms (
    id VARCHAR(36) PRIMARY KEY,
    attribute_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    value VARCHAR(255),
    color VARCHAR(7),
    image VARCHAR(500),
    position INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_product_attribute_terms_attribute_id (attribute_id),
    INDEX idx_product_attribute_terms_name (name),
    INDEX idx_product_attribute_terms_slug (slug),
    INDEX idx_product_attribute_terms_position (position),
    
    FOREIGN KEY (attribute_id) REFERENCES product_attributes(id) ON DELETE CASCADE
);

-- Create product_attribute_values table
CREATE TABLE product_attribute_values (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL,
    attribute_id VARCHAR(36) NOT NULL,
    term_id VARCHAR(36),
    value VARCHAR(255),
    position INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_product_attribute_values_product_id (product_id),
    INDEX idx_product_attribute_values_attribute_id (attribute_id),
    INDEX idx_product_attribute_values_term_id (term_id),
    INDEX idx_product_attribute_values_position (position),
    
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (attribute_id) REFERENCES product_attributes(id) ON DELETE CASCADE,
    FOREIGN KEY (term_id) REFERENCES product_attribute_terms(id) ON DELETE SET NULL
);

-- Create product_variants table
CREATE TABLE product_variants (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL,
    sku VARCHAR(255) UNIQUE NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    compare_price DECIMAL(10,2),
    cost_price DECIMAL(10,2),
    stock INT DEFAULT 0,
    weight DECIMAL(8,3),
    length DECIMAL(8,3),
    width DECIMAL(8,3),
    height DECIMAL(8,3),
    image VARCHAR(500),
    position INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_product_variants_product_id (product_id),
    INDEX idx_product_variants_sku (sku),
    INDEX idx_product_variants_price (price),
    INDEX idx_product_variants_stock (stock),
    INDEX idx_product_variants_position (position),
    INDEX idx_product_variants_is_active (is_active),
    
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Create product_variant_attributes table
CREATE TABLE product_variant_attributes (
    id VARCHAR(36) PRIMARY KEY,
    variant_id VARCHAR(36) NOT NULL,
    attribute_id VARCHAR(36) NOT NULL,
    term_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_product_variant_attributes_variant_id (variant_id),
    INDEX idx_product_variant_attributes_attribute_id (attribute_id),
    INDEX idx_product_variant_attributes_term_id (term_id),
    
    UNIQUE KEY unique_variant_attribute (variant_id, attribute_id),
    
    FOREIGN KEY (variant_id) REFERENCES product_variants(id) ON DELETE CASCADE,
    FOREIGN KEY (attribute_id) REFERENCES product_attributes(id) ON DELETE CASCADE,
    FOREIGN KEY (term_id) REFERENCES product_attribute_terms(id) ON DELETE CASCADE
);

-- Create product_relations table for related products
CREATE TABLE product_relations (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL,
    related_product_id VARCHAR(36) NOT NULL,
    relation_type ENUM('related', 'upsell', 'cross_sell') DEFAULT 'related',
    position INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_product_relations_product_id (product_id),
    INDEX idx_product_relations_related_product_id (related_product_id),
    INDEX idx_product_relations_type (relation_type),
    INDEX idx_product_relations_position (position),
    
    UNIQUE KEY unique_product_relation (product_id, related_product_id, relation_type),
    
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (related_product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Add foreign key constraint for brand_id
ALTER TABLE products 
ADD CONSTRAINT fk_products_brand_id 
FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE SET NULL;

-- Update existing products to have slugs (temporary solution)
UPDATE products SET slug = LOWER(REPLACE(REPLACE(REPLACE(name, ' ', '-'), '.', ''), '/', '-')) WHERE slug IS NULL OR slug = '';

-- Now make slug NOT NULL and UNIQUE
ALTER TABLE products
ALTER COLUMN slug SET NOT NULL,
ADD CONSTRAINT products_slug_unique UNIQUE (slug);

-- Add some sample data
INSERT INTO brands (id, name, slug, description, is_active) VALUES
(UUID(), 'Apple', 'apple', 'Technology company', TRUE),
(UUID(), 'Samsung', 'samsung', 'Electronics manufacturer', TRUE),
(UUID(), 'Nike', 'nike', 'Sportswear brand', TRUE),
(UUID(), 'Adidas', 'adidas', 'Athletic apparel', TRUE),
(UUID(), 'Sony', 'sony', 'Electronics and entertainment', TRUE);

INSERT INTO product_attributes (id, name, slug, type, is_visible) VALUES
(UUID(), 'Color', 'color', 'color', TRUE),
(UUID(), 'Size', 'size', 'select', TRUE),
(UUID(), 'Material', 'material', 'text', TRUE),
(UUID(), 'Brand', 'brand', 'select', TRUE);
