package schema

import (
	"opscore/backend/internal/execution_record/application/dto"
)

// ToCreateExecutionRecordDTO converts API schema to application DTO.
func ToCreateExecutionRecordDTO(req CreateExecutionRecordRequest, executorID string) *dto.CreateExecutionRecordRequest {
	variableValues := make([]dto.VariableValueDTO, len(req.VariableValues))
	for i, vv := range req.VariableValues {
		variableValues[i] = dto.VariableValueDTO{
			Name:  vv.Name,
			Value: vv.Value,
		}
	}

	return &dto.CreateExecutionRecordRequest{
		DocumentID:        req.DocumentID,
		DocumentVersionID: req.DocumentVersionID,
		ExecutorID:        executorID,
		Title:             req.Title,
		VariableValues:    variableValues,
	}
}

// FromExecutionRecordDTO converts application DTO to API schema.
func FromExecutionRecordDTO(dtoResp *dto.ExecutionRecordResponse) ExecutionRecordResponse {
	variableValues := make([]VariableValueResponseSchema, len(dtoResp.VariableValues))
	for i, vv := range dtoResp.VariableValues {
		variableValues[i] = VariableValueResponseSchema{
			Name:  vv.Name,
			Value: vv.Value,
		}
	}

	steps := make([]ExecutionStepResponseSchema, len(dtoResp.Steps))
	for i, step := range dtoResp.Steps {
		steps[i] = ExecutionStepResponseSchema{
			ID:                step.ID,
			ExecutionRecordID: step.ExecutionRecordID,
			StepNumber:        step.StepNumber,
			Description:       step.Description,
			Notes:             step.Notes,
			ExecutedAt:        step.ExecutedAt,
		}
	}

	return ExecutionRecordResponse{
		ID:                dtoResp.ID,
		DocumentID:        dtoResp.DocumentID,
		DocumentVersionID: dtoResp.DocumentVersionID,
		ExecutorID:        dtoResp.ExecutorID,
		Title:             dtoResp.Title,
		VariableValues:    variableValues,
		Notes:             dtoResp.Notes,
		Status:            dtoResp.Status,
		AccessScope:       dtoResp.AccessScope,
		Steps:             steps,
		StartedAt:         dtoResp.StartedAt,
		CompletedAt:       dtoResp.CompletedAt,
		CreatedAt:         dtoResp.CreatedAt,
		UpdatedAt:         dtoResp.UpdatedAt,
	}
}

// ToAddStepDTO converts API schema to application DTO.
func ToAddStepDTO(req AddStepRequest, recordID string) *dto.AddStepRequest {
	return &dto.AddStepRequest{
		ExecutionRecordID: recordID,
		StepNumber:        req.StepNumber,
		Description:       req.Description,
	}
}

// ToUpdateStepNotesDTO converts API schema to application DTO.
func ToUpdateStepNotesDTO(req UpdateStepNotesRequest, recordID string, stepNumber int) *dto.UpdateStepNotesRequest {
	return &dto.UpdateStepNotesRequest{
		ExecutionRecordID: recordID,
		StepNumber:        stepNumber,
		Notes:             req.Notes,
	}
}

// ToUpdateNotesDTO converts API schema to application DTO.
func ToUpdateNotesDTO(req UpdateNotesRequest, recordID string) *dto.UpdateNotesRequest {
	return &dto.UpdateNotesRequest{
		ExecutionRecordID: recordID,
		Notes:             req.Notes,
	}
}

// ToUpdateTitleDTO converts API schema to application DTO.
func ToUpdateTitleDTO(req UpdateTitleRequest, recordID string) *dto.UpdateTitleRequest {
	return &dto.UpdateTitleRequest{
		ExecutionRecordID: recordID,
		Title:             req.Title,
	}
}

// ToUpdateAccessScopeDTO converts API schema to application DTO.
func ToUpdateAccessScopeDTO(req UpdateAccessScopeRequest, recordID string) *dto.UpdateAccessScopeRequest {
	return &dto.UpdateAccessScopeRequest{
		ExecutionRecordID: recordID,
		AccessScope:       req.AccessScope,
	}
}

// ToSearchExecutionRecordDTO converts API schema to application DTO.
func ToSearchExecutionRecordDTO(req SearchExecutionRecordRequest) *dto.SearchExecutionRecordRequest {
	return &dto.SearchExecutionRecordRequest{
		ExecutorID:  req.ExecutorID,
		DocumentID:  req.DocumentID,
		Status:      req.Status,
		StartedFrom: req.StartedFrom,
		StartedTo:   req.StartedTo,
	}
}

// FromAttachmentDTO converts application DTO to API schema.
func FromAttachmentDTO(dtoResp *dto.AttachmentResponse) AttachmentResponse {
	return AttachmentResponse{
		ID:                dtoResp.ID,
		ExecutionRecordID: dtoResp.ExecutionRecordID,
		ExecutionStepID:   dtoResp.ExecutionStepID,
		FileName:          dtoResp.FileName,
		FileSize:          dtoResp.FileSize,
		MimeType:          dtoResp.MimeType,
		StorageType:       dtoResp.StorageType,
		UploadedBy:        dtoResp.UploadedBy,
		UploadedAt:        dtoResp.UploadedAt,
	}
}
