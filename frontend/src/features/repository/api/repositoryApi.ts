import { Document, Repository } from "../types";
import { Pagenation } from "../../../shared/data-types";

export interface Repositories {
  data: Repository[];
  pagenation: Pagenation;
}

export interface RepositoryManagementAdapter {
  listRepositories(): Promise<Repositories>;
  getFileContent(repoId: string, filePath: string): Promise<Document>;
}
