package usecase

import (
	"context"

	"opscore/backend/internal/execution_record/application/dto"
	apperror "opscore/backend/internal/execution_record/application/error"
	"opscore/backend/internal/execution_record/domain/repository"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// updateAccessScopeUseCase updates access scope of a record.
type updateAccessScopeUseCase struct {
	repo repository.ExecutionRecordRepository
}

func newUpdateAccessScopeUseCase(repo repository.ExecutionRecordRepository) updateAccessScopeUseCase {
	return updateAccessScopeUseCase{repo: repo}
}

func (uc updateAccessScopeUseCase) Execute(ctx context.Context, req *dto.UpdateAccessScopeRequest) (*dto.ExecutionRecordResponse, error) {
	id, err := value_object.NewExecutionRecordID(req.ExecutionRecordID)
	if err != nil {
		return nil, &apperror.ValidationError{Field: "executionRecordID", Message: "invalid execution record ID format"}
	}

	scope, err := value_object.NewAccessScope(req.AccessScope)
	if err != nil {
		return nil, &apperror.ValidationError{Field: "accessScope", Message: err.Error()}
	}

	record, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, &apperror.NotFoundError{ResourceType: "ExecutionRecord", ResourceID: req.ExecutionRecordID}
	}

	record.UpdateAccessScope(scope)

	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return toExecutionRecordResponse(record), nil
}
