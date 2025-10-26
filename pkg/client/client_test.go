package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	config := Config{
		URL:        "https://portainer.example.com",
		Token:      "test-token",
		EndpointID: 1,
		AuthType:   "jwt",
	}

	client := NewClient(config)

	if client.BaseURL != config.URL {
		t.Errorf("NewClient() BaseURL = %q, want %q", client.BaseURL, config.URL)
	}
	if client.Token != config.Token {
		t.Errorf("NewClient() Token = %q, want %q", client.Token, config.Token)
	}
	if client.EndpointID != config.EndpointID {
		t.Errorf("NewClient() EndpointID = %d, want %d", client.EndpointID, config.EndpointID)
	}
	if client.AuthType != config.AuthType {
		t.Errorf("NewClient() AuthType = %q, want %q", client.AuthType, config.AuthType)
	}
	if client.HTTPClient == nil {
		t.Errorf("NewClient() HTTPClient should not be nil")
	}
}

func TestNewClientDefaultAuthType(t *testing.T) {
	config := Config{
		URL:        "https://portainer.example.com",
		Token:      "test-token",
		EndpointID: 1,
		// AuthType not set
	}

	client := NewClient(config)

	if client.AuthType != "api-key" {
		t.Errorf("NewClient() AuthType = %q, want %q", client.AuthType, "api-key")
	}
}

func TestNewAuthClient(t *testing.T) {
	baseURL := "https://portainer.example.com"
	client := NewAuthClient(baseURL)

	if client.BaseURL != baseURL {
		t.Errorf("NewAuthClient() BaseURL = %q, want %q", client.BaseURL, baseURL)
	}
	if client.AuthType != "none" {
		t.Errorf("NewAuthClient() AuthType = %q, want %q", client.AuthType, "none")
	}
	if client.HTTPClient == nil {
		t.Errorf("NewAuthClient() HTTPClient should not be nil")
	}
}

func TestMakeRequestJWT(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify JWT auth header
		auth := r.Header.Get("Authorization")
		expected := "Bearer test-jwt-token"
		if auth != expected {
			t.Errorf("Authorization header = %q, want %q", auth, expected)
		}

		// Verify content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Content-Type header = %q, want %q", contentType, "application/json")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := &Client{
		BaseURL:  server.URL,
		Token:    "test-jwt-token",
		AuthType: "jwt",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	resp, err := client.makeRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("makeRequest() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("makeRequest() status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestMakeRequestAPIKey(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify API key header
		apiKey := r.Header.Get("X-API-Key")
		expected := "test-api-key"
		if apiKey != expected {
			t.Errorf("X-API-Key header = %q, want %q", apiKey, expected)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := &Client{
		BaseURL:  server.URL,
		Token:    "test-api-key",
		AuthType: "api-key",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	resp, err := client.makeRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("makeRequest() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("makeRequest() status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestMakeRequestNoAuth(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify no auth headers
		if auth := r.Header.Get("Authorization"); auth != "" {
			t.Errorf("Should not have Authorization header, got %q", auth)
		}
		if apiKey := r.Header.Get("X-API-Key"); apiKey != "" {
			t.Errorf("Should not have X-API-Key header, got %q", apiKey)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := &Client{
		BaseURL:  server.URL,
		AuthType: "none",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	resp, err := client.makeRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("makeRequest() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("makeRequest() status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestMakeRequestWithBody(t *testing.T) {
	requestBody := map[string]string{"key": "value"}

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request body
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if body["key"] != "value" {
			t.Errorf("Request body = %v, want %v", body, requestBody)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := &Client{
		BaseURL:  server.URL,
		AuthType: "none",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	resp, err := client.makeRequest("POST", "/test", requestBody)
	if err != nil {
		t.Fatalf("makeRequest() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("makeRequest() status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestMakeRequestError(t *testing.T) {
	// Mock server returning error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	client := &Client{
		BaseURL:  server.URL,
		AuthType: "none",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	_, err := client.makeRequest("GET", "/test", nil)
	if err == nil {
		t.Errorf("makeRequest() should return error for 400 status")
	}

	expectedError := "API request failed with status 400: Bad Request"
	if err.Error() != expectedError {
		t.Errorf("makeRequest() error = %q, want %q", err.Error(), expectedError)
	}
}

func TestSetJWTToken(t *testing.T) {
	client := &Client{}
	token := "test-jwt-token"

	client.SetJWTToken(token)

	if client.Token != token {
		t.Errorf("SetJWTToken() token = %q, want %q", client.Token, token)
	}
	if client.AuthType != "jwt" {
		t.Errorf("SetJWTToken() authType = %q, want %q", client.AuthType, "jwt")
	}
}

func TestSetAPIKey(t *testing.T) {
	client := &Client{}
	token := "test-api-key"

	client.SetAPIKey(token)

	if client.Token != token {
		t.Errorf("SetAPIKey() token = %q, want %q", client.Token, token)
	}
	if client.AuthType != "api-key" {
		t.Errorf("SetAPIKey() authType = %q, want %q", client.AuthType, "api-key")
	}
}

func TestIsAuthPath(t *testing.T) {
	client := &Client{}

	tests := []struct {
		path     string
		expected bool
	}{
		{"/api/auth", true},
		{"/api/auth/login", true},
		{"/api/auth/oauth", true},
		{"/api/auth/oauth/login", true},
		{"/api/stacks", false},
		{"/api/users", false},
		{"/health", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := client.IsAuthPath(tt.path)
			if result != tt.expected {
				t.Errorf("IsAuthPath(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}
