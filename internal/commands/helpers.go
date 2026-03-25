package commands

import (
	"fmt"
	"os"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
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
