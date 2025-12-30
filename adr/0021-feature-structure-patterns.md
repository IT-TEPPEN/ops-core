# ADR 0021: Feature Structure Patterns

## Status

Accepted

## Date

2025-12-30

## Context

Features vary greatly in complexity. A simple display-only feature may have 3-5 files, while a complex feature with business logic, local state management, and multiple external integrations may have 30+ files.

**Challenge**: Using the same structure for all features leads to either:

- **Over-engineering**: Small features have unnecessary directories and layers
- **Under-engineering**: Complex features become difficult to maintain as they grow

We need **two structural patterns** that developers can choose based on feature complexity, with a clear migration path between them.

## Decision

We adopt **two feature structure patterns** based on complexity:

1. **Pattern 1 (Simplified)**: For simple features (< 5-10 files)
2. **Pattern 2 (Layered)**: For complex features (10+ files)

### Pattern 1: Simplified Structure

**Use when**:

- Feature has < 10 files total
- Minimal business logic (mostly UI and data fetching)
- No local state management (no LocalStorage, drafts, etc.)
- Simple CRUD operations
- Feature is purely presentational

**Structure**:

```text
features/simpleFeature/
├── components/          # React components
│   ├── FeatureList.tsx
│   ├── FeatureCard.tsx
│   └── index.ts
├── hooks/              # Custom hooks
│   ├── useFeatureList.ts
│   └── index.ts
├── types/              # Type definitions
│   ├── feature.ts
│   └── index.ts
└── index.ts            # Public API
```

**Characteristics**:

- **Flat structure**: No nested layers
- **Simple organization**: Group by technical concern (components, hooks, types)
- **Minimal overhead**: Easy to understand and navigate
- **Quick setup**: Fast to create new features

**Example**: A simple notification display feature

```typescript
// features/notification/components/NotificationList.tsx
export function NotificationList() {
  const { notifications } = useNotifications();
  return (
    <div>
      {notifications.map(n => <NotificationCard key={n.id} notification={n} />)}
    </div>
  );
}

// features/notification/hooks/useNotifications.ts
export function useNotifications() {
  return useQuery({
    queryKey: ['notifications'],
    queryFn: () => fetch('/api/notifications').then(r => r.json()),
  });
}

// features/notification/types/notification.ts
export interface Notification {
  id: string;
  title: string;
  message: string;
  type: 'info' | 'success' | 'warning' | 'error';
}

// features/notification/index.ts
export { NotificationList } from './components/NotificationList';
export { useNotifications } from './hooks/useNotifications';
export type { Notification } from './types/notification';
```

### Pattern 2: Layered Structure

**Use when**:

- Feature has 10+ files
- Complex business logic requiring validation
- Local state management (LocalStorage, drafts, workflows)
- Multiple external integrations
- Requires testing business logic in isolation

**Structure**:

```text
features/complexFeature/
├── domain/                    # Business logic (no React, no HTTP)
│   ├── models/               # Immutable Domain entities
│   │   ├── Feature.ts
│   │   └── index.ts
│   ├── repositories/         # Repository interfaces (local state)
│   │   ├── FeatureRepository.ts
│   │   └── index.ts
│   └── errors/
│       └── FeatureErrors.ts
├── application/              # Use cases (no React hooks)
│   ├── services/            # Query/Command Service interfaces
│   │   ├── FeatureQueryService.ts
│   │   ├── FeatureCommandService.ts
│   │   └── index.ts
│   ├── usecases/            # Business workflows
│   │   ├── CreateFeatureUsecase.ts
│   │   ├── UpdateFeatureUsecase.ts
│   │   └── index.ts
│   └── dto/                 # Application layer DTOs
│       ├── FeatureViewData.ts
│       └── index.ts
├── infrastructure/           # External implementations
│   ├── services/            # HTTP Service implementations
│   │   ├── HttpFeatureQueryService.ts
│   │   ├── HttpFeatureCommandService.ts
│   │   └── index.ts
│   └── repositories/        # LocalStorage/useReducer implementations
│       ├── LocalStorageFeatureRepository.ts
│       └── index.ts
└── presentation/             # React UI layer
    ├── components/          # UI components
    │   ├── FeatureList.tsx
    │   ├── FeatureForm.tsx
    │   └── index.ts
    ├── hooks/               # UI state + Usecase invocation
    │   ├── useFeatureCreate.ts
    │   ├── useFeatureList.ts
    │   └── index.ts
    ├── contexts/            # Dependency Injection
    │   ├── FeatureQueryServiceContext.tsx
    │   └── index.ts
    └── index.ts
```

**Characteristics**:

- **Layered architecture**: Clear separation of concerns (see ADR 0018)
- **Testable business logic**: Domain/Application layers are pure TypeScript
- **Flexible**: Can swap implementations (HTTP → WebSocket → LocalStorage)
- **Scalable**: Can grow without restructuring

