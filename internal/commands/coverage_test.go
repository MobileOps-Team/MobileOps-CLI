package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// TestAPICoverage checks the CLI command tree against the REST API's OpenAPI
// spec. It needs the spec on disk; point MOBILEOPS_OPENAPI_SPEC at it
// (`make api-check` downloads the published one). Without the variable the
// test is skipped so `go test ./...` stays offline.
func TestAPICoverage(t *testing.T) {
	specPath := os.Getenv("MOBILEOPS_OPENAPI_SPEC")
	if specPath == "" {
		t.Skip("set MOBILEOPS_OPENAPI_SPEC=/path/to/openapi.yaml (or run `make api-check`)")
	}

	available, unavailable, err := loadSpecOperations(specPath)
	if err != nil {
		t.Fatalf("load spec: %v", err)
	}

	var problems []string

	// 1. Every API-backed leaf command must be in apiCoverage, and every
	//    operation it names must exist and be available in the spec.
	covered := map[string]bool{}
	for _, path := range leafCommandPaths(rootCmd) {
		if nonAPICommands[path] {
			continue
		}
		ops, ok := apiCoverage[path]
		if !ok {
			problems = append(problems, fmt.Sprintf("command %q has no entry in apiCoverage (api_map.go)", path))
			continue
		}
		for _, op := range ops {
			covered[op] = true
			switch {
			case available[op]:
			case unavailable[op]:
				problems = append(problems, fmt.Sprintf("command %q calls %s, which the spec marks as not currently available", path, op))
			default:
				problems = append(problems, fmt.Sprintf("command %q calls %s, which is not in the spec", path, op))
			}
		}
	}

	// 2. apiCoverage must not name commands that no longer exist.
	existing := map[string]bool{}
	for _, path := range leafCommandPaths(rootCmd) {
		existing[path] = true
	}
	for path := range apiCoverage {
		if !existing[path] {
			problems = append(problems, fmt.Sprintf("apiCoverage names %q but no such command exists", path))
		}
	}

	// 3. Every available spec operation must be covered or explicitly ignored.
	for op := range available {
		if covered[op] {
			continue
		}
		if _, ignored := apiIgnored[op]; ignored {
			continue
		}
		problems = append(problems, fmt.Sprintf("spec operation %s has no CLI command (add one and list it in api_map.go, or add it to apiIgnored)", op))
	}
	for op := range apiIgnored {
		if !available[op] && !unavailable[op] {
			problems = append(problems, fmt.Sprintf("apiIgnored lists %s but the spec no longer has it", op))
		}
	}

	if len(problems) > 0 {
		sort.Strings(problems)
		t.Errorf("CLI is out of sync with %s:\n  - %s", specPath, strings.Join(problems, "\n  - "))
	}
}

// loadSpecOperations returns the "VERB /path" operations in an OpenAPI file,
// split into those that work and those whose description says
// "Not currently available".
func loadSpecOperations(path string) (available, unavailable map[string]bool, err error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	// Path items also carry non-operation keys (e.g. "parameters"), so decode
	// loosely and only look at the HTTP verbs.
	var spec struct {
		Paths map[string]map[string]interface{} `yaml:"paths"`
	}
	if err := yaml.Unmarshal(raw, &spec); err != nil {
		return nil, nil, err
	}

	available = map[string]bool{}
	unavailable = map[string]bool{}
	for p, item := range spec.Paths {
		for verb, op := range item {
			switch verb {
			case "get", "post", "put", "patch", "delete":
			default:
				continue
			}
			key := strings.ToUpper(verb) + " " + p
			desc := ""
			if m, ok := op.(map[string]interface{}); ok {
				desc, _ = m["description"].(string)
			}
			if strings.HasPrefix(strings.TrimSpace(desc), "Not currently available") {
				unavailable[key] = true
			} else {
				available[key] = true
			}
		}
	}
	return available, unavailable, nil
}

// leafCommandPaths lists runnable commands as "group sub" paths (root omitted).
func leafCommandPaths(root *cobra.Command) []string {
	var out []string
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		for _, child := range c.Commands() {
			if child.Name() == "help" || child.Name() == "completion" {
				continue
			}
			if child.Runnable() {
				out = append(out, strings.TrimPrefix(child.CommandPath(), root.Name()+" "))
			}
			walk(child)
		}
	}
	walk(root)
	sort.Strings(out)
	return out
}
