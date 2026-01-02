/**
 * Repository feature public API.
 * Following ADR 0012 - Features should be self-contained with clear public APIs.
 * Following ADR 0018/0019/0021 - Layered architecture for complex features.
 */

// Presentation layer exports (components and hooks)
export {
  useRepositoryQueryService,
  useRepositoryCommandService,
  RepositoryQueryServiceProvider,
  RepositoryCommandServiceProvider,
} from "./presentation";

// Application layer exports (ViewData types for external usage)
export type { RepositoryViewData, PagedResponse } from "./application";

// Type exports from existing types
export type { DocumentVariable } from "./types/repository";
