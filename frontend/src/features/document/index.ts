// Document feature public API

// Application layer (DTOs and Service interfaces)
export type {
  DocumentQueryService,
  DocumentCommandService,
  CreateDocumentRequest,
  UpdateDocumentRequest,
  DocumentViewData,
  DocumentListItem,
  VersionHistoryItem,
  PagedResponse,
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
