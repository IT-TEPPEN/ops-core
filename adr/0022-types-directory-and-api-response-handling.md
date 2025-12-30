# ADR 0022: types/ Directory and API Response Handling

## Status

Accepted

## Date

2025-12-30

## Context

TypeScript type definitions are used throughout the frontend application for different purposes:

- **API response shapes**: Raw data from backend APIs (snake_case, string dates)
- **View data**: Display-ready data for UI components (camelCase, Date objects)
- **Domain model data**: Internal data structures for business logic
- **Simple feature types**: Type definitions for features using simplified structure

**Challenge**: Where should each type be placed to maintain clear boundaries and avoid confusion?

Without clear guidelines, developers mix concerns:

- API types leak into UI components
- View data types duplicate across features
- Domain model data is confused with API responses

## Decision

We adopt **layer-based type placement** with clear rules for each type purpose.

### Type Placement Rules

| Type Purpose | Location | Layer | Example |
| --- | --- | --- | --- |
| **View-oriented data (display-ready)** | `application/dto/` | Application | `DocumentViewData.ts` |
| **Raw API responses** | `infrastructure/services/types.ts` or inline | Infrastructure | `ApiDocumentResponse` |
| **Domain model data** | Part of Domain Model (no separate types/) | Domain | `document.getId()` |
| **Simple feature types** | `types/` (feature root, flat structure) | Feature root | `notification.ts` |
| **Shared domain types** | `features/common/types/` | Shared | `Permission.ts` |
| **Generic types** | `shared/types/` | Shared | `ApiResponse<T>` |

### Pattern A: Application Layer DTOs (View-Oriented Data)

**Purpose**: Data Transfer Objects for data that Presentation layer consumes

**Location**: `application/dto/`

**Characteristics**:

- Display-ready format (camelCase, Date objects, formatted strings)
- Returned by Query/Command Services
- Used by Presentation components
- No business logic (pure data)

**Example**:

```typescript
// application/dto/DocumentViewData.ts
export interface DocumentViewData {
  id: string;
  title: string;
  content: string;
  status: 'draft' | 'published';
  createdAt: Date;        // ✅ Transformed from string
  updatedAt: Date;        // ✅ Transformed from string
  ownerName: string;      // ✅ Display-ready format
  tagLabels: string[];    // ✅ Processed for display
}

// application/dto/CreateDocumentInput.ts
export interface CreateDocumentInput {
  repositoryId: string;
  filePath: string;
  title: string;
  content: string;
  accessScope: 'public' | 'private';
}
```

**Usage**:

```typescript
// Presentation layer consumes directly
function DocumentCard({ document }: { document: DocumentViewData }) {
  return (
    <div>
      <h3>{document.title}</h3>
      {/* ✅ No transformation needed */}
      <p>Created: {document.createdAt.toLocaleDateString()}</p>
      <p>Owner: {document.ownerName}</p>
      <div>{document.tagLabels.join(', ')}</div>
    </div>
  );
}
```

### Pattern B: Infrastructure Layer Types (Raw API Responses)

**Purpose**: Internal type safety within Infrastructure implementations

**Location**: `infrastructure/services/types.ts` (or inline in service file)

**Characteristics**:

- Raw API format (snake_case, string dates)
- Internal to Infrastructure layer
- Not exported from feature
- Transformed to ViewData before returning

**Example**:

```typescript
// infrastructure/services/types.ts
interface ApiDocumentResponse {
  id: string;
  title: string;
  content: string;
  status: 'draft' | 'published';
  created_at: string;     // ✅ Raw from API
  updated_at: string;     // ✅ Raw from API
  owner_name: string;     // ✅ snake_case
  tags: string[];         // ✅ Raw tags
}

// infrastructure/services/HttpDocumentQueryService.ts
export class HttpDocumentQueryService implements DocumentQueryService {
  async getById(id: string): Promise<DocumentViewData> {
    const response = await this.apiClient.get<ApiDocumentResponse>(
      `/documents/${id}`
    );

    // ✅ Transform to ViewData
    return {
      id: response.id,
      title: response.title,
      content: response.content,
      status: response.status,
      createdAt: new Date(response.created_at),
      updatedAt: new Date(response.updated_at),
      ownerName: response.owner_name,
      tagLabels: response.tags.map(tag => `#${tag}`),
    };
  }
}
```

**Rules**:

- Keep API types internal to Infrastructure
- Never export API types from feature
- Always transform to ViewData before returning

### Pattern C: Domain Model Data

**Purpose**: Data managed by Domain entities with business logic

**Location**: Part of Domain Model (no separate `types/`)

**Characteristics**:

- Accessed via getter methods
- Encapsulated within Domain entities
- Contains business logic
- Immutable

**Example**:

```typescript
// domain/models/Document.ts
export class DocumentEntity {
  private constructor(
    private readonly id: string,
    private readonly title: string,
    private readonly content: string,
    private readonly status: 'draft' | 'published'
  ) {}

