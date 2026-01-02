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
