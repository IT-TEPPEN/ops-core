import matter from "gray-matter";
import { RepositoryManagementAdapter } from "../features/repository/adapters/RepositoryManagementAdapter";
import {
  Document,
  DocumentMeta,
  Repositories,
  Repository,
} from "../features/repository/data-types";
import { Pagenation } from "../shared/data-types";
import { V1ApiClient } from "./client";

/**
 * Custom function to parse frontmatter from markdown content
 * This avoids using gray-matter which depends on Node.js Buffer
 */
function parseFrontmatter(markdown: string): {
  content: string;
  meta: DocumentMeta;
} {
  const result = matter(markdown);

  return {
    content: result.content,
    meta: result.data as DocumentMeta,
  };
}

interface ResponseListRepositories {
  repositories: {
    created_at: string;
    id: string;
    name: string;
    updated_at: string;
    url: string;
  }[];
}

interface ResponseGetFileContent {
  repoId: string;
  filePath: string;
  content: string;
}

class RepositoryImpl implements Repository {
  private readonly id: string;
  private readonly name: string;
  private readonly url: string;
  private readonly createdAt: Date;
  private readonly updatedAt: Date;

  constructor(data: {
    id: string;
    name: string;
    url: string;
    created_at: string;
    updated_at: string;
  }) {
    this.id = data.id;
    this.name = data.name;
    this.url = data.url;
    this.createdAt = new Date(data.created_at);
    this.updatedAt = new Date(data.updated_at);
  }

  getId(): string {
    return this.id;
  }

  getName(): string {
    return this.name;
  }

  getUrl(): string {
    return this.url;
  }

  getCreatedAt(): Date {
    return this.createdAt;
  }

  getUpdatedAt(): Date {
    return this.updatedAt;
  }
}

class DocumentImpl implements Document {
  private readonly repoId: string;
  private readonly filePath: string;
  private readonly meta: DocumentMeta;
  private readonly content: string;

  constructor(data: { repoId: string; filePath: string; content: string }) {
    this.repoId = data.repoId;
    this.filePath = data.filePath;

    const { content, meta } = parseFrontmatter(data.content);
    this.content = content;
    this.meta = meta;
  }

  getRepoId(): string {
    return this.repoId;
  }

  getFilePath(): string {
    return this.filePath;
  }

  getMeta(): DocumentMeta {
    return this.meta;
  }

  getContent(): string {
    return this.content;
  }
}

class PagenationImpl implements Pagenation {
  private readonly currentPage: number;
  private readonly previousPage: number | null;
  private readonly nextPage: number | null;
  private readonly totalPages: number;
  private readonly perPage: number;
  private readonly currentItems: number;
  private readonly totalItems: number;

  constructor(data: {
    currentPage: number;
    previousPage: number | null;
    nextPage: number | null;
    totalPages: number;
    perPage: number;
    currentItems: number;
    totalItems: number;
  }) {
    this.currentPage = data.currentPage;
    this.previousPage = data.previousPage;
    this.nextPage = data.nextPage;
    this.totalPages = data.totalPages;
    this.perPage = data.perPage;
    this.currentItems = data.currentItems;
    this.totalItems = data.totalItems;
  }

  getCurrentPage(): number {
    return this.currentPage;
  }
  getPreviousPage(): number | null {
    return this.previousPage;
  }
  getNextPage(): number | null {
    return this.nextPage;
  }
  getTotalPages(): number {
    return this.totalPages;
  }
  getParPage(): number {
    return this.perPage;
  }
  getCurrentItems(): number {
    return this.currentItems;
  }
  getTotalItems(): number {
    return this.totalItems;
  }
}

export class RepositoryApi
  extends V1ApiClient
  implements RepositoryManagementAdapter
{
  constructor() {
    super("/repositories");
  }

  async listRepositories(): Promise<Repositories> {
    const res = await this.get<ResponseListRepositories>("")
      .catch((error) => {
        console.log("Error fetching repositories:", error);
        throw error;
      })
      .finally(() => {
        console.log("Finished fetching repositories");
      });

    const repositories = res.repositories.map(
      (repoData) => new RepositoryImpl(repoData)
    );

    const pagenation = new PagenationImpl({
      currentPage: 1,
      previousPage: null,
      nextPage: null,
      totalPages: 1,
      perPage: repositories.length,
      currentItems: repositories.length,
      totalItems: repositories.length,
    });

    return {
      data: repositories,
      pagenation: pagenation,
    };
  }

  async getFileContent(repoId: string, filePath: string): Promise<Document> {
    const res = await this.get<ResponseGetFileContent>(
      `/${repoId}/files/content?path=${encodeURIComponent(filePath)}`
    );

    return new DocumentImpl(res);
  }
}
