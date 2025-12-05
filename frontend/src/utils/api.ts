/**
 * API communication utility
 */

import type { ApiResponse, ApiError, RequestConfig } from "../types/api";

/** Base API URL */
const getApiBaseUrl = (): string => {
  const apiHost = import.meta.env.VITE_API_HOST;
  return apiHost ? `http://${apiHost}/api/v1` : "/api";
};

/** Custom error class for API errors */
export class ApiRequestError extends Error {
  constructor(
    public code: string,
    message: string,
    public details?: Record<string, unknown>
  ) {
    super(message);
    this.name = "ApiRequestError";
  }

  toApiError(): ApiError {
    return {
      code: this.code,
      message: this.message,
      details: this.details,
    };
  }
}

/**
 * Makes an API request with the given configuration
 */
export async function apiRequest<T>(
  endpoint: string,
  config: RequestConfig = {}
): Promise<ApiResponse<T>> {
  const { method = "GET", headers = {}, body, signal } = config;

  const url = `${getApiBaseUrl()}${endpoint}`;

  const requestHeaders: Record<string, string> = {
    "Content-Type": "application/json",
    ...headers,
  };

  const requestOptions: RequestInit = {
    method,
    headers: requestHeaders,
    signal,
  };

  if (body !== undefined && method !== "GET") {
    requestOptions.body = JSON.stringify(body);
  }

  const response = await fetch(url, requestOptions);

  const data = await response.json();

  if (!response.ok) {
    throw new ApiRequestError(
      data.code || "API_ERROR",
      data.message || "An error occurred",
      data.details
    );
  }

  return {
    data,
    status: response.status,
    message: data.message,
  };
}

/**
 * GET request helper
 */
export async function get<T>(
  endpoint: string,
  signal?: AbortSignal
): Promise<T> {
  const response = await apiRequest<T>(endpoint, { method: "GET", signal });
  return response.data;
}

/**
 * POST request helper
 */
export async function post<T>(
  endpoint: string,
  body: unknown,
  signal?: AbortSignal
): Promise<T> {
  const response = await apiRequest<T>(endpoint, { method: "POST", body, signal });
  return response.data;
}

/**
 * PUT request helper
 */
export async function put<T>(
  endpoint: string,
  body: unknown,
  signal?: AbortSignal
): Promise<T> {
  const response = await apiRequest<T>(endpoint, { method: "PUT", body, signal });
  return response.data;
}

/**
 * DELETE request helper
 */
export async function del<T>(
  endpoint: string,
  signal?: AbortSignal
): Promise<T> {
  const response = await apiRequest<T>(endpoint, { method: "DELETE", signal });
  return response.data;
}
