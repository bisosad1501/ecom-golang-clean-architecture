-- Enhanced Smart Autocomplete Features Migration
-- This migration adds smart features to the autocomplete system

-- Add new columns to autocomplete_entries table for smart features
ALTER TABLE autocomplete_entries ADD COLUMN IF NOT EXISTS is_trending BOOLEAN DEFAULT FALSE;
ALTER TABLE autocomplete_entries ADD COLUMN IF NOT EXISTS is_personalized BOOLEAN DEFAULT FALSE;
ALTER TABLE autocomplete_entries ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE autocomplete_entries ADD COLUMN IF NOT EXISTS synonyms TEXT[];
ALTER TABLE autocomplete_entries ADD COLUMN IF NOT EXISTS tags TEXT[];
ALTER TABLE autocomplete_entries ADD COLUMN IF NOT EXISTS score DECIMAL(10,4) DEFAULT 0;
ALTER TABLE autocomplete_entries ADD COLUMN IF NOT EXISTS language VARCHAR(10) DEFAULT 'en';
ALTER TABLE autocomplete_entries ADD COLUMN IF NOT EXISTS region VARCHAR(10);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_autocomplete_trending ON autocomplete_entries(is_trending, score DESC) WHERE is_trending = true;
CREATE INDEX IF NOT EXISTS idx_autocomplete_personalized ON autocomplete_entries(user_id, score DESC) WHERE is_personalized = true;
CREATE INDEX IF NOT EXISTS idx_autocomplete_language ON autocomplete_entries(language);
CREATE INDEX IF NOT EXISTS idx_autocomplete_score ON autocomplete_entries(score DESC);
CREATE INDEX IF NOT EXISTS idx_autocomplete_synonyms ON autocomplete_entries USING GIN(synonyms);
CREATE INDEX IF NOT EXISTS idx_autocomplete_tags ON autocomplete_entries USING GIN(tags);

-- Create composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_autocomplete_type_score ON autocomplete_entries(type, score DESC, is_active) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_autocomplete_user_type ON autocomplete_entries(user_id, type, updated_at DESC) WHERE user_id IS NOT NULL;

