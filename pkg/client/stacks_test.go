package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListStacks(t *testing.T) {
	mockStacks := []Stack{
		{
			ID:           1,
			Name:         "test-stack-1",
			Type:         1,
			EndpointID:   1,
			Status:       1,
			CreationDate: time.Now().Unix(),
			CreatedBy:    "admin",
		},
		{
			ID:           2,
			Name:         "test-stack-2",
			Type:         1,
			EndpointID:   1,
			Status:       1,
			CreationDate: time.Now().Unix(),
			CreatedBy:    "user",
		},
	}

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/stacks" {
			t.Errorf("Expected path /api/stacks, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockStacks)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:  server.URL,
		AuthType: "none",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	stacks, err := client.ListStacks()
	if err != nil {
		t.Fatalf("ListStacks() error = %v", err)
	}

	if len(stacks) != len(mockStacks) {
		t.Errorf("ListStacks() returned %d stacks, want %d", len(stacks), len(mockStacks))
	}

	for i, stack := range stacks {
		if stack.ID != mockStacks[i].ID {
			t.Errorf("Stack %d ID = %d, want %d", i, stack.ID, mockStacks[i].ID)
		}
		if stack.Name != mockStacks[i].Name {
			t.Errorf("Stack %d Name = %q, want %q", i, stack.Name, mockStacks[i].Name)
		}
	}
}

func TestFindStackByName(t *testing.T) {
	mockStacks := []Stack{
		{ID: 1, Name: "test-stack-1"},
		{ID: 2, Name: "test-stack-2"},
		{ID: 3, Name: "another-stack"},
	}

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockStacks)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:  server.URL,
		AuthType: "none",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	// Test finding existing stack
	stack, err := client.FindStackByName("test-stack-2")
	if err != nil {
		t.Fatalf("FindStackByName() error = %v", err)
	}

	if stack == nil {
		t.Fatalf("FindStackByName() returned nil for existing stack")
	}

	if stack.ID != 2 {
		t.Errorf("FindStackByName() ID = %d, want %d", stack.ID, 2)
	}
	if stack.Name != "test-stack-2" {
		t.Errorf("FindStackByName() Name = %q, want %q", stack.Name, "test-stack-2")
	}

	// Test finding non-existing stack
	stack, err = client.FindStackByName("non-existing")
	if err != nil {
		t.Fatalf("FindStackByName() error = %v", err)
	}

	if stack != nil {
		t.Errorf("FindStackByName() should return nil for non-existing stack")
	}
}

