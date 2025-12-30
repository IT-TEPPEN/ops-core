# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

OpsCore is an operations management system that aggregates operational procedure documents (Markdown files) from external repositories (GitHub/GitLab) and makes them viewable on the web. The system supports variable input, execution tracking, and evidence management.

**Tech Stack:**

- Backend: Go 1.21+ (Echo framework, Wire DI, zap logging, PostgreSQL)
- Frontend: React 19 + TypeScript + Vite
- Database: PostgreSQL 14+
- Architecture: Onion Architecture with Domain-Driven Design

## Essential Commands

### Root Level

```bash
# Start both backend and frontend in development mode
npm run dev

# Start backend only
npm run dev:backend

# Start frontend only
npm run dev:frontend

# Start development environment with Docker
docker compose up -d
```

### Backend (`backend/`)

```bash
# Run the server
go run cmd/server/main.go

# Build
go build -o ./bin/server ./cmd/server

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests for specific context
go test ./internal/document/...

# Run with race detection
go test -race ./...

# Database migrations
go run ./cmd/migrate up          # Apply all pending migrations
go run ./cmd/migrate status      # Check migration status
go run ./cmd/migrate down        # Rollback last migration

# Generate Swagger API documentation
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
# View at: http://localhost:8080/swagger/index.html

# API integration testing
go run ./cmd/api_tester/main.go
```

### Frontend (`frontend/`)

```bash
# Development server (http://localhost:5173)
npm run dev

# Production build
npm run build
npm run preview

# Linting
npm run lint

# Testing
npm test                  # Run once
npm run test:watch        # Watch mode
npm run test:coverage     # With coverage report

# Storybook
npm run storybook         # Start Storybook dev server
npm run build-storybook   # Build Storybook
```

## Architecture

### Onion Architecture with Bounded Contexts

The backend follows Onion Architecture principles with clear separation of concerns across layers. **Dependencies always point inward** (Interfaces → Application → Domain).

#### Directory Structure

```text
backend/internal/<context>/
├── domain/              # Business logic (no external dependencies)
│   ├── entity/         # Domain entities (unexported fields, factory functions)
│   ├── value_object/   # Immutable value objects
│   ├── domain_service/ # Multi-entity domain logic
│   ├── repository/     # Repository interfaces
│   └── error/          # Domain-specific errors
├── application/         # Use cases (depends on Domain)
│   ├── usecase/        # Use case implementations
│   ├── dto/            # Data Transfer Objects (no API-specific tags)
│   └── error/          # Application-specific errors
├── infrastructure/      # Technical implementations (depends on Domain/Application interfaces)
│   ├── persistence/    # Database implementations
│   ├── external/       # External API clients (e.g., git/)
│   └── error/          # Infrastructure-specific errors
└── interfaces/          # External interfaces (depends on Application)
    ├── api/
    │   ├── handlers/   # HTTP handlers
    │   └── schema/     # API schemas (with json/binding/swagger tags)
    └── error/          # HTTP error responses
```

#### Bounded Contexts

The system is divided into the following contexts:

1. **git_repository**: External Git repository management (GitHub/GitLab)
   - Repository registration, access token encryption (AES-256-GCM)
   - See: `backend/internal/git_repository/infrastructure/encryption/README.md`

2. **document**: Document lifecycle management
   - Document publishing, versioning (commit hash + version number)
   - Variable definitions in Frontmatter, access scope (public/private)
   - See: ADR 0013, 0016, 0017

3. **execution_record**: Execution tracking and evidence management
   - Execution sessions, step-by-step notes, screenshot attachments
   - Storage abstraction (local/S3/MinIO)
   - See: ADR 0014, `backend/infrastructure/storage/README.md`

4. **user**: User and group management
   - User CRUD, group management, role-based access (admin/user)

5. **view_history**: Document viewing history tracking

6. **view_statistics**: Aggregated view statistics

7. **auth**: Authentication and authorization

8. **oauth**: OAuth connection handling

