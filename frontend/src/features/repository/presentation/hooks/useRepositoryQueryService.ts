import { useContext } from "react";
import { RepositoryQueryServiceContext } from "../contexts";
import type { RepositoryQueryService } from "../../application";

export function useRepositoryQueryService(): RepositoryQueryService {
  const context = useContext(RepositoryQueryServiceContext);

  if (!context) {
    throw new Error(
      "useRepositoryQueryService must be used within a RepositoryQueryServiceProvider"
    );
  }

  return context;
}
