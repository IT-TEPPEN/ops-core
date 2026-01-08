import { createContext } from "react";
import type { DocumentCommandService } from "../../application";

/**
 * Context for DocumentCommandService dependency injection.
 * Following ADR 0023 - Context-based DI pattern.
 */
export const DocumentCommandServiceContext =
  createContext<DocumentCommandService | null>(null);
