import type { GitRepository } from "../../types";
import { Connection } from "../dto/connection";

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
  listRepositories(connectionId: string): Promise<GitRepository[]>;
}
