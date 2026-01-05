import type { VariableValue } from "./DocumentViewData";

/**
 * DTO for creating a new document from existing repository file.
 */
export interface CreateDocumentDto {
  repositoryId: string;
  owner: string;
  repository: string;
  providerRepositoryId: string;
  filePath: string;
  commitHash?: string; // Optional: if empty, use latest commit
  accessScope: "public" | "private";
  isAutoUpdate: boolean;
}

/**
 * DTO for updating a document (creates new version).
 */
export interface UpdateDocumentDto {
  docId: string;
  filePath: string;
  commitHash?: string; // Optional: if empty, use latest commit
}

/**
 * DTO for updating document metadata.
 */
export interface UpdateDocumentMetadataDto {
  docId: string;
  accessScope?: "public" | "private";
  isAutoUpdate?: boolean;
}

/**
 * DTO for publishing a specific version of a document.
 */
export interface PublishDocumentVersionDto {
  docId: string;
  versionNumber: number;
}

/**
 * DTO for rolling back document to a previous version.
 */
export interface RollbackDocumentVersionDto {
  docId: string;
  versionNumber: number;
}

/**
 * DTO for publishing a document from OAuth connection.
 */
export interface PublishDocumentFromOAuthDto {
  connectionId: string;
  owner: string;
  repository: string;
  providerRepositoryId: string;
  filePath: string;
  ref: string; // branch name
  accessScope: "public" | "private";
  isAutoUpdate: boolean;
}

/**
 * DTO for validating variable values.
 */
export interface ValidateVariablesDto {
  docId: string;
  values: VariableValue[];
}
