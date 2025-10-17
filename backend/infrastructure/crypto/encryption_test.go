package crypto

import (
	"crypto/rand"
	"testing"
)

func TestNewAESEncryptor(t *testing.T) {
	tests := []struct {
		name    string
		keySize int
		wantErr bool
	}{
		{
			name:    "valid 32-byte key",
			keySize: 32,
			wantErr: false,
		},
		{
			name:    "invalid 16-byte key",
			keySize: 16,
			wantErr: true,
		},
		{
			name:    "invalid 24-byte key",
			keySize: 24,
			wantErr: true,
		},
		{
			name:    "invalid 0-byte key",
			keySize: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := make([]byte, tt.keySize)
			_, err := NewAESEncryptor(key)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAESEncryptor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAESEncryptor_EncryptDecrypt(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("failed to generate random key: %v", err)
	}

	encryptor, err := NewAESEncryptor(key)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple text",
			plaintext: "hello world",
		},
		{
			name:      "access token",
			plaintext: "ghp_1234567890abcdefghijklmnopqrstuvwxyz",
		},
		{
			name:      "long text",
			plaintext: "This is a longer text that contains multiple words and special characters: !@#$%^&*()_+-=[]{}|;:',.<>?/`~",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "unicode text",
			plaintext: "日本語のテキスト 🔐",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciphertext, err := encryptor.Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}

			if tt.plaintext == "" {
				if ciphertext != "" {
					t.Errorf("expected empty ciphertext for empty plaintext, got %q", ciphertext)
				}
				return
			}

			if ciphertext == tt.plaintext {
				t.Error("ciphertext should not equal plaintext")
			}

			decrypted, err := encryptor.Decrypt(ciphertext)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}

			if decrypted != tt.plaintext {
				t.Errorf("Decrypt() = %q, want %q", decrypted, tt.plaintext)
			}
		})
	}
}

func TestAESEncryptor_EncryptProducesDifferentCiphertexts(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("failed to generate random key: %v", err)
	}

	encryptor, err := NewAESEncryptor(key)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	plaintext := "test data"

	ciphertext1, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	ciphertext2, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	if ciphertext1 == ciphertext2 {
		t.Error("encrypting the same plaintext twice should produce different ciphertexts")
	}

	decrypted1, err := encryptor.Decrypt(ciphertext1)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if decrypted1 != plaintext {
		t.Errorf("decrypted1 = %q, want %q", decrypted1, plaintext)
	}

	decrypted2, err := encryptor.Decrypt(ciphertext2)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if decrypted2 != plaintext {
		t.Errorf("decrypted2 = %q, want %q", decrypted2, plaintext)
	}
}

func TestAESEncryptor_DecryptInvalidData(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("failed to generate random key: %v", err)
	}

	encryptor, err := NewAESEncryptor(key)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	tests := []struct {
		name       string
		ciphertext string
		wantErr    bool
	}{
		{
			name:       "invalid base64",
			ciphertext: "not-valid-base64!@#$",
			wantErr:    true,
		},
		{
			name:       "too short data",
			ciphertext: "YWJj", // "abc" in base64, too short for nonce
			wantErr:    true,
		},
		{
			name:       "corrupted data",
			ciphertext: "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo=",
			wantErr:    true,
		},
		{
			name:       "empty string",
			ciphertext: "",
			wantErr:    false, // Empty ciphertext should return empty plaintext
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := encryptor.Decrypt(tt.ciphertext)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAESEncryptor_DecryptWithWrongKey(t *testing.T) {
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	if _, err := rand.Read(key1); err != nil {
		t.Fatalf("failed to generate key1: %v", err)
	}
	if _, err := rand.Read(key2); err != nil {
		t.Fatalf("failed to generate key2: %v", err)
	}

	encryptor1, err := NewAESEncryptor(key1)
	if err != nil {
		t.Fatalf("failed to create encryptor1: %v", err)
	}

	encryptor2, err := NewAESEncryptor(key2)
	if err != nil {
		t.Fatalf("failed to create encryptor2: %v", err)
	}

	plaintext := "secret data"

	ciphertext, err := encryptor1.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	_, err = encryptor2.Decrypt(ciphertext)
	if err == nil {
		t.Error("Decrypt() with wrong key should return an error")
	}
}