  // ✅ No separate interface, data is encapsulated
  getId(): string { return this.id; }
  getTitle(): string { return this.title; }
  getContent(): string { return this.content; }
  getStatus(): 'draft' | 'published' { return this.status; }

  // Business logic
  canPublish(): boolean {
    return this.status === 'draft' && this.content.length > 0;
  }
}
```

**Note**: Do NOT create separate types for Domain model data. Domain models ARE the types.

### Pattern D: Simple Feature types/ (Flat Structure)

**Purpose**: Type definitions for features using simplified structure (Pattern 1)

**Location**: `types/` (at feature root)

**Characteristics**:

- Used in simplified features (< 10 files)
- Flat structure, no layers
- May include both API types and view types

**Example**:

```typescript
// features/notification/types/notification.ts
export interface Notification {
  id: string;
  title: string;
  message: string;
  type: 'info' | 'success' | 'warning' | 'error';
  createdAt: Date;
}

export interface NotificationListResponse {
  notifications: Notification[];
  unreadCount: number;
}
```

**When migrating to layered structure**, these types move to:

- View data → `application/dto/`
- API types → `infrastructure/services/types.ts`
- Domain data → Domain models

### Decision Matrix

**Question flow**: "Where should I place this type?"

```text
1. Is the feature using layered structure (Pattern 2)?
   NO → Place in `types/` (feature root)
   YES ↓

2. Is this raw API response shape?
   YES → `infrastructure/services/types.ts` (internal only)
   NO ↓

3. Is this display-ready data for UI?
   YES → `application/dto/`
   NO ↓

4. Is this part of Domain model?
   YES → Encapsulate in Domain model (no separate type)
   NO ↓

5. Is this shared across features?
   YES → `features/common/types/`
   NO → `application/dto/` (default)
```

### API Response Transformation Pattern

**Standard pattern** for transforming API responses to ViewData:

```typescript
// 1. Define API response type (Infrastructure)
interface ApiUserResponse {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  created_at: string;
  role: string;
}

// 2. Define ViewData (Application)
export interface UserViewData {
  id: string;
  fullName: string;      // ✅ Computed from first + last
  email: string;
  createdAt: Date;       // ✅ Transformed
  roleLabel: string;     // ✅ Formatted for display
}

// 3. Transform in Service (Infrastructure)
export class HttpUserQueryService implements UserQueryService {
  private toViewData(response: ApiUserResponse): UserViewData {
    return {
      id: response.id,
      fullName: `${response.first_name} ${response.last_name}`,
      email: response.email,
      createdAt: new Date(response.created_at),
      roleLabel: this.formatRole(response.role),
    };
  }

  private formatRole(role: string): string {
    return role === 'admin' ? 'Administrator' : 'User';
  }

  async getById(id: string): Promise<UserViewData> {
    const response = await this.apiClient.get<ApiUserResponse>(`/users/${id}`);
    return this.toViewData(response);
  }
}
```

**Benefits**:

- **Consistency**: All components receive the same format
- **No duplicate transformation**: Done once in Service
- **Type safety**: ViewData prevents using raw API format
- **Testability**: Easy to test transformation logic

### Shared Types

**`shared/types/`** (Generic, reusable):

```typescript
// shared/types/ApiResponse.ts
export interface ApiResponse<T> {
  data: T;
  status: number;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  perPage: number;
  total: number;
}

export type Result<T, E = Error> =
  | { ok: true; value: T }
  | { ok: false; error: E };
```

**`features/common/types/`** (Domain-specific, shared):

```typescript
// features/common/types/Permission.ts
export type Permission =
  | 'document:read'
  | 'document:write'
  | 'document:delete'
  | 'user:manage';

export type UserRole = 'admin' | 'user' | 'guest';

// features/common/types/DocumentStatus.ts
export type DocumentStatus = 'draft' | 'published' | 'archived';

