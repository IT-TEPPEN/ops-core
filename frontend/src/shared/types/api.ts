/**
 * API type definitions
 */

/** Generic API response wrapper */
export interface ApiResponse<T> {
  data: T;
  message?: string;
  status: number;
}

/** API error response */
export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}

/** Pagination parameters */
export interface PaginationParams {
  page: number;
  limit: number;
  sortBy?: string;
  sortOrder?: "asc" | "desc";
}

/** Paginated response */
export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

/** Request configuration */
export interface RequestConfig {
  method?: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
  headers?: Record<string, string>;
  body?: unknown;
  signal?: AbortSignal;
}

/** API state for async operations */
export interface ApiState<T> {
  data: T | null;
  isLoading: boolean;
  error: ApiError | null;
}
