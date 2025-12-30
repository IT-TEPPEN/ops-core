import type { OAuthQueryService } from "../../application";
import type { GitProvider, OAuthConnection, GitRepository } from "../../types";
import { V1ApiClient } from "@/shared/api/client";

interface ListConnectionsResponse {
  connections: OAuthConnection[];
}

interface ListRepositoriesResponse {
  repositories: GitRepository[];
}

/**
 * HTTP implementation of OAuthQueryService.
 * Handles all read operations (GET) for OAuth management.
 * Following ADR 0019 - Query Service pattern.
 */
export class HttpOAuthQueryService
  extends V1ApiClient
  implements OAuthQueryService
{
  constructor() {
    super("");
  }

  async listConnections(): Promise<OAuthConnection[]> {
    const response = await this.get<ListConnectionsResponse>(
      "/auth/oauth/connections"
    );
    return response.connections || [];
  }

  async isConnected(provider: GitProvider): Promise<boolean> {
    try {
      const connections = await this.listConnections();
      return connections.some((conn) => conn.provider === provider);
    } catch {
      return false;
    }
  }

  async listRepositories(provider: GitProvider): Promise<GitRepository[]> {
    const response = await this.get<ListRepositoriesResponse>(
      `/git-providers/${provider}/repositories`
    );
    return response.repositories || [];
  }
}
