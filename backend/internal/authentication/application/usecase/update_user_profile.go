package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// UpdateUserProfile updates the user's customizable profile information
type UpdateUserProfile interface {
	// Execute updates the user's profile
	// Nil values mean "use the primary identity's value"
	Execute(ctx context.Context, profile dto.UserProfile) (*dto.UserProfile, error)
}
