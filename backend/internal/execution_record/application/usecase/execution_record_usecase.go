package usecase

import (
	"context"

	docvo "opscore/backend/internal/document/domain/value_object"
	apperror "opscore/backend/internal/execution_record/application/error"
	"opscore/backend/internal/execution_record/application/dto"
	"opscore/backend/internal/execution_record/domain/entity"
	"opscore/backend/internal/execution_record/domain/repository"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// ExecutionRecordUsecase handles execution record business logic.
type ExecutionRecordUsecase struct {
	repo repository.ExecutionRecordRepository
}

// NewExecutionRecordUsecase creates a new ExecutionRecordUsecase.
func NewExecutionRecordUsecase(repo repository.ExecutionRecordRepository) *ExecutionRecordUsecase {
	return &ExecutionRecordUsecase{repo: repo}
}

// CreateExecutionRecord creates a new execution record.
func (uc *ExecutionRecordUsecase) CreateExecutionRecord(
	ctx context.Context,
	req *dto.CreateExecutionRecordRequest,
) (*dto.ExecutionRecordResponse, error) {
	// Generate ID
	id := value_object.GenerateExecutionRecordID()

	// Parse document ID
	documentID, err := docvo.NewDocumentID(req.DocumentID)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "documentID",
			Message: "invalid document ID format",
		}
	}

	// Parse version ID
	versionID, err := docvo.NewVersionID(req.DocumentVersionID)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "documentVersionID",
			Message: "invalid version ID format",
		}
	}

	// Convert variable values
	variableValues := make([]value_object.VariableValue, 0, len(req.VariableValues))
	for _, vv := range req.VariableValues {
		v, err := value_object.NewVariableValue(vv.Name, vv.Value)
		if err != nil {
			return nil, &apperror.ValidationError{
				Field:   "variableValues",
				Message: "invalid variable value: " + err.Error(),
			}
		}
		variableValues = append(variableValues, v)
	}

	// Create the execution record
	record, err := entity.NewExecutionRecord(
		id,
		documentID,
		versionID,
		req.ExecutorID,
		req.Title,
		variableValues,
	)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "executionRecord",
			Message: err.Error(),
		}
	}

	// Save to repository
	if err := uc.repo.Save(ctx, record); err != nil {
		return nil, err
	}

	return toExecutionRecordResponse(record), nil
}

// GetExecutionRecord retrieves an execution record by ID.
func (uc *ExecutionRecordUsecase) GetExecutionRecord(
	ctx context.Context,
	recordID string,
) (*dto.ExecutionRecordResponse, error) {
	id, err := value_object.NewExecutionRecordID(recordID)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "recordID",
			Message: "invalid execution record ID format",
		}
	}

	record, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, &apperror.NotFoundError{
			ResourceType: "ExecutionRecord",
			ResourceID:   recordID,
		}
	}

	return toExecutionRecordResponse(record), nil
}

// AddStep adds a step to an execution record.
func (uc *ExecutionRecordUsecase) AddStep(
	ctx context.Context,
	req *dto.AddStepRequest,
) (*dto.ExecutionRecordResponse, error) {
	id, err := value_object.NewExecutionRecordID(req.ExecutionRecordID)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "executionRecordID",
			Message: "invalid execution record ID format",
		}
	}

	record, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, &apperror.NotFoundError{
			ResourceType: "ExecutionRecord",
			ResourceID:   req.ExecutionRecordID,
		}
	}

	if err := record.AddStep(req.StepNumber, req.Description); err != nil {
		return nil, &apperror.ValidationError{
			Field:   "step",
			Message: err.Error(),
		}
	}

	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return toExecutionRecordResponse(record), nil
}

// UpdateStepNotes updates notes for a specific step.
func (uc *ExecutionRecordUsecase) UpdateStepNotes(
	ctx context.Context,
	req *dto.UpdateStepNotesRequest,
) (*dto.ExecutionRecordResponse, error) {
	id, err := value_object.NewExecutionRecordID(req.ExecutionRecordID)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "executionRecordID",
			Message: "invalid execution record ID format",
		}
	}

	record, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, &apperror.NotFoundError{
			ResourceType: "ExecutionRecord",
			ResourceID:   req.ExecutionRecordID,
		}
	}

	if err := record.UpdateStepNotes(req.StepNumber, req.Notes); err != nil {
		return nil, &apperror.NotFoundError{
			ResourceType: "ExecutionStep",
			ResourceID:   "step " + string(rune(req.StepNumber)),
		}
	}

	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return toExecutionRecordResponse(record), nil
}

// UpdateNotes updates the overall notes.
func (uc *ExecutionRecordUsecase) UpdateNotes(
	ctx context.Context,
	req *dto.UpdateNotesRequest,
) (*dto.ExecutionRecordResponse, error) {
	id, err := value_object.NewExecutionRecordID(req.ExecutionRecordID)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "executionRecordID",
			Message: "invalid execution record ID format",
		}
	}

	record, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, &apperror.NotFoundError{
			ResourceType: "ExecutionRecord",
			ResourceID:   req.ExecutionRecordID,
		}
	}

	record.UpdateNotes(req.Notes)

	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return toExecutionRecordResponse(record), nil
}

