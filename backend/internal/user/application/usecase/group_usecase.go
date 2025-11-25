package usecase

import (
	"context"
	"fmt"

	apperror "opscore/backend/internal/user/application/error"
	"opscore/backend/internal/user/application/dto"
	"opscore/backend/internal/user/domain/entity"
	"opscore/backend/internal/user/domain/repository"
	"opscore/backend/internal/user/domain/value_object"

	"github.com/google/uuid"
)

// GroupUseCase defines the interface for group related use cases
type GroupUseCase interface {
	// Create creates a new group
	Create(ctx context.Context, req dto.CreateGroupRequest) (*dto.GroupResponse, error)
	// GetByID retrieves a group by its ID
	GetByID(ctx context.Context, groupID string) (*dto.GroupResponse, error)
	// GetAll retrieves all groups
	GetAll(ctx context.Context) ([]dto.GroupResponse, error)
	// Update updates an existing group
	Update(ctx context.Context, groupID string, req dto.UpdateGroupRequest) (*dto.GroupResponse, error)
	// Delete removes a group by its ID
	Delete(ctx context.Context, groupID string) error
	// AddMember adds a user to a group
	AddMember(ctx context.Context, groupID string, req dto.AddMemberRequest) (*dto.GroupResponse, error)
	// RemoveMember removes a user from a group
	RemoveMember(ctx context.Context, groupID string, req dto.RemoveMemberRequest) (*dto.GroupResponse, error)
	// GetByMemberID retrieves all groups that contain the specified user as a member
	GetByMemberID(ctx context.Context, userID string) ([]dto.GroupResponse, error)
}

// groupUseCase implements the GroupUseCase interface
type groupUseCase struct {
	groupRepo repository.GroupRepository
	userRepo  repository.UserRepository
}

// NewGroupUseCase creates a new instance of groupUseCase
func NewGroupUseCase(groupRepo repository.GroupRepository, userRepo repository.UserRepository) GroupUseCase {
	return &groupUseCase{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// Create creates a new group
func (uc *groupUseCase) Create(ctx context.Context, req dto.CreateGroupRequest) (*dto.GroupResponse, error) {
	groupID, err := value_object.NewGroupID(uuid.NewString())
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "id", Message: err.Error()},
		})
	}

	group, err := entity.NewGroup(groupID, req.Name, req.Description)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "name", Message: err.Error()},
		})
	}

	err = uc.groupRepo.Save(ctx, group)
	if err != nil {
		return nil, fmt.Errorf("failed to save group: %w", err)
	}

	response := dto.ToGroupResponse(group)
	return &response, nil
}

// GetByID retrieves a group by its ID
func (uc *groupUseCase) GetByID(ctx context.Context, groupID string) (*dto.GroupResponse, error) {
	id, err := value_object.NewGroupID(groupID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "id", Message: err.Error()},
		})
	}

	group, err := uc.groupRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve group: %w", err)
	}
	if group == nil {
		return nil, apperror.NewNotFoundError("Group", groupID, nil)
	}

	response := dto.ToGroupResponse(group)
	return &response, nil
}

// GetAll retrieves all groups
func (uc *groupUseCase) GetAll(ctx context.Context) ([]dto.GroupResponse, error) {
	groups, err := uc.groupRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve groups: %w", err)
	}

	return dto.ToGroupResponseList(groups), nil
}

// Update updates an existing group
func (uc *groupUseCase) Update(ctx context.Context, groupID string, req dto.UpdateGroupRequest) (*dto.GroupResponse, error) {
	id, err := value_object.NewGroupID(groupID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "id", Message: err.Error()},
		})
	}

	group, err := uc.groupRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve group: %w", err)
	}
	if group == nil {
		return nil, apperror.NewNotFoundError("Group", groupID, nil)
	}

	err = group.UpdateInfo(req.Name, req.Description)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "name", Message: err.Error()},
		})
	}

	err = uc.groupRepo.Update(ctx, group)
	if err != nil {
		return nil, fmt.Errorf("failed to update group: %w", err)
	}

	response := dto.ToGroupResponse(group)
	return &response, nil
}

// Delete removes a group by its ID
func (uc *groupUseCase) Delete(ctx context.Context, groupID string) error {
	id, err := value_object.NewGroupID(groupID)
	if err != nil {
		return apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "id", Message: err.Error()},
		})
	}

	group, err := uc.groupRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to retrieve group: %w", err)
	}
	if group == nil {
		return apperror.NewNotFoundError("Group", groupID, nil)
	}

	err = uc.groupRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	return nil
}

// AddMember adds a user to a group
func (uc *groupUseCase) AddMember(ctx context.Context, groupID string, req dto.AddMemberRequest) (*dto.GroupResponse, error) {
	gid, err := value_object.NewGroupID(groupID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "id", Message: err.Error()},
		})
	}

	group, err := uc.groupRepo.FindByID(ctx, gid)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve group: %w", err)
	}
	if group == nil {
		return nil, apperror.NewNotFoundError("Group", groupID, nil)
	}

	userID, err := value_object.NewUserID(req.UserID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "user_id", Message: err.Error()},
		})
	}

	// Verify user exists
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("User", req.UserID, nil)
	}

	err = group.AddMember(userID)
	if err != nil {
		return nil, apperror.NewConflictError("Group", groupID, err.Error(), nil)
	}

	// Also add group to user
	err = user.JoinGroup(gid)
	if err != nil {
		// Already in group - this is expected since group.AddMember succeeded
	}

	err = uc.groupRepo.Update(ctx, group)
	if err != nil {
		return nil, fmt.Errorf("failed to update group: %w", err)
	}

	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := dto.ToGroupResponse(group)
	return &response, nil
}

// RemoveMember removes a user from a group
func (uc *groupUseCase) RemoveMember(ctx context.Context, groupID string, req dto.RemoveMemberRequest) (*dto.GroupResponse, error) {
	gid, err := value_object.NewGroupID(groupID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "id", Message: err.Error()},
		})
	}

	group, err := uc.groupRepo.FindByID(ctx, gid)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve group: %w", err)
	}
	if group == nil {
		return nil, apperror.NewNotFoundError("Group", groupID, nil)
	}

	userID, err := value_object.NewUserID(req.UserID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "user_id", Message: err.Error()},
		})
	}

	// Verify user exists
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("User", req.UserID, nil)
	}

	err = group.RemoveMember(userID)
	if err != nil {
		return nil, apperror.NewConflictError("Group", groupID, err.Error(), nil)
	}

	// Also remove group from user
	err = user.LeaveGroup(gid)
	if err != nil {
		// Not in group - this is expected since group.RemoveMember succeeded
	}

	err = uc.groupRepo.Update(ctx, group)
	if err != nil {
		return nil, fmt.Errorf("failed to update group: %w", err)
	}

	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := dto.ToGroupResponse(group)
	return &response, nil
}

// GetByMemberID retrieves all groups that contain the specified user as a member
func (uc *groupUseCase) GetByMemberID(ctx context.Context, userID string) ([]dto.GroupResponse, error) {
	id, err := value_object.NewUserID(userID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "user_id", Message: err.Error()},
		})
	}

	// Verify user exists
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("User", userID, nil)
	}

	groups, err := uc.groupRepo.FindByMemberID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve groups: %w", err)
	}

	return dto.ToGroupResponseList(groups), nil
}
