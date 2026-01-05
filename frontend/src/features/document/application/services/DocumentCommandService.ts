import type {
  DocumentViewData,
  ValidationResult,
  CreateDocumentDto,
  UpdateDocumentDto,
  UpdateDocumentMetadataDto,
  PublishDocumentVersionDto,
  RollbackDocumentVersionDto,
  PublishDocumentFromOAuthDto,
  ValidateVariablesDto,
} from "../dto";

/**
 * Document Command Service interface for write operations (POST/PUT/DELETE/PATCH).
 * Following ADR 0019 - Query/Command separation pattern.
 */
export interface DocumentCommandService {
  /**
   * Create a new document from existing repository file.
   */
  create(dto: CreateDocumentDto): Promise<DocumentViewData>;

  /**
   * Update an existing document (creates new version from repository file).
   */
  update(dto: UpdateDocumentDto): Promise<DocumentViewData>;

  /**
   * Update document metadata (access scope, auto-update setting).
   */
  updateMetadata(dto: UpdateDocumentMetadataDto): Promise<DocumentViewData>;

  /**
   * Publish a specific version of a document.
   */
  publishVersion(dto: PublishDocumentVersionDto): Promise<DocumentViewData>;

  /**
   * Rollback document to a previous version.
   */
  rollbackVersion(dto: RollbackDocumentVersionDto): Promise<DocumentViewData>;

  /**
   * Publish a document from an OAuth connection (fetches from Git provider).
   */
  publishFromOAuth(dto: PublishDocumentFromOAuthDto): Promise<DocumentViewData>;

  /**
   * Validate variable values against document's variable definitions.
   */
  validateVariables(dto: ValidateVariablesDto): Promise<ValidationResult>;
}
