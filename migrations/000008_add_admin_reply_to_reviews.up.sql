-- Add admin reply fields to reviews table
ALTER TABLE reviews 
ADD COLUMN admin_reply TEXT,
ADD COLUMN admin_reply_at TIMESTAMP;

-- Add index for admin_reply_at for performance
CREATE INDEX idx_reviews_admin_reply_at ON reviews(admin_reply_at) WHERE admin_reply_at IS NOT NULL;
