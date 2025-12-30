# ADR 0020: Frontend Validation Strategy

## Status

Accepted

## Date

2025-12-30

## Context

The frontend application needs to validate user input at multiple levels:

1. **Format validation**: Is the input in the correct format? (email, URL, number range)
2. **Business validation**: Does the input satisfy business rules? (uniqueness, state transitions, permissions)

**Challenge**: Where should each type of validation occur?

Without clear separation:

- Business logic leaks into UI components
- Validation is duplicated across components
- Testing becomes difficult (need to render components to test business rules)
- Backend and frontend validation drift apart

We need a **two-layer validation strategy** that separates concerns while maintaining a good user experience.

## Decision

We adopt a **two-layer validation strategy**:

1. **Input Validation** (Presentation layer): Format, required, length
2. **Business Validation** (Domain/Application layers): Business rules, constraints, permissions

### Layer 1: Input Validation (Presentation Layer)

**Purpose**: Ensure user input meets format requirements for immediate feedback

**Implementation**: `react-hook-form` + `zod`

**Responsibility**: Presentation layer

**Validates**:

- Required fields
- Data types (string, number, boolean, date)
- Format patterns (email, URL, regex)
- Length constraints (min/max characters)
- Enum values
- Basic rules (min/max numbers)

**Timing**: Real-time (on blur, on change) for immediate user feedback

**Example**:

```typescript
// Presentation layer: Input validation schema
import { z } from 'zod';

const documentFormSchema = z.object({
  title: z
    .string()
    .min(1, 'Title is required')
    .max(200, 'Title must be 200 characters or less'),

  filePath: z
    .string()
    .regex(/\.md$/, 'Must be a Markdown file (.md)'),

  repositoryId: z
    .string()
    .uuid('Invalid repository ID'),

  accessScope: z
    .enum(['public', 'private'], {
      errorMap: () => ({ message: 'Must be public or private' }),
    }),

  isAutoUpdate: z
    .boolean(),

  tags: z
    .array(z.string())
    .max(10, 'Maximum 10 tags allowed')
    .optional(),

  content: z
    .string()
    .max(50000, 'Content too large')
    .optional(),
});

type DocumentFormData = z.infer<typeof documentFormSchema>;
```

**React Hook Form integration**:

```typescript
// presentation/components/DocumentRegistrationDialog.tsx
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';

export function DocumentRegistrationDialog() {
  const form = useForm<DocumentFormData>({
    resolver: zodResolver(documentFormSchema),
    mode: 'onBlur', // Validate on blur for better UX
  });

  const { registerDocument, isLoading, error } = useDocumentRegistration();

  const onSubmit = async (data: DocumentFormData) => {
    try {
      // ✅ Input validation already passed
      await registerDocument(data);
    } catch (err) {
      // ❌ Business validation errors from backend
      if (err instanceof ApplicationError) {
        form.setError('root', { message: err.message });
      }
    }
  };

  return (
    <form onSubmit={form.handleSubmit(onSubmit)}>
      <input
        {...form.register('title')}
        placeholder="Document title"
      />
      {form.formState.errors.title && (
        <span className="error">{form.formState.errors.title.message}</span>
      )}

      <input
        {...form.register('filePath')}
        placeholder="path/to/document.md"
      />
      {form.formState.errors.filePath && (
        <span className="error">{form.formState.errors.filePath.message}</span>
      )}

      <button type="submit" disabled={isLoading}>
        Register
      </button>

      {error && <div className="error">{error}</div>}
    </form>
  );
}
```

**Benefits**:

- **Immediate feedback**: Users see errors as they type
- **Type-safe**: TypeScript enforces schema
- **Reusable**: Schema can be shared across forms
- **Standard library**: Using established tools

### Layer 2: Business Validation (Domain/Application Layers)

**Purpose**: Enforce business rules and domain invariants

**Implementation**: Domain entities and Application usecases

**Responsibility**: Domain/Application layers

**Validates**:

