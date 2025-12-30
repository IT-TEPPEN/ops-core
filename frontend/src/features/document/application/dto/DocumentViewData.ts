/**
 * Display-ready document data.
 * Following ADR 0019 - ViewData pattern.
 */
export interface DocumentViewData {
  id: string;
  title: string;
  content: string;
  docType: string;
  tags: string[];
  status: string;
  accessScope: string;
  createdAt: Date;
  updatedAt: Date;
  ownerName?: string;
}

/**
 * List item for document list display.
 */
export interface DocumentListItem {
  id: string;
  title: string;
  docType: string;
  tags: string[];
  status: string;
  accessScope: string;
  createdAt: Date;
  updatedAt: Date;
}

/**
 * Version history item.
 */
export interface VersionHistoryItem {
  id: string;
  versionNumber: number;
  commitHash: string;
  createdAt: Date;
  isCurrent: boolean;
}

/**
 * Paged response structure.
 */
export interface PagedResponse<T> {
  data: T[];
  pagination: {
    currentPage: number;
    previousPage: number | null;
    nextPage: number | null;
    totalPages: number;
    perPage: number;
    currentItems: number;
    totalItems: number;
  };
}
