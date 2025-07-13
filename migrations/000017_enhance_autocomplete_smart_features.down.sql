-- Rollback Enhanced Smart Autocomplete Features Migration

-- Drop triggers
DROP TRIGGER IF EXISTS update_autocomplete_score_trigger ON autocomplete_entries;

-- Drop functions
DROP FUNCTION IF EXISTS trigger_update_autocomplete_score();
DROP FUNCTION IF EXISTS cleanup_autocomplete_analytics();
DROP FUNCTION IF EXISTS get_smart_autocomplete_suggestions(TEXT, TEXT[], UUID, BOOLEAN, BOOLEAN, BOOLEAN, VARCHAR(10), INTEGER);
DROP FUNCTION IF EXISTS update_trending_autocomplete();
DROP FUNCTION IF EXISTS update_autocomplete_scores();

-- Drop autocomplete analytics table
DROP TABLE IF EXISTS autocomplete_analytics;

-- Drop indexes
DROP INDEX IF EXISTS idx_autocomplete_trending;
DROP INDEX IF EXISTS idx_autocomplete_personalized;
DROP INDEX IF EXISTS idx_autocomplete_language;
DROP INDEX IF EXISTS idx_autocomplete_score;
DROP INDEX IF EXISTS idx_autocomplete_synonyms;
DROP INDEX IF EXISTS idx_autocomplete_tags;
DROP INDEX IF EXISTS idx_autocomplete_type_score;
DROP INDEX IF EXISTS idx_autocomplete_user_type;
DROP INDEX IF EXISTS idx_autocomplete_value_trgm;
DROP INDEX IF EXISTS idx_autocomplete_display_trgm;

-- Remove columns from autocomplete_entries table
ALTER TABLE autocomplete_entries DROP COLUMN IF EXISTS is_trending;
ALTER TABLE autocomplete_entries DROP COLUMN IF EXISTS is_personalized;
ALTER TABLE autocomplete_entries DROP COLUMN IF EXISTS user_id;
ALTER TABLE autocomplete_entries DROP COLUMN IF EXISTS synonyms;
ALTER TABLE autocomplete_entries DROP COLUMN IF EXISTS tags;
ALTER TABLE autocomplete_entries DROP COLUMN IF EXISTS score;
ALTER TABLE autocomplete_entries DROP COLUMN IF EXISTS language;
ALTER TABLE autocomplete_entries DROP COLUMN IF EXISTS region;
