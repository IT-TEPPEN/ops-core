# ADR 0019: Data Access Patterns (Query Service, Command Service, Repository)

## Status

Accepted

## Date

2025-12-30

## Context

Frontend applications interact with data from two distinct sources:

1. **External APIs**: Data managed by backend services (CRUD operations)
2. **Local state**: Data managed within the frontend (LocalStorage, drafts, workflows)

**Challenge**: Traditional backend patterns (Repository for all data access) don't fit well in frontend:

- **Backend Repository**: Abstracts database operations, manages persistence
- **Frontend "Repository"**: What does it manage? API calls? LocalStorage? Both?

Mixing concerns leads to confusion:

- HTTP clients called "Repositories"
- Business logic scattered between API calls and local state
- Unclear responsibility boundaries

We need **clear patterns** that distinguish between external API operations and local state management.

## Decision

We adopt **three distinct patterns** for data access:

1. **Query Service**: External read operations (GET)
2. **Command Service**: External write operations (POST/PUT/DELETE)
3. **Repository**: Local state management (LocalStorage, useReducer)

### Pattern A: Query Service (External Read Operations)

**Purpose**: Fetch data from external APIs and provide display-ready format

**Use when**:

- Simple data fetching from backend
- Read-only operations
- Data for display purposes

**Key characteristics**:

- Interface defined in **Application layer**
- Returns **ViewData** (display-ready format)
- Handles data transformations (date strings → Date, snake_case → camelCase)
- **No business logic** (simple fetching and transformation)
- Works well with React Query/SWR

**Example**:

```typescript
// application/services/DocumentQueryService.ts
export interface DocumentQueryService {
  getById(id: string): Promise<DocumentViewData>;
  list(filters: DocumentFilters): Promise<DocumentViewData[]>;
  search(query: string): Promise<DocumentViewData[]>;
}

// application/dto/DocumentViewData.ts
export interface DocumentViewData {
  id: string;
  title: string;
  content: string;
  createdAt: Date;        // ✅ Transformed from string
  updatedAt: Date;        // ✅ Transformed from string
  ownerName: string;      // ✅ Display-ready
}

// infrastructure/services/HttpDocumentQueryService.ts
export class HttpDocumentQueryService implements DocumentQueryService {
  constructor(private readonly apiClient: ApiClient) {}

  async getById(id: string): Promise<DocumentViewData> {
    const response = await this.apiClient.get<ApiDocumentResponse>(
      `/documents/${id}`
    );

    // ✅ Transform to display-ready format
    return {
      id: response.id,
      title: response.title,
      content: response.content,
      createdAt: new Date(response.created_at),
      updatedAt: new Date(response.updated_at),
      ownerName: response.owner_name,
    };
  }

  async list(filters: DocumentFilters): Promise<DocumentViewData[]> {
    const response = await this.apiClient.get<ApiDocumentListResponse>(
      `/documents`,
      { params: filters }
    );

    return response.documents.map(doc => this.toViewData(doc));
  }

  private toViewData(response: ApiDocumentResponse): DocumentViewData {
    return {
      id: response.id,
      title: response.title,
      content: response.content,
      createdAt: new Date(response.created_at),
      updatedAt: new Date(response.updated_at),
      ownerName: response.owner_name,
    };
  }
}

// presentation/hooks/useDocumentList.ts
export function useDocumentList(filters: DocumentFilters) {
  const queryService = useDocumentQueryService();

  // ✅ Simple usage with React Query (Usecase not needed for simple queries)
  return useQuery({
    queryKey: ['documents', filters],
    queryFn: () => queryService.list(filters),
  });
}

// presentation/components/DocumentList.tsx
function DocumentList() {
  const { data: documents, isLoading } = useDocumentList({});

  if (isLoading) return <Loading />;

  return (
    <ul>
      {documents?.map(doc => (
        <li key={doc.id}>
          <h3>{doc.title}</h3>
          {/* ✅ No transformation needed, data is display-ready */}
          <time>{doc.createdAt.toLocaleDateString()}</time>
          <p>By {doc.ownerName}</p>
        </li>
      ))}
    </ul>
  );
}
```

