import type {
  DocumentQueryService,
  ListDocumentsDto,
  GetDocumentByIdDto,
  GetDocumentVersionHistoryDto,
  GetDocumentVersionDto,
  GetDocumentVariablesDto,
  DocumentViewData,
  DocumentVersionViewData,
  DocumentListItem,
  VersionHistoryItem,
  VariableDefinitionViewData,
} from "../../application";
import { V1ApiClient } from "@/shared/api/client";
import type {
  ApiDocumentResponse,
  ApiDocumentListResponse,
  ApiDocumentListItemResponse,
  ApiVersionHistoryResponse,
  ApiDocumentVersionResponse,
  ApiGetVariableDefinitionsResponse,
  ApiVariableDefinitionResponse,
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

  async list(dto?: ListDocumentsDto): Promise<DocumentListItem[]> {
    const params = new URLSearchParams();
    if (dto?.repositoryId) {
      params.append("repository_id", dto.repositoryId);
    }
    if (dto?.providerRepositoryId) {
      params.append("provider_repository_id", dto.providerRepositoryId);
    }
    if (dto?.owner) {
      params.append("owner", dto.owner);
    }
    if (dto?.repository) {
      params.append("repository", dto.repository);
    }

    const queryString = params.toString();
    const url = queryString ? `?${queryString}` : "";

    const response = await this.get<ApiDocumentListResponse>(url);
    return response.documents.map((doc) => this.toDocumentListItem(doc));
  }

  async getById(dto: GetDocumentByIdDto): Promise<DocumentViewData> {
    const response = await this.get<ApiDocumentResponse>(`/${dto.docId}`);
    return this.toDocumentViewData(response);
  }

  async getVersionHistory(
    dto: GetDocumentVersionHistoryDto
  ): Promise<VersionHistoryItem[]> {
    const response = await this.get<ApiVersionHistoryResponse>(
      `/${dto.docId}/versions`
    );
    return response.versions.map((version) =>
      this.toVersionHistoryItem(version)
    );
  }

  async getVersion(dto: GetDocumentVersionDto): Promise<DocumentVersionViewData> {
    const response = await this.get<ApiDocumentVersionResponse>(
      `/${dto.docId}/versions/${dto.versionNumber}`
    );
    return this.toDocumentVersionViewData(response);
  }

  async getVariables(
    dto: GetDocumentVariablesDto
  ): Promise<VariableDefinitionViewData[]> {
    const response = await this.get<ApiGetVariableDefinitionsResponse>(
      `/${dto.docId}/variables`
    );
    return response.variables.map((variable) =>
      this.toVariableDefinitionViewData(variable)
    );
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
      owner: response.owner,
      repository: response.repository,
      providerRepositoryId: response.provider_repository_id,
      repositoryId: response.repository_id,
      accessScope: response.access_scope,
      isAutoUpdate: response.is_auto_update,
      isPublished: response.is_published,
      versionCount: response.version_count,
      currentVersion: this.toDocumentVersionViewData(response.current_version),
      createdAt: new Date(response.created_at),
      updatedAt: new Date(response.updated_at),
    };
  }

  /**
   * Transform API response to DocumentVersionViewData.
   */
  private toDocumentVersionViewData(
    response: ApiDocumentVersionResponse
  ): DocumentVersionViewData {
    return {
      id: response.id,
      documentId: response.document_id,
      versionNumber: response.version_number,
      filePath: response.file_path,
      commitHash: response.commit_hash,
      title: response.title,
      content: response.content,
      docType: response.doc_type,
      tags: response.tags,
      variables: response.variables.map((v) =>
        this.toVariableDefinitionViewData(v)
      ),
      isCurrent: response.is_current,
      publishedAt: new Date(response.published_at),
      unpublishedAt: response.unpublished_at
        ? new Date(response.unpublished_at)
        : null,
    };
  }

  /**
   * Transform API response to DocumentListItem.
   */
  private toDocumentListItem(
    response: ApiDocumentListItemResponse
  ): DocumentListItem {
    return {
      id: response.id,
      owner: response.owner,
      repository: response.repository,
      providerRepositoryId: response.provider_repository_id,
      repositoryId: response.repository_id,
      title: response.title,
      docType: response.doc_type,
      tags: response.tags,
      isPublished: response.is_published,
      versionCount: response.version_count,
      createdAt: new Date(response.created_at),
      updatedAt: new Date(response.updated_at),
    };
  }

  /**
   * Transform API response to VersionHistoryItem.
   */
  private toVersionHistoryItem(
    response: ApiDocumentVersionResponse
  ): VersionHistoryItem {
    return {
      id: response.id,
      versionNumber: response.version_number,
      commitHash: response.commit_hash,
      title: response.title,
      docType: response.doc_type,
      isCurrent: response.is_current,
      publishedAt: new Date(response.published_at),
      unpublishedAt: response.unpublished_at
        ? new Date(response.unpublished_at)
        : null,
    };
  }

  /**
   * Transform API response to VariableDefinitionViewData.
   */
  private toVariableDefinitionViewData(
    response: ApiVariableDefinitionResponse
  ): VariableDefinitionViewData {
    return {
      name: response.name,
      label: response.label,
      description: response.description,
      type: response.type as "string" | "number" | "boolean" | "date",
      required: response.required,
      defaultValue: response.default_value,
    };
  }
}
