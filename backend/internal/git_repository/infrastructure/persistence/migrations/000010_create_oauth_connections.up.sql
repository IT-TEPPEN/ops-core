-- OAuth connections table for storing user's Git provider tokens
CREATE TABLE IF NOT EXISTS oauth_connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    provider_username VARCHAR(255),
    access_token_encrypted TEXT NOT NULL,
    refresh_token_encrypted TEXT,
    token_expires_at TIMESTAMP WITH TIME ZONE,
    scopes TEXT[],
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, provider)
);

CREATE INDEX IF NOT EXISTS idx_oauth_connections_user_id ON oauth_connections(user_id);
CREATE INDEX IF NOT EXISTS idx_oauth_connections_provider ON oauth_connections(provider);

-- Add comment for documentation
COMMENT ON TABLE oauth_connections IS 'Stores OAuth tokens for Git provider connections';
COMMENT ON COLUMN oauth_connections.provider IS 'Git provider: github or gitlab';
COMMENT ON COLUMN oauth_connections.access_token_encrypted IS 'AES-256-GCM encrypted access token';
COMMENT ON COLUMN oauth_connections.refresh_token_encrypted IS 'AES-256-GCM encrypted refresh token';