**Note**: For detailed layer responsibilities, see **ADR 0018: Frontend Layered Architecture**.

### Pattern Selection Decision Tree

```text
START
  │
  ├─ Does feature have < 10 files?
  │   YES → Pattern 1 (Simplified)
  │   NO ↓
  │
  ├─ Does feature contain complex business logic?
  │   YES → Pattern 2 (Layered)
  │   NO ↓
  │
  ├─ Does feature manage local state (LocalStorage, drafts)?
  │   YES → Pattern 2 (Layered)
  │   NO ↓
  │
  ├─ Does feature require isolated testing of business logic?
  │   YES → Pattern 2 (Layered)
  │   NO → Pattern 1 (Simplified)
```

### Pattern Selection Guidelines

| Criteria | Pattern 1 (Simplified) | Pattern 2 (Layered) |
| -------- | --------------------- | ------------------- |
| **File count** | < 10 files | 10+ files |
| **Business logic** | Minimal | Complex validation/rules |
| **Local state** | None | LocalStorage, drafts, workflows |
| **External APIs** | Simple GET/POST | Multiple integrations |
| **Testing needs** | UI testing only | Isolated business logic tests |
| **Complexity** | Low | High |

### Migration Path: Pattern 1 → Pattern 2

Features should **start simple** (Pattern 1) and **migrate to layered** (Pattern 2) when complexity demands it.

**Migration triggers**:

1. **File count threshold**: Feature grows beyond 10 files
2. **Business logic emergence**: Complex validation or business rules appear
3. **Local state needs**: Need to manage LocalStorage, drafts, or workflows
4. **Testing challenges**: Difficulty testing logic mixed with UI

**Migration steps**:

1. **Create layer directories**:

   ```text
   features/myFeature/
   ├── domain/
   ├── application/
   ├── infrastructure/
   └── presentation/
   ```

2. **Move existing code**:
   - `components/` → `presentation/components/`
   - `hooks/` → `presentation/hooks/`
   - `types/` → `application/dto/` (if view data) or `domain/models/` (if domain models)

3. **Extract business logic**:
   - Create Domain models with business rules
   - Create Usecases for workflows
   - Create Services for API calls

4. **Update imports**: Update import paths to reflect new structure

5. **Add tests**: Add tests for Domain/Application layers

**Migration example**:

```typescript
// BEFORE (Pattern 1)
// features/document/hooks/useDocumentCreate.ts
export function useDocumentCreate() {
  const [isLoading, setIsLoading] = useState(false);

  const create = async (data: CreateDocumentInput) => {
    setIsLoading(true);

    // ❌ Business validation mixed with UI logic
    if (!data.title || data.title.length > 200) {
      throw new Error('Invalid title');
    }

    // ❌ HTTP call directly in hook
    const response = await fetch('/api/documents', {
      method: 'POST',
      body: JSON.stringify(data),
    });

    setIsLoading(false);
    return response.json();
  };

  return { create, isLoading };
}

// AFTER (Pattern 2)
// domain/models/Document.ts
export class Document {
  constructor(private readonly title: string) {
    // ✅ Business validation in Domain
    if (!title || title.length > 200) {
      throw new DomainError('INVALID_TITLE');
    }
  }

  getTitle(): string {
    return this.title;
  }
}

// application/usecases/CreateDocumentUsecase.ts
export class CreateDocumentUsecase {
  constructor(private readonly commandService: DocumentCommandService) {}

  async execute(input: CreateDocumentInput): Promise<DocumentViewData> {
    // ✅ Business logic in Usecase
    const document = new Document(input.title);
    return await this.commandService.create(input);
  }
}

// presentation/hooks/useDocumentCreate.ts
export function useDocumentCreate() {
  const [isLoading, setIsLoading] = useState(false);
  const commandService = useDocumentCommandService();

  const create = async (data: CreateDocumentInput) => {
    setIsLoading(true);
    try {
      // ✅ UI logic only, delegates to Usecase
      const usecase = new CreateDocumentUsecase(commandService);
      return await usecase.execute(data);
    } finally {
      setIsLoading(false);
    }
  };

  return { create, isLoading };
}
```

### Three Strikes Rule

**Rule**: Only refactor to Pattern 2 (or extract to `features/common/`) after encountering the same need **three times**.

**Rationale**:

- **Avoid premature abstraction**: Wait until the pattern is clear
- **YAGNI principle**: "You Aren't Gonna Need It" until you actually do
- **Reduce over-engineering**: Keep code simple until complexity demands otherwise

**Example**:

1. **First time**: Implement form validation in `features/document/components/DocumentForm.tsx`
2. **Second time**: Need similar validation in `features/repository/components/RepositoryForm.tsx` → Copy code
3. **Third time**: Need validation in `features/user/components/UserForm.tsx` → **Now extract** to `features/common/utils/validateForm.ts`

## Consequences

### Positive

1. **Flexibility**: Choose appropriate structure for each feature
2. **Avoid over-engineering**: Small features stay simple
3. **Scalability**: Complex features have room to grow
4. **Clear migration path**: Start simple, add layers when needed
5. **Team autonomy**: Teams can decide based on clear criteria

### Negative

1. **Inconsistency**: Features have different structures
2. **Learning curve**: Developers need to know both patterns
3. **Migration effort**: Refactoring Pattern 1 → Pattern 2 takes time
4. **Judgment calls**: Deciding when to migrate requires experience

### Mitigation

1. **Clear criteria**: This ADR provides decision tree and guidelines
2. **Documentation**: Document pattern choice in feature README
3. **Code reviews**: Validate pattern choice and migration decisions
4. **Pair programming**: Experienced developers guide pattern selection
5. **Gradual migration**: Migrate incrementally, not all at once

## Examples

### Example 1: Simple Feature (Pattern 1)

**Feature**: Display list of user notifications

**Justification**:

- 5 files total
- No business logic (just display)
- Simple data fetching
- No local state

**Structure**:

```text
features/notification/
├── components/
│   ├── NotificationList.tsx
│   ├── NotificationCard.tsx
│   └── index.ts
├── hooks/
│   ├── useNotifications.ts
│   └── index.ts
├── types/
│   ├── notification.ts
│   └── index.ts
└── index.ts
```

### Example 2: Complex Feature (Pattern 2)

**Feature**: Document draft management with auto-save to LocalStorage

**Justification**:

- 20+ files
- Complex business logic (validation, state transitions)
- Local state (LocalStorage for drafts)
- Multiple workflows (create, edit, auto-save, submit)

**Structure**:

```text
features/documentDraft/
├── domain/
│   ├── models/
│   │   ├── DraftDocument.ts      # With validation
│   │   └── DraftWorkflow.ts      # State machine
│   └── repositories/
│       └── DraftRepository.ts     # Interface
├── application/
│   ├── services/
│   │   └── DocumentCommandService.ts
│   ├── usecases/
│   │   ├── SaveDraftUsecase.ts
│   │   ├── AutoSaveDraftUsecase.ts
│   │   └── SubmitDraftUsecase.ts
│   └── dto/
│       └── DraftViewData.ts
├── infrastructure/
│   ├── services/
│   │   └── HttpDocumentCommandService.ts
│   └── repositories/
│       └── LocalStorageDraftRepository.ts
└── presentation/
    ├── components/
    │   ├── DraftEditor.tsx
    │   └── DraftList.tsx
    ├── hooks/
    │   ├── useDraftAutoSave.ts
    │   ├── useDraftSubmit.ts
    │   └── useD

raftList.ts
    └── contexts/
        └── DraftRepositoryContext.tsx
```

### Example 3: Migration Scenario

**Initial state** (Pattern 1):

```text
features/document/
├── components/
│   └── DocumentForm.tsx          # 3 files
├── hooks/
│   └── useDocumentCreate.ts
└── types/
    └── document.ts
```

**After growth** (Pattern 2):

```text
features/document/
├── domain/
│   └── models/
│       └── Document.ts            # Business validation added
├── application/
│   ├── services/
│   │   └── DocumentCommandService.ts
│   └── usecases/
│       ├── CreateDocumentUsecase.ts
│       └── ValidateDocumentUsecase.ts
├── infrastructure/
│   └── services/
│       └── HttpDocumentCommandService.ts
└── presentation/
    ├── components/
    │   ├── DocumentForm.tsx
    │   ├── DocumentList.tsx       # 15 files now
    │   └── DocumentDetail.tsx
    └── hooks/
        ├── useDocumentCreate.ts
        ├── useDocumentList.ts
        └── useDocumentDetail.ts
```

**Triggers**:

- File count grew from 3 → 15
- Business validation requirements appeared
- Multiple workflows needed (create, validate, submit)

## Related ADRs

- **ADR 0012**: Frontend Project Structure - Defines project-level structure
- **ADR 0018**: Frontend Layered Architecture - Defines Pattern 2 layers in detail
- **ADR 0023**: React-Specific Patterns - Patterns applicable to both structures

## References

- [YAGNI (You Aren't Gonna Need It)](https://martinfowler.com/bliki/Yagni.html)
- [Rule of Three (refactoring)](https://en.wikipedia.org/wiki/Rule_of_three_(computer_programming))
- [Evolutionary Architecture](https://www.thoughtworks.com/insights/blog/evolutionary-architectures)