- Uniqueness constraints (e.g., duplicate document names)
- State transition rules (e.g., can only publish draft documents)
- Authorization (e.g., only owner can delete)
- Cross-aggregate rules (e.g., repository must exist)
- Domain invariants (e.g., published document cannot be empty)
- Complex business logic

**Timing**: On usecase execution (async, may require backend validation)

**Example - Domain Layer**:

```typescript
// domain/models/Document.ts
export class DocumentEntity {
  private constructor(
    private readonly id: string,
    private readonly title: string,
    private readonly content: string,
    private readonly status: 'draft' | 'published',
    private readonly ownerId: string
  ) {}

  static create(data: {
    title: string;
    content: string;
    ownerId: string;
  }): DocumentEntity {
    // ✅ Business validation: domain invariants
    if (data.title.length === 0) {
      throw new DomainError('EMPTY_TITLE', 'Title cannot be empty');
    }

    if (data.title.length > 200) {
      throw new DomainError('TITLE_TOO_LONG', 'Title too long');
    }

    return new DocumentEntity(
      crypto.randomUUID(),
      data.title,
      data.content,
      'draft',
      data.ownerId
    );
  }

  // ✅ Business logic: state transition rules
  publish(userId: string): DocumentEntity {
    if (this.status !== 'draft') {
      throw new DomainError(
        'DOCUMENT_NOT_PUBLISHABLE',
        'Only draft documents can be published'
      );
    }

    if (this.ownerId !== userId) {
      throw new DomainError(
        'UNAUTHORIZED',
        'Only the owner can publish this document'
      );
    }

    if (this.content.length === 0) {
      throw new DomainError(
        'DOCUMENT_EMPTY',
        'Cannot publish an empty document'
      );
    }

    return new DocumentEntity(
      this.id,
      this.title,
      this.content,
      'published',
      this.ownerId
    );
  }

  canPublish(userId: string): boolean {
    return (
      this.status === 'draft' &&
      this.ownerId === userId &&
      this.content.length > 0
    );
  }
}
```

**Example - Application Layer**:

```typescript
// application/usecases/RegisterDocumentUsecase.ts
export class RegisterDocumentUsecase {
  constructor(
    private readonly documentRepository: DocumentRepository,
    private readonly repositoryRepository: RepositoryRepository
  ) {}

  async execute(input: {
    repositoryId: string;
    filePath: string;
    title: string;
    content: string;
    accessScope: 'public' | 'private';
    currentUserId: string;
  }): Promise<{ documentId: string }> {
    // ✅ Business validation: repository exists
    const repository = await this.repositoryRepository.findById(input.repositoryId);
    if (!repository) {
      throw new ApplicationError(
        'REPOSITORY_NOT_FOUND',
        'Repository does not exist'
      );
    }

    // ✅ Business validation: file type (business rule, not just format)
    if (!input.filePath.endsWith('.md')) {
      throw new ApplicationError(
        'INVALID_FILE_TYPE',
        'Only Markdown files can be registered as documents'
      );
    }

    // ✅ Business validation: duplicate check
    const existing = await this.documentRepository.findByRepositoryAndPath(
      input.repositoryId,
      input.filePath
    );
    if (existing) {
      throw new ApplicationError(
        'DOCUMENT_ALREADY_EXISTS',
        'A document already exists at this path'
      );
    }

    // ✅ Create domain entity (domain validation happens here)
    const document = DocumentEntity.create({
      title: input.title,
      content: input.content,
      ownerId: input.currentUserId,
    });

    // ✅ Persist
    const saved = await this.documentRepository.save(document);

    return { documentId: saved.getId() };
  }
}
```

**Presentation layer handling**:

