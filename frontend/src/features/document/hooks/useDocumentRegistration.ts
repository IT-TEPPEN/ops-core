import { DocumentApiClient } from "@/adapters";
import { useState } from "react";

export interface DocumentRegistrationOptions {
  accessScope: "public" | "private";
  isAutoUpdate: boolean;
}

export interface DocumentRegistrationResult {
  id: string;
  repository_id: string;
}

export function useDocumentRegistration() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const apiClient = new DocumentApiClient();

  const registerDocument = async (
    repositoryId: string,
    filePath: string,
    options: DocumentRegistrationOptions,
    commitHash?: string
  ): Promise<DocumentRegistrationResult> => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await apiClient.createDocument({
        repository_id: repositoryId,
        file_path: filePath,
        commit_hash: commitHash, // Optional: backend uses latest if not provided
        access_scope: options.accessScope,
        is_auto_update: options.isAutoUpdate,
      });

      return {
        id: response.id,
        repository_id: response.repository_id,
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
