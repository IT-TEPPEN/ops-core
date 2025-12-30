/**
 * Internal API response types.
 * These types are NOT exported from the public API.
 * Only used within Infrastructure layer for API communication.
 */

export interface ApiRepositoryResponse {
  id: string;
  name: string;
  url: string;
  created_at: string;
  updated_at: string;
}

export interface ApiListRepositoriesResponse {
  repositories: ApiRepositoryResponse[];
}

export interface ApiFileNodeResponse {
  path: string;
  type: "file" | "dir";
}

export interface ApiListFilesResponse {
  files: ApiFileNodeResponse[];
}

export interface ApiFileContentResponse {
  repoId: string;
  filePath: string;
  content: string;
}

export interface ApiUpdateTokenResponse {
  message: string;
  repoId: string;
}
