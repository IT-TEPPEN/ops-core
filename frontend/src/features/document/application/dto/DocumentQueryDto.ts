/**
 * DTO for listing documents with optional filtering.
 */
export interface ListDocumentsDto {
  repositoryId?: string;
  providerRepositoryId?: string;
  owner?: string;
  repository?: string;
}

/**
 * DTO for getting document by ID.
 */
export interface GetDocumentByIdDto {
  docId: string;
}

/**
 * DTO for getting document version history.
 */
export interface GetDocumentVersionHistoryDto {
  docId: string;
}

/**
 * DTO for getting specific document version.
 */
export interface GetDocumentVersionDto {
  docId: string;
  versionNumber: number;
}

/**
 * DTO for getting document variables.
 */
export interface GetDocumentVariablesDto {
  docId: string;
}
