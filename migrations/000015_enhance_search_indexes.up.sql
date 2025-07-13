-- Enhanced search indexes and full-text search optimization
-- Create GIN indexes for full-text search performance

-- Create search vector column for products
ALTER TABLE products ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- Create function to update search vector
CREATE OR REPLACE FUNCTION update_product_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector := to_tsvector('english', 
        coalesce(NEW.name, '') || ' ' || 
        coalesce(NEW.description, '') || ' ' || 
        coalesce(NEW.short_description, '') || ' ' || 
        coalesce(NEW.sku, '') || ' ' || 
        coalesce(NEW.keywords, '') || ' ' ||
        coalesce(NEW.meta_title, '') || ' ' ||
        coalesce(NEW.meta_description, '')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update search vector
DROP TRIGGER IF EXISTS update_product_search_vector_trigger ON products;
CREATE TRIGGER update_product_search_vector_trigger
    BEFORE INSERT OR UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_product_search_vector();

-- Update existing products with search vector
UPDATE products SET search_vector = to_tsvector('english', 
    coalesce(name, '') || ' ' || 
    coalesce(description, '') || ' ' || 
    coalesce(short_description, '') || ' ' || 
    coalesce(sku, '') || ' ' || 
    coalesce(keywords, '') || ' ' ||
    coalesce(meta_title, '') || ' ' ||
    coalesce(meta_description, '')
);

-- Create GIN index for full-text search
CREATE INDEX IF NOT EXISTS idx_products_search_vector ON products USING GIN(search_vector);

-- Create additional search indexes
CREATE INDEX IF NOT EXISTS idx_products_name_gin ON products USING GIN(to_tsvector('english', name));
CREATE INDEX IF NOT EXISTS idx_products_description_gin ON products USING GIN(to_tsvector('english', description));

-- Create trigram indexes for fuzzy search
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS idx_products_name_trgm ON products USING GIN(name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_products_sku_trgm ON products USING GIN(sku gin_trgm_ops);

-- Create composite indexes for common search patterns
CREATE INDEX IF NOT EXISTS idx_products_search_composite ON products(status, category_id, price) WHERE status = 'active';
CREATE INDEX IF NOT EXISTS idx_products_featured_search ON products(featured, status, created_at) WHERE status = 'active';
CREATE INDEX IF NOT EXISTS idx_products_price_range ON products(price, status) WHERE status = 'active';

-- Create search synonyms table
CREATE TABLE IF NOT EXISTS search_synonyms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    term VARCHAR(255) NOT NULL,
    synonyms TEXT[] NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for synonyms
CREATE INDEX IF NOT EXISTS idx_search_synonyms_term ON search_synonyms(term);
CREATE INDEX IF NOT EXISTS idx_search_synonyms_active ON search_synonyms(is_active);

-- Create search analytics table
CREATE TABLE IF NOT EXISTS search_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    query TEXT NOT NULL,
    result_count INTEGER NOT NULL DEFAULT 0,
    click_through_rate DECIMAL(5,4) DEFAULT 0,
    conversion_rate DECIMAL(5,4) DEFAULT 0,
    avg_position_clicked DECIMAL(5,2) DEFAULT 0,
    search_date DATE NOT NULL DEFAULT CURRENT_DATE,
    total_searches INTEGER DEFAULT 1,
    total_clicks INTEGER DEFAULT 0,
    total_conversions INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(query, search_date)
);

-- Create indexes for search analytics
CREATE INDEX IF NOT EXISTS idx_search_analytics_query ON search_analytics(query);
CREATE INDEX IF NOT EXISTS idx_search_analytics_date ON search_analytics(search_date);
CREATE INDEX IF NOT EXISTS idx_search_analytics_ctr ON search_analytics(click_through_rate);

-- Create search suggestions table with better structure
CREATE TABLE IF NOT EXISTS search_suggestions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    query TEXT NOT NULL UNIQUE,
    frequency INTEGER DEFAULT 1,
    last_searched TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    result_count INTEGER DEFAULT 0,
    is_trending BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for search suggestions
CREATE INDEX IF NOT EXISTS idx_search_suggestions_query ON search_suggestions(query);
CREATE INDEX IF NOT EXISTS idx_search_suggestions_frequency ON search_suggestions(frequency DESC);
CREATE INDEX IF NOT EXISTS idx_search_suggestions_trending ON search_suggestions(is_trending, frequency DESC);

-- Insert some default synonyms
INSERT INTO search_synonyms (term, synonyms) VALUES
('phone', ARRAY['mobile', 'smartphone', 'cell phone', 'cellular']),
('laptop', ARRAY['notebook', 'computer', 'pc']),
('tv', ARRAY['television', 'monitor', 'screen']),
('shoes', ARRAY['footwear', 'sneakers', 'boots']),
('shirt', ARRAY['top', 'blouse', 'tee']),
('pants', ARRAY['trousers', 'jeans', 'slacks'])
ON CONFLICT DO NOTHING;

-- Create function to update search suggestions
CREATE OR REPLACE FUNCTION update_search_suggestion(search_query TEXT, result_cnt INTEGER DEFAULT 0)
RETURNS VOID AS $$
BEGIN
    INSERT INTO search_suggestions (query, frequency, result_count, last_searched)
    VALUES (search_query, 1, result_cnt, CURRENT_TIMESTAMP)
    ON CONFLICT (query) DO UPDATE SET
        frequency = search_suggestions.frequency + 1,
        last_searched = CURRENT_TIMESTAMP,
        result_count = EXCLUDED.result_count,
        updated_at = CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql;

-- Create function to get search suggestions with synonyms
CREATE OR REPLACE FUNCTION get_search_suggestions_with_synonyms(search_query TEXT, suggestion_limit INTEGER DEFAULT 10)
RETURNS TABLE(suggestion TEXT, frequency INTEGER, result_count INTEGER) AS $$
BEGIN
    RETURN QUERY
    SELECT
        s.query as suggestion,
        s.frequency,
        s.result_count
    FROM search_suggestions s
    WHERE s.query ILIKE '%' || search_query || '%'
       OR s.query % search_query  -- trigram similarity
    ORDER BY
        s.frequency DESC,
        similarity(s.query, search_query) DESC
    LIMIT suggestion_limit;
END;
$$ LANGUAGE plpgsql;
