import type {
  DocumentViewData,
  DocumentListItem,
  VersionHistoryItem,
  PagedResponse,
} from "../dto";

/**
 * Document Query Service interface for read operations (GET).
 * Following ADR 0019 - Query/Command separation pattern.
 */
export interface DocumentQueryService {
  /**
   * List all documents with pagination.
   */
  list(): Promise<PagedResponse<DocumentListItem>>;

  /**
   * Get document details by ID.
   */
  getById(docId: string): Promise<DocumentViewData>;

  /**
   * Get version history for a document.
   */
  getVersionHistory(docId: string): Promise<VersionHistoryItem[]>;

  /**
   * Get specific version of a document.
   */
  getVersion(docId: string, versionId: string): Promise<DocumentViewData>;
}
