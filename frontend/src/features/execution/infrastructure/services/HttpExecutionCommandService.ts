import type {
  ExecutionCommandService,
  ExecutionRecordViewData,
  CreateExecutionRecordRequest,
  AddExecutionStepRequest,
  UpdateStepNotesRequest,
} from "../../application";
import { V1ApiClient } from "@/shared/api/client";
import type {
  ApiExecutionRecordResponse,
  ApiCreateExecutionRecordRequest,
  ApiAddExecutionStepRequest,
  ApiUpdateStepNotesRequest,
} from "../types";

/**
 * HTTP implementation of ExecutionCommandService.
 * Handles all write operations (POST/PUT/DELETE) for execution record management.
 * Following ADR 0019 - Command Service pattern.
 */
export class HttpExecutionCommandService
  extends V1ApiClient
  implements ExecutionCommandService
{
  constructor() {
    super("/execution-records");
  }

  async create(
    request: CreateExecutionRecordRequest
  ): Promise<ExecutionRecordViewData> {
    const apiRequest: ApiCreateExecutionRecordRequest = {
      document_id: request.documentId,
      document_version_id: request.documentVersionId,
      title: request.title,
      variable_values: request.variableValues,
    };

    const response = await this.post<ApiExecutionRecordResponse, ApiCreateExecutionRecordRequest>(
      "",
      apiRequest
    );
    return this.toExecutionRecordViewData(response);
  }

  async updateTitle(
    recordId: string,
    title: string
  ): Promise<ExecutionRecordViewData> {
    const response = await this.put<
      ApiExecutionRecordResponse,
      { title: string }
    >(`/${recordId}`, { title });
    return this.toExecutionRecordViewData(response);
  }

  async updateNotes(
    recordId: string,
    notes: string
  ): Promise<ExecutionRecordViewData> {
    const response = await this.put<
      ApiExecutionRecordResponse,
      { notes: string }
    >(`/${recordId}`, { notes });
    return this.toExecutionRecordViewData(response);
  }

  async addStep(
    recordId: string,
    request: AddExecutionStepRequest
  ): Promise<ExecutionRecordViewData> {
    const apiRequest: ApiAddExecutionStepRequest = {
      step_number: request.stepNumber,
      description: request.description,
    };

    const response = await this.post<ApiExecutionRecordResponse, ApiAddExecutionStepRequest>(
      `/${recordId}/steps`,
      apiRequest
    );
    return this.toExecutionRecordViewData(response);
  }

  async updateStepNotes(
    recordId: string,
    request: UpdateStepNotesRequest
  ): Promise<ExecutionRecordViewData> {
    const apiRequest: ApiUpdateStepNotesRequest = {
      step_number: request.stepNumber,
      notes: request.notes,
    };

    const response = await this.put<
      ApiExecutionRecordResponse,
      ApiUpdateStepNotesRequest
    >(`/${recordId}/steps/${request.stepNumber}`, apiRequest);
    return this.toExecutionRecordViewData(response);
  }

  async complete(recordId: string): Promise<ExecutionRecordViewData> {
    const response = await this.post<ApiExecutionRecordResponse, {}>(
      `/${recordId}/complete`,
      {}
    );
    return this.toExecutionRecordViewData(response);
  }

  async fail(recordId: string): Promise<ExecutionRecordViewData> {
    const response = await this.post<ApiExecutionRecordResponse, {}>(
      `/${recordId}/fail`,
      {}
    );
    return this.toExecutionRecordViewData(response);
  }

  /**
   * Transform API response to ExecutionRecordViewData.
   * Converts snake_case to camelCase and string dates to Date objects.
   */
  private toExecutionRecordViewData(
    response: ApiExecutionRecordResponse
  ): ExecutionRecordViewData {
    return {
      id: response.id,
      documentId: response.document_id,
      documentVersionId: response.document_version_id,
      title: response.title,
      notes: response.notes,
      status: response.status,
      variableValues: response.variable_values,
      steps: response.steps.map((step) => ({
        stepNumber: step.step_number,
        description: step.description,
        notes: step.notes,
        status: step.status,
        completedAt: step.completed_at ? new Date(step.completed_at) : null,
      })),
      createdAt: new Date(response.created_at),
      updatedAt: new Date(response.updated_at),
      completedAt: response.completed_at
        ? new Date(response.completed_at)
        : null,
    };
  }
}
