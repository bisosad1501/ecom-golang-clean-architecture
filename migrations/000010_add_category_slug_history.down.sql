-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_category_seo_analytics ON categories;
DROP TRIGGER IF EXISTS trigger_track_category_slug_changes ON categories;

-- Drop functions
DROP FUNCTION IF EXISTS update_category_seo_analytics();
DROP FUNCTION IF EXISTS track_category_slug_changes();

-- Drop indexes
DROP INDEX IF EXISTS idx_category_slug_history_active_unique;
DROP INDEX IF EXISTS idx_category_seo_analytics_last_validated_at;
DROP INDEX IF EXISTS idx_category_seo_analytics_created_at;
DROP INDEX IF EXISTS idx_category_seo_analytics_seo_score;
DROP INDEX IF EXISTS idx_category_seo_analytics_category_id;
DROP INDEX IF EXISTS idx_category_redirects_is_active;
DROP INDEX IF EXISTS idx_category_redirects_category_id;
DROP INDEX IF EXISTS idx_category_redirects_to_slug;
DROP INDEX IF EXISTS idx_category_redirects_from_slug;
DROP INDEX IF EXISTS idx_category_slug_history_is_active;
DROP INDEX IF EXISTS idx_category_slug_history_created_at;
DROP INDEX IF EXISTS idx_category_slug_history_category_id;

-- Drop tables
DROP TABLE IF EXISTS category_seo_analytics;
DROP TABLE IF EXISTS category_redirects;
DROP TABLE IF EXISTS category_slug_history;
