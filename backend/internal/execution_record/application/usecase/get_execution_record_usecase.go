package usecase

import (
	"context"

	"opscore/backend/internal/execution_record/application/dto"
	apperror "opscore/backend/internal/execution_record/application/error"
	"opscore/backend/internal/execution_record/domain/repository"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// getExecutionRecordUseCase retrieves a record by ID.
type getExecutionRecordUseCase struct {
	repo repository.ExecutionRecordRepository
}

func newGetExecutionRecordUseCase(repo repository.ExecutionRecordRepository) getExecutionRecordUseCase {
	return getExecutionRecordUseCase{repo: repo}
}

func (uc getExecutionRecordUseCase) Execute(ctx context.Context, recordID string) (*dto.ExecutionRecordResponse, error) {
	id, err := value_object.NewExecutionRecordID(recordID)
	if err != nil {
		return nil, &apperror.ValidationError{Field: "recordID", Message: "invalid execution record ID format"}
	}

	record, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, &apperror.NotFoundError{ResourceType: "ExecutionRecord", ResourceID: recordID}
	}

	return toExecutionRecordResponse(record), nil
}
