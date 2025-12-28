package domain

import (
	"testing"
	"time"
)

func TestNewOAuthConnection(t *testing.T) {
	t.Run("valid connection", func(t *testing.T) {
		expiresAt := time.Now().Add(time.Hour)
		conn, err := NewOAuthConnection(
			"test-id",
			"user-id",
			ProviderGitHub,
			"12345",
			"testuser",
			"access-token",
			"refresh-token",
			&expiresAt,
			[]string{"repo", "read:user"},
		)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if conn.ID() != "test-id" {
			t.Errorf("expected id 'test-id', got '%s'", conn.ID())
		}
		if conn.UserID() != "user-id" {
			t.Errorf("expected user_id 'user-id', got '%s'", conn.UserID())
		}
		if conn.Provider() != ProviderGitHub {
			t.Errorf("expected provider 'github', got '%s'", conn.Provider())
		}
		if conn.ProviderUserID() != "12345" {
			t.Errorf("expected provider_user_id '12345', got '%s'", conn.ProviderUserID())
		}
		if conn.ProviderUsername() != "testuser" {
			t.Errorf("expected provider_username 'testuser', got '%s'", conn.ProviderUsername())
		}
		if conn.AccessToken() != "access-token" {
			t.Errorf("expected access_token 'access-token', got '%s'", conn.AccessToken())
		}
		if conn.RefreshToken() != "refresh-token" {
			t.Errorf("expected refresh_token 'refresh-token', got '%s'", conn.RefreshToken())
		}
		if len(conn.Scopes()) != 2 {
			t.Errorf("expected 2 scopes, got %d", len(conn.Scopes()))
		}
	})

	t.Run("missing id", func(t *testing.T) {
		_, err := NewOAuthConnection("", "user-id", ProviderGitHub, "12345", "testuser", "token", "", nil, nil)
		if err == nil {
			t.Error("expected error for missing id")
		}
		if err.Error() != "id is required" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})

	t.Run("missing user_id", func(t *testing.T) {
		_, err := NewOAuthConnection("id", "", ProviderGitHub, "12345", "testuser", "token", "", nil, nil)
		if err == nil {
			t.Error("expected error for missing user_id")
		}
		if err.Error() != "user_id is required" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})

	t.Run("invalid provider", func(t *testing.T) {
		_, err := NewOAuthConnection("id", "user-id", "invalid", "12345", "testuser", "token", "", nil, nil)
		if err == nil {
			t.Error("expected error for invalid provider")
		}
		if err.Error() != "invalid provider" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})

	t.Run("missing access_token", func(t *testing.T) {
		_, err := NewOAuthConnection("id", "user-id", ProviderGitHub, "12345", "testuser", "", "", nil, nil)
		if err == nil {
			t.Error("expected error for missing access_token")
		}
		if err.Error() != "access_token is required" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})
}

func TestOAuthConnection_IsTokenExpired(t *testing.T) {
	t.Run("token not expired", func(t *testing.T) {
		expiresAt := time.Now().Add(time.Hour)
		conn, _ := NewOAuthConnection("id", "user-id", ProviderGitHub, "12345", "testuser", "token", "", &expiresAt, nil)
		if conn.IsTokenExpired() {
			t.Error("expected token to not be expired")
		}
	})

	t.Run("token expired", func(t *testing.T) {
		expiresAt := time.Now().Add(-time.Hour)
		conn, _ := NewOAuthConnection("id", "user-id", ProviderGitHub, "12345", "testuser", "token", "", &expiresAt, nil)
		if !conn.IsTokenExpired() {
			t.Error("expected token to be expired")
		}
	})

	t.Run("token expires soon (within 5 minute buffer)", func(t *testing.T) {
		expiresAt := time.Now().Add(3 * time.Minute) // Less than 5 minute buffer
		conn, _ := NewOAuthConnection("id", "user-id", ProviderGitHub, "12345", "testuser", "token", "", &expiresAt, nil)
		if !conn.IsTokenExpired() {
			t.Error("expected token to be considered expired (within buffer)")
		}
	})

	t.Run("no expiration set", func(t *testing.T) {
		conn, _ := NewOAuthConnection("id", "user-id", ProviderGitHub, "12345", "testuser", "token", "", nil, nil)
		if conn.IsTokenExpired() {
			t.Error("expected token without expiration to not be expired")
		}
	})
}

func TestOAuthConnection_UpdateTokens(t *testing.T) {
	t.Run("update tokens successfully", func(t *testing.T) {
		expiresAt := time.Now().Add(time.Hour)
		conn, _ := NewOAuthConnection("id", "user-id", ProviderGitHub, "12345", "testuser", "old-token", "old-refresh", &expiresAt, nil)

		newExpiresAt := time.Now().Add(2 * time.Hour)
		err := conn.UpdateTokens("new-token", "new-refresh", &newExpiresAt)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if conn.AccessToken() != "new-token" {
			t.Errorf("expected access_token 'new-token', got '%s'", conn.AccessToken())
		}
		if conn.RefreshToken() != "new-refresh" {
			t.Errorf("expected refresh_token 'new-refresh', got '%s'", conn.RefreshToken())
		}
	})

	t.Run("update with empty access token fails", func(t *testing.T) {
		conn, _ := NewOAuthConnection("id", "user-id", ProviderGitHub, "12345", "testuser", "token", "", nil, nil)
		err := conn.UpdateTokens("", "", nil)
		if err == nil {
			t.Error("expected error for empty access token")
		}
	})

	t.Run("update keeps old refresh token if new one is empty", func(t *testing.T) {
		conn, _ := NewOAuthConnection("id", "user-id", ProviderGitHub, "12345", "testuser", "token", "old-refresh", nil, nil)
		err := conn.UpdateTokens("new-token", "", nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if conn.RefreshToken() != "old-refresh" {
			t.Errorf("expected refresh_token to remain 'old-refresh', got '%s'", conn.RefreshToken())
		}
	})
}

func TestProvider_IsValid(t *testing.T) {
	tests := []struct {
		provider Provider
		valid    bool
	}{
		{ProviderGitHub, true},
		{ProviderGitLab, true},
		{ProviderGitLabSelfHosted, true},
		{Provider("bitbucket"), false},
		{Provider(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.provider), func(t *testing.T) {
			if got := tt.provider.IsValid(); got != tt.valid {
				t.Errorf("Provider(%s).IsValid() = %v, want %v", tt.provider, got, tt.valid)
			}
		})
	}
}
