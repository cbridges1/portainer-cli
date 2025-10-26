package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

func (c *Client) ListStacks() ([]Stack, error) {
	resp, err := c.makeRequest("GET", "/api/stacks", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var stacks []Stack
	if err := json.NewDecoder(resp.Body).Decode(&stacks); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return stacks, nil
}

func (c *Client) FindStackByName(name string) (*Stack, error) {
	resp, err := c.makeRequest("GET", "/api/stacks", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var stacks []Stack
	if err := json.NewDecoder(resp.Body).Decode(&stacks); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	for _, stack := range stacks {
		if stack.Name == name {
			return &stack, nil
		}
	}

	return nil, nil
}

func (c *Client) GetStack(stackID int) (*Stack, error) {
	path := fmt.Sprintf("/api/stacks/%d", stackID)
	resp, err := c.makeRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var stack Stack
	if err := json.NewDecoder(resp.Body).Decode(&stack); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &stack, nil
}

func (c *Client) CreateStackFromRepository(req CreateStackRequest) (*Stack, error) {
	params := url.Values{}
	params.Add("method", "repository")
	params.Add("type", "1") // Docker Compose
	params.Add("endpointId", strconv.Itoa(c.EndpointID))

	path := "/api/stacks?" + params.Encode()
	resp, err := c.makeRequest("POST", path, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var stack Stack
	if err := json.NewDecoder(resp.Body).Decode(&stack); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &stack, nil
}

func (c *Client) CreateStackFromString(req CreateStackRequest) (*Stack, error) {
	params := url.Values{}
	params.Add("method", "string")
	params.Add("type", "1") // Docker Compose
	params.Add("endpointId", strconv.Itoa(c.EndpointID))

	path := "/api/stacks/create/standalone/string?" + params.Encode()

	resp, err := c.makeRequest("POST", path, req)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	var stack Stack
	if err := json.NewDecoder(resp.Body).Decode(&stack); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	println(resp)

	return &stack, nil
}

func (c *Client) UpdateStack(stackID int, req UpdateStackRequest) (*Stack, error) {
	params := url.Values{}
	params.Add("endpointId", strconv.Itoa(c.EndpointID))

	path := fmt.Sprintf("/api/stacks/%d?%s", stackID, params.Encode())
	resp, err := c.makeRequest("PUT", path, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var stack Stack
	if err := json.NewDecoder(resp.Body).Decode(&stack); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &stack, nil
}

func (c *Client) DeleteStack(stackID int) error {
	params := url.Values{}
	params.Add("endpointId", strconv.Itoa(c.EndpointID))

	path := fmt.Sprintf("/api/stacks/%d?%s", stackID, params.Encode())
	resp, err := c.makeRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read and discard the response body
	_, err = io.ReadAll(resp.Body)
	return err
}

func (c *Client) GetStackFile(stackID int) (string, error) {
	path := fmt.Sprintf("/api/stacks/%d/file", stackID)
	resp, err := c.makeRequest("GET", path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		StackFileContent string `json:"StackFileContent"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.StackFileContent, nil
}
