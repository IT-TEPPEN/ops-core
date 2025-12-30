import { useState } from "react";
import { useDocumentCommandService } from "../contexts";
import type { CreateDocumentRequest } from "../../application";

export interface DocumentRegistrationOptions {
  accessScope: "public" | "private";
  isAutoUpdate: boolean;
}

export interface DocumentRegistrationResult {
  id: string;
  repositoryId: string;
}

/**
 * Hook for registering new documents.
 * Refactored to use DocumentCommandService following ADR 0019.
 */
export function useDocumentRegistration() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const commandService = useDocumentCommandService();

  const registerDocument = async (
    repositoryId: string,
    filePath: string,
    options: DocumentRegistrationOptions,
    commitHash?: string
  ): Promise<DocumentRegistrationResult> => {
    setIsLoading(true);
    setError(null);

    try {
      const request: CreateDocumentRequest = {
        repositoryId,
        filePath,
        title: filePath.split("/").pop() || "Untitled",
        content: "", // Will be populated from repository
        docType: "procedure",
        tags: [],
        accessScope: options.accessScope,
      };

      const response = await commandService.create(request);

      return {
        id: response.id,
        repositoryId: response.id, // Assuming repositoryId is in the response
      };
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to register document";
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  return { registerDocument, isLoading, error };
}
