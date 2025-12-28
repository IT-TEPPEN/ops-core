# ADR 0020: User Authentication with Multi-Provider OpenID Connect (OIDC)

## Status

Accepted (2025-12-27)
Updated (2025-12-28): Extended to support multiple OIDC providers

## Context

The application currently lacks user authentication, which is required for:
- Associating OAuth tokens (GitHub/GitLab) with specific users
- Tracking who created/executed documents and procedures
- Implementing access control and audit logs
- Securing API endpoints

We need a simple, reliable authentication mechanism that:
- Requires minimal setup for users (no password management)
- Provides strong security guarantees
- Integrates easily with existing OAuth infrastructure
- Follows industry standards

## Decision

We will implement **Multi-Provider OpenID Connect (OIDC)** for user authentication with **JWT** for session management.

### Core Principles

1. **OpScore Internal User Management**: Each user has a unique UUID managed by OpScore, independent of external providers
2. **Multiple Identity Support**: Users can link multiple OIDC providers (Google, GitHub, GitLab, Microsoft) to a single OpScore account
3. **Provider Abstraction**: Common OIDC logic abstracted to support any compliant provider
4. **Graceful Migration**: Existing users with one provider can add additional providers

### Authentication Flow

```
┌─────────┐                ┌──────────┐              ┌──────────┐
│ Frontend│                │ Backend  │              │ Provider │
└────┬────┘                └────┬─────┘              └────┬─────┘
     │                          │                         │
     │ 1. Click "Sign in with X"│                         │
     │─────────────────────────>│                         │
     │                          │                         │
     │ 2. Redirect to Provider  │                         │
     │<─────────────────────────│                         │
     │                          │                         │
     │ 3. User authenticates    │                         │
     │─────────────────────────────────────────────────────>│
     │                          │                         │
     │ 4. Redirect with code    │                         │
     │<─────────────────────────────────────────────────────│
     │                          │                         │
     │ 5. Send code to backend  │                         │
     │─────────────────────────>│                         │
     │                          │ 6. Exchange code for    │
     │                          │    ID Token             │
     │                          │────────────────────────>│
     │                          │<────────────────────────│
     │                          │ 7. Verify ID Token      │
     │                          │    signature            │
     │                          │                         │
     │                          │ 8. Find user_identity   │
     │                          │    by provider+sub      │
     │                          │                         │
     │                          │ 9. If not found, create │
     │                          │    user + user_identity │
     │                          │                         │
     │ 10. Return JWT (OpScore) │                         │
     │<─────────────────────────│                         │
     │                          │                         │
     │ 11. API calls with JWT   │                         │
     │─────────────────────────>│                         │
```

### Technology Stack

#### Backend (Go)
- **github.com/coreos/go-oidc/v3**: Generic OIDC client with automatic provider discovery
- **github.com/golang-jwt/jwt/v5**: JWT generation and validation
- **golang.org/x/oauth2**: OAuth 2.0 client library

#### Frontend (React)
- JWT stored in `localStorage`
- Automatic token refresh before expiration
- AuthContext for global authentication state

### OIDC Implementation Details

#### 1. Supported Providers

| Provider      | Issuer URL                                      | Notes                             |
| ------------- | ----------------------------------------------- | --------------------------------- |
| **Google**    | `https://accounts.google.com`                   | Most widely used                  |
| **GitHub**    | `https://token.actions.githubusercontent.com`   | Requires OAuth App setup          |
| **GitLab**    | `https://gitlab.com`                            | Self-hosted GitLab also supported |
| **Microsoft** | `https://login.microsoftonline.com/common/v2.0` | Azure AD / Microsoft Account      |

#### 2. Provider Discovery

All providers support OIDC Discovery via `/.well-known/openid-configuration`:

```go
provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
// Automatically discovers:
// - Authorization endpoint
// - Token endpoint
// - JWKS URI
// - Supported scopes and claims
```

#### 3. OIDC Scopes

```
openid email profile
```

- `openid`: Required for OIDC
- `email`: User's email address
- `profile`: User's name and profile picture

#### 4. ID Token Structure

All OIDC providers return a JWT ID Token with standard claims:

```json
{
  "iss": "https://provider.com",
  "sub": "provider-user-id",
  "email": "user@example.com",
  "email_verified": true,
  "name": "User Name",
  "picture": "https://...",
  "aud": "your-client-id",
  "exp": 1640000000,
  "iat": 1639999000
}
```

Provider-specific claims are stored but not required for authentication.

#### 5. ID Token Verification

Backend must verify:
1. **Signature**: Using provider's public keys (JWKS) via automatic discovery
2. **Issuer**: `iss` must match the configured provider issuer URL
3. **Audience**: `aud` must match our Client ID for that provider
4. **Expiration**: `exp` must be in the future
5. **Issued At**: `iat` must be in the past
6. **Email Verification**: `email_verified` should be true (provider-dependent)

### JWT Session Token

After successful OIDC authentication, backend issues its own JWT for session management.

#### JWT Claims

