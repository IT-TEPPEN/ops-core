import type { RepositoryViewData } from "../dto";

/**
 * Repository Command Service interface for write operations (POST/PUT/DELETE).
 * Following ADR 0019 - Query/Command separation pattern.
 */
export interface RepositoryCommandService {
  /**
   * Create a new repository.
   */
  create(data: CreateRepositoryRequest): Promise<RepositoryViewData>;

  /**
   * Update repository information.
   */
  update(
    repoId: string,
    data: UpdateRepositoryRequest
  ): Promise<RepositoryViewData>;

  /**
   * Remove a repository.
   */
  remove(repoId: string): Promise<void>;

  /**
   * Update repository access token.
   */
  updateAccessToken(
    repoId: string,
    accessToken: string
  ): Promise<UpdateTokenResult>;
}

export interface CreateRepositoryRequest {
  name: string;
  url: string;
  provider: string;
  accessToken?: string;
}

export interface UpdateRepositoryRequest {
  name?: string;
  url?: string;
}

export interface UpdateTokenResult {
  message: string;
  repoId: string;
}
