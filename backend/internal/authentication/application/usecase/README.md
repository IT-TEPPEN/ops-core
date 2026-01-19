# Authentication Usecase Implementation Plan

## Overview

This directory contains the usecase interfaces for the authentication bounded context.
Usecases are designed to be technology-agnostic and focus on business logic.

## Implementation Priority

### Phase 1: Core Authentication (最優先 - システムの基盤)

These usecases are essential for the system to function at all.

1. **AuthenticateUser** - User authentication/registration
   - Authenticates user with external provider information
   - Creates new user if identity not found
   - Updates last login time for existing users
   - **Status**: Not implemented

2. **CreateSession** - Session creation
   - Creates authentication session for authenticated user
   - Supports short-term (7 days) and long-term (30 days) sessions
   - **Status**: Not implemented

3. **GetAuthenticatedUser** - Retrieve authenticated user
   - Validates session and returns current user information
   - Used for every authenticated API call
   - **Status**: Not implemented

### Phase 2: Session Management (高優先 - 実用的なシステムに必要)

These usecases are required for practical, production-ready system.

4. **RefreshSession** - Session refresh
   - Extends session validity without re-authentication
   - Essential for long-running user sessions
   - **Status**: Not implemented

5. **RevokeSession** - Session revocation
   - Invalidates a specific session (logout)
   - Important for security
   - **Status**: Not implemented

---

## Implementation Note

We will implement **usecases 1-5** first, as they form the minimum viable authentication system.

After Phase 1 and 2 are complete, we can proceed to:

- Phase 3: Multi-provider identity management (LinkNewIdentity, UnlinkIdentity, SetPrimaryIdentity, ListUserIdentities)
- Phase 4: Profile and security features (UpdateUserProfile, GetUserSessions, RevokeAllUserSessions)

## Architecture

All usecases follow DDD principles:

- **Technology-independent**: No OIDC, JWT, or other implementation details in usecase interfaces
- **Business logic focused**: Defines "what" to do, not "how" to do it
- **Interface-based**: Usecases are interfaces, implementations go in separate files
- **DTO separation**: Data structures are defined in `../dto/` package

## Next Steps

For each usecase implementation, we need to create:

1. Domain layer: Entities, Value Objects, Repository interfaces
2. Infrastructure layer: Repository implementations, external service adapters
3. Application layer: Usecase implementations
4. Interface layer: HTTP handlers, request/response schemas
