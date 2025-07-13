-- Rollback search_suggestions table updates

-- Drop indexes
DROP INDEX IF EXISTS idx_search_suggestions_result_count;
DROP INDEX IF EXISTS idx_search_suggestions_trending;
DROP INDEX IF EXISTS idx_search_suggestions_last_searched;
DROP INDEX IF EXISTS idx_search_suggestions_frequency_new;

-- Drop columns
ALTER TABLE search_suggestions DROP COLUMN IF EXISTS is_trending;
ALTER TABLE search_suggestions DROP COLUMN IF EXISTS last_searched;
ALTER TABLE search_suggestions DROP COLUMN IF EXISTS result_count;
ALTER TABLE search_suggestions DROP COLUMN IF EXISTS frequency;
