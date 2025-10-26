package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthenticateUser(t *testing.T) {
	expectedJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.signature"

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/auth" {
			t.Errorf("Expected path /api/auth, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Verify request body
		var authReq AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if authReq.Username != "testuser" {
			t.Errorf("Expected username 'testuser', got %s", authReq.Username)
		}
		if authReq.Password != "testpass" {
			t.Errorf("Expected password 'testpass', got %s", authReq.Password)
		}

		// Return mock JWT
		resp := AuthResponse{JWT: expectedJWT}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewAuthClient(server.URL)

	authResp, err := client.AuthenticateUser("testuser", "testpass")
	if err != nil {
		t.Fatalf("AuthenticateUser() error = %v", err)
	}

	if authResp.JWT != expectedJWT {
		t.Errorf("AuthenticateUser() JWT = %q, want %q", authResp.JWT, expectedJWT)
	}
}

func TestAuthenticateUserError(t *testing.T) {
	// Mock server returning error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid credentials"))
	}))
	defer server.Close()

	client := NewAuthClient(server.URL)

	_, err := client.AuthenticateUser("wronguser", "wrongpass")
	if err == nil {
		t.Errorf("AuthenticateUser() should return error for invalid credentials")
	}
}

func TestAuthenticateWithCode(t *testing.T) {
	expectedJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.signature"

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/auth" {
			t.Errorf("Expected path /api/auth, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Verify request body
		var authReq AccessCodeAuthRequest
		if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if authReq.AccessCode != "ABCD1234" {
			t.Errorf("Expected access code 'ABCD1234', got %s", authReq.AccessCode)
		}

		// Return mock JWT
		resp := AuthResponse{JWT: expectedJWT}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewAuthClient(server.URL)

	authResp, err := client.AuthenticateWithCode("ABCD1234")
	if err != nil {
		t.Fatalf("AuthenticateWithCode() error = %v", err)
	}

	if authResp.JWT != expectedJWT {
		t.Errorf("AuthenticateWithCode() JWT = %q, want %q", authResp.JWT, expectedJWT)
	}
}

func TestAuthenticateWithCodeError(t *testing.T) {
	// Mock server returning error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid access code"))
	}))
	defer server.Close()

	client := NewAuthClient(server.URL)

	_, err := client.AuthenticateWithCode("INVALID")
	if err == nil {
		t.Errorf("AuthenticateWithCode() should return error for invalid code")
	}
}

func TestValidateToken(t *testing.T) {
	expectedUser := User{
		ID:       1,
		Username: "testuser",
		Role:     1,
	}

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/users/me" {
			t.Errorf("Expected path /api/users/me, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Verify Authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-jwt" {
			t.Errorf("Expected Authorization 'Bearer test-jwt', got %s", auth)
		}

		// Return mock user
		resp := AuthStatus{User: expectedUser}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:  server.URL,
		Token:    "test-jwt",
		AuthType: "jwt",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	status, err := client.ValidateToken()
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if status.User.ID != expectedUser.ID {
		t.Errorf("ValidateToken() User.ID = %d, want %d", status.User.ID, expectedUser.ID)
	}
	if status.User.Username != expectedUser.Username {
		t.Errorf("ValidateToken() User.Username = %q, want %q", status.User.Username, expectedUser.Username)
	}
}

func TestValidateTokenError(t *testing.T) {
	// Mock server returning error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token"))
	}))
	defer server.Close()

	client := &Client{
		BaseURL:  server.URL,
		Token:    "invalid-jwt",
		AuthType: "jwt",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	_, err := client.ValidateToken()
	if err == nil {
		t.Errorf("ValidateToken() should return error for invalid token")
	}
}

func TestParseJWTExpiry(t *testing.T) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.signature"

	expiry, err := ParseJWTExpiry(jwt)
	if err != nil {
		t.Fatalf("ParseJWTExpiry() error = %v", err)
	}

	// Should return a time in the future (current implementation adds 8 hours)
	if expiry.Before(time.Now()) {
		t.Errorf("ParseJWTExpiry() should return future time")
	}

	// Should be approximately 8 hours from now (within 1 minute tolerance)
	expectedExpiry := time.Now().Add(8 * time.Hour)
	diff := expiry.Sub(expectedExpiry)
	if diff < -time.Minute || diff > time.Minute {
		t.Errorf("ParseJWTExpiry() expiry time difference too large: %v", diff)
	}
}

func TestAuthRequestStructure(t *testing.T) {
	req := AuthRequest{
		Username: "testuser",
		Password: "testpass",
	}

	// Test JSON marshaling
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal AuthRequest: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled AuthRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal AuthRequest: %v", err)
	}

	if unmarshaled.Username != req.Username {
		t.Errorf("Unmarshaled Username = %q, want %q", unmarshaled.Username, req.Username)
	}
	if unmarshaled.Password != req.Password {
		t.Errorf("Unmarshaled Password = %q, want %q", unmarshaled.Password, req.Password)
	}
}

func TestAccessCodeAuthRequestStructure(t *testing.T) {
	req := AccessCodeAuthRequest{
		AccessCode: "ABCD1234",
	}

	// Test JSON marshaling
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal AccessCodeAuthRequest: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled AccessCodeAuthRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal AccessCodeAuthRequest: %v", err)
	}

	if unmarshaled.AccessCode != req.AccessCode {
		t.Errorf("Unmarshaled AccessCode = %q, want %q", unmarshaled.AccessCode, req.AccessCode)
	}
}