-- Create trigram indexes for fuzzy search
CREATE INDEX IF NOT EXISTS idx_autocomplete_value_trgm ON autocomplete_entries USING GIN(value gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_autocomplete_display_trgm ON autocomplete_entries USING GIN(display_text gin_trgm_ops);

-- Create function to update autocomplete scores
CREATE OR REPLACE FUNCTION update_autocomplete_scores()
RETURNS void AS $$
BEGIN
    UPDATE autocomplete_entries 
    SET score = (
        (search_count * 0.4) + 
        (click_count * 0.3) + 
        (priority * 0.2) + 
        (CASE 
            WHEN updated_at >= NOW() - INTERVAL '7 days' THEN 10
            WHEN updated_at >= NOW() - INTERVAL '30 days' THEN 5
            ELSE 0
        END * 0.1)
    )
    WHERE is_active = true;
END;
$$ LANGUAGE plpgsql;

-- Create function to update trending status
CREATE OR REPLACE FUNCTION update_trending_autocomplete()
RETURNS void AS $$
BEGIN
    -- Reset all trending flags
    UPDATE autocomplete_entries SET is_trending = false WHERE is_trending = true;
    
    -- Mark entries as trending based on recent activity (last 24 hours)
    UPDATE autocomplete_entries 
    SET is_trending = true
    WHERE updated_at >= NOW() - INTERVAL '24 hours' 
    AND (search_count > 10 OR click_count > 5)
    AND is_active = true;
END;
$$ LANGUAGE plpgsql;

-- Create function for smart autocomplete suggestions
CREATE OR REPLACE FUNCTION get_smart_autocomplete_suggestions(
    search_query TEXT,
    suggestion_types TEXT[] DEFAULT ARRAY['product', 'category', 'brand', 'query'],
    user_id_param UUID DEFAULT NULL,
    include_trending BOOLEAN DEFAULT FALSE,
    include_personalized BOOLEAN DEFAULT FALSE,
    include_popular BOOLEAN DEFAULT FALSE,
    language_param VARCHAR(10) DEFAULT 'en',
    result_limit INTEGER DEFAULT 10
)
RETURNS TABLE (
    id UUID,
    type VARCHAR(50),
    value TEXT,
    display_text TEXT,
    entity_id UUID,
    priority INTEGER,
    score DECIMAL(10,4),
    is_trending BOOLEAN,
    is_personalized BOOLEAN,
    metadata JSONB,
    synonyms TEXT[],
    tags TEXT[],
    reason TEXT
) AS $$
BEGIN
    RETURN QUERY
    WITH scored_suggestions AS (
        -- Exact matches (highest priority)
        SELECT 
            ae.id, ae.type, ae.value, ae.display_text, ae.entity_id,
            ae.priority, ae.score, ae.is_trending, ae.is_personalized,
            ae.metadata, ae.synonyms, ae.tags,
            'exact_match' as reason,
            100.0 as match_score
        FROM autocomplete_entries ae
        WHERE ae.is_active = true
        AND (array_length(suggestion_types, 1) IS NULL OR ae.type = ANY(suggestion_types))
        AND (language_param IS NULL OR ae.language = language_param)
        AND (ae.value ILIKE search_query OR ae.display_text ILIKE search_query)
        
        UNION ALL
        
        -- Fuzzy matches
        SELECT 
            ae.id, ae.type, ae.value, ae.display_text, ae.entity_id,
            ae.priority, ae.score, ae.is_trending, ae.is_personalized,
            ae.metadata, ae.synonyms, ae.tags,
            'fuzzy_match' as reason,
            similarity(ae.value, search_query) * 80.0 as match_score
        FROM autocomplete_entries ae
        WHERE ae.is_active = true
        AND (array_length(suggestion_types, 1) IS NULL OR ae.type = ANY(suggestion_types))
        AND (language_param IS NULL OR ae.language = language_param)
        AND (ae.value % search_query OR ae.display_text % search_query)
        
        UNION ALL
        
        -- Personalized suggestions (if user is authenticated)
        SELECT 
            ae.id, ae.type, ae.value, ae.display_text, ae.entity_id,
            ae.priority, ae.score, ae.is_trending, ae.is_personalized,
            ae.metadata, ae.synonyms, ae.tags,
            'personalized' as reason,
            90.0 as match_score
        FROM autocomplete_entries ae
        WHERE include_personalized = true
        AND user_id_param IS NOT NULL
        AND ae.user_id = user_id_param
        AND ae.is_active = true
        AND (array_length(suggestion_types, 1) IS NULL OR ae.type = ANY(suggestion_types))
        AND (search_query IS NULL OR ae.value ILIKE '%' || search_query || '%')
        
        UNION ALL
        
        -- Trending suggestions
        SELECT 
            ae.id, ae.type, ae.value, ae.display_text, ae.entity_id,
            ae.priority, ae.score, ae.is_trending, ae.is_personalized,
            ae.metadata, ae.synonyms, ae.tags,
            'trending' as reason,
            70.0 as match_score
        FROM autocomplete_entries ae
        WHERE include_trending = true
        AND ae.is_trending = true
        AND ae.is_active = true
        AND (array_length(suggestion_types, 1) IS NULL OR ae.type = ANY(suggestion_types))
        AND (search_query IS NULL OR ae.value ILIKE '%' || search_query || '%')
        
        UNION ALL
        
        -- Popular suggestions
        SELECT 
            ae.id, ae.type, ae.value, ae.display_text, ae.entity_id,
            ae.priority, ae.score, ae.is_trending, ae.is_personalized,
            ae.metadata, ae.synonyms, ae.tags,
            'popular' as reason,
            60.0 as match_score
        FROM autocomplete_entries ae
        WHERE include_popular = true
        AND ae.search_count > 50
        AND ae.is_active = true
        AND (array_length(suggestion_types, 1) IS NULL OR ae.type = ANY(suggestion_types))
        AND (search_query IS NULL OR ae.value ILIKE '%' || search_query || '%')
    )
    SELECT DISTINCT ON (ss.value, ss.type)
        ss.id, ss.type, ss.value, ss.display_text, ss.entity_id,
        ss.priority, ss.score, ss.is_trending, ss.is_personalized,
        ss.metadata, ss.synonyms, ss.tags, ss.reason
    FROM scored_suggestions ss
    ORDER BY ss.value, ss.type, ss.match_score DESC, ss.score DESC, ss.priority DESC
    LIMIT result_limit;
END;
$$ LANGUAGE plpgsql;

-- Create autocomplete analytics table for tracking interactions
CREATE TABLE IF NOT EXISTS autocomplete_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id UUID NOT NULL REFERENCES autocomplete_entries(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    session_id VARCHAR(255),
    interaction_type VARCHAR(50) NOT NULL, -- 'impression', 'click', 'search'
    query TEXT,
    position INTEGER, -- position in suggestion list
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for autocomplete analytics
CREATE INDEX IF NOT EXISTS idx_autocomplete_analytics_entry_id ON autocomplete_analytics(entry_id);
CREATE INDEX IF NOT EXISTS idx_autocomplete_analytics_user_id ON autocomplete_analytics(user_id);
CREATE INDEX IF NOT EXISTS idx_autocomplete_analytics_session_id ON autocomplete_analytics(session_id);
CREATE INDEX IF NOT EXISTS idx_autocomplete_analytics_interaction_type ON autocomplete_analytics(interaction_type);
CREATE INDEX IF NOT EXISTS idx_autocomplete_analytics_created_at ON autocomplete_analytics(created_at);

-- Insert some sample smart autocomplete data
INSERT INTO autocomplete_entries (type, value, display_text, priority, synonyms, tags, language, score) VALUES
('product', 'iPhone', 'iPhone - Apple Smartphone', 90, ARRAY['phone', 'smartphone', 'mobile'], ARRAY['electronics', 'apple', 'popular'], 'en', 95.0),
('product', 'Samsung Galaxy', 'Samsung Galaxy Series', 85, ARRAY['phone', 'smartphone', 'android'], ARRAY['electronics', 'samsung', 'popular'], 'en', 90.0),
('category', 'Electronics', 'Electronics & Technology', 80, ARRAY['tech', 'gadgets', 'devices'], ARRAY['category', 'popular'], 'en', 85.0),
('brand', 'Apple', 'Apple Inc.', 95, ARRAY['iphone', 'mac', 'ipad'], ARRAY['brand', 'premium'], 'en', 98.0),
('query', 'best smartphone', 'Best Smartphone Deals', 70, ARRAY['top phone', 'mobile deals'], ARRAY['query', 'popular'], 'en', 75.0)
ON CONFLICT DO NOTHING;

-- Create trigger to automatically update scores when entries are modified
CREATE OR REPLACE FUNCTION trigger_update_autocomplete_score()
RETURNS TRIGGER AS $$
BEGIN
    NEW.score = (
        (NEW.search_count * 0.4) + 
        (NEW.click_count * 0.3) + 
        (NEW.priority * 0.2) + 
        (CASE 
            WHEN NEW.updated_at >= NOW() - INTERVAL '7 days' THEN 10
            WHEN NEW.updated_at >= NOW() - INTERVAL '30 days' THEN 5
            ELSE 0
        END * 0.1)
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_autocomplete_score_trigger
    BEFORE UPDATE ON autocomplete_entries
    FOR EACH ROW
    EXECUTE FUNCTION trigger_update_autocomplete_score();

-- Schedule periodic updates (this would typically be done via cron or application scheduler)
-- For now, we'll create the functions that can be called manually or via scheduler

-- Function to cleanup old autocomplete analytics (keep last 90 days)
CREATE OR REPLACE FUNCTION cleanup_autocomplete_analytics()
RETURNS void AS $$
BEGIN
    DELETE FROM autocomplete_analytics 
    WHERE created_at < NOW() - INTERVAL '90 days';
END;
$$ LANGUAGE plpgsql;