9. **shared**: Common utilities and middleware

### Key Architectural Patterns

#### Entity Pattern

- Entities have **unexported fields** to enforce encapsulation
- Created via **factory functions** (`New<Entity>`) with validation
- Reconstructed via **reconstructor functions** (`Reconstruct<Entity>`) for persistence
- Access fields via **getter methods**
- Mutate state via **behavior methods**

```go
// Factory for new entities (with validation)
func NewDocument(id DocumentID, title string) (*Document, error) { ... }

// Reconstructor for entities from DB (minimal validation)
func ReconstructDocument(id DocumentID, title string) *Document { ... }
```

#### DTO vs Schema Separation

- **Application DTOs** (`application/dto/`): Use case data structures, no API-specific tags
- **API Schemas** (`interfaces/api/schema/`): API request/response with json/binding/swagger tags
- **Conversion**: Schema → DTO (handler) → Use case → DTO → Schema (handler)
- This ensures application layer is independent of API format

#### Custom Error Design

Each layer has its own error types with specific prefixes:

- `DOMAIN_XXX`: Domain layer errors
- `APP_XXX`: Application layer errors
- `INFRA_XXX`: Infrastructure layer errors
- `API_XXX`: API layer errors

See: ADR 0015 for error design details

#### Dependency Injection

Uses Wire for compile-time DI. All dependencies configured in `cmd/server/di.go`.

### Frontend Architecture (React + TypeScript)

The frontend follows a **feature-based architecture** with clear separation between application foundation, features, and UI components. **Dependencies always point inward** for complex features.

#### Project Structure

```text
frontend/src/
├── app/                    # Application configuration
│   ├── providers/          # Global providers (auth, theme, etc.)
│   ├── routes/             # Routing configuration
│   └── App.tsx            # Root component
├── pages/                  # Page components (route handlers)
│   └── <PageName>/
├── features/               # Feature modules (Bounded Contexts)
│   ├── <featureName>/     # Individual feature
│   └── common/            # Shared across features (domain-specific)
├── ui/                     # UI component library (design system)
│   └── <ComponentName>/
└── shared/                 # Application foundation (domain-agnostic)
    ├── api/               # API client configuration
    ├── constants/         # Environment variables, routes
    ├── types/             # Global type definitions
    ├── utils/             # Generic utilities
    └── lib/               # External library wrappers
```

See: ADR 0012 for project structure details

#### Feature Structure Patterns

Features use one of two patterns based on complexity:

**Pattern 1: Simplified Structure** (< 10 files)

```text
features/simpleFeature/
├── components/          # React components
├── hooks/              # Custom hooks
├── types/              # Type definitions
└── index.ts            # Public API
```

**Pattern 2: Layered Structure** (10+ files, complex business logic)

```text
features/complexFeature/
├── domain/                    # Business logic (no React, no HTTP)
│   ├── models/               # Domain entities (immutable)
│   └── repositories/         # Repository interfaces (local state)
├── application/              # Use cases (no React hooks)
│   ├── services/            # Query/Command Service interfaces
│   ├── usecases/            # Business workflows
│   └── dto/                 # Data Transfer Objects (ViewData)
├── infrastructure/           # External implementations
│   ├── services/            # HTTP Service implementations
│   └── repositories/        # LocalStorage implementations
└── presentation/             # React UI
    ├── components/          # UI components
    ├── hooks/               # UI state + Usecase invocation
    └── contexts/            # Dependency injection
```

**Migration Path**: Start with Pattern 1 (simplified), migrate to Pattern 2 (layered) when:

- Feature grows beyond 10 files
- Complex business logic emerges
- Local state management needed (LocalStorage, drafts)

**Three Strikes Rule**: Only extract to `features/common/` after encountering the same need three times.

See: ADR 0021 for feature structure patterns

#### Key Frontend Patterns

##### Data Access Patterns

Three distinct patterns for data operations:

