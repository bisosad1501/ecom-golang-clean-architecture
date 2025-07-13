-- Enhance existing search_analytics table

-- Add missing columns to search_analytics
ALTER TABLE search_analytics 
ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS session_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS ip_address INET,
ADD COLUMN IF NOT EXISTS user_agent TEXT,
ADD COLUMN IF NOT EXISTS clicked_result BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS click_position INTEGER,
ADD COLUMN IF NOT EXISTS response_time_ms INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS search_type VARCHAR(50) DEFAULT 'full_text',
ADD COLUMN IF NOT EXISTS filters JSONB,
ADD COLUMN IF NOT EXISTS sort_by VARCHAR(50),
ADD COLUMN IF NOT EXISTS language VARCHAR(10) DEFAULT 'en';

-- Create missing indexes
CREATE INDEX IF NOT EXISTS idx_search_analytics_user_id ON search_analytics(user_id);
CREATE INDEX IF NOT EXISTS idx_search_analytics_session_id ON search_analytics(session_id);
CREATE INDEX IF NOT EXISTS idx_search_analytics_search_type ON search_analytics(search_type);
CREATE INDEX IF NOT EXISTS idx_search_analytics_language ON search_analytics(language);

-- Fix search_synonyms table structure
DO $$
BEGIN
    -- Check if language column exists, if not add it
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'search_synonyms' AND column_name = 'language') THEN
        ALTER TABLE search_synonyms ADD COLUMN language VARCHAR(10) DEFAULT 'en';
        CREATE INDEX idx_search_synonyms_language ON search_synonyms(language);
    END IF;
END $$;

-- Insert sample synonyms with proper error handling
DO $$
BEGIN
    INSERT INTO search_synonyms (term, synonyms, language) VALUES
    ('phone', ARRAY['mobile', 'smartphone', 'cell phone', 'cellular'], 'en'),
    ('laptop', ARRAY['notebook', 'computer', 'pc'], 'en'),
    ('tv', ARRAY['television', 'monitor', 'display'], 'en'),
    ('headphones', ARRAY['earphones', 'earbuds', 'headset'], 'en'),
    ('watch', ARRAY['timepiece', 'smartwatch', 'wristwatch'], 'en')
    ON CONFLICT (term) DO NOTHING;
EXCEPTION
    WHEN others THEN
        -- If there's any error, just continue
        NULL;
END $$;

-- Create enhanced search analytics summary view
DROP MATERIALIZED VIEW IF EXISTS search_analytics_summary;
CREATE MATERIALIZED VIEW search_analytics_summary AS
SELECT 
    query,
    COUNT(*) as search_count,
    AVG(result_count) as avg_result_count,
    COALESCE(AVG(response_time_ms), 0) as avg_response_time,
    COUNT(CASE WHEN clicked_result THEN 1 END) as click_count,
    ROUND(
        COALESCE(
            COUNT(CASE WHEN clicked_result THEN 1 END)::NUMERIC / NULLIF(COUNT(*), 0)::NUMERIC * 100, 
            0
        ), 2
    ) as click_through_rate,
    DATE_TRUNC('day', created_at) as date
FROM search_analytics
GROUP BY query, DATE_TRUNC('day', created_at);

-- Create indexes on materialized view
CREATE INDEX IF NOT EXISTS idx_search_analytics_summary_query ON search_analytics_summary(query);
CREATE INDEX IF NOT EXISTS idx_search_analytics_summary_date ON search_analytics_summary(date);

-- Create enhanced search analytics table
CREATE TABLE IF NOT EXISTS enhanced_search_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    query TEXT NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    session_id VARCHAR(255),
    ip_address INET,
    user_agent TEXT,
    result_count INTEGER DEFAULT 0,
    clicked_result BOOLEAN DEFAULT FALSE,
    click_position INTEGER,
    response_time_ms INTEGER DEFAULT 0,
    search_type VARCHAR(50) DEFAULT 'full_text',
    filters JSONB,
    sort_by VARCHAR(50),
    language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for enhanced search analytics
CREATE INDEX IF NOT EXISTS idx_enhanced_search_analytics_query ON enhanced_search_analytics(query);
CREATE INDEX IF NOT EXISTS idx_enhanced_search_analytics_user_id ON enhanced_search_analytics(user_id);
CREATE INDEX IF NOT EXISTS idx_enhanced_search_analytics_session_id ON enhanced_search_analytics(session_id);
CREATE INDEX IF NOT EXISTS idx_enhanced_search_analytics_search_type ON enhanced_search_analytics(search_type);
CREATE INDEX IF NOT EXISTS idx_enhanced_search_analytics_language ON enhanced_search_analytics(language);
CREATE INDEX IF NOT EXISTS idx_enhanced_search_analytics_created_at ON enhanced_search_analytics(created_at);

-- Function to track enhanced search analytics
CREATE OR REPLACE FUNCTION track_enhanced_search_analytics(
    p_query TEXT,
    p_user_id UUID DEFAULT NULL,
    p_session_id VARCHAR(255) DEFAULT NULL,
    p_ip_address INET DEFAULT NULL,
    p_user_agent TEXT DEFAULT NULL,
    p_result_count INTEGER DEFAULT 0,
    p_response_time_ms INTEGER DEFAULT 0,
    p_search_type VARCHAR(50) DEFAULT 'full_text',
    p_filters JSONB DEFAULT NULL,
    p_sort_by VARCHAR(50) DEFAULT NULL,
    p_language VARCHAR(10) DEFAULT 'en'
)
RETURNS UUID AS $$
DECLARE
    analytics_id UUID;
BEGIN
    INSERT INTO enhanced_search_analytics (
        query, user_id, session_id, ip_address, user_agent,
        result_count, response_time_ms, search_type, filters, sort_by, language
    ) VALUES (
        p_query, p_user_id, p_session_id, p_ip_address, p_user_agent,
        p_result_count, p_response_time_ms, p_search_type, p_filters, p_sort_by, p_language
    ) RETURNING id INTO analytics_id;

    RETURN analytics_id;
END;
$$ LANGUAGE plpgsql;

-- Function to track search clicks
CREATE OR REPLACE FUNCTION track_enhanced_search_click(
    p_analytics_id UUID,
    p_click_position INTEGER
)
RETURNS void AS $$
BEGIN
    UPDATE enhanced_search_analytics
    SET clicked_result = TRUE, click_position = p_click_position
    WHERE id = p_analytics_id;
END;
$$ LANGUAGE plpgsql;
