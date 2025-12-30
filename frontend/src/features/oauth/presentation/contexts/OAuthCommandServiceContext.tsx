import { createContext, useContext, useMemo, ReactNode } from "react";
import type { OAuthCommandService } from "../../application";
import { HttpOAuthCommandService } from "../../infrastructure/services";

/**
 * Context for OAuthCommandService dependency injection.
 * Following ADR 0023 - Context-based DI pattern.
 */
const OAuthCommandServiceContext = createContext<OAuthCommandService | null>(
  null
);

export function OAuthCommandServiceProvider({
  children,
}: {
  children: ReactNode;
}) {
  const service = useMemo(() => {
    return new HttpOAuthCommandService();
  }, []);

  return (
    <OAuthCommandServiceContext.Provider value={service}>
      {children}
    </OAuthCommandServiceContext.Provider>
  );
}

/**
 * Hook to access OAuthCommandService from context.
 */
export function useOAuthCommandService(): OAuthCommandService {
  const context = useContext(OAuthCommandServiceContext);
  if (!context) {
    throw new Error(
      "useOAuthCommandService must be used within OAuthCommandServiceProvider"
    );
  }
  return context;
}
