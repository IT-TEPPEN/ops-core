// Execution feature public API

// Application layer (DTOs and Service interfaces)
export type {
  ExecutionQueryService,
  ExecutionCommandService,
  ExecutionRecordViewData,
  ExecutionStepViewData,
  CreateExecutionRecordRequest,
  UpdateExecutionRecordRequest,
  AddExecutionStepRequest,
  UpdateStepNotesRequest,
} from "./application";

// Infrastructure layer (Service implementations)
export {
  HttpExecutionQueryService,
  HttpExecutionCommandService,
} from "./infrastructure";

// Presentation layer (Components, Hooks, Contexts)
export * from "./presentation";

// Legacy exports (for backward compatibility)
// TODO: Remove these after updating all imports
export { useExecutionRecord } from "./hooks/useExecutionRecord";
export { DocumentContentDisplay } from "./components/DocumentContentDisplay";
export { ExecutionHeader } from "./components/ExecutionHeader";
export { ExecutionSteps } from "./components/ExecutionSteps";
