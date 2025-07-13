-- Create category slug history table for tracking slug changes
CREATE TABLE category_slug_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    old_slug VARCHAR(255) NOT NULL,
    new_slug VARCHAR(255) NOT NULL,
    reason TEXT DEFAULT '',
    changed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT FALSE
);

-- Create indexes for efficient querying
CREATE INDEX idx_category_slug_history_category_id ON category_slug_history(category_id);
CREATE INDEX idx_category_slug_history_created_at ON category_slug_history(created_at);
CREATE INDEX idx_category_slug_history_is_active ON category_slug_history(is_active);

-- Create category redirects table for SEO-friendly redirects
CREATE TABLE category_redirects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_slug VARCHAR(255) NOT NULL UNIQUE,
    to_slug VARCHAR(255) NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    redirect_type VARCHAR(20) DEFAULT '301' CHECK (redirect_type IN ('301', '302', '307', '308')),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for redirects
CREATE INDEX idx_category_redirects_from_slug ON category_redirects(from_slug);
CREATE INDEX idx_category_redirects_to_slug ON category_redirects(to_slug);
CREATE INDEX idx_category_redirects_category_id ON category_redirects(category_id);
CREATE INDEX idx_category_redirects_is_active ON category_redirects(is_active);

-- Create SEO analytics tracking table
CREATE TABLE category_seo_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    seo_score INTEGER DEFAULT 0 CHECK (seo_score >= 0 AND seo_score <= 100),
    meta_title_length INTEGER DEFAULT 0,
    meta_description_length INTEGER DEFAULT 0,
    has_meta_title BOOLEAN DEFAULT FALSE,
    has_meta_description BOOLEAN DEFAULT FALSE,
    has_meta_keywords BOOLEAN DEFAULT FALSE,
    has_canonical_url BOOLEAN DEFAULT FALSE,
    has_open_graph BOOLEAN DEFAULT FALSE,
    has_twitter_cards BOOLEAN DEFAULT FALSE,
    has_schema_markup BOOLEAN DEFAULT FALSE,
    issues_count INTEGER DEFAULT 0,
    suggestions_count INTEGER DEFAULT 0,
    last_validated_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for SEO analytics
CREATE INDEX idx_category_seo_analytics_category_id ON category_seo_analytics(category_id);
CREATE INDEX idx_category_seo_analytics_seo_score ON category_seo_analytics(seo_score);
CREATE INDEX idx_category_seo_analytics_created_at ON category_seo_analytics(created_at);
CREATE INDEX idx_category_seo_analytics_last_validated_at ON category_seo_analytics(last_validated_at);

