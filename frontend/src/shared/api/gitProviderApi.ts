/**
 * Git Provider API client
 * OAuth経由でGitプロバイダーのリポジトリにアクセスするためのAPI
 */

import { V1ApiClient } from "./client";

export type GitProvider = "github" | "gitlab" | "gitlab-self-hosted";

export interface OAuthConnection {
  id: string;
  provider: GitProvider;
  providerUserId: string;
  providerUsername: string;
  scopes: string[];
  connectedAt: string;
}

export interface GitRepository {
  id: number;
  name: string;
  fullName: string;
  description: string;
  htmlUrl: string;
  cloneUrl: string;
  private: boolean;
  defaultBranch: string;
  owner: {
    login: string;
    avatarUrl: string;
  };
}

export interface ListRepositoriesResponse {
  repositories: GitRepository[];
}

export interface ListConnectionsResponse {
  connections: OAuthConnection[];
}

class GitProviderApiClient extends V1ApiClient {
  constructor() {
    super("");
  }

  /**
   * OAuth接続一覧を取得
   */
  async listConnections(): Promise<OAuthConnection[]> {
    const response = await this.get<ListConnectionsResponse>(
      "/auth/oauth/connections"
    );
    return response.connections || [];
  }

  /**
   * 特定のプロバイダーとの接続状態を確認
   */
  async isConnected(provider: GitProvider): Promise<boolean> {
    try {
      const connections = await this.listConnections();
      return connections.some((conn) => conn.provider === provider);
    } catch {
      return false;
    }
  }

  /**
   * プロバイダーからリポジトリ一覧を取得
   */
  async listRepositories(provider: GitProvider): Promise<GitRepository[]> {
    const response = await this.get<ListRepositoriesResponse>(
      `/git-providers/${provider}/repositories`
    );
    return response.repositories || [];
  }

  /**
   * プロバイダーとの接続を解除
   */
  async disconnect(provider: GitProvider): Promise<void> {
    await this.delete(`/auth/oauth/connections/${provider}`);
  }
}

// シングルトンインスタンス
export const gitProviderApi = new GitProviderApiClient();
