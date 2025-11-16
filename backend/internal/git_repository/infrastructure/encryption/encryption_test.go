package encryption

import (
	"crypto/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEncryptor(t *testing.T) {
	tests := []struct {
		name    string
		key     []byte
		wantErr bool
		errType error
	}{
		{
			name:    "valid 32-byte key",
			key:     make([]byte, 32),
			wantErr: false,
		},
		{
			name:    "invalid 16-byte key",
			key:     make([]byte, 16),
			wantErr: true,
			errType: ErrInvalidKey,
		},
		{
			name:    "invalid 24-byte key",
			key:     make([]byte, 24),
			wantErr: true,
			errType: ErrInvalidKey,
		},
		{
			name:    "invalid empty key",
			key:     []byte{},
			wantErr: true,
			errType: ErrInvalidKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := NewEncryptor(tt.key)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.errType)
				assert.Nil(t, enc)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, enc)
			}
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	// Generate a random 32-byte key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	enc, err := NewEncryptor(key)
	require.NoError(t, err)

	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple text",
			plaintext: "hello world",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "access token",
			plaintext: "ghp_1234567890abcdefghijklmnopqrstuvwxyz",
		},
		{
			name:      "long text",
			plaintext: strings.Repeat("a", 1000),
		},
		{
			name:      "special characters",
			plaintext: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
		{
			name:      "unicode text",
			plaintext: "こんにちは世界",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			ciphertext, err := enc.Encrypt(tt.plaintext)
			assert.NoError(t, err)

			if tt.plaintext == "" {
				assert.Equal(t, "", ciphertext)
			} else {
				assert.NotEmpty(t, ciphertext)
				assert.NotEqual(t, tt.plaintext, ciphertext)
			}

			// Decrypt
			decrypted, err := enc.Decrypt(ciphertext)
			assert.NoError(t, err)
			assert.Equal(t, tt.plaintext, decrypted)
		})
	}
}

func TestEncrypt_DifferentCiphertexts(t *testing.T) {
	// Each encryption should produce a different ciphertext due to random nonce
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	enc, err := NewEncryptor(key)
	require.NoError(t, err)

	plaintext := "test message"

	ciphertext1, err := enc.Encrypt(plaintext)
	require.NoError(t, err)

	ciphertext2, err := enc.Encrypt(plaintext)
	require.NoError(t, err)

	// Different ciphertexts due to different nonces
	assert.NotEqual(t, ciphertext1, ciphertext2)

	// But both should decrypt to the same plaintext
	decrypted1, err := enc.Decrypt(ciphertext1)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted1)

	decrypted2, err := enc.Decrypt(ciphertext2)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted2)
}

func TestDecrypt_InvalidCiphertext(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	enc, err := NewEncryptor(key)
	require.NoError(t, err)

	tests := []struct {
		name       string
		ciphertext string
		wantErr    bool
	}{
		{
			name:       "invalid base64",
			ciphertext: "not-base64!!!",
			wantErr:    true,
		},
		{
			name:       "too short ciphertext",
			ciphertext: "YWJj", // "abc" in base64, which is too short
			wantErr:    true,
		},
		{
			name:       "valid base64 but invalid ciphertext",
			ciphertext: "dGVzdCBtZXNzYWdl", // "test message" in base64, not encrypted
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := enc.Decrypt(tt.ciphertext)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	// Encrypt with one key
	key1 := make([]byte, 32)
	_, err := rand.Read(key1)
	require.NoError(t, err)

	enc1, err := NewEncryptor(key1)
	require.NoError(t, err)

	plaintext := "secret message"
	ciphertext, err := enc1.Encrypt(plaintext)
	require.NoError(t, err)

	// Try to decrypt with a different key
	key2 := make([]byte, 32)
	_, err = rand.Read(key2)
	require.NoError(t, err)

	enc2, err := NewEncryptor(key2)
	require.NoError(t, err)

	_, err = enc2.Decrypt(ciphertext)
	assert.Error(t, err)
}
