package usecase

import (
	"context"

	docvo "opscore/backend/internal/document/domain/value_object"
	"opscore/backend/internal/execution_record/application/dto"
	apperror "opscore/backend/internal/execution_record/application/error"
	"opscore/backend/internal/execution_record/domain/repository"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// searchExecutionRecordsUseCase searches for execution records by criteria.
type searchExecutionRecordsUseCase struct {
	repo repository.ExecutionRecordRepository
}

func newSearchExecutionRecordsUseCase(repo repository.ExecutionRecordRepository) searchExecutionRecordsUseCase {
	return searchExecutionRecordsUseCase{repo: repo}
}

func (uc searchExecutionRecordsUseCase) Execute(ctx context.Context, req *dto.SearchExecutionRecordRequest) ([]*dto.ExecutionRecordResponse, error) {
	var documentID *docvo.DocumentID
	if req.DocumentID != nil {
		id, err := docvo.NewDocumentID(*req.DocumentID)
		if err != nil {
			return nil, &apperror.ValidationError{Field: "documentID", Message: "invalid document ID format"}
		}
		documentID = &id
	}

	var status *value_object.ExecutionStatus
	if req.Status != nil {
		st, err := value_object.NewExecutionStatus(*req.Status)
		if err != nil {
			return nil, &apperror.ValidationError{Field: "status", Message: err.Error()}
		}
		status = &st
	}

	variableFilters := make([]repository.VariableFilter, 0, len(req.VariableFilters))
	for _, vf := range req.VariableFilters {
		variableFilters = append(variableFilters, repository.VariableFilter{Name: vf.Name, Value: vf.Value})
	}

	criteria := repository.SearchCriteria{
		ExecutorID:      req.ExecutorID,
		DocumentID:      documentID,
		Status:          status,
		StartedFrom:     req.StartedFrom,
		StartedTo:       req.StartedTo,
		VariableFilters: variableFilters,
	}

	records, err := uc.repo.Search(ctx, criteria)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ExecutionRecordResponse, 0, len(records))
	for _, record := range records {
		responses = append(responses, toExecutionRecordResponse(record))
	}

	return responses, nil
}
