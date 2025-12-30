import type { GitProvider, OAuthConnection, GitRepository } from "../../types";

/**
 * OAuth Query Service interface for read operations (GET).
 * Following ADR 0019 - Query/Command separation pattern.
 */
export interface OAuthQueryService {
  /**
   * List all OAuth connections.
   */
  listConnections(): Promise<OAuthConnection[]>;

  /**
   * Check if connected to a specific provider.
   */
  isConnected(provider: GitProvider): Promise<boolean>;

  /**
   * List repositories for a specific provider.
   */
  listRepositories(provider: GitProvider): Promise<GitRepository[]>;
}
