# Portainer CLI

A command-line interface for managing Docker Compose stacks through the Portainer API with secure authentication.

## Installation

Build from source:

```bash
go build -o portainer-cli
```

## Quick Start

1. **Configure Portainer URL:**
   ```bash
   portainer-cli config set url https://your-portainer-instance.com
   ```

2. **Set default endpoint (optional):**
   ```bash
   portainer-cli config set endpoint 1
   ```

3. **Login with username/password:**
   ```bash
   portainer-cli login password
   ```
   
   Or login with access code:
   ```bash
   portainer-cli login code
   ```

4. **Use stack commands:**
   ```bash
   portainer-cli stacks list
   ```

## Authentication

The CLI supports two authentication methods:

### JWT Authentication (Recommended)
Uses username/password or access code to obtain a JWT token that's automatically refreshed:

```bash
# Login with username/password
portainer-cli login password --username myuser

# Login with access code
portainer-cli login code --code ABCD1234
```

Credentials are stored securely encrypted on your local machine and tokens are automatically refreshed when they expire.

### Legacy API Key Authentication
For backward compatibility, you can still use API keys:

```bash
portainer-cli --url https://portainer.com --token your-api-key --endpoint 1 stacks list
```

## Configuration

### Secure Configuration (Recommended)
Store configuration securely using the config commands:

```bash
# Set Portainer URL
portainer-cli config set url https://your-portainer-instance.com

# Set default endpoint ID
portainer-cli config set endpoint 1

# View current configuration
portainer-cli config get

# Login and store credentials
portainer-cli login password
```

### Environment Variables (Legacy)
- `PORTAINER_URL`: Portainer URL
- `PORTAINER_TOKEN`: Portainer API token
- `PORTAINER_ENDPOINT`: Portainer endpoint ID

### Configuration File (Legacy)
Create a file at `~/.portainer-cli.yaml`:

```yaml
url: https://your-portainer-instance.com
token: your-api-token
endpoint: 1
```

## Commands

### List stacks
```bash
portainer-cli stacks list
```

### Create a stack

From a local file:
```bash
portainer-cli stacks create --name my-stack --file docker-compose.yml
```

From a Git repository:
```bash
portainer-cli stacks create --name my-stack --repository https://github.com/user/repo --compose-path docker-compose.yml
```

From inline content:
```bash
portainer-cli stacks create --name my-stack --content "version: '3.8'
services:
  app:
    image: nginx:alpine
    ports:
      - '80:80'"
```

### Update a stack
```bash
portainer-cli stacks update 123 --file docker-compose.yml --prune --pull-image
```

### Delete a stack
```bash
portainer-cli stacks delete 123
```

Force delete without confirmation:
```bash
portainer-cli stacks delete 123 --force
```

### Inspect a stack
```bash
portainer-cli stacks inspect 123
```

Show stack file content:
```bash
portainer-cli stacks inspect 123 --show-file
```

## Security Features

- **Encrypted Storage**: Credentials and tokens are encrypted using AES-256-GCM with machine-specific keys
- **Automatic Token Refresh**: JWT tokens are automatically refreshed when they expire
- **Secure Password Input**: Passwords are entered securely without echoing to the terminal
- **Local Storage**: All sensitive data is stored locally in `~/.config/portainer-cli/`

## Logout

To clear stored credentials and tokens:

```bash
portainer-cli logout
```

## Legacy API Token Usage

For environments that require API tokens, you can still create them in Portainer:

1. Log in to your Portainer instance
2. Go to "User settings" → "Access tokens"
3. Click "Add access token"
4. Copy the generated token and use it with the `--token` flag or `PORTAINER_TOKEN` environment variable

## Test Environment

The `test/` directory contains sample Docker Compose files for testing the CLI:

- `docker-compose.yml`: Full podinfo stack with Redis
- `docker-compose.simple.yml`: Simple podinfo deployment
- `test-cli.sh`: Test script demonstrating CLI usage

```bash
# Run the test script
cd test
./test-cli.sh
```

## Requirements

- Portainer CE/EE with API access
- Valid credentials or API token with appropriate permissions
- Network access to your Portainer instance