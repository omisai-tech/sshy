# Configuration Examples

This directory contains example configuration files for sshy in both YAML and JSON formats.

## File Overview

### Shared Server Configuration
Files that are typically stored in a Git repository and shared with your team:

- `servers.yaml` - Shared servers list (YAML format)
- `servers.json` - Shared servers list (JSON format)

### Local Configuration
Files stored in `~/.sshy/` that are NOT shared (contain private servers and local overrides):

- `local.yaml` - Local overrides and private servers (YAML format)
- `local.json` - Local overrides and private servers (JSON format)

### Global Configuration
The main sshy configuration file stored in `~/.sshy/`:

- `config.yaml` - Global sshy settings (YAML format)
- `config.json` - Global sshy settings (JSON format)

## Usage

When you run `sshy init`, you'll be prompted to choose your preferred format (YAML or JSON). The tool will then create the appropriate configuration files.

### Server Properties

| Property | Description | Required |
|----------|-------------|----------|
| `name` | Unique identifier for the server | Yes |
| `host` | Hostname or IP address | Yes |
| `user` | SSH username | No |
| `port` | SSH port (default: 22) | No |
| `tags` | List of tags for filtering | No |
| `key` | Path to SSH private key | No |
| `options` | Additional SSH options | No |

### Local Configuration

The local configuration file supports two sections:

1. **servers** - Override properties of shared servers (matched by name)
2. **private** - Define servers that are not shared with your team

This allows you to:
- Add SSH keys to shared servers without committing them
- Override user or host for your local environment
- Keep personal servers separate from team servers
