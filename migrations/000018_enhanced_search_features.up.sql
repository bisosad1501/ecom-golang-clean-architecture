-- Enhanced Search Features Migration
-- This migration adds search analytics, synonyms, and enhanced search capabilities

-- Create search_analytics table for tracking search behavior
CREATE TABLE IF NOT EXISTS search_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    query VARCHAR(255) NOT NULL,
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
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for search_analytics
CREATE INDEX IF NOT EXISTS idx_search_analytics_query ON search_analytics(query);
CREATE INDEX IF NOT EXISTS idx_search_analytics_user_id ON search_analytics(user_id);
CREATE INDEX IF NOT EXISTS idx_search_analytics_session_id ON search_analytics(session_id);
CREATE INDEX IF NOT EXISTS idx_search_analytics_created_at ON search_analytics(created_at);
CREATE INDEX IF NOT EXISTS idx_search_analytics_search_type ON search_analytics(search_type);
CREATE INDEX IF NOT EXISTS idx_search_analytics_language ON search_analytics(language);

-- Create search_synonyms table for query expansion
CREATE TABLE IF NOT EXISTS search_synonyms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    term VARCHAR(255) NOT NULL,
    synonyms TEXT[] NOT NULL,
    language VARCHAR(10) DEFAULT 'en',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for search_synonyms
CREATE INDEX IF NOT EXISTS idx_search_synonyms_term ON search_synonyms(term);
CREATE INDEX IF NOT EXISTS idx_search_synonyms_language ON search_synonyms(language);
CREATE INDEX IF NOT EXISTS idx_search_synonyms_is_active ON search_synonyms(is_active);

-- Create GIN index for synonyms array
CREATE INDEX IF NOT EXISTS idx_search_synonyms_synonyms_gin ON search_synonyms USING GIN(synonyms);

-- Add enhanced search vector to products if not exists
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'products' AND column_name = 'search_vector_enhanced') THEN
        ALTER TABLE products ADD COLUMN search_vector_enhanced tsvector;
    END IF;
END $$;

-- Create enhanced search vector index
CREATE INDEX IF NOT EXISTS idx_products_search_vector_enhanced ON products USING GIN(search_vector_enhanced);

-- Function to update enhanced search vector
CREATE OR REPLACE FUNCTION update_enhanced_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    -- Enhanced search vector includes more fields and weights
    NEW.search_vector_enhanced := 
        setweight(to_tsvector('english', COALESCE(NEW.name, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.description, '')), 'B') ||
        setweight(to_tsvector('english', COALESCE(NEW.sku, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.short_description, '')), 'C') ||
        setweight(to_tsvector('english', 
            COALESCE((SELECT string_agg(name, ' ') FROM categories WHERE id = NEW.category_id), '')
        ), 'B') ||
        setweight(to_tsvector('english', 
            COALESCE((SELECT string_agg(name, ' ') FROM brands WHERE id = NEW.brand_id), '')
        ), 'B');
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for enhanced search vector
DROP TRIGGER IF EXISTS trigger_update_enhanced_search_vector ON products;
CREATE TRIGGER trigger_update_enhanced_search_vector
    BEFORE INSERT OR UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_enhanced_search_vector();

-- Update existing products with enhanced search vector
UPDATE products SET search_vector_enhanced = 
    setweight(to_tsvector('english', COALESCE(name, '')), 'A') ||
    setweight(to_tsvector('english', COALESCE(description, '')), 'B') ||
    setweight(to_tsvector('english', COALESCE(sku, '')), 'A') ||
    setweight(to_tsvector('english', COALESCE(short_description, '')), 'C');

-- Function for enhanced search with synonyms
CREATE OR REPLACE FUNCTION enhanced_search_with_synonyms(
    search_query TEXT,
    search_language VARCHAR(10) DEFAULT 'en'
)
RETURNS TABLE(
    product_id UUID,
    relevance_score FLOAT,
    match_type VARCHAR(50)
) AS $$
DECLARE
    expanded_query TEXT;
    synonym_terms TEXT[];
    term TEXT;
