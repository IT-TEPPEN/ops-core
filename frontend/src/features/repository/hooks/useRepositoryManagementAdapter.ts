import { useContext } from "react";
import { RepositoryManagementAdapterContext } from "../contexts";
import { RepositoryManagementAdapter } from "../api";

export function useRepositoryManagementAdapter(): RepositoryManagementAdapter {
  const context = useContext(RepositoryManagementAdapterContext);

  if (!context) {
    throw new Error(
      "useRepositoryManagementAdapter must be used within a RepositoryManagementAdapterProvider"
    );
  }

  return context;
}
