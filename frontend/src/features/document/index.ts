// Document feature public API

// Application layer (DTOs and Service interfaces)
export type {
  // Service interfaces
  DocumentQueryService,
  DocumentCommandService,
  // ViewData types
  DocumentViewData,
  DocumentVersionViewData,
  DocumentListItem,
  VersionHistoryItem,
  VariableDefinitionViewData,
  VariableValue,
  ValidationError,
  ValidationResult,
  // Query DTOs
  ListDocumentsDto,
  GetDocumentByIdDto,
  GetDocumentVersionHistoryDto,
  GetDocumentVersionDto,
  GetDocumentVariablesDto,
  // Command DTOs
  CreateDocumentDto,
  UpdateDocumentDto,
  UpdateDocumentMetadataDto,
  PublishDocumentVersionDto,
  RollbackDocumentVersionDto,
  PublishDocumentFromOAuthDto,
  ValidateVariablesDto,
} from "./application";

// Infrastructure layer (Service implementations)
export {
  HttpDocumentQueryService,
  HttpDocumentCommandService,
} from "./infrastructure";

// Presentation layer (Components, Hooks, Contexts)
export * from "./presentation";

// Legacy exports (for backward compatibility)
// TODO: Remove these after updating all imports
export { useDocumentExecution } from "./hooks/useDocumentExecution";
export { DocumentRegistrationDialog } from "./components/DocumentRegistrationDialog";
