import type {
  ExecutionRecordViewData,
  CreateExecutionRecordRequest,
  UpdateExecutionRecordRequest,
  AddExecutionStepRequest,
  UpdateStepNotesRequest,
} from "../dto";

/**
 * Execution Command Service interface for write operations (POST/PUT/DELETE).
 * Following ADR 0019 - Query/Command separation pattern.
 */
export interface ExecutionCommandService {
  /**
   * Create a new execution record.
   */
  create(
    request: CreateExecutionRecordRequest
  ): Promise<ExecutionRecordViewData>;

  /**
   * Update execution record title.
   */
  updateTitle(
    recordId: string,
    title: string
  ): Promise<ExecutionRecordViewData>;

  /**
   * Update execution record notes.
   */
  updateNotes(
    recordId: string,
    notes: string
  ): Promise<ExecutionRecordViewData>;

  /**
   * Add a new execution step.
   */
  addStep(
    recordId: string,
    request: AddExecutionStepRequest
  ): Promise<ExecutionRecordViewData>;

  /**
   * Update step notes.
   */
  updateStepNotes(
    recordId: string,
    request: UpdateStepNotesRequest
  ): Promise<ExecutionRecordViewData>;

  /**
   * Mark execution record as completed.
   */
  complete(recordId: string): Promise<ExecutionRecordViewData>;

  /**
   * Mark execution record as failed.
   */
  fail(recordId: string): Promise<ExecutionRecordViewData>;
}
