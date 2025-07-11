-- Drop triggers
DROP TRIGGER IF EXISTS update_user_preferences_updated_at ON user_preferences;
DROP TRIGGER IF EXISTS update_user_verifications_updated_at ON user_verifications;

-- Drop indexes
DROP INDEX IF EXISTS idx_user_preferences_user_id;
DROP INDEX IF EXISTS idx_user_verifications_user_id;
DROP INDEX IF EXISTS idx_user_verifications_token;
DROP INDEX IF EXISTS idx_user_verifications_code;
DROP INDEX IF EXISTS idx_user_verifications_type;
DROP INDEX IF EXISTS idx_user_verifications_expires_at;

-- Drop tables
DROP TABLE IF EXISTS user_verifications;
DROP TABLE IF EXISTS user_preferences;
