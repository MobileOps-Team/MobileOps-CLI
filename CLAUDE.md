# mobileops-cli

Go CLI (cobra) over the MobileOps REST API. The API's OpenAPI spec is the
source of truth: `public/openapi.yaml` in MobileOps-Web, published at
https://www.mobileops.at/openapi.yaml.

## Keeping the CLI in sync with the REST API

1. `make api-check` (or `make api-check SPEC_FILE=../MobileOps-Web/public/openapi.yaml`
   for an unmerged branch). It downloads the spec and runs `TestAPICoverage`,
   which lists every spec operation without a command, every command that
   calls an operation the spec no longer has, and every command that calls an
   operation whose description starts with "Not currently available".
2. For each reported gap:
   - New resource: add `internal/commands/<resource>.go` (copy the closest
     existing file: `wheelhouse_logs.go` for list-only, `codes.go` for CRUD),
     register it in `root.go`, and add one line per subcommand to
     `apiCoverage` in `internal/commands/api_map.go`.
   - New filter or body field on an existing endpoint: add the flag and its
     `flag -> api_param` mapping in that resource's file. Query filters use
     `listFilters`; body fields use `bodyFromFlags` / `bodyFromBoolFlags` /
     `bodyFromArrayFlags` / `bodyFromJSONFlags`.
   - Endpoint removed or marked not available: delete the subcommand and its
     `apiCoverage` line.
   - Operation that should never be a command (webhooks): add it to
     `apiIgnored` with a reason.
3. Read the operation's description and parameters in the spec, not the Rails
   controller: some payloads go beside the wrapper instead of inside it
   (customers `add_vessels`, position reports, `jobs/vessel-availability`).
4. Update `skills/mobileops/SKILL.md`: the resource section, the "Write
   Operations Summary" table and the terminology aliases. Agents read the
   skill, so a command that is not documented there is effectively invisible.
5. `make build && make test && make api-check`, then run a read-only command
   against production (`bin/mobileops <resource> list --limit 2 --json`).
6. Bump `internal/version/version.go` in the same PR (CI fails otherwise).
   Merging to master tags `vX.Y.Z` and runs GoReleaser (`release.yml`); users'
   binaries then self-update within a day and reinstall the skill with
   `npx skills update MobileOps-Team/mobileops-cli`. Never tag by hand unless
   the workflow failed.

`mobileops --help --agent` is generated from the cobra tree plus
`apiCoverage`, so it never needs a manual edit.

## Conventions

- Every list command takes `--page` / `--limit` via `paginationParams`.
- Vessel filters are called `--vessel-id` on the CLI even when the API
  parameter is `asset_id` or `job_asset_id`.
- Date-window filters are `--from` / `--to`.
- Breadcrumbs point at the next command an agent would want.
- Commits and tags are made by the maintainer; do not commit.
