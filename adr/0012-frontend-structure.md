# ADR 001: Frontend Folder Structure

## Status

Accepted

## Date

2025-12-28

## Context

We are building the IT-TEPPEN frontend application using React + Vite. We need to establish a clear and scalable folder structure that: 

- Separates concerns effectively
- Supports feature-based development
- Distinguishes between types and domain models
- Clarifies the difference between shared utilities and feature-specific shared code
- Maintains good developer experience with Fast Refresh

## Decision

### Overall Structure

```
src/
├── app/                    # Application-wide configuration
│   ├── providers/          # Global providers
│   ├── routes/             # Routing configuration
│   └── App.tsx
│
├── pages/                  # Page components
│   ├── Home/
│   ├── UserList/
│   └── ... 
│
├── features/               # Feature-based modules
│   ├── <featureName>/
│   │   ├── components/     # Feature-specific components
│   │   ├── hooks/          # Feature-specific hooks (including context hooks)
│   │   ├── contexts/       # Context + Provider (components only)
│   │   ├── types/          # Type definitions (interfaces, types)
│   │   ├── models/         # Domain models (Entities, Value Objects with logic)
│   │   ├── api/            # Feature-specific API calls
│   │   └── index.ts        # Public API
│   │
│   └── common/             # Shared across features (domain-specific)
│       ├── components/     # Shared components with business logic
│       ├── hooks/          # Shared hooks (auth, pagination, permission)
│       ├── types/          # Shared types (domain-specific)
│       └── utils/          # Shared utilities (business logic)
│
├── ui/                     # UI component library (design system)
│   ├── Button/
│   ├── Input/
│   ├── Typography/
│   └── index.ts
│
└── shared/                 # Application foundation (domain-agnostic)
    ├── api/                # API client, base fetch configuration
    ├── constants/          # Environment variables, route definitions
    ├── types/              # Global type definitions
    ├── utils/              # Generic utilities (date, string, validation)
    └── lib/                # External library wrappers
```

### Key Principles

#### 1. Separation of Types and Models

- **`types/`**: Pure TypeScript type definitions (interfaces, types)
  - Data structure definitions only
  - No behavior or logic

- **`models/`**: Domain models with business logic
  - Classes or objects with methods
  - Entities and Value Objects (DDD concepts)
  - Encapsulate business rules

**Example:**

```typescript
// features/user/types/user.ts
export interface UserData {
  id: string;
  name: string;
  email: string;
  age: number;
}

// features/user/models/User.ts
import { UserData } from '../types/user';

export class User {
  constructor(private data: UserData) {}

  get displayName(): string {
    return this.data.name;
  }

  isAdult(): boolean {
    return this.data.age >= 18;
  }
}
```

#### 2. Context and Hook Separation

To maintain Fast Refresh compatibility, we separate Context/Provider (components) from their corresponding hooks:

- **`contexts/UserContext. tsx`**: Context + Provider (components only)
- **`hooks/useUserContext.ts`**: Custom hook for consuming the context

**Example:**

```typescript
// features/user/contexts/UserContext.tsx
import { createContext, ReactNode, useState } from 'react';
import { User } from '../models/User';

interface UserContextValue {
  currentUser: User | null;
  login: (user: User) => void;
  logout: () => void;
}

export const UserContext = createContext<UserContextValue | undefined>(undefined);

export function UserProvider({ children }: { children: ReactNode }) {
  const [currentUser, setCurrentUser] = useState<User | null>(null);

  const login = (user: User) => setCurrentUser(user);
  const logout = () => setCurrentUser(null);

  return (
    <UserContext.Provider value={{ currentUser, login, logout }}>
      {children}
    </UserContext.Provider>
  );
}
```

```typescript
// features/user/hooks/useUserContext.ts
import { useContext } from 'react';
import { UserContext } from '../contexts/UserContext';

export function useUserContext() {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error('useUserContext must be used within UserProvider');
  }
  return context;
}
```

#### 3. `shared/` vs `features/common/`

**`shared/`** - Domain-agnostic application foundation: 
- No business logic
- Potentially reusable across different projects
- Generic utilities and configurations
- "Could be extracted into an npm package"

**`features/common/`** - Domain-specific shared code:
- Contains business logic
- Specific to IT-TEPPEN application
- Shared across multiple features
- "Requires knowledge of IT-TEPPEN business rules"

**Decision criteria:**

| Aspect            | `shared/`                        | `features/common/`                    |
| ----------------- | -------------------------------- | ------------------------------------- |
| **Nature**        | Generic, reusable                | App-specific, domain-dependent        |
| **Utils example** | `formatDate`, `debounce`         | `formatUserName`, `calculateDiscount` |
| **Components**    | - (place in `ui/`)               | `DataTable`, `StatusBadge`            |
| **Hooks**         | `useLocalStorage`, `useDebounce` | `useAuth`, `usePermission`            |
| **Types**         | `ApiResponse<T>`, `Result<T, E>` | `Permission`, `UserRole`              |
| **Portability**   | Usable in other projects         | IT-TEPPEN specific                    |

**Simple rule:**
- "Could this be used in another project?" → YES → `shared/`
- "Does this require knowledge of IT-TEPPEN business logic?" → YES → `features/common/`

#### 4. Feature Public API

Each feature should export its public API through `index.ts`:

```typescript
// features/user/index.ts
export { UserTable } from './components/UserTable';
export { useUser } from './hooks/useUser';
export { useUserContext } from './hooks/useUserContext';
export type { UserData, UserRole } from './types/user';
export { User } from './models/User';
```

This encapsulation: 
- Makes dependencies explicit
- Prevents internal implementation details from leaking
- Improves refactoring safety

### File Organization Within Features

```
features/<featureName>/
├── components/          # Feature-specific UI components
├── hooks/               # Custom hooks (including context hooks)
├── contexts/            # React Context + Provider (components only)
├── types/               # TypeScript type definitions
├── models/              # Domain models with business logic
├── api/                 # API calls specific to this feature
└── index.ts             # Public API exports
```

## Consequences

### Positive

- **Clear separation of concerns**:  Types vs models, shared vs common
- **Scalable**:  Easy to add new features without affecting existing code
- **Maintainable**:  Clear boundaries and responsibilities
- **Fast Refresh compatible**: Context/hook separation prevents HMR issues
- **Type-safe**: Explicit public APIs improve type inference
- **Testable**: Isolated features and clear dependencies

### Negative

- **Initial overhead**: Requires understanding of the structure
- **More files**: Separation increases file count
- **Potential over-engineering**: For very small features, this might be excessive

### Mitigation

- Document the structure clearly (this ADR)
- Use code generation/templates for new features
- Follow the "Three Strikes Rule":  Only extract to `features/common/` after seeing the same code three times

## References

- [Feature-Sliced Design](https://feature-sliced.design/)
- [Bulletproof React](https://github.com/alan2207/bulletproof-react)
- [Domain-Driven Design (DDD)](https://en.wikipedia.org/wiki/Domain-driven_design)
- [React Fast Refresh](https://github.com/facebook/react/tree/main/packages/react-refresh)