import type {
  RepositoryQueryService,
  RepositoryViewData,
  FileCommitViewData,
  PagedResponse,
} from "../../application";
import { V1ApiClient } from "@/shared/api/client";
import type {
  ApiListRepositoriesResponse,
  ApiRepositoryResponse,
  ApiFileHistoryResponse,
} from "../types";

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

  async list(): Promise<PagedResponse<RepositoryViewData>> {
    const response = await this.get<ApiListRepositoriesResponse>("")
      .catch((error) => {
        console.error("Error fetching repositories:", error);
        throw error;
      })
      .finally(() => {
        console.log("Finished fetching repositories");
      });

    const repositories = response.repositories.map((repo) =>
      this.toRepositoryViewData(repo)
    );

    return {
      data: repositories,
      pagination: {
        currentPage: 1,
        previousPage: null,
        nextPage: null,
        totalPages: 1,
        perPage: repositories.length,
        currentItems: repositories.length,
        totalItems: repositories.length,
      },
    };
  }

  async getById(repoId: string): Promise<RepositoryViewData> {
    const response = await this.get<ApiRepositoryResponse>(`/${repoId}`);
    return this.toRepositoryViewData(response);
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

  /**
   * Transform API response to RepositoryViewData.
   * Converts snake_case to camelCase and string dates to Date objects.
   */
  private toRepositoryViewData(
    response: ApiRepositoryResponse
  ): RepositoryViewData {
    return {
      id: response.id,
      name: response.name,
      url: response.url,
      createdAt: new Date(response.created_at),
      updatedAt: new Date(response.updated_at),
    };
  }
}
