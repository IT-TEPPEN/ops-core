import { RepositoryManagementAdapter } from "../features/repository/adapters/RepositoryManagementAdapter";
import { Repositories, Repository } from "../features/repository/data-types";
import { Pagenation } from "../shared/data-types";
import { V1ApiClient } from "./client";

interface ResponseListRepositories {
  repositories: {
    created_at: string;
    id: string;
    name: string;
    updated_at: string;
    url: string;
  }[];
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
}
