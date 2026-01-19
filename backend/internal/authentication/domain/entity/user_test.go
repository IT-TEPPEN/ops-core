package entity

import (
	"testing"
	"time"

	"opscore/backend/internal/authentication/domain/value_object"
)

func TestNewUser(t *testing.T) {
	userID := value_object.NewUserID()

	user := NewUser(
		userID,
		"test@example.com",
		"Test User",
		"https://example.com/picture.jpg",
	)

	if user == nil {
		t.Fatal("expected user but got nil")
	}

	// Verify initial state
	if user.Email() != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", user.Email())
	}
	if user.DisplayName() != "Test User" {
		t.Errorf("expected display name 'Test User', got %s", user.DisplayName())
	}

	// New user should have empty identities (loaded separately)
	if len(user.Identities()) != 0 {
		t.Errorf("expected 0 identities, got %d", len(user.Identities()))
	}

	// Verify domain event
	events := user.GetEvents()
	if len(events) != 1 {
		t.Errorf("expected 1 domain event, got %d", len(events))
	}
	if _, ok := events[0].(*UserCreatedEvent); !ok {
		t.Error("expected UserCreatedEvent")
	}
}

func TestReconstructUser(t *testing.T) {
	userID := value_object.NewUserID()

	user := ReconstructUser(
		userID,
		"test@example.com",
		"Test User",
		"https://example.com/picture.jpg",
		testTime(),
		testTime(),
		testTime(),
	)

	if user == nil {
		t.Fatal("expected user but got nil")
	}

	// Identities should not be loaded yet
	if user.HasIdentitiesLoaded() {
		t.Error("identities should not be loaded until LoadIdentities is called")
	}

	// No domain events should be generated
	events := user.GetEvents()
	if len(events) != 0 {
		t.Errorf("expected 0 domain events for reconstructed user, got %d", len(events))
	}
}

func TestUser_AddIdentity(t *testing.T) {
	user := createTestUser(t)

	// Load the initial identity
	initialIdentity := createTestIdentityForUser(t, user.ID(), "google", "google-user-123")
	user.LoadIdentities([]*Identity{initialIdentity})
	user.ClearEvents()

	initialCount := len(user.Identities())

	newIdentity, err := NewIdentity(
		value_object.NewIdentityID(),
		user.ID(),
		"github",
		"github-user-456",
		"test@github.com",
		"Test User",
		"https://github.com/picture.jpg",
	)
	if err != nil {
		t.Fatalf("failed to create identity: %v", err)
	}

	err = user.AddIdentity(newIdentity)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(user.Identities()) != initialCount+1 {
		t.Errorf("expected %d identities, got %d", initialCount+1, len(user.Identities()))
	}

	// Verify domain event
	events := user.GetEvents()
	foundEvent := false
	for _, event := range events {
		if _, ok := event.(*IdentityAddedEvent); ok {
			foundEvent = true
			break
		}
	}
	if !foundEvent {
		t.Error("expected IdentityAddedEvent")
	}
}

func TestUser_AddIdentity_Duplicate(t *testing.T) {
	user := createTestUser(t)

	// Load an initial identity
	initialIdentity := createTestIdentityForUser(t, user.ID(), "google", "google-user-123")
	user.LoadIdentities([]*Identity{initialIdentity})

	// Try to add the same provider identity again
	duplicateIdentity, err := NewIdentity(
		value_object.NewIdentityID(),
		user.ID(),
		initialIdentity.Provider(),
		initialIdentity.ProviderUserID(),
		"test2@example.com",
		"Test User 2",
		"https://example.com/picture2.jpg",
	)
	if err != nil {
		t.Fatalf("failed to create identity: %v", err)
	}

	err = user.AddIdentity(duplicateIdentity)
	if err == nil {
		t.Error("expected error when adding duplicate provider identity")
	}
}

func TestUser_AddIdentity_FirstIsSetAsPrimary(t *testing.T) {
	user := createTestUser(t)

	identity := createTestIdentityForUser(t, user.ID(), "google", "google-user-123")
	err := user.AddIdentity(identity)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !identity.IsPrimary() {
		t.Error("first identity should be set as primary")
	}
}

func TestUser_RemoveIdentity(t *testing.T) {
	user := createTestUserWithMultipleIdentities(t)
	identities := user.Identities()
	if len(identities) < 2 {
		t.Fatal("test requires at least 2 identities")
	}

	// Remove a non-primary identity
	var toRemove *Identity
	for _, identity := range identities {
		if !identity.IsPrimary() {
			toRemove = identity
			break
		}
	}

	err := user.RemoveIdentity(toRemove.ID())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify identity was removed
	remainingIdentities := user.Identities()
	for _, identity := range remainingIdentities {
		if identity.ID().Equals(toRemove.ID()) {
			t.Error("identity should have been removed")
		}
	}

	// Verify domain event
	events := user.GetEvents()
	foundEvent := false
	for _, event := range events {
		if _, ok := event.(*IdentityRemovedEvent); ok {
			foundEvent = true
			break
		}
	}
	if !foundEvent {
		t.Error("expected IdentityRemovedEvent")
	}
}

