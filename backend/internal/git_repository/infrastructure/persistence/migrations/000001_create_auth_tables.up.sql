-- Filepath: backend/internal/git_repository/infrastructure/persistence/migrations/000001_create_auth_tables.up.sql
-- Create authentication and authorization related tables

-- Users table (OpScore internal user management)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    primary_email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    picture_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ
);

CREATE INDEX idx_users_primary_email ON users(primary_email);

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

-- Groups table
CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_groups_name ON groups(name);

-- User-Groups junction table (many-to-many)
CREATE TABLE user_groups (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, group_id)
);

CREATE INDEX idx_user_groups_user_id ON user_groups(user_id);
CREATE INDEX idx_user_groups_group_id ON user_groups(group_id);

-- OAuth connections table (Git provider tokens for repository access)
CREATE TABLE oauth_connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_host VARCHAR(255) NOT NULL DEFAULT '',
    provider_metadata JSONB NOT NULL,
    access_token_encrypted TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, provider, provider_host)
);

CREATE INDEX idx_oauth_connections_user_id ON oauth_connections(user_id);
CREATE INDEX idx_oauth_connections_provider ON oauth_connections(provider);
CREATE INDEX idx_oauth_connections_provider_host ON oauth_connections(provider_host);

-- Comments for documentation
COMMENT ON TABLE users IS 'OpScore internal user accounts';
COMMENT ON TABLE user_identities IS 'Links between OpScore users and external authentication providers (Google, etc.)';
COMMENT ON TABLE groups IS 'User groups for access control';
COMMENT ON TABLE oauth_connections IS 'OAuth tokens for accessing Git repositories (GitHub, GitLab)';
COMMENT ON COLUMN oauth_connections.provider IS 'Git provider: github, gitlab, or gitlab-self-hosted';
COMMENT ON COLUMN oauth_connections.provider_host IS 'Provider host: github.com, gitlab.com, 192.168.0.1, etc.';
COMMENT ON COLUMN oauth_connections.provider_metadata IS 'Provider-specific metadata (JSONB): provider_user_id, provider_username, scopes, refresh_token_encrypted, token_expires_at, client_id_encrypted, client_secret_encrypted';
COMMENT ON COLUMN oauth_connections.access_token_encrypted IS 'AES-256-GCM encrypted access token';
