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

export interface ApiFileCommitInfo {
  commit_hash: string;
  message: string;
  author: string;
  author_email: string;
  date: string;
}

export interface ApiFileHistoryResponse {
  file_path: string;
  commits: ApiFileCommitInfo[];
}
