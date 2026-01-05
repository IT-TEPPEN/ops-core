import matter from "gray-matter";
import type { OAuthQueryService } from "../../application";
import {
  Content,
  Document,
  Connection,
  Repository,
  DocumentMeta,
  Commit,
} from "../../application/dto";
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
  count: number;
  repositories: Repository[];
}

interface ContentResponse {
  name: string;
  path: string;
  type: string; // "file" or "directory"
  size: number;
  url: string;
}

interface ListContentsResponse {
  files: ContentResponse[];
  path: string;
}

interface GetFileContentResponse {
  content: string;
  path: string;
  sha: string;
  encoding: string;
}

interface FileCommitInfo {
  hash: string;
  message: string;
  author: string;
  authorEmail: string;
  date: string;
}

interface GetFileCommitHistoryResponse {
  file_path: string;
  commits: FileCommitInfo[];
  count: number;
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
    super("/auth/oauth");
  }

  async listConnections(): Promise<Connection[]> {
    const response = await this.get<ListConnectionsResponse>("/connections");

    return response.connections.map((conn) => ({
      id: conn.id,
      provider: conn.provider,
      providerHost: conn.provider_host,
      providerUsername: conn.provider_username,
      connectedAt: new Date(conn.connected_at),
    }));
  }

  async listRepositories(connectionId: string): Promise<Repository[]> {
    const response = await this.get<ListRepositoriesResponse>(
      `/connections/${connectionId}/repositories`
    );
    return response.repositories || [];
  }

  async listRepositoryContents(
    connectionId: string,
    repositoryId: string,
    repositoryFullName: string,
    path = ""
  ): Promise<Content[]> {
    const response = await this.get<ListContentsResponse>(
      `/connections/${connectionId}/repositories/${encodeURIComponent(
        repositoryFullName
      )}/contents?repositoryId=${encodeURIComponent(repositoryId)}${
        path ? `&path=${encodeURIComponent(path)}` : ""
      }`
    );

    return response.files.map((file) => ({
      name: file.name,
      path: file.path,
      type: file.type,
      size: file.size,
      url: file.url,
    }));
  }

  async getFileContent(
    connectionId: string,
    repositoryId: string,
    repositoryFullName: string,
    filePath: string,
    commitSha?: string
  ): Promise<Document> {
    const response = await this.get<GetFileContentResponse>(
      `/connections/${connectionId}/repositories/${encodeURIComponent(
        repositoryFullName
      )}/files/${encodeURIComponent(
        filePath
      )}?repositoryId=${encodeURIComponent(repositoryId)}${
        commitSha ? `&commit_sha=${encodeURIComponent(commitSha)}` : ""
      }`
    );

    const decodedContent =
      response.encoding === "base64"
        ? new TextDecoder("utf-8").decode(
            Uint8Array.from(atob(response.content), (c) => c.charCodeAt(0))
          )
        : response.content;

    const { content, data: meta } = matter(decodedContent);

    return {
      repositoryFullName: repositoryFullName,
      filePath: filePath,
      content,
      meta: meta as DocumentMeta,
      commitHash: commitSha,
    };
  }

  async getFileCommitHistory(
    connectionId: string,
    repositoryId: string,
    repositoryFullName: string,
    filePath: string
  ): Promise<Commit[]> {
    const response = await this.get<GetFileCommitHistoryResponse>(
      `/connections/${connectionId}/repositories/${encodeURIComponent(
        repositoryFullName
      )}/file-histories/${encodeURIComponent(
        filePath
      )}?repositoryId=${encodeURIComponent(repositoryId)}`
    );

    return response.commits.map((commit) => ({
      hash: commit.hash,
      message: commit.message,
      author: commit.author,
      authorEmail: commit.authorEmail,
      date: new Date(commit.date),
    }));
  }
}
