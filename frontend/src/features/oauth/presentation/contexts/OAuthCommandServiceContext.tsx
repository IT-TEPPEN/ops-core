import { createContext, useMemo, ReactNode } from "react";
import type { OAuthCommandService } from "../../application";
import { HttpOAuthCommandService } from "../../infrastructure/services";

/**
 * Context for OAuthCommandService dependency injection.
 * Following ADR 0023 - Context-based DI pattern.
 */
export const OAuthCommandServiceContext =
  createContext<OAuthCommandService | null>(null);
