# ADR 0023: React-Specific Implementation Patterns

## Status

Accepted

## Date

2025-12-30

## Context

React has specific requirements and best practices that affect how we implement domain models and organize code:

1. **React Strict Mode**: Requires pure functions and immutable data
2. **Fast Refresh (HMR)**: Can fail when component files mix concerns
3. **Context API**: Idiomatic way to provide dependencies
4. **Module boundaries**: Need clear public APIs for features

These React-specific concerns require patterns that may differ from traditional OOP or backend architectures.

## Decision

We adopt the following React-specific patterns across the frontend application.

### Pattern 1: Immutable Domain Models

**Decision**: All Domain models are strictly immutable to comply with React Strict Mode.

**Rationale**:

- React Strict Mode requires pure functions
- React state updates detect changes by reference equality
- Immutability prevents accidental mutations
- Predictable state management

**Implementation**:

```typescript
export class DocumentEntity {
  // ✅ Private readonly constructor
  private constructor(
    private readonly id: string,
    private readonly title: string,
    private readonly content: string,
    private readonly status: 'draft' | 'published'
  ) {}

  // ✅ Static factory method (validation)
  static create(data: { title: string; content: string }): DocumentEntity {
    if (data.title.length === 0) {
      throw new DomainError('EMPTY_TITLE', 'Title cannot be empty');
    }
    return new DocumentEntity(
      crypto.randomUUID(),
      data.title,
      data.content,
      'draft'
    );
  }

  // ✅ Static reconstructor (from persistence)
  static reconstruct(data: {
    id: string;
    title: string;
    content: string;
    status: 'draft' | 'published';
  }): DocumentEntity {
    return new DocumentEntity(data.id, data.title, data.content, data.status);
  }

  // ✅ Getter methods
  getId(): string { return this.id; }
  getTitle(): string { return this.title; }
  getContent(): string { return this.content; }
  getStatus(): 'draft' | 'published' { return this.status; }

  // ✅ Update methods return NEW instances
  updateTitle(newTitle: string): DocumentEntity {
    if (newTitle.length === 0) {
      throw new DomainError('EMPTY_TITLE', 'Title cannot be empty');
    }
    return new DocumentEntity(this.id, newTitle, this.content, this.status);
  }

  updateContent(newContent: string): DocumentEntity {
    return new DocumentEntity(this.id, this.title, newContent, this.status);
  }

  publish(): DocumentEntity {
    if (this.status !== 'draft') {
      throw new DomainError('ALREADY_PUBLISHED', 'Document is already published');
    }
    return new DocumentEntity(this.id, this.title, this.content, 'published');
  }
}
```

**React usage**:

```typescript
function DocumentEditor() {
  const [document, setDocument] = useState(() =>
    DocumentEntity.create({ title: 'New Doc', content: '' })
  );

  const handleTitleChange = (newTitle: string) => {
    // ✅ Returns new instance, React detects change
    setDocument(prev => prev.updateTitle(newTitle));
  };

  const handlePublish = () => {
    // ✅ Returns new instance
    setDocument(prev => prev.publish());
  };

  return (
    <div>
      <input
        value={document.getTitle()}
        onChange={e => handleTitleChange(e.target.value)}
      />
      <button onClick={handlePublish}>Publish</button>
    </div>
  );
}
```

**Key principles**:

- **Private constructor**: Enforce factory methods
- **Readonly fields**: Prevent mutations
- **Factory methods**: `create()` for new entities, `reconstruct()` for existing
- **Getter methods**: Access data
- **Update methods**: Return new instances

### Pattern 2: Fast Refresh Compatible Context/Hook Separation

**Decision**: Separate Context/Provider (components) from consumption hooks to maintain Fast Refresh compatibility.

**Rationale**:

- React Fast Refresh can fail when a file mixes component definitions with hooks
- Separating them ensures reliable Hot Module Replacement (HMR)
- Improves development experience

**Implementation**:

