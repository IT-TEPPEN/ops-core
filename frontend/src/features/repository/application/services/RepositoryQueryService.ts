import type { FileCommitViewData } from "../dto";

/**
 * Repository Query Service interface for read operations (GET).
 * Following ADR 0019 - Query/Command separation pattern.
 */
export interface RepositoryQueryService {
  /**
   * Get file commit history for a specific file.
   */
  getFileHistory(
    repoId: string,
    filePath: string
  ): Promise<FileCommitViewData[]>;
}
