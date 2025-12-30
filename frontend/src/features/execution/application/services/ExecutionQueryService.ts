import type { ExecutionRecordViewData } from "../dto";

/**
 * Execution Query Service interface for read operations (GET).
 * Following ADR 0019 - Query/Command separation pattern.
 */
export interface ExecutionQueryService {
  /**
   * Get execution record by ID.
   */
  getById(recordId: string): Promise<ExecutionRecordViewData>;

  /**
   * List execution records for a document.
   */
  listByDocument(documentId: string): Promise<ExecutionRecordViewData[]>;
}
