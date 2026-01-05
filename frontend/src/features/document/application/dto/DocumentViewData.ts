/**
 * Display-ready document data.
 * Following ADR 0019 - ViewData pattern.
 */
export interface DocumentViewData {
  id: string;
  owner: string;
  repository: string;
  providerRepositoryId: string;
  repositoryId: string;
  accessScope: string;
  isAutoUpdate: boolean;
  isPublished: boolean;
  versionCount: number;
  currentVersion: DocumentVersionViewData;
  createdAt: Date;
  updatedAt: Date;
}

/**
 * Document version detail view data.
 */
export interface DocumentVersionViewData {
  id: string;
  documentId: string;
  versionNumber: number;
  filePath: string;
  commitHash: string;
  title: string;
  content: string;
  docType: string;
  tags: string[];
  variables: VariableDefinitionViewData[];
  isCurrent: boolean;
  publishedAt: Date;
  unpublishedAt: Date | null;
}

/**
 * Variable definition in document.
 */
export interface VariableDefinitionViewData {
  name: string;
  label: string;
  description: string;
  type: "string" | "number" | "boolean" | "date";
  required: boolean;
  defaultValue: unknown;
}

/**
 * List item for document list display.
 */
export interface DocumentListItem {
  id: string;
  owner: string;
  repository: string;
  providerRepositoryId: string;
  repositoryId: string;
  title: string;
  docType: string;
  tags: string[];
  isPublished: boolean;
  versionCount: number;
  createdAt: Date;
  updatedAt: Date;
}

/**
 * Version history item.
 */
export interface VersionHistoryItem {
  id: string;
  versionNumber: number;
  commitHash: string;
  title: string;
  docType: string;
  isCurrent: boolean;
  publishedAt: Date;
  unpublishedAt: Date | null;
}

/**
 * Variable value for execution.
 */
export interface VariableValue {
  name: string;
  value: unknown;
}

/**
 * Validation error for variable value.
 */
export interface ValidationError {
  name: string;
  message: string;
}

/**
 * Validation result for variable values.
 */
export interface ValidationResult {
  valid: boolean;
  errors: ValidationError[];
}
