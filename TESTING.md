# Testing Guide

This document describes the test suite for the Portainer CLI project.

## Test Coverage

The project includes comprehensive unit tests for all major components:

### Package Coverage
- **pkg/storage**: 69.8% coverage - Tests for secure credential storage and encryption
- **pkg/client**: 68.6% coverage - Tests for HTTP client and API interactions  
- **cmd**: 22.5% coverage - Tests for CLI command validation and configuration
- **main**: 0.0% coverage - Simple entry point with minimal logic

## Test Structure

### Storage Package Tests (`pkg/storage/*_test.go`)

**crypto_test.go**
- `TestEncryptDecrypt` - Tests encryption/decryption with various data types
- `TestEncryptConsistency` - Verifies different encryptions produce different ciphertext
- `TestDecryptInvalidData` - Tests error handling for invalid encrypted data
- `TestHashPassword` - Tests password hashing functionality
- `TestGetMachineID` - Tests machine ID generation for key derivation
- `TestDeriveKey` - Tests cryptographic key derivation

**config_test.go** 
- `TestConfigSetGetURL` - Tests URL configuration storage and retrieval
- `TestConfigSetGetEndpoint` - Tests endpoint ID configuration
- `TestConfigCredentials` - Tests encrypted credential storage
- `TestConfigToken` - Tests JWT token storage and validation
- `TestConfigClearCredentials` - Tests credential clearing functionality
- `TestLoadConfigNewFile` - Tests loading configuration from non-existent file
- `TestConfigFilePermissions` - Tests secure file permissions (0600)

### Client Package Tests (`pkg/client/*_test.go`)

**client_test.go**
- `TestNewClient` - Tests client initialization with various configurations
- `TestMakeRequestJWT` - Tests JWT authentication in HTTP requests
- `TestMakeRequestAPIKey` - Tests API key authentication
- `TestMakeRequestNoAuth` - Tests unauthenticated requests
- `TestMakeRequestError` - Tests error handling for failed requests
- `TestSetJWTToken/SetAPIKey` - Tests authentication method switching
- `TestIsAuthPath` - Tests authentication path detection

**auth_test.go**
- `TestAuthenticateUser` - Tests username/password authentication
- `TestAuthenticateWithCode` - Tests access code authentication  
- `TestValidateToken` - Tests JWT token validation
- `TestParseJWTExpiry` - Tests JWT expiration parsing
- Error handling tests for authentication failures

**stacks_test.go**
- `TestListStacks` - Tests stack listing functionality
- `TestFindStackByName` - Tests stack lookup by name
- `TestGetStack` - Tests individual stack retrieval
- `TestCreateStackFromString` - Tests stack creation from inline content
- `TestCreateStackFromRepository` - Tests stack creation from Git repository
- `TestUpdateStack` - Tests stack updating
- `TestDeleteStack` - Tests stack deletion
- `TestGetStackFile` - Tests stack file content retrieval

### Command Package Tests (`cmd/*_test.go`)

**config_test.go**
- `TestConfigSetUrlCmd` - Tests URL setting command with validation
- `TestConfigSetEndpointCmd` - Tests endpoint setting command
- `TestConfigGetCmd` - Tests configuration display command

## Running Tests

### Basic Test Execution
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run specific package tests
go test ./pkg/storage -v
go test ./pkg/client -v
go test ./cmd -v
```

### Using Make Commands
```bash
# Run all tests
make test

# Run tests with verbose output  
make test-verbose

# Run tests with coverage report
make test-coverage

# Run tests for specific packages
make test-storage
make test-client
make test-cmd
```

### Coverage Analysis
```bash
# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Or use make command
make test-coverage
```

## Test Patterns

### Mock HTTP Servers
Tests use `httptest.NewServer` to create mock Portainer API servers for testing HTTP client functionality without requiring a real Portainer instance.

### Temporary Directories
Storage tests use `t.TempDir()` and environment variable manipulation to test file operations in isolation.

### Table-Driven Tests
Many tests use table-driven patterns to test multiple scenarios:

```go
tests := []struct {
    name        string
    input       string
    expected    string
    expectError bool
}{
    // test cases...
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // test implementation
    })
}
```

### Error Testing
Tests verify both success and failure scenarios, ensuring proper error handling and user feedback.

## Security Testing

### Encryption Tests
- Verify encrypted data differs from plaintext
- Test decryption of various data types
- Validate error handling for tampered data
- Test key derivation consistency

### File Permission Tests  
- Verify config files are created with secure permissions (0600)
- Test that sensitive data is properly protected

### Authentication Tests
- Test various authentication methods (JWT, API key, none)
- Verify proper header setting for different auth types
- Test token validation and refresh scenarios

## Integration Testing

While the current test suite focuses on unit tests, integration tests can be added by:

1. Setting up a test Portainer instance
2. Using build tags to separate integration tests
3. Running tests against real API endpoints

```bash
# Run integration tests (when implemented)
go test ./test/... -tags=integration -v
```

## Continuous Integration

The project includes GitHub Actions workflow (`.github/workflows/test.yml`) that:
- Runs on push and pull requests
- Tests on multiple Go versions
- Generates coverage reports
- Runs linting checks

## Best Practices

1. **Test Isolation**: Each test runs in isolation with temporary directories
2. **Comprehensive Coverage**: Tests cover both success and error scenarios  
3. **Mock Dependencies**: External dependencies are mocked for unit tests
4. **Clear Test Names**: Test names clearly describe what is being tested
5. **Table-Driven Tests**: Multiple scenarios tested efficiently
6. **Security Focus**: Encryption and authentication are thoroughly tested