-- Create function to automatically update SEO analytics when category is updated
CREATE OR REPLACE FUNCTION update_category_seo_analytics()
RETURNS TRIGGER AS $$
BEGIN
    -- Insert or update SEO analytics
    INSERT INTO category_seo_analytics (
        category_id,
        meta_title_length,
        meta_description_length,
        has_meta_title,
        has_meta_description,
        has_meta_keywords,
        has_canonical_url,
        has_open_graph,
        has_twitter_cards,
        has_schema_markup,
        updated_at
    ) VALUES (
        NEW.id,
        LENGTH(COALESCE(NEW.meta_title, '')),
        LENGTH(COALESCE(NEW.meta_description, '')),
        NEW.meta_title IS NOT NULL AND NEW.meta_title != '',
        NEW.meta_description IS NOT NULL AND NEW.meta_description != '',
        NEW.meta_keywords IS NOT NULL AND NEW.meta_keywords != '',
        NEW.canonical_url IS NOT NULL AND NEW.canonical_url != '',
        (NEW.og_title IS NOT NULL AND NEW.og_title != '') OR (NEW.og_description IS NOT NULL AND NEW.og_description != ''),
        (NEW.twitter_title IS NOT NULL AND NEW.twitter_title != '') OR (NEW.twitter_description IS NOT NULL AND NEW.twitter_description != ''),
        NEW.schema_markup IS NOT NULL AND NEW.schema_markup != '',
        CURRENT_TIMESTAMP
    )
    ON CONFLICT (category_id) DO UPDATE SET
        meta_title_length = LENGTH(COALESCE(NEW.meta_title, '')),
        meta_description_length = LENGTH(COALESCE(NEW.meta_description, '')),
        has_meta_title = NEW.meta_title IS NOT NULL AND NEW.meta_title != '',
        has_meta_description = NEW.meta_description IS NOT NULL AND NEW.meta_description != '',
        has_meta_keywords = NEW.meta_keywords IS NOT NULL AND NEW.meta_keywords != '',
        has_canonical_url = NEW.canonical_url IS NOT NULL AND NEW.canonical_url != '',
        has_open_graph = (NEW.og_title IS NOT NULL AND NEW.og_title != '') OR (NEW.og_description IS NOT NULL AND NEW.og_description != ''),
        has_twitter_cards = (NEW.twitter_title IS NOT NULL AND NEW.twitter_title != '') OR (NEW.twitter_description IS NOT NULL AND NEW.twitter_description != ''),
        has_schema_markup = NEW.schema_markup IS NOT NULL AND NEW.schema_markup != '',
        updated_at = CURRENT_TIMESTAMP;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update SEO analytics
CREATE TRIGGER trigger_update_category_seo_analytics
    AFTER INSERT OR UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION update_category_seo_analytics();

-- Create function to track slug changes
CREATE OR REPLACE FUNCTION track_category_slug_changes()
RETURNS TRIGGER AS $$
BEGIN
    -- Only track if slug actually changed
    IF OLD.slug IS DISTINCT FROM NEW.slug THEN
        INSERT INTO category_slug_history (
            category_id,
            old_slug,
            new_slug,
            reason,
            is_active
        ) VALUES (
            NEW.id,
            OLD.slug,
            NEW.slug,
            'Slug updated',
            TRUE
        );

        -- Mark previous entries as inactive
        UPDATE category_slug_history 
        SET is_active = FALSE 
        WHERE category_id = NEW.id AND old_slug != OLD.slug;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to track slug changes
CREATE TRIGGER trigger_track_category_slug_changes
    AFTER UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION track_category_slug_changes();

-- Add unique constraint to ensure one active slug history per category
CREATE UNIQUE INDEX idx_category_slug_history_active_unique 
ON category_slug_history(category_id) 
WHERE is_active = TRUE;

-- Initialize SEO analytics for existing categories
INSERT INTO category_seo_analytics (
    category_id,
    meta_title_length,
    meta_description_length,
    has_meta_title,
    has_meta_description,
    has_meta_keywords,
    has_canonical_url,
    has_open_graph,
    has_twitter_cards,
    has_schema_markup
)
SELECT 
    id,
    LENGTH(COALESCE(meta_title, '')),
    LENGTH(COALESCE(meta_description, '')),
    meta_title IS NOT NULL AND meta_title != '',
    meta_description IS NOT NULL AND meta_description != '',
    meta_keywords IS NOT NULL AND meta_keywords != '',
    canonical_url IS NOT NULL AND canonical_url != '',
    (og_title IS NOT NULL AND og_title != '') OR (og_description IS NOT NULL AND og_description != ''),
    (twitter_title IS NOT NULL AND twitter_title != '') OR (twitter_description IS NOT NULL AND twitter_description != ''),
    schema_markup IS NOT NULL AND schema_markup != ''
FROM categories
ON CONFLICT (category_id) DO NOTHING;

-- Initialize slug history for existing categories
INSERT INTO category_slug_history (
    category_id,
    old_slug,
    new_slug,
    reason,
    is_active
)
SELECT 
    id,
    slug,
    slug,
    'Initial creation',
    TRUE
FROM categories
ON CONFLICT DO NOTHING;
