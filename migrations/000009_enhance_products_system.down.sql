-- Remove foreign key constraint
ALTER TABLE products DROP FOREIGN KEY fk_products_brand_id;

-- Drop new tables (in reverse order of creation)
DROP TABLE IF EXISTS product_relations;
DROP TABLE IF EXISTS product_variant_attributes;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS product_attribute_values;
DROP TABLE IF EXISTS product_attribute_terms;
DROP TABLE IF EXISTS product_attributes;
DROP TABLE IF EXISTS brands;

-- Drop indexes
DROP INDEX IF EXISTS idx_products_slug ON products;
DROP INDEX IF EXISTS idx_products_featured ON products;
DROP INDEX IF EXISTS idx_products_visibility ON products;
DROP INDEX IF EXISTS idx_products_brand_id ON products;
DROP INDEX IF EXISTS idx_products_product_type ON products;
DROP INDEX IF EXISTS idx_products_stock_status ON products;
DROP INDEX IF EXISTS idx_products_sale_price ON products;

-- Remove new columns from products table
ALTER TABLE products 
DROP COLUMN IF EXISTS short_description,
DROP COLUMN IF EXISTS slug,
DROP COLUMN IF EXISTS meta_title,
DROP COLUMN IF EXISTS meta_description,
DROP COLUMN IF EXISTS keywords,
DROP COLUMN IF EXISTS featured,
DROP COLUMN IF EXISTS visibility,
DROP COLUMN IF EXISTS sale_price,
DROP COLUMN IF EXISTS sale_start_date,
DROP COLUMN IF EXISTS sale_end_date,
DROP COLUMN IF EXISTS low_stock_threshold,
DROP COLUMN IF EXISTS track_quantity,
DROP COLUMN IF EXISTS allow_backorder,
DROP COLUMN IF EXISTS stock_status,
DROP COLUMN IF EXISTS requires_shipping,
DROP COLUMN IF EXISTS shipping_class,
DROP COLUMN IF EXISTS tax_class,
DROP COLUMN IF EXISTS country_of_origin,
DROP COLUMN IF EXISTS brand_id,
DROP COLUMN IF EXISTS product_type;
