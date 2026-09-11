# Contributing

`mobileops` is a Go CLI built on [cobra](https://github.com/spf13/cobra). It is a thin layer over the MobileOps REST API, whose contract is the OpenAPI spec published at https://www.mobileops.at/openapi.yaml (source: `public/openapi.yaml` in the MobileOps-Web repository).

## Building and testing

You need Go (the version in `go.mod`) and `make`.

```bash
make build        # bin/mobileops
make test         # go test ./...
make api-check    # compare the command tree with the published API spec (see below)
go vet ./...
```

To try a build against production, log in once with `bin/mobileops auth login` and stick to `list`/`get` commands. Never point a `create`, `update` or `delete` at a customer's company while developing.

## Layout

```
cmd/mobileops/            main()
internal/commands/        one file per resource (vessels.go, jobs.go, ...), root.go wires them up
internal/commands/api_map.go          command -> REST operation map (drives the coverage tests and --agent)
internal/commands/coverage_test.go    TestAPICoverage: operations
internal/commands/field_coverage_test.go  TestAPIFieldCoverage: query params and body fields
internal/client/          HTTP client, auth headers, error mapping
internal/envelope/        the {ok, data, summary, breadcrumbs} envelope
internal/formatter/       table vs JSON output
internal/config/          credentials files and environments
internal/updater/         release check and self-update
internal/agent/           the --help --agent manifest, generated from the command tree
internal/version/         the single version constant
skills/mobileops/SKILL.md the agent skill (hand-written)
scripts/install.sh        the curl installer (served at https://www.mobileops.at/install-cli)
```

## Keeping the CLI in sync with the REST API

The spec is the contract: implement what it documents, not what the Rails controllers happen to accept. Two tests enforce this and both run in `make api-check`, which downloads the published spec (use `SPEC_FILE=../MobileOps-Web/public/openapi.yaml` to check against an unmerged branch):

- **`TestAPICoverage`** (operations): every spec operation must have a command in `apiCoverage` (`api_map.go`), every command must call an operation the spec has, and no command may call an operation whose description starts with "Not currently available". Operations that should never be commands (webhooks) go in `apiIgnored` with a reason.
- **`TestAPIFieldCoverage`** (fields): runs every command against a fake API with all flags set and diffs the query parameters and body fields it sends against the spec, in both directions. A contract field you deliberately don't expose goes in `apiFieldIgnored` with a reason.

The same check runs weekly in this repository (`.github/workflows/api-coverage.yml`) and on every REST API pull request in MobileOps-Web, so a new endpoint shows up as a red check without anyone remembering to look.

### Adding or changing a command

1. Read the operation in the spec, not the controller. Note in particular payloads that go beside the wrapper rather than inside it (customers `add_vessels`, position reports, `jobs/vessel-availability`, `jobs/job-export`, `events/search`).
2. New resource: copy the closest existing file (`wheelhouse_logs.go` for list-only, `codes.go` for full CRUD), register the command in `root.go`, and add one `apiCoverage` line per subcommand.
3. New filter or field on an existing endpoint: add the flag and its `flag -> api_param` entry in that resource's file. Query filters use `listFilters`; body fields use `bodyFromFlags`, `bodyFromBoolFlags`, `bodyFromArrayFlags` (comma-separated) or `bodyFromJSONFlags` (objects and arrays of objects).
4. Endpoint removed or marked not available: delete the subcommand and its `apiCoverage` line.
5. Update `skills/mobileops/SKILL.md`: the resource section, the "Write Operations Summary" table and, for a new concept, the terminology aliases. Agents only know what the skill tells them.
6. `make build && make test && make api-check`, then run the new command read-only against production with `--json` and check the summary and breadcrumbs.

### Conventions

- Every `list` command takes `--page` and `--limit` through `paginationParams`.
- The vessel filter is always `--vessel-id`, even when the API parameter is `asset_id` or `job_asset_id`.
- Date windows are `--from` and `--to`.
- Boolean API fields are boolean flags; comma-separated flags map to arrays; nested objects take JSON.
- Breadcrumbs point at the next command a person or agent would want, using real IDs from the response when possible.
- Errors go through `handleClientError` so `--json` callers get the envelope and a non-zero exit.

## Releasing

Merging to `master` is the release. `release.yml` reads `internal/version/version.go`, and if no tag exists for that version it creates `vX.Y.Z` and runs GoReleaser, which publishes the archives and Linux packages. Installed CLIs then self-update within a day and reinstall the skill.

So every pull request that changes `cmd/` or `internal/` must bump `internal/version/version.go`; `ci.yml` fails otherwise, and also fails if the version is already tagged. Do not tag by hand unless the workflow failed.

Skill-only changes (`skills/`) take effect on merge without a release, because `npx skills add` reads `SKILL.md` from the default branch. Merge skill changes together with the command changes they describe so an installed binary is never behind its skill.
