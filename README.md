# MobileOps CLI

`mobileops` is the command-line interface for the [MobileOps](https://www.mobileops.at) platform. Manage vessels, jobs, crew, maintenance, and more from your terminal or through AI agents.

- Works standalone or with any AI agent (Claude, Codex, Copilot, Gemini)
- JSON output with breadcrumbs for easy navigation
- API key authentication
- Includes agent skill for AI discovery

## Quick Start

```bash
curl -fsSL https://www.mobileops.at/install-cli | bash
```

That's it. You now have full access to MobileOps from your terminal.

<details>
<summary>Other installation methods</summary>

**Linux (deb/rpm/apk):**
```bash
# Download from https://github.com/MobileOps-Team/mobileops-cli/releases/latest
sudo apt install ./mobileops-cli_*_linux_amd64.deb            # Debian/Ubuntu
sudo dnf install ./mobileops-cli_*_linux_amd64.rpm            # Fedora/RHEL
sudo apk add --allow-untrusted ./mobileops-cli_*_linux_amd64.apk  # Alpine
```
Arm64: substitute `arm64` for `amd64` in the filename.

**Go install:**
```bash
go install github.com/MobileOps-Team/mobileops-cli/cmd/mobileops@latest
```

**GitHub Release:** download from [Releases](https://github.com/MobileOps-Team/mobileops-cli/releases).

</details>

## Usage

```bash
mobileops vessels list                          # List vessels
mobileops vessels get 123                       # Get vessel details
mobileops jobs list --vessel-id 123             # Jobs for a vessel
mobileops crew list                             # List crew members
mobileops components list --vessel-id 123       # Components for a vessel
mobileops work-requests list --vessel-id 123    # Work requests
mobileops deficiencies list --vessel-id 123     # Deficiencies
mobileops tree                                  # See all commands
```

### Output Formats

```bash
mobileops vessels list              # Styled output in terminal
mobileops vessels list --json       # JSON with envelope and breadcrumbs
```

### JSON Envelope

Every command supports `--json` for structured output:

```json
{
  "ok": true,
  "data": [...],
  "summary": "6 vessels found",
  "breadcrumbs": [
    "mobileops vessels get 123",
    "mobileops components list --vessel-id 123"
  ]
}
```

Breadcrumbs suggest next commands, making it easy for humans and agents to navigate.

## Authentication

Get your API keys from **Settings > REST API** in MobileOps, then:

```bash
mobileops auth login               # Authenticate with MobileOps
mobileops auth status              # Check current auth status
mobileops auth logout              # Remove stored credentials
```

Credentials are stored in `~/.config/mobileops/credentials.json` (mode 0600).

## AI Agent Integration

`mobileops` works with any AI agent that can run shell commands.

**Install the skill:**
```bash
npx skills add MobileOps-Team/mobileops-cli
```

**Agent discovery:** Every command supports `--help --agent` for structured JSON output (flags, subcommands). Use `mobileops tree` for the full command catalog.

See [`skills/mobileops/SKILL.md`](skills/mobileops/SKILL.md) for the full skill reference.

## Configuration

```
~/.config/mobileops/              # Your MobileOps identity
├── credentials.json              #   Production API keys
└── credentials.<env>.json        #   Environment-specific keys
```

## License

[MIT](LICENSE)
