import type {
  RepositoryQueryService,
  FileCommitViewData,
} from "../../application";
import { V1ApiClient } from "@/shared/api/client";
import type { ApiFileHistoryResponse } from "../types";

/**
 * HTTP implementation of RepositoryQueryService.
 * Handles all read operations (GET) for repository management.
 * Following ADR 0019 - Query Service pattern.
 */
export class HttpRepositoryQueryService
  extends V1ApiClient
  implements RepositoryQueryService
{
  constructor() {
    super("/repositories");
  }

  async getFileHistory(
    repoId: string,
    filePath: string
  ): Promise<FileCommitViewData[]> {
    const response = await this.get<ApiFileHistoryResponse>(
      `/${repoId}/files/history?path=${encodeURIComponent(filePath)}`
    );

    return response.commits.map((commit) => ({
      commitHash: commit.commit_hash,
      message: commit.message,
      author: commit.author,
      authorEmail: commit.author_email,
      date: new Date(commit.date),
    }));
  }
}
