import type {
  DocumentCommandService,
  DocumentViewData,
  CreateDocumentRequest,
  UpdateDocumentRequest,
} from "../../application";
import { V1ApiClient } from "@/shared/api/client";
import type {
  ApiDocumentResponse,
  ApiCreateDocumentRequest,
  ApiUpdateDocumentRequest,
} from "../types";

/**
 * HTTP implementation of DocumentCommandService.
 * Handles all write operations (POST/PUT/DELETE) for document management.
 * Following ADR 0019 - Command Service pattern.
 */
export class HttpDocumentCommandService
  extends V1ApiClient
  implements DocumentCommandService
{
  constructor() {
    super("/documents");
  }

  async create(request: CreateDocumentRequest): Promise<DocumentViewData> {
    const apiRequest: ApiCreateDocumentRequest = {
      repository_id: request.repositoryId,
      file_path: request.filePath,
      title: request.title,
      content: request.content,
      doc_type: request.docType,
      tags: request.tags,
      access_scope: request.accessScope,
    };

    const response = await this.post<ApiDocumentResponse, ApiCreateDocumentRequest>("", apiRequest);
    return this.toDocumentViewData(response);
  }

  async update(
    docId: string,
    request: UpdateDocumentRequest
  ): Promise<DocumentViewData> {
    const apiRequest: ApiUpdateDocumentRequest = {
      title: request.title,
      content: request.content,
      doc_type: request.docType,
      tags: request.tags,
      access_scope: request.accessScope,
    };

    const response = await this.put<ApiDocumentResponse>(
      `/${docId}`,
      apiRequest
    );
    return this.toDocumentViewData(response);
  }

  async deleteDocument(docId: string): Promise<void> {
    await super.delete(`/${docId}`);
  }

  async publish(docId: string): Promise<DocumentViewData> {
    const response = await this.post<ApiDocumentResponse, {}>(
      `/${docId}/publish`,
      {}
    );
    return this.toDocumentViewData(response);
  }

  /**
   * Transform API response to DocumentViewData.
   * Converts snake_case to camelCase and string dates to Date objects.
   */
  private toDocumentViewData(
    response: ApiDocumentResponse
  ): DocumentViewData {
    return {
      id: response.id,
      title: response.title,
      content: response.content,
      docType: response.doc_type,
      tags: response.tags,
      status: response.status,
      accessScope: response.access_scope,
      createdAt: new Date(response.created_at),
      updatedAt: new Date(response.updated_at),
      ownerName: response.owner_name,
    };
  }
}
