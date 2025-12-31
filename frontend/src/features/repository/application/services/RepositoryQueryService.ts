import type {
  RepositoryViewData,
  FileNodeViewData,
  DocumentContentViewData,
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
  list(connectionId: string): Promise<PagedResponse<RepositoryViewData>>;

  /**
   * Get repository details by ID.
   */
  getById(repoId: string): Promise<RepositoryViewData>;

  /**
   * List files in a repository.
   */
  listFiles(repoId: string): Promise<FileNodeViewData[]>;

  /**
   * Get file content from a repository.
   * Optionally specify a commit hash to retrieve a specific version.
   */
  getFileContent(
    repoId: string,
    filePath: string,
    commit?: string
  ): Promise<DocumentContentViewData>;

  /**
   * Get file commit history for a specific file.
   */
  getFileHistory(
    repoId: string,
    filePath: string
  ): Promise<FileCommitViewData[]>;
}
