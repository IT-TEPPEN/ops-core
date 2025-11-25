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

// UserUseCase defines the interface for user related use cases
type UserUseCase interface {
	// Create creates a new user
	Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	// GetByID retrieves a user by its ID
	GetByID(ctx context.Context, userID string) (*dto.UserResponse, error)
	// GetAll retrieves all users
	GetAll(ctx context.Context) ([]dto.UserResponse, error)
	// Update updates an existing user
	Update(ctx context.Context, userID string, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	// Delete removes a user by its ID
	Delete(ctx context.Context, userID string) error
	// ChangeRole changes the user's role
	ChangeRole(ctx context.Context, userID string, req dto.ChangeRoleRequest) (*dto.UserResponse, error)
	// JoinGroup adds a user to a group
	JoinGroup(ctx context.Context, userID string, req dto.JoinGroupRequest) (*dto.UserResponse, error)
	// LeaveGroup removes a user from a group
	LeaveGroup(ctx context.Context, userID string, req dto.LeaveGroupRequest) (*dto.UserResponse, error)
}

// userUseCase implements the UserUseCase interface
type userUseCase struct {
	userRepo  repository.UserRepository
	groupRepo repository.GroupRepository
}

// NewUserUseCase creates a new instance of userUseCase
func NewUserUseCase(userRepo repository.UserRepository, groupRepo repository.GroupRepository) UserUseCase {
	return &userUseCase{
		userRepo:  userRepo,
		groupRepo: groupRepo,
	}
}

// Create creates a new user
func (uc *userUseCase) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Validate and create value objects
	userID, err := value_object.NewUserID(uuid.NewString())
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "id", Message: err.Error()},
		})
	}

	email, err := value_object.NewEmail(req.Email)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "email", Message: err.Error()},
		})
	}

	role, err := value_object.NewRole(req.Role)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "role", Message: err.Error()},
		})
	}

	// Check if email already exists
	existingUser, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing user: %w", err)
	}
	if existingUser != nil {
		return nil, apperror.NewConflictError("User", req.Email, "email already registered", nil)
	}

	// Create user entity
	user, err := entity.NewUser(userID, req.Name, email, role)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "name", Message: err.Error()},
		})
	}

	// Persist user
	err = uc.userRepo.Save(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

// GetByID retrieves a user by its ID
func (uc *userUseCase) GetByID(ctx context.Context, userID string) (*dto.UserResponse, error) {
	id, err := value_object.NewUserID(userID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "id", Message: err.Error()},
		})
	}

	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("User", userID, nil)
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

// GetAll retrieves all users
func (uc *userUseCase) GetAll(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := uc.userRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	return dto.ToUserResponseList(users), nil
}

// Update updates an existing user
func (uc *userUseCase) Update(ctx context.Context, userID string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	id, err := value_object.NewUserID(userID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "id", Message: err.Error()},
		})
	}

	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("User", userID, nil)
	}

	email, err := value_object.NewEmail(req.Email)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "email", Message: err.Error()},
		})
	}

	// Check if email changed and is already in use
	if !user.Email().Equals(email) {
		existingUser, err := uc.userRepo.FindByEmail(ctx, email)
		if err != nil {
			return nil, fmt.Errorf("failed to check for existing user: %w", err)
		}
		if existingUser != nil {
			return nil, apperror.NewConflictError("User", req.Email, "email already registered", nil)
		}
	}

	err = user.UpdateProfile(req.Name, email)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "name", Message: err.Error()},
		})
	}

	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

// Delete removes a user by its ID
func (uc *userUseCase) Delete(ctx context.Context, userID string) error {
	id, err := value_object.NewUserID(userID)
	if err != nil {
		return apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "id", Message: err.Error()},
		})
	}

	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return apperror.NewNotFoundError("User", userID, nil)
	}

	err = uc.userRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ChangeRole changes the user's role
func (uc *userUseCase) ChangeRole(ctx context.Context, userID string, req dto.ChangeRoleRequest) (*dto.UserResponse, error) {
	id, err := value_object.NewUserID(userID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "id", Message: err.Error()},
		})
	}

	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("User", userID, nil)
	}

	role, err := value_object.NewRole(req.Role)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "role", Message: err.Error()},
		})
	}

	err = user.ChangeRole(role)
	if err != nil {
		return nil, fmt.Errorf("failed to change role: %w", err)
	}

	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

// JoinGroup adds a user to a group
func (uc *userUseCase) JoinGroup(ctx context.Context, userID string, req dto.JoinGroupRequest) (*dto.UserResponse, error) {
	id, err := value_object.NewUserID(userID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "id", Message: err.Error()},
		})
	}

	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("User", userID, nil)
	}

	groupID, err := value_object.NewGroupID(req.GroupID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "group_id", Message: err.Error()},
		})
	}

	// Verify group exists
	group, err := uc.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve group: %w", err)
	}
	if group == nil {
		return nil, apperror.NewNotFoundError("Group", req.GroupID, nil)
	}

	err = user.JoinGroup(groupID)
	if err != nil {
		return nil, apperror.NewConflictError("User", userID, err.Error(), nil)
	}

	// Also add user to group
	err = group.AddMember(id)
	if err != nil {
		// Already a member - this is expected since user.JoinGroup succeeded
	}

	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	err = uc.groupRepo.Update(ctx, group)
	if err != nil {
		return nil, fmt.Errorf("failed to update group: %w", err)
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

// LeaveGroup removes a user from a group
func (uc *userUseCase) LeaveGroup(ctx context.Context, userID string, req dto.LeaveGroupRequest) (*dto.UserResponse, error) {
	id, err := value_object.NewUserID(userID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "id", Message: err.Error()},
		})
	}

	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("User", userID, nil)
	}

	groupID, err := value_object.NewGroupID(req.GroupID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "group_id", Message: err.Error()},
		})
	}

	// Verify group exists
	group, err := uc.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve group: %w", err)
	}
	if group == nil {
		return nil, apperror.NewNotFoundError("Group", req.GroupID, nil)
	}

	err = user.LeaveGroup(groupID)
	if err != nil {
		return nil, apperror.NewConflictError("User", userID, err.Error(), nil)
	}

	// Also remove user from group
	err = group.RemoveMember(id)
	if err != nil {
		// Not a member - this is expected since user.LeaveGroup succeeded
	}

	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	err = uc.groupRepo.Update(ctx, group)
	if err != nil {
		return nil, fmt.Errorf("failed to update group: %w", err)
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}
