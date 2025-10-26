package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	URL         string                `json:"url"`
	EndpointID  int                   `json:"endpoint_id,omitempty"`
	Credentials *EncryptedCredentials `json:"credentials,omitempty"`
	Token       *EncryptedToken       `json:"token,omitempty"`
}

type EncryptedCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"` // Encrypted
	Type     string `json:"type"`     // "password" or "code"
}

type EncryptedToken struct {
	JWT       string    `json:"jwt"` // Encrypted
	ExpiresAt time.Time `json:"expires_at"`
	Refresh   string    `json:"refresh,omitempty"` // Encrypted refresh token if available
}

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "portainer-cli")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}

func getConfigPath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (c *Config) SetURL(url string) error {
	c.URL = url
	return c.Save()
}

func (c *Config) SetEndpoint(endpointID int) error {
	c.EndpointID = endpointID
	return c.Save()
}

func (c *Config) SetCredentials(username, password, credType string) error {
	encryptedPassword, err := Encrypt(password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	c.Credentials = &EncryptedCredentials{
		Username: username,
		Password: encryptedPassword,
		Type:     credType,
	}

	return c.Save()
}

func (c *Config) GetCredentials() (string, string, string, error) {
	if c.Credentials == nil {
		return "", "", "", fmt.Errorf("no credentials stored")
	}

	password, err := Decrypt(c.Credentials.Password)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to decrypt password: %w", err)
	}

	return c.Credentials.Username, password, c.Credentials.Type, nil
}

func (c *Config) SetToken(jwt string, expiresAt time.Time) error {
	encryptedJWT, err := Encrypt(jwt)
	if err != nil {
		return fmt.Errorf("failed to encrypt JWT: %w", err)
	}

	c.Token = &EncryptedToken{
		JWT:       encryptedJWT,
		ExpiresAt: expiresAt,
	}

	return c.Save()
}

func (c *Config) GetToken() (string, time.Time, error) {
	if c.Token == nil {
		return "", time.Time{}, fmt.Errorf("no token stored")
	}

	jwt, err := Decrypt(c.Token.JWT)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to decrypt JWT: %w", err)
	}

	return jwt, c.Token.ExpiresAt, nil
}

func (c *Config) IsTokenValid() bool {
	if c.Token == nil {
		return false
	}
	return time.Now().Before(c.Token.ExpiresAt.Add(-5 * time.Minute)) // 5 minute buffer
}

func (c *Config) ClearCredentials() error {
	c.Credentials = nil
	c.Token = nil
	return c.Save()
}
