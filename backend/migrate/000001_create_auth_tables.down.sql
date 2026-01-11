-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS primary_user_identity CASCADE;
DROP TABLE IF EXISTS user_profiles CASCADE;
DROP TABLE IF EXISTS user_identities CASCADE;
DROP TABLE IF EXISTS refresh_tokens CASCADE;
DROP TABLE IF EXISTS users CASCADE;