func TestGetStack(t *testing.T) {
	mockStack := Stack{
		ID:           1,
		Name:         "test-stack",
		Type:         1,
		EndpointID:   1,
		Status:       1,
		CreationDate: time.Now().Unix(),
		CreatedBy:    "admin",
	}

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/stacks/1" {
			t.Errorf("Expected path /api/stacks/1, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockStack)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:  server.URL,
		AuthType: "none",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	stack, err := client.GetStack(1)
	if err != nil {
		t.Fatalf("GetStack() error = %v", err)
	}

	if stack.ID != mockStack.ID {
		t.Errorf("GetStack() ID = %d, want %d", stack.ID, mockStack.ID)
	}
	if stack.Name != mockStack.Name {
		t.Errorf("GetStack() Name = %q, want %q", stack.Name, mockStack.Name)
	}
}

func TestCreateStackFromString(t *testing.T) {
	mockStack := Stack{
		ID:   1,
		Name: "test-stack",
	}

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/stacks/create/standalone/string" {
			t.Errorf("Expected path /api/stacks/create/standalone/string, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Verify query parameters
		query := r.URL.Query()
		if query.Get("method") != "string" {
			t.Errorf("Expected method=string, got %s", query.Get("method"))
		}
		if query.Get("type") != "1" {
			t.Errorf("Expected type=1, got %s", query.Get("type"))
		}
		if query.Get("endpointId") != "1" {
			t.Errorf("Expected endpointId=1, got %s", query.Get("endpointId"))
		}

		// Verify request body
		var req CreateStackRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Name != "test-stack" {
			t.Errorf("Expected Name=test-stack, got %s", req.Name)
		}
		if req.StackFileContent != "version: '3'\nservices:\n  app:\n    image: nginx" {
			t.Errorf("Unexpected StackFileContent: %s", req.StackFileContent)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockStack)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		EndpointID: 1,
		AuthType:   "none",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	req := CreateStackRequest{
		Name:             "test-stack",
		StackFileContent: "version: '3'\nservices:\n  app:\n    image: nginx",
	}

	stack, err := client.CreateStackFromString(req)
	if err != nil {
		t.Fatalf("CreateStackFromString() error = %v", err)
	}

	if stack.ID != mockStack.ID {
		t.Errorf("CreateStackFromString() ID = %d, want %d", stack.ID, mockStack.ID)
	}
	if stack.Name != mockStack.Name {
		t.Errorf("CreateStackFromString() Name = %q, want %q", stack.Name, mockStack.Name)
	}
}

func TestCreateStackFromRepository(t *testing.T) {
	mockStack := Stack{
		ID:   1,
		Name: "test-stack",
	}

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/stacks" {
			t.Errorf("Expected path /api/stacks, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Verify query parameters
		query := r.URL.Query()
		if query.Get("method") != "repository" {
			t.Errorf("Expected method=repository, got %s", query.Get("method"))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockStack)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		EndpointID: 1,
		AuthType:   "none",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	req := CreateStackRequest{
		Name:                        "test-stack",
		RepositoryURL:               "https://github.com/user/repo",
		ComposeFilePathInRepository: "docker-compose.yml",
	}

	stack, err := client.CreateStackFromRepository(req)
	if err != nil {
		t.Fatalf("CreateStackFromRepository() error = %v", err)
	}

	if stack.ID != mockStack.ID {
		t.Errorf("CreateStackFromRepository() ID = %d, want %d", stack.ID, mockStack.ID)
	}
}

func TestUpdateStack(t *testing.T) {
	mockStack := Stack{
		ID:   1,
		Name: "updated-stack",
	}

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/stacks/1" {
			t.Errorf("Expected path /api/stacks/1, got %s", r.URL.Path)
		}
		if r.Method != "PUT" {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		// Verify query parameters
		query := r.URL.Query()
		if query.Get("endpointId") != "1" {
			t.Errorf("Expected endpointId=1, got %s", query.Get("endpointId"))
		}

		// Verify request body
		var req UpdateStackRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if !req.Prune {
			t.Errorf("Expected Prune=true, got %v", req.Prune)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockStack)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		EndpointID: 1,
		AuthType:   "none",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	req := UpdateStackRequest{
		StackFileContent: "version: '3'\nservices:\n  app:\n    image: nginx:latest",
		Prune:            true,
		PullImage:        false,
	}

	stack, err := client.UpdateStack(1, req)
	if err != nil {
		t.Fatalf("UpdateStack() error = %v", err)
	}

	if stack.Name != mockStack.Name {
		t.Errorf("UpdateStack() Name = %q, want %q", stack.Name, mockStack.Name)
	}
}

func TestDeleteStack(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/stacks/1" {
			t.Errorf("Expected path /api/stacks/1, got %s", r.URL.Path)
		}
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}

		// Verify query parameters
		query := r.URL.Query()
		if query.Get("endpointId") != "1" {
			t.Errorf("Expected endpointId=1, got %s", query.Get("endpointId"))
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		EndpointID: 1,
		AuthType:   "none",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	err := client.DeleteStack(1)
	if err != nil {
		t.Fatalf("DeleteStack() error = %v", err)
	}
}

func TestGetStackFile(t *testing.T) {
	expectedContent := "version: '3'\nservices:\n  app:\n    image: nginx"

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/stacks/1/file" {
			t.Errorf("Expected path /api/stacks/1/file, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		response := struct {
			StackFileContent string `json:"StackFileContent"`
		}{
			StackFileContent: expectedContent,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:  server.URL,
		AuthType: "none",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	content, err := client.GetStackFile(1)
	if err != nil {
		t.Fatalf("GetStackFile() error = %v", err)
	}

	if content != expectedContent {
		t.Errorf("GetStackFile() = %q, want %q", content, expectedContent)
	}
}
