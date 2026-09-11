# mobileops-cli

Go CLI (cobra) over the MobileOps REST API. The API's OpenAPI spec is the
contract: `public/openapi.yaml` in MobileOps-Web, published at
https://www.mobileops.at/openapi.yaml.

Follow [CONTRIBUTING.md](CONTRIBUTING.md) for the layout, the two coverage
tests (`make api-check`), the procedure for adding or changing a command, the
flag conventions, and releasing (merge to master with a bumped
`internal/version/version.go`; never tag by hand).

Notes for Claude specifically:

- Commits and tags are made by the maintainer; do not commit.
- Verify against production with read-only commands only
  (`bin/mobileops <resource> list --limit 2 --json`); never run `create`,
  `update` or `delete` against real data.
- `mobileops --help --agent` is generated from the cobra tree plus
  `apiCoverage`; it never needs a manual edit. `SKILL.md` does.
- After any command change, run `make api-check` against the local spec:
  `make api-check SPEC_FILE=../MobileOps-Web/public/openapi.yaml`.
