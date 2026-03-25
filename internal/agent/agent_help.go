package agent

import (
	"encoding/json"
	"fmt"
)

const version = "0.1.0"

type option struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Default     interface{} `json:"default,omitempty"`
	Max         int         `json:"max,omitempty"`
	Description string      `json:"description,omitempty"`
}

type outputShape struct {
	Type      string `json:"type"`
	DataShape string `json:"data_shape,omitempty"`
}

type command struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
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
	Commands    []command `json:"commands"`
}

// Generate prints the machine-readable agent help manifest.
func Generate() {
	m := manifest{
		Name:        "mobileops",
		Version:     version,
		Description: "CLI for managing vessels, jobs, crew, and operations via the MobileOps API",
		Auth: authInfo{
			Type:    "api_key_pair",
			Setup:   "mobileops auth login",
			Headers: []string{"X-Api-Access-Key", "X-Api-Secret-Key"},
		},
		Commands: []command{
			{
				Name:        "auth login",
				Description: "Authenticate with MobileOps API using access key and secret key",
				Options: []option{
					{Name: "--host", Type: "string", Default: "https://www.mobileops.at", Description: "API host URL"},
				},
				Output: outputShape{Type: "message"},
			},
			{
				Name:        "auth logout",
				Description: "Remove stored credentials",
				Options:     []option{},
				Output:      outputShape{Type: "message"},
			},
			{
				Name:        "auth status",
				Description: "Check authentication status and API connectivity",
				Options:     []option{},
				Output:      outputShape{Type: "envelope", DataShape: "auth_status"},
			},
			{
				Name:        "vessels list",
				Description: "List all vessels/assets with pagination",
				Options: []option{
					{Name: "--page", Type: "integer", Default: 1},
					{Name: "--limit", Type: "integer", Default: 10, Max: 100},
					{Name: "--json", Type: "boolean", Default: false},
				},
				Output: outputShape{Type: "envelope", DataShape: "array<vessel>"},
			},
			{
				Name:        "vessels get ID",
				Description: "Get a single vessel by ID",
				Options: []option{
					{Name: "--json", Type: "boolean", Default: false},
				},
				Output: outputShape{Type: "envelope", DataShape: "vessel"},
			},
			{
				Name:        "vessels specs ID",
				Description: "List specs for a vessel",
				Options: []option{
					{Name: "--json", Type: "boolean", Default: false},
				},
				Output: outputShape{Type: "envelope", DataShape: "array<vessel_spec>"},
			},
			{
				Name:        "jobs list",
				Description: "List jobs with optional filters",
				Options: []option{
					{Name: "--page", Type: "integer", Default: 1},
					{Name: "--limit", Type: "integer", Default: 10, Max: 100},
					{Name: "--vessel-id", Type: "string", Description: "Filter by vessel ID"},
					{Name: "--from", Type: "string", Description: "Start date (YYYY-MM-DD)"},
					{Name: "--to", Type: "string", Description: "End date (YYYY-MM-DD)"},
					{Name: "--active-only", Type: "boolean", Default: false},
					{Name: "--json", Type: "boolean", Default: false},
				},
				Output: outputShape{Type: "envelope", DataShape: "array<job>"},
			},
			{
				Name:        "jobs get ID",
				Description: "Get a single job by ID",
				Options: []option{
					{Name: "--json", Type: "boolean", Default: false},
				},
				Output: outputShape{Type: "envelope", DataShape: "job"},
			},
			{
				Name:        "crew list",
				Description: "List crew members/users",
				Options: []option{
					{Name: "--page", Type: "integer", Default: 1},
					{Name: "--limit", Type: "integer", Default: 10, Max: 100},
					{Name: "--json", Type: "boolean", Default: false},
				},
				Output: outputShape{Type: "envelope", DataShape: "array<user>"},
			},
			{
				Name:        "crew get ID",
				Description: "Get a single crew member by ID",
				Options: []option{
					{Name: "--json", Type: "boolean", Default: false},
				},
				Output: outputShape{Type: "envelope", DataShape: "user"},
			},
			{
				Name:        "crew find-by-employee-number NUMBER",
				Description: "Find a crew member by employee number",
				Options: []option{
					{Name: "--json", Type: "boolean", Default: false},
				},
				Output: outputShape{Type: "envelope", DataShape: "user"},
			},
			{
				Name:        "components list",
				Description: "List components with optional vessel filter",
				Options: []option{
					{Name: "--page", Type: "integer", Default: 1},
					{Name: "--limit", Type: "integer", Default: 10, Max: 100},
					{Name: "--vessel-id", Type: "string", Description: "Filter by vessel ID"},
					{Name: "--json", Type: "boolean", Default: false},
				},
				Output: outputShape{Type: "envelope", DataShape: "array<component>"},
			},
			{
				Name:        "components get ID",
				Description: "Get a single component by ID",
				Options: []option{
					{Name: "--json", Type: "boolean", Default: false},
				},
				Output: outputShape{Type: "envelope", DataShape: "component"},
			},
			{
				Name:        "work-requests list",
				Description: "List work requests with optional vessel filter",
				Options: []option{
					{Name: "--page", Type: "integer", Default: 1},
					{Name: "--limit", Type: "integer", Default: 10, Max: 100},
					{Name: "--vessel-id", Type: "string", Description: "Filter by vessel ID"},
					{Name: "--json", Type: "boolean", Default: false},
				},
				Output: outputShape{Type: "envelope", DataShape: "array<work_request>"},
			},
			{
				Name:        "work-requests get ID",
				Description: "Get a single work request by ID",
				Options: []option{
					{Name: "--json", Type: "boolean", Default: false},
				},
				Output: outputShape{Type: "envelope", DataShape: "work_request"},
			},
		},
	}

	data, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println(string(data))
}
