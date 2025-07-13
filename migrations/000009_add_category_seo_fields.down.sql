-- Remove indexes
DROP INDEX IF EXISTS idx_categories_meta_title;
DROP INDEX IF EXISTS idx_categories_canonical_url;

-- Remove SEO fields from categories table
ALTER TABLE categories 
DROP COLUMN IF EXISTS meta_title,
DROP COLUMN IF EXISTS meta_description,
DROP COLUMN IF EXISTS meta_keywords,
DROP COLUMN IF EXISTS canonical_url,
DROP COLUMN IF EXISTS og_title,
DROP COLUMN IF EXISTS og_description,
DROP COLUMN IF EXISTS og_image,
DROP COLUMN IF EXISTS twitter_title,
DROP COLUMN IF EXISTS twitter_description,
DROP COLUMN IF EXISTS twitter_image,
DROP COLUMN IF EXISTS schema_markup;
