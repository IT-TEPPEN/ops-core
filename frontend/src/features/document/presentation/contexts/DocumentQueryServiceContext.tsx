import { createContext } from "react";
import type { DocumentQueryService } from "../../application";

/**
 * Context for DocumentQueryService dependency injection.
 * Following ADR 0023 - Context-based DI pattern.
 */
export const DocumentQueryServiceContext =
  createContext<DocumentQueryService | null>(null);
