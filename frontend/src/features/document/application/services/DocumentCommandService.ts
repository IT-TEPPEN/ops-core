import type { DocumentViewData } from "../dto";

/**
 * Request for creating a new document.
 */
export interface CreateDocumentRequest {
  repositoryId: string;
  filePath: string;
  title: string;
  content: string;
  docType: string;
  tags: string[];
  accessScope: "public" | "private";
}

/**
 * Request for updating a document.
 */
export interface UpdateDocumentRequest {
  title?: string;
  content?: string;
  docType?: string;
  tags?: string[];
  accessScope?: "public" | "private";
}

/**
 * Document Command Service interface for write operations (POST/PUT/DELETE).
 * Following ADR 0019 - Query/Command separation pattern.
 */
export interface DocumentCommandService {
  /**
   * Create a new document.
   */
  create(request: CreateDocumentRequest): Promise<DocumentViewData>;

  /**
   * Update an existing document.
   */
  update(docId: string, request: UpdateDocumentRequest): Promise<DocumentViewData>;

  /**
   * Delete a document.
   */
  deleteDocument(docId: string): Promise<void>;

  /**
   * Publish a document.
   */
  publish(docId: string): Promise<DocumentViewData>;
}
