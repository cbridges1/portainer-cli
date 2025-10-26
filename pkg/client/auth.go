package client

import (
	"encoding/json"
	"fmt"
	"time"
)

type AuthRequest struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type AuthResponse struct {
	JWT string `json:"jwt"`
}

type AccessCodeAuthRequest struct {
	AccessCode string `json:"AccessCode"`
}

type User struct {
	ID       int    `json:"Id"`
	Username string `json:"Username"`
	Role     int    `json:"Role"`
}

type AuthStatus struct {
	User User `json:"User"`
}

func (c *Client) AuthenticateUser(username, password string) (*AuthResponse, error) {
	authReq := AuthRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.makeRequest("POST", "/api/auth", authReq)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, fmt.Errorf("failed to decode auth response: %w", err)
	}

	return &authResp, nil
}

func (c *Client) AuthenticateWithCode(accessCode string) (*AuthResponse, error) {
	authReq := AccessCodeAuthRequest{
		AccessCode: accessCode,
	}

	resp, err := c.makeRequest("POST", "/api/auth", authReq)
	if err != nil {
		return nil, fmt.Errorf("authentication with code failed: %w", err)
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, fmt.Errorf("failed to decode auth response: %w", err)
	}

	return &authResp, nil
}

func (c *Client) ValidateToken() (*AuthStatus, error) {
	resp, err := c.makeRequest("GET", "/api/users/me", nil)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}
	defer resp.Body.Close()

	var status AuthStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode auth status: %w", err)
	}

	return &status, nil
}

func ParseJWTExpiry(jwt string) (time.Time, error) {
	// For production use, you'd want to properly parse the JWT
	// For now, we'll estimate 8 hours from creation
	return time.Now().Add(8 * time.Hour), nil
}
