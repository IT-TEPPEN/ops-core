# ADR 0012: Frontend Project Structure

## Status

Accepted

## Date

2025-12-28 (Updated: 2025-12-30)

## Context

The frontend application (React + TypeScript + Vite) needs a clear project-level folder structure that:

- Separates concerns at the application level
- Supports feature-based development (Bounded Contexts)
- Clarifies the difference between generic utilities and domain-specific code
- Provides a foundation for scalable architecture

## Decision

### Project-Level Structure

We adopt a feature-based organization with clear separation between application foundation, features, and UI components:

```text
src/
├── app/                    # Application configuration
│   ├── providers/          # Global providers (auth, theme, etc.)
│   ├── routes/             # Routing configuration
│   └── App.tsx            # Root component
│
├── pages/                  # Page components (route handlers)
│   └── <PageName>/
│
├── features/               # Feature modules (Bounded Contexts)
│   ├── <featureName>/     # Individual feature
│   └── common/            # Shared across features (domain-specific)
│
├── ui/                     # UI component library (design system)
│   └── <ComponentName>/
│
└── shared/                 # Application foundation (domain-agnostic)
    ├── api/               # API client, fetch configuration
    ├── constants/         # Environment variables, routes
    ├── types/             # Global type definitions
    ├── utils/             # Generic utilities
    └── lib/               # External library wrappers
```

### Responsibility of Each Directory

#### `app/` - Application Configuration

**Purpose**: Application-wide configuration and setup

**Contains**:

- **`providers/`**: Global React Context Providers (AuthProvider, ThemeProvider, NotificationProvider, etc.)
- **`routes/`**: Routing configuration (React Router setup)
- **`App.tsx`**: Root application component

**Examples**:

- `app/providers/AuthProvider.tsx`
- `app/routes/index.tsx`
- `app/App.tsx`

#### `pages/` - Page Components

**Purpose**: Page-level components that correspond to routes

**Contains**:

- One directory per page/route
- Composes features into cohesive pages
- Handles route-level concerns (params, query strings)

**Rules**:

- Pages should be thin - they compose features, not implement logic
- Business logic belongs in features, not pages

**Examples**:

- `pages/HomePage/`
- `pages/DocumentListPage/`
- `pages/UserSettingsPage/`

#### `features/` - Feature Modules (Bounded Contexts)

**Purpose**: Self-contained feature modules representing Bounded Contexts

**Contains**:

- **`<featureName>/`**: Individual features (document, user, repository, etc.)
- **`common/`**: Code shared across multiple features (domain-specific)

**Characteristics**:

- Each feature is independently developable
- Features communicate through public APIs (exported via `index.ts`)
- Features should minimize dependencies on other features

**Note**: For internal structure of features, see **ADR 0021: Feature Structure Patterns**.

**Examples**:

- `features/document/`
- `features/repository/`
- `features/common/hooks/useAuth.ts`

#### `ui/` - UI Component Library (Design System)

**Purpose**: Reusable UI components (design system)

**Contains**:

- Generic UI components (Button, Input, Card, Typography, etc.)
- Design system primitives
- No business logic

**Rules**:

- Components should be pure presentational
- Reusable across any project
- No domain-specific logic

**Examples**:

- `ui/Button/`
- `ui/Input/`
- `ui/Card/`

#### `shared/` - Application Foundation (Domain-Agnostic)

**Purpose**: Generic, reusable utilities and configurations

**Contains**:

- **`api/`**: Base API client, fetch configuration
- **`constants/`**: Environment variables, route definitions
- **`types/`**: Global type definitions (e.g., `ApiResponse<T>`, `Result<T, E>`)
- **`utils/`**: Generic utilities (formatDate, debounce, validation)
- **`lib/`**: External library wrappers

**Rules**:

- No business logic specific to this application
- Should be portable to other projects
- "Could be extracted into an npm package"

**Examples**:

- `shared/utils/formatDate.ts`
- `shared/api/client.ts`
- `shared/types/ApiResponse.ts`

#### `features/common/` - Shared Domain Code

**Purpose**: Code shared across features that contains domain-specific logic

**Contains**:

- Shared components with business logic (DataTable, StatusBadge)
- Shared hooks (useAuth, usePagination, usePermission)
- Shared types (Permission, UserRole)
- Shared utilities with business logic

**Rules**:

- Contains application-specific business logic
- Requires knowledge of this application's domain
- Not portable to other projects

**Examples**:

- `features/common/components/DataTable.tsx`
- `features/common/hooks/useAuth.ts`
- `features/common/types/Permission.ts`

### `shared/` vs `features/common/` Decision Criteria

**Key Question**: "Does this code contain business logic specific to this application?"