```typescript
// ✅ presentation/contexts/DocumentRepositoryContext.tsx
import { createContext, ReactNode, useMemo } from 'react';
import { DocumentRepository } from '../../domain/repositories/DocumentRepository';
import { LocalStorageDocumentRepository } from '../../infrastructure/repositories/LocalStorageDocumentRepository';

// 1. Context definition
export const DocumentRepositoryContext = createContext<DocumentRepository | null>(null);

// 2. Provider component (components only in this file)
export function DocumentRepositoryProvider({ children }: { children: ReactNode }) {
  const repository = useMemo(() => {
    return new LocalStorageDocumentRepository();
  }, []);

  return (
    <DocumentRepositoryContext.Provider value={repository}>
      {children}
    </DocumentRepositoryContext.Provider>
  );
}

// ✅ presentation/hooks/useDocumentRepository.ts
import { useContext } from 'react';
import { DocumentRepositoryContext } from '../contexts/DocumentRepositoryContext';
import { DocumentRepository } from '../../domain/repositories/DocumentRepository';

// 3. Consumption hook (separate file)
export function useDocumentRepository(): DocumentRepository {
  const context = useContext(DocumentRepositoryContext);
  if (!context) {
    throw new Error(
      'useDocumentRepository must be used within DocumentRepositoryProvider'
    );
  }
  return context;
}
```

**Usage**:

```typescript
// App-level setup
function App() {
  return (
    <DocumentRepositoryProvider>
      <Routes />
    </DocumentRepositoryProvider>
  );
}

// Component usage
function DocumentList() {
  const repository = useDocumentRepository();
  // Use repository...
}
```

**Alternative (co-located)**: For simple cases where Fast Refresh isn't an issue:

```typescript
// contexts/DocumentRepositoryContext.tsx
export const DocumentRepositoryContext = createContext<DocumentRepository | null>(null);

export function DocumentRepositoryProvider({ children }: { children: ReactNode }) {
  // Provider implementation
}

// Export hook from same file
export function useDocumentRepository(): DocumentRepository {
  const context = useContext(DocumentRepositoryContext);
  if (!context) {
    throw new Error('useDocumentRepository must be used within DocumentRepositoryProvider');
  }
  return context;
}
```

### Pattern 3: Feature Public API (index.ts)

**Decision**: Each feature exports its public API through `index.ts` at the feature root.

**Rationale**:

- **Encapsulation**: Internal implementation details remain private
- **Explicit dependencies**: Other features can only import through public API
- **Refactoring safety**: Can change internal structure without breaking consumers
- **Clear boundaries**: Easy to see what's "public" vs "internal"

**Implementation**:

```typescript
// features/document/index.ts

// ✅ Export public components
export { DocumentTable } from './presentation/components/DocumentTable';
export { DocumentForm } from './presentation/components/DocumentForm';
export { DocumentDetail } from './presentation/components/DocumentDetail';

// ✅ Export public hooks
export { useDocumentCreate } from './presentation/hooks/useDocumentCreate';
export { useDocumentList } from './presentation/hooks/useDocumentList';
export { useDocumentDetail } from './presentation/hooks/useDocumentDetail';

// ✅ Export public types (DTOs only, not internal types)
export type { DocumentViewData } from './application/dto/DocumentViewData';
export type { CreateDocumentInput } from './application/dto/CreateDocumentInput';

// ✅ Export Context Providers (for app-level setup)
export { DocumentRepositoryProvider } from './presentation/contexts/DocumentRepositoryContext';
export { DocumentQueryServiceProvider } from './presentation/contexts/DocumentQueryServiceContext';

// ❌ Do NOT export internal implementation details
// ❌ Do NOT export: Domain models, Usecases, Infrastructure implementations
```

**Usage in other features/pages**:

```typescript
// ✅ Import from feature root
import {
  DocumentTable,
  useDocumentList,
  type DocumentViewData
} from '@/features/document';

// ❌ Avoid: Direct imports from internal structure
import { DocumentTable } from '@/features/document/presentation/components/DocumentTable';
```

**Benefits**:

- Changing `presentation/components/DocumentTable.tsx` to `presentation/components/DocumentList.tsx` only requires updating `index.ts`
- Consumers don't need to know internal structure
- TypeScript can enforce boundaries (using path mapping)

### Pattern 4: Context-Based Dependency Injection

**Decision**: Use React Context for dependency injection of repositories and services.

