# MobileOps CLI

`mobileops` is the command-line interface for the [MobileOps](https://www.mobileops.at) maritime operations platform. Manage vessels, crew, jobs, maintenance, inspections, fuel and compliance from your terminal, or hand it to an AI agent.

```bash
mobileops auth login                                   # paste your API keys once
mobileops vessels list                                 # the fleet
mobileops jobs list --status LAUNCHED,IN_PROGRESS      # what is underway right now
mobileops work-requests create --vessel-id 123 --description "Pump leaking" --priority high
mobileops audits list                                  # SIRE inspections, surveys, audits
mobileops bunker-partitions list --vessel-id 123       # fuel lifts and the burn drawn from them
mobileops tree                                         # every command
```

Every command has `--json` output with a summary and breadcrumbs (the next commands you or an agent would want), and the CLI ships a skill so agents can discover it.

## Install

**macOS / Linux / WSL2**

```bash
curl -fsSL https://www.mobileops.at/install-cli | bash
```

**Windows:** download the zip for your architecture from [Releases](https://github.com/MobileOps-Team/mobileops-cli/releases/latest) and put `mobileops.exe` somewhere on your `PATH`.

<details>
<summary>Other installation methods</summary>

**Linux packages (deb/rpm/apk):**
```bash
# Download from https://github.com/MobileOps-Team/mobileops-cli/releases/latest
sudo apt install ./mobileops-cli_*_linux_amd64.deb                # Debian/Ubuntu
sudo dnf install ./mobileops-cli_*_linux_amd64.rpm                # Fedora/RHEL
sudo apk add --allow-untrusted ./mobileops-cli_*_linux_amd64.apk  # Alpine
```
Arm64: substitute `arm64` for `amd64` in the filename.

**Go install:**
```bash
go install github.com/MobileOps-Team/mobileops-cli/cmd/mobileops@latest
```

**GitHub Release:** download the archive for your platform from [Releases](https://github.com/MobileOps-Team/mobileops-cli/releases).

</details>

Once installed, the CLI keeps itself current. See [Updating](#updating).

## Authentication

API keys are issued per company. In MobileOps, open your company settings, choose **API Keys**, and create a key. Then:

```bash
mobileops auth login               # prompts for the access key and secret key
mobileops auth status              # verifies the stored keys against the API
mobileops auth logout              # removes stored credentials
```

Read commands work with any key. `create`, `update` and `delete` commands need a key with **write** permission; with a read-only key the API answers with a permission error.

Credentials are stored in `~/.config/mobileops/credentials.json` (mode 0600). To keep a second set of keys for a non-production environment, log in with `--env gamma` and pass `--env gamma` (or set `MOBILEOPS_ENV=gamma`) on later commands.

## Usage

```bash
mobileops vessels list                          # List vessels
mobileops vessels get 123                       # One vessel
mobileops crew list                             # Crew members
mobileops jobs list --vessel-id 123 --from 2026-09-01 --to 2026-09-30
mobileops jobs vessel-availability --vessel-id 123 --from 2026-10-01 --to 2026-10-07
mobileops components list --vessel-id 123       # Equipment on a vessel
mobileops routine-calculations list --vessel-id 123 --value-lte 0   # Overdue maintenance
mobileops deficiencies list --vessel-id 123     # Findings
mobileops observations list --audit-type "SIRE Inspection"
mobileops fuel-level-readings list --vessel-id 123 --from 2026-09-01T00:00:00Z
mobileops events search --job-ids <job_id>      # Wheelhouse timeline of a job
mobileops tree                                  # Full command tree
```

Each resource follows the same pattern: `list`, `get ID`, and where the API allows it `create`, `update ID` and `delete ID`. Filters are flags on `list`; writable fields are flags on `create` and `update`. Run any command with `--help` for its flags.

### Output

```bash
mobileops vessels list              # Table in the terminal
mobileops vessels list --json       # JSON envelope, for scripts and agents
```

### JSON envelope

```json
{
  "ok": true,
  "data": [ ... ],
  "summary": "Showing 10 of 62 vessels",
  "breadcrumbs": [
    "mobileops vessels get 123",
    "mobileops components list --vessel-id 123"
  ]
}
```

Errors use the same shape with `ok: false` and an `error` message, and the command exits non-zero:

```json
{
  "ok": false,
  "error": "API key does not have permission for this operation."
}
```

### Pagination

Every `list` command takes `--page` and `--limit` (default 10, maximum 100). The summary line tells you how many records exist in total.

## AI agents

`mobileops` works with any agent that can run shell commands. Install the skill so the agent knows the commands, the terminology and the common workflows:

```bash
npx skills add MobileOps-Team/mobileops-cli
```

For structured discovery, `mobileops --help --agent` prints a JSON manifest of every command, its flags, and the REST API operation it calls. The skill itself is at [`skills/mobileops/SKILL.md`](skills/mobileops/SKILL.md).

Agent runners that must stay on a fixed version should set `MOBILEOPS_AUTO_UPDATE=0` (see below).

## Updating

The CLI keeps itself current. Every command checks GitHub Releases in the background (once a day) and, when a newer version exists, replaces its own binary and refreshes the installed skill. You'll see one line on stderr:

```
✓ mobileops updated to v0.3.0 (takes effect on the next command)
```

This applies to any install whose binary the current user can overwrite, including the installer script and `go install`. If the binary lives somewhere the user cannot write (for example a root-owned `/usr/local/bin`), the CLI prints a reminder instead:

```
⚠ Update available: v0.2.1 → v0.3.0
  Run: mobileops update
```

```bash
mobileops update                   # Update now (binary + skill)
mobileops version                  # Current version
```

Set `MOBILEOPS_AUTO_UPDATE=0` to pin a version.

## Configuration

| Setting | Flag | Environment variable | Default |
|---------|------|----------------------|---------|
| Environment (which credentials file to use) | `--env` | `MOBILEOPS_ENV` | `production` |
| API host | `auth login --host` | | per environment |
| Unattended self-update | | `MOBILEOPS_AUTO_UPDATE` | `1` (set `0` to pin) |
| Daily release check | | `MOBILEOPS_SKIP_VERSION_CHECK` | unset (set `1` to skip) |

```
~/.config/mobileops/
├── credentials.json              # Production API keys
├── credentials.<env>.json        # Keys for --env <env>
└── version-check.json            # Cache for the daily release check
```

## REST API

The CLI is a thin layer over the MobileOps REST API. The API contract is the OpenAPI spec at [www.mobileops.at/openapi.yaml](https://www.mobileops.at/openapi.yaml), documented at [docs.mobileops.app/api](https://docs.mobileops.app/api). `mobileops --help --agent` shows which operation each command calls, and the CLI's CI checks the command set against the spec so the two stay in step.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for building, testing, adding commands and releasing.

## License

[MIT](LICENSE.txt)
