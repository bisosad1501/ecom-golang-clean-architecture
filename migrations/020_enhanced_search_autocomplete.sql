-- Enhanced Search Autocomplete Migration
-- This migration adds enhanced autocomplete and search features

-- Create autocomplete_entries table
CREATE TABLE IF NOT EXISTS autocomplete_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(50) NOT NULL,
    value TEXT NOT NULL,
    display_text TEXT NOT NULL,
    entity_id UUID,
    priority INTEGER DEFAULT 0,
    search_count INTEGER DEFAULT 0,
    click_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for autocomplete_entries
CREATE INDEX IF NOT EXISTS idx_autocomplete_entries_type ON autocomplete_entries(type);
CREATE INDEX IF NOT EXISTS idx_autocomplete_entries_value ON autocomplete_entries(value);
CREATE INDEX IF NOT EXISTS idx_autocomplete_entries_entity_id ON autocomplete_entries(entity_id);
CREATE INDEX IF NOT EXISTS idx_autocomplete_entries_priority ON autocomplete_entries(priority DESC);
CREATE INDEX IF NOT EXISTS idx_autocomplete_entries_search_count ON autocomplete_entries(search_count DESC);
CREATE INDEX IF NOT EXISTS idx_autocomplete_entries_is_active ON autocomplete_entries(is_active);
CREATE INDEX IF NOT EXISTS idx_autocomplete_entries_value_text ON autocomplete_entries USING gin(to_tsvector('english', value));
CREATE INDEX IF NOT EXISTS idx_autocomplete_entries_display_text ON autocomplete_entries USING gin(to_tsvector('english', display_text));

