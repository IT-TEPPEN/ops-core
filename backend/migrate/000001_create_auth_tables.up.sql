-- Create authentication related tables
-- Aligned with authentication domain model (DDD)

-- Users table (Aggregate Root: User)
-- Represents an authenticated user in OpScore
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    picture_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

-- User identities table (Child Entity: Identity)
-- Links between OpScore users and external authentication providers
CREATE TABLE user_identities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    picture_url TEXT,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_provider_identity UNIQUE(provider, provider_user_id)
);

CREATE INDEX idx_user_identities_user_id ON user_identities(user_id);
CREATE INDEX idx_user_identities_provider ON user_identities(provider);
CREATE INDEX idx_user_identities_email ON user_identities(email);
CREATE INDEX idx_user_identities_is_primary ON user_identities(user_id, is_primary) WHERE is_primary = TRUE;

-- Refresh tokens table (Entity: Session)
-- JWT refresh tokens for session management
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    jti VARCHAR(36) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    remember_me BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_jti ON refresh_tokens(jti);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX idx_refresh_tokens_is_revoked ON refresh_tokens(is_revoked);

-- Comments for documentation
COMMENT ON TABLE users IS 'OpScore user accounts (Aggregate Root)';
COMMENT ON COLUMN users.email IS 'User email address (from primary identity or first identity)';
COMMENT ON COLUMN users.display_name IS 'User display name (from primary identity or first identity)';
COMMENT ON COLUMN users.picture_url IS 'User profile picture URL (from primary identity or first identity)';

COMMENT ON TABLE user_identities IS 'External provider authentication links (Child Entity of User)';
COMMENT ON COLUMN user_identities.provider IS 'Provider name: google, github, gitlab, microsoft';
COMMENT ON COLUMN user_identities.provider_user_id IS 'User ID from the provider (sub claim from OIDC)';
COMMENT ON COLUMN user_identities.is_primary IS 'Whether this is the primary identity (used for display info)';

COMMENT ON TABLE refresh_tokens IS 'JWT refresh tokens for session management (Entity: Session)';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'SHA-256 hash of the refresh token';
COMMENT ON COLUMN refresh_tokens.jti IS 'JWT ID (unique identifier for the JWT access token)';
COMMENT ON COLUMN refresh_tokens.remember_me IS 'Whether this token has extended expiration (30 days vs 7 days)';

