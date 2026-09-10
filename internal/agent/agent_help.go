package agent

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/MobileOps-Team/mobileops-cli/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type option struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Default     string `json:"default,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Description string `json:"description,omitempty"`
}

type outputShape struct {
	Type      string `json:"type"`
	DataShape string `json:"data_shape,omitempty"`
}

type command struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	API         []string    `json:"api,omitempty"`
	Options     []option    `json:"options"`
	Output      outputShape `json:"output"`
}

type authInfo struct {
	Type    string   `json:"type"`
	Setup   string   `json:"setup"`
	Headers []string `json:"headers"`
}

type manifest struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Auth        authInfo  `json:"auth"`
	GlobalFlags []option  `json:"global_flags"`
	Commands    []command `json:"commands"`
}

// Generate prints the manifest for the given root command. apiMap gives, per
// command path (without the root name), the REST operations the command calls.
func Generate(root *cobra.Command, apiMap map[string][]string) {
	m := manifest{
		Name:        root.Name(),
		Version:     version.Version,
		Description: root.Long,
		Auth: authInfo{
			Type:    "api_key_pair",
			Setup:   root.Name() + " auth login",
			Headers: []string{"X-Api-Access-Key", "X-Api-Secret-Key"},
		},
		GlobalFlags: flagOptions(root.PersistentFlags()),
	}

	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		for _, child := range c.Commands() {
			if child.Name() == "help" || child.Name() == "completion" {
				continue
			}
			if child.Runnable() {
				path := strings.TrimPrefix(child.CommandPath(), root.Name()+" ")
				m.Commands = append(m.Commands, command{
					Name:        strings.TrimPrefix(child.UseLine(), root.Name()+" "),
					Description: child.Short,
					API:         apiMap[path],
					Options:     flagOptions(child.LocalFlags()),
					Output:      shapeFor(child, apiMap[path]),
				})
			}
			walk(child)
		}
	}
	walk(root)
	sort.Slice(m.Commands, func(i, j int) bool { return m.Commands[i].Name < m.Commands[j].Name })

	data, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println(string(data))
}

func flagOptions(fs *pflag.FlagSet) []option {
	opts := []option{}
	fs.VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == "help" {
			return
		}
		o := option{Name: "--" + f.Name, Type: f.Value.Type(), Description: f.Usage}
		if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" {
			o.Default = f.DefValue
		}
		if f.Annotations[cobra.BashCompOneRequiredFlag] != nil {
			o.Required = true
		}
		opts = append(opts, o)
	})
	return opts
}

func shapeFor(c *cobra.Command, ops []string) outputShape {
	if len(ops) == 0 {
		return outputShape{Type: "message"}
	}
	name := c.Name()
	resource := strings.ReplaceAll(c.Parent().Name(), "-", "_")
	switch {
	case name == "list" || name == "search" || name == "specs" || name == "attachments" || name == "export" || name == "vessel-availability":
		return outputShape{Type: "envelope", DataShape: "array<" + strings.TrimSuffix(resource, "s") + ">"}
	case name == "delete":
		return outputShape{Type: "envelope", DataShape: "message"}
	default:
		return outputShape{Type: "envelope", DataShape: strings.TrimSuffix(resource, "s")}
	}
}