1. **Query Service** (External API reads)
   - Interface in `application/services/`
   - Returns **ViewData** (display-ready format: camelCase, Date objects)
   - Implemented as HTTP client in `infrastructure/services/`
   - No business logic, just fetch and transform
   - Works with React Query/SWR for caching

2. **Command Service** (External API writes)
   - Interface in `application/services/`
   - Executes POST/PUT/DELETE operations
   - Called by Usecases after business validation
   - Backend is source of truth for validation

3. **Repository** (Local state management)
   - Interface in `domain/repositories/`
   - Manages **Domain Models** (not ViewData)
   - Contains business logic and validation
   - Implemented with LocalStorage, SessionStorage, or useReducer
   - Used for: drafts, multi-step forms, offline support, undo/redo

**Decision flow**:

```text
Is data managed by backend?
  YES → Query Service (GET) or Command Service (POST/PUT/DELETE)
  NO → Does it need complex frontend business logic?
    YES → Repository with Domain Models
    NO → Query Service (simple display)
```

See: ADR 0019 for data access patterns

##### Type Placement

- **ViewData** (`application/dto/`): Display-ready data for Presentation layer (camelCase, Date objects)
- **API Response Types** (`infrastructure/services/types.ts`): Raw API format (snake_case, string dates, internal only)
- **Domain Model Data**: Part of Domain entity (no separate types, accessed via getters)
- **Simple Feature Types** (`types/` at feature root): For Pattern 1 features
- **Shared Domain Types** (`features/common/types/`): Domain-specific, shared (Permission, UserRole)
- **Generic Types** (`shared/types/`): Domain-agnostic (`ApiResponse<T>`, Result<T, E>)

**API Response Transformation**:

```typescript
// Infrastructure: Raw API type (internal)
interface ApiDocumentResponse {
  id: string;
  created_at: string;  // Raw from API
  owner_name: string;  // snake_case
}

// Application: ViewData (public)
export interface DocumentViewData {
  id: string;
  createdAt: Date;      // Transformed
  ownerName: string;    // camelCase
}

// Infrastructure: Transform in Service
class HttpDocumentQueryService {
  async getById(id: string): Promise<DocumentViewData> {
    const response = await this.apiClient.get<ApiDocumentResponse>(`/documents/${id}`);
    return {
      id: response.id,
      createdAt: new Date(response.created_at),
      ownerName: response.owner_name,
    };
  }
}
```

See: ADR 0022 for type placement rules

##### Immutable Domain Models

All Domain models are strictly immutable (React Strict Mode requirement):

```typescript
export class DocumentEntity {
  private constructor(
    private readonly id: string,
    private readonly title: string,
    private readonly status: 'draft' | 'published'
  ) {}

  // Factory for new entities
  static create(data: { title: string }): DocumentEntity {
    if (data.title.length === 0) {
      throw new DomainError('EMPTY_TITLE', 'Title cannot be empty');
    }
    return new DocumentEntity(crypto.randomUUID(), data.title, 'draft');
  }

  // Reconstructor for persistence
  static reconstruct(data: { id: string; title: string; status: 'draft' | 'published' }): DocumentEntity {
    return new DocumentEntity(data.id, data.title, data.status);
  }

  // Getters
  getId(): string { return this.id; }
  getTitle(): string { return this.title; }

  // Update methods return NEW instances
  updateTitle(newTitle: string): DocumentEntity {
    return new DocumentEntity(this.id, newTitle, this.status);
  }

  publish(): DocumentEntity {
    return new DocumentEntity(this.id, this.title, 'published');
  }
}
```

See: ADR 0023 for React-specific patterns

##### Context-Based Dependency Injection

Use React Context for injecting repositories and services:

