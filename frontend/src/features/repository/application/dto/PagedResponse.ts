/**
 * Generic paged response wrapper.
 * Replaces PagenationImpl class with a simpler interface.
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
