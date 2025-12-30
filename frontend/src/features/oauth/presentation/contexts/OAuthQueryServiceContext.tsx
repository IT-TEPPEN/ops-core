import { createContext, useContext, useMemo, ReactNode } from "react";
import type { OAuthQueryService } from "../../application";
import { HttpOAuthQueryService } from "../../infrastructure/services";

/**
 * Context for OAuthQueryService dependency injection.
 * Following ADR 0023 - Context-based DI pattern.
 */
const OAuthQueryServiceContext = createContext<OAuthQueryService | null>(null);

export function OAuthQueryServiceProvider({
  children,
}: {
  children: ReactNode;
}) {
  const service = useMemo(() => {
    return new HttpOAuthQueryService();
  }, []);

  return (
    <OAuthQueryServiceContext.Provider value={service}>
      {children}
    </OAuthQueryServiceContext.Provider>
  );
}

/**
 * Hook to access OAuthQueryService from context.
 */
export function useOAuthQueryService(): OAuthQueryService {
  const context = useContext(OAuthQueryServiceContext);
  if (!context) {
    throw new Error(
      "useOAuthQueryService must be used within OAuthQueryServiceProvider"
    );
  }
  return context;
}