**Rationale**:

- **React-idiomatic**: Uses built-in Context API
- **Testable**: Easy to provide mock implementations
- **Flexible**: Can swap implementations (HTTP → WebSocket → LocalStorage)
- **Type-safe**: TypeScript ensures correct types

**Implementation**:

```typescript
// 1. Define interface in Domain/Application layer
// domain/repositories/DocumentRepository.ts
export interface DocumentRepository {
  get(id: string): Document | null;
  getAll(): Document[];
  add(document: Document): void;
  update(document: Document): void;
  remove(id: string): void;
}

// 2. Implement in Infrastructure layer
// infrastructure/repositories/LocalStorageDocumentRepository.ts
export class LocalStorageDocumentRepository implements DocumentRepository {
  private static STORAGE_KEY = 'documents';

  get(id: string): Document | null {
    const stored = localStorage.getItem(LocalStorageDocumentRepository.STORAGE_KEY);
    if (!stored) return null;
    const documents = JSON.parse(stored);
    return documents.find((d: any) => d.id === id) || null;
  }

  // ... other methods
}

// 3. Create Context and Provider in Presentation layer
// presentation/contexts/DocumentRepositoryContext.tsx
const DocumentRepositoryContext = createContext<DocumentRepository | null>(null);

export function DocumentRepositoryProvider({ children }: { children: ReactNode }) {
  const repository = useMemo(
    () => new LocalStorageDocumentRepository(),
    []
  );

  return (
    <DocumentRepositoryContext.Provider value={repository}>
      {children}
    </DocumentRepositoryContext.Provider>
  );
}

// 4. Create consumption hook
// presentation/hooks/useDocumentRepository.ts
export function useDocumentRepository(): DocumentRepository {
  const context = useContext(DocumentRepositoryContext);
  if (!context) {
    throw new Error('useDocumentRepository must be used within DocumentRepositoryProvider');
  }
  return context;
}

// 5. Use in components/hooks
// presentation/hooks/useDocumentCreate.ts
export function useDocumentCreate() {
  const repository = useDocumentRepository(); // ✅ Get via Context

  const create = async (data: CreateDocumentInput) => {
    const usecase = new CreateDocumentUsecase(repository);
    return await usecase.execute(data);
  };

  return { create };
}
```

**Testing with mocks**:

```typescript
// __tests__/useDocumentCreate.test.ts
import { renderHook } from '@testing-library/react';
import { useDocumentCreate } from '../useDocumentCreate';
import { DocumentRepositoryContext } from '../contexts/DocumentRepositoryContext';

describe('useDocumentCreate', () => {
  it('should create document', async () => {
    // ✅ Provide mock repository
    const mockRepository: DocumentRepository = {
      add: vi.fn(),
      // ... other methods
    };

    const wrapper = ({ children }) => (
      <DocumentRepositoryContext.Provider value={mockRepository}>
        {children}
      </DocumentRepositoryContext.Provider>
    );

    const { result } = renderHook(() => useDocumentCreate(), { wrapper });

    await result.current.create({ title: 'Test', content: 'Content' });

    expect(mockRepository.add).toHaveBeenCalled();
  });
});
```

**Environment-specific implementations**:

```typescript
// app/providers/index.tsx
export function AppProviders({ children }: { children: ReactNode }) {
  return (
    <DocumentRepositoryProvider implementation={
      // ✅ Choose implementation based on environment
      import.meta.env.VITE_STORAGE_TYPE === 'local'
        ? 'localStorage'
        : 'indexedDB'
    }>
      {children}
    </DocumentRepositoryProvider>
  );
}
```

## Consequences

### Positive

1. **React compatibility**: Patterns work seamlessly with React features
2. **Immutability**: Prevents bugs from accidental mutations
3. **Fast Refresh**: Reliable HMR during development
4. **Testability**: Easy to mock dependencies
5. **Type safety**: TypeScript enforces correct usage
6. **Flexibility**: Can swap implementations easily

### Negative

1. **Boilerplate**: More files and setup code
2. **Learning curve**: Developers need to understand patterns
3. **Verbosity**: Immutable models require more method definitions
4. **Context overhead**: Many Providers in App root