**Benefits**:

- **Consistency**: All components receive the same data format
- **No duplicate transformation**: Done once in Service
- **Simple**: No business logic, just fetch and transform
- **Cacheable**: Works perfectly with React Query/SWR

**When to skip Usecase**: For simple queries with no business logic, Presentation can call Query Service directly.

### Pattern B: Command Service (External Write Operations)

**Purpose**: Execute data changes on external APIs after business validation

**Use when**:

- Sending changes to backend (POST/PUT/DELETE)
- Operations that modify server state
- After Usecase has performed business validation

**Key characteristics**:

- Interface defined in **Application layer**
- **No business logic** (Usecase validates before calling)
- Simple request/response handling
- Returns ViewData
- Backend is responsible for detailed business logic

**Example**:

```typescript
// application/services/DocumentCommandService.ts
export interface DocumentCommandService {
  create(request: CreateDocumentRequest): Promise<DocumentViewData>;
  update(id: string, request: UpdateDocumentRequest): Promise<DocumentViewData>;
  delete(id: string): Promise<void>;
  publish(id: string): Promise<DocumentViewData>;
}

export interface CreateDocumentRequest {
  repositoryId: string;
  filePath: string;
  title: string;
  content: string;
  accessScope: 'public' | 'private';
}

// infrastructure/services/HttpDocumentCommandService.ts
export class HttpDocumentCommandService implements DocumentCommandService {
  constructor(private readonly apiClient: ApiClient) {}

  async create(request: CreateDocumentRequest): Promise<DocumentViewData> {
    // ✅ Simple API call, no business logic
    const response = await this.apiClient.post<ApiDocumentResponse>(
      '/documents',
      request
    );

    return this.toViewData(response);
  }

  async publish(id: string): Promise<DocumentViewData> {
    // ✅ Simple API call
    const response = await this.apiClient.post<ApiDocumentResponse>(
      `/documents/${id}/publish`,
      {}
    );

    return this.toViewData(response);
  }

  private toViewData(response: ApiDocumentResponse): DocumentViewData {
    return {
      id: response.id,
      title: response.title,
      content: response.content,
      createdAt: new Date(response.created_at),
      updatedAt: new Date(response.updated_at),
      ownerName: response.owner_name,
    };
  }
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
    // ✅ Business validation: Fetch document
    const document = await this.queryService.getById(input.documentId);

    // ✅ Business validation: Check status (lightweight frontend check)
    if (document.status !== 'draft') {
      throw new ApplicationError(
        'DOCUMENT_NOT_DRAFT',
        'Only draft documents can be published'
      );
    }

    // ✅ Command Service executes the API call
    // Backend will also validate (source of truth)
    return await this.commandService.publish(input.documentId);
  }
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
      const usecase = new PublishDocumentUsecase(queryService, commandService);
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
```

**Important**: Backend APIs enforce complete business rules. Frontend validation is for **UX only** (immediate feedback).

**Benefits**:

- **Clear separation**: Business logic in Usecase, API calls in Service
- **Simple Service**: Easy to test, no complex logic
- **Backend trust**: Rely on backend for authoritative validation

### Pattern C: Repository (Local State Management)

**Purpose**: Manage state that lives **only in the frontend** with full domain logic

**Use when**:

- Frontend-only state (drafts, temporary state)
- LocalStorage or SessionStorage persistence
- Complex workflows (multi-step forms, state machines)
- Offline support
- Undo/redo functionality

**Key characteristics**:

- Interface defined in **Domain layer**
- Manages **Domain Models** (not ViewData)
- Contains **business logic and validation**
- Typically implemented with LocalStorage, SessionStorage, or useReducer

**Example**:

