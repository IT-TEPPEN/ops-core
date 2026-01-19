package entity

import (
	"testing"

	"opscore/backend/internal/authentication/domain/value_object"
)

func TestNewIdentity(t *testing.T) {
	id := value_object.NewIdentityID()
	userID := value_object.NewUserID()

	tests := []struct {
		name           string
		providerUserID string
		wantError      bool
	}{
		{
			name:           "valid identity",
			providerUserID: "google-user-123",
			wantError:      false,
		},
		{
			name:           "empty provider user ID",
			providerUserID: "",
			wantError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity, err := NewIdentity(
				id,
				userID,
				"google",
				tt.providerUserID,
				"test@example.com",
				"Test User",
				"https://example.com/picture.jpg",
			)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				if identity != nil {
					t.Error("expected nil identity on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if identity == nil {
					t.Fatal("expected identity but got nil")
				}

				// Verify initial state
				if identity.IsPrimary() {
					t.Error("new identity should not be primary by default")
				}
				if identity.Provider() != "google" {
					t.Errorf("expected provider 'google', got %s", identity.Provider())
				}
				if identity.ProviderUserID() != tt.providerUserID {
					t.Errorf("expected provider user ID %s, got %s", tt.providerUserID, identity.ProviderUserID())
				}
			}
		})
	}
}

func TestReconstructIdentity(t *testing.T) {
	id := value_object.NewIdentityID()
	userID := value_object.NewUserID()

	// ReconstructIdentity should not validate
	identity := ReconstructIdentity(
		id,
		userID,
		"google",
		"", // empty provider user ID is allowed in reconstruction
		"test@example.com",
		"Test User",
		"https://example.com/picture.jpg",
		true,
		testTime(),
		testTime(),
		testTime(),
	)

	if identity == nil {
		t.Fatal("expected identity but got nil")
	}

	if !identity.IsPrimary() {
		t.Error("reconstructed identity should respect isPrimary flag")
	}
}

func TestIdentity_UpdateLastUsed(t *testing.T) {
	identity := createTestIdentity(t)
	initialLastUsedAt := identity.LastUsedAt()
	initialUpdatedAt := identity.UpdatedAt()

	identity.UpdateLastUsed()

	if !identity.LastUsedAt().After(initialLastUsedAt) {
		t.Error("LastUsedAt should be updated")
	}
	if !identity.UpdatedAt().After(initialUpdatedAt) {
		t.Error("UpdatedAt should be updated")
	}
}

func TestIdentity_UpdateInfo(t *testing.T) {
	identity := createTestIdentity(t)

	newEmail := "new@example.com"
	newName := "New Name"
	newPicture := "https://example.com/new.jpg"

	identity.UpdateInfo(newEmail, newName, newPicture)

	if identity.Email() != newEmail {
		t.Errorf("expected email %s, got %s", newEmail, identity.Email())
	}
	if identity.Name() != newName {
		t.Errorf("expected name %s, got %s", newName, identity.Name())
	}
	if identity.PictureURL() != newPicture {
		t.Errorf("expected picture URL %s, got %s", newPicture, identity.PictureURL())
	}
}

func TestIdentity_IsSameProvider(t *testing.T) {
	identity := createTestIdentity(t)

	tests := []struct {
		name           string
		provider       string
		providerUserID string
		expected       bool
	}{
		{
			name:           "same provider and user ID",
			provider:       "google",
			providerUserID: "google-user-123",
			expected:       true,
		},
		{
			name:           "different provider",
			provider:       "github",
			providerUserID: "google-user-123",
			expected:       false,
		},
		{
			name:           "different user ID",
			provider:       "google",
			providerUserID: "google-user-456",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := identity.IsSameProvider(tt.provider, tt.providerUserID)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIdentity_SetAsPrimary(t *testing.T) {
	identity := createTestIdentity(t)

	if identity.IsPrimary() {
		t.Error("identity should not be primary initially")
	}

	identity.SetAsPrimary()

	if !identity.IsPrimary() {
		t.Error("identity should be primary after SetAsPrimary")
	}
}

func TestIdentity_UnsetPrimary(t *testing.T) {
	identity := createTestIdentity(t)
	identity.SetAsPrimary()

	if !identity.IsPrimary() {
		t.Error("identity should be primary")
	}

	identity.UnsetPrimary()

	if identity.IsPrimary() {
		t.Error("identity should not be primary after UnsetPrimary")
	}
}

// Helper functions

func createTestIdentity(t *testing.T) *Identity {
	t.Helper()
	identity, err := NewIdentity(
		value_object.NewIdentityID(),
		value_object.NewUserID(),
		"google",
		"google-user-123",
		"test@example.com",
		"Test User",
		"https://example.com/picture.jpg",
	)
	if err != nil {
		t.Fatalf("failed to create test identity: %v", err)
	}
	return identity
}
