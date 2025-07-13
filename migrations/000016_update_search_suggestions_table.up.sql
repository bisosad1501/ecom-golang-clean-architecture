-- Update search_suggestions table to match new entity structure

-- Add new columns to search_suggestions table
ALTER TABLE search_suggestions ADD COLUMN IF NOT EXISTS frequency INTEGER DEFAULT 1;
ALTER TABLE search_suggestions ADD COLUMN IF NOT EXISTS result_count INTEGER DEFAULT 0;
ALTER TABLE search_suggestions ADD COLUMN IF NOT EXISTS last_searched TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE search_suggestions ADD COLUMN IF NOT EXISTS is_trending BOOLEAN DEFAULT FALSE;

-- Update existing data
UPDATE search_suggestions SET 
    frequency = COALESCE(search_count, 1),
    last_searched = COALESCE(updated_at, created_at)
WHERE frequency IS NULL OR last_searched IS NULL;

-- Create indexes for new columns
CREATE INDEX IF NOT EXISTS idx_search_suggestions_frequency_new ON search_suggestions(frequency DESC);
CREATE INDEX IF NOT EXISTS idx_search_suggestions_last_searched ON search_suggestions(last_searched DESC);
CREATE INDEX IF NOT EXISTS idx_search_suggestions_trending ON search_suggestions(is_trending, frequency DESC);
CREATE INDEX IF NOT EXISTS idx_search_suggestions_result_count ON search_suggestions(result_count DESC);
