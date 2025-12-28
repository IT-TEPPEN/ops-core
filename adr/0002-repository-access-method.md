# ADR 0002: Access Method for External Git Repositories

## Status

Accepted (Updated: 2025-01)

## Context

Following ADR 0001, OpsCore requires read access to external private GitLab and GitHub repositories to fetch operational procedure Markdown documents. A secure, reliable, and manageable method is needed to authenticate and authorize OpsCore's access to these repositories.

## Decision

We use **OAuth 2.0** as the primary method for accessing external Git repositories. This provides a secure, user-friendly authentication flow with proper token lifecycle management.

### OAuth 2.0 Flow

1. **Authorization Flow**
   - User initiates GitHub/GitLab connection from OpsCore UI
   - OpsCore redirects to provider's OAuth authorization endpoint
   - User grants permissions to OpsCore application
   - Provider redirects back with authorization code
   - OpsCore exchanges code for access token and refresh token
   - Tokens are encrypted and stored in database, associated with user

2. **Token Management**
   - Access tokens are used for API calls to Git providers
   - Refresh tokens are used to obtain new access tokens when expired
   - Tokens are encrypted using AES-256-GCM before storage
   - Token refresh is handled automatically when API calls fail with 401

3. **Required Scopes**
   - **GitHub**: `repo` (for private repos), `read:user`
   - **GitLab**: `read_repository`, `read_user`

### Database Schema

```sql
CREATE TABLE oauth_connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,  -- 'github' or 'gitlab'
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

CREATE INDEX idx_oauth_connections_user_id ON oauth_connections(user_id);
CREATE INDEX idx_oauth_connections_provider ON oauth_connections(provider);
```

### Repository Access

When accessing repositories:

1. Look up user's OAuth connection for the repository's provider
2. Use the stored access token to make API calls
3. If token is expired, use refresh token to obtain new access token
4. If refresh fails, prompt user to re-authenticate

## Consequences

### Pros

- **User-friendly**: Standard OAuth flow familiar to users
- **Secure**: No manual token handling by users
- **Automatic refresh**: Token lifecycle managed by application
- **Audit trail**: All actions associated with authenticated user
- **Granular permissions**: Users explicitly grant permissions

### Cons

- **Complexity**: Requires OAuth flow implementation
- **Provider registration**: Need to register OpsCore as OAuth app with each provider
- **Token expiry**: Need to handle token refresh and re-authentication