```json
{
  "sub": "user-uuid",
  "email": "user@example.com",
  "name": "User Name",
  "iat": 1640000000,
  "exp": 1640086400
}
```

#### JWT Configuration

- **Algorithm**: HS256 (HMAC-SHA256)
- **Secret**: Environment variable `JWT_SECRET` (min 32 bytes)
- **Expiration**: 24 hours
- **Issuer**: `opscore`

### Database Schema

#### Users Table (OpScore Internal)

Stores OpScore's internal user representation:

```sql
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
```

#### User Identities Table (OIDC Provider Links)

Stores links between OpScore users and OIDC provider identities:

```sql
CREATE TABLE user_identities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,  -- 'google', 'github', 'gitlab', 'microsoft'
    provider_user_id VARCHAR(255) NOT NULL,  -- Provider's sub claim
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
```

#### Authentication Logic

1. **First Login**: Creates both `users` and `user_identities` records
2. **Subsequent Logins**: Finds `user_identities` by `(provider, provider_user_id)`, returns associated `user`
3. **Account Linking**: Adds new `user_identities` record for existing `user`
4. **Email Matching**: Optional feature to auto-link accounts with same verified email

### API Endpoints

#### Authentication Endpoints

##### `GET /api/v1/auth/{provider}/login`
Initiates OIDC flow by redirecting to the provider's authorization endpoint.

**Path Parameters:**
- `provider`: One of `google`, `github`, `gitlab`, `microsoft`

**Query Parameters:**
- `redirect_uri` (optional): Frontend URL to redirect after successful login

**Response:** 302 Redirect to Provider

##### `GET /api/v1/auth/{provider}/callback`
Handles provider's redirect with authorization code.

**Path Parameters:**
- `provider`: One of `google`, `github`, `gitlab`, `microsoft`

**Query Parameters:**
- `code`: Authorization code from provider
- `state`: CSRF protection token

**Response:**
```json
{
  "token": "jwt-token-here",
  "user": {
    "id": "opscore-uuid",
    "email": "user@example.com",
    "name": "User Name",
    "picture": "https://..."
  }
}
```

##### `GET /api/v1/auth/me`
Returns current authenticated user information.

**Headers:**
- `Authorization: Bearer <jwt-token>`

