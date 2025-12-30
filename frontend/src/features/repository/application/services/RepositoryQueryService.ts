import type {
  RepositoryViewData,
  FileNodeViewData,
  DocumentContentViewData,
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
   * List files in a repository.
   */
  listFiles(repoId: string): Promise<FileNodeViewData[]>;

  /**
   * Get file content from a repository.
   */
  getFileContent(
    repoId: string,
    filePath: string
  ): Promise<DocumentContentViewData>;
}
