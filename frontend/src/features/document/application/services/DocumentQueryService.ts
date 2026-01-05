import type {
  DocumentViewData,
  DocumentVersionViewData,
  DocumentListItem,
  VersionHistoryItem,
  VariableDefinitionViewData,
  ListDocumentsDto,
  GetDocumentByIdDto,
  GetDocumentVersionHistoryDto,
  GetDocumentVersionDto,
  GetDocumentVariablesDto,
} from "../dto";

/**
 * Document Query Service interface for read operations (GET).
 * Following ADR 0019 - Query/Command separation pattern.
 */
export interface DocumentQueryService {
  /**
   * List all documents with optional filtering.
   */
  list(dto?: ListDocumentsDto): Promise<DocumentListItem[]>;

  /**
   * Get document details by ID.
   */
  getById(dto: GetDocumentByIdDto): Promise<DocumentViewData>;

  /**
   * Get version history for a document.
   */
  getVersionHistory(dto: GetDocumentVersionHistoryDto): Promise<VersionHistoryItem[]>;

  /**
   * Get specific version of a document by version number.
   */
  getVersion(dto: GetDocumentVersionDto): Promise<DocumentVersionViewData>;

  /**
   * Get variable definitions from current version of a document.
   */
  getVariables(dto: GetDocumentVariablesDto): Promise<VariableDefinitionViewData[]>;
}
