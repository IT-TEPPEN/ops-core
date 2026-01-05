import type {
  DocumentCommandService,
  CreateDocumentDto,
  UpdateDocumentDto,
  UpdateDocumentMetadataDto,
  PublishDocumentVersionDto,
  RollbackDocumentVersionDto,
  PublishDocumentFromOAuthDto,
  ValidateVariablesDto,
  DocumentViewData,
  ValidationResult,
} from "../../application";
import { V1ApiClient } from "@/shared/api/client";
import type {
  ApiDocumentResponse,
  ApiDocumentVersionResponse,
  ApiCreateDocumentRequest,
  ApiUpdateDocumentRequest,
  ApiUpdateDocumentMetadataRequest,
  ApiPublishDocumentRequest,
  ApiValidateVariableValuesRequest,
  ApiValidateVariableValuesResponse,
  ApiVariableDefinitionResponse,
} from "../types";

/**
 * HTTP implementation of DocumentCommandService.
 * Handles all write operations (POST/PUT/DELETE/PATCH) for document management.
 * Following ADR 0019 - Command Service pattern.
 */
export class HttpDocumentCommandService
  extends V1ApiClient
  implements DocumentCommandService
{
  constructor() {
    super("/documents");
  }

  async create(dto: CreateDocumentDto): Promise<DocumentViewData> {
    const apiRequest: ApiCreateDocumentRequest = {
      repository_id: dto.repositoryId,
      owner: dto.owner,
      repository: dto.repository,
      provider_repository_id: dto.providerRepositoryId,
      file_path: dto.filePath,
      commit_hash: dto.commitHash,
      access_scope: dto.accessScope,
      is_auto_update: dto.isAutoUpdate,
    };

    const response = await this.post<
      ApiDocumentResponse,
      ApiCreateDocumentRequest
    >("", apiRequest);
    return this.toDocumentViewData(response);
  }

  async update(dto: UpdateDocumentDto): Promise<DocumentViewData> {
    const apiRequest: ApiUpdateDocumentRequest = {
      file_path: dto.filePath,
      commit_hash: dto.commitHash,
    };

    const response = await this.put<
      ApiDocumentResponse,
      ApiUpdateDocumentRequest
    >(`/${dto.docId}`, apiRequest);
    return this.toDocumentViewData(response);
  }

  async updateMetadata(
    dto: UpdateDocumentMetadataDto
  ): Promise<DocumentViewData> {
    const apiRequest: ApiUpdateDocumentMetadataRequest = {
      access_scope: dto.accessScope,
      is_auto_update: dto.isAutoUpdate,
    };

    const response = await this.patch<
      ApiDocumentResponse,
      ApiUpdateDocumentMetadataRequest
    >(`/${dto.docId}/metadata`, apiRequest);
    return this.toDocumentViewData(response);
  }

  async publishVersion(dto: PublishDocumentVersionDto): Promise<DocumentViewData> {
    const response = await this.post<ApiDocumentResponse, Record<string, never>>(
      `/${dto.docId}/versions/${dto.versionNumber}/publish`,
      {}
    );
    return this.toDocumentViewData(response);
  }

  async rollbackVersion(
    dto: RollbackDocumentVersionDto
  ): Promise<DocumentViewData> {
    const response = await this.post<ApiDocumentResponse, Record<string, never>>(
      `/${dto.docId}/versions/${dto.versionNumber}/rollback`,
      {}
    );
    return this.toDocumentViewData(response);
  }

  async publishFromOAuth(
    dto: PublishDocumentFromOAuthDto
  ): Promise<DocumentViewData> {
    const apiRequest: ApiPublishDocumentRequest = {
      connection_id: dto.connectionId,
      owner: dto.owner,
      repository: dto.repository,
      provider_repository_id: dto.providerRepositoryId,
      file_path: dto.filePath,
      ref: dto.ref,
      access_scope: dto.accessScope,
      is_auto_update: dto.isAutoUpdate,
    };

    const response = await this.post<
      ApiDocumentResponse,
      ApiPublishDocumentRequest
    >("/publish", apiRequest);
    return this.toDocumentViewData(response);
  }

  async validateVariables(dto: ValidateVariablesDto): Promise<ValidationResult> {
    const apiRequest: ApiValidateVariableValuesRequest = {
      values: dto.values.map((v) => ({
        name: v.name,
        value: v.value,
      })),
    };

    const response = await this.post<
      ApiValidateVariableValuesResponse,
      ApiValidateVariableValuesRequest
    >(`/${dto.docId}/validate-variables`, apiRequest);

    return {
      valid: response.valid,
      errors: response.errors.map((error) => ({
        name: error.name,
        message: error.message,
      })),
    };
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
  private toDocumentVersionViewData(response: ApiDocumentVersionResponse) {
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
   * Transform API response to VariableDefinitionViewData.
   */
  private toVariableDefinitionViewData(
    response: ApiVariableDefinitionResponse
  ) {
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
