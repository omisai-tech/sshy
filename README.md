# sshy

A Go-based CLI tool for managing and connecting to SSH servers via YAML or JSON configuration.

## Features

- **YAML & JSON configuration** - Define servers in YAML or JSON files
- **Remote URL support** - Fetch server configurations from a remote URL (HTTP/HTTPS)
- **Local overrides** - Override shared configurations with local settings
- **Private servers** - Keep private servers separate from shared configs
- **Fuzzy search** - Quickly find and connect to servers interactively
- **SSH passthrough** - Pass any SSH flags directly to the underlying command
- **SCP & SFTP support** - Copy files and start SFTP sessions using server names

## Installation

### Quick install (Linux/macOS)

```bash
# Linux (amd64)
curl -sL https://github.com/omisai-tech/sshy/releases/latest/download/sshy_$(curl -s https://api.github.com/repos/omisai-tech/sshy/releases/latest | grep '"tag_name"' | cut -d'"' -f4)_linux_amd64.tar.gz | tar xz -C /usr/local/bin sshy

# Linux (arm64)
curl -sL https://github.com/omisai-tech/sshy/releases/latest/download/sshy_$(curl -s https://api.github.com/repos/omisai-tech/sshy/releases/latest | grep '"tag_name"' | cut -d'"' -f4)_linux_arm64.tar.gz | tar xz -C /usr/local/bin sshy

# macOS (Apple Silicon)
curl -sL https://github.com/omisai-tech/sshy/releases/latest/download/sshy_$(curl -s https://api.github.com/repos/omisai-tech/sshy/releases/latest | grep '"tag_name"' | cut -d'"' -f4)_darwin_arm64.tar.gz | tar xz -C /usr/local/bin sshy

# macOS (Intel)
curl -sL https://github.com/omisai-tech/sshy/releases/latest/download/sshy_$(curl -s https://api.github.com/repos/omisai-tech/sshy/releases/latest | grep '"tag_name"' | cut -d'"' -f4)_darwin_amd64.tar.gz | tar xz -C /usr/local/bin sshy
```

Or use the install script:

```bash
curl -sSL https://raw.githubusercontent.com/omisai-tech/sshy/master/install.sh | bash
```

### From releases

Download the latest binary for your platform from the [releases page](https://github.com/omisai-tech/sshy/releases).

### From source

```bash
go install github.com/omisai-tech/sshy@latest
```

### Docker

```bash
docker pull ghcr.io/omisai-tech/sshy:latest
```

## Quick Start

```bash
# Initialize configuration
sshy init

# Add a server
sshy add

# List all servers
sshy list

# Connect to a server (interactive)
sshy

# Connect to a specific server
sshy connect my-server
```

## Configuration

sshy uses two configuration files and supports both YAML and JSON formats.

### Shared servers (`~/.sshy/servers.yaml` or `servers.json`)

Shared servers can be configured from a local file or a remote URL.

**YAML format:**

```yaml
- name: production-server
  host: 192.168.1.100
  user: admin
  port: 22
  tags:
    - prod
    - web

- name: staging-server
  host: 192.168.1.101
  user: deploy
  tags:
    - staging
```

**JSON format:**

```json
[
  {
    "name": "production-server",
    "host": "192.168.1.100",
    "user": "admin",
    "port": 22,
    "tags": ["prod", "web"]
  },
  {
    "name": "staging-server",
    "host": "192.168.1.101",
    "user": "deploy",
    "tags": ["staging"]
  }
]
```

### Remote URL Configuration

You can configure sshy to fetch shared servers from a remote URL instead of a local file. This is useful when:

- **VPN-protected environments**: Access sensitive server lists only when connected to your corporate VPN
- **Dynamic server lists**: Server configurations change frequently and you want to always fetch the latest
- **Centralized management**: Maintain server configurations in a central location accessible to your team

To configure URL-based servers, run `sshy init` and choose option 2 for remote URL:

```bash
sshy init
# Choose format (YAML/JSON)
# Choose source type: 2) Remote URL
# Enter the URL for your servers configuration
```

The global config (`~/.sshy/config.yaml`) will contain:

```yaml
servers_url: https://internal.company.com/api/servers.yaml
config_path: /home/user/.sshy
```

The URL endpoint should return valid YAML or JSON in the same format as the local `servers.yaml` file. sshy will automatically detect the format based on:

1. Response content (JSON starts with `{` or `[`)
2. Content-Type header (`application/json`, `application/x-yaml`)
3. URL file extension (`.json`, `.yaml`, `.yml`)

### Local overrides (`~/.sshy/local.yaml` or `local.json`)

```yaml
servers:
  production-server:
    key: ~/.ssh/prod_key

private:
  - name: personal-server
    host: my-private-host.com
    user: me
    port: 22
    tags:
      - personal
```

## Usage

### Connect to a server

```bash
# Interactive selection with fuzzy search
sshy

# Direct connection
sshy connect server-name

# Pass SSH flags
sshy connect server-name -v

# Execute remote command
sshy connect server-name -- ls -la
```

### Manage servers

```bash
# Add a new server
sshy add
sshy add server-name host.com user

# Edit a server
sshy edit server-name

# Remove a server
sshy rm server-name

# List servers
sshy list
sshy list --tags prod,web
```

### File operations

```bash
# Copy file to server
sshy scp local-file.txt server-name:/remote/path/

# Copy file from server
sshy scp server-name:/remote/file.txt ./local/

# Start SFTP session
sshy sftp server-name
```

### View configuration

```bash
sshy view
```

### Edit local overrides

```bash
sshy local
```

### Update sshy

```bash
# Check for updates and update to the latest version
sshy update
```

## Server sources

When listing servers, prefixes indicate the source:

- `[S]` - Shared servers from `servers.yaml`/`servers.json` or remote URL
- `[L]` - Local private servers from `local.yaml`/`local.json`
- `[O]` - Shared servers with local overrides

## License

MIT License - see [LICENSE](LICENSE) for details.
