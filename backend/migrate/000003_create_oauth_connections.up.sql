-- Filepath: backend/migrate/000003_create_oauth_connections.up.sql
-- Create oauth_connections table for managing user OAuth tokens

CREATE TABLE oauth_connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    git_provider_id UUID REFERENCES git_providers(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_host VARCHAR(255) NOT NULL DEFAULT '',
    provider_metadata JSONB NOT NULL,
    access_token_encrypted TEXT NOT NULL,
    refresh_token_encrypted TEXT,
    token_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ,
    UNIQUE(user_id, provider, provider_host)
);

CREATE INDEX idx_oauth_connections_user_id ON oauth_connections(user_id);
CREATE INDEX idx_oauth_connections_git_provider_id ON oauth_connections(git_provider_id);
CREATE INDEX idx_oauth_connections_provider ON oauth_connections(provider);
CREATE INDEX idx_oauth_connections_provider_host ON oauth_connections(provider_host);
CREATE INDEX idx_oauth_connections_last_used_at ON oauth_connections(last_used_at);

-- Comments
COMMENT ON TABLE oauth_connections IS 'User OAuth connections to Git providers for repository access';
COMMENT ON COLUMN oauth_connections.git_provider_id IS 'Reference to git_provider configuration (NULL for SaaS providers)';
COMMENT ON COLUMN oauth_connections.provider IS 'Provider type (github, gitlab, etc.)';
COMMENT ON COLUMN oauth_connections.provider_host IS 'Provider host (empty for SaaS, custom for self-hosted)';
COMMENT ON COLUMN oauth_connections.access_token_encrypted IS 'Encrypted OAuth access token';
COMMENT ON COLUMN oauth_connections.refresh_token_encrypted IS 'Encrypted OAuth refresh token (if supported)';
