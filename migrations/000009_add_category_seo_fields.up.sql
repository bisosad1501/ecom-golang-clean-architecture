-- Add SEO fields to categories table
ALTER TABLE categories 
ADD COLUMN meta_title VARCHAR(255) DEFAULT '',
ADD COLUMN meta_description TEXT DEFAULT '',
ADD COLUMN meta_keywords TEXT DEFAULT '',
ADD COLUMN canonical_url VARCHAR(500) DEFAULT '',
ADD COLUMN og_title VARCHAR(255) DEFAULT '',
ADD COLUMN og_description TEXT DEFAULT '',
ADD COLUMN og_image VARCHAR(500) DEFAULT '',
ADD COLUMN twitter_title VARCHAR(255) DEFAULT '',
ADD COLUMN twitter_description TEXT DEFAULT '',
ADD COLUMN twitter_image VARCHAR(500) DEFAULT '',
ADD COLUMN schema_markup TEXT DEFAULT '';

-- Create indexes for SEO fields that might be searched
CREATE INDEX idx_categories_meta_title ON categories(meta_title);
CREATE INDEX idx_categories_canonical_url ON categories(canonical_url);

-- Update existing categories with basic SEO data
UPDATE categories 
SET 
    meta_title = CASE 
        WHEN LENGTH(name) <= 60 THEN name || ' - Shop Online'
        ELSE name
    END,
    meta_description = CASE 
        WHEN description IS NOT NULL AND description != '' THEN 
            CASE 
                WHEN LENGTH(description) <= 160 THEN description
                ELSE LEFT(description, 157) || '...'
            END
        ELSE 'Shop ' || name || ' products online. Find the best deals and latest products in ' || name || ' category.'
    END,
    meta_keywords = name || ', shop ' || name || ', buy ' || name || ' online',
    canonical_url = '/categories/' || slug,
    og_title = CASE 
        WHEN LENGTH(name) <= 60 THEN name || ' - Shop Online'
        ELSE name
    END,
    og_description = CASE 
        WHEN description IS NOT NULL AND description != '' THEN 
            CASE 
                WHEN LENGTH(description) <= 160 THEN description
                ELSE LEFT(description, 157) || '...'
            END
        ELSE 'Shop ' || name || ' products online. Find the best deals and latest products in ' || name || ' category.'
    END,
    twitter_title = CASE 
        WHEN LENGTH(name) <= 60 THEN name || ' - Shop Online'
        ELSE name
    END,
    twitter_description = CASE 
        WHEN description IS NOT NULL AND description != '' THEN 
            CASE 
                WHEN LENGTH(description) <= 160 THEN description
                ELSE LEFT(description, 157) || '...'
            END
        ELSE 'Shop ' || name || ' products online. Find the best deals and latest products in ' || name || ' category.'
    END
WHERE meta_title = '' OR meta_title IS NULL;
