# ADR 0016: Document Domain Model Design

## Status

Accepted

## Context

OpsCore needs a robust document management system to handle operational procedure documents and knowledge documents. The system must support version control, variable definitions for parameterized procedures, and access control.

Key requirements:
1. Track document versions with their source location (file path and commit hash)
2. Store document metadata (title, type, tags, variables) per version since these can change when the file is updated
3. Support public and private access with granular sharing for private documents
4. Enable automatic updates when the source repository changes
5. Allow rollback to previous versions

## Decision

We will implement the Document aggregate with the following design:

### 1. Document Aggregate Structure

The Document aggregate consists of:
- **Document** (aggregate root): Manages the document lifecycle and versions
- **DocumentVersion** (entity): Represents a specific version snapshot

### 2. Field Placement

**Document (Aggregate Root)** contains repository-level metadata:
- `id`: DocumentID - Unique identifier for the document
- `repositoryID`: RepositoryID - The source repository
- `owner`: string - Document owner (user who created/manages it)
- `isPublished`: boolean - Publication status
- `isAutoUpdate`: boolean - Whether to auto-update on repository changes
- `accessScope`: AccessScope - Access control (public/private)
- `currentVersion`: DocumentVersion - Reference to current version
- `versions`: []DocumentVersion - Version history
- `createdAt`, `updatedAt`: DateTime - Lifecycle timestamps

**DocumentVersion** contains file-specific metadata:
- `id`: VersionID - Unique version identifier
- `documentID`: DocumentID - Parent document reference
- `versionNumber`: VersionNumber - Sequential number (1, 2, 3...)
- `source`: DocumentSource - File path and commit hash combination
- `title`: string - Document title from frontmatter
- `docType`: DocumentType - procedure or knowledge
- `tags`: []Tag - Tags from frontmatter
- `variables`: []VariableDefinition - Variable definitions from frontmatter
- `content`: string - Markdown content
- `publishedAt`, `unpublishedAt`: DateTime - Publication timestamps
- `isCurrentVersion`: boolean - Whether this is the current version

**Rationale for field placement:**
- File path, title, type, tags, and variables are stored in DocumentVersion because:
  1. These fields are specified in each Markdown file's frontmatter (per ADR 0013)
  2. Files can be renamed/moved in the repository (file path changes)
  3. Frontmatter metadata can change between versions
  4. This enables complete version history with all metadata changes
- Owner and access control remain at Document level because they're management concerns, not file content

### 3. Value Objects

**DocumentSource**: Combines FilePath and CommitHash into a single value object
- Rationale: These two fields together uniquely identify a document's source location
- Provides a cohesive representation of "where this version came from"
- Format: `{filePath}@{commitHash}` (e.g., `docs/backup.md@abc1234`)

**AccessScope**: Simplified to two values
- `public`: Accessible to all users
- `private`: Owner-only by default, can be shared with specific users/groups
- Rationale: Simplifies the model while maintaining flexibility through a separate sharing mechanism
- Sharing will be implemented as a separate concern (e.g., DocumentShare entity)

**Removed Value Objects:**
- **Category**: Removed because it overlaps with DocumentType and Tags
  - DocumentType provides broad categorization (procedure vs knowledge)
  - Tags provide fine-grained categorization
  - No ADR specified Category, so it was an unnecessary abstraction
  
- **VariableValue**: Removed from Document aggregate
  - Per ADR 0013 and ADR 0014, variable values belong to ExecutionRecord aggregate
  - Variable values are execution-time data, not document metadata
  - Will be implemented in ExecutionRecord aggregate (issue #32)

### 4. Repository Interface

```go
type DocumentRepository interface {
    Save(ctx context.Context, document Document) error
    FindByID(ctx context.Context, id DocumentID) (Document, error)
    FindByRepositoryID(ctx context.Context, repoID RepositoryID) ([]Document, error)
    FindPublished(ctx context.Context, filters ...Filter) ([]Document, error)
    Update(ctx context.Context, document Document) error
    Delete(ctx context.Context, id DocumentID) error
    
    // Version management
    SaveVersion(ctx context.Context, version DocumentVersion) error
    FindVersionsByDocumentID(ctx context.Context, docID DocumentID) ([]DocumentVersion, error)
    FindVersionByNumber(ctx context.Context, docID DocumentID, versionNumber VersionNumber) (DocumentVersion, error)
}
```

### 5. Domain Behaviors

**Document:**
- `Publish(source, title, docType, tags, variables, content)`: Create and publish a new version
- `Unpublish()`: Unpublish the current version
- `UpdateAccessScope(scope)`: Change access scope
- `EnableAutoUpdate()` / `DisableAutoUpdate()`: Control automatic updates
- `RollbackToVersion(versionNumber)`: Revert to a previous version
- `AddVersion(version)`: Add a version (for reconstruction from persistence)

**DocumentVersion:**
- `MarkAsCurrent()`: Mark this version as current
- `Unpublish()`: Unpublish this version
- `IsPublished()`: Check publication status

## Consequences

### Pros

- **Complete version history**: All metadata changes are tracked in versions
- **Flexibility**: Documents can evolve (rename, move, change type) while maintaining history
- **Simplified access control**: Two-tier approach (public/private + sharing) is easier to understand
- **Source tracking**: DocumentSource provides clear lineage from repository
- **Aligned with ADRs**: Follows specifications from ADR 0013 and ADR 0014

### Cons

- **More complex version entity**: DocumentVersion has more fields than before
- **Metadata duplication**: Title, type, etc. stored per version instead of once per document
- **Migration needed**: Existing code expecting metadata at Document level needs updates

### Migration Path

1. Update infrastructure layer to persist new structure
2. Create migration to move existing metadata to versions
3. Update application layer to work with new Publish signature
4. Update API handlers to accept new parameters

## Related ADRs

- ADR 0007: Backend Architecture - Onion Architecture
- ADR 0013: Document Variable Definition and Substitution
- ADR 0014: Execution Record and Evidence Management

## References

- Issue #31: Document集約の実装
- Issue #30: ドメインモデル再設計の実装（親Issue）
