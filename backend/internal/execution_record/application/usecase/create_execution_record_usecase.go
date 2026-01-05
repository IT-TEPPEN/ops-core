package usecase

import (
	"context"

	docvo "opscore/backend/internal/document/domain/value_object"
	"opscore/backend/internal/execution_record/application/dto"
	apperror "opscore/backend/internal/execution_record/application/error"
	"opscore/backend/internal/execution_record/domain/entity"
	"opscore/backend/internal/execution_record/domain/repository"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// createExecutionRecordUseCase handles creation of execution records.
type createExecutionRecordUseCase struct {
	repo repository.ExecutionRecordRepository
}

func newCreateExecutionRecordUseCase(repo repository.ExecutionRecordRepository) createExecutionRecordUseCase {
	return createExecutionRecordUseCase{repo: repo}
}

func (uc createExecutionRecordUseCase) Execute(ctx context.Context, req *dto.CreateExecutionRecordRequest) (*dto.ExecutionRecordResponse, error) {
	id := value_object.GenerateExecutionRecordID()

	documentID, err := docvo.NewDocumentID(req.DocumentID)
	if err != nil {
		return nil, &apperror.ValidationError{Field: "documentID", Message: "invalid document ID format"}
	}

	versionID, err := docvo.NewVersionID(req.DocumentVersionID)
	if err != nil {
		return nil, &apperror.ValidationError{Field: "documentVersionID", Message: "invalid version ID format"}
	}

	variableValues := make([]value_object.VariableValue, 0, len(req.VariableValues))
	for _, vv := range req.VariableValues {
		v, convErr := value_object.NewVariableValue(vv.Name, vv.Value)
		if convErr != nil {
			return nil, &apperror.ValidationError{Field: "variableValues", Message: "invalid variable value: " + convErr.Error()}
		}
		variableValues = append(variableValues, v)
	}

	record, err := entity.NewExecutionRecord(id, documentID, versionID, req.ExecutorID, req.Title, variableValues)
	if err != nil {
		return nil, &apperror.ValidationError{Field: "executionRecord", Message: err.Error()}
	}

	if err := uc.repo.Save(ctx, record); err != nil {
		return nil, err
	}

	return toExecutionRecordResponse(record), nil
}
