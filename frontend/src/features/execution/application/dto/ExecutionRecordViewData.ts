/**
 * Display-ready execution record data.
 * Following ADR 0019 - ViewData pattern.
 */

export interface ExecutionStepViewData {
  stepNumber: number;
  description: string;
  notes: string;
  status: "pending" | "in_progress" | "completed" | "skipped";
  completedAt: Date | null;
}

export interface ExecutionRecordViewData {
  id: string;
  documentId: string;
  documentVersionId: string;
  title: string;
  notes: string;
  status: "in_progress" | "completed" | "failed";
  variableValues: Array<{
    name: string;
    value: string | number | boolean;
  }>;
  steps: ExecutionStepViewData[];
  createdAt: Date;
  updatedAt: Date;
  completedAt: Date | null;
}

export interface CreateExecutionRecordRequest {
  documentId: string;
  documentVersionId: string;
  title: string;
  variableValues: Array<{
    name: string;
    value: string | number | boolean;
  }>;
}

export interface UpdateExecutionRecordRequest {
  title?: string;
  notes?: string;
}

export interface AddExecutionStepRequest {
  stepNumber: number;
  description: string;
}

export interface UpdateStepNotesRequest {
  stepNumber: number;
  notes: string;
}
