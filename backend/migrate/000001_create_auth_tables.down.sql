-- Drop authentication tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS refresh_tokens CASCADE;
DROP TABLE IF EXISTS user_identities CASCADE;
DROP TABLE IF EXISTS users CASCADE;

