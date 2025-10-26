package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConfigSetGetURL(t *testing.T) {
	// Create temporary config for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	config := &Config{}

	testURL := "https://portainer.example.com"
	err := config.SetURL(testURL)
	if err != nil {
		t.Fatalf("SetURL() error = %v", err)
	}

	if config.URL != testURL {
		t.Errorf("SetURL() did not set URL correctly, got %q, want %q", config.URL, testURL)
	}

	// Test loading from file
	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if loadedConfig.URL != testURL {
		t.Errorf("LoadConfig() URL = %q, want %q", loadedConfig.URL, testURL)
	}
}

func TestConfigSetGetEndpoint(t *testing.T) {
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	config := &Config{}

	testEndpoint := 5
	err := config.SetEndpoint(testEndpoint)
	if err != nil {
		t.Fatalf("SetEndpoint() error = %v", err)
	}

	if config.EndpointID != testEndpoint {
		t.Errorf("SetEndpoint() did not set endpoint correctly, got %d, want %d", config.EndpointID, testEndpoint)
	}

	// Test loading from file
	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if loadedConfig.EndpointID != testEndpoint {
		t.Errorf("LoadConfig() EndpointID = %d, want %d", loadedConfig.EndpointID, testEndpoint)
	}
}

func TestConfigCredentials(t *testing.T) {
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	config := &Config{}

	testUsername := "testuser"
	testPassword := "testpassword"
	testType := "password"

	err := config.SetCredentials(testUsername, testPassword, testType)
	if err != nil {
		t.Fatalf("SetCredentials() error = %v", err)
	}

	// Test getting credentials
	username, password, credType, err := config.GetCredentials()
	if err != nil {
		t.Fatalf("GetCredentials() error = %v", err)
	}

	if username != testUsername {
		t.Errorf("GetCredentials() username = %q, want %q", username, testUsername)
	}
	if password != testPassword {
		t.Errorf("GetCredentials() password = %q, want %q", password, testPassword)
	}
	if credType != testType {
		t.Errorf("GetCredentials() type = %q, want %q", credType, testType)
	}

	// Test loading from file
	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	username2, password2, credType2, err := loadedConfig.GetCredentials()
	if err != nil {
		t.Fatalf("GetCredentials() from loaded config error = %v", err)
	}

	if username2 != testUsername || password2 != testPassword || credType2 != testType {
		t.Errorf("Loaded credentials don't match: got (%q, %q, %q), want (%q, %q, %q)",
			username2, password2, credType2, testUsername, testPassword, testType)
	}
}

func TestConfigToken(t *testing.T) {
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	config := &Config{}

	testJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.token"
	testExpiry := time.Now().Add(1 * time.Hour)

	err := config.SetToken(testJWT, testExpiry)
	if err != nil {
		t.Fatalf("SetToken() error = %v", err)
	}

	// Test getting token
	jwt, expiry, err := config.GetToken()
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}

	if jwt != testJWT {
		t.Errorf("GetToken() jwt = %q, want %q", jwt, testJWT)
	}
	if !expiry.Equal(testExpiry) {
		t.Errorf("GetToken() expiry = %v, want %v", expiry, testExpiry)
	}

	// Test token validity
	if !config.IsTokenValid() {
		t.Errorf("IsTokenValid() should return true for future expiry")
	}

	// Test with expired token
	pastExpiry := time.Now().Add(-1 * time.Hour)
	err = config.SetToken(testJWT, pastExpiry)
	if err != nil {
		t.Fatalf("SetToken() with past expiry error = %v", err)
	}

	if config.IsTokenValid() {
		t.Errorf("IsTokenValid() should return false for past expiry")
	}
}

func TestConfigClearCredentials(t *testing.T) {
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	config := &Config{}

	// Set some credentials and token
	err := config.SetCredentials("user", "pass", "password")
	if err != nil {
		t.Fatalf("SetCredentials() error = %v", err)
	}

	err = config.SetToken("token", time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("SetToken() error = %v", err)
	}

	// Clear credentials
	err = config.ClearCredentials()
	if err != nil {
		t.Fatalf("ClearCredentials() error = %v", err)
	}

	// Verify credentials are cleared
	if config.Credentials != nil {
		t.Errorf("ClearCredentials() should set Credentials to nil")
	}
	if config.Token != nil {
		t.Errorf("ClearCredentials() should set Token to nil")
	}

	// Test getting credentials after clearing
	_, _, _, err = config.GetCredentials()
	if err == nil {
		t.Errorf("GetCredentials() should return error after clearing")
	}
}

func TestLoadConfigNewFile(t *testing.T) {
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Load config when file doesn't exist
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Should return empty config
	if config.URL != "" || config.EndpointID != 0 || config.Credentials != nil || config.Token != nil {
		t.Errorf("LoadConfig() should return empty config when file doesn't exist")
	}
}

func TestGetConfigDir(t *testing.T) {
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir() error = %v", err)
	}

	expectedDir := filepath.Join(tempDir, ".config", "portainer-cli")
	if configDir != expectedDir {
		t.Errorf("getConfigDir() = %q, want %q", configDir, expectedDir)
	}

	// Verify directory was created
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("getConfigDir() should create the directory")
	}
}

func TestConfigFilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	config := &Config{}
	err := config.SetURL("https://test.com")
	if err != nil {
		t.Fatalf("SetURL() error = %v", err)
	}

	configPath, err := getConfigPath()
	if err != nil {
		t.Fatalf("getConfigPath() error = %v", err)
	}

	// Check file permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}

	perm := info.Mode().Perm()
	expected := os.FileMode(0600)
	if perm != expected {
		t.Errorf("Config file permissions = %v, want %v", perm, expected)
	}
}
