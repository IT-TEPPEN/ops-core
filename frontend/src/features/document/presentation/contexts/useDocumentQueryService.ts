import { useContext } from "react";
import { DocumentQueryService } from "../../application";
import { DocumentQueryServiceContext } from "./DocumentQueryServiceContext";

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
