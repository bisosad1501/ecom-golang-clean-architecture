-- Remove admin reply fields from reviews table
DROP INDEX IF EXISTS idx_reviews_admin_reply_at;

ALTER TABLE reviews 
DROP COLUMN IF EXISTS admin_reply,
DROP COLUMN IF EXISTS admin_reply_at;
