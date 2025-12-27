import { Document, Repositories } from "../data-types";

export interface RepositoryManagementAdapter {
  listRepositories(): Promise<Repositories>;
  getFileContent(repoId: string, filePath: string): Promise<Document>;
}
