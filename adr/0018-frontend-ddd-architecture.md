# ADR 0018: Frontend Layered Architecture

## Status

Accepted

## Date

2025-12-30

## Context

Complex features in the frontend application require clear separation between business logic and UI. Without structure:

- Business logic mixes with UI code
- Testing business logic requires rendering components
- Difficult to change implementations (HTTP → WebSocket → LocalStorage)
- Unclear where to place code (hooks? components? utils?)

For complex features (10+ files, business logic, local state management), we need a layered architecture that:

- Separates business logic from UI
- Makes business logic independently testable
- Provides clear boundaries and responsibilities
- Supports dependency injection for flexibility

**Note**: For simple features (< 10 files), use simplified structure (see ADR 0021).

## Decision

We adopt a **four-layer architecture** for complex features, inspired by Onion Architecture and Clean Architecture:

```text
features/complexFeature/
├── domain/                    # Core business logic
│   ├── models/               # Domain entities
│   └── repositories/         # Repository interfaces (local state)
├── application/              # Use cases
│   ├── services/            # Query/Command Service interfaces
│   ├── usecases/            # Business workflows
│   └── dto/                 # Data Transfer Objects
├── infrastructure/           # External implementations
│   ├── services/            # HTTP Service implementations
│   └── repositories/        # LocalStorage implementations
└── presentation/             # React UI
    ├── components/          # UI components
    ├── hooks/               # UI state + Usecase invocation
    └── contexts/            # Dependency injection
```

**Dependency direction**: Always **inward**

```text
Presentation → Application → Domain
Infrastructure → Application/Domain
```

- **Domain** has no dependencies (pure TypeScript/JavaScript)
- **Application** depends on Domain only
- **Infrastructure** implements interfaces from Domain/Application
- **Presentation** depends on Application (and uses Infrastructure via DI)

### Layer 1: Domain Layer

**Responsibility**: Core business entities and rules (frontend-specific)

**Location**: `domain/`

**Contains**:

- **`models/`**: Immutable Domain entities with business logic
- **`repositories/`**: Repository interfaces for local state management
- **`errors/`**: Domain-specific errors

**Rules**:

- ❌ No dependencies on React, UI libraries, or HTTP clients
- ❌ No external service calls
- ✅ Pure TypeScript/JavaScript logic only
- ✅ Contains business validation (domain invariants)
- ✅ Immutable data structures (return new instances on updates)

**Example**:

```typescript
// domain/models/DraftDocument.ts
export class DraftDocument {
  private constructor(
    private readonly id: string,
    private readonly title: string,
    private readonly content: string,
    private readonly lastSavedAt: Date | null
  ) {}

  static create(data: { title: string }): DraftDocument {
    return new DraftDocument(crypto.randomUUID(), data.title, '', null);
  }

  static reconstruct(data: {
    id: string;
    title: string;
    content: string;
    lastSavedAt: string | null;
  }): DraftDocument {
    return new DraftDocument(
      data.id,
      data.title,
      data.content,
      data.lastSavedAt ? new Date(data.lastSavedAt) : null
    );
  }

  getId(): string { return this.id; }
  getTitle(): string { return this.title; }
  getContent(): string { return this.content; }
  getLastSavedAt(): Date | null { return this.lastSavedAt; }

  // ✅ Business logic: validation
  canSubmit(): boolean {
    return this.title.length > 0 && this.content.length > 0;
  }

  // ✅ Immutable update
  updateTitle(newTitle: string): DraftDocument {
    if (newTitle.length > 200) {
      throw new DomainError('TITLE_TOO_LONG', 'Title must be 200 characters or less');
    }
    return new DraftDocument(this.id, newTitle, this.content, this.lastSavedAt);
  }

  markAsSaved(): DraftDocument {
    return new DraftDocument(this.id, this.title, this.content, new Date());
  }
}

// domain/repositories/DraftDocumentRepository.ts
export interface DraftDocumentRepository {
  get(id: string): DraftDocument | null;
  getAll(): DraftDocument[];
  add(draft: DraftDocument): void;
  update(draft: DraftDocument): void;
  remove(id: string): void;
}
```

