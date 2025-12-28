-- Rollback: Rename users table back
ALTER TABLE users RENAME TO users_new;

-- Recreate old users table structure
CREATE TABLE users (
    id UUID PRIMARY KEY,
    google_id VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    picture_url TEXT,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Migrate data back
INSERT INTO users (id, google_id, email, name, picture_url, last_login_at, created_at, updated_at)
SELECT u.id, ui.provider_user_id, u.primary_email, u.display_name, u.picture_url, u.last_login_at, u.created_at, u.updated_at
FROM users_new u
LEFT JOIN user_identities ui ON ui.user_id = u.id AND ui.provider = 'google';

-- Drop new tables
DROP TABLE user_identities;
DROP TABLE users_new CASCADE;

-- Recreate indexes
CREATE INDEX idx_users_google_id ON users(google_id);
CREATE INDEX idx_users_email ON users(email);
