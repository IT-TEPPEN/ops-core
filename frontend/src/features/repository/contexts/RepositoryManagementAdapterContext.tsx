import { createContext } from "react";
import { RepositoryManagementAdapter } from "../api";

export const RepositoryManagementAdapterContext =
  createContext<RepositoryManagementAdapter | null>(null);
