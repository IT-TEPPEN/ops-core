import { createContext } from "react";
import type { OAuthCommandService } from "../../application";

/**
 * Context for OAuthCommandService dependency injection.
 * Following ADR 0023 - Context-based DI pattern.
 */
export const OAuthCommandServiceContext =
  createContext<OAuthCommandService | null>(null);