export interface StatusTransition {
  from: DocumentStatus;
  to: DocumentStatus;
  requiredPermission: Permission;
}
```

## Consequences

### Positive

1. **Clear boundaries**: Each type has a well-defined location
2. **No confusion**: ViewData vs API types are clearly separated
3. **Consistency**: All components receive the same format
4. **Transformation centralized**: Done once in Services
5. **Type safety**: TypeScript enforces correct usage
6. **Discoverability**: Easy to find relevant types

### Negative

1. **Duplication**: API types and ViewData may look similar
2. **Transformation overhead**: Every API response needs transformation
3. **Verbosity**: More types to maintain

### Mitigation

1. **Code generation**: Generate ViewData from OpenAPI specs
2. **Shared utilities**: Create reusable transformation helpers
3. **Documentation**: This ADR clarifies placement rules
4. **Code reviews**: Ensure transformations are correct

## Examples

### Example 1: Complete Flow (API → ViewData → Component)

```typescript
// 1. API type (Infrastructure, internal)
interface ApiDocumentResponse {
  id: string;
  title: string;
  created_at: string;
  owner: {
    id: string;
    name: string;
  };
}

// 2. ViewData (Application, public)
export interface DocumentViewData {
  id: string;
  title: string;
  createdAt: Date;
  ownerName: string;
}

// 3. Service transformation (Infrastructure)
export class HttpDocumentQueryService {
  async getById(id: string): Promise<DocumentViewData> {
    const response = await this.apiClient.get<ApiDocumentResponse>(`/documents/${id}`);

    return {
      id: response.id,
      title: response.title,
      createdAt: new Date(response.created_at),
      ownerName: response.owner.name,
    };
  }
}

// 4. Component usage (Presentation)
function DocumentCard({ document }: { document: DocumentViewData }) {
  return (
    <div>
      <h3>{document.title}</h3>
      <p>By {document.ownerName}</p>
      <time>{document.createdAt.toLocaleDateString()}</time>
    </div>
  );
}
```

### Example 2: Migration from Simple to Layered

**Before** (Simple structure):

```typescript
// features/user/types/user.ts
export interface User {
  id: string;
  name: string;
  email: string;
  created_at: string;  // ❌ Mixing API format with view
}
```

**After** (Layered structure):

```typescript
// infrastructure/services/types.ts (internal)
interface ApiUserResponse {
  id: string;
  name: string;
  email: string;
  created_at: string;
}

// application/dto/UserViewData.ts (public)
export interface UserViewData {
  id: string;
  name: string;
  email: string;
  createdAt: Date;  // ✅ Transformed
}

// infrastructure/services/HttpUserQueryService.ts
export class HttpUserQueryService {
  async getById(id: string): Promise<UserViewData> {
    const response = await this.apiClient.get<ApiUserResponse>(`/users/${id}`);
    return {
      id: response.id,
      name: response.name,
      email: response.email,
      createdAt: new Date(response.created_at),
    };
  }
}
```

### Example 3: Nested API Response Transformation

```typescript
// API type (Infrastructure)
interface ApiDocumentDetailResponse {
  id: string;
  title: string;
  content: string;
  author: {
    id: string;
    first_name: string;
    last_name: string;
    avatar_url: string;
  };
  tags: Array<{
    id: string;
    name: string;
    color: string;
  }>;
  created_at: string;
  updated_at: string;
}

// ViewData (Application)
export interface DocumentDetailViewData {
  id: string;
  title: string;
  content: string;
  author: {
    id: string;
    fullName: string;
    avatarUrl: string;
  };
  tags: Array<{
    id: string;
    label: string;
    color: string;
  }>;
  createdAt: Date;
  updatedAt: Date;
}

// Transformation (Infrastructure)
export class HttpDocumentQueryService {
  private toDetailViewData(response: ApiDocumentDetailResponse): DocumentDetailViewData {
    return {
      id: response.id,
      title: response.title,
      content: response.content,
      author: {
        id: response.author.id,
        fullName: `${response.author.first_name} ${response.author.last_name}`,
        avatarUrl: response.author.avatar_url,
      },
      tags: response.tags.map(tag => ({
        id: tag.id,
        label: `#${tag.name}`,
        color: tag.color,
      })),
      createdAt: new Date(response.created_at),
      updatedAt: new Date(response.updated_at),
    };
  }

  async getDetailById(id: string): Promise<DocumentDetailViewData> {
    const response = await this.apiClient.get<ApiDocumentDetailResponse>(
      `/documents/${id}/detail`
    );
    return this.toDetailViewData(response);
  }
}
```

## Related ADRs

- **ADR 0018**: Frontend Layered Architecture - Layer definitions
- **ADR 0019**: Data Access Patterns - Query/Command Services that use these types
- **ADR 0021**: Feature Structure Patterns - When to use flat vs layered

## References

- [TypeScript Handbook: Types](https://www.typescriptlang.org/docs/handbook/2/everyday-types.html)
- [DTO Pattern](https://martinfowler.com/eaaCatalog/dataTransferObject.html)
- [Anti-Corruption Layer](https://docs.microsoft.com/en-us/azure/architecture/patterns/anti-corruption-layer)
