import { useContext } from "react";
import { DocumentCommandService } from "../../application";
import { DocumentCommandServiceContext } from "./DocumentCommandServiceContext";

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