-- Create search_trends table
CREATE TABLE IF NOT EXISTS search_trends (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    query TEXT NOT NULL,
    search_count INTEGER DEFAULT 0,
    period VARCHAR(20) NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for search_trends
CREATE INDEX IF NOT EXISTS idx_search_trends_query ON search_trends(query);
CREATE INDEX IF NOT EXISTS idx_search_trends_period ON search_trends(period);
CREATE INDEX IF NOT EXISTS idx_search_trends_date ON search_trends(date);
CREATE INDEX IF NOT EXISTS idx_search_trends_search_count ON search_trends(search_count DESC);
CREATE UNIQUE INDEX IF NOT EXISTS idx_search_trends_unique ON search_trends(query, period, date);

-- Create user_search_preferences table
CREATE TABLE IF NOT EXISTS user_search_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    preferred_categories JSONB DEFAULT '[]',
    preferred_brands JSONB DEFAULT '[]',
    search_language VARCHAR(10) DEFAULT 'en',
    autocomplete_enabled BOOLEAN DEFAULT true,
    search_history_enabled BOOLEAN DEFAULT true,
    personalized_results BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for user_search_preferences
CREATE INDEX IF NOT EXISTS idx_user_search_preferences_user_id ON user_search_preferences(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_search_preferences_unique ON user_search_preferences(user_id);

-- Create search_sessions table
CREATE TABLE IF NOT EXISTS search_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(255) NOT NULL,
    user_id UUID,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    search_count INTEGER DEFAULT 0,
    click_count INTEGER DEFAULT 0,
    conversion_count INTEGER DEFAULT 0,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for search_sessions
CREATE INDEX IF NOT EXISTS idx_search_sessions_session_id ON search_sessions(session_id);
CREATE INDEX IF NOT EXISTS idx_search_sessions_user_id ON search_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_search_sessions_start_time ON search_sessions(start_time);

-- Create function to update autocomplete entries from products
CREATE OR REPLACE FUNCTION update_autocomplete_from_products()
RETURNS void AS $$
BEGIN
    -- Insert product suggestions
    INSERT INTO autocomplete_entries (type, value, display_text, entity_id, priority, metadata)
    SELECT 
        'product',
        p.name,
        p.name,
        p.id,
        50,
        json_build_object('product_id', p.id, 'price', p.price, 'image', COALESCE(pi.url, ''))::jsonb
    FROM products p
    LEFT JOIN product_images pi ON p.id = pi.product_id AND pi.position = 0
    WHERE p.status = 'active'
    ON CONFLICT DO NOTHING;
END;
$$ LANGUAGE plpgsql;

-- Create function to update autocomplete entries from categories
CREATE OR REPLACE FUNCTION update_autocomplete_from_categories()
RETURNS void AS $$
BEGIN
    -- Insert category suggestions
    INSERT INTO autocomplete_entries (type, value, display_text, entity_id, priority, metadata)
    SELECT 
        'category',
        c.name,
        c.name,
        c.id,
        70,
        json_build_object('category_id', c.id, 'slug', c.slug, 'image', COALESCE(c.image, ''))::jsonb
    FROM categories c
    WHERE c.is_active = true
    ON CONFLICT DO NOTHING;
END;
$$ LANGUAGE plpgsql;

-- Create function to update autocomplete entries from brands
CREATE OR REPLACE FUNCTION update_autocomplete_from_brands()
RETURNS void AS $$
BEGIN
    -- Insert brand suggestions
    INSERT INTO autocomplete_entries (type, value, display_text, entity_id, priority, metadata)
    SELECT 
        'brand',
        b.name,
        b.name,
        b.id,
        60,
        json_build_object('brand_id', b.id, 'slug', b.slug, 'logo', COALESCE(b.logo, ''))::jsonb
    FROM brands b
    WHERE b.is_active = true
    ON CONFLICT DO NOTHING;
END;
$$ LANGUAGE plpgsql;

-- Create function to get enhanced autocomplete suggestions
CREATE OR REPLACE FUNCTION get_enhanced_autocomplete_suggestions(
    search_query TEXT,
    suggestion_types TEXT[] DEFAULT ARRAY['product', 'category', 'brand', 'query'],
    result_limit INTEGER DEFAULT 10
)
RETURNS TABLE (
    id UUID,
    type TEXT,
    value TEXT,
    display_text TEXT,
    entity_id UUID,
    priority INTEGER,
    search_count INTEGER,
    click_count INTEGER,
    metadata JSONB,
    relevance_score FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        ae.id,
        ae.type,
        ae.value,
        ae.display_text,
        ae.entity_id,
        ae.priority,
        ae.search_count,
        ae.click_count,
        ae.metadata,
        CASE 
            WHEN ae.value ILIKE search_query || '%' THEN 1.0
            WHEN ae.value ILIKE '%' || search_query || '%' THEN 0.8
            WHEN ae.display_text ILIKE search_query || '%' THEN 0.9
            WHEN ae.display_text ILIKE '%' || search_query || '%' THEN 0.7
            ELSE 0.5
        END as relevance_score
    FROM autocomplete_entries ae
    WHERE ae.is_active = true
        AND ae.type = ANY(suggestion_types)
        AND (
            ae.value ILIKE '%' || search_query || '%' 
            OR ae.display_text ILIKE '%' || search_query || '%'
            OR search_query = ''
        )
    ORDER BY 
        relevance_score DESC,
        ae.priority DESC,
        ae.search_count DESC,
        ae.click_count DESC
    LIMIT result_limit;
END;
$$ LANGUAGE plpgsql;

-- Create function to update search trends
CREATE OR REPLACE FUNCTION update_search_trend(
    search_query TEXT,
    trend_period TEXT DEFAULT 'daily'
)
RETURNS void AS $$
DECLARE
    trend_date DATE;
BEGIN
    trend_date := CURRENT_DATE;
    
    -- Update or insert search trend
    INSERT INTO search_trends (query, search_count, period, date)
    VALUES (search_query, 1, trend_period, trend_date)
    ON CONFLICT (query, period, date)
    DO UPDATE SET 
        search_count = search_trends.search_count + 1;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to update autocomplete entries when products change
CREATE OR REPLACE FUNCTION trigger_update_product_autocomplete()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        -- Update or insert autocomplete entry for this product
        INSERT INTO autocomplete_entries (type, value, display_text, entity_id, priority, metadata)
        VALUES (
            'product',
            NEW.name,
            NEW.name,
            NEW.id,
            50,
            json_build_object('product_id', NEW.id, 'price', NEW.price)::jsonb
        )
        ON CONFLICT (type, entity_id) 
        DO UPDATE SET 
            value = NEW.name,
            display_text = NEW.name,
            metadata = json_build_object('product_id', NEW.id, 'price', NEW.price)::jsonb,
            updated_at = CURRENT_TIMESTAMP;
            
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        -- Remove autocomplete entry for this product
        DELETE FROM autocomplete_entries 
        WHERE type = 'product' AND entity_id = OLD.id;
        
        RETURN OLD;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
DROP TRIGGER IF EXISTS trigger_product_autocomplete ON products;
CREATE TRIGGER trigger_product_autocomplete
    AFTER INSERT OR UPDATE OR DELETE ON products
    FOR EACH ROW
    EXECUTE FUNCTION trigger_update_product_autocomplete();

-- Initial data population
SELECT update_autocomplete_from_products();
SELECT update_autocomplete_from_categories();
SELECT update_autocomplete_from_brands();

-- Add some initial trending searches (example data)
INSERT INTO search_trends (query, search_count, period, date) VALUES
('phone', 150, 'daily', CURRENT_DATE),
('laptop', 120, 'daily', CURRENT_DATE),
('headphones', 90, 'daily', CURRENT_DATE),
('watch', 80, 'daily', CURRENT_DATE),
('camera', 70, 'daily', CURRENT_DATE)
ON CONFLICT DO NOTHING;
