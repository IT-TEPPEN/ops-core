import { V1ApiClient } from "@/shared/api/client";

export interface VariableDefinitionRequest {
  name: string;
  label: string;
  description: string;
  type: string;
  required: boolean;
  default_value: unknown;
}

// Simplified CreateDocumentRequest - frontmatter fields are extracted by backend
export interface CreateDocumentRequest {
  repository_id: string;
  file_path: string;
  commit_hash?: string; // Optional: if empty, backend uses latest commit
  access_scope: string;
  is_auto_update: boolean;
}

export interface DocumentVersionResponse {
  id: string;
  document_id: string;
  version_number: number;
  file_path: string;
  commit_hash: string;
  title: string;
  doc_type: string;
  tags: string[];
  variables: VariableDefinitionRequest[];
  content: string;
  published_at: string;
  unpublished_at: string | null;
  is_current: boolean;
}

export interface DocumentResponse {
  id: string;
  repository_id: string;
  owner: string;
  is_published: boolean;
  is_auto_update: boolean;
  access_scope: string;
  current_version: DocumentVersionResponse | null;
  version_count: number;
  created_at: string;
  updated_at: string;
}

export class DocumentApiClient extends V1ApiClient {
  constructor() {
    super("/documents");
  }

  async createDocument(
    request: CreateDocumentRequest
  ): Promise<DocumentResponse> {
    return this.post<DocumentResponse, CreateDocumentRequest>("", request);
  }
}