**What belongs in Domain**:

- **Frontend-specific business logic**: Validation, state transitions for frontend-only entities
- **Local state entities**: Drafts, temporary state, workflows
- **Domain invariants**: Rules that must always be true

**What does NOT belong in Domain**:

- **Backend-managed entities**: These are represented as ViewData (Application layer)
- **API-driven logic**: Backend enforces this, frontend just displays

### Layer 2: Application Layer

**Responsibility**: Use cases and business workflows

**Location**: `application/`

**Contains**:

- **`services/`**: Query/Command Service interfaces (external API)
- **`usecases/`**: Business workflows, coordinate services/repositories
- **`dto/`**: Data Transfer Objects (ViewData for display)
- **`errors/`**: Application-specific errors

**Rules**:

- ❌ No React hooks (useState, useEffect, useContext, etc.)
- ❌ No direct HTTP calls (use Service interfaces)
- ✅ Depends on Domain layer only
- ✅ Contains business validation (cross-entity rules, external checks)
- ✅ Pure async functions, fully testable without UI

**Example**:

```typescript
// application/services/DocumentQueryService.ts (interface)
export interface DocumentQueryService {
  getById(id: string): Promise<DocumentViewData>;
  list(filters: DocumentFilters): Promise<DocumentViewData[]>;
}

// application/services/DocumentCommandService.ts (interface)
export interface DocumentCommandService {
  create(request: CreateDocumentRequest): Promise<DocumentViewData>;
  publish(id: string): Promise<DocumentViewData>;
}

// application/dto/DocumentViewData.ts
export interface DocumentViewData {
  id: string;
  title: string;
  createdAt: Date;       // Display-ready
  ownerName: string;     // Display-ready
}

// application/usecases/PublishDocumentUsecase.ts
export class PublishDocumentUsecase {
  constructor(
    private readonly queryService: DocumentQueryService,
    private readonly commandService: DocumentCommandService
  ) {}

  async execute(input: {
    documentId: string;
    userId: string;
  }): Promise<DocumentViewData> {
    // ✅ Business validation: fetch document
    const document = await this.queryService.getById(input.documentId);

    // ✅ Business validation: check status (lightweight frontend check)
    if (document.status !== 'draft') {
      throw new ApplicationError('DOCUMENT_NOT_DRAFT', 'Only drafts can be published');
    }

    // ✅ Execute operation via Command Service
    return await this.commandService.publish(input.documentId);
  }
}
```

**What belongs in Application**:

- **Usecases**: Business workflows (publish document, submit draft, etc.)
- **Service interfaces**: Contracts for external operations (Query/Command)
- **ViewData**: Display-ready data for Presentation layer
- **Cross-entity validation**: Rules spanning multiple entities

**What does NOT belong in Application**:

- **UI state**: Loading, error states (belongs in Presentation)
- **Service implementations**: HTTP/LocalStorage details (belongs in Infrastructure)

### Layer 3: Infrastructure Layer

**Responsibility**: External service implementations

**Location**: `infrastructure/`

**Contains**:

- **`services/`**: HTTP implementations of Query/Command Services
- **`repositories/`**: LocalStorage/SessionStorage implementations of Repositories
- **`types.ts`**: Internal API response types (not exported)

**Rules**:

- ✅ Implements interfaces from Domain/Application
- ✅ Handles HTTP communication, LocalStorage, etc.
- ✅ Transforms API responses to ViewData or Domain models
- ❌ No business logic (just data access)

**Example**:

```typescript
// infrastructure/services/HttpDocumentQueryService.ts
export class HttpDocumentQueryService implements DocumentQueryService {
  constructor(private readonly apiClient: ApiClient) {}

  async getById(id: string): Promise<DocumentViewData> {
    const response = await this.apiClient.get<ApiDocumentResponse>(`/documents/${id}`);

    // ✅ Transform API response to ViewData
    return {
      id: response.id,
      title: response.title,
      createdAt: new Date(response.created_at),
      ownerName: response.owner_name,
    };
  }
}

// infrastructure/repositories/LocalStorageDraftRepository.ts
export class LocalStorageDraftRepository implements DraftDocumentRepository {
  private static STORAGE_KEY = 'draft_documents';

  get(id: string): DraftDocument | null {
    const stored = localStorage.getItem(LocalStorageDraftRepository.STORAGE_KEY);
    if (!stored) return null;

    const drafts: string[] = JSON.parse(stored);
    const draftJson = drafts.find(json => {
      const draft = DraftDocument.deserialize(json);
      return draft.getId() === id;
    });

    return draftJson ? DraftDocument.deserialize(draftJson) : null;
  }

  add(draft: DraftDocument): void {
    const drafts = this.getAll();
    drafts.push(draft);
    this.persist(drafts);
  }

  private persist(drafts: DraftDocument[]): void {
    const serialized = drafts.map(d => d.serialize());
    localStorage.setItem(LocalStorageDraftRepository.STORAGE_KEY, JSON.stringify(serialized));
  }

  // ... other methods
}
```

**What belongs in Infrastructure**:

- **HTTP implementations**: API client wrappers
- **LocalStorage implementations**: Persistence logic
- **Data transformation**: API response → ViewData/Domain Model

**What does NOT belong in Infrastructure**:

- **Business logic**: Validation, business rules (belongs in Domain/Application)
- **UI logic**: Component state, rendering (belongs in Presentation)

### Layer 4: Presentation Layer

**Responsibility**: UI rendering, user interaction, UI state management

**Location**: `presentation/`

**Contains**:

- **`components/`**: React components
- **`hooks/`**: UI state management + Usecase invocation
- **`contexts/`**: Dependency injection (Context Providers)

**Rules**:

- ✅ Uses React hooks freely (useState, useEffect, useContext, etc.)
- ✅ Manages UI state (loading, error, form state)
- ✅ Invokes Usecases (via hooks)
- ✅ Handles input validation (react-hook-form + zod)
- ❌ No business logic (delegates to Usecases)
- ❌ No direct business validation (uses Domain models/Usecases)

**Example**:

```typescript
// presentation/contexts/DocumentQueryServiceContext.tsx
const DocumentQueryServiceContext = createContext<DocumentQueryService | null>(null);

export function DocumentQueryServiceProvider({ children }: { children: ReactNode }) {
  const service = useMemo(() => {
    const apiClient = new ApiClient('/api/v1');
    return new HttpDocumentQueryService(apiClient);
  }, []);

  return (
    <DocumentQueryServiceContext.Provider value={service}>
      {children}
    </DocumentQueryServiceContext.Provider>
  );
}

// presentation/hooks/useDocumentQueryService.ts
export function useDocumentQueryService(): DocumentQueryService {
  const context = useContext(DocumentQueryServiceContext);
  if (!context) {
    throw new Error('useDocumentQueryService must be used within provider');
  }
  return context;
}

// presentation/hooks/useDocumentPublish.ts
export function useDocumentPublish() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const queryService = useDocumentQueryService();
  const commandService = useDocumentCommandService();
  const { currentUser } = useAuth();

  const publish = async (documentId: string) => {
    setIsLoading(true);
    setError(null);

    try {
      // ✅ Instantiate usecase with dependencies
      const usecase = new PublishDocumentUsecase(queryService, commandService);

      // ✅ Execute usecase
      const result = await usecase.execute({
        documentId,
        userId: currentUser!.id,
      });

      return result;
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  return { publish, isLoading, error };
}

// presentation/components/DocumentCard.tsx
function DocumentCard({ document }: { document: DocumentViewData }) {
  const { publish, isLoading } = useDocumentPublish();

  const handlePublish = async () => {
    try {
      await publish(document.id);
    } catch (err) {
      // Error already handled by hook
    }
  };

  return (
    <div>
      <h3>{document.title}</h3>
      <p>By {document.ownerName}</p>
      <button onClick={handlePublish} disabled={isLoading}>
        Publish
      </button>
    </div>
  );
}
```