```typescript
// domain/models/DraftDocument.ts (Frontend-specific Domain Model)
export class DraftDocument {
  private constructor(
    private readonly id: string,
    private readonly title: string,
    private readonly content: string,
    private readonly repositoryId: string,
    private readonly filePath: string,
    private readonly lastSavedAt: Date | null
  ) {}

  static create(data: {
    repositoryId: string;
    filePath: string;
  }): DraftDocument {
    return new DraftDocument(
      crypto.randomUUID(),
      '',
      '',
      data.repositoryId,
      data.filePath,
      null
    );
  }

  static reconstruct(data: {
    id: string;
    title: string;
    content: string;
    repositoryId: string;
    filePath: string;
    lastSavedAt: string | null;
  }): DraftDocument {
    return new DraftDocument(
      data.id,
      data.title,
      data.content,
      data.repositoryId,
      data.filePath,
      data.lastSavedAt ? new Date(data.lastSavedAt) : null
    );
  }

  getId(): string { return this.id; }
  getTitle(): string { return this.title; }
  getContent(): string { return this.content; }
  getRepositoryId(): string { return this.repositoryId; }
  getFilePath(): string { return this.filePath; }
  getLastSavedAt(): Date | null { return this.lastSavedAt; }

  // ✅ Domain logic: Business validation
  canSubmit(): boolean {
    return this.title.length > 0 && this.content.length > 0;
  }

  // ✅ Domain logic: Validation
  updateTitle(newTitle: string): DraftDocument {
    if (newTitle.length > 200) {
      throw new DomainError('TITLE_TOO_LONG', 'Title must be 200 characters or less');
    }
    return new DraftDocument(
      this.id,
      newTitle,
      this.content,
      this.repositoryId,
      this.filePath,
      this.lastSavedAt
    );
  }

  updateContent(newContent: string): DraftDocument {
    return new DraftDocument(
      this.id,
      this.title,
      newContent,
      this.repositoryId,
      this.filePath,
      this.lastSavedAt
    );
  }

  markAsSaved(): DraftDocument {
    return new DraftDocument(
      this.id,
      this.title,
      this.content,
      this.repositoryId,
      this.filePath,
      new Date()
    );
  }

  // Serialization for LocalStorage
  serialize(): string {
    return JSON.stringify({
      id: this.id,
      title: this.title,
      content: this.content,
      repositoryId: this.repositoryId,
      filePath: this.filePath,
      lastSavedAt: this.lastSavedAt?.toISOString() ?? null,
    });
  }

  static deserialize(json: string): DraftDocument {
    const data = JSON.parse(json);
    return DraftDocument.reconstruct(data);
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

// infrastructure/repositories/LocalStorageDraftRepository.ts
export class LocalStorageDraftRepository implements DraftDocumentRepository {
  private static STORAGE_KEY = 'draft_documents';

  get(id: string): DraftDocument | null {
    const stored = localStorage.getItem(LocalStorageDraftRepository.STORAGE_KEY);
    if (!stored) return null;

    try {
      const drafts: string[] = JSON.parse(stored);
      const draftJson = drafts.find(json => {
        const draft = DraftDocument.deserialize(json);
        return draft.getId() === id;
      });

      return draftJson ? DraftDocument.deserialize(draftJson) : null;
    } catch (err) {
      console.error('Failed to load draft from LocalStorage', err);
      return null;
    }
  }

  getAll(): DraftDocument[] {
    const stored = localStorage.getItem(LocalStorageDraftRepository.STORAGE_KEY);
    if (!stored) return [];

    try {
      const drafts: string[] = JSON.parse(stored);
      return drafts.map(json => DraftDocument.deserialize(json));
    } catch (err) {
      console.error('Failed to load drafts from LocalStorage', err);
      return [];
    }
  }

  add(draft: DraftDocument): void {
    const drafts = this.getAll();
    drafts.push(draft);
    this.persist(drafts);
  }

  update(draft: DraftDocument): void {
    const drafts = this.getAll();
    const index = drafts.findIndex(d => d.getId() === draft.getId());
    if (index === -1) {
      throw new Error('Draft not found');
    }
    drafts[index] = draft;
    this.persist(drafts);
  }

  remove(id: string): void {
    const drafts = this.getAll().filter(d => d.getId() !== id);
    this.persist(drafts);
  }

  private persist(drafts: DraftDocument[]): void {
    const serialized = drafts.map(d => d.serialize());
    localStorage.setItem(
      LocalStorageDraftRepository.STORAGE_KEY,
      JSON.stringify(serialized)
    );
  }
}

// application/usecases/SaveDraftUsecase.ts
export class SaveDraftUsecase {
  constructor(private readonly draftRepository: DraftDocumentRepository) {}

  execute(input: {
    draftId: string;
    title: string;
    content: string;
  }): DraftDocument {
    // Get from local state
    const draft = this.draftRepository.get(input.draftId);
    if (!draft) {
      throw new ApplicationError('DRAFT_NOT_FOUND', 'Draft not found');
    }

    // ✅ Domain logic: validation and transformation
    const updated = draft
      .updateTitle(input.title)
      .updateContent(input.content)
      .markAsSaved();

    // Update local state
    this.draftRepository.update(updated);

    return updated;
  }
}

// application/usecases/SubmitDraftUsecase.ts
export class SubmitDraftUsecase {
  constructor(
    private readonly draftRepository: DraftDocumentRepository,
    private readonly commandService: DocumentCommandService
  ) {}

  async execute(input: { draftId: string }): Promise<DocumentViewData> {
    // Get from local state
    const draft = this.draftRepository.get(input.draftId);
    if (!draft) {
      throw new ApplicationError('DRAFT_NOT_FOUND', 'Draft not found');
    }

    // ✅ Domain validation
    if (!draft.canSubmit()) {
      throw new ApplicationError(
        'DRAFT_INCOMPLETE',
        'Draft must have title and content'
      );
    }

    // ✅ Send to external API (Command Service)
    const created = await this.commandService.create({
      repositoryId: draft.getRepositoryId(),
      filePath: draft.getFilePath(),
      title: draft.getTitle(),
      content: draft.getContent(),
      accessScope: 'private',
    });

    // ✅ Clean up local state
    this.draftRepository.remove(draft.getId());

    return created;
  }
}

// presentation/hooks/useDraftDocument.ts
export function useDraftDocument(draftId: string) {
  const draftRepository = useDraftRepository();
  const [draft, setDraft] = useState<DraftDocument | null>(null);

  useEffect(() => {
    const loaded = draftRepository.get(draftId);
    setDraft(loaded);
  }, [draftId, draftRepository]);

  const saveDraft = (title: string, content: string) => {
    if (!draft) return;

    const usecase = new SaveDraftUsecase(draftRepository);
    const updated = usecase.execute({ draftId: draft.getId(), title, content });
    setDraft(updated);
  };

  return { draft, saveDraft };
}
```

