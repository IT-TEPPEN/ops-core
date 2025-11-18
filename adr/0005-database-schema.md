# ADR 0005: Database Schema for Repository Configurations and Document Management

## Status

Accepted (Updated: 2025-11-17)

## Context

ADR 0004 established PostgreSQL as the chosen database for OpsCore. This ADR originally defined the initial database schema for external repository configurations. 

**Domain Model Redesign (2025-11)**

Following the domain model redesign documented in `domain-model-redesign-summary.md`, this ADR has been updated to include:
- Corrections to align with actual implementation
- New tables for document management, execution records, and user/group management
- JSONB field schemas for flexible data structures
- Comprehensive indexing strategy

Related ADRs:
- ADR 0013: Document Variable Definition and Substitution
- ADR 0014: Execution Record and Evidence Management

## Decision

### Existing Tables (Corrected from Initial Implementation)

#### `repositories` Table

Stores details about configured external Git repositories.

**Discrepancies from Original ADR:**
- Table name: `repositories` (not `repository_configurations`) - shortened for better usability
- Column name: `url` (not `repository_url`) - simplified naming convention
- Column name: `access_token` (not `credentials`) - more specific and accurate
- No `provider_type` column - provider is derived from URL, avoiding redundancy
- URL length: VARCHAR(255) (not 2048) - sufficient for practical Git repository URLs

**Schema:**
```sql
CREATE TABLE repositories (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL UNIQUE,
    access_token TEXT, -- Encrypted using AES-256-GCM at application level
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Rationale for Changes:**
- **Table name:** Shorter, more conventional name following PostgreSQL best practices
- **URL length:** 255 characters is sufficient for GitHub/GitLab URLs and improves index performance
- **access_token:** Replaces `credentials`; encrypted at application level using AES-256-GCM
- **No provider_type:** Git provider can be determined from URL pattern (github.com, gitlab.com, etc.)

#### `managed_files` Table

Tracks files within repositories that are managed by OpsCore.

**Schema:**
```sql
CREATE TABLE managed_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repository_id UUID NOT NULL,
    file_path TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE,
    UNIQUE (repository_id, file_path)
);
```

### New Tables (Domain Model Redesign)

#### `documents` Table

Stores document metadata for published operational procedures and knowledge articles.

**Schema:**
```sql
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repository_id UUID NOT NULL,
    file_path TEXT NOT NULL,
    title VARCHAR(255) NOT NULL,
    owner VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('procedure', 'knowledge')),
    tags TEXT[] NOT NULL DEFAULT '{}',
    category VARCHAR(255),
    variables JSONB, -- VariableDefinition[] schema
    is_published BOOLEAN NOT NULL DEFAULT false,
    is_auto_update BOOLEAN NOT NULL DEFAULT false,
    access_scope JSONB NOT NULL, -- AccessScope schema
    current_version_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE,
    FOREIGN KEY (current_version_id) REFERENCES document_versions(id),
    UNIQUE (repository_id, file_path)
);

CREATE INDEX idx_documents_repository_id ON documents(repository_id);
CREATE INDEX idx_documents_type ON documents(type);
CREATE INDEX idx_documents_is_published ON documents(is_published);
CREATE INDEX idx_documents_tags ON documents USING GIN(tags);
CREATE INDEX idx_documents_category ON documents(category);
```

#### `document_versions` Table

Stores version history for documents, linking commit hashes to sequential version numbers.

**Schema:**
```sql
CREATE TABLE document_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL,
    version_number INTEGER NOT NULL,
    commit_hash VARCHAR(64) NOT NULL,
    content TEXT NOT NULL,
    published_at TIMESTAMPTZ NOT NULL,
    unpublished_at TIMESTAMPTZ,
    is_current_version BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE,
    UNIQUE (document_id, version_number),
    UNIQUE (document_id, commit_hash)
);

CREATE INDEX idx_document_versions_document_id ON document_versions(document_id);
CREATE INDEX idx_document_versions_is_current ON document_versions(is_current_version);
```

#### `users` Table

Stores user information for authentication and authorization.

**Schema:**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'user')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
```

#### `groups` Table

Stores user groups for organizing users and managing access control.

**Schema:**
```sql
CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_groups_name ON groups(name);
```

#### `user_groups` Table

Junction table for many-to-many relationship between users and groups.

