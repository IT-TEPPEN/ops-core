/**
 * Execution Record API utilities
 */

import { get, post, put, del } from "./api";
import type {
  ExecutionRecord,
  CreateExecutionRecordRequest,
  SearchExecutionRecordRequest,
} from "../types/domain";

/**
 * Create a new execution record
 */
export async function createExecutionRecord(
  request: CreateExecutionRecordRequest
): Promise<ExecutionRecord> {
  return post<ExecutionRecord>("/execution-records", request);
}

/**
 * Get execution record by ID
 */
export async function getExecutionRecord(
  id: string
): Promise<ExecutionRecord> {
  return get<ExecutionRecord>(`/execution-records/${id}`);
}

/**
 * Search execution records
 */
export async function searchExecutionRecords(
  params: SearchExecutionRecordRequest = {}
): Promise<{ execution_records: ExecutionRecord[] }> {
  const queryParams = new URLSearchParams();
  if (params.executor_id) queryParams.append("executor_id", params.executor_id);
  if (params.document_id) queryParams.append("document_id", params.document_id);
  if (params.status) queryParams.append("status", params.status);
  if (params.started_from)
    queryParams.append("started_from", params.started_from);
  if (params.started_to) queryParams.append("started_to", params.started_to);

  const query = queryParams.toString();
  const endpoint = query
    ? `/execution-records?${query}`
    : "/execution-records";
  return get<{ execution_records: ExecutionRecord[] }>(endpoint);
}

/**
 * Update execution record title
 */
export async function updateExecutionRecordTitle(
  id: string,
  title: string
): Promise<ExecutionRecord> {
  return put<ExecutionRecord>(`/execution-records/${id}/title`, { title });
}

/**
 * Update execution record notes
 */
export async function updateExecutionRecordNotes(
  id: string,
  notes: string
): Promise<ExecutionRecord> {
  return put<ExecutionRecord>(`/execution-records/${id}/notes`, { notes });
}

/**
 * Add a step to execution record
 */
export async function addExecutionStep(
  id: string,
  stepNumber: number,
  description: string
): Promise<ExecutionRecord> {
  return post<ExecutionRecord>(`/execution-records/${id}/steps`, {
    step_number: stepNumber,
    description,
  });
}

/**
 * Update step notes
 */
export async function updateStepNotes(
  id: string,
  stepNumber: number,
  notes: string
): Promise<ExecutionRecord> {
  return put<ExecutionRecord>(
    `/execution-records/${id}/steps/${stepNumber}/notes`,
    { notes }
  );
}

/**
 * Complete execution record
 */
export async function completeExecutionRecord(
  id: string
): Promise<ExecutionRecord> {
  return post<ExecutionRecord>(`/execution-records/${id}/complete`, {});
}

/**
 * Mark execution record as failed
 */
export async function failExecutionRecord(
  id: string
): Promise<ExecutionRecord> {
  return post<ExecutionRecord>(`/execution-records/${id}/fail`, {});
}

/**
 * Update access scope
 */
export async function updateExecutionRecordAccessScope(
  id: string,
  accessScope: "public" | "private"
): Promise<ExecutionRecord> {
  return put<ExecutionRecord>(`/execution-records/${id}/access-scope`, {
    access_scope: accessScope,
  });
}

/**
 * Delete execution record
 */
export async function deleteExecutionRecord(id: string): Promise<void> {
  return del<void>(`/execution-records/${id}`);
}