**Use cases for Repository**:

- Draft documents with auto-save to LocalStorage
- Multi-step form wizards (temporary state between steps)
- Offline-first features (sync when online)
- Undo/redo functionality
- Complex local state with business rules

### Pattern Selection Guidelines

**Decision flow**:

```text
1. Is the data managed by backend?
   YES → Use Query Service (GET) or Command Service (POST/PUT/DELETE)
   NO ↓

2. Does the operation need complex frontend business logic?
   YES → Repository with Domain Models
   NO ↓

3. Is it just displaying data?
   YES → Query Service (possibly skip Usecase)
   NO → Usecase + Command Service
```

**Comparison table**:

| Pattern | Use When | Domain Logic | Data Location | Interface Location |
| ------- | -------- | ------------ | ------------- | ------------------ |
| **Query Service** | Fetching external data for display | ❌ No | Backend API | Application layer |
| **Command Service** | Sending changes to external API | ⚠️ Minimal (lightweight checks) | Backend API | Application layer |
| **Repository** | Managing frontend-only state | ✅ Yes (full validation) | LocalStorage / useReducer | Domain layer |

**Examples**:

- **Document list display** → Query Service (simple fetch)
- **Publish document** → Usecase + Command Service (business validation + API call)
- **Draft document editor** → Repository (local state with business logic)
- **Auto-save draft** → Repository (periodic local persistence)

