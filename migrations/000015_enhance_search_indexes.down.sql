-- Rollback enhanced search indexes and full-text search optimization

-- Drop functions
DROP FUNCTION IF EXISTS get_search_suggestions_with_synonyms(TEXT, INTEGER);
DROP FUNCTION IF EXISTS update_search_suggestion(TEXT, INTEGER);

-- Drop triggers
DROP TRIGGER IF EXISTS update_product_search_vector_trigger ON products;
DROP FUNCTION IF EXISTS update_product_search_vector();

-- Drop tables
DROP TABLE IF EXISTS search_suggestions;
DROP TABLE IF EXISTS search_analytics;
DROP TABLE IF EXISTS search_synonyms;

-- Drop indexes
DROP INDEX IF EXISTS idx_search_suggestions_trending;
DROP INDEX IF EXISTS idx_search_suggestions_frequency;
DROP INDEX IF EXISTS idx_search_suggestions_query;
DROP INDEX IF EXISTS idx_search_analytics_ctr;
DROP INDEX IF EXISTS idx_search_analytics_date;
DROP INDEX IF EXISTS idx_search_analytics_query;
DROP INDEX IF EXISTS idx_search_synonyms_active;
DROP INDEX IF EXISTS idx_search_synonyms_term;
DROP INDEX IF EXISTS idx_products_price_range;
DROP INDEX IF EXISTS idx_products_featured_search;
DROP INDEX IF EXISTS idx_products_search_composite;
DROP INDEX IF EXISTS idx_products_sku_trgm;
DROP INDEX IF EXISTS idx_products_name_trgm;
DROP INDEX IF EXISTS idx_products_description_gin;
DROP INDEX IF EXISTS idx_products_name_gin;
DROP INDEX IF EXISTS idx_products_search_vector;

-- Drop column
ALTER TABLE products DROP COLUMN IF EXISTS search_vector;

-- Drop extension (be careful with this in production)
-- DROP EXTENSION IF EXISTS pg_trgm;
