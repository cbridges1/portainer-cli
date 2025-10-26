package storage

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple string",
			plaintext: "hello world",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "complex password",
			plaintext: "myP@ssw0rd!123",
		},
		{
			name:      "json data",
			plaintext: `{"username":"test","password":"secret"}`,
		},
		{
			name:      "unicode characters",
			plaintext: "héllo wørld 🌍",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt the plaintext
			encrypted, err := Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}

			// Verify encrypted text is different from plaintext
			if encrypted == tt.plaintext {
				t.Errorf("Encrypted text should be different from plaintext")
			}

			// Decrypt the encrypted text
			decrypted, err := Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}

			// Verify decrypted text matches original plaintext
			if decrypted != tt.plaintext {
				t.Errorf("Decrypt() = %q, want %q", decrypted, tt.plaintext)
			}
		})
	}
}

func TestEncryptConsistency(t *testing.T) {
	plaintext := "test string"

	// Encrypt the same string multiple times
	encrypted1, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("First Encrypt() error = %v", err)
	}

	encrypted2, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Second Encrypt() error = %v", err)
	}

	// Encrypted values should be different (due to random nonce)
	if encrypted1 == encrypted2 {
		t.Errorf("Multiple encryptions of same text should produce different results")
	}

	// But both should decrypt to the same plaintext
	decrypted1, err := Decrypt(encrypted1)
	if err != nil {
		t.Fatalf("First Decrypt() error = %v", err)
	}

	decrypted2, err := Decrypt(encrypted2)
	if err != nil {
		t.Fatalf("Second Decrypt() error = %v", err)
	}

	if decrypted1 != plaintext || decrypted2 != plaintext {
		t.Errorf("Both decryptions should equal original plaintext")
	}
}

func TestDecryptInvalidData(t *testing.T) {
	tests := []struct {
		name        string
		ciphertext  string
		expectError bool
	}{
		{
			name:        "invalid base64",
			ciphertext:  "not-valid-base64!@#",
			expectError: true,
		},
		{
			name:        "empty string",
			ciphertext:  "",
			expectError: true,
		},
		{
			name:        "too short data",
			ciphertext:  "YWJj", // "abc" in base64
			expectError: true,
		},
		{
			name:        "random data",
			ciphertext:  "cmFuZG9tZGF0YXRoYXRpc250dmFsaWQ=", // "randomdatathatisntvalid" in base64
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decrypt(tt.ciphertext)
			if (err != nil) != tt.expectError {
				t.Errorf("Decrypt() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected string
	}{
		{
			name:     "simple password",
			password: "password",
			expected: "XohImNooBHFR0OVvjcYpJ3NgPQ1qq73WKhHvch0VQtg=", // SHA256 of "password"
		},
		{
			name:     "empty password",
			password: "",
			expected: "47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=", // SHA256 of empty string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HashPassword(tt.password)
			if result != tt.expected {
				t.Errorf("HashPassword() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestHashPasswordConsistency(t *testing.T) {
	password := "test-password"

	hash1 := HashPassword(password)
	hash2 := HashPassword(password)

	if hash1 != hash2 {
		t.Errorf("HashPassword should be deterministic, got different hashes: %q vs %q", hash1, hash2)
	}
}

func TestGetMachineID(t *testing.T) {
	id, err := getMachineID()
	if err != nil {
		t.Fatalf("getMachineID() error = %v", err)
	}

	if id == "" {
		t.Errorf("getMachineID() should return non-empty string")
	}

	// Test consistency
	id2, err := getMachineID()
	if err != nil {
		t.Fatalf("getMachineID() second call error = %v", err)
	}

	if id != id2 {
		t.Errorf("getMachineID() should be consistent across calls")
	}
}

func TestDeriveKey(t *testing.T) {
	key, err := deriveKey()
	if err != nil {
		t.Fatalf("deriveKey() error = %v", err)
	}

	if len(key) != 32 {
		t.Errorf("deriveKey() should return 32-byte key, got %d bytes", len(key))
	}

	// Test consistency
	key2, err := deriveKey()
	if err != nil {
		t.Fatalf("deriveKey() second call error = %v", err)
	}

	if string(key) != string(key2) {
		t.Errorf("deriveKey() should be consistent across calls")
	}
}
