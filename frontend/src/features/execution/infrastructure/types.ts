/**
 * Internal API response types (not exported outside infrastructure layer).
 * Following ADR 0022 - API response types belong in infrastructure.
 */

export interface ApiExecutionStepResponse {
  step_number: number;
  description: string;
  notes: string;
  status: "pending" | "in_progress" | "completed" | "skipped";
  completed_at: string | null;
}

export interface ApiExecutionRecordResponse {
  id: string;
  document_id: string;
  document_version_id: string;
  title: string;
  notes: string;
  status: "in_progress" | "completed" | "failed";
  variable_values: Array<{
    name: string;
    value: string | number | boolean;
  }>;
  steps: ApiExecutionStepResponse[];
  created_at: string;
  updated_at: string;
  completed_at: string | null;
}

export interface ApiCreateExecutionRecordRequest {
  document_id: string;
  document_version_id: string;
  title: string;
  variable_values: Array<{
    name: string;
    value: string | number | boolean;
  }>;
}

export interface ApiUpdateExecutionRecordRequest {
  title?: string;
  notes?: string;
}

export interface ApiAddExecutionStepRequest {
  step_number: number;
  description: string;
}

export interface ApiUpdateStepNotesRequest {
  step_number: number;
  notes: string;
}
