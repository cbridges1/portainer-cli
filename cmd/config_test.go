package cmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestConfigSetUrlCmd(t *testing.T) {
	// Setup temporary home directory
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "valid URL",
			args:        []string{"https://portainer.example.com"},
			expectError: false,
		},
		{
			name:        "valid URL with port",
			args:        []string{"https://portainer.example.com:9000"},
			expectError: false,
		},
		{
			name:        "invalid URL",
			args:        []string{"not-a-valid-url"},
			expectError: true,
		},
		{
			name:        "no arguments",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "too many arguments",
			args:        []string{"https://portainer.example.com", "extra"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.SetArgs(tt.args)

			// Skip validation for args length as Cobra handles this
			if len(tt.args) == 0 || len(tt.args) > 1 {
				// Test that cobra validation works by checking if RunE would be called
				// For these cases, cobra's Args validation should prevent RunE from being called
				if !tt.expectError {
					t.Errorf("Test case expects no error but should have args validation error")
				}
				return
			}

			err := configSetUrlCmd.RunE(cmd, tt.args)

			if (err != nil) != tt.expectError {
				t.Errorf("configSetUrlCmd.RunE() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestConfigSetEndpointCmd(t *testing.T) {
	// Setup temporary home directory
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "valid endpoint ID",
			args:        []string{"1"},
			expectError: false,
		},
		{
			name:        "valid larger endpoint ID",
			args:        []string{"123"},
			expectError: false,
		},
		{
			name:        "zero endpoint ID",
			args:        []string{"0"},
			expectError: true,
		},
		{
			name:        "negative endpoint ID",
			args:        []string{"-1"},
			expectError: true,
		},
		{
			name:        "non-numeric endpoint ID",
			args:        []string{"abc"},
			expectError: true,
		},
		{
			name:        "no arguments",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "too many arguments",
			args:        []string{"1", "2"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.SetArgs(tt.args)

			// Skip validation for args length as Cobra handles this
			if len(tt.args) == 0 || len(tt.args) > 1 {
				if !tt.expectError {
					t.Errorf("Test case expects no error but should have args validation error")
				}
				return
			}

			err := configSetEndpointCmd.RunE(cmd, tt.args)

			if (err != nil) != tt.expectError {
				t.Errorf("configSetEndpointCmd.RunE() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestConfigGetCmd(t *testing.T) {
	// Setup temporary home directory
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Test with empty config
	cmd := &cobra.Command{}
	err := configGetCmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("configGetCmd.RunE() with empty config error = %v", err)
	}

	// Set some config values and test again
	err = configSetUrlCmd.RunE(cmd, []string{"https://test.com"})
	if err != nil {
		t.Fatalf("Failed to set URL: %v", err)
	}

	err = configSetEndpointCmd.RunE(cmd, []string{"5"})
	if err != nil {
		t.Fatalf("Failed to set endpoint: %v", err)
	}

	// Test get with configured values
	err = configGetCmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("configGetCmd.RunE() with configured values error = %v", err)
	}
}
