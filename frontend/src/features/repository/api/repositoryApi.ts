import { Document, Repository } from "../types";
import { Pagenation } from "@/shared/types";

export interface Repositories {
  data: Repository[];
  pagenation: Pagenation;
}

export interface FileNode {
  path: string;
  type: "file" | "dir";
}

export interface RepositoryDetail {
  id: string;
  name: string;
  url: string;
  createdAt: string;
  updatedAt: string;
}

export interface UpdateTokenResult {
  message: string;
  repoId: string;
}

export interface RepositoryManagementAdapter {
  listRepositories(): Promise<Repositories>;
  getFileContent(repoId: string, filePath: string): Promise<Document>;
  getRepository(repoId: string): Promise<RepositoryDetail>;
  listFiles(repoId: string): Promise<FileNode[]>;
  updateAccessToken(
    repoId: string,
    accessToken: string
  ): Promise<UpdateTokenResult>;
}