func TestUser_RemoveIdentity_LastOne(t *testing.T) {
	user := createTestUser(t)
	identity := createTestIdentityForUser(t, user.ID(), "google", "google-user-123")
	user.LoadIdentities([]*Identity{identity})

	if len(user.Identities()) != 1 {
		t.Fatal("test requires exactly 1 identity")
	}

	err := user.RemoveIdentity(identity.ID())

	if err == nil {
		t.Error("expected error when removing the last identity")
	}
}

func TestUser_RemoveIdentity_NotFound(t *testing.T) {
	user := createTestUser(t)
	identity := createTestIdentityForUser(t, user.ID(), "google", "google-user-123")
	user.LoadIdentities([]*Identity{identity})

	nonExistentID := value_object.NewIdentityID()

	err := user.RemoveIdentity(nonExistentID)
	if err == nil {
		t.Error("expected error when removing non-existent identity")
	}
}

func TestUser_SetPrimaryIdentity(t *testing.T) {
	user := createTestUserWithMultipleIdentities(t)
	identities := user.Identities()

	// Find a non-primary identity
	var newPrimary *Identity
	for _, identity := range identities {
		if !identity.IsPrimary() {
			newPrimary = identity
			break
		}
	}

	if newPrimary == nil {
		t.Fatal("test requires at least one non-primary identity")
	}

	err := user.SetPrimaryIdentity(newPrimary.ID())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify new primary
	if !user.GetPrimaryIdentity().ID().Equals(newPrimary.ID()) {
		t.Error("primary identity should be updated")
	}

	// Verify only one primary exists
	primaryCount := 0
	for _, identity := range user.Identities() {
		if identity.IsPrimary() {
			primaryCount++
		}
	}
	if primaryCount != 1 {
		t.Errorf("expected exactly 1 primary identity, got %d", primaryCount)
	}

	// Verify domain event
	events := user.GetEvents()
	foundEvent := false
	for _, event := range events {
		if _, ok := event.(*PrimaryIdentityChangedEvent); ok {
			foundEvent = true
			break
		}
	}
	if !foundEvent {
		t.Error("expected PrimaryIdentityChangedEvent")
	}
}

func TestUser_SetPrimaryIdentity_NotFound(t *testing.T) {
	user := createTestUser(t)
	identity := createTestIdentityForUser(t, user.ID(), "google", "google-user-123")
	user.LoadIdentities([]*Identity{identity})

	nonExistentID := value_object.NewIdentityID()

	err := user.SetPrimaryIdentity(nonExistentID)
	if err == nil {
		t.Error("expected error when setting non-existent identity as primary")
	}
}

func TestUser_ClearEvents(t *testing.T) {
	user := createTestUser(t)

	// Add an identity to generate an event
	identity := createTestIdentityForUser(t, user.ID(), "google", "google-user-123")
	err := user.AddIdentity(identity)
	if err != nil {
		t.Fatalf("failed to add identity: %v", err)
	}

	events := user.GetEvents()
	if len(events) == 0 {
		t.Fatal("expected domain events to be present")
	}

	user.ClearEvents()

	events = user.GetEvents()
	if len(events) != 0 {
		t.Errorf("expected 0 domain events after clearing, got %d", len(events))
	}
}

// Helper functions

func createTestUser(t *testing.T) *User {
	t.Helper()
	userID := value_object.NewUserID()

	user := NewUser(
		userID,
		"test@example.com",
		"Test User",
		"https://example.com/picture.jpg",
	)

	// Clear initial events for cleaner test assertions
	user.ClearEvents()

	return user
}

func createTestIdentityForUser(t *testing.T, userID value_object.UserID, provider, providerUserID string) *Identity {
	t.Helper()
	identityID := value_object.NewIdentityID()

	identity, err := NewIdentity(
		identityID,
		userID,
		provider,
		providerUserID,
		"test@example.com",
		"Test User",
		"https://example.com/picture.jpg",
	)
	if err != nil {
		t.Fatalf("failed to create identity: %v", err)
	}

	return identity
}

func createTestUserWithMultipleIdentities(t *testing.T) *User {
	t.Helper()
	user := createTestUser(t)

	// Add first identity (will be set as primary)
	identity1 := createTestIdentityForUser(t, user.ID(), "google", "google-user-123")
	err := user.AddIdentity(identity1)
	if err != nil {
		t.Fatalf("failed to add first identity: %v", err)
	}

	// Add second identity
	identity2 := createTestIdentityForUser(t, user.ID(), "github", "github-user-456")
	err = user.AddIdentity(identity2)
	if err != nil {
		t.Fatalf("failed to add second identity: %v", err)
	}

	// Clear events for cleaner test assertions
	user.ClearEvents()

	return user
}

func testTime() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}
