package usecase

import (
	"context"
	"strconv"

	"opscore/backend/internal/execution_record/application/dto"
	apperror "opscore/backend/internal/execution_record/application/error"
	"opscore/backend/internal/execution_record/domain/repository"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// updateStepNotesUseCase updates notes for a specific step.
type updateStepNotesUseCase struct {
	repo repository.ExecutionRecordRepository
}

func newUpdateStepNotesUseCase(repo repository.ExecutionRecordRepository) updateStepNotesUseCase {
	return updateStepNotesUseCase{repo: repo}
}

func (uc updateStepNotesUseCase) Execute(ctx context.Context, req *dto.UpdateStepNotesRequest) (*dto.ExecutionRecordResponse, error) {
	id, err := value_object.NewExecutionRecordID(req.ExecutionRecordID)
	if err != nil {
		return nil, &apperror.ValidationError{Field: "executionRecordID", Message: "invalid execution record ID format"}
	}

	record, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, &apperror.NotFoundError{ResourceType: "ExecutionRecord", ResourceID: req.ExecutionRecordID}
	}

	if err := record.UpdateStepNotes(req.StepNumber, req.Notes); err != nil {
		return nil, &apperror.NotFoundError{ResourceType: "ExecutionStep", ResourceID: "step " + strconv.Itoa(req.StepNumber)}
	}

	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return toExecutionRecordResponse(record), nil
}
