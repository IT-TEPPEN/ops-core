import { useContext } from "react";
import { RepositoryCommandServiceContext } from "../contexts";
import type { RepositoryCommandService } from "../../application";

export function useRepositoryCommandService(): RepositoryCommandService {
  const context = useContext(RepositoryCommandServiceContext);

  if (!context) {
    throw new Error(
      "useRepositoryCommandService must be used within a RepositoryCommandServiceProvider"
    );
  }

  return context;
}
