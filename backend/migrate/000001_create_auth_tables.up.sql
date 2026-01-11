-- Filepath: backend/internal/git_repository/infrastructure/persistence/migrations/000001_create_auth_tables.up.sql
-- Create authentication and authorization related tables

-- Users table (OpScore internal user management - authentication only)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ
);

-- Refresh tokens table (JWT refresh token management)
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

-- User identities table (provider-based authentication links)
CREATE TABLE user_identities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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

-- User profiles table (user customizable profile information)
CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    display_name VARCHAR(255),
    picture_url TEXT,
    bio TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Primary user identity table (default identity selection)
CREATE TABLE primary_user_identity (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    identity_id UUID NOT NULL REFERENCES user_identities(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_primary_user_identity_identity_id ON primary_user_identity(identity_id);

-- Comments for documentation
COMMENT ON TABLE users IS 'OpScore internal user accounts (authentication only)';
COMMENT ON TABLE refresh_tokens IS 'JWT refresh tokens for session management';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'SHA-256 hash of the refresh token';
COMMENT ON COLUMN refresh_tokens.jti IS 'JWT ID (unique identifier for the refresh token)';
COMMENT ON COLUMN refresh_tokens.remember_me IS 'Whether this token has extended expiration (30 days vs 7 days)';
COMMENT ON TABLE user_identities IS 'External provider authentication links (GitHub, GitLab, etc.)';
COMMENT ON TABLE user_profiles IS 'User customizable profile information (overrides identity data)';
COMMENT ON COLUMN user_profiles.display_name IS 'Custom display name (NULL = use primary identity name)';
COMMENT ON COLUMN user_profiles.picture_url IS 'Custom avatar URL (NULL = use primary identity picture)';
COMMENT ON TABLE primary_user_identity IS 'Default identity selection for each user';
COMMENT ON COLUMN primary_user_identity.identity_id IS 'The identity to use as default for display and notification';
COMMENT ON TABLE user_identities IS 'Links between OpScore users and external authentication providers (Google, etc.)';
COMMENT ON TABLE groups IS 'User groups for access control';
COMMENT ON TABLE oauth_connections IS 'OAuth tokens for accessing Git repositories (GitHub, GitLab)';
COMMENT ON COLUMN oauth_connections.provider IS 'Git provider: github, gitlab, or gitlab-self-hosted';
COMMENT ON COLUMN oauth_connections.provider_host IS 'Provider host: github.com, gitlab.com, 192.168.0.1, etc.';
COMMENT ON COLUMN oauth_connections.provider_metadata IS 'Provider-specific metadata (JSONB): provider_user_id, provider_username, scopes, refresh_token_encrypted, token_expires_at, client_id_encrypted, client_secret_encrypted';
COMMENT ON COLUMN oauth_connections.access_token_encrypted IS 'AES-256-GCM encrypted access token';
