/**
 * API communication utility
 */

import axios, { AxiosInstance, AxiosResponse } from "axios";
import type { ApiResponse, ApiError, RequestConfig } from "../types/api";

export class V1ApiClient {
  private client: AxiosInstance;

  constructor(resourcePath?: string) {
    const apiHost = import.meta.env.VITE_API_HOST;
    const baseURL = apiHost
      ? `http://${apiHost}/api/v1${resourcePath ?? ""}`
      : `/api/v1${resourcePath ?? ""}`;

    console.log(`API Base URL: ${baseURL}`);

    this.client = axios.create({
      baseURL,
      headers: {
        "Content-Type": "application/json",
      },
    });
  }

  protected async get<ResponseData>(endpoint: string): Promise<ResponseData> {
    console.log(`GET request to: ${endpoint}`);
    const response = await this.client.get<ResponseData>(endpoint);
    return response.data;
  }

  protected async post<ResponseData, RequestBody>(
    endpoint: string,
    body: RequestBody
  ): Promise<ResponseData> {
    const response = await this.client.post<
      ResponseData,
      AxiosResponse<ResponseData>,
      RequestBody
    >(endpoint, body);
    return response.data;
  }

  protected async put<ResponseData>(
    endpoint: string,
    body: unknown
  ): Promise<ResponseData> {
    const response = await this.client.put<ResponseData>(endpoint, body);
    return response.data;
  }

  protected async delete<ResponseData>(
    endpoint: string
  ): Promise<ResponseData> {
    const response = await this.client.delete<ResponseData>(endpoint);
    return response.data;
  }
}

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
  const response = await apiRequest<T>(endpoint, {
    method: "POST",
    body,
    signal,
  });
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
  const response = await apiRequest<T>(endpoint, {
    method: "PUT",
    body,
    signal,
  });
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
