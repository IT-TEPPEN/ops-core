import type { OAuthQueryService } from "../../application";
import { Connection } from "../../application/dto/connection";
import type { GitRepository } from "../../types";
import { V1ApiClient } from "@/shared/api/client";

interface OAuthConnection {
  id: string;
  provider: string;
  provider_host: string;
  provider_username: string;
  scopes: string[];
  connected_at: string;
}

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

  async listConnections(): Promise<Connection[]> {
    const response = await this.get<ListConnectionsResponse>(
      "/auth/oauth/connections"
    );

    return response.connections.map((conn) => ({
      id: conn.id,
      provider: conn.provider,
      providerHost: conn.provider_host,
      providerUsername: conn.provider_username,
      connectedAt: new Date(conn.connected_at),
    }));
  }

  async listRepositories(connectionId: string): Promise<GitRepository[]> {
    const response = await this.get<ListRepositoriesResponse>(
      `/auth/oauth/connections/${connectionId}/repositories`
    );
    return response.repositories || [];
  }
}