## Consequences

### Positive

1. **Clear responsibilities**: Each pattern has a distinct purpose
2. **Separation of concerns**: External API vs local state
3. **Testability**: Can test each pattern independently
4. **Flexibility**: Can swap implementations (HTTP → WebSocket → LocalStorage)
5. **Maintainability**: Easy to understand where logic lives
6. **React Query friendly**: Query Service works seamlessly with caching libraries

### Negative

1. **Three patterns to learn**: More concepts than traditional Repository
2. **Judgment calls**: Deciding which pattern to use
3. **More interfaces**: Query + Command + Repository instead of just Repository
4. **Potential confusion**: "Why not just use Repository for everything?"

### Mitigation

1. **Clear guidelines**: This ADR provides decision tree
2. **Code examples**: Complete examples for each pattern
3. **Code reviews**: Validate pattern choice
4. **Start simple**: Use Query/Command Services first, add Repository when needed

## Examples

### Example 1: Simple Query (No Usecase)

```typescript
// Presentation directly uses Query Service for simple display
function DocumentList() {
  const queryService = useDocumentQueryService();

  const { data, isLoading } = useQuery({
    queryKey: ['documents'],
    queryFn: () => queryService.list({}),
  });

  // Simple display, no business logic needed
  return <ul>{data?.map(doc => <li key={doc.id}>{doc.title}</li>)}</ul>;
}
```

### Example 2: Command with Business Logic (Usecase)

```typescript
// Complex operation requires Usecase
class DeleteDocumentUsecase {
  constructor(
    private readonly queryService: DocumentQueryService,
    private readonly commandService: DocumentCommandService
  ) {}

  async execute(input: { documentId: string; userId: string }): Promise<void> {
    // Business validation
    const document = await this.queryService.getById(input.documentId);

    if (document.ownerId !== input.userId) {
      throw new ApplicationError('UNAUTHORIZED', 'Only owner can delete');
    }

    if (document.status === 'published') {
      throw new ApplicationError('CANNOT_DELETE_PUBLISHED', 'Cannot delete published document');
    }

    // Execute deletion
    await this.commandService.delete(input.documentId);
  }
}
```

### Example 3: Local State with Repository

```typescript
// Complex local state management
function DraftEditor({ draftId }: { draftId: string }) {
  const { draft, saveDraft } = useDraftDocument(draftId);
  const { submit, isSubmitting } = useSubmitDraft();

  // Auto-save every 5 seconds
  useEffect(() => {
    if (!draft) return;

    const timer = setInterval(() => {
      saveDraft(draft.getTitle(), draft.getContent());
    }, 5000);

    return () => clearInterval(timer);
  }, [draft, saveDraft]);

  const handleSubmit = async () => {
    if (!draft) return;

    try {
      await submit(draft.getId());
      // Redirect to document list
    } catch (err) {
      // Show error
    }
  };

  if (!draft) return <div>Draft not found</div>;

  return (
    <div>
      <input
        value={draft.getTitle()}
        onChange={e => saveDraft(e.target.value, draft.getContent())}
      />
      <textarea
        value={draft.getContent()}
        onChange={e => saveDraft(draft.getTitle(), e.target.value)}
      />
      <button
        onClick={handleSubmit}
        disabled={!draft.canSubmit() || isSubmitting}
      >
        Submit
      </button>
      <p>Last saved: {draft.getLastSavedAt()?.toLocaleString() ?? 'Never'}</p>
    </div>
  );
}
```

## Related ADRs

- **ADR 0018**: Frontend Layered Architecture - Layer definitions
- **ADR 0022**: types/ Directory - ViewData and API response types
- **ADR 0023**: React-Specific Patterns - Immutable models, Context-based DI

## References

- [CQRS Pattern](https://martinfowler.com/bliki/CQRS.html) - Command Query Responsibility Segregation
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html) - Traditional repository
- [React Query](https://tanstack.com/query/latest) - Query Service integration
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) - Ports and Adapters
