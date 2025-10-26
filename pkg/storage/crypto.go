package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"runtime"

	"golang.org/x/crypto/scrypt"
)

func getMachineID() (string, error) {
	var machineID string

	switch runtime.GOOS {
	case "darwin":
		// Use hardware UUID on macOS
		if data, err := os.ReadFile("/proc/sys/kernel/random/uuid"); err == nil {
			machineID = string(data)
		} else {
			// Fallback to hostname
			if hostname, err := os.Hostname(); err == nil {
				machineID = hostname
			} else {
				return "", fmt.Errorf("failed to get machine identifier")
			}
		}
	case "linux":
		// Use machine-id on Linux
		if data, err := os.ReadFile("/etc/machine-id"); err == nil {
			machineID = string(data)
		} else if data, err := os.ReadFile("/var/lib/dbus/machine-id"); err == nil {
			machineID = string(data)
		} else {
			// Fallback to hostname
			if hostname, err := os.Hostname(); err == nil {
				machineID = hostname
			} else {
				return "", fmt.Errorf("failed to get machine identifier")
			}
		}
	case "windows":
		// Use hostname on Windows (could be enhanced with registry)
		if hostname, err := os.Hostname(); err == nil {
			machineID = hostname
		} else {
			return "", fmt.Errorf("failed to get machine identifier")
		}
	default:
		// Fallback for other systems
		if hostname, err := os.Hostname(); err == nil {
			machineID = hostname
		} else {
			return "", fmt.Errorf("failed to get machine identifier")
		}
	}

	return machineID, nil
}

func deriveKey() ([]byte, error) {
	machineID, err := getMachineID()
	if err != nil {
		return nil, err
	}

	// Add a constant salt to make the key derivation more secure
	salt := []byte("portainer-cli-salt-2024")

	// Derive a 32-byte key using scrypt
	key, err := scrypt.Key([]byte(machineID), salt, 32768, 8, 1, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to derive encryption key: %w", err)
	}

	return key, nil
}

func Encrypt(plaintext string) (string, error) {
	key, err := deriveKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertext string) (string, error) {
	key, err := deriveKey()
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(hash[:])
}
