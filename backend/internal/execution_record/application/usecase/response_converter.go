package usecase

import (
	"opscore/backend/internal/execution_record/application/dto"
	"opscore/backend/internal/execution_record/domain/entity"
)

func toExecutionRecordResponse(record entity.ExecutionRecord) *dto.ExecutionRecordResponse {
	steps := make([]dto.ExecutionStepResponse, len(record.Steps()))
	for i, step := range record.Steps() {
		steps[i] = dto.ExecutionStepResponse{
			ID:                step.ID().String(),
			ExecutionRecordID: step.ExecutionRecordID().String(),
			StepNumber:        step.StepNumber(),
			Description:       step.Description(),
			Notes:             step.Notes(),
			ExecutedAt:        step.ExecutedAt(),
		}
	}

	variableValues := make([]dto.VariableValueDTO, len(record.VariableValues()))
	for i, vv := range record.VariableValues() {
		variableValues[i] = dto.VariableValueDTO{
			Name:  vv.Name(),
			Value: vv.Value(),
		}
	}

	return &dto.ExecutionRecordResponse{
		ID:                record.ID().String(),
		DocumentID:        record.DocumentID().String(),
		DocumentVersionID: record.DocumentVersionID().String(),
		ExecutorID:        record.ExecutorID(),
		Title:             record.Title(),
		VariableValues:    variableValues,
		Notes:             record.Notes(),
		Status:            record.Status().String(),
		AccessScope:       record.AccessScope().String(),
		Steps:             steps,
		StartedAt:         record.StartedAt(),
		CompletedAt:       record.CompletedAt(),
		CreatedAt:         record.CreatedAt(),
		UpdatedAt:         record.UpdatedAt(),
	}
}
