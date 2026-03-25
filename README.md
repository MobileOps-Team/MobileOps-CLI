# MobileOps CLI

Command-line interface for the [MobileOps](https://www.mobileops.at) platform. Manage vessels, jobs, crew, components, and work requests from your terminal.

Built for **humans and AI agents alike** — every command supports `--json` output with breadcrumb navigation, and `--help --agent` provides machine-readable command discovery.

## Install

### Quick install (macOS / Linux)

```bash
curl -fsSL https://www.mobileops.at/install-cli | bash
```

### Homebrew (macOS)

```bash
brew install MobileOps-Team/tap/mobileops
```

### Linux packages

```bash
# Debian / Ubuntu
sudo apt install ./mobileops-cli_*_linux_amd64.deb

# Fedora / RHEL
sudo dnf install ./mobileops-cli_*_linux_amd64.rpm

# Alpine
sudo apk add --allow-untrusted ./mobileops-cli_*_linux_amd64.apk
```

For ARM64 systems, substitute `arm64` for `amd64` in filenames.

### Scoop (Windows)

```bash
scoop bucket add mobileops https://github.com/MobileOps-Team/homebrew-tap
scoop install mobileops-cli
```

### Go

```bash
go install github.com/MobileOps-Team/mobileops-cli/cmd/mobileops@latest
```

### GitHub Releases

Pre-built binaries for macOS, Linux, and Windows are available on the [Releases](https://github.com/MobileOps-Team/mobileops-cli/releases) page.

## Authenticate

Get your API keys from **Settings > REST API** in MobileOps, then:

```bash
mobileops auth login
# Access Key: ********
# Secret Key: ********
# Verifying credentials... OK
```

Credentials are stored in `~/.config/mobileops/credentials.json` (mode 0600).

### Environments

```bash
mobileops auth login --env development   # http://localhost:3000
mobileops auth login --env gamma         # https://gamma.mobileops.at/
mobileops auth login --env production    # https://www.mobileops.at (default)
```

Check your auth status:

```bash
mobileops auth status
```

## Usage

```bash
# Vessels
mobileops vessels list
mobileops vessels get VESSEL_ID
mobileops vessels specs VESSEL_ID

# Jobs
mobileops jobs list
mobileops jobs list --vessel-id VESSEL_ID --from 2026-01-01 --to 2026-03-31
mobileops jobs list --active-only
mobileops jobs get JOB_ID

# Crew
mobileops crew list
mobileops crew get USER_ID
mobileops crew find-by-employee-number 12345

# Components
mobileops components list --vessel-id VESSEL_ID
mobileops components get COMPONENT_ID

# Work Requests
mobileops work-requests list --vessel-id VESSEL_ID
mobileops work-requests get WORK_REQUEST_ID

# See all commands
mobileops tree
```

## JSON Output

Add `--json` to any command for structured output with breadcrumbs:

```bash
mobileops vessels get abc123 --json
```

```json
{
  "ok": true,
  "data": { "name": "MV Atlantic", "status": "active" },
  "summary": "Vessel abc123",
  "breadcrumbs": [
    "mobileops vessels specs abc123",
    "mobileops components list --vessel-id abc123",
    "mobileops jobs list --vessel-id abc123"
  ]
}
```

## AI Agent Integration

For AI agents (Claude, GPT, etc.), get a machine-readable command manifest:

```bash
mobileops --help --agent
```

This returns a JSON schema of all commands, their parameters, and expected output shapes — so agents can discover and use the CLI without reading docs.

## Pagination

All list commands support pagination:

```bash
mobileops vessels list --page 2 --limit 25
```

## Shell Completion

```bash
# Bash
mobileops completion bash > /etc/bash_completion.d/mobileops

# Zsh
mobileops completion zsh > "${fpath[1]}/_mobileops"

# Fish
mobileops completion fish > ~/.config/fish/completions/mobileops.fish
```

## Configuration

| Setting | Location |
|---|---|
| Credentials | `~/.config/mobileops/credentials.json` |
| Custom host | `mobileops auth login --host https://your-instance.mobileops.at` |
| Environment | `--env` flag or `MOBILEOPS_ENV` env var |

## Development

```bash
git clone https://github.com/MobileOps-Team/mobileops-cli.git
cd mobileops-cli
make build
make test
```

## License

MIT
