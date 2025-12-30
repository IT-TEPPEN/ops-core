import type {
  ExecutionQueryService,
  ExecutionRecordViewData,
  ExecutionStepViewData,
} from "../../application";
import { V1ApiClient } from "@/shared/api/client";
import type {
  ApiExecutionRecordResponse,
  ApiExecutionStepResponse,
} from "../types";

/**
 * HTTP implementation of ExecutionQueryService.
 * Handles all read operations (GET) for execution record management.
 * Following ADR 0019 - Query Service pattern.
 */
export class HttpExecutionQueryService
  extends V1ApiClient
  implements ExecutionQueryService
{
  constructor() {
    super("/execution-records");
  }

  async getById(recordId: string): Promise<ExecutionRecordViewData> {
    const response = await this.get<ApiExecutionRecordResponse>(
      `/${recordId}`
    );
    return this.toExecutionRecordViewData(response);
  }

  async listByDocument(
    documentId: string
  ): Promise<ExecutionRecordViewData[]> {
    const response = await this.get<{ records: ApiExecutionRecordResponse[] }>(
      `?document_id=${documentId}`
    );
    return response.records.map((record) =>
      this.toExecutionRecordViewData(record)
    );
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
      steps: response.steps.map((step) => this.toExecutionStepViewData(step)),
      createdAt: new Date(response.created_at),
      updatedAt: new Date(response.updated_at),
      completedAt: response.completed_at
        ? new Date(response.completed_at)
        : null,
    };
  }

  /**
   * Transform API response to ExecutionStepViewData.
   */
  private toExecutionStepViewData(
    response: ApiExecutionStepResponse
  ): ExecutionStepViewData {
    return {
      stepNumber: response.step_number,
      description: response.description,
      notes: response.notes,
      status: response.status,
      completedAt: response.completed_at
        ? new Date(response.completed_at)
        : null,
    };
  }
}
