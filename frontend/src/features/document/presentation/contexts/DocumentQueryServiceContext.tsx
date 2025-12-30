import { createContext, useContext, useMemo, ReactNode } from "react";
import type { DocumentQueryService } from "../../application";
import { HttpDocumentQueryService } from "../../infrastructure";

/**
 * Context for DocumentQueryService dependency injection.
 * Following ADR 0023 - Context-based DI pattern.
 */
const DocumentQueryServiceContext = createContext<DocumentQueryService | null>(
  null
);

export function DocumentQueryServiceProvider({
  children,
}: {
  children: ReactNode;
}) {
  const service = useMemo(() => {
    return new HttpDocumentQueryService();
  }, []);

  return (
    <DocumentQueryServiceContext.Provider value={service}>
      {children}
    </DocumentQueryServiceContext.Provider>
  );
}

/**
 * Hook to access DocumentQueryService from context.
 */
export function useDocumentQueryService(): DocumentQueryService {
  const context = useContext(DocumentQueryServiceContext);
  if (!context) {
    throw new Error(
      "useDocumentQueryService must be used within DocumentQueryServiceProvider"
    );
  }
  return context;
}
