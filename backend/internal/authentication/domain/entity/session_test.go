package entity

import (
	"testing"
	"time"

	"opscore/backend/internal/authentication/domain/value_object"
)

func TestNewSession(t *testing.T) {
	userID := value_object.NewUserID()

	tests := []struct {
		name       string
		duration   time.Duration
		rememberMe bool
	}{
		{
			name:       "short session",
			duration:   7 * 24 * time.Hour,
			rememberMe: false,
		},
		{
			name:       "long session",
			duration:   30 * 24 * time.Hour,
			rememberMe: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, token, err := NewSession(
				userID,
				tt.duration,
				tt.rememberMe,
			)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if session == nil {
				t.Fatal("expected session but got nil")
			}
			if token == "" {
				t.Fatal("expected token but got empty string")
			}

			// Verify initial state
			if session.IsRevoked() {
				t.Error("new session should not be revoked")
			}
			if session.TokenHash() == "" {
				t.Error("session should have a token hash")
			}
			if token == session.TokenHash() {
				t.Error("token and token hash should be different")
			}
			if session.RememberMe() != tt.rememberMe {
				t.Errorf("expected RememberMe=%v, got %v", tt.rememberMe, session.RememberMe())
			}

			// Verify expiration time
			expectedExpiry := session.CreatedAt().Add(tt.duration)
			if !session.ExpiresAt().Equal(expectedExpiry) {
				t.Errorf("expected expiry %v, got %v", expectedExpiry, session.ExpiresAt())
			}
		})
	}
}

func TestReconstructSession(t *testing.T) {
	sessionID := value_object.NewSessionID()
	userID := value_object.NewUserID()

	session := ReconstructSession(
		sessionID,
		userID,
		"token-hash-123",
		"jti-123",
		testTime(),
		testTime(),
		nil,
		nil,
		false,
		false,
	)

	if session == nil {
		t.Fatal("expected session but got nil")
	}

	if session.TokenHash() != "token-hash-123" {
		t.Errorf("expected token hash 'token-hash-123', got %s", session.TokenHash())
	}
}

func TestSession_IsValid(t *testing.T) {
	tests := []struct {
		name           string
		setupSession   func() *Session
		expectedResult bool
	}{
		{
			name: "valid session",
			setupSession: func() *Session {
				session, _, _ := NewSession(
					value_object.NewUserID(),
					7*24*time.Hour,
					false,
				)
				return session
			},
			expectedResult: true,
		},
		{
			name: "expired session",
			setupSession: func() *Session {
				sessionID := value_object.NewSessionID()
				userID := value_object.NewUserID()
				pastTime := time.Now().Add(-8 * 24 * time.Hour) // 8 days ago

				return ReconstructSession(
					sessionID,
					userID,
					"token-hash",
					"jti-123",
					pastTime,
					pastTime.Add(-1*time.Hour),
					nil,
					nil,
					false,
					false,
				)
			},
			expectedResult: false,
		},
		{
			name: "revoked session",
			setupSession: func() *Session {
				session, _, _ := NewSession(
					value_object.NewUserID(),
					7*24*time.Hour,
					false,
				)
				session.Revoke()
				return session
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := tt.setupSession()
			result := session.IsValid()

			if result != tt.expectedResult {
				t.Errorf("expected IsValid() = %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestSession_UpdateLastUsed(t *testing.T) {
	session := createTestSession(t)

	if session.LastUsedAt() != nil {
		t.Error("LastUsedAt should be nil initially")
	}

	// Small delay to ensure time difference
	time.Sleep(1 * time.Millisecond)

	session.UpdateLastUsed()

	if session.LastUsedAt() == nil {
		t.Error("LastUsedAt should be set after UpdateLastUsed")
	}
	if !session.LastUsedAt().After(session.CreatedAt()) {
		t.Error("LastUsedAt should be after CreatedAt")
	}
}

func TestSession_Revoke(t *testing.T) {
	session := createTestSession(t)

	if session.IsRevoked() {
		t.Error("session should not be revoked initially")
	}
	if session.RevokedAt() != nil {
		t.Error("RevokedAt should be nil initially")
	}

	session.Revoke()

	if !session.IsRevoked() {
		t.Error("session should be revoked after Revoke()")
	}
	if session.RevokedAt() == nil {
		t.Error("RevokedAt should be set after Revoke()")
	}
	if session.IsValid() {
		t.Error("revoked session should not be valid")
	}
}

// Helper functions

func createTestSession(t *testing.T) *Session {
	t.Helper()
	session, _, err := NewSession(
		value_object.NewUserID(),
		7*24*time.Hour,
		false,
	)
	if err != nil {
		t.Fatalf("failed to create test session: %v", err)
	}
	return session
}
