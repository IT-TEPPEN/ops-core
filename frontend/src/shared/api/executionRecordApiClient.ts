/**
 * Execution Record API Client
 */

import { V1ApiClient } from "./client";
import type {
  ExecutionRecord,
  CreateExecutionRecordRequest,
  SearchExecutionRecordRequest,
} from "../types/domain";

/** Search execution records response */
interface SearchExecutionRecordsResponse {
  execution_records: ExecutionRecord[];
}

export class ExecutionRecordApi extends V1ApiClient {
  constructor() {
    super("/execution-records");
  }

  /**
   * Create a new execution record
   */
  async createExecutionRecord(
    request: CreateExecutionRecordRequest
  ): Promise<ExecutionRecord> {
    return this.post<ExecutionRecord, CreateExecutionRecordRequest>(
      "",
      request
    );
  }

  /**
   * Get execution record by ID
   */
  async getExecutionRecord(id: string): Promise<ExecutionRecord> {
    return this.get<ExecutionRecord>(`/${id}`);
  }

  /**
   * Search execution records
   */
  async searchExecutionRecords(
    params: SearchExecutionRecordRequest = {}
  ): Promise<SearchExecutionRecordsResponse> {
    const queryParams = new URLSearchParams();
    if (params.executor_id)
      queryParams.append("executor_id", params.executor_id);
    if (params.document_id)
      queryParams.append("document_id", params.document_id);
    if (params.status) queryParams.append("status", params.status);
    if (params.started_from)
      queryParams.append("started_from", params.started_from);
    if (params.started_to) queryParams.append("started_to", params.started_to);

    const query = queryParams.toString();
    const endpoint = query ? `?${query}` : "";
    return this.get<SearchExecutionRecordsResponse>(endpoint);
  }

  /**
   * Update execution record title
   */
  async updateExecutionRecordTitle(
    id: string,
    title: string
  ): Promise<ExecutionRecord> {
    return this.put<ExecutionRecord, { title: string }>(`/${id}/title`, {
      title,
    });
  }

  /**
   * Update execution record notes
   */
  async updateExecutionRecordNotes(
    id: string,
    notes: string
  ): Promise<ExecutionRecord> {
    return this.put<ExecutionRecord, { notes: string }>(`/${id}/notes`, {
      notes,
    });
  }

  /**
   * Add a step to execution record
   */
  async addExecutionStep(
    id: string,
    stepNumber: number,
    description: string
  ): Promise<ExecutionRecord> {
    return this.post<
      ExecutionRecord,
      { step_number: number; description: string }
    >(`/${id}/steps`, { step_number: stepNumber, description });
  }

  /**
   * Update step notes
   */
  async updateStepNotes(
    id: string,
    stepNumber: number,
    notes: string
  ): Promise<ExecutionRecord> {
    return this.put<ExecutionRecord, { notes: string }>(
      `/${id}/steps/${stepNumber}/notes`,
      { notes }
    );
  }

  /**
   * Complete execution record
   */
  async completeExecutionRecord(id: string): Promise<ExecutionRecord> {
    return this.post<ExecutionRecord, Record<string, never>>(
      `/${id}/complete`,
      {}
    );
  }

  /**
   * Mark execution record as failed
   */
  async failExecutionRecord(id: string): Promise<ExecutionRecord> {
    return this.post<ExecutionRecord, Record<string, never>>(`/${id}/fail`, {});
  }

  /**
   * Update access scope
   */
  async updateExecutionRecordAccessScope(
    id: string,
    accessScope: "public" | "private"
  ): Promise<ExecutionRecord> {
    return this.put<ExecutionRecord, { access_scope: "public" | "private" }>(
      `/${id}/access-scope`,
      { access_scope: accessScope }
    );
  }

  /**
   * Delete execution record
   */
  async deleteExecutionRecord(id: string): Promise<void> {
    return this.delete<void>(`/${id}`);
  }
}

// Export singleton instance for convenience
export const executionRecordApi = new ExecutionRecordApi();
