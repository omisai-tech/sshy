# sshy

A Go-based CLI tool for managing and connecting to SSH servers via YAML configuration.

## Features

- **YAML-based configuration** - Define servers in simple YAML files
- **Local overrides** - Override shared configurations with local settings
- **Private servers** - Keep private servers separate from shared configs
- **Fuzzy search** - Quickly find and connect to servers interactively
- **SSH passthrough** - Pass any SSH flags directly to the underlying command
- **SCP & SFTP support** - Copy files and start SFTP sessions using server names

## Installation

### From releases

Download the latest binary from the [releases page](https://github.com/omisai-tech/sshy/releases).

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

sshy uses two configuration files:

### Shared servers (`~/.sshy/servers.yaml`)

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

### Local overrides (`~/.sshy/local.yaml`)

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

## Server sources

When listing servers, prefixes indicate the source:

- `[S]` - Shared servers from `servers.yaml`
- `[L]` - Local private servers from `local.yaml`
- `[O]` - Shared servers with local overrides

## License

MIT License - see [LICENSE](LICENSE) for details.
