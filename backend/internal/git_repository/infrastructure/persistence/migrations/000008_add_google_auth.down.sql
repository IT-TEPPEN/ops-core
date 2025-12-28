-- Remove Google authentication columns from users table
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_google_id;

ALTER TABLE users
DROP COLUMN IF EXISTS last_login_at,
DROP COLUMN IF EXISTS picture_url,
DROP COLUMN IF EXISTS google_id;
