package usecase

import (
	"context"

	"opscore/backend/internal/execution_record/application/dto"
	"opscore/backend/internal/execution_record/domain/repository"
)

// getByExecutorIDUseCase fetches records by executor ID.
type getByExecutorIDUseCase struct {
	repo repository.ExecutionRecordRepository
}

func newGetByExecutorIDUseCase(repo repository.ExecutionRecordRepository) getByExecutorIDUseCase {
	return getByExecutorIDUseCase{repo: repo}
}

func (uc getByExecutorIDUseCase) Execute(ctx context.Context, executorID string) ([]*dto.ExecutionRecordResponse, error) {
	records, err := uc.repo.FindByExecutorID(ctx, executorID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ExecutionRecordResponse, 0, len(records))
	for _, record := range records {
		responses = append(responses, toExecutionRecordResponse(record))
	}

	return responses, nil
}
