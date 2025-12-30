import type {
  DocumentQueryService,
  DocumentViewData,
  DocumentListItem,
  VersionHistoryItem,
  PagedResponse,
} from "../../application";
import { V1ApiClient } from "@/shared/api/client";
import type {
  ApiDocumentResponse,
  ApiDocumentListResponse,
  ApiVersionHistoryResponse,
} from "../types";

/**
 * HTTP implementation of DocumentQueryService.
 * Handles all read operations (GET) for document management.
 * Following ADR 0019 - Query Service pattern.
 */
export class HttpDocumentQueryService
  extends V1ApiClient
  implements DocumentQueryService
{
  constructor() {
    super("/documents");
  }

  async list(): Promise<PagedResponse<DocumentListItem>> {
    const response = await this.get<ApiDocumentListResponse>("")
      .catch((error) => {
        console.error("Error fetching documents:", error);
        throw error;
      })
      .finally(() => {
        console.log("Finished fetching documents");
      });

    const documents = response.documents.map((doc) =>
      this.toDocumentListItem(doc)
    );

    const pagination = response.pagination
      ? {
          currentPage: response.pagination.current_page,
          previousPage: response.pagination.previous_page,
          nextPage: response.pagination.next_page,
          totalPages: response.pagination.total_pages,
          perPage: response.pagination.per_page,
          currentItems: response.pagination.current_items,
          totalItems: response.pagination.total_items,
        }
      : {
          currentPage: 1,
          previousPage: null,
          nextPage: null,
          totalPages: 1,
          perPage: documents.length,
          currentItems: documents.length,
          totalItems: documents.length,
        };

    return {
      data: documents,
      pagination,
    };
  }

  async getById(docId: string): Promise<DocumentViewData> {
    const response = await this.get<ApiDocumentResponse>(`/${docId}`);
    return this.toDocumentViewData(response);
  }

  async getVersionHistory(docId: string): Promise<VersionHistoryItem[]> {
    const response = await this.get<ApiVersionHistoryResponse>(
      `/${docId}/versions`
    );
    return response.versions.map((version) => ({
      id: version.id,
      versionNumber: version.version_number,
      commitHash: version.commit_hash,
      createdAt: new Date(version.created_at),
      isCurrent: version.is_current,
    }));
  }

  async getVersion(
    docId: string,
    versionId: string
  ): Promise<DocumentViewData> {
    const response = await this.get<ApiDocumentResponse>(
      `/${docId}/versions/${versionId}`
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

  /**
   * Transform API response to DocumentListItem.
   */
  private toDocumentListItem(
    response: ApiDocumentResponse
  ): DocumentListItem {
    return {
      id: response.id,
      title: response.title,
      docType: response.doc_type,
      tags: response.tags,
      status: response.status,
      accessScope: response.access_scope,
      createdAt: new Date(response.created_at),
      updatedAt: new Date(response.updated_at),
    };
  }
}