```typescript
// 1. Interface in Domain/Application layer
export interface DocumentRepository {
  get(id: string): Document | null;
  getAll(): Document[];
}

// 2. Implementation in Infrastructure layer
export class LocalStorageDocumentRepository implements DocumentRepository { ... }

// 3. Context Provider in Presentation layer
const DocumentRepositoryContext = createContext<DocumentRepository | null>(null);

export function DocumentRepositoryProvider({ children }: { children: ReactNode }) {
  const repository = useMemo(() => new LocalStorageDocumentRepository(), []);
  return (
    <DocumentRepositoryContext.Provider value={repository}>
      {children}
    </DocumentRepositoryContext.Provider>
  );
}

// 4. Consumption hook (separate file for Fast Refresh compatibility)
export function useDocumentRepository(): DocumentRepository {
  const context = useContext(DocumentRepositoryContext);
  if (!context) throw new Error('useDocumentRepository must be used within provider');
  return context;
}
```

##### Feature Public API

Each feature exports its public API through `index.ts`:

```typescript
// features/document/index.ts

// Export public components
export { DocumentTable } from './presentation/components/DocumentTable';
export { DocumentForm } from './presentation/components/DocumentForm';

// Export public hooks
export { useDocumentList } from './presentation/hooks/useDocumentList';

// Export public types (ViewData only, not internal types)
export type { DocumentViewData } from './application/dto/DocumentViewData';

// Export providers for app setup
export { DocumentRepositoryProvider } from './presentation/contexts/DocumentRepositoryContext';

// Do NOT export: Domain models, Usecases, Infrastructure implementations
```

**Usage**:

```typescript
// Import from feature root (preferred)
import { DocumentTable, useDocumentList } from '@/features/document';

// Avoid: Direct imports from internal structure
```

See: ADR 0023 for React-specific implementation patterns

#### `shared/` vs `features/common/` Decision Criteria

**Key Question**: "Does this code contain business logic specific to this application?"

| Aspect             | `shared/`                        | `features/common/`                    |
| ------------------ | -------------------------------- | ------------------------------------- |
| **Nature**         | Generic, reusable                | App-specific, domain-dependent        |
| **Business Logic** | None                             | Yes                                   |
| **Utils Example**  | `formatDate`, `debounce`         | `formatUserName`, `calculateDiscount` |
| **Components**     | - (place in `ui/`)               | `DataTable`, `StatusBadge`            |
| **Hooks**          | `useLocalStorage`, `useDebounce` | `useAuth`, `usePermission`            |
| **Types**          | `ApiResponse<T>`, `Result<T, E>` | `Permission`, `UserRole`              |
| **Portability**    | Usable in other projects         | Application-specific                  |

## Critical System Details

### Security

#### User Authentication

Multi-Provider OpenID Connect (OIDC) for user authentication:

**Supported Providers**: Google, GitHub, GitLab, Microsoft

**Database Tables**:

- `users`: OpScore internal user (id UUID, primary_email, display_name, picture_url)
- `user_identities`: Links to OIDC providers (user_id, provider, provider_user_id, email)

**Flow**:

1. User authenticates with chosen provider (Google/GitHub/GitLab/Microsoft)
2. Backend verifies ID Token from provider
3. Creates/updates `user` and `user_identities` records
4. Issues JWT session token (24-hour expiration)

**Environment Variables** (required for each provider):

```bash
GOOGLE_CLIENT_ID / GOOGLE_CLIENT_SECRET
GITHUB_CLIENT_ID / GITHUB_CLIENT_SECRET
GITLAB_CLIENT_ID / GITLAB_CLIENT_SECRET
MICROSOFT_CLIENT_ID / MICROSOFT_CLIENT_SECRET
JWT_SECRET=<32-byte-secret>  # For session token signing
```

**Features**:

- Users can link multiple providers to one account
- Account linking via verified email matching
- OpScore manages internal user UUIDs (provider-independent)

See: ADR 0050 for authentication details

#### Git Repository Integration

**OAuth 2.0 Access**: Uses OAuth for GitHub/GitLab repository access

**Database Table**:

- `oauth_connections`: Links users to Git providers (user_id, provider, access_token_encrypted, refresh_token_encrypted)

**Token Management**:

