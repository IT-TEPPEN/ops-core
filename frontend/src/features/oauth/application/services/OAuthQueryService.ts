import { Content, Connection, Repository, Document } from "../dto";

/**
 * OAuth Query Service interface for read operations (GET).
 * Following ADR 0019 - Query/Command separation pattern.
 */
export interface OAuthQueryService {
  /**
   * List all OAuth connections.
   */
  listConnections(): Promise<Connection[]>;

  /**
   * List repositories for a specific OAuth connection.
   */
  listRepositories(connectionId: string): Promise<Repository[]>;

  /**
   * List Contents of a specific OAuth connection's repository.
   */
  listRepositoryContents(
    connectionId: string,
    repositoryId: string,
    repositoryFullName: string,
    path?: string
  ): Promise<Content[]>;

  /**
   * Get file content from a specific OAuth connection's repository.
   */
  getFileContent(
    connectionId: string,
    repositoryId: string,
    repositoryFullName: string,
    filePath: string,
    commitSha?: string
  ): Promise<Document>;
}