BEGIN
    -- Start with original query
    expanded_query := search_query;
    
    -- Get synonyms for query terms
    FOR term IN SELECT unnest(string_to_array(lower(search_query), ' '))
    LOOP
        SELECT array_agg(unnest(synonyms)) INTO synonym_terms
        FROM search_synonyms 
        WHERE lower(term) = ANY(synonyms) 
        AND language = search_language 
        AND is_active = true;
        
        IF synonym_terms IS NOT NULL THEN
            expanded_query := expanded_query || ' ' || array_to_string(synonym_terms, ' ');
        END IF;
    END LOOP;
    
    -- Return enhanced search results
    RETURN QUERY
    SELECT 
        p.id as product_id,
        ts_rank(p.search_vector_enhanced, plainto_tsquery('english', expanded_query)) as relevance_score,
        'enhanced_search' as match_type
    FROM products p
    WHERE p.search_vector_enhanced @@ plainto_tsquery('english', expanded_query)
    ORDER BY relevance_score DESC;
END;
$$ LANGUAGE plpgsql;

-- Function for search result highlighting
CREATE OR REPLACE FUNCTION get_search_highlights(
    content TEXT,
    search_query TEXT,
    max_fragments INTEGER DEFAULT 3,
    fragment_size INTEGER DEFAULT 150
)
RETURNS TEXT[] AS $$
DECLARE
    highlighted_fragments TEXT[];
    query_tsquery tsquery;
BEGIN
    -- Convert search query to tsquery
    query_tsquery := plainto_tsquery('english', search_query);
    
    -- Generate highlighted fragments
    SELECT array_agg(
        ts_headline('english', content, query_tsquery, 
            'MaxFragments=' || max_fragments || 
            ', FragmentDelimiter=|, MaxWords=' || fragment_size ||
            ', MinWords=10, StartSel=<mark>, StopSel=</mark>'
        )
    ) INTO highlighted_fragments;
    
    RETURN highlighted_fragments;
END;
$$ LANGUAGE plpgsql;

-- Add some sample synonyms
INSERT INTO search_synonyms (term, synonyms, language) VALUES
('phone', ARRAY['mobile', 'smartphone', 'cell phone', 'cellular'], 'en'),
('laptop', ARRAY['notebook', 'computer', 'pc'], 'en'),
('tv', ARRAY['television', 'monitor', 'display'], 'en'),
('headphones', ARRAY['earphones', 'earbuds', 'headset'], 'en'),
('watch', ARRAY['timepiece', 'smartwatch', 'wristwatch'], 'en')
ON CONFLICT DO NOTHING;

-- Create materialized view for search analytics summary
CREATE MATERIALIZED VIEW IF NOT EXISTS search_analytics_summary AS
SELECT 
    query,
    COUNT(*) as search_count,
    AVG(result_count) as avg_result_count,
    AVG(response_time_ms) as avg_response_time,
    COUNT(CASE WHEN clicked_result THEN 1 END) as click_count,
    ROUND(
        COUNT(CASE WHEN clicked_result THEN 1 END)::NUMERIC / COUNT(*)::NUMERIC * 100, 2
    ) as click_through_rate,
    DATE_TRUNC('day', created_at) as date
FROM search_analytics
GROUP BY query, DATE_TRUNC('day', created_at);

-- Create index on materialized view
CREATE INDEX IF NOT EXISTS idx_search_analytics_summary_query ON search_analytics_summary(query);
CREATE INDEX IF NOT EXISTS idx_search_analytics_summary_date ON search_analytics_summary(date);

-- Refresh materialized view function
CREATE OR REPLACE FUNCTION refresh_search_analytics_summary()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW search_analytics_summary;
END;
$$ LANGUAGE plpgsql;
