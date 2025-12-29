/**
 * API communication utility
 */

import axios, { AxiosInstance, AxiosResponse } from "axios";
import type { ApiError } from "../types/api";

const TOKEN_KEY = "auth_token";

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
