import type {
  RepositoryCommandService,
  RepositoryViewData,
  CreateRepositoryRequest,
  UpdateRepositoryRequest,
  UpdateTokenResult,
} from "../../application";
import { V1ApiClient } from "@/shared/api/client";
import type { ApiRepositoryResponse, ApiUpdateTokenResponse } from "../types";

/**
 * HTTP implementation of RepositoryCommandService.
 * Handles all write operations (POST/PUT/DELETE) for repository management.
 * Following ADR 0019 - Command Service pattern.
 */
export class HttpRepositoryCommandService
  extends V1ApiClient
  implements RepositoryCommandService
{
  constructor() {
    super("/repositories");
  }

  async create(data: CreateRepositoryRequest): Promise<RepositoryViewData> {
    const response = await this.post<
      ApiRepositoryResponse,
      CreateRepositoryRequest
    >("", data);
    return this.toRepositoryViewData(response);
  }

  async update(
    repoId: string,
    data: UpdateRepositoryRequest
  ): Promise<RepositoryViewData> {
    const response = await this.put<ApiRepositoryResponse>(`/${repoId}`, data);
    return this.toRepositoryViewData(response);
  }

  async remove(repoId: string): Promise<void> {
    await super.delete<void>(`/${repoId}`);
  }

  async updateAccessToken(
    repoId: string,
    accessToken: string
  ): Promise<UpdateTokenResult> {
    const response = await this.put<ApiUpdateTokenResponse>(
      `/${repoId}/token`,
      { accessToken }
    );
    return {
      message: response.message,
      repoId: response.repoId,
    };
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
