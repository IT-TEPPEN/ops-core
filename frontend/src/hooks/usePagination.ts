/**
 * Custom hook for pagination
 */

import { useState, useCallback, useMemo } from "react";
import type { PaginationState } from "../types/ui";

export interface UsePaginationOptions {
  /** Initial page (1-indexed) */
  initialPage?: number;
  /** Items per page */
  pageSize?: number;
  /** Total number of items */
  totalItems: number;
  /** Callback when page changes */
  onPageChange?: (page: number) => void;
}

export interface UsePaginationResult extends PaginationState {
  /** Go to a specific page */
  goToPage: (page: number) => void;
  /** Go to the next page */
  nextPage: () => void;
  /** Go to the previous page */
  prevPage: () => void;
  /** Go to the first page */
  firstPage: () => void;
  /** Go to the last page */
  lastPage: () => void;
  /** Set the page size */
  setPageSize: (size: number) => void;
  /** Total number of pages */
  totalPages: number;
  /** Whether there's a next page */
  hasNextPage: boolean;
  /** Whether there's a previous page */
  hasPrevPage: boolean;
  /** Get the slice of items for the current page */
  getPageItems: <T>(items: T[]) => T[];
  /** Get page numbers for display */
  getPageNumbers: (maxVisible?: number) => (number | "...")[];
}

/**
 * Hook for managing pagination state
 */
export function usePagination(options: UsePaginationOptions): UsePaginationResult {
  const { initialPage = 1, pageSize: initialPageSize = 10, totalItems, onPageChange } = options;

  const [currentPage, setCurrentPage] = useState(initialPage);
  const [pageSize, setPageSizeState] = useState(initialPageSize);

  const totalPages = useMemo(
    () => Math.max(1, Math.ceil(totalItems / pageSize)),
    [totalItems, pageSize]
  );

  const hasNextPage = currentPage < totalPages;
  const hasPrevPage = currentPage > 1;

  const goToPage = useCallback(
    (page: number) => {
      const validPage = Math.max(1, Math.min(page, totalPages));
      setCurrentPage(validPage);
      onPageChange?.(validPage);
    },
    [totalPages, onPageChange]
  );

  const nextPage = useCallback(() => {
    if (hasNextPage) {
      goToPage(currentPage + 1);
    }
  }, [hasNextPage, currentPage, goToPage]);

  const prevPage = useCallback(() => {
    if (hasPrevPage) {
      goToPage(currentPage - 1);
    }
  }, [hasPrevPage, currentPage, goToPage]);

  const firstPage = useCallback(() => {
    goToPage(1);
  }, [goToPage]);

  const lastPage = useCallback(() => {
    goToPage(totalPages);
  }, [goToPage, totalPages]);

  const setPageSize = useCallback(
    (size: number) => {
      setPageSizeState(size);
      // Reset to first page when page size changes
      setCurrentPage(1);
      onPageChange?.(1);
    },
    [onPageChange]
  );

  const getPageItems = useCallback(
    <T,>(items: T[]): T[] => {
      const start = (currentPage - 1) * pageSize;
      const end = start + pageSize;
      return items.slice(start, end);
    },
    [currentPage, pageSize]
  );

  const getPageNumbers = useCallback(
    (maxVisible = 7): (number | "...")[] => {
      if (totalPages <= maxVisible) {
        return Array.from({ length: totalPages }, (_, i) => i + 1);
      }

      const halfVisible = Math.floor(maxVisible / 2);
      const pages: (number | "...")[] = [];

      if (currentPage <= halfVisible + 1) {
        // Near the start
        for (let i = 1; i <= maxVisible - 2; i++) {
          pages.push(i);
        }
        pages.push("...");
        pages.push(totalPages);
      } else if (currentPage >= totalPages - halfVisible) {
        // Near the end
        pages.push(1);
        pages.push("...");
        for (let i = totalPages - maxVisible + 3; i <= totalPages; i++) {
          pages.push(i);
        }
      } else {
        // In the middle
        pages.push(1);
        pages.push("...");
        for (let i = currentPage - 1; i <= currentPage + 1; i++) {
          pages.push(i);
        }
        pages.push("...");
        pages.push(totalPages);
      }

      return pages;
    },
    [currentPage, totalPages]
  );

  return {
    currentPage,
    pageSize,
    totalItems,
    totalPages,
    hasNextPage,
    hasPrevPage,
    goToPage,
    nextPage,
    prevPage,
    firstPage,
    lastPage,
    setPageSize,
    getPageItems,
    getPageNumbers,
  };
}
