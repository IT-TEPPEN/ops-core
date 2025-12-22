import { createContext, useContext } from "react";
import { RepositoryManagementAdapter } from "../adapters/RepositoryManagementAdapter";

export const RepositoryManagementAdapterContext =
  createContext<RepositoryManagementAdapter | null>(null);

export function useRepositoryManagementAdapter(): RepositoryManagementAdapter {
  const context = useContext(RepositoryManagementAdapterContext);

  if (!context) {
    throw new Error(
      "useRepositoryManagementAdapter must be used within a RepositoryManagementAdapterProvider"
    );
  }

  return context;
}
