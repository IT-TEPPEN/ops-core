/**
 * API communication utility
 */

import axios, {
  AxiosInstance,
  AxiosResponse,
  AxiosError,
  InternalAxiosRequestConfig,
} from "axios";
import type { ApiError } from "../types/api";

const TOKEN_KEY = "access_token";
const REFRESH_TOKEN_KEY = "refresh_token";
const USER_KEY = "auth_user";

// Global state for token refresh (to prevent duplicate refresh calls)
let isRefreshing = false;
let refreshSubscribers: Array<(token: string) => void> = [];

function subscribeTokenRefresh(cb: (token: string) => void) {
  refreshSubscribers.push(cb);
}

function onTokenRefreshed(token: string) {
  refreshSubscribers.forEach((cb) => cb(token));
  refreshSubscribers = [];
}

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

    // Add request interceptor to include auth token
    this.client.interceptors.request.use((config) => {
      const token = localStorage.getItem(TOKEN_KEY);
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    // Add response interceptor to handle 401 errors with token refresh
    this.client.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        const originalRequest = error.config as InternalAxiosRequestConfig & {
          _retry?: boolean;
        };

        // If error is 401 and we haven't retried yet
        if (error.response?.status === 401 && !originalRequest._retry) {
          if (isRefreshing) {
            // Token refresh is already in progress, queue this request
            return new Promise((resolve) => {
              subscribeTokenRefresh((token: string) => {
                originalRequest.headers.Authorization = `Bearer ${token}`;
                resolve(this.client(originalRequest));
              });
            });
          }

          originalRequest._retry = true;
          isRefreshing = true;

          try {
            const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);

            if (!refreshToken) {
              isRefreshing = false;
              refreshSubscribers = [];
              this.handleLogout();
              return Promise.reject(error);
            }

            // Attempt to refresh the token
            // Extract base URL without resource path (e.g., "/api/v1/auth" -> "/api/v1")
            const baseUrl =
              this.client.defaults.baseURL?.replace(/\/api\/v1.*$/, "") || "";
            const response = await axios.post(
              `${baseUrl}/api/v1/auth/refresh`,
              {
                refresh_token: refreshToken,
              }
            );

            const { token: newAccessToken, refresh_token: newRefreshToken } =
              response.data;

            // Update stored tokens
            localStorage.setItem(TOKEN_KEY, newAccessToken);
            localStorage.setItem(REFRESH_TOKEN_KEY, newRefreshToken);

            // Update authorization header for the original request
            originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;

            // Notify all subscribers with the new token
            onTokenRefreshed(newAccessToken);

            isRefreshing = false;

            // Retry the original request
            return this.client(originalRequest);
          } catch (refreshError) {
            // Token refresh failed, logout
            isRefreshing = false;
            refreshSubscribers = [];
            this.handleLogout();
            return Promise.reject(refreshError);
          }
        }

        return Promise.reject(error);
      }
    );
  }

  private handleLogout() {
    // Clear tokens and user data
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    localStorage.removeItem(USER_KEY);

    if (window.location.pathname === "/login") {
      return;
    }

    window.location.href = "/login";
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

  protected async deleteWithBody<ResponseData, RequestBody>(
    endpoint: string,
    body: RequestBody
  ): Promise<ResponseData> {
    const response = await this.client.delete<ResponseData>(endpoint, {
      data: body,
    });
    return response.data;
  }
}

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
