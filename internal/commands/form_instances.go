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

	formInstancesCmd.AddCommand(formInstancesListCmd)
}

func runFormInstancesList(cmd *cobra.Command, args []string) error {
	c, err := client.New("", "", "", envFlag)
	if err != nil {
		return handleClientError(err)
	}

	params := paginationParams(cmd)
	if v, _ := cmd.Flags().GetString("template-id"); v != "" {
		params["form_template_id"] = v
	}
	if v, _ := cmd.Flags().GetString("job-id"); v != "" {
		params["job_id"] = v
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
