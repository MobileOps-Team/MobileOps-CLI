package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

// handleClientError prints client errors in a user-friendly way.
func handleClientError(err error) error {
	if jsonOutput {
		env := envelope.WrapError(err.Error(), nil)
		formatter.Output(env, true)
		os.Exit(1)
	}

	switch err.(type) {
	case *client.AuthError:
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	case *client.NotFoundError:
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	case *client.PermissionError:
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	case *client.ApiError:
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	default:
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
	return nil
}

// collectionBreadcrumbs extracts IDs from the first N records in a collection response
// and formats them as breadcrumb commands.
func collectionBreadcrumbs(response map[string]interface{}, cmdTemplate string, max int) []string {
	data, ok := response["data"].([]interface{})
	if !ok || len(data) == 0 {
		return nil
	}

	limit := max
	if len(data) < limit {
		limit = len(data)
	}

	var breadcrumbs []string
	for i := 0; i < limit; i++ {
		record, ok := data[i].(map[string]interface{})
		if !ok {
			continue
		}
		id := envelope.ExtractID(record)
		if id != "" {
			breadcrumbs = append(breadcrumbs, fmt.Sprintf(cmdTemplate, id))
		}
	}
	return breadcrumbs
}

// bodyFromFlags builds a JSON body from string flags on a command.
// flagMap maps flag-name -> api_param_name.
func bodyFromFlags(cmd *cobra.Command, flagMap map[string]string) map[string]interface{} {
	body := make(map[string]interface{})
	for flag, param := range flagMap {
		if cmd.Flags().Changed(flag) {
			val, _ := cmd.Flags().GetString(flag)
			body[param] = val
		}
	}
	return body
}

// bodyFromBoolFlags adds boolean flags to an existing body.
func bodyFromBoolFlags(cmd *cobra.Command, body map[string]interface{}, flagMap map[string]string) {
	for flag, param := range flagMap {
		if cmd.Flags().Changed(flag) {
			val, _ := cmd.Flags().GetBool(flag)
			body[param] = val
		}
	}
}

// bodyFromArrayFlags adds array flags (comma-separated) to an existing body.
func bodyFromArrayFlags(cmd *cobra.Command, body map[string]interface{}, flagMap map[string]string) {
	for flag, param := range flagMap {
		if cmd.Flags().Changed(flag) {
			val, _ := cmd.Flags().GetString(flag)
			parts := strings.Split(val, ",")
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
			}
			body[param] = parts
		}
	}
}

// wrapBody wraps the body in a resource key (e.g., {"asset": {...}}).
func wrapBody(key string, body map[string]interface{}) map[string]interface{} {
	if len(body) == 0 {
		return body
	}
	return map[string]interface{}{key: body}
}

func bodyFromJSONFlags(cmd *cobra.Command, body map[string]interface{}, flagMap map[string]string) error {
	for flag, param := range flagMap {
		if !cmd.Flags().Changed(flag) {
			continue
		}
		raw, _ := cmd.Flags().GetString(flag)
		var parsed interface{}
		if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
			return fmt.Errorf("--%s must be valid JSON: %v", flag, err)
		}
		body[param] = parsed
	}
	return nil
}

func listFilters(cmd *cobra.Command, params map[string]string, flagMap map[string]string) {
	for flag, param := range flagMap {
		if v, _ := cmd.Flags().GetString(flag); v != "" {
			params[param] = v
		}
	}
}