```typescript
// presentation/hooks/useDocumentRegistration.ts
export function useDocumentRegistration() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const documentRepository = useDocumentRepository();
  const repositoryRepository = useRepositoryRepository();
  const { currentUser } = useAuth();

  const registerDocument = async (input: DocumentFormData) => {
    setIsLoading(true);
    setError(null);

    try {
      const usecase = new RegisterDocumentUsecase(
        documentRepository,
        repositoryRepository
      );

      const result = await usecase.execute({
        ...input,
        currentUserId: currentUser!.id,
      });

      return result;
    } catch (err) {
      // ✅ Handle business validation errors
      if (err instanceof ApplicationError || err instanceof DomainError) {
        setError(err.message);
      } else {
        setError('An unexpected error occurred');
      }
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  return { registerDocument, isLoading, error };
}
```

### Validation Flow

```text
User Input
  │
  ├─▶ [Input Validation] (Presentation - react-hook-form + zod)
  │     │
  │     ├─▶ ❌ Format error → Show immediately
  │     └─▶ ✅ Format valid
  │           │
  ├─────────▶ [Submit to Usecase]
  │           │
  ├─────────▶ [Business Validation] (Application/Domain)
  │           │
  │           ├─▶ ❌ Business error → Return to UI
  │           └─▶ ✅ Business valid
  │                 │
  └───────────────▶ [Execute business logic]
```

### Validation Responsibility Matrix

| Validation Type | Example | Layer | Tool | Timing |
| --- | --- | --- | --- | --- |
| **Required fields** | Title required | Presentation | zod | On blur |
| **Format** | Email format | Presentation | zod | On blur |
| **Length** | Max 200 chars | Presentation | zod | On change |
| **Enum values** | Status: draft/published | Presentation | zod | On change |
| **Uniqueness** | Document path not taken | Application | Usecase | On submit |
| **State transitions** | Only draft → published | Domain | Entity | On submit |
| **Authorization** | Only owner can delete | Domain | Entity | On submit |
| **Domain invariants** | Published doc not empty | Domain | Entity | On submit |
| **Cross-aggregate** | Repository exists | Application | Usecase | On submit |

### Error Handling Pattern

**Error types**:

```typescript
// domain/errors/DomainErrors.ts
export class DomainError extends Error {
  constructor(
    public readonly code: string,
    message: string
  ) {
    super(message);
    this.name = 'DomainError';
  }
}

// application/errors/ApplicationErrors.ts
export class ApplicationError extends Error {
  constructor(
    public readonly code: string,
    message: string
  ) {
    super(message);
    this.name = 'ApplicationError';
  }
}
```

**Presentation error handling**:

```typescript
// presentation/components/DocumentForm.tsx
const onSubmit = async (data: DocumentFormData) => {
  try {
    await registerDocument(data);
    // Success handling
  } catch (err) {
    if (err instanceof DomainError) {
      // Handle domain errors (e.g., show modal)
      showErrorDialog(err.message);
    } else if (err instanceof ApplicationError) {
      // Handle application errors (e.g., show toast)
      showToast(err.message, 'error');
    } else {
      // Handle unexpected errors
      showToast('An unexpected error occurred', 'error');
    }
  }
};
```

### Backend Validation

**Important**: Backend APIs should **always** enforce complete business rules. Frontend validation is for **UX only**.

- Frontend: **Lightweight checks** for immediate feedback
- Backend: **Source of truth** for all business logic

**Example**:

```typescript
// Frontend: Check status before calling API (UX optimization)
const { publish } = useDocumentPublish();

const handlePublish = async () => {
  // ✅ Quick client-side check for better UX
  if (!document.canPublish(currentUser.id)) {
    showToast('Document cannot be published', 'error');
    return;
  }

  try {
    // ✅ Backend will also validate (source of truth)
    await publish(document.getId());
  } catch (err) {
    // Backend rejected (e.g., someone else published it)
    showToast(err.message, 'error');
  }
};
```

## Consequences

### Positive

1. **Clear separation**: Format vs business validation
2. **Immediate feedback**: Users see format errors instantly
3. **Testable**: Business logic can be tested without UI
4. **Reusable**: Validation schemas can be shared
5. **Type-safe**: zod + TypeScript prevent errors
6. **Backend alignment**: Frontend mirrors backend validation

### Negative

