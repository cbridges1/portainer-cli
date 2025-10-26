package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jalenbridges/portainer-cli/pkg/storage"
)

func NewAuthenticatedClient() (*Client, error) {
	config, err := storage.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if config.URL == "" {
		return nil, fmt.Errorf("Portainer URL not configured. Use 'portainer-cli config set url <url>' first")
	}

	// Check if we have a valid token
	if config.IsTokenValid() {
		jwt, _, err := config.GetToken()
		if err != nil {
			return nil, fmt.Errorf("failed to get stored token: %w", err)
		}

		// Create client with JWT
		client := &Client{
			BaseURL:    config.URL,
			Token:      jwt,
			AuthType:   "jwt",
			EndpointID: getEndpointFromConfig(config),
			HTTPClient: defaultHTTPClient(),
		}

		// Validate token by making a test request
		if _, err := client.ValidateToken(); err == nil {
			return client, nil
		}
		// If validation fails, continue to refresh logic below
	}

	// Token is invalid or missing, try to refresh using stored credentials
	if config.Credentials == nil {
		return nil, fmt.Errorf("no valid token or stored credentials found. Please login using 'portainer-cli login'")
	}

	username, password, credType, err := config.GetCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to get stored credentials: %w", err)
	}

	// Create auth client for re-authentication
	authClient := NewAuthClient(config.URL)

	var authResp *AuthResponse
	switch credType {
	case "password":
		authResp, err = authClient.AuthenticateUser(username, password)
	case "code":
		authResp, err = authClient.AuthenticateWithCode(password) // password field stores the code
	default:
		return nil, fmt.Errorf("unknown credential type: %s", credType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to refresh authentication: %w", err)
	}

	// Parse and store new token
	expiresAt, err := ParseJWTExpiry(authResp.JWT)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT expiry: %w", err)
	}

	if err := config.SetToken(authResp.JWT, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to store refreshed token: %w", err)
	}

	// Create authenticated client
	client := &Client{
		BaseURL:    config.URL,
		Token:      authResp.JWT,
		AuthType:   "jwt",
		EndpointID: getEndpointFromConfig(config),
		HTTPClient: defaultHTTPClient(),
	}

	return client, nil
}

func getEndpointFromConfig(config *storage.Config) int {
	if config.EndpointID != 0 {
		return config.EndpointID
	}
	// Default to 1 if not configured
	return 1
}

func defaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 300 * time.Second,
	}
}