**What belongs in Presentation**:

- **UI components**: React components
- **UI state**: Loading, error, form state
- **Input validation**: react-hook-form + zod
- **Usecase invocation**: Calling business logic
- **Context providers**: Dependency injection setup

**What does NOT belong in Presentation**:

- **Business logic**: Validation, business rules
- **Data transformation**: API response conversion
- **Direct API calls**: Use Services via Usecases

## Layer Interaction Patterns

### Pattern 1: Query (Read Display Data)

```text
Presentation
  ↓ (useQuery)
Query Service (Infrastructure)
  ↓ (HTTP GET)
Backend API
  ↓ (return ViewData)
Presentation (display)
```

**No Usecase needed** for simple queries.

### Pattern 2: Command (Write with Business Logic)

```text
Presentation
  ↓ (invoke Usecase)
Usecase (Application)
  ↓ (validate via Query Service)
Query Service (Infrastructure)
  ↓ (execute via Command Service)
Command Service (Infrastructure)
  ↓ (HTTP POST/PUT)
Backend API
```

**Usecase coordinates** validation and execution.

### Pattern 3: Local State Management

```text
Presentation
  ↓ (invoke Usecase)
Usecase (Application)
  ↓ (get from Repository)
Repository Interface (Domain)
  ↓ (LocalStorage implementation)
LocalStorageDraftRepository (Infrastructure)
  ↓ (return Domain Model)
Domain Model (Domain)
```

**Usecase manages** local state via Repository.

## Consequences

### Positive

1. **Clear separation**: Business logic isolated from UI
2. **Testable**: Domain/Application layers are pure TypeScript
3. **Flexible**: Swap implementations (HTTP → WebSocket → LocalStorage)
4. **Maintainable**: Easy to locate responsibilities
5. **Scalable**: Features can grow without restructuring

### Negative

1. **More files**: Layered structure requires more files
2. **Indirection**: More layers between UI and data
3. **Learning curve**: Team needs to understand layers
4. **Over-engineering risk**: May be excessive for simple features

### Mitigation

1. **Use simplified structure for small features**: See ADR 0021
2. **Clear migration path**: Start simple, add layers when needed
3. **Documentation**: This ADR + related ADRs
4. **Code reviews**: Ensure layer boundaries are maintained

## When to Use This Architecture

**Use layered architecture when**:

- Feature has 10+ files
- Complex business logic requiring validation
- Local state management (LocalStorage, drafts, workflows)
- Multiple external integrations
- Requires testing business logic in isolation

**Use simplified structure when**:

- Feature has < 10 files
- Minimal business logic (mostly UI)
- Simple CRUD operations
- No local state management

See **ADR 0021: Feature Structure Patterns** for detailed criteria.

## Complete Example

See ADR 0019, 0020, 0021, 0022, 0023 for complete examples of:

- Query/Command Service patterns (ADR 0019)
- Validation strategy (ADR 0020)
- Feature structure migration (ADR 0021)
- Type placement (ADR 0022)
- Immutability and DI (ADR 0023)

## Related ADRs

- **ADR 0012**: Frontend Project Structure - Project-level organization
- **ADR 0019**: Data Access Patterns - Query/Command/Repository details
- **ADR 0020**: Validation Strategy - Where to validate
- **ADR 0021**: Feature Structure Patterns - When to use this architecture
- **ADR 0022**: types/ Directory - Type placement rules
- **ADR 0023**: React-Specific Patterns - Immutability, Context/Hook separation

## References

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Onion Architecture](https://jeffreypalermo.com/2008/07/the-onion-architecture-part-1/)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
