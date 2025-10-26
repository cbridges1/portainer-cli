package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	Token      string
	EndpointID int
	HTTPClient *http.Client
	AuthType   string // "jwt" or "api-key"
}

type Config struct {
	URL        string
	Token      string
	EndpointID int
	AuthType   string
}

func NewClient(config Config) *Client {
	authType := config.AuthType
	if authType == "" {
		authType = "api-key" // Default to API key for backward compatibility
	}

	return &Client{
		BaseURL:    config.URL,
		Token:      config.Token,
		EndpointID: config.EndpointID,
		AuthType:   authType,
		HTTPClient: &http.Client{
			Timeout: 300 * time.Second,
		},
	}
}

func NewAuthClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		AuthType: "none", // No auth for authentication endpoints
	}
}

func (c *Client) makeRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authentication header based on auth type
	switch c.AuthType {
	case "jwt":
		if c.Token != "" {
			req.Header.Set("Authorization", "Bearer "+c.Token)
		}
	case "api-key":
		if c.Token != "" {
			req.Header.Set("X-API-Key", c.Token)
		}
	case "none":
		// No authentication header
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

func (c *Client) SetJWTToken(token string) {
	c.Token = token
	c.AuthType = "jwt"
}

func (c *Client) SetAPIKey(token string) {
	c.Token = token
	c.AuthType = "api-key"
}

func (c *Client) IsAuthPath(path string) bool {
	authPaths := []string{"/api/auth", "/api/auth/oauth"}
	for _, authPath := range authPaths {
		if strings.HasPrefix(path, authPath) {
			return true
		}
	}
	return false
}
