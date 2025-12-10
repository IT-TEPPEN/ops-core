/**
 * Domain model type definitions
 */

/** Repository entity */
export interface Repository {
  id: string;
  name: string;
  url: string;
  createdAt: string;
  updatedAt: string;
}

/** File node in repository */
export interface FileNode {
  path: string;
  type: "file" | "dir";
}

/** Variable definition for documents */
export interface VariableDefinition {
  name: string;
  label: string;
  description?: string;
  type: "string" | "number" | "boolean" | "date";
  required: boolean;
  default_value?: string | number | boolean;
}

/** Document version entity */
export interface DocumentVersion {
  id: string;
  document_id: string;
  version_number: number;
  file_path: string;
  commit_hash: string;
  title: string;
  doc_type: DocumentType;
  tags: string[];
  variables: VariableDefinition[];
  content: string;
  published_at: string;
  unpublished_at?: string;
  is_current: boolean;
}

/** Document entity */
export interface Document {
  id: string;
  repository_id: string;
  owner: string;
  is_published: boolean;
  is_auto_update: boolean;
  access_scope: AccessScope;
  current_version?: DocumentVersion;
  version_count: number;
  created_at: string;
  updated_at: string;
}

/** Document list item (simplified for list display) */
export interface DocumentListItem {
  id: string;
  repository_id: string;
  title: string;
  owner: string;
  doc_type: DocumentType;
  tags: string[];
  is_published: boolean;
  version_count: number;
  created_at: string;
  updated_at: string;
}

/** Document type */
export type DocumentType = "procedure" | "knowledge";

/** Access scope */
export type AccessScope = "public" | "private";

/** Document metadata */
export interface DocumentMetadata {
  title?: string;
  author?: string;
  date?: string;
  tags?: string[];
  [key: string]: unknown;
}

/** User entity */
export interface User {
  id: string;
  name: string;
  email: string;
  role: UserRole;
}

/** User roles */
export type UserRole = "admin" | "editor" | "viewer";

/** View History entity */
export interface ViewHistory {
  id: string;
  document_id: string;
  user_id: string;
  viewed_at: string;
  view_duration: number;
}

/** Document Statistics entity */
export interface DocumentStatistics {
  document_id: string;
  total_views: number;
  unique_viewers: number;
  last_viewed_at: string;
  average_view_duration: number;
}

/** User Statistics entity */
export interface UserStatistics {
  user_id: string;
  total_views: number;
  unique_documents: number;
}

/** Popular Document entity */
export interface PopularDocument {
  document_id: string;
  total_views: number;
  unique_viewers: number;
  last_viewed_at: string;
}

/** Recent Document entity */
export interface RecentDocument {
  document_id: string;
  last_viewed_at: string;
  total_views: number;
}

/** Execution record */
export interface ExecutionRecord {
  id: string;
  document_id: string;
  document_version_id: string;
  executor_id: string;
  title: string;
  variable_values: VariableValue[];
  notes: string;
  status: ExecutionStatus;
  access_scope: AccessScope;
  steps: ExecutionStep[];
  started_at: string;
  completed_at?: string;
  created_at: string;
  updated_at: string;
}

/** Variable value */
export interface VariableValue {
  name: string;
  value: string | number | boolean;
}

/** Execution step */
export interface ExecutionStep {
  id: string;
  execution_record_id: string;
  step_number: number;
  description: string;
  notes: string;
  executed_at: string;
}

/** Execution status */
export type ExecutionStatus = "in_progress" | "completed" | "failed";

/** Evidence entry */
export interface Evidence {
  id: string;
  execution_record_id: string;
  type: EvidenceType;
  content: string;
  created_at: string;
}

/** Evidence types */
export type EvidenceType = "screenshot" | "log" | "output" | "note";

/** Create execution record request */
export interface CreateExecutionRecordRequest {
  document_id: string;
  document_version_id: string;
  title: string;
  variable_values: VariableValue[];
}

/** Search execution record request */
export interface SearchExecutionRecordRequest {
  executor_id?: string;
  document_id?: string;
  status?: ExecutionStatus;
  started_from?: string;
  started_to?: string;
}