### Mitigation

1. **Code generation**: Templates for common patterns
2. **Documentation**: This ADR + code examples
3. **IDE snippets**: Quick scaffolding
4. **Code reviews**: Ensure patterns are followed

## Examples

### Example 1: Immutable Model with Complex State

```typescript
// domain/models/DraftDocument.ts
export class DraftDocument {
  private constructor(
    private readonly id: string,
    private readonly title: string,
    private readonly content: string,
    private readonly lastSavedAt: Date | null,
    private readonly isDirty: boolean
  ) {}

  static create(data: { title: string }): DraftDocument {
    return new DraftDocument(
      crypto.randomUUID(),
      data.title,
      '',
      null,
      false
    );
  }

  getId(): string { return this.id; }
  getTitle(): string { return this.title; }
  getContent(): string { return this.content; }
  getLastSavedAt(): Date | null { return this.lastSavedAt; }
  isDirty(): boolean { return this.isDirty; }

  updateTitle(newTitle: string): DraftDocument {
    return new DraftDocument(
      this.id,
      newTitle,
      this.content,
      this.lastSavedAt,
      true // Mark as dirty
    );
  }

  updateContent(newContent: string): DraftDocument {
    return new DraftDocument(
      this.id,
      this.title,
      newContent,
      this.lastSavedAt,
      true
    );
  }

  markAsSaved(): DraftDocument {
    return new DraftDocument(
      this.id,
      this.title,
      this.content,
      new Date(),
      false // Clear dirty flag
    );
  }

  // Serialization for LocalStorage
  serialize(): string {
    return JSON.stringify({
      id: this.id,
      title: this.title,
      content: this.content,
      lastSavedAt: this.lastSavedAt?.toISOString() ?? null,
      isDirty: this.isDirty,
    });
  }

  static deserialize(json: string): DraftDocument {
    const data = JSON.parse(json);
    return new DraftDocument(
      data.id,
      data.title,
      data.content,
      data.lastSavedAt ? new Date(data.lastSavedAt) : null,
      data.isDirty
    );
  }
}
```

### Example 2: Multiple Providers Composition

```typescript
// app/providers/index.tsx
export function AppProviders({ children }: { children: ReactNode }) {
  return (
    <AuthProvider>
      <ThemeProvider>
        <DocumentRepositoryProvider>
          <DocumentQueryServiceProvider>
            <NotificationProvider>
              {children}
            </NotificationProvider>
          </DocumentQueryServiceProvider>
        </DocumentRepositoryProvider>
      </ThemeProvider>
    </AuthProvider>
  );
}

// app/App.tsx
function App() {
  return (
    <AppProviders>
      <RouterProvider router={router} />
    </AppProviders>
  );
}
```

### Example 3: Feature Public API with Re-exports

```typescript
// features/document/index.ts

// Re-export from sub-features
export * from './draft';
export * from './version';

// Main feature exports
export { DocumentTable } from './presentation/components/DocumentTable';
export { DocumentForm } from './presentation/components/DocumentForm';

export { useDocumentCreate } from './presentation/hooks/useDocumentCreate';
export { useDocumentList } from './presentation/hooks/useDocumentList';

export type { DocumentViewData } from './application/dto/DocumentViewData';

export { DocumentRepositoryProvider } from './presentation/contexts/DocumentRepositoryContext';

// features/document/draft/index.ts (sub-feature)
export { DraftEditor } from './presentation/components/DraftEditor';
export { useDraftAutoSave } from './presentation/hooks/useDraftAutoSave';
export type { DraftViewData } from './application/dto/DraftViewData';
```

## Related ADRs

- **ADR 0018**: Frontend Layered Architecture - Layer definitions
- **ADR 0021**: Feature Structure Patterns - Feature organization
- **ADR 0022**: types/ Directory - Type placement rules

## References

- [React Strict Mode](https://react.dev/reference/react/StrictMode)
- [React Fast Refresh](https://github.com/facebook/react/tree/main/packages/react-refresh)
- [React Context](https://react.dev/learn/passing-data-deeply-with-context)
- [Immutability in React](https://react.dev/learn/updating-objects-in-state)
