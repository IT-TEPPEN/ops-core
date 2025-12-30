/**
 * Internal API response types (not exported outside infrastructure layer).
 * Following ADR 0022 - API response types belong in infrastructure.
 */

export interface ApiDocumentResponse {
  id: string;
  title: string;
  content: string;
  doc_type: string;
  tags: string[];
  status: string;
  access_scope: string;
  created_at: string;
  updated_at: string;
  owner_name?: string;
}

export interface ApiDocumentListResponse {
  documents: ApiDocumentResponse[];
  pagination?: {
    current_page: number;
    previous_page: number | null;
    next_page: number | null;
    total_pages: number;
    per_page: number;
    current_items: number;
    total_items: number;
  };
}

export interface ApiVersionHistoryResponse {
  versions: Array<{
    id: string;
    version_number: number;
    commit_hash: string;
    created_at: string;
    is_current: boolean;
  }>;
}

export interface ApiCreateDocumentRequest {
  repository_id: string;
  file_path: string;
  title: string;
  content: string;
  doc_type: string;
  tags: string[];
  access_scope: "public" | "private";
}

export interface ApiUpdateDocumentRequest {
  title?: string;
  content?: string;
  doc_type?: string;
  tags?: string[];
  access_scope?: "public" | "private";
}