**Response:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "User Name",
  "picture": "https://...",
  "created_at": "2025-12-27T10:00:00Z"
}
```

##### `POST /api/v1/auth/logout`
Logout endpoint (client-side token deletion primarily).

**Headers:**
- `Authorization: Bearer <jwt-token>`

**Response:**
```json
{
  "message": "Logged out successfully"
}
```

##### `GET /api/v1/auth/identities`
Returns all linked provider identities for the current user.

**Headers:**
- `Authorization: Bearer <jwt-token>`

**Response:**
```json
{
  "identities": [
    {
      "id": "uuid",
      "provider": "google",
      "email": "user@gmail.com",
      "name": "User Name",
      "linked_at": "2025-12-27T10:00:00Z",
      "last_used_at": "2025-12-28T09:00:00Z"
    },
    {
      "id": "uuid",
      "provider": "github",
      "email": "user@github.com",
      "name": "User Name",
      "linked_at": "2025-12-28T08:00:00Z",
      "last_used_at": "2025-12-28T08:00:00Z"
    }
  ]
}
```

##### `POST /api/v1/auth/{provider}/link`
Links a new provider identity to the current authenticated user.

**Path Parameters:**
- `provider`: One of `google`, `github`, `gitlab`, `microsoft`

**Headers:**
- `Authorization: Bearer <jwt-token>`

**Query Parameters:**
- `code`: Authorization code from provider
- `state`: CSRF protection token

**Response:**
```json
{
  "identity": {
    "id": "uuid",
    "provider": "gitlab",
    "email": "user@gitlab.com",
    "linked_at": "2025-12-28T10:00:00Z"
  }
}
```

##### `DELETE /api/v1/auth/identities/{id}`
Unlinks a provider identity from the current user.

**Headers:**
- `Authorization: Bearer <jwt-token>`

**Response:**
```json
{
  "message": "Identity unlinked successfully"
}
```

**Error Cases:**
- 400: Cannot unlink the last remaining identity
- 404: Identity not found or does not belong to user

### Authentication Middleware

Protected endpoints will use JWT middleware:

```go
// Verify JWT token
// Extract user ID from claims
// Set user ID in Gin context: c.Set("user_id", userID)
```

Protected routes:
- All `/api/v1/repositories` endpoints
- All `/api/v1/documents` endpoints
- All `/api/v1/execution-records` endpoints
- All other resource management endpoints

Public routes:
- `/api/v1/auth/*` (authentication endpoints)
- `/swagger/*` (API documentation)
- Health check endpoints

### Backend Architecture

#### Provider Abstraction Layer

```go
type OIDCProvider interface {
    GetAuthURL(state string, redirectURI string) string
    ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error)
    VerifyIDToken(ctx context.Context, rawIDToken string) (*IDTokenClaims, error)
}

type IDTokenClaims struct {
    Issuer          string
    Subject         string
    Email           string
    EmailVerified   bool
    Name            string
    Picture         string
}

// Provider Factory
func NewOIDCProvider(providerName string) (OIDCProvider, error) {
    switch providerName {
    case "google":
        return newGoogleProvider()
    case "github":
        return newGitHubProvider()
    case "gitlab":
        return newGitLabProvider()
    case "microsoft":
        return newMicrosoftProvider()
    default:
        return nil, fmt.Errorf("unsupported provider: %s", providerName)
    }
}
```

#### Environment Variables

Each provider requires its own OAuth credentials:

```bash
# Google
GOOGLE_CLIENT_ID=xxx
GOOGLE_CLIENT_SECRET=xxx

# GitHub
GITHUB_CLIENT_ID=xxx
GITHUB_CLIENT_SECRET=xxx

# GitLab
GITLAB_CLIENT_ID=xxx
GITLAB_CLIENT_SECRET=xxx

# Microsoft
MICROSOFT_CLIENT_ID=xxx
MICROSOFT_CLIENT_SECRET=xxx

# JWT
JWT_SECRET=xxx
```

### Frontend Implementation

#### AuthContext

```typescript
interface Identity {
  id: string;
  provider: string;
  email: string;
  name: string;
  linkedAt: string;
  lastUsedAt: string;
}

interface AuthContext {
  user: User | null;
  identities: Identity[];
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (provider: string) => void;
  logout: () => void;
  linkProvider: (provider: string) => void;
  unlinkIdentity: (identityId: string) => Promise<void>;
}
```

#### Login Page

```typescript
// Multiple provider buttons
<button onClick={() => login('google')}>Sign in with Google</button>
<button onClick={() => login('github')}>Sign in with GitHub</button>
<button onClick={() => login('gitlab')}>Sign in with GitLab</button>
<button onClick={() => login('microsoft')}>Sign in with Microsoft</button>
```

#### Protected Routes

```typescript
<ProtectedRoute>
  <RepositoryListPage />
</ProtectedRoute>
```

Redirects to `/login` if not authenticated.

#### Account Settings Page

Users can manage their linked identities:

```typescript
<IdentityList>
  {identities.map(identity => (
    <IdentityCard
      key={identity.id}
      provider={identity.provider}
      email={identity.email}
      onUnlink={() => unlinkIdentity(identity.id)}
      canUnlink={identities.length > 1}
    />
  ))}
</IdentityList>

<button onClick={() => linkProvider('gitlab')}>
  + Add GitLab Account
</button>
```

### Security Considerations

#### JWT Secret Management
- **Production**: Use strong random secret (min 32 bytes)
- Store in environment variable `JWT_SECRET`
- Never commit to version control

#### Token Storage
- Frontend stores JWT in `localStorage`
- XSS mitigation: Content Security Policy (CSP)
- HTTPS required in production

#### CSRF Protection
- Use `state` parameter in OIDC flow
- Generate cryptographically secure random state
- Validate state on callback

#### Token Expiration
- JWT expires after 24 hours
- Frontend checks expiration before API calls
- Automatic redirect to login on 401 responses

## Consequences

### Positive

✅ **No Password Management**: Users authenticate via trusted providers
✅ **Provider Choice**: Users can choose their preferred identity provider
✅ **Multiple Accounts**: Users can link multiple providers to one OpScore account
✅ **Industry Standard**: OIDC is a proven authentication protocol
✅ **Flexible UX**: Users can switch providers if one is unavailable
✅ **Secure**: ID Token verification prevents impersonation
✅ **Scalable**: JWT allows stateless authentication
✅ **Provider Independence**: OpScore manages its own user IDs, independent of providers
✅ **Migration Ready**: Easy to add new OIDC providers in the future

### Negative

❌ **Complexity**: More complex than single-provider authentication
❌ **Multiple Configs**: Each provider requires separate OAuth app setup
❌ **Privacy Concerns**: Providers know when users log in
❌ **Token Revocation**: JWTs cannot be revoked until expiration (mitigated by short expiration)
❌ **Identity Confusion**: Users may forget which provider they used

### Migration Path from Single Provider

1. **Phase 1**: Deploy new schema with `user_identities` table
2. **Phase 2**: Migrate existing `google_id` to `user_identities` records
3. **Phase 3**: Add support for additional providers (GitHub, GitLab, Microsoft)
4. **Phase 4**: Remove deprecated `google_id` column

### Account Merging Strategy

When a user authenticates with a new provider:

1. **Exact Email Match**: If another user has the same verified email, offer to link accounts
2. **User Confirmation**: Require user to confirm account linking via email verification
3. **Admin Override**: Admins can manually merge accounts if needed

### Related ADRs

- ADR 0004: Database Specification
- ADR 0005: Database Schema
- ADR 0008: Backend Logging Strategy

## References

- [OpenID Connect Core 1.0](https://openid.net/specs/openid-connect-core-1_0.html)
- [Google Identity Platform - OpenID Connect](https://developers.google.com/identity/protocols/oauth2/openid-connect)
- [RFC 7519: JSON Web Token (JWT)](https://tools.ietf.org/html/rfc7519)
