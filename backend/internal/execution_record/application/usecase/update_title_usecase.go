package usecase

import (
	"context"

	"opscore/backend/internal/execution_record/application/dto"
	apperror "opscore/backend/internal/execution_record/application/error"
	"opscore/backend/internal/execution_record/domain/repository"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// updateTitleUseCase updates execution title.
type updateTitleUseCase struct {
	repo repository.ExecutionRecordRepository
}

func newUpdateTitleUseCase(repo repository.ExecutionRecordRepository) updateTitleUseCase {
	return updateTitleUseCase{repo: repo}
}

func (uc updateTitleUseCase) Execute(ctx context.Context, req *dto.UpdateTitleRequest) (*dto.ExecutionRecordResponse, error) {
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

	if err := record.UpdateTitle(req.Title); err != nil {
		return nil, &apperror.ValidationError{Field: "title", Message: err.Error()}
	}

	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return toExecutionRecordResponse(record), nil
}
