-- Create new users table with OpScore internal user management
CREATE TABLE IF NOT EXISTS users_new (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    primary_email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    picture_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ
);

CREATE INDEX idx_users_new_primary_email ON users_new(primary_email);

-- Create user_identities table for provider links
CREATE TABLE IF NOT EXISTS user_identities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    picture_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ,
    CONSTRAINT uq_provider_identity UNIQUE(provider, provider_user_id)
);

CREATE INDEX idx_user_identities_user_id ON user_identities(user_id);
CREATE INDEX idx_user_identities_provider ON user_identities(provider);
CREATE INDEX idx_user_identities_email ON user_identities(email);

-- Migrate existing data from old users table
INSERT INTO users_new (id, primary_email, display_name, picture_url, created_at, updated_at, last_login_at)
SELECT id, email, name, picture_url, created_at, updated_at, last_login_at
FROM users
WHERE email IS NOT NULL AND name IS NOT NULL
ON CONFLICT DO NOTHING;

-- Migrate existing Google authentication info to user_identities
INSERT INTO user_identities (user_id, provider, provider_user_id, email, name, picture_url, created_at, updated_at, last_used_at)
SELECT id, 'google', google_id, email, name, picture_url, created_at, updated_at, last_login_at
FROM users
WHERE google_id IS NOT NULL
ON CONFLICT DO NOTHING;

-- Drop old users table and rename new one
DROP TABLE IF EXISTS users CASCADE;
ALTER TABLE users_new RENAME TO users;

-- Add foreign key constraint after table rename
ALTER TABLE user_identities ADD CONSTRAINT fk_user_identities_user 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
