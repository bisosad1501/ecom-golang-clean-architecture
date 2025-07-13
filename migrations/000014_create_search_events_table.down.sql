-- Drop search tables
DROP TABLE IF EXISTS search_filters;
DROP TABLE IF EXISTS search_history;
DROP TABLE IF EXISTS search_events;

-- Drop triggers
DROP TRIGGER IF EXISTS update_search_events_updated_at ON search_events;
DROP TRIGGER IF EXISTS update_search_history_updated_at ON search_history;
DROP TRIGGER IF EXISTS update_search_filters_updated_at ON search_filters;