- Access tokens encrypted with **AES-256-GCM** (same as below)
- Automatic token refresh when expired
- Tokens associated with authenticated user (not global)

**Markdown Document Structure**: Documents must have YAML frontmatter:

```yaml
---
title: "Procedure Title"
owner: "Team Name"
version: "1.0"
type: "procedure"  # or "knowledge"
tags: ["tag1", "tag2"]
variables:  # Optional: for parameterized procedures
  - name: server_name
    label: "サーバー名"
    type: string
    required: true
    defaultValue: "prod-db-01"
---
```

See: ADR 0001 (Markdown structure), ADR 0002 (OAuth access)

#### Access Token Encryption

- All Git provider tokens encrypted with **AES-256-GCM**
- Encryption key must be set via `ENCRYPTION_KEY` environment variable (32 bytes)
- Development: `export ENCRYPTION_KEY="dev-key-12345678901234567890123"`
- Production: Generate secure 32-byte key
- See: `backend/internal/git_repository/infrastructure/encryption/README.md`

#### Data Validation

- Application layer validates all input data
- Domain entities enforce business rules via factory functions
- See: ADR 0017

### Database

#### Schema

Key tables:

**Authentication**:

- `users`: OpScore internal users (id UUID, primary_email, display_name)
- `user_identities`: OIDC provider links (user_id, provider, provider_user_id, email)

**Git Integration**:

- `repositories`: Git repository information
- `oauth_connections`: OAuth tokens for Git providers (user_id, provider, access_token_encrypted)

**Documents**:

- `documents`, `document_versions`: Document and version management

**Execution**:

- `execution_records`, `execution_steps`: Execution tracking
- `attachments`: Screenshot/evidence files

**Access Control**:

- `groups`, `user_groups`: Group management

**Analytics**:

- `view_histories`, `view_statistics`: View tracking

See: ADR 0004, 0005, 0006 for database specification and migration strategy

#### Migrations

- Located in `backend/migrations/`
- Use `golang-migrate` library
- Migration files follow format: `<version>_<description>.up.sql` and `.down.sql`
- See: ADR 0006

### Storage Abstraction

Attachment file storage supports multiple backends:

- **Local filesystem**: Development (`STORAGE_TYPE=local`)
- **AWS S3**: Production (`STORAGE_TYPE=s3`)
- **MinIO**: On-premise (`STORAGE_TYPE=minio`)

Configuration via environment variables. See: `backend/infrastructure/storage/README.md`

## Development Guidelines

### Follow the ADRs

All architectural decisions are documented in `adr/` directory.

**Key Backend ADRs**:

- **0007**: Backend Onion Architecture
- **0009**: Backend Testing Strategy
- **0013**: Document Variable Definition
- **0014**: Execution Record and Evidence Management
- **0015**: Backend Custom Error Design
- **0016**: Document Domain Model Design
- **0017**: Application Data Validation

**Key Frontend ADRs**:

- **0012**: Frontend Project Structure
- **0018**: Frontend Layered Architecture
- **0019**: Data Access Patterns (Query/Command/Repository)
- **0020**: Frontend Validation Strategy
- **0021**: Feature Structure Patterns
- **0022**: Types Directory and API Response Handling
- **0023**: React-Specific Implementation Patterns

**Authentication & Integration ADRs**:

- **0001**: External Repository Markdown Structure
- **0002**: Repository Access Method (OAuth 2.0)
- **0050**: User Authentication (Multi-Provider OIDC)

### Minimize Dependencies

Use minimal external libraries. Justify any new dependency additions.

### Testing Requirements

#### Test Coverage Targets

- **Overall**: 80%+
- **Domain layer**: 90%+ (business logic core)
- **Application layer**: 85%+
- **Infrastructure layer**: 70%+
- **Interface layer**: 75%+

#### Test Structure

- Backend: `<file>_test.go` colocated with source files
- Frontend: `<file>.test.ts(x)` colocated with source files
- Use table-driven tests for Go
- Use Vitest + React Testing Library for frontend

