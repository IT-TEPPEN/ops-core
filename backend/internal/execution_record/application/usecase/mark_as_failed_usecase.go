package usecase

import (
	"context"

	"opscore/backend/internal/execution_record/application/dto"
	apperror "opscore/backend/internal/execution_record/application/error"
	"opscore/backend/internal/execution_record/domain/repository"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// markAsFailedUseCase marks an execution as failed.
type markAsFailedUseCase struct {
	repo repository.ExecutionRecordRepository
}

func newMarkAsFailedUseCase(repo repository.ExecutionRecordRepository) markAsFailedUseCase {
	return markAsFailedUseCase{repo: repo}
}

func (uc markAsFailedUseCase) Execute(ctx context.Context, req *dto.MarkAsFailedRequest) (*dto.ExecutionRecordResponse, error) {
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

	if err := record.MarkAsFailed(); err != nil {
		return nil, &apperror.ValidationError{Field: "status", Message: err.Error()}
	}

	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return toExecutionRecordResponse(record), nil
}