**Schema:**
```sql
CREATE TABLE user_groups (
    user_id UUID NOT NULL,
    group_id UUID NOT NULL,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_groups_user_id ON user_groups(user_id);
CREATE INDEX idx_user_groups_group_id ON user_groups(group_id);
```

#### `execution_records` Table

Stores execution records for procedure executions with evidence tracking.

**Schema:**
```sql
CREATE TABLE execution_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL,
    document_version_id UUID NOT NULL,
    executor_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    variable_values JSONB, -- VariableValue[] schema
    notes TEXT,
    status VARCHAR(50) NOT NULL CHECK (status IN ('in_progress', 'completed', 'failed')),
    access_scope JSONB NOT NULL, -- AccessScope schema
    started_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE,
    FOREIGN KEY (document_version_id) REFERENCES document_versions(id),
    FOREIGN KEY (executor_id) REFERENCES users(id)
);

CREATE INDEX idx_execution_records_document_id ON execution_records(document_id);
CREATE INDEX idx_execution_records_executor_id ON execution_records(executor_id);
CREATE INDEX idx_execution_records_status ON execution_records(status);
CREATE INDEX idx_execution_records_started_at ON execution_records(started_at);
CREATE INDEX idx_execution_records_variable_values ON execution_records USING GIN(variable_values);
```

#### `execution_steps` Table

Stores individual execution steps within an execution record.

**Schema:**
```sql
CREATE TABLE execution_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    execution_record_id UUID NOT NULL,
    step_number INTEGER NOT NULL,
    description TEXT NOT NULL,
    notes TEXT,
    executed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (execution_record_id) REFERENCES execution_records(id) ON DELETE CASCADE,
    UNIQUE (execution_record_id, step_number)
);

CREATE INDEX idx_execution_steps_execution_record_id ON execution_steps(execution_record_id);
```

#### `attachments` Table

Stores metadata for files attached to execution steps (screenshots, evidence files).

**Schema:**
```sql
CREATE TABLE attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    execution_record_id UUID NOT NULL,
    execution_step_id UUID NOT NULL,
    step_number INTEGER NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(127) NOT NULL,
    storage_type VARCHAR(50) NOT NULL CHECK (storage_type IN ('local', 's3', 'minio')),
    storage_path TEXT NOT NULL,
    uploaded_by UUID NOT NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (execution_record_id) REFERENCES execution_records(id) ON DELETE CASCADE,
    FOREIGN KEY (execution_step_id) REFERENCES execution_steps(id) ON DELETE CASCADE,
    FOREIGN KEY (uploaded_by) REFERENCES users(id)
);

CREATE INDEX idx_attachments_execution_record_id ON attachments(execution_record_id);
CREATE INDEX idx_attachments_execution_step_id ON attachments(execution_step_id);
CREATE INDEX idx_attachments_uploaded_by ON attachments(uploaded_by);
```

#### `view_history` Table

Records document view history for analytics and user convenience.

**Schema:**
```sql
CREATE TABLE view_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL,
    user_id UUID, -- NULL for anonymous users
    ip_address INET,
    user_agent TEXT,
    viewed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_view_history_document_id ON view_history(document_id);
CREATE INDEX idx_view_history_user_id ON view_history(user_id);
CREATE INDEX idx_view_history_viewed_at ON view_history(viewed_at);
```

#### `view_statistics` Table

Aggregates view statistics for documents.

**Schema:**
```sql
CREATE TABLE view_statistics (
    document_id UUID PRIMARY KEY,
    total_views BIGINT NOT NULL DEFAULT 0,
    unique_users INTEGER NOT NULL DEFAULT 0,
    last_viewed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE
);
```

### JSONB Field Schemas

#### `variables` Field (documents.variables)

Stores array of VariableDefinition objects for parameterized procedures.

**JSON Schema:**
```json
[
  {
    "name": "string",           // Variable identifier (alphanumeric + underscore)
    "label": "string",          // Human-readable label for UI
    "description": "string",    // Optional detailed explanation
    "type": "string",           // "string" | "number" | "boolean" | "date"
    "required": boolean,        // Whether variable is mandatory
    "defaultValue": any         // Optional default value matching type
  }
]
```

**Example:**
```json
[
  {
    "name": "server_name",
    "label": "サーバー名",
    "description": "バックアップ対象のサーバー名",
    "type": "string",
    "required": true,
    "defaultValue": "prod-db-01"
  },
  {
    "name": "retention_days",
    "label": "保持期間（日数）",
    "type": "number",
    "required": false,
    "defaultValue": 30
  }
]
```