See: `docs/development/TESTING.md` for detailed testing guide

### Code Conventions

#### Backend

- Package names match directory names
- Entity factory functions: `NewEntity()` with validation
- Reconstructor functions: `ReconstructEntity()` for persistence
- Repository interfaces in `domain/repository/`
- Repository implementations in `infrastructure/persistence/`
- Use context.Context for cancellation/timeouts
- Error handling: wrap errors with context, use custom error types

#### Frontend

- **Components**: Functional components with TypeScript, use React 19 features
- **State Management**:
  - Server state: React Query (external APIs)
  - Local state: useState/useReducer for UI, Repository pattern for complex local state
- **Form Handling**: React Hook Form + Zod for validation
- **Styling**: Tailwind CSS with component-based design
- **Markdown Processing**: Gray-matter (frontmatter), Unified/Remark/Rehype pipeline
- **File Organization**:
  - Start with simplified structure (Pattern 1) for new features
  - Migrate to layered structure (Pattern 2) when feature exceeds 10 files or has complex logic
  - Export public API through feature `index.ts`
- **Naming Conventions**:
  - Components: PascalCase (`DocumentTable.tsx`)
  - Hooks: camelCase with `use` prefix (`useDocumentList.ts`)
  - Types: PascalCase for interfaces/types (`DocumentViewData`)
  - Files: Match primary export name
- **Imports**:
  - Use feature public API: `import { X } from '@/features/document'`
  - Avoid deep imports: `import { X } from '@/features/document/presentation/components/X'`
- **Context/Hook Separation**: Separate Context providers from consumption hooks for Fast Refresh compatibility
- **Immutability**: All domain models must be immutable (return new instances on updates)

### API Documentation

- All API endpoints documented with Swagger annotations
- Generate with `swag init` command
- See: ADR 0010 for API definition generation spec

## Environment Variables

### Backend Required

```bash
# Core
DATABASE_URL=<postgres-url>       # PostgreSQL connection string
PORT=8080                         # Server port (default: 8080)

# Security
ENCRYPTION_KEY=<32-byte-key>      # Token encryption key (REQUIRED, 32 bytes)
JWT_SECRET=<32-byte-secret>       # JWT signing secret (REQUIRED, 32 bytes)

# Storage
STORAGE_TYPE=local|s3|minio       # Attachment storage backend
```

### Authentication (OIDC Providers)

At least one provider must be configured:

```bash
# Google
GOOGLE_CLIENT_ID=<client-id>
GOOGLE_CLIENT_SECRET=<secret>

# GitHub
GITHUB_CLIENT_ID=<client-id>
GITHUB_CLIENT_SECRET=<secret>

# GitLab
GITLAB_CLIENT_ID=<client-id>
GITLAB_CLIENT_SECRET=<secret>

# Microsoft
MICROSOFT_CLIENT_ID=<client-id>
MICROSOFT_CLIENT_SECRET=<secret>
```

### Storage-Specific (when using S3/MinIO)

```bash
AWS_REGION=<region>               # AWS region for S3
AWS_ACCESS_KEY_ID=<key>           # S3/MinIO access key
AWS_SECRET_ACCESS_KEY=<secret>    # S3/MinIO secret key
S3_BUCKET=<bucket-name>           # S3/MinIO bucket name
S3_ENDPOINT=<endpoint>            # MinIO endpoint (MinIO only)
```

### Frontend Environment Variables

```bash
# Development server
VITE_API_BASE_URL=http://localhost:8080  # Backend API URL

# Feature flags (optional)
VITE_ENABLE_STORYBOOK=true
```

## Additional Documentation

- **User Guide**: `docs/user-guide/README.md`
- **Development Guide**: `docs/development/CONTRIBUTING.md`, `API.md`, `TESTING.md`
- **Deployment**: `docs/deployment/README.md`
- **Operations**: `docs/operations/MONITORING.md`, `BACKUP.md`
- **ADRs**: `adr/` directory for all architectural decisions
