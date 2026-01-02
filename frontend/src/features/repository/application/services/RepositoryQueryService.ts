import type {
  RepositoryViewData,
  FileCommitViewData,
  PagedResponse,
} from "../dto";

/**
 * Repository Query Service interface for read operations (GET).
 * Following ADR 0019 - Query/Command separation pattern.
 */
export interface RepositoryQueryService {
  /**
   * List all repositories with pagination.
   */
  list(): Promise<PagedResponse<RepositoryViewData>>;

  /**
   * Get repository details by ID.
   */
  getById(repoId: string): Promise<RepositoryViewData>;

  /**
   * Get file commit history for a specific file.
   */
  getFileHistory(
    repoId: string,
    filePath: string
  ): Promise<FileCommitViewData[]>;
}
