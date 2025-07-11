-- Rollback migration: Restore username unique constraint
-- This will fail if there are duplicate usernames in the database

-- Remove the optional index
DROP INDEX IF EXISTS idx_users_username_optional;

-- Remove comment
COMMENT ON COLUMN users.username IS NULL;

-- Make username required again (this may fail if there are NULL values)
ALTER TABLE users ALTER COLUMN username SET NOT NULL;

-- Restore unique constraint (this may fail if there are duplicate values)
CREATE UNIQUE INDEX idx_users_username ON users(username);
