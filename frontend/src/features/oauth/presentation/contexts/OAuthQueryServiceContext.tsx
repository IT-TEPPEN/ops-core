import { createContext } from "react";
import type { OAuthQueryService } from "../../application";

/**
 * Context for OAuthQueryService dependency injection.
 * Following ADR 0023 - Context-based DI pattern.
 */
export const OAuthQueryServiceContext = createContext<OAuthQueryService | null>(
  null
);
