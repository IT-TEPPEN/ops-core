package usecase

import (
	"context"

	apperror "opscore/backend/internal/execution_record/application/error"
	"opscore/backend/internal/execution_record/domain/repository"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// deleteExecutionRecordUseCase deletes an execution record.
type deleteExecutionRecordUseCase struct {
	repo repository.ExecutionRecordRepository
}

func newDeleteExecutionRecordUseCase(repo repository.ExecutionRecordRepository) deleteExecutionRecordUseCase {
	return deleteExecutionRecordUseCase{repo: repo}
}

func (uc deleteExecutionRecordUseCase) Execute(ctx context.Context, recordID string) error {
	id, err := value_object.NewExecutionRecordID(recordID)
	if err != nil {
		return &apperror.ValidationError{Field: "recordID", Message: "invalid execution record ID format"}
	}

	record, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if record == nil {
		return &apperror.NotFoundError{ResourceType: "ExecutionRecord", ResourceID: recordID}
	}

	return uc.repo.Delete(ctx, id)
}