// UpdateTitle updates the execution title.
func (uc *ExecutionRecordUsecase) UpdateTitle(
	ctx context.Context,
	req *dto.UpdateTitleRequest,
) (*dto.ExecutionRecordResponse, error) {
	id, err := value_object.NewExecutionRecordID(req.ExecutionRecordID)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "executionRecordID",
			Message: "invalid execution record ID format",
		}
	}

	record, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, &apperror.NotFoundError{
			ResourceType: "ExecutionRecord",
			ResourceID:   req.ExecutionRecordID,
		}
	}

	if err := record.UpdateTitle(req.Title); err != nil {
		return nil, &apperror.ValidationError{
			Field:   "title",
			Message: err.Error(),
		}
	}

	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return toExecutionRecordResponse(record), nil
}

// Complete marks an execution as completed.
func (uc *ExecutionRecordUsecase) Complete(
	ctx context.Context,
	req *dto.CompleteExecutionRequest,
) (*dto.ExecutionRecordResponse, error) {
	id, err := value_object.NewExecutionRecordID(req.ExecutionRecordID)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "executionRecordID",
			Message: "invalid execution record ID format",
		}
	}

	record, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, &apperror.NotFoundError{
			ResourceType: "ExecutionRecord",
			ResourceID:   req.ExecutionRecordID,
		}
	}

	if err := record.Complete(); err != nil {
		return nil, &apperror.ConflictError{
			ResourceType: "ExecutionRecord",
			Identifier:   req.ExecutionRecordID,
			Reason:       err.Error(),
		}
	}

	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return toExecutionRecordResponse(record), nil
}

// MarkAsFailed marks an execution as failed.
func (uc *ExecutionRecordUsecase) MarkAsFailed(
	ctx context.Context,
	req *dto.MarkAsFailedRequest,
) (*dto.ExecutionRecordResponse, error) {
	id, err := value_object.NewExecutionRecordID(req.ExecutionRecordID)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "executionRecordID",
			Message: "invalid execution record ID format",
		}
	}

	record, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, &apperror.NotFoundError{
			ResourceType: "ExecutionRecord",
			ResourceID:   req.ExecutionRecordID,
		}
	}

	if err := record.MarkAsFailed(); err != nil {
		return nil, &apperror.ConflictError{
			ResourceType: "ExecutionRecord",
			Identifier:   req.ExecutionRecordID,
			Reason:       err.Error(),
		}
	}

	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return toExecutionRecordResponse(record), nil
}

// UpdateAccessScope updates the access scope.
func (uc *ExecutionRecordUsecase) UpdateAccessScope(
	ctx context.Context,
	req *dto.UpdateAccessScopeRequest,
) (*dto.ExecutionRecordResponse, error) {
	id, err := value_object.NewExecutionRecordID(req.ExecutionRecordID)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "executionRecordID",
			Message: "invalid execution record ID format",
		}
	}

	scope, err := value_object.NewAccessScope(req.AccessScope)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "accessScope",
			Message: "invalid access scope",
		}
	}

	record, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, &apperror.NotFoundError{
			ResourceType: "ExecutionRecord",
			ResourceID:   req.ExecutionRecordID,
		}
	}

	record.UpdateAccessScope(scope)

	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return toExecutionRecordResponse(record), nil
}

// SearchExecutionRecords searches for execution records.
func (uc *ExecutionRecordUsecase) SearchExecutionRecords(
	ctx context.Context,
	req *dto.SearchExecutionRecordRequest,
) ([]*dto.ExecutionRecordResponse, error) {
	criteria := repository.SearchCriteria{}

	if req.ExecutorID != nil {
		criteria.ExecutorID = req.ExecutorID
	}

	if req.DocumentID != nil {
		docID, err := docvo.NewDocumentID(*req.DocumentID)
		if err != nil {
			return nil, &apperror.ValidationError{
				Field:   "documentID",
				Message: "invalid document ID format",
			}
		}
		criteria.DocumentID = &docID
	}

	if req.Status != nil {
		status, err := value_object.NewExecutionStatus(*req.Status)
		if err != nil {
			return nil, &apperror.ValidationError{
				Field:   "status",
				Message: "invalid status",
			}
		}
		criteria.Status = &status
	}

	criteria.StartedFrom = req.StartedFrom
	criteria.StartedTo = req.StartedTo

	// Convert variable filters
	for _, vf := range req.VariableFilters {
		criteria.VariableFilters = append(criteria.VariableFilters, repository.VariableFilter{
			Name:  vf.Name,
			Value: vf.Value,
		})
	}

	records, err := uc.repo.Search(ctx, criteria)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ExecutionRecordResponse, len(records))
	for i, record := range records {
		responses[i] = toExecutionRecordResponse(record)
	}

	return responses, nil
}

// GetByExecutorID retrieves execution records by executor ID.
func (uc *ExecutionRecordUsecase) GetByExecutorID(
	ctx context.Context,
	executorID string,
) ([]*dto.ExecutionRecordResponse, error) {
	records, err := uc.repo.FindByExecutorID(ctx, executorID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ExecutionRecordResponse, len(records))
	for i, record := range records {
		responses[i] = toExecutionRecordResponse(record)
	}

	return responses, nil
}

// DeleteExecutionRecord deletes an execution record.
func (uc *ExecutionRecordUsecase) DeleteExecutionRecord(
	ctx context.Context,
	recordID string,
) error {
	id, err := value_object.NewExecutionRecordID(recordID)
	if err != nil {
		return &apperror.ValidationError{
			Field:   "recordID",
			Message: "invalid execution record ID format",
		}
	}

	record, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if record == nil {
		return &apperror.NotFoundError{
			ResourceType: "ExecutionRecord",
			ResourceID:   recordID,
		}
	}

	return uc.repo.Delete(ctx, id)
}

// Helper function to convert entity to DTO response
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
