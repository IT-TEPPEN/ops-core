package usecase

import (
	"context"

	"opscore/backend/internal/execution_record/application/dto"
	apperror "opscore/backend/internal/execution_record/application/error"
	"opscore/backend/internal/execution_record/domain/repository"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// updateNotesUseCase updates overall notes.
type updateNotesUseCase struct {
	repo repository.ExecutionRecordRepository
}

func newUpdateNotesUseCase(repo repository.ExecutionRecordRepository) updateNotesUseCase {
	return updateNotesUseCase{repo: repo}
}

func (uc updateNotesUseCase) Execute(ctx context.Context, req *dto.UpdateNotesRequest) (*dto.ExecutionRecordResponse, error) {
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

	record.UpdateNotes(req.Notes)

	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return toExecutionRecordResponse(record), nil
}
