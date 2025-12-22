import { Repositories } from "../data-types";

export interface RepositoryManagementAdapter {
  listRepositories(): Promise<Repositories>;
}
