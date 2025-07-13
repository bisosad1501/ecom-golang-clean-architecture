-- Down migration for Enhanced Search Features

-- Drop materialized view and related functions
DROP FUNCTION IF EXISTS refresh_search_analytics_summary();
DROP MATERIALIZED VIEW IF EXISTS search_analytics_summary;

-- Drop search functions
DROP FUNCTION IF EXISTS get_search_highlights(TEXT, TEXT, INTEGER, INTEGER);
DROP FUNCTION IF EXISTS enhanced_search_with_synonyms(TEXT, VARCHAR);

-- Drop triggers and functions
DROP TRIGGER IF EXISTS trigger_update_enhanced_search_vector ON products;
DROP FUNCTION IF EXISTS update_enhanced_search_vector();

-- Drop enhanced search vector column
ALTER TABLE products DROP COLUMN IF EXISTS search_vector_enhanced;

-- Drop search_synonyms table
DROP INDEX IF EXISTS idx_search_synonyms_synonyms_gin;
DROP INDEX IF EXISTS idx_search_synonyms_is_active;
DROP INDEX IF EXISTS idx_search_synonyms_language;
DROP INDEX IF EXISTS idx_search_synonyms_term;
DROP TABLE IF EXISTS search_synonyms;

-- Drop search_analytics table
DROP INDEX IF EXISTS idx_search_analytics_language;
DROP INDEX IF EXISTS idx_search_analytics_search_type;
DROP INDEX IF EXISTS idx_search_analytics_created_at;
DROP INDEX IF EXISTS idx_search_analytics_session_id;
DROP INDEX IF EXISTS idx_search_analytics_user_id;
DROP INDEX IF EXISTS idx_search_analytics_query;
DROP TABLE IF EXISTS search_analytics;
