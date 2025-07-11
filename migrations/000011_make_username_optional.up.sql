-- Remove unique constraint on username and make it nullable
-- Migration: Make username optional and non-unique for e-commerce platform

-- Drop the unique index on username
DROP INDEX IF EXISTS idx_users_username;

-- Make username nullable 
ALTER TABLE users ALTER COLUMN username DROP NOT NULL;

-- Create a regular (non-unique) index for performance
CREATE INDEX idx_users_username_optional ON users(username) WHERE username IS NOT NULL;

-- Add comment to clarify the design decision
COMMENT ON COLUMN users.username IS 'Optional display name for user, non-unique. Email is the primary identifier.';