#### `access_scope` Field (documents.access_scope, execution_records.access_scope)

Defines access control for documents and execution records.

**JSON Schema:**
```json
{
  "type": "string",              // "private" | "shared"
  "sharedWith": [                // Array of user/group IDs (when type="shared")
    {
      "id": "uuid",
      "type": "string"           // "user" | "group"
    }
  ]
}
```

**Examples:**
```json
// Private (only creator and admins)
{
  "type": "private"
}

// Shared with specific users and groups
{
  "type": "shared",
  "sharedWith": [
    {"id": "user-uuid-1", "type": "user"},
    {"id": "group-uuid-1", "type": "group"}
  ]
}
```

#### `variable_values` Field (execution_records.variable_values)

Stores the actual values used during procedure execution.

**JSON Schema:**
```json
[
  {
    "name": "string",          // Variable name matching definition
    "value": any               // Actual value used
  }
]
```

**Example:**
```json
[
  {
    "name": "server_name",
    "value": "prod-db-01"
  },
  {
    "name": "retention_days",
    "value": 30
  },
  {
    "name": "enable_compression",
    "value": true
  }
]
```

### Indexing Strategy

#### Standard B-tree Indexes

Used for:
- Primary keys (automatic)
- Foreign keys (for join performance)
- Frequently queried columns (type, status, email, etc.)
- Timestamp columns for range queries

**Rationale:** B-tree indexes are optimal for equality and range queries on scalar values.

#### GIN Indexes (Generalized Inverted Index)

Used for:
- **Array columns:** `documents.tags` - enables efficient tag filtering with `@>` operator
- **JSONB columns:** `documents.variables`, `execution_records.variable_values` - enables queries on JSON structure

**Example GIN Index Usage:**
```sql
-- Find documents with specific tag
SELECT * FROM documents WHERE tags @> ARRAY['backup'];

-- Find execution records with specific variable value
SELECT * FROM execution_records 
WHERE variable_values @> '[{"name": "server_name", "value": "prod-db-01"}]';
```

**Rationale:** GIN indexes are optimized for composite values where elements can be searched independently.

#### Index Maintenance Considerations

- GIN indexes have higher update costs but excellent query performance
- Suitable for read-heavy workloads (document browsing, search)
- JSONB GIN indexes support both containment (`@>`) and existence (`?`) operators
- Regular VACUUM operations maintain index efficiency

### Data Type Choices

- **UUID:** All primary keys use UUID for distributed system compatibility
- **TIMESTAMPTZ:** All timestamps include timezone for global deployment
- **TEXT vs VARCHAR:** 
  - TEXT for unbounded content (notes, content, file_paths)
  - VARCHAR with limits for controlled fields (names, categories)
- **JSONB vs JSON:** JSONB for better query performance and indexing
- **INET:** Specialized type for IP addresses with built-in validation
- **TEXT[]:** Native array type for tags (better than junction table for simple lists)

## Consequences

### Pros

- **Comprehensive Coverage:** Supports full document management lifecycle
- **Flexible Data Structures:** JSONB enables schema evolution without migrations
- **Strong Referential Integrity:** Foreign keys ensure data consistency
- **Performance Optimized:** Strategic indexing for common query patterns
- **Audit Trail:** Execution records provide complete traceability
- **Scalable Access Control:** Flexible sharing model for documents and records
- **Version Management:** Full version history with rollback capability

### Cons

- **Complexity:** Significantly more tables and relationships to manage
- **Storage Requirements:** Version history and attachments increase storage needs
- **JSONB Validation:** Application-level validation required for JSONB structures
- **Migration Effort:** Requires careful migration planning for production deployment
- **Index Overhead:** GIN indexes consume more space and update time
- **Query Complexity:** Cross-aggregate queries may require careful optimization

### Migration Considerations

- New tables should be created in numbered migration files (following ADR 0006)
- Existing `repositories` and `managed_files` tables are unchanged
- Consider creating tables in dependency order to avoid foreign key issues
- Initial deployment may use empty tables populated later
- Production deployment should include data validation and rollback plan

### Security Considerations

- `access_token` encrypted at application level (AES-256-GCM)
- Access control enforced at application layer using `access_scope`
- Administrators have implicit access to all resources
- Attachment storage paths should be validated to prevent directory traversal
- JSONB validation critical to prevent injection attacks
