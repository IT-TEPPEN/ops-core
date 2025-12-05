-- Filepath: backend/internal/git_repository/infrastructure/persistence/migrations/000006_create_execution_records_tables.up.sql
-- Create execution_records, execution_steps, and attachments tables

-- execution_records table
CREATE TABLE execution_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    document_version_id UUID NOT NULL REFERENCES document_versions(id),
    executor_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    variable_values JSONB,
    notes TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed', 'failed')),
    access_scope VARCHAR(50) NOT NULL CHECK (access_scope IN ('public', 'private')),
    started_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_execution_records_document_id ON execution_records(document_id);
CREATE INDEX idx_execution_records_executor_id ON execution_records(executor_id);
CREATE INDEX idx_execution_records_status ON execution_records(status);
CREATE INDEX idx_execution_records_started_at ON execution_records(started_at);
CREATE INDEX idx_execution_records_completed_at ON execution_records(completed_at);
CREATE INDEX idx_execution_records_variable_values ON execution_records USING GIN(variable_values);

-- execution_steps table
CREATE TABLE execution_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    execution_record_id UUID NOT NULL REFERENCES execution_records(id) ON DELETE CASCADE,
    step_number INTEGER NOT NULL,
    description TEXT NOT NULL,
    notes TEXT,
    executed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (execution_record_id, step_number)
);

CREATE INDEX idx_execution_steps_execution_record_id ON execution_steps(execution_record_id);

-- attachments table
CREATE TABLE attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    execution_record_id UUID NOT NULL REFERENCES execution_records(id) ON DELETE CASCADE,
    execution_step_id UUID NOT NULL REFERENCES execution_steps(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(127) NOT NULL,
    storage_type VARCHAR(50) NOT NULL CHECK (storage_type IN ('local', 's3', 'minio')),
    storage_path TEXT NOT NULL,
    uploaded_by UUID NOT NULL REFERENCES users(id),
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_attachments_execution_record_id ON attachments(execution_record_id);
CREATE INDEX idx_attachments_execution_step_id ON attachments(execution_step_id);
CREATE INDEX idx_attachments_uploaded_by ON attachments(uploaded_by);
