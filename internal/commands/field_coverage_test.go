package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/MobileOps-Team/mobileops-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

func TestAPIFieldCoverage(t *testing.T) {
	specPath := os.Getenv("MOBILEOPS_OPENAPI_SPEC")
	if specPath == "" {
		t.Skip("set MOBILEOPS_OPENAPI_SPEC=/path/to/openapi.yaml (or run `make api-check`)")
	}
	ops, err := loadSpecFields(specPath)
	if err != nil {
		t.Fatalf("load spec: %v", err)
	}

	// Fake API: record what arrives, answer with a shape every command accepts.
	type request struct {
		op    string
		query []string
		body  map[string]interface{}
	}
	var seen []request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		var query []string
		for k := range r.URL.Query() {
			query = append(query, k)
		}
		seen = append(seen, request{
			op:    r.Method + " " + templatePath(strings.TrimPrefix(r.URL.Path, "/api")),
			query: query,
			body:  body,
		})
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[],"total":0,"has_more":false,"jobs":[]}`)
	}))
	defer srv.Close()

	// Point the CLI at it and keep it quiet.
	t.Setenv("HOME", t.TempDir())
	t.Setenv("MOBILEOPS_ENV", config.EnvTest)
	t.Setenv("MOBILEOPS_SKIP_VERSION_CHECK", "1")
	t.Setenv("MOBILEOPS_AUTO_UPDATE", "0")
	if err := config.Save(srv.URL, "test-access", "test-secret", config.EnvTest); err != nil {
		t.Fatal(err)
	}
	devnull, _ := os.Open(os.DevNull)
	realStdout := os.Stdout
	os.Stdout = devnull
	defer func() { os.Stdout = realStdout }()

	sent := map[string]map[string]bool{} // op -> fields sent by any command
	var problems []string

	for _, path := range leafCommandPaths(rootCmd) {
		if nonAPICommands[path] {
			continue
		}
		cmd := findCommand(rootCmd, path)
		start := len(seen)
		rootCmd.SetArgs(argsWithEveryFlag(cmd, path))
		if err := rootCmd.Execute(); err != nil {
			problems = append(problems, fmt.Sprintf("command %q failed against the fake API: %v", path, err))
			continue
		}
		if len(seen) == start {
			problems = append(problems, fmt.Sprintf("command %q made no API request", path))
			continue
		}

		for _, req := range seen[start:] {
			spec, ok := ops[req.op]
			if !ok {
				problems = append(problems, fmt.Sprintf("command %q called %s, which is not in the spec", path, req.op))
				continue
			}
			if sent[req.op] == nil {
				sent[req.op] = map[string]bool{}
			}
			for _, q := range req.query {
				if q == "page" || q == "limit" {
					continue
				}
				sent[req.op]["?"+q] = true
				if !spec.query[q] {
					problems = append(problems, fmt.Sprintf("command %q sends query param %q to %s, which the spec does not define", path, q, req.op))
				}
			}
			if spec.freeFormBody {
				continue
			}
			for _, f := range flattenBody(req.body) {
				sent[req.op][f] = true
				if !spec.body[f] {
					problems = append(problems, fmt.Sprintf("command %q sends body field %q to %s, which the spec does not define", path, f, req.op))
				}
			}
		}
	}

	// Reverse direction: contract fields no command exposes.
	for op, spec := range ops {
		if sent[op] == nil {
			continue // uncovered operations are TestAPICoverage's job
		}
		for q := range spec.query {
			if q == "page" || q == "limit" || sent[op]["?"+q] {
				continue
			}
			if _, ok := apiFieldIgnored[op+" ?"+q]; ok {
				continue
			}
			problems = append(problems, fmt.Sprintf("spec query param %q on %s has no CLI flag (add one, or add %q to apiFieldIgnored)", q, op, op+" ?"+q))
		}
		for f := range spec.body {
			if sent[op][f] {
				continue
			}
			if _, ok := apiFieldIgnored[op+" "+f]; ok {
				continue
			}
			problems = append(problems, fmt.Sprintf("spec body field %q on %s has no CLI flag (add one, or add %q to apiFieldIgnored)", f, op, op+" "+f))
		}
	}
	for key := range apiFieldIgnored {
		// key is "VERB /path field": the operation is the first two words.
		parts := strings.SplitN(key, " ", 3)
		if len(parts) != 3 {
			problems = append(problems, fmt.Sprintf("apiFieldIgnored key %q is not \"VERB /path field\"", key))
			continue
		}
		op, field := parts[0]+" "+parts[1], parts[2]
		spec, ok := ops[op]
		if !ok || !(spec.query[strings.TrimPrefix(field, "?")] || spec.body[field]) {
			problems = append(problems, fmt.Sprintf("apiFieldIgnored lists %q but the spec no longer has it", key))
		}
	}

	if len(problems) > 0 {
		sort.Strings(problems)
		t.Errorf("CLI flags are out of sync with %s:\n  - %s", specPath, strings.Join(problems, "\n  - "))
	}
}

type specOperation struct {
	query        map[string]bool
	body         map[string]bool // "wrapper.field" for wrapped fields, "field" for top-level ones
	freeFormBody bool            // request schema has no declared properties
}

// loadSpecFields returns, per "VERB /path", the query parameter names and the
// request body fields (one level below a wrapper object) the spec declares.
func loadSpecFields(path string) (map[string]specOperation, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var root map[string]interface{}
	if err := yaml.Unmarshal(raw, &root); err != nil {
		return nil, err
	}
	resolve := func(v interface{}) map[string]interface{} {
		for i := 0; i < 10; i++ {
			m, ok := v.(map[string]interface{})
			if !ok {
				return nil
			}
			ref, ok := m["$ref"].(string)
			if !ok {
				return m
			}
			var cur interface{} = root
			for _, seg := range strings.Split(strings.TrimPrefix(ref, "#/"), "/") {
				cm, ok := cur.(map[string]interface{})
				if !ok {
					return nil
				}
				cur = cm[seg]
			}
			v = cur
		}
		return nil
	}

	ops := map[string]specOperation{}
	paths, _ := root["paths"].(map[string]interface{})
	for p, item := range paths {
		pm, _ := item.(map[string]interface{})
		for verb, raw := range pm {
			switch verb {
			case "get", "post", "put", "patch", "delete":
			default:
				continue
			}
			op := resolve(raw)
			if desc, _ := op["description"].(string); strings.HasPrefix(strings.TrimSpace(desc), "Not currently available") {
				continue
			}
			so := specOperation{query: map[string]bool{}, body: map[string]bool{}}
			if params, ok := op["parameters"].([]interface{}); ok {
				for _, pr := range params {
					prm := resolve(pr)
					if prm["in"] == "query" {
						if name, ok := prm["name"].(string); ok {
							so.query[name] = true
						}
					}
				}
			}
			if rb := resolve(op["requestBody"]); rb != nil {
				content, _ := rb["content"].(map[string]interface{})
				appJSON, _ := content["application/json"].(map[string]interface{})
				schema := resolve(appJSON["schema"])
				props, _ := schema["properties"].(map[string]interface{})
				if len(props) == 0 {
					so.freeFormBody = true
				}
				for name, sub := range props {
					so.body[name] = true
					if subProps, ok := resolve(sub)["properties"].(map[string]interface{}); ok {
						for subName := range subProps {
							so.body[name+"."+subName] = true
						}
					}
				}
			}
			ops[strings.ToUpper(verb)+" "+p] = so
		}
	}
	return ops, nil
}

// templatePath turns the concrete path a command requested into the spec's
// template form: the test passes "1" for every positional ID.
func templatePath(p string) string {
	segs := strings.Split(p, "/")
	for i, s := range segs {
		if s == "1" {
			segs[i] = "{id}"
		}
	}
	return strings.Join(segs, "/")
}

// flattenBody lists the fields a request body carries as "wrapper.field" for
// nested objects and "field" for top-level values.
func flattenBody(body map[string]interface{}) []string {
	var out []string
	for k, v := range body {
		if m, ok := v.(map[string]interface{}); ok {
			out = append(out, k)
			for sub := range m {
				out = append(out, k+"."+sub)
			}
			continue
		}
		out = append(out, k)
	}
	return out
}

func findCommand(root *cobra.Command, path string) *cobra.Command {
	cmd, _, err := root.Find(strings.Fields(path))
	if err != nil {
		return nil
	}
	return cmd
}

// argsWithEveryFlag builds argv for a command with "1" for each positional
// argument and every local flag set (strings to "1", JSON flags to "{}",
// booleans to true) so the request carries every field the command can send.
func argsWithEveryFlag(cmd *cobra.Command, path string) []string {
	args := strings.Fields(path)
	for _, tok := range strings.Fields(cmd.Use)[1:] {
		if strings.ToUpper(tok) == tok && !strings.HasPrefix(tok, "[") {
			args = append(args, "1")
		}
	}
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if f.Name == "help" {
			return
		}
		switch f.Value.Type() {
		case "bool":
			args = append(args, "--"+f.Name+"=true")
		case "string":
			val := "1"
			if strings.Contains(f.Usage, "JSON") {
				val = "{}"
			}
			args = append(args, "--"+f.Name+"="+val)
		}
	})
	return append(args, "--json")
}
