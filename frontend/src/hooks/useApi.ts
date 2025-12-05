/**
 * Custom hook for API calls
 */

import { useState, useCallback, useRef, useEffect } from "react";
import type { ApiError, ApiState } from "../types/api";
import { ApiRequestError } from "../utils/api";

export interface UseApiOptions<T> {
  /** Initial data value */
  initialData?: T;
  /** Execute immediately on mount */
  immediate?: boolean;
}

export interface UseApiResult<T> extends ApiState<T> {
  /** Execute the API call */
  execute: () => Promise<T | null>;
  /** Reset the state */
  reset: () => void;
  /** Set data manually */
  setData: (data: T | null) => void;
}

/**
 * Hook for managing API calls with loading, error, and data states
 */
export function useApi<T>(
  fetcher: (signal: AbortSignal) => Promise<T>,
  options: UseApiOptions<T> = {}
): UseApiResult<T> {
  const { initialData = null, immediate = false } = options;

  const [data, setData] = useState<T | null>(initialData);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<ApiError | null>(null);

  const abortControllerRef = useRef<AbortController | null>(null);
  const isMountedRef = useRef(true);

  useEffect(() => {
    isMountedRef.current = true;
    return () => {
      isMountedRef.current = false;
      abortControllerRef.current?.abort();
    };
  }, []);

  const execute = useCallback(async (): Promise<T | null> => {
    // Abort previous request
    abortControllerRef.current?.abort();
    abortControllerRef.current = new AbortController();

    setIsLoading(true);
    setError(null);

    try {
      const result = await fetcher(abortControllerRef.current.signal);

      if (isMountedRef.current) {
        setData(result);
        setIsLoading(false);
      }
      return result;
    } catch (err) {
      // Ignore abort errors
      if (err instanceof DOMException && err.name === "AbortError") {
        return null;
      }

      if (isMountedRef.current) {
        const apiError =
          err instanceof ApiRequestError
            ? err.toApiError()
            : {
                code: "UNKNOWN_ERROR",
                message: err instanceof Error ? err.message : "An unknown error occurred",
              };
        setError(apiError);
        setIsLoading(false);
      }
      return null;
    }
  }, [fetcher]);

  const reset = useCallback(() => {
    setData(initialData);
    setError(null);
    setIsLoading(false);
  }, [initialData]);

  useEffect(() => {
    if (immediate) {
      execute();
    }
  }, [immediate, execute]);

  return {
    data,
    isLoading,
    error,
    execute,
    reset,
    setData,
  };
}
