# Backend OAuth2.0 Setup Guide

## Overview

This backend supports OAuth2.0 authentication for:
- GitHub
- GitLab (gitlab.com)
- GitLab Self-Hosted (custom instances)

## Environment Variables

Create a `.env` file in the backend root directory:

```bash
# GitHub OAuth
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret

# GitLab OAuth  
GITLAB_CLIENT_ID=your_gitlab_application_id
GITLAB_CLIENT_SECRET=your_gitlab_secret

# OAuth Redirect URI
OAUTH_REDIRECT_URI=http://localhost:5173/oauth/callback

# Database
DATABASE_URL=postgres://opscore_user:opscore_password@db:5432/opscore_db?sslmode=disable

# Encryption Key (32 characters)
ENCRYPTION_KEY=your-32-character-encryption-key-here

# Logging
LOG_LEVEL=INFO
```

## OAuth Flow

### 1. Frontend initiates OAuth
Frontend sends user to OAuth provider (GitHub/GitLab) with:
- `client_id`
- `redirect_uri`
- `scope`
- `state` (CSRF protection)

### 2. User authenticates
User logs in and authorizes the application on the OAuth provider's site.

### 3. OAuth provider redirects back
Provider redirects to: `http://localhost:5173/oauth/callback?code=...&state=...`

### 4. Frontend sends code to backend
Frontend calls: `POST /api/v1/auth/oauth/callback`

Request body:
```json
{
  "provider": "github",
  "code": "authorization_code",
  "state": "random_state_string"
}
```

For self-hosted GitLab:
```json
{
  "provider": "gitlab-self-hosted",
  "code": "authorization_code",
  "state": "random_state_string",
  "gitlabUrl": "https://gitlab.example.com",
  "clientSecret": "client_secret_for_self_hosted"
}
```

### 5. Backend exchanges code for token
Backend:
1. Validates the request
2. Exchanges `code` for `access_token` with OAuth provider
3. Returns `access_token` to frontend

Response:
```json
{
  "access_token": "gho_xxxxxxxxxxxx",
  "token_type": "bearer",
  "scope": "repo,read:user",
  "expires_in": 28800
}
```

## API Endpoint

### POST /api/v1/auth/oauth/callback

**Request Body:**
```json
{
  "provider": "github" | "gitlab" | "gitlab-self-hosted",
  "code": "string",
  "state": "string",
  "gitlabUrl": "string (optional, required for gitlab-self-hosted)",
  "clientSecret": "string (optional, required for gitlab-self-hosted)"
}
```

**Success Response (200 OK):**
```json
{
  "access_token": "string",
  "token_type": "bearer",
  "scope": "string",
  "refresh_token": "string (optional)",
  "expires_in": number
}
```

**Error Responses:**

400 Bad Request - Invalid request format
```json
{
  "error": "invalid_request",
  "message": "Invalid request format"
}
```

400 Bad Request - Unsupported provider
```json
{
  "error": "invalid_provider",
  "message": "Unsupported Git provider"
}
```

500 Internal Server Error - Token exchange failed
```json
{
  "error": "token_exchange_failed",
  "message": "Error details..."
}
```

## Security Considerations

1. **Client Secret Storage**
   - Store in environment variables
   - Never commit to version control
   - Use secret management services in production

2. **State Parameter**
   - Frontend generates random state
   - Backend doesn't validate state (frontend responsibility)
   - Prevents CSRF attacks

3. **HTTPS**
   - Use HTTPS in production
   - Protects tokens in transit

4. **Token Storage**
   - Backend returns token to frontend
   - Frontend stores securely (localStorage with caution)
   - Consider using HttpOnly cookies for better security

## Testing

### Manual Test with curl

1. Get authorization code from OAuth provider manually
2. Test the callback endpoint:

```bash
curl -X POST http://localhost:8080/api/v1/auth/oauth/callback \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "github",
    "code": "your_authorization_code",
    "state": "random_state"
  }'
```

## Troubleshooting

### "OAuth credentials not configured"
- Check environment variables are set correctly
- Verify variable names match exactly

### "Token exchange failed"
- Verify Client ID and Secret are correct
- Check OAuth app settings on provider
- Ensure redirect URI matches exactly

### Self-hosted GitLab not working
- Verify GitLab URL is accessible from backend
- Check GitLab version supports OAuth2.0
- Ensure Client ID and Secret are from the correct GitLab instance