| Aspect | `shared/` | `features/common/` |
| ------ | --------- | ------------------ |
| **Nature** | Generic, reusable | App-specific, domain-dependent |
| **Business Logic** | None | Yes |
| **Utils Example** | `formatDate`, `debounce`, `clamp` | `formatUserName`, `calculateDiscount` |
| **Components** | - (place in `ui/`) | `DataTable`, `StatusBadge` |
| **Hooks** | `useLocalStorage`, `useDebounce` | `useAuth`, `usePermission` |
| **Types** | `ApiResponse<T>`, `Result<T, E>` | `Permission`, `UserRole` |
| **Portability** | Usable in other projects | Application-specific |
| **Dependencies** | No domain knowledge required | Requires domain knowledge |

**Decision Rules**:

1. **"Could this be used in another project without modification?"**
   - YES → `shared/`
   - NO → `features/common/`

2. **"Does this require knowledge of this application's business rules?"**
   - YES → `features/common/`
   - NO → `shared/`

3. **"Is this a pure UI component?"**
   - YES → `ui/`
   - NO → Check rules 1 & 2

**Examples**:

```typescript
// ✅ shared/utils/formatDate.ts (generic)
export function formatDate(date: Date): string {
  return date.toLocaleDateString();
}

// ✅ features/common/utils/formatUserName.ts (domain-specific)
export function formatUserName(user: User): string {
  // Business logic: admin users have special formatting
  return user.role === 'admin' ? `[ADMIN] ${user.name}` : user.name;
}
```

### Feature-Based Organization Principles

**Bounded Context**: Each feature represents a Bounded Context from Domain-Driven Design:

- **Autonomy**: Features should be independently developable
- **Encapsulation**: Internal implementation details are hidden
- **Explicit APIs**: Features export public APIs via `index.ts`
- **Minimal coupling**: Features minimize dependencies on each other

**When to create a new feature**:

- The functionality represents a distinct business domain
- The feature can be understood and developed independently
- The feature has clear boundaries and responsibilities

**When NOT to create a new feature**:

- The functionality is too small (< 5 files) - keep it simple first
- The functionality is tightly coupled to an existing feature - extend that feature instead
- The functionality is purely UI - consider `ui/` or `features/common/components/`

## Consequences

### Positive

1. **Clear boundaries**: Each directory has a well-defined purpose
2. **Scalability**: Easy to add new features without affecting existing ones
3. **Reusability**: `shared/` and `ui/` can be extracted to packages
4. **Team collaboration**: Features are independently developable
5. **Discoverability**: Clear structure makes it easy to find code
6. **Testability**: Features can be tested in isolation

### Negative

1. **Initial setup**: Requires understanding the structure
2. **Judgment calls**: Deciding between `shared/` and `features/common/` requires thought
3. **Potential duplication**: Strict separation may lead to some duplication

### Mitigation

1. **Documentation**: This ADR provides clear decision criteria
2. **Code reviews**: Validate placement decisions during review
3. **Start specific, extract later**: When in doubt, start in feature-specific location and extract to `features/common/` or `shared/` when the need becomes clear
4. **Three Strikes Rule**: Only extract to `features/common/` after seeing the same code three times (see ADR 0021)

## Examples

### Example 1: Adding a new utility

**Scenario**: You need a function to format currency values.

**Decision process**:

1. Is it generic? → YES (currency formatting is universal)
2. Does it need domain knowledge? → NO
3. → Place in `shared/utils/formatCurrency.ts`

### Example 2: Adding authentication hook

**Scenario**: You need a hook to check user permissions.

**Decision process**:

1. Is it generic? → NO (permissions are app-specific)
2. Does it need domain knowledge? → YES (Permission types, role logic)
3. Is it shared across features? → YES
4. → Place in `features/common/hooks/usePermission.ts`

### Example 3: Adding a status badge component

**Scenario**: You need a badge to show document status (draft/published).

**Decision process**:

1. Is it pure UI? → NO (status logic is domain-specific)
2. Is it generic? → NO (status values are app-specific)
3. Is it shared across features? → YES
4. → Place in `features/common/components/StatusBadge.tsx`

## Related ADRs

- **ADR 0021**: Feature Structure Patterns - Defines internal structure of features (simplified vs layered)
- **ADR 0018**: Frontend Layered Architecture - Defines layered structure within complex features

## References

- [Feature-Sliced Design](https://feature-sliced.design/) - Feature-based architecture methodology
- [Bulletproof React](https://github.com/alan2207/bulletproof-react) - React application architecture
- [Domain-Driven Design](https://en.wikipedia.org/wiki/Domain-driven_design) - Bounded Context concept
