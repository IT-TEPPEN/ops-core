import { createContext, useContext, useMemo, ReactNode } from "react";
import type { DocumentCommandService } from "../../application";
import { HttpDocumentCommandService } from "../../infrastructure";

/**
 * Context for DocumentCommandService dependency injection.
 * Following ADR 0023 - Context-based DI pattern.
 */
const DocumentCommandServiceContext = createContext<DocumentCommandService | null>(
  null
);

export function DocumentCommandServiceProvider({
  children,
}: {
  children: ReactNode;
}) {
  const service = useMemo(() => {
    return new HttpDocumentCommandService();
  }, []);

  return (
    <DocumentCommandServiceContext.Provider value={service}>
      {children}
    </DocumentCommandServiceContext.Provider>
  );
}

/**
 * Hook to access DocumentCommandService from context.
 */
export function useDocumentCommandService(): DocumentCommandService {
  const context = useContext(DocumentCommandServiceContext);
  if (!context) {
    throw new Error(
      "useDocumentCommandService must be used within DocumentCommandServiceProvider"
    );
  }
  return context;
}