1. **Duplication**: Some validation exists in both layers
2. **Complexity**: Two validation systems to maintain
3. **Learning curve**: Developers need to understand both
4. **Potential mismatch**: Frontend and backend validation can drift

### Mitigation

1. **Share schemas**: Generate zod schemas from backend OpenAPI specs (future)
2. **Clear documentation**: This ADR clarifies responsibilities
3. **Code reviews**: Ensure validation is in correct layer
4. **Testing**: Test both input and business validation separately
5. **Error codes**: Use consistent error codes across frontend/backend

## Examples

### Example 1: Complete Validation Flow

```typescript
// 1. Input validation (Presentation)
const userFormSchema = z.object({
  email: z.string().email('Invalid email format'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
  age: z.number().min(0).max(150),
});

// 2. Form component
function UserRegistrationForm() {
  const form = useForm({
    resolver: zodResolver(userFormSchema),
  });

  const { register: registerUser } = useUserRegistration();

  const onSubmit = async (data: UserFormData) => {
    try {
      // ✅ Input validation passed
      await registerUser(data);
    } catch (err) {
      // ❌ Business validation failed
      if (err instanceof ApplicationError) {
        if (err.code === 'EMAIL_ALREADY_EXISTS') {
          form.setError('email', { message: 'Email already registered' });
        }
      }
    }
  };

  return <form onSubmit={form.handleSubmit(onSubmit)}>...</form>;
}

// 3. Business validation (Application)
class RegisterUserUsecase {
  async execute(input: RegisterUserInput): Promise<User> {
    // ✅ Business validation: email uniqueness
    const existing = await this.userRepository.findByEmail(input.email);
    if (existing) {
      throw new ApplicationError(
        'EMAIL_ALREADY_EXISTS',
        'Email is already registered'
      );
    }

    // ✅ Business validation: age restriction
    if (input.age < 18) {
      throw new ApplicationError(
        'UNDERAGE',
        'Must be 18 or older to register'
      );
    }

    // Create user...
  }
}
```

### Example 2: Conditional Validation

```typescript
// Input validation with conditional rules
const documentFormSchema = z.object({
  type: z.enum(['procedure', 'knowledge']),
  title: z.string().min(1).max(200),

  // ✅ Conditional validation based on type
  variables: z.array(z.object({
    name: z.string(),
    type: z.string(),
  })).optional(),

  executionTime: z.number().optional(),
}).refine(
  (data) => {
    // If type is 'procedure', variables are required
    if (data.type === 'procedure' && !data.variables) {
      return false;
    }
    return true;
  },
  {
    message: 'Procedure documents must define variables',
    path: ['variables'],
  }
);
```

### Example 3: Async Business Validation

```typescript
// Business validation that requires async checks
class CreateDocumentUsecase {
  async execute(input: CreateDocumentInput): Promise<Document> {
    // ✅ Async business validation
    const [repository, existingDocument, userPermissions] = await Promise.all([
      this.repositoryRepository.findById(input.repositoryId),
      this.documentRepository.findByPath(input.path),
      this.permissionService.getUserPermissions(input.userId),
    ]);

    if (!repository) {
      throw new ApplicationError('REPOSITORY_NOT_FOUND');
    }

    if (existingDocument) {
      throw new ApplicationError('DOCUMENT_ALREADY_EXISTS');
    }

    if (!userPermissions.includes('document:create')) {
      throw new ApplicationError('INSUFFICIENT_PERMISSIONS');
    }

    // Create document...
  }
}
```

## Related ADRs

- **ADR 0017**: Application Data Validation (backend) - Backend validation approach
- **ADR 0018**: Frontend Layered Architecture - Where validation lives in each layer
- **ADR 0023**: React-Specific Patterns - Immutable models with validation

## References

- [react-hook-form](https://react-hook-form.com/) - Form validation library
- [zod](https://zod.dev/) - TypeScript schema validation
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html) - Domain validation
- [Client-side vs Server-side Validation](https://owasp.org/www-project-proactive-controls/v3/en/c5-validate-inputs) - Security perspective
