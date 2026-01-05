/**
 * Internal API response types (not exported outside infrastructure layer).
 * Following ADR 0022 - API response types belong in infrastructure.
 * Based on swagger.yaml specification.
 */

/**
 * Variable definition in document frontmatter
 */
export interface ApiVariableDefinitionResponse {
  name: string;
  label: string;
  description: string;
  type: string; // "string", "number", "boolean", "date"
  required: boolean;
  default_value: unknown;
}

/**
 * Document version detail response
 */
export interface ApiDocumentVersionResponse {
  id: string;
  document_id: string;
  version_number: number;
  file_path: string;
  commit_hash: string;
  title: string;
  content: string;
  doc_type: string;
  tags: string[];
  variables: ApiVariableDefinitionResponse[];
  is_current: boolean;
  published_at: string;
  unpublished_at: string | null;
}

/**
 * Document detail response (includes current version)
 */
export interface ApiDocumentResponse {
  id: string;
  owner: string;
  repository: string;
  provider_repository_id: string;
  repository_id: string;
  access_scope: string;
  is_auto_update: boolean;
  is_published: boolean;
  version_count: number;
  current_version: ApiDocumentVersionResponse;
  created_at: string;
  updated_at: string;
}

/**
 * Document list item response
 */
export interface ApiDocumentListItemResponse {
  id: string;
  owner: string;
  repository: string;
  provider_repository_id: string;
  repository_id: string;
  title: string;
  doc_type: string;
  tags: string[];
  is_published: boolean;
  version_count: number;
  created_at: string;
  updated_at: string;
}

/**
 * Document list response
 */
export interface ApiDocumentListResponse {
  documents: ApiDocumentListItemResponse[];
}

/**
 * Version history response
 */
export interface ApiVersionHistoryResponse {
  document_id: string;
  versions: ApiDocumentVersionResponse[];
}

/**
 * Create document request (from existing repository file)
 */
export interface ApiCreateDocumentRequest {
  repository_id: string;
  owner: string;
  repository: string;
  provider_repository_id: string;
  file_path: string;
  commit_hash?: string; // Optional: if empty, use latest commit
  access_scope: "public" | "private";
  is_auto_update: boolean;
}

/**
 * Update document request (create new version)
 */
export interface ApiUpdateDocumentRequest {
  file_path: string;
  commit_hash?: string; // Optional: if empty, use latest commit
}

/**
 * Update document metadata request
 */
export interface ApiUpdateDocumentMetadataRequest {
  access_scope?: "public" | "private";
  is_auto_update?: boolean;
}

/**
 * Publish document request (from OAuth connection)
 */
export interface ApiPublishDocumentRequest {
  connection_id: string;
  owner: string;
  repository: string;
  provider_repository_id: string;
  file_path: string;
  ref: string; // branch name
  access_scope: "public" | "private";
  is_auto_update: boolean;
}

/**
 * Variable value for validation
 */
export interface ApiVariableValueDto {
  name: string;
  value: unknown;
}

/**
 * Validate variable values request
 */
export interface ApiValidateVariableValuesRequest {
  values: ApiVariableValueDto[];
}

/**
 * Validation error detail
 */
export interface ApiValidationErrorDto {
  name: string;
  message: string;
}

/**
 * Validate variable values response
 */
export interface ApiValidateVariableValuesResponse {
  valid: boolean;
  errors: ApiValidationErrorDto[];
}

/**
 * Get variable definitions response
 */
export interface ApiGetVariableDefinitionsResponse {
  documentId: string;
  variables: ApiVariableDefinitionResponse[];
}
