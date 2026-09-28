package commands

import (
	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var formInstancesCmd = &cobra.Command{
	Use:   "form-instances",
	Short: "List form instances",
}

var formInstancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List form instances",
	RunE:  runFormInstancesList,
}

func init() {
	formInstancesListCmd.Flags().Int("page", 1, "Page number")
	formInstancesListCmd.Flags().Int("limit", 10, "Items per page (max 100)")
	formInstancesListCmd.Flags().String("template-id", "", "Filter by form template ID")
	formInstancesListCmd.Flags().String("job-id", "", "Filter by job ID")
	formInstancesListCmd.Flags().String("created-after", "", "Created on or after this time (ISO 8601)")
	formInstancesListCmd.Flags().String("created-before", "", "Created on or before this time (ISO 8601)")
	formInstancesListCmd.Flags().String("updated-after", "", "Updated on or after this time (ISO 8601)")
	formInstancesListCmd.Flags().String("updated-before", "", "Updated on or before this time (ISO 8601)")
	formInstancesListCmd.Flags().String("sort", "", "Sort by created_at (default) or updated_at")
	formInstancesListCmd.Flags().String("order", "", "Sort direction: desc (default) or asc")

	formInstancesCmd.AddCommand(formInstancesListCmd)
}

func runFormInstancesList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	listFilters := map[string]string{
		"template-id":    "form_template_id",
		"job-id":         "job_id",
		"created-after":  "created_at_gte",
		"created-before": "created_at_lte",
		"updated-after":  "updated_at_gte",
		"updated-before": "updated_at_lte",
		"sort":           "sort",
		"order":          "order",
	}
	for flag, paramKey := range listFilters {
		if v, _ := cmd.Flags().GetString(flag); v != "" {
			params[paramKey] = v
		}
	}

	response, err := c.Get("form-instances", params)
	if err != nil {
		return handleClientError(err)
	}

	env := envelope.WrapCollection(response, "form instances", []string{
		"mobileops forms list",
		"mobileops jobs list",
	})
	formatter.Output(env, jsonOutput)
	return nil
}
