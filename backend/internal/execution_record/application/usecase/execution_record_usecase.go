package usecase

import (
	"context"

	"opscore/backend/internal/execution_record/application/dto"
	"opscore/backend/internal/execution_record/domain/repository"
)

// ExecutionRecordUsecase is a facade delegating to per-usecase executors.
type ExecutionRecordUsecase struct {
	create          createExecutionRecordUseCase
	get             getExecutionRecordUseCase
	addStep         addStepUseCase
	updateStepNotes updateStepNotesUseCase
	updateNotes     updateNotesUseCase
	updateTitle     updateTitleUseCase
	complete        completeExecutionUseCase
	markAsFailed    markAsFailedUseCase
	updateAccess    updateAccessScopeUseCase
	search          searchExecutionRecordsUseCase
	getByExecutor   getByExecutorIDUseCase
	delete          deleteExecutionRecordUseCase
}

// NewExecutionRecordUsecase wires executors into the facade.
func NewExecutionRecordUsecase(repo repository.ExecutionRecordRepository) *ExecutionRecordUsecase {
	return &ExecutionRecordUsecase{
		create:          newCreateExecutionRecordUseCase(repo),
		get:             newGetExecutionRecordUseCase(repo),
		addStep:         newAddStepUseCase(repo),
		updateStepNotes: newUpdateStepNotesUseCase(repo),
		updateNotes:     newUpdateNotesUseCase(repo),
		updateTitle:     newUpdateTitleUseCase(repo),
		complete:        newCompleteExecutionUseCase(repo),
		markAsFailed:    newMarkAsFailedUseCase(repo),
		updateAccess:    newUpdateAccessScopeUseCase(repo),
		search:          newSearchExecutionRecordsUseCase(repo),
		getByExecutor:   newGetByExecutorIDUseCase(repo),
		delete:          newDeleteExecutionRecordUseCase(repo),
	}
}

func (uc *ExecutionRecordUsecase) CreateExecutionRecord(ctx context.Context, req *dto.CreateExecutionRecordRequest) (*dto.ExecutionRecordResponse, error) {
	return uc.create.Execute(ctx, req)
}

func (uc *ExecutionRecordUsecase) GetExecutionRecord(ctx context.Context, recordID string) (*dto.ExecutionRecordResponse, error) {
	return uc.get.Execute(ctx, recordID)
}

func (uc *ExecutionRecordUsecase) AddStep(ctx context.Context, req *dto.AddStepRequest) (*dto.ExecutionRecordResponse, error) {
	return uc.addStep.Execute(ctx, req)
}

func (uc *ExecutionRecordUsecase) UpdateStepNotes(ctx context.Context, req *dto.UpdateStepNotesRequest) (*dto.ExecutionRecordResponse, error) {
	return uc.updateStepNotes.Execute(ctx, req)
}

func (uc *ExecutionRecordUsecase) UpdateNotes(ctx context.Context, req *dto.UpdateNotesRequest) (*dto.ExecutionRecordResponse, error) {
	return uc.updateNotes.Execute(ctx, req)
}

func (uc *ExecutionRecordUsecase) UpdateTitle(ctx context.Context, req *dto.UpdateTitleRequest) (*dto.ExecutionRecordResponse, error) {
	return uc.updateTitle.Execute(ctx, req)
}

func (uc *ExecutionRecordUsecase) Complete(ctx context.Context, req *dto.CompleteExecutionRequest) (*dto.ExecutionRecordResponse, error) {
	return uc.complete.Execute(ctx, req)
}

func (uc *ExecutionRecordUsecase) MarkAsFailed(ctx context.Context, req *dto.MarkAsFailedRequest) (*dto.ExecutionRecordResponse, error) {
	return uc.markAsFailed.Execute(ctx, req)
}

func (uc *ExecutionRecordUsecase) UpdateAccessScope(ctx context.Context, req *dto.UpdateAccessScopeRequest) (*dto.ExecutionRecordResponse, error) {
	return uc.updateAccess.Execute(ctx, req)
}

func (uc *ExecutionRecordUsecase) SearchExecutionRecords(ctx context.Context, req *dto.SearchExecutionRecordRequest) ([]*dto.ExecutionRecordResponse, error) {
	return uc.search.Execute(ctx, req)
}

func (uc *ExecutionRecordUsecase) GetByExecutorID(ctx context.Context, executorID string) ([]*dto.ExecutionRecordResponse, error) {
	return uc.getByExecutor.Execute(ctx, executorID)
}

func (uc *ExecutionRecordUsecase) DeleteExecutionRecord(ctx context.Context, recordID string) error {
	return uc.delete.Execute(ctx, recordID)
}